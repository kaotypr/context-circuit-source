package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

func TestCopyDesignRejectsPathsOutsideTheWorkspace(t *testing.T) {
	dir := t.TempDir()
	ws := filepath.Join(dir, "ws")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(dir, "outside.txt")
	if err := os.WriteFile(outside, []byte("private"), 0o644); err != nil {
		t.Fatal(err)
	}
	rt := &runtime{}
	for _, rel := range []string{"../outside.txt", outside, "", "sub/../../outside.txt", "../*"} {
		sc := scenario.Scenario{DesignArtifacts: []string{rel}}
		err := rt.copyDesign(world.World{WS: ws, RunRoot: dir}, sc, filepath.Join(dir, "out"), nil)
		if err == nil {
			t.Fatalf("design artifact %q was accepted", rel)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "out")); err == nil {
		t.Fatal("something was written to the output directory")
	}
}

func TestCopyDesignSkipsSymlinkThatLeavesTheWorkspace(t *testing.T) {
	dir := t.TempDir()
	ws := filepath.Join(dir, "ws")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(dir, "outside.txt")
	if err := os.WriteFile(outside, []byte("private"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(ws, "link.md")); err != nil {
		t.Skip("symlinks unavailable")
	}
	var drops []string
	rt := &runtime{}
	sc := scenario.Scenario{DesignArtifacts: []string{"link.md"}}
	if err := rt.copyDesign(world.World{WS: ws, RunRoot: dir}, sc, filepath.Join(dir, "out"), &drops); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "out", "link.md")); err == nil {
		t.Fatal("a symlink to a file outside the workspace was copied")
	}
	if len(drops) != 1 {
		t.Fatalf("drops = %v", drops)
	}
}

func TestCopyDesignGlobAndDirectory(t *testing.T) {
	dir := t.TempDir()
	ws := filepath.Join(dir, "ws")
	notes := filepath.Join(ws, "notes")
	plans := filepath.Join(ws, "plans", "item")
	if err := os.MkdirAll(filepath.Join(notes, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(plans, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(path, body string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(notes, "a.txt"), "a")
	write(filepath.Join(notes, "sub", "b.txt"), "b")
	write(filepath.Join(plans, "plan.md"), "plan")
	write(filepath.Join(plans, "extra.txt"), "extra")
	write(filepath.Join(ws, "skip.txt"), "skip")
	outside := filepath.Join(dir, "outside.txt")
	write(outside, "private-marker")
	if err := os.Symlink(outside, filepath.Join(notes, "sub", "leak")); err != nil {
		t.Skip("symlinks unavailable")
	}
	var drops []string
	rt := &runtime{}
	sc := scenario.Scenario{DesignArtifacts: []string{"notes", "plans/*", "missing/*.md", "["}}
	err := rt.copyDesign(world.World{WS: ws, RunRoot: dir}, sc, filepath.Join(dir, "out"), &drops)
	if err == nil || !strings.Contains(err.Error(), "[") {
		t.Fatalf("bad pattern: %v", err)
	}
	sc.DesignArtifacts = []string{"notes", "plans/*", "missing/*.md"}
	drops = nil
	if err := rt.copyDesign(world.World{WS: ws, RunRoot: dir}, sc, filepath.Join(dir, "out"), &drops); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	for _, rel := range []string{"notes/a.txt", "notes/sub/b.txt", "plans/item/plan.md", "plans/item/extra.txt"} {
		body, err := os.ReadFile(filepath.Join(out, rel))
		if err != nil {
			t.Fatalf("%s: %v", rel, err)
		}
		if len(body) == 0 {
			t.Fatalf("%s empty", rel)
		}
	}
	if _, err := os.Stat(filepath.Join(out, "skip.txt")); err == nil {
		t.Fatal("a file outside the pattern was copied")
	}
	if _, err := os.Stat(filepath.Join(out, "notes", "sub", "leak")); err == nil {
		t.Fatal("a symlink out of the workspace was copied")
	}
	copied := treeText(t, out)
	if strings.Contains(copied, "private-marker") {
		t.Fatal("outside file contents were copied")
	}
	found := false
	for _, line := range drops {
		if line == "design artifact missing/*.md: no match" {
			found = true
		}
	}
	if !found {
		t.Fatalf("drops %v", drops)
	}
}

func TestRedactorRemovesTheCallersLiteralKeys(t *testing.T) {
	rt := &runtime{opts: RunOptions{Caller: map[string]string{
		"CURSOR_API_KEY": "key_0123456789abcdef",
		"PATH":           "/usr/bin",
	}}}
	red := rt.redact()
	got := red.text("error: bad key key_0123456789abcdef in request")
	if strings.Contains(got, "key_0123456789abcdef") {
		t.Fatalf("key survived: %s", got)
	}
	path := filepath.Join(t.TempDir(), "note.txt")
	var drops []string
	if err := writeBytes(path, []byte("leaked key_0123456789abcdef"), red, &drops); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err == nil || len(drops) != 1 {
		t.Fatalf("a file holding the literal key was written (drops=%v)", drops)
	}
}
