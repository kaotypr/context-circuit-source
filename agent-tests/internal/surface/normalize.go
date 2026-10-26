package surface

import (
	"regexp"
	"sort"
	"strings"
)

var (
	reTime     = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})`)
	reDate     = regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}\b`)
	reIntentID = regexp.MustCompile(`\bi[0-9]{3,}\b`)
	rePlanID   = regexp.MustCompile(`\bp[0-9]{4,}\b`)
)

// Norm replaces values that differ between two fresh workspaces even
// when the template did not: timestamps, calendar dates, record ids,
// the built version, and absolute paths under the run.
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
	s = reTime.ReplaceAllString(s, "<time>")
	s = reDate.ReplaceAllString(s, "<date>")
	s = reIntentID.ReplaceAllString(s, "iNNN")
	s = rePlanID.ReplaceAllString(s, "pNNNN")
	if n.pin != nil {
		s = n.pin.ReplaceAllString(s, "<version>")
	}
	return s
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
