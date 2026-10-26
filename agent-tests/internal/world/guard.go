package world

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Guard refuses a run root that could see or overwrite a checkout, a home
// directory, or any repository. Existing paths are resolved through
// symlinks before the checks.
func Guard(runRoot, checkout, home string) error {
	runRoot, err := resolve(runRoot)
	if err != nil {
		return err
	}
	checkout, err = resolve(checkout)
	if err != nil {
		return err
	}
	home, err = resolve(home)
	if err != nil {
		return err
	}
	if runRoot == "" {
		return fmt.Errorf("run root is empty")
	}
	if checkout != "" && (runRoot == checkout || inside(checkout, runRoot) || inside(runRoot, checkout)) {
		return fmt.Errorf("run root %s overlaps checkout %s", runRoot, checkout)
	}
	if home != "" && (runRoot == home || inside(runRoot, home)) {
		return fmt.Errorf("run root %s is %s or an ancestor of it", runRoot, home)
	}
	if git, ok := ancestorGit(runRoot); ok {
		return fmt.Errorf("run root is inside a repository or workspace (%s)", git)
	}
	info, err := os.Stat(runRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("run root %s is not a directory", runRoot)
	}
	entries, err := os.ReadDir(runRoot)
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("run root %s is not empty", runRoot)
	}
	return nil
}

// NewRunRoot creates a private directory under parent, or under the system
// temp directory when parent is empty. The directory mode is 0700.
func NewRunRoot(parent string) (string, error) {
	if parent == "" {
		parent = os.TempDir()
	}
	dir, err := os.MkdirTemp(parent, "agent-test-")
	if err != nil {
		return "", err
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

func resolve(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	abs = filepath.Clean(abs)
	if _, err := os.Lstat(abs); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		parent, base := filepath.Split(abs)
		parent = strings.TrimRight(parent, string(os.PathSeparator))
		if parent == "" || parent == abs {
			return abs, nil
		}
		resolvedParent, err := resolve(parent)
		if err != nil {
			return "", err
		}
		return filepath.Join(resolvedParent, base), nil
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func inside(parent, child string) bool {
	if parent == "" || child == "" || parent == child {
		return false
	}
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func ancestorGit(runRoot string) (string, bool) {
	dir := runRoot
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		git := filepath.Join(parent, ".git")
		if _, err := os.Lstat(git); err == nil {
			return git, true
		}
		dir = parent
	}
}
