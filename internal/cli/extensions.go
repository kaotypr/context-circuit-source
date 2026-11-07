package cli

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// withStageHooks adds the extension blocks a CLI-fired stage returns to the
// command's result. Nothing is added where the workspace lists no extensions
// or none applied, so the result keeps the shape it always had.
func withStageHooks(ctx context.Context, s *workspace.Store, stage string, in *input, result any) any {
	repositories, scoped := stageRepositories(s, stage, in, result)
	blocks := s.StageHooks(ctx, stage, repositories, scoped)
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
		var ids []string
		if moved, ok := result.(map[string]any)["archived"].([]workspace.ArchivedRecord); ok {
			for _, record := range moved {
				ids = append(ids, record.ID)
			}
		}
		return planRepositories(ids...)
	}
	return nil, false
}

// runExtensionHooks resolves one hook point for a stage the CLI does not see.
func runExtensionHooks(ctx context.Context, s *workspace.Store, in *input) (any, error) {
	point, err := workspace.NewHookPoint(in.get("stage"), in.get("when"))
	if err != nil {
		return nil, err
	}
	if command, fired := workspace.CLIFiredStages[point.Stage]; fired && point.When == "after" {
		return nil, fmt.Errorf("%s is fired by `%s`: read the extensions block it returns rather than firing it again", point, command)
	}
	query := workspace.HookQuery{Point: point, Fetch: true}
	query.Member, _ = s.ActiveMember()
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
		if repo != "" {
			if !slices.Contains(record.Repositories, repo) {
				return nil, fmt.Errorf("%s does not name repository %s", plan, repo)
			}
			query.Repositories = []string{repo}
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
	case repo != "":
		cfg, err := s.Config()
		if err != nil {
			return nil, err
		}
		if _, ok := cfg.Repositories[repo]; !ok {
			return nil, fmt.Errorf("unknown repository: %s", repo)
		}
		query.Repositories, query.Scoped = []string{repo}, true
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
