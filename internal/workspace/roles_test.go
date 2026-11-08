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

// The shipped roles keep their briefs and access whatever extensions declare:
// a shipped role's return section is its instructions alone until an addition
// appends headings after them, and only the worker edits.
func TestShippedRolesStayAsTheyShip(t *testing.T) {
	for _, name := range ShippedRoles {
		r := Roles[name]
		if got := returnSection(r, nil); got != r.Instructions {
			t.Errorf("%s: the return section differs from the shipped instructions", name)
		}
		if r.edits(name) != (name == "worker") {
			t.Errorf("%s: edits is %v", name, r.edits(name))
		}
		withAddition := returnSection(r, []roleAddition{{Extension: "parity", Returns: []string{"Parity"}}})
		if !strings.HasPrefix(withAddition, r.Instructions) || !strings.HasSuffix(withAddition, "- Parity") {
			t.Errorf("%s: an addition's heading is not after the core text:\n%s", name, withAddition)
		}
	}
	if agentType("codex", "x-parity-parity-verifier") != "cc_x_parity_parity_verifier" || agentType("cursor", "worker") != "cc-worker" {
		t.Fatal("agent type names")
	}
}

func TestExtensionRolesAreValidated(t *testing.T) {
	valid := ManifestRole{ID: "verifier", Description: "Checks.", Access: AccessReadOnly, Like: "reviewer", Instructions: "v.md", Returns: []string{"Verdict"}}
	if err := validateManifestRoles(Manifest{Roles: []ManifestRole{valid}}); err != nil {
		t.Fatalf("a valid role was refused: %v", err)
	}
	// An extension role's own Verdict is no planner heading, so it stands;
	// an addition to the planner may not repeat one, in any case.
	if err := validateManifestRoles(Manifest{RoleAdditions: []RoleAddition{{Role: "planner", Returns: []string{"verdict"}}}}); err == nil || !strings.Contains(err.Error(), "core heading") {
		t.Fatalf("a planner core heading was accepted: %v", err)
	}
	if err := validateManifestRoles(Manifest{RoleAdditions: []RoleAddition{{Role: "worker", Returns: []string{"Parity", "parity"}}}}); err == nil || !strings.Contains(err.Error(), "twice") {
		t.Fatalf("a repeated heading was accepted: %v", err)
	}
	edit := valid
	edit.Access = "write"
	if err := validateManifestRoles(Manifest{Roles: []ManifestRole{edit}}); err == nil {
		t.Fatal("an unknown access was accepted")
	}
	limits := extensionRoleLimits(Role{Extension: "parity"})
	for _, want := range []string{"no push", "no approval", "starts neither planning nor execution", "launches no other agent", "Commit only on the branch"} {
		if !strings.Contains(limits, want) && !strings.Contains(limits, strings.Replace(want, "no approval", "gives no approval", 1)) {
			t.Errorf("an edit role's limits miss %q", want)
		}
	}
}
