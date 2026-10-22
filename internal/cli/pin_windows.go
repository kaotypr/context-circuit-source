//go:build windows

package cli

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

// runChild runs the pinned CLI and waits for it; a test replaces it.
var runChild = func(child *exec.Cmd) error { return child.Run() }

// handOff runs the pinned CLI as a child, since Windows has no exec, with this
// process's input, its streams, arguments and working directory, and returns
// the child's exit code. exec.Command rather than CommandContext: Ctrl-C
// reaches both processes through the console, and this one keeps waiting for
// the child to finish rather than killing it.
func handOff(target string, args, env []string, out, errOut io.Writer) (int, error) {
	child := exec.Command(target, args...)
	child.Env = env
	child.Stdin = os.Stdin
	child.Stdout = out
	child.Stderr = errOut
	err := runChild(child)
	if exit := (*exec.ExitError)(nil); errors.As(err, &exit) {
		return exit.ExitCode(), nil
	}
	if err != nil {
		return 1, err
	}
	return 0, nil
}
