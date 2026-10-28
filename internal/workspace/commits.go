package workspace

import (
	"context"
	"errors"
	"os/exec"
	"slices"
	"strings"
)

// UncommittedFile is one record or ID file Git reports as not committed, named
// by its slash-separated path relative to the workspace root.
type UncommittedFile struct {
	Path   string `yaml:"path" json:"path"`
	Change string `yaml:"change" json:"change"`
}

// UncommittedRecords lists the files of the named records, wherever they live,
// every ID file, and every record an archive rename in this run moved, that
// Git reports as not committed. A record left uncommitted is invisible to
// every other clone, and an ID file left uncommitted is how two people come to
// reserve the same number.
//
// It never fails: it runs after a command has already succeeded, and Git's
// absence must not turn that success into a failure. Whatever kept Git from
// answering comes back as the unchecked reason instead.
func (s *Store) UncommittedRecords(ctx context.Context, ids ...string) (files []UncommittedFile, unchecked string) {
	files = []UncommittedFile{}
	// Git computes the prefix of a workspace nested inside another repository,
	// so symlinked paths and Windows short names never enter the arithmetic.
	prefix, err := Git(ctx, s.Root, "rev-parse", "--show-prefix")
	if err != nil {
		return files, gitUnavailable(err)
	}
	// --no-optional-locks keeps a read from taking the index lock out from
	// under the person's own Git use.
	text, err := Git(ctx, s.Root, "--no-optional-locks", "status", "--porcelain=v1", "-z", "--untracked-files=all", "--", ".")
	if err != nil {
		return files, gitUnavailable(err)
	}
	var matchers []func(string) bool
	for _, id := range ids {
		matchers = append(matchers, placement.GitPathMatcher(id))
	}
	concerns := func(rel string) bool {
		if rel == "" {
			return false
		}
		if IsIDFile(rel) {
			return true
		}
		for _, folder := range s.archivesMoved {
			if strings.HasPrefix(rel, folder+"/") {
				return true
			}
		}
		for _, matches := range matchers {
			if matches(rel) {
				return true
			}
		}
		return false
	}
	inside := func(path string) string {
		rel, ok := strings.CutPrefix(path, prefix)
		if !ok {
			return ""
		}
		return rel
	}
	changes := map[string]string{}
	for _, entry := range ParseStatus(text) {
		path := inside(entry.Path)
		if entry.From != "" {
			if concerns(path) {
				changes[path] = "new"
			}
			// A copy leaves its source in place, so only a rename's old path
			// is gone.
			if from := inside(entry.From); strings.Contains(entry.Code, "R") && concerns(from) {
				changes[from] = "deleted"
			}
			continue
		}
		if concerns(path) {
			changes[path] = statusChange(entry.Code)
		}
	}
	for path, change := range changes {
		files = append(files, UncommittedFile{Path: path, Change: change})
	}
	slices.SortFunc(files, func(a, b UncommittedFile) int { return strings.Compare(a.Path, b.Path) })
	return files, ""
}

// statusChange reads a two-letter porcelain code as new, modified, or deleted.
// An unresolved conflict is a modification whatever its letters, since the file
// exists on at least one side.
func statusChange(code string) string {
	switch {
	case strings.Contains(code, "U") || code == "AA" || code == "DD":
		return "modified"
	case code == "??" || code[0] == 'A':
		return "new"
	case strings.Contains(code, "D"):
		return "deleted"
	default:
		return "modified"
	}
}

func gitUnavailable(err error) string {
	if errors.Is(err, exec.ErrNotFound) {
		return "Git is not available"
	}
	if strings.Contains(strings.ToLower(err.Error()), "not a git repository") {
		return "the workspace is not in a Git repository"
	}
	return err.Error()
}
