package workspace

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Git(ctx context.Context, path string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, "git", append([]string{"-C", path}, args...)...)
	// A coding host may itself be running inside a Git operation. Its checkout
	// overrides must not redirect commands away from the selected repository.
	for _, item := range os.Environ() {
		key, _, _ := strings.Cut(item, "=")
		switch strings.ToUpper(key) {
		case "GIT_DIR", "GIT_WORK_TREE", "GIT_COMMON_DIR", "GIT_INDEX_FILE", "GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES", "GIT_TERMINAL_PROMPT":
			continue
		}
		command.Env = append(command.Env, item)
	}
	command.Env = append(command.Env, "GIT_TERMINAL_PROMPT=0")
	data, err := command.Output()
	if err != nil {
		if failure, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(string(failure.Stderr)))
		}
		return "", fmt.Errorf("git %s: %w", args[0], err)
	}
	return strings.TrimRight(string(data), "\r\n"), nil
}

func (s *Store) LocalPath(path string) (string, error) {
	if path == "" {
		return "", errors.New("provide a local path")
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(s.Root, path)
	}
	return filepath.Abs(path)
}

func checkBranch(ctx context.Context, path, branch string) error {
	if branch == "" || strings.HasPrefix(branch, "-") {
		return errors.New("provide an explicit branch name")
	}
	checked, err := Git(ctx, path, "check-ref-format", "--branch", branch)
	if err != nil || checked != branch {
		return fmt.Errorf("invalid branch name: %s", branch)
	}
	return nil
}

// checkRemote rejects a remote a shared record must not carry: text that is not
// one usable line, and a URL embedding a password, which would put a secret in a
// committed file.
func checkRemote(remote string) error {
	if err := Text(remote); err != nil {
		return err
	}
	if u, err := url.Parse(remote); err == nil && u.User != nil {
		if _, password := u.User.Password(); password {
			return errors.New("use Git's credential helper instead of a URL containing a password")
		}
	}
	return nil
}

func rootOf(ctx context.Context, path string) (string, error) {
	root, err := Git(ctx, path, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(root)
}

// Checkout is one repository as this machine sees it: the shared description,
// the bound working path, and the branch work starts from here.
type Checkout struct {
	Repository
	Path       string
	BaseBranch string
}

func (s *Store) Repository(ctx context.Context, id string) (Checkout, error) {
	cfg, err := s.Config()
	if err != nil {
		return Checkout{}, err
	}
	repo, ok := cfg.Repositories[id]
	if !ok {
		return Checkout{}, fmt.Errorf("unknown repository: %s", id)
	}
	bindings, err := s.Bindings()
	if err != nil {
		return Checkout{}, err
	}
	binding, ok := bindings.Bindings[id]
	if !ok {
		return Checkout{}, fmt.Errorf("connect this machine's checkout for repository %s", id)
	}
	path, err := s.LocalPath(binding.Path)
	if err != nil {
		return Checkout{}, err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return Checkout{}, err
	}
	root, err := rootOf(ctx, path)
	if err != nil {
		return Checkout{}, err
	}
	if root != path {
		return Checkout{}, fmt.Errorf("repository binding must name its checkout root: %s", id)
	}
	return Checkout{Repository: repo, Path: path, BaseBranch: bindings.BaseBranch(cfg, id)}, nil
}

// Connect binds a repository ID to this machine's checkout and to the branch
// work starts from here. The shared record is written only when the ID is new:
// connecting a second machine describes that machine, never the repository.
func (s *Store) Connect(ctx context.Context, id, input, base, remote, defaultBranch string) error {
	if err := Name(id); err != nil {
		return err
	}
	cfg, err := s.Config()
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
	root, err := rootOf(ctx, path)
	if err != nil {
		return err
	}
	if root != path {
		return errors.New("connect the root of a Git checkout")
	}
	if err := checkBranch(ctx, path, base); err != nil {
		return err
	}
	bindings, err := s.Bindings()
	if err != nil {
		return err
	}
	common, err := Git(ctx, path, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return err
	}
	for other, binding := range bindings.Bindings {
		if other == id {
			continue
		}
		otherPath, err := s.LocalPath(binding.Path)
		if err != nil {
			return err
		}
		otherCommon, err := Git(ctx, otherPath, "rev-parse", "--path-format=absolute", "--git-common-dir")
		if err == nil && otherCommon == common {
			return fmt.Errorf("this repository is already connected as %s", other)
		}
	}
	if existing, registered := cfg.Repositories[id]; registered {
		if remote != "" && remote != existing.URL {
			return fmt.Errorf("%s already records another URL; use repo remote to change it explicitly", id)
		}
		if defaultBranch != "" && defaultBranch != existing.DefaultBranch {
			return fmt.Errorf("%s already records another default branch; use repo remote to change it explicitly", id)
		}
	} else {
		if defaultBranch == "" {
			defaultBranch = base
		}
		if err := checkBranch(ctx, path, defaultBranch); err != nil {
			return err
		}
		if remote == "" {
			// Git already records where this checkout came from. A repository
			// with no origin simply has no shared URL to record.
			remote, _ = Git(ctx, path, "remote", "get-url", "origin")
		}
		if remote != "" {
			if err := checkRemote(remote); err != nil {
				return err
			}
		}
		if err := s.Update("workspace.yaml", []string{"repositories", id}, Repository{remote, defaultBranch}, 0644); err != nil {
			return err
		}
	}
	local := path
	if path == s.Root {
		local = "."
	}
	binding := Binding{local, base}
	if _, err := s.Read("repositories.local.yaml"); os.IsNotExist(err) {
		return s.WriteYAML("repositories.local.yaml", Bindings{map[string]Binding{id: binding}}, 0600)
	} else if err != nil {
		return err
	}
	return s.Update("repositories.local.yaml", []string{"bindings", id}, binding, 0600)
}

func (s *Store) CreateRepository(ctx context.Context, id, destination, base, remote, defaultBranch string) error {
	if err := Name(id); err != nil {
		return err
	}
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	if _, ok := cfg.Repositories[id]; ok {
		return errors.New("repository ID is already registered; connect its checkout instead")
	}
	if err := checkBranch(ctx, s.Root, base); err != nil {
		return err
	}
	path, err := s.LocalPath(destination)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		return errors.New("repository creation requires a new destination directory")
	}
	if remote != "" {
		if err := checkRemote(remote); err != nil {
			return err
		}
		_, err = Git(ctx, s.Root, "clone", "--", remote, path)
	} else {
		_, err = Git(ctx, s.Root, "init", "-b", base, "--", path)
	}
	if err != nil {
		return fmt.Errorf("repository creation failed; inspect any partial destination at %s: %w", path, err)
	}
	if err := s.Connect(ctx, id, path, base, remote, defaultBranch); err != nil {
		return fmt.Errorf("repository created at %s but registration failed: %w", path, err)
	}
	return nil
}

// SetBase records the branch this machine starts work from. It is local: one
// checkout's base says nothing about anyone else's, and the shared default in
// workspace.yaml describes the repository rather than a machine.
func (s *Store) SetBase(ctx context.Context, id, branch string) error {
	checkout, err := s.Repository(ctx, id)
	if err != nil {
		return err
	}
	if err := checkBranch(ctx, checkout.Path, branch); err != nil {
		return err
	}
	return s.Update("repositories.local.yaml", []string{"bindings", id, "base_branch"}, branch, 0600)
}

// SetRemote changes what the workspace shares about a repository: where it lives
// and which branch it defaults to. Neither value needs a checkout, so this works
// on a machine that has not connected one.
func (s *Store) SetRemote(ctx context.Context, id, remote, defaultBranch string) error {
	if remote == "" && defaultBranch == "" {
		return errors.New("provide --url or --default-branch")
	}
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	if _, ok := cfg.Repositories[id]; !ok {
		return fmt.Errorf("unknown repository: %s", id)
	}
	if remote != "" {
		if err := checkRemote(remote); err != nil {
			return err
		}
		if err := s.Update("workspace.yaml", []string{"repositories", id, "url"}, remote, 0644); err != nil {
			return err
		}
	}
	if defaultBranch == "" {
		return nil
	}
	if err := checkBranch(ctx, s.Root, defaultBranch); err != nil {
		return err
	}
	return s.Update("workspace.yaml", []string{"repositories", id, "default_branch"}, defaultBranch, 0644)
}

type Snapshot struct {
	Path    string   `json:"path"`
	Branch  string   `json:"branch"`
	Head    string   `json:"head"`
	Changes []string `json:"changes"`
}

func Inspect(ctx context.Context, path string) (Snapshot, error) {
	result := Snapshot{Path: path, Changes: []string{}}
	if _, err := rootOf(ctx, path); err != nil {
		return result, err
	}
	result.Branch, _ = Git(ctx, path, "symbolic-ref", "--quiet", "--short", "HEAD")
	result.Head, _ = Git(ctx, path, "rev-parse", "--verify", "HEAD")
	changes, err := Git(ctx, path, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	if err != nil {
		return result, err
	}
	if changes != "" {
		result.Changes = strings.Split(strings.TrimSuffix(changes, "\x00"), "\x00")
	}
	return result, nil
}
