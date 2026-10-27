package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

func TestRunStubHosts(t *testing.T) {
	home := os.Getenv("HOME")
	dumps, outText := runAgainstStubs(t, home, "")
	for host, dump := range dumps {
		if home != "" && strings.Contains(dump, home) {
			t.Fatalf("%s env dump contains the test process home", host)
		}
	}
	if home != "" && strings.Contains(outText, home) {
		t.Fatal("output contains the test process home")
	}
}

func TestRunIsolationCanary(t *testing.T) {
	base := tempOutsideHome(t)
	fake := filepath.Join(base, "fakehome")
	const canary = "ISOLATION-CANARY-9c1d"
	for _, rel := range []string{
		"CLAUDE.md",
		"AGENTS.md",
		".codex/config.toml",
		".cursor/cli-config.json",
		".claude/CLAUDE.md",
	} {
		path := filepath.Join(fake, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(canary+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	dumps, outText := runAgainstStubs(t, fake, base)
	for host, dump := range dumps {
		if strings.Contains(dump, canary) {
			t.Fatalf("%s capture contains the isolation canary", host)
		}
		if strings.Contains(dump, fake) {
			t.Fatalf("%s capture contains the caller home", host)
		}
	}
	if strings.Contains(outText, canary) || strings.Contains(outText, fake) {
		t.Fatal("output contains the isolation canary or the caller home")
	}
}

type stubHost struct {
	name    string
	model   string
	input   int64
	output  int64
	basis   string
	costNil bool
	version string
	effort  string
	price   string
	auth    []string
	stream  string
	reads   string
}

func runAgainstStubs(t *testing.T, home, rootDir string) (map[string]string, string) {
	t.Helper()
	if testing.Short() {
		t.Skip("short")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go missing")
	}
	root := repoRoot(t)
	before := porcelain(t, root)
	t.Cleanup(func() {
		after := porcelain(t, root)
		if before != after {
			t.Errorf("checkout status changed\nbefore:\n%s\nafter:\n%s", before, after)
		}
	})

	base := rootDir
	if base == "" {
		base = tempOutsideHome(t)
	}
	bin := filepath.Join(base, "bin")
	parent := filepath.Join(base, "parent")
	out := filepath.Join(base, "out")
	stubs := filepath.Join(base, "stubs")
	for _, dir := range []string{parent, out, stubs} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	path := toolPath(t, bin)
	secret := "sk-" + strings.Repeat("a", 40)
	streams := filepath.Join(root, "agent-tests", "testdata", "streams")
	hosts := []stubHost{
		{name: "claude", model: "claude-opus-5-5", input: 2 + 31194 + 11304, output: 4, basis: "reported", version: "stub-claude", auth: []string{"ANTHROPIC_API_KEY", "CLAUDE_CODE_OAUTH_TOKEN"}, stream: "claude-result.jsonl"},
		{name: "codex", model: "gpt-5.6-luna", input: 15974, output: 5, basis: "price-table", version: "stub-codex", effort: "medium", price: "gpt-5.6-luna", auth: []string{"OPENAI_API_KEY"}, stream: "codex-exec.jsonl"},
		{name: "cursor", model: "Claude Opus 4.8 300K High", input: 2 + 21364, output: 73, basis: "price-table", version: "stub-cursor", price: "composer-2.5", auth: []string{"CURSOR_API_KEY"}, stream: "cursor-stream.jsonl", reads: "reads list not available"},
	}
	var rows []hostRow
	for _, host := range hosts {
		rows = append(rows, hostRow{
			name: host.name, command: writeStreamStub(t, stubs, host.name, filepath.Join(streams, host.stream), secret),
			model: configModel(host), version: host.version, timeout: 120,
			effort: host.effort, price: host.price, auth: host.auth,
		})
	}
	work := filepath.Join(base, "work")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatal(err)
	}
	hostsFile := writeHostsFile(t, work, rows)
	scenarios := writeScenario(t, work, "marker")
	fixtures := filepath.Join(root, "agent-tests", "fixtures")
	wantDigest, err := record.Digest(fixtures)
	if err != nil {
		t.Fatal(err)
	}
	pin, err := worldPin(root)
	if err != nil {
		t.Fatal(err)
	}
	src := revParseTest(t, root, "HEAD")
	caller := map[string]string{
		"PATH":                  path,
		"HOME":                  home,
		"GITHUB_TOKEN":          "PLANTED-GITHUB-TOKEN-VALUE",
		"SSH_AUTH_SOCK":         "/planted/ssh-agent.sock",
		"AWS_SECRET_ACCESS_KEY": "PLANTED-AWS-SECRET-VALUE",
		// run checks auth readiness before the first host starts.
		"ANTHROPIC_API_KEY": "stub-anthropic-auth-value",
		"OPENAI_API_KEY":    "stub-openai-auth-value",
		"CURSOR_API_KEY":    "stub-cursor-auth-value",
	}
	dumps := map[string]string{}
	for _, host := range hosts {
		var buf bytes.Buffer
		err := Run(context.Background(), RunOptions{
			Checkout:     root,
			Ref:          "HEAD",
			ScenarioIDs:  []string{"marker"},
			HostNames:    []string{host.name},
			Runs:         1,
			HostsFile:    hostsFile,
			PricesFile:   filepath.Join(root, "agent-tests", "prices.yaml"),
			ScenariosDir: scenarios,
			FixturesDir:  fixtures,
			OutDir:       out,
			Caller:       caller,
			ParentTemp:   parent,
			Confirm:      true,
			Stdout:       &buf,
		})
		if err != nil {
			t.Fatalf("%s: %v\n%s", host.name, err, buf.String())
		}
		dumpPath := filepath.Join(stubs, host.name+".capture")
		raw, err := os.ReadFile(dumpPath)
		if err != nil {
			t.Fatal(err)
		}
		dump := string(raw)
		dumps[host.name] = dump
		pwd, envText, argv := splitDump(t, dump)
		ws := envValue(envText, "WS")
		if !strings.HasSuffix(pwd, "/ws") || !samePath(ws, pwd) {
			t.Fatalf("%s WS %s pwd %s", host.name, ws, pwd)
		}
		if strings.Contains(argv, "$WS") || (!strings.Contains(argv, ws) && !strings.Contains(argv, pwd)) {
			t.Fatalf("%s argv does not contain the workspace", host.name)
		}
		if pathExists(ws) || pathExists(pwd) {
			t.Fatalf("%s workspace still exists", host.name)
		}
		for _, banned := range []string{"GITHUB_TOKEN", "SSH_AUTH_SOCK", "AWS_SECRET_ACCESS_KEY", "PLANTED-GITHUB-TOKEN-VALUE", "PLANTED-AWS-SECRET-VALUE", "/planted/ssh-agent.sock"} {
			if strings.Contains(dump, banned) {
				t.Fatalf("%s capture contains %s", host.name, banned)
			}
		}
		assertHostOutput(t, out, host, secret, wantDigest, pin, src)
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
	return dumps, treeText(t, out)
}

func configModel(host stubHost) string {
	if host.name == "cursor" {
		return "auto"
	}
	return host.model
}

func assertHostOutput(t *testing.T, out string, host stubHost, secret, digest, pin, src string) {
	t.Helper()
	line := findRun(t, out, host.name)
	if line.Scenario != "marker" || line.Side != "branch" || line.Run != 1 || line.Verdict != "pass" {
		t.Fatalf("%s record %+v", host.name, line)
	}
	// The stub prints stub-9 for --version; the configured value is not recorded.
	if line.Model != host.model || line.HostVersion != "stub-9" || line.HostVersionSource != "command" || line.UsageSource != "reported" {
		t.Fatalf("%s model %s version %s (%s) usage %s", host.name, line.Model, line.HostVersion, line.HostVersionSource, line.UsageSource)
	}
	if line.Tokens.InputTotal == nil || *line.Tokens.InputTotal != host.input {
		t.Fatalf("%s input %v", host.name, line.Tokens.InputTotal)
	}
	if line.Tokens.Output == nil || *line.Tokens.Output != host.output {
		t.Fatalf("%s output %v", host.name, line.Tokens.Output)
	}
	if line.CostBasis != host.basis || (line.CostUSD == nil) != host.costNil {
		t.Fatalf("%s cost %v basis %s", host.name, line.CostUSD, line.CostBasis)
	}
	if line.FixtureDigest != digest || line.CLIVersion != pin || line.SrcCommit != src {
		t.Fatalf("%s digest/pin/commit %s %s %s", host.name, line.FixtureDigest, line.CLIVersion, line.SrcCommit)
	}
	if line.ScenarioVersion != "1" {
		t.Fatalf("%s version %s", host.name, line.ScenarioVersion)
	}
	marker := globOne(t, out, "marker", host.name, filepath.Join("artifacts", "files", "marker.txt"))
	body, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "marker") {
		t.Fatalf("%s marker %q", host.name, body)
	}
	text := treeText(t, filepath.Dir(filepath.Dir(filepath.Dir(marker))))
	if strings.Contains(text, secret) {
		t.Fatalf("%s output contains the planted key", host.name)
	}
	drop := globOne(t, out, "marker", host.name, "drop.log")
	dropBody, err := os.ReadFile(drop)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(dropBody), "stderr.txt") {
		t.Fatalf("%s drop log %q", host.name, dropBody)
	}
	if host.reads != "" {
		reads := globOne(t, out, "marker", host.name, "reads.txt")
		got, err := os.ReadFile(reads)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(got), host.reads) {
			t.Fatalf("%s reads %q", host.name, got)
		}
	}
	decoy := globOne(t, out, "marker", host.name, "decoy.log")
	got, err := os.ReadFile(decoy)
	if err != nil {
		t.Fatal(err)
	}
	if len(bytes.TrimSpace(got)) != 0 {
		t.Fatalf("%s decoy %q", host.name, got)
	}
}

type runLine struct {
	Scenario          string   `json:"scenario"`
	ScenarioVersion   string   `json:"scenario_version"`
	FixtureDigest     string   `json:"fixture_digest"`
	SrcCommit         string   `json:"src_commit"`
	Side              string   `json:"side"`
	CLIVersion        string   `json:"cli_version"`
	Host              string   `json:"host"`
	HostVersion       string   `json:"host_version"`
	HostVersionSource string   `json:"host_version_source"`
	Model             string   `json:"model"`
	Run               int      `json:"run"`
	Verdict           string   `json:"verdict"`
	UsageSource       string   `json:"usage_source"`
	CostUSD           *float64 `json:"cost_usd"`
	CostBasis         string   `json:"cost_basis"`
	Tokens            struct {
		InputTotal *int64 `json:"input_total"`
		Output     *int64 `json:"output"`
	} `json:"tokens"`
}

func findRun(t *testing.T, out, host string) runLine {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(out, "*", "runs.jsonl"))
	if err != nil || len(matches) == 0 {
		t.Fatalf("runs.jsonl: %v %v", matches, err)
	}
	var found []runLine
	for _, match := range matches {
		f, err := os.Open(match)
		if err != nil {
			t.Fatal(err)
		}
		dec := json.NewDecoder(f)
		for dec.More() {
			var line runLine
			if err := dec.Decode(&line); err != nil {
				f.Close()
				t.Fatal(err)
			}
			if line.Host == host {
				found = append(found, line)
			}
		}
		f.Close()
	}
	if len(found) != 1 {
		t.Fatalf("%s lines %d", host, len(found))
	}
	return found[0]
}

func globOne(t *testing.T, out, scenario, host, rel string) string {
	t.Helper()
	pattern := filepath.Join(out, "*", "runs", scenario, host, "branch", "r1", rel)
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) != 1 {
		t.Fatalf("glob %s: %v %v", pattern, matches, err)
	}
	return matches[0]
}

func treeText(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !d.Type().IsRegular() {
			return err
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		b.Write(body)
		b.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func samePath(a, b string) bool {
	return resolvePath(a) == resolvePath(b)
}

// resolvePath canonicalises a path whose leaf may already have been removed.
// Ancestors that still exist are resolved, so a symlink such as /tmp matches.
func resolvePath(path string) string {
	path = filepath.Clean(path)
	var tail []string
	for {
		if _, err := os.Lstat(path); err == nil {
			resolved, err := filepath.EvalSymlinks(path)
			if err != nil {
				resolved = path
			}
			return filepath.Clean(filepath.Join(append([]string{resolved}, tail...)...))
		}
		parent := filepath.Dir(path)
		if parent == path {
			return filepath.Clean(filepath.Join(append([]string{path}, tail...)...))
		}
		tail = append([]string{filepath.Base(path)}, tail...)
		path = parent
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func splitDump(t *testing.T, dump string) (pwd, envText, argv string) {
	t.Helper()
	line, _, ok := strings.Cut(dump, "\n")
	if !ok || !strings.HasPrefix(line, "PWD=") {
		t.Fatalf("dump missing pwd")
	}
	pwd = strings.TrimPrefix(line, "PWD=")
	envText = between(t, dump, "ENV_START\n", "\nENV_END")
	argv = between(t, dump, "ARGV_START\n", "\nARGV_END")
	return pwd, envText, argv
}

func between(t *testing.T, s, start, end string) string {
	t.Helper()
	_, rest, ok := strings.Cut(s, start)
	if !ok {
		t.Fatalf("missing %s", start)
	}
	got, _, ok := strings.Cut(rest, end)
	if !ok {
		t.Fatalf("missing %s", end)
	}
	return got
}

func envValue(envText, key string) string {
	for _, line := range strings.Split(envText, "\n") {
		name, val, ok := strings.Cut(line, "=")
		if ok && name == key {
			return val
		}
	}
	return ""
}

func writeStreamStub(t *testing.T, dir, name, stream, secret string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	var b strings.Builder
	b.WriteString("#!/bin/sh\n")
	b.WriteString("dump=$(dirname \"$0\")/" + name + ".capture\n")
	b.WriteString("{\n")
	b.WriteString("echo \"PWD=$(pwd)\"\n")
	b.WriteString("echo ARGV_START\n")
	b.WriteString("printf '%s\\n' \"$0\" \"$@\"\n")
	b.WriteString("echo ARGV_END\n")
	b.WriteString("echo ENV_START\n")
	b.WriteString("env\n")
	b.WriteString("echo ENV_END\n")
	b.WriteString("echo HOMEFILES_START\n")
	b.WriteString("if [ -n \"${HOME:-}\" ]; then\n")
	b.WriteString("for f in \"$HOME/CLAUDE.md\" \"$HOME/AGENTS.md\" \"$HOME/.codex/config.toml\" \"$HOME/.cursor/cli-config.json\" \"$HOME/.claude/CLAUDE.md\"; do\n")
	b.WriteString("if [ -f \"$f\" ]; then echo \"FILE $f\"; cat \"$f\"; fi\n")
	b.WriteString("done\n")
	b.WriteString("fi\n")
	b.WriteString("echo HOMEFILES_END\n")
	b.WriteString("} > \"$dump\"\n")
	b.WriteString("if [ \"${1:-}\" = \"--version\" ]; then echo stub-9; exit 0; fi\n")
	b.WriteString("if [ \"${1:-}\" = \"login\" ]; then cat >/dev/null; exit 0; fi\n")
	b.WriteString("printf '%s\\n' " + shSingle(secret) + " >&2\n")
	b.WriteString("cat " + shSingle(stream) + "\n")
	if err := os.WriteFile(path, []byte(b.String()), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func toolPath(t *testing.T, bin string) string {
	t.Helper()
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"sh", "git"} {
		found, err := exec.LookPath(name)
		if err != nil {
			t.Skip(name + " missing")
		}
		if err := os.Symlink(found, filepath.Join(bin, name)); err != nil {
			t.Fatal(err)
		}
	}
	return bin + ":/usr/bin:/bin"
}

func tempOutsideHome(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "agent-tests-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	home := os.Getenv("HOME")
	if home != "" && strings.Contains(dir, home) {
		t.Fatal("temp dir is inside the test home")
	}
	return dir
}

func revParseTest(t *testing.T, root, ref string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	got, err := revParse(ctx, root, ref, os.Getenv("PATH"))
	if err != nil {
		t.Fatal(err)
	}
	return got
}
