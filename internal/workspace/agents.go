package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
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
	"explorer": {"Investigate a bounded codebase question and return evidence.", "Inspect only the assigned repository question. Return file locations, observed behavior, and uncertainties. You cannot run anything, so report behavior you read, never behavior you confirmed. Do not edit files or launch other agents.", true},
	"planner":  {"Plan an approved intent using evidence from the repositories.", "Read the approved intent and relevant code. Answer under these headings: Verdict, one of feasible, feasible-with-changes, or not-feasible; Approach; Tasks and order, each naming its repository, paths, and the tasks it depends on; Risks and checks, naming check commands and where they are defined; Evidence, the paths that established each conclusion; Uncertainties, what could not be determined and what would settle it; Plan shape, single or a proposed split with reasons. You cannot run anything, so report a check as found, never as passing. Not-feasible is a complete answer: return it with its evidence instead of a plan. The coordinator writes the plan and decides any split. Do not implement or request a separate plan approval. Do not edit files or launch other agents.", true},
	"worker":   {"Implement a bounded part of an approved plan and run normal checks.", "Implement the assigned approved plan in the supplied working directory. Own only the assigned task and paths. Do not write or reconcile durable project knowledge; report what the coordinator should record. Run normal tests and report changes, results, and remaining work. Do not independently verify your implementation, launch agents, or deliver changes without explicit authorization.", false},
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

func (s *Store) DispatchAgent(host, role, task, directory, plan string, shared, reviewRequested bool) (Dispatch, error) {
	setting, err := s.agentSetting(host, role)
	if err != nil {
		return Dispatch{}, err
	}
	if err := Text(task); err != nil {
		return Dispatch{}, err
	}
	if role == "reviewer" && !reviewRequested {
		return Dispatch{}, fmt.Errorf("independent review requires an explicit user request; use --review-requested only to represent that request")
	}
	if shared && role != "worker" {
		return Dispatch{}, errors.New("--shared describes workers sharing one worktree")
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
	prompt := r.Instructions + "\n\nWorkspace: " + s.Root + "\nWorking directory: " + work.Root
	if !r.ReadOnly {
		// Isolation differs by parallelism axis. One worker per plan owns its
		// whole worktree; several workers inside one worktree genuinely share
		// files. Saying the wrong one invites a worker to guess at edits it
		// cannot see, or to overwrite edits it can.
		if shared {
			prompt += "\n\nOwnership: other workers are editing this same working directory. Preserve their edits, adapt your changes, and stay within your assigned paths."
		} else {
			prompt += "\n\nOwnership: you are the sole owner of this working directory. Work only here. Do not create, switch, merge, push, or delete branches, and do not run git worktree."
		}
	}
	if plan != "" {
		record, err := s.FindRecord(plan)
		if err != nil {
			return Dispatch{}, err
		}
		if !strings.HasPrefix(record.ID, "p") {
			return Dispatch{}, errors.New("--plan needs a plan ID")
		}
		prompt += "\n\n" + planBrief(record)
	}
	prompt += "\n\nTask: " + task
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

// planBrief states the facts the CLI can verify. A dependent plan's branch is
// prepared from its predecessors, so their work is already in its ancestry;
// saying so prevents a worker from reimplementing what it inherited.
func planBrief(record Record) string {
	brief := "Plan " + record.ID
	if record.Intent != "" {
		brief += " of intent " + record.Intent
	}
	brief += "\nRepositories: " + strings.Join(record.Repositories, ", ")
	if len(record.DependsOn) > 0 {
		brief += "\nDepends on: " + strings.Join(record.DependsOn, ", ") +
			"\nTheir completed work is expected in this branch's ancestry. Build on it rather than reimplementing it." +
			"\nAny plan running beside you is invisible in a separate worktree and is reconciled later by the coordinator. Do not guess at its changes; report an interface you assume it may also be changing."
	}
	_, body, err := splitRecord([]byte(record.Content))
	if err == nil && len(strings.TrimSpace(string(body))) > 0 {
		brief += "\n\nPlan record (authoritative):\n" + strings.TrimSpace(string(body))
	}
	return brief
}
