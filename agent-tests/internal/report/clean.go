package report

import (
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// cleanPaths rewrites home, run-root, and user path elements.
// Key-shaped strings are left in place so a later check can refuse them.
func cleanPaths(s, home, user string, roots []string) string {
	found := detectRunRoots(s)
	all := make([]string, 0, len(roots)+len(found))
	all = append(all, roots...)
	all = append(all, found...)
	all = uniqueLongestFirst(all)
	for _, root := range all {
		s = record.RedactRunRoot(s, root)
	}
	return record.RedactPaths(s, home, user)
}

// publicText cleans paths and replaces key-shaped strings.
// Design excerpts use it. The zip does not: a key there is a refusal.
func publicText(s, home, user string, roots []string) string {
	return record.Redact(cleanPaths(s, home, user, roots), "", "")
}

func uniqueLongestFirst(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		s = strings.TrimRight(s, `/\`)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && len(out[j]) > len(out[j-1]) {
			out[j], out[j-1] = out[j-1], out[j]
			j--
		}
	}
	return out
}

// detectRunRoots finds absolute paths whose last element is an agent-test
// directory, the run root the world package creates.
func detectRunRoots(s string) []string {
	const mark = "agent-test-"
	var out []string
	for i := 0; i < len(s); {
		j := strings.Index(s[i:], mark)
		if j < 0 {
			break
		}
		j += i
		end := j + len(mark)
		for end < len(s) && isRunToken(s[end]) {
			end++
		}
		start := j
		for start > 0 && isRunPath(s[start-1]) {
			start--
		}
		root := s[start:end]
		if strings.Contains(root, "/") || strings.Contains(root, `\`) {
			out = append(out, root)
		}
		if end == i {
			end++
		}
		i = end
	}
	return out
}

func isRunToken(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

func isRunPath(b byte) bool {
	return isRunToken(b) || b == '/' || b == '\\' || b == '.' || b == '_' || b == '-' || b == '+'
}
