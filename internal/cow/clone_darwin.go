package cow

import "golang.org/x/sys/unix"

func cloneFile(source, destination string) error {
	return unix.Clonefile(source, destination, unix.CLONE_NOFOLLOW)
}
