//go:build !windows

package cli_test

// shortPath reports that 8.3 short names exist only on Windows.
func shortPath(string) (string, bool) { return "", false }
