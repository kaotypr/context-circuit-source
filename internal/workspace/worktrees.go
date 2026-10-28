package workspace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
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

	// Start is what a new worktree actually started from: --start as given, the
	// base branch or origin/<base> it fell back to, or the existing branch it
	// checked out or resumed. BaseBranch is this machine's base, the branch work on it is
	// delivered to, which is not always where it started.
	Start       string       `json:"start,omitempty"`
	BaseBranch  string       `json:"base_branch,omitempty"`
	Reused      bool         `json:"reused,omitempty"`
	Environment []ReuseEntry `json:"environment,omitempty"`

	// The base branch this worktree started from was behind its fetched remote.
	// Preparation succeeded, so this is reported rather than refused: the work
	// is recoverable, and only the caller knows whether the missing commits
	// matter. Naming it here is what stops it being discovered at review.
	SyncRequired string `json:"sync_required,omitempty"`

	// Plans holding work in this repository that this worktree's starting point
	// does not contain. Reported, never refused: parallel plans on one
	// repository are ordinary, and only the caller knows which it meant to
	// build on.
	UnmergedPlans string `json:"unmerged_plans,omitempty"`

	// planCreatedAt is the created_at of the plan record this worktree was
	// prepared for, kept in the association and never printed.
	planCreatedAt string
}
type Association struct {
	Repository  string `yaml:"repository" json:"repository"`
	Plan        string `yaml:"plan,omitempty" json:"plan,omitempty"`
	Path        string `yaml:"path" json:"path"`
	Branch      string `yaml:"branch" json:"branch"`
	StartCommit string `yaml:"start_commit" json:"start_commit"`

	// PlanCreatedAt tells the plan record this worktree was prepared for from
	// an earlier one given the same number after a deletion released it.
	// Associations written before 2.4.0 lack it.
	PlanCreatedAt ISOTime `yaml:"plan_created_at,omitempty" json:"plan_created_at,omitempty"`
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

// errNotWorktree is the refusal for a path Git does not list as a worktree of the repository, kept apart
// so that a caller can tell a worktree that is gone from a repository that cannot be read.
var errNotWorktree = errors.New("path is not a registered worktree of this repository")

func (s *Store) Worktrees(ctx context.Context, repoID string) ([]Worktree, error) {
	checkout, err := s.Repository(ctx, repoID)
	if err != nil {
		return nil, err
	}
	text, err := Git(ctx, checkout.Path, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, err
	}
	list := []Worktree{}
	var current *Worktree
	for field := range strings.SplitSeq(text, "\x00") {
		if after, ok := strings.CutPrefix(field, "worktree "); ok {
			// Git prints forward slashes on Windows. Return the same canonical
			// native path for a new worktree, listing, and later reuse.
			path := filepath.Clean(after)
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
				list[i].Plan, list[i].StartCommit, list[i].planCreatedAt = assoc.Plan, assoc.StartCommit, string(assoc.PlanCreatedAt)
			}
		}
	}
	return list, nil
}

func (s *Store) Prepare(ctx context.Context, repoID, plan, branch, start, destination string, reuse bool, reuseOptions ...ReuseOptions) (Worktree, error) {
	options := ReuseOptions{Mode: "auto"}
	if len(reuseOptions) > 0 {
		options = reuseOptions[0]
	}
	if err := options.Validate(); err != nil {
		return Worktree{}, err
	}
	checkout, err := s.Repository(ctx, repoID)
	if err != nil {
		return Worktree{}, err
	}
	repo := checkout.Path
	record, err := s.planRecord(plan, repoID)
	if err != nil {
		return Worktree{}, err
	}
	if branch == "" && plan != "" {
		branch = PlanBranch(plan, repoID)
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
			return s.resumeWorktree(ctx, repoID, plan, branch, reuse, options, checkout, record, tree)
		}
		if tree.Branch == branch {
			return Worktree{}, checkedOutElsewhere(repoID, plan, branch, tree)
		}
	}
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		return Worktree{}, fmt.Errorf("destination already exists: %s", path)
	}
	var commit, stale, unmerged, started string
	if _, err := Git(ctx, repo, "show-ref", "--verify", "refs/heads/"+branch); err == nil {
		started = branch
		commit, err = checkoutExistingBranch(ctx, repo, checkout, record, plan, branch, path, reuse)
		if err != nil {
			return Worktree{}, err
		}
	} else {
		commit, started, stale, unmerged, err = s.addNewBranch(ctx, repoID, plan, branch, start, path, checkout, repo, record)
		if err != nil {
			return Worktree{}, err
		}
	}
	actual, err := filepath.EvalSymlinks(path)
	if err != nil {
		return Worktree{}, err
	}
	result := Worktree{Path: actual, Branch: branch, Head: commit, Plan: plan, StartCommit: commit, Start: started, BaseBranch: checkout.BaseBranch, SyncRequired: stale, UnmergedPlans: unmerged, planCreatedAt: string(record.CreatedAt)}
	if err := s.saveAssociation(repoID, result); err != nil {
		return result, fmt.Errorf("worktree exists at %s but local association could not be saved: %w", actual, err)
	}
	result.Environment, err = ReuseEnvironment(ctx, repo, result.Path, options)
	return result, err
}

// planRecord finds the plan a worktree is prepared for; no plan means no record.
func (s *Store) planRecord(plan, repoID string) (Record, error) {
	if plan == "" {
		return Record{}, nil
	}
	record, err := s.FindRecord(plan)
	if err != nil {
		return Record{}, err
	}
	if !strings.HasPrefix(plan, "p") || !slices.Contains(record.Repositories, repoID) {
		return Record{}, errors.New("plan does not include this repository")
	}
	return record, nil
}

// resumeWorktree validates and adopts a worktree that already sits at the
// destination, saving its association when it is new or changed.
func (s *Store) resumeWorktree(ctx context.Context, repoID, plan, branch string, reuse bool, options ReuseOptions, checkout Checkout, record Record, tree Worktree) (Worktree, error) {
	repo := checkout.Path
	if tree.Main || tree.Branch != branch || tree.Prunable {
		return Worktree{}, errors.New("destination is an existing different or unavailable worktree")
	}
	if (tree.Plan != plan || tree.StartCommit == "") && !reuse {
		return Worktree{}, errors.New("worktree has another association; use --reuse to select it explicitly")
	}
	// The same plan ID is not the same plan once a deletion released
	// the number: another machine's leftover worktree for the old plan
	// matches on path, branch, and ID alike.
	if !reuse {
		if why := earlierPlan(ctx, repo, record, tree.planCreatedAt, tree.StartCommit, branch); why != "" {
			return Worktree{}, earlierPlanError(repoID, plan, branch, tree.Path, why)
		}
	}
	snapshot, err := Inspect(ctx, tree.Path)
	if err != nil {
		return Worktree{}, err
	}
	if !reuse && record.ID != "" && tree.planCreatedAt == "" {
		if why := unprovenWorktree(ctx, repo, checkout.BaseBranch, record, tree, snapshot); why != "" {
			return Worktree{}, earlierPlanError(repoID, plan, branch, tree.Path, why)
		}
	}
	if tree.Plan != plan || tree.StartCommit == "" || tree.planCreatedAt != string(record.CreatedAt) {
		tree.Plan, tree.planCreatedAt = plan, string(record.CreatedAt)
		if tree.StartCommit == "" {
			tree.StartCommit = tree.Head
		}
		if err := s.saveAssociation(repoID, tree); err != nil {
			return tree, err
		}
	}
	// A resumed worktree carries on from its own branch, which is what it
	// starts from now whatever it was first prepared from.
	tree.Reused, tree.Start = true, tree.Branch
	tree.Environment, err = ReuseEnvironment(ctx, repo, tree.Path, options)
	return tree, err
}

// checkoutExistingBranch checks out a branch that already exists, which is only
// allowed when the caller chose to reuse it.
func checkoutExistingBranch(ctx context.Context, repo string, checkout Checkout, record Record, plan, branch, path string, reuse bool) (commit string, err error) {
	if !reuse {
		// A clone with only origin's base branch would otherwise date the
		// branch's whole history.
		base, _ := baseRef(ctx, repo, checkout.BaseBranch)
		if why := earlierPlan(ctx, repo, record, "", base, branch); why != "" {
			return "", fmt.Errorf("branch %s already exists and %s, so it is likely left from an earlier plan given the number %s. Delete that branch once its work no longer matters, or pass --reuse to start from its commits deliberately. Nothing was checked out", branch, why, plan)
		}
		return "", errors.New("branch already exists; use --reuse to preserve and check out its existing commits")
	}
	commit, err = Git(ctx, repo, "rev-parse", "--verify", "--end-of-options", "refs/heads/"+branch+"^{commit}")
	if err != nil {
		return "", err
	}
	_, err = Git(ctx, repo, "worktree", "add", "--", path, branch)
	if err != nil {
		// The inventory above catches a branch Git lists anywhere, but a
		// worktree mid-rebase or mid-bisect lists as detached while Git
		// still refuses to check its branch out twice. Git's text names the
		// place and varies by version and locale, so it is kept, not parsed.
		return "", fmt.Errorf("branch %s could not be checked out; if Git says it is used by another worktree, finish or abort the operation there, switch that checkout to another branch, or remove that worktree (dirty work is kept unless discard is explicitly authorized), then prepare again. Nothing was checked out, switched, or reset: %w", branch, err)
	}
	return commit, nil
}

// addNewBranch creates the branch and its worktree from the explicit start or
// the derived base, reporting a derived start that is stale or unmerged.
func (s *Store) addNewBranch(ctx context.Context, repoID, plan, branch, start, path string, checkout Checkout, repo string, record Record) (commit, started, stale, unmerged string, err error) {
	derived := false
	if start == "" {
		derived = true
		// Preparation does not derive order, so a dependent plan started
		// from the base branch gets a worktree missing the work it was
		// meant to build on — which reads as the dependency having produced
		// nothing. Name the start rather than choose it: with several
		// predecessors the extras are integration merges, and that is the
		// coordinator's to perform.
		if dependency := s.dependencyStart(record, repoID); dependency.Base != "" {
			return "", "", "", "", dependencyStartError(ctx, repo, plan, dependency)
		}
		start, started = "refs/heads/"+checkout.BaseBranch, checkout.BaseBranch
		if _, err := Git(ctx, repo, "rev-parse", "--verify", "--end-of-options", start+"^{commit}"); err != nil {
			start, started = "refs/remotes/origin/"+checkout.BaseBranch, "origin/"+checkout.BaseBranch
		} else if n, remote := behindRemote(ctx, repo, checkout.BaseBranch); n > 0 {
			// A local base branch nobody synchronized is the quietest way to
			// get a worktree missing work that already exists: preparation
			// reports the branch it used and nothing about how old it is, so
			// the gap is found at review or not at all.
			stale = fmt.Sprintf("%s is %d commit(s) behind %s and this worktree starts from it, so that work is not here. Fast-forward %s in the bound checkout and rebase %s onto it, or remove this worktree and prepare again.",
				checkout.BaseBranch, n, remote, checkout.BaseBranch, branch)
		}
	} else {
		started = start
	}
	commit, err = Git(ctx, repo, "rev-parse", "--verify", "--end-of-options", start+"^{commit}")
	if err != nil {
		return "", "", "", "", fmt.Errorf("cannot resolve starting point %s (new repositories need a first commit): %w", start, err)
	}
	if !derived {
		stale = staleStart(ctx, repo, repoID, start, commit, branch)
	}
	// An explicit --start is the coordinator saying what they meant, so only
	// a start this command derived is questioned.
	if derived {
		remoteBase := remoteCounterpart(ctx, repo, checkout.BaseBranch)
		unmerged = describeUnmerged(s.unmergedSiblings(ctx, repo, repoID, plan, commit, remoteBase), repoID, plan, checkout.BaseBranch)
	}
	_, err = Git(ctx, repo, "worktree", "add", "-b", branch, "--", path, commit)
	if err != nil {
		return "", "", "", "", fmt.Errorf("worktree preparation failed; inspect branch %s and path %s: %w", branch, path, err)
	}
	return commit, started, stale, unmerged, nil
}

// staleStart questions an explicit --start that names a local branch behind its
// remote-tracking branch, or diverged from it. The worktree still starts where
// the caller said, since that is what they asked for and nothing here fetches;
// the report says what is missing and how to start from the remote instead,
// the way a stale derived base is reported. A start that is not a local
// branch, such as origin/main or a commit, is not questioned.
func staleStart(ctx context.Context, repo, repoID, start, commit, branch string) string {
	local := strings.TrimPrefix(start, "refs/heads/")
	if local == "" || strings.HasPrefix(local, "-") {
		return ""
	}
	// Only when the start resolved to that local branch: a tag of the same name
	// would have won, and then the branch is not what was used.
	head, err := Git(ctx, repo, "rev-parse", "--verify", "--end-of-options", "refs/heads/"+local+"^{commit}")
	if err != nil || head != commit {
		return ""
	}
	ahead, behind, remote := aheadBehind(ctx, repo, local)
	if behind == 0 {
		return ""
	}
	state := fmt.Sprintf("is %d commit(s) behind %s", behind, remote)
	if ahead > 0 {
		state = fmt.Sprintf("has diverged from %s: %d commit(s) ahead and %d behind", remote, ahead, behind)
	}
	return fmt.Sprintf("--start %s names the local branch, which %s, and this worktree starts from the local branch, so that work is not here. To start from %s, remove this worktree and prepare again with --start %s; or bring %s up to date in the bound checkout and rebase %s onto it. %s is as recent as the last fetch, which this command never makes; `repo fetch --id %s` updates it.",
		local, state, remote, remote, local, branch, remote, repoID)
}

// checkedOutElsewhere explains a plan branch Git already has checked out in
// another working copy. A path alone left the caller to work out whether that
// copy was theirs to resume or someone's checkout to leave alone, so the
// refusal says which, and names only steps that preserve the work there.
// Nothing is checked out, switched, or reset to make preparation succeed.
func checkedOutElsewhere(repoID, plan, branch string, tree Worktree) error {
	var text strings.Builder
	fmt.Fprintf(&text, "branch %s is already checked out at %s. ", branch, tree.Path)
	switch {
	case tree.Main:
		fmt.Fprintf(&text, "That is the bound checkout of %s, not a worktree for this plan. Switch that checkout to another branch, then prepare again.", repoID)
	case plan != "" && tree.Plan == plan:
		// The path is quoted plainly rather than escaped, so a Windows path
		// keeps single backslashes and still pastes into --path whole.
		quoted := `"` + tree.Path + `"`
		fmt.Fprintf(&text, "That is the worktree this workspace prepared for %s. Resume it with `worktree prepare --repo %s --plan %s --path %s`, or remove it with `worktree remove --repo %s --path %s` (dirty work is kept unless discard is explicitly authorized) and prepare again.",
			plan, repoID, plan, quoted, repoID, quoted)
	default:
		if tree.Plan != "" {
			fmt.Fprintf(&text, "That worktree is associated with %s, not this plan. ", tree.Plan)
		} else {
			text.WriteString("That is not a worktree this workspace knows for this plan. ")
		}
		text.WriteString("Resume the work there, switch that checkout to another branch, or remove that worktree (dirty work is kept unless discard is explicitly authorized), then prepare again.")
	}
	if tree.Prunable {
		quoted := `"` + tree.Path + `"`
		fmt.Fprintf(&text, " Git marks that worktree prunable: run `worktree repair --repo %s --path %s` if it moved; if it is gone, run `git worktree prune`, then `worktree prune --repo %s`.",
			repoID, quoted, repoID)
	}
	text.WriteString(" Nothing was checked out, switched, or reset.")
	return errors.New(text.String())
}

// earlierPlan says why a worktree or branch found for a plan belongs to an
// earlier plan that held the same number, or returns "" when nothing says so.
// An association records the plan's created_at, which settles it. One written
// before that, or a bare branch, is judged by its commits: work committed
// before this plan record existed was not done for it.
func earlierPlan(ctx context.Context, repo string, plan Record, recorded, start, branch string) string {
	if plan.ID == "" {
		return ""
	}
	if recorded != "" {
		if recorded == string(plan.CreatedAt) {
			return ""
		}
		return fmt.Sprintf("it was prepared for a %s created at %s, while this %s was created at %s", plan.ID, recorded, plan.ID, plan.CreatedAt)
	}
	created, err := time.Parse(TimeLayout, string(plan.CreatedAt))
	if err != nil {
		return ""
	}
	span := branch
	if start != "" {
		span = start + ".." + branch
	}
	times, err := Git(ctx, repo, "log", "--format=%ct", "--end-of-options", span)
	if err != nil {
		return ""
	}
	for field := range strings.FieldsSeq(times) {
		if at, err := strconv.ParseInt(field, 10, 64); err == nil && at < created.Unix() {
			return fmt.Sprintf("it holds commits made before this %s was created at %s", plan.ID, plan.CreatedAt)
		}
	}
	return ""
}

// unprovenWorktree says why a worktree whose association predates
// plan_created_at, and whose branch holds no commit past its start, is not safe
// to resume, or returns "". With no commit there is nothing to date, so it is
// resumed only when it is exactly what a fresh prepare would create: clean, and
// at the tip of the base branch. Anything else may be the earlier plan's.
func unprovenWorktree(ctx context.Context, repo, base string, plan Record, tree Worktree, snapshot Snapshot) string {
	count, err := Git(ctx, repo, "rev-list", "--count", "--end-of-options", tree.StartCommit+"..refs/heads/"+tree.Branch)
	if err != nil {
		return "its association predates 2.4.0 and its commits could not be read"
	}
	if count != "0" {
		return ""
	}
	if len(snapshot.Changes) > 0 {
		return "its association predates 2.4.0, it holds no commit to date it by, and its uncommitted changes are tied to no " + plan.ID
	}
	if ref, commit := baseRef(ctx, repo, base); ref == "" || snapshot.Head != tree.StartCommit || commit != snapshot.Head {
		return fmt.Sprintf("its association predates 2.4.0, it holds no commit to date it by, and it is not at the tip of %s, where a fresh prepare for this %s would start", base, plan.ID)
	}
	return ""
}

// baseRef is the ref a fresh prepare starts from, and its commit: the local
// base branch, or origin's when this clone has none. Both are "" when neither
// exists.
func baseRef(ctx context.Context, repo, base string) (string, string) {
	for _, ref := range []string{"refs/heads/" + base, "refs/remotes/origin/" + base} {
		if commit, err := Git(ctx, repo, "rev-parse", "--verify", "--end-of-options", ref+"^{commit}"); err == nil {
			return ref, commit
		}
	}
	return "", ""
}

// earlierPlanError refuses to resume another plan's worktree under a reused
// number and names only steps that keep the old work until a person drops it.
func earlierPlanError(repoID, plan, branch, path, why string) error {
	quoted := `"` + path + `"`
	return fmt.Errorf("the worktree at %s on %s is not this plan's to resume: %s, so it is likely left from an earlier plan given the number %s. Remove it with `worktree remove --repo %s --path %s` (dirty work is kept unless discard is explicitly authorized) and delete branch %s once its work no longer matters, or pass --reuse to continue from its commits deliberately. Nothing was checked out, switched, or reset",
		path, branch, why, plan, repoID, quoted, branch)
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
	kept = append(kept, Association{repo, tree.Plan, tree.Path, tree.Branch, tree.StartCommit, ISOTime(tree.planCreatedAt)})
	return s.WriteYAML(".context-circuit/local/worktrees.yaml", Associations{kept}, 0600)
}

func (s *Store) InspectWorktree(ctx context.Context, repoID, input string) (Snapshot, error) {
	_, tree, err := s.selectedTree(ctx, repoID, input)
	if err != nil {
		return Snapshot{}, err
	}
	snapshot, err := Inspect(ctx, tree.Path)
	if err != nil {
		return snapshot, err
	}
	// Delivery happens from a worktree, long after preparation, so the branch
	// this work delivers back to belongs in the description of the worktree.
	checkout, err := s.Repository(ctx, repoID)
	if err != nil {
		return snapshot, err
	}
	snapshot.BaseBranch = checkout.BaseBranch
	return snapshot, nil
}

func (s *Store) selectedTree(ctx context.Context, repoID, input string) (string, Worktree, error) {
	checkout, err := s.Repository(ctx, repoID)
	if err != nil {
		return "", Worktree{}, err
	}
	repo := checkout.Path
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
	return "", Worktree{}, errNotWorktree
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

// RemoveWorktree removes one worktree through Git and drops its association. It then removes the
// directories under the workspace's .worktrees/ that the removal left empty, and returns them, never nil,
// with a warning for each one that could not be removed; neither fails the command, because the worktree
// is already gone.
func (s *Store) RemoveWorktree(ctx context.Context, repoID, input string, discard bool) ([]string, []string, error) {
	repo, tree, err := s.selectedTree(ctx, repoID, input)
	if err != nil {
		return nil, nil, err
	}
	if tree.Main || tree.Locked {
		return nil, nil, errors.New("primary or locked worktrees cannot be removed by this command")
	}
	changes, err := Git(ctx, tree.Path, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignored")
	if err != nil {
		return nil, nil, err
	}
	if changes != "" && !discard {
		return nil, nil, errors.New("worktree contains changes, untracked, or ignored files; preserve them or explicitly authorize --discard")
	}
	args := []string{"worktree", "remove"}
	if discard {
		args = append(args, "--force")
	}
	args = append(args, "--", tree.Path)
	if _, err := Git(ctx, repo, args...); err != nil {
		return nil, nil, err
	}
	if err := s.updateAssociation(repoID, tree.Path, ""); err != nil {
		return nil, nil, fmt.Errorf("worktree removed; local association update failed: %w", err)
	}
	removed, warnings := s.removeEmptyWorktreeParents(tree.Path)
	return removed, warnings, nil
}

// removeEmptyWorktreeParents removes the directories that removing the worktree at path left empty under
// the workspace's .worktrees/, from its parent up to and including .worktrees/ itself, and stops at the
// first one that is not empty. Both sides are compared resolved, so a worktree outside .worktrees/ gets no
// cleanup at all. Only os.Remove is used, which refuses a directory that is not empty, and .worktrees/ is
// kept when it is a symbolic link rather than removing what it points to. A directory that cannot be
// removed stops the walk and is reported as a warning.
func (s *Store) removeEmptyWorktreeParents(path string) ([]string, []string) {
	removed, warnings := []string{}, []string(nil)
	link, err := s.LocalPath(".worktrees")
	if err != nil {
		return removed, warnings
	}
	info, err := os.Lstat(link)
	if err != nil {
		return removed, warnings
	}
	root, err := filepath.EvalSymlinks(link)
	if err != nil {
		return removed, warnings
	}
	dir, err := filepath.EvalSymlinks(filepath.Dir(filepath.Clean(path)))
	if err != nil {
		return removed, warnings
	}
	for {
		rel, err := filepath.Rel(root, dir)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
			return removed, warnings
		}
		if rel == "." && info.Mode()&os.ModeSymlink != 0 {
			return removed, warnings
		}
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			return removed, warnings
		}
		if err := os.Remove(dir); err != nil {
			return removed, append(warnings, fmt.Sprintf("the empty directory %s was left in place: %v", dir, err))
		}
		removed = append(removed, dir)
		if rel == "." {
			return removed, warnings
		}
		dir = filepath.Dir(dir)
	}
}

// PruneWorktreeAssociations drops the plan-to-worktree associations of one repository whose path Git no
// longer lists as a worktree, which is what remains when a worktree is removed with `git worktree remove`
// or by hand. It reads Git's inventory and rewrites only the local association file: no worktree, branch,
// or file is touched, and an association whose worktree Git still lists (even one it marks prunable, or one
// now on another branch) is left for `worktree repair` and `git worktree prune` to judge. It returns what
// it dropped, never nil.
func (s *Store) PruneWorktreeAssociations(ctx context.Context, repoID string) ([]Association, error) {
	list, err := s.Worktrees(ctx, repoID)
	if err != nil {
		return nil, err
	}
	listed := make(map[string]bool, len(list))
	for _, tree := range list {
		listed[filepath.Clean(tree.Path)] = true
	}
	local, err := s.associations()
	if err != nil {
		return nil, err
	}
	kept, dropped := []Association{}, []Association{}
	for _, assoc := range local.Worktrees {
		path := filepath.Clean(assoc.Path)
		if resolved, e := filepath.EvalSymlinks(path); e == nil {
			path = resolved
		}
		if assoc.Repository == repoID && !listed[path] {
			dropped = append(dropped, assoc)
			continue
		}
		kept = append(kept, assoc)
	}
	if len(dropped) == 0 {
		return dropped, nil
	}
	if err := s.WriteYAML(".context-circuit/local/worktrees.yaml", Associations{kept}, 0600); err != nil {
		return nil, err
	}
	return dropped, nil
}

func (s *Store) RepairWorktree(ctx context.Context, repoID, input string) error {
	checkout, err := s.Repository(ctx, repoID)
	if err != nil {
		return err
	}
	repo := checkout.Path
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

// dependencyStart reports the branch a dependent plan's worktree must start from
// in one repository, plus any integration merges a second predecessor adds. An
// empty base means the plan depends on nothing that touches this repository.
func (s *Store) dependencyStart(plan Record, repo string) OrderStart {
	if plan.ID == "" || len(plan.DependsOn) == 0 {
		return OrderStart{}
	}
	return s.startFor(map[string]Record{}, plan, repo, "")
}

// unmergedSiblings names the plans holding work in this repository that a given
// starting point does not contain. It answers the question a declared
// dependency answers explicitly, for the far more common case where nobody
// declared one: the guard in preparation is armed by `depends_on`, and a plan
// that builds on another without saying so gets a worktree missing its work.
//
// Containment is the test rather than completion, because the two disagree in
// both directions — a plan can be completed and never merged, or merged while
// still open — and it is containment that decides what a worktree will hold.
//
// This reports and never refuses. Two plans touching one repository are
// ordinary parallel work, so a refusal here would block what waves exist to
// allow; only the coordinator knows whether the missing work is wanted.
type unmergedPlan struct {
	Plan   string
	Branch string
	Count  string
}

func (s *Store) unmergedSiblings(ctx context.Context, repo, repoID, plan, start, remoteBase string) []unmergedPlan {
	records, err := s.ListRecords(false)
	if err != nil || start == "" {
		return nil
	}
	found := []unmergedPlan{}
	for _, other := range records {
		if other.ID == plan || !strings.HasPrefix(other.ID, "p") || !slices.Contains(other.Repositories, repoID) {
			continue
		}
		branch := PlanBranch(other.ID, repoID)
		if _, err := Git(ctx, repo, "show-ref", "--verify", "refs/heads/"+branch); err != nil {
			continue
		}
		if _, err := Git(ctx, repo, "merge-base", "--is-ancestor", "--end-of-options", branch, start); err == nil {
			continue
		}
		// Work already merged upstream is not work this plan has to declare a
		// dependency on; it is work this checkout has not pulled. Saying
		// "depend on it" there would record a relationship that does not exist
		// and branch from a merged branch instead of the updated base. The
		// staleness report names the real act, so leave that to it.
		if remoteBase != "" {
			if _, err := Git(ctx, repo, "merge-base", "--is-ancestor", "--end-of-options", branch, remoteBase); err == nil {
				continue
			}
		}
		count, err := Git(ctx, repo, "rev-list", "--count", "--end-of-options", start+".."+branch)
		if err != nil {
			continue
		}
		found = append(found, unmergedPlan{other.ID, branch, count})
	}
	return found
}

// describeUnmerged reports work a starting point lacks. One sibling has a
// direct answer and is given it. Several do not: the correct start is a base
// plus integration merges, and `record order` derives that from `depends_on`
// alone — which is empty here, or this report would not exist. Pointing at it
// now would answer "start from the base", the very thing being questioned. So
// the several-sibling case names the declaration that makes ordering possible
// rather than one arbitrary branch that would quietly exclude the rest.
func describeUnmerged(found []unmergedPlan, repoID, plan, start string) string {
	if len(found) == 0 {
		return ""
	}
	parts := make([]string, 0, len(found))
	for _, item := range found {
		parts = append(parts, fmt.Sprintf("%s (%s, %s commit(s))", item.Plan, item.Branch, item.Count))
	}
	lead := fmt.Sprintf("%s does not contain work in %s from %s.", start, repoID, strings.Join(parts, ", "))
	if len(found) == 1 {
		return lead + fmt.Sprintf(" Start from %s to build on it, or continue if this plan is independent of it.", found[0].Branch)
	}
	return lead + fmt.Sprintf(" Record the ones this plan builds on with `record dependencies --id %s --depends-on ID`, then `record order` derives the start and any integration merges. Continue if this plan is independent of them.", plan)
}

// dependencyStartError names the exact starting point to pass. A predecessor
// whose branch does not exist yet has not been implemented and committed, which
// is a different fault than forgetting the flag, so it is reported as one.
func dependencyStartError(ctx context.Context, repo, plan string, start OrderStart) error {
	missing := []string{}
	for _, branch := range append([]string{start.Base}, start.Merge...) {
		if _, err := Git(ctx, repo, "show-ref", "--verify", "refs/heads/"+branch); err != nil {
			missing = append(missing, branch)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%s depends on work this repository does not have yet: %s missing; implement and commit those plans first, or name another starting point with --start",
			plan, strings.Join(missing, ", "))
	}
	if len(start.Merge) > 0 {
		return fmt.Errorf("%s depends on several plans in this repository: prepare with --start %s, then merge %s; run `record order` for the full shape",
			plan, start.Base, strings.Join(start.Merge, ", "))
	}
	return fmt.Errorf("%s depends on work in this repository: prepare with --start %s, or name another starting point explicitly; the base branch does not contain it",
		plan, start.Base)
}
