package cli

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// stageStopped is a CLI-fired stage whose before point stopped it: the command
// was not run and wrote nothing, and its result is the extension block that
// says why, which is still shown.
type stageStopped struct {
	result  any
	command string
	point   string
}

func (e *stageStopped) Error() string {
	return e.command + " was not run, and nothing was changed: an extension stopped " + e.point + "; take the extensions block to the person"
}

// stageGate refuses, before any before hook runs, what the command itself
// would refuse first: an incomplete setup and its precheck. A hook therefore
// never runs for a stage that was never going to happen.
func stageGate(ctx context.Context, s *workspace.Store, cmd command, in *input) error {
	if !cmd.readOnly && !cmd.setup {
		if err := s.RequireSetup(); err != nil {
			return fmt.Errorf("setup is incomplete, so nothing was changed: %w", err)
		}
	}
	if cmd.precheck != nil {
		return cmd.precheck(ctx, s, in)
	}
	return nil
}

// stageHook fires one point of the stage a CLI command fires, outside the
// lock. The before point reads only the command line; the after point reads
// only what the command returned and the command line, never state the
// command's precheck saw, which another command may have changed since.
func stageHook(ctx context.Context, s *workspace.Store, stage, when string, in *input, result any) *workspace.HookResult {
	repositories, scoped := stageRepositories(s, stage, in, result)
	script := workspace.ScriptContext{Worktrees: map[string]string{}}
	switch stage {
	case "approve":
		script.Record = s.ScriptRecordFor(in.get("id"))
	case "prepare":
		if plan := in.get("plan"); plan != "" {
			script.Record = s.ScriptRecordFor(plan)
		}
		if tree, ok := result.(workspace.Worktree); ok && tree.Path != "" {
			script.Worktrees[in.get("repo")] = tree.Path
			script.Dir = tree.Path
		}
	case "complete":
		script.Record = s.ScriptRecordFor(in.get("id"))
		script.Worktrees = s.PlanWorktrees(in.get("id"))
	case "archive":
		for _, id := range archiveIDs(in, result) {
			if record := s.ScriptRecordFor(id); record != nil {
				script.Records = append(script.Records, *record)
			}
		}
		if len(script.Records) == 1 {
			script.Record = &script.Records[0]
		}
	}
	return s.StageHook(ctx, workspace.HookPoint{When: when, Stage: stage}, repositories, scoped, script)
}

// withStageHooks adds the extension blocks a CLI-fired stage returned to the
// command's result. Nothing is added where none applied, so the result keeps
// the shape it always had.
func withStageHooks(result any, blocks []workspace.HookResult) any {
	switch typed := result.(type) {
	case map[string]any:
		if len(blocks) > 0 {
			typed["extensions"] = blocks
		}
		return typed
	case workspace.Worktree:
		return workspace.PreparedWorktree{Worktree: typed, Extensions: blocks}
	}
	return result
}

// archiveIDs is what an archive fires for: the records it moved once it ran,
// and before that the IDs it was given.
func archiveIDs(in *input, result any) []string {
	if typed, ok := result.(map[string]any); ok {
		var ids []string
		if moved, ok := typed["archived"].([]workspace.ArchivedRecord); ok {
			for _, record := range moved {
				ids = append(ids, record.ID)
			}
		}
		return ids
	}
	return in.recordIDs
}

// stageRepositories is the repository list a stage's hooks are scoped by. An
// intent carries no structured list, so approval is unscoped; a plan's are its
// frontmatter's; a worktree's is the one it was prepared in.
func stageRepositories(s *workspace.Store, stage string, in *input, result any) ([]string, bool) {
	planRepositories := func(ids ...string) ([]string, bool) {
		var repositories []string
		for _, id := range ids {
			if record, err := s.FindRecord(id); err == nil {
				for _, repo := range record.Repositories {
					if !slices.Contains(repositories, repo) {
						repositories = append(repositories, repo)
					}
				}
			}
		}
		return repositories, len(repositories) > 0
	}
	switch stage {
	case "prepare":
		return []string{in.get("repo")}, true
	case "complete":
		return planRepositories(in.get("id"))
	case "archive":
		return planRepositories(archiveIDs(in, result)...)
	}
	return nil, false
}

// repositoryStages are the stages whose scripts run in the worktree when the
// point fires for one repository.
var repositoryStages = []string{"prepare", "implement", "review", "deliver"}

// runExtensionHooks resolves one hook point for a stage the CLI does not see.
func runExtensionHooks(ctx context.Context, s *workspace.Store, in *input) (any, error) {
	point, err := workspace.NewHookPoint(in.get("stage"), in.get("when"))
	if err != nil {
		return nil, err
	}
	if command, fired := workspace.CLIFiredStages[point.Stage]; fired && point.When == "after" {
		return nil, fmt.Errorf("%s is fired by `%s`: read the extensions block it returns rather than firing it again", point, command)
	}
	query := workspace.HookQuery{Point: point, Fetch: true, Run: true, Confirm: in.get("confirm")}
	query.Member, _ = s.ActiveMember()
	if command, fired := workspace.CLIFiredStages[point.Stage]; fired {
		// Its scripts run when the command fires the point, outside the lock
		// and before it writes; here it returns the instructions only.
		if query.Confirm != "" {
			return nil, fmt.Errorf("%s is fired by `%s`, which runs its scripts; --confirm is for a stage the agent fires", point, command)
		}
		query.Run, query.NotRun = false, "script not run here: it runs when `"+command+"` fires "+point.String()
	}
	plan, intent, repo := in.get("plan"), in.get("intent"), in.get("repo")
	if plan != "" && intent != "" {
		return nil, errors.New("pass --plan or --intent, not both")
	}
	switch {
	case plan != "":
		record, err := s.FindRecord(plan)
		if err != nil {
			return nil, err
		}
		if record.ID[0] != 'p' {
			return nil, errors.New("--plan needs a plan ID")
		}
		query.Repositories, query.Scoped = record.Repositories, len(record.Repositories) > 0
		query.Script.Record = s.ScriptRecordFor(record.ID)
		query.Script.Worktrees = s.PlanWorktrees(record.ID)
		if repo != "" {
			if !slices.Contains(record.Repositories, repo) {
				return nil, fmt.Errorf("%s does not name repository %s", plan, repo)
			}
			query.Repositories = []string{repo}
			if tree, ok := query.Script.Worktrees[repo]; ok {
				query.Script.Worktrees = map[string]string{repo: tree}
			} else {
				query.Script.Worktrees = map[string]string{}
			}
		}
	case intent != "":
		record, err := s.FindRecord(intent)
		if err != nil {
			return nil, err
		}
		if record.ID[0] != 'i' {
			return nil, errors.New("--intent needs an intent ID")
		}
		if repo != "" {
			return nil, errors.New("an intent carries no repository list, so --repo goes with --plan or alone")
		}
		query.Script.Record = s.ScriptRecordFor(record.ID)
	case repo != "":
		cfg, err := s.Config()
		if err != nil {
			return nil, err
		}
		if _, ok := cfg.Repositories[repo]; !ok {
			return nil, fmt.Errorf("unknown repository: %s", repo)
		}
		query.Repositories, query.Scoped = []string{repo}, true
		// A change made directly works in the bound checkout, which is then
		// the repository's worktree.
		if checkout, err := s.Repository(ctx, repo); err == nil {
			query.Script.Worktrees = map[string]string{repo: checkout.Path}
		}
	}
	if slices.Contains(repositoryStages, point.Stage) && len(query.Script.Worktrees) == 1 {
		for _, tree := range query.Script.Worktrees {
			query.Script.Dir = tree
		}
	}
	result, listed, err := s.ResolveHooks(ctx, query)
	if err != nil {
		return nil, err
	}
	out := map[string]any{"extensions": []workspace.HookResult{result}}
	if !listed {
		out["note"] = "this workspace lists no extensions, so no hook applies"
	}
	return out, nil
}
