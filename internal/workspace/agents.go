package workspace

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
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

// Role is one agent role a brief and a native definition are made from. The
// four shipped roles are fixed; an extension role is read from the extension's
// manifest and instructions file on this machine, and its name is
// x-<extension>-<id>, which no shipped role's can equal.
type Role struct {
	Description, Instructions string
	ReadOnly                  bool
	// Headings are the fixed headings the role's result uses. A planner's are
	// transcribed into the plan record, so no extension may add one of them
	// again; an extension role's are the ones its manifest returns.
	Headings []string
	// Extension and ID name an extension role's source; both are empty for a
	// shipped role. Stages are the stages it serves, whose hook instructions
	// its brief carries, and Like the shipped role whose model and effort it
	// takes where the tiering does not name it.
	Extension, ID string
	Stages        []string
	Like          string
}

// edits reports whether a role changes its working directory: the worker, and
// an extension role with edit access. Such a role is dispatched against a plan
// and briefed as the worktree's owner.
func (r Role) edits(name string) bool {
	return !r.ReadOnly && (name == "worker" || r.Extension != "")
}

// ShippedRoles is the order setup writes the shipped definitions in.
var ShippedRoles = []string{"explorer", "planner", "worker", "reviewer"}

var Roles = map[string]Role{
	"explorer": {Description: "Investigate a bounded codebase question and return evidence.", Instructions: "Inspect only the assigned repository question. Return file locations, observed behavior, and uncertainties. Run no tests, linters, or builds: report behavior you read, never behavior you confirmed. Do not edit files or launch other agents.", ReadOnly: true},
	"planner": {Description: "Plan an approved intent or specified standalone outcome using repository evidence.", Instructions: `Read the approved intent above, or the specified standalone outcome in the task, and the relevant code. Then answer under these headings:

- Verdict: feasible, feasible-with-changes, or not-feasible.
- Approach: how the outcome is reached.
- Tasks and order: each task naming its repository, its paths, and the tasks it depends on.
- Risks and checks: the check commands and where they are defined.
- Plan shape: single, or a proposed split with reasons.

Run no tests, linters, or builds: a check you name is one you read, never one you saw pass.
Not-feasible is a complete answer: return it with its reasons instead of a plan.
The coordinator writes the plan record and decides any split, so no plan exists yet and none is yours to number.
Do not implement, do not request a separate plan approval, do not edit files, and do not launch other agents.`, ReadOnly: true,
		Headings: []string{"Verdict", "Approach", "Tasks and order", "Risks and checks", "Plan shape"}},
	"worker": {Description: "Implement a bounded part of a requested plan and run normal checks.", Instructions: `Implement the plan above in the working directory above. Own only your assignment and the paths it names.

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

Do not write or reconcile durable project knowledge; report what the coordinator should record. Do not launch other agents. Committing on your own branch is part of the work; pushing, opening a pull request, merging into the base branch, or delivering changes any other way needs explicit authorization.`},
	"reviewer": {Description: "Independently review a diff only when the user requests review.", Instructions: "Perform the explicitly requested independent read-only review. Inspect the supplied diff and current revision against the intent's success criteria and relevant surrounding code. Report actionable findings with file locations and limitations. Never edit files, run commands that change files, dispatch repairs, or post external comments. Do not launch other agents.", ReadOnly: true},
}

func validHost(host string) error {
	if host != "codex" && host != "claude-code" && host != "cursor" {
		return errors.New("host must be codex, claude-code, or cursor")
	}
	return nil
}
func validateSetting(host string, setting AgentSetting) error {
	if err := Text(setting.Model); err != nil {
		return err
	}
	if strings.ContainsAny(setting.Model, "\t[] ") {
		return errors.New("use a plain model ID or inherit; effort has its own field")
	}
	for _, r := range setting.Model {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return errors.New("invalid model ID")
		}
	}
	switch setting.Effort {
	case "inherit", "none", "minimal", "low", "medium", "high", "xhigh", "max", "ultra":
	default:
		return fmt.Errorf("unsupported effort: %s", setting.Effort)
	}
	if host == "claude-code" && (setting.Effort == "none" || setting.Effort == "minimal" || setting.Effort == "ultra") {
		return errors.New("Claude Code does not expose this effort setting")
	}
	if host == "cursor" && setting.Model == "inherit" && setting.Effort != "inherit" {
		return errors.New("Cursor effort selection needs an explicit model ID")
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
	slices.Sort(overrides)
	for host, settings := range cfg.Hosts {
		if err := validHost(host); err != nil {
			return cfg, nil, err
		}
		for name, setting := range settings {
			if !tieringRole(name) {
				return cfg, nil, fmt.Errorf("unknown agent role: %s", name)
			}
			if err := validateSetting(host, setting); err != nil {
				return cfg, nil, fmt.Errorf("%s/%s: %w", host, name, err)
			}
		}
	}
	return cfg, overrides, nil
}

// tieringRole reports whether the tiering may name a role: a shipped one, or
// an extension role's x- name. The tiering is shared and an extension may be
// off for one member and on for another, so an x- name is accepted whether or
// not its extension is on here.
func tieringRole(name string) bool {
	if _, ok := Roles[name]; ok {
		return true
	}
	return strings.HasPrefix(name, "x-") && Name(name) == nil
}

// ConfigureAgent writes one host/role setting. Local keeps the change on this
// machine; without it the change is to the shared file every member reads, so
// which of the two is being edited is stated rather than inferred.
func (s *Store) ConfigureAgent(host, role, model, effort string, local bool) error {
	if err := validHost(host); err != nil {
		return err
	}
	if _, err := s.resolveRole(role); err != nil {
		return err
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

// agentSetting is the model and effort in force for a role on host. An
// extension role the tiering does not name takes the setting of the shipped
// role its manifest says it is like, and inherits without one.
func (s *Store) agentSetting(host, name string, r Role) (AgentSetting, error) {
	if err := validHost(host); err != nil {
		return AgentSetting{}, err
	}
	cfg, err := s.AgentSettings()
	if err != nil {
		return AgentSetting{}, err
	}
	if setting, ok := cfg.Hosts[host][name]; ok {
		return setting, nil
	}
	if r.Like != "" {
		if setting, ok := cfg.Hosts[host][r.Like]; ok {
			return setting, nil
		}
	}
	return AgentSetting{"inherit", "inherit"}, nil
}

// agentType is the name a host invokes a role by: cc-<role>, and on Codex
// cc_<role> with every hyphen an underscore, so cc_x_parity_parity_verifier.
func agentType(host, role string) string {
	if host == "codex" {
		return strings.ReplaceAll("cc_"+role, "-", "_")
	}
	return "cc-" + role
}

// roleFile is a role's native definition for host. An extension role's body is
// what its brief returns under "What to return", limits included, because a
// host may run the definition with a prompt that carries none of the brief.
func roleFile(host, role string, r Role, setting AgentSetting) (string, []byte, error) {
	body := r.Instructions
	if r.Extension != "" {
		body = returnSection(r, nil)
	}
	name := agentType(host, role)
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
		text += "developer_instructions = " + strconv.Quote(body) + "\n"
		return ".codex/agents/cc-" + role + ".toml", []byte(text), nil
	}
	fields := map[string]any{"name": name, "description": r.Description, "model": setting.Model}
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
	return directory + "/agents/cc-" + role + ".md", []byte("---\n" + string(front) + "---\n\n" + body + "\n"), err
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

// AgentSetup is what setup did: the files it wrote, the extension role
// definitions it removed because their extension is gone or off for this
// member, and each such definition it kept because a person edited it.
type AgentSetup struct {
	Written []string    `json:"written"`
	Removed []string    `json:"removed,omitempty"`
	Kept    []KeptAgent `json:"kept,omitempty"`
}

// KeptAgent is a definition setup no longer produces and left in place.
type KeptAgent struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

// SetupAgents writes native role definitions. An empty host covers every host,
// because the same workspace is opened in different ones and a definition
// installed for only the host that happened to run setup leaves the next one
// with nothing to invoke.
func (s *Store) SetupAgents(host string) (AgentSetup, error) {
	result := AgentSetup{Written: []string{}}
	hosts := Hosts
	if host != "" {
		hosts = []string{host}
	}
	for _, each := range hosts {
		done, err := s.setupHost(each)
		result.Written = append(result.Written, done.Written...)
		result.Removed = append(result.Removed, done.Removed...)
		result.Kept = append(result.Kept, done.Kept...)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// definitionRole is one role setup writes a definition for.
type definitionRole struct {
	Name string
	Role Role
}

// definitionRoles are the roles setup writes: the shipped four, then each
// extension role on for this machine's member, in extensions.yaml order.
func (s *Store) definitionRoles() ([]definitionRole, error) {
	var roles []definitionRole
	for _, name := range ShippedRoles {
		roles = append(roles, definitionRole{name, Roles[name]})
	}
	member, _ := s.ActiveMember()
	extension, err := s.extensionRoles(member)
	if err != nil {
		return nil, err
	}
	for _, r := range extension {
		roles = append(roles, definitionRole{extensionRoleName(r.Extension, r.ID), r})
	}
	return roles, nil
}

func agentInventoryPath(host string) string { return ".context-circuit/local/agents-" + host + ".yaml" }

// definitions is every definition setup would write for host, by path.
func (s *Store) definitions(host string) (map[string][]byte, []string, error) {
	roles, err := s.definitionRoles()
	if err != nil {
		return nil, nil, err
	}
	files := map[string][]byte{}
	var paths []string
	for _, each := range roles {
		setting, err := s.agentSetting(host, each.Name, each.Role)
		if err != nil {
			return nil, nil, err
		}
		path, content, err := roleFile(host, each.Name, each.Role, setting)
		if err != nil {
			return nil, nil, err
		}
		if _, seen := files[path]; seen {
			continue
		}
		files[path], paths = content, append(paths, path)
	}
	return files, paths, nil
}

func (s *Store) setupHost(host string) (AgentSetup, error) {
	var result AgentSetup
	if err := validHost(host); err != nil {
		return result, err
	}
	inventoryPath := agentInventoryPath(host)
	old := agentInventory{Files: map[string]string{}}
	if err := s.YAML(inventoryPath, &old); err != nil && !os.IsNotExist(err) {
		return result, err
	}
	files, paths, err := s.definitions(host)
	if err != nil {
		return result, err
	}
	for _, path := range paths {
		current, err := s.Read(path)
		if err != nil && !os.IsNotExist(err) {
			return result, err
		}
		if err == nil && digest(current) != digest(files[path]) && digest(current) != old.Files[path] {
			return result, fmt.Errorf("preserve customized agent definition: %s", path)
		}
	}
	// A definition the inventory records and setup no longer produces belongs
	// to an extension that is gone or off for this member. Setup removes only a
	// file it wrote and nobody changed since; an edited one is the person's.
	var remove []string
	for _, path := range slices.Sorted(maps.Keys(old.Files)) {
		if _, produced := files[path]; produced {
			continue
		}
		current, err := s.Read(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return result, err
		}
		if digest(current) == old.Files[path] {
			remove = append(remove, path)
			continue
		}
		result.Kept = append(result.Kept, KeptAgent{path, "edited since setup wrote it, so it is kept although no extension on for this member produces it; delete it by hand if it is no longer wanted"})
	}
	// The host setting is checked before anything is written, so a settings
	// file setup refuses to touch leaves the role definitions as they were too.
	setting, settingChanged, settingMode, err := s.hostSetting(host)
	if err != nil {
		return result, err
	}
	next := agentInventory{Files: map[string]string{}}
	for _, path := range paths {
		if err := s.Write(path, files[path], 0644); err != nil {
			return result, err
		}
		result.Written = append(result.Written, path)
		next.Files[path] = digest(files[path])
	}
	for _, path := range remove {
		resolved, err := s.Path(path)
		if err != nil {
			return result, err
		}
		if err := os.Remove(resolved); err != nil && !os.IsNotExist(err) {
			return result, err
		}
		result.Removed = append(result.Removed, path)
	}
	if err := s.WriteYAML(inventoryPath, next, 0600); err != nil {
		return result, err
	}
	if !settingChanged {
		return result, nil
	}
	path := hostSettings[host].Path
	result.Written = append(result.Written, path)
	return result, s.Write(path, setting, settingMode)
}

// AgentSetupRequired names the `agent setup` run that brings this machine's
// extension role definitions in line with the extensions on for its member, or
// is empty when they already match. Only extension definitions are compared:
// a shipped definition never changes with an extension, and a host that was
// never set up and has no extension role to write is not behind.
func (s *Store) AgentSetupRequired() string {
	for _, host := range Hosts {
		old := agentInventory{Files: map[string]string{}}
		if err := s.YAML(agentInventoryPath(host), &old); err != nil && !os.IsNotExist(err) {
			continue
		}
		files, _, err := s.definitions(host)
		if err != nil {
			continue
		}
		want, have := map[string]string{}, map[string]string{}
		for path, content := range files {
			if extensionDefinition(path) {
				want[path] = digest(content)
			}
		}
		for path, sum := range old.Files {
			if extensionDefinition(path) {
				have[path] = sum
			}
		}
		if !maps.Equal(want, have) {
			return "this machine's extension role definitions no longer match the extensions on for this member: run `" + s.agentSetupCommand("") + "`, then reload the host's agent definitions if it does not pick them up"
		}
	}
	return ""
}

// extensionDefinition reports whether a definition path is an extension role's.
func extensionDefinition(path string) bool {
	return strings.HasPrefix(filepath.Base(filepath.FromSlash(path)), "cc-x-")
}

// hostSettingFile is one top-level key setup adds to a host's own local settings
// file. A host that writes an attribution trailer into every commit it makes
// works against the rule that commits carry none, and a small model keeps the
// trailer its host wrote, so where the host has a project setting that turns it
// off, setup writes it. The file is the host's machine-local one, gitignored,
// because a person's settings there are theirs: setup adds the key when it is
// absent and otherwise leaves the file alone, a value the person chose included.
//
// Codex has no such setting. Cursor's is read only from its user configuration;
// its project file rejects the key and stops the CLI, so nothing is written.
type hostSettingFile struct {
	Path, Key string
	Value     json.RawMessage
}

var hostSettings = map[string]hostSettingFile{
	"claude-code": {".claude/settings.local.json", "attribution", json.RawMessage(`{"commit": "", "pr": ""}`)},
}

// hostSetting returns the merged settings file for a host and whether it
// differs from what is on disk. The person's bytes are kept: the key is
// inserted before the closing brace rather than re-encoding the file, which
// would reorder and reformat everything else in it.
func (s *Store) hostSetting(host string) ([]byte, bool, fs.FileMode, error) {
	want, ok := hostSettings[host]
	if !ok {
		return nil, false, 0, nil
	}
	current, err := s.Read(want.Path)
	if os.IsNotExist(err) {
		return []byte("{\n  " + strconv.Quote(want.Key) + ": " + string(want.Value) + "\n}\n"), true, 0644, nil
	}
	if err != nil {
		return nil, false, 0, err
	}
	resolved, err := s.Path(want.Path)
	if err != nil {
		return nil, false, 0, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return nil, false, 0, err
	}
	var top map[string]any
	decoder := json.NewDecoder(bytes.NewReader(current))
	decoder.UseNumber()
	if err := decoder.Decode(&top); err != nil || top == nil || decoder.Decode(new(any)) != io.EOF {
		return nil, false, 0, fmt.Errorf("preserve unreadable host settings: %s is not one JSON object", want.Path)
	}
	if _, present := top[want.Key]; present {
		return current, false, 0, nil
	}
	end := bytes.LastIndexByte(current, '}')
	body := bytes.TrimRight(current[:end], " \t\r\n")
	entry := "  " + strconv.Quote(want.Key) + ": " + string(want.Value) + "\n"
	separator := ",\n"
	if len(top) == 0 {
		separator = "\n"
	}
	merged := append(append(append([]byte{}, body...), separator...), entry...)
	merged = append(merged, current[end:]...)
	return merged, true, info.Mode().Perm(), nil
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

// DispatchOverride names the override a dispatch runs under, and the round of
// it, 0 when none is named.
type DispatchOverride struct {
	Name  string
	Round int
}

func (s *Store) DispatchAgent(host, role, task, directory, plan, intent string, shared, reviewRequested bool, override DispatchOverride, repository ...string) (Dispatch, error) {
	repo := ""
	if len(repository) > 0 {
		repo = repository[0]
	}
	if err := validHost(host); err != nil {
		return Dispatch{}, err
	}
	r, err := s.resolveRole(role)
	if err != nil {
		return Dispatch{}, err
	}
	if repo != "" && !r.edits(role) {
		return Dispatch{}, errors.New("--repo assigns a worker, or an extension role with edit access")
	}
	setting, err := s.agentSetting(host, role, r)
	if err != nil {
		return Dispatch{}, err
	}
	var unlocked *HookOverride
	if override.Name == "" && override.Round != 0 {
		return Dispatch{}, errors.New("--round goes with --override")
	}
	if override.Name != "" {
		if plan == "" {
			return Dispatch{}, fmt.Errorf("--override %s runs a step of a plan's execution: pass --plan", override.Name)
		}
		record, err := s.FindRecord(plan)
		if err != nil {
			return Dispatch{}, err
		}
		repositories := record.Repositories
		if repo != "" {
			repositories = []string{repo}
		}
		resolved, err := s.dispatchOverride(override.Name, role, repositories)
		if err != nil {
			return Dispatch{}, err
		}
		unlocked = &resolved
	}
	if err := checkDispatchArgs(role, r, task, plan, intent, repo, shared, reviewRequested, unlocked, override.Round); err != nil {
		return Dispatch{}, err
	}
	work, err := Open(directory)
	if err != nil {
		return Dispatch{}, err
	}
	name := agentType(host, role)
	// The brief is read top to bottom by an agent that has seen nothing else, so
	// it states where the work happens, then what the work is, and only then
	// what to return. Every record it names is also quoted in full: an agent
	// that has to go and find a file reads the whole repository on the way.
	prompt := "# Working directory\n\n" + work.Root + "\n\nWorkspace root: " + s.Root
	ownership := ownershipBrief(r, shared)
	prompt += "\n\n# Ownership\n\n" + ownership
	// Whoever wrote the record wrote it in their own language, so the brief
	// reports that one rather than the language of whoever is dispatching.
	author := ""
	var planRepositories []string
	if plan != "" {
		record, err := s.FindRecord(plan)
		if err != nil {
			return Dispatch{}, err
		}
		if !strings.HasPrefix(record.ID, "p") {
			return Dispatch{}, errors.New("--plan needs a plan ID")
		}
		author = record.CreatedBy
		planRepositories = record.Repositories
		prompt += "\n\n" + planBrief(record, s.Root)
		if r.edits(role) {
			assignment, err := s.workerAssignment(record, repo)
			if err != nil {
				return Dispatch{}, err
			}
			prompt += assignment
		}
	}
	if intent != "" {
		record, err := s.approvedIntent(intent)
		if err != nil {
			return Dispatch{}, err
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
	if unlocked != nil {
		prompt += "\n\n" + s.overrideBrief(*unlocked, override.Round)
	}
	if repo != "" {
		planRepositories = []string{repo}
	}
	member, _ := s.ActiveMember()
	additions := s.roleAdditions(member, role, planRepositories, role == "planner")
	if section := s.extensionBrief(role, r, planRepositories, additions); section != "" {
		prompt += "\n\n" + section
	}
	prompt += "\n\n" + languageBrief(s.authorLanguage(author))
	prompt += "\n\n# What to return\n\n" + returnSection(r, additions)
	definition, _, err := roleFile(host, role, r, setting)
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
		setup = s.agentSetupCommand(host)
	}
	return Dispatch{host, role, name, setting, work.Root, r.ReadOnly, prompt, true, definition, installed, setup}, nil
}

// agentSetupCommand names the call that writes a host's role definitions. A
// workspace that carries a launcher gets the launcher, which runs the version it
// pins whichever CLI happens to be running; one without names the running CLI's
// own command.
func (s *Store) agentSetupCommand(host string) string {
	// An empty host is every host, which is what an extension change needs.
	narrow := ""
	if host != "" {
		narrow = " --host " + host
	}
	if launcher, err := s.Path(".context-circuit/cli"); err == nil {
		if info, err := os.Stat(launcher); err == nil && info.Mode().IsRegular() {
			if runtime.GOOS == "windows" {
				return "powershell -NoProfile -ExecutionPolicy Bypass -File " + launcher + ".ps1 agent setup" + narrow
			}
			return "sh " + launcher + " agent setup" + narrow
		}
	}
	command := s.Command
	if command == "" {
		command = "context-circuit-cli"
	}
	return command + " --workspace " + s.Root + " agent setup" + narrow
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
	for i := range len(text) {
		if text[i] == char {
			run++
			longest = max(longest, run)
		} else {
			run = 0
		}
	}
	return longest
}

// checkDispatchArgs rejects flag combinations that do not describe one agent's
// assignment.
//
// An allowed override unlocks exactly the dispatches its default needs: the
// verifier during execution, a worker carrying the verifier's findings, or a
// reviewer nobody asked for. Each is otherwise refused, so the check grows one
// condition rather than a second path.
func checkDispatchArgs(role string, r Role, task, plan, intent, repo string, shared, reviewRequested bool, override *HookOverride, round int) error {
	if override != nil {
		if err := checkOverrideUse(role, r, *override, task, plan, round); err != nil {
			return err
		}
		// The override stands in for the request a reviewer otherwise needs.
		reviewRequested = true
	}
	// An extension role is held to the shipped role its access makes it: one
	// that edits is assigned like the worker, through a plan, and a read-only
	// one like the explorer, through a task or a plan to read.
	if r.Extension != "" {
		edits := r.edits(role)
		if intent != "" {
			return errors.New("--intent dispatches a planner; an extension role takes --plan or --task")
		}
		if task != "" {
			if err := Text(task); err != nil {
				return err
			}
		}
		if shared && !edits {
			return errors.New("--shared describes roles sharing one worktree, and " + role + " is read-only")
		}
		if edits && plan == "" {
			return errors.New(role + " has edit access, so it is dispatched against the plan whose worktree it works in: pass --plan")
		}
		if plan == "" && task == "" {
			return errors.New(role + " needs --plan, or --task naming its assignment")
		}
		return nil
	}
	// The record is the assignment. A planner is dispatched against a whole
	// approved intent and a worker against the plan quoted in its brief, so a
	// task on either can only restate what the brief already carries, and a
	// required flag leaves the coordinator nothing to write but that
	// restatement. A worker splitting one plan with others is the exception:
	// its slice is named in the task and nowhere else.
	if role == "planner" && intent != "" && task != "" {
		return errors.New("an intent planner takes the whole approved intent, so it takes no --task")
	}
	briefed := role == "planner" || role == "worker"
	if task != "" || !briefed {
		if err := Text(task); err != nil {
			return err
		}
	}
	if role == "reviewer" && !reviewRequested {
		return errors.New("independent review requires an explicit user request; use --review-requested only to represent that request, or --override with an override this workspace allowed")
	}
	if shared && role != "worker" {
		return errors.New("--shared describes workers sharing one worktree")
	}
	// A planner decides the plan shape, so a plan cannot already exist to hand
	// it: numbering one first settles the split the planner was dispatched to
	// propose, and an intent holding several plans has no single ID to pass.
	if role == "planner" {
		if plan != "" {
			return errors.New("a planner proposes the plan shape, so it is dispatched against --intent, not an already numbered --plan")
		}
		if intent == "" && task == "" {
			return errors.New("a planner needs --intent or --task specifying the standalone outcome")
		}
	} else if intent != "" {
		return errors.New("--intent dispatches a planner; other roles take --plan")
	}
	// The brief tells a worker to implement the plan above it, so a worker with
	// neither a plan nor a task would be launched with nothing to implement.
	if role == "worker" && plan == "" && task == "" {
		return errors.New("a worker needs --plan, or --task naming its assignment")
	}
	if repo != "" && plan == "" {
		return errors.New("--repo picks one repository of the worker's --plan; pass --plan")
	}
	return nil
}

// ownershipBrief says what a role may do in its working directory.
func ownershipBrief(r Role, shared bool) string {
	// A read-only role was told where its working directory is and nothing
	// about what it may do there, which leaves how to read it an open question
	// on a host that also offers GUI automation and web search. One planner
	// answered it by driving an editor's interface to read a local file and
	// searching the public web for a repository sitting on disk.
	if r.ReadOnly {
		return "You are reading, not changing. Change nothing in this working directory. Read and search the files here directly rather than driving another application, and consult no external source: what this task needs is this directory and this brief."
	}
	// Isolation differs by parallelism axis. One worker per plan owns its
	// whole worktree; several workers inside one worktree genuinely share
	// files. Saying the wrong one invites a worker to guess at edits it
	// cannot see, or to overwrite edits it can.
	if shared {
		return "Other workers are editing this same working directory. Preserve their edits, adapt your changes, and stay within your assigned paths."
	}
	return "You are the sole owner of this working directory. Work only here. Do not create, switch, merge, push, or delete branches, and do not run git worktree."
}

// workerAssignment names the repository a worker owns and the files it must
// read before editing.
func (s *Store) workerAssignment(record Record, repo string) (string, error) {
	if repo == "" {
		if len(record.Repositories) != 1 {
			return "", errors.New("multi-repository worker dispatch needs --repo")
		}
		repo = record.Repositories[0]
	}
	files, err := s.requiredPlanFiles(record, repo)
	if err != nil {
		return "", err
	}
	brief := "\n\nRepository assignment: " + repo
	if len(files) > 0 {
		brief += "\n\n# Required files relative to workspace root\n\n- " + strings.Join(files, "\n- ") + "\n"
		brief += "\nResolve these paths against the workspace root above and read every required file before editing."
	}
	return brief + "\n\nWork only in your assigned repository worktree. Links in the plan body to other repositories' files are for human review; do not open them unless assigned here.", nil
}

// approvedIntent finds an intent a planner may be dispatched against.
func (s *Store) approvedIntent(intent string) (Record, error) {
	record, err := s.FindRecord(intent)
	if err != nil {
		return Record{}, err
	}
	if !strings.HasPrefix(record.ID, "i") {
		return Record{}, errors.New("--intent needs an intent ID")
	}
	// Planning an unapproved intent plans an outcome nobody agreed to.
	if record.ApprovedAt == "" {
		return Record{}, fmt.Errorf("intent %s is not approved; approve it before planning it", record.ID)
	}
	return record, nil
}
