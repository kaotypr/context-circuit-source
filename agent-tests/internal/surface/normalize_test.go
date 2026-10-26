package surface

import "testing"

func TestNormalize(t *testing.T) {
	n := NewNorm([]PathToken{
		{Path: "/tmp/run", Token: "$RUN"},
		{Path: "/tmp/run/ws", Token: "$WS"},
	}, "9.9.9")
	in := "created 2026-09-15T10:53:00Z on 2026-01-02 at /tmp/run/ws/intent/i100-demo.md plan p1000 version 9.9.9 under /tmp/run/gocache"
	got := n.Text(in)
	want := "created <time> on <date> at $WS/intent/iNNN-demo.md plan pNNNN version <version> under $RUN/gocache"
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
