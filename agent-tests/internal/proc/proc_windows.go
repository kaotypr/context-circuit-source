//go:build windows

package proc

import (
	"os"
	"os/exec"
)

func prepareCmd(cmd *exec.Cmd) {}

// signal stops this process only. Windows has no POSIX process group here.
func signal(proc *os.Process, hard bool) {
	_ = hard
	if proc != nil {
		_ = proc.Kill()
	}
}
