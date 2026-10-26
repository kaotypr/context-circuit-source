package surface

import "strings"

import "testing"

func TestRenderEmpty(t *testing.T) {
	text := Render(Report{
		BaseRef:   "main",
		Base:      "abcdef1234567890",
		BranchRef: "HEAD",
		Branch:    "abcdef1234567890",
		Empty:     true,
		Message:   emptyMessage,
		Mapping:   mappingNote,
	})
	if !strings.Contains(text, emptyMessage) {
		t.Fatalf("%s", text)
	}
	if strings.Contains(text, "Looks larger") {
		t.Fatal(text)
	}
}

func TestRenderReasonTags(t *testing.T) {
	text := Render(Report{
		Message: "Surface diff: 1 workspace files and 0 CLI outputs differ.",
		Larger:  true,
		LargerReasons: []Reason{
			{Kind: "layout", Detail: "moved", Exact: true},
			{Kind: "skill", Detail: "record skeleton headings or frontmatter changed"},
		},
		Sizes:       []SizeChange{{Path: "AGENTS.md"}},
		OrientHosts: []string{"claude", "codex", "cursor"},
	})
	for _, part := range []string{"Looks larger", "- layout [exact]: moved", "- skill [heuristic]: record skeleton headings or frontmatter changed", "orient on all three hosts: claude, codex, cursor"} {
		if !strings.Contains(text, part) {
			t.Errorf("missing %q\n%s", part, text)
		}
	}
}

func TestRenderLarger(t *testing.T) {
	text := Render(Report{
		Message: "Surface diff: 1 workspace files and 0 CLI outputs differ.",
		Files:   []FileChange{{Status: "moved", From: "a.md", Path: "b.md"}},
		Sizes: []SizeChange{{
			Path: "AGENTS.md", BaseBytes: 8, BranchBytes: 8, BaseTokens: 2, BranchTokens: 2,
		}},
		Larger: true,
		LargerReasons: []Reason{{
			Kind:   "gate",
			Detail: "keyword heuristic on added and removed lines",
		}},
		Hosts: []string{},
	})
	for _, part := range []string{"Looks larger", "[heuristic]", "A person decides", "approx", "any one host", "moved a.md -> b.md"} {
		if !strings.Contains(text, part) {
			t.Errorf("missing %q\n%s", part, text)
		}
	}
}
