package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A registry the installer could not use is worse than none, because the
// default it would otherwise have applied is one that works.
func TestConfigRejectsAnUnusableCLIRegistry(t *testing.T) {
	for _, tc := range []struct {
		name, registry, want string
	}{
		{"unknown source", "cli_registry:\n  source: bitbucket\n  repository: acme/source\n", "expected github or gitlab"},
		{"gitlab without api", "cli_registry:\n  source: gitlab\n  repository: acme/source\n", "needs api"},
		{"plain http api", "cli_registry:\n  source: gitlab\n  repository: acme/source\n  api: http://gitlab.example.com/api/v4\n", "https"},
		{"no repository", "cli_registry:\n  source: github\n  repository: \"\"\n", "cli_registry repository"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := &Store{Root: t.TempDir()}
			document := "version: 2\nname: Acme\npurpose: Billing\nrepositories: {}\nrelationships: []\n" + tc.registry
			if err := os.WriteFile(filepath.Join(store.Root, "workspace.yaml"), []byte(document), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := store.Config()
			if err == nil {
				t.Fatal("an unusable registry read as a valid workspace")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error does not name %q: %v", tc.want, err)
			}
		})
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
