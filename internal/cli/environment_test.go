package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

func TestWorktreeReusesIgnoredEnvironmentAndPreservesUpdates(t *testing.T) {
	f := setup(t)
	repo := f.repository("web")
	write(t, filepath.Join(repo, ".gitignore"), "node_modules/\n.env*\n.cache/\n")
	write(t, filepath.Join(repo, "package.json"), "{\"name\":\"fixture\"}\n")
	git(t, repo, "add", ".gitignore", "package.json")
	git(t, repo, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "test: environment fixture")
	write(t, filepath.Join(repo, "node_modules", "module.js"), "installed dependency")
	write(t, filepath.Join(repo, ".env"), "FIXTURE_VALUE=synthetic-only\n")
	write(t, filepath.Join(repo, ".cache", "result"), "cached")
	path := filepath.Join(f.home, "worktree")
	args := []string{"worktree", "prepare", "--repo", "web", "--branch", "feature", "--path", path, "--copy-path", ".cache"}
	out := f.ok(args...)
	if strings.Contains(out, "synthetic-only") {
		t.Fatal("environment contents leaked")
	}
	for _, rel := range []string{"node_modules/module.js", ".env", ".cache/result"} {
		if read(t, filepath.Join(path, rel)) != read(t, filepath.Join(repo, rel)) {
			t.Fatal("missing reused entry", rel)
		}
	}
	write(t, filepath.Join(path, ".env"), "LOCAL_CHANGE=keep\n")
	f.ok(args...)
	if read(t, filepath.Join(path, ".env")) != "LOCAL_CHANGE=keep\n" {
		t.Fatal("resume overwrote worktree environment")
	}
	if read(t, filepath.Join(repo, ".env")) != "FIXTURE_VALUE=synthetic-only\n" {
		t.Fatal("shared mutable environment")
	}
	// A different dependency manifest must not receive stale installed modules.
	write(t, filepath.Join(repo, "package.json"), "{\"name\":\"changed\"}\n")
	different := filepath.Join(f.home, "different")
	out = f.ok("worktree", "prepare", "--repo", "web", "--branch", "different", "--path", different)
	if !strings.Contains(out, "dependency inputs differ") {
		t.Fatal(out)
	}
	if _, err := os.Stat(filepath.Join(different, "node_modules")); !os.IsNotExist(err) {
		t.Fatal("copied mismatched dependencies")
	}
	if read(t, filepath.Join(different, ".env")) == "" {
		t.Fatal("lost independent env reuse")
	}
	// Branch ignore rules protect env files from becoming untracked secrets.
	git(t, repo, "rm", ".gitignore")
	git(t, repo, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "test: branch without ignored runtime")
	unignored := filepath.Join(f.home, "unignored")
	f.ok("worktree", "prepare", "--repo", "web", "--branch", "unignored", "--path", unignored, "--copy-path", ".env")
	if _, err := os.Stat(filepath.Join(unignored, ".env")); !os.IsNotExist(err) {
		t.Fatal("copied unignored environment")
	}
}

func TestRoleSettingsApplyToNativeDefinitions(t *testing.T) {
	f := setup(t)
	for _, host := range []string{"codex", "claude-code", "cursor"} {
		f.ok("agent", "setup", "--host", host)
		f.ok("agent", "configure", "--host", host, "--role", "worker", "--model", "fixture-model", "--effort", "high")
		f.ok("agent", "setup", "--host", host)
		var file string
		switch host {
		case "codex":
			file = ".codex/agents/cc-worker.toml"
		case "claude-code":
			file = ".claude/agents/cc-worker.md"
		case "cursor":
			file = ".cursor/agents/cc-worker.md"
		}
		text := read(t, filepath.Join(f.root, file))
		if !strings.Contains(text, "fixture-model") || !strings.Contains(text, "high") {
			t.Fatal("settings were not applied", text)
		}
		out := f.ok("agent", "dispatch", "--host", host, "--role", "worker", "--task", "Implement the assigned fixture", "--path", f.root)
		if !strings.Contains(out, "fixture-model") || !strings.Contains(out, `"launch_required": true`) {
			t.Fatal(out)
		}
		write(t, filepath.Join(f.root, file), "custom user definition")
		f.fail("agent", "setup", "--host", host)
		if read(t, filepath.Join(f.root, file)) != "custom user definition" {
			t.Fatal("overwrote customized agent")
		}
	}
	f.fail("agent", "dispatch", "--host", "codex", "--role", "reviewer", "--task", "Review the requested diff", "--path", f.root)
	f.ok("agent", "dispatch", "--host", "codex", "--role", "reviewer", "--task", "Review the requested diff", "--path", f.root, "--review-requested")
	f.fail("agent", "configure", "--host", "codex", "--role", "worker", "--model", "fixture", "--effort", "nonsense")
	f.fail("agent", "configure", "--host", "cursor", "--role", "worker", "--model", "inherit", "--effort", "high")
}

// One roster's members run different hosts and pay for different models, so a
// machine answers the shared tiering without editing what everyone else reads.
func TestLocalRoleTieringOverridesTheSharedSetting(t *testing.T) {
	f := setup(t)
	shared := filepath.Join(f.root, workspace.SharedTiering)
	local := filepath.Join(f.root, workspace.LocalTiering)
	f.ok("agent", "configure", "--host", "codex", "--role", "worker", "--model", "team-model", "--effort", "medium")
	f.ok("agent", "configure", "--host", "codex", "--role", "planner", "--model", "team-planner", "--effort", "high")
	if _, err := os.Stat(local); !os.IsNotExist(err) {
		t.Fatal("a shared configure must not write the local file")
	}

	f.ok("agent", "configure", "--host", "codex", "--role", "worker", "--model", "my-model", "--effort", "low", "--local")
	if !strings.Contains(read(t, shared), "team-model") {
		t.Fatal("a local configure must not edit the shared file", read(t, shared))
	}
	var effective workspace.EffectiveTiering
	if err := json.Unmarshal([]byte(f.ok("agent", "settings")), &effective); err != nil {
		t.Fatal(err)
	}
	if effective.Hosts["codex"]["worker"].Model != "my-model" {
		t.Fatal("the local override is not in force", effective.Hosts["codex"]["worker"])
	}
	// An untouched role keeps tracking the shared file rather than being frozen
	// by the presence of an override elsewhere.
	if effective.Hosts["codex"]["planner"].Model != "team-planner" {
		t.Fatal("an unoverridden role must follow the shared file", effective.Hosts["codex"]["planner"])
	}
	if strings.Join(effective.LocalOverrides, ",") != "codex/worker" {
		t.Fatal("the merged view must name what this machine overrode", effective.LocalOverrides)
	}

	// A second host has no node in the local file yet; opening one must not
	// depend on which host happened to create the file.
	f.ok("agent", "configure", "--host", "claude-code", "--role", "planner", "--model", "other-model", "--effort", "high", "--local")
	if err := json.Unmarshal([]byte(f.ok("agent", "settings")), &effective); err != nil {
		t.Fatal(err)
	}
	if effective.Hosts["claude-code"]["planner"].Model != "other-model" || effective.Hosts["codex"]["worker"].Model != "my-model" {
		t.Fatal("a second host must join the local file, not replace it", effective.Hosts)
	}

	// The override reaches the native definitions the host actually loads.
	f.ok("agent", "setup", "--host", "codex")
	if text := read(t, filepath.Join(f.root, ".codex/agents/cc-worker.toml")); !strings.Contains(text, "my-model") {
		t.Fatal("setup ignored the local override", text)
	}

	// A local file is validated like the shared one, and cannot invent a host.
	sound := read(t, local)
	write(t, local, strings.Replace(sound, "effort: low", "effort: nonsense", 1))
	f.fail("agent", "settings")
	write(t, local, "hosts:\n  invented-host:\n    worker:\n      model: x\n      effort: low\n")
	f.fail("agent", "settings")
	write(t, local, sound)
	f.ok("agent", "settings")
}

func TestIndependentTemplateVersionAndCopyOptionPreflight(t *testing.T) {
	f := setup(t)
	blank := filepath.Join(f.home, "newer-template")
	f.ok("template", "export", "--path", blank)
	write(t, filepath.Join(blank, ".context-circuit", "VERSION"), "2.5.0\n")
	write(t, filepath.Join(blank, "AGENTS.md"), "Customized workspace instructions\n")
	g := fixture{t, blank, f.home}
	g.ok("init", "--name", "Other", "--purpose", "Fixture", "--member", "maya", "--member-name", "Maya")
	if read(t, filepath.Join(blank, ".context-circuit", "VERSION")) != "2.5.0\n" || read(t, filepath.Join(blank, "AGENTS.md")) != "Customized workspace instructions\n" {
		t.Fatal("CLI initialization rewrote template product")
	}
	g.ok("check")
	f.repository("web")
	for _, options := range [][]string{{"--copy-mode", "invalid"}, {"--copy-path", "../outside"}, {"--copy-path", ".git/config"}} {
		path := filepath.Join(f.home, "invalid")
		args := append([]string{"worktree", "prepare", "--repo", "web", "--branch", "invalid", "--path", path}, options...)
		f.fail(args...)
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatal("invalid copy options created a worktree")
		}
	}
}

// `npm link` leaves node_modules holding a link outside the checkout, which
// reuse refuses to copy. Reuse is a convenience, so prepare reports the entry
// and still returns the worktree, unless the caller required the copy.
func TestWorktreeSkipsAnEntryItCannotReuse(t *testing.T) {
	f := setup(t)
	repo := f.repository("web")
	write(t, filepath.Join(repo, ".gitignore"), "node_modules/\n")
	git(t, repo, "add", ".gitignore")
	git(t, repo, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "test: ignore modules")
	linked := filepath.Join(f.home, "linked-package")
	write(t, filepath.Join(linked, "index.js"), "linked")
	if err := os.MkdirAll(filepath.Join(repo, "node_modules"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(linked, filepath.Join(repo, "node_modules", "linked")); err != nil {
		t.Skip("symlinks unavailable:", err)
	}
	path := filepath.Join(f.home, "worktree")
	out := f.ok("worktree", "prepare", "--repo", "web", "--branch", "feature", "--path", path)
	if !strings.Contains(out, "external symlink") {
		t.Fatalf("the skipped entry was not reported: %s", out)
	}
	if _, err := os.Lstat(filepath.Join(path, "node_modules")); !os.IsNotExist(err) {
		t.Fatal("a refused entry was partly copied")
	}
	f.fail("worktree", "prepare", "--repo", "web", "--branch", "strict", "--path", filepath.Join(f.home, "strict"), "--copy-mode", "required")
}

const claudeLocalSettings = ".claude/settings.local.json"

// Claude Code adds an attribution trailer to every commit unless its settings
// turn that off, so setup writes the setting into the machine-local file a
// person also keeps their own settings in.
func TestAgentSetupTurnsOffClaudeCodeAttribution(t *testing.T) {
	f := setup(t)
	path := filepath.Join(f.root, claudeLocalSettings)
	// Initialization runs setup, so a new workspace already carries it.
	var settings map[string]map[string]string
	if err := json.Unmarshal([]byte(read(t, path)), &settings); err != nil {
		t.Fatal(err)
	}
	if got := settings["attribution"]; len(got) != 2 || got["commit"] != "" || got["pr"] != "" {
		t.Fatalf("attribution is not off: %s", read(t, path))
	}
	if !strings.Contains(strings.ReplaceAll(read(t, filepath.Join(f.root, ".gitignore")), "\r\n", "\n"), "/.claude/settings.local.json\n") {
		t.Fatal("the local settings file is not ignored")
	}
	for _, other := range []string{".cursor/cli.json", ".codex/config.toml"} {
		if _, err := os.Lstat(filepath.Join(f.root, other)); !os.IsNotExist(err) {
			t.Fatal("setup wrote a settings file the host does not read:", other)
		}
	}

	// A re-run changes nothing.
	before := read(t, path)
	f.ok("agent", "setup")
	if read(t, path) != before {
		t.Fatal("a re-run rewrote the settings file")
	}

	// An existing file keeps every byte it had and gains only the key.
	existing := "{\n  \"permissions\": {\"allow\": [\"Bash(make:*)\"]},\n  \"model\": \"fixture\",\n  \"limit\": 12345678901234567890\n}\n"
	write(t, path, existing)
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	f.ok("agent", "setup", "--host", "claude-code")
	merged := read(t, path)
	want := "{\n  \"permissions\": {\"allow\": [\"Bash(make:*)\"]},\n  \"model\": \"fixture\",\n  \"limit\": 12345678901234567890,\n  \"attribution\": {\"commit\": \"\", \"pr\": \"\"}\n}\n"
	if merged != want {
		t.Fatalf("merge changed the person's settings:\n%s\nwant:\n%s", merged, want)
	}
	if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0600 {
		t.Fatal("merge changed the file mode", info, err)
	}
	f.ok("agent", "setup", "--host", "claude-code")
	if read(t, path) != merged {
		t.Fatal("a re-run rewrote the merged file")
	}

	// A value the person set is theirs, whatever it is.
	own := "{\"attribution\": {\"commit\": \"Reviewed by the team\"}}"
	write(t, path, own)
	f.ok("agent", "setup", "--host", "claude-code")
	if read(t, path) != own {
		t.Fatal("setup replaced the person's attribution setting")
	}

	// An empty object gains the key; an unreadable file is refused and kept.
	write(t, path, "{}\n")
	f.ok("agent", "setup", "--host", "claude-code")
	if err := json.Unmarshal([]byte(read(t, path)), &settings); err != nil || settings["attribution"] == nil {
		t.Fatalf("an empty file did not gain the key: %s", read(t, path))
	}
	for _, malformed := range []string{"{\"model\": \"fixture\",}\n", "[]\n", "null\n", "{} {}\n", "// comment\n{}\n"} {
		write(t, path, malformed)
		roleFile := filepath.Join(f.root, ".claude/agents/cc-worker.md")
		role := read(t, roleFile)
		if out := f.fail("agent", "setup", "--host", "claude-code"); !strings.Contains(out, "preserve") {
			t.Fatalf("refusal does not say the file is preserved: %s", out)
		}
		if read(t, path) != malformed {
			t.Fatalf("a refused file was changed: %q", malformed)
		}
		if read(t, roleFile) != role {
			t.Fatal("a refused setup still rewrote role definitions")
		}
	}

	// A link is followed nowhere, so a file outside the workspace is never edited.
	outside := filepath.Join(f.home, "outside-settings.json")
	write(t, outside, "{}\n")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, path); err != nil {
		t.Skip("symlinks unavailable:", err)
	}
	f.fail("agent", "setup", "--host", "claude-code")
	if read(t, outside) != "{}\n" {
		t.Fatal("setup wrote through a symlink")
	}
}

// sources/ is passive evidence. The seed ships a root .ignore that keeps it out
// of search tools reading that file, anchored so a sources/ folder inside a
// cloned repository or a worktree stays searchable, and ignores the host's
// local settings file setup merges into.
func TestTemplateExportKeepsSourcesOutOfSearch(t *testing.T) {
	f := setup(t)
	blank := filepath.Join(f.home, "blank")
	f.ok("template", "export", "--path", blank)
	var rules []string
	for line := range strings.SplitSeq(strings.ReplaceAll(read(t, filepath.Join(blank, ".ignore")), "\r\n", "\n"), "\n") {
		if line != "" && !strings.HasPrefix(line, "#") {
			rules = append(rules, line)
		}
	}
	if strings.Join(rules, ",") != "/sources/" {
		t.Fatalf(".ignore rules are %q, want only the anchored /sources/", rules)
	}
	if !strings.Contains(strings.ReplaceAll(read(t, filepath.Join(blank, ".gitignore")), "\r\n", "\n"), "\n/.claude/settings.local.json\n") {
		t.Fatal("the seed .gitignore does not ignore the local Claude Code settings")
	}
	for _, rel := range []string{".ignore", "sources/README.md"} {
		if _, err := os.Stat(filepath.Join(f.root, rel)); err != nil {
			t.Fatal("an initialized workspace lacks", rel, err)
		}
	}
}
