package cli

import (
	"path/filepath"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// The pin step's internals, for the tests in pin_test.go.
var (
	ReleaseLine    = releaseLine
	DifferentLines = differentLines
	ValidPin       = validPin
	FindStore      = findStore
	Installed      = installed
)

const (
	ReservedLabels = reservedLabels
	LabelPattern   = labelPattern
)

// SetExecutable makes the pin step see path as this binary's location.
func SetExecutable(path string) (restore func()) {
	old := executable
	executable = func() (string, error) { return path, nil }
	return func() { executable = old }
}

// PinSteps renders the next steps a warning or refusal names, for any platform.
// target is the installed pinned binary, "" when it is not installed.
func PinSteps(s *workspace.Store, pin, running, goos, bin, target string, args []string) string {
	p := pinState{root: s.Root, pin: pin, running: running, goos: goos, bin: bin, target: target, args: args}
	if bin != "" {
		p.store = filepath.Join(bin, versionStore)
	}
	return p.nextSteps(s)
}
