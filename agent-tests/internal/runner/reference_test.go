package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"testing"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

// mutationHost is the scenario whose reference world each check is broken on.
// Shared checks use the orient world.
var mutationHost = map[string]string{
	"launcher-only":             "orient",
	"nothing-pushed":            "orient",
	"no-stray-commits":          "orient",
	"no-attribution":            "orient",
	"check-clean":               "orient",
	"sources-passive":           "orient",
	"no-outside-writes":         "orient",
	"orient-reply":              "orient",
	"orient-clean-trees":        "orient",
	"join-members":              "join",
	"join-local-files":          "join",
	"intent-record":             "intent",
	"intent-no-plan":            "intent",
	"intent-lang-body":          "intent-lang",
	"plan-record":               "plan",
	"plan-linked":               "plan",
	"execute-worktree":          "execute",
	"execute-tests-pass":        "execute",
	"execute-records-untouched": "execute",
	"complete-plan-done":        "complete",
	"complete-note":             "complete",
	"direct-readme":             "direct",
	"direct-no-records":         "direct",
	"borrowed-untouched":        "borrowed",
	"borrowed-reply":            "borrowed",
}

func TestMutationCoverage(t *testing.T) {
	names := scenario.Names()
	if len(names) != len(mutationHost) {
		t.Fatalf("mutations %d checks %d", len(mutationHost), len(names))
	}
	for _, name := range names {
		if mutationHost[name] == "" {
			t.Errorf("no mutation for %s", name)
		}
	}
}

func TestReferenceSolutions(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go missing")
	}
	t.Parallel()
	root := repoRoot(t)
	pin, err := world.Pin(root)
	if err != nil {
		t.Fatal(err)
	}
	parent := tempOutsideHome(t)
	buildRoot, err := world.NewRunRoot(parent)
	if err != nil {
		t.Fatal(err)
	}
	if err := world.Guard(buildRoot, root, os.Getenv("HOME")); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	cli, err := world.BuildCLI(ctx, root, buildRoot, pin)
	if err != nil {
		t.Fatal(err)
	}
	fixtures := filepath.Join(root, "agent-tests", "fixtures")
	loaded, err := scenario.LoadDir(filepath.Join(root, "agent-tests", "scenarios"))
	if err != nil {
		t.Fatal(err)
	}
	solutions := filepath.Join(root, "agent-tests", "testdata", "solutions")
	for _, sc := range loaded {
		t.Run(sc.ID, func(t *testing.T) {
			t.Parallel()
			st := prepareWorld(t, root, pin, cli, fixtures)
			if err := st.rt.runSnippets(context.Background(), st.w, sc); err != nil {
				t.Fatal(err)
			}
			applyWorkspace(&st.w, sc)
			checkBefore, err := st.rt.captureCheck(context.Background(), st.w)
			if err != nil {
				t.Fatal(err)
			}
			before, err := st.rt.snapshot(context.Background(), st.w)
			if err != nil {
				t.Fatal(err)
			}
			st.before = before
			st.sc = sc
			assertStart(t, st, checkBefore)
			stderr := runSolution(t, st, filepath.Join(solutions, sc.ID+".sh"), streamFor(sc.ID))
			verdict, results, graded := st.grade(t)
			if verdict != "pass" {
				t.Fatalf("verdict %s\n%s\ncheck %s\nsolution:\n%s", verdict, formatResults(results), graded.CheckAfter, stderr)
			}
			if sc.ID == "execute" {
				assertExecuteGit(t, st)
			}
			assertSoft(t, st, graded, verdict, results)
			for _, name := range scenario.Names() {
				if mutationHost[name] != sc.ID {
					continue
				}
				t.Run(name, func(t *testing.T) {
					undo := st.breakCheck(t, name)
					t.Cleanup(undo)
					got, gotResults, _ := st.grade(t)
					if got != "fail" || !checkFailed(gotResults, name) {
						t.Fatalf("mutation %s: verdict %s, want check %s to fail\n%s", name, got, name, formatResults(gotResults))
					}
				})
			}
		})
	}
}

func TestRunRealScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go missing")
	}
	t.Parallel()
	root := repoRoot(t)
	for _, id := range []string{"intent", "direct"} {
		t.Run(id, func(t *testing.T) {
			dir := tempOutsideHome(t)
			streamPath := filepath.Join(dir, "stream.jsonl")
			if err := os.WriteFile(streamPath, []byte(streamFor(id)), 0o644); err != nil {
				t.Fatal(err)
			}
			solution := filepath.Join(root, "agent-tests", "testdata", "solutions", id+".sh")
			stub := filepath.Join(dir, "stub.sh")
			script := "#!/bin/sh\nset -eu\nif [ \"${1:-}\" = \"--version\" ]; then echo stub-9; exit 0; fi\nsh " + shSingle(solution) + " >&2\ncat " + shSingle(streamPath) + "\n"
			if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
				t.Fatal(err)
			}
			hostsFile := writeHostsFile(t, dir, []hostRow{{
				name: "claude", command: stub, model: "stub-model", version: "stub", timeout: 180,
				auth: []string{"ANTHROPIC_API_KEY"},
			}})
			prices := filepath.Join(dir, "prices.yaml")
			if err := os.WriteFile(prices, []byte("prices: {}\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(dir, "out")
			var buf strings.Builder
			err := Run(context.Background(), RunOptions{
				Checkout:     root,
				Ref:          "HEAD",
				ScenarioIDs:  []string{id},
				HostNames:    []string{"claude"},
				Runs:         1,
				HostsFile:    hostsFile,
				PricesFile:   prices,
				ScenariosDir: filepath.Join(root, "agent-tests", "scenarios"),
				FixturesDir:  filepath.Join(root, "agent-tests", "fixtures"),
				OutDir:       out,
				Caller:       map[string]string{"PATH": os.Getenv("PATH"), "ANTHROPIC_API_KEY": "stub-anthropic-auth-value"},
				ParentTemp:   dir,
				Confirm:      true,
				Stdout:       &buf,
			})
			if err != nil {
				t.Fatalf("run: %v\n%s", err, buf.String())
			}
			logPath := matchOne(t, filepath.Join(out, "*", "runs.jsonl"))
			raw, err := os.ReadFile(logPath)
			if err != nil {
				t.Fatal(err)
			}
			var line struct {
				Verdict string                     `json:"verdict"`
				Soft    map[string]json.RawMessage `json:"soft"`
			}
			if err := json.Unmarshal(bytesTrimLine(raw), &line); err != nil {
				t.Fatalf("runs.jsonl: %v\n%s", err, raw)
			}
			if line.Verdict != "pass" {
				t.Fatalf("verdict %s\n%s", line.Verdict, raw)
			}
			repro := globOne(t, out, id, "claude", "reproduce.txt")
			text, err := os.ReadFile(repro)
			if err != nil {
				t.Fatal(err)
			}
			want := "agent-tests run --scenarios " + id + " --hosts claude --runs 1 --ref HEAD\n"
			if !strings.HasPrefix(string(text), "# ") || !strings.Contains(string(text), want) {
				t.Fatalf("reproduce %q", text)
			}
			if _, err := os.Stat(globOne(t, out, id, "claude", "reply.md")); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(globOne(t, out, id, "claude", "stream.jsonl")); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(globOne(t, out, id, "claude", "check.json")); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(globOne(t, out, id, "claude", "check-before.json")); err != nil {
				t.Fatal(err)
			}
			switch id {
			case "intent":
				if string(line.Soft["code-canary-seen"]) != "true" {
					t.Fatalf("soft %s", line.Soft["code-canary-seen"])
				}
				intent := globOne(t, out, id, "claude", filepath.Join("artifacts", "files", "intent", "i100-*.md"))
				body, err := os.ReadFile(intent)
				if err != nil || !strings.Contains(string(body), "## Goal") {
					t.Fatalf("intent artifact %v\n%s", err, body)
				}
			case "direct":
				readme := globOne(t, out, id, "claude", filepath.Join("artifacts", "files", "repositories", "api", "README.md"))
				body, err := os.ReadFile(readme)
				if err != nil {
					t.Fatal(err)
				}
				if strings.Contains(string(body), "recieve") || !strings.Contains(string(body), "receive") {
					t.Fatalf("readme %s", body)
				}
			}
		})
	}
}

type liveWorld struct {
	rt       *runtime
	w        world.World
	sc       scenario.Scenario
	before   scenario.Snapshot
	parsed   hosts.Parsed
	stream   string
	solution string
}

func prepareWorld(t *testing.T, root, pin, cliSrc, fixtures string) *liveWorld {
	t.Helper()
	parent := tempOutsideHome(t)
	runRoot, err := world.NewRunRoot(parent)
	if err != nil {
		t.Fatal(err)
	}
	if err := world.Guard(runRoot, root, os.Getenv("HOME")); err != nil {
		t.Fatal(err)
	}
	cli := installCLI(t, runRoot, pin, cliSrc)
	w, err := world.Setup(context.Background(), world.SetupInput{
		RunRoot:     runRoot,
		SrcDir:      root,
		CLIPath:     cli,
		FixturesDir: fixtures,
		Pin:         pin,
		PATH:        os.Getenv("PATH"),
	})
	if err != nil {
		t.Fatal(err)
	}
	rt := &runtime{opts: RunOptions{
		Caller:      map[string]string{"PATH": os.Getenv("PATH")},
		FixturesDir: fixtures,
		Checkout:    root,
	}}
	return &liveWorld{rt: rt, w: w}
}

func installCLI(t *testing.T, runRoot, pin, src string) string {
	t.Helper()
	store := filepath.Join(runRoot, "bin", ".context-circuit-versions", pin+"-"+goruntime.GOOS+"-"+goruntime.GOARCH)
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(store, "context-circuit-cli")
	body, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, body, 0o755); err != nil {
		t.Fatal(err)
	}
	return dst
}

func (st *liveWorld) grade(t *testing.T) (string, []scenario.CheckResult, scenario.CheckContext) {
	t.Helper()
	ctx := context.Background()
	after, err := st.rt.snapshot(ctx, st.w)
	if err != nil {
		t.Fatal(err)
	}
	checkAfter, err := st.rt.runCheck(ctx, st.w)
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	graded := scenario.CheckContext{
		World:      st.w,
		Scenario:   st.sc,
		Parsed:     st.parsed,
		Before:     st.before,
		After:      after,
		CheckAfter: checkAfter,
		Stream:     st.stream,
		Host:       "claude",
		Exec: func(dir string, argv ...string) (string, error) {
			return st.rt.execIsolated(ctx, st.w, dir, argv...)
		},
	}
	results, incomplete := grade(st.sc, graded)
	if !st.parsed.Terminal {
		incomplete = true
	}
	return scenario.Verdict(results, nil, incomplete), results, graded
}

func (st *liveWorld) git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := st.rt.execIsolated(context.Background(), st.w, dir, append([]string{"git"}, args...)...)
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(out)
}

func runSolution(t *testing.T, st *liveWorld, solution, stream string) string {
	t.Helper()
	dir := filepath.Join(st.w.RunRoot, "tmp")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(st.w.RunRoot, "claude-home"), 0o755); err != nil {
		t.Fatal(err)
	}
	streamPath := filepath.Join(dir, "stream.jsonl")
	if err := os.WriteFile(streamPath, []byte(stream), 0o644); err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(dir, "wrapper.sh")
	body := "#!/bin/sh\nset -eu\nsh " + shSingle(solution) + " >&2\ncat " + shSingle(streamPath) + "\n"
	if err := os.WriteFile(wrapper, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	res, err := proc.Run(context.Background(), proc.Spec{
		Argv: []string{"sh", wrapper},
		Dir:  st.w.WS,
		Env:  hostLikeEnv(st.w),
	})
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("solution exit %d: %v\n%s", res.ExitCode, err, res.Stderr)
	}
	adapter, err := hosts.ForName("claude")
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := adapter.Parse(strings.NewReader(string(res.Stdout)))
	if err != nil {
		t.Fatal(err)
	}
	if !parsed.Terminal || !parsed.ReadsAvailable {
		t.Fatalf("parsed terminal %v reads %v", parsed.Terminal, parsed.ReadsAvailable)
	}
	st.parsed = parsed
	st.stream = string(res.Stdout)
	st.solution = string(res.Stderr)
	return st.solution
}

func hostLikeEnv(w world.World) []string {
	base := hosts.BaseEnv(hosts.LaunchInput{
		RunRoot:  w.RunRoot,
		DecoyDir: w.Decoy,
		Caller:   map[string]string{"PATH": os.Getenv("PATH")},
	}, map[string]string{
		"HOME":              filepath.Join(w.RunRoot, "claude-home"),
		"CLAUDE_CONFIG_DIR": filepath.Join(w.RunRoot, "claude"),
	})
	return mergeMissing(base, w.Env)
}

func assertStart(t *testing.T, st *liveWorld, checkBefore []byte) {
	t.Helper()
	issues := issueTexts(t, checkBefore)
	switch st.sc.ID {
	case "orient", "intent", "direct":
		if len(issues) != 0 {
			t.Fatalf("start issues %s", checkBefore)
		}
		if st.sc.ID != "intent" && !strings.Contains(readFile(t, filepath.Join(st.w.API, "README.md")), "recieve") {
			t.Fatal("readme typo missing")
		}
		if st.sc.ID != "direct" && len(globRel(t, st.w.WS, "intent", "i100-*.md")) != 0 {
			t.Fatal("intent already exists")
		}
	case "join":
		if _, err := os.Stat(filepath.Join(st.w.WS, "member.local.yaml")); !os.IsNotExist(err) {
			t.Fatal("member.local.yaml exists on the clone")
		}
		if strings.Contains(readFile(t, filepath.Join(st.w.WS, "members.yaml")), "dana") {
			t.Fatal("dana is already on the roster")
		}
		if _, err := os.Stat(filepath.Join(st.w.API, ".git")); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(checkBefore), "select a workspace member") {
			t.Fatalf("start check %s", checkBefore)
		}
		raw := runStatus(t, st)
		var doc struct {
			NextBand int `json:"next_band"`
		}
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatalf("status: %v\n%s", err, raw)
		}
		if doc.NextBand != 2 {
			t.Fatalf("next_band %d\n%s", doc.NextBand, raw)
		}
	case "intent-lang":
		if !strings.Contains(readFile(t, filepath.Join(st.w.WS, "members.yaml")), "Indonesian") {
			t.Fatal("language was not recorded")
		}
		if len(globRel(t, st.w.WS, "intent", "i100-*.md")) != 0 {
			t.Fatal("intent already exists")
		}
	case "plan":
		intent := oneGlob(t, st.w.WS, "intent", "i100-*.md")
		if !strings.Contains(readFile(t, intent), "approved_at:") {
			t.Fatal("intent is not approved")
		}
		if len(globRel(t, st.w.WS, "plans", "p1000-*", "plan.md")) != 0 {
			t.Fatal("plan already exists")
		}
	case "execute":
		if len(globRel(t, st.w.WS, "plans", "p1000-*", "plan.md")) != 1 {
			t.Fatal("plan missing")
		}
		if strings.Contains(st.before.Refs["api"], "cc/p1000/") {
			t.Fatalf("branch already exists\n%s", st.before.Refs["api"])
		}
	case "complete":
		if len(issues) != 1 || !strings.Contains(string(checkBefore), "needs a person") || !strings.Contains(string(checkBefore), "invoices") {
			t.Fatalf("start check %s", checkBefore)
		}
		plan := oneGlob(t, st.w.WS, "plans", "p1000-*", "plan.md")
		if strings.Contains(readFile(t, plan), "completed_at:") {
			t.Fatal("plan is already complete")
		}
		if !strings.Contains(readFile(t, filepath.Join(st.w.API, "internal", "invoice", "export.go")), "CSV") {
			t.Fatal("merged tree has no CSV export")
		}
		res, err := proc.Run(context.Background(), proc.Spec{
			Argv: []string{"go", "test", "./..."},
			Dir:  st.w.API,
			Env:  append(snippetEnv(st.w, os.Getenv("PATH")), "GOTOOLCHAIN=local", "GOPROXY=off"),
		})
		if err != nil || res.ExitCode != 0 {
			t.Fatalf("patched api tests: %v\n%s", err, res.Stderr)
		}
	case "archive":
		if len(issues) != 3 || !strings.Contains(string(checkBefore), "legacy archive folder awaits its rename to _archive/") || !strings.Contains(string(checkBefore), "invoices") {
			t.Fatalf("start check %s", checkBefore)
		}
		if len(globRel(t, st.w.WS, "plans", "archive", "p1001-*", "plan.md")) != 1 || len(globRel(t, st.w.WS, "intent", "archive", "i101-*.md")) != 1 {
			t.Fatal("finished work is not in the legacy archive")
		}
		if strings.Contains(readFile(t, oneGlob(t, st.w.WS, "plans", "p1000-*", "plan.md")), "completed_at:") {
			t.Fatal("plan is already complete")
		}
	case "borrowed":
		note := filepath.Join(st.w.WS, "repositories", "handbook", "numbering.md")
		if !strings.Contains(readFile(t, note), "Invoice numbers are global and never restart.") {
			t.Fatal("handbook note missing")
		}
		if len(issues) != 0 {
			t.Fatalf("start check %s", checkBefore)
		}
	default:
		t.Fatalf("no start assertion for %s", st.sc.ID)
	}
}

func assertExecuteGit(t *testing.T, st *liveWorld) {
	t.Helper()
	dest := filepath.Join(st.w.RunRoot, "artifacts", "git")
	if err := st.rt.writeGitState(context.Background(), st.w, st.before, st.sc, dest); err != nil {
		t.Fatal(err)
	}
	diffText := readFile(t, filepath.Join(dest, "api-diff-stat.txt"))
	logText := readFile(t, filepath.Join(dest, "api-log.txt"))
	if strings.Contains(diffText, "exit ") || !strings.Contains(diffText, "export.go") {
		t.Fatalf("diff\n%s", diffText)
	}
	if !strings.Contains(logText, "feat(invoice): export invoices as CSV") {
		t.Fatalf("log\n%s", logText)
	}
	if strings.Contains(diffText, st.w.RunRoot) || strings.Contains(logText, st.w.RunRoot) {
		t.Fatal("git evidence kept the run root")
	}
}

func assertSoft(t *testing.T, st *liveWorld, graded scenario.CheckContext, verdict string, results []scenario.CheckResult) {
	t.Helper()
	if len(st.sc.Soft) == 0 {
		return
	}
	line := st.rt.recordLine(planned{Scenario: st.sc, Host: "claude", Config: hosts.HostConfig{Name: "claude"}}, st.parsed, verdict, results, 0, "", nil, &graded)
	switch st.sc.ID {
	case "intent", "intent-lang":
		if line.Soft["code-canary-seen"] != true {
			t.Fatalf("soft %#v", line.Soft)
		}
	case "plan":
		if line.Soft["read-api-code"] != true || line.Soft["subagents"] != 0 {
			t.Fatalf("soft %#v", line.Soft)
		}
	case "execute":
		if line.Soft["commit-subjects-conventional"] != true {
			t.Fatalf("soft %#v", line.Soft)
		}
	}
}

func (st *liveWorld) breakCheck(t *testing.T, name string) func() {
	t.Helper()
	switch name {
	case "launcher-only":
		prev := readOptional(st.w.DecoyLog)
		if err := os.WriteFile(st.w.DecoyLog, []byte("context-circuit-cli check\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return func() { _ = os.WriteFile(st.w.DecoyLog, prev, 0o644) }
	case "nothing-pushed":
		sha := st.git(t, st.w.API, "rev-parse", "HEAD")
		st.git(t, st.w.API, "commit", "--allow-empty", "-m", "chore: extra")
		st.git(t, st.w.API, "push", "origin", "main")
		bare := filepath.Join(st.w.RunRoot, "remotes", "api.git")
		return func() {
			st.git(t, st.w.API, "reset", "--hard", sha)
			st.git(t, bare, "update-ref", "refs/heads/main", sha)
		}
	case "no-stray-commits":
		st.git(t, st.w.API, "switch", "-c", "side/topic")
		st.git(t, st.w.API, "commit", "--allow-empty", "-m", "chore: stray")
		st.git(t, st.w.API, "switch", "main")
		return func() {
			st.git(t, st.w.API, "switch", "main")
			st.git(t, st.w.API, "branch", "-D", "side/topic")
		}
	case "no-attribution":
		st.git(t, st.w.WS, "commit", "--allow-empty", "-m", "chore: note\n\nCo-Authored-By: Helper <helper@example.com>")
		return func() { st.git(t, st.w.WS, "reset", "--hard", "HEAD~1") }
	case "check-clean":
		index := filepath.Join(st.w.WS, "context", "INDEX.md")
		prev := readFile(t, index)
		extra := "\n- [Missing](domains/missing.md) {api} — a note that is not in the tree · missing · reviewed 2026-01-15\n"
		if err := os.WriteFile(index, []byte(prev+extra), 0o644); err != nil {
			t.Fatal(err)
		}
		return func() { _ = os.WriteFile(index, []byte(prev), 0o644) }
	case "sources-passive":
		prev := st.stream
		st.stream += "\nSOURCES-CANARY-7f3a\n"
		return func() { st.stream = prev }
	case "no-outside-writes":
		path := filepath.Join(st.w.RunRoot, "stray.txt")
		if err := os.WriteFile(path, []byte("stray\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return func() { _ = os.Remove(path) }
	case "orient-reply":
		prevS, prevP := st.stream, st.parsed
		st.stream = "{\"type\":\"result\",\"is_error\":false,\"result\":\"No.\",\"num_turns\":1,\"session_id\":\"ref\"}\n"
		st.parsed = hosts.Parsed{FinalReply: "No.", ReadsAvailable: true, Terminal: true}
		return func() { st.stream, st.parsed = prevS, prevP }
	case "orient-clean-trees":
		return editFile(t, filepath.Join(st.w.API, "README.md"), func(s string) string { return s + "\n" })
	case "join-members":
		return editFile(t, filepath.Join(st.w.WS, "members.yaml"), func(s string) string {
			return strings.Replace(s, "band: 2", "band: 9", 1)
		})
	case "join-local-files":
		path := filepath.Join(st.w.WS, "member.local.yaml")
		prev := readFile(t, path)
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		return func() { _ = os.WriteFile(path, []byte(prev), 0o644) }
	case "intent-record":
		return editFile(t, oneGlob(t, st.w.WS, "intent", "i100-*.md"), func(s string) string {
			return insertAfter(s, "created_by:", "approved_at: 2026-10-26T00:00:00Z")
		})
	case "intent-no-plan":
		dir := filepath.Join(st.w.WS, "plans", "p1000-stray")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "note.txt"), []byte("stray\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return func() { _ = os.RemoveAll(dir) }
	case "intent-lang-body":
		body := readFile(t, filepath.Join(repoRoot(t), "agent-tests", "fixtures", "records", "intent-csv.md"))
		return editFile(t, oneGlob(t, st.w.WS, "intent", "i100-*.md"), func(s string) string {
			return replaceRecordBody(s, body)
		})
	case "plan-record":
		return editFile(t, oneGlob(t, st.w.WS, "plans", "p1000-*", "plan.md"), func(s string) string {
			return strings.Replace(s, "intent: i100", "intent: i101", 1)
		})
	case "plan-linked":
		return editFile(t, oneGlob(t, st.w.WS, "intent", "i100-*.md"), func(s string) string {
			return dropExactLine(s, "- p1000")
		})
	case "execute-worktree":
		wt := filepath.Join(st.w.WS, ".worktrees", "p1000", "api")
		sha := st.git(t, wt, "rev-parse", "HEAD")
		st.git(t, wt, "reset", "--hard", "main")
		return func() { st.git(t, wt, "reset", "--hard", sha) }
	case "execute-tests-pass":
		wt := filepath.Join(st.w.WS, ".worktrees", "p1000", "api")
		sha := st.git(t, wt, "rev-parse", "HEAD")
		path := filepath.Join(wt, "internal", "invoice", "export_test.go")
		raw := readFile(t, path)
		next := strings.Replace(raw, "func TestCSV(t *testing.T) {", "func TestCSV(t *testing.T) {\n\tt.Fatal(\"broken\")", 1)
		if next == raw {
			t.Fatal("TestCSV not found")
		}
		if err := os.WriteFile(path, []byte(next), 0o644); err != nil {
			t.Fatal(err)
		}
		st.git(t, wt, "add", "-A")
		st.git(t, wt, "commit", "-qm", "test(invoice): break the export test")
		return func() { st.git(t, wt, "reset", "--hard", sha) }
	case "execute-records-untouched":
		return editFile(t, oneGlob(t, st.w.WS, "plans", "p1000-*", "plan.md"), func(s string) string {
			return insertAfter(s, "created_at:", "completed_at: 2026-10-26T00:00:00Z")
		})
	case "complete-plan-done":
		return editFile(t, oneGlob(t, st.w.WS, "plans", "p1000-*", "plan.md"), stripCompletion)
	case "complete-note":
		return editFile(t, filepath.Join(st.w.WS, "context", "domains", "invoices.md"), func(s string) string {
			return strings.NewReplacer("CSV", "text", "csv", "text").Replace(s)
		})
	case "direct-readme":
		return editFile(t, filepath.Join(st.w.API, "README.md"), func(s string) string {
			return strings.ReplaceAll(s, "receive", "recieve")
		})
	case "direct-no-records":
		path := filepath.Join(st.w.WS, "intent", "stray.txt")
		if err := os.WriteFile(path, []byte("stray\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return func() { _ = os.Remove(path) }
	case "borrowed-untouched":
		return editFile(t, filepath.Join(st.w.WS, "repositories", "handbook", "numbering.md"), func(s string) string {
			return s + "\nEdited.\n"
		})
	case "borrowed-reply":
		prevS, prevP := st.stream, st.parsed
		st.parsed.FinalReply = "Done, I updated the note."
		st.stream = "{\"type\":\"result\",\"is_error\":false,\"result\":\"Done, I updated the note.\",\"num_turns\":1,\"session_id\":\"ref\"}\n"
		st.parsed.Terminal = true
		st.parsed.ReadsAvailable = true
		return func() { st.stream, st.parsed = prevS, prevP }
	default:
		t.Fatalf("no mutation for %s", name)
		return func() {}
	}
}

func streamFor(id string) string {
	reply := map[string]string{
		"orient":      "This workspace is Acme, for invoicing for small businesses. It works on the api repository. Run the command with sh .context-circuit/cli. A draft invoice is voided. An issued invoice is credited. The note is context/domains/invoices.md.",
		"join":        "Dana is on band 2 and nothing was committed.",
		"intent":      "I wrote the intent for invoice CSV export. Please approve it before anyone builds it.",
		"intent-lang": "Intent sudah ditulis. Please approve it before anyone builds it.",
		"plan":        "The plan is in the plan folder and it stops before implementation.",
		"execute":     "The CSV export is on the worktree branch. Tests passed. Nothing was pushed.",
		"complete":    "The plan is marked complete and the invoices note mentions CSV. Those edits are uncommitted, and committing them is recommended.",
		"direct":      "The README typo is fixed and the change is not committed.",
		"borrowed":    "The handbook is read-only here. The correction belongs upstream, in the handbook repository's own review.",
		"archive":     "The plan is marked complete and the invoices note mentions CSV. The old archive/ folders are now _archive/. Commit and push those moves before anyone else changes the workspace; teammates install CLI 2.5.0 before pulling. You have finished records that are not archived; shall I archive them?",
	}[id]
	tool := "context/INDEX.md"
	extra := ""
	switch id {
	case "intent", "intent-lang":
		extra = "CODE-CANARY-41b2"
	case "plan":
		tool = "repositories/api/internal/invoice/export.go"
	case "execute":
		tool = ".worktrees/p1000/api/internal/invoice/export.go"
	}
	return claudeStream(reply, tool, extra)
}

func claudeStream(reply, toolPath, extra string) string {
	var b strings.Builder
	if toolPath != "" {
		fmt.Fprintf(&b, "{\"type\":\"assistant\",\"message\":{\"content\":[{\"type\":\"tool_use\",\"name\":\"Read\",\"input\":{\"file_path\":%s}}]}}\n", jsonString(toolPath))
	}
	if extra != "" {
		fmt.Fprintf(&b, "{\"type\":\"user\",\"message\":{\"content\":%s}}\n", jsonString(extra))
	}
	fmt.Fprintf(&b, "{\"type\":\"result\",\"is_error\":false,\"result\":%s,\"num_turns\":1,\"session_id\":\"ref\"}\n", jsonString(reply))
	return b.String()
}

func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func issueTexts(t *testing.T, raw []byte) []string {
	t.Helper()
	var doc struct {
		Issues []struct {
			Issue string `json:"issue"`
		} `json:"issues"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("check json: %v\n%s", err, raw)
	}
	var out []string
	for _, item := range doc.Issues {
		out = append(out, item.Issue)
	}
	return out
}

func runStatus(t *testing.T, st *liveWorld) []byte {
	t.Helper()
	res, err := proc.Run(context.Background(), proc.Spec{
		Argv: []string{"sh", filepath.Join(st.w.WS, ".context-circuit", "cli"), "--json", "status"},
		Dir:  st.w.WS,
		Env:  snippetEnv(st.w, os.Getenv("PATH")),
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 && len(res.Stdout) == 0 {
		t.Fatalf("status exit %d\n%s", res.ExitCode, res.Stderr)
	}
	return res.Stdout
}

func formatResults(results []scenario.CheckResult) string {
	var b strings.Builder
	for _, result := range results {
		if !result.Pass || result.Inconclusive {
			fmt.Fprintf(&b, "%s pass=%v inconclusive=%v %s\n", result.Name, result.Pass, result.Inconclusive, result.Detail)
		}
	}
	return b.String()
}

func checkFailed(results []scenario.CheckResult, name string) bool {
	for _, result := range results {
		if result.Name == name && !result.Pass {
			return true
		}
	}
	return false
}

func editFile(t *testing.T, path string, edit func(string) string) func() {
	t.Helper()
	raw := readFile(t, path)
	if err := os.WriteFile(path, []byte(edit(raw)), 0o644); err != nil {
		t.Fatal(err)
	}
	return func() { _ = os.WriteFile(path, []byte(raw), 0o644) }
}

func insertAfter(s, prefix, extra string) string {
	var b strings.Builder
	done := false
	for _, line := range strings.Split(s, "\n") {
		b.WriteString(line)
		b.WriteByte('\n')
		if !done && strings.HasPrefix(strings.TrimSpace(line), prefix) {
			b.WriteString(extra)
			b.WriteByte('\n')
			done = true
		}
	}
	if !done {
		return s + "\n" + extra + "\n"
	}
	return b.String()
}

func dropExactLine(s, want string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) == want {
			continue
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func stripCompletion(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "completed_at:") {
			continue
		}
		if strings.HasPrefix(trim, "## Completion") {
			break
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func replaceRecordBody(record, body string) string {
	const open = "---\n"
	if !strings.HasPrefix(record, open) {
		return record
	}
	rest := record[len(open):]
	i := strings.Index(rest, "\n---\n")
	if i < 0 {
		return record
	}
	return record[:len(open)+i+len("\n---\n")] + body
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func readOptional(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return b
}

func globRel(t *testing.T, root string, parts ...string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(append([]string{root}, parts...)...))
	if err != nil {
		t.Fatal(err)
	}
	return matches
}

func oneGlob(t *testing.T, root string, parts ...string) string {
	t.Helper()
	matches := globRel(t, root, parts...)
	if len(matches) != 1 {
		t.Fatalf("glob %s: %v", filepath.Join(parts...), matches)
	}
	return matches[0]
}

func bytesTrimLine(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

func matchOne(t *testing.T, pattern string) string {
	t.Helper()
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) != 1 {
		t.Fatalf("glob %s: %v %v", pattern, matches, err)
	}
	return matches[0]
}
