package workspace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Knowledge a workspace owns lives in context/ and is reconciled when the work
// that changed it completes. Knowledge it borrows cannot work that way: an
// organization's knowledge center is shared by many workspaces and changed
// through its own repository, so no plan completed here will ever invalidate it,
// and no edit made here could reach anyone else.
//
// So a borrowed repository is mounted rather than copied. Copying is the failure
// this exists to prevent: a note duplicated into context/ goes stale silently
// while still reading as current, and nothing in the workspace can tell that it
// has. Mounting keeps one copy, upstream, and makes the workspace read it.
//
// Two rules follow, and both are enforced rather than asked for. It is read-only
// here, because an improvement that lands only in this checkout helps nobody and
// is lost on the next fast-forward. And it is never validated the way context/
// is: the catalog rules exist so a person here can fix what they break, and
// neither half is true of a repository this workspace does not own.

// knowledgePushDisabled replaces the push URL of a borrowed checkout. Git echoes
// it back when a push is attempted, so the refusal explains itself at the moment
// someone tries rather than in documentation they did not read.
const knowledgePushDisabled = "READ-ONLY-KNOWLEDGE-REPOSITORY"

// knowledgeIndexNames are tried in order when a connect does not name an index.
// A borrowed repository chooses its own entry point and nothing here may rename
// it, so this only guesses the common spellings and records what it found.
var knowledgeIndexNames = []string{"index.md", "INDEX.md", "context/INDEX.md", "README.md"}

// KnowledgeState is one borrowed repository as this machine sees it. Status
// carries what a caller should do about it, and Detail says why whenever the
// answer is not "nothing".
type KnowledgeState struct {
	ID       string `yaml:"id" json:"id"`
	URL      string `yaml:"url,omitempty" json:"url,omitempty"`
	Tracking string `yaml:"tracking" json:"tracking"`
	Path     string `yaml:"path,omitempty" json:"path,omitempty"`
	Branch   string `yaml:"branch,omitempty" json:"branch,omitempty"`
	Index    string `yaml:"index,omitempty" json:"index,omitempty"`
	Status   string `yaml:"status" json:"status"`
	Detail   string `yaml:"detail,omitempty" json:"detail,omitempty"`
}

// KnowledgeCheckout resolves a borrowed repository to this machine's checkout.
// It reports the shared record separately from the binding so a caller can tell
// "nobody described this" from "this machine has not obtained it".
func (s *Store) KnowledgeCheckout(ctx context.Context, id string) (KnowledgeRepository, string, error) {
	cfg, err := s.Config()
	if err != nil {
		return KnowledgeRepository{}, "", err
	}
	repo, ok := cfg.KnowledgeRepositories[id]
	if !ok {
		return KnowledgeRepository{}, "", fmt.Errorf("unknown knowledge repository: %s", id)
	}
	bindings, err := s.Bindings()
	if err != nil {
		return repo, "", err
	}
	binding, ok := bindings.Knowledge[id]
	if !ok {
		return repo, "", fmt.Errorf("obtain this machine's checkout for knowledge repository %s with knowledge clone --id %s", id, id)
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
		return repo, "", fmt.Errorf("knowledge binding must name its checkout root: %s", id)
	}
	return repo, path, nil
}

// disableKnowledgePush makes a push fail rather than half-succeed. It is applied
// on every sync, not only at clone, so a checkout obtained before this existed
// and one whose push URL was restored by hand both end up read-only again.
func disableKnowledgePush(ctx context.Context, path string) {
	if _, err := Git(ctx, path, "remote", "get-url", "origin"); err != nil {
		return
	}
	_, _ = Git(ctx, path, "remote", "set-url", "--push", "origin", knowledgePushDisabled)
}

// detectKnowledgeIndex records which file maps a borrowed repository's contents.
// An unrecognized layout records nothing rather than guessing wrong; `check`
// reports the gap and a person names the file.
func detectKnowledgeIndex(root string) string {
	for _, name := range knowledgeIndexNames {
		info, err := os.Stat(filepath.Join(root, filepath.FromSlash(name)))
		if err == nil && info.Mode().IsRegular() {
			return name
		}
	}
	return ""
}

// checkKnowledgeIndex rejects an index path that would read outside the borrowed
// checkout or name something other than a file in it.
func checkKnowledgeIndex(value string) error {
	if value == "" {
		return nil
	}
	cleaned := filepath.ToSlash(filepath.Clean(value))
	if filepath.IsAbs(value) || cleaned == "." || strings.HasPrefix(cleaned, "../") || strings.HasPrefix(cleaned, "/") {
		return fmt.Errorf("name the index as a path inside the knowledge repository: %s", value)
	}
	return nil
}

// ConnectKnowledge records a borrowed repository and binds this machine's
// checkout of it. The shared half is written only when the ID is new, the same
// split `repo connect` makes: connecting a second machine describes that
// machine, never the repository.
func (s *Store) ConnectKnowledge(ctx context.Context, id, input, url, defaultBranch, index string) error {
	if err := Name(id); err != nil {
		return err
	}
	if err := checkKnowledgeIndex(index); err != nil {
		return err
	}
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	// One namespace, because a note's anchor names a repository by ID and an ID
	// meaning two things makes every anchor carrying it ambiguous.
	if _, clash := cfg.Repositories[id]; clash {
		return fmt.Errorf("%s already names a repository where work happens; give the knowledge repository another ID", id)
	}
	path, err := s.LocalPath(input)
	if err != nil {
		return err
	}
	if path, err = filepath.EvalSymlinks(path); err != nil {
		return err
	}
	root, err := rootOf(ctx, path)
	if err != nil {
		return err
	}
	if root != path {
		return errors.New("connect the root of a Git checkout")
	}
	if within(path, s.Root) {
		return errors.New("this checkout holds the workspace itself; describe it with workspace connect")
	}
	bindings, err := s.Bindings()
	if err != nil {
		return err
	}
	common, err := Git(ctx, path, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return err
	}
	// A checkout registered both as a repository and as knowledge would be
	// fast-forwarded under work someone is doing in it.
	for other, binding := range bindings.Bindings {
		otherPath, err := s.LocalPath(binding.Path)
		if err != nil {
			continue
		}
		if otherCommon, err := Git(ctx, otherPath, "rev-parse", "--path-format=absolute", "--git-common-dir"); err == nil && otherCommon == common {
			return fmt.Errorf("this checkout is already connected as repository %s, where work happens; knowledge is read-only", other)
		}
	}
	if existing, registered := cfg.KnowledgeRepositories[id]; registered {
		if url != "" && url != existing.URL {
			return fmt.Errorf("%s already records another URL; use knowledge remote to change it explicitly", id)
		}
		if defaultBranch != "" && defaultBranch != existing.DefaultBranch {
			return fmt.Errorf("%s already records another default branch; use knowledge remote to change it explicitly", id)
		}
		if index != "" && index != existing.Index {
			return fmt.Errorf("%s already records another index; use knowledge remote to change it explicitly", id)
		}
	} else {
		if defaultBranch == "" {
			defaultBranch, err = Git(ctx, path, "symbolic-ref", "--quiet", "--short", "HEAD")
			if err != nil || defaultBranch == "" {
				return errors.New("provide --default-branch; this checkout is not on a branch")
			}
		}
		if err := checkBranch(ctx, path, defaultBranch); err != nil {
			return err
		}
		if url == "" {
			url, _ = Git(ctx, path, "remote", "get-url", "origin")
		}
		if url != "" {
			if err := checkRemote(url); err != nil {
				return err
			}
		}
		if index == "" {
			index = detectKnowledgeIndex(path)
		}
		// A workspace that has never borrowed knowledge carries no such key,
		// including every workspace initialized before this existed. Writing the
		// whole map creates it; an entry is added to it once it is there.
		record := KnowledgeRepository{url, defaultBranch, index}
		keys, value := []string{"knowledge_repositories", id}, any(record)
		if len(cfg.KnowledgeRepositories) == 0 {
			keys, value = []string{"knowledge_repositories"}, map[string]KnowledgeRepository{id: record}
		}
		if err := s.Update("workspace.yaml", keys, value, 0644); err != nil {
			return err
		}
	}
	disableKnowledgePush(ctx, path)
	binding := Binding{Path: path}
	if _, err := s.Read("repositories.local.yaml"); os.IsNotExist(err) {
		return s.WriteYAML("repositories.local.yaml", Bindings{Bindings: map[string]Binding{}, Knowledge: map[string]Binding{id: binding}}, 0600)
	} else if err != nil {
		return err
	}
	if len(bindings.Knowledge) == 0 {
		return s.Update("repositories.local.yaml", []string{"knowledge"}, map[string]Binding{id: binding}, 0600)
	}
	return s.Update("repositories.local.yaml", []string{"knowledge", id}, binding, 0600)
}

// CloneKnowledge obtains a borrowed repository and binds it. Cloning an ID the
// workspace already describes is the ordinary case for a member joining one: the
// URL is in the shared record, so the caller need not repeat it.
func (s *Store) CloneKnowledge(ctx context.Context, id, destination, url, defaultBranch, index string) error {
	if err := Name(id); err != nil {
		return err
	}
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	if existing, registered := cfg.KnowledgeRepositories[id]; registered {
		switch {
		case url == "":
			url = existing.URL
		case url != existing.URL:
			return fmt.Errorf("%s already records another URL; use knowledge remote to change it explicitly", id)
		}
		if url == "" {
			return fmt.Errorf("%s records no URL to clone from; obtain the checkout yourself and use knowledge connect", id)
		}
		if defaultBranch == "" {
			defaultBranch = existing.DefaultBranch
		}
	} else if url == "" {
		return errors.New("provide --url, or --id naming a knowledge repository this workspace already describes")
	}
	if destination == "" {
		destination = "knowledge/" + id
	}
	path, err := s.LocalPath(destination)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		return errors.New("obtaining a knowledge repository requires a new destination directory")
	}
	if err := checkRemote(url); err != nil {
		return err
	}
	args := []string{"clone"}
	if defaultBranch != "" {
		args = append(args, "--branch", defaultBranch)
	}
	if _, err := Git(ctx, s.Root, append(args, "--", url, path)...); err != nil {
		return fmt.Errorf("knowledge clone failed; inspect any partial destination at %s: %w", path, err)
	}
	if err := s.ConnectKnowledge(ctx, id, path, url, defaultBranch, index); err != nil {
		return fmt.Errorf("knowledge repository cloned to %s but registration failed: %w", path, err)
	}
	return nil
}

// SyncKnowledge fetches a borrowed repository and fast-forwards it to the branch
// the workspace shares, which is the only way it ever moves. It never merges,
// rebases, resets, or discards: a checkout holding local work is reported and
// left exactly as it is, because the work is somebody's and this command has no
// standing to decide what happens to it.
func (s *Store) SyncKnowledge(ctx context.Context, id string) (KnowledgeState, error) {
	repo, path, err := s.KnowledgeCheckout(ctx, id)
	state := KnowledgeState{ID: id, URL: repo.URL, Tracking: repo.DefaultBranch, Index: repo.Index}
	if err != nil {
		state.Status, state.Detail = "unbound", err.Error()
		return state, err
	}
	state.Path = path
	disableKnowledgePush(ctx, path)
	if _, err := Git(ctx, path, "fetch", "--prune", "--quiet", "origin"); err != nil {
		state.Status, state.Detail = "unreachable", err.Error()
		return state, err
	}
	snapshot, err := Inspect(ctx, path)
	if err != nil {
		state.Status, state.Detail = "unreachable", err.Error()
		return state, err
	}
	state.Branch = snapshot.Branch
	if len(snapshot.Changes) > 0 {
		state.Status = "dirty"
		state.Detail = fmt.Sprintf("%d uncommitted paths; nothing was changed. Move the work to a separate checkout of %s and open a merge request there, then sync again", len(snapshot.Changes), id)
		return state, nil
	}
	if snapshot.Branch != repo.DefaultBranch {
		state.Status = "off-branch"
		state.Detail = fmt.Sprintf("on %q, and the workspace reads %q; nothing was changed", snapshot.Branch, repo.DefaultBranch)
		return state, nil
	}
	before := snapshot.Head
	if _, err := Git(ctx, path, "merge", "--ff-only", "--quiet", "--end-of-options", "refs/remotes/origin/"+repo.DefaultBranch); err != nil {
		state.Status = "diverged"
		state.Detail = fmt.Sprintf("this checkout holds commits that are not on origin/%s; nothing was changed or discarded. Push them from a separate checkout of %s, then sync again", repo.DefaultBranch, id)
		return state, nil
	}
	after, _ := Git(ctx, path, "rev-parse", "--verify", "HEAD")
	if after != before {
		state.Status = "updated"
	} else {
		state.Status = "current"
	}
	return state, nil
}

// KnowledgeStates describes every borrowed repository without touching the
// network, so orientation and `check` can report them without waiting.
func (s *Store) KnowledgeStates(ctx context.Context) ([]KnowledgeState, error) {
	cfg, err := s.Config()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(cfg.KnowledgeRepositories))
	for id := range cfg.KnowledgeRepositories {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	states := make([]KnowledgeState, 0, len(ids))
	for _, id := range ids {
		repo := cfg.KnowledgeRepositories[id]
		state := KnowledgeState{ID: id, URL: repo.URL, Tracking: repo.DefaultBranch, Index: repo.Index}
		_, path, err := s.KnowledgeCheckout(ctx, id)
		if err != nil {
			state.Status, state.Detail = "unbound", err.Error()
			states = append(states, state)
			continue
		}
		state.Path = path
		snapshot, err := Inspect(ctx, path)
		if err != nil {
			state.Status, state.Detail = "unreachable", err.Error()
			states = append(states, state)
			continue
		}
		state.Branch = snapshot.Branch
		switch {
		case len(snapshot.Changes) > 0:
			state.Status = "dirty"
			state.Detail = fmt.Sprintf("%d uncommitted paths in a repository that is read-only here", len(snapshot.Changes))
		case snapshot.Branch != repo.DefaultBranch:
			state.Status = "off-branch"
			state.Detail = fmt.Sprintf("on %q, and the workspace reads %q", snapshot.Branch, repo.DefaultBranch)
		default:
			state.Status = "bound"
			if behind, remote := behindRemote(ctx, path, snapshot.Branch); behind > 0 {
				state.Status = "behind"
				state.Detail = fmt.Sprintf("%d commits behind %s as of the last fetch; run knowledge sync", behind, remote)
			}
		}
		states = append(states, state)
	}
	return states, nil
}

// borrowedIssues reports what stands between this machine and the knowledge it
// is supposed to be reading. It never inspects the borrowed content itself: the
// catalog rules exist so a person here can fix what they break, and a repository
// this workspace does not own fails neither half of that. Reporting a fault
// nobody here may correct would be noise in a diagnostic that has to stay worth
// reading.
func borrowedIssues(states []KnowledgeState) []string {
	var issues []string
	for _, state := range states {
		switch state.Status {
		case "unbound":
			issues = append(issues, fmt.Sprintf("knowledge repository %s is not obtained on this machine: run knowledge clone --id %s", state.ID, state.ID))
		case "unreachable":
			issues = append(issues, fmt.Sprintf("knowledge repository %s: %s", state.ID, state.Detail))
		case "dirty":
			issues = append(issues, fmt.Sprintf("knowledge repository %s is read-only here yet has local changes (%s): move them to a separate checkout and open a merge request there", state.ID, state.Detail))
		case "off-branch":
			issues = append(issues, fmt.Sprintf("knowledge repository %s is %s: check out %s so this machine reads what everyone else reads", state.ID, state.Detail, state.Tracking))
		case "behind":
			issues = append(issues, fmt.Sprintf("knowledge repository %s is %s", state.ID, state.Detail))
		}
		if state.Status == "unbound" {
			continue
		}
		if state.Index == "" {
			issues = append(issues, fmt.Sprintf("knowledge repository %s records no index, so nothing here can retrieve from it: name it with knowledge remote --id %s --index <path>", state.ID, state.ID))
			continue
		}
		if state.Path != "" {
			if _, err := readBorrowed(state.Path, state.Index); err != nil {
				issues = append(issues, fmt.Sprintf("knowledge repository %s records index %s, which cannot be read: %s", state.ID, state.Index, err.Error()))
			}
		}
	}
	return issues
}

// SetKnowledgeRemote changes what the workspace shares about a borrowed
// repository. None of the three values needs a checkout, so this works on a
// machine that has not obtained one.
func (s *Store) SetKnowledgeRemote(ctx context.Context, id, url, defaultBranch, index string) error {
	if url == "" && defaultBranch == "" && index == "" {
		return errors.New("provide --url, --default-branch, or --index")
	}
	if err := checkKnowledgeIndex(index); err != nil {
		return err
	}
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	if _, ok := cfg.KnowledgeRepositories[id]; !ok {
		return fmt.Errorf("unknown knowledge repository: %s", id)
	}
	if url != "" {
		if err := checkRemote(url); err != nil {
			return err
		}
		if err := s.Update("workspace.yaml", []string{"knowledge_repositories", id, "url"}, url, 0644); err != nil {
			return err
		}
	}
	if defaultBranch != "" {
		if err := checkBranch(ctx, s.Root, defaultBranch); err != nil {
			return err
		}
		if err := s.Update("workspace.yaml", []string{"knowledge_repositories", id, "default_branch"}, defaultBranch, 0644); err != nil {
			return err
		}
	}
	if index == "" {
		return nil
	}
	return s.Update("workspace.yaml", []string{"knowledge_repositories", id, "index"}, index, 0644)
}

// readBorrowed reads a file from inside a borrowed checkout. Nothing here goes
// through Store.Read, which resolves against the workspace root and refuses to
// leave it; a borrowed checkout is deliberately outside. The guards it would
// have applied are applied here instead, against the checkout root.
func readBorrowed(root, relative string) ([]byte, error) {
	if err := checkKnowledgeIndex(relative); err != nil {
		return nil, err
	}
	target := filepath.Join(root, filepath.FromSlash(relative))
	resolved, err := filepath.EvalSymlinks(target)
	if err != nil {
		return nil, err
	}
	if !within(root, resolved) {
		return nil, fmt.Errorf("index resolves outside the knowledge repository: %s", relative)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return nil, fmt.Errorf("expected a regular file under 16 MiB: %s", relative)
	}
	return os.ReadFile(resolved)
}

// KnowledgeIndex returns one borrowed repository's index as text, for retrieval
// across it. The content is never parsed as a catalog: a borrowed repository
// writes entries its own way, and holding it to this product's shape would
// report faults in a file nobody here may fix.
func (s *Store) KnowledgeIndex(ctx context.Context, id string) (string, []byte, error) {
	repo, path, err := s.KnowledgeCheckout(ctx, id)
	if err != nil {
		return "", nil, err
	}
	if repo.Index == "" {
		return "", nil, fmt.Errorf("%s records no index; name it with knowledge remote --id %s --index <path>", id, id)
	}
	data, err := readBorrowed(path, repo.Index)
	if err != nil {
		return repo.Index, nil, err
	}
	return repo.Index, data, nil
}
