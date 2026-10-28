package runner

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/report"
)

// TestExportBaselineFromStubRuns exports real stub-host run directories and
// reads them back as the before side of report --design.
func TestExportBaselineFromStubRuns(t *testing.T) {
	home := os.Getenv("HOME")
	_, out := runStubs(t, home, "")
	logs, err := filepath.Glob(filepath.Join(out, "*", "runs.jsonl"))
	if err != nil || len(logs) == 0 {
		t.Fatalf("run logs %v %v", logs, err)
	}
	var dirs []string
	for _, log := range logs {
		dirs = append(dirs, filepath.Dir(log))
	}
	exported := filepath.Join(tempOutsideHome(t), "baseline")
	var buf bytes.Buffer
	if err := report.ExportBaseline(report.ExportOptions{Runs: dirs, Out: exported, Home: home, User: os.Getenv("USER"), Stdout: &buf}); err != nil {
		t.Fatalf("export: %v\n%s", err, buf.String())
	}
	hostsYAML, err := os.ReadFile(filepath.Join(exported, "hosts.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(hostsYAML), `host_version: "stub-9"`) != 3 {
		t.Fatalf("hosts.yaml does not record the reported version\n%s", hostsYAML)
	}
	var design strings.Builder
	for _, dir := range dirs {
		if err := report.Run(context.Background(), report.Options{Run: dir, Baseline: exported, Design: true, NoZip: true, Stdout: &bytes.Buffer{}}); err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(filepath.Join(dir, "design.md"))
		if err != nil {
			t.Fatal(err)
		}
		design.Write(b)
	}
	text := design.String()
	for _, host := range []string{"claude", "codex", "cursor"} {
		head := "### `marker` · " + host + " · version 1"
		// Each run directory holds some hosts; a host's pair is in the
		// design.md of the directory that ran it.
		section, ok := "", false
		for _, part := range strings.Split(text, "\n### `") {
			if strings.HasPrefix("### `"+part, head) && !strings.Contains(part, "only on") {
				section, ok = part, true
				break
			}
		}
		if !ok {
			t.Fatalf("no pair with a baseline run for %s\n%s", host, text)
		}
		if !strings.Contains(section, "`files/marker.txt`") {
			t.Fatalf("%s produced files missing\n%s", host, section)
		}
		if strings.Count(section, "(none)") != 0 {
			t.Fatalf("%s reply missing\n%s", host, section)
		}
		// The stub streams carry no tool call, so each side's reads list is
		// empty or not available; the exported file is what design read.
		reads, err := os.ReadFile(filepath.Join(exported, "artifacts", "marker", host, "reads.txt"))
		if err != nil || (host == "cursor" && !strings.Contains(string(reads), "reads list not available")) {
			t.Fatalf("%s exported reads %q %v", host, reads, err)
		}
		if !strings.Contains(section, "### Reads\n\nreads list not available") {
			t.Fatalf("%s reads\n%s", host, section)
		}
	}
}
