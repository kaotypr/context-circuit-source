package cli

import (
	"regexp"
	"slices"
	"strings"
)

// The release-line rule, kept identical to the copies in
// scripts/release-version.sh, scripts/build-cli.sh, and the cc-cli installers.
const (
	reservedLabels = "alpha beta rc dev test pre preview snapshot"
	labelPattern   = `^[a-z]{1,20}$`
)

// publicCommand is the PATH command every version outside a variant line
// installs as.
const publicCommand = "context-circuit-cli"

var (
	stableVersion = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
	variantLabel  = regexp.MustCompile(labelPattern)
	variantCount  = regexp.MustCompile(`^[1-9][0-9]*$`)
)

// releaseLine reports which line a version belongs to: stable for X.Y.Z,
// variant with its label for X.Y.Z-<label>.N, and unknown for anything else,
// developer builds and candidates included.
func releaseLine(version string) (kind, label string) {
	if stableVersion.MatchString(version) {
		return "stable", ""
	}
	base, rest, found := strings.Cut(version, "-")
	if !found || !stableVersion.MatchString(base) {
		return "unknown", ""
	}
	label, counter, found := strings.Cut(rest, ".")
	if !found || !variantLabel.MatchString(label) || !variantCount.MatchString(counter) ||
		slices.Contains(strings.Fields(reservedLabels), label) {
		return "unknown", ""
	}
	return "variant", label
}

// commandName is the PATH command the installers put this version under, so a
// variant build never shares a command with the public one.
func commandName(version string) string {
	if kind, label := releaseLine(version); kind == "variant" {
		return label + "-" + publicCommand
	}
	return publicCommand
}
