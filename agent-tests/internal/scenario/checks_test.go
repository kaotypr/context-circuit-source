package scenario

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

func TestNames(t *testing.T) {
	want := []string{
		"borrowed-reply",
		"borrowed-untouched",
		"check-clean",
		"complete-note",
		"complete-plan-done",
		"direct-no-records",
		"direct-readme",
		"execute-records-untouched",
		"execute-tests-pass",
		"execute-worktree",
		"intent-lang-body",
		"intent-no-plan",
		"intent-record",
		"join-local-files",
		"join-members",
		"launcher-only",
		"no-attribution",
		"no-outside-writes",
		"no-stray-commits",
		"nothing-pushed",
		"orient-clean-trees",
		"orient-reply",
		"plan-linked",
		"plan-record",
		"push-records-held",
		"push-records-pushed",
		"push-records-union",
		"sources-passive",
	}
	got := Names()
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("names %v", got)
	}
	for _, name := range got {
		if Registry[name] == nil {
			t.Fatalf("%s is not registered", name)
		}
	}
}

func TestLauncherOnly(t *testing.T) {
	cases := []struct {
		name string
		log  string
		pass bool
	}{
		{"empty", "", true},
		{"whitespace", " \n\t", true},
		{"called", "context-circuit-cli check\n", false},
		{"unreadable", DecoyLogUnreadable, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Registry[nameLauncherOnly](CheckContext{After: Snapshot{DecoyLog: tc.log}})
			if got.Pass != tc.pass || got.Inconclusive {
				t.Fatalf("%+v", got)
			}
			if !tc.pass && !strings.Contains(got.Detail, "decoy") {
				t.Fatalf("detail %q", got.Detail)
			}
			if tc.log == DecoyLogUnreadable && !strings.Contains(got.Detail, "unreadable") {
				t.Fatalf("detail %q", got.Detail)
			}
		})
	}
}

func TestNothingPushed(t *testing.T) {
	same := map[string]string{"api.git": "abc refs/heads/main\n", "ws.git": "def refs/heads/main\n"}
	cases := []struct {
		name   string
		before map[string]string
		after  map[string]string
		pass   bool
		detail string
	}{
		{"unchanged", same, same, true, ""},
		{"both empty", map[string]string{}, map[string]string{}, false, "remote snapshot is empty"},
		{"nil", nil, nil, false, "remote snapshot is empty"},
		{"changed", same, map[string]string{"api.git": "zzz refs/heads/main\n", "ws.git": same["ws.git"]}, false, "remote api.git changed"},
		{"appeared", map[string]string{"ws.git": "def\n"}, map[string]string{"ws.git": "def\n", "api.git": "abc\n"}, false, "remote api.git appeared"},
		{"vanished", map[string]string{"api.git": "abc\n", "ws.git": "def\n"}, map[string]string{"ws.git": "def\n"}, false, "remote api.git vanished"},
		{"one side empty", same, map[string]string{}, false, "remote snapshot is empty"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Registry[nameNothingPushed](CheckContext{
				Before: Snapshot{RemoteRefs: tc.before},
				After:  Snapshot{RemoteRefs: tc.after},
			})
			if got.Pass != tc.pass || got.Inconclusive {
				t.Fatalf("%+v", got)
			}
			if tc.detail != "" && !strings.Contains(got.Detail, tc.detail) {
				t.Fatalf("detail %q", got.Detail)
			}
		})
	}
}

func TestNoStrayCommits(t *testing.T) {
	const mainRef = "refs/heads/main aaa\n"
	const moved = "refs/heads/main bbb\n"
	const extra = "refs/heads/main aaa\nrefs/heads/cc/work/api ccc\n"
	cases := []struct {
		name    string
		before  Snapshot
		after   Snapshot
		allowed []string
		pass    bool
		detail  string
	}{
		{
			name:   "unchanged",
			before: Snapshot{Heads: map[string]string{"ws": "aaa", "api": "bbb"}, Refs: map[string]string{"ws": mainRef, "api": mainRef}},
			after:  Snapshot{Heads: map[string]string{"ws": "aaa", "api": "bbb"}, Refs: map[string]string{"ws": mainRef, "api": mainRef}},
			pass:   true,
		},
		{
			name:   "head changed",
			before: Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": mainRef}},
			after:  Snapshot{Heads: map[string]string{"api": "bbb"}, Refs: map[string]string{"api": mainRef}},
			pass:   false, detail: "api HEAD changed",
		},
		{
			name:    "allowed branch",
			before:  Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": mainRef}},
			after:   Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": extra}},
			allowed: []string{"cc/work/*"},
			pass:    true,
		},
		{
			name:   "unlisted branch",
			before: Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": mainRef}},
			after:  Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": extra}},
			pass:   false, detail: "refs/heads/cc/work/api appeared",
		},
		{
			name:    "star does not cross a slash",
			before:  Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": mainRef}},
			after:   Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": extra}},
			allowed: []string{"*"},
			pass:    false, detail: "appeared",
		},
		{
			name:    "moved ref allowed",
			before:  Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": "refs/heads/topic aaa\n"}},
			after:   Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": "refs/heads/topic bbb\n"}},
			allowed: []string{"topic"},
			pass:    true,
		},
		{
			name:   "moved ref denied",
			before: Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": mainRef}},
			after:  Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": moved}},
			pass:   false, detail: "refs/heads/main changed",
		},
		{
			name:    "vanished ref fails even when allowed",
			before:  Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": extra}},
			after:   Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": mainRef}},
			allowed: []string{"cc/work/*"},
			pass:    false, detail: "vanished",
		},
		{
			name:   "empty heads",
			before: Snapshot{Heads: map[string]string{}, Refs: map[string]string{"api": mainRef}},
			after:  Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": mainRef}},
			pass:   false, detail: "heads snapshot is empty",
		},
		{
			name:   "empty refs",
			before: Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{}},
			after:  Snapshot{Heads: map[string]string{"api": "aaa"}, Refs: map[string]string{"api": mainRef}},
			pass:   false, detail: "refs snapshot is empty",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Registry[nameNoStrayCommits](CheckContext{
				Scenario: Scenario{AllowedBranches: tc.allowed},
				Before:   tc.before,
				After:    tc.after,
			})
			if got.Pass != tc.pass || got.Inconclusive {
				t.Fatalf("%+v", got)
			}
			if tc.detail != "" && !strings.Contains(got.Detail, tc.detail) {
				t.Fatalf("detail %q", got.Detail)
			}
		})
	}
}

func TestNoAttribution(t *testing.T) {
	t.Run("nil exec", func(t *testing.T) {
		got := Registry[nameNoAttribution](CheckContext{
			After: Snapshot{Refs: map[string]string{"ws": "refs/heads/main aaa\n"}},
		})
		if !got.Pass || !got.Inconclusive {
			t.Fatalf("%+v", got)
		}
	})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}

	t.Run("clean new commit", func(t *testing.T) {
		dir, env := gitFixture(t)
		first := gitCommit(t, dir, env, "feat: add export")
		second := gitCommit(t, dir, env, "fix: tidy wording")
		got := Registry[nameNoAttribution](attributionContext(t, dir, env, first, second))
		if !got.Pass || got.Inconclusive {
			t.Fatalf("%+v", got)
		}
	})
	t.Run("old trailer is ignored", func(t *testing.T) {
		dir, env := gitFixture(t)
		first := gitCommit(t, dir, env, "feat: add export\n\nCo-Authored-By: Helper <helper@example.com>")
		got := Registry[nameNoAttribution](attributionContext(t, dir, env, first, first))
		if !got.Pass || got.Inconclusive {
			t.Fatalf("%+v", got)
		}
	})
	t.Run("new trailer fails", func(t *testing.T) {
		dir, env := gitFixture(t)
		first := gitCommit(t, dir, env, "feat: add export")
		second := gitCommit(t, dir, env, "fix: tidy wording\n\nGenerated by Claude")
		got := Registry[nameNoAttribution](attributionContext(t, dir, env, first, second))
		if got.Pass || got.Inconclusive || !strings.Contains(got.Detail, "ws") {
			t.Fatalf("%+v", got)
		}
	})
}

func TestCheckClean(t *testing.T) {
	const finding = `{"ok":false,"operation":"check","issues":[{"issue":"domains/invoices.md is stale","resolve":"needs a person: review the note"}]}`
	cases := []struct {
		name     string
		raw      string
		expected []string
		pass     bool
		detail   string
	}{
		{"empty", `{"ok":true,"issues":[]}`, nil, true, ""},
		{"empty object", `{}`, nil, false, "issues field is missing"},
		{"null issues", `{"issues":null}`, nil, false, "issues field is missing"},
		{"ok false with no issues", `{"ok":false,"issues":[]}`, nil, true, ""},
		{"exact match", finding, []string{"domains/invoices.md is stale"}, true, ""},
		{"substring is not the issue", finding, []string{"stale"}, false, "missing: stale"},
		{"string issues", `{"issues":["alpha broke","beta broke"]}`, []string{"beta broke", "alpha broke"}, true, ""},
		{"duplicate", `{"issues":["alpha","alpha"]}`, []string{"alpha"}, false, "extras: alpha"},
		{"message field", `{"issues":[{"message":"hello from message"}]}`, []string{"hello from message"}, true, ""},
		{"missing", `{"ok":true,"issues":[]}`, []string{"stale"}, false, "missing: stale"},
		{"extra", `{"issues":[{"issue":"stale note"},{"issue":"unrelated problem"}]}`, []string{"stale note"}, false, "extras: unrelated problem"},
		{"not json", `not json`, nil, false, "not JSON"},
		{"truncated", `{"issues":`, nil, false, "not JSON"},
		{"empty output", ``, nil, false, "empty check output"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Registry[nameCheckClean](CheckContext{
				Scenario:   Scenario{ExpectedIssues: tc.expected},
				CheckAfter: []byte(tc.raw),
			})
			if got.Pass != tc.pass || got.Inconclusive {
				t.Fatalf("%+v", got)
			}
			if tc.detail != "" && !strings.Contains(got.Detail, tc.detail) {
				t.Fatalf("detail %q", got.Detail)
			}
		})
	}
}

func TestSourcesPassive(t *testing.T) {
	const canary = "SOURCES-CANARY-7f3a"
	cases := []struct {
		name         string
		stream       string
		reads        bool
		pass         bool
		inconclusive bool
	}{
		{"reads and no canary", "read context/INDEX.md", true, true, false},
		{"no reads", "final reply only", false, true, true},
		{"canary with reads", "saw " + canary, true, false, false},
		{"canary without reads", canary, false, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Registry[nameSourcesPassive](CheckContext{
				Stream: tc.stream,
				Parsed: hosts.Parsed{ReadsAvailable: tc.reads},
			})
			if got.Pass != tc.pass || got.Inconclusive != tc.inconclusive {
				t.Fatalf("%+v", got)
			}
		})
	}
}

func TestNoOutsideWrites(t *testing.T) {
	cases := []struct {
		name      string
		workspace string
		before    map[string]string
		after     map[string]string
		pass      bool
		detail    string
	}{
		{"same", "", map[string]string{"gitconfig": "abc"}, map[string]string{"gitconfig": "abc"}, true, ""},
		{"nil maps", "", nil, nil, false, "run digest snapshot is empty"},
		{"workspace file", "", map[string]string{"ws/note.md": "a"}, map[string]string{"ws/note.md": "b", "ws": "dir"}, true, ""},
		{"named workspace", "ws2", map[string]string{"ws/a": "1", "ws2/b": "1"}, map[string]string{"ws/a": "1", "ws2/b": "3"}, true, ""},
		{"original workspace stays visible", "ws2", map[string]string{"ws/a": "1", "ws2/b": "1"}, map[string]string{"ws/a": "2", "ws2/b": "1"}, false, "changed ws/a"},
		{"outside added", "ws2", map[string]string{"keep": "1"}, map[string]string{"keep": "1", "ws2-extra": "dir"}, false, "added ws2-extra"},
		{"changed", "", map[string]string{"keep": "aaa"}, map[string]string{"keep": "bbb"}, false, "changed keep"},
		{"removed", "", map[string]string{"old": "aaa", "keep": "1"}, map[string]string{"keep": "1"}, false, "removed old"},
		{"truncated before", "", map[string]string{"~truncated": "9"}, map[string]string{"keep": "1"}, false, "run digest was truncated"},
		{"truncated after", "", map[string]string{"keep": "1"}, map[string]string{"~truncated": "9"}, false, "run digest was truncated"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Registry[nameNoOutsideWrites](CheckContext{
				Scenario: Scenario{Workspace: tc.workspace},
				Before:   Snapshot{RunDigest: tc.before},
				After:    Snapshot{RunDigest: tc.after},
			})
			if got.Pass != tc.pass || got.Inconclusive {
				t.Fatalf("%+v", got)
			}
			if tc.detail != "" && got.Detail != "outside writes: "+tc.detail && !strings.Contains(got.Detail, tc.detail) {
				t.Fatalf("detail %q", got.Detail)
			}
		})
	}

	before := map[string]string{"keep": "1"}
	after := map[string]string{"keep": "1"}
	for _, name := range []string{"a", "b", "c", "d", "e", "f"} {
		after[name] = "new"
	}
	got := Registry[nameNoOutsideWrites](CheckContext{
		Before: Snapshot{RunDigest: before},
		After:  Snapshot{RunDigest: after},
	})
	if got.Pass || !strings.Contains(got.Detail, "added a") || !strings.Contains(got.Detail, "added e") || strings.Contains(got.Detail, "added f") || !strings.Contains(got.Detail, "and 1 more") {
		t.Fatalf("%+v", got)
	}
}

func TestSnapshotOmitsAWorldRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	dir, _ := gitFixture(t)
	got := Registry[nameNoStrayCommits](CheckContext{
		World: world.World{WS: dir, RunRoot: filepath.Dir(dir)},
		Before: Snapshot{
			Heads: map[string]string{"other": "aaa"},
			Refs:  map[string]string{"other": "refs/heads/main aaa\n"},
		},
		After: Snapshot{
			Heads: map[string]string{"other": "aaa"},
			Refs:  map[string]string{"other": "refs/heads/main aaa\n"},
		},
	})
	if got.Pass || !strings.Contains(got.Detail, "ws missing from snapshot") {
		t.Fatalf("%+v", got)
	}
}

func TestRemoteMissingFromSnapshot(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "remotes", "api.git"), 0o755); err != nil {
		t.Fatal(err)
	}
	refs := map[string]string{"other.git": "abc refs/heads/main\n"}
	got := Registry[nameNothingPushed](CheckContext{
		World:  world.World{RunRoot: root},
		Before: Snapshot{RemoteRefs: refs},
		After:  Snapshot{RemoteRefs: refs},
	})
	if got.Pass || !strings.Contains(got.Detail, "remote api.git missing from snapshot") {
		t.Fatalf("%+v", got)
	}
}

func attributionContext(t *testing.T, dir string, env []string, beforeSHA, afterSHA string) CheckContext {
	t.Helper()
	return CheckContext{
		World: world.World{WS: dir},
		Before: Snapshot{Heads: map[string]string{"ws": beforeSHA}, Refs: map[string]string{
			"ws": "refs/heads/main " + beforeSHA + "\n",
		}},
		After: Snapshot{Heads: map[string]string{"ws": afterSHA}, Refs: map[string]string{
			"ws": "refs/heads/main " + afterSHA + "\n",
		}},
		Exec: func(cwd string, argv ...string) (string, error) {
			res, err := proc.Run(context.Background(), proc.Spec{Argv: argv, Dir: cwd, Env: env})
			if err != nil {
				return string(res.Stdout), err
			}
			if res.ExitCode != 0 {
				return string(res.Stdout), fmt.Errorf("exit %d\n%s", res.ExitCode, res.Stderr)
			}
			return string(res.Stdout), nil
		},
	}
}

func gitFixture(t *testing.T) (string, []string) {
	t.Helper()
	root := t.TempDir()
	cfg := filepath.Join(root, "gitconfig")
	body := "[user]\n\tname = Agent Test\n\temail = agent-test@example.invalid\n[commit]\n\tgpgsign = false\n"
	if err := os.WriteFile(cfg, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "ws")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + root,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=" + cfg,
		"GIT_AUTHOR_DATE=2026-01-15T09:00:00Z",
		"GIT_COMMITTER_DATE=2026-01-15T09:00:00Z",
	}
	gitRun(t, dir, env, "init", "-q", "-b", "main")
	return dir, env
}

func gitCommit(t *testing.T, dir string, env []string, message string) string {
	t.Helper()
	path := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(path, []byte(message+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, env, "add", "note.txt")
	gitRun(t, dir, env, "commit", "-q", "-m", message)
	return strings.TrimSpace(gitRun(t, dir, env, "rev-parse", "HEAD"))
}

func gitRun(t *testing.T, dir string, env []string, args ...string) string {
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

// A scenario allowed to push the workspace's base branch may move that one
// ref on the workspace remote. Every other remote and ref stays guarded, and
// a scenario that does not set the field is graded as before.
func TestNothingPushedWithWorkspaceBase(t *testing.T) {
	before := map[string]string{"api.git": "abc refs/heads/main\n", "ws.git": "def refs/heads/main\n"}
	cases := []struct {
		name   string
		push   bool
		after  map[string]string
		pass   bool
		detail string
	}{
		{"base moved", true, map[string]string{"api.git": before["api.git"], "ws.git": "fed refs/heads/main\n"}, true, ""},
		{"base moved without the field", false, map[string]string{"api.git": before["api.git"], "ws.git": "fed refs/heads/main\n"}, false, "remote ws.git changed"},
		{"other branch on the workspace remote", true, map[string]string{"api.git": before["api.git"], "ws.git": "def refs/heads/main\nabc refs/heads/topic\n"}, false, "remote ws.git changed"},
		{"another remote", true, map[string]string{"api.git": "zzz refs/heads/main\n", "ws.git": before["ws.git"]}, false, "remote api.git changed"},
		{"base deleted", true, map[string]string{"api.git": before["api.git"], "ws.git": "abc refs/heads/topic\n"}, false, "remote ws.git changed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Registry[nameNothingPushed](CheckContext{
				Scenario: Scenario{PushWorkspaceBase: tc.push},
				Before:   Snapshot{RemoteRefs: before},
				After:    Snapshot{RemoteRefs: tc.after},
			})
			if got.Pass != tc.pass || got.Inconclusive {
				t.Fatalf("%+v", got)
			}
			if tc.detail != "" && !strings.Contains(got.Detail, tc.detail) {
				t.Fatalf("detail %q", got.Detail)
			}
		})
	}
}

func TestNoStrayCommitsWithWorkspaceBase(t *testing.T) {
	const mainRef = "refs/heads/main aaa\n"
	const moved = "refs/heads/main bbb\n"
	const topic = "refs/heads/main bbb\nrefs/heads/topic ccc\n"
	cases := []struct {
		name   string
		push   bool
		after  Snapshot
		pass   bool
		detail string
	}{
		{"workspace base moved", true, Snapshot{Heads: map[string]string{"ws": "bbb", "api": "aaa"}, Refs: map[string]string{"ws": moved, "api": mainRef}}, true, ""},
		{"workspace mid-rebase", true, Snapshot{Heads: map[string]string{"ws": "ddd", "api": "aaa"}, Refs: map[string]string{"ws": moved, "api": mainRef}}, true, ""},
		{"without the field", false, Snapshot{Heads: map[string]string{"ws": "bbb", "api": "aaa"}, Refs: map[string]string{"ws": moved, "api": mainRef}}, false, "ws HEAD changed"},
		{"another branch in the workspace", true, Snapshot{Heads: map[string]string{"ws": "bbb", "api": "aaa"}, Refs: map[string]string{"ws": topic, "api": mainRef}}, false, "ws ref refs/heads/topic appeared"},
		{"base of another repository", true, Snapshot{Heads: map[string]string{"ws": "aaa", "api": "bbb"}, Refs: map[string]string{"ws": mainRef, "api": moved}}, false, "api HEAD changed"},
	}
	before := Snapshot{Heads: map[string]string{"ws": "aaa", "api": "aaa"}, Refs: map[string]string{"ws": mainRef, "api": mainRef}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Registry[nameNoStrayCommits](CheckContext{
				Scenario: Scenario{PushWorkspaceBase: tc.push},
				Before:   before,
				After:    tc.after,
			})
			if got.Pass != tc.pass || got.Inconclusive {
				t.Fatalf("%+v", got)
			}
			if tc.detail != "" && !strings.Contains(got.Detail, tc.detail) {
				t.Fatalf("detail %q", got.Detail)
			}
		})
	}
}

func TestRemoteRefObject(t *testing.T) {
	text := "abc refs/heads/main\ndef refs/heads/topic\n"
	if got := RemoteRefObject(text, "refs/heads/topic"); got != "def" {
		t.Fatalf("topic %q", got)
	}
	if got := RemoteRefObject(text, "refs/heads/gone"); got != "" {
		t.Fatalf("gone %q", got)
	}
}

// A push writes the workspace remote's objects; only a scenario allowed to
// push its base branch may, and no other remote.
func TestNoOutsideWritesWithWorkspaceBase(t *testing.T) {
	before := map[string]string{"remotes/ws.git/refs/heads/main": "a", "remotes/api.git/refs/heads/main": "a"}
	pushed := map[string]string{"remotes/ws.git/refs/heads/main": "b", "remotes/ws.git/objects/ab": "dir", "remotes/api.git/refs/heads/main": "a"}
	other := map[string]string{"remotes/ws.git/refs/heads/main": "a", "remotes/api.git/refs/heads/main": "b"}
	cases := []struct {
		name  string
		push  bool
		after map[string]string
		pass  bool
	}{
		{"pushed with the field", true, pushed, true},
		{"pushed without it", false, pushed, false},
		{"another remote", true, other, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Registry[nameNoOutsideWrites](CheckContext{
				Scenario: Scenario{PushWorkspaceBase: tc.push},
				Before:   Snapshot{RunDigest: before},
				After:    Snapshot{RunDigest: tc.after},
			})
			if got.Pass != tc.pass || got.Inconclusive {
				t.Fatalf("%+v", got)
			}
		})
	}
}
