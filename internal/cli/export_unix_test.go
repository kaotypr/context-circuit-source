//go:build !windows

package cli

// StubHandOff records a hand-off instead of replacing the test process.
func StubHandOff(record func(target string, argv, env []string)) (restore func()) {
	old := execCall
	execCall = func(target string, argv, env []string) error {
		record(target, argv, env)
		return nil
	}
	return func() { execCall = old }
}
