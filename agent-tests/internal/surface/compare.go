package surface

import (
	"crypto/sha256"
	"sort"
)

// FileChange is one workspace path that is not identical on both sides.
// Status is added, removed, changed, or moved. For a move, Path is the
// branch path and From is the base path. A removed file and an added
// file with the same bytes are one move.
type FileChange struct {
	Status string `json:"status"`
	Path   string `json:"path"`
	From   string `json:"from,omitempty"`
}

// CompareTrees diffs two path-to-content maps. Content should already
// be normalized so a fresh timestamp is not a change.
func CompareTrees(base, branch map[string][]byte) []FileChange {
	baseHash := digestTree(base)
	branchHash := digestTree(branch)
	var changed []FileChange
	removed := map[string][]string{}
	added := map[string][]string{}
	for path, sum := range baseHash {
		if other, ok := branchHash[path]; ok {
			if sum != other {
				changed = append(changed, FileChange{Status: "changed", Path: path})
			}
			continue
		}
		removed[sum] = append(removed[sum], path)
	}
	for path, sum := range branchHash {
		if _, ok := baseHash[path]; ok {
			continue
		}
		added[sum] = append(added[sum], path)
	}
	for sum, paths := range removed {
		sort.Strings(paths)
		extra := append([]string{}, added[sum]...)
		sort.Strings(extra)
		n := len(paths)
		if len(extra) < n {
			n = len(extra)
		}
		for i := 0; i < n; i++ {
			changed = append(changed, FileChange{Status: "moved", Path: extra[i], From: paths[i]})
		}
		for _, path := range paths[n:] {
			changed = append(changed, FileChange{Status: "removed", Path: path})
		}
		if n > 0 {
			added[sum] = extra[n:]
		}
	}
	for _, paths := range added {
		for _, path := range paths {
			changed = append(changed, FileChange{Status: "added", Path: path})
		}
	}
	sort.Slice(changed, func(i, j int) bool {
		if changed[i].Path != changed[j].Path {
			return changed[i].Path < changed[j].Path
		}
		return changed[i].Status < changed[j].Status
	})
	return changed
}

func digestTree(files map[string][]byte) map[string]string {
	out := make(map[string]string, len(files))
	for path, body := range files {
		sum := sha256.Sum256(body)
		out[path] = string(sum[:])
	}
	return out
}
