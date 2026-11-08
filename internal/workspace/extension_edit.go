package workspace

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

// extensionsHeader opens an extensions.yaml the CLI creates.
const extensionsHeader = "# The extensions this workspace uses, in the order their hooks run at each\n" +
	"# point. `extension add` writes an entry; settings are edited here by hand and\n" +
	"# validated by `check`. Never put a credential in this file.\n"

// AddExtension validates an extension and appends its entry. A remote source is
// resolved to the commit its ref names and fetched, so the manifest is read
// before anything is written.
func (s *Store) AddExtension(ctx context.Context, id, path, url, ref, subpath string, repositories []string) (ExtensionEntry, error) {
	if err := Name(id); err != nil {
		return ExtensionEntry{}, err
	}
	file, _, err := s.Extensions()
	if err != nil {
		return ExtensionEntry{}, err
	}
	if file.extensionIndex(id) >= 0 {
		return ExtensionEntry{}, fmt.Errorf("extensions.yaml already lists %s; change its source with `extension update`, or remove it first", id)
	}
	entry := ExtensionEntry{ID: id, Repositories: unique(repositories)}
	switch {
	case path != "" && url == "":
		if ref != "" || subpath != "" {
			return ExtensionEntry{}, errors.New("--ref and --subpath belong to a --url source")
		}
		entry.Source = ExtensionSource{Path: path}
	case url != "" && path == "":
		entry.Source = ExtensionSource{URL: url, Ref: ref, Path: subpath}
	default:
		return ExtensionEntry{}, errors.New("pass --path for a folder in the workspace, or --url with --ref for a Git repository")
	}
	if err := checkSource(entry.Source); err != nil {
		return ExtensionEntry{}, err
	}
	if err := s.knownRepositories(entry.Repositories); err != nil {
		return ExtensionEntry{}, err
	}
	if entry.Source.Remote() {
		commit, err := s.resolveExtensionRef(ctx, url, ref)
		if err != nil {
			return ExtensionEntry{}, err
		}
		entry.Source.Commit = commit
		if _, err := s.fetchExtension(ctx, entry); err != nil {
			return ExtensionEntry{}, err
		}
	}
	if err := s.usableExtension(entry); err != nil {
		return ExtensionEntry{}, err
	}
	if _, err := s.Read(ExtensionsPath); os.IsNotExist(err) {
		encoded, err := yaml.Marshal(ExtensionsFile{Extensions: []ExtensionEntry{entry}})
		if err != nil {
			return ExtensionEntry{}, err
		}
		return entry, s.Write(ExtensionsPath, append([]byte(extensionsHeader), encoded...), 0644)
	}
	return entry, s.AppendItem(ExtensionsPath, []string{"extensions"}, entry, 0644)
}

// usableExtension refuses an entry whose manifest this CLI cannot use.
func (s *Store) usableExtension(entry ExtensionEntry) error {
	loaded := s.loadExtension(entry, 0)
	if loaded.Fault != extensionUsable {
		return fmt.Errorf("%s cannot be used: %s", entry.ID, loaded.Problem)
	}
	if loaded.Manifest.ID != entry.ID {
		return fmt.Errorf("the manifest at %s is extension %s, not %s; add it under its own id", loaded.Folder, loaded.Manifest.ID, entry.ID)
	}
	return nil
}

func (s *Store) knownRepositories(repositories []string) error {
	if len(repositories) == 0 {
		return nil
	}
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	for _, repo := range repositories {
		if _, ok := cfg.Repositories[repo]; !ok {
			return fmt.Errorf("unknown repository: %s", repo)
		}
	}
	return nil
}

// extensionEntry finds an entry by ID or says it is not listed.
func (s *Store) extensionEntry(id string) (ExtensionsFile, int, error) {
	file, present, err := s.Extensions()
	if err != nil {
		return file, -1, err
	}
	index := file.extensionIndex(id)
	if !present || index < 0 {
		return file, -1, fmt.Errorf("extensions.yaml lists no extension %s", id)
	}
	return file, index, nil
}

// UpdateExtension moves a remote extension to a new ref and records the commit
// it names now.
func (s *Store) UpdateExtension(ctx context.Context, id, ref string) (ExtensionEntry, error) {
	file, index, err := s.extensionEntry(id)
	if err != nil {
		return ExtensionEntry{}, err
	}
	entry := file.Extensions[index]
	if !entry.Source.Remote() {
		return ExtensionEntry{}, fmt.Errorf("%s is a folder in the workspace; it changes through a workspace commit, not a ref", id)
	}
	commit, err := s.resolveExtensionRef(ctx, entry.Source.URL, ref)
	if err != nil {
		return ExtensionEntry{}, err
	}
	entry.Source.Ref, entry.Source.Commit = ref, commit
	if _, err := s.fetchExtension(ctx, entry); err != nil {
		return ExtensionEntry{}, err
	}
	if err := s.usableExtension(entry); err != nil {
		return ExtensionEntry{}, err
	}
	keys := []string{"extensions", ListIndex(index), "source"}
	if err := s.Update(ExtensionsPath, append(keys, "ref"), ref, 0644); err != nil {
		return ExtensionEntry{}, err
	}
	return entry, s.Update(ExtensionsPath, append(keys, "commit"), commit, 0644)
}

// RemoveExtension drops an entry. Its cached commits stay until a prune.
func (s *Store) RemoveExtension(id string) error {
	_, index, err := s.extensionEntry(id)
	if err != nil {
		return err
	}
	return s.RemoveItem(ExtensionsPath, []string{"extensions"}, index, 0644)
}

// SetExtensionSwitch writes one member's on or off for an extension. A required
// extension cannot be switched off, so that is refused rather than written.
func (s *Store) SetExtensionSwitch(id, member string, on bool) error {
	file, index, err := s.extensionEntry(id)
	if err != nil {
		return err
	}
	members, err := s.Members()
	if err != nil {
		return err
	}
	if member == "" {
		if member, err = s.activeMember(members); err != nil {
			return err
		}
	}
	if _, ok := members.Members[member]; !ok {
		return fmt.Errorf("member %s is not on the roster", member)
	}
	value := ExtensionSwitch("off")
	if on {
		value = "on"
	}
	if !on && s.loadExtension(file.Extensions[index], index).required() {
		return fmt.Errorf("%s is required by its manifest, so it cannot be switched off; whether it should stay required is a question for its author", id)
	}
	keys := []string{"extensions", ListIndex(index), "members"}
	if len(file.Extensions[index].Members) == 0 {
		return s.Update(ExtensionsPath, keys, map[string]string{member: string(value)}, 0644)
	}
	return s.Update(ExtensionsPath, append(keys, member), string(value), 0644)
}

// ExtensionState is one extension as this member and machine see it.
type ExtensionState struct {
	ID           string          `json:"id"`
	Name         string          `json:"name,omitempty"`
	Version      string          `json:"version,omitempty"`
	Summary      string          `json:"summary,omitempty"`
	Source       ExtensionSource `json:"source"`
	Repositories []string        `json:"repositories,omitempty"`
	Members      string          `json:"members,omitempty"`
	// State is on, off, unavailable, or invalid; Reason says why when not on.
	State  string   `json:"state"`
	Reason string   `json:"reason,omitempty"`
	Cached *bool    `json:"cached,omitempty"`
	Hooks  []string `json:"hooks,omitempty"`
	// Roles names each agent role the extension declares, by the name it is
	// dispatched as, and Adds each shipped role it adds to.
	Roles []string `json:"roles,omitempty"`
	Adds  []string `json:"role_additions,omitempty"`
	// Overrides names each overridable default the manifest declares; status
	// says which a person allowed.
	Overrides []string `json:"overrides,omitempty"`
	// Skills names each skill the extension ships, by its installed name, and
	// Stages each named stage they declare, as other extensions hook it.
	Skills []string `json:"skills,omitempty"`
	Stages []string `json:"stages,omitempty"`
	// Scripts names the points that run a script, and Trust whether this
	// member allowed the extension's current version on this machine: allowed,
	// not allowed, or changed since allowed.
	Scripts []string `json:"scripts,omitempty"`
	Trust   string   `json:"trust,omitempty"`
	// Allow is the command that clears a trust other than allowed, which the
	// person decides to run.
	Allow string `json:"allow,omitempty"`
}

// ExtensionStates lists every extension for the active member on this machine.
func (s *Store) ExtensionStates() ([]ExtensionState, error) {
	file, _, err := s.Extensions()
	if err != nil {
		return nil, err
	}
	member, _ := s.ActiveMember()
	states := []ExtensionState{}
	for index, entry := range file.Extensions {
		loaded := s.loadExtension(entry, index)
		state := ExtensionState{ID: entry.ID, Source: entry.Source, Repositories: entry.Repositories, Members: loaded.membersDefault()}
		if entry.Source.Remote() {
			cached := loaded.Fault != extensionUncached && loaded.Fault != extensionInvalid
			state.Cached = &cached
		}
		switch loaded.Fault {
		case extensionUsable:
			state.Name, state.Version, state.Summary = loaded.Manifest.Name, loaded.Manifest.Version, loaded.Manifest.Summary
			on, why := loaded.memberState(member)
			state.State, state.Reason = "off", why
			if on {
				state.State = "on"
			}
			for _, hook := range loaded.Manifest.Hooks {
				point, _ := ParseHookPoint(hook.At)
				state.Hooks = append(state.Hooks, point.String())
				if hook.Script != nil {
					state.Scripts = append(state.Scripts, point.String())
				}
			}
			for _, declared := range loaded.Manifest.Roles {
				state.Roles = append(state.Roles, extensionRoleName(entry.ID, declared.ID))
			}
			for _, addition := range loaded.Manifest.RoleAdditions {
				if !slices.Contains(state.Adds, addition.Role) {
					state.Adds = append(state.Adds, addition.Role)
				}
			}
			for _, o := range loaded.Manifest.Overrides {
				state.Overrides = append(state.Overrides, o.Default)
			}
			for _, skill := range loaded.Manifest.Skills {
				state.Skills = append(state.Skills, extensionSkillName(entry.ID, skill.ID))
				if skill.Stage != "" {
					state.Stages = append(state.Stages, entry.ID+"/"+skill.Stage)
				}
			}
			if len(state.Scripts) > 0 {
				version, err := s.extensionVersion(loaded)
				if err == nil {
					state.Trust, err = s.trustState(member, entry.ID, version)
				}
				if err != nil {
					return nil, err
				}
				if state.Trust != TrustAllowed {
					state.Allow = "extension allow --id " + entry.ID
				}
			}
		case extensionUncached, extensionUnavailable:
			state.State, state.Reason = "unavailable", loaded.Problem+"; "+extensionRemedy(loaded)
		default:
			state.State, state.Reason = "invalid", loaded.Problem+"; "+extensionRemedy(loaded)
		}
		states = append(states, state)
	}
	return states, nil
}

// ExtensionSummary is what status says about extensions: how many are listed,
// and how many are on for the active member, so an agent knows whether the
// extension rule applies without reading extensions.yaml.
type ExtensionSummary struct {
	Listed      int `json:"listed"`
	On          int `json:"on"`
	Unavailable int `json:"unavailable,omitempty"`
}

func (s *Store) extensionSummary(member string) (*ExtensionSummary, error) {
	file, present, err := s.Extensions()
	if err != nil || !present {
		return nil, err
	}
	summary := &ExtensionSummary{Listed: len(file.Extensions)}
	for index, entry := range file.Extensions {
		loaded := s.loadExtension(entry, index)
		if loaded.Fault != extensionUsable {
			summary.Unavailable++
			continue
		}
		if on, _ := loaded.memberState(member); on {
			summary.On++
		}
	}
	return summary, nil
}

// extensionIssues reports what is wrong in extensions.yaml and the manifests
// it names, each with what clears it. An unreadable file is reported by
// status, which check includes. unchecked names each remote this machine has
// not fetched, which check never contacts the network to read.
func (s *Store) extensionIssues() (issues, unchecked []Finding) {
	file, present, err := s.Extensions()
	if err != nil || !present {
		return nil, nil
	}
	cfg, cfgErr := s.Config()
	members, membersErr := s.Members()
	seen := map[string]int{}
	for index, entry := range file.Extensions {
		label := fmt.Sprintf("extensions.yaml entry %d (%s)", index+1, entry.ID)
		seen[entry.ID]++
		if seen[entry.ID] == 2 {
			issues = append(issues, found("extensions.yaml lists "+entry.ID+" more than once",
				needsAPerson+"decide which entry stays, and delete the other from extensions.yaml"))
		}
		if cfgErr == nil {
			for _, repo := range entry.Repositories {
				if _, ok := cfg.Repositories[repo]; !ok {
					issues = append(issues, found(label+" names repository "+repo+", which workspace.yaml does not list",
						"edit its repositories in extensions.yaml, or connect the repository with `repo connect`"))
				}
			}
		}
		if membersErr == nil {
			for _, member := range slices.Sorted(maps.Keys(entry.Members)) {
				if _, ok := members.Members[member]; !ok {
					issues = append(issues, found(label+" switches member "+member+", who is not on the roster",
						"remove the switch from extensions.yaml, or add the member with `member add`"))
				}
			}
		}
		if entry.Source.Remote() && entry.Source.Commit == "" {
			issues = append(issues, found(label+" records no resolved commit",
				fmt.Sprintf("`extension update --id %s --ref %s`", entry.ID, orPlaceholder(entry.Source.Ref, "REF"))))
			continue
		}
		loaded := s.loadExtension(entry, index)
		switch loaded.Fault {
		case extensionUsable:
		case extensionUncached:
			unchecked = append(unchecked, found(label+": "+loaded.Problem+", so its manifest was not checked",
				"`extension sync`"))
			continue
		case extensionNewerCLI:
			issues = append(issues, found(label+": "+loaded.Problem, "upgrade the CLI through `.agents/skills/cc-cli/SKILL.md`"))
			continue
		case extensionEscapes:
			issues = append(issues, found(label+": "+loaded.Problem,
				needsAPerson+"the extension's author must keep every path inside its folder; until then remove it with `extension remove --id "+entry.ID+"`"))
			continue
		default:
			issues = append(issues, found(label+": "+loaded.Problem,
				"fix the source, or remove it with `extension remove --id "+entry.ID+"`"))
			continue
		}
		manifest := loaded.Manifest
		if manifest.ID != entry.ID {
			issues = append(issues, found(fmt.Sprintf("%s: its manifest is extension %s", label, manifest.ID),
				"correct the entry's id in extensions.yaml to "+manifest.ID+", or point the entry at the extension it meant"))
		}
		if manifest.Members == ExtensionRequired && membersErr == nil {
			for _, member := range slices.Sorted(maps.Keys(entry.Members)) {
				if entry.Members[member] == "off" {
					issues = append(issues, found(label+" switches off "+member+", but its manifest makes it required, so the switch has no effect",
						"remove the switch from extensions.yaml; whether it should stay required is a question for the extension's author"))
				}
			}
		}
		for _, key := range slices.Sorted(maps.Keys(manifest.Settings)) {
			if _, ok := entry.Settings[key]; manifest.Settings[key].Required && !ok {
				issues = append(issues, found(label+" is missing the required setting "+key,
					"add settings."+key+" to the entry in extensions.yaml"))
			}
		}
		for _, key := range slices.Sorted(maps.Keys(entry.Settings)) {
			if _, ok := manifest.Settings[key]; !ok {
				issues = append(issues, found(label+" sets "+key+", which its manifest does not declare",
					"remove settings."+key+" from the entry in extensions.yaml, or correct its name"))
			}
		}
	}
	issues = append(issues, s.extensionRoleIssues(file)...)
	issues = append(issues, s.extensionSkillIssues(file)...)
	issues = append(issues, s.overrideIssues(file)...)
	return issues, unchecked
}

// UncheckedExtensions names the remote extensions check could not read on this
// machine. It never fails: an unreadable file is already a finding.
func (s *Store) UncheckedExtensions() []Finding {
	_, unchecked := s.extensionIssues()
	return unchecked
}

func orPlaceholder(value, placeholder string) string {
	if strings.TrimSpace(value) == "" {
		return placeholder
	}
	return value
}
