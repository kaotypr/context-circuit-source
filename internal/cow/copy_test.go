package cow

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

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
	if _, err := Copy(context.Background(), filepath.Join(source, "node_modules"), bad, source, filepath.Join(root, "new"), "auto"); err == nil {
		t.Fatal("accepted external link")
	}
	if _, err := os.Lstat(bad); !os.IsNotExist(err) {
		t.Fatal("published a partial entry")
	}
}
