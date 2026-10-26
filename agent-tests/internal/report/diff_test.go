package report

import (
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

func TestUnified(t *testing.T) {
	got := unified("a\nb\nc\n", "a\nb\nd\n", "files/note.txt")
	want := "--- before/files/note.txt\n+++ after/files/note.txt\n@@ -1,3 +1,3 @@\n a\n b\n-c\n+d\n"
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUnifiedNewFile(t *testing.T) {
	got := unified("", "one\n", "files/new.txt")
	want := "--- before/files/new.txt\n+++ after/files/new.txt\n@@ -0,0 +1,1 @@\n+one\n"
	if got != want {
		t.Fatalf("got:\n%s", got)
	}
}

func TestPickRunMedian(t *testing.T) {
	runs := []record.Run{
		passRun(1, 30),
		passRun(2, 10),
		passRun(3, 20),
		passRun(4, 40),
	}
	// Passing runs sorted by tokens: 10, 20, 30, 40. Lower median is 20, run 3.
	got, how, ok := pickRun(runs)
	if !ok || got.Run != 3 || *got.Tokens.InputTotal != 20 {
		t.Fatalf("got r%d tokens %v ok %v", got.Run, got.Tokens.InputTotal, ok)
	}
	if how != "the passing run with the median total input tokens" {
		t.Fatalf("how %q", how)
	}
	tie := []record.Run{passRun(2, 5), passRun(1, 5)}
	got, _, ok = pickRun(tie)
	if !ok || got.Run != 1 {
		t.Fatalf("tie r%d", got.Run)
	}
	failed := []record.Run{
		{Scenario: "plan", Host: "claude", Verdict: "fail", Run: 1, Tokens: record.Tokens{InputTotal: intPtr(1)}},
		{Scenario: "plan", Host: "claude", Verdict: "fail", Run: 2, Tokens: record.Tokens{InputTotal: intPtr(9)}},
		{Scenario: "plan", Host: "claude", Verdict: "fail", Run: 3, Tokens: record.Tokens{InputTotal: intPtr(5)}},
	}
	got, how, ok = pickRun(failed)
	if !ok || got.Run != 3 || how != "the median run; none passed" {
		t.Fatalf("r%d %s", got.Run, how)
	}
}

func passRun(n int, tokens int64) record.Run {
	return record.Run{
		Scenario: "intent",
		Host:     "claude",
		Verdict:  "pass",
		Run:      n,
		Tokens:   record.Tokens{InputTotal: intPtr(tokens)},
	}
}

func intPtr(n int64) *int64 { return &n }
