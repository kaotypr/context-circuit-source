//go:build !darwin && !linux && !windows

package cow

import "errors"

func cloneFile(_, _ string) error { return errors.New("native CoW is unavailable on this platform") }
