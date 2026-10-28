package triggers

import (
	"fmt"
	"regexp"
	"strings"
)

// compileGlob compiles a path glob. * stays inside one path segment.
// ** crosses segments, and ** at the start may match an empty prefix.
// The pattern is matched against the full slash-separated path.
func compileGlob(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, fmt.Errorf("empty glob")
	}
	if strings.Contains(pattern, "***") {
		return nil, fmt.Errorf("glob %q is not valid", pattern)
	}
	var b strings.Builder
	b.WriteString("^")
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				i++
				if i+1 < len(pattern) && pattern[i+1] == '/' {
					i++
					b.WriteString("(?:.*/)?")
					continue
				}
				b.WriteString(".*")
				continue
			}
			b.WriteString("[^/]*")
		case '?':
			b.WriteString("[^/]")
		case '[':
			return nil, fmt.Errorf("glob %q: character classes are not supported", pattern)
		default:
			b.WriteString(regexp.QuoteMeta(pattern[i : i+1]))
		}
	}
	b.WriteString("$")
	re, err := regexp.Compile(b.String())
	if err != nil {
		return nil, fmt.Errorf("glob %q: %w", pattern, err)
	}
	return re, nil
}
