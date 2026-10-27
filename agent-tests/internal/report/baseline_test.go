package report

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

type fakeRun struct {
	run   record.Run
	note  string
	reads string
	reply string
	extra map[string]string
}

func baseRecord(scenario, host, verdict, date string, n int) record.Run {
	in, out := int64(1000), int64(10)
	return record.Run{
		Scenario: scenario, ScenarioVersion: "1", Host: host, HostVersion: "1.0.0", HostVersionSource: "command",
		Model: "model-a", Effort: "medium", SrcCommit: "abc123", CLIVersion: "2.5.0", FixtureDigest: "sha256:00",
		Side: "branch", Run: n, Verdict: verdict, Date: date, FailedChecks: []string{},
		Tokens: record.Tokens{InputTotal: &in, Output: &out},
	}
}

// writeRunDir lays out a finished run directory the way the runner does.
func writeRunDir(t *testing.T, parent, id string, runs []fakeRun) string {
	t.Helper()
	dir := filepath.Join(parent, id)
	for _, item := range runs {
		if err := record.Append(filepath.Join(dir, "runs.jsonl"), item.run); err != nil {
			t.Fatal(err)
		}
		if item.run.Verdict == "harness-error" {
			continue
		}
		side := item.run.Side
		runDir := filepath.Join(dir, "runs", item.run.Scenario, item.run.Host, side, fmt.Sprintf("r%d", item.run.Run))
		files := map[string]string{
			"artifacts/files/note.md":      item.note,
			"artifacts/design/plan.md":     "# Plan\n",
			"artifacts/git/workspace.txt":  "clean\n",
			"reads.txt":                    item.reads,
			"reply.md":                     item.reply,
			"check.json":                   "{\"findings\":[]}\n",
			"stream.jsonl":                 "{\"type\":\"result\"}\n",
			"stderr.txt":                   "stderr\n",
			"snapshots/before.txt":         "before\n",
			"artifacts/files/sub/other.md": "other\n",
		}
		for rel, body := range item.extra {
			files[rel] = body
		}
		for rel, body := range files {
			path := filepath.Join(runDir, filepath.FromSlash(rel))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return dir
}

func noLeftovers(t *testing.T, parent string) {
	t.Helper()
	entries, err := os.ReadDir(parent)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".export-baseline-") || entry.Name() == "baseline" {
			t.Fatalf("export left %s behind", entry.Name())
		}
	}
}

func TestExportBaselineMergesAndDesignReadsIt(t *testing.T) {
	root := t.TempDir()
	runsRoot := filepath.Join(root, "runs")
	older := writeRunDir(t, runsRoot, "20261001T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "pass", "2026-10-01", 1), note: "old note\n", reads: "context/INDEX.md\n", reply: "old reply\n"},
		{run: baseRecord("overhead-empty", "claude", "pass", "2026-10-01", 1), note: "", reads: "", reply: "OK\n"},
	})
	harness := writeRunDir(t, runsRoot, "20261003T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "harness-error", "2026-10-03", 1)},
		{run: baseRecord("plan", "codex", "harness-error", "2026-10-03", 1)},
	})
	newer := writeRunDir(t, runsRoot, "20261002T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "fail", "2026-10-02", 1), note: "new note\n", reads: "context/INDEX.md\ncontext/decisions.md\n", reply: "newest reply\n"},
	})
	out := filepath.Join(root, "baseline")
	var buf bytes.Buffer
	err := ExportBaseline(ExportOptions{Runs: []string{older, harness, newer}, Out: out, Stdout: &buf})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "exported 2 run(s); 1 scenario and host pair(s) missing") {
		t.Fatalf("stdout %q", buf.String())
	}

	runs, err := LoadBaseline(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 2 {
		t.Fatalf("runs %d", len(runs))
	}
	for _, run := range runs {
		if run.Side != "base" {
			t.Fatalf("side %s", run.Side)
		}
		if run.Scenario == "intent" && (run.Verdict != "fail" || run.Date != "2026-10-02") {
			t.Fatalf("newest non-harness-error run was not chosen: %+v", run)
		}
	}
	reply, err := os.ReadFile(filepath.Join(out, "artifacts", "intent", "claude", "reply.md"))
	if err != nil || string(reply) != "newest reply\n" {
		t.Fatalf("reply %q %v", reply, err)
	}
	for _, rel := range []string{"files/note.md", "files/sub/other.md", "design/plan.md", "git/workspace.txt", "reads.txt", "check.json"} {
		if _, err := os.Stat(filepath.Join(out, "artifacts", "intent", "claude", filepath.FromSlash(rel))); err != nil {
			t.Fatalf("%s: %v", rel, err)
		}
	}
	for _, rel := range []string{"stream.jsonl", "stderr.txt", "snapshots"} {
		if _, err := os.Stat(filepath.Join(out, "artifacts", "intent", "claude", rel)); !os.IsNotExist(err) {
			t.Fatalf("%s was exported", rel)
		}
	}
	if _, err := os.Stat(filepath.Join(out, "artifacts", "overhead-empty", "claude", "reply.md")); err != nil {
		t.Fatalf("overhead run was not kept: %v", err)
	}

	var summary struct {
		SrcCommit     string `json:"src_commit"`
		Rows          []summaryRow
		Missing       []pairRef    `json:"missing"`
		HarnessErrors []harnessRef `json:"harness_errors"`
	}
	b, err := os.ReadFile(filepath.Join(out, "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &summary); err != nil {
		t.Fatal(err)
	}
	if summary.SrcCommit != "abc123" || len(summary.Rows) != 2 || len(summary.HarnessErrors) != 2 {
		t.Fatalf("summary %+v", summary)
	}
	if len(summary.Missing) != 1 || summary.Missing[0] != (pairRef{Scenario: "plan", Host: "codex"}) {
		t.Fatalf("missing %+v", summary.Missing)
	}
	// summary.json alone is a baseline LoadBaseline accepts.
	alone := filepath.Join(root, "summary-only")
	if err := os.MkdirAll(alone, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(alone, "summary.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	if rows, err := LoadBaseline(alone); err != nil || len(rows) != 2 || rows[0].ScenarioVersion != "1" {
		t.Fatalf("summary baseline %v %+v", err, rows)
	}

	hostsYAML, err := os.ReadFile(filepath.Join(out, "hosts.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"claude:", `host_version: "1.0.0"`, `host_version_source: "command"`, `model: "model-a"`, `effort: "medium"`, `src_commit: "abc123"`, `cli_version: "2.5.0"`, `fixture_digest: "sha256:00"`, "isolation:", `"2026-10-01"`, `"2026-10-02"`, `"intent": "1"`} {
		if !strings.Contains(string(hostsYAML), part) {
			t.Fatalf("hosts.yaml missing %q\n%s", part, hostsYAML)
		}
	}
	readme, err := os.ReadFile(filepath.Join(out, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"| `intent` | claude | fail |", "`plan` · codex", "## Notable behaviour"} {
		if !strings.Contains(string(readme), part) {
			t.Fatalf("README missing %q\n%s", part, readme)
		}
	}

	branchDir := writeRunDir(t, filepath.Join(root, "branch"), "20261010T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "pass", "2026-10-10", 1), note: "branch note\n", reads: "context/INDEX.md\n", reply: "branch reply\n"},
	})
	if err := Run(context.Background(), Options{Run: branchDir, Baseline: out, Design: true, NoZip: true, Stdout: &bytes.Buffer{}}); err != nil {
		t.Fatal(err)
	}
	design, err := os.ReadFile(filepath.Join(branchDir, "design.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(design)
	for _, part := range []string{"### `intent` · claude", "-new note", "+branch note", "files/sub/other.md", "-context/decisions.md", "newest reply", "branch reply"} {
		if !strings.Contains(text, part) {
			t.Fatalf("design.md missing %q\n%s", part, text)
		}
	}
	// The overhead run has no branch side; the intent pair must be paired.
	if strings.Contains(text, "only on branch") || strings.Contains(text, "before: reads list not available") {
		t.Fatalf("baseline run was not found\n%s", text)
	}
}

func TestExportBaselineRefusesAbsolutePath(t *testing.T) {
	root := t.TempDir()
	dir := writeRunDir(t, filepath.Join(root, "runs"), "20261001T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "pass", "2026-10-01", 1), note: "note\n", reads: "a\n", reply: "see /var/folders/xy/T/run/ws\n"},
	})
	out := filepath.Join(root, "baseline")
	err := ExportBaseline(ExportOptions{Runs: []string{dir}, Out: out, Stdout: &bytes.Buffer{}})
	if err == nil || !strings.Contains(err.Error(), "artifacts/intent/claude/reply.md") || !strings.Contains(err.Error(), "absolute path (/var/folders/)") {
		t.Fatalf("err %v", err)
	}
	noLeftovers(t, root)
}

func TestExportBaselineRedactsHomeBeforeChecking(t *testing.T) {
	root := t.TempDir()
	dir := writeRunDir(t, filepath.Join(root, "runs"), "20261001T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "pass", "2026-10-01", 1), note: "note\n", reads: "/home/sample-user/project/a.md\n", reply: "ok\n"},
	})
	out := filepath.Join(root, "baseline")
	if err := ExportBaseline(ExportOptions{Runs: []string{dir}, Out: out, Home: "/home/sample-user", User: "sample-user", Stdout: &bytes.Buffer{}}); err != nil {
		t.Fatal(err)
	}
	reads, err := os.ReadFile(filepath.Join(out, "artifacts", "intent", "claude", "reads.txt"))
	if err != nil || string(reads) != "$HOME/project/a.md\n" {
		t.Fatalf("reads %q %v", reads, err)
	}
}

func TestExportBaselineRefusesDenyTerm(t *testing.T) {
	for _, tc := range []struct {
		name  string
		item  fakeRun
		where string
	}{
		{"content", fakeRun{note: "Built for the ZephyrQuill team.\n", reads: "a\n", reply: "ok\n"}, "artifacts/intent/claude/files/note.md contains a deny-listed term"},
		{"file name", fakeRun{note: "n\n", reads: "a\n", reply: "ok\n", extra: map[string]string{"artifacts/files/zephyrquill-notes.md": "x\n"}}, "a file name under artifacts/intent/claude contains a deny-listed term"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			tc.item.run = baseRecord("intent", "claude", "pass", "2026-10-01", 1)
			dir := writeRunDir(t, filepath.Join(root, "runs"), "20261001T000000Z", []fakeRun{tc.item})
			deny := filepath.Join(root, "deny.txt")
			if err := os.WriteFile(deny, []byte("# local terms\n\nzephyrquill\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(root, "baseline")
			var buf bytes.Buffer
			err := ExportBaseline(ExportOptions{Runs: []string{dir}, Out: out, DenyFile: deny, Stdout: &buf})
			if err == nil || !strings.Contains(err.Error(), tc.where) {
				t.Fatalf("err %v", err)
			}
			if strings.Contains(strings.ToLower(err.Error()+buf.String()), "zephyrquill") {
				t.Fatal("the deny-listed term was printed")
			}
			noLeftovers(t, root)
		})
	}
}

func TestExportBaselineDenyFileUnreadable(t *testing.T) {
	root := t.TempDir()
	dir := writeRunDir(t, filepath.Join(root, "runs"), "20261001T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "pass", "2026-10-01", 1), note: "n\n", reads: "a\n", reply: "ok\n"},
	})
	missing := filepath.Join(root, "private-terms-file")
	err := ExportBaseline(ExportOptions{Runs: []string{dir}, Out: filepath.Join(root, "baseline"), DenyFile: missing})
	if err == nil || strings.Contains(err.Error(), "private-terms-file") {
		t.Fatalf("err %v", err)
	}
	noLeftovers(t, root)
}

func TestExportBaselineRefusesMixedSourceCommits(t *testing.T) {
	root := t.TempDir()
	other := baseRecord("plan", "codex", "pass", "2026-10-01", 1)
	other.SrcCommit = "def456"
	dir := writeRunDir(t, filepath.Join(root, "runs"), "20261001T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "pass", "2026-10-01", 1), note: "n\n", reads: "a\n", reply: "ok\n"},
		{run: other, note: "n\n", reads: "a\n", reply: "ok\n"},
	})
	err := ExportBaseline(ExportOptions{Runs: []string{dir}, Out: filepath.Join(root, "baseline")})
	if err == nil || !strings.Contains(err.Error(), "2 source commits") {
		t.Fatalf("err %v", err)
	}
	noLeftovers(t, root)
}

func TestExportBaselineRefusesHostDisagreement(t *testing.T) {
	root := t.TempDir()
	other := baseRecord("plan", "claude", "pass", "2026-10-01", 1)
	other.HostVersion = "1.0.1"
	dir := writeRunDir(t, filepath.Join(root, "runs"), "20261001T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "pass", "2026-10-01", 1), note: "n\n", reads: "a\n", reply: "ok\n"},
		{run: other, note: "n\n", reads: "a\n", reply: "ok\n"},
	})
	err := ExportBaseline(ExportOptions{Runs: []string{dir}, Out: filepath.Join(root, "baseline")})
	if err == nil || !strings.Contains(err.Error(), "host claude") || !strings.Contains(err.Error(), "host_version") {
		t.Fatalf("err %v", err)
	}
	noLeftovers(t, root)
}

func TestExportBaselineOutputTarget(t *testing.T) {
	root := t.TempDir()
	dir := writeRunDir(t, filepath.Join(root, "runs"), "20261001T000000Z", []fakeRun{
		{run: baseRecord("intent", "claude", "pass", "2026-10-01", 1), note: "n\n", reads: "a\n", reply: "ok\n"},
	})
	full := filepath.Join(root, "full")
	if err := os.MkdirAll(full, 0o755); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(full, "keep.txt")
	if err := os.WriteFile(keep, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ExportBaseline(ExportOptions{Runs: []string{dir}, Out: full}); err == nil || !strings.Contains(err.Error(), "not empty") {
		t.Fatalf("err %v", err)
	}
	if b, err := os.ReadFile(keep); err != nil || string(b) != "keep\n" {
		t.Fatal("existing file was touched")
	}
	empty := filepath.Join(root, "empty")
	if err := os.MkdirAll(empty, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ExportBaseline(ExportOptions{Runs: []string{dir}, Out: empty, Stdout: &bytes.Buffer{}}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(empty, "runs.jsonl")); err != nil {
		t.Fatal(err)
	}
}
