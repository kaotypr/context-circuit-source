package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	assets "github.com/kaotypr/context-circuit-source"
	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

const Help = `Context Circuit — shared workspace operations

Usage: context-circuit-cli [--workspace PATH] [--json] COMMAND [OPTIONS]

init                  --name NAME --purpose TEXT --member ID --member-name NAME
status                inspect workspace, members, bindings, and Git state
check                 report record, binding, dependency, and worktree issues
member add            --id ID --name NAME [--band N]
member band           --id ID --band N|0 (allocation block; 0 clears it)
member use            --id ID
member list
repo connect          --id ID --path PATH --base BRANCH
repo clone            --id ID --path NEW_PATH --base BRANCH --url URL
repo init             --id ID --path NEW_PATH --base BRANCH
repo base             --id ID --branch BRANCH
repo relate           --from ID --to ID --description TEXT
repo fetch            --id ID [--remote origin]
repo inspect          --id ID
record create         --kind intent|plan --slug SLUG --title TITLE
                      [--intent ID --repo ID ... --depends-on ID ...]
record show           --id ID
record list           [--archived]
record note           --id ID --text TEXT
record approve        --id INTENT_ID --text USER_APPROVAL
record complete       --id PLAN_ID --text RESULT
record dependencies   --id PLAN_ID [--depends-on ID ...]
record order          [--intent ID] [--mode auto|waves|linear]
                      derive dependency waves, start refs, and integration merges
context find          --query TEXT (optional index only)
worktree prepare      --repo ID [--plan ID] [--branch NAME] [--start REF]
                      [--path PATH] [--reuse] [--copy-mode auto|required|copy|off]
                      [--copy-path IGNORED_PATH ...]
worktree list          --repo ID
worktree inspect       --repo ID --path PATH
worktree move          --repo ID --path PATH --to NEW_PATH
worktree repair        --repo ID --path PATH
worktree remove        --repo ID --path PATH [--discard]
template export       --path NEW_DIRECTORY (blank workspace files)
agent settings        inspect per-host role model and effort preferences
agent configure       --host HOST --role ROLE --model MODEL|inherit --effort LEVEL|inherit
agent setup           --host codex|claude-code|cursor
agent dispatch        --host HOST --role ROLE --task TEXT --path DIRECTORY
                      [--plan ID] [--shared] [--review-requested]
                      (host must launch returned specification)
version

Approval, completion, fetch, and removal are explicit operations. The caller must
hold the user's authorization. No command executes plans, starts review, commits,
pushes, merges, deploys, or automatically removes completed worktrees.
`

type listFlag []string

func (v *listFlag) String() string         { return strings.Join(*v, ",") }
func (v *listFlag) Set(value string) error { *v = append(*v, value); return nil }

// pinnedVersionWarning reports when the running CLI differs from the version a
// workspace pins in .context-circuit/CLI_VERSION. Workspaces are installed and
// pinned independently, so a mismatch is a caller mistake rather than a
// workspace fault: warn on the error stream and continue, leaving structured
// output on the result stream untouched.
func pinnedVersionWarning(s *workspace.Store, running string) string {
	data, err := s.Read(".context-circuit/CLI_VERSION")
	if err != nil {
		return ""
	}
	pinned := strings.TrimSpace(string(data))
	if pinned == "" || pinned == running {
		return ""
	}
	return fmt.Sprintf("warning: this workspace pins CLI %s but %s is running; "+
		"run the version this workspace pins, or ask the cc-cli skill to install it\n", pinned, running)
}

func Run(ctx context.Context, args []string, out, errOut io.Writer, version string) int {
	global := flag.NewFlagSet("context-circuit-cli", flag.ContinueOnError)
	global.SetOutput(errOut)
	root := global.String("workspace", ".", "workspace root")
	asJSON := global.Bool("json", false, "structured JSON output")
	versionFlag := global.Bool("version", false, "show version")
	if err := global.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(out, Help)
			return 0
		}
		return 2
	}
	args = global.Args()
	if *versionFlag || len(args) > 0 && args[0] == "version" {
		fmt.Fprintln(out, version)
		return 0
	}
	if len(args) == 0 || args[0] == "help" {
		fmt.Fprint(out, Help)
		return 0
	}
	command := args[0]
	args = args[1:]
	if command == "member" || command == "repo" || command == "record" || command == "worktree" || command == "context" || command == "template" || command == "agent" {
		if len(args) == 0 {
			fmt.Fprint(errOut, Help)
			return 2
		}
		command += " " + args[0]
		args = args[1:]
	}
	f := flag.NewFlagSet(command, flag.ContinueOnError)
	f.SetOutput(errOut)
	values := map[string]*string{}
	add := func(names ...string) {
		for _, name := range names {
			values[name] = f.String(name, "", name)
		}
	}
	var repos, dependencies, copyPaths listFlag
	var band int
	var archived, reuse, discard, reviewRequested, shared bool
	switch command {
	case "init":
		add("name", "purpose", "member", "member-name")
	case "member add":
		add("id", "name")
		f.IntVar(&band, "band", 0, "allocation block for this member")
	case "member band":
		add("id")
		f.IntVar(&band, "band", 0, "allocation block, or 0 to clear it")
	case "member use", "repo inspect":
		add("id")
	case "repo connect", "repo init":
		add("id", "path", "base")
	case "repo clone":
		add("id", "path", "base", "url")
	case "repo base":
		add("id", "branch")
	case "repo fetch":
		add("id", "remote")
	case "repo relate":
		add("from", "to", "description")
	case "record create":
		add("kind", "slug", "title", "intent")
		f.Var(&repos, "repo", "repository (repeatable)")
		f.Var(&dependencies, "depends-on", "plan dependency (repeatable)")
	case "record show":
		add("id")
	case "record list":
		f.BoolVar(&archived, "archived", false, "include archived records")
	case "record note", "record approve", "record complete":
		add("id", "text")
	case "record dependencies":
		add("id")
		f.Var(&dependencies, "depends-on", "plan dependency (repeatable)")
	case "record order":
		add("intent")
		values["mode"] = f.String("mode", "auto", "auto, waves, or linear")
	case "context find":
		add("query")
	case "worktree prepare":
		add("repo", "plan", "branch", "start", "path")
		values["copy-mode"] = f.String("copy-mode", "auto", "runtime file reuse: auto, required, copy, off")
		f.Var(&copyPaths, "copy-path", "additional ignored runtime path (repeatable)")
		f.BoolVar(&reuse, "reuse", false, "explicitly reuse existing work")
	case "worktree list":
		add("repo")
	case "worktree inspect", "worktree repair":
		add("repo", "path")
	case "worktree move":
		add("repo", "path", "to")
	case "worktree remove":
		add("repo", "path")
		f.BoolVar(&discard, "discard", false, "explicitly discard worktree files")
	case "template export":
		add("path")
	case "agent setup":
		add("host")
	case "agent configure":
		add("host", "role", "model", "effort")
	case "agent dispatch":
		add("host", "role", "task", "path", "plan")
		f.BoolVar(&shared, "shared", false, "several workers share this worktree")
		f.BoolVar(&reviewRequested, "review-requested", false, "user explicitly requested independent review")
	case "agent settings":
	case "status", "check", "member list":
	default:
		fmt.Fprintf(errOut, "unknown command: %s\n", command)
		return 2
	}
	if err := f.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if f.NArg() != 0 {
		fmt.Fprintln(errOut, "unexpected positional arguments")
		return 2
	}
	get := func(key string) string {
		if p := values[key]; p != nil {
			return *p
		}
		return ""
	}
	optional := map[string]bool{"remote": true}
	optional["intent"] = command == "record create" || command == "record order"
	if command == "worktree prepare" {
		for _, k := range []string{"plan", "branch", "start", "path"} {
			optional[k] = true
		}
	}
	if command == "agent dispatch" {
		optional["plan"] = true
	}
	for key, value := range values {
		if *value == "" && !optional[key] {
			fmt.Fprintf(errOut, "missing --%s\n", key)
			return 2
		}
	}
	if command == "init" || command == "template export" {
		if command == "template export" {
			*root = get("path")
		}
		if err := os.MkdirAll(*root, 0755); err != nil {
			fmt.Fprintln(errOut, err)
			return 1
		}
	}
	s, err := workspace.Open(*root)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	if command != "template export" {
		if warning := pinnedVersionWarning(s, version); warning != "" {
			fmt.Fprint(errOut, warning)
		}
	}
	readOnly := command == "agent settings" || command == "agent dispatch" || command == "template export" || command == "status" || command == "check" || command == "member list" || command == "record show" || command == "record list" || command == "repo inspect" || command == "context find" || command == "worktree list" || command == "worktree inspect"
	checkFailed := false
	action := func() (any, error) {
		var err error
		switch command {
		case "init", "template export":
			files, e := assets.Files()
			if e != nil {
				return nil, e
			}
			if command == "init" {
				err = s.Init(files, get("name"), get("purpose"), get("member"), get("member-name"))
			} else {
				err = s.Export(files)
			}
		case "status":
			return s.Status(ctx)
		case "agent settings":
			return s.AgentSettings()
		case "agent configure":
			err = s.ConfigureAgent(get("host"), get("role"), get("model"), get("effort"))
		case "agent setup":
			return s.SetupAgents(get("host"))
		case "agent dispatch":
			return s.DispatchAgent(get("host"), get("role"), get("task"), get("path"), get("plan"), shared, reviewRequested)
		case "check":
			issues, e := s.Check(ctx)
			checkFailed = len(issues) > 0
			return map[string]any{"issues": issues, "ok": !checkFailed}, e
		case "member list":
			return s.Members()
		case "member add":
			err = s.AddMember(get("id"), get("name"), band)
		case "member band":
			err = s.SetMemberBand(get("id"), band)
		case "member use":
			err = s.UseMember(get("id"))
		case "repo connect":
			err = s.Connect(ctx, get("id"), get("path"), get("base"))
		case "repo clone", "repo init":
			err = s.CreateRepository(ctx, get("id"), get("path"), get("base"), get("url"))
		case "repo base":
			err = s.SetBase(ctx, get("id"), get("branch"))
		case "repo relate":
			err = s.Relate(get("from"), get("to"), get("description"))
		case "repo inspect", "repo fetch":
			_, path, e := s.Repository(ctx, get("id"))
			if e != nil {
				return nil, e
			}
			if command == "repo inspect" {
				return workspace.Inspect(ctx, path)
			}
			remote := get("remote")
			if remote == "" {
				remote = "origin"
			}
			if e := workspace.Name(remote); e != nil {
				return nil, e
			}
			_, err = workspace.Git(ctx, path, "fetch", "--", remote)
		case "record create":
			return s.CreateRecord(get("kind"), get("slug"), get("title"), get("intent"), repos, dependencies)
		case "record show":
			return s.FindRecord(get("id"))
		case "record list":
			return s.ListRecords(archived)
		case "record note", "record approve":
			err = s.Note(get("id"), get("text"), strings.TrimPrefix(command, "record "))
		case "record complete":
			if err := s.Note(get("id"), get("text"), "complete"); err != nil {
				return nil, err
			}
			entries, e := s.KnowledgeCandidates(get("id"))
			return map[string]any{"completed": get("id"), "knowledge_candidates": entries}, e
		case "record dependencies":
			err = s.SetDependencies(get("id"), dependencies)
		case "record order":
			return s.Order(get("intent"), get("mode"))
		case "context find":
			return s.FindContext(get("query"))
		case "worktree prepare":
			return s.Prepare(ctx, get("repo"), get("plan"), get("branch"), get("start"), get("path"), reuse, workspace.ReuseOptions{Mode: get("copy-mode"), Paths: copyPaths})
		case "worktree list":
			return s.Worktrees(ctx, get("repo"))
		case "worktree inspect":
			return s.InspectWorktree(ctx, get("repo"), get("path"))
		case "worktree move":
			err = s.MoveWorktree(ctx, get("repo"), get("path"), get("to"))
		case "worktree remove":
			err = s.RemoveWorktree(ctx, get("repo"), get("path"), discard)
		case "worktree repair":
			err = s.RepairWorktree(ctx, get("repo"), get("path"))
		}
		return map[string]any{"operation": command, "workspace": s.Root, "ok": err == nil}, err
	}
	var result any
	if readOnly {
		result, err = action()
	} else {
		result, err = s.WithLock(ctx, action)
	}
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	if command == "record show" && !*asJSON {
		_, err = fmt.Fprint(out, result.(workspace.Record).Content)
	} else if *asJSON {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		err = enc.Encode(result)
	} else {
		var data []byte
		// Output uses the public JSON names, including record paths; the YAML
		// tags on records describe frontmatter and deliberately omit content.
		data, err = json.Marshal(result)
		if err == nil {
			data, err = yaml.JSONToYAML(data)
		}
		if err == nil {
			_, err = out.Write(data)
		}
	}
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	if checkFailed {
		return 1
	}
	return 0
}
