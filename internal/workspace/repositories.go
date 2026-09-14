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

func rootOf(ctx context.Context, path string) (string, error) {
	root, err := Git(ctx, path, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(root)
}

func (s *Store) Repository(ctx context.Context, id string) (Repository, string, error) {
	cfg, err := s.Config()
	if err != nil {
		return Repository{}, "", err
	}
	repo, ok := cfg.Repositories[id]
	if !ok {
		return repo, "", fmt.Errorf("unknown repository: %s", id)
	}
	bindings, err := s.Bindings()
	if err != nil {
		return repo, "", err
	}
	binding, ok := bindings.Bindings[id]
	if !ok {
		return repo, "", fmt.Errorf("connect this machine's checkout for repository %s", id)
	}
	path, err := s.LocalPath(binding.Path)
	if err != nil {
		return repo, "", err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return repo, "", err
	}
	root, err := rootOf(ctx, path)
	if err != nil {
		return repo, "", err
	}
	if root != path {
		return repo, "", fmt.Errorf("repository binding must name its checkout root: %s", id)
	}
	return repo, path, nil
}

func (s *Store) Connect(ctx context.Context, id, input, base string) error {
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
	if existing, ok := cfg.Repositories[id]; ok && existing.BaseBranch != base {
		return fmt.Errorf("%s has another default base; use repo base to change it explicitly", id)
	}
	local := path
	if path == s.Root {
		local = "."
	}
	if _, exists := cfg.Repositories[id]; !exists {
		if err := s.Update("workspace.yaml", []string{"repositories", id}, Repository{base}, 0644); err != nil {
			return err
		}
	}
	if _, err := s.Read("repositories.local.yaml"); os.IsNotExist(err) {
		return s.WriteYAML("repositories.local.yaml", Bindings{map[string]Binding{id: {local}}}, 0600)
	} else if err != nil {
		return err
	}
	return s.Update("repositories.local.yaml", []string{"bindings", id}, Binding{local}, 0600)
}

func (s *Store) CreateRepository(ctx context.Context, id, destination, base, remote string) error {
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
		if u, err := url.Parse(remote); err == nil && u.User != nil {
			if _, password := u.User.Password(); password {
				return errors.New("use Git's credential helper instead of a URL containing a password")
			}
		}
		_, err = Git(ctx, s.Root, "clone", "--", remote, path)
	} else {
		_, err = Git(ctx, s.Root, "init", "-b", base, "--", path)
	}
	if err != nil {
		return fmt.Errorf("repository creation failed; inspect any partial destination at %s: %w", path, err)
	}
	if err := s.Connect(ctx, id, path, base); err != nil {
		return fmt.Errorf("repository created at %s but registration failed: %w", path, err)
	}
	return nil
}

func (s *Store) SetBase(ctx context.Context, id, branch string) error {
	_, path, err := s.Repository(ctx, id)
	if err != nil {
		return err
	}
	if err := checkBranch(ctx, path, branch); err != nil {
		return err
	}
	return s.Update("workspace.yaml", []string{"repositories", id, "base_branch"}, branch, 0644)
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
