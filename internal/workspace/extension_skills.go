package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

// An extension may ship skills. `agent setup` installs each skill of each
// extension on for this machine's member as x-<extension>-<skill>: a copy of
// its folder under .agents/skills/, which is what Codex reads, with the name in
// its frontmatter set to the installed name, and a one-line wrapper under
// .claude/skills/ and .cursor/skills/ pointing to it, in the shipped wrapper
// form. No shipped skill starts with x-, so an extension skill never replaces
// or shadows one, and no template version ships an x- path, so a template
// update never writes or removes one. All of them are gitignored: the source
// of truth is the extension's folder or pinned commit.
//
// A skill may declare a named stage. Its before and after points are hooked
// by other extensions as <extension>/<stage>, and the skill's own text runs
// them through `extension hooks`. A named stage never enters the core
// lifecycle: record order, AGENTS.md, and the core skills know nothing of it.

// ManifestSkill is a skill an extension ships: a folder inside the extension
// holding SKILL.md and whatever it references, and optionally the named stage
// it declares.
type ManifestSkill struct {
	ID    string `yaml:"id" json:"id"`
	Path  string `yaml:"path" json:"path"`
	Stage string `yaml:"stage,omitempty" json:"stage,omitempty"`
}

// skillInventoryPath records the skill folders setup copied under
// .agents/skills/, shared by every host, the way agents-<host>.yaml records a
// host's definitions.
const skillInventoryPath = ".context-circuit/local/agents-skills.yaml"

// skillFileLimit bounds how many files one skill folder may hold.
const skillFileLimit = 500

// extensionSkillName is an extension skill's installed name.
func extensionSkillName(extension, id string) string { return "x-" + extension + "-" + id }

// validateManifestSkills checks the skills a manifest declares. Their folders
// are checked by loadExtension, which has the extension folder.
func validateManifestSkills(m Manifest) error {
	ids, stages := map[string]bool{}, map[string]string{}
	for _, skill := range m.Skills {
		if err := Name(skill.ID); err != nil {
			return fmt.Errorf("skill id: %w", err)
		}
		if ids[skill.ID] {
			return fmt.Errorf("skill %s appears twice; each skill id appears once", skill.ID)
		}
		ids[skill.ID] = true
		if skill.Path == "" {
			return fmt.Errorf("skill %s names no path; a skill is a folder holding SKILL.md", skill.ID)
		}
		if err := insidePath(skill.Path); err != nil {
			return fmt.Errorf("skill %s path: %w: %w", skill.ID, err, errEscapes)
		}
		if skill.Stage == "" {
			continue
		}
		if err := Name(skill.Stage); err != nil {
			return fmt.Errorf("skill %s stage: %w", skill.ID, err)
		}
		if other, ok := stages[skill.Stage]; ok {
			return fmt.Errorf("skills %s and %s both declare stage %s; a stage is declared once", other, skill.ID, skill.Stage)
		}
		stages[skill.Stage] = skill.ID
	}
	return nil
}

// skillFrontmatter reads SKILL.md's frontmatter, requiring name and
// description, and returns where the frontmatter text ends.
func skillFrontmatter(data []byte) (name, description string, end int, err error) {
	text := string(data)
	rest, ok := strings.CutPrefix(text, "---\n")
	if !ok {
		rest, ok = strings.CutPrefix(text, "---\r\n")
	}
	if !ok {
		return "", "", 0, errors.New("SKILL.md does not open with --- frontmatter")
	}
	start := len(text) - len(rest)
	closing := regexp.MustCompile(`(?m)^---\r?$`).FindStringIndex(rest)
	if closing == nil {
		return "", "", 0, errors.New("SKILL.md frontmatter is not closed")
	}
	var fields map[string]any
	if err := yaml.Unmarshal([]byte(rest[:closing[0]]), &fields); err != nil {
		return "", "", 0, fmt.Errorf("SKILL.md frontmatter: %w", err)
	}
	name, _ = fields["name"].(string)
	description, _ = fields["description"].(string)
	if strings.TrimSpace(name) == "" || strings.TrimSpace(description) == "" {
		return "", "", 0, errors.New("SKILL.md frontmatter needs a name and a description")
	}
	return name, strings.TrimSpace(description), start + closing[0], nil
}

// nameLine is a top-level name field in frontmatter.
var nameLine = regexp.MustCompile(`(?m)^name:[^\r\n]*`)

// renameSkill sets the name in SKILL.md's frontmatter to the installed one,
// leaving every other byte as the author wrote it.
func renameSkill(data []byte, installed string) ([]byte, error) {
	_, _, end, err := skillFrontmatter(data)
	if err != nil {
		return nil, err
	}
	front := data[:end]
	loc := nameLine.FindIndex(front)
	if loc == nil {
		return nil, errors.New("SKILL.md frontmatter names the skill in a form setup cannot rewrite; write it as a top-level `name:` line")
	}
	out := append([]byte{}, front[:loc[0]]...)
	out = append(out, "name: "+installed...)
	out = append(out, front[loc[1]:]...)
	return append(out, data[end:]...), nil
}

// skillFile is one file of a skill folder, relative to that folder.
type skillFile struct {
	Path string
	Data []byte
	Mode fs.FileMode
}

// readSkillFolder reads every regular file of a skill folder inside an
// extension folder, refusing a path or link that leaves the extension.
func readSkillFolder(extensionRoot, path string) ([]skillFile, error) {
	if err := insidePath(path); err != nil {
		return nil, fmt.Errorf("%w: %w", err, errEscapes)
	}
	folder, err := resolveWithin(extensionRoot, path)
	if err != nil {
		return nil, err
	}
	if info, err := os.Stat(folder); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("%s is not a folder", path)
	}
	var files []skillFile
	err = filepath.WalkDir(folder, func(current string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(folder, current)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == ".DS_Store" || strings.HasSuffix(rel, "/.DS_Store") {
			return nil
		}
		if len(files) >= skillFileLimit {
			return fmt.Errorf("%s holds more than %d files", path, skillFileLimit)
		}
		data, err := readWithin(extensionRoot, path+"/"+rel)
		if err != nil {
			return err
		}
		mode := fs.FileMode(0o644)
		if info, err := os.Stat(current); err == nil && info.Mode().Perm()&0o111 != 0 {
			mode = 0o755
		}
		files = append(files, skillFile{rel, data, mode})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !slices.ContainsFunc(files, func(f skillFile) bool { return f.Path == "SKILL.md" }) {
		return nil, fmt.Errorf("%s holds no SKILL.md", path)
	}
	return files, nil
}

// checkManifestSkills checks each skill's folder and SKILL.md frontmatter.
func checkManifestSkills(m Manifest, extensionRoot string) error {
	for _, skill := range m.Skills {
		files, err := readSkillFolder(extensionRoot, skill.Path)
		if err != nil {
			return fmt.Errorf("skill %s: %w", skill.ID, err)
		}
		for _, file := range files {
			if file.Path != "SKILL.md" {
				continue
			}
			if _, _, _, err := skillFrontmatter(file.Data); err != nil {
				return fmt.Errorf("skill %s: %w", skill.ID, err)
			}
		}
	}
	return nil
}

// ExtensionSkill is one skill an extension on for the member ships, as setup
// installs it.
type ExtensionSkill struct {
	Name        string
	Extension   string
	ID          string
	Description string
	// Stage is the qualified named stage it declares, or empty.
	Stage string
	Files []skillFile
}

// extensionSkills is every skill of every extension on for member that this
// machine can read, in extensions.yaml order. A full name two extensions
// produce is installed for the first only, which check reports.
func (s *Store) extensionSkills(member string) ([]ExtensionSkill, error) {
	file, present, err := s.Extensions()
	if err != nil || !present {
		return nil, err
	}
	var skills []ExtensionSkill
	seen := map[string]bool{}
	for index, entry := range file.Extensions {
		loaded := s.loadExtension(entry, index)
		if loaded.Fault != extensionUsable {
			continue
		}
		if on, _ := loaded.memberState(member); !on {
			continue
		}
		root, err := s.Path(loaded.Folder)
		if err != nil {
			continue
		}
		for _, declared := range loaded.Manifest.Skills {
			name := extensionSkillName(entry.ID, declared.ID)
			if seen[name] {
				continue
			}
			files, err := readSkillFolder(root, declared.Path)
			if err != nil {
				continue
			}
			skill := ExtensionSkill{Name: name, Extension: entry.ID, ID: declared.ID}
			if declared.Stage != "" {
				skill.Stage = entry.ID + "/" + declared.Stage
			}
			ok := true
			for i, f := range files {
				if f.Path != "SKILL.md" {
					continue
				}
				_, description, _, err := skillFrontmatter(f.Data)
				renamed, renameErr := renameSkill(f.Data, name)
				if err != nil || renameErr != nil {
					ok = false
					break
				}
				skill.Description, files[i].Data = description, renamed
			}
			if !ok {
				continue
			}
			skill.Files = files
			seen[name] = true
			skills = append(skills, skill)
		}
	}
	return skills, nil
}

// skillCopies is every file setup copies under .agents/skills/, by path, with
// its mode, in a stable order.
func (s *Store) skillCopies() (map[string][]byte, map[string]fs.FileMode, []string, error) {
	member, _ := s.ActiveMember()
	skills, err := s.extensionSkills(member)
	if err != nil {
		return nil, nil, nil, err
	}
	files, modes := map[string][]byte{}, map[string]fs.FileMode{}
	var paths []string
	for _, skill := range skills {
		for _, f := range skill.Files {
			path := ".agents/skills/" + skill.Name + "/" + f.Path
			files[path], modes[path] = f.Data, f.Mode
			paths = append(paths, path)
		}
	}
	return files, modes, paths, nil
}

// skillWrapperHosts are the hosts that read a wrapper pointing at the copy;
// Codex reads .agents/skills/ itself.
var skillWrapperHosts = map[string]string{"claude-code": ".claude", "cursor": ".cursor"}

// skillWrapper is the one-line wrapper a host reads for an installed skill, in
// the form the shipped skills' route files take.
func skillWrapper(skill ExtensionSkill) ([]byte, error) {
	front, err := yaml.Marshal(yaml.MapSlice{{Key: "name", Value: skill.Name}, {Key: "description", Value: skill.Description}})
	if err != nil {
		return nil, err
	}
	return []byte("---\n" + string(front) + "---\n\nRead and follow the canonical workspace skill at\n`.agents/skills/" + skill.Name + "/SKILL.md`.\n"), nil
}

// skillWrappers is each wrapper setup writes for host, by path.
func (s *Store) skillWrappers(host string) (map[string][]byte, []string, error) {
	directory, ok := skillWrapperHosts[host]
	if !ok {
		return nil, nil, nil
	}
	member, _ := s.ActiveMember()
	skills, err := s.extensionSkills(member)
	if err != nil {
		return nil, nil, err
	}
	files := map[string][]byte{}
	var paths []string
	for _, skill := range skills {
		content, err := skillWrapper(skill)
		if err != nil {
			return nil, nil, err
		}
		path := directory + "/skills/" + skill.Name + "/SKILL.md"
		files[path], paths = content, append(paths, path)
	}
	return files, paths, nil
}

// extensionSkillPath reports whether a path is an installed extension skill's.
func extensionSkillPath(path string) bool {
	return strings.Contains(filepath.ToSlash(path), "/skills/x-")
}

// setupSkills copies the skill folders under .agents/skills/, once for every
// host. Like a definition, a copied file a person edited is never
// overwritten, and a copy no extension on for the member produces any longer
// is removed only when nobody edited it.
func (s *Store) setupSkills() (AgentSetup, error) {
	var result AgentSetup
	old := agentInventory{Files: map[string]string{}}
	if err := s.YAML(skillInventoryPath, &old); err != nil && !os.IsNotExist(err) {
		return result, err
	}
	files, modes, paths, err := s.skillCopies()
	if err != nil {
		return result, err
	}
	if len(paths) == 0 && len(old.Files) == 0 {
		return result, nil
	}
	remove, kept, err := s.planSync(old, files, paths, "preserve customized extension skill file: %s")
	if err != nil {
		return result, err
	}
	result.Kept = kept
	written, removed, err := s.applySync(skillInventoryPath, files, modes, paths, remove)
	result.Written, result.Removed = written, removed
	return result, err
}

// skillsBehind reports whether the copied skill folders differ from what
// setup would write now.
func (s *Store) skillsBehind() bool {
	old := agentInventory{Files: map[string]string{}}
	if err := s.YAML(skillInventoryPath, &old); err != nil && !os.IsNotExist(err) {
		return false
	}
	files, _, _, err := s.skillCopies()
	if err != nil {
		return false
	}
	want := map[string]string{}
	for path, content := range files {
		want[path] = digest(content)
	}
	return !maps.Equal(want, old.Files)
}

// pruneSkillFolders removes the folders a removed skill file leaves empty,
// up to the skills folder itself.
func (s *Store) pruneSkillFolders(path string) {
	if !extensionSkillPath(path) {
		return
	}
	dir := filepath.ToSlash(filepath.Dir(filepath.FromSlash(path)))
	for strings.Contains(dir, "/skills/x-") {
		resolved, err := s.Path(dir)
		if err != nil {
			return
		}
		entries, err := os.ReadDir(resolved)
		if err != nil || len(entries) > 0 {
			return
		}
		if os.Remove(resolved) != nil {
			return
		}
		dir = filepath.ToSlash(filepath.Dir(filepath.FromSlash(dir)))
	}
}

// ExtensionSkillState is one extension skill on for the member, as status
// lists it.
type ExtensionSkillState struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	// Stage is the named stage it declares, as other extensions hook it.
	Stage string `json:"stage,omitempty"`
	// Installed reports whether setup copied it on this machine.
	Installed bool `json:"installed"`
}

// extensionSkillStates lists the extension skills on for member, and whether
// each is installed here.
func (s *Store) extensionSkillStates(member string) []ExtensionSkillState {
	skills, err := s.extensionSkills(member)
	if err != nil {
		return nil
	}
	var states []ExtensionSkillState
	for _, skill := range skills {
		state := ExtensionSkillState{Name: skill.Name, Extension: skill.Extension, Stage: skill.Stage}
		if _, err := s.Read(".agents/skills/" + skill.Name + "/SKILL.md"); err == nil {
			state.Installed = true
		}
		states = append(states, state)
	}
	return states
}

// namedStages is every named stage a listed extension readable here declares,
// qualified, mapped to the extension declaring it, whether or not it is on.
func (s *Store) namedStages(file ExtensionsFile) map[string]string {
	stages := map[string]string{}
	for index, entry := range file.Extensions {
		loaded := s.loadExtension(entry, index)
		if loaded.Fault != extensionUsable {
			continue
		}
		for _, skill := range loaded.Manifest.Skills {
			if skill.Stage != "" {
				stages[entry.ID+"/"+skill.Stage] = entry.ID
			}
		}
	}
	return stages
}

// namedStageSkip is why a named point resolves no hook, or empty when its
// extension is listed, readable, on for member, and declares the stage.
func (s *Store) namedStageSkip(file ExtensionsFile, point HookPoint, member string) (owner, reason string) {
	owner, stage, _ := point.Named()
	index := file.extensionIndex(owner)
	if index < 0 {
		return owner, "named stage " + point.Stage + ": extension " + owner + " is not listed in extensions.yaml"
	}
	loaded := s.loadExtension(file.Extensions[index], index)
	if loaded.Fault != extensionUsable {
		return owner, "named stage " + point.Stage + ": " + loaded.Fault + ": " + loaded.Problem
	}
	if on, why := loaded.memberState(member); !on {
		return owner, "named stage " + point.Stage + ": " + why
	}
	if !slices.ContainsFunc(loaded.Manifest.Skills, func(skill ManifestSkill) bool { return skill.Stage == stage }) {
		return owner, "named stage " + point.Stage + ": " + owner + " declares no stage " + stage
	}
	return "", ""
}

// extensionSkillIssues reports what no single manifest shows: two extensions
// whose skills install under one name, and a hook at a named stage no listed
// extension declares.
func (s *Store) extensionSkillIssues(file ExtensionsFile) []Finding {
	var issues []Finding
	names := map[string]string{}
	stages := s.namedStages(file)
	listed := map[string]bool{}
	for _, entry := range file.Extensions {
		listed[entry.ID] = true
	}
	for index, entry := range file.Extensions {
		loaded := s.loadExtension(entry, index)
		if loaded.Fault != extensionUsable {
			continue
		}
		for _, skill := range loaded.Manifest.Skills {
			name := extensionSkillName(entry.ID, skill.ID)
			if other, ok := names[name]; ok && other != entry.ID {
				issues = append(issues, found("extensions "+other+" and "+entry.ID+" both ship skill "+name,
					needsAPerson+"only the first is installed; decide which extension stays, or ask an author to rename the skill"))
				continue
			}
			names[name] = entry.ID
		}
		for _, hook := range loaded.Manifest.Hooks {
			point, err := ParseHookPoint(hook.At)
			if err != nil {
				continue
			}
			owner, _, named := point.Named()
			if !named {
				continue
			}
			if _, declared := stages[point.Stage]; declared {
				continue
			}
			if index := file.extensionIndex(owner); index >= 0 {
				// A listed extension this machine cannot read is reported on
				// its own; what it declares is not known here.
				if s.loadExtension(file.Extensions[index], index).Fault != extensionUsable {
					continue
				}
			}
			why := "no listed extension declares that stage"
			resolve := needsAPerson + "add the extension that declares stage " + point.Stage + " with `extension add`, or ask the author of " + entry.ID + " to drop the hook; until then it is skipped"
			if listed[owner] {
				why = owner + " declares no stage " + strings.TrimPrefix(point.Stage, owner+"/")
				resolve = needsAPerson + "ask the author of " + entry.ID + " to hook a stage " + owner + " declares; until then it is skipped"
			}
			issues = append(issues, found(fmt.Sprintf("extension %s hooks %s, but %s", entry.ID, point, why), resolve))
		}
	}
	return issues
}
