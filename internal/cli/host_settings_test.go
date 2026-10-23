package cli_test

import (
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	assets "github.com/kaotypr/context-circuit-source"
)

// The shipped Claude Code settings auto-approve the workspace launcher and
// nothing else. A rule added here runs without a person seeing the prompt, so
// the exact list is pinned rather than a pattern it must contain.
var launcherAllowRules = []string{
	"Bash(sh .context-circuit/cli:*)",
	"Bash(powershell -NoProfile -ExecutionPolicy Bypass -File .context-circuit/cli.ps1:*)",
}

func TestClaudeSettingsAllowOnlyTheLauncher(t *testing.T) {
	files, err := assets.Files()
	if err != nil {
		t.Fatal(err)
	}
	data, ok := files[".claude/settings.json"]
	if !ok {
		t.Fatal("the seed does not ship .claude/settings.json")
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		t.Fatalf("settings are not valid JSON: %v", err)
	}
	if len(top) != 1 || top["permissions"] == nil {
		t.Fatalf("settings hold more than permissions: %s", data)
	}
	var permissions map[string]json.RawMessage
	if err := json.Unmarshal(top["permissions"], &permissions); err != nil {
		t.Fatal(err)
	}
	if len(permissions) != 1 || permissions["allow"] == nil {
		t.Fatalf("permissions hold more than allow: %s", top["permissions"])
	}
	var allow []string
	if err := json.Unmarshal(permissions["allow"], &allow); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(allow, launcherAllowRules) {
		t.Fatalf("allow rules are %q, want %q", allow, launcherAllowRules)
	}
}

func TestTemplateExportShipsClaudeSettings(t *testing.T) {
	f := setup(t)
	blank := filepath.Join(f.home, "blank")
	f.ok("template", "export", "--path", blank)
	got := strings.ReplaceAll(read(t, filepath.Join(blank, ".claude", "settings.json")), "\r\n", "\n")
	if want := productFile(t, "host/claude-settings.json"); got != want {
		t.Fatalf("exported settings differ from the product file:\n%s\nwant:\n%s", got, want)
	}
}
