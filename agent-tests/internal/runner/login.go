package runner

import (
	"errors"
	"os"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
)

// afterLoginCopy, when set, runs after a subscription login file has been
// copied and before host setup. Tests use it to cover a panic on that path.
var afterLoginCopy func()

var errLoginRemove = errors.New("warning: subscription login file could not be removed")

func loginCopyPaths(copies []hosts.FileCopy) []string {
	var out []string
	for _, copy := range copies {
		if copy.Dst != "" {
			out = append(out, copy.Dst)
		}
	}
	return out
}

// removeLoginFiles deletes copied login files. A missing file is success.
// os.Remove does not follow a symlink. Errors do not include the path or
// the file contents.
func removeLoginFiles(paths []string) error {
	failed := false
	for _, path := range paths {
		if path == "" {
			continue
		}
		info, err := os.Lstat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			failed = true
			continue
		}
		if info.IsDir() {
			failed = true
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			failed = true
		}
	}
	if failed {
		return errLoginRemove
	}
	return nil
}
