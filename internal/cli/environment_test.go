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
