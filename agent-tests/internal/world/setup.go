package world

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
)

// Fixture commit timestamps. They apply only while the fixture commits
// are created, and are absent from the environment a host receives.
const (
	fixtureAuthorDate    = "2026-01-15T09:00:00Z"
	fixtureCommitterDate = "2026-01-15T09:00:00Z"
)

// SetupInput is everything the common world needs besides the process
// environment of whoever started the runner.
type SetupInput struct {
	RunRoot     string
	SrcDir      string
	CLIPath     string
	FixturesDir string
	Pin         string
	PATH        string
}

// World is one isolated fixture workspace.
// Env is the extra map passed to hosts.BaseEnv. Together with that
// allowlist it is the environment a host receives from the world.
// Adapters add their own home directories on top.
type World struct {
	RunRoot  string
	WS       string
	API      string
	Decoy    string
	DecoyLog string
	Env      map[string]string
}

// Setup builds the common world under RunRoot. Every command is a
// proc.Run with an explicit environment. A non-zero exit stops the
// setup and includes that step's stderr. Nothing is written outside
// RunRoot.
func Setup(ctx context.Context, p SetupInput) (World, error) {
	if p.RunRoot == "" {
		return World{}, fmt.Errorf("run root is empty")
	}
	if p.CLIPath == "" {
		return World{}, fmt.Errorf("cli path is empty")
	}
	if p.FixturesDir == "" {
		if p.SrcDir == "" {
			return World{}, fmt.Errorf("fixtures directory is empty")
		}
		p.FixturesDir = filepath.Join(p.SrcDir, "agent-tests", "fixtures")
	}
	run := p.RunRoot
	ws := filepath.Join(run, "ws")
	api := filepath.Join(ws, "repositories", "api")
	decoy := filepath.Join(run, "decoy")
	decoyLog := filepath.Join(run, "decoy.log")
	w := World{
		RunRoot:  run,
		WS:       ws,
		API:      api,
		Decoy:    decoy,
		DecoyLog: decoyLog,
		Env: map[string]string{
			"CONTEXT_CIRCUIT_BIN_DIR": filepath.Join(run, "bin"),
			"RUN":                     run,
			"WS":                      ws,
		},
	}
	dirs := []string{
		filepath.Join(run, "bin", ".context-circuit-versions"),
		filepath.Join(run, "remotes"),
		decoy,
		filepath.Join(run, "gocache"),
		filepath.Join(run, "out"),
		filepath.Join(run, "tmp"),
		filepath.Join(run, "setup-home"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return World{}, err
		}
	}
	gitconfig := "[user]\n\tname = Agent Test\n\temail = agent-test@example.invalid\n[commit]\n\tgpgsign = false\n"
	if err := os.WriteFile(filepath.Join(run, "gitconfig"), []byte(gitconfig), 0o644); err != nil {
		return World{}, err
	}
	script := "#!/bin/sh\necho \"$0 $*\" >> '" + decoyLog + "'\nexit 127\n"
	if err := os.WriteFile(filepath.Join(decoy, "context-circuit-cli"), []byte(script), 0o755); err != nil {
		return World{}, err
	}

	base := envBase(p, run)
	dated := append(append([]string{}, base...),
		"GIT_AUTHOR_DATE="+fixtureAuthorDate,
		"GIT_COMMITTER_DATE="+fixtureCommitterDate,
	)

	steps := []struct {
		name string
		spec proc.Spec
	}{
		{"git init ws remote", proc.Spec{
			Argv: []string{"git", "init", "-q", "--bare", "-b", "main", filepath.Join(run, "remotes", "ws.git")},
			Env:  dated, Dir: run,
		}},
		{"git init workspace", proc.Spec{
			Argv: []string{"git", "init", "-q", "-b", "main", ws},
			Env:  dated, Dir: run,
		}},
		{"git remote workspace", proc.Spec{
			Argv: []string{"git", "-C", ws, "remote", "add", "origin", filepath.Join(run, "remotes", "ws.git")},
			Env:  dated, Dir: run,
		}},
		{"cli init", proc.Spec{
			Argv: []string{p.CLIPath, "--workspace", ws, "init", "--name", "Acme", "--purpose", "Invoicing for small businesses", "--member", "maya", "--member-name", "Maya"},
			Env:  dated, Dir: run,
		}},
		{"git add workspace", proc.Spec{
			Argv: []string{"git", "-C", ws, "add", "-A"},
			Env:  dated, Dir: run,
		}},
		{"git commit workspace", proc.Spec{
			Argv: []string{"git", "-C", ws, "commit", "-qm", "chore: initialize workspace"},
			Env:  dated, Dir: run,
		}},
		{"workspace connect", proc.Spec{
			Argv: []string{"sh", filepath.Join(ws, ".context-circuit", "cli"), "workspace", "connect", "--base", "main"},
			Env:  dated, Dir: ws,
		}},
		{"git init api remote", proc.Spec{
			Argv: []string{"git", "init", "-q", "--bare", "-b", "main", filepath.Join(run, "remotes", "api.git")},
			Env:  dated, Dir: run,
		}},
		{"git init api", proc.Spec{
			Argv: []string{"git", "init", "-q", "-b", "main", api},
			Env:  dated, Dir: run,
		}},
	}
	for _, step := range steps {
		if err := runStep(ctx, step.name, step.spec); err != nil {
			return World{}, err
		}
	}
	if err := copyTree(filepath.Join(p.FixturesDir, "api"), api); err != nil {
		return World{}, fmt.Errorf("copy api fixture: %w", err)
	}
	apiSteps := []struct {
		name string
		spec proc.Spec
	}{
		{"git add api", proc.Spec{Argv: []string{"git", "-C", api, "add", "-A"}, Env: dated, Dir: run}},
		{"git commit api", proc.Spec{Argv: []string{"git", "-C", api, "commit", "-qm", "feat(invoice): text summary export"}, Env: dated, Dir: run}},
		{"git remote api", proc.Spec{Argv: []string{"git", "-C", api, "remote", "add", "origin", filepath.Join(run, "remotes", "api.git")}, Env: dated, Dir: run}},
		{"git push api", proc.Spec{Argv: []string{"git", "-C", api, "push", "-q", "origin", "main"}, Env: dated, Dir: run}},
		{"repo connect", proc.Spec{Argv: []string{"sh", filepath.Join(ws, ".context-circuit", "cli"), "repo", "connect", "--id", "api", "--path", api, "--base", "main"}, Env: dated, Dir: ws}},
	}
	for _, step := range apiSteps {
		if err := runStep(ctx, step.name, step.spec); err != nil {
			return World{}, err
		}
	}
	if err := os.MkdirAll(filepath.Join(ws, "context", "domains"), 0o755); err != nil {
		return World{}, err
	}
	if err := copyFile(filepath.Join(p.FixturesDir, "context", "invoices.md"), filepath.Join(ws, "context", "domains", "invoices.md")); err != nil {
		return World{}, fmt.Errorf("copy invoices note: %w", err)
	}
	if err := appendFile(filepath.Join(ws, "context", "INDEX.md"), filepath.Join(p.FixturesDir, "context", "index-append.md")); err != nil {
		return World{}, fmt.Errorf("append catalog: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(ws, "sources"), 0o755); err != nil {
		return World{}, err
	}
	if err := copyFile(filepath.Join(p.FixturesDir, "sources", "old-brief.md"), filepath.Join(ws, "sources", "old-brief.md")); err != nil {
		return World{}, fmt.Errorf("copy old brief: %w", err)
	}
	final := []struct {
		name string
		spec proc.Spec
	}{
		{"git add fixture", proc.Spec{Argv: []string{"git", "-C", ws, "add", "-A"}, Env: dated, Dir: run}},
		{"git commit fixture", proc.Spec{Argv: []string{"git", "-C", ws, "commit", "-qm", "chore: fixture"}, Env: dated, Dir: run}},
		{"git push workspace", proc.Spec{Argv: []string{"git", "-C", ws, "push", "-q", "origin", "main"}, Env: dated, Dir: run}},
		// The date variables are fixture-only. check runs without them.
		{"check", proc.Spec{
			Argv: []string{"sh", filepath.Join(ws, ".context-circuit", "cli"), "--json", "check"},
			Env:  base, Dir: ws,
		}},
	}
	for _, step := range final {
		if step.name == "check" {
			res, err := proc.Run(ctx, step.spec)
			if err != nil {
				return World{}, fmt.Errorf("check: %w\n%s", err, res.Stderr)
			}
			if res.ExitCode != 0 {
				return World{}, fmt.Errorf("check: exit %d\n%s\n%s", res.ExitCode, res.Stdout, res.Stderr)
			}
			if err := os.WriteFile(filepath.Join(run, "out", "check-before.json"), res.Stdout, 0o644); err != nil {
				return World{}, err
			}
			continue
		}
		if err := runStep(ctx, step.name, step.spec); err != nil {
			return World{}, err
		}
	}
	return w, nil
}

func envBase(p SetupInput, run string) []string {
	path := p.PATH
	if decoy := filepath.Join(run, "decoy"); path == "" {
		path = decoy
	} else {
		path = decoy + ":" + path
	}
	return []string{
		"PATH=" + path,
		"HOME=" + filepath.Join(run, "setup-home"),
		"TMPDIR=" + filepath.Join(run, "tmp"),
		"GOCACHE=" + filepath.Join(run, "gocache"),
		"LANG=C.UTF-8",
		"GIT_CONFIG_GLOBAL=" + filepath.Join(run, "gitconfig"),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0",
		"NO_COLOR=1",
		"CONTEXT_CIRCUIT_BIN_DIR=" + filepath.Join(run, "bin"),
		"RUN=" + run,
		"WS=" + filepath.Join(run, "ws"),
	}
}

func runStep(ctx context.Context, name string, spec proc.Spec) error {
	res, err := proc.Run(ctx, spec)
	if err != nil {
		return fmt.Errorf("%s: %w\n%s", name, err, res.Stderr)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("%s: exit %d\n%s", name, res.ExitCode, res.Stderr)
	}
	return nil
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("fixture %s is not a regular file", path)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func appendFile(dst, src string) error {
	extra, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if len(extra) > 0 && extra[0] != '\n' {
		if _, err := f.Write([]byte("\n")); err != nil {
			return err
		}
	}
	text := string(extra)
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	_, err = f.WriteString(text)
	return err
}
