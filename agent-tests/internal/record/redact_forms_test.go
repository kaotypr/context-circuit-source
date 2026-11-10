package record

import (
	"strings"
	"testing"
)

// A path with a space and a tilde is printed raw, shell-escaped, URL-encoded
// and JSON-escaped; every form is replaced.
func TestReplacePathForms(t *testing.T) {
	root := "/Volumes/My Disk~/tmp/agent-test-42"
	for _, in := range []string{
		`cd /Volumes/My Disk~/tmp/agent-test-42/ws`,
		`cd /Volumes/My\ Disk~/tmp/agent-test-42/ws`,
		`cd /Volumes/My\ Disk\~/tmp/agent-test-42/ws`,
		`"cd /Volumes/My\\ Disk\\~/tmp/agent-test-42/ws"`,
		`"cd /Volumes/My\\\\ Disk\\\\~/tmp/agent-test-42/ws"`,
		`[a.go](/Volumes/My%20Disk~/tmp/agent-test-42/ws/a.go)`,
	} {
		got := RedactRunRoot(in, root)
		if strings.Contains(got, "Disk") || !strings.Contains(got, "$RUN/ws") {
			t.Errorf("%s -> %s", in, got)
		}
	}
}

func TestRedactPrivate(t *testing.T) {
	private := PrivateEntries(map[string]string{"TMPDIR": "/Volumes/My Disk/tmp/"}, "acme-internal", "initech")
	in := "Update with /plugin update tool@acme-internal. INITECH notice. ls /Volumes/My\\ Disk/tmp/x"
	got := RedactPrivate(in, private)
	if leak := PrivateLeak(got, "", "", private); leak != "" {
		t.Fatalf("left %q in %s", leak, got)
	}
	if !strings.Contains(got, "$PRIVATE/x") || !strings.Contains(got, "tool@[private]") {
		t.Fatalf("got %s", got)
	}
	if PrivateLeak(in, "", "", private) == "" {
		t.Fatal("raw text reported no leak")
	}
}

// A host may shorten a tool line in the middle of a path; the private
// prefix that remains is still replaced.
func TestRedactPrivateCutPath(t *testing.T) {
	private := []string{"/Volumes/My Disk/tmp"}
	got := RedactPrivate(`rg --files /Volumes/My\\ Disk/Work/.cc-agent-te`, private)
	if leak := PrivateLeak(got, "", "", private); leak != "" || strings.Contains(got, "Disk") {
		t.Fatalf("got %s", got)
	}
}
