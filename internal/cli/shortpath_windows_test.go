package cli_test

import (
	"syscall"
)

// shortPath is path's 8.3 short form, and false where the volume keeps no
// short names for it.
func shortPath(path string) (string, bool) {
	long, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", false
	}
	buf := make([]uint16, syscall.MAX_PATH)
	for {
		n, err := syscall.GetShortPathName(long, &buf[0], uint32(len(buf)))
		if err != nil || n == 0 {
			return "", false
		}
		// A buffer too small gets back the size it needs, terminator included.
		if int(n) < len(buf) {
			short := syscall.UTF16ToString(buf[:n])
			return short, short != path
		}
		buf = make([]uint16, n)
	}
}
