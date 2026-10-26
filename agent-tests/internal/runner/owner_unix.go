//go:build unix

package runner

import (
	"os"
	"syscall"
)

func ownedByCurrentUser(info os.FileInfo) bool {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok || st == nil {
		return false
	}
	return st.Uid == uint32(os.Getuid())
}
