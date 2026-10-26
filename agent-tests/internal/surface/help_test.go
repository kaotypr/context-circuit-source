package surface

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

func TestHelpTopicScenarios(t *testing.T) {
	got := helpScenariosFor("internal/cli/topics.go", []OutputChange{
		{Name: "help", Changed: true},
		{Name: "help-record", Changed: true},
		{Name: "help-worktree", Changed: false},
		{Name: "help-knowledge", Changed: true},
		{Name: "help-context", Changed: true},
		{Name: "help-agent", Changed: true},
		{Name: "status", Changed: true},
	})
	want := []string{"borrowed", "complete", "intent", "join", "plan"}
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Fatalf("%v", got)
	}
	if helpScenariosFor("internal/cli/commands.go", []OutputChange{{Name: "help-record", Changed: true}}) != nil {
		t.Fatal("topics mapping applied to another file")
	}
}

func TestArchiveCLIVersion(t *testing.T) {
	checkout := t.TempDir()
	archive := t.TempDir()
	if err := os.WriteFile(filepath.Join(checkout, "CLI_VERSION"), []byte("1.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(archive, "CLI_VERSION"), []byte("2.5.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := archiveCLIVersion(archive)
	if err != nil || got != "2.5.0" {
		t.Fatalf("%q %v", got, err)
	}
	checkoutPin, err := world.Pin(checkout)
	if err != nil || checkoutPin == got {
		t.Fatalf("checkout pin %q err %v", checkoutPin, err)
	}
}
