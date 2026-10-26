package surface

import (
	"strings"
	"testing"
)

func TestCompareTreesMoveAndChange(t *testing.T) {
	base := map[string][]byte{
		"keep.md": []byte("same"),
		"edit.md": []byte("old"),
		"a.md":    []byte("moved-body"),
		"b.md":    []byte("pair"),
		"c.md":    []byte("pair"),
		"gone.md": []byte("only-base"),
	}
	branch := map[string][]byte{
		"keep.md":  []byte("same"),
		"edit.md":  []byte("new"),
		"dest.md":  []byte("moved-body"),
		"d.md":     []byte("pair"),
		"e.md":     []byte("pair"),
		"fresh.md": []byte("only-branch"),
	}
	got := CompareTrees(base, branch)
	want := map[string]FileChange{
		"edit.md":  {Status: "changed", Path: "edit.md"},
		"dest.md":  {Status: "moved", Path: "dest.md", From: "a.md"},
		"d.md":     {Status: "moved", Path: "d.md", From: "b.md"},
		"e.md":     {Status: "moved", Path: "e.md", From: "c.md"},
		"gone.md":  {Status: "removed", Path: "gone.md"},
		"fresh.md": {Status: "added", Path: "fresh.md"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %#v", got)
	}
	for _, file := range got {
		exp, ok := want[file.Path]
		if !ok || file != exp {
			t.Fatalf("file %#v", file)
		}
	}
}

func TestUnifiedShowsTheNewLine(t *testing.T) {
	diff := Unified("help", "one\ntwo\n", "one\ntwo\nthree\n")
	if !strings.Contains(diff, "+three") {
		t.Fatalf("%s", diff)
	}
}
