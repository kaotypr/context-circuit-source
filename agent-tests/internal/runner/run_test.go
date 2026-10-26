package runner

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
)

func TestSourcesDoNotInheritCallerEnv(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	checks := []struct {
		dir  string
		main bool
	}{
		{filepath.Join(root, "internal", "runner"), false},
		{filepath.Join(root, "internal", "world"), false},
		{filepath.Join(root, "cmd", "agent-tests"), true},
	}
	environ := "os.Environ("
	home := "os.Getenv(\"HOME\""
	for _, check := range checks {
		entries, err := os.ReadDir(check.dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			body, err := os.ReadFile(filepath.Join(check.dir, name))
			if err != nil {
				t.Fatal(err)
			}
			text := string(body)
			if strings.Contains(text, home) {
				t.Fatalf("%s reads HOME from the process environment", name)
			}
			n := strings.Count(text, environ)
			if check.main {
				if n != 1 {
					t.Fatalf("%s calls os.Environ %d times, want 1", name, n)
				}
				continue
			}
			if n != 0 {
				t.Fatalf("%s calls os.Environ", name)
			}
		}
	}
}

func TestDigestFixturesMatchesRecord(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := digestFixtures(dir)
	if err != nil {
		t.Fatal(err)
	}
	want, err := record.Digest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != want || !strings.HasPrefix(got, "sha256:") {
		t.Fatalf("digest %s want %s", got, want)
	}
}

func TestWorldPinReadsVersionFile(t *testing.T) {
	root := repoRoot(t)
	got, err := worldPin(root)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, "CLI_VERSION"))
	if err != nil {
		t.Fatal(err)
	}
	if got != strings.TrimSpace(string(b)) || got == "" {
		t.Fatalf("pin %q", got)
	}
}

func TestDryRunNoScenarios(t *testing.T) {
	dir := t.TempDir()
	hostsFile := writeHostsFile(t, dir, []hostRow{{
		name: "claude", command: filepath.Join(dir, "missing-stub"), model: "claude-opus-5-5", version: "0", timeout: 5,
	}})
	for _, scenarios := range []string{filepath.Join(dir, "missing"), filepath.Join(dir, "empty")} {
		if strings.HasSuffix(scenarios, "empty") {
			if err := os.MkdirAll(scenarios, 0o755); err != nil {
				t.Fatal(err)
			}
		}
		var buf bytes.Buffer
		err := Run(context.Background(), RunOptions{
			Runs:         1,
			HostsFile:    hostsFile,
			ScenariosDir: scenarios,
			DryRun:       true,
			Stdout:       &buf,
		})
		if err == nil || !strings.Contains(err.Error(), "no scenarios found") {
			t.Fatalf("%s: %v", scenarios, err)
		}
		if !strings.Contains(buf.String(), "no scenarios found") {
			t.Fatalf("stdout %q", buf.String())
		}
	}
}

func TestDryRunPrintsPlan(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "started")
	stub := writeTouchStub(t, dir, marker)
	hostsFile := writeHostsFile(t, dir, []hostRow{{
		name: "claude", command: stub, model: "claude-opus-5-5", version: "0", timeout: 5,
	}})
	scenarios := writeScenario(t, dir, "marker")
	var buf bytes.Buffer
	err := Run(context.Background(), RunOptions{
		Ref:          "HEAD",
		BaseRef:      "base-ref",
		Runs:         2,
		HostsFile:    hostsFile,
		ScenariosDir: scenarios,
		DryRun:       true,
		Stdout:       &buf,
	})
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	want := strings.Join([]string{
		"plan: 4 combination(s)",
		"marker x claude x base x r1",
		"marker x claude x branch x r1",
		"marker x claude x branch x r2",
		"marker x claude x base x r2",
		"",
	}, "\n")
	if got != want {
		t.Fatalf("plan\n%s", got)
	}
	if strings.Contains(got, "re-run with --yes") {
		t.Fatal("dry-run asked for confirmation")
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("dry-run started a host")
	}
}

func TestRunRefusesWithoutYes(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "started")
	stub := writeTouchStub(t, dir, marker)
	hostsFile := writeHostsFile(t, dir, []hostRow{{
		name: "claude", command: stub, model: "claude-opus-5-5", version: "0", timeout: 5,
	}})
	scenarios := writeScenario(t, dir, "marker")
	var buf bytes.Buffer
	err := Run(context.Background(), RunOptions{
		Runs:         1,
		HostsFile:    hostsFile,
		ScenariosDir: scenarios,
		Stdout:       &buf,
	})
	const msg = "this run would use model usage on your account; re-run with --yes to start"
	if err == nil || !strings.Contains(err.Error(), msg) || !strings.Contains(buf.String(), msg) {
		t.Fatalf("err %v\n%s", err, buf.String())
	}
	if _, statErr := os.Stat(marker); !os.IsNotExist(statErr) {
		t.Fatal("host started without --yes")
	}
}

func TestRunRejectsUnknownHost(t *testing.T) {
	dir := t.TempDir()
	hostsFile := writeHostsFile(t, dir, []hostRow{{
		name: "claude", command: "claude", model: "claude-opus-5-5", version: "0", timeout: 5,
	}})
	scenarios := writeScenario(t, dir, "marker")
	var buf bytes.Buffer
	err := Run(context.Background(), RunOptions{
		Runs:         1,
		HostNames:    []string{"nope"},
		HostsFile:    hostsFile,
		ScenariosDir: scenarios,
		DryRun:       true,
		Stdout:       &buf,
	})
	if err == nil || !strings.Contains(err.Error(), "nope") || !strings.Contains(buf.String(), "not in hosts.yaml") {
		t.Fatalf("%v\n%s", err, buf.String())
	}
}

func TestRunsMustBePositive(t *testing.T) {
	err := Run(context.Background(), RunOptions{Runs: 0})
	if err == nil || !strings.Contains(err.Error(), "at least 1") {
		t.Fatal(err)
	}
}

func TestSnapshotSkipsNonGit(t *testing.T) {
	runRoot := t.TempDir()
	w, err := emptyWorld(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	rt := &runtime{opts: RunOptions{Caller: map[string]string{"PATH": "/usr/bin:/bin"}}}
	snap, err := rt.snapshot(context.Background(), w)
	if err != nil {
		t.Fatal(err)
	}
	if snap.Refs == nil || snap.Worktrees == nil || snap.RunDigest == nil {
		t.Fatalf("nil maps %+v", snap)
	}
	if len(snap.Heads) != 0 || len(snap.Status) != 0 || len(snap.Refs) != 0 {
		t.Fatalf("snapshot %+v", snap)
	}
}

func TestWriteBytesDropsSecrets(t *testing.T) {
	dir := t.TempDir()
	var drops []string
	secret := "sk-" + strings.Repeat("b", 40)
	if err := writeBytes(filepath.Join(dir, "stderr.txt"), []byte("noise "+secret+"\n"), textRedactor{}, &drops); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "stderr.txt")); !os.IsNotExist(err) {
		t.Fatal("secret file was written")
	}
	if err := writeBytes(filepath.Join(dir, "auth.json"), []byte("{}\n"), textRedactor{}, &drops); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "auth.json")); !os.IsNotExist(err) {
		t.Fatal("auth file was written")
	}
	if len(drops) != 2 {
		t.Fatalf("drops %v", drops)
	}
	note := filepath.Join(dir, "note.txt")
	if err := writeBytes(note, []byte("see /home/example/notes\n"), textRedactor{home: "/home/example"}, nil); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(note)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "/home/example") || !strings.Contains(string(body), "$HOME/notes") {
		t.Fatalf("redaction %q", body)
	}
}

func TestEmptyWorkspaceStub(t *testing.T) {
	base := tempOutsideHome(t)
	stubs := filepath.Join(base, "stubs")
	if err := os.MkdirAll(stubs, 0o755); err != nil {
		t.Fatal(err)
	}
	parent := filepath.Join(base, "parent")
	out := filepath.Join(base, "out")
	home := filepath.Join(base, "caller-home")
	for _, dir := range []string{parent, out, home} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	secret := "sk-" + strings.Repeat("c", 40)
	root := repoRoot(t)
	stub := writeStreamStub(t, stubs, "claude", filepath.Join(root, "agent-tests", "testdata", "streams", "claude-result.jsonl"), secret)
	rt := &runtime{
		opts: RunOptions{
			Checkout:   root,
			OutDir:     out,
			ParentTemp: parent,
			Caller:     map[string]string{"PATH": "/usr/bin:/bin", "HOME": home},
		},
		id:  "empty",
		pin: "0",
	}
	err := rt.one(context.Background(), planned{
		Scenario: scenario.Scenario{ID: "overhead-empty", Prompt: "Reply with exactly: OK"},
		Host:     "claude",
		Config: hosts.HostConfig{
			Name: "claude", Command: stub, Model: "claude-opus-5-5",
			HostVersion: "stub", TimeoutSeconds: 30,
		},
		Side: "branch",
		Ref:  "HEAD",
		RunN: 1,
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	left, err := os.ReadDir(parent)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range left {
		if strings.HasPrefix(entry.Name(), "agent-test-") {
			t.Fatalf("run root left behind: %s", entry.Name())
		}
	}
	body := treeText(t, out)
	if strings.Contains(body, secret) || strings.Contains(body, home) {
		t.Fatal("empty run kept a secret or the caller home")
	}
	line := findRun(t, out, "claude")
	if line.Scenario != "overhead-empty" || line.Verdict != "pass" || line.Tokens.InputTotal == nil || *line.Tokens.InputTotal != 2+31194+11304 {
		t.Fatalf("%+v", line)
	}
}

func TestProbePrintsOverheadWithoutStarting(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "prompt-started")
	stub := writeVersionStub(t, dir, marker)
	hostsFile := writeHostsFile(t, dir, []hostRow{{
		name: "claude", command: stub, model: "claude-opus-5-5", version: "2.0.0", timeout: 5,
		auth: []string{"ANTHROPIC_API_KEY"},
	}})
	const planted = "planted-anthropic-value"
	var buf bytes.Buffer
	err := Probe(context.Background(), ProbeOptions{
		HostNames: []string{"claude"},
		HostsFile: hostsFile,
		Caller:    map[string]string{"PATH": "/usr/bin:/bin", "ANTHROPIC_API_KEY": planted},
		Stdout:    &buf,
	})
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	for _, part := range []string{
		"host: claude",
		"configured version: 2.0.0",
		"command version: stub-9",
		"model: claude-opus-5-5 pinned",
		"auth env: ANTHROPIC_API_KEY",
		"auth: environment variables",
		"auth: ready",
		"measured with a developer's own host config loaded; the isolated figure is lower",
		`would run the overhead prompt "Reply with exactly: OK"`,
		"overhead-empty",
		"overhead-world",
		"42500 input tokens",
	} {
		if !strings.Contains(got, part) {
			t.Fatalf("missing %q\n%s", part, got)
		}
	}
	if strings.Contains(got, planted) {
		t.Fatal("probe printed an auth value")
	}
	if _, statErr := os.Stat(marker); !os.IsNotExist(statErr) {
		t.Fatal("probe started the prompt")
	}
}

type hostRow struct {
	name    string
	command string
	model   string
	version string
	timeout int
	effort  string
	price   string
	auth    []string
}

func writeHostsFile(t *testing.T, dir string, rows []hostRow) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("hosts:\n")
	for _, row := range rows {
		b.WriteString("  " + row.name + ":\n")
		b.WriteString("    command: " + yamlString(row.command) + "\n")
		b.WriteString("    model: " + yamlString(row.model) + "\n")
		b.WriteString("    host_version: " + yamlString(row.version) + "\n")
		b.WriteString("    timeout_seconds: " + strconv.Itoa(row.timeout) + "\n")
		if row.effort != "" {
			b.WriteString("    effort: " + yamlString(row.effort) + "\n")
		}
		if row.price != "" {
			b.WriteString("    price_row: " + yamlString(row.price) + "\n")
		}
		if len(row.auth) > 0 {
			b.WriteString("    auth_env:\n")
			for _, name := range row.auth {
				b.WriteString("      - " + name + "\n")
			}
		}
	}
	path := filepath.Join(dir, "hosts.yaml")
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeScenario(t *testing.T, dir, id string) string {
	t.Helper()
	scenarios := filepath.Join(dir, "scenarios")
	if err := os.MkdirAll(scenarios, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "id: " + id + "\nversion: 1\ntitle: Marker\nprompt: \"Work in $WS and reply OK\"\nsetup:\n  - \"echo marker > marker.txt\"\n"
	if err := os.WriteFile(filepath.Join(scenarios, id+".yaml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return scenarios
}

func writeTouchStub(t *testing.T, dir, marker string) string {
	t.Helper()
	path := filepath.Join(dir, "stub")
	script := "#!/bin/sh\necho started > " + shSingle(marker) + "\nexit 0\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeVersionStub(t *testing.T, dir, marker string) string {
	t.Helper()
	path := filepath.Join(dir, "stub-version")
	script := "#!/bin/sh\nif [ \"${1:-}\" = \"--version\" ]; then echo stub-9; exit 0; fi\necho started > " + shSingle(marker) + "\nexit 0\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func yamlString(s string) string {
	return "\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\""
}

func shSingle(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "cmd", "context-circuit", "main.go")); err != nil {
		t.Fatal(err)
	}
	return root
}

func porcelain(t *testing.T, root string) string {
	t.Helper()
	res, err := proc.Run(context.Background(), proc.Spec{
		Argv: []string{"git", "-C", root, "status", "--porcelain"},
		Env: []string{
			"PATH=" + os.Getenv("PATH"),
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_CONFIG_GLOBAL=/dev/null",
		},
	})
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("git status: %v\n%s", err, res.Stderr)
	}
	return string(res.Stdout)
}
