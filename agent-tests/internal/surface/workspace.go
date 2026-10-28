package surface

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

// helpTopics are the top-level help texts. The full `help` output still
// shows a new command's signature when this list does not name it yet.
var helpTopics = []string{
	"init", "status", "check", "member", "knowledge", "repo", "record",
	"context", "worktree", "template", "agent", "workspace", "version",
}

var intentIDPattern = regexp.MustCompile(`(?m)^id:\s*"?(i[0-9]+)"?\s*$`)

type captured struct {
	files   map[string][]byte
	outputs []capturedOutput
	norm    Norm
}

type capturedOutput struct {
	name   string
	text   string
	stdout string
	code   int
}

func captureSide(ctx context.Context, checkout, ref, pathEnv, home string) (captured, error) {
	run, err := world.NewRunRoot("")
	if err != nil {
		return captured{}, err
	}
	defer os.RemoveAll(run)
	ws, cli, pin, err := initSide(ctx, checkout, ref, pathEnv, home, run)
	if err != nil {
		return captured{}, err
	}
	outputs, err := runScript(ctx, cli, ws, sideEnv(run, pathEnv))
	if err != nil {
		return captured{}, err
	}
	files, err := snapshotTree(ws)
	if err != nil {
		return captured{}, err
	}
	return captured{
		files:   files,
		outputs: outputs,
		norm:    NewNorm(pathTokens(run, ws), pin),
	}, nil
}

// OpenInitWorkspace archives ref, builds the CLI, and inits a workspace
// the same way surface-diff does. cleanup removes the run root, including
// the CLI binary. The binary is stamped with that ref's CLI_VERSION.
func OpenInitWorkspace(ctx context.Context, checkout, ref, pathEnv, home string) (ws, cli string, cleanup func(), err error) {
	if pathEnv == "" {
		pathEnv = os.Getenv("PATH")
	}
	run, err := world.NewRunRoot("")
	if err != nil {
		return "", "", nil, err
	}
	cleanup = func() { _ = os.RemoveAll(run) }
	ws, cli, _, err = initSide(ctx, checkout, ref, pathEnv, home, run)
	if err != nil {
		cleanup()
		return "", "", nil, err
	}
	return ws, cli, cleanup, nil
}

func initSide(ctx context.Context, checkout, ref, pathEnv, home, run string) (ws, cli, pin string, err error) {
	if err := world.Guard(run, checkout, home); err != nil {
		return "", "", "", err
	}
	for _, dir := range []string{"home", "tmp", "gocache"} {
		if err := os.MkdirAll(filepath.Join(run, dir), 0o755); err != nil {
			return "", "", "", err
		}
	}
	if err := writeGitConfig(run); err != nil {
		return "", "", "", err
	}
	src := filepath.Join(run, "src")
	if err := world.ExtractArchive(ctx, checkout, ref, src); err != nil {
		return "", "", "", err
	}
	pin, err = archiveCLIVersion(src)
	if err != nil {
		return "", "", "", err
	}
	cli, err = world.BuildCLI(ctx, src, run, pin)
	if err != nil {
		return "", "", "", err
	}
	ws = filepath.Join(run, "ws")
	if err := prepareWorkspace(ctx, cli, ws, sideEnv(run, pathEnv)); err != nil {
		return "", "", "", err
	}
	return ws, cli, pin, nil
}

// archiveCLIVersion reads CLI_VERSION from an extracted archive.
func archiveCLIVersion(dir string) (string, error) {
	return world.Pin(dir)
}

func pathTokens(run, ws string) []PathToken {
	var out []PathToken
	add := func(p, token string) {
		if p == "" {
			return
		}
		out = append(out, PathToken{Path: p, Token: token})
		resolved, err := filepath.EvalSymlinks(p)
		if err == nil && resolved != "" && resolved != p {
			out = append(out, PathToken{Path: resolved, Token: token})
		}
	}
	add(ws, "$WS")
	add(run, "$RUN")
	return out
}

func sideEnv(run, pathEnv string) []string {
	return []string{
		"PATH=" + pathEnv,
		"HOME=" + filepath.Join(run, "home"),
		"TMPDIR=" + filepath.Join(run, "tmp"),
		"GOCACHE=" + filepath.Join(run, "gocache"),
		"LANG=C.UTF-8",
		"LC_ALL=C.UTF-8",
		"GIT_CONFIG_GLOBAL=" + filepath.Join(run, "gitconfig"),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0",
		"GIT_AUTHOR_NAME=Agent Test",
		"GIT_AUTHOR_EMAIL=agent-test@example.invalid",
		"GIT_COMMITTER_NAME=Agent Test",
		"GIT_COMMITTER_EMAIL=agent-test@example.invalid",
		"GIT_AUTHOR_DATE=2026-01-15T09:00:00Z",
		"GIT_COMMITTER_DATE=2026-01-15T09:00:00Z",
		"NO_COLOR=1",
	}
}

func prepareWorkspace(ctx context.Context, cli, ws string, env []string) error {
	if err := runCmd(ctx, "git init", env, "", "git", "init", "-q", "-b", "main", ws); err != nil {
		return err
	}
	if err := runCmd(ctx, "init", env, filepath.Dir(ws), cli, "--workspace", ws, "init",
		"--name", "Acme", "--purpose", "Invoicing for small businesses",
		"--member", "maya", "--member-name", "Maya"); err != nil {
		return err
	}
	api := filepath.Join(ws, "repositories", "api")
	return runCmd(ctx, "repo init", env, ws, cli, "--workspace", ws, "repo", "init",
		"--id", "api", "--path", api, "--base", "main")
}

func runScript(ctx context.Context, cli, ws string, env []string) ([]capturedOutput, error) {
	var outputs []capturedOutput
	run := func(name string, args ...string) (capturedOutput, error) {
		argv := append([]string{cli, "--workspace", ws}, args...)
		text, stdout, code, err := runCapture(ctx, env, ws, argv...)
		if err != nil {
			return capturedOutput{}, fmt.Errorf("%s: %w", name, err)
		}
		return capturedOutput{name: name, text: text, stdout: stdout, code: code}, nil
	}
	full, err := run("help", "help")
	if err != nil {
		return nil, err
	}
	outputs = append(outputs, full)
	for _, topic := range helpTopics {
		out, err := run("help-"+topic, "help", topic)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, out)
	}
	intent, err := run("record-create-intent", "record", "create", "--kind", "intent", "--slug", "surface-intent", "--title", "Surface intent")
	if err != nil {
		return nil, err
	}
	outputs = append(outputs, intent)
	if id, ok := parseIntentID(intent.stdout); ok {
		plan, err := run("record-create-plan", "record", "create", "--kind", "plan", "--slug", "surface-plan", "--title", "Surface plan", "--intent", id, "--repo", "api")
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, plan)
	}
	for _, step := range []struct {
		name string
		args []string
	}{
		{"check", []string{"check"}},
		{"check-json", []string{"--json", "check"}},
		{"agent-setup", []string{"agent", "setup"}},
		{"status", []string{"status"}},
		{"status-json", []string{"--json", "status"}},
	} {
		out, err := run(step.name, step.args...)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, out)
	}
	return outputs, nil
}

func parseIntentID(stdout string) (string, bool) {
	m := intentIDPattern.FindStringSubmatch(stdout)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func runCmd(ctx context.Context, name string, env []string, dir string, argv ...string) error {
	text, _, code, err := runCapture(ctx, env, dir, argv...)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if code != 0 {
		return fmt.Errorf("%s: exit %d\n%s", name, code, text)
	}
	return nil
}

func runCapture(ctx context.Context, env []string, dir string, argv ...string) (text, stdout string, code int, err error) {
	res, err := proc.Run(ctx, proc.Spec{
		Argv:    argv,
		Dir:     dir,
		Env:     env,
		Timeout: 2 * time.Minute,
	})
	if err != nil {
		return "", "", -1, err
	}
	stdout = string(res.Stdout)
	return formatCapture(stdout, string(res.Stderr), res.ExitCode), stdout, res.ExitCode, nil
}

func formatCapture(stdout, stderr string, code int) string {
	var b strings.Builder
	b.WriteString(stdout)
	if strings.TrimSpace(stderr) != "" {
		if stdout != "" && !strings.HasSuffix(stdout, "\n") {
			b.WriteByte('\n')
		}
		b.WriteString("--- stderr ---\n")
		b.WriteString(stderr)
		if !strings.HasSuffix(stderr, "\n") {
			b.WriteByte('\n')
		}
	}
	fmt.Fprintf(&b, "--- exit %d ---\n", code)
	return b.String()
}

func snapshotTree(root string) (map[string][]byte, error) {
	out := map[string][]byte{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		body, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(rel)] = body
		return nil
	})
	return out, err
}

func writeTree(dir string, files map[string][]byte) error {
	for rel, body := range files {
		target := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, body, 0o644); err != nil {
			return err
		}
	}
	return nil
}
