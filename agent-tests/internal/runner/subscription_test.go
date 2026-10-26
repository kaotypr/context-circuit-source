package runner

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
)

const (
	loginMarker  = "MARKER-subscription-login-9f3c-zzzz"
	loginJWT     = "eyJhbGciOiJub25lIn0.eyJzdWIiOiJ0ZXN0In0.signatureXXXX"
	loginKeyName = "planted-openai-key-value"
)

func loginRefresh() string { return "rt_" + strings.Repeat("q", 24) }

func loginBody() []byte {
	return []byte(`{"access_token":"` + loginJWT + `","refresh_token":"` + loginRefresh() + `","note":"` + loginMarker + `"}` + "\n")
}

func TestSubscriptionWithoutFlagDoesNotTouchLogin(t *testing.T) {
	// The source path does not exist. A stat or a copy would fail the run.
	r := runLogin(t, loginCase{
		flag: false, auth: "subscription", source: "missing", behavior: "ok",
	})
	if r.err != nil {
		t.Fatal(r.err)
	}
	if r.saw != "0" {
		t.Fatalf("saw %s", r.saw)
	}
	if strings.Contains(r.repro, "--subscription-auth") {
		t.Fatalf("reproduce %s", r.repro)
	}
	if strings.Contains(r.outText, loginMarker) || strings.Contains(r.outText, loginJWT) {
		t.Fatal("output contains login bytes without the flag")
	}
	if len(r.roots) != 0 || len(r.copies) != 0 {
		t.Fatalf("roots %v copies %v", r.roots, r.copies)
	}
}

func TestSubscriptionCopyAndDelete(t *testing.T) {
	r := runLogin(t, loginCase{flag: true, auth: "env", source: "file", behavior: "ok"})
	if r.err != nil {
		t.Fatal(r.err)
	}
	assertLoginSeen(t, r)
	if !strings.Contains(r.repro, "--subscription-auth codex") {
		t.Fatalf("reproduce %s", r.repro)
	}
	if len(r.roots) != 0 || len(r.copies) != 0 {
		t.Fatalf("roots %v copies %v", r.roots, r.copies)
	}
	assertRedacted(t, r.out)
	if !strings.Contains(r.outText, "[redacted]") {
		t.Fatal("stream and reply were not redacted")
	}
}

func TestSubscriptionHomeFallback(t *testing.T) {
	r := runLogin(t, loginCase{flag: true, auth: "env", source: "home", behavior: "ok"})
	if r.err != nil {
		t.Fatal(r.err)
	}
	assertLoginSeen(t, r)
	if strings.Contains(r.capture, "caller-home-9f3c") {
		t.Fatal("host saw the caller home")
	}
}

func TestSubscriptionKeepDeletesLogin(t *testing.T) {
	r := runLogin(t, loginCase{flag: true, keep: true, source: "file", behavior: "ok"})
	if r.err != nil {
		t.Fatal(r.err)
	}
	assertLoginSeen(t, r)
	if len(r.roots) != 1 || len(r.copies) != 0 {
		t.Fatalf("roots %d copies %v", len(r.roots), r.copies)
	}
	assertRedacted(t, r.out)
	if strings.Contains(treeText(t, r.roots[0]), loginMarker) {
		t.Fatal("kept run root still has the login marker")
	}
}

func TestSubscriptionDeletedAfterExit(t *testing.T) {
	r := runLogin(t, loginCase{flag: true, source: "file", behavior: "exit"})
	if r.err == nil || !strings.Contains(r.err.Error(), "exit") {
		t.Fatal(r.err)
	}
	assertLoginSeen(t, r)
	if len(r.roots) != 1 || len(r.copies) != 0 {
		t.Fatalf("roots %d copies %v", len(r.roots), r.copies)
	}
	if r.verdict != "harness-error" {
		t.Fatalf("verdict %s", r.verdict)
	}
	assertRedacted(t, r.out)
}

func TestSubscriptionDeletedAfterTimeout(t *testing.T) {
	r := runLogin(t, loginCase{flag: true, source: "file", behavior: "sleep", timeout: 1})
	if r.err == nil || !strings.Contains(r.err.Error(), "timed out") {
		t.Fatal(r.err)
	}
	assertLoginSeen(t, r)
	if len(r.roots) != 1 || len(r.copies) != 0 {
		t.Fatalf("roots %d copies %v", len(r.roots), r.copies)
	}
	assertRedacted(t, r.out)
}

func TestSubscriptionDeletedAfterCancel(t *testing.T) {
	r := runLogin(t, loginCase{flag: true, source: "file", behavior: "sleep", cancel: true, timeout: 30})
	if r.err == nil {
		t.Fatal("expected cancel to stop the host")
	}
	assertLoginSeen(t, r)
	if len(r.roots) != 1 || len(r.copies) != 0 {
		t.Fatalf("roots %d copies %v", len(r.roots), r.copies)
	}
	if strings.Contains(r.err.Error(), loginMarker) {
		t.Fatal("error contains the login marker")
	}
}

func TestSubscriptionDeletedAfterPanic(t *testing.T) {
	afterLoginCopy = func() { panic("login copy panic") }
	t.Cleanup(func() { afterLoginCopy = nil })
	r := runLogin(t, loginCase{flag: true, source: "file", behavior: "ok", panic: true})
	if r.err == nil || !strings.Contains(r.err.Error(), "panic") {
		t.Fatal(r.err)
	}
	if len(r.roots) != 1 || len(r.copies) != 0 {
		t.Fatalf("roots %d copies %v", len(r.roots), r.copies)
	}
	if strings.Contains(r.err.Error(), loginMarker) || strings.Contains(r.outText, loginMarker) {
		t.Fatal("panic path leaked the login marker")
	}
}

func TestSubscriptionDeletedAfterStartFailure(t *testing.T) {
	r := runLogin(t, loginCase{flag: true, source: "file", missingCmd: true})
	if r.err == nil {
		t.Fatal("expected the host to fail to start")
	}
	if len(r.roots) != 1 || len(r.copies) != 0 {
		t.Fatalf("roots %d copies %v", len(r.roots), r.copies)
	}
	if strings.Contains(r.err.Error(), loginMarker) {
		t.Fatal("error contains the login marker")
	}
}

func TestSubscriptionDeleteFailureKeepsRoot(t *testing.T) {
	r := runLogin(t, loginCase{flag: true, source: "file", behavior: "chmod"})
	if r.err == nil || !strings.Contains(r.err.Error(), "subscription login file") {
		t.Fatal(r.err)
	}
	if r.verdict != "harness-error" {
		t.Fatalf("verdict %s", r.verdict)
	}
	if len(r.copies) != 1 {
		t.Fatalf("copy count %d", len(r.copies))
	}
	if !strings.Contains(r.outText, "subscription login file could not be removed") {
		t.Fatalf("drop text %s", r.outText)
	}
	assertRedacted(t, r.out)
}

func TestSubscriptionMissingAndSymlink(t *testing.T) {
	missing := runLogin(t, loginCase{flag: true, source: "missing", behavior: "ok"})
	if missing.err == nil || !strings.Contains(missing.err.Error(), "subscription login file") || !strings.Contains(missing.err.Error(), "not found") {
		t.Fatal(missing.err)
	}
	if strings.Contains(missing.err.Error(), "missing-login-src") || strings.Contains(missing.err.Error(), loginMarker) {
		t.Fatalf("err %v", missing.err)
	}
	link := runLogin(t, loginCase{flag: true, source: "symlink", behavior: "ok"})
	if link.err == nil || !strings.Contains(link.err.Error(), "subscription login file") || !strings.Contains(link.err.Error(), "symlink") {
		t.Fatal(link.err)
	}
	if strings.Contains(link.err.Error(), loginMarker) || strings.Contains(link.err.Error(), loginJWT) {
		t.Fatalf("err %v", link.err)
	}
}

func TestSweepSubscriptionFiles(t *testing.T) {
	base := t.TempDir()
	parent := filepath.Join(base, "parent")
	if err := os.MkdirAll(parent, 0o700); err != nil {
		t.Fatal(err)
	}
	staleDir := filepath.Join(parent, "agent-test-stale", "codex")
	if err := os.MkdirAll(staleDir, 0o700); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(staleDir, "auth.json")
	if err := os.WriteFile(stale, []byte("stale-copy\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(staleDir, "keep.txt")
	if err := os.WriteFile(keep, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(parent, "agent-test-stale", "other.txt")
	if err := os.WriteFile(other, []byte("other\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	unrelated := filepath.Join(parent, "notes", "codex")
	if err := os.MkdirAll(unrelated, 0o700); err != nil {
		t.Fatal(err)
	}
	unrelatedFile := filepath.Join(unrelated, "auth.json")
	if err := os.WriteFile(unrelatedFile, []byte("unrelated\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(base, "outside", "codex")
	if err := os.MkdirAll(outside, 0o700); err != nil {
		t.Fatal(err)
	}
	outsideFile := filepath.Join(outside, "auth.json")
	if err := os.WriteFile(outsideFile, []byte("outside-copy\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(base, "outside"), filepath.Join(parent, "agent-test-link")); err != nil {
		t.Fatal(err)
	}
	symDir := filepath.Join(parent, "agent-test-sym", "codex")
	if err := os.MkdirAll(symDir, 0o700); err != nil {
		t.Fatal(err)
	}
	symTarget := filepath.Join(base, "sym-target")
	if err := os.WriteFile(symTarget, []byte("sym-target-bytes\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(symTarget, filepath.Join(symDir, "auth.json")); err != nil {
		t.Fatal(err)
	}
	homeFile := filepath.Join(base, "caller-home-9f3c", ".codex", "auth.json")
	if err := os.MkdirAll(filepath.Dir(homeFile), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(homeFile, []byte("home-copy\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := sweepBeforeHosts(&buf, parent); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "swept 1 subscription login file(s)\n" {
		t.Fatalf("report %q", buf.String())
	}
	if _, err := os.Lstat(stale); !os.IsNotExist(err) {
		t.Fatal("stale copy remains")
	}
	for _, path := range []string{keep, other, unrelatedFile, outsideFile, symTarget, homeFile, filepath.Join(symDir, "auth.json")} {
		if _, err := os.Lstat(path); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
	}
}

func TestProbeSubscriptionMode(t *testing.T) {
	dir := t.TempDir()
	started := filepath.Join(dir, "prompt-started")
	stub := writeVersionStub(t, dir, started)
	hostsFile := writeHostsFile(t, dir, []hostRow{{
		name: "codex", command: stub, model: "gpt-6-sol", version: "stub", timeout: 5,
		auth: []string{"OPENAI_API_KEY"},
	}})
	codexHome := filepath.Join(dir, "codex-home")
	if err := os.MkdirAll(codexHome, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(codexHome, "auth.json"), loginBody(), 0o600); err != nil {
		t.Fatal(err)
	}
	caller := map[string]string{
		"PATH":       "/usr/bin:/bin",
		"HOME":       filepath.Join(dir, "caller-home-9f3c"),
		"CODEX_HOME": codexHome,
	}
	var present bytes.Buffer
	if err := Probe(context.Background(), ProbeOptions{
		HostNames: []string{"codex"}, HostsFile: hostsFile, Caller: caller,
		SubscriptionAuth: []string{"codex"}, Stdout: &present,
	}); err != nil {
		t.Fatal(err)
	}
	got := present.String()
	if !strings.Contains(got, "auth: subscription login file (present)") {
		t.Fatalf("%s", got)
	}
	if strings.Contains(got, loginMarker) || strings.Contains(got, loginJWT) || strings.Contains(got, loginRefresh()) {
		t.Fatal("probe printed login contents")
	}
	if _, err := os.Stat(started); !os.IsNotExist(err) {
		t.Fatal("probe started a prompt")
	}
	var missing bytes.Buffer
	caller["CODEX_HOME"] = filepath.Join(dir, "missing-login-src")
	if err := Probe(context.Background(), ProbeOptions{
		HostNames: []string{"codex"}, HostsFile: hostsFile, Caller: caller,
		SubscriptionAuth: []string{"codex"}, Stdout: &missing,
	}); err != nil {
		t.Fatal(err)
	}
	text := missing.String()
	if !strings.Contains(text, "auth: subscription login file (missing)") || !strings.Contains(text, "no Codex login file found") {
		t.Fatalf("%s", text)
	}
	if strings.Contains(text, "missing-login-src") || strings.Contains(text, loginMarker) {
		t.Fatal("probe printed the login path or contents")
	}
	var envMode bytes.Buffer
	if err := Probe(context.Background(), ProbeOptions{
		HostNames: []string{"codex"}, HostsFile: hostsFile, Caller: caller, Stdout: &envMode,
	}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(envMode.String(), "auth: environment variables") || strings.Contains(envMode.String(), "no Codex login file found") {
		t.Fatalf("%s", envMode.String())
	}
}

func TestValidateSubscriptionAuth(t *testing.T) {
	if err := ValidateSubscriptionAuth(nil); err != nil {
		t.Fatal(err)
	}
	if err := ValidateSubscriptionAuth([]string{"codex", "codex"}); err != nil {
		t.Fatal(err)
	}
	err := ValidateSubscriptionAuth([]string{"codex", "claude"})
	if err == nil || !strings.Contains(err.Error(), "codex") || !strings.Contains(err.Error(), "claude") {
		t.Fatal(err)
	}
}

type loginCase struct {
	flag       bool
	keep       bool
	auth       string
	source     string
	behavior   string
	timeout    int
	cancel     bool
	panic      bool
	missingCmd bool
}

type loginResult struct {
	err     error
	saw     string
	mode    string
	hash    string
	capture string
	repro   string
	out     string
	outText string
	roots   []string
	copies  []string
	verdict string
	source  []byte
}

func runLogin(t *testing.T, spec loginCase) loginResult {
	t.Helper()
	base := t.TempDir()
	t.Cleanup(func() { chmodTree(base) })
	parent := filepath.Join(base, "parent")
	out := filepath.Join(base, "out")
	for _, dir := range []string{parent, out} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	body := loginBody()
	var sourcePath string
	callerHome := filepath.Join(base, "caller-home-9f3c")
	if err := os.MkdirAll(callerHome, 0o700); err != nil {
		t.Fatal(err)
	}
	switch spec.source {
	case "home":
		sourcePath = filepath.Join(callerHome, ".codex", "auth.json")
	case "missing":
		sourcePath = filepath.Join(base, "missing-login-src", "auth.json")
	case "symlink":
		real := filepath.Join(base, "real-login", "auth.json")
		if err := os.MkdirAll(filepath.Dir(real), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(real, body, 0o644); err != nil {
			t.Fatal(err)
		}
		sourcePath = filepath.Join(base, "link-login", "auth.json")
		if err := os.MkdirAll(filepath.Dir(sourcePath), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(real, sourcePath); err != nil {
			t.Fatal(err)
		}
	default:
		sourcePath = filepath.Join(base, "login-source-9f3c", "auth.json")
		spec.source = "file"
	}
	if spec.source == "file" || spec.source == "home" {
		if err := os.MkdirAll(filepath.Dir(sourcePath), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(sourcePath, body, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	before := append([]byte(nil), body...)
	if spec.source == "file" || spec.source == "home" || spec.source == "symlink" {
		got, err := os.ReadFile(sourcePath)
		if spec.source == "symlink" {
			got, err = os.ReadFile(filepath.Join(base, "real-login", "auth.json"))
		}
		if err != nil {
			t.Fatal(err)
		}
		before = got
	}
	capture := filepath.Join(base, "capture.txt")
	ready := filepath.Join(base, "ready")
	stub := filepath.Join(base, "stub.sh")
	if err := os.WriteFile(stub, []byte(loginStub(capture, ready, spec.behavior)), 0o755); err != nil {
		t.Fatal(err)
	}
	command := stub
	if spec.missingCmd {
		command = filepath.Join(base, "missing-stub")
	}
	timeout := spec.timeout
	if timeout == 0 {
		timeout = 30
	}
	caller := map[string]string{
		"PATH":           "/usr/bin:/bin",
		"HOME":           callerHome,
		"OPENAI_API_KEY": loginKeyName,
		"USER":           "zztester",
	}
	if spec.source != "home" {
		caller["CODEX_HOME"] = filepath.Dir(sourcePath)
	}
	ctx := context.Background()
	cancel := func() {}
	if spec.cancel {
		var stop context.CancelFunc
		ctx, stop = context.WithCancel(ctx)
		cancel = stop
		go func() {
			for i := 0; i < 100; i++ {
				if _, err := os.Stat(ready); err == nil {
					stop()
					return
				}
				time.Sleep(20 * time.Millisecond)
			}
			stop()
		}()
	}
	defer cancel()
	auth := spec.auth
	if auth == "" {
		auth = "env"
	}
	var subs []string
	if spec.flag {
		subs = []string{"codex"}
	}
	rt := &runtime{
		opts: RunOptions{
			OutDir:           out,
			ParentTemp:       parent,
			Keep:             spec.keep,
			Caller:           caller,
			SubscriptionAuth: subs,
		},
		id:     "login",
		pin:    "0",
		stdout: &bytes.Buffer{},
	}
	err := rt.one(ctx, planned{
		Scenario: scenario.Scenario{ID: "login", Version: 1, Title: "Login", Prompt: "Reply OK"},
		Host:     "codex",
		Config: hosts.HostConfig{
			Name: "codex", Command: command, Model: "gpt-6-sol",
			HostVersion: "stub", TimeoutSeconds: timeout, Auth: auth,
			AuthEnv: []string{"OPENAI_API_KEY"},
		},
		Side: "branch", Ref: "HEAD", RunN: 1,
	}, true)
	res := loginResult{err: err, out: out, source: before}
	if spec.source == "file" || spec.source == "home" {
		got, readErr := os.ReadFile(sourcePath)
		if readErr != nil {
			t.Fatal(readErr)
		}
		if !bytes.Equal(got, before) {
			t.Fatal("source bytes changed")
		}
	}
	if spec.source == "symlink" {
		got, readErr := os.ReadFile(filepath.Join(base, "real-login", "auth.json"))
		if readErr != nil || !bytes.Equal(got, before) {
			t.Fatal("symlink target changed")
		}
	}
	raw, readErr := os.ReadFile(capture)
	if readErr == nil {
		res.capture = string(raw)
		res.saw = captureValue(res.capture, "SAW")
		res.mode = strings.TrimLeft(captureValue(res.capture, "MODE"), "0")
		res.hash = captureValue(res.capture, "HASH")
	}
	if repro, readErr := os.ReadFile(filepath.Join(out, "login", "runs", "login", "codex", "branch", "r1", "reproduce.txt")); readErr == nil {
		res.repro = string(repro)
	}
	res.outText = treeText(t, out)
	res.roots = agentTestDirs(parent)
	res.copies = loginCopies(parent)
	if line := readVerdict(out); line != "" {
		res.verdict = line
	}
	realHome := os.Getenv("HOME")
	if realHome != "" && !strings.Contains(base, realHome) {
		if strings.Contains(res.capture, realHome) || strings.Contains(res.outText, realHome) {
			t.Fatal("output contains the process home")
		}
	}
	if strings.Contains(res.capture, "OPENAI_API_KEY") || strings.Contains(res.capture, loginKeyName) || strings.Contains(res.outText, loginKeyName) {
		t.Fatal("auth variable or its value leaked")
	}
	return res
}

func assertLoginSeen(t *testing.T, r loginResult) {
	t.Helper()
	if r.saw != "1" {
		t.Fatalf("saw %q\n%s", r.saw, r.capture)
	}
	if r.mode != "600" {
		t.Fatalf("mode %q", r.mode)
	}
	sum := sha256.Sum256(r.source)
	if r.hash != hex.EncodeToString(sum[:]) {
		t.Fatalf("hash %s", r.hash)
	}
	if strings.Contains(r.capture, "login-source-9f3c") || strings.Contains(r.capture, "caller-home-9f3c") {
		t.Fatal("host env contains the caller login path")
	}
	for _, banned := range []string{"OPENAI_API_KEY=", "HOME=" + "caller"} {
		if strings.Contains(r.capture, banned) {
			t.Fatalf("env leak %s", banned)
		}
	}
}

func assertRedacted(t *testing.T, out string) {
	t.Helper()
	text := treeText(t, out)
	for _, banned := range []string{loginMarker, loginJWT, loginRefresh(), string(loginBody())} {
		if strings.Contains(text, banned) {
			t.Fatalf("output contains login bytes %q", banned)
		}
	}
	stream := filepath.Join(out, "login", "runs", "login", "codex", "branch", "r1", "stream.jsonl")
	reply := filepath.Join(out, "login", "runs", "login", "codex", "branch", "r1", "reply.md")
	for _, path := range []string{stream, reply} {
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if len(bytes.TrimSpace(body)) == 0 {
			continue
		}
		if !strings.Contains(string(body), "[redacted]") {
			t.Fatalf("%s missing redaction", path)
		}
	}
}

func loginStub(capture, ready, behavior string) string {
	var b strings.Builder
	b.WriteString("#!/bin/sh\n")
	b.WriteString("auth=\"${CODEX_HOME:-}/auth.json\"\n")
	b.WriteString("{\n")
	b.WriteString("echo ENV_START\n")
	b.WriteString("env\n")
	b.WriteString("echo ENV_END\n")
	b.WriteString("if [ -f \"$auth\" ]; then\n")
	b.WriteString("echo SAW=1\n")
	b.WriteString("if stat -f %OLp \"$auth\" >/dev/null 2>&1; then echo \"MODE=$(stat -f %OLp \"$auth\")\"; else echo \"MODE=$(stat -c %a \"$auth\")\"; fi\n")
	b.WriteString("if command -v shasum >/dev/null 2>&1; then echo \"HASH=$(shasum -a 256 \"$auth\" | awk '{print $1}')\"; elif command -v sha256sum >/dev/null 2>&1; then echo \"HASH=$(sha256sum \"$auth\" | awk '{print $1}')\"; else echo HASH=missing; fi\n")
	b.WriteString("else echo SAW=0; fi\n")
	b.WriteString("} > " + shSingle(capture) + "\n")
	b.WriteString("if [ \"${1:-}\" = \"login\" ]; then cat >/dev/null; exit 0; fi\n")
	reply := loginMarker + " " + loginJWT + " " + loginRefresh()
	b.WriteString("printf '%s\\n' '{\"type\":\"thread.started\",\"thread_id\":\"t-sub\"}'\n")
	b.WriteString("if [ -f \"$auth\" ]; then\n")
	b.WriteString("printf '%s\\n' " + shSingle(`{"type":"item.completed","item":{"type":"agent_message","text":"`+reply+`"}}`) + "\n")
	b.WriteString("cat \"$auth\"; echo\n")
	b.WriteString("else\n")
	b.WriteString("printf '%s\\n' '{\"type\":\"item.completed\",\"item\":{\"type\":\"agent_message\",\"text\":\"OK\"}}'\n")
	b.WriteString("fi\n")
	b.WriteString("printf '%s\\n' '{\"type\":\"turn.completed\",\"usage\":{\"input_tokens\":4,\"output_tokens\":3}}'\n")
	switch behavior {
	case "exit":
		b.WriteString("exit 1\n")
	case "sleep":
		b.WriteString("echo ready > " + shSingle(ready) + "\nsleep 60\n")
	case "chmod":
		b.WriteString("if [ -f \"$auth\" ]; then chmod 555 \"$(dirname \"$auth\")\"; fi\nexit 0\n")
	default:
		b.WriteString("exit 0\n")
	}
	return b.String()
}

func captureValue(text, key string) string {
	for _, line := range strings.Split(text, "\n") {
		name, val, ok := strings.Cut(line, "=")
		if ok && name == key {
			return strings.TrimSpace(val)
		}
	}
	return ""
}

func agentTestDirs(parent string) []string {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil
	}
	var out []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "agent-test-") {
			out = append(out, filepath.Join(parent, entry.Name()))
		}
	}
	return out
}

func loginCopies(parent string) []string {
	var out []string
	for _, dir := range agentTestDirs(parent) {
		path := filepath.Join(dir, "codex", "auth.json")
		info, err := os.Lstat(path)
		if err == nil && info.Mode().IsRegular() {
			out = append(out, path)
		}
	}
	return out
}

func readVerdict(out string) string {
	matches, err := filepath.Glob(filepath.Join(out, "*", "runs.jsonl"))
	if err != nil || len(matches) == 0 {
		return ""
	}
	body, err := os.ReadFile(matches[0])
	if err != nil {
		return ""
	}
	if strings.Contains(string(body), `"verdict":"harness-error"`) {
		return "harness-error"
	}
	if strings.Contains(string(body), `"verdict":"pass"`) {
		return "pass"
	}
	return ""
}

func chmodTree(root string) {
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err == nil {
			_ = os.Chmod(path, 0o755)
		}
		return nil
	})
}
