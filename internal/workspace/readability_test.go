package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Shape is the part of readability a diagnostic can see. Each case is a note
// body; want is a fragment of the one finding it should produce, or empty for a
// note the pass leaves alone.
func TestNoteReadability(t *testing.T) {
	for _, tc := range []struct {
		name, body, want string
	}{
		{"anchor in a sentence", "The limit lives in `api@internal/support/limits.go`.", "code anchor outside"},
		{"anchor in an owner block", "The limit is a support concern.\n\nOwner:\n\n- `api@internal/support/` — the goodwill limit.", ""},
		{"anchor in a table", "| Path | What it owns |\n| --- | --- |\n| `api@internal/support/` | The goodwill limit |", ""},
		{"anchor in a fenced example", "Shape:\n\n```\n- `api@internal/support/` — the limit.\n```", ""},
		{"an address is not an anchor", "Ask billing@example.com before voiding.", ""},
		{"a clone URL is not an anchor", "Clone it from git@github.com:example/api.git today.", ""},
		{"a root file with an extension is an anchor", "Embedding happens in `api@assets.go` at build time.", "code anchor outside"},
		{"eight prose lines pass", strings.TrimSuffix(strings.Repeat("A short line of prose.\n", 8), "\n"), ""},
		{"nine prose lines are carrying a list", strings.TrimSuffix(strings.Repeat("A short line of prose.\n", 9), "\n"), "may be hiding in it"},
		{"a blank line ends a paragraph", strings.Repeat("A short line of prose.\n\n", 9), ""},
		{"a sentence past sixty words", strings.Repeat("word ", 61) + "ends.", "a sentence of"},
		{"mermaid with a known type", "```mermaid\nflowchart TD\n  a --> b\n```", ""},
		{"mermaid with an unknown type", "```mermaid\nflowChart TD\n  a --> b\n```", "not recognized"},
		{"an unclosed fence", "```sh\nmake build", "never closed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			note := "# Billing\n\nWhen an invoice is voided rather than credited.\n\n" + tc.body + "\n"
			issues := noteReadability("context/domains/billing.md", []byte(note))
			if tc.want == "" {
				if len(issues) != 0 {
					t.Fatalf("expected nothing to report: %v", issues)
				}
				return
			}
			if len(issues) != 1 || !strings.Contains(issues[0], tc.want) {
				t.Fatalf("expected %q: %v", tc.want, issues)
			}
			if !strings.Contains(issues[0], "readability:") {
				t.Fatalf("a shape finding says so: %v", issues)
			}
		})
	}
}

// The line under the title is what a reader sees before deciding to read on.
func TestNoteReadabilityWantsAThesisLine(t *testing.T) {
	for _, tc := range []struct {
		name, note string
		want       bool
	}{
		{"title then a line saying what it is for", "# Billing\n\nWhen an invoice is voided rather than credited.\n", false},
		{"title straight into a heading", "# Billing\n\n## Voiding\n\nVoiding happens here.\n", true},
		{"title straight into a list", "# Billing\n\n- Voiding happens here.\n", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			issues := noteReadability("context/domains/billing.md", []byte(tc.note))
			found := false
			for _, issue := range issues {
				if strings.Contains(issue, "no line under the title") {
					found = true
				}
			}
			if found != tc.want {
				t.Fatalf("want=%v: %v", tc.want, issues)
			}
		})
	}
}

// An actor note carries standing needs as stories, numbered once through the
// note so a number stays a handle worth citing.
func TestActorNoteCarriesStories(t *testing.T) {
	body := "# Support agent\n\nAnswers customer contact after an order is placed.\n\n"
	for _, tc := range []struct {
		name, note, path string
		want             bool
	}{
		{"a story in the documented shape", body + "## Refunds\n\n3. **As a support agent, I want to refund small amounts,** so that a call ends once.\n", "context/actors/support-agent.md", false},
		{"prose where stories belong", body + "## Refunds\n\nThis actor can refund small amounts without review.\n", "context/actors/support-agent.md", true},
		{"an unnumbered story is not a handle", body + "## Refunds\n\n- **As a support agent, I want to refund small amounts,** so that a call ends once.\n", "context/actors/support-agent.md", true},
		{"other notes carry no stories", body + "## Refunds\n\nThis is how refunds work.\n", "context/domains/refunds.md", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found := false
			for _, issue := range noteReadability(tc.path, []byte(tc.note)) {
				if strings.Contains(issue, "stories") {
					found = true
				}
			}
			if found != tc.want {
				t.Fatalf("want=%v for %s", tc.want, tc.path)
			}
		})
	}
}

// The glossary is a table of terms mapped to identifiers, which is the one place
// an anchor belongs outside an Owner block.
func TestGlossaryIsExemptFromNoteShape(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "context"), 0o755); err != nil {
		t.Fatal(err)
	}
	glossary := "# Glossary\n\n| Term | Meaning | Code |\n| --- | --- | --- |\n| Dunning | Retrying a failed charge | `api@internal/billing/dunning/` |\n"
	if err := os.WriteFile(filepath.Join(root, "context", "glossary.md"), []byte(glossary), 0o644); err != nil {
		t.Fatal(err)
	}
	issues, err := (&Store{Root: root}).knowledgeIssues()
	if err != nil {
		t.Fatal(err)
	}
	for _, issue := range issues {
		if strings.Contains(issue, "readability:") {
			t.Fatalf("the glossary is not a note: %v", issues)
		}
	}
}
