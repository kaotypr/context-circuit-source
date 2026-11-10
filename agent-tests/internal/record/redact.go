package record

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var (
	reSK     = regexp.MustCompile(`sk-[A-Za-z0-9_-]{16,}`)
	reGHP    = regexp.MustCompile(`ghp_[A-Za-z0-9]{20,}`)
	reGHPat  = regexp.MustCompile(`github_pat_[A-Za-z0-9_]{20,}`)
	reAKIA   = regexp.MustCompile(`AKIA[0-9A-Z]{16}`)
	reBearer = regexp.MustCompile(`Bearer [A-Za-z0-9._~+/=-]{20,}`)
	reJWT    = regexp.MustCompile(`eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`)
	reRT     = regexp.MustCompile(`rt_[A-Za-z0-9._-]{20,}`)
)

// LooksLikeSecret reports whether s contains an API-key shaped value.
func LooksLikeSecret(s string) bool {
	return reSK.MatchString(s) ||
		reGHP.MatchString(s) ||
		reGHPat.MatchString(s) ||
		reAKIA.MatchString(s) ||
		reBearer.MatchString(s) ||
		reJWT.MatchString(s) ||
		reRT.MatchString(s)
}

// ShapeSecrets returns JWT-shaped and Codex refresh-token-shaped substrings.
func ShapeSecrets(s string) []string {
	var out []string
	out = append(out, reJWT.FindAllString(s, -1)...)
	out = append(out, reRT.FindAllString(s, -1)...)
	return out
}

// Redact replaces home with $HOME, a user path element with $USER, and
// key-shaped values with [redacted].
func Redact(s, home, user string) string {
	return redactSecrets(RedactPaths(s, home, user))
}

// RedactPaths replaces a home prefix with $HOME and the user name, as a
// whole word, with $USER. A longer word that contains the user name is
// left alone. Key-shaped strings are left unchanged.
func RedactPaths(s, home, user string) string {
	s = ReplacePath(s, home, "$HOME")
	return redactUserElement(s, user)
}

// RedactRunRoot replaces the run-root prefix with $RUN.
func RedactRunRoot(s, runRoot string) string {
	return ReplacePath(s, runRoot, "$RUN")
}

// ReplacePath replaces a path prefix with repl in every form a host or
// shell prints it: as written, with its spaces and other shell characters
// escaped by a backslash, URL-encoded, and with those backslashes escaped
// again by JSON, up to three levels deep.
func ReplacePath(s, prefix, repl string) string {
	for _, form := range PathForms(prefix) {
		s = replacePathPrefix(s, form, repl)
	}
	return s
}

// PathForms returns the forms of a path that ReplacePath recognizes,
// longest first, so a longer escaped form is replaced before a shorter one
// it contains.
func PathForms(p string) []string {
	p = strings.TrimRight(p, `/\`)
	if p == "" {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	add := func(f string) {
		if f != "" && !seen[f] {
			seen[f] = true
			out = append(out, f)
		}
	}
	base := []string{
		p,
		shellEscape(p),
		strings.ReplaceAll(p, " ", `\ `),
		strings.ReplaceAll(p, " ", "%20"),
		strings.ReplaceAll(strings.ReplaceAll(p, " ", "%20"), "~", "%7E"),
	}
	for _, f := range base {
		add(f)
		for level := 1; level <= 3; level++ {
			f = strings.ReplaceAll(f, `\`, `\\`)
			add(f)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return len(out[i]) > len(out[j]) })
	return out
}

// shellEscape puts a backslash before each character a shell would split
// or expand, as zsh and bash print a path they complete or quote.
func shellEscape(p string) string {
	var b strings.Builder
	for _, r := range p {
		if strings.ContainsRune(" ~'\"$`()&;|<>*?[]{}!#", r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func redactSecrets(s string) string {
	s = reSK.ReplaceAllString(s, "[redacted]")
	s = reGHP.ReplaceAllString(s, "[redacted]")
	s = reGHPat.ReplaceAllString(s, "[redacted]")
	s = reAKIA.ReplaceAllString(s, "[redacted]")
	s = reBearer.ReplaceAllString(s, "Bearer [redacted]")
	s = reJWT.ReplaceAllString(s, "[redacted]")
	s = reRT.ReplaceAllString(s, "[redacted]")
	return s
}

func redactUserElement(s, user string) string {
	if user == "" || s == "" {
		return s
	}
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(user) + `\b`)
	return re.ReplaceAllLiteralString(s, "$USER")
}

func replacePathPrefix(s, prefix, repl string) string {
	prefix = strings.TrimRight(prefix, `/\`)
	if prefix == "" || s == "" {
		return s
	}
	var b strings.Builder
	i := 0
	for i < len(s) {
		j := strings.Index(s[i:], prefix)
		if j < 0 {
			b.WriteString(s[i:])
			break
		}
		j += i
		end := j + len(prefix)
		if boundaryBefore(s, j) && boundaryAfter(s, end) {
			b.WriteString(s[i:j])
			b.WriteString(repl)
			i = end
			continue
		}
		b.WriteString(s[i:end])
		i = end
	}
	return b.String()
}

func boundaryBefore(s string, i int) bool {
	if i == 0 {
		return true
	}
	// A JSON escape such as \n or \t ends with a letter, but the path
	// after it starts a new word.
	if i >= 2 && s[i-2] == '\\' && strings.IndexByte("ntr", s[i-1]) >= 0 {
		return true
	}
	return !isPathChar(rune(s[i-1]))
}

func boundaryAfter(s string, i int) bool {
	if i >= len(s) {
		return true
	}
	return !isPathChar(rune(s[i]))
}

func isPathChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// RedactPrivate replaces each private entry: a path, in every form
// PathForms knows, with $PRIVATE, and a word, in any case, with
// [private]. The entries come from the person running the tests, such as
// the machine's temporary directory, the checkout, and names that must not
// be published.
func RedactPrivate(s string, private []string) string {
	for _, p := range longestFirst(private) {
		if isPathEntry(p) {
			// Longest first: the whole path, then each prefix of two or
			// more elements, so a path a host cut short is still replaced.
			for _, prefix := range pathPrefixes(strings.TrimRight(p, `/\`)) {
				s = ReplacePath(s, prefix, "$PRIVATE")
			}
			continue
		}
		s = privateWord(p).ReplaceAllLiteralString(s, "[private]")
	}
	return s
}

// PrivateLeak names the first private entry, home path or user name that
// is still in s after cleaning, or "" when none is.
func PrivateLeak(s, home, user string, private []string) string {
	entries := append([]string{home}, private...)
	for _, p := range entries {
		p = strings.TrimRight(p, `/\`)
		if p == "" {
			continue
		}
		if isPathEntry(p) {
			// A path cut short by an earlier replacement still names the
			// machine, so every prefix of two or more elements counts.
			for _, prefix := range pathPrefixes(p) {
				for _, form := range PathForms(prefix) {
					if containsPath(s, form) {
						return p
					}
				}
			}
			continue
		}
		if privateWord(p).MatchString(s) {
			return p
		}
	}
	if user != "" && regexp.MustCompile(`\b`+regexp.QuoteMeta(user)+`\b`).MatchString(s) {
		return user
	}
	return ""
}

// pathPrefixes returns p and each shorter prefix of it that still has two
// or more elements, such as /Volumes/Disk for /Volumes/Disk/tmp.
func pathPrefixes(p string) []string {
	out := []string{p}
	for {
		i := strings.LastIndexAny(p, `/\`)
		if i <= 0 {
			break
		}
		p = p[:i]
		if strings.Count(strings.Trim(p, `/\`), "/") < 1 {
			break
		}
		out = append(out, p)
	}
	return out
}

// containsPath reports whether form occurs in s ending at a path boundary.
func containsPath(s, form string) bool {
	for i := 0; i < len(s); {
		j := strings.Index(s[i:], form)
		if j < 0 {
			return false
		}
		if boundaryAfter(s, i+j+len(form)) {
			return true
		}
		i += j + 1
	}
	return false
}

func isPathEntry(p string) bool {
	return strings.ContainsAny(p, `/\`)
}

func privateWord(w string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)` + regexp.QuoteMeta(w))
}

// PrivateEntries returns the caller's TMPDIR, which holds every run root,
// this machine's host name, which the CLI writes into local files such as
// the extension trust file, and the extra paths and terms given, such as
// the source checkout and the terms of a deny file.
func PrivateEntries(caller map[string]string, extra ...string) []string {
	var out []string
	if h, err := os.Hostname(); err == nil {
		if h = strings.TrimSpace(h); len(h) >= 4 {
			out = append(out, h)
			if short, _, ok := strings.Cut(h, "."); ok && len(short) >= 4 {
				out = append(out, short)
			}
		}
	}
	if t := strings.TrimRight(caller["TMPDIR"], `/\`); t != "" {
		out = append(out, t)
	}
	for _, p := range extra {
		if p = strings.TrimRight(p, `/\`); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func longestFirst(in []string) []string {
	out := append([]string(nil), in...)
	sort.SliceStable(out, func(i, j int) bool { return len(out[i]) > len(out[j]) })
	return out
}
