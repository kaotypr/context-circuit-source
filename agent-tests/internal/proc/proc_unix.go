//go:build unix

package proc

import (
	"os"
	"os/exec"
	"syscall"
)

func prepareCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func signal(proc *os.Process, hard bool) {
	if proc == nil {
		return
	}
	sig := syscall.SIGTERM
	if hard {
		sig = syscall.SIGKILL
	}
	if err := syscall.Kill(-proc.Pid, sig); err != nil {
		if hard {
			_ = proc.Kill()
			return
		}
		_ = proc.Signal(syscall.SIGTERM)
	}
}
