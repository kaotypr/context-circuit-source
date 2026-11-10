package workspace

import (
	"strings"
	"testing"
)

// YAML 1.1 reads a bare on or off as a boolean. A switch is read from its text,
// so both spellings land as the same value and a boolean never does.
func TestMemberSwitchesReadOnAndOffAsWritten(t *testing.T) {
	for _, text := range []string{
		"extensions:\n  - id: parity\n    source: {path: x}\n    members:\n      maya: off\n      lee: on\n",
		"extensions:\n  - id: parity\n    source: {path: x}\n    members:\n      maya: \"off\"\n      lee: 'on'\n",
	} {
		var file ExtensionsFile
		if err := Decode([]byte(text), &file); err != nil {
			t.Fatalf("%v\n%s", err, text)
		}
		members := file.Extensions[0].Members
		if members["maya"] != "off" || members["lee"] != "on" {
			t.Fatalf("switches read as %v from:\n%s", members, text)
		}
	}
	for _, value := range []string{"true", "false", "yes", "enabled"} {
		var file ExtensionsFile
		text := "extensions:\n  - id: parity\n    source: {path: x}\n    members:\n      maya: " + value + "\n"
		if err := Decode([]byte(text), &file); err == nil {
			t.Errorf("a switch of %s was accepted", value)
		}
	}
}

func TestEditAddressesOneListItem(t *testing.T) {
	src := "# kept\nextensions:\n  # first\n  - id: a\n    source: {path: a}\n  # second\n  - id: b\n    source:\n      url: https://example.invalid/b.git\n      ref: v1\n"
	out, err := Edit([]byte(src), []string{"extensions", ListIndex(1), "source", "ref"}, "v2")
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, want := range []string{"# kept", "# first", "# second", "ref: v2", "id: a"} {
		if !strings.Contains(text, want) {
			t.Fatalf("edit lost %q:\n%s", want, text)
		}
	}
	if _, err := Edit([]byte(src), []string{"extensions", ListIndex(5), "id"}, "c"); err == nil {
		t.Fatal("an edit of an item that does not exist succeeded")
	}
	appended, err := AppendListItem(out, []string{"extensions"}, map[string]any{"id": "c", "source": map[string]any{"path": "c"}})
	if err != nil {
		t.Fatal(err)
	}
	var file ExtensionsFile
	if err := Decode(appended, &file); err != nil || len(file.Extensions) != 3 || file.Extensions[2].ID != "c" {
		t.Fatalf("append: %v %+v\n%s", err, file, appended)
	}
	removed, err := RemoveListItem(appended, []string{"extensions"}, 1)
	if err != nil {
		t.Fatal(err)
	}
	file = ExtensionsFile{}
	if err := Decode(removed, &file); err != nil || len(file.Extensions) != 2 || file.Extensions[1].ID != "c" {
		t.Fatalf("remove: %v %+v\n%s", err, file, removed)
	}
	if !strings.Contains(string(removed), "# kept") {
		t.Fatalf("remove lost the file's comment:\n%s", removed)
	}
}

func TestRequiresCLIComparesReleases(t *testing.T) {
	for _, c := range []struct {
		running, requirement string
		ok                   bool
	}{
		{"2.6.0", ">=2.6.0", true},
		{"2.6.1", ">=2.6.0", true},
		{"2.5.9", ">=2.6.0", false},
		{"2.10.0", ">=2.9.0", true},
		{"2.0.0-dev", ">=9.0.0", true},
		{"", ">=9.0.0", true},
	} {
		ok, err := cliSatisfies(c.running, c.requirement)
		if err != nil || ok != c.ok {
			t.Errorf("%s against %s: %v %v", c.running, c.requirement, ok, err)
		}
	}
	if _, err := cliSatisfies("2.6.0", "2.6.0"); err == nil {
		t.Error("a requirement without >= was accepted")
	}
}

func TestManifestValidation(t *testing.T) {
	valid := Manifest{ID: "parity", Members: "required", Hooks: []ManifestHook{{At: "after implement", Instructions: "a.md"}}}
	if err := validateManifest(valid); err != nil {
		t.Fatal(err)
	}
	script := &HookScript{Run: []string{"sh", "scripts/run.sh"}, OnFailure: "stop", Consent: "standing", Action: "Creates a task", Timeout: "10m"}
	withScript := valid
	withScript.Hooks = []ManifestHook{{At: "after plan", Script: script}}
	if err := validateManifest(withScript); err != nil {
		t.Fatalf("a well-formed script was refused: %v", err)
	}
	for name, broken := range map[string]Manifest{
		"bad members":   {ID: "parity", Members: "always", Hooks: valid.Hooks},
		"bad point":     {ID: "parity", Members: "opt-in", Hooks: []ManifestHook{{At: "during implement", Instructions: "a.md"}}},
		"bad stage":     {ID: "parity", Members: "opt-in", Hooks: []ManifestHook{{At: "after build", Instructions: "a.md"}}},
		"twice":         {ID: "parity", Members: "opt-in", Hooks: []ManifestHook{{At: "after plan", Instructions: "a.md"}, {At: "after plan", Instructions: "b.md"}}},
		"empty hook":    {ID: "parity", Members: "opt-in", Hooks: []ManifestHook{{At: "after plan"}}},
		"no on_failure": {ID: "parity", Members: "opt-in", Hooks: []ManifestHook{{At: "after plan", Script: &HookScript{Run: []string{"sh"}}}}},
		"no action":     {ID: "parity", Members: "opt-in", Hooks: []ManifestHook{{At: "after plan", Script: &HookScript{Run: []string{"sh"}, OnFailure: "report", Consent: "each-time"}}}},
		"escaping arg":  {ID: "parity", Members: "opt-in", Hooks: []ManifestHook{{At: "after plan", Script: &HookScript{Run: []string{"sh", "../x.sh"}, OnFailure: "report"}}}},
	} {
		if err := validateManifest(broken); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}
