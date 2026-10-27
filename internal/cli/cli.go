package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

const Help = `Context Circuit — shared workspace operations

Usage: context-circuit-cli [--workspace PATH] [--json] COMMAND [OPTIONS]
The global options --workspace PATH and --json may also follow the full command,
as in member list --json, but not sit between a group and its subcommand.

init                  --name NAME --purpose TEXT --member ID --member-name NAME
                      [--member-band N] [--member-role ROLE ...] (the first member,
                      selected here; band 1 unless another is given)
status                your open work first, then workspace, members, bindings,
                      and Git state
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
record list           [--member ID | --all] [--archived] (default: your own records;
                      --all adds every member's and records with no author)
record archive        [--id ID ...] (move under _archive/: your finished records,
                      or exactly the named ones; IDs stay reserved)
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
// not say. It never stages, commits or pushes, because committing and pushing
// workspace records is the person's decision.
func recommendCommit(ctx context.Context, s *workspace.Store, result map[string]any, ids ...string) {
	files, unchecked := s.UncommittedRecords(ctx, ids...)
	if unchecked != "" {
		result["commit_unchecked"] = unchecked + "; whether these records and the ID files are committed was not checked"
		return
	}
	if len(files) > 0 {
		result["commit_recommended"] = commitRecommendation{
			Files: files,
			Note:  "these records and ID files are not committed, so other clones cannot see them yet: ask the person whether to commit and push them, or commit only; committing and pushing need the person's answer as their authorization, and nothing was staged, committed or pushed",
		}
	}
}

// usage is the full help with its usage line naming the command this build is
// installed as. Help itself stays the public text, which signatures reads.
func usage(cliCommand string) string {
	return strings.Replace(Help, "Usage: "+publicCommand+" ", "Usage: "+cliCommand+" ", 1)
}

// helpTopic answers `help` and `help TOPIC`.
func helpTopic(args []string, out, errOut io.Writer, cliCommand string) int {
	if len(args) > 1 {
		topic := strings.Join(args[1:], " ")
		text, ok := CommandHelp(topic)
		if !ok {
			// Printing the whole manual for a question about one command
			// answers something the caller did not ask and reads as success.
			fmt.Fprintf(errOut, "no help topic: %s\ntry: %s help\n", topic, cliCommand)
			return 2
		}
		fmt.Fprint(out, text)
		return 0
	}
	fmt.Fprint(out, usage(cliCommand))
	return 0
}

// resolve names the command a line asks for, joining a group to its
// subcommand. ok is false when it printed why not, with the exit code to use.
func resolve(args []string, errOut io.Writer, cliCommand string) (name string, cmd command, rest []string, code int, ok bool) {
	name, rest = args[0], args[1:]
	if groups[name] {
		if len(rest) == 0 {
			fmt.Fprint(errOut, usage(cliCommand))
			return "", command{}, nil, 2, false
		}
		// A group with an option where its subcommand belongs is a caller
		// looking for help, not an unknown command. Name the real fault.
		if strings.HasPrefix(rest[0], "-") {
			fmt.Fprintf(errOut, "unknown option for %s: %s\ntry: %s help %s\n", name, rest[0], cliCommand, name)
			return "", command{}, nil, 2, false
		}
		name += " " + rest[0]
		rest = rest[1:]
	}
	cmd, found := commands[name]
	if !found {
		fmt.Fprintf(errOut, "unknown command: %s\n", name)
		// An empty argument is a command name with no fields at all, so there
		// is no group to offer a topic for. Reporting it is still right; going
		// looking for its first word is what crashes.
		if topic := strings.Fields(name); len(topic) > 0 && groups[topic[0]] {
			fmt.Fprintf(errOut, "try: %s help %s\n", cliCommand, topic[0])
		}
		return "", command{}, nil, 2, false
	}
	return name, cmd, rest, 0, true
}

// newInput defines the command's flags, and the global options after them.
func newInput(name string, cmd command, errOut io.Writer, root *string, asJSON *bool) *input {
	in := &input{name: name, f: flag.NewFlagSet(name, flag.ContinueOnError), values: map[string]*string{}}
	in.f.SetOutput(errOut)
	in.add(cmd.flags...)
	if cmd.extra != nil {
		cmd.extra(in)
	}
	// The global options are accepted after the command as well as before it,
	// because a caller adding --json or --workspace to a command line they
	// already typed puts it at the end. A command defining a flag of the same
	// name keeps its own meaning; none does today.
	if in.f.Lookup("workspace") == nil {
		in.f.StringVar(root, "workspace", *root, "workspace root")
	}
	if in.f.Lookup("json") == nil {
		in.f.BoolVar(asJSON, "json", *asJSON, "structured JSON output")
	}
	// The flag package still reports the faulty flag itself; the usage that
	// follows is this command's, printed by parse.
	in.f.Usage = func() {}
	return in
}

// parse reads the command's options. stop is true when the run is over, with
// the exit code to return.
func (in *input) parse(args []string, out, errOut io.Writer, cliCommand string) (code int, stop bool) {
	if err := in.f.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			text, _ := CommandHelp(in.name)
			fmt.Fprint(out, text)
			return 0, true
		}
		fmt.Fprint(errOut, commandUsage(cliCommand, in.name))
		return 2, true
	}
	if in.f.NArg() != 0 {
		fmt.Fprintln(errOut, "unexpected positional arguments")
		return 2, true
	}
	return 0, false
}

// missingFlag names the first required option the line left out, or "".
func (in *input) missingFlag(cmd command) string {
	// Sorted, so a call missing several flags is told the same one every time.
	for _, key := range slices.Sorted(maps.Keys(in.values)) {
		if *in.values[key] == "" && !slices.Contains(alwaysOptional, key) && !slices.Contains(cmd.optional, key) {
			return key
		}
	}
	if cmd.requiredNumber != "" && !provided(in.f, cmd.requiredNumber) {
		return cmd.requiredNumber
	}
	return ""
}

// openStore opens the workspace a command works on, creating its root first for
// the commands that make one, and applies the workspace's CLI pin. stop is true
// when the run is over before the command, with the exit code to return: the
// pinned CLI ran instead, the pin step refused, or the workspace did not open.
func openStore(cmd command, in *input, root *string, original []string, out, errOut io.Writer, version string) (s *workspace.Store, code int, stop bool) {
	if cmd.makeRoot {
		if cmd.fromPath {
			*root = in.get("path")
		}
		if err := os.MkdirAll(*root, 0755); err != nil {
			fmt.Fprintln(errOut, err)
			return nil, 1, true
		}
	}
	s, err := workspace.Open(*root)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return nil, 1, true
	}
	s.Command = commandName(version)
	// An exported blank seed is no workspace yet, so --path names no pin.
	if !cmd.fromPath {
		if code, stop := applyPin(s, cmd, in, original, version, out, errOut); stop {
			return nil, code, true
		}
	}
	return s, 0, false
}

// splitter makes the migrations a changing command owes the workspace: the ID
// split and the rename of legacy archive/ folders to _archive/.
type splitter struct {
	s       *workspace.Store
	enabled bool
	split   *workspace.IDSplit
	rename  *workspace.ArchiveRename
	warning string
}

// migrationError says which migration could not be made.
type migrationError struct {
	what string
	err  error
}

func (e *migrationError) Error() string { return e.what + ": " + e.err.Error() }
func (e *migrationError) Unwrap() error { return e.err }

// try makes each migration once setup is complete and it is still pending. It
// is tried before the command, so a change lands in the migrated layout, and
// again after it, so the setup step that completed setup makes the migrations
// itself rather than leaving them to whichever command comes next.
func (p *splitter) try() *migrationError {
	if !p.enabled {
		return nil
	}
	state, err := p.s.Setup()
	if err != nil {
		return nil
	}
	if p.split == nil {
		pending, err := p.s.IDSplitPending()
		if err == nil && pending {
			p.split, err = p.s.SplitIDLedger(state.Active)
		}
		if err != nil {
			return &migrationError{"the ID files could not be moved to one file per member", err}
		}
	}
	if p.rename == nil {
		pending, err := p.s.ArchiveRenamePending()
		if err == nil && pending {
			p.rename, err = p.s.RenameLegacyArchives()
		}
		if err != nil {
			return &migrationError{"the archive/ folders could not be renamed to _archive/", err}
		}
	}
	return nil
}

// locked runs action under the write lock and reports the files it changed.
func (p *splitter) locked(ctx context.Context, name string, action func() (any, error)) (any, error) {
	return p.s.WithLock(ctx, func() (any, error) {
		// Taken before anything is written, so the notice can name every
		// shared file this command changed, members.yaml included, rather
		// than the files a split is expected to touch.
		before, snapErr := p.s.SharedFiles()
		result, err := action()
		if err == nil {
			// The command did what it was asked, so a split it could not make
			// afterwards is a warning: the next changing command refuses with
			// the same reason until a person fixes it.
			if e := p.try(); e != nil {
				p.warning = fmt.Sprintf("warning: %s was done, but %s yet: %v\n", name, e.what, e.err)
			}
		}
		if p.split != nil && snapErr == nil {
			if changed, e := p.s.ChangedSince(before); e == nil {
				p.split.Files = changed
			}
		}
		return result, err
	})
}

// splitEnabled reports whether this call tries the ID split: it changes the
// workspace, and its command does not say the split is no part of it.
func splitEnabled(cmd command, in *input) bool {
	return !cmd.readOnly && (cmd.splits == nil || cmd.splits(in))
}

// execute runs the command, under the write lock unless it only reads.
func execute(ctx context.Context, s *workspace.Store, cmd command, in *input, errOut io.Writer) (any, error) {
	p := &splitter{s: s, enabled: splitEnabled(cmd, in)}
	action := func() (any, error) {
		if !cmd.readOnly && !cmd.setup {
			if err := s.RequireSetup(); err != nil {
				return nil, fmt.Errorf("setup is incomplete, so nothing was changed: %w", err)
			}
		}
		if cmd.precheck != nil {
			if err := cmd.precheck(ctx, s, in); err != nil {
				return nil, err
			}
		}
		if _, err := workspace.MemberRoles(in.roles); err != nil {
			return nil, err
		}
		if e := p.try(); e != nil {
			// A split already made is reported on its own; only the command
			// itself is then left undone.
			if p.split != nil {
				return nil, fmt.Errorf("%s, so the command was not run: %w", e.what, e.err)
			}
			return nil, fmt.Errorf("%s, so nothing was changed: %w", e.what, e.err)
		}
		return cmd.run(ctx, s, in)
	}
	var result any
	var err error
	if cmd.readOnly {
		result, err = action()
	} else {
		result, err = p.locked(ctx, in.name, action)
	}
	// The split is reported on the error stream in both output modes, the way
	// the pinned-version warning is, so no command's result changes shape. It
	// is reported even when the command then failed, because it still happened.
	if p.split != nil {
		fmt.Fprint(errOut, p.split.Notice())
	}
	if p.rename != nil {
		fmt.Fprint(errOut, p.rename.Notice())
	}
	fmt.Fprint(errOut, p.warning)
	return result, err
}

// render writes a command's result to the result stream.
func render(out io.Writer, asJSON bool, name string, result any) error {
	switch {
	case asJSON:
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	case name == "record show":
		_, err := fmt.Fprint(out, result.(workspace.Record).Content)
		return err
	case name == "record create":
		// A record's body is a document. Rendered as a YAML scalar it arrives as
		// one escaped line, which is the shape `record show` already refuses to
		// print. The reservation and the path lead, because they are what the
		// caller asked the CLI for, and the body follows as itself.
		created := result.(workspace.Record)
		body := created.Content
		created.Content = ""
		err := writeYAML(out, created)
		if err == nil && body != "" {
			_, err = fmt.Fprintf(out, "\n%s", body)
		}
		return err
	default:
		return writeYAML(out, result)
	}
}

func Run(ctx context.Context, args []string, out, errOut io.Writer, version string) int {
	// Hints name the command this version is installed as, so a variant build
	// never points its caller at the public command.
	cliCommand := commandName(version)
	global := flag.NewFlagSet(cliCommand, flag.ContinueOnError)
	global.SetOutput(errOut)
	root := global.String("workspace", ".", "workspace root")
	asJSON := global.Bool("json", false, "structured JSON output")
	versionFlag := global.Bool("version", false, "show version")
	// A hand-off to the pinned CLI passes the line on exactly as it was given.
	original := slices.Clone(args)
	if err := global.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(out, usage(cliCommand))
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
		return helpTopic(args, out, errOut, cliCommand)
	}
	name, cmd, args, code, ok := resolve(args, errOut, cliCommand)
	if !ok {
		return code
	}
	in := newInput(name, cmd, errOut, root, asJSON)
	if code, stop := in.parse(args, out, errOut, cliCommand); stop {
		return code
	}
	if key := in.missingFlag(cmd); key != "" {
		fmt.Fprintf(errOut, "missing --%s\n", key)
		return 2
	}
	s, code, stop := openStore(cmd, in, root, original, out, errOut, version)
	if stop {
		return code
	}
	result, err := execute(ctx, s, cmd, in, errOut)
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
	if err := render(out, *asJSON, name, result); err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	fmt.Fprint(errOut, in.notice)
	if refused != nil {
		fmt.Fprintln(errOut, refused)
		return 1
	}
	if in.checkFailed {
		return 1
	}
	return 0
}
