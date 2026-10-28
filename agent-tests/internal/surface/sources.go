package surface

import (
	"path"
	"sort"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/triggers"
)

// SourceInput is what SelectSources needs besides the trigger table.
// ManifestBase and ManifestBranch map a workspace path to a source path.
type SourceInput struct {
	GitPaths       []string
	Files          []FileChange
	ManifestBase   map[string]string
	ManifestBranch map[string]string
	OutputsChanged bool
}

// SourceResult is the attribution of a non-empty surface.
type SourceResult struct {
	Paths         []string
	Ignored       []string
	EveryScenario bool
	MovedShipped  bool
	DestChanged   bool
}

// SelectSources decides which changed source paths explain the surface.
//
// A workspace path is tied to a source through the release manifest of
// either ref, plus the README exception in WithReadme. The source is
// kept when it also appears in git diff --name-status. Paths under
// internal/ are kept when CLI output changed, or when that file writes
// the role definitions, the record skeletons, or the workspace README
// and those bytes changed. scripts/release-manifest.txt and assets.go
// are kept when a destination was added, removed, or renamed, or a
// shipped file moved. Paths on the none list are ignored.
func SelectSources(table triggers.Table, in SourceInput) SourceResult {
	destChanged := !sameStringSet(destSet(in.ManifestBase), destSet(in.ManifestBranch))
	var result SourceResult
	result.DestChanged = destChanged
	explained := map[string]bool{}
	var unexplained []FileChange
	for _, file := range in.Files {
		src, ok := sourceOf(file, in.ManifestBase, in.ManifestBranch)
		if ok && containsPath(in.GitPaths, src) {
			explained[src] = true
			continue
		}
		unexplained = append(unexplained, file)
	}
	roleChanged := false
	recordChanged := false
	readmeChanged := false
	movedShipped := false
	for _, file := range in.Files {
		if isRolePath(file.Path) || (file.From != "" && isRolePath(file.From)) {
			roleChanged = true
		}
		if isGeneratedRecord(file.Path) || (file.From != "" && isGeneratedRecord(file.From)) {
			recordChanged = true
		}
		if file.Path == "README.md" || file.From == "README.md" {
			readmeChanged = true
		}
		if file.Status == "moved" && shippedPath(file.Path) && shippedPath(file.From) {
			movedShipped = true
		}
	}
	result.MovedShipped = movedShipped
	result.EveryScenario = destChanged || movedShipped

	for _, gitPath := range in.GitPaths {
		gitPath = path.Clean(gitPath)
		if gitPath == "." {
			continue
		}
		if table.NoScenarios(gitPath) {
			result.Ignored = append(result.Ignored, gitPath)
			continue
		}
		switch {
		case gitPath == "assets.go" || gitPath == "scripts/release-manifest.txt":
			if destChanged || movedShipped {
				result.Paths = append(result.Paths, gitPath)
			}
		case strings.HasPrefix(gitPath, "internal/"):
			if in.OutputsChanged ||
				(gitPath == "internal/workspace/agents.go" && roleChanged) ||
				(gitPath == "internal/workspace/records.go" && recordChanged) ||
				(gitPath == "internal/workspace/readme.go" && readmeChanged) {
				result.Paths = append(result.Paths, gitPath)
			}
		default:
			if explained[gitPath] {
				result.Paths = append(result.Paths, gitPath)
			}
		}
	}
	// A generated file can change with no matching git path, for
	// example role text, when the writer itself changed. The switch
	// above already kept that writer. Unexplained shipped paths whose
	// source is not in the diff stay out of Paths; the file list still
	// shows them.
	_ = unexplained
	result.Paths = uniquePaths(result.Paths)
	result.Ignored = uniquePaths(result.Ignored)
	return result
}

func sourceOf(file FileChange, base, branch map[string]string) (string, bool) {
	if src, ok := branch[file.Path]; ok {
		return src, true
	}
	if src, ok := base[file.Path]; ok {
		return src, true
	}
	if file.From != "" {
		if src, ok := base[file.From]; ok {
			return src, true
		}
		if src, ok := branch[file.From]; ok {
			return src, true
		}
	}
	return "", false
}

func destSet(manifest map[string]string) map[string]bool {
	out := map[string]bool{}
	for dest := range manifest {
		if dest == "README.md" && manifest[dest] == "product/workspace-README.md" {
			// The front page is added by WithReadme on both sides, so it
			// is not a manifest destination.
			continue
		}
		out[dest] = true
	}
	return out
}

func sameStringSet(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for key := range a {
		if !b[key] {
			return false
		}
	}
	return true
}

func containsPath(paths []string, want string) bool {
	for _, path := range paths {
		if path == want {
			return true
		}
	}
	return false
}

func uniquePaths(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range in {
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}

func isRolePath(p string) bool {
	return strings.HasPrefix(p, ".claude/agents/") ||
		strings.HasPrefix(p, ".cursor/agents/") ||
		strings.HasPrefix(p, ".codex/agents/") ||
		strings.HasPrefix(p, ".context-circuit/local/agents-")
}

func isGeneratedRecord(p string) bool {
	base := path.Base(p)
	if strings.HasPrefix(p, "intent/") && strings.HasPrefix(base, "i") && strings.Contains(base, "-") && strings.HasSuffix(base, ".md") {
		return true
	}
	return strings.HasPrefix(p, "plans/") && strings.HasPrefix(path.Base(path.Dir(p)), "p")
}

func shippedPath(p string) bool {
	if p == "" || isGeneratedRecord(p) || isIDLayout(p) || isRolePath(p) {
		return false
	}
	return true
}

func isIDLayout(p string) bool {
	return p == "ids.yaml" || p == ".context-circuit/ids.yaml" || strings.HasPrefix(p, ".context-circuit/ids/")
}
