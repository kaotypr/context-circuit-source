package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
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
	"explorer": {"Investigate a bounded codebase question and return evidence.", "Inspect only the assigned repository question. Return file locations, observed behavior, and uncertainties. Do not edit files or launch other agents.", true},
	"planner":  {"Plan an approved intent using evidence from the repositories.", "Read the approved intent and relevant code. Return task order, cross-repository dependencies, risks, and expected checks. The coordinator writes the plan. Do not implement or request a separate plan approval. Do not edit files or launch other agents.", true},
	"worker":   {"Implement a bounded part of an approved plan and run normal checks.", "Implement the assigned approved plan in the supplied working directory. You are not alone in the codebase: preserve others' edits and adapt your changes. Own only the assigned task and paths. Run normal tests and report changes, results, and remaining work. Do not independently verify your implementation, launch agents, or deliver changes without explicit authorization.", false},
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

func (s *Store) AgentSettings() (RoleTiering, error) {
	var cfg RoleTiering
	if err := s.YAML(".context-circuit/role-tiering.yaml", &cfg); err != nil {
		return cfg, err
	}
	for host, settings := range cfg.Hosts {
		if err := validHost(host); err != nil {
			return cfg, err
		}
		for name, setting := range settings {
			if _, ok := Roles[name]; !ok {
				return cfg, fmt.Errorf("unknown agent role: %s", name)
			}
			if err := validateSetting(host, setting); err != nil {
				return cfg, fmt.Errorf("%s/%s: %w", host, name, err)
			}
		}
	}
	return cfg, nil
}
func (s *Store) ConfigureAgent(host, role, model, effort string) error {
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
	cfg, err := s.AgentSettings()
	if err != nil {
		return err
	}
	if cfg.Hosts[host] == nil {
		return fmt.Errorf("host is absent from role-tiering.yaml: %s", host)
	}
	return s.Update(".context-circuit/role-tiering.yaml", []string{"hosts", host, role}, setting, 0644)
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
func (s *Store) SetupAgents(host string) ([]string, error) {
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
}

func (s *Store) DispatchAgent(host, role, task, directory string, reviewRequested bool) (Dispatch, error) {
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
	work, err := Open(directory)
	if err != nil {
		return Dispatch{}, err
	}
	r := Roles[role]
	name := "cc-" + role
	if host == "codex" {
		name = "cc_" + role
	}
	prompt := r.Instructions + "\n\nWorkspace: " + s.Root + "\nWorking directory: " + work.Root + "\nTask: " + task
	return Dispatch{host, role, name, setting, work.Root, r.ReadOnly, prompt, true}, nil
}
