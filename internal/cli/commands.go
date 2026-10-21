package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"

	assets "github.com/kaotypr/context-circuit-source"
	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// input is one command line's parsed options. A command's flags land here, and
// its run function reads them back, so a command is declared once in commands.
type input struct {
	name   string
	f      *flag.FlagSet
	values map[string]*string

	repos, dependencies, copyPaths, plans, roles, recordIDs listFlag
	band                                                    int
	archived, reuse, discard, reviewRequested, shared       bool
	localTiering, clearTone, clearRoles, all, confirm       bool

	// checkFailed is set by check, whose findings make the exit code 1.
	checkFailed bool
}

// get returns a string flag's value, or "" for a flag the command never defined.
func (in *input) get(key string) string {
	if p := in.values[key]; p != nil {
		return *p
	}
	return ""
}

// add defines string flags that default to empty, which makes them required
// unless the command lists them as optional.
func (in *input) add(names ...string) {
	for _, name := range names {
		in.values[name] = in.f.String(name, "", name)
	}
}

// command declares everything one command line needs in one place.
type command struct {
	flags []string // string flags; empty means missing unless listed in optional
	// extra defines the flags that are not plain strings: numbers, booleans,
	// lists, and strings with a default.
	extra func(in *input)
	// optional names string flags that may be empty. remote, default-branch and
	// url are optional for every command, see alwaysOptional.
	optional []string
	readOnly bool
	// setup marks the steps that settle setup, so they run before it is
	// complete; every other change waits for it. A command that is neither
	// setup nor readOnly is refused until setup is complete, which is the safe
	// direction for one added later.
	setup bool
	// splits decides whether the ID split is tried; nil means every changing
	// command makes it. The split is a change of its own, so a deletion preview
	// and init (which writes the new layout itself) must say they do not.
	splits func(in *input) bool
	// requiredNumber names a number flag that must be passed even when zero is
	// a real value. The string-flag check cannot see an omitted number.
	requiredNumber string
	makeRoot       bool // create the workspace root before opening it
	fromPath       bool // the root is --path, and no workspace pins a version there
	// precheck refuses a call its own command would refuse before the ID split
	// runs, so a mistyped command never moves a schema 2 workspace on its way to
	// failing.
	precheck func(ctx context.Context, s *workspace.Store, in *input) error
	run      func(ctx context.Context, s *workspace.Store, in *input) (any, error)
}

// alwaysOptional is true of every command that defines the flag. A URL is one
// shared detail the checkout usually already knows. Even a clone may omit it,
// because an ID this workspace already describes carries the source in the
// shared record, which is what joining one looks like.
var alwaysOptional = []string{"remote", "default-branch", "url"}

// done is the result of a command that has nothing to report but its outcome.
func done(s *workspace.Store, in *input, err error) (any, error) {
	return map[string]any{"operation": in.name, "workspace": s.Root, "ok": err == nil}, err
}

// act adapts a store method that only reports success into a run function.
func act(f func(ctx context.Context, s *workspace.Store, in *input) error) func(context.Context, *workspace.Store, *input) (any, error) {
	return func(ctx context.Context, s *workspace.Store, in *input) (any, error) {
		return done(s, in, f(ctx, s, in))
	}
}

func noSplit(*input) bool { return false }

var commands = map[string]command{
	"init": {
		flags: []string{"name", "purpose", "member", "member-name"},
		extra: func(in *input) {
			in.f.IntVar(&in.band, "member-band", 1, "allocation block for the first member")
			in.f.Var(&in.roles, "member-role", "informational role of the first member (repeatable)")
		},
		setup:    true,
		splits:   noSplit,
		makeRoot: true,
		run:      runInit,
	},
	"member add": {
		flags: []string{"id", "name", "language", "tone"},
		extra: func(in *input) {
			in.f.IntVar(&in.band, "band", 0, "allocation block for this member")
			in.f.Var(&in.roles, "role", "informational role of this member (repeatable)")
		},
		// A member writing in the workspace's usual language records nothing, and
		// a team happy with the general composition guidance records no register;
		// both fields exist for the member who is neither.
		optional: []string{"language", "tone"},
		setup:    true,
		precheck: func(_ context.Context, s *workspace.Store, in *input) error {
			// Nothing is written without a band the person chose, but the
			// refusal offers one so choosing costs nothing.
			if provided(in.f, "band") {
				return nil
			}
			suggestions, e := s.SuggestBands()
			if e != nil {
				return fmt.Errorf("pass --band N: %w", e)
			}
			return fmt.Errorf("pass --band N; the next free band is %d", suggestions.Next)
		},
		run: act(func(_ context.Context, s *workspace.Store, in *input) error {
			return s.AddMember(in.get("id"), in.get("name"), in.band, in.get("language"), in.get("tone"), in.roles)
		}),
	},
	"member roles": {
		flags: []string{"id"},
		extra: func(in *input) {
			in.f.Var(&in.roles, "role", "informational role, replacing those recorded (repeatable)")
			in.f.BoolVar(&in.clearRoles, "clear", false, "remove this member's recorded roles")
		},
		precheck: func(_ context.Context, _ *workspace.Store, in *input) error {
			return workspace.RolesCall(in.roles, in.clearRoles)
		},
		run: act(func(_ context.Context, s *workspace.Store, in *input) error {
			return s.SetMemberRoles(in.get("id"), in.roles, in.clearRoles)
		}),
	},
	"member language": {
		flags: []string{"id", "language"},
		run: act(func(_ context.Context, s *workspace.Store, in *input) error {
			return s.SetMemberLanguage(in.get("id"), in.get("language"))
		}),
	},
	"member tone": {
		flags: []string{"id", "tone"},
		extra: func(in *input) {
			in.f.BoolVar(&in.clearTone, "clear", false, "remove this member's recorded register")
		},
		optional: []string{"tone"},
		run: act(func(_ context.Context, s *workspace.Store, in *input) error {
			return s.SetMemberTone(in.get("id"), in.get("tone"), in.clearTone)
		}),
	},
	"member band": {
		flags: []string{"id"},
		extra: func(in *input) {
			in.f.IntVar(&in.band, "band", 0, "allocation block; 0 is a band like any other")
		},
		// Omitting the flag would otherwise move the member to band 0 silently,
		// which a caller who forgot it cannot have meant.
		requiredNumber: "band",
		setup:          true,
		run: act(func(_ context.Context, s *workspace.Store, in *input) error {
			return s.SetMemberBand(in.get("id"), in.band)
		}),
	},
	"member use": {
		flags: []string{"id"},
		setup: true,
		run: act(func(_ context.Context, s *workspace.Store, in *input) error {
			return s.UseMember(in.get("id"))
		}),
	},
	"member list": {
		readOnly: true,
		run: func(_ context.Context, s *workspace.Store, _ *input) (any, error) {
			return s.Members()
		},
	},
	"repo connect": {
		flags: []string{"id", "path", "base", "url", "default-branch"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.Connect(ctx, in.get("id"), in.get("path"), in.get("base"), in.get("url"), in.get("default-branch"))
		}),
	},
	"repo clone": {
		flags: []string{"id", "path", "base", "url", "default-branch"},
		// A clone needs a source: --url, or a shared record already naming one.
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.CreateRepository(ctx, in.get("id"), in.get("path"), in.get("base"), in.get("url"), in.get("default-branch"), true)
		}),
	},
	"repo init": {
		flags: []string{"id", "path", "base", "default-branch"},
		// An init never takes a source, so the store tells it from a clone.
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.CreateRepository(ctx, in.get("id"), in.get("path"), in.get("base"), in.get("url"), in.get("default-branch"), false)
		}),
	},
	"repo remote": {
		flags: []string{"id", "url", "default-branch"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.SetRemote(ctx, in.get("id"), in.get("url"), in.get("default-branch"))
		}),
	},
	"repo base": {
		flags: []string{"id", "branch"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.SetBase(ctx, in.get("id"), in.get("branch"))
		}),
	},
	"repo fetch": {
		flags: []string{"id", "remote"},
		run:   runRepoFetch,
	},
	"repo inspect": {
		flags:    []string{"id"},
		readOnly: true,
		run:      runRepoInspect,
	},
	"repo relate": {
		flags: []string{"from", "to", "description"},
		run: act(func(_ context.Context, s *workspace.Store, in *input) error {
			return s.Relate(in.get("from"), in.get("to"), in.get("description"))
		}),
	},
	"workspace connect": {
		flags: []string{"path", "base", "url", "default-branch"},
		// The workspace's own checkout is the workspace root unless a workspace
		// kept inside a larger repository names the root above it.
		optional: []string{"path"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.ConnectWorkspace(ctx, in.get("path"), in.get("base"), in.get("url"), in.get("default-branch"))
		}),
	},
	"workspace base": {
		flags: []string{"branch"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.SetWorkspaceBase(ctx, in.get("branch"))
		}),
	},
	"workspace remote": {
		flags: []string{"url", "default-branch"},
		// Either half of the shared record is a whole call; an empty one is
		// refused by the store, where the alternative can be named.
		optional: []string{"default-branch"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.SetWorkspaceRemote(ctx, in.get("url"), in.get("default-branch"))
		}),
	},
	// A borrowed repository names its own index, and a layout this product
	// recognizes is detected rather than demanded. Obtaining one lands under
	// repositories/<id> unless the caller keeps it somewhere else, so only a
	// connect must name its path.
	"knowledge connect": {
		flags:    []string{"id", "path", "url", "default-branch", "index"},
		optional: []string{"index"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.ConnectKnowledge(ctx, in.get("id"), in.get("path"), in.get("url"), in.get("default-branch"), in.get("index"))
		}),
	},
	"knowledge clone": {
		flags:    []string{"id", "path", "url", "default-branch", "index"},
		optional: []string{"index", "path"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.CloneKnowledge(ctx, in.get("id"), in.get("path"), in.get("url"), in.get("default-branch"), in.get("index"))
		}),
	},
	"knowledge remote": {
		flags:    []string{"id", "url", "default-branch", "index"},
		optional: []string{"index"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.SetKnowledgeRemote(ctx, in.get("id"), in.get("url"), in.get("default-branch"), in.get("index"))
		}),
	},
	"knowledge sync": {
		flags: []string{"id"},
		run: func(ctx context.Context, s *workspace.Store, in *input) (any, error) {
			return s.SyncKnowledge(ctx, in.get("id"))
		},
	},
	"knowledge list": {
		readOnly: true,
		run: func(ctx context.Context, s *workspace.Store, _ *input) (any, error) {
			return s.KnowledgeStates(ctx)
		},
	},
	"record create": {
		flags: []string{"kind", "slug", "title", "intent"},
		extra: func(in *input) {
			in.f.Var(&in.repos, "repo", "repository (repeatable)")
			in.f.Var(&in.dependencies, "depends-on", "plan dependency (repeatable)")
		},
		optional: []string{"intent"},
		run: func(ctx context.Context, s *workspace.Store, in *input) (any, error) {
			return s.CreateRecord(ctx, in.get("kind"), in.get("slug"), in.get("title"), in.get("intent"), in.repos, in.dependencies)
		},
	},
	"record show": {
		flags:    []string{"id"},
		readOnly: true,
		run: func(_ context.Context, s *workspace.Store, in *input) (any, error) {
			return s.FindRecord(in.get("id"))
		},
	},
	"record list": {
		extra: func(in *input) {
			in.f.BoolVar(&in.archived, "archived", false, "include archived records")
		},
		readOnly: true,
		run: func(_ context.Context, s *workspace.Store, in *input) (any, error) {
			return s.ListRecords(in.archived)
		},
	},
	"record approve": {
		flags: []string{"id", "text"},
		run:   runApprove,
	},
	"record complete": {
		flags: []string{"id", "text"},
		run:   runComplete,
	},
	"record delete": {
		flags: []string{"member"},
		extra: func(in *input) {
			in.f.BoolVar(&in.all, "all", false, "every record, while the roster lists one member")
			in.f.Var(&in.recordIDs, "id", "intent or plan ID to delete (repeatable)")
			in.f.BoolVar(&in.confirm, "confirm", false, "delete, after the person agreed to the previewed list")
		},
		// Deleting every record names no member; the store refuses a half-given
		// pair where it can say why.
		optional: []string{"member"},
		// A deletion preview promises to write nothing.
		splits: func(in *input) bool { return in.confirm },
		// The preview writes nothing and reads either ID-file layout, so a
		// deletion it refuses is refused before the split.
		precheck: func(ctx context.Context, s *workspace.Store, in *input) error {
			if !in.confirm {
				return nil
			}
			_, err := s.DeleteRecords(ctx, in.get("member"), in.all, in.recordIDs, false)
			return err
		},
		run: func(ctx context.Context, s *workspace.Store, in *input) (any, error) {
			return s.DeleteRecords(ctx, in.get("member"), in.all, in.recordIDs, in.confirm)
		},
	},
	"record dependencies": {
		flags: []string{"id"},
		extra: func(in *input) {
			in.f.Var(&in.dependencies, "depends-on", "plan dependency (repeatable)")
		},
		run: act(func(_ context.Context, s *workspace.Store, in *input) error {
			return s.SetDependencies(in.get("id"), in.dependencies)
		}),
	},
	"record order": {
		flags: []string{"intent"},
		extra: func(in *input) {
			in.f.Var(&in.plans, "plan", "selected plan (repeatable)")
			in.values["mode"] = in.f.String("mode", "auto", "auto, waves, or linear")
		},
		optional: []string{"intent"},
		run: func(_ context.Context, s *workspace.Store, in *input) (any, error) {
			return s.Order(in.get("intent"), in.get("mode"), in.plans...)
		},
	},
	"context find": {
		flags: []string{"query", "repo"},
		// A search takes a query and a reconciliation takes a repository. Either
		// one is a whole call, so both are optional here and an empty call is
		// refused by the run function, where the alternative can be named.
		optional: []string{"query", "repo"},
		readOnly: true,
		run:      runContextFind,
	},
	"worktree prepare": {
		flags: []string{"repo", "plan", "branch", "start", "path"},
		extra: func(in *input) {
			in.values["copy-mode"] = in.f.String("copy-mode", "auto", "runtime file reuse: auto, required, copy, off")
			in.f.Var(&in.copyPaths, "copy-path", "additional ignored runtime path (repeatable)")
			in.f.BoolVar(&in.reuse, "reuse", false, "explicitly reuse existing work")
		},
		optional: []string{"plan", "branch", "start", "path"},
		run: func(ctx context.Context, s *workspace.Store, in *input) (any, error) {
			return s.Prepare(ctx, in.get("repo"), in.get("plan"), in.get("branch"), in.get("start"), in.get("path"), in.reuse, workspace.ReuseOptions{Mode: in.get("copy-mode"), Paths: in.copyPaths})
		},
	},
	"worktree list": {
		flags:    []string{"repo"},
		readOnly: true,
		run: func(ctx context.Context, s *workspace.Store, in *input) (any, error) {
			return s.Worktrees(ctx, in.get("repo"))
		},
	},
	"worktree inspect": {
		flags:    []string{"repo", "path"},
		readOnly: true,
		run: func(ctx context.Context, s *workspace.Store, in *input) (any, error) {
			return s.InspectWorktree(ctx, in.get("repo"), in.get("path"))
		},
	},
	"worktree repair": {
		flags: []string{"repo", "path"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.RepairWorktree(ctx, in.get("repo"), in.get("path"))
		}),
	},
	"worktree move": {
		flags: []string{"repo", "path", "to"},
		run: act(func(ctx context.Context, s *workspace.Store, in *input) error {
			return s.MoveWorktree(ctx, in.get("repo"), in.get("path"), in.get("to"))
		}),
	},
	"worktree remove": {
		flags: []string{"repo", "path"},
		extra: func(in *input) {
			in.f.BoolVar(&in.discard, "discard", false, "explicitly discard worktree files")
		},
		run: runWorktreeRemove,
	},
	"worktree prune": {
		flags: []string{"repo"},
		run:   runWorktreePrune,
	},
	"template export": {
		flags:    []string{"path"},
		readOnly: true,
		makeRoot: true,
		fromPath: true,
		run:      runExport,
	},
	"agent settings": {
		readOnly: true,
		run: func(_ context.Context, s *workspace.Store, _ *input) (any, error) {
			return s.EffectiveAgentSettings()
		},
	},
	"agent configure": {
		flags: []string{"host", "role", "model", "effort"},
		extra: func(in *input) {
			in.f.BoolVar(&in.localTiering, "local", false, "write this machine's override instead of the shared setting")
		},
		run: act(func(_ context.Context, s *workspace.Store, in *input) error {
			return s.ConfigureAgent(in.get("host"), in.get("role"), in.get("model"), in.get("effort"), in.localTiering)
		}),
	},
	"agent setup": {
		flags: []string{"host"},
		// A workspace is opened in more than one host, so setup covers them all
		// unless the caller narrows it to one.
		optional: []string{"host"},
		run: func(_ context.Context, s *workspace.Store, in *input) (any, error) {
			return s.SetupAgents(in.get("host"))
		},
	},
	"agent dispatch": {
		flags: []string{"host", "role", "task", "path", "plan", "intent", "repo"},
		extra: func(in *input) {
			in.f.BoolVar(&in.shared, "shared", false, "several workers share this worktree")
			in.f.BoolVar(&in.reviewRequested, "review-requested", false, "user explicitly requested independent review")
		},
		// A planner takes --intent and a worker takes --plan; neither is
		// universal, so the command accepts both as optional and refuses the
		// wrong pairing by role, where the reason can be stated. A planner and a
		// whole-plan worker are assigned by the record quoted in the brief; only
		// a role without one, or a worker taking a slice of a shared plan, has
		// something left to say, so --task is optional too.
		optional: []string{"intent", "plan", "repo", "task"},
		readOnly: true,
		run: func(_ context.Context, s *workspace.Store, in *input) (any, error) {
			return s.DispatchAgent(in.get("host"), in.get("role"), in.get("task"), in.get("path"), in.get("plan"), in.get("intent"), in.shared, in.reviewRequested, in.get("repo"))
		},
	},
	"status": {
		readOnly: true,
		run: func(ctx context.Context, s *workspace.Store, _ *input) (any, error) {
			return s.Status(ctx)
		},
	},
	"check": {
		readOnly: true,
		run: func(ctx context.Context, s *workspace.Store, in *input) (any, error) {
			issues, e := s.Check(ctx)
			in.checkFailed = len(issues) > 0
			return map[string]any{"issues": issues, "ok": !in.checkFailed}, e
		},
	},
}

func runInit(_ context.Context, s *workspace.Store, in *input) (any, error) {
	files, err := assets.Files()
	if err != nil {
		return nil, err
	}
	readme, err := assets.Readme()
	if err != nil {
		return nil, err
	}
	return done(s, in, s.Init(files, readme, in.get("name"), in.get("purpose"), in.get("member"), in.get("member-name"), in.band, in.roles))
}

func runExport(_ context.Context, s *workspace.Store, in *input) (any, error) {
	files, err := assets.Files()
	if err != nil {
		return nil, err
	}
	return done(s, in, s.Export(files))
}

func runRepoInspect(ctx context.Context, s *workspace.Store, in *input) (any, error) {
	checkout, err := s.Repository(ctx, in.get("id"))
	if err != nil {
		return nil, err
	}
	snapshot, err := workspace.Inspect(ctx, checkout.Path)
	if err != nil {
		return nil, err
	}
	snapshot.BaseBranch = checkout.BaseBranch
	return snapshot, nil
}

func runRepoFetch(ctx context.Context, s *workspace.Store, in *input) (any, error) {
	checkout, err := s.Repository(ctx, in.get("id"))
	if err != nil {
		return nil, err
	}
	remote := in.get("remote")
	if remote == "" {
		remote = "origin"
	}
	if err := workspace.Name(remote); err != nil {
		return nil, err
	}
	_, err = workspace.Git(ctx, checkout.Path, "fetch", "--", remote)
	return done(s, in, err)
}

func runApprove(ctx context.Context, s *workspace.Store, in *input) (any, error) {
	if err := s.Note(in.get("id"), in.get("text"), "approve"); err != nil {
		return nil, err
	}
	// Approval hands work back rather than finishing it, and a bare
	// `ok: true` reads as the request being done. The next step to
	// recommend is said here, where the caller is standing, rather
	// than only in instructions read once at the start of a session.
	result := map[string]any{
		"approved":             in.get("id"),
		"planning_recommended": "approval is recorded; planning is the next step: recommend creating linked pNNNN-slug/plan.md folders to the person and stop, unless the same request already asked for plans, in which case inspect real code, create them, present the complete plans and stop",
	}
	recommendCommit(ctx, s, result, in.get("id"))
	return result, nil
}

func runComplete(ctx context.Context, s *workspace.Store, in *input) (any, error) {
	if err := s.Note(in.get("id"), in.get("text"), "complete"); err != nil {
		return nil, err
	}
	entries, e := s.KnowledgeCandidates(in.get("id"))
	result := map[string]any{"completed": in.get("id"), "knowledge_candidates": entries}
	// A bare list of catalog lines reads as information, and a caller
	// that has just been told "completed" treats the request as done.
	// Naming the outstanding act beside the data is what `setup_required`
	// does for a dispatch, so completion says it the same way.
	if len(entries) > 0 {
		result["reconcile_required"] = "judge each entry above against what this plan changed: edit the note and its catalog entry together and move its reviewed date, or record in the completion note that it changed nothing"
	}
	ids := []string{in.get("id")}
	if plan, found := s.FindRecord(in.get("id")); found == nil && plan.Intent != "" {
		ids = append(ids, plan.Intent)
	}
	recommendCommit(ctx, s, result, ids...)
	return result, e
}

func runContextFind(ctx context.Context, s *workspace.Store, in *input) (any, error) {
	// Completion names the knowledge a plan made due; a change made
	// without a plan has no record to name it, so the repository asks
	// for the same candidates and is told the same act is outstanding.
	query, repository := in.get("query"), in.get("repo")
	if query == "" && repository == "" {
		return nil, errors.New("pass --query to search, or --repo for the entries claiming one repository")
	}
	if repository == "" {
		return s.FindContext(ctx, query)
	}
	if query != "" {
		return nil, errors.New("pass --query or --repo, not both")
	}
	cfg, err := s.Config()
	if err != nil {
		return nil, err
	}
	if _, known := cfg.Repositories[repository]; !known {
		return nil, fmt.Errorf("unknown repository: %s", repository)
	}
	entries, err := s.CatalogEntriesFor([]string{repository})
	result := map[string]any{"repository": repository, "knowledge_candidates": entries}
	if len(entries) > 0 {
		result["reconcile_required"] = "judge each entry above against what this change actually altered: edit the note and its catalog entry together and move its reviewed date, or report that it changed nothing"
	}
	return result, err
}

func runWorktreeRemove(ctx context.Context, s *workspace.Store, in *input) (any, error) {
	// The empty .worktrees/ folders the removal left are removed with it and
	// listed, the way prune lists what it dropped. One that could not be
	// removed is a warning, not a failure: the worktree is already gone.
	removed, warnings, err := s.RemoveWorktree(ctx, in.get("repo"), in.get("path"), in.discard)
	if err != nil {
		return nil, err
	}
	result := map[string]any{"operation": in.name, "workspace": s.Root, "ok": true, "removed_directories": removed}
	if len(warnings) > 0 {
		result["warnings"] = warnings
	}
	return result, nil
}

func runWorktreePrune(ctx context.Context, s *workspace.Store, in *input) (any, error) {
	dropped, err := s.PruneWorktreeAssociations(ctx, in.get("repo"))
	if err != nil {
		return nil, err
	}
	return map[string]any{"operation": in.name, "workspace": s.Root, "ok": true, "repository": in.get("repo"), "dropped": dropped}, nil
}
