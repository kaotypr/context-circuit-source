package workspace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Worktree struct {
	Path        string `json:"path"`
	Branch      string `json:"branch"`
	Head        string `json:"head"`
	Main        bool   `json:"main"`
	Locked      bool   `json:"locked"`
	Prunable    bool   `json:"prunable"`
	Plan        string `json:"plan,omitempty"`
	StartCommit string `json:"start_commit,omitempty"`
	Reused      bool   `json:"reused,omitempty"`
}
type Association struct {
	Repository  string `yaml:"repository"`
	Plan        string `yaml:"plan,omitempty"`
	Path        string `yaml:"path"`
	Branch      string `yaml:"branch"`
	StartCommit string `yaml:"start_commit"`
}
type Associations struct {
	Worktrees []Association `yaml:"worktrees"`
}

func (s *Store) associations() (Associations, error) {
	result := Associations{[]Association{}}
	err := s.YAML(".context-circuit/local/worktrees.yaml", &result)
	if os.IsNotExist(err) {
		return result, nil
	}
	return result, err
}

func (s *Store) Worktrees(ctx context.Context, repoID string) ([]Worktree, error) {
	_, repo, err := s.Repository(ctx, repoID)
	if err != nil {
		return nil, err
	}
	text, err := Git(ctx, repo, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, err
	}
	list := []Worktree{}
	var current *Worktree
	for _, field := range strings.Split(text, "\x00") {
		if strings.HasPrefix(field, "worktree ") {
			// Git prints forward slashes on Windows. Return the same canonical
			// native path for a new worktree, listing, and later reuse.
			path := filepath.Clean(strings.TrimPrefix(field, "worktree "))
			if resolved, err := filepath.EvalSymlinks(path); err == nil {
				path = resolved
			}
			list = append(list, Worktree{Path: path, Main: len(list) == 0})
			current = &list[len(list)-1]
		}
		if current == nil {
			continue
		}
		switch {
		case strings.HasPrefix(field, "branch "):
			current.Branch = strings.TrimPrefix(field, "branch refs/heads/")
		case strings.HasPrefix(field, "HEAD "):
			current.Head = strings.TrimPrefix(field, "HEAD ")
		case field == "locked" || strings.HasPrefix(field, "locked "):
			current.Locked = true
		case field == "prunable" || strings.HasPrefix(field, "prunable "):
			current.Prunable = true
		}
	}
	local, err := s.associations()
	if err != nil {
		return nil, err
	}
	for i := range list {
		for _, assoc := range local.Worktrees {
			if assoc.Repository == repoID && filepath.Clean(assoc.Path) == filepath.Clean(list[i].Path) && assoc.Branch == list[i].Branch {
				list[i].Plan, list[i].StartCommit = assoc.Plan, assoc.StartCommit
			}
		}
	}
	return list, nil
}

func (s *Store) Prepare(ctx context.Context, repoID, plan, branch, start, destination string, reuse bool) (Worktree, error) {
	cfg, repo, err := s.Repository(ctx, repoID)
	if err != nil {
		return Worktree{}, err
	}
	if plan != "" {
		r, err := s.FindRecord(plan)
		if err != nil {
			return Worktree{}, err
		}
		if !strings.HasPrefix(plan, "p") || !slices.Contains(r.Repositories, repoID) {
			return Worktree{}, errors.New("plan does not include this repository")
		}
	}
	if branch == "" && plan != "" {
		branch = "cc/" + plan + "/" + repoID
	}
	if err := checkBranch(ctx, repo, branch); err != nil {
		return Worktree{}, err
	}
	if destination == "" {
		if plan == "" {
			return Worktree{}, errors.New("provide a plan or explicit destination")
		}
		destination = filepath.Join(".worktrees", plan, repoID)
	}
	path, err := s.LocalPath(destination)
	if err != nil {
		return Worktree{}, err
	}
	if resolved, e := filepath.EvalSymlinks(path); e == nil {
		path = resolved
	}
	list, err := s.Worktrees(ctx, repoID)
	if err != nil {
		return Worktree{}, err
	}
	for _, tree := range list {
		if filepath.Clean(tree.Path) == filepath.Clean(path) {
			if tree.Main || tree.Branch != branch || tree.Prunable {
				return Worktree{}, errors.New("destination is an existing different or unavailable worktree")
			}
			if (tree.Plan != plan || tree.StartCommit == "") && !reuse {
				return Worktree{}, errors.New("worktree has another association; use --reuse to select it explicitly")
			}
			if _, err := Inspect(ctx, tree.Path); err != nil {
				return Worktree{}, err
			}
			if tree.Plan != plan || tree.StartCommit == "" {
				tree.Plan = plan
				if tree.StartCommit == "" {
					tree.StartCommit = tree.Head
				}
				if err := s.saveAssociation(repoID, tree); err != nil {
					return tree, err
				}
			}
			tree.Reused = true
			return tree, nil
		}
		if tree.Branch == branch {
			return Worktree{}, fmt.Errorf("branch is already checked out at %s; select that worktree to resume it", tree.Path)
		}
	}
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		return Worktree{}, fmt.Errorf("destination already exists: %s", path)
	}
	var commit string
	if _, err := Git(ctx, repo, "show-ref", "--verify", "refs/heads/"+branch); err == nil {
		if !reuse {
			return Worktree{}, errors.New("branch already exists; use --reuse to preserve and check out its existing commits")
		}
		commit, err = Git(ctx, repo, "rev-parse", "--verify", "--end-of-options", "refs/heads/"+branch+"^{commit}")
		if err != nil {
			return Worktree{}, err
		}
		_, err = Git(ctx, repo, "worktree", "add", "--", path, branch)
		if err != nil {
			return Worktree{}, err
		}
	} else {
		if start == "" {
			start = "refs/heads/" + cfg.BaseBranch
			if _, err := Git(ctx, repo, "rev-parse", "--verify", "--end-of-options", start+"^{commit}"); err != nil {
				start = "refs/remotes/origin/" + cfg.BaseBranch
			}
		}
		commit, err = Git(ctx, repo, "rev-parse", "--verify", "--end-of-options", start+"^{commit}")
		if err != nil {
			return Worktree{}, fmt.Errorf("cannot resolve starting point %s (new repositories need a first commit): %w", start, err)
		}
		_, err = Git(ctx, repo, "worktree", "add", "-b", branch, "--", path, commit)
		if err != nil {
			return Worktree{}, fmt.Errorf("worktree preparation failed; inspect branch %s and path %s: %w", branch, path, err)
		}
	}
	actual, err := filepath.EvalSymlinks(path)
	if err != nil {
		return Worktree{}, err
	}
	result := Worktree{Path: actual, Branch: branch, Head: commit, Plan: plan, StartCommit: commit}
	if err := s.saveAssociation(repoID, result); err != nil {
		return result, fmt.Errorf("worktree exists at %s but local association could not be saved: %w", actual, err)
	}
	return result, nil
}

func (s *Store) saveAssociation(repo string, tree Worktree) error {
	local, err := s.associations()
	if err != nil {
		return err
	}
	kept := []Association{}
	for _, assoc := range local.Worktrees {
		if assoc.Repository != repo || assoc.Path != tree.Path {
			kept = append(kept, assoc)
		}
	}
	kept = append(kept, Association{repo, tree.Plan, tree.Path, tree.Branch, tree.StartCommit})
	return s.WriteYAML(".context-circuit/local/worktrees.yaml", Associations{kept}, 0600)
}

func (s *Store) InspectWorktree(ctx context.Context, repoID, input string) (Snapshot, error) {
	_, tree, err := s.selectedTree(ctx, repoID, input)
	if err != nil {
		return Snapshot{}, err
	}
	return Inspect(ctx, tree.Path)
}

func (s *Store) selectedTree(ctx context.Context, repoID, input string) (string, Worktree, error) {
	_, repo, err := s.Repository(ctx, repoID)
	if err != nil {
		return "", Worktree{}, err
	}
	path, err := s.LocalPath(input)
	if err != nil {
		return "", Worktree{}, err
	}
	if resolved, e := filepath.EvalSymlinks(path); e == nil {
		path = resolved
	}
	list, err := s.Worktrees(ctx, repoID)
	if err != nil {
		return "", Worktree{}, err
	}
	for _, tree := range list {
		if filepath.Clean(tree.Path) == filepath.Clean(path) {
			return repo, tree, nil
		}
	}
	return "", Worktree{}, errors.New("path is not a registered worktree of this repository")
}

func (s *Store) updateAssociation(repo, old, next string) error {
	local, err := s.associations()
	if err != nil {
		return err
	}
	updated := []Association{}
	for _, assoc := range local.Worktrees {
		if assoc.Repository == repo && assoc.Path == old {
			if next == "" {
				continue
			}
			assoc.Path = next
		}
		updated = append(updated, assoc)
	}
	return s.WriteYAML(".context-circuit/local/worktrees.yaml", Associations{updated}, 0600)
}

func (s *Store) MoveWorktree(ctx context.Context, repoID, input, destination string) error {
	repo, tree, err := s.selectedTree(ctx, repoID, input)
	if err != nil {
		return err
	}
	if tree.Main {
		return errors.New("the primary checkout is not moved by worktree move")
	}
	path, err := s.LocalPath(destination)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		return errors.New("destination already exists")
	}
	if _, err := Git(ctx, repo, "worktree", "move", "--", tree.Path, path); err != nil {
		return err
	}
	actual, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	if err := s.updateAssociation(repoID, tree.Path, actual); err != nil {
		return fmt.Errorf("moved worktree to %s; association update failed: %w", actual, err)
	}
	return nil
}

func (s *Store) RemoveWorktree(ctx context.Context, repoID, input string, discard bool) error {
	repo, tree, err := s.selectedTree(ctx, repoID, input)
	if err != nil {
		return err
	}
	if tree.Main || tree.Locked {
		return errors.New("primary or locked worktrees cannot be removed by this command")
	}
	changes, err := Git(ctx, tree.Path, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignored")
	if err != nil {
		return err
	}
	if changes != "" && !discard {
		return errors.New("worktree contains changes, untracked, or ignored files; preserve them or explicitly authorize --discard")
	}
	args := []string{"worktree", "remove"}
	if discard {
		args = append(args, "--force")
	}
	args = append(args, "--", tree.Path)
	if _, err := Git(ctx, repo, args...); err != nil {
		return err
	}
	if err := s.updateAssociation(repoID, tree.Path, ""); err != nil {
		return fmt.Errorf("worktree removed; local association update failed: %w", err)
	}
	return nil
}

func (s *Store) RepairWorktree(ctx context.Context, repoID, input string) error {
	_, repo, err := s.Repository(ctx, repoID)
	if err != nil {
		return err
	}
	path, err := s.LocalPath(input)
	if err != nil {
		return err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	if _, err := Git(ctx, repo, "worktree", "repair", "--", path); err != nil {
		return err
	}
	_, tree, err := s.selectedTree(ctx, repoID, path)
	if err != nil {
		return err
	}
	local, err := s.associations()
	if err != nil {
		return err
	}
	for _, assoc := range local.Worktrees {
		if assoc.Repository == repoID && assoc.Branch == tree.Branch && assoc.Path != tree.Path {
			return s.updateAssociation(repoID, assoc.Path, tree.Path)
		}
	}
	return nil
}
