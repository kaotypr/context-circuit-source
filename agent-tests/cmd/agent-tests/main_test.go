package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestUsageMentionsSurfaceDiff(t *testing.T) {
	var buf bytes.Buffer
	usage(&buf)
	text := buf.String()
	for _, part := range []string{"surface-diff", "--base", "--ref", "--checkout", "--out", "does not start a host", "needs no key"} {
		if !strings.Contains(text, part) {
			t.Errorf("usage missing %q\n%s", part, text)
		}
	}
}

func TestSurfaceDiffUsageErrors(t *testing.T) {
	if code := surfaceDiffCmd(nil, map[string]string{}); code != 2 {
		t.Fatalf("missing base %d", code)
	}
	if code := surfaceDiffCmd([]string{"--nope"}, map[string]string{}); code != 2 {
		t.Fatalf("unknown flag %d", code)
	}
	if code := surfaceDiffCmd([]string{"--base", "HEAD", "extra"}, map[string]string{}); code != 2 {
		t.Fatalf("extra arg %d", code)
	}
}
