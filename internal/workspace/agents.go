package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/goccy/go-yaml"
)

type AgentSetting struct {
	Model  string `yaml:"model" json:"model"`
	Effort string `yaml:"effort" json:"effort"`
}
type RoleTiering struct {
	Hosts map[string]map[string]AgentSetting `yaml:"hosts" json:"hosts"`
}
type Role struct {
	Description, Instructions string
	ReadOnly                  bool
}

var Roles = map[string]Role{
	"explorer": {"Investigate a bounded codebase question and return evidence.", "Inspect only the assigned repository question. Return file locations, observed behavior, and uncertainties. Run no tests, linters, or builds: report behavior you read, never behavior you confirmed. Do not edit files or launch other agents.", true},
	"planner": {"Plan an approved intent or specified standalone outcome using repository evidence.", `Read the approved intent above, or the specified standalone outcome in the task, and the relevant code. Then answer under these headings:

- Verdict: feasible, feasible-with-changes, or not-feasible.
- Approach: how the outcome is reached.
- Tasks and order: each task naming its repository, its paths, and the tasks it depends on.
- Risks and checks: the check commands and where they are defined.
- Plan shape: single, or a proposed split with reasons.

Run no tests, linters, or builds: a check you name is one you read, never one you saw pass.
Not-feasible is a complete answer: return it with its reasons instead of a plan.
The coordinator writes the plan record and decides any split, so no plan exists yet and none is yours to number.
Do not implement, do not request a separate plan approval, do not edit files, and do not launch other agents.`, true},
	"worker": {"Implement a bounded part of a requested plan and run normal checks.", `Implement the plan above in the working directory above. Own only your assignment and the paths it names.

The working directory is a repository checkout. Read the AGENTS.md or CLAUDE.md
at its root if there is one, and follow it where it is more specific than this
brief. If the repository ships skills, match your assignment against their
descriptions and open only the ones that match. A repository with neither needs
nothing extra.

Run the repositories' ordinary checks, then commit on the branch this working
directory is already on — leave nothing uncommitted when you report. Take as many
commits as the work naturally needs, and commit failing or partial work too: what
you leave uncommitted is work the coordinator's next step cannot see. Follow the
repository's commit convention and add no attribution of any kind.

Report, as your result:

- every file you changed, and what changed in it;
- the exact check commands you ran and their real outcome, failures included;
- what remains unfinished, and anything you had to assume.

The coordinator integrates your work from that report and does not re-run your checks, so a check you did not run is one nobody ran. Report a failure plainly rather than working around it.

Do not write or reconcile durable project knowledge; report what the coordinator should record. Do not launch other agents. Committing on your own branch is part of the work; pushing, opening a pull request, merging into the base branch, or delivering changes any other way needs explicit authorization.`, false},
	"reviewer": {"Independently review a diff only when the user requests review.", "Perform the explicitly requested independent read-only review. Inspect the supplied diff and current revision against the intent's success criteria and relevant surrounding code. Report actionable findings with file locations and limitations. Never edit files, run commands that change files, dispatch repairs, or post external comments. Do not launch other agents.", true},
}

func validHost(host string) error {
	if host != "codex" && host != "claude-code" && host != "cursor" {
		return fmt.Errorf("host must be codex, claude-code, or cursor")
	}
	return nil
}
func validateSetting(host string, setting AgentSetting) error {
	if err := Text(setting.Model); err != nil {
		return err
	}
	if strings.ContainsAny(setting.Model, "\t[] ") {
		return fmt.Errorf("use a plain model ID or inherit; effort has its own field")
	}
	for _, r := range setting.Model {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return fmt.Errorf("invalid model ID")
		}
	}
	switch setting.Effort {
	case "inherit", "none", "minimal", "low", "medium", "high", "xhigh", "max", "ultra":
	default:
		return fmt.Errorf("unsupported effort: %s", setting.Effort)
	}
	if host == "claude-code" && (setting.Effort == "none" || setting.Effort == "minimal" || setting.Effort == "ultra") {
		return fmt.Errorf("Claude Code does not expose this effort setting")
	}
	if host == "cursor" && setting.Model == "inherit" && setting.Effort != "inherit" {
		return fmt.Errorf("Cursor effort selection needs an explicit model ID")
	}
	return nil
}

// The shared file is the workspace's agreed tiering; the local one is this
// machine's answer to it, because one roster's members run different hosts and
// pay for different models. Local overrides a single host/role rather than
// replacing the file, so a member who changes one role still tracks every
// later change the team makes to the rest.
const (
	SharedTiering = ".context-circuit/role-tiering.yaml"
	LocalTiering  = ".context-circuit/role-tiering.local.yaml"
)

// AgentSettings reports the effective tiering: shared, with each host/role the
// local file names replacing its counterpart. Overrides names those keys, since
// a merged view that cannot say which half supplied a value claims more than it
// established.
func (s *Store) AgentSettings() (RoleTiering, error) {
	cfg, _, err := s.agentTiering()
	return cfg, err
}

// EffectiveTiering is what a reader of the settings command needs: the values in
// force, and which of them this machine overrode.
type EffectiveTiering struct {
	Hosts          map[string]map[string]AgentSetting `yaml:"hosts" json:"hosts"`
	LocalOverrides []string                           `yaml:"local_overrides,omitempty" json:"local_overrides,omitempty"`
}

func (s *Store) EffectiveAgentSettings() (EffectiveTiering, error) {
	cfg, overrides, err := s.agentTiering()
	return EffectiveTiering{cfg.Hosts, overrides}, err
}

func (s *Store) agentTiering() (RoleTiering, []string, error) {
	var cfg RoleTiering
	if err := s.YAML(SharedTiering, &cfg); err != nil {
		return cfg, nil, err
	}
	var local RoleTiering
	if err := s.YAML(LocalTiering, &local); err != nil && !os.IsNotExist(err) {
		return cfg, nil, fmt.Errorf("%s: %w", LocalTiering, err)
	}
	var overrides []string
	for host, settings := range local.Hosts {
		if cfg.Hosts[host] == nil {
			return cfg, nil, fmt.Errorf("%s: host is absent from the shared tiering: %s", LocalTiering, host)
		}
		for role, setting := range settings {
			cfg.Hosts[host][role] = setting
			overrides = append(overrides, host+"/"+role)
		}
	}
	sort.Strings(overrides)
	for host, settings := range cfg.Hosts {
		if err := validHost(host); err != nil {
			return cfg, nil, err
		}
		for name, setting := range settings {
			if _, ok := Roles[name]; !ok {
				return cfg, nil, fmt.Errorf("unknown agent role: %s", name)
			}
			if err := validateSetting(host, setting); err != nil {
				return cfg, nil, fmt.Errorf("%s/%s: %w", host, name, err)
			}
		}
	}
	return cfg, overrides, nil
}

// ConfigureAgent writes one host/role setting. Local keeps the change on this
// machine; without it the change is to the shared file every member reads, so
// which of the two is being edited is stated rather than inferred.
func (s *Store) ConfigureAgent(host, role, model, effort string, local bool) error {
	if err := validHost(host); err != nil {
		return err
	}
	if _, ok := Roles[role]; !ok {
		return fmt.Errorf("unknown agent role: %s", role)
	}
	setting := AgentSetting{model, effort}
	if err := validateSetting(host, setting); err != nil {
		return err
	}
	var shared RoleTiering
	if err := s.YAML(SharedTiering, &shared); err != nil {
		return err
	}
	// The shared file names the hosts a workspace supports; a local override
	// answers one of them rather than introducing one nobody else has.
	if shared.Hosts[host] == nil {
		return fmt.Errorf("host is absent from role-tiering.yaml: %s", host)
	}
	if !local {
		return s.Update(SharedTiering, []string{"hosts", host, role}, setting, 0644)
	}
	if _, err := s.Read(LocalTiering); os.IsNotExist(err) {
		seed := "# Machine-local role settings. Each host/role here replaces its counterpart\n" +
			"# in role-tiering.yaml; everything absent keeps tracking the shared file.\n" +
			"hosts:\n  " + host + ": {}\n"
		if err := s.Write(LocalTiering, []byte(seed), 0600); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// A later override of a second host has no node to attach to yet, so the
	// host is opened before the role is written into it.
	var current RoleTiering
	if err := s.YAML(LocalTiering, &current); err != nil {
		return err
	}
	if current.Hosts[host] == nil {
		return s.Update(LocalTiering, []string{"hosts", host}, map[string]AgentSetting{role: setting}, 0600)
	}
	return s.Update(LocalTiering, []string{"hosts", host, role}, setting, 0600)
}

func (s *Store) agentSetting(host, role string) (AgentSetting, error) {
	if err := validHost(host); err != nil {
		return AgentSetting{}, err
	}
	if _, ok := Roles[role]; !ok {
		return AgentSetting{}, fmt.Errorf("unknown agent role: %s", role)
	}
	cfg, err := s.AgentSettings()
	if err != nil {
		return AgentSetting{}, err
	}
	if setting, ok := cfg.Hosts[host][role]; ok {
		return setting, nil
	}
	return AgentSetting{"inherit", "inherit"}, nil
}

func roleFile(host, role string, setting AgentSetting) (string, []byte, error) {
	r := Roles[role]
	name := "cc_" + role
	if host == "codex" {
		text := "name = " + strconv.Quote(name) + "\ndescription = " + strconv.Quote(r.Description) + "\n"
		if setting.Model != "inherit" {
			text += "model = " + strconv.Quote(setting.Model) + "\n"
		}
		if setting.Effort != "inherit" {
			text += "model_reasoning_effort = " + strconv.Quote(setting.Effort) + "\n"
		}
		if r.ReadOnly {
			text += "sandbox_mode = \"read-only\"\n"
		}
		text += "developer_instructions = " + strconv.Quote(r.Instructions) + "\n"
		return ".codex/agents/cc-" + role + ".toml", []byte(text), nil
	}
	fields := map[string]any{"name": "cc-" + role, "description": r.Description, "model": setting.Model}
	directory := ".claude"
	if host == "claude-code" {
		if setting.Effort != "inherit" {
			fields["effort"] = setting.Effort
		}
		if r.ReadOnly {
			fields["tools"] = "Read, Glob, Grep"
		}
	} else {
		directory = ".cursor"
		fields["readonly"] = r.ReadOnly
		if setting.Effort != "inherit" {
			fields["model"] = setting.Model + "[effort=" + setting.Effort + "]"
		}
	}
	front, err := yaml.Marshal(fields)
	return directory + "/agents/cc-" + role + ".md", []byte("---\n" + string(front) + "---\n\n" + r.Instructions + "\n"), err
}

type agentInventory struct {
	Files map[string]string `yaml:"files"`
}

func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

// SetupAgents only refreshes files previously written by this command, or files
// already identical to the requested output. User edits are never overwritten.
// Hosts names every coding host a workspace can be opened in. A workspace is
// used from more than one, so setup writes them all unless one is named.
var Hosts = []string{"codex", "claude-code", "cursor"}

// SetupAgents writes native role definitions. An empty host covers every host,
// because the same workspace is opened in different ones and a definition
// installed for only the host that happened to run setup leaves the next one
// with nothing to invoke.
func (s *Store) SetupAgents(host string) ([]string, error) {
	if host != "" {
		return s.setupHost(host)
	}
	var paths []string
	for _, each := range Hosts {
		written, err := s.setupHost(each)
		paths = append(paths, written...)
		if err != nil {
			return paths, err
		}
	}
	return paths, nil
}

func (s *Store) setupHost(host string) ([]string, error) {
	if err := validHost(host); err != nil {
		return nil, err
	}
	inventoryPath := ".context-circuit/local/agents-" + host + ".yaml"
	old := agentInventory{Files: map[string]string{}}
	if err := s.YAML(inventoryPath, &old); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	files := map[string][]byte{}
	var paths []string
	for _, role := range []string{"explorer", "planner", "worker", "reviewer"} {
		setting, err := s.agentSetting(host, role)
		if err != nil {
			return nil, err
		}
		path, content, err := roleFile(host, role, setting)
		if err != nil {
			return nil, err
		}
		current, err := s.Read(path)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if err == nil && digest(current) != digest(content) && digest(current) != old.Files[path] {
			return nil, fmt.Errorf("preserve customized agent definition: %s", path)
		}
		files[path], paths = content, append(paths, path)
	}
	next := agentInventory{Files: map[string]string{}}
	for _, path := range paths {
		if err := s.Write(path, files[path], 0644); err != nil {
			return paths, err
		}
		next.Files[path] = digest(files[path])
	}
	return paths, s.WriteYAML(inventoryPath, next, 0600)
}

type Dispatch struct {
	Host             string       `json:"host"`
	Role             string       `json:"role"`
	AgentType        string       `json:"agent_type"`
	Setting          AgentSetting `json:"setting"`
	WorkingDirectory string       `json:"working_directory"`
	ReadOnly         bool         `json:"read_only"`
	Prompt           string       `json:"prompt"`
	LaunchRequired   bool         `json:"launch_required"`
	// A specification names an agent type the host may not have registered,
	// because role files are host-local, gitignored, and written only by setup.
	// Reporting the absent definition turns a silent non-launch into a fact the
	// caller can act on; it is not a refusal, since the prompt still carries
	// everything a live spawn tool needs when native roles are unavailable.
	DefinitionPath      string `json:"definition_path"`
	DefinitionInstalled bool   `json:"definition_installed"`
	SetupRequired       string `json:"setup_required,omitempty"`
}

func (s *Store) DispatchAgent(host, role, task, directory, plan, intent string, shared, reviewRequested bool, repository ...string) (Dispatch, error) {
	repo := ""
	if len(repository) > 0 {
		repo = repository[0]
	}
	if repo != "" && role != "worker" {
		return Dispatch{}, errors.New("--repo assigns a worker")
	}
	setting, err := s.agentSetting(host, role)
	if err != nil {
		return Dispatch{}, err
	}
	// The record is the assignment. A planner is dispatched against a whole
	// approved intent and a worker against the plan quoted in its brief, so a
	// task on either can only restate what the brief already carries, and a
	// required flag leaves the coordinator nothing to write but that
	// restatement. A worker splitting one plan with others is the exception:
	// its slice is named in the task and nowhere else.
	if role == "planner" && intent != "" && task != "" {
		return Dispatch{}, errors.New("an intent planner takes the whole approved intent, so it takes no --task")
	}
	briefed := role == "planner" || role == "worker"
	if task != "" || !briefed {
		if err := Text(task); err != nil {
			return Dispatch{}, err
		}
	}
	if role == "reviewer" && !reviewRequested {
		return Dispatch{}, fmt.Errorf("independent review requires an explicit user request; use --review-requested only to represent that request")
	}
	if shared && role != "worker" {
		return Dispatch{}, errors.New("--shared describes workers sharing one worktree")
	}
	// A planner decides the plan shape, so a plan cannot already exist to hand
	// it: numbering one first settles the split the planner was dispatched to
	// propose, and an intent holding several plans has no single ID to pass.
	if role == "planner" {
		if plan != "" {
			return Dispatch{}, errors.New("a planner proposes the plan shape, so it is dispatched against --intent, not an already numbered --plan")
		}
		if intent == "" && task == "" {
			return Dispatch{}, errors.New("a planner needs --intent or --task specifying the standalone outcome")
		}
	} else if intent != "" {
		return Dispatch{}, errors.New("--intent dispatches a planner; other roles take --plan")
	}
	work, err := Open(directory)
	if err != nil {
		return Dispatch{}, err
	}
	r := Roles[role]
	name := "cc-" + role
	if host == "codex" {
		name = "cc_" + role
	}
	// The brief is read top to bottom by an agent that has seen nothing else, so
	// it states where the work happens, then what the work is, and only then
	// what to return. Every record it names is also quoted in full: an agent
	// that has to go and find a file reads the whole repository on the way.
	prompt := "# Working directory\n\n" + work.Root + "\n\nWorkspace root: " + s.Root
	// A read-only role was told where its working directory is and nothing
	// about what it may do there, which leaves how to read it an open question
	// on a host that also offers GUI automation and web search. One planner
	// answered it by driving an editor's interface to read a local file and
	// searching the public web for a repository sitting on disk.
	ownership := "You are reading, not changing. Change nothing in this working directory. Read and search the files here directly rather than driving another application, and consult no external source: what this task needs is this directory and this brief."
	if !r.ReadOnly {
		// Isolation differs by parallelism axis. One worker per plan owns its
		// whole worktree; several workers inside one worktree genuinely share
		// files. Saying the wrong one invites a worker to guess at edits it
		// cannot see, or to overwrite edits it can.
		if shared {
			ownership = "Other workers are editing this same working directory. Preserve their edits, adapt your changes, and stay within your assigned paths."
		} else {
			ownership = "You are the sole owner of this working directory. Work only here. Do not create, switch, merge, push, or delete branches, and do not run git worktree."
		}
	}
	prompt += "\n\n# Ownership\n\n" + ownership
	// Whoever wrote the record wrote it in their own language, so the brief
	// reports that one rather than the language of whoever is dispatching.
	author := ""
	if plan != "" {
		record, err := s.FindRecord(plan)
		if err != nil {
			return Dispatch{}, err
		}
		if !strings.HasPrefix(record.ID, "p") {
			return Dispatch{}, errors.New("--plan needs a plan ID")
		}
		author = record.CreatedBy
		prompt += "\n\n" + planBrief(record, s.Root)
		if role == "worker" {
			if repo == "" {
				if len(record.Repositories) != 1 {
					return Dispatch{}, errors.New("multi-repository worker dispatch needs --repo")
				}
				repo = record.Repositories[0]
			}
			files, err := s.requiredPlanFiles(record, repo)
			if err != nil {
				return Dispatch{}, err
			}
			prompt += "\n\nRepository assignment: " + repo
			if len(files) > 0 {
				prompt += "\n\n# Required files relative to workspace root\n\n"
				for _, file := range files {
					prompt += "- " + file + "\n"
				}
				prompt += "\nResolve these paths against the workspace root above and read every required file before editing."
			}
			prompt += "\n\nWork only in your assigned repository worktree. Links in the plan body to other repositories' files are for human review; do not open them unless assigned here."
		}
	}
	if intent != "" {
		record, err := s.FindRecord(intent)
		if err != nil {
			return Dispatch{}, err
		}
		if !strings.HasPrefix(record.ID, "i") {
			return Dispatch{}, errors.New("--intent needs an intent ID")
		}
		// Planning an unapproved intent plans an outcome nobody agreed to.
		if record.ApprovedAt == "" {
			return Dispatch{}, fmt.Errorf("intent %s is not approved; approve it before planning it", record.ID)
		}
		author = record.CreatedBy
		prompt += "\n\n" + intentBrief(record, s.Root)
	}
	if role == "planner" && intent == "" {
		author, _ = s.ActiveMember()
	}
	if task != "" {
		prompt += "\n\n# Task\n\n" + task
	}
	prompt += "\n\n" + languageBrief(s.authorLanguage(author))
	prompt += "\n\n# What to return\n\n" + r.Instructions
	definition, _, err := roleFile(host, role, setting)
	if err != nil {
		return Dispatch{}, err
	}
	installed := false
	if resolved, e := s.Path(definition); e == nil {
		if info, e := os.Stat(resolved); e == nil && info.Mode().IsRegular() {
			installed = true
		}
	}
	setup := ""
	if !installed {
		setup = "context-circuit-cli --workspace " + s.Root + " agent setup --host " + host
	}
	return Dispatch{host, role, name, setting, work.Root, r.ReadOnly, prompt, true, definition, installed, setup}, nil
}

// authorLanguage reports the language a record's prose is written in: the
// authoring member's setting, and English for a member who records none, which
// is what every record written before the field assumed. A member the roster no
// longer carries reads as English rather than failing the dispatch, because a
// missing name is not a reason to refuse to start work.
func (s *Store) authorLanguage(member string) string {
	if member != "" {
		if members, err := s.Members(); err == nil {
			if recorded, ok := members.Members[member]; ok && recorded.Language != "" {
				return recorded.Language
			}
		}
	}
	return "English"
}

// languageBrief draws the boundary a worker cannot infer. A record in one
// language is quoted into the brief as authoritative, which invites a worker to
// answer it in that language — in comments, in commit messages, and worst of
// all in identifiers, where a domain term rendered into English silently
// disagrees with the glossary that named it.
//
// The worker never reads the workspace's own instructions: it works inside a
// repository worktree, under that repository's. So the boundary travels in the
// brief or nowhere.
func languageBrief(language string) string {
	return "# Language\n\nThe work described here is stated in " + language +
		". Write everything you put into the repository in English — code, comments, identifiers, tests, commit messages — unless this repository's own instructions say otherwise." +
		"\n\nNames are quoted, never translated, in either direction. Domain vocabulary keeps the form used here, and code identifiers keep the form the code uses. Translating either breaks the thing it names."
}

// planBrief states the facts the CLI can verify. A dependent plan's branch is
// prepared from its predecessors, so their work is already in its ancestry;
// saying so prevents a worker from reimplementing what it inherited.
func planBrief(record Record, root string) string {
	brief := "# Plan " + record.ID
	if record.Intent != "" {
		brief += " of intent " + record.Intent
	}
	brief += "\n\nFile: " + filepath.Join(root, record.Path)
	brief += "\nRepositories: " + strings.Join(record.Repositories, ", ")
	if len(record.DependsOn) > 0 {
		brief += "\nDepends on: " + strings.Join(record.DependsOn, ", ") +
			"\n\nTheir completed work is expected in this branch's ancestry. Build on it rather than reimplementing it." +
			"\nAny plan running beside you is invisible in a separate worktree and is reconciled later by the coordinator. Do not guess at its changes; report an interface you assume it may also be changing."
	}
	return brief + recordBody(record, "plan")
}

// intentBrief carries the approved outcome a planner plans against. No plan is
// named, because deciding whether there is one plan or several is the planner's
// answer to return rather than a number it is handed.
func intentBrief(record Record, root string) string {
	brief := "# Intent " + record.ID +
		"\n\nFile: " + filepath.Join(root, record.Path) +
		"\nApproved: " + string(record.ApprovedAt)
	if len(record.Repositories) > 0 {
		brief += "\nRepositories: " + strings.Join(record.Repositories, ", ")
	}
	return brief + recordBody(record, "intent")
}

// A quoted record is the copy the agent works from. Sending the agent to open
// the file instead costs a search of the whole repository to find it.
//
// The record is fenced because its own headings would otherwise read as sections
// of the brief that contains it, and a plan quoting code carries fences of its
// own, so the fence is always longer than the longest run already inside.
func recordBody(record Record, kind string) string {
	_, raw, err := splitRecord([]byte(record.Content))
	body := strings.TrimSpace(string(raw))
	if err != nil || body == "" {
		return ""
	}
	fence := strings.Repeat("`", max(3, longestRun(body, '`')+1))
	intro := "The full " + kind + " follows and is authoritative. You do not need to open the file."
	if kind == "plan" {
		intro = "The plan entry body follows and is authoritative. You do not need to open the entry file. If a required-files list follows, read those files before editing."
	}
	return "\n\n" + intro + "\n\n" +
		fence + "markdown\n" + body + "\n" + fence
}

func longestRun(text string, char byte) int {
	longest, run := 0, 0
	for i := 0; i < len(text); i++ {
		if text[i] == char {
			run++
			longest = max(longest, run)
		} else {
			run = 0
		}
	}
	return longest
}
