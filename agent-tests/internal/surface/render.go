package surface

import (
	"fmt"
	"sort"
	"strings"
)

// emptyMessage is printed when the workspace, the CLI outputs, and the
// measured sizes are all equal.
const emptyMessage = "No agent tests needed: surface diff is empty"

// mappingNote is how a changed workspace path becomes a source path.
const mappingNote = "Source paths come from git diff --name-status between the two refs. A workspace path is tied to a source through scripts/release-manifest.txt. README.md is tied to product/workspace-README.md, which is written at init and is not a manifest line. Paths under internal/ are used when CLI output changed, or when that file writes role definitions, record skeletons, or the workspace README and those bytes changed. scripts/release-manifest.txt and assets.go are used when a destination changed or a shipped file moved. .git directories are not part of the tree. Comparison normalizes the id, created_at, and approved_at fields the fixed script writes, temporary paths, and the CLI version from that ref's archive. Other ids, dates, and timestamps stay. JSON key sets are taken from the raw stdout."

// OutputChange is one captured command. Diff is empty when the
// normalized texts match. Stdout is kept off the JSON file.
type OutputChange struct {
	Name         string `json:"name"`
	Changed      bool   `json:"changed"`
	Diff         string `json:"diff,omitempty"`
	StdoutBase   string `json:"-"`
	StdoutBranch string `json:"-"`
	Base         string `json:"-"`
	Branch       string `json:"-"`
}

// SizeChange is the byte and approximate-token size of one always-loaded
// file or one skill. Approximate tokens are bytes/4.
type SizeChange struct {
	Path         string `json:"path"`
	BaseBytes    int    `json:"base_bytes"`
	BranchBytes  int    `json:"branch_bytes"`
	DeltaBytes   int    `json:"delta_bytes"`
	BaseTokens   int    `json:"base_tokens_approx"`
	BranchTokens int    `json:"branch_tokens_approx"`
	DeltaTokens  int    `json:"delta_tokens_approx"`
}

// SourceHit is one attributed source path.
type SourceHit struct {
	Path      string   `json:"path"`
	Scenarios []string `json:"scenarios"`
	Hosts     []string `json:"hosts"`
	Kinds     []string `json:"kinds"`
	Notes     []string `json:"notes,omitempty"`
}

// Report is surface.json and the text of surface.md.
type Report struct {
	ID            string         `json:"id"`
	Base          string         `json:"base"`
	Branch        string         `json:"branch"`
	BaseRef       string         `json:"base_ref"`
	BranchRef     string         `json:"branch_ref"`
	Empty         bool           `json:"empty"`
	Message       string         `json:"message"`
	Files         []FileChange   `json:"files"`
	Outputs       []OutputChange `json:"outputs"`
	Sizes         []SizeChange   `json:"sizes"`
	Scenarios     []string       `json:"scenarios"`
	Hosts         []string       `json:"hosts"`
	Notes         []string       `json:"notes"`
	Kinds         []string       `json:"kinds"`
	Larger        bool           `json:"larger"`
	LargerReasons []Reason       `json:"larger_reasons"`
	NoRow         []string       `json:"no_row"`
	OrientHosts   []string       `json:"orient_hosts"`
	Ignored       []string       `json:"ignored"`
	Sources       []SourceHit    `json:"sources"`
	Mapping       string         `json:"mapping"`
	Dir           string         `json:"dir"`
}

// Render writes the printed report. Approximate tokens are bytes/4.
func Render(r Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Base: %s (%s)\nBranch: %s (%s)\n\n", r.BaseRef, short(r.Base), r.BranchRef, short(r.Branch))
	b.WriteString(r.Message)
	b.WriteString("\n")
	if r.Empty {
		b.WriteString("\n")
		b.WriteString(r.Mapping)
		b.WriteString("\n")
		return b.String()
	}
	b.WriteString("\n## Files\n\n")
	if len(r.Files) == 0 {
		b.WriteString("No workspace file changes.\n")
	}
	for _, file := range r.Files {
		switch file.Status {
		case "moved":
			fmt.Fprintf(&b, "- moved %s -> %s\n", file.From, file.Path)
		default:
			fmt.Fprintf(&b, "- %s %s\n", file.Status, file.Path)
		}
	}
	b.WriteString("\n## Outputs\n\n")
	if len(r.Outputs) == 0 {
		b.WriteString("No CLI output changes.\n")
	}
	for _, out := range r.Outputs {
		fmt.Fprintf(&b, "### %s\n\n```diff\n%s```\n\n", out.Name, out.Diff)
	}
	b.WriteString("## Size\n\nApproximate tokens are bytes/4.\n\n")
	b.WriteString("| File | Base bytes | Branch bytes | Delta bytes | Base tokens (approx) | Branch tokens (approx) | Delta tokens (approx) |\n")
	b.WriteString("|---|---:|---:|---:|---:|---:|---:|\n")
	for _, size := range r.Sizes {
		fmt.Fprintf(&b, "| `%s` | %d | %d | %d | %d | %d | %d |\n",
			size.Path, size.BaseBytes, size.BranchBytes, size.DeltaBytes, size.BaseTokens, size.BranchTokens, size.DeltaTokens)
	}
	b.WriteString("\n## Scenarios\n\n")
	b.WriteString(listOr(r.Scenarios, "none"))
	b.WriteString("\n\n## Hosts\n\n")
	if len(r.Hosts) == 0 {
		b.WriteString("any one host\n")
	} else {
		b.WriteString(strings.Join(r.Hosts, ", "))
		b.WriteString("\n")
	}
	if len(r.OrientHosts) > 0 {
		b.WriteString("orient on all three hosts: ")
		b.WriteString(strings.Join(r.OrientHosts, ", "))
		b.WriteString("\n")
	}
	if len(r.Notes) > 0 {
		b.WriteString("\n## Notes\n\n")
		for _, note := range r.Notes {
			fmt.Fprintf(&b, "- %s\n", note)
		}
	}
	b.WriteString("\n## Kinds\n\n")
	b.WriteString(listOr(r.Kinds, "none"))
	b.WriteString("\n\n## Larger\n\n")
	if r.Larger {
		b.WriteString("Looks larger. A person decides.\n")
		for _, reason := range r.LargerReasons {
			fmt.Fprintf(&b, "- %s [%s]: %s\n", reason.Kind, reasonTag(reason), reason.Detail)
		}
	} else {
		b.WriteString("Does not look larger.\n")
	}
	b.WriteString("\n## No row\n\n")
	if len(r.NoRow) == 0 {
		b.WriteString("Every attributed path has a trigger row.\n")
	} else {
		b.WriteString("No row: the author picks the closest scenario.\n")
		for _, path := range r.NoRow {
			fmt.Fprintf(&b, "- %s\n", path)
		}
	}
	if len(r.Ignored) > 0 {
		b.WriteString("\n## Ignored\n\n")
		b.WriteString("These changed paths are on the no-scenarios list.\n")
		for _, path := range r.Ignored {
			fmt.Fprintf(&b, "- %s\n", path)
		}
	}
	b.WriteString("\n## Mapping\n\n")
	b.WriteString(r.Mapping)
	b.WriteString("\n")
	return b.String()
}

func reasonTag(reason Reason) string {
	if reason.Exact {
		return "exact"
	}
	return "heuristic"
}

func listOr(items []string, empty string) string {
	if len(items) == 0 {
		return empty
	}
	return strings.Join(items, ", ")
}

func short(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

func summaryMessage(files []FileChange, outputs []OutputChange) string {
	nOut := 0
	for _, out := range outputs {
		if out.Changed {
			nOut++
		}
	}
	return fmt.Sprintf("Surface diff: %d workspace files and %d CLI outputs differ.", len(files), nOut)
}

func changedOutputs(outputs []OutputChange) []OutputChange {
	var out []OutputChange
	for _, item := range outputs {
		if item.Changed {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
