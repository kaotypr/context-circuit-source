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
