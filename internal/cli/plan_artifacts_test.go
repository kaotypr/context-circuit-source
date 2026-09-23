package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStandaloneFolderPlanAndSelectedOrder(t *testing.T) {
	f := setup(t)
	f.repository("api")
	first := record(t, f.ok("record", "create", "--kind", "plan", "--slug", "first", "--title", "First", "--repo", "api"))
	second := record(t, f.ok("record", "create", "--kind", "plan", "--slug", "second", "--title", "Second", "--repo", "api", "--depends-on", first.ID))
	if first.Intent != "" || !strings.HasSuffix(first.Path, "/plan.md") {
		t.Fatal(first)
	}
	if strings.Contains(first.Content, "intent:") {
		t.Fatal("standalone plan wrote an intent field")
	}
	f.ok("record", "show", "--id", first.ID)
	result := order(t, f.ok("record", "order", "--plan", second.ID))
	if len(result.PlanIDs) != 1 || result.PlanIDs[0] != second.ID || result.Selection != "plans" || len(result.Blocked) != 1 {
		t.Fatal(result)
	}
	f.fail("record", "order", "--plan", second.ID, "--intent", "i001")
	f.ok("record", "complete", "--id", first.ID, "--text", "Completed")
	result = order(t, f.ok("record", "order", "--plan", second.ID))
	if len(result.Blocked) != 0 || result.Remaining != 1 {
		t.Fatal(result)
	}
	f.ok("check")
	planner := dispatch(t, f, "agent", "dispatch", "--host", "codex", "--role", "planner", "--task", "Plan the specified API outcome.", "--path", f.root)
	if !strings.Contains(planner.Prompt, "Plan the specified API outcome.") || strings.Contains(planner.Prompt, "# Intent ") {
		t.Fatal(planner.Prompt)
	}
}

func TestPlanDetailsDispatchAndDiagnostics(t *testing.T) {
	f := setup(t)
	f.repository("api")
	f.repository("web")
	p := record(t, f.ok("record", "create", "--kind", "plan", "--slug", "shared", "--title", "Shared", "--repo", "api", "--repo", "web"))
	base := filepath.Dir(filepath.Join(f.root, p.Path))
	write(t, filepath.Join(base, "model.md"), "# Data model\n")
	write(t, filepath.Join(base, "api.md"), "# API\n")
	write(t, filepath.Join(base, "web.md"), "# Web\n")
	entry := read(t, filepath.Join(f.root, p.Path))
	entry = strings.Replace(entry, "---\n\n# Shared", "required_files:\n  shared:\n    - model.md\n  by_repository:\n    api:\n      - api.md\n    web:\n      - web.md\n---\n\n# Shared", 1)
	entry += "\n[Model](model.md)\n[API](api.md)\n[Web](web.md)\n"
	write(t, filepath.Join(f.root, p.Path), entry)
	f.ok("check")
	f.fail("agent", "dispatch", "--host", "codex", "--role", "worker", "--plan", p.ID, "--path", f.root)
	api := dispatch(t, f, "agent", "dispatch", "--host", "codex", "--role", "worker", "--plan", p.ID, "--repo", "api", "--path", f.root)
	if !strings.Contains(api.Prompt, "plans/"+filepath.Base(base)+"/model.md") || !strings.Contains(api.Prompt, "/api.md") || strings.Contains(api.Prompt, "- plans/"+filepath.Base(base)+"/web.md") {
		t.Fatal(api.Prompt)
	}
	if strings.Contains(api.Prompt, "by_repository:") {
		t.Fatal("frontmatter leaked into brief")
	}
	write(t, filepath.Join(base, "extra.md"), "# Unassigned\n")
	f.fail("check")
	if err := os.Remove(filepath.Join(base, "extra.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(base, "api.md")); err != nil {
		t.Fatal(err)
	}
	f.fail("check")
	f.fail("agent", "dispatch", "--host", "codex", "--role", "worker", "--plan", p.ID, "--repo", "api", "--path", f.root)
}

func TestLegacyAndFolderPlanLookup(t *testing.T) {
	f := setup(t)
	f.repository("api")
	legacy := record(t, f.ok("record", "create", "--kind", "plan", "--slug", "old", "--title", "Old", "--repo", "api"))
	oldEntry := filepath.Join(f.root, legacy.Path)
	oldFile := filepath.Join(f.root, "plans", legacy.ID+"-old.md")
	if err := os.Rename(oldEntry, oldFile); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Dir(oldEntry)); err != nil {
		t.Fatal(err)
	}
	show := record(t, f.ok("record", "show", "--id", legacy.ID))
	if show.Path != "plans/"+legacy.ID+"-old.md" {
		t.Fatal(show.Path)
	}
	f.ok("record", "complete", "--id", legacy.ID, "--text", "Done")
	current := record(t, f.ok("record", "create", "--kind", "plan", "--slug", "new", "--title", "New", "--repo", "api", "--depends-on", legacy.ID))
	result := order(t, f.ok("record", "order", "--plan", legacy.ID, "--plan", current.ID))
	if len(result.PlanIDs) != 2 || result.Remaining != 1 {
		t.Fatal(result)
	}
	f.ok("check")
	duplicate := filepath.Join(f.root, "plans", legacy.ID+"-duplicate", "plan.md")
	write(t, duplicate, read(t, oldFile))
	f.fail("record", "show", "--id", legacy.ID)
	f.fail("check")
}

func TestPlanRequiredPathSafetyAndMissingEntry(t *testing.T) {
	f := setup(t)
	f.repository("api")
	p := record(t, f.ok("record", "create", "--kind", "plan", "--slug", "safe", "--title", "Safe", "--repo", "api"))
	entryPath := filepath.Join(f.root, p.Path)
	entry := read(t, entryPath)
	unsafe := strings.Replace(entry, "---\n\n# Safe", "required_files:\n  shared: [../outside.md]\n---\n\n# Safe", 1)
	write(t, entryPath, unsafe)
	f.fail("check")
	f.fail("agent", "dispatch", "--host", "codex", "--role", "worker", "--plan", p.ID, "--path", f.root)
	write(t, entryPath, entry)
	if err := os.Remove(entryPath); err != nil {
		t.Fatal(err)
	}
	if out := f.fail("check"); !strings.Contains(out, "missing or invalid plan entry") {
		t.Fatal(out)
	}
}
