package surface

import (
	"regexp"
	"sort"
	"strings"
)

var (
	reYAMLID = regexp.MustCompile(`(?m)^(id:[ \t]*)"?(i[0-9]+|p[0-9]+)"?[ \t]*$`)
	reYAMLAt = regexp.MustCompile(`(?m)^((?:created_at|approved_at):[ \t]*)(?:"[^"\n]*"|\S+)[ \t]*$`)
	reJSONID = regexp.MustCompile(`"(id)"[ \t]*:[ \t]*"(i[0-9]+|p[0-9]+)"`)
	reJSONAt = regexp.MustCompile(`"(created_at|approved_at)"[ \t]*:[ \t]*"[^"]*"`)
)

// Norm replaces values the fixed surface script generates on a fresh
// workspace: the id, created_at, and approved_at fields it writes, temporary
// paths, and the CLI version built for that ref. Other ids, dates, and
// timestamps stay, so a help line that mentions an example id still differs.
type Norm struct {
	replacer *strings.Replacer
	pin      *regexp.Regexp
}

// NewNorm builds a normalizer. Longer paths are replaced first so a
// workspace path is not partly rewritten as the run root.
func NewNorm(paths []PathToken, pin string) Norm {
	type pair struct{ path, token string }
	var pairs []pair
	for _, item := range paths {
		if item.Path == "" {
			continue
		}
		pairs = append(pairs, pair{item.Path, item.Token})
		slash := strings.ReplaceAll(item.Path, "\\", "/")
		if slash != item.Path {
			pairs = append(pairs, pair{slash, item.Token})
		}
	}
	sort.Slice(pairs, func(i, j int) bool { return len(pairs[i].path) > len(pairs[j].path) })
	var flat []string
	seen := map[string]bool{}
	for _, pair := range pairs {
		if seen[pair.path] {
			continue
		}
		seen[pair.path] = true
		flat = append(flat, pair.path, pair.token)
	}
	n := Norm{replacer: strings.NewReplacer(flat...)}
	if pin != "" {
		n.pin = regexp.MustCompile(`\b` + regexp.QuoteMeta(pin) + `\b`)
	}
	return n
}

// PathToken is one absolute path and the placeholder that replaces it.
type PathToken struct {
	Path  string
	Token string
}

// Text normalizes s.
func (n Norm) Text(s string) string {
	if n.replacer != nil {
		s = n.replacer.Replace(s)
	}
	s = reYAMLID.ReplaceAllStringFunc(s, func(line string) string {
		m := reYAMLID.FindStringSubmatch(line)
		if m == nil {
			return line
		}
		return m[1] + idToken(m[2])
	})
	s = reYAMLAt.ReplaceAllString(s, "${1}<time>")
	s = reJSONID.ReplaceAllStringFunc(s, func(match string) string {
		m := reJSONID.FindStringSubmatch(match)
		if m == nil {
			return match
		}
		return `"` + m[1] + `":"` + idToken(m[2]) + `"`
	})
	s = reJSONAt.ReplaceAllString(s, `"${1}":"<time>"`)
	if n.pin != nil {
		s = n.pin.ReplaceAllString(s, "<version>")
	}
	return s
}

func idToken(id string) string {
	if strings.HasPrefix(id, "p") {
		return "pNNNN"
	}
	return "iNNN"
}

// Bytes normalizes text. A file with a NUL byte is returned unchanged.
func (n Norm) Bytes(b []byte) []byte {
	for _, c := range b {
		if c == 0 {
			return b
		}
	}
	return []byte(n.Text(string(b)))
}
