package workspace

import (
	"strings"
	"testing"
)

// Setup renames a copied skill and leaves every other byte as written.
func TestRenameSkillKeepsTheRest(t *testing.T) {
	in := "---\ndescription: Runs QA.\nname: qa\nallowed-tools: Read\n---\n\nname: not frontmatter\n"
	out, err := renameSkill([]byte(in), "x-qa-pass-qa")
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != strings.Replace(in, "name: qa\n", "name: x-qa-pass-qa\n", 1) {
		t.Fatalf("renamed:\n%s", out)
	}
	crlf := "---\r\nname: qa\r\ndescription: Runs QA.\r\n---\r\n\r\nBody\r\n"
	out, err = renameSkill([]byte(crlf), "x-qa-pass-qa")
	if err != nil || string(out) != strings.Replace(crlf, "name: qa", "name: x-qa-pass-qa", 1) {
		t.Fatalf("CRLF: %v %q", err, out)
	}
	if _, err := renameSkill([]byte("# no frontmatter\n"), "x"); err == nil {
		t.Fatal("a skill without frontmatter was renamed")
	}
}

func TestNamedHookPoints(t *testing.T) {
	point, err := ParseHookPoint("after qa-pass/qa")
	if err != nil {
		t.Fatal(err)
	}
	if extension, stage, ok := point.Named(); !ok || extension != "qa-pass" || stage != "qa" {
		t.Fatalf("named: %q %q %v", extension, stage, ok)
	}
	if _, _, ok := (HookPoint{"after", "implement"}).Named(); ok {
		t.Fatal("a core stage reads as named")
	}
	for _, bad := range []string{"after QA/qa", "after qa-pass/", "after /qa", "after qa-pass/qa/x", "after ship"} {
		if _, err := ParseHookPoint(bad); err == nil {
			t.Errorf("%q was accepted", bad)
		}
	}
}
