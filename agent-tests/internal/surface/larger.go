package surface

import (
	"regexp"
	"sort"
	"strings"
)

// Reason is one way the diff looks like a larger change.
// Kind is layout, skill, cli, gate, or migration.
type Reason struct {
	Kind   string `json:"kind"`
	Detail string `json:"detail"`
}

var gatePatterns = []struct {
	name string
	re   *regexp.Regexp
}{
	{"push", regexp.MustCompile(`(?i)\bpush\b`)},
	{"commit", regexp.MustCompile(`(?i)\bcommit\b`)},
	{"attribution", regexp.MustCompile(`(?i)\battribution\b`)},
	{"approval", regexp.MustCompile(`(?i)\bapprov`)},
	{"planning stop", regexp.MustCompile(`(?i)planning stop|planning_recommended`)},
	{"sources", regexp.MustCompile(`(?i)\bsources\b`)},
	{"borrowed", regexp.MustCompile(`(?i)\bborrowed\b`)},
}

// LayoutReasons reports a shipped file added, removed, or moved.
// Generated records and the ids directory are not layout: those are
// migration or record output.
func LayoutReasons(files []FileChange) []Reason {
	var paths []string
	for _, file := range files {
		if file.Status == "changed" {
			continue
		}
		if !shippedPath(file.Path) || (file.From != "" && !shippedPath(file.From)) {
			continue
		}
		if file.Status == "moved" {
			paths = append(paths, file.From+" -> "+file.Path)
			continue
		}
		paths = append(paths, file.Status+" "+file.Path)
	}
	if len(paths) == 0 {
		return nil
	}
	sort.Strings(paths)
	return []Reason{{Kind: "layout", Detail: "shipped file added, removed, or moved: " + strings.Join(paths, ", ")}}
}

// SkillReasons reports a skill file added or removed, or a change to
// its headings or frontmatter. A body-only edit is not a shape change.
func SkillReasons(base, branch map[string][]byte) []Reason {
	names := map[string]bool{}
	for _, tree := range []map[string][]byte{base, branch} {
		for path := range tree {
			if name, ok := skillName(path); ok {
				names[name] = true
			}
		}
	}
	var reasons []Reason
	for name := range names {
		rel := ".agents/skills/" + name + "/SKILL.md"
		_, inBase := base[rel]
		_, inBranch := branch[rel]
		switch {
		case inBase && !inBranch:
			reasons = append(reasons, Reason{Kind: "skill", Detail: "skill removed: " + name})
		case !inBase && inBranch:
			reasons = append(reasons, Reason{Kind: "skill", Detail: "skill added: " + name})
		default:
			bf, bh := skillShape(base[rel])
			ff, fh := skillShape(branch[rel])
			if bf != ff || strings.Join(bh, "\n") != strings.Join(fh, "\n") {
				reasons = append(reasons, Reason{Kind: "skill", Detail: "skill headings or frontmatter changed: " + name})
			}
		}
	}
	sort.Slice(reasons, func(i, j int) bool { return reasons[i].Detail < reasons[j].Detail })
	return reasons
}

func skillName(p string) (string, bool) {
	const prefix = ".agents/skills/"
	const suffix = "/SKILL.md"
	if !strings.HasPrefix(p, prefix) || !strings.HasSuffix(p, suffix) {
		return "", false
	}
	name := strings.TrimSuffix(strings.TrimPrefix(p, prefix), suffix)
	if name == "" || strings.Contains(name, "/") {
		return "", false
	}
	return name, true
}

func skillShape(b []byte) (front string, headings []string) {
	text := string(b)
	if strings.HasPrefix(text, "---\n") {
		rest := strings.TrimPrefix(text, "---\n")
		if i := strings.Index(rest, "\n---\n"); i >= 0 {
			front = rest[:i]
			text = rest[i+len("\n---\n"):]
		}
	}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#") {
			headings = append(headings, line)
		}
	}
	return front, headings
}

// CLIReasons reports help text, JSON keys, check findings, and record
// skeleton shape. Other output changes stay in the diff and are not
// called a CLI action.
func CLIReasons(outputs []OutputChange, base, branch map[string][]byte) []Reason {
	help, check, keys := false, false, false
	var keyNames []string
	for _, out := range outputs {
		if !out.Changed {
			continue
		}
		if strings.HasPrefix(out.Name, "help") {
			help = true
		}
		if out.Name == "check" || out.Name == "check-json" {
			check = true
		}
		if strings.HasSuffix(out.Name, "-json") && jsonKeysChanged(out.StdoutBase, out.StdoutBranch) {
			keys = true
			keyNames = append(keyNames, out.Name)
		}
	}
	var reasons []Reason
	if help {
		reasons = append(reasons, Reason{Kind: "cli", Detail: "help topics output changed"})
	}
	if keys {
		sort.Strings(keyNames)
		reasons = append(reasons, Reason{Kind: "cli", Detail: "JSON keys changed: " + strings.Join(keyNames, ", ")})
	}
	if check {
		reasons = append(reasons, Reason{Kind: "cli", Detail: "check findings changed"})
	}
	if skeletonChanged(base, branch) {
		reasons = append(reasons, Reason{Kind: "cli", Detail: "record skeleton headings or frontmatter changed"})
	}
	return reasons
}

func skeletonChanged(base, branch map[string][]byte) bool {
	paths := map[string]bool{}
	for _, tree := range []map[string][]byte{base, branch} {
		for path := range tree {
			if isGeneratedRecord(path) {
				paths[path] = true
			}
		}
	}
	for path := range paths {
		if recordShape(base[path]) != recordShape(branch[path]) {
			return true
		}
	}
	return false
}

func recordShape(b []byte) string {
	text := string(b)
	var keys, headings []string
	if strings.HasPrefix(text, "---\n") {
		rest := strings.TrimPrefix(text, "---\n")
		if i := strings.Index(rest, "\n---\n"); i >= 0 {
			for _, line := range strings.Split(rest[:i], "\n") {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "#") {
					continue
				}
				if name, _, ok := strings.Cut(line, ":"); ok && !strings.HasPrefix(line, " ") {
					keys = append(keys, name)
				}
			}
			text = rest[i+len("\n---\n"):]
		}
	}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "## ") {
			headings = append(headings, line)
		}
	}
	return strings.Join(keys, ",") + "\n" + strings.Join(headings, "\n")
}

// GateReasons is a keyword heuristic over added and removed lines.
// It is not a reading of the change: a line that mentions push, commit,
// attribution, approval, the planning stop, sources, or borrowed
// knowledge is reported, and a person decides whether it is a gate.
func GateReasons(lines []string) []Reason {
	found := map[string]bool{}
	var names []string
	for _, line := range lines {
		for _, pattern := range gatePatterns {
			if pattern.re.MatchString(line) && !found[pattern.name] {
				found[pattern.name] = true
				names = append(names, pattern.name)
			}
		}
	}
	if len(names) == 0 {
		return nil
	}
	return []Reason{{
		Kind:   "gate",
		Detail: "keyword heuristic on added and removed lines (push, commit, attribution, approval, planning stop, sources, borrowed knowledge); matched: " + strings.Join(names, ", "),
	}}
}

// ChangedLines returns lines that were added or removed, using a
// multiset so an unchanged line is not reported.
func ChangedLines(before, after string) []string {
	oldCount := lineCounts(before)
	newCount := lineCounts(after)
	var out []string
	for line, n := range newCount {
		if extra := n - oldCount[line]; extra > 0 {
			for i := 0; i < extra; i++ {
				out = append(out, line)
			}
		}
	}
	for line, n := range oldCount {
		if extra := n - newCount[line]; extra > 0 {
			for i := 0; i < extra; i++ {
				out = append(out, line)
			}
		}
	}
	sort.Strings(out)
	return out
}

func lineCounts(s string) map[string]int {
	out := map[string]int{}
	if s == "" {
		return out
	}
	for _, line := range strings.Split(s, "\n") {
		out[line]++
	}
	return out
}

// MigrationReasons reports a workspace.yaml version change or an ids
// path layout change. A VERSION or CLI_VERSION bump is not migration:
// those paths are on the none list and an existing workspace takes them
// with the template delta. The version field here is the schema number
// in workspace.yaml.
func MigrationReasons(base, branch map[string][]byte) []Reason {
	var reasons []Reason
	oldSchema := schemaField(base["workspace.yaml"])
	newSchema := schemaField(branch["workspace.yaml"])
	if oldSchema != "" && newSchema != "" && oldSchema != newSchema {
		reasons = append(reasons, Reason{Kind: "migration", Detail: "workspace.yaml version changed from " + oldSchema + " to " + newSchema})
	}
	if !sameStringSet(idPaths(base), idPaths(branch)) {
		reasons = append(reasons, Reason{Kind: "migration", Detail: "ids file layout changed"})
	}
	return reasons
}

func schemaField(b []byte) string {
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "version:") {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "version:")), `"'`)
		}
	}
	return ""
}

func idPaths(tree map[string][]byte) map[string]bool {
	out := map[string]bool{}
	for path := range tree {
		if isIDLayout(path) {
			out[path] = true
		}
	}
	return out
}

func jsonKeysChanged(before, after string) bool {
	b, bok := jsonKeyList(before)
	a, aok := jsonKeyList(after)
	if !bok || !aok {
		return false
	}
	return strings.Join(b, "\n") != strings.Join(a, "\n")
}
