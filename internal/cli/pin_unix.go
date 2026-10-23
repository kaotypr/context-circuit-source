//go:build !windows

package cli

import (
	"io"
	"syscall"
)

// execCall replaces this process; a test replaces it, because the real one
// would replace the test.
var execCall = syscall.Exec

// handOff runs the pinned CLI in place of this one: the same process, working
// directory, streams, and arguments. It returns only when that failed, or when
// a test stood in for the exec.
func handOff(target string, args, env []string, _, _ io.Writer) (int, error) {
	if err := execCall(target, append([]string{target}, args...), env); err != nil {
		return 1, err
	}
	return 0, nil
}
