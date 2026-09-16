package workspace

import (
	"fmt"
	"regexp"
	"strings"
)

// Readability findings report the shape of a note rather than its truth. A note
// that is accurate and unreadable still fails the reader it was written for, and
// shape is the part a diagnostic can see: where anchors sit, whether a paragraph
// is carrying a list it never made, whether a fence closes. Whether a sentence is
// any good is not mechanical, and stays with the skill that writes it.
const (
	proseRun     = 8
	longSentence = 60
	headlessNote = 40
)

// An anchor is a logical repository ID and a repository-relative path. Requiring
// a slash or a source file extension keeps an email address and a clone URL from
// reading as one; a root file with no extension is missed rather than guessed at,
// because a false report costs more than a quiet one here.
var noteAnchor = regexp.MustCompile(`(^|[^A-Za-z0-9_@.-])([a-z0-9][a-z0-9-]*@(?:[A-Za-z0-9_.<>-]+/[A-Za-z0-9_./<>-]*|[A-Za-z0-9_-]+\.(?:go|ts|tsx|js|jsx|py|rb|rs|java|kt|swift|php|cs|cpp|sql|sh|ps1|ya?ml|json|toml|md|txt|in)))`)

// The documented actor story shape: numbered once through the note, so the number
// stays a handle, and opening with the actor so the need is readable on its own.
var noteStory = regexp.MustCompile(`(?i)^\s*\d+\.\s+\*\*as an? `)

var listNumber = regexp.MustCompile(`^\d+[.)]\s`)

var sentenceEnd = regexp.MustCompile(`[.!?]["')\]]?\s+`)

// Mermaid renders nothing when it cannot recognize the first word, so an
// unknown type is a broken diagram rather than a matter of taste.
var mermaidDiagrams = map[string]bool{
	"flowchart": true, "graph": true, "sequenceDiagram": true,
	"stateDiagram": true, "stateDiagram-v2": true, "erDiagram": true,
	"classDiagram": true, "journey": true, "timeline": true,
	"gantt": true, "pie": true, "mindmap": true, "gitGraph": true,
}

// A structural line is anything that already gives the eye somewhere to stop.
func structuralLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return true
	}
	switch trimmed[0] {
	case '#', '-', '*', '+', '|', '>', '`':
		return true
	}
	return listNumber.MatchString(trimmed)
}

func noteReadability(relative string, data []byte) []string {
	var issues []string
	report := func(index int, format string, args ...any) {
		issues = append(issues, fmt.Sprintf("%s:%d: readability: %s", relative, index+1, fmt.Sprintf(format, args...)))
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	fenced, mermaid, fenceAt := false, false, 0
	owner, headings, title, thesis := false, 0, -1, true
	story, paragraph, paragraphAt := false, []string{}, 0

	flush := func() {
		if len(paragraph) > proseRun {
			report(paragraphAt, "a paragraph of %d lines; a list or a table may be hiding in it", len(paragraph))
		}
		for _, sentence := range sentenceEnd.Split(strings.Join(paragraph, " "), -1) {
			if words := len(strings.Fields(sentence)); words > longSentence {
				report(paragraphAt, "a sentence of %d words; a reader holds about half that", words)
			}
		}
		paragraph = paragraph[:0]
	}

	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			flush()
			if fenced {
				fenced, mermaid = false, false
				continue
			}
			fenced, fenceAt = true, index
			mermaid = strings.TrimSpace(strings.TrimLeft(trimmed, "`")) == "mermaid"
			continue
		}
		if fenced {
			if mermaid && trimmed != "" {
				mermaid = false
				if kind, _, _ := strings.Cut(trimmed, " "); !mermaidDiagrams[kind] {
					report(index, "mermaid diagram type not recognized (%s)", kind)
				}
			}
			continue
		}
		if trimmed == "Owner:" {
			owner = true
		}
		// A table cell already pairs a path with what it owns, which is the
		// property the Owner block exists to give; a sentence does not.
		if !owner && !strings.HasPrefix(trimmed, "|") {
			if match := noteAnchor.FindStringSubmatch(line); match != nil {
				report(index, "code anchor outside the Owner: block (%s)", match[2])
			}
		}
		if noteStory.MatchString(line) {
			story = true
		}
		if title < 0 && strings.HasPrefix(trimmed, "# ") {
			title, thesis = index, false
			continue
		}
		if !thesis && trimmed != "" {
			thesis = true
			if structuralLine(trimmed) {
				report(title, "no line under the title saying what the note is for")
			}
		}
		if strings.HasPrefix(trimmed, "## ") {
			headings++
		}
		if structuralLine(line) {
			flush()
			continue
		}
		if len(paragraph) == 0 {
			paragraphAt = index
		}
		paragraph = append(paragraph, trimmed)
	}
	flush()

	if fenced {
		report(fenceAt, "code fence is never closed")
	}
	if headings == 0 && len(lines) > headlessNote {
		issues = append(issues, fmt.Sprintf("%s: readability: %d lines and no section headings", relative, len(lines)))
	}
	if strings.Contains(relative, "/actors/") && !story {
		issues = append(issues, relative+": readability: an actor note carries numbered `As a ...` stories")
	}
	return issues
}
