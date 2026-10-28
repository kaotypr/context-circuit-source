package main

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/surface"
)

func TestSurfaceDiff(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go missing")
	}
	repo := sourceRoot(t)
	before := gitPorcelain(t, repo)
	t.Cleanup(func() {
		after := gitPorcelain(t, repo)
		if before != after {
			t.Errorf("source checkout status changed\nbefore:\n%s\nafter:\n%s", before, after)
		}
	})

	parent := t.TempDir()
	cfg := filepath.Join(parent, "gitconfig")
	if err := os.WriteFile(cfg, []byte("[user]\n\tname = Agent Test\n\temail = agent-test@example.invalid\n[commit]\n\tgpgsign = false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	clone := filepath.Join(parent, "clone")
	gitAt(t, parent, cfg, "clone", "--local", repo, clone)
	triggers, err := os.ReadFile(filepath.Join(repo, "agent-tests", "triggers.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(clone, "agent-tests", "triggers.yaml"), triggers, 0o644); err != nil {
		t.Fatal(err)
	}
	base := strings.TrimSpace(gitAt(t, clone, cfg, "rev-parse", "HEAD"))
	caller := map[string]string{"PATH": os.Getenv("PATH"), "HOME": os.Getenv("HOME")}

	t.Run("refactor", func(t *testing.T) {
		reset(t, clone, cfg, base)
		appendFile(t, filepath.Join(clone, "internal", "cli", "topics.go"), "\n// Renamed nothing.\n")
		commit(t, clone, cfg, "internal/cli/topics.go", "test: comment a help topic file")
		report := runSurface(t, clone, caller)
		if !report.Empty || report.Message != "No agent tests needed: surface diff is empty" || report.Larger || len(report.Scenarios) != 0 {
			t.Fatalf("%s\n%s", report.Message, readReport(t, report.Dir))
		}
	})

	t.Run("skill", func(t *testing.T) {
		reset(t, clone, cfg, base)
		appendFile(t, filepath.Join(clone, "product", "skills", "cc-intent", "SKILL.md"), "\nContributors read this extra sentence.\n")
		commit(t, clone, cfg, "product/skills/cc-intent/SKILL.md", "test: extend intent skill text")
		report := runSurface(t, clone, caller)
		if report.Empty || report.Larger || !slices.Equal(report.Scenarios, []string{"intent"}) || len(report.Hosts) != 0 || !slices.Contains(report.Kinds, "skill-text") {
			t.Fatalf("%s\nscenarios %v hosts %v kinds %v larger %v\n%s", report.Message, report.Scenarios, report.Hosts, report.Kinds, report.Larger, readReport(t, report.Dir))
		}
	})

	t.Run("agents", func(t *testing.T) {
		reset(t, clone, cfg, base)
		path := filepath.Join(clone, "product", "AGENTS.md.in")
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		old := "## Language\n"
		if !strings.Contains(string(body), old) {
			t.Fatal("language heading missing")
		}
		next := strings.Replace(string(body), old, old+"\nAn extra sentence sits in this section.\n", 1)
		if err := os.WriteFile(path, []byte(next), 0o644); err != nil {
			t.Fatal(err)
		}
		commit(t, clone, cfg, "product/AGENTS.md.in", "test: extend the language section")
		report := runSurface(t, clone, caller)
		if !slices.Contains(report.Scenarios, "orient") || !slices.Contains(report.Scenarios, "intent-lang") {
			t.Fatalf("scenarios %v\n%s", report.Scenarios, readReport(t, report.Dir))
		}
	})

	t.Run("route", func(t *testing.T) {
		reset(t, clone, cfg, base)
		appendFile(t, filepath.Join(clone, "product", "skills", "cc-plan", "route.md"), "\nOpen the canonical skill again.\n")
		commit(t, clone, cfg, "product/skills/cc-plan/route.md", "test: extend the plan route")
		report := runSurface(t, clone, caller)
		if !slices.Equal(report.Hosts, []string{"claude", "cursor"}) || !slices.Contains(report.Scenarios, "plan") {
			t.Fatalf("hosts %v scenarios %v\n%s", report.Hosts, report.Scenarios, readReport(t, report.Dir))
		}
	})

	t.Run("manifest", func(t *testing.T) {
		reset(t, clone, cfg, base)
		path := filepath.Join(clone, "scripts", "release-manifest.txt")
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		old := ".context-circuit/docs/how-it-works.md"
		next := strings.Replace(string(body), old, ".context-circuit/docs/how-it-works-renamed.md", 1)
		if next == string(body) {
			t.Fatal("destination missing")
		}
		if err := os.WriteFile(path, []byte(next), 0o644); err != nil {
			t.Fatal(err)
		}
		commit(t, clone, cfg, "scripts/release-manifest.txt", "test: rename a shipped doc destination")
		report := runSurface(t, clone, caller)
		if !report.Larger || !slices.Equal(report.Scenarios, scenarioList(t, filepath.Join(clone, "agent-tests", "scenarios"))) {
			t.Fatalf("larger %v scenarios %v\n%s", report.Larger, report.Scenarios, readReport(t, report.Dir))
		}
	})
}

func runSurface(t *testing.T, clone string, caller map[string]string) surface.Report {
	t.Helper()
	out := t.TempDir()
	stdout, stderr := os.Stdout, os.Stderr
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout, os.Stderr = wOut, wErr
	code := surfaceDiffCmd([]string{"--base", "HEAD~1", "--ref", "HEAD", "--checkout", clone, "--out", out}, caller)
	wOut.Close()
	wErr.Close()
	os.Stdout, os.Stderr = stdout, stderr
	printed, _ := io.ReadAll(rOut)
	errs, _ := io.ReadAll(rErr)
	if code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, errs, printed)
	}
	var found string
	err = filepath.WalkDir(out, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == "surface.json" {
			found = p
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if found == "" {
		t.Fatalf("no surface.json\n%s\n%s", errs, printed)
	}
	raw, err := os.ReadFile(found)
	if err != nil {
		t.Fatal(err)
	}
	var report surface.Report
	if err := json.Unmarshal(raw, &report); err != nil {
		t.Fatal(err)
	}
	return report
}

func readReport(t *testing.T, dir string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, "surface.md"))
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func reset(t *testing.T, clone, cfg, base string) {
	t.Helper()
	gitAt(t, clone, cfg, "checkout", "-B", "surface-case", base)
}

func commit(t *testing.T, clone, cfg, rel, message string) {
	t.Helper()
	gitAt(t, clone, cfg, "add", "--", rel)
	gitAt(t, clone, cfg, "commit", "-m", message)
}

func appendFile(t *testing.T, path, extra string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(extra); err != nil {
		t.Fatal(err)
	}
}

func scenarioList(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") {
			out = append(out, strings.TrimSuffix(name, ".yaml"))
		}
	}
	slices.Sort(out)
	return out
}

func sourceRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "cmd", "context-circuit", "main.go")); err != nil {
		t.Fatal(err)
	}
	return root
}

func gitPorcelain(t *testing.T, root string) string {
	t.Helper()
	cmd := exec.Command("git", "-C", root, "status", "--porcelain")
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("status: %v\n%s", err, out)
	}
	return string(out)
}

func gitAt(t *testing.T, dir, cfg string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=" + cfg,
		"GIT_AUTHOR_NAME=Agent Test",
		"GIT_AUTHOR_EMAIL=agent-test@example.invalid",
		"GIT_COMMITTER_NAME=Agent Test",
		"GIT_COMMITTER_EMAIL=agent-test@example.invalid",
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}
