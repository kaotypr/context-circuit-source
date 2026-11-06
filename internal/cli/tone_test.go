package cli_test

import (
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
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
		{"both values", []string{"member", "tone", "--id", "maya", "--tone", "semi-formal", "--clear"}, "only one"},
		{"a preset and a line", []string{"member", "tone", "--id", "maya", "--preset", "casual", "--tone", "semi-formal"}, "only one"},
		{"a preset and clear", []string{"member", "tone", "--id", "maya", "--preset", "casual", "--clear"}, "only one"},
		{"an unknown preset", []string{"member", "tone", "--id", "maya", "--preset", "chatty"}, "choose colleague, casual, formal, non-technical"},
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

	// The procedure behind it lives in the writing skill, which the skills that
	// write records point to by its workspace path, since Codex has no route
	// file and reaches it only through those pointers.
	for _, skill := range []string{"skills/cc-intent/SKILL.md", "skills/cc-plan/SKILL.md", "skills/cc-complete/SKILL.md"} {
		if body := productFile(t, skill); !strings.Contains(body, "`.agents/skills/cc-writing/SKILL.md`") {
			t.Errorf("%s does not point to the cc-writing skill", skill)
		}
	}

	// The rules name operations rather than virtues, because "write naturally"
	// gives a model nothing to apply. They hold for a language with no reference
	// folder of its own, and the folders add the specifics for theirs.
	writing := productFile(t, "skills/cc-writing/SKILL.md")
	for _, phrase := range []string{"clause order", "nominalize", "engineers in that language actually", "translating"} {
		if !strings.Contains(writing, phrase) {
			t.Errorf("cc-writing does not name the operation %q", phrase)
		}
	}
}
	}
}

// A preset is a name for a line, not a second kind of tone: what members.yaml
// records is the same one line a hand-written tone fills, so a CLI or agent
// that knows nothing of presets still reads it as a plain instruction.
func TestAPresetRecordsItsLine(t *testing.T) {
	f := setup(t)
	roster := filepath.Join(f.root, "members.yaml")
	for _, preset := range workspace.TonePresets() {
		f.ok("member", "tone", "--id", "maya", "--preset", preset.Name)
		var listed struct {
			Members map[string]struct {
				Tone string `json:"tone"`
			} `json:"members"`
		}
		if err := json.Unmarshal([]byte(f.ok("member", "list")), &listed); err != nil {
			t.Fatal(err)
		}
		// The line holds ": " after its name, so it only survives the round
		// trip if members.yaml quotes it.
		if got := listed.Members["maya"].Tone; got != preset.Line {
			t.Fatalf("preset %s recorded %q, want %q:\n%s", preset.Name, got, preset.Line, read(t, roster))
		}
	}
	f.ok("member", "tone", "--id", "maya", "--clear")
	if body := read(t, roster); strings.Contains(body, "tone:") {
		t.Fatalf("expected --clear to remove a preset too:\n%s", body)
	}
}

func TestPresetsAreListedWithTheirLines(t *testing.T) {
	f := setup(t)
	var listed struct {
		Presets []workspace.TonePreset `json:"presets"`
	}
	if err := json.Unmarshal([]byte(f.ok("member", "presets")), &listed); err != nil {
		t.Fatal(err)
	}
	want := workspace.TonePresets()
	if len(listed.Presets) != len(want) {
		t.Fatalf("listed %d presets, want %d", len(listed.Presets), len(want))
	}
	for i, preset := range want {
		if listed.Presets[i] != preset {
			t.Errorf("preset %d listed as %+v, want %+v", i, listed.Presets[i], preset)
		}
	}
	if human := f.human("member", "presets"); !strings.Contains(human, want[0].Line) {
		t.Fatalf("expected the human listing to carry each line:\n%s", human)
	}
}

// Every preset is a valid tone in its own right, starts with its own name, and
// exactly one is the default the skill falls back to.
func TestEveryPresetIsAUsableTone(t *testing.T) {
	defaults := 0
	for _, preset := range workspace.TonePresets() {
		if err := workspace.Tone(preset.Line); err != nil {
			t.Errorf("preset %s: %v", preset.Name, err)
		}
		if !strings.HasPrefix(preset.Line, preset.Name+": ") {
			t.Errorf("preset %s does not start with its name: %q", preset.Name, preset.Line)
		}
		if preset.Default {
			defaults++
			if preset.Name != workspace.DefaultTonePreset {
				t.Errorf("preset %s is marked default, want %s", preset.Name, workspace.DefaultTonePreset)
			}
		}
	}
	if defaults != 1 {
		t.Errorf("%d presets are marked default, want 1", defaults)
	}
}

// The docs and the skill quote the preset lines rather than paraphrasing them,
// so a member reading either sees exactly what the CLI records.
func TestPresetLinesAreQuotedVerbatim(t *testing.T) {
	for _, rel := range []string{"docs/workspace.md", "docs/commands.md", "skills/cc-writing/SKILL.md"} {
		body := productFile(t, rel)
		for _, preset := range workspace.TonePresets() {
			if !strings.Contains(body, preset.Line) {
				t.Errorf("%s does not quote the %s line verbatim", rel, preset.Name)
			}
		}
	}
}

// A member's tone shapes what a person reads and decides on, and never the
// knowledge the next agent reads. Nothing mechanical can check a note's
// register, so the instruction is pinned at every place one skill hands over
// to the next, including completion, where a note in the member's tone and the
// knowledge it reconciles are written in one turn.
func TestToneStopsAtKnowledge(t *testing.T) {
	for rel, phrases := range map[string][]string{
		"AGENTS.md.in":                {"| `context/` notes, the catalog, the glossary | English, always |", "the glossary take no tone"},
		"skills/cc-writing/SKILL.md":  {"neither this skill nor the member's tone applies", "`.agents/skills/cc-knowledge/SKILL.md`"},
		"skills/cc-knowledge/SKILL.md": {"tone and `.agents/skills/cc-writing/SKILL.md` never apply here"},
		"skills/cc-complete/SKILL.md": {"English regardless", "`.agents/skills/cc-knowledge/SKILL.md`, with no tone"},
	} {
		body := productFile(t, rel)
		for _, phrase := range phrases {
			if !strings.Contains(body, phrase) {
				t.Errorf("%s does not carry %q", rel, phrase)
			}
		}
	}
}

// An agent opens only the reference files it is sent to, so a reference file
// the skill never names is shipped and never read.
func TestEveryWritingReferenceIsNamed(t *testing.T) {
	root := filepath.Join("..", "..", "product", "skills", "cc-writing")
	skill := productFile(t, "skills/cc-writing/SKILL.md")
	var found int
	err := filepath.WalkDir(filepath.Join(root, "references"), func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		found++
		rel, err := filepath.Rel(filepath.Join(root, "references"), path)
		if err != nil {
			return err
		}
		if !strings.Contains(skill, filepath.ToSlash(rel)) {
			t.Errorf("cc-writing never names references/%s", filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if found < 9 {
		t.Fatalf("expected the nine reference files, found %d", found)
	}
}
