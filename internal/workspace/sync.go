package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// SyncState is one repository's outcome from repo sync or workspace sync: where
// its base branch was, where it is now, and, when nothing moved, why. Every
// status other than updated and current names what was left alone, so one
// repository's state never hides another's.
type SyncState struct {
	ID          string   `json:"id"`
	Path        string   `json:"path,omitempty"`
	BaseBranch  string   `json:"base_branch,omitempty"`
	Counterpart string   `json:"counterpart,omitempty"`
	Status      string   `json:"status"`
	From        string   `json:"from,omitempty"`
	To          string   `json:"to,omitempty"`
	Commits     int      `json:"commits,omitempty"`
	Detail      string   `json:"detail,omitempty"`
	Changes     []string `json:"changes,omitempty"`
}

// VersionChange names the CLI version a workspace pinned before and after it
// moved, so the person hears that the launcher now wants another release.
type VersionChange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// WorkspaceSyncResult is workspace sync's answer: the workspace repository's
// state, and the next steps that apply once it moved.
type WorkspaceSyncResult struct {
	SyncState
	CheckRecommended  string         `json:"check_recommended,omitempty"`
	CLIVersionChanged *VersionChange `json:"cli_version_changed,omitempty"`
}

// RepoSync fast-forwards the base branch of each named repository, or of every
// repository bound on this machine when none is named. An unknown ID refuses the
// whole call before anything is fetched; a repository described but not bound
// here is reported as unbound and the rest go on.
func (s *Store) RepoSync(ctx context.Context, ids []string) ([]SyncState, error) {
	cfg, err := s.Config()
	if err != nil {
		return nil, err
	}
	bindings, err := s.Bindings()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		// A binding whose repository a teammate removed from workspace.yaml is
		// left out rather than refusing the whole call; naming it still refuses.
		for id := range bindings.Bindings {
			if _, ok := cfg.Repositories[id]; ok {
				ids = append(ids, id)
			}
		}
		slices.Sort(ids)
	}
	seen := map[string]bool{}
	var selected []string
	for _, id := range ids {
		if _, ok := cfg.Repositories[id]; !ok {
			return nil, fmt.Errorf("unknown repository: %s; nothing was fetched", id)
		}
		if !seen[id] {
			seen[id] = true
			selected = append(selected, id)
		}
	}
	states := make([]SyncState, 0, len(selected))
	for _, id := range selected {
		state := SyncState{ID: id}
		checkout, err := s.Repository(ctx, id)
		if err != nil {
			state.Status, state.Detail = "unbound", err.Error()+"; connect it with repo connect"
			states = append(states, state)
			continue
		}
		state.Path, state.BaseBranch = checkout.Path, checkout.BaseBranch
		fastForwardBase(ctx, checkout.Path, &state)
		states = append(states, state)
	}
	return states, nil
}

// WorkspaceSync fast-forwards the workspace repository's own base branch. It is
// stricter than RepoSync: any uncommitted path, tracked or untracked, leaves it
// untouched and unfetched, because those paths are records a person has not
// committed and a moved branch would land other people's records among them.
func (s *Store) WorkspaceSync(ctx context.Context) (WorkspaceSyncResult, error) {
	result := WorkspaceSyncResult{SyncState: SyncState{ID: "workspace"}}
	checkout, described, err := s.WorkspaceCheckout(ctx)
	if !described {
		if err != nil {
			return result, err
		}
		result.Status = "unbound"
		result.Detail = "the workspace's own repository is not described; describe it with workspace connect"
		return result, nil
	}
	if err != nil {
		result.Status, result.Detail = "unbound", err.Error()
		return result, nil
	}
	result.Path, result.BaseBranch = checkout.Path, checkout.BaseBranch
	snapshot, err := Inspect(ctx, checkout.Path)
	if err != nil {
		result.Status, result.Detail = "unreachable", err.Error()
		return result, nil
	}
	if len(snapshot.Changes) > 0 {
		result.Status = "dirty"
		result.Changes = snapshot.Changes
		result.Detail = fmt.Sprintf("%d uncommitted path(s) in the workspace repository; nothing was fetched or changed. Commit or remove them, then sync again", len(snapshot.Changes))
		return result, nil
	}
	before := s.pinnedCLI()
	fastForwardBase(ctx, checkout.Path, &result.SyncState)
	if result.Status != "updated" {
		return result, nil
	}
	result.CheckRecommended = "the workspace moved: run check and report what came in, and any finding the pull brought"
	if after := s.pinnedCLI(); after != before && after != "" {
		result.CLIVersionChanged = &VersionChange{From: before, To: after}
	}
	return result, nil
}

// pinnedCLI reads the version the workspace pins, or "" when it cannot.
func (s *Store) pinnedCLI() string {
	data, err := os.ReadFile(filepath.Join(s.Root, ".context-circuit", "CLI_VERSION"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// fastForwardBase fetches a checkout's base branch's remote and moves the local
// branch to its counterpart, fast-forward only. It never switches the branch a
// checkout has checked out and never merges, rebases, resets, stashes, or
// discards: anything it cannot do as a pure fast-forward is reported with the
// state left exactly as it was.
func fastForwardBase(ctx context.Context, path string, state *SyncState) {
	branch := state.BaseBranch
	if branch == "" || strings.HasPrefix(branch, "-") || strings.ContainsAny(branch, " \t\n\r") {
		state.Status, state.Detail = "no-local-branch", fmt.Sprintf("no usable base branch is recorded (%q); set one with repo base or workspace base", branch)
		return
	}
	from, err := Git(ctx, path, "rev-parse", "--verify", "--quiet", "--end-of-options", "refs/heads/"+branch+"^{commit}")
	if err != nil || from == "" {
		state.Status, state.Detail = "no-local-branch", fmt.Sprintf("this checkout has no local branch %q; nothing was changed", branch)
		return
	}
	state.From = from

	remote, _ := Git(ctx, path, "config", "--get", "branch."+branch+".remote")
	if remote == "" {
		remote = "origin"
	}
	// An upstream that is another local branch (remote = .) has nothing to fetch.
	if remote != "." {
		if strings.HasPrefix(remote, "-") || strings.ContainsAny(remote, " \t\n\r\x00") {
			state.Status, state.Detail = "no-upstream", fmt.Sprintf("branch %q names remote %q, which cannot be fetched safely; nothing was changed", branch, remote)
			return
		}
		if _, err := Git(ctx, path, "remote", "get-url", "--", remote); err != nil {
			state.Status, state.Detail = "no-upstream", fmt.Sprintf("the checkout has no remote %q; nothing was changed. Record one with repo remote or workspace remote", remote)
			return
		}
		if _, err := Git(ctx, path, "fetch", "--prune", "--quiet", "--", remote); err != nil {
			state.Status, state.Detail = "unreachable", fmt.Sprintf("fetching %s failed; nothing was changed: %v", remote, err)
			return
		}
	}

	ahead, behind, counterpart := aheadBehind(ctx, path, branch)
	if counterpart == "" {
		state.Status, state.Detail = "no-upstream", fmt.Sprintf("%q has no upstream and %s holds no %s/%s; nothing was changed", branch, remote, remote, branch)
		return
	}
	state.Counterpart = counterpart
	to, err := Git(ctx, path, "rev-parse", "--verify", "--end-of-options", counterpart+"^{commit}")
	if err != nil {
		state.Status, state.Detail = "no-upstream", fmt.Sprintf("%s cannot be read; nothing was changed", counterpart)
		return
	}
	switch {
	case ahead > 0 && behind > 0:
		state.Status = "diverged"
		state.Detail = fmt.Sprintf("%s holds %d commit(s) %s lacks, and %s holds %d it lacks; nothing was changed or discarded. Settle it by hand", branch, ahead, counterpart, counterpart, behind)
		return
	case ahead > 0:
		state.Status, state.To = "ahead", from
		state.Detail = fmt.Sprintf("%s holds %d commit(s) not on %s and nothing to take; deliver them through the usual push", branch, ahead, counterpart)
		return
	case behind == 0:
		state.Status, state.To = "current", from
		return
	}

	here, elsewhere := checkedOut(ctx, path, branch)
	switch {
	case here:
		if _, err := Git(ctx, path, "merge", "--ff-only", "--quiet", "--end-of-options", to); err != nil {
			state.Status = "blocked"
			if snapshot, e := Inspect(ctx, path); e == nil {
				state.Changes = snapshot.Changes
			}
			state.Detail = fmt.Sprintf("uncommitted changes would be overwritten by moving %s to %s; nothing was changed. Commit or move them, then sync again: %v", branch, counterpart, err)
			return
		}
	case elsewhere != "":
		state.Status = "checked-out-elsewhere"
		state.Detail = fmt.Sprintf("%s is checked out in the worktree at %s; nothing was changed. Sync from that worktree", branch, elsewhere)
		return
	default:
		// The old value guards the move: a branch changed since it was read
		// fails here rather than losing what changed it.
		if _, err := Git(ctx, path, "update-ref", "-m", "context-circuit: sync "+branch+" to "+counterpart, "refs/heads/"+branch, to, from); err != nil {
			state.Status, state.Detail = "blocked", fmt.Sprintf("%s changed while it was being moved; nothing was changed: %v", branch, err)
			return
		}
	}
	state.Status, state.To, state.Commits = "updated", to, behind
}

// checkedOut reports whether branch is checked out in the worktree at path, or
// else the path of another worktree of the same repository that has it.
func checkedOut(ctx context.Context, path, branch string) (bool, string) {
	list, err := Git(ctx, path, "worktree", "list", "--porcelain")
	if err != nil {
		return false, ""
	}
	self := canonical(path)
	var current string
	for line := range strings.SplitSeq(list, "\n") {
		if p, ok := strings.CutPrefix(line, "worktree "); ok {
			current = p
			continue
		}
		if ref, ok := strings.CutPrefix(line, "branch "); ok && ref == "refs/heads/"+branch {
			if canonical(current) == self {
				return true, ""
			}
			return false, current
		}
	}
	return false, ""
}

func canonical(path string) string {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
	return filepath.Clean(path)
}
