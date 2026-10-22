//go:build windows

package cli

import "os/exec"

// StubHandOff records a hand-off instead of starting the pinned CLI.
func StubHandOff(record func(target string, argv, env []string)) (restore func()) {
	old := runChild
	runChild = func(child *exec.Cmd) error {
		record(child.Path, child.Args, child.Env)
		return nil
	}
	return func() { runChild = old }
}
