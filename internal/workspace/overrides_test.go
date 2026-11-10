package workspace

import (
	"strings"
	"testing"
)

func TestOverrideManifestValidation(t *testing.T) {
	verifier := ManifestRole{ID: "verifier", Description: "Checks.", Access: AccessReadOnly, Instructions: "v.md"}
	fixer := ManifestRole{ID: "fixer", Description: "Fixes.", Access: AccessEdit, Instructions: "f.md"}
	base := Manifest{ID: "parity", Members: ExtensionDefaultOn, Roles: []ManifestRole{verifier, fixer}}
	with := func(overrides ...ManifestOverride) Manifest {
		m := base
		m.Overrides = overrides
		return m
	}
	valid := with(
		ManifestOverride{Default: OverrideRepairLoop, Verifier: "verifier", Rounds: 10, Instructions: "loop.md"},
		ManifestOverride{Default: OverrideVerify, Verifier: "reviewer"},
		ManifestOverride{Default: OverrideRiskReview, Instructions: "risk.md"},
		ManifestOverride{Default: OverrideDelegation, Stage: "implement", Role: "fixer"},
	)
	if err := validateManifest(valid); err != nil {
		t.Fatal(err)
	}
	for name, tc := range map[string]struct {
		manifest Manifest
		want     string
	}{
		// The fixed rules are not on the list, so naming one is invalid.
		"approve":            {with(ManifestOverride{Default: "approve"}), "not an overridable default"},
		"push authorization": {with(ManifestOverride{Default: "push"}), "approving an intent, authorizing push"},
		"twice":              {with(ManifestOverride{Default: OverrideRiskReview, Instructions: "a.md"}, ManifestOverride{Default: OverrideRiskReview, Instructions: "b.md"}), "appears twice"},
		"no verifier":        {with(ManifestOverride{Default: OverrideRepairLoop, Rounds: 2}), "names no verifier"},
		"edit verifier":      {with(ManifestOverride{Default: OverrideVerify, Verifier: "fixer"}), "a verifier is read-only"},
		"unknown verifier":   {with(ManifestOverride{Default: OverrideVerify, Verifier: "worker"}), "neither the shipped reviewer"},
		"no rounds":          {with(ManifestOverride{Default: OverrideRepairLoop, Verifier: "verifier"}), "rounds is between 1 and 10"},
		"too many rounds":    {with(ManifestOverride{Default: OverrideRepairLoop, Verifier: "verifier", Rounds: 11}), "rounds is between 1 and 10"},
		"verify rounds":      {with(ManifestOverride{Default: OverrideVerify, Verifier: "verifier", Rounds: 2}), "takes no rounds"},
		"risk without text":  {with(ManifestOverride{Default: OverrideRiskReview}), "names no instructions"},
		"edit reviewer":      {with(ManifestOverride{Default: OverrideRiskReview, Instructions: "r.md", Reviewer: "fixer"}), "a verifier is read-only"},
		"delegation stage":   {with(ManifestOverride{Default: OverrideDelegation, Stage: "deliver", Role: "worker"}), "stage is plan, implement, review"},
		"delegation access":  {with(ManifestOverride{Default: OverrideDelegation, Stage: "review", Role: "worker"}), "needs a role with read-only access"},
		"delegation unknown": {with(ManifestOverride{Default: OverrideDelegation, Stage: "plan", Role: "boss"}), "neither a shipped role"},
		"delegation extra":   {with(ManifestOverride{Default: OverrideDelegation, Stage: "plan", Role: "planner", Rounds: 2}), "takes no rounds"},
	} {
		err := validateManifest(tc.manifest)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: want %q, got %v", name, tc.want, err)
		}
	}
}

func TestRemoveMapKeyKeepsComments(t *testing.T) {
	data := []byte("# head\nextensions:\n  - id: parity\n    # allowed by a person\n    overrides:\n      repair-loop: {rounds: 2}\n      risk-review: {}\n    repositories: [api]\n")
	out, err := RemoveMapKey(data, []string{"extensions", ListIndex(0), "overrides", "risk-review"})
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	if strings.Contains(text, "risk-review") || !strings.Contains(text, "repair-loop") || !strings.Contains(text, "# head") || !strings.Contains(text, "repositories: [api]") {
		t.Fatalf("removing one key:\n%s", text)
	}
	out, err = RemoveMapKey(out, []string{"extensions", ListIndex(0), "overrides"})
	if err != nil {
		t.Fatal(err)
	}
	if text := string(out); strings.Contains(text, "overrides") || strings.Contains(text, "repair-loop") || !strings.Contains(text, "id: parity") {
		t.Fatalf("removing the map:\n%s", text)
	}
	if _, err := RemoveMapKey(out, []string{"extensions", ListIndex(0), "absent"}); err == nil {
		t.Fatal("an absent key was removed")
	}
}
