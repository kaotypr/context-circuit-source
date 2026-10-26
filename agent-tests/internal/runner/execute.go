package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

func (rt *runtime) one(ctx context.Context, item planned, emptyWorkspace bool) (err error) {
	verdict := "harness-error"
	var parsed hosts.Parsed
	var results []scenario.CheckResult
	var incomplete bool
	var wall time.Duration
	var drops []string
	var runRoot string
	runDir := filepath.Join(rt.opts.OutDir, rt.id, "runs", item.Scenario.ID, item.Host, item.Side, fmt.Sprintf("r%d", item.RunN))
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic: %v", rec)
		}
		if err != nil {
			verdict = scenario.Verdict(nil, err, false)
		}
		if len(drops) > 0 {
			_ = writeText(filepath.Join(runDir, "drop.log"), strings.Join(drops, "\n")+"\n", rt.redact())
		}
		digest, _ := record.Digest(runDir)
		line := rt.recordLine(item, parsed, verdict, results, wall, digest, err)
		if werr := record.Append(filepath.Join(rt.opts.OutDir, rt.id, "runs.jsonl"), line); werr != nil && err == nil {
			err = werr
		}
		if runRoot != "" && !rt.opts.Keep && verdict != "harness-error" {
			_ = os.RemoveAll(runRoot)
		}
	}()
	if err = os.MkdirAll(runDir, 0o755); err != nil {
		return err
	}

	runRoot, err = world.NewRunRoot(rt.opts.ParentTemp)
	if err != nil {
		return err
	}
	home := ""
	if rt.opts.Caller != nil {
		home = rt.opts.Caller["HOME"]
	}
	if err = world.Guard(runRoot, rt.opts.Checkout, home); err != nil {
		return err
	}
	var w world.World
	if emptyWorkspace {
		w, err = emptyWorld(runRoot)
		if err != nil {
			return err
		}
	} else {
		src := filepath.Join(runRoot, "src")
		if err = world.ExtractArchive(ctx, rt.opts.Checkout, item.Ref, src); err != nil {
			return err
		}
		cli, buildErr := world.BuildCLI(ctx, src, runRoot, rt.pin)
		if buildErr != nil {
			return buildErr
		}
		w, err = world.Setup(ctx, world.SetupInput{
			RunRoot:     runRoot,
			SrcDir:      src,
			CLIPath:     cli,
			FixturesDir: rt.opts.FixturesDir,
			Pin:         rt.pin,
			PATH:        callerPath(rt.opts.Caller),
		})
		if err != nil {
			return err
		}
		if err = rt.runSnippets(ctx, w, item.Scenario); err != nil {
			return err
		}
	}
	before, err := rt.snapshot(ctx, w)
	if err != nil {
		return err
	}
	adapter, err := hosts.ForName(item.Config.Name)
	if err != nil {
		return err
	}
	prompt := expandPrompt(item.Scenario.Prompt, w)
	launch, err := adapter.Build(hosts.LaunchInput{
		Prompt:    prompt,
		Model:     item.Config.Model,
		Effort:    item.Config.Effort,
		Workspace: w.WS,
		RunRoot:   w.RunRoot,
		SessionID: newSessionID(),
		Cfg:       item.Config,
		Caller:    rt.opts.Caller,
		DecoyDir:  w.Decoy,
	})
	if err != nil {
		return err
	}
	env := mergeMissing(launch.Env, w.Env)
	for _, dir := range launch.MkDirs {
		if err = os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	timeout := time.Duration(item.Config.TimeoutSeconds) * time.Second
	for i, argv := range launch.SetupArgv {
		stdin := ""
		if i < len(launch.SetupStdin) {
			stdin = launch.SetupStdin[i]
		}
		res, runErr := proc.Run(ctx, proc.Spec{Argv: argv, Dir: launch.Dir, Env: env, Stdin: stdin, Timeout: timeout})
		if runErr != nil {
			return fmt.Errorf("host setup: %w\n%s", runErr, res.Stderr)
		}
		if res.ExitCode != 0 {
			return fmt.Errorf("host setup: exit %d\n%s", res.ExitCode, res.Stderr)
		}
	}
	streamPath := filepath.Join(runRoot, "host-stream")
	res, err := proc.Run(ctx, proc.Spec{
		Argv:       launch.Argv,
		Dir:        launch.Dir,
		Env:        env,
		Stdin:      launch.Stdin,
		Timeout:    timeout,
		StreamPath: streamPath,
	})
	if err != nil {
		return fmt.Errorf("host: %w\n%s", err, res.Stderr)
	}
	wall = res.Wall
	if res.TimedOut {
		_ = rt.publishHost(runDir, w, res, streamPath, nil, &drops)
		return fmt.Errorf("host timed out")
	}
	if res.ExitCode != 0 {
		_ = rt.publishHost(runDir, w, res, streamPath, nil, &drops)
		return fmt.Errorf("host: exit %d\n%s", res.ExitCode, res.Stderr)
	}
	parsed, err = parseStream(adapter, streamPath, res.Stdout)
	if err != nil {
		return err
	}
	after, err := rt.snapshot(ctx, w)
	if err != nil {
		return err
	}
	checkOut, checkErr := rt.runCheck(ctx, w)
	if checkErr != nil {
		return checkErr
	}
	results, incomplete = grade(item.Scenario, scenario.CheckContext{
		World: w, Scenario: item.Scenario, Parsed: parsed, Before: before, After: after,
	})
	if !parsed.Terminal {
		incomplete = true
	}
	verdict = scenario.Verdict(results, nil, incomplete)
	red := rt.redact()
	red.runRoot = w.RunRoot
	if err = writeText(filepath.Join(runDir, "snapshots", "before.txt"), formatSnap(before), red); err != nil {
		return err
	}
	if err = writeText(filepath.Join(runDir, "snapshots", "after.txt"), formatSnap(after), red); err != nil {
		return err
	}
	if err = writeBytes(filepath.Join(runDir, "check.json"), checkOut, red, &drops); err != nil {
		return err
	}
	if err = rt.publishHost(runDir, w, res, streamPath, &parsed, &drops); err != nil {
		return err
	}
	if err = rt.copyChanged(ctx, w, filepath.Join(runDir, "artifacts", "files"), &drops); err != nil {
		return err
	}
	if err = rt.copyDesign(w, item.Scenario, filepath.Join(runDir, "artifacts", "design"), &drops); err != nil {
		return err
	}
	return nil
}

func (rt *runtime) runSnippets(ctx context.Context, w world.World, sc scenario.Scenario) error {
	env := snippetEnv(w, callerPath(rt.opts.Caller))
	for i, snippet := range sc.Setup {
		res, err := proc.Run(ctx, proc.Spec{
			Argv: []string{"sh", "-c", snippet},
			Dir:  w.WS,
			Env:  env,
		})
		if err != nil {
			return fmt.Errorf("setup snippet %d: %w\n%s", i+1, err, res.Stderr)
		}
		if res.ExitCode != 0 {
			return fmt.Errorf("setup snippet %d: exit %d\n%s", i+1, res.ExitCode, res.Stderr)
		}
	}
	return nil
}

func snippetEnv(w world.World, path string) []string {
	extra := map[string]string{}
	for key, val := range w.Env {
		extra[key] = val
	}
	extra["HOME"] = filepath.Join(w.RunRoot, "setup-home")
	return hosts.BaseEnv(hosts.LaunchInput{
		RunRoot:  w.RunRoot,
		DecoyDir: w.Decoy,
		Caller:   map[string]string{"PATH": path},
	}, extra)
}

func expandPrompt(prompt string, w world.World) string {
	prompt = strings.ReplaceAll(prompt, "$WS", w.WS)
	prompt = strings.ReplaceAll(prompt, "$RUN", w.RunRoot)
	return prompt
}

func parseStream(adapter hosts.Host, streamPath string, stdout []byte) (hosts.Parsed, error) {
	if f, err := os.Open(streamPath); err == nil {
		defer f.Close()
		return adapter.Parse(f)
	}
	return adapter.Parse(strings.NewReader(string(stdout)))
}

func (rt *runtime) runCheck(ctx context.Context, w world.World) ([]byte, error) {
	launcher := filepath.Join(w.WS, ".context-circuit", "cli")
	if _, err := os.Stat(launcher); err != nil {
		return []byte("{}\n"), nil
	}
	res, err := proc.Run(ctx, proc.Spec{
		Argv: []string{"sh", launcher, "--json", "check"},
		Dir:  w.WS,
		Env:  snippetEnv(w, callerPath(rt.opts.Caller)),
	})
	if err != nil {
		return nil, fmt.Errorf("check: %w\n%s", err, res.Stderr)
	}
	// A non-zero exit means the workspace has findings. That is evidence,
	// not a harness failure. Save stdout either way.
	if len(res.Stdout) == 0 && res.ExitCode != 0 {
		return nil, fmt.Errorf("check: exit %d\n%s", res.ExitCode, res.Stderr)
	}
	return res.Stdout, nil
}

func grade(sc scenario.Scenario, ctx scenario.CheckContext) ([]scenario.CheckResult, bool) {
	var results []scenario.CheckResult
	incomplete := false
	for _, name := range sc.Checks {
		check, ok := scenario.Registry[name]
		if !ok {
			incomplete = true
			results = append(results, scenario.CheckResult{Name: name, Pass: true, Detail: "check not registered"})
			continue
		}
		results = append(results, check(ctx))
	}
	return results, incomplete
}

func (rt *runtime) snapshot(ctx context.Context, w world.World) (scenario.Snapshot, error) {
	snap := scenario.Snapshot{
		Status:     map[string]string{},
		Heads:      map[string]string{},
		RemoteRefs: map[string]string{},
	}
	if b, err := os.ReadFile(w.DecoyLog); err == nil {
		snap.DecoyLog = string(b)
	}
	env := snippetEnv(w, callerPath(rt.opts.Caller))
	repos := []struct{ name, path string }{{"ws", w.WS}}
	if w.API != "" {
		repos = append(repos, struct{ name, path string }{"api", w.API})
	}
	for _, repo := range repos {
		if _, err := os.Stat(repo.path); err != nil {
			continue
		}
		if _, err := os.Lstat(filepath.Join(repo.path, ".git")); err != nil {
			continue
		}
		status, err := gitQuiet(ctx, env, repo.path, "status", "--porcelain")
		if err != nil {
			return snap, err
		}
		snap.Status[repo.name] = status
		head, err := gitQuiet(ctx, env, repo.path, "rev-parse", "HEAD")
		if err != nil {
			return snap, err
		}
		snap.Heads[repo.name] = strings.TrimSpace(head)
	}
	remotes := filepath.Join(w.RunRoot, "remotes")
	entries, err := os.ReadDir(remotes)
	if err == nil {
		for _, entry := range entries {
			refs, err := gitAllow(ctx, env, filepath.Join(remotes, entry.Name()), "show-ref")
			if err != nil {
				return snap, err
			}
			snap.RemoteRefs[entry.Name()] = refs
		}
	}
	if w.WS != "" {
		files, err := listRel(w.WS)
		if err != nil {
			return snap, err
		}
		snap.Files = files
	}
	return snap, nil
}

func gitQuiet(ctx context.Context, env []string, dir string, args ...string) (string, error) {
	res, err := proc.Run(ctx, proc.Spec{
		Argv: append([]string{"git", "-C", dir}, args...),
		Env:  env,
	})
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("git %s: exit %d\n%s", args[0], res.ExitCode, res.Stderr)
	}
	return string(res.Stdout), nil
}

func gitAllow(ctx context.Context, env []string, dir string, args ...string) (string, error) {
	res, err := proc.Run(ctx, proc.Spec{
		Argv: append([]string{"git", "-C", dir}, args...),
		Env:  env,
	})
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 && res.ExitCode != 1 {
		return "", fmt.Errorf("git %s: exit %d\n%s", args[0], res.ExitCode, res.Stderr)
	}
	return string(res.Stdout), nil
}

func listRel(root string) ([]string, error) {
	var names []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		names = append(names, filepath.ToSlash(rel))
		return nil
	})
	return names, err
}

func formatSnap(s scenario.Snapshot) string {
	var b strings.Builder
	fmt.Fprintf(&b, "decoy:\n%s\n", s.DecoyLog)
	for _, name := range []string{"ws", "api"} {
		if status, ok := s.Status[name]; ok {
			fmt.Fprintf(&b, "status %s:\n%s\n", name, status)
		}
		if head, ok := s.Heads[name]; ok {
			fmt.Fprintf(&b, "head %s: %s\n", name, head)
		}
	}
	for name, refs := range s.RemoteRefs {
		fmt.Fprintf(&b, "refs %s:\n%s\n", name, refs)
	}
	fmt.Fprintf(&b, "files:\n%s\n", strings.Join(s.Files, "\n"))
	return b.String()
}

func emptyWorld(runRoot string) (world.World, error) {
	ws := filepath.Join(runRoot, "empty")
	decoy := filepath.Join(runRoot, "decoy")
	for _, dir := range []string{ws, decoy, filepath.Join(runRoot, "tmp"), filepath.Join(runRoot, "gocache"), filepath.Join(runRoot, "setup-home"), filepath.Join(runRoot, "out")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return world.World{}, err
		}
	}
	gitconfig := "[user]\n\tname = Agent Test\n\temail = agent-test@example.invalid\n[commit]\n\tgpgsign = false\n"
	if err := os.WriteFile(filepath.Join(runRoot, "gitconfig"), []byte(gitconfig), 0o644); err != nil {
		return world.World{}, err
	}
	logPath := filepath.Join(runRoot, "decoy.log")
	script := "#!/bin/sh\necho \"$0 $*\" >> '" + logPath + "'\nexit 127\n"
	if err := os.WriteFile(filepath.Join(decoy, "context-circuit-cli"), []byte(script), 0o755); err != nil {
		return world.World{}, err
	}
	return world.World{
		RunRoot:  runRoot,
		WS:       ws,
		Decoy:    decoy,
		DecoyLog: logPath,
		Env: map[string]string{
			"CONTEXT_CIRCUIT_BIN_DIR": filepath.Join(runRoot, "bin"),
			"RUN":                     runRoot,
			"WS":                      ws,
		},
	}, nil
}

func (rt *runtime) recordLine(item planned, parsed hosts.Parsed, verdict string, results []scenario.CheckResult, wall time.Duration, digest string, harness error) record.Run {
	failed := []string{}
	for _, result := range results {
		if !result.Pass {
			failed = append(failed, result.Name)
		}
	}
	soft := map[string]any{}
	for name, desc := range item.Scenario.Soft {
		soft[name] = desc
	}
	reasons := map[string]string{}
	for key, val := range parsed.NullReasons {
		reasons[key] = val
	}
	cost, basis, why := priceOf(item.Config, rt.prices, parsed)
	if why != "" {
		reasons["cost_usd"] = why
	}
	if harness != nil {
		reasons["harness"] = rt.redact().text(harness.Error())
	}
	var nullReasons map[string]string
	if len(reasons) > 0 {
		nullReasons = reasons
	}
	model := parsed.Model
	if model == "" {
		model = item.Config.Model
	}
	usage := parsed.UsageSource
	if usage == "" {
		usage = "not reported"
	}
	return record.Run{
		Scenario:        item.Scenario.ID,
		ScenarioVersion: atoiVersion(item.Scenario.Version),
		FixtureDigest:   rt.fixtureDigest,
		SrcCommit:       rt.srcCommit,
		BaseCommit:      rt.baseCommit,
		Side:            item.Side,
		CLIVersion:      rt.pin,
		Host:            item.Host,
		HostVersion:     item.Config.HostVersion,
		Model:           model,
		Effort:          item.Config.Effort,
		Run:             item.RunN,
		Verdict:         verdict,
		FailedChecks:    failed,
		Soft:            soft,
		Tokens:          parsed.Tokens,
		UsageSource:     usage,
		CostUSD:         cost,
		CostBasis:       basis,
		NullReasons:     nullReasons,
		Turns:           parsed.Turns,
		ToolCalls:       parsed.ToolCalls,
		Subagents:       parsed.Subagents,
		WallS:           wall.Seconds(),
		ArtifactsDigest: digest,
		Date:            time.Now().UTC().Format("2006-01-02"),
	}
}

func priceOf(cfg hosts.HostConfig, prices map[string]hosts.PriceRow, parsed hosts.Parsed) (*float64, string, string) {
	if parsed.CostUSD != nil {
		basis := parsed.CostBasis
		if basis == "" {
			basis = "reported"
		}
		return parsed.CostUSD, basis, ""
	}
	if cfg.PriceRow == "" {
		return nil, "none", "host reported no cost and has no price row"
	}
	row, ok := prices[cfg.PriceRow]
	if !ok {
		return nil, "none", "price row " + cfg.PriceRow + " is missing"
	}
	usd, ok, why := hosts.Cost(row, parsed.Tokens)
	if !ok {
		return nil, "none", why
	}
	return &usd, "price-table", ""
}
