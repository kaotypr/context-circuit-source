package surface

import (
	"strings"
	"testing"
)

func TestNormalize(t *testing.T) {
	n := NewNorm([]PathToken{
		{Path: "/tmp/run", Token: "$RUN"},
		{Path: "/tmp/run/ws", Token: "$WS"},
	}, "9.9.9")
	in := "created 2026-09-15T10:53:00Z on 2026-01-02 at /tmp/run/ws/intent/i100-demo.md plan p1000 version 9.9.9 under /tmp/run/gocache"
	got := n.Text(in)
	want := "created 2026-09-15T10:53:00Z on 2026-01-02 at $WS/intent/i100-demo.md plan p1000 version <version> under $RUN/gocache"
	if got != want {
		t.Fatalf("got %q", got)
	}
	raw := []byte("bin\x00" + in)
	if string(n.Bytes(raw)) != string(raw) {
		t.Fatal("NUL file was rewritten")
	}
}

func TestNormalizeShorterPathDoesNotEatWorkspace(t *testing.T) {
	n := NewNorm([]PathToken{
		{Path: "/tmp/run", Token: "$RUN"},
		{Path: "/tmp/run/workspace", Token: "$WS"},
	}, "")
	got := n.Text("see /tmp/run/workspace/README.md and /tmp/run/tmp")
	want := "see $WS/README.md and $RUN/tmp"
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestNormalizeKeepsExampleIDs(t *testing.T) {
	n := NewNorm(nil, "")
	if n.Text("see i100") == n.Text("see i101") {
		t.Fatal("help line ids were rewritten")
	}
	base := "---\nid: i100\ncreated_at: 2026-01-02T03:04:05Z\napproved_at: \"2026-01-02T03:04:06Z\"\n---\n## Request\nsee i100\n"
	branch := "---\nid: i101\ncreated_at: 2026-02-03T04:05:06Z\napproved_at: \"2026-02-03T04:05:07Z\"\n---\n## Request\nsee i100\n"
	if n.Text(base) != n.Text(branch) {
		t.Fatalf("id field still differs\n%s\n%s", n.Text(base), n.Text(branch))
	}
	prose := map[string][]byte{"help.txt": []byte("see i101\n")}
	files := CompareTrees(normalizeTree(map[string][]byte{"help.txt": []byte("see i100\n")}, n), normalizeTree(prose, n))
	if len(files) != 1 {
		t.Fatalf("help ids %#v", files)
	}
	skelBase := map[string][]byte{"intent/i100-demo.md": []byte(base)}
	skelBranch := map[string][]byte{"intent/i100-demo.md": []byte(branch)}
	if files := CompareTrees(normalizeTree(skelBase, n), normalizeTree(skelBranch, n)); len(files) != 0 {
		t.Fatalf("skeleton %#v", files)
	}
	if n.Text(`{"id":"i100","created_at":"2026-01-02T03:04:05Z"}`) != n.Text(`{"id":"i200","created_at":"2026-03-04T05:06:07Z"}`) {
		t.Fatal("json id field still differs")
	}
}

func TestJSONKeysUseRawStdout(t *testing.T) {
	if !jsonKeysChanged(`{"i100":1}`, `{"i200":1}`) {
		t.Fatal("object keys compared equal")
	}
	got := CLIReasons([]OutputChange{{
		Name:         "status-json",
		Changed:      true,
		StdoutBase:   `{"i100":1}`,
		StdoutBranch: `{"i200":1}`,
	}})
	if len(got) != 1 || !strings.Contains(got[0].Detail, "JSON keys") {
		t.Fatalf("%#v", got)
	}
}
