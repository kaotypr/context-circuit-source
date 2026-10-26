package report

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

func TestFileDeltaSkips(t *testing.T) {
	base := t.TempDir()
	branch := t.TempDir()
	files := map[string]string{
		"README.md":                 "workspace readme\n",
		"LICENSE":                   "terms\n",
		"assets/logo.txt":           "logo\n",
		".context-circuit/ids.yaml": "next: 1\n",
		"workspace.yaml":            "name: Acme\nversion: 3\n",
		"AGENTS.md":                 "old\n",
		"CODE_OF_CONDUCT.md":        "conduct\n",
		"CONTRIBUTING.md":           "contribute\n",
		"SECURITY.md":               "security\n",
		"CHANGELOG.md":              "log\n",
	}
	for rel, body := range files {
		mustWrite(t, filepath.Join(base, rel), body)
		mustWrite(t, filepath.Join(branch, rel), body)
	}
	mustWrite(t, filepath.Join(branch, "README.md"), "changed readme\n")
	mustWrite(t, filepath.Join(branch, "assets/logo.txt"), "new logo\n")
	mustWrite(t, filepath.Join(branch, ".context-circuit/ids.yaml"), "next: 2\n")
	mustWrite(t, filepath.Join(branch, "workspace.yaml"), "name: Acme\nversion: 9\n")
	mustWrite(t, filepath.Join(branch, "AGENTS.md"), "new\n")
	mustWrite(t, filepath.Join(branch, "notes/extra.md"), "added\n")
	changes, err := FileDelta(base, branch)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, change := range changes {
		got[change.Path] = change.Op
	}
	if len(got) != 2 || got["AGENTS.md"] != "change" || got["notes/extra.md"] != "add" {
		t.Fatalf("%v", got)
	}
	for _, change := range changes {
		if change.Path == "workspace.yaml" && strings.Contains(string(change.Body), "version: 9") {
			t.Fatal("version change was applied")
		}
	}
	mustWrite(t, filepath.Join(branch, "workspace.yaml"), "name: Other\nversion: 9\n")
	changes, err = FileDelta(base, branch)
	if err != nil {
		t.Fatal(err)
	}
	var yaml []byte
	for _, change := range changes {
		if change.Path == "workspace.yaml" {
			yaml = change.Body
		}
	}
	if !strings.Contains(string(yaml), "name: Other") || !strings.Contains(string(yaml), "version: 3") || strings.Contains(string(yaml), "version: 9") {
		t.Fatalf("workspace.yaml\n%s", yaml)
	}
}

func TestUpgradeCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go missing")
	}
	ctx := context.Background()
	root := sourceRoot(t)
	pin, err := world.Pin(root)
	if err != nil {
		t.Fatal(err)
	}
	buildRoot := t.TempDir()
	cli, err := world.BuildCLI(ctx, root, buildRoot, pin)
	if err != nil {
		t.Fatal(err)
	}
	base := initWorkspace(t, cli, os.Getenv("PATH"))
	plain, err := CheckWorkspace(ctx, cli, base, os.Getenv("PATH"))
	if err != nil {
		t.Fatal(err)
	}
	if !plain.OK {
		t.Fatalf("fresh workspace is not check-clean: %+v", plain.Findings)
	}

	t.Run("trivial", func(t *testing.T) {
		branch := filepath.Join(t.TempDir(), "branch")
		if err := CopyTree(base, branch); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(branch, "AGENTS.md")
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(body, []byte("\nContributors read this extra sentence.\n")...), 0o644); err != nil {
			t.Fatal(err)
		}
		dest := filepath.Join(t.TempDir(), "dest")
		if err := CopyTree(base, dest); err != nil {
			t.Fatal(err)
		}
		changes, err := FileDelta(base, branch)
		if err != nil {
			t.Fatal(err)
		}
		if err := ApplyDelta(dest, changes); err != nil {
			t.Fatal(err)
		}
		got, err := CheckWorkspace(ctx, cli, dest, os.Getenv("PATH"))
		if err != nil {
			t.Fatal(err)
		}
		if !got.OK {
			t.Fatalf("trivial delta: %+v", got.Findings)
		}
		section := renderUpgrade(changes, &got, nil)
		if !strings.Contains(section, "file-level delta") || !strings.Contains(section, "not the host's judgement") {
			t.Fatalf("%s", section)
		}
	})

	t.Run("frontmatter", func(t *testing.T) {
		branch := filepath.Join(t.TempDir(), "branch")
		if err := CopyTree(base, branch); err != nil {
			t.Fatal(err)
		}
		mustWrite(t, filepath.Join(branch, "intent", "i100-demo.md"), "---\nid: i100\ncreated_by: maya\ncreated_at: 2026-01-15T09:00:00Z\n\n# Request\n")
		dest := filepath.Join(t.TempDir(), "dest")
		if err := CopyTree(base, dest); err != nil {
			t.Fatal(err)
		}
		changes, err := FileDelta(base, branch)
		if err != nil {
			t.Fatal(err)
		}
		if err := ApplyDelta(dest, changes); err != nil {
			t.Fatal(err)
		}
		got, err := CheckWorkspace(ctx, cli, dest, os.Getenv("PATH"))
		if err != nil {
			t.Fatal(err)
		}
		var text strings.Builder
		for _, finding := range got.Findings {
			text.WriteString(finding.Issue)
			text.WriteByte('\n')
		}
		if !strings.Contains(text.String(), "frontmatter") {
			t.Fatalf("findings:\n%s", text.String())
		}
		section := renderUpgrade(changes, &got, nil)
		if !strings.Contains(section, "frontmatter") || !strings.Contains(section, "file-level delta") {
			t.Fatalf("%s", section)
		}
	})
}

func initWorkspace(t *testing.T, cli, pathEnv string) string {
	t.Helper()
	parent := t.TempDir()
	ws := filepath.Join(parent, "ws")
	home := filepath.Join(parent, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	gitconfig := filepath.Join(parent, "gitconfig")
	if err := os.WriteFile(gitconfig, []byte("[user]\n\tname = Agent Test\n\temail = agent-test@example.invalid\n[commit]\n\tgpgsign = false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	env := []string{
		"PATH=" + pathEnv,
		"HOME=" + home,
		"GIT_CONFIG_GLOBAL=" + gitconfig,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=Agent Test",
		"GIT_AUTHOR_EMAIL=agent-test@example.invalid",
		"GIT_COMMITTER_NAME=Agent Test",
		"GIT_COMMITTER_EMAIL=agent-test@example.invalid",
		"GIT_AUTHOR_DATE=2026-01-15T09:00:00Z",
		"GIT_COMMITTER_DATE=2026-01-15T09:00:00Z",
	}
	run := func(dir string, argv ...string) {
		t.Helper()
		res, err := proc.Run(context.Background(), proc.Spec{Argv: argv, Dir: dir, Env: env, Timeout: 2 * time.Minute})
		if err != nil || res.ExitCode != 0 {
			t.Fatalf("%v: %v exit %d\n%s\n%s", argv, err, res.ExitCode, res.Stdout, res.Stderr)
		}
	}
	run(parent, "git", "init", "-q", "-b", "main", ws)
	run(parent, cli, "--workspace", ws, "init", "--name", "Acme", "--purpose", "Invoicing for small businesses", "--member", "maya", "--member-name", "Maya")
	run(ws, "git", "add", "-A")
	run(ws, "git", "commit", "-qm", "chore: initialize workspace")
	run(ws, cli, "--workspace", ws, "workspace", "connect", "--base", "main")
	return ws
}

func sourceRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}
