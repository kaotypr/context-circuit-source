package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSubscriptionAuthFlag(t *testing.T) {
	if code := runCmd([]string{"--subscription-auth", "cursor"}, map[string]string{}); code != 2 {
		t.Fatalf("run %d", code)
	}
	if code := probeCmd([]string{"--subscription-auth", "claude"}, map[string]string{}); code != 2 {
		t.Fatalf("probe %d", code)
	}
	names, err := subscriptionNames("codex")
	if err != nil || len(names) != 1 || names[0] != "codex" {
		t.Fatalf("%v %v", names, err)
	}
	names, err = subscriptionNames("")
	if err != nil || names != nil {
		t.Fatalf("%v %v", names, err)
	}
	if _, err := subscriptionNames("codex,cursor"); err == nil || !strings.Contains(err.Error(), "codex") {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	usage(&buf)
	if !strings.Contains(buf.String(), "--subscription-auth") {
		t.Fatalf("%s", buf.String())
	}
}

func TestUsageMentionsSurfaceDiff(t *testing.T) {
	var buf bytes.Buffer
	usage(&buf)
	text := buf.String()
	for _, part := range []string{"surface-diff", "--base", "--ref", "--checkout", "--out", "does not start a host", "needs no key"} {
		if !strings.Contains(text, part) {
			t.Errorf("usage missing %q\n%s", part, text)
		}
	}
}

func TestReportUsage(t *testing.T) {
	var buf bytes.Buffer
	usage(&buf)
	text := buf.String()
	for _, part := range []string{"report", "--run", "--baseline", "--design", "--upgrade-check", "--no-zip", "does not start a host"} {
		if !strings.Contains(text, part) {
			t.Errorf("usage missing %q\n%s", part, text)
		}
	}
	if code := reportCmd([]string{"--nope"}, map[string]string{}); code != 2 {
		t.Fatalf("unknown flag %d", code)
	}
	if code := reportCmd([]string{"extra"}, map[string]string{"PATH": os.Getenv("PATH")}); code != 2 {
		t.Fatalf("extra arg %d", code)
	}
}

func TestReportCommandWritesMarkdown(t *testing.T) {
	out := t.TempDir()
	run := filepath.Join(out, "20260115T090000Z")
	if err := os.MkdirAll(run, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `{"scenario":"intent","host":"claude","side":"branch","verdict":"pass","run":1,"host_version":"2","model":"m","tokens":{"input_total":1000,"input_uncached":null,"cache_write":null,"cache_read":null,"output":10,"reasoning":null},"cost_usd":0.1,"cost_basis":"reported","failed_checks":[]}
`
	if err := os.WriteFile(filepath.Join(run, "runs.jsonl"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	code := reportCmd([]string{"--out", out, "--run", "20260115T090000Z", "--no-zip"}, map[string]string{})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	got, err := os.ReadFile(filepath.Join(run, "report.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "## Agent tests") || !strings.Contains(string(got), "intent") {
		t.Fatalf("%s", got)
	}
}

func TestSurfaceDiffUsageErrors(t *testing.T) {
	if code := surfaceDiffCmd(nil, map[string]string{}); code != 2 {
		t.Fatalf("missing base %d", code)
	}
	if code := surfaceDiffCmd([]string{"--nope"}, map[string]string{}); code != 2 {
		t.Fatalf("unknown flag %d", code)
	}
	if code := surfaceDiffCmd([]string{"--base", "HEAD", "extra"}, map[string]string{}); code != 2 {
		t.Fatalf("extra arg %d", code)
	}
}

func TestExportBaselineFlags(t *testing.T) {
	var buf bytes.Buffer
	usage(&buf)
	for _, part := range []string{"export-baseline", "--deny-file", "--redact-deny", "--run", "--out"} {
		if !strings.Contains(buf.String(), part) {
			t.Errorf("usage missing %q", part)
		}
	}
	if code := exportBaselineCmd([]string{"--out", t.TempDir()}, map[string]string{}); code != 2 {
		t.Fatalf("missing --run %d", code)
	}
	if code := exportBaselineCmd([]string{"--run", t.TempDir()}, map[string]string{}); code != 2 {
		t.Fatalf("missing --out %d", code)
	}
	if code := exportBaselineCmd([]string{"--run", t.TempDir(), "--out", t.TempDir(), "--redact-deny"}, map[string]string{}); code != 2 {
		t.Fatalf("--redact-deny without --deny-file %d", code)
	}
	out := filepath.Join(t.TempDir(), "baseline")
	if code := exportBaselineCmd([]string{"--run", t.TempDir(), "--out", out}, map[string]string{}); code != 1 {
		t.Fatalf("run without runs.jsonl %d", code)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatal("a failed export left its output directory")
	}
}
