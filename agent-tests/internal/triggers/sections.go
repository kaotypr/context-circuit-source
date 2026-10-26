package triggers

import (
	"fmt"
	"strings"
)

// Sections returns the level-2 headings touched by diff.
// diff is a unified diff of product/AGENTS.md.in. A heading is touched
// when its line is added or removed, when a changed line sits under it,
// or when a hunk header names it as @@ ## Heading @@.
// SectionDiff builds that form from the two full texts, which is what
// surface-diff passes: a short unified diff often starts mid-section and
// would not contain the heading.
func Sections(diff string) []string {
	var out []string
	seen := map[string]bool{}
	add := func(heading string) {
		heading = strings.TrimSpace(heading)
		if heading == "" || seen[heading] {
			return
		}
		seen[heading] = true
		out = append(out, heading)
	}
	current := ""
	for _, line := range strings.Split(diff, "\n") {
		if heading, ok := hunkHeading(line); ok {
			add(heading)
			current = heading
			continue
		}
		if line == "" || strings.HasPrefix(line, "diff ") || strings.HasPrefix(line, "index ") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
			continue
		}
		prefix, body, ok := diffBody(line)
		if !ok {
			continue
		}
		if heading, isHeading := level2(body); isHeading {
			if prefix == '+' || prefix == '-' {
				add(heading)
			}
			if prefix != '-' {
				current = heading
			}
			continue
		}
		if (prefix == '+' || prefix == '-') && current != "" {
			add(current)
		}
	}
	return out
}

// SectionDiff builds a diff whose hunk headers name each level-2 heading
// whose body or presence differs. Sections reads that header back.
func SectionDiff(before, after string) string {
	headings := changedHeadings(before, after)
	if len(headings) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("--- product/AGENTS.md.in\n+++ product/AGENTS.md.in\n")
	for _, heading := range headings {
		fmt.Fprintf(&b, "@@ ## %s @@\n", heading)
	}
	return b.String()
}

func changedHeadings(before, after string) []string {
	old := indexSections(splitSections(before))
	newSecs := splitSections(after)
	newIndex := indexSections(newSecs)
	var out []string
	seen := map[string]bool{}
	add := func(heading string) {
		if heading == "" || seen[heading] {
			return
		}
		seen[heading] = true
		out = append(out, heading)
	}
	for _, sec := range newSecs {
		if sec.heading == "" {
			continue
		}
		body, ok := old[sec.heading]
		if !ok || body != sec.body {
			add(sec.heading)
		}
	}
	for _, sec := range splitSections(before) {
		if _, ok := newIndex[sec.heading]; !ok {
			add(sec.heading)
		}
	}
	return out
}

type section struct {
	heading string
	body    string
}

func splitSections(text string) []section {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	var out []section
	var cur section
	var buf []string
	flush := func() {
		cur.body = strings.Join(buf, "\n")
		out = append(out, cur)
		buf = nil
	}
	for _, line := range strings.Split(text, "\n") {
		if heading, ok := level2(line); ok {
			flush()
			cur = section{heading: heading}
			continue
		}
		buf = append(buf, line)
	}
	flush()
	return out
}

func indexSections(sections []section) map[string]string {
	out := map[string]string{}
	for _, sec := range sections {
		if sec.heading == "" {
			continue
		}
		out[sec.heading] = sec.body
	}
	return out
}

func level2(line string) (string, bool) {
	if !strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "###") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, "## ")), true
}

func hunkHeading(line string) (string, bool) {
	if !strings.HasPrefix(line, "@@") || !strings.HasSuffix(line, "@@") {
		return "", false
	}
	inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "@@"), "@@"))
	return level2(inner)
}

func diffBody(line string) (prefix byte, body string, ok bool) {
	if line == "" {
		return 0, "", false
	}
	switch line[0] {
	case '+', '-', ' ':
		return line[0], line[1:], true
	default:
		return 0, "", false
	}
}
