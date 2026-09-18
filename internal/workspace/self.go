package workspace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The workspace carries the shared records every member reads: intents, plans,
// the ID ledger, and the knowledge tree. That makes its own Git state part of
// what orientation has to answer — a member reading a plan from a branch three
// commits behind is reading a stale plan — and nothing reported it, because
// status and check only walk the repositories map.
//
// These operations describe the workspace's repository the way `repo connect`,
// `repo base`, and `repo remote` describe an ordinary one, and split the same
// way: where it lives and what it defaults to are shared, while the checkout
// root and the branch this machine works from are local.

// selfPath resolves the local binding's path to the Git root holding this
// workspace, defaulting to the workspace root when no path is recorded.
func (s *Store) selfPath(binding *Binding) (string, error) {
	local := "."
	if binding != nil && binding.Path != "" {
		local = binding.Path
	}
	path, err := s.LocalPath(local)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}

// WorkspaceCheckout describes the workspace's own repository as this machine
// sees it. It reports no checkout when nothing is recorded, so a caller can
// tell "not described" from "described and broken".
func (s *Store) WorkspaceCheckout(ctx context.Context) (Checkout, bool, error) {
	cfg, err := s.Config()
	if err != nil {
		return Checkout{}, false, err
	}
	if cfg.WorkspaceRepository == nil {
		return Checkout{}, false, nil
	}
	bindings, err := s.Bindings()
	if err != nil {
		return Checkout{}, false, err
	}
	if bindings.Workspace == nil {
		return Checkout{}, true, errors.New("connect this machine's workspace checkout with workspace connect")
	}
	path, err := s.selfPath(bindings.Workspace)
	if err != nil {
		return Checkout{}, true, err
	}
	root, err := rootOf(ctx, path)
	if err != nil {
		return Checkout{}, true, err
	}
	if root != path {
		return Checkout{}, true, errors.New("workspace binding must name its checkout root")
	}
	return Checkout{Repository: *cfg.WorkspaceRepository, Path: path, BaseBranch: bindings.WorkspaceBase(cfg)}, true, nil
}

// sameCheckout reports whether a path belongs to the Git repository already
// bound to some repository ID, and names it. Recording one checkout twice would
// let `repo base` and `workspace base` disagree about one machine's branch.
func (s *Store) sameCheckout(ctx context.Context, path string) (string, error) {
	common, err := Git(ctx, path, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", err
	}
	bindings, err := s.Bindings()
	if err != nil {
		return "", err
	}
	for id, binding := range bindings.Bindings {
		other, err := s.LocalPath(binding.Path)
		if err != nil {
			continue
		}
		if otherCommon, err := Git(ctx, other, "rev-parse", "--path-format=absolute", "--git-common-dir"); err == nil && otherCommon == common {
			return id, nil
		}
	}
	return "", nil
}

// ConnectWorkspace records the workspace's own repository and binds this
// machine's checkout of it. The shared half is written only when nothing is
// recorded yet: connecting a second machine describes that machine.
func (s *Store) ConnectWorkspace(ctx context.Context, input, base, remote, defaultBranch string) error {
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	if input == "" {
		input = "."
	}
	path, err := s.LocalPath(input)
	if err != nil {
		return err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	root, err := rootOf(ctx, path)
	if err != nil {
		return err
	}
	if root != path {
		// A workspace kept inside a larger repository has a root above it, and
		// saying so is more useful than repeating the rule.
		return fmt.Errorf("connect the root of the Git checkout holding this workspace; it is %s", root)
	}
	if !within(path, s.Root) {
		return errors.New("the workspace repository must contain the workspace")
	}
	if id, err := s.sameCheckout(ctx, path); err == nil && id != "" {
		return fmt.Errorf("this checkout is already registered as repository %s; remove that entry before describing it as the workspace", id)
	}
	if err := checkBranch(ctx, path, base); err != nil {
		return err
	}
	if cfg.WorkspaceRepository != nil {
		if remote != "" && remote != cfg.WorkspaceRepository.URL {
			return errors.New("the workspace repository already records another URL; use workspace remote to change it explicitly")
		}
		if defaultBranch != "" && defaultBranch != cfg.WorkspaceRepository.DefaultBranch {
			return errors.New("the workspace repository already records another default branch; use workspace remote to change it explicitly")
		}
	} else {
		if defaultBranch == "" {
			defaultBranch = base
		}
		if err := checkBranch(ctx, path, defaultBranch); err != nil {
			return err
		}
		if remote == "" {
			// Git already records where this checkout came from. A workspace
			// with no origin simply has no shared URL to record.
			remote, _ = Git(ctx, path, "remote", "get-url", "origin")
		}
		if remote != "" {
			if err := checkRemote(remote); err != nil {
				return err
			}
		}
		if err := s.Update("workspace.yaml", []string{"workspace_repository"}, Repository{remote, defaultBranch}, 0644); err != nil {
			return err
		}
	}
	local := "."
	if path != s.Root {
		relative, err := filepath.Rel(s.Root, path)
		if err != nil {
			return err
		}
		local = filepath.ToSlash(relative)
	}
	binding := Binding{local, base}
	if _, err := s.Read("repositories.local.yaml"); os.IsNotExist(err) {
		return s.WriteYAML("repositories.local.yaml", Bindings{Workspace: &binding, Bindings: map[string]Binding{}}, 0600)
	} else if err != nil {
		return err
	}
	return s.Update("repositories.local.yaml", []string{"workspace"}, binding, 0600)
}

// within reports whether root contains path, which is how the workspace's own
// repository is told from an unrelated checkout that happens to be named.
func within(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	if err != nil || filepath.IsAbs(relative) {
		return false
	}
	return relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

// SetWorkspaceBase records the branch workspace edits start from here. It is
// local for the reason a repository's base is: one machine may keep the
// workspace on a release branch while another stays on the default.
func (s *Store) SetWorkspaceBase(ctx context.Context, branch string) error {
	checkout, _, err := s.WorkspaceCheckout(ctx)
	if err != nil {
		return err
	}
	if checkout.Path == "" {
		return errors.New("describe the workspace repository first with workspace connect")
	}
	if err := checkBranch(ctx, checkout.Path, branch); err != nil {
		return err
	}
	return s.Update("repositories.local.yaml", []string{"workspace", "base_branch"}, branch, 0600)
}

// SetWorkspaceRemote changes what the workspace shares about its own
// repository. Neither value needs a checkout, so this works on a machine that
// has not connected one.
func (s *Store) SetWorkspaceRemote(ctx context.Context, remote, defaultBranch string) error {
	if remote == "" && defaultBranch == "" {
		return errors.New("provide --url or --default-branch")
	}
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	if cfg.WorkspaceRepository == nil {
		return errors.New("describe the workspace repository first with workspace connect")
	}
	if remote != "" {
		if err := checkRemote(remote); err != nil {
			return err
		}
		if err := s.Update("workspace.yaml", []string{"workspace_repository", "url"}, remote, 0644); err != nil {
			return err
		}
	}
	if defaultBranch == "" {
		return nil
	}
	if err := checkBranch(ctx, s.Root, defaultBranch); err != nil {
		return err
	}
	return s.Update("workspace.yaml", []string{"workspace_repository", "default_branch"}, defaultBranch, 0644)
}
