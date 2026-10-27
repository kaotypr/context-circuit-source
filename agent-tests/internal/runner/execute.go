package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	var graded *scenario.CheckContext
	var loginPaths []string
	loginNoted := false
	noteLoginRemoval := func(derr error) {
		if derr == nil || loginNoted {
			return
		}
		loginNoted = true
		drops = append(drops, errLoginRemove.Error())
	}
	rt.loginSecrets = nil
	runDir := filepath.Join(rt.opts.OutDir, rt.id, "runs", item.Scenario.ID, item.Host, item.Side, fmt.Sprintf("r%d", item.RunN))
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic: %v", rec)
		}
		if derr := removeLoginFiles(loginPaths); derr != nil {
			noteLoginRemoval(derr)
			if err == nil {
				err = derr
			}
		}
		if err != nil {
			verdict = scenario.Verdict(nil, err, false)
		}
		if len(drops) > 0 {
			_ = writeText(filepath.Join(runDir, "drop.log"), strings.Join(drops, "\n")+"\n", rt.redact())
		}
		digest, _ := record.Digest(runDir)
		line := rt.recordLine(item, parsed, verdict, results, wall, digest, err, graded)
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
	sub := ""
	if containsHost(rt.opts.SubscriptionAuth, item.Host) {
		sub = item.Host
	}
	if err = writeText(filepath.Join(runDir, "reproduce.txt"), reproduceLine(item, rt.opts.Ref, rt.opts.BaseRef, time.Now().UTC().Format("2006-01-02"), sub), rt.redact()); err != nil {
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
		_ = os.Remove(runRoot) // still empty: nothing has been written yet
		runRoot = ""
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
	applyWorkspace(&w, item.Scenario)
	checkBefore, err := rt.captureCheck(ctx, w)
	if err != nil {
		if !errors.Is(err, errNoLauncher) {
			return err
		}
		checkBefore = nil
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
		Prompt:           prompt,
		Model:            item.Config.Model,
		Effort:           item.Config.Effort,
		Workspace:        w.WS,
		RunRoot:          w.RunRoot,
		SessionID:        newSessionID(),
		Cfg:              item.Config,
		Caller:           rt.opts.Caller,
		DecoyDir:         w.Decoy,
		SubscriptionAuth: containsHost(rt.opts.SubscriptionAuth, item.Host),
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
	loginPaths = loginCopyPaths(launch.CopyIn)
	defer func() {
		if derr := removeLoginFiles(loginPaths); derr != nil {
			noteLoginRemoval(derr)
			if err == nil {
				err = derr
			}
		}
	}()
	secrets, copyErr := hosts.ApplyCopies(launch.CopyIn)
	if copyErr != nil {
		return copyErr
	}
	launch.Secrets = secrets
	rt.loginSecrets = append(rt.loginSecrets, secrets...)
	if len(launch.CopyIn) > 0 && afterLoginCopy != nil {
		afterLoginCopy()
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
	res, runErr := proc.Run(ctx, proc.Spec{
		Argv:       launch.Argv,
		Dir:        launch.Dir,
		Env:        env,
		Stdin:      launch.Stdin,
		Timeout:    timeout,
		StreamPath: streamPath,
	})
	rt.refreshLoginSecrets(loginPaths)
	if derr := removeLoginFiles(loginPaths); derr != nil {
		noteLoginRemoval(derr)
		if runErr == nil && !res.TimedOut && res.ExitCode == 0 {
			rt.scrubFile(streamPath)
			return derr
		}
	}
	rt.scrubFile(streamPath)
	if runErr != nil {
		return fmt.Errorf("host: %w\n%s", runErr, res.Stderr)
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
	streamText := string(res.Stdout)
	if raw, readErr := os.ReadFile(streamPath); readErr == nil {
		streamText = string(raw)
	}
	after, err := rt.snapshot(ctx, w)
	if err != nil {
		return err
	}
	checkOut, checkErr := rt.runCheck(ctx, w)
	if checkErr != nil {
		if !errors.Is(checkErr, errNoLauncher) {
			return checkErr
		}
		checkOut = nil
	}
	graded = &scenario.CheckContext{
		World:       w,
		Scenario:    item.Scenario,
		Parsed:      parsed,
		Before:      before,
		After:       after,
		CheckBefore: checkBefore,
		CheckAfter:  checkOut,
		Stream:      streamText,
		Host:        item.Host,
		Exec: func(dir string, argv ...string) (string, error) {
			return rt.execIsolated(ctx, w, dir, argv...)
		},
	}
	results, incomplete = grade(item.Scenario, *graded)
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
	if err = publishCheckFiles(runDir, checkBefore, checkOut, red, &drops); err != nil {
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
	if err = rt.writeGitState(ctx, w, before, item.Scenario, filepath.Join(runDir, "artifacts", "git")); err != nil {
		return err
	}
	return nil
}

func (rt *runtime) runSnippets(ctx context.Context, w world.World, sc scenario.Scenario) error {
	// FIXTURES is for setup snippets only. The host environment is built
	// later and does not receive it.
	env := snippetEnv(w, callerPath(rt.opts.Caller))
	if rt.opts.FixturesDir != "" {
		env = append(env, "FIXTURES="+rt.opts.FixturesDir)
	}
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

// errNoLauncher means the workspace has no launcher to run check with.
// Callers leave the check document empty so check-clean fails, and they
// do not substitute an empty JSON object.
var errNoLauncher = errors.New("workspace launcher is missing")

func (rt *runtime) runCheck(ctx context.Context, w world.World) ([]byte, error) {
	launcher := filepath.Join(w.WS, ".context-circuit", "cli")
	info, err := os.Stat(launcher)
	if err != nil || info.IsDir() {
		return nil, fmt.Errorf("check: %w", errNoLauncher)
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
			results = append(results, scenario.CheckResult{Name: name, Pass: false, Detail: "check not registered"})
			continue
		}
		result := check(ctx)
		// Inconclusive is not a failure and is not a pass. FailedChecks
		// lists only results with Pass false.
		if result.Inconclusive {
			incomplete = true
		}
		results = append(results, result)
	}
	return results, incomplete
}

func (rt *runtime) snapshot(ctx context.Context, w world.World) (scenario.Snapshot, error) {
	snap := scenario.Snapshot{
		Status:     map[string]string{},
		Heads:      map[string]string{},
		RemoteRefs: map[string]string{},
		Refs:       map[string]string{},
		Worktrees:  map[string]string{},
	}
	b, err := os.ReadFile(w.DecoyLog)
	switch {
	case err == nil:
		snap.DecoyLog = string(b)
	case os.IsNotExist(err):
	default:
		snap.DecoyLog = scenario.DecoyLogUnreadable
	}
	env := snippetEnv(w, callerPath(rt.opts.Caller))
	for _, repo := range snapshotRepos(w) {
		if !repoHasGit(repo.path) {
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
		refs, err := gitQuiet(ctx, env, repo.path, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads")
		if err != nil {
			return snap, err
		}
		snap.Refs[repo.name] = sortedRefLines(refs)
		trees, err := gitQuiet(ctx, env, repo.path, "worktree", "list", "--porcelain")
		if err != nil {
			return snap, err
		}
		snap.Worktrees[repo.name] = trees
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
	digest, err := runDigest(w.RunRoot)
	if err != nil {
		return snap, err
	}
	snap.RunDigest = digest
	return snap, nil
}

type repoPath struct {
	name string
	path string
}

// snapshotRepos lists the workspace and every directory under
// repositories/. World.API is included when it is not already in that list.
func snapshotRepos(w world.World) []repoPath {
	listed := scenario.WorldRepos(w)
	repos := make([]repoPath, len(listed))
	for i, repo := range listed {
		repos[i] = repoPath{name: repo.Name, path: repo.Path}
	}
	return repos
}

func repoHasGit(path string) bool {
	if path == "" {
		return false
	}
	if _, err := os.Stat(path); err != nil {
		return false
	}
	_, err := os.Lstat(filepath.Join(path, ".git"))
	return err == nil
}

func sortedRefLines(text string) string {
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	sort.Strings(lines)
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

// applyWorkspace points the host at <run root>/<workspace> after setup
// snippets have run in the default workspace.
func applyWorkspace(w *world.World, sc scenario.Scenario) {
	if sc.Workspace == "" || w == nil {
		return
	}
	w.WS = filepath.Join(w.RunRoot, sc.Workspace)
	w.API = filepath.Join(w.WS, "repositories", "api")
	if w.Env == nil {
		w.Env = map[string]string{}
	}
	w.Env["WS"] = w.WS
}

func (rt *runtime) captureCheck(ctx context.Context, w world.World) ([]byte, error) {
	out, err := rt.runCheck(ctx, w)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(w.RunRoot, "out")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "check-before.json"), out, 0o644); err != nil {
		return nil, err
	}
	return out, nil
}

func (rt *runtime) execIsolated(ctx context.Context, w world.World, dir string, argv ...string) (string, error) {
	res, err := proc.Run(ctx, proc.Spec{
		Argv: argv,
		Dir:  dir,
		Env:  snippetEnv(w, callerPath(rt.opts.Caller)),
	})
	if err != nil {
		return string(res.Stdout), fmt.Errorf("%w\n%s", err, res.Stderr)
	}
	if res.ExitCode != 0 {
		return string(res.Stdout), fmt.Errorf("exit %d\n%s", res.ExitCode, res.Stderr)
	}
	return string(res.Stdout), nil
}

// reproduceLine is the command that repeats this run. It names no local path.
// A base-side line keeps the configured branch ref on --ref and adds --base.
func reproduceLine(item planned, branchRef, baseRef, date, subscription string) string {
	ref := branchRef
	if ref == "" {
		ref = item.Ref
	}
	line := fmt.Sprintf("agent-tests run --scenarios %s --hosts %s --runs 1 --ref %s", item.Scenario.ID, item.Host, ref)
	if item.Side == "base" {
		base := baseRef
		if base == "" {
			base = item.Ref
		}
		line += " --base " + base
	}
	if subscription != "" {
		line += " --subscription-auth " + subscription
	}
	return "# " + date + "\n" + line + "\n"
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
	names := map[string]struct{}{}
	for _, section := range []map[string]string{s.Status, s.Heads, s.Refs, s.Worktrees} {
		for name := range section {
			names[name] = struct{}{}
		}
	}
	keys := make([]string, 0, len(names))
	for name := range names {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	for _, name := range keys {
		if status, ok := s.Status[name]; ok {
			fmt.Fprintf(&b, "status %s:\n%s\n", name, status)
		}
		if head, ok := s.Heads[name]; ok {
			fmt.Fprintf(&b, "head %s: %s\n", name, head)
		}
		if refs, ok := s.Refs[name]; ok {
			fmt.Fprintf(&b, "local-refs %s:\n%s\n", name, refs)
		}
		if trees, ok := s.Worktrees[name]; ok {
			fmt.Fprintf(&b, "worktrees %s:\n%s\n", name, trees)
		}
	}
	remoteNames := make([]string, 0, len(s.RemoteRefs))
	for name := range s.RemoteRefs {
		remoteNames = append(remoteNames, name)
	}
	sort.Strings(remoteNames)
	for _, name := range remoteNames {
		fmt.Fprintf(&b, "remote-refs %s:\n%s\n", name, s.RemoteRefs[name])
	}
	fmt.Fprintf(&b, "files:\n%s\n", strings.Join(s.Files, "\n"))
	fmt.Fprintf(&b, "digest-entries: %d\n", len(s.RunDigest))
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

func (rt *runtime) recordLine(item planned, parsed hosts.Parsed, verdict string, results []scenario.CheckResult, wall time.Duration, digest string, harness error, graded *scenario.CheckContext) record.Run {
	// FailedChecks lists checks with Pass false. An inconclusive result
	// has Pass set, so it is left out.
	failed := []string{}
	for _, result := range results {
		if !result.Pass {
			failed = append(failed, result.Name)
		}
	}
	soft := map[string]any{}
	for name, desc := range item.Scenario.Soft {
		if graded != nil {
			if fn, ok := scenario.SoftRegistry[name]; ok {
				soft[name] = fn(*graded)
				continue
			}
		}
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
	hostVer, hostVerSource := item.Config.HostVersion, "config"
	if hv, ok := rt.hostVersions[item.Host]; ok {
		hostVer = hv.value
		if hv.measured {
			hostVerSource = "command"
		} else {
			reasons["host_version"] = "the host command did not report --version; the configured version is recorded"
		}
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
		Scenario:          item.Scenario.ID,
		ScenarioVersion:   atoiVersion(item.Scenario.Version),
		FixtureDigest:     rt.fixtureDigest,
		SrcCommit:         rt.srcCommit,
		BaseCommit:        rt.baseCommit,
		Side:              item.Side,
		CLIVersion:        rt.pin,
		Host:              item.Host,
		HostVersion:       hostVer,
		HostVersionSource: hostVerSource,
		Model:             model,
		Effort:            item.Config.Effort,
		Run:               item.RunN,
		Verdict:           verdict,
		FailedChecks:      failed,
		Soft:              soft,
		Tokens:            parsed.Tokens,
		UsageSource:       usage,
		CostUSD:           cost,
		CostBasis:         basis,
		NullReasons:       nullReasons,
		Turns:             parsed.Turns,
		ToolCalls:         parsed.ToolCalls,
		Subagents:         parsed.Subagents,
		WallS:             wall.Seconds(),
		ArtifactsDigest:   digest,
		Date:              time.Now().UTC().Format("2006-01-02"),
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
