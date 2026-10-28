//go:build unix

package proc

import "syscall"

func processAlive(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}

func killPid(pid int) {
	_ = syscall.Kill(pid, syscall.SIGKILL)
}
