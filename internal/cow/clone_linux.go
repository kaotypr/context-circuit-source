package cow

import (
	"golang.org/x/sys/unix"
	"os"
)

func cloneFile(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	return unix.IoctlFileClone(int(out.Fd()), int(in.Fd()))
}
