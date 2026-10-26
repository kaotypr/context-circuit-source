package world

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
)

// BuildCLI compiles ./cmd/context-circuit from srcDir into the launcher
// store under runRoot. The child environment is the allowlist below.
// GOMODCACHE is the caller's value when one is set, otherwise `go env
// GOMODCACHE`. Using that cache only reads module downloads that are
// already there.
func BuildCLI(ctx context.Context, srcDir, runRoot, pin string) (string, error) {
	if pin == "" {
		return "", fmt.Errorf("pin is empty")
	}
	cache, err := moduleCache(ctx)
	if err != nil {
		return "", err
	}
	store := filepath.Join(runRoot, "bin", ".context-circuit-versions", pin+"-"+runtime.GOOS+"-"+runtime.GOARCH)
	if err := os.MkdirAll(store, 0o755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(runRoot, "gocache"), 0o755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(runRoot, "gopath"), 0o755); err != nil {
		return "", err
	}
	cliPath := filepath.Join(store, "context-circuit-cli")
	res, err := proc.Run(ctx, proc.Spec{
		Argv: []string{
			"go", "build", "-trimpath",
			"-ldflags", "-X main.version=" + pin,
			"-o", cliPath,
			"./cmd/context-circuit",
		},
		Dir: srcDir,
		Env: []string{
			"PATH=" + os.Getenv("PATH"),
			"CGO_ENABLED=0",
			"GOCACHE=" + filepath.Join(runRoot, "gocache"),
			"GOFLAGS=",
			"GOPATH=" + filepath.Join(runRoot, "gopath"),
			"GOMODCACHE=" + cache,
		},
	})
	if err != nil {
		return "", fmt.Errorf("go build: %w\n%s", err, res.Stderr)
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("go build: exit %d\n%s", res.ExitCode, res.Stderr)
	}
	return cliPath, nil
}

// Pin reads the checkout's CLI_VERSION file and trims surrounding whitespace.
func Pin(checkout string) (string, error) {
	b, err := os.ReadFile(filepath.Join(checkout, "CLI_VERSION"))
	if err != nil {
		return "", err
	}
	pin := strings.TrimSpace(string(b))
	if pin == "" {
		return "", fmt.Errorf("CLI_VERSION is empty")
	}
	return pin, nil
}

func moduleCache(ctx context.Context) (string, error) {
	if v := strings.TrimSpace(os.Getenv("GOMODCACHE")); v != "" {
		return v, nil
	}
	res, err := proc.Run(ctx, proc.Spec{
		Argv: []string{"go", "env", "GOMODCACHE"},
		Env:  []string{"PATH=" + os.Getenv("PATH")},
	})
	if err != nil {
		return "", fmt.Errorf("go env GOMODCACHE: %w", err)
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("go env GOMODCACHE: exit %d\n%s", res.ExitCode, res.Stderr)
	}
	cache := strings.TrimSpace(string(res.Stdout))
	if cache == "" {
		return "", fmt.Errorf("go env GOMODCACHE returned an empty path")
	}
	return cache, nil
}
