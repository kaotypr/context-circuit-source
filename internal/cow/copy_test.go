package cow

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// volumeCluster reports the cluster size of the volume holding dir. It is set
// only on Windows, where block cloning works in whole clusters.
var volumeCluster func(t *testing.T, dir string) int64

// TestNativeCloneOnCapableVolume runs cloneFile itself. It needs a directory on
// a clone-capable volume in CC_COW_NATIVE_DIR, and once given one it fails
// rather than skips when a clone does not happen.
func TestNativeCloneOnCapableVolume(t *testing.T) {
	dir := os.Getenv("CC_COW_NATIVE_DIR")
	if dir == "" {
		t.Skip("CC_COW_NATIVE_DIR does not name a clone-capable directory")
	}
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		t.Fatalf("CC_COW_NATIVE_DIR is not a directory: %s (%v)", dir, err)
	}
	// prepare writes data as source/module.js in a fresh directory under dir and
	// returns the source and the target Copy must create.
	prepare := func(t *testing.T, data []byte) (string, string) {
		t.Helper()
		root, err := os.MkdirTemp(dir, "cow-*")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.RemoveAll(root) })
		source, target := filepath.Join(root, "source"), filepath.Join(root, "target")
		if err := os.Mkdir(source, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(source, "module.js"), data, 0600); err != nil {
			t.Fatal(err)
		}
		return source, target
	}
	patterned := func(size int64) []byte {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 251)
		}
		return data
	}
	windowsOnly := func(t *testing.T) {
		t.Helper()
		if runtime.GOOS != "windows" {
			t.Skip("ReFS cluster alignment applies only on Windows")
		}
	}

	t.Run("clone is independent", func(t *testing.T) {
		original := patterned(128 * 1024)
		source, target := prepare(t, original)
		stats, err := Copy(context.Background(), source, target, source, target, "required")
		if err != nil || stats.Cloned != 1 {
			t.Fatalf("required mode did not clone: %+v %v", stats, err)
		}
		clone, err := os.OpenFile(filepath.Join(target, "module.js"), os.O_WRONLY, 0)
		if err != nil {
			t.Fatal(err)
		}
		_, writeErr := clone.WriteAt(bytes.Repeat([]byte{0xff}, len(original)), 0)
		if err := errors.Join(writeErr, clone.Close()); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(source, "module.js"))
		if err != nil || !bytes.Equal(data, original) {
			t.Fatal("overwriting the clone changed the source", err)
		}
	})

	t.Run("aligned body and copied tail", func(t *testing.T) {
		windowsOnly(t)
		original := patterned(volumeCluster(t, dir) + 100)
		source, target := prepare(t, original)
		stats, err := Copy(context.Background(), source, target, source, target, "required")
		if err != nil || stats.Cloned != 1 {
			t.Fatalf("required mode did not clone: %+v %v", stats, err)
		}
		data, err := os.ReadFile(filepath.Join(target, "module.js"))
		if err != nil || !bytes.Equal(data, original) {
			t.Fatal("clone differs from its source", err)
		}
	})

	t.Run("file below one cluster", func(t *testing.T) {
		windowsOnly(t)
		source, target := prepare(t, patterned(100))
		if _, err := Copy(context.Background(), source, target, source, target, "required"); err == nil || !strings.Contains(err.Error(), "filesystem CoW unavailable") {
			t.Fatalf("required mode must refuse a file below one cluster, got %v", err)
		}
		source, target = prepare(t, patterned(100))
		stats, err := Copy(context.Background(), source, target, source, target, "auto")
		if err != nil || stats.Copied != 1 {
			t.Fatalf("auto mode did not fall back to copying: %+v %v", stats, err)
		}
	})
}

func TestCopyIsolationAndFallback(t *testing.T) {
	for _, mode := range []string{"copy", "auto", "required"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			source, target := filepath.Join(root, "source"), filepath.Join(root, "target")
			if err := os.Mkdir(source, 0700); err != nil {
				t.Fatal(err)
			}
			original := make([]byte, 128*1024)
			for i := range original {
				original[i] = byte(i % 251)
			}
			if err := os.WriteFile(filepath.Join(source, "module.js"), original, 0600); err != nil {
				t.Fatal(err)
			}
			stats, err := Copy(context.Background(), source, target, source, target, mode)
			if err != nil && mode == "required" {
				t.Skip("filesystem has no native clone support:", err)
			}
			if err != nil {
				t.Fatal(err)
			}
			if mode == "required" && stats.Cloned != 1 {
				t.Fatal("required mode did not clone", stats)
			}
			if mode == "copy" && stats.Copied != 1 {
				t.Fatal("copy mode did not copy", stats)
			}
			t.Logf("mode=%s cloned=%d copied=%d", mode, stats.Cloned, stats.Copied)
			if err := os.WriteFile(filepath.Join(target, "module.js"), []byte("target-only edit"), 0600); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(source, "module.js"))
			if err != nil || len(data) != len(original) || data[17] != original[17] {
				t.Fatal("destination writes changed source")
			}
			if _, err := Copy(context.Background(), source, target, source, target, mode); err == nil {
				t.Fatal("overwrote existing destination")
			}
		})
	}
}

func TestSymlinksAreIndependentAndExternalLinksRejected(t *testing.T) {
	root := t.TempDir()
	source, target := filepath.Join(root, "source"), filepath.Join(root, "target")
	if err := os.MkdirAll(filepath.Join(source, "node_modules", ".bin"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "node_modules", "tool.js"), []byte("tool"), 0755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(source, "node_modules", ".bin", "tool")
	if err := os.Symlink("../tool.js", link); err != nil {
		t.Skip("symlinks unavailable:", err)
	}
	stats, err := Copy(context.Background(), filepath.Join(source, "node_modules"), filepath.Join(target, "node_modules"), source, target, "auto")
	if err != nil || stats.Links != 1 {
		t.Fatal(stats, err)
	}
	if err := os.WriteFile(filepath.Join(target, "node_modules", ".bin", "tool"), []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(source, "node_modules", "tool.js"))
	if string(data) != "tool" {
		t.Fatal("symlink points into original checkout")
	}
	if err := os.Symlink(filepath.Join(root, "outside"), filepath.Join(source, "node_modules", "external")); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(root, "bad")
	if _, err := Copy(context.Background(), filepath.Join(source, "node_modules"), bad, source, filepath.Join(root, "new"), "auto"); !errors.Is(err, ErrNotReusable) {
		t.Fatalf("an external link must be refused as not reusable, got %v", err)
	}
	if _, err := os.Lstat(bad); !os.IsNotExist(err) {
		t.Fatal("published a partial entry")
	}
}

// A failure of the machine is not a property of the entry: it must not read as
// "not reusable", or a caller skipping such entries would swallow a full disk.
func TestAnIOFailureIsNotReportedAsNotReusable(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.Mkdir(source, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "module.js"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "target")
	if err := os.WriteFile(target, []byte("a file where the destination's parent must be"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := Copy(context.Background(), source, filepath.Join(target, "entry"), source, root, "auto")
	if err == nil || errors.Is(err, ErrNotReusable) {
		t.Fatalf("expected an ordinary I/O error, got %v", err)
	}
}
