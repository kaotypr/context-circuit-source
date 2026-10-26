package world

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGuard(t *testing.T) {
	parent := t.TempDir()
	checkout := filepath.Join(parent, "checkout")
	home := filepath.Join(parent, "home")
	run := filepath.Join(parent, "run")
	for _, dir := range []string{checkout, home, run} {
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("pass", func(t *testing.T) {
		if err := Guard(run, checkout, home); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("run root inside checkout", func(t *testing.T) {
		nested := filepath.Join(checkout, "nested")
		if err := os.Mkdir(nested, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := Guard(nested, checkout, home); err == nil {
			t.Fatal("expected refusal")
		}
	})

	t.Run("checkout inside run root", func(t *testing.T) {
		outer := filepath.Join(parent, "outer")
		inner := filepath.Join(outer, "src")
		if err := os.MkdirAll(inner, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := Guard(outer, inner, home); err == nil {
			t.Fatal("expected refusal")
		}
	})

	t.Run("symlink into checkout", func(t *testing.T) {
		link := filepath.Join(parent, "linked-run")
		target := filepath.Join(checkout, "via-link")
		if err := os.Mkdir(target, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
		if err := Guard(link, checkout, home); err == nil {
			t.Fatal("expected refusal")
		}
	})

	t.Run("run root is home", func(t *testing.T) {
		if err := Guard(home, checkout, home); err == nil {
			t.Fatal("expected refusal")
		}
	})

	t.Run("run root is ancestor of home", func(t *testing.T) {
		root := filepath.Join(parent, "ancestor")
		nestedHome := filepath.Join(root, "user")
		if err := os.MkdirAll(nestedHome, 0o755); err != nil {
			t.Fatal(err)
		}
		other := filepath.Join(parent, "other-checkout")
		if err := os.Mkdir(other, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := Guard(root, other, nestedHome); err == nil {
			t.Fatal("expected refusal")
		}
	})

	t.Run("ancestor has git", func(t *testing.T) {
		repo := filepath.Join(parent, "repo")
		nested := filepath.Join(repo, "work")
		if err := os.MkdirAll(nested, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := Guard(nested, checkout, home); err == nil {
			t.Fatal("expected refusal")
		}
	})

	t.Run("not empty", func(t *testing.T) {
		full := filepath.Join(parent, "full")
		if err := os.Mkdir(full, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(full, "file"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := Guard(full, checkout, home); err == nil {
			t.Fatal("expected refusal")
		}
	})
}

func TestNewRunRoot(t *testing.T) {
	parent := t.TempDir()
	dir, err := NewRunRoot(parent)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatal("not a directory")
	}
	if info.Mode().Perm() != 0o700 {
		t.Fatalf("mode %o", info.Mode().Perm())
	}
	if filepath.Dir(dir) != parent {
		t.Fatalf("parent %s", filepath.Dir(dir))
	}
}
