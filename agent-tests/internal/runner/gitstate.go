package runner

import (
	"context"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

// writeGitState writes artifacts/git/<key>.txt for each repository. One
// command failing is written into that file and does not stop the run.
func (rt *runtime) writeGitState(ctx context.Context, w world.World, before scenario.Snapshot, sc scenario.Scenario, dest string) error {
	red := rt.redact()
	red.runRoot = w.RunRoot
	for _, repo := range snapshotRepos(w) {
		if !repoHasGit(repo.path) {
			continue
		}
		body := rt.gitStateText(ctx, w, repo.path, before.Refs[repo.name])
		if err := writeText(filepath.Join(dest, repo.name+".txt"), body, red); err != nil {
			return err
		}
	}
	if sc.ID == "execute" {
		return rt.writeExecuteGit(ctx, w, dest, red)
	}
	return nil
}

// writeExecuteGit records the plan branch's diff stat and full commit
// messages. The branch name is the one worktree prepare creates.
func (rt *runtime) writeExecuteGit(ctx context.Context, w world.World, dest string, red textRedactor) error {
	dir := w.API
	if dir == "" {
		dir = filepath.Join(w.WS, "repositories", "api")
	}
	sections := []struct {
		file  string
		title string
		argv  []string
	}{
		{
			file:  "api-diff-stat.txt",
			title: "git diff --stat main...cc/p1000/api",
			argv:  []string{"git", "-C", dir, "diff", "--stat", "main...cc/p1000/api"},
		},
		{
			file:  "api-log.txt",
			title: "git log main..cc/p1000/api --format=%B",
			argv:  []string{"git", "-C", dir, "log", "main..cc/p1000/api", "--format=%B"},
		},
	}
	for _, section := range sections {
		var b strings.Builder
		b.WriteString(section.title)
		b.WriteByte('\n')
		out, err := rt.execIsolated(ctx, w, dir, section.argv...)
		if err != nil {
			b.WriteString(err.Error())
			if !strings.HasSuffix(err.Error(), "\n") {
				b.WriteByte('\n')
			}
		} else {
			b.WriteString(out)
			if out != "" && !strings.HasSuffix(out, "\n") {
				b.WriteByte('\n')
			}
		}
		if err := writeText(filepath.Join(dest, section.file), b.String(), red); err != nil {
			return err
		}
	}
	return nil
}

func (rt *runtime) gitStateText(ctx context.Context, w world.World, dir, beforeRefs string) string {
	logTitle, logArgv := gitLogCommand(dir, beforeObjects(beforeRefs))
	sections := []struct {
		title string
		argv  []string
	}{
		{"git status --porcelain", []string{"git", "-C", dir, "status", "--porcelain"}},
		{"git worktree list", []string{"git", "-C", dir, "worktree", "list"}},
		{"git branch -a -vv", []string{"git", "-C", dir, "branch", "-a", "-vv"}},
		{logTitle, logArgv},
		{"git diff --stat", []string{"git", "-C", dir, "diff", "--stat"}},
	}
	var b strings.Builder
	for i, section := range sections {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(section.title)
		b.WriteByte('\n')
		out, err := rt.execIsolated(ctx, w, dir, section.argv...)
		if err != nil {
			text := err.Error()
			b.WriteString(text)
			if !strings.HasSuffix(text, "\n") {
				b.WriteByte('\n')
			}
			continue
		}
		b.WriteString(out)
		if out != "" && !strings.HasSuffix(out, "\n") {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func gitLogCommand(dir string, shas []string) (string, []string) {
	argv := []string{"git", "-C", dir, "log", "--all"}
	title := "git log --all"
	if len(shas) > 0 {
		argv = append(argv, "--not")
		argv = append(argv, shas...)
		title += " --not " + strings.Join(shas, " ")
	}
	argv = append(argv, "--format=%h %s")
	title += " --format=%h %s"
	return title, argv
}

func beforeObjects(text string) []string {
	seen := map[string]struct{}{}
	var objects []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		_, object, ok := strings.Cut(line, " ")
		if !ok {
			continue
		}
		object = strings.TrimSpace(object)
		if object == "" {
			continue
		}
		if _, ok := seen[object]; ok {
			continue
		}
		seen[object] = struct{}{}
		objects = append(objects, object)
	}
	sort.Strings(objects)
	return objects
}
