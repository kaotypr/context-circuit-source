package workspace

import (
	"context"
	"strings"
	"testing"
)

// Roles in a schema 2 roster are reported, because a 2.3.0 CLI reading that
// workspace stops at the unknown field instead of asking for an upgrade. The
// same roster at schema 3 is ordinary.
func TestASchema2RosterWithRolesIsReported(t *testing.T) {
	const roster = "members:\n  maya:\n    name: Maya\n    band: 1\n    roles: [maintainer]\n  alex:\n    name: Alex\n    band: 2\n"
	reported := func(store *Store) []string {
		t.Helper()
		if err := store.Write("members.yaml", []byte(roster), 0o644); err != nil {
			t.Fatal(err)
		}
		state, err := store.Status(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, finding := range state.Issues {
			if strings.Contains(finding.Issue, "roles") {
				out = append(out, finding.Issue)
			}
		}
		return out
	}

	got := reported(legacy(t, Ledger{}))
	if len(got) != 1 || !strings.Contains(got[0], "member maya records roles in a schema 2 workspace") {
		t.Fatalf("expected one finding naming maya, got %q", got)
	}
	if got := reported(banded(t, map[string]*int{"maya": BandOf(1), "alex": BandOf(2)}, Ledger{})); len(got) != 0 {
		t.Fatalf("a schema 3 roster with roles was reported: %q", got)
	}
}
