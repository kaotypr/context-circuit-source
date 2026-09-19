package workspace

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// behindRemote reports how many commits a local branch lacks from the remote
// counterpart this checkout has already fetched, and names that counterpart. It
// never fetches: a check that reached the network would make every caller wait
// on it and fail where the network is absent. So it answers from what the
// checkout already knows, and a branch with no fetched counterpart reports
// nothing rather than claiming the branch is current.
func remoteCounterpart(ctx context.Context, repo, branch string) string {
	if branch == "" || strings.HasPrefix(branch, "-") {
		return ""
	}
	// --end-of-options is echoed back by --abbrev-ref rather than consumed, so
	// this call is guarded by the branch-name check above instead.
	remote, err := Git(ctx, repo, "rev-parse", "--abbrev-ref", branch+"@{upstream}")
	if err == nil && remote != "" && !strings.ContainsAny(remote, "\n\r") {
		return remote
	}
	remote = "origin/" + branch
	if _, err := Git(ctx, repo, "rev-parse", "--verify", "--end-of-options", "refs/remotes/"+remote+"^{commit}"); err != nil {
		return ""
	}
	return remote
}

func behindRemote(ctx context.Context, repo, branch string) (int, string) {
	remote := remoteCounterpart(ctx, repo, branch)
	if remote == "" {
		return 0, ""
	}
	count, err := Git(ctx, repo, "rev-list", "--count", "--end-of-options", "refs/heads/"+branch+".."+remote)
	if err != nil {
		return 0, ""
	}
	n, err := strconv.Atoi(count)
	if err != nil || n <= 0 {
		return 0, ""
	}
	return n, remote
}

func Git(ctx context.Context, path string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, "git", append([]string{"-C", path}, args...)...)
	// A coding host may itself be running inside a Git operation. Its checkout
	// overrides must not redirect commands away from the selected repository.
	for _, item := range os.Environ() {
		key, _, _ := strings.Cut(item, "=")
		key = strings.ToUpper(key)
		switch key {
		case "GIT_DIR", "GIT_WORK_TREE", "GIT_COMMON_DIR", "GIT_INDEX_FILE", "GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES", "GIT_TERMINAL_PROMPT":
			continue
		// One-shot configuration travels the same way and redirects a command
		// just as effectively: `git -c` exports GIT_CONFIG_PARAMETERS to
		// everything it runs, and the count/key/value trio is the same
		// mechanism spelled out. A checkout's own config files are not dropped
		// with them — GIT_CONFIG_GLOBAL and GIT_CONFIG_SYSTEM are set by a
		// person isolating their environment, never by a Git operation.
		case "GIT_CONFIG_PARAMETERS", "GIT_CONFIG_COUNT":
			continue
		}
		if strings.HasPrefix(key, "GIT_CONFIG_KEY_") || strings.HasPrefix(key, "GIT_CONFIG_VALUE_") {
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
	// One namespace across both maps: a note's anchor names a repository by ID,
	// so an ID meaning two things makes every anchor carrying it ambiguous.
	if _, clash := cfg.KnowledgeRepositories[id]; clash {
		return fmt.Errorf("%s already names a knowledge repository, which is read-only here; give this repository another ID", id)
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
	// The workspace's own repository is described by `workspace connect`, not
	// here: an entry in this map is somewhere work happens, and a plan, a
	// relationship, or a worktree cut from the workspace is never what the
	// caller meant. Connecting it here once read as the way to record it, so
	// name the command that replaced it rather than only refusing.
	if within(path, s.Root) {
		return errors.New("this checkout holds the workspace itself; describe it with workspace connect")
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
	// A worktree of the workspace's repository sits outside the workspace, so
	// the containment check above does not see it, but it shares the object
	// store and every commit made in it lands in the workspace's history.
	if selfPath, err := s.selfPath(bindings.Workspace); err == nil && bindings.Workspace != nil {
		if selfCommon, err := Git(ctx, selfPath, "rev-parse", "--path-format=absolute", "--git-common-dir"); err == nil && selfCommon == common {
			return errors.New("this checkout belongs to the workspace's own repository; describe it with workspace connect")
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
	binding := Binding{path, base}
	if _, err := s.Read("repositories.local.yaml"); os.IsNotExist(err) {
		return s.WriteYAML("repositories.local.yaml", Bindings{Bindings: map[string]Binding{id: binding}}, 0600)
	} else if err != nil {
		return err
	}
	return s.Update("repositories.local.yaml", []string{"bindings", id}, binding, 0600)
}

// CreateRepository obtains a working copy and binds it. Cloning an ID the
// workspace already describes is the ordinary case for a member joining one:
// the URL is in the shared record, so the caller need not repeat it, and
// nothing shared is rewritten. Initializing is refused there instead, because
// an ID that already names a repository somewhere is not one to start empty.
func (s *Store) CreateRepository(ctx context.Context, id, destination, base, remote, defaultBranch string, clone bool) error {
	if err := Name(id); err != nil {
		return err
	}
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	if existing, registered := cfg.Repositories[id]; registered {
		if !clone {
			return errors.New("repository ID is already registered; clone or connect its checkout instead")
		}
		switch {
		case remote == "":
			remote = existing.URL
		case remote != existing.URL:
			return fmt.Errorf("%s already records another URL; use repo remote to change it explicitly", id)
		}
		if remote == "" {
			return fmt.Errorf("%s records no URL to clone from; obtain the checkout yourself and use repo connect", id)
		}
		if defaultBranch != "" && defaultBranch != existing.DefaultBranch {
			return fmt.Errorf("%s already records another default branch; use repo remote to change it explicitly", id)
		}
	} else if clone && remote == "" {
		return errors.New("provide --url, or --id naming a repository this workspace already describes")
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

// A snapshot says which branch is checked out now; BaseBranch says which branch
// work here starts from and delivers back to. The two differ routinely, and the
// shared configuration records only a default_branch, which is a plausible and
// wrong answer to "what does this deliver to" — so the recorded base is reported
// wherever a repository or worktree is described rather than left to be found.
type Snapshot struct {
	Path       string   `json:"path"`
	Branch     string   `json:"branch"`
	BaseBranch string   `json:"base_branch,omitempty"`
	Head       string   `json:"head"`
	Changes    []string `json:"changes"`
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
