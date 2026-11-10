package report

import (
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// cleanPaths rewrites home, run-root, and the user name as a whole word.
// Key-shaped strings are left in place so a later check can refuse them.
// Auth files are removed by the zip builder before this text is scanned.
func cleanPaths(s, home, user string, roots, private []string) string {
	found := detectRunRoots(s)
	all := make([]string, 0, len(roots)+len(found))
	all = append(all, roots...)
	all = append(all, found...)
	all = uniqueLongestFirst(all)
	for _, root := range all {
		s = record.RedactRunRoot(s, root)
	}
	return record.RedactPrivate(record.RedactPaths(s, home, user), private)
}

// publicText cleans paths and replaces key-shaped strings.
// Design excerpts use it. The zip does not: a key there is a refusal.
func publicText(s, home, user string, roots, private []string) string {
	return record.Redact(cleanPaths(s, home, user, roots, private), "", "")
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
		start = skipEscape(s, start, j)
		// A path with a space in it, such as a volume name, is printed
		// with the space as written or escaped. Walk back over one space
		// at a time until the path reaches its leading slash.
		if s[start] != '/' {
			wide := start
			for wide > 1 && s[wide] != '/' && s[wide-1] == ' ' && isRunPath(s[wide-2]) {
				wide--
				for wide > 0 && isRunPath(s[wide-1]) {
					wide--
				}
				wide = skipEscape(s, wide, j)
			}
			if s[wide] == '/' {
				start = wide
			}
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

// skipEscape moves start past the last JSON or shell escape such as \n,
// \t or \" that the backward walk took in with the path: a path printed
// after a newline inside a JSON string begins after the escape, not at the
// word before it.
func skipEscape(s string, start, limit int) int {
	out := start
	for i := start; i < limit; i++ {
		if s[i] != '\\' {
			continue
		}
		k := i
		for k < limit && s[k] == '\\' {
			k++
		}
		if k < limit && strings.IndexByte("ntr\"", s[k]) >= 0 {
			out = k + 1
		}
		i = k
	}
	return out
}

func isRunToken(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

func isRunPath(b byte) bool {
	return isRunToken(b) || b == '/' || b == '\\' || b == '.' || b == '_' || b == '-' || b == '+' || b == '~' || b == '%'
}
