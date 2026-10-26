package report

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

func TestReportGolden(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "2.5.0")
	run := filepath.Join(dir, "run")
	if err := os.MkdirAll(filepath.Join(run, "surface"), 0o755); err != nil {
		t.Fatal(err)
	}
	surface := []byte(`{
  "base": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "files": [{"status": "changed", "path": "AGENTS.md"}, {"status": "added", "path": ".agents/skills/cc-plan/SKILL.md"}],
  "outputs": [{"name": "help", "changed": true}, {"name": "status", "changed": false}],
  "scenarios": ["plan", "intent"],
  "kinds": ["skill-text", "always-loaded"]
}
`)
	if err := os.WriteFile(filepath.Join(run, "surface", "surface.json"), surface, 0o644); err != nil {
		t.Fatal(err)
	}
	writeRuns(t, filepath.Join(base, "runs.jsonl"),
		line("intent", "claude", "base", "pass", 1, 198000, 3900, 0.57, nil),
		line("plan", "claude", "base", "pass", 1, 860000, 12000, 2.02, nil),
		line("execute", "claude", "base", "pass", 1, 100000, 1000, 0.25, nil),
		line("direct", "cursor", "base", "pass", 1, 0, 0, 0, nil),
		line("borrowed", "claude", "base", "inconclusive", 1, 80000, 2000, 0.20, nil),
		line("complete", "claude", "base", "pass", 1, 100000, 2000, 0.40, nil),
	)
	// direct's zeros above are placeholders; nil tokens are written next.
	fixNull(t, filepath.Join(base, "runs.jsonl"), "direct")
	writeRuns(t, filepath.Join(run, "runs.jsonl"),
		line("intent", "claude", "branch", "pass", 1, 212000, 4100, 0.61, nil),
		line("plan", "claude", "branch", "fail", 1, 880000, 12000, 2.10, []string{"plan-record", "nothing-pushed"}),
		line("execute", "claude", "branch", "pass", 1, 120000, 1000, 0.30, nil),
		line("direct", "cursor", "branch", "pass", 1, 0, 0, 0, nil),
		line("join", "codex", "branch", "harness-error", 1, 0, 0, 0, nil),
		line("borrowed", "claude", "branch", "inconclusive", 1, 80000, 2000, 0.20, nil),
		line("complete", "claude", "branch", "pass", 1, 100000, 2000, 0.40, nil),
		line("complete", "claude", "branch", "pass", 2, 120000, 4000, 0.60, nil),
	)
	fixNull(t, filepath.Join(run, "runs.jsonl"), "direct")
	fixNull(t, filepath.Join(run, "runs.jsonl"), "join")

	branch, err := LoadRuns(filepath.Join(run, "runs.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	before, err := LoadBaseline(base)
	if err != nil {
		t.Fatal(err)
	}
	gotSurface, has, err := loadSurface(run)
	if err != nil || !has {
		t.Fatal(err, has)
	}
	got := Render(Build(Input{
		Branch:       branch,
		Before:       before,
		BaselineName: "2.5.0",
		Surface:      gotSurface,
		HasSurface:   true,
	}))
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile("testdata/report.md", []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile("testdata/report.md")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Fatalf("report.md mismatch\n%s", got)
	}
}

func TestTokenWarnBoundary(t *testing.T) {
	if overTokenWarn(1000, 1150) {
		t.Fatal("15% is a warning")
	}
	if !overTokenWarn(1000, 1151) {
		t.Fatal("above 15% is not a warning")
	}
	if overTokenWarn(0, 10) {
		t.Fatal("zero before has a percentage")
	}
}

func TestNullTokenHasNoPercent(t *testing.T) {
	row := Row{
		Scenario: "direct",
		Host:     "cursor",
		Branch:   sideStats{Count: 1, InputMissing: true, Verdicts: map[string]int{"pass": 1}},
		Before:   sideStats{Count: 1, InputMissing: true, Verdicts: map[string]int{"pass": 1}},
	}
	got := inputCell(row)
	if got != "n/a / n/a" {
		t.Fatalf("%q", got)
	}
}

func TestOverheadRow(t *testing.T) {
	emptyIn, worldIn := int64(1000), int64(2500)
	beforeEmpty, beforeWorld := int64(900), int64(2000)
	in := Input{
		Branch: []record.Run{
			{Scenario: "overhead-empty", Host: "claude", Side: "branch", Verdict: "pass", Run: 1, Tokens: record.Tokens{InputTotal: &emptyIn}},
			{Scenario: "overhead-world", Host: "claude", Side: "branch", Verdict: "pass", Run: 1, Tokens: record.Tokens{InputTotal: &worldIn}, CostBasis: "reported"},
		},
		Before: []record.Run{
			{Scenario: "overhead-empty", Host: "claude", Side: "base", Verdict: "pass", Run: 1, Tokens: record.Tokens{InputTotal: &beforeEmpty}},
			{Scenario: "overhead-world", Host: "claude", Side: "base", Verdict: "pass", Run: 1, Tokens: record.Tokens{InputTotal: &beforeWorld}, CostBasis: "reported"},
		},
	}
	m := Build(in)
	if len(m.Rows) != 1 || m.Rows[0].Scenario != "overhead (S0)" {
		t.Fatalf("%+v", m.Rows)
	}
	got := inputCell(m.Rows[0])
	if got != "ws-empty: 1.5k / 1.1k / +36%" {
		t.Fatalf("%q", got)
	}
	if branchCell(m.Rows[0]) != "—" || outputCell(m.Rows[0]) != "—" {
		t.Fatalf("measurement cells %s %s", branchCell(m.Rows[0]), outputCell(m.Rows[0]))
	}
}

func TestSummaryBaseline(t *testing.T) {
	dir := t.TempDir()
	body := []byte(`{"rows":[{"scenario":"intent","host":"claude","verdict":"pass","input_total":1000,"output":10,"cost_usd":0.1,"cost_basis":"reported","runs":2}]}`)
	if err := os.WriteFile(filepath.Join(dir, "summary.json"), body, 0o644); err != nil {
		t.Fatal(err)
	}
	runs, err := LoadBaseline(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 2 || runs[0].Side != "base" || runs[1].Run != 2 {
		t.Fatalf("%+v", runs)
	}
	m := Build(Input{Before: runs, BaselineName: "summary"})
	if len(m.Rows) != 1 || m.Rows[0].Before.Count != 2 {
		t.Fatalf("%+v", m.Rows)
	}
}

func line(scenario, host, side, verdict string, n int, in, out int64, cost float64, failed []string) record.Run {
	return record.Run{
		Scenario:     scenario,
		Host:         host,
		Side:         side,
		Verdict:      verdict,
		Run:          n,
		HostVersion:  hostVersion(host),
		Model:        hostModel(host),
		Effort:       hostEffort(host),
		Tokens:       record.Tokens{InputTotal: &in, Output: &out},
		CostUSD:      &cost,
		CostBasis:    "reported",
		FailedChecks: failed,
	}
}

func hostVersion(host string) string {
	switch host {
	case "codex":
		return "0.162.1"
	case "cursor":
		return "2026.10.22-27dcf5f"
	default:
		return "2.1.304"
	}
}

func hostModel(host string) string {
	switch host {
	case "codex":
		return "gpt-6-sol"
	case "cursor":
		return "auto"
	default:
		return "claude-opus-5-5"
	}
}

func hostEffort(host string) string {
	if host == "codex" {
		return "medium"
	}
	return ""
}

func writeRuns(t *testing.T, path string, runs ...record.Run) {
	t.Helper()
	for _, run := range runs {
		if err := record.Append(path, run); err != nil {
			t.Fatal(err)
		}
	}
}

func fixNull(t *testing.T, path, scenario string) {
	t.Helper()
	runs, err := LoadRuns(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	for _, run := range runs {
		if run.Scenario == scenario {
			run.Tokens = record.Tokens{}
			run.CostUSD = nil
			run.CostBasis = "none"
		}
		if err := record.Append(path, run); err != nil {
			t.Fatal(err)
		}
	}
}
