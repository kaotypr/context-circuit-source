package workspace

import (
	"fmt"
	"html"
	"net/url"
	"os"
	"strings"
)

// The shipped README is the product guide, and an initialized workspace keeps it
// as that workspace's own front page. Two spans of it are therefore generated
// rather than fixed, and each is delimited so the rewrite has an exact anchor
// instead of a guess about where a heading ends.
const (
	readmeTitleOpen   = "<!-- context-circuit:title -->"
	readmeTitleClose  = "<!-- /context-circuit:title -->"
	readmeBadgesOpen  = "<!-- context-circuit:badges -->"
	readmeBadgesClose = "<!-- /context-circuit:badges -->"
	productWebsite    = "https://context-circuit.kaotypr.com"
)

// TitleReadme names the workspace in its README heading and pins the version
// badges to what this workspace actually received.
//
// The shipped badges query the newest published release, which is the right
// answer on the product repository and the wrong one here: a workspace runs the
// versions recorded in .context-circuit, not whatever was published since. The
// title is generated for the same reason the name exists — a repository shared
// with a team should say whose workspace it is on its front page.
//
// A README missing its markers has been edited by whoever owns this workspace and
// is left exactly as it is. Initialization records data; rewriting somebody's
// prose to find an anchor is not a judgment it may make.
func (s *Store) TitleReadme(name string) error {
	data, err := s.Read("README.md")
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	text := string(data)
	// A Windows checkout converts the shipped seed to CRLF, and writing LF into
	// it would leave the file mixed. The markers hold no line break either way.
	nl := "\n"
	if strings.Contains(text, "\r\n") {
		nl = "\r\n"
	}
	badges, err := s.pinnedBadges(nl)
	if err != nil {
		return err
	}
	text, titled := replaceMarked(text, readmeTitleOpen, readmeTitleClose,
		strings.Join([]string{
			readmeTitleOpen,
			fmt.Sprintf("<h1 align=%q>Context Circuit - %s</h1>", "center", html.EscapeString(name)),
			readmeTitleClose,
		}, nl))
	text, badged := replaceMarked(text, readmeBadgesOpen, readmeBadgesClose, badges)
	if !titled && !badged {
		return nil
	}
	p, err := s.Path("README.md")
	if err != nil {
		return err
	}
	return os.WriteFile(p, []byte(text), 0644)
}

// replaceMarked swaps the whole span from open through close, reporting whether
// a complete pair was found. A lone opening marker leaves the text untouched,
// because a truncated replacement is worse than none.
func replaceMarked(text, open, closing, replacement string) (string, bool) {
	start := strings.Index(text, open)
	if start < 0 {
		return text, false
	}
	end := strings.Index(text[start:], closing)
	if end < 0 {
		return text, false
	}
	return text[:start] + replacement + text[start+end+len(closing):], true
}

func (s *Store) pinnedBadges(nl string) (string, error) {
	var rows []string
	for _, badge := range []struct{ alt, label, stamp string }{
		{"Workspace template version", "workspace", ".context-circuit/VERSION"},
		{"CLI version", "cli", ".context-circuit/CLI_VERSION"},
	} {
		version, err := s.stamp(badge.stamp)
		if err != nil {
			return "", err
		}
		// A version this workspace does not record is reported by omitting its
		// badge rather than by rendering an empty one.
		if version == "" {
			continue
		}
		rows = append(rows, fmt.Sprintf(
			"  <img alt=%q src=\"https://img.shields.io/static/v1?label=%s&message=v%s&color=1f6feb\">",
			badge.alt, url.QueryEscape(badge.label), url.QueryEscape(version)))
	}
	rows = append(rows, fmt.Sprintf(
		"  <a href=%q><img alt=\"Website\" src=\"https://img.shields.io/badge/website-context--circuit.kaotypr.com-0b7285\"></a>",
		productWebsite))
	return strings.Join(append(append(
		[]string{readmeBadgesOpen, `<p align="center">`}, rows...),
		"</p>", readmeBadgesClose), nl), nil
}

// stamp reads a single-line version file, treating an absent one as unrecorded.
func (s *Store) stamp(rel string) (string, error) {
	data, err := s.Read(rel)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
