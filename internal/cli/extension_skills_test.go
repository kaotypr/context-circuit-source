package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	contextcircuit "github.com/kaotypr/context-circuit-source"
	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

const qaManifest = `id: qa-pass
version: 1.0.0
name: QA pass
summary: Runs a QA pass on a plan's worktree.
members: default-on
hooks: []
skills:
  - id: qa
    path: skills/qa
    stage: qa
`

const qaSkill = `---
name: qa
description: Run the QA pass on a plan's worktree when a person asks for it.
---

# QA pass

Run ` + "`extension hooks --stage qa-pass/qa --when before`" + ` first.
`

// qaExtension writes the qa-pass extension, whose skill declares stage qa,
// with a reference file and a script beside its SKILL.md.
func qaExtension(t *testing.T, root string) {
	t.Helper()
	dir := filepath.Join(root, "extensions", "qa-pass")
	write(t, filepath.Join(dir, "extension.yaml"), qaManifest)
	write(t, filepath.Join(dir, "skills", "qa", "SKILL.md"), qaSkill)
	write(t, filepath.Join(dir, "skills", "qa", "references", "checklist.md"), "# Checklist\n\n- Run the tests.\n")
}

// A second extension hooks the named stage.
const reportManifest = `id: qa-report
members: default-on
hooks:
  - at: after qa-pass/qa
    instructions: instructions/after-qa.md
`

func TestExtensionSkillsAreInstalledRefreshedAndRemoved(t *testing.T) {
	f := setup(t)
	qaExtension(t, f.root)
	added := decodeMap(t, f.ok("extension", "add", "--id", "qa-pass", "--path", "extensions/qa-pass"))
	if note, _ := added["setup_required"].(string); !strings.Contains(note, "agent setup") || !strings.Contains(note, "skills") {
		t.Fatalf("add did not name agent setup for the skill: %v", added)
	}
	status := decodeMap(t, f.ok("status"))
	if skills, _ := status["extension_skills"].([]any); len(skills) != 1 {
		t.Fatalf("status does not list the skill: %v", status)
	} else if skill := skills[0].(map[string]any); skill["name"] != "x-qa-pass-qa" || skill["installed"] != false || skill["stage"] != "qa-pass/qa" {
		t.Fatalf("status skill: %v", skill)
	}
	if stages, _ := status["named_stages"].([]any); len(stages) != 1 || stages[0] != "qa-pass/qa" {
		t.Fatalf("status named stages: %v", status["named_stages"])
	}

	var result workspace.AgentSetup
	if err := json.Unmarshal([]byte(f.ok("agent", "setup", "--host", "claude-code")), &result); err != nil {
		t.Fatal(err)
	}
	// Naming one host still copies the shared folder its wrapper points at.
	copied := filepath.Join(f.root, ".agents", "skills", "x-qa-pass-qa")
	skill := read(t, filepath.Join(copied, "SKILL.md"))
	if !strings.HasPrefix(skill, "---\nname: x-qa-pass-qa\ndescription: Run the QA pass") || !strings.Contains(skill, "# QA pass") {
		t.Fatalf("the copy's frontmatter was not renamed, or its body changed:\n%s", skill)
	}
	if read(t, filepath.Join(copied, "references", "checklist.md")) != "# Checklist\n\n- Run the tests.\n" {
		t.Fatal("a reference file was not copied")
	}
	wrapper := read(t, filepath.Join(f.root, ".claude", "skills", "x-qa-pass-qa", "SKILL.md"))
	if wrapper != "---\nname: x-qa-pass-qa\ndescription: Run the QA pass on a plan's worktree when a person asks for it.\n---\n\nRead and follow the canonical workspace skill at\n`.agents/skills/x-qa-pass-qa/SKILL.md`.\n" {
		t.Fatalf("the wrapper is not in the shipped form:\n%s", wrapper)
	}
	if _, err := os.Stat(filepath.Join(f.root, ".cursor", "skills", "x-qa-pass-qa")); !os.IsNotExist(err) {
		t.Fatal("setup for claude-code wrote cursor's wrapper")
	}
	f.ok("agent", "setup")
	read(t, filepath.Join(f.root, ".cursor", "skills", "x-qa-pass-qa", "SKILL.md"))
	if _, err := os.Stat(filepath.Join(f.root, ".codex", "skills")); !os.IsNotExist(err) {
		t.Fatal("codex was given a wrapper; it reads .agents/skills itself")
	}
	if _, ok := decodeMap(t, f.ok("status"))["setup_required"]; ok {
		t.Fatal("status still reports setup after it ran")
	}
	// A shipped skill is never touched.
	shipped := read(t, filepath.Join(f.root, ".agents", "skills", "cc-plan", "SKILL.md"))

	// The installed paths are ignored, so a commit never carries them.
	git(t, f.root, "init", "-b", "main")
	for _, path := range []string{".agents/skills/x-qa-pass-qa/SKILL.md", ".agents/skills/x-qa-pass-qa/references/checklist.md", ".claude/skills/x-qa-pass-qa/SKILL.md", ".cursor/skills/x-qa-pass-qa/SKILL.md"} {
		git(t, f.root, "check-ignore", path)
	}
	if out := git(t, f.root, "status", "--porcelain", "--untracked-files=all", "--", ".agents/skills/cc-plan"); !strings.Contains(out, "cc-plan/SKILL.md") {
		t.Fatalf("a shipped skill is ignored: %q", out)
	}

	// The extension changes: setup refreshes the copy.
	write(t, filepath.Join(f.root, "extensions", "qa-pass", "skills", "qa", "SKILL.md"), qaSkill+"\nThen run the tests.\n")
	if note, _ := decodeMap(t, f.ok("status"))["setup_required"].(string); note == "" {
		t.Fatal("a changed skill left the copy current")
	}
	f.ok("agent", "setup")
	if !strings.Contains(read(t, filepath.Join(copied, "SKILL.md")), "Then run the tests.") {
		t.Fatal("setup did not refresh the copy")
	}

	// A copy a person edited is never overwritten.
	write(t, filepath.Join(copied, "references", "checklist.md"), "my own list\n")
	write(t, filepath.Join(f.root, "extensions", "qa-pass", "skills", "qa", "references", "checklist.md"), "# Checklist\n\n- Run the tests twice.\n")
	if out := f.fail("agent", "setup"); !strings.Contains(out, "preserve customized extension skill file") {
		t.Fatalf("setup overwrote an edited copy: %s", out)
	}
	if err := os.Remove(filepath.Join(copied, "references", "checklist.md")); err != nil {
		t.Fatal(err)
	}
	f.ok("agent", "setup")

	// Off for this member, the skill leaves on all three paths, and its
	// folders go with it.
	f.ok("extension", "disable", "--id", "qa-pass")
	result = workspace.AgentSetup{}
	if err := json.Unmarshal([]byte(f.ok("agent", "setup")), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Removed) != 4 || len(result.Kept) != 0 {
		t.Fatalf("setup after disable: %+v", result)
	}
	for _, gone := range []string{copied, filepath.Join(f.root, ".claude", "skills", "x-qa-pass-qa"), filepath.Join(f.root, ".cursor", "skills", "x-qa-pass-qa")} {
		if _, err := os.Stat(gone); !os.IsNotExist(err) {
			t.Fatalf("%s was not removed", gone)
		}
	}
	if read(t, filepath.Join(f.root, ".agents", "skills", "cc-plan", "SKILL.md")) != shipped {
		t.Fatal("setup changed a shipped skill")
	}
	if _, ok := decodeMap(t, f.ok("status"))["extension_skills"]; ok {
		t.Fatal("status lists the skill of an extension that is off")
	}

	// An edited wrapper is kept when its extension goes.
	f.ok("extension", "enable", "--id", "qa-pass")
	f.ok("agent", "setup")
	claudeWrapper := filepath.Join(f.root, ".claude", "skills", "x-qa-pass-qa", "SKILL.md")
	write(t, claudeWrapper, "my wrapper\n")
	f.ok("extension", "remove", "--id", "qa-pass")
	result = workspace.AgentSetup{}
	if err := json.Unmarshal([]byte(f.ok("agent", "setup")), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Kept) != 1 || result.Kept[0].Path != ".claude/skills/x-qa-pass-qa/SKILL.md" || len(result.Removed) != 3 {
		t.Fatalf("setup after remove: %+v", result)
	}
	if read(t, claudeWrapper) != "my wrapper\n" {
		t.Fatal("an edited wrapper was removed")
	}
}

func TestNamedStageHooksResolveAndSkip(t *testing.T) {
	f := setup(t)
	f.repository("api")
	qaExtension(t, f.root)
	extensionFolder(t, filepath.Join(f.root, "extensions", "qa-report"), reportManifest, "instructions/after-qa.md")

	// Hooked but not listed: check says no listed extension declares it, and
	// the point skips with that reason.
	f.ok("extension", "add", "--id", "qa-report", "--path", "extensions/qa-report")
	if out := f.fail("check"); !strings.Contains(out, "extension qa-report hooks after qa-pass/qa, but no listed extension declares that stage") {
		t.Fatalf("check did not report the undeclared stage: %s", out)
	}
	h := hooks(t, f.ok("extension", "hooks", "--stage", "qa-pass/qa", "--when", "after"))
	if len(h.Instructions) != 0 || len(h.Skipped) != 1 || h.Skipped[0].ID != "qa-pass" || !strings.Contains(h.Skipped[0].Reason, "not listed in extensions.yaml") {
		t.Fatalf("an unlisted stage's hooks: %+v", h)
	}

	f.ok("extension", "add", "--id", "qa-pass", "--path", "extensions/qa-pass")
	if out := f.ok("check"); strings.Contains(out, "qa-pass/qa") {
		t.Fatalf("a declared stage is a finding: %s", out)
	}
	i := f.intent("billing")
	p := f.plan(i.ID, "orders", "api")
	h = hooks(t, f.ok("extension", "hooks", "--stage", "qa-pass/qa", "--when", "after", "--plan", p.ID))
	if h.Point != "after qa-pass/qa" || len(h.Instructions) != 1 || h.Instructions[0].ID != "qa-report" || h.Instructions[0].File != "extensions/qa-report/instructions/after-qa.md" {
		t.Fatalf("the named stage's hooks: %+v", h)
	}
	if h = hooks(t, f.ok("extension", "hooks", "--stage", "qa-pass/qa", "--when", "before")); len(h.Instructions) != 0 || len(h.Skipped) != 0 {
		t.Fatalf("a point nobody hooks: %+v", h)
	}
	list := f.ok("extension", "list")
	if !strings.Contains(list, `"skills":["x-qa-pass-qa"]`) && !strings.Contains(list, `"x-qa-pass-qa"`) {
		t.Fatalf("list does not name the skill: %s", list)
	}

	// Off for the member, the stage is skipped with that reason.
	f.ok("extension", "disable", "--id", "qa-pass")
	h = hooks(t, f.ok("extension", "hooks", "--stage", "qa-pass/qa", "--when", "after"))
	if len(h.Instructions) != 0 || len(h.Skipped) != 1 || !strings.Contains(h.Skipped[0].Reason, "off for member maya") {
		t.Fatalf("a stage of an extension that is off: %+v", h)
	}
	f.ok("extension", "enable", "--id", "qa-pass")

	// A stage the listed extension does not declare.
	h = hooks(t, f.ok("extension", "hooks", "--stage", "qa-pass/smoke", "--when", "after"))
	if len(h.Skipped) != 1 || !strings.Contains(h.Skipped[0].Reason, "declares no stage smoke") {
		t.Fatalf("an undeclared stage: %+v", h)
	}
	for _, bad := range []string{"QA/qa", "qa-pass/", "ship"} {
		if out := f.fail("extension", "hooks", "--stage", bad, "--when", "after"); !strings.Contains(out, "stage") {
			t.Fatalf("%s was accepted: %s", bad, out)
		}
	}

	// A named stage never enters the lifecycle: a core stage's hooks do not
	// carry it, and a plan's ordering does not name it.
	if h = hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--plan", p.ID)); len(h.Instructions) != 0 {
		t.Fatalf("a core point carries a named stage's hook: %+v", h)
	}
	if out := f.ok("record", "order"); strings.Contains(out, "qa-pass") {
		t.Fatalf("record order names a named stage: %s", out)
	}
}

func TestExtensionSkillCheckFindings(t *testing.T) {
	base := "id: qa-pass\nmembers: default-on\nhooks: []\n"
	for _, tc := range []struct{ name, manifest, skill, want string }{
		{"id not a slug", base + "skills:\n  - id: QA Pass\n    path: skills/qa\n", qaSkill, "skill id"},
		{"id twice", base + "skills:\n  - id: qa\n    path: skills/qa\n  - id: qa\n    path: skills/qa\n", qaSkill, "skill qa appears twice"},
		{"stage not a slug", base + "skills:\n  - id: qa\n    path: skills/qa\n    stage: Q A\n", qaSkill, "skill qa stage"},
		{"stage twice", base + "skills:\n  - id: qa\n    path: skills/qa\n    stage: qa\n  - id: other\n    path: skills/qa\n    stage: qa\n", qaSkill, "both declare stage qa"},
		{"no path", base + "skills:\n  - id: qa\n", qaSkill, "names no path"},
		{"escaping path", base + "skills:\n  - id: qa\n    path: ../qa\n", qaSkill, "leaves its folder"},
		{"missing folder", base + "skills:\n  - id: qa\n    path: skills/absent\n", qaSkill, "skill qa"},
		{"no frontmatter", base + "skills:\n  - id: qa\n    path: skills/qa\n", "# QA\n", "frontmatter"},
		{"no description", base + "skills:\n  - id: qa\n    path: skills/qa\n", "---\nname: qa\n---\n\n# QA\n", "needs a name and a description"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := setup(t)
			write(t, filepath.Join(f.root, "extensions", "qa-pass", "extension.yaml"), tc.manifest)
			write(t, filepath.Join(f.root, "extensions", "qa-pass", "skills", "qa", "SKILL.md"), tc.skill)
			write(t, filepath.Join(f.root, "extensions.yaml"), "extensions:\n  - id: qa-pass\n    source: {path: extensions/qa-pass}\n")
			if out := f.fail("check"); !strings.Contains(out, tc.want) {
				t.Fatalf("check did not report %q: %s", tc.want, out)
			}
		})
	}

	// A hook at a stage its listed extension does not declare.
	f := setup(t)
	qaExtension(t, f.root)
	extensionFolder(t, filepath.Join(f.root, "extensions", "qa-report"), strings.Replace(reportManifest, "qa-pass/qa", "qa-pass/smoke", 1), "instructions/after-qa.md")
	write(t, filepath.Join(f.root, "extensions.yaml"), "extensions:\n  - id: qa-pass\n    source: {path: extensions/qa-pass}\n  - id: qa-report\n    source: {path: extensions/qa-report}\n")
	if out := f.fail("check"); !strings.Contains(out, "extension qa-report hooks after qa-pass/smoke, but qa-pass declares no stage smoke") {
		t.Fatalf("check did not report the undeclared stage: %s", out)
	}
}

// No template version ships a path starting with x- under any skills folder,
// so a delta between two template tags can never write or remove an
// installed extension skill.
func TestNoShippedSkillStartsWithX(t *testing.T) {
	underSkills := func(path string) bool {
		parts := strings.Split(filepath.ToSlash(path), "/")
		for i := 0; i+1 < len(parts); i++ {
			if parts[i] == "skills" && strings.HasPrefix(parts[i+1], "x-") {
				return true
			}
		}
		return false
	}
	manifest, err := os.ReadFile(filepath.Join("..", "..", "scripts", "release-manifest.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for line := range strings.SplitSeq(string(manifest), "\n") {
		for _, field := range strings.Fields(line) {
			if underSkills(field) {
				t.Errorf("the release manifest ships an x- skill path: %s", line)
			}
		}
	}
	files, err := contextcircuit.Files()
	if err != nil {
		t.Fatal(err)
	}
	for path := range files {
		if underSkills(path) {
			t.Errorf("a workspace is seeded with an x- skill path: %s", path)
		}
	}
	ignore := productFile(t, filepath.Join("..", "template", "gitignore.in"))
	for _, want := range []string{"/.agents/skills/x-*/", "/.claude/skills/x-*/", "/.cursor/skills/x-*/"} {
		if !strings.Contains(ignore, want) {
			t.Errorf("template/gitignore.in does not ignore %s", want)
		}
	}
}
