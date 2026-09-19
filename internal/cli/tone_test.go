package cli_test

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestAMemberRecordsTheRegisterTheyWriteIn(t *testing.T) {
	f := setup(t)
	roster := filepath.Join(f.root, "members.yaml")

	f.ok("member", "language", "--id", "maya", "--language", "Bahasa Indonesia")
	f.ok("member", "tone", "--id", "maya", "--tone", "semi-formal; keep technical terms in English")
	if body := read(t, roster); !strings.Contains(body, "tone: semi-formal; keep technical terms in English") {
		t.Fatalf("expected the register beside the language:\n%s", body)
	}
	if listed := f.ok("member", "list"); !strings.Contains(listed, "semi-formal") {
		t.Fatalf("expected the register to be readable back:\n%s", listed)
	}

	// A team that changes its mind wants the general guidance back, and prose
	// has no sentinel the way band 0 clears a band.
	f.ok("member", "tone", "--id", "maya", "--clear")
	body := read(t, roster)
	if strings.Contains(body, "tone:") {
		t.Fatalf("expected --clear to remove the register:\n%s", body)
	}
	if !strings.Contains(body, "language: Bahasa Indonesia") {
		t.Fatalf("clearing the register must leave the language alone:\n%s", body)
	}

	// Recorded at the same time as the member, for a roster written in one go.
	f.ok("member", "add", "--id", "alex", "--name", "Alex", "--band", "2",
		"--language", "Bahasa Indonesia", "--tone", "semi-formal")
	if body := read(t, roster); !strings.Contains(body, "tone: semi-formal") {
		t.Fatalf("expected member add to carry the register:\n%s", body)
	}
}

func TestARecordedRegisterIsOneUsableLine(t *testing.T) {
	f := setup(t)
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"neither value", []string{"member", "tone", "--id", "maya"}, "--clear"},
		{"both values", []string{"member", "tone", "--id", "maya", "--tone", "semi-formal", "--clear"}, "not both"},
		{"unknown member", []string{"member", "tone", "--id", "nobody", "--tone", "semi-formal"}, "unknown member"},
		{"an essay", []string{"member", "tone", "--id", "maya", "--tone", strings.Repeat("formal ", 40)}, "one line"},
		{"several lines", []string{"member", "tone", "--id", "maya", "--tone", "semi-formal\nand also"}, "single-line"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if out := f.fail(tc.args...); !strings.Contains(out, tc.want) {
				t.Fatalf("expected a refusal naming %q:\n%s", tc.want, out)
			}
		})
	}
}

// The structure rule is a standing one: it holds whether or not a skill was
// loaded, so it belongs in the always-loaded instruction and is asserted there
// by name, the way the gates and prohibitions are.
func TestRecordStructureStaysEnglishInTheEntryInstruction(t *testing.T) {
	entry := productFile(t, "AGENTS.md.in")
	for _, phrase := range []string{
		"A record's headings and field names | English, always",
		"headings are English",
	} {
		if !strings.Contains(entry, phrase) {
			t.Errorf("the entry instruction does not carry %q", phrase)
		}
	}

	// The procedure behind it lives in the skills that write records, and both
	// of them write in the author's language.
	for _, skill := range []string{"skills/cc-intent/SKILL.md", "skills/cc-plan/SKILL.md"} {
		body := productFile(t, skill)
		for _, phrase := range []string{"tone", "translating"} {
			if !strings.Contains(strings.ToLower(body), phrase) {
				t.Errorf("%s does not carry %q", skill, phrase)
			}
		}
	}

	// The rules name operations rather than virtues, because "write naturally"
	// gives a model nothing to apply. The skill stays language-agnostic: a
	// worked example in one language is noise to a member writing another, and
	// a team that wants specifics records a tone instead.
	intent := productFile(t, "skills/cc-intent/SKILL.md")
	for _, phrase := range []string{"clause order", "nominalize", "engineers in that language actually"} {
		if !strings.Contains(intent, phrase) {
			t.Errorf("cc-intent does not name the operation %q", phrase)
		}
	}
}
