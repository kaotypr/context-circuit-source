package runner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

func TestGradeInconclusiveLeavesFailedChecksEmpty(t *testing.T) {
	sc := scenario.Scenario{Checks: []string{"sources-passive", "launcher-only"}}
	ctx := scenario.CheckContext{
		Parsed: hosts.Parsed{ReadsAvailable: false},
		Stream: "final reply",
		After:  scenario.Snapshot{DecoyLog: "context-circuit-cli status\n"},
	}
	results, incomplete := grade(sc, ctx)
	if !incomplete {
		t.Fatal("expected an inconclusive run")
	}
	if scenario.Verdict(results, nil, incomplete) != "fail" {
		t.Fatalf("verdict %s results %+v", scenario.Verdict(results, nil, incomplete), results)
	}
	line := (&runtime{}).recordLine(planned{Scenario: sc, Config: hosts.HostConfig{}}, hosts.Parsed{}, "fail", results, 0, "", nil, nil)
	if len(line.FailedChecks) != 1 || line.FailedChecks[0] != "launcher-only" {
		t.Fatalf("failed checks %v", line.FailedChecks)
	}

	only := scenario.Scenario{Checks: []string{"sources-passive"}}
	results, incomplete = grade(only, scenario.CheckContext{
		Parsed: hosts.Parsed{ReadsAvailable: false},
		Stream: "final reply",
	})
	if !incomplete || !results[0].Pass || !results[0].Inconclusive {
		t.Fatalf("incomplete %v results %+v", incomplete, results)
	}
	if scenario.Verdict(results, nil, incomplete) != "inconclusive" {
		t.Fatal(scenario.Verdict(results, nil, incomplete))
	}
	line = (&runtime{}).recordLine(planned{Scenario: only, Config: hosts.HostConfig{}}, hosts.Parsed{}, "inconclusive", results, 0, "", nil, nil)
	if len(line.FailedChecks) != 0 {
		t.Fatalf("failed checks %v", line.FailedChecks)
	}
}

func TestRecordLineSoftComputed(t *testing.T) {
	sc := scenario.Scenario{Soft: map[string]string{
		"code-canary-seen": "CODE-CANARY-41b2 appeared in the host stream",
		"subagents":        "subagent count",
		"unknown-metric":   "keep me",
	}}
	graded := &scenario.CheckContext{
		Stream: "saw CODE-CANARY-41b2 in a tool result",
		Parsed: hosts.Parsed{Subagents: 0},
	}
	line := (&runtime{}).recordLine(planned{Scenario: sc}, hosts.Parsed{}, "pass", nil, 0, "", nil, graded)
	if line.Soft["code-canary-seen"] != true {
		t.Fatalf("canary %#v", line.Soft["code-canary-seen"])
	}
	if line.Soft["subagents"] != 0 {
		t.Fatalf("subagents %#v", line.Soft["subagents"])
	}
	if line.Soft["unknown-metric"] != "keep me" {
		t.Fatalf("unknown %#v", line.Soft["unknown-metric"])
	}
	plain := (&runtime{}).recordLine(planned{Scenario: sc}, hosts.Parsed{}, "pass", nil, 0, "", nil, nil)
	if plain.Soft["code-canary-seen"] != sc.Soft["code-canary-seen"] {
		t.Fatalf("nil context stored %#v", plain.Soft["code-canary-seen"])
	}
}

func TestReproduceLine(t *testing.T) {
	branch := reproduceLine(planned{
		Scenario: scenario.Scenario{ID: "orient"},
		Host:     "claude",
		Side:     "branch",
		Ref:      "abc123",
	}, "abc123", "", "2026-10-26", "")
	want := "# 2026-10-26\nagent-tests run --scenarios orient --hosts claude --runs 1 --ref abc123\n"
	if branch != want {
		t.Fatalf("%q", branch)
	}
	base := reproduceLine(planned{
		Scenario: scenario.Scenario{ID: "orient"},
		Host:     "claude",
		Side:     "base",
		Ref:      "base123",
	}, "abc123", "base123", "2026-10-26", "")
	if strings.Contains(base, "/") || !strings.Contains(base, "--ref abc123 --base base123") {
		t.Fatalf("%q", base)
	}
}

func TestApplyWorkspace(t *testing.T) {
	w := world.World{
		RunRoot: "/run",
		WS:      "/run/ws",
		API:     "/run/ws/repositories/api",
		Env:     map[string]string{"WS": "/run/ws"},
	}
	applyWorkspace(&w, scenario.Scenario{})
	if w.WS != "/run/ws" || w.Env["WS"] != "/run/ws" {
		t.Fatalf("default %+v", w)
	}
	applyWorkspace(&w, scenario.Scenario{Workspace: "ws2"})
	if w.WS != "/run/ws2" || w.API != "/run/ws2/repositories/api" || w.Env["WS"] != "/run/ws2" {
		t.Fatalf("retarget %+v", w)
	}
}

func TestSetupSnippetsReceiveFixtures(t *testing.T) {
	runRoot := t.TempDir()
	w, err := emptyWorld(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	fix := filepath.Join(runRoot, "fixtures")
	if err := os.MkdirAll(fix, 0o755); err != nil {
		t.Fatal(err)
	}
	rt := &runtime{opts: RunOptions{
		FixturesDir: fix,
		Caller:      map[string]string{"PATH": os.Getenv("PATH")},
	}}
	err = rt.runSnippets(context.Background(), w, scenario.Scenario{
		Setup: []string{`printf '%s' "$FIXTURES" > seen`},
	})
	if err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(w.WS, "seen"))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != fix {
		t.Fatalf("FIXTURES %q", body)
	}
}

func TestHostEnvOmitsFixtures(t *testing.T) {
	base := tempOutsideHome(t)
	stubs := filepath.Join(base, "stubs")
	parent := filepath.Join(base, "parent")
	out := filepath.Join(base, "out")
	home := filepath.Join(base, "caller-home")
	for _, dir := range []string{stubs, parent, out, home} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	root := repoRoot(t)
	stub := writeStreamStub(t, stubs, "claude", filepath.Join(root, "agent-tests", "testdata", "streams", "claude-result.jsonl"), "not-a-key")
	fixtures := filepath.Join(base, "fixture-pack-token")
	if err := os.MkdirAll(fixtures, 0o755); err != nil {
		t.Fatal(err)
	}
	rt := &runtime{
		opts: RunOptions{
			Checkout:    root,
			Ref:         "abc123",
			OutDir:      out,
			ParentTemp:  parent,
			FixturesDir: fixtures,
			Keep:        true,
			Caller:      map[string]string{"PATH": "/usr/bin:/bin", "HOME": home},
		},
		id:  "case1",
		pin: "0",
	}
	err := rt.one(context.Background(), planned{
		Scenario: scenario.Scenario{ID: "marker", Prompt: "Reply with exactly: OK"},
		Host:     "claude",
		Config: hosts.HostConfig{
			Name: "claude", Command: stub, Model: "claude-opus-5-5",
			HostVersion: "stub", TimeoutSeconds: 30,
		},
		Side: "branch",
		Ref:  "abc123",
		RunN: 1,
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	dump, err := os.ReadFile(filepath.Join(stubs, "claude.capture"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(dump), "FIXTURES") || strings.Contains(string(dump), "fixture-pack-token") {
		t.Fatal("host environment contains FIXTURES")
	}
	body, err := os.ReadFile(filepath.Join(out, "case1", "runs", "marker", "claude", "branch", "r1", "reproduce.txt"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(body), "\n")
	if len(lines) < 3 || !nearUTCDate(strings.TrimPrefix(lines[0], "# "), time.Now()) {
		t.Fatalf("reproduce date %q", body)
	}
	if lines[1] != "agent-tests run --scenarios marker --hosts claude --runs 1 --ref abc123" || lines[2] != "" {
		t.Fatalf("reproduce %q", body)
	}
	if strings.Contains(string(body), base) || strings.Contains(string(body), "Users") {
		t.Fatalf("reproduce names a local path: %s", body)
	}
	matches, err := filepath.Glob(filepath.Join(parent, "agent-test-*", "out", "check-before.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		raw, readErr := os.ReadFile(matches[0])
		t.Fatalf("missing launcher wrote check-before %v %s", readErr, raw)
	}
	published := filepath.Join(out, "case1", "runs", "marker", "claude", "branch", "r1", "check.json")
	if raw, readErr := os.ReadFile(published); readErr == nil && strings.TrimSpace(string(raw)) == "{}" {
		t.Fatal("substituted {}")
	}
}

func nearUTCDate(got string, now time.Time) bool {
	now = now.UTC()
	for _, delta := range []int{-1, 0, 1} {
		if got == now.AddDate(0, 0, delta).Format("2006-01-02") {
			return true
		}
	}
	return false
}

func TestRunDigestSkipsAndTruncates(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, dir := range []string{"out", "gocache", "gopath", "tmp", "bin", "src", "setup-home", "home", "claude", "claude-home", "codex", "codex-home", "cursor-home"} {
		write(dir+"/secret.txt", "secret")
	}
	write("host-stream", "stream")
	write("decoy.log", "decoy")
	write("visible.txt", "hello")
	write("keep/tmp/inside.txt", "inside")
	repo := filepath.Join(root, "repo")
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	write("repo/.git/SECRET", "do-not-hash")
	write("repo/readme.txt", "readme")
	if err := os.Symlink("readme.txt", filepath.Join(repo, "linked")); err != nil {
		t.Skip("symlinks unavailable")
	}
	got, err := runDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte("hello"))
	if got["visible.txt"] != hex.EncodeToString(sum[:]) {
		t.Fatalf("visible %q", got["visible.txt"])
	}
	if got["keep/tmp/inside.txt"] == "" || got["repo/readme.txt"] == "" {
		t.Fatalf("missing kept files %+v", got)
	}
	if got["repo"] != "dir" || got["repo/.git"] != "dir" {
		t.Fatalf("dirs %+v", got)
	}
	if _, ok := got["repo/.git/SECRET"]; ok {
		t.Fatal("hashed a file inside .git")
	}
	if got["repo/linked"] != "link:readme.txt" {
		t.Fatalf("link %q", got["repo/linked"])
	}
	if got["home/secret.txt"] == "" {
		t.Fatal("home was skipped")
	}
	for _, dir := range []string{"out", "gocache", "gopath", "tmp", "bin", "src", "setup-home", "claude", "claude-home", "codex", "codex-home", "cursor-home"} {
		if _, ok := got[dir+"/secret.txt"]; ok {
			t.Fatalf("recorded %s", dir)
		}
	}
	for key := range got {
		top, _, _ := strings.Cut(key, "/")
		if digestSkipDir[top] || key == "host-stream" || key == "decoy.log" {
			t.Fatalf("skipped path recorded: %s", key)
		}
	}

	old := runDigestLimit
	runDigestLimit = 2
	t.Cleanup(func() { runDigestLimit = old })
	many := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		if err := os.WriteFile(filepath.Join(many, name), []byte(name), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	truncated, err := runDigest(many)
	if err != nil {
		t.Fatal(err)
	}
	if len(truncated) != 1 || truncated["~truncated"] != "3" {
		t.Fatalf("truncated %+v", truncated)
	}
}

func TestSnapshotReposAndGitState(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	runRoot := t.TempDir()
	w, err := emptyWorld(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	env := worldGitEnv(t, runRoot)
	api := filepath.Join(w.WS, "repositories", "api")
	handbook := filepath.Join(w.WS, "repositories", "handbook")
	notes := filepath.Join(w.WS, "repositories", "notes")
	for _, dir := range []string{api, handbook, notes} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	w.API = api
	initRepo(t, w.WS, env)
	initRepo(t, api, env)
	initRepo(t, handbook, env)
	first := commitRepo(t, w.WS, env, "feat: first")
	commitRepo(t, api, env, "feat: api")
	runGit(t, api, env, "checkout", "-q", "-b", "feature")
	commitRepo(t, api, env, "feat: side")
	runGit(t, api, env, "checkout", "-q", "main")
	commitRepo(t, handbook, env, "feat: handbook")
	if err := os.WriteFile(filepath.Join(runRoot, "visible.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runRoot, "out", "hidden.txt"), []byte("hidden"), 0o644); err != nil {
		t.Fatal(err)
	}
	rt := &runtime{opts: RunOptions{Caller: map[string]string{"PATH": os.Getenv("PATH")}}}
	snap, err := rt.snapshot(context.Background(), w)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"ws", "api", "handbook"} {
		if !strings.Contains(snap.Refs[name], "refs/heads/main ") {
			t.Fatalf("%s refs %q", name, snap.Refs[name])
		}
		if !strings.Contains(snap.Worktrees[name], "branch refs/heads/main") {
			t.Fatalf("%s worktrees %q", name, snap.Worktrees[name])
		}
		if snap.Heads[name] == "" {
			t.Fatalf("%s head empty", name)
		}
	}
	if _, ok := snap.Heads["notes"]; ok {
		t.Fatal("a directory without .git was snapshotted")
	}
	if !strings.Contains(snap.Refs["api"], "refs/heads/feature ") {
		t.Fatalf("api refs %q", snap.Refs["api"])
	}
	apiHead := strings.TrimSpace(runGit(t, api, env, "rev-parse", "HEAD"))
	if snap.Heads["api"] != apiHead {
		t.Fatalf("api head %s want %s", snap.Heads["api"], apiHead)
	}
	if _, ok := snap.RunDigest["out/hidden.txt"]; ok {
		t.Fatal("digest recorded out/")
	}
	if len(snap.RunDigest["visible.txt"]) != 64 {
		t.Fatalf("visible digest %q", snap.RunDigest["visible.txt"])
	}

	second := commitRepo(t, w.WS, env, "feat: second")
	if second == first {
		t.Fatal("second commit did not move HEAD")
	}
	if err := os.WriteFile(filepath.Join(w.WS, "note.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(runRoot, "artifacts", "git")
	if err := rt.writeGitState(context.Background(), w, snap, scenario.Scenario{}, dest); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(dest, "ws.txt"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, part := range []string{"git status --porcelain", "git worktree list", "git branch -a -vv", "git diff --stat", "feat: second", "note.txt", "$RUN"} {
		if !strings.Contains(text, part) {
			t.Fatalf("missing %q\n%s", part, text)
		}
	}
	if strings.Contains(text, runRoot) {
		t.Fatal("git state kept the run root path")
	}
	if _, err := os.Stat(filepath.Join(dest, "api.txt")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dest, "notes.txt")); err == nil {
		t.Fatal("wrote git state for a directory without .git")
	}
}

func TestGitStateKeepsCommandErrors(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	runRoot := t.TempDir()
	w, err := emptyWorld(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.WS, ".git"), []byte("gitdir: /nowhere\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rt := &runtime{opts: RunOptions{Caller: map[string]string{"PATH": os.Getenv("PATH")}}}
	dest := filepath.Join(runRoot, "artifacts", "git")
	if err := rt.writeGitState(context.Background(), w, scenario.Snapshot{}, scenario.Scenario{}, dest); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(dest, "ws.txt"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, title := range []string{"git status --porcelain", "git worktree list", "git branch -a -vv", "git log --all", "git diff --stat"} {
		if !strings.Contains(text, title) {
			t.Fatalf("missing %s\n%s", title, text)
		}
	}
	if !strings.Contains(text, "exit ") {
		t.Fatalf("expected a command error\n%s", text)
	}
}

func worldGitEnv(t *testing.T, runRoot string) []string {
	t.Helper()
	return []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + filepath.Join(runRoot, "setup-home"),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=" + filepath.Join(runRoot, "gitconfig"),
		"GIT_AUTHOR_DATE=2026-01-15T09:00:00Z",
		"GIT_COMMITTER_DATE=2026-01-15T09:00:00Z",
	}
}

func initRepo(t *testing.T, dir string, env []string) {
	t.Helper()
	runGit(t, dir, env, "init", "-q", "-b", "main")
}

func commitRepo(t *testing.T, dir string, env []string, message string) string {
	t.Helper()
	path := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(path, []byte(message+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, env, "add", "note.txt")
	runGit(t, dir, env, "commit", "-q", "-m", message)
	return strings.TrimSpace(runGit(t, dir, env, "rev-parse", "HEAD"))
}

func TestGradeUnknownCheckFails(t *testing.T) {
	results, incomplete := grade(scenario.Scenario{Checks: []string{"not-a-check"}}, scenario.CheckContext{})
	if incomplete || len(results) != 1 || results[0].Pass || results[0].Name != "not-a-check" {
		t.Fatalf("incomplete %v %+v", incomplete, results)
	}
	if scenario.Verdict(results, nil, incomplete) != "fail" {
		t.Fatal(scenario.Verdict(results, nil, incomplete))
	}
}

func TestUnreadableDecoyFailsLauncherCheck(t *testing.T) {
	runRoot := t.TempDir()
	w, err := emptyWorld(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(w.DecoyLog, 0o755); err != nil {
		t.Fatal(err)
	}
	rt := &runtime{opts: RunOptions{Caller: map[string]string{"PATH": "/usr/bin:/bin"}}}
	snap, err := rt.snapshot(context.Background(), w)
	if err != nil {
		t.Fatal(err)
	}
	if snap.DecoyLog != scenario.DecoyLogUnreadable {
		t.Fatalf("decoy %q", snap.DecoyLog)
	}
	got := scenario.Registry["launcher-only"](scenario.CheckContext{After: snap})
	if got.Pass || !strings.Contains(got.Detail, "unreadable") {
		t.Fatalf("%+v", got)
	}
}

func TestRunCheckMissingLauncher(t *testing.T) {
	runRoot := t.TempDir()
	w, err := emptyWorld(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	rt := &runtime{opts: RunOptions{Caller: map[string]string{"PATH": os.Getenv("PATH")}}}
	out, err := rt.runCheck(context.Background(), w)
	if !errors.Is(err, errNoLauncher) || out != nil {
		t.Fatalf("err %v out %q", err, out)
	}
	got := scenario.Registry["check-clean"](scenario.CheckContext{CheckAfter: out})
	if got.Pass {
		t.Fatal("empty check output passed")
	}
}

func TestPublishCheckFiles(t *testing.T) {
	dir := t.TempDir()
	before := []byte("{\"issues\":[]}\n")
	after := []byte("{\"issues\":[{\"issue\":\"stale\"}]}\n")
	if err := publishCheckFiles(dir, before, after, textRedactor{}, nil); err != nil {
		t.Fatal(err)
	}
	if raw, err := os.ReadFile(filepath.Join(dir, "check-before.json")); err != nil || string(raw) != string(before) {
		t.Fatalf("before %v %s", err, raw)
	}
	if raw, err := os.ReadFile(filepath.Join(dir, "check.json")); err != nil || string(raw) != string(after) {
		t.Fatalf("after %v %s", err, raw)
	}
	empty := t.TempDir()
	if err := publishCheckFiles(empty, nil, nil, textRedactor{}, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(empty, "check.json")); !os.IsNotExist(err) {
		t.Fatal("wrote check.json without a document")
	}
	if _, err := os.Stat(filepath.Join(empty, "check-before.json")); !os.IsNotExist(err) {
		t.Fatal("wrote check-before.json without a document")
	}
}

func TestJoinSnapshotsOriginalCheckout(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	runRoot := t.TempDir()
	w, err := emptyWorld(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	env := worldGitEnv(t, runRoot)
	original := filepath.Join(runRoot, "ws")
	clone := filepath.Join(runRoot, "ws2")
	for _, dir := range []string{
		original, filepath.Join(original, "repositories", "api"),
		clone, filepath.Join(clone, "repositories", "api"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		initRepo(t, dir, env)
		commitRepo(t, dir, env, "feat: seed")
	}
	w.WS = clone
	w.API = filepath.Join(clone, "repositories", "api")
	w.RunRoot = runRoot
	rt := &runtime{opts: RunOptions{Caller: map[string]string{"PATH": os.Getenv("PATH")}}}
	before, err := rt.snapshot(context.Background(), w)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"ws", "api", "ws0", "ws0/api"} {
		if before.Heads[name] == "" {
			t.Fatalf("missing %s in %+v", name, before.Heads)
		}
	}
	commitRepo(t, original, env, "chore: hidden commit")
	after, err := rt.snapshot(context.Background(), w)
	if err != nil {
		t.Fatal(err)
	}
	got := scenario.Registry["no-stray-commits"](scenario.CheckContext{
		World: w, Before: before, After: after,
	})
	if got.Pass || !strings.Contains(got.Detail, "ws0 HEAD changed") {
		t.Fatalf("%+v", got)
	}
	if err := os.WriteFile(filepath.Join(original, "stray.txt"), []byte("stray\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dirtied, err := rt.snapshot(context.Background(), w)
	if err != nil {
		t.Fatal(err)
	}
	outside := scenario.Registry["no-outside-writes"](scenario.CheckContext{
		Scenario: scenario.Scenario{Workspace: "ws2"},
		Before:   before,
		After:    dirtied,
	})
	if outside.Pass || !strings.Contains(outside.Detail, "ws/") {
		t.Fatalf("%+v", outside)
	}
}

func TestExecuteGitEvidence(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	runRoot := t.TempDir()
	w, err := emptyWorld(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	env := worldGitEnv(t, runRoot)
	api := filepath.Join(w.WS, "repositories", "api")
	if err := os.MkdirAll(filepath.Join(api, "internal", "invoice"), 0o755); err != nil {
		t.Fatal(err)
	}
	initRepo(t, api, env)
	if err := os.WriteFile(filepath.Join(api, "internal", "invoice", "export.go"), []byte("package invoice\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	commitRepo(t, api, env, "feat: seed")
	runGit(t, api, env, "switch", "-qc", "cc/p1000/api")
	if err := os.WriteFile(filepath.Join(api, "internal", "invoice", "export.go"), []byte("package invoice\n\nfunc CSV() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, api, env, "add", "-A")
	runGit(t, api, env, "commit", "-qm", "feat(invoice): export invoices as CSV", "-m", "Rows are number, customer, status, and cents.")
	w.API = api
	rt := &runtime{opts: RunOptions{Caller: map[string]string{"PATH": os.Getenv("PATH")}}}
	dest := filepath.Join(runRoot, "artifacts", "git")
	if err := rt.writeGitState(context.Background(), w, scenario.Snapshot{}, scenario.Scenario{ID: "execute"}, dest); err != nil {
		t.Fatal(err)
	}
	diffText := readTestFile(t, filepath.Join(dest, "api-diff-stat.txt"))
	logText := readTestFile(t, filepath.Join(dest, "api-log.txt"))
	if !strings.Contains(diffText, "git diff --stat main...cc/p1000/api") || !strings.Contains(diffText, "export.go") {
		t.Fatalf("diff\n%s", diffText)
	}
	if !strings.Contains(logText, "feat(invoice): export invoices as CSV") || !strings.Contains(logText, "Rows are number, customer, status, and cents.") {
		t.Fatalf("log\n%s", logText)
	}
	if strings.Contains(diffText, runRoot) || strings.Contains(logText, runRoot) {
		t.Fatal("git evidence kept the run root")
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func runGit(t *testing.T, dir string, env []string, args ...string) string {
	t.Helper()
	res, err := proc.Run(context.Background(), proc.Spec{
		Argv: append([]string{"git"}, args...),
		Dir:  dir,
		Env:  env,
	})
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("git %v: %v\n%s", args, err, res.Stderr)
	}
	return string(res.Stdout)
}
