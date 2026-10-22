package cli

import (
	"strings"
	"testing"
)

func TestReleaseLine(t *testing.T) {
	cases := []struct {
		version, kind, label string
	}{
		{"2.5.0", "stable", ""},
		{"0.0.0", "stable", ""},
		{"10.20.30", "stable", ""},
		{"2.5.0-acme.1", "variant", "acme"},
		{"2.5.0-sicepat.3", "variant", "sicepat"},
		{"2.5.0-abcdefghijklmnopqrst.12", "variant", "abcdefghijklmnopqrst"},
		{"2.5.0-abcdefghijklmnopqrstu.1", "unknown", ""},
		{"2.5.0-Acme.1", "unknown", ""},
		{"2.5.0-acme2.1", "unknown", ""},
		{"2.5.0-ac-me.1", "unknown", ""},
		{"2.5.0-acme.0", "unknown", ""},
		{"2.5.0-acme.01", "unknown", ""},
		{"2.5.0-acme.1.2", "unknown", ""},
		{"2.5.0-acme", "unknown", ""},
		{"2.5.0-.1", "unknown", ""},
		{"02.5.0", "unknown", ""},
		{"2.05.0-acme.1", "unknown", ""},
		{"2.5.00", "unknown", ""},
		{"2.0.0-dev", "unknown", ""},
		{"2.6.0-rc.1", "unknown", ""},
		{"2.5.0+x", "unknown", ""},
		{"2.5", "unknown", ""},
		{"", "unknown", ""},
	}
	for _, word := range strings.Fields(reservedLabels) {
		cases = append(cases, struct{ version, kind, label string }{"2.5.0-" + word + ".1", "unknown", ""})
	}
	for _, c := range cases {
		kind, label := releaseLine(c.version)
		if kind != c.kind || label != c.label {
			t.Errorf("releaseLine(%q) = %q, %q; want %q, %q", c.version, kind, label, c.kind, c.label)
		}
	}
}

func TestCommandName(t *testing.T) {
	for version, want := range map[string]string{
		"2.5.0":        "context-circuit-cli",
		"2.5.0-acme.1": "acme-context-circuit-cli",
		"2.5.0-rc.1":   "context-circuit-cli",
		"2.5.0-Acme.1": "context-circuit-cli",
		"2.0.0-test":   "context-circuit-cli",
	} {
		if got := commandName(version); got != want {
			t.Errorf("commandName(%q) = %q; want %q", version, got, want)
		}
	}
}
