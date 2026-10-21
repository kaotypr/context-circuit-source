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
The global options --workspace PATH and --json may also follow the command.

init                  --name NAME --purpose TEXT --member ID --member-name NAME
                      [--member-band N] [--member-role ROLE ...] (the first member,
                      selected here; band 1 unless another is given)
status                inspect workspace, members, bindings, and Git state
check                 report record, binding, dependency, and worktree issues
member add            --id ID --name NAME --band N [--language NAME] [--tone TEXT]
                      [--role ROLE ...]
member band           --id ID --band N (allocation block, 0-9999; 0 is a band)
member language       --id ID --language NAME (the language this member's
                      intents and plans are written in; knowledge stays English)
member tone           --id ID --tone TEXT | --clear (the register that language
                      is written in, such as 'semi-formal; keep technical terms
                      in English')
member roles          --id ID --role ROLE ... | --clear (what this member
                      contributes, such as maintainer; informational, decides nothing)
member use            --id ID
member list
repo connect          --id ID --path PATH --base BRANCH [--url URL]
                      [--default-branch BRANCH] (--base is this machine's)
repo clone            --id ID --path NEW_PATH --base BRANCH --url URL
                      [--default-branch BRANCH]
repo init             --id ID --path NEW_PATH --base BRANCH [--default-branch BRANCH]
repo base             --id ID --branch BRANCH (this machine's base branch)
repo remote           --id ID [--url URL] [--default-branch BRANCH] (shared)
repo relate           --from ID --to ID --description TEXT
workspace connect     [--path PATH] --base BRANCH [--url URL]
                      [--default-branch BRANCH] (the workspace's own repository;
                      --path defaults to the workspace root, --base is this machine's)
workspace base        --branch BRANCH (this machine's workspace base branch)
workspace remote      [--url URL] [--default-branch BRANCH] (shared)
knowledge connect     --id ID --path PATH [--url URL] [--default-branch BRANCH]
                      [--index PATH] (a repository this workspace reads and never writes)
knowledge clone       --id ID [--path NEW_PATH] [--url URL] [--default-branch BRANCH]
                      [--index PATH] (--path defaults to repositories/<id>)
knowledge sync        --id ID (fetch, then fast-forward only when clean and
                      on the shared branch; never merges, resets, or discards)
knowledge remote      --id ID [--url URL] [--default-branch BRANCH] [--index PATH] (shared)
knowledge list        every borrowed repository and this machine's state of it
repo fetch            --id ID [--remote origin]
repo inspect          --id ID
record create         --kind intent|plan --slug SLUG --title TITLE
                      [--intent ID --repo ID ... --depends-on ID ...]
record show           --id ID
record list           [--archived]
record approve        --id INTENT_ID --text USER_APPROVAL
record complete       --id PLAN_ID --text RESULT
record dependencies   --id PLAN_ID [--depends-on ID ...]
record delete         --member ID | --all | --id ID [--id ID ...] [--confirm]
                      (a member's intents and plans, every record while the roster
                      lists one member, or the named records, releasing their IDs;
                      without --confirm it previews and deletes nothing)
record order          [--intent ID | --plan ID ...] [--mode auto|waves|linear]
                      derive dependency waves, start refs, and integration merges
context find          --query TEXT | --repo ID (optional index only)
worktree prepare      --repo ID [--plan ID] [--branch NAME] [--start REF]
                      [--path PATH] [--reuse] [--copy-mode auto|required|copy|off]
                      [--copy-path IGNORED_PATH ...]
worktree list          --repo ID
worktree inspect       --repo ID --path PATH
worktree move          --repo ID --path PATH --to NEW_PATH
worktree repair        --repo ID --path PATH
worktree remove        --repo ID --path PATH [--discard]
worktree prune         --repo ID (drops associations of worktrees Git no longer lists)
template export       --path NEW_DIRECTORY (blank workspace files)
agent settings        inspect effective per-host role model and effort preferences
                      (shared, with any role-tiering.local.yaml override applied)
agent configure       --host HOST --role ROLE --model MODEL|inherit --effort LEVEL|inherit
                      [--local] (write this machine's override, not the shared setting)
agent setup           [--host codex|claude-code|cursor] (default: every host)
agent dispatch        --host HOST --role ROLE --path DIRECTORY [--task TEXT]
                      [--intent ID] [--plan ID] [--repo ID] [--shared] [--review-requested]
                      (--intent plans an approved intent; planner --task plans a specified standalone outcome; --plan implements one.
                      --task is the assignment where no record carries it.
                      Host must launch the returned prompt unmodified.)
version

Every command that changes something is refused until the roster lists this
machine's selected member and every member holds a band; init, member add,
member use, and member band are the steps that settle it.

Approval, completion, fetch, and removal are explicit operations. The caller must
hold the user's authorization. No command executes plans, starts review, commits,
pushes, merges, deploys, or automatically removes completed worktrees.
`

type listFlag []string

func (v *listFlag) String() string         { return strings.Join(*v, ",") }
func (v *listFlag) Set(value string) error { *v = append(*v, value); return nil }

// provided reports whether the caller actually passed a flag, which is the only
// way to tell an omitted number from one deliberately set to its zero value.
func provided(f *flag.FlagSet, name string) bool {
	found := false
	f.Visit(func(set *flag.Flag) {
		if set.Name == name {
			found = true
		}
	})
	return found
}

// writeYAML renders a result for a human through the public JSON names, so a
// field is called the same thing whichever stream a caller reads.
func writeYAML(out io.Writer, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if data, err = yaml.JSONToYAML(data); err != nil {
		return err
	}
	_, err = out.Write(data)
	return err
}

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

// commitRecommendation names the records and ID files a gate leaves
// uncommitted. Until they are committed and pushed no other clone sees them,
// and an ID file left behind is how two people reserve the same number.
type commitRecommendation struct {
	Files []workspace.UncommittedFile `json:"files"`
	Note  string                      `json:"note"`
}

// recommendCommit adds what approval and completion leave uncommitted to their
// result. It runs after the command has succeeded and only ever adds a key: a
// recommendation when there is something to commit, or the reason Git could
// not say. It never stages or commits, because committing workspace records is
// the person's decision.
func recommendCommit(ctx context.Context, s *workspace.Store, result map[string]any, ids ...string) {
	files, unchecked := s.UncommittedRecords(ctx, ids...)
	if unchecked != "" {
		result["commit_unchecked"] = unchecked + "; whether these records and the ID files are committed was not checked"
		return
	}
	if len(files) > 0 {
		result["commit_recommended"] = commitRecommendation{
			Files: files,
			Note:  "these records and ID files are not committed, so other clones cannot see them yet: recommend committing and pushing them to the person; committing still needs the person's authorization, and nothing was staged or committed",
		}
	}
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
		if len(args) > 1 {
			topic := strings.Join(args[1:], " ")
			text, ok := CommandHelp(topic)
			if !ok {
				// Printing the whole manual for a question about one command
				// answers something the caller did not ask and reads as success.
				fmt.Fprintf(errOut, "no help topic: %s\ntry: context-circuit-cli help\n", topic)
				return 2
			}
			fmt.Fprint(out, text)
			return 0
		}
		fmt.Fprint(out, Help)
		return 0
	}
	command := args[0]
	args = args[1:]
	if groups[command] {
		if len(args) == 0 {
			fmt.Fprint(errOut, Help)
			return 2
		}
		// A group with an option where its subcommand belongs is a caller
		// looking for help, not an unknown command. Name the real fault.
		if strings.HasPrefix(args[0], "-") {
			fmt.Fprintf(errOut, "unknown option for %s: %s\ntry: context-circuit-cli help %s\n", command, args[0], command)
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
	var repos, dependencies, copyPaths, plans, roles, recordIDs listFlag
	var band int
	var archived, reuse, discard, reviewRequested, shared, localTiering, clearTone, clearRoles, all, confirm bool
	switch command {
	case "init":
		add("name", "purpose", "member", "member-name")
		f.IntVar(&band, "member-band", 1, "allocation block for the first member")
		f.Var(&roles, "member-role", "informational role of the first member (repeatable)")
	case "member add":
		add("id", "name", "language", "tone")
		f.IntVar(&band, "band", 0, "allocation block for this member")
		f.Var(&roles, "role", "informational role of this member (repeatable)")
	case "member roles":
		add("id")
		f.Var(&roles, "role", "informational role, replacing those recorded (repeatable)")
		f.BoolVar(&clearRoles, "clear", false, "remove this member's recorded roles")
	case "member language":
		add("id", "language")
	case "member tone":
		add("id", "tone")
		f.BoolVar(&clearTone, "clear", false, "remove this member's recorded register")
	case "member band":
		add("id")
		f.IntVar(&band, "band", 0, "allocation block; 0 is a band like any other")
	case "member use", "repo inspect":
		add("id")
	case "repo connect", "repo clone":
		add("id", "path", "base", "url", "default-branch")
	case "repo init":
		add("id", "path", "base", "default-branch")
	case "repo remote":
		add("id", "url", "default-branch")
	case "repo base":
		add("id", "branch")
	case "repo fetch":
		add("id", "remote")
	case "repo relate":
		add("from", "to", "description")
	case "workspace connect":
		add("path", "base", "url", "default-branch")
	case "workspace base":
		add("branch")
	case "workspace remote":
		add("url", "default-branch")
	case "knowledge connect", "knowledge clone":
		add("id", "path", "url", "default-branch", "index")
	case "knowledge remote":
		add("id", "url", "default-branch", "index")
	case "knowledge sync":
		add("id")
	case "knowledge list":
	case "record create":
		add("kind", "slug", "title", "intent")
		f.Var(&repos, "repo", "repository (repeatable)")
		f.Var(&dependencies, "depends-on", "plan dependency (repeatable)")
	case "record show":
		add("id")
	case "record list":
		f.BoolVar(&archived, "archived", false, "include archived records")
	case "record approve", "record complete":
		add("id", "text")
	case "record delete":
		add("member")
		f.BoolVar(&all, "all", false, "every record, while the roster lists one member")
		f.Var(&recordIDs, "id", "intent or plan ID to delete (repeatable)")
		f.BoolVar(&confirm, "confirm", false, "delete, after the person agreed to the previewed list")
	case "record dependencies":
		add("id")
		f.Var(&dependencies, "depends-on", "plan dependency (repeatable)")
	case "record order":
		add("intent")
		f.Var(&plans, "plan", "selected plan (repeatable)")
		values["mode"] = f.String("mode", "auto", "auto, waves, or linear")
	case "context find":
		add("query", "repo")
	case "worktree prepare":
		add("repo", "plan", "branch", "start", "path")
		values["copy-mode"] = f.String("copy-mode", "auto", "runtime file reuse: auto, required, copy, off")
		f.Var(&copyPaths, "copy-path", "additional ignored runtime path (repeatable)")
		f.BoolVar(&reuse, "reuse", false, "explicitly reuse existing work")
	case "worktree list", "worktree prune":
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
		f.BoolVar(&localTiering, "local", false, "write this machine's override instead of the shared setting")
	case "agent dispatch":
		add("host", "role", "task", "path", "plan", "intent", "repo")
		f.BoolVar(&shared, "shared", false, "several workers share this worktree")
		f.BoolVar(&reviewRequested, "review-requested", false, "user explicitly requested independent review")
	case "agent settings":
	case "status", "check", "member list":
	default:
		fmt.Fprintf(errOut, "unknown command: %s\n", command)
		// An empty argument is a command name with no fields at all, so there
		// is no group to offer a topic for. Reporting it is still right; going
		// looking for its first word is what crashes.
		if topic := strings.Fields(command); len(topic) > 0 && groups[topic[0]] {
			fmt.Fprintf(errOut, "try: context-circuit-cli help %s\n", topic[0])
		}
		return 2
	}
	// The global options are accepted after the command as well as before it,
	// because a caller adding --json or --workspace to a command line they
	// already typed puts it at the end. A command defining a flag of the same
	// name keeps its own meaning; none does today.
	if f.Lookup("workspace") == nil {
		f.StringVar(root, "workspace", *root, "workspace root")
	}
	if f.Lookup("json") == nil {
		f.BoolVar(asJSON, "json", *asJSON, "structured JSON output")
	}
	// The flag package still reports the faulty flag itself; the usage that
	// follows is this command's, printed below.
	f.Usage = func() {}
	if err := f.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			text, _ := CommandHelp(command)
			fmt.Fprint(out, text)
			return 0
		}
		fmt.Fprint(errOut, commandUsage(command))
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
	optional := map[string]bool{"remote": true, "default-branch": true}
	// A member writing in the workspace's usual language records nothing, and a
	// team happy with the general composition guidance records no register; both
	// fields exist for the member who is neither.
	optional["language"] = command == "member add"
	optional["tone"] = command == "member add" || command == "member tone"
	// The workspace's own checkout is the workspace root unless a workspace
	// kept inside a larger repository names the root above it.
	if command == "workspace connect" {
		optional["path"] = true
	}
	// Either half of the shared record is a whole call; an empty one is refused
	// below, where the alternative can be named.
	if command == "workspace remote" {
		optional["url"], optional["default-branch"] = true, true
	}
	// A borrowed repository names its own index, and a layout this product
	// recognizes is detected rather than demanded. Obtaining one lands under
	// repositories/<id> unless the caller keeps it somewhere else.
	if strings.HasPrefix(command, "knowledge ") {
		optional["index"] = true
		optional["path"] = command != "knowledge connect"
	}
	// A URL is one shared detail the checkout usually already knows. Even a
	// clone may omit it, because an ID this workspace already describes
	// carries the source in the shared record — which is what joining one
	// looks like.
	optional["url"] = true
	// Deleting every record names no member; the store refuses a half-given
	// pair where it can say why.
	optional["member"] = command == "record delete"
	optional["intent"] = command == "record create" || command == "record order" || command == "agent dispatch"
	if command == "worktree prepare" {
		for _, k := range []string{"plan", "branch", "start", "path"} {
			optional[k] = true
		}
	}
	// A search takes a query and a reconciliation takes a repository. Either one
	// is a whole call, so both are optional here and an empty call is refused
	// below, where the alternative can be named.
	if command == "context find" {
		optional["query"], optional["repo"] = true, true
	}
	// A planner takes --intent and a worker takes --plan; neither is universal,
	// so the command accepts both as optional and refuses the wrong pairing by
	// role, where the reason can be stated.
	if command == "agent dispatch" {
		optional["plan"], optional["intent"], optional["repo"] = true, true, true
		// A planner and a whole-plan worker are assigned by the record quoted
		// in the brief; only a role without one, or a worker taking a slice of
		// a shared plan, has something left to say.
		optional["task"] = true
	}
	// A workspace is opened in more than one host, so setup covers them all
	// unless the caller narrows it to one.
	if command == "agent setup" {
		optional["host"] = true
	}
	for key, value := range values {
		if *value == "" && !optional[key] {
			fmt.Fprintf(errOut, "missing --%s\n", key)
			return 2
		}
	}
	// The required check above reads string flags, and a band is a number whose
	// zero value is a real band. Omitting the flag would otherwise move the
	// member to band 0 silently, which a caller who forgot it cannot have meant.
	if command == "member band" && !provided(f, "band") {
		fmt.Fprintln(errOut, "missing --band")
		return 2
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
	readOnly := command == "agent settings" || command == "agent dispatch" || command == "template export" || command == "status" || command == "check" || command == "member list" || command == "record show" || command == "record list" || command == "repo inspect" || command == "context find" || command == "knowledge list" || command == "worktree list" || command == "worktree inspect"
	// These settle setup, so they must run before it is complete; every other
	// change waits for it. A command missing from readOnly is refused here by
	// default, which is the safe direction for one added later.
	setupStep := command == "init" || command == "member add" || command == "member use" || command == "member band"
	// The ID split is a change of its own, so only a changing command makes it,
	// and only once setup has named the member who receives authorless
	// reservations. Init writes the new layout itself, and a deletion preview
	// promises to write nothing.
	splits := !readOnly && command != "init" && !(command == "record delete" && !confirm)
	var split *workspace.IDSplit
	splitWarning := ""
	// trySplit makes the split once setup is complete and it is still pending.
	// It is tried before the command, so a change lands in the per-member layout,
	// and again after it, so the setup step that completed setup makes the split
	// itself rather than leaving it to whichever command comes next.
	trySplit := func() error {
		if !splits || split != nil {
			return nil
		}
		state, err := s.Setup()
		if err != nil {
			return nil
		}
		pending, err := s.IDSplitPending()
		if err != nil || !pending {
			return err
		}
		split, err = s.SplitIDLedger(state.Active)
		return err
	}
	checkFailed := false
	action := func() (any, error) {
		if !readOnly && !setupStep {
			if err := s.RequireSetup(); err != nil {
				return nil, fmt.Errorf("setup is incomplete, so nothing was changed: %w", err)
			}
		}
		// A call its own command would refuse is refused before the split, so a
		// mistyped command never moves a schema 2 workspace on its way to failing.
		switch command {
		case "member add":
			// Nothing is written without a band the person chose, but the
			// refusal offers one so choosing costs nothing.
			if !provided(f, "band") {
				suggestions, e := s.SuggestBands()
				if e != nil {
					return nil, fmt.Errorf("pass --band N: %w", e)
				}
				return nil, fmt.Errorf("pass --band N; the next free band is %d", suggestions.Next)
			}
		case "member roles":
			if err := workspace.RolesCall(roles, clearRoles); err != nil {
				return nil, err
			}
		case "record delete":
			// The preview writes nothing and reads either ID-file layout, so a
			// deletion it refuses is refused before the split.
			if confirm {
				if preview, err := s.DeleteRecords(ctx, get("member"), all, recordIDs, false); err != nil {
					return preview, err
				}
			}
		}
		if _, err := workspace.MemberRoles(roles); err != nil {
			return nil, err
		}
		if err := trySplit(); err != nil {
			return nil, fmt.Errorf("the ID files could not be moved to one file per member, so nothing was changed: %w", err)
		}
		var err error
		switch command {
		case "init", "template export":
			files, e := assets.Files()
			if e != nil {
				return nil, e
			}
			if command == "init" {
				readme, e := assets.Readme()
				if e != nil {
					return nil, e
				}
				err = s.Init(files, readme, get("name"), get("purpose"), get("member"), get("member-name"), band, roles)
			} else {
				err = s.Export(files)
			}
		case "status":
			return s.Status(ctx)
		case "agent settings":
			return s.EffectiveAgentSettings()
		case "agent configure":
			err = s.ConfigureAgent(get("host"), get("role"), get("model"), get("effort"), localTiering)
		case "agent setup":
			return s.SetupAgents(get("host"))
		case "agent dispatch":
			return s.DispatchAgent(get("host"), get("role"), get("task"), get("path"), get("plan"), get("intent"), shared, reviewRequested, get("repo"))
		case "check":
			issues, e := s.Check(ctx)
			checkFailed = len(issues) > 0
			return map[string]any{"issues": issues, "ok": !checkFailed}, e
		case "member list":
			return s.Members()
		case "member add":
			err = s.AddMember(get("id"), get("name"), band, get("language"), get("tone"), roles)
		case "member band":
			err = s.SetMemberBand(get("id"), band)
		case "member language":
			err = s.SetMemberLanguage(get("id"), get("language"))
		case "member tone":
			err = s.SetMemberTone(get("id"), get("tone"), clearTone)
		case "member roles":
			err = s.SetMemberRoles(get("id"), roles, clearRoles)
		case "member use":
			err = s.UseMember(get("id"))
		case "repo connect":
			err = s.Connect(ctx, get("id"), get("path"), get("base"), get("url"), get("default-branch"))
		case "repo clone", "repo init":
			// A clone needs a source: --url, or a shared record already naming
			// one. An init never takes one, so the store tells them apart.
			err = s.CreateRepository(ctx, get("id"), get("path"), get("base"), get("url"), get("default-branch"), command == "repo clone")
		case "repo base":
			err = s.SetBase(ctx, get("id"), get("branch"))
		case "repo remote":
			err = s.SetRemote(ctx, get("id"), get("url"), get("default-branch"))
		case "repo relate":
			err = s.Relate(get("from"), get("to"), get("description"))
		case "workspace connect":
			err = s.ConnectWorkspace(ctx, get("path"), get("base"), get("url"), get("default-branch"))
		case "workspace base":
			err = s.SetWorkspaceBase(ctx, get("branch"))
		case "workspace remote":
			err = s.SetWorkspaceRemote(ctx, get("url"), get("default-branch"))
		case "knowledge connect":
			err = s.ConnectKnowledge(ctx, get("id"), get("path"), get("url"), get("default-branch"), get("index"))
		case "knowledge clone":
			err = s.CloneKnowledge(ctx, get("id"), get("path"), get("url"), get("default-branch"), get("index"))
		case "knowledge remote":
			err = s.SetKnowledgeRemote(ctx, get("id"), get("url"), get("default-branch"), get("index"))
		case "knowledge sync":
			return s.SyncKnowledge(ctx, get("id"))
		case "knowledge list":
			return s.KnowledgeStates(ctx)
		case "repo inspect", "repo fetch":
			checkout, e := s.Repository(ctx, get("id"))
			if e != nil {
				return nil, e
			}
			path := checkout.Path
			if command == "repo inspect" {
				snapshot, e := workspace.Inspect(ctx, path)
				if e != nil {
					return nil, e
				}
				snapshot.BaseBranch = checkout.BaseBranch
				return snapshot, nil
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
			return s.CreateRecord(ctx, get("kind"), get("slug"), get("title"), get("intent"), repos, dependencies)
		case "record show":
			return s.FindRecord(get("id"))
		case "record list":
			return s.ListRecords(archived)
		case "record approve":
			if err := s.Note(get("id"), get("text"), "approve"); err != nil {
				return nil, err
			}
			// Approval hands work back rather than finishing it, and a bare
			// `ok: true` reads as the request being done. The next step to
			// recommend is said here, where the caller is standing, rather
			// than only in instructions read once at the start of a session.
			result := map[string]any{
				"approved":             get("id"),
				"planning_recommended": "approval is recorded; planning is the next step: recommend creating linked pNNNN-slug/plan.md folders to the person and stop, unless the same request already asked for plans, in which case inspect real code, create them, present the complete plans and stop",
			}
			recommendCommit(ctx, s, result, get("id"))
			return result, nil
		case "record complete":
			if err := s.Note(get("id"), get("text"), "complete"); err != nil {
				return nil, err
			}
			entries, e := s.KnowledgeCandidates(get("id"))
			result := map[string]any{"completed": get("id"), "knowledge_candidates": entries}
			// A bare list of catalog lines reads as information, and a caller
			// that has just been told "completed" treats the request as done.
			// Naming the outstanding act beside the data is what `setup_required`
			// does for a dispatch, so completion says it the same way.
			if len(entries) > 0 {
				result["reconcile_required"] = "judge each entry above against what this plan changed: edit the note and its catalog entry together and move its reviewed date, or record in the completion note that it changed nothing"
			}
			ids := []string{get("id")}
			if plan, found := s.FindRecord(get("id")); found == nil && plan.Intent != "" {
				ids = append(ids, plan.Intent)
			}
			recommendCommit(ctx, s, result, ids...)
			return result, e
		case "record delete":
			return s.DeleteRecords(ctx, get("member"), all, recordIDs, confirm)
		case "record dependencies":
			err = s.SetDependencies(get("id"), dependencies)
		case "record order":
			return s.Order(get("intent"), get("mode"), plans...)
		case "context find":
			// Completion names the knowledge a plan made due; a change made
			// without a plan has no record to name it, so the repository asks
			// for the same candidates and is told the same act is outstanding.
			if get("query") == "" && get("repo") == "" {
				return nil, fmt.Errorf("pass --query to search, or --repo for the entries claiming one repository")
			}
			if repository := get("repo"); repository != "" {
				if get("query") != "" {
					return nil, fmt.Errorf("pass --query or --repo, not both")
				}
				cfg, e := s.Config()
				if e != nil {
					return nil, e
				}
				if _, known := cfg.Repositories[repository]; !known {
					return nil, fmt.Errorf("unknown repository: %s", repository)
				}
				entries, e := s.CatalogEntriesFor([]string{repository})
				result := map[string]any{"repository": repository, "knowledge_candidates": entries}
				if len(entries) > 0 {
					result["reconcile_required"] = "judge each entry above against what this change actually altered: edit the note and its catalog entry together and move its reviewed date, or report that it changed nothing"
				}
				return result, e
			}
			return s.FindContext(ctx, get("query"))
		case "worktree prepare":
			return s.Prepare(ctx, get("repo"), get("plan"), get("branch"), get("start"), get("path"), reuse, workspace.ReuseOptions{Mode: get("copy-mode"), Paths: copyPaths})
		case "worktree list":
			return s.Worktrees(ctx, get("repo"))
		case "worktree inspect":
			return s.InspectWorktree(ctx, get("repo"), get("path"))
		case "worktree move":
			err = s.MoveWorktree(ctx, get("repo"), get("path"), get("to"))
		case "worktree remove":
			// The empty .worktrees/ folders the removal left are removed with it and
			// listed, the way prune lists what it dropped. One that could not be
			// removed is a warning, not a failure: the worktree is already gone.
			removed, warnings, e := s.RemoveWorktree(ctx, get("repo"), get("path"), discard)
			if e != nil {
				return nil, e
			}
			result := map[string]any{"operation": command, "workspace": s.Root, "ok": true, "removed_directories": removed}
			if len(warnings) > 0 {
				result["warnings"] = warnings
			}
			return result, nil
		case "worktree repair":
			err = s.RepairWorktree(ctx, get("repo"), get("path"))
		case "worktree prune":
			dropped, e := s.PruneWorktreeAssociations(ctx, get("repo"))
			if e != nil {
				return nil, e
			}
			return map[string]any{"operation": command, "workspace": s.Root, "ok": true, "repository": get("repo"), "dropped": dropped}, nil
		}
		return map[string]any{"operation": command, "workspace": s.Root, "ok": err == nil}, err
	}
	var result any
	if readOnly {
		result, err = action()
	} else {
		result, err = s.WithLock(ctx, func() (any, error) {
			// Taken before anything is written, so the notice can name every
			// shared file this command changed, members.yaml included, rather
			// than the files a split is expected to touch.
			before, snapErr := s.SharedFiles()
			result, err := action()
			if err == nil {
				// The command did what it was asked, so a split it could not make
				// afterwards is a warning: the next changing command refuses with
				// the same reason until a person fixes it.
				if e := trySplit(); e != nil {
					splitWarning = fmt.Sprintf("warning: %s was done, but the ID files could not be moved to one file per member yet: %v\n", command, e)
				}
			}
			if split != nil && snapErr == nil {
				if changed, e := s.ChangedSince(before); e == nil {
					split.Files = changed
				}
			}
			return result, err
		})
	}
	// The split is reported on the error stream in both output modes, the way
	// the pinned-version warning is, so no command's result changes shape. It
	// is reported even when the command then failed, because it still happened.
	if split != nil {
		fmt.Fprint(errOut, split.Notice())
	}
	fmt.Fprint(errOut, splitWarning)
	// A refused deletion still shows its preview: the records asked for and the
	// ones blocking them, which is what the person needs to decide what to add.
	// It deleted nothing and stays a failure.
	var refused *workspace.DeletionRefused
	if errors.As(err, &refused) {
		result = refused.Preview
	} else if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	switch {
	case *asJSON:
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		err = enc.Encode(result)
	case command == "record show":
		_, err = fmt.Fprint(out, result.(workspace.Record).Content)
	case command == "record create":
		// A record's body is a document. Rendered as a YAML scalar it arrives as
		// one escaped line, which is the shape `record show` already refuses to
		// print. The reservation and the path lead, because they are what the
		// caller asked the CLI for, and the body follows as itself.
		created := result.(workspace.Record)
		body := created.Content
		created.Content = ""
		if err = writeYAML(out, created); err == nil && body != "" {
			_, err = fmt.Fprintf(out, "\n%s", body)
		}
	default:
		err = writeYAML(out, result)
	}
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	if refused != nil {
		fmt.Fprintln(errOut, refused)
		return 1
	}
	if checkFailed {
		return 1
	}
	return 0
}
