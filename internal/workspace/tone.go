package workspace

import (
	"fmt"
	"strings"
)

// TonePreset is one of the registers the product names, so a member can pick
// one instead of describing a register from scratch. Its line is what
// members.yaml records, the same field a hand-written tone fills, so a CLI or
// agent that knows nothing of presets still reads it as a plain instruction.
// The line starts with the preset's name and a colon, which is how the
// cc-writing skill recognizes the preset even after its wording changes.
type TonePreset struct {
	Name    string `yaml:"name" json:"name"`
	Line    string `yaml:"line" json:"line"`
	Default bool   `yaml:"default" json:"default"`
}

// DefaultTonePreset is the register an agent writes in for a member who
// records no tone.
const DefaultTonePreset = "colleague"

// TonePresets lists the presets in the order they are offered. These lines are
// the single source of the wording: the shipped docs and the cc-writing skill
// quote them verbatim, and a test holds them to it.
func TonePresets() []TonePreset {
	return []TonePreset{
		{Name: "colleague", Line: "colleague: semi-formal, the way you'd explain it to a teammate; technical terms as the team says them", Default: true},
		{Name: "casual", Line: "casual: chat register, short and relaxed; still complete sentences"},
		{Name: "formal", Line: "formal: for managers or people outside the team; no slang or contractions, still plain"},
		{Name: "non-technical", Line: "non-technical: for readers who don't read code; goal, non-goals and questions say what people will see or can do, and the rest stays technical for the team"},
	}
}

// TonePresetLine returns the line a preset records, or an error naming every
// preset when the name is not one of them.
func TonePresetLine(name string) (string, error) {
	names := make([]string, 0, 4)
	for _, preset := range TonePresets() {
		if preset.Name == name {
			return preset.Line, nil
		}
		names = append(names, preset.Name)
	}
	return "", fmt.Errorf("unknown tone preset %q; choose %s", name, strings.Join(names, ", "))
}
