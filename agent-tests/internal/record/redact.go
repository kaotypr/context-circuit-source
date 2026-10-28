package record

import (
	"regexp"
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
	s = replacePathPrefix(s, home, "$HOME")
	return redactUserElement(s, user)
}

// RedactRunRoot replaces the run-root prefix with $RUN.
func RedactRunRoot(s, runRoot string) string {
	return replacePathPrefix(s, runRoot, "$RUN")
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
