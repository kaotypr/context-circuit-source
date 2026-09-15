package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// cli_registry names the GitLab project a workspace installs its CLI from. It
// is shared, so one member records it for everyone, and a value that cannot be
// resolved must fail here rather than on someone else's machine at install
// time. The accepted shape is the one the installers parse.
func TestConfigValidatesCLIRegistry(t *testing.T) {
	for _, tc := range []struct {
		name, value string
		reject      bool
	}{
		{"unset", "", false},
		{"project URL", "https://gitlab.example.com/acme/context-circuit", false},
		{"nested group", "https://gitlab.example.com/acme/platform/context-circuit", false},
		{"host with a port", "https://gitlab.example.com:8443/acme/context-circuit", false},
		{"trailing slash", "https://gitlab.example.com/acme/context-circuit/", false},
		{"plain http", "http://gitlab.example.com/acme/context-circuit", true},
		{"no scheme", "gitlab.example.com/acme/context-circuit", true},
		{"scheme only", "https://", true},
		{"host without a project", "https://gitlab.example.com", true},
		{"empty project path", "https://gitlab.example.com/", true},
		{"ssh remote", "git@gitlab.example.com:acme/context-circuit.git", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := &Store{Root: t.TempDir()}
			document := "version: 2\nname: Acme\npurpose: Billing\nrepositories: {}\nrelationships: []\n"
			if tc.value != "" {
				document += "cli_registry: \"" + tc.value + "\"\n"
			}
			if err := os.WriteFile(filepath.Join(store.Root, "workspace.yaml"), []byte(document), 0o644); err != nil {
				t.Fatal(err)
			}
			cfg, err := store.Config()
			if tc.reject {
				if err == nil {
					t.Fatalf("accepted %q", tc.value)
				}
				if !strings.Contains(err.Error(), "cli_registry must be an https project URL") {
					t.Fatalf("unhelpful failure for %q: %v", tc.value, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("rejected %q: %v", tc.value, err)
			}
			if cfg.CLIRegistry != tc.value {
				t.Fatalf("recorded %q, read back %q", tc.value, cfg.CLIRegistry)
			}
		})
	}
}

// An unknown field is an error, so a misspelled key cannot silently leave a
// workspace installing from the wrong place.
func TestConfigRejectsMisspelledCLIRegistry(t *testing.T) {
	store := &Store{Root: t.TempDir()}
	document := "version: 2\nname: Acme\npurpose: Billing\nrepositories: {}\nrelationships: []\n" +
		"cli_registery: \"https://gitlab.example.com/acme/context-circuit\"\n"
	if err := os.WriteFile(filepath.Join(store.Root, "workspace.yaml"), []byte(document), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Config(); err == nil {
		t.Fatal("accepted a misspelled cli_registry key")
	}
}

// Candidates through 2.0.0-rc.6 recorded base_branch in the shared file. Every
// workspace initialized by one hits this on its first read after upgrading, and
// a bare unknown-key error reads as a corrupt file rather than a moved value.
func TestConfigNamesWhereTheBaseBranchWent(t *testing.T) {
	for _, tc := range []struct {
		name, repositories string
		moved              bool
	}{
		{"earlier workspace", "repositories:\n  api:\n    base_branch: main\n", true},
		{"current workspace", "repositories:\n  api:\n    default_branch: main\n", false},
		{"no repositories", "repositories: {}\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := &Store{Root: t.TempDir()}
			document := "version: 2\nname: Acme\npurpose: Billing\n" + tc.repositories + "relationships: []\n"
			if err := os.WriteFile(filepath.Join(store.Root, "workspace.yaml"), []byte(document), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := store.Config()
			if !tc.moved {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("a moved key read as a valid workspace")
			}
			// The remedy is useless without both halves: where the base branch
			// goes, and what the shared entry holds in its place.
			for _, want := range []string{"repositories.local.yaml", "url and default_branch"} {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("error does not name %q: %v", want, err)
				}
			}
		})
	}
}
