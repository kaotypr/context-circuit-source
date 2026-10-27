package runner

import (
	"strings"
	"testing"
)

// A tool line longer than the cut, with the run root straddling the cut
// point, must come out with $RUN and no fragment of the real path.
func TestReadsLinesAreRedactedBeforeTheyAreShortened(t *testing.T) {
	root := "/private/var/folders/ab/cdefghij/T/agent-test-123456"
	line := "Tool: Bash " + strings.Repeat("x", 190) + " cat " + root + "/ws/README.md"
	red := textRedactor{runRoot: root}
	got := shortenLines(red.text(line), 240)
	if strings.Contains(got, "/private/") || strings.Contains(got, "/var/folders") {
		t.Fatalf("a fragment of the run root survived: %q", got)
	}
	if !strings.Contains(got, "$RUN") {
		t.Fatalf("run root was not replaced: %q", got)
	}
	if n := len([]rune(got)); n > 240 {
		t.Fatalf("line not shortened: %d runes", n)
	}
}
