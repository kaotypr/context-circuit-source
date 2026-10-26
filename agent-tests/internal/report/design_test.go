package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

func TestDesignGolden(t *testing.T) {
	dir := t.TempDir()
	beforeNote := "alpha\nbeta\ngamma\n"
	afterNote := "alpha\nbeta\ndelta\n"
	beforePlan := "# Request\n\nCustomers need a CSV download.\n"
	afterPlan := "# Request\n\nCustomers need a CSV download.\n\n# Open questions\n\nNone known.\n"
	writePair(t, dir, "base", beforeNote, beforePlan, "context/INDEX.md\n", "The before reply.\nIt asks for approval.\n")
	writePair(t, dir, "branch", afterNote, afterPlan, "context/INDEX.md\ncontext/domains/invoices.md\n", "The after reply.\nIt still asks for approval.\n")
	inTok, outTok := int64(198000), int64(3900)
	brIn, brOut := int64(212000), int64(4100)
	before := record.Run{Scenario: "intent", Host: "claude", Side: "base", Verdict: "pass", Run: 1, Tokens: record.Tokens{InputTotal: &inTok, Output: &outTok}}
	branch := record.Run{Scenario: "intent", Host: "claude", Side: "branch", Verdict: "pass", Run: 1, Tokens: record.Tokens{InputTotal: &brIn, Output: &brOut}}
	got, err := RenderDesign(DesignInput{RunDir: dir, Branch: []record.Run{branch}, Before: []record.Run{before}})
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile("testdata/design.md", []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile("testdata/design.md")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Fatalf("design.md mismatch\n%s", got)
	}
	if strings.Contains(got, "reads list not available") {
		t.Fatal("reads were present")
	}
}

func TestPairKeyIncludesVersionAndOneSide(t *testing.T) {
	base := record.Run{Scenario: "intent", Host: "claude", ScenarioVersion: "1", Side: "base", Verdict: "pass", Run: 1}
	branch := record.Run{Scenario: "intent", Host: "claude", ScenarioVersion: "2", Side: "branch", Verdict: "pass", Run: 1}
	only := record.Run{Scenario: "plan", Host: "codex", ScenarioVersion: "1", Side: "branch", Verdict: "pass", Run: 1}
	got, err := RenderDesign(DesignInput{Branch: []record.Run{branch, only}, Before: []record.Run{base}})
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"only on base", "only on branch", "version 1", "version 2"} {
		if !strings.Contains(got, part) {
			t.Fatalf("missing %q\n%s", part, got)
		}
	}
	if strings.Contains(got, "Pair:") {
		t.Fatalf("different versions were paired\n%s", got)
	}
}

func TestReadsUnavailable(t *testing.T) {
	dir := t.TempDir()
	for _, side := range []string{"base", "branch"} {
		root := filepath.Join(dir, "runs", "orient", "claude", side, "r1")
		if err := os.MkdirAll(root, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "reads.txt"), []byte("reads list not available\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run := record.Run{Scenario: "orient", Host: "claude", Verdict: "pass", Run: 1}
	base := run
	base.Side = "base"
	branch := run
	branch.Side = "branch"
	got, err := RenderDesign(DesignInput{RunDir: dir, Branch: []record.Run{branch}, Before: []record.Run{base}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "reads list not available") {
		t.Fatalf("%s", got)
	}
}

func writePair(t *testing.T, dir, side, note, plan, reads, reply string) {
	t.Helper()
	root := filepath.Join(dir, "runs", "intent", "claude", side, "r1")
	files := filepath.Join(root, "artifacts", "files")
	design := filepath.Join(root, "artifacts", "design", "intent")
	if err := os.MkdirAll(files, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(design, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(files, "note.txt"), []byte(note), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(design, "i100-demo.md"), []byte(plan), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "reads.txt"), []byte(reads), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "reply.md"), []byte(reply), 0o644); err != nil {
		t.Fatal(err)
	}
}
