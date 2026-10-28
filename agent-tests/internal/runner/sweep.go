package runner

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const runRootPrefix = "agent-test-"

// sweepBeforeHosts removes leftover subscription login files under parent
// and reports how many it removed. It does not print paths.
func sweepBeforeHosts(out io.Writer, parent string) error {
	n, err := sweepSubscriptionFiles(parent)
	if err != nil {
		return err
	}
	if out != nil {
		fmt.Fprintf(out, "swept %d subscription login file(s)\n", n)
	}
	return nil
}

// sweepSubscriptionFiles deletes auth.json at <parent>/agent-test-*/codex/auth.json
// when that run directory is owned by the current user. Symlinks are left
// in place. parent empty means the system temp directory.
func sweepSubscriptionFiles(parent string) (int, error) {
	if parent == "" {
		parent = os.TempDir()
	}
	entries, err := os.ReadDir(parent)
	if err != nil {
		return 0, err
	}
	n := 0
	failed := false
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), runRootPrefix) {
			continue
		}
		removed, err := sweepOne(filepath.Join(parent, entry.Name()))
		if err != nil {
			failed = true
			continue
		}
		if removed {
			n++
		}
	}
	if failed {
		return n, errLoginRemove
	}
	return n, nil
}

func sweepOne(dir string) (bool, error) {
	info, err := os.Lstat(dir)
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return false, nil
	}
	if !ownedByCurrentUser(info) {
		return false, nil
	}
	codexDir := filepath.Join(dir, "codex")
	codexInfo, err := os.Lstat(codexDir)
	if err != nil || codexInfo.Mode()&os.ModeSymlink != 0 || !codexInfo.IsDir() {
		return false, nil
	}
	target := filepath.Join(codexDir, "auth.json")
	tinfo, err := os.Lstat(target)
	if err != nil || tinfo.Mode()&os.ModeSymlink != 0 || !tinfo.Mode().IsRegular() {
		return false, nil
	}
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return false, errLoginRemove
	}
	return true, nil
}
