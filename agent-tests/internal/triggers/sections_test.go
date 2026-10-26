package triggers

import (
	"slices"
	"strings"
	"testing"
)

func TestSectionDiffLanguageOnly(t *testing.T) {
	before := "## Who is working\n\nRoster.\n\n## Language\n\nTalk.\n\n## Shared knowledge\n\nNotes.\n"
	after := strings.Replace(before, "## Language\n\nTalk.\n", "## Language\n\nTalk.\nAn extra sentence sits in this section.\n", 1)
	diff := SectionDiff(before, after)
	got := Sections(diff)
	if !slices.Equal(got, []string{"Language"}) {
		t.Fatalf("sections %v\n%s", got, diff)
	}
}

func TestSectionsUnifiedContext(t *testing.T) {
	diff := "" +
		"--- product/AGENTS.md.in\n" +
		"+++ product/AGENTS.md.in\n" +
		"@@ -1,3 +1,3 @@\n" +
		" ## Language\n" +
		"-Talk.\n" +
		"+Talk more.\n"
	got := Sections(diff)
	if !slices.Equal(got, []string{"Language"}) {
		t.Fatalf("%v", got)
	}
}

func TestSectionDiffRemovedHeading(t *testing.T) {
	before := "## Language\n\nTalk.\n\n## Shared knowledge\n\nNotes.\n"
	after := "## Language\n\nTalk.\n"
	got := Sections(SectionDiff(before, after))
	if !slices.Equal(got, []string{"Shared knowledge"}) {
		t.Fatalf("%v", got)
	}
}
