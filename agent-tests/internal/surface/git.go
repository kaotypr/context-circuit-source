package surface

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
)

func gitRun(ctx context.Context, checkout, pathEnv string, args ...string) (string, int, error) {
	res, err := proc.Run(ctx, proc.Spec{
		Argv:    append([]string{"git", "-C", checkout}, args...),
		Env:     gitEnv(pathEnv),
		Timeout: 2 * time.Minute,
	})
	if err != nil {
		return "", -1, err
	}
	return string(res.Stdout), res.ExitCode, nil
}

func gitEnv(pathEnv string) []string {
	return []string{
		"PATH=" + pathEnv,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_TERMINAL_PROMPT=0",
	}
}

func gitText(ctx context.Context, checkout, pathEnv string, args ...string) (string, error) {
	out, code, err := gitRun(ctx, checkout, pathEnv, args...)
	if err != nil {
		return "", err
	}
	if code != 0 {
		return "", fmt.Errorf("git %s: exit %d\n%s", strings.Join(args, " "), code, out)
	}
	return out, nil
}

func revParse(ctx context.Context, checkout, pathEnv, ref string) (string, error) {
	out, err := gitText(ctx, checkout, pathEnv, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("rev-parse %s: %w", ref, err)
	}
	return strings.TrimSpace(out), nil
}

func nameStatus(ctx context.Context, checkout, pathEnv, base, branch string) ([]string, error) {
	out, code, err := gitRun(ctx, checkout, pathEnv, "diff", "--name-status", "--find-renames", base, branch)
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return nil, fmt.Errorf("git diff --name-status: exit %d\n%s", code, out)
	}
	return parseNameStatus(out), nil
}

func parseNameStatus(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 2 {
			continue
		}
		status := fields[0]
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			out = append(out, fields[1:]...)
			continue
		}
		out = append(out, fields[1])
	}
	return out
}

func gitShow(ctx context.Context, checkout, pathEnv, rev, file string) (string, bool, error) {
	out, code, err := gitRun(ctx, checkout, pathEnv, "show", rev+":"+file)
	if err != nil {
		return "", false, err
	}
	if code != 0 {
		return "", false, nil
	}
	return out, true, nil
}

func writeGitConfig(run string) error {
	text := "[user]\n\tname = Agent Test\n\temail = agent-test@example.invalid\n[commit]\n\tgpgsign = false\n"
	return os.WriteFile(filepath.Join(run, "gitconfig"), []byte(text), 0o644)
}
