package workspace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// An extension adds work before and after the stages an agent goes through,
// without changing a stage. It is a folder holding an extension.yaml manifest
// and the instruction files it names, and a workspace lists the extensions it
// uses in extensions.yaml at its root, in the order their hooks run.
//
// extensions.yaml is a file of its own rather than a key of workspace.yaml
// because workspace.yaml is decoded strictly: a key an older CLI does not know
// would break every one of them. A workspace without the file behaves exactly as
// one from before extensions existed, and nothing here reads anything else.
//
// A hook returns instructions for the agent, and may run a script. A script
// runs only after this machine's member allowed the extension's exact version
// (extension_trust.go), only with the consent its hook declares, and is launched
// as extension_run.go describes. Every hook that does not run says why.

// ExtensionsPath is the workspace file listing the extensions a workspace uses.
const ExtensionsPath = "extensions.yaml"

// extensionCache holds the remote sources this machine fetched, one folder per
// extension and commit. It sits under .context-circuit/local/, which every
// workspace already ignores, so a fetched source is never committed.
const extensionCache = ".context-circuit/local/extensions"

// manifestName is the one fixed name in an extension folder.
const manifestName = "extension.yaml"

// ExtensionStages are the stages a hook binds to, in lifecycle order.
var ExtensionStages = []string{"intent", "approve", "plan", "prepare", "implement", "review", "deliver", "complete", "archive"}

// CLIFiredStages are the stages a CLI command fires itself, by the command that
// fires each. Their after point is returned by that command and refused by
// `extension hooks`, so an agent cannot fire it twice or without the stage.
var CLIFiredStages = map[string]string{
	"approve":  "record approve",
	"prepare":  "worktree prepare",
	"complete": "record complete",
	"archive":  "record archive",
}

// Member defaults an extension's manifest may declare.
const (
	ExtensionRequired  = "required"
	ExtensionDefaultOn = "default-on"
	ExtensionOptIn     = "opt-in"
)

// ExtensionSwitch is one member's on or off for an extension. YAML 1.1 reads a
// bare on or off as a boolean, and a boolean is not what anyone writing a
// switch meant, so the value is read from its text: on or off, quoted or not,
// and nothing else.
type ExtensionSwitch string

func (v *ExtensionSwitch) UnmarshalYAML(data []byte) error {
	text := strings.TrimSpace(string(data))
	if len(text) >= 2 && (text[0] == '"' && text[len(text)-1] == '"' || text[0] == '\'' && text[len(text)-1] == '\'') {
		text = text[1 : len(text)-1]
	}
	switch text {
	case "on", "off":
		*v = ExtensionSwitch(text)
		return nil
	}
	return fmt.Errorf("a member switch is on or off, not %s", text)
}

// ExtensionSource is where an extension comes from: a folder in the workspace
// (Path alone), or a Git repository pinned to Commit, with Ref the name it was
// resolved from and Path an optional subfolder of that repository.
type ExtensionSource struct {
	Path   string `yaml:"path,omitempty" json:"path,omitempty"`
	URL    string `yaml:"url,omitempty" json:"url,omitempty"`
	Ref    string `yaml:"ref,omitempty" json:"ref,omitempty"`
	Commit string `yaml:"commit,omitempty" json:"commit,omitempty"`
}

// Remote reports whether the source is a Git repository.
func (src ExtensionSource) Remote() bool { return src.URL != "" }

// ExtensionEntry is one extension a workspace uses. Repositories limits it to
// those logical repositories, and none means every one. Settings are passed to
// its instructions and hold no credentials. Members overrides the manifest's
// default for each member named.
type ExtensionEntry struct {
	ID           string                     `yaml:"id" json:"id"`
	Source       ExtensionSource            `yaml:"source" json:"source"`
	Repositories []string                   `yaml:"repositories,omitempty" json:"repositories,omitempty"`
	Settings     map[string]any             `yaml:"settings,omitempty" json:"settings,omitempty"`
	Members      map[string]ExtensionSwitch `yaml:"members,omitempty" json:"members,omitempty"`
}

// ExtensionsFile is extensions.yaml. Order is the order hooks at one point run.
type ExtensionsFile struct {
	Extensions []ExtensionEntry `yaml:"extensions" json:"extensions"`
}

// Manifest is an extension's extension.yaml. Version is informative: what a
// workspace pins is the commit, or the folder it holds.
type Manifest struct {
	ID          string                     `yaml:"id" json:"id"`
	Version     string                     `yaml:"version,omitempty" json:"version,omitempty"`
	Name        string                     `yaml:"name,omitempty" json:"name,omitempty"`
	Summary     string                     `yaml:"summary,omitempty" json:"summary,omitempty"`
	Members     string                     `yaml:"members" json:"members"`
	Settings    map[string]ManifestSetting `yaml:"settings,omitempty" json:"settings,omitempty"`
	Hooks       []ManifestHook             `yaml:"hooks" json:"hooks"`
	RequiresCLI string                     `yaml:"requires_cli,omitempty" json:"requires_cli,omitempty"`
	// Roles are agent roles the extension adds, installed beside the shipped
	// four as x-<extension>-<id>. RoleAdditions add instructions and return
	// headings to a shipped role's brief, never to its definition or access.
	Roles         []ManifestRole `yaml:"roles,omitempty" json:"roles,omitempty"`
	RoleAdditions []RoleAddition `yaml:"role_additions,omitempty" json:"role_additions,omitempty"`
}

// Role access a manifest may declare.
const (
	AccessReadOnly = "read-only"
	AccessEdit     = "edit"
)

// ManifestRole is an agent role an extension declares.
type ManifestRole struct {
	ID          string `yaml:"id" json:"id"`
	Description string `yaml:"description" json:"description"`
	// Access is read-only, briefed like the explorer, or edit, briefed like
	// the worker.
	Access string `yaml:"access" json:"access"`
	// Stages are the stages whose hook instructions its brief carries.
	Stages []string `yaml:"stages,omitempty" json:"stages,omitempty"`
	// Like names the shipped role whose model and effort it takes where the
	// tiering does not name it.
	Like         string   `yaml:"like,omitempty" json:"like,omitempty"`
	Instructions string   `yaml:"instructions" json:"instructions"`
	Returns      []string `yaml:"returns,omitempty" json:"returns,omitempty"`
}

// RoleAddition adds instructions and return headings to a shipped role.
type RoleAddition struct {
	Role         string   `yaml:"role" json:"role"`
	Instructions string   `yaml:"instructions,omitempty" json:"instructions,omitempty"`
	Returns      []string `yaml:"returns,omitempty" json:"returns,omitempty"`
}

// ManifestSetting declares one key an entry's settings may carry.
type ManifestSetting struct {
	Required bool   `yaml:"required,omitempty" json:"required,omitempty"`
	Summary  string `yaml:"summary,omitempty" json:"summary,omitempty"`
}

// ManifestHook is the work an extension adds at one point: an instruction file
// for the agent, a script, or both.
type ManifestHook struct {
	At           string      `yaml:"at" json:"at"`
	Instructions string      `yaml:"instructions,omitempty" json:"instructions,omitempty"`
	Script       *HookScript `yaml:"script,omitempty" json:"script,omitempty"`
}

// HookScript is a script hook as a manifest declares it. Consent none means
// the script acts only locally; standing means enabling and allowing the
// extension authorized Action for this hook; each-time means the person is
// asked before every run.
type HookScript struct {
	Run        []string `yaml:"run" json:"run"`
	RunWindows []string `yaml:"run_windows,omitempty" json:"run_windows,omitempty"`
	OnFailure  string   `yaml:"on_failure" json:"on_failure"`
	Consent    string   `yaml:"consent,omitempty" json:"consent,omitempty"`
	Action     string   `yaml:"action,omitempty" json:"action,omitempty"`
	Timeout    string   `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

// manifestHeader is what can be read from a manifest this CLI cannot decode
// strictly: which CLI it needs, and whether members may switch it off.
type manifestHeader struct {
	ID          string `yaml:"id"`
	Members     string `yaml:"members"`
	RequiresCLI string `yaml:"requires_cli"`
}

// Extensions reads extensions.yaml. present is false when the workspace has
// none, which is every workspace that does not use extensions.
func (s *Store) Extensions() (file ExtensionsFile, present bool, err error) {
	data, err := s.Read(ExtensionsPath)
	if os.IsNotExist(err) {
		return ExtensionsFile{}, false, nil
	}
	if err != nil {
		return ExtensionsFile{}, true, err
	}
	if err := Decode(data, &file); err != nil {
		return ExtensionsFile{}, true, fmt.Errorf("%s: %w", ExtensionsPath, err)
	}
	return file, true, nil
}

// extensionIndex finds an entry by ID, reporting -1 when it is absent.
func (file ExtensionsFile) extensionIndex(id string) int {
	return slices.IndexFunc(file.Extensions, func(entry ExtensionEntry) bool { return entry.ID == id })
}

// HookPoint is a stage's before or after point, written "after implement".
type HookPoint struct {
	When  string
	Stage string
}

func (p HookPoint) String() string { return p.When + " " + p.Stage }

// ParseHookPoint reads "<before|after> <stage>".
func ParseHookPoint(text string) (HookPoint, error) {
	fields := strings.Fields(text)
	if len(fields) != 2 {
		return HookPoint{}, fmt.Errorf("a hook point is \"<before|after> <stage>\", not %q", text)
	}
	return NewHookPoint(fields[1], fields[0])
}

// NewHookPoint checks a stage and when.
func NewHookPoint(stage, when string) (HookPoint, error) {
	if when != "before" && when != "after" {
		return HookPoint{}, fmt.Errorf("a hook fires before or after a stage, not %q", when)
	}
	if !slices.Contains(ExtensionStages, stage) {
		return HookPoint{}, fmt.Errorf("unknown stage %q; the stages are %s", stage, strings.Join(ExtensionStages, ", "))
	}
	return HookPoint{when, stage}, nil
}

// commitPattern is a full object name, SHA-1 or SHA-256.
var commitPattern = regexp.MustCompile(`^[0-9a-f]{40}([0-9a-f]{24})?$`)

// refPattern keeps a ref to the characters a branch or tag name uses, and away
// from anything Git would read as an option.
var refPattern = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9._/@+-]*$`)

func checkRef(ref string) error {
	if !refPattern.MatchString(ref) || strings.Contains(ref, "..") || len(ref) > 200 {
		return fmt.Errorf("name a branch, tag, or commit as the ref: %q", ref)
	}
	return nil
}

// insidePath rejects a path that would leave the folder it is relative to.
func insidePath(value string) error {
	if value == "" {
		return errors.New("empty path")
	}
	cleaned := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	if filepath.IsAbs(value) || strings.HasPrefix(value, "/") || cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return fmt.Errorf("%s leaves its folder", value)
	}
	return nil
}

// readWithin reads a regular file inside root, refusing a path or a symlink
// that resolves outside it: the guard borrowed knowledge applies, for a folder
// the workspace may not own.
func readWithin(root, relative string) ([]byte, error) {
	if err := insidePath(relative); err != nil {
		return nil, err
	}
	resolved, err := resolveWithin(root, relative)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return nil, fmt.Errorf("expected a regular file under 16 MiB: %s", relative)
	}
	return os.ReadFile(resolved)
}

// errEscapes marks a path that resolves outside its extension folder, which
// only the extension's author can fix.
var errEscapes = errors.New("resolves outside the extension folder")

func resolveWithin(root, relative string) (string, error) {
	target := filepath.Join(root, filepath.FromSlash(relative))
	resolved, err := filepath.EvalSymlinks(target)
	if err != nil {
		return "", err
	}
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	if !within(realRoot, resolved) {
		return "", fmt.Errorf("%s %w", relative, errEscapes)
	}
	return resolved, nil
}

// checkSource reports what makes a source unusable as written.
func checkSource(src ExtensionSource) error {
	if !src.Remote() {
		if src.Ref != "" || src.Commit != "" {
			return errors.New("a folder source takes path alone; ref and commit belong to a source with a url")
		}
		if err := insidePath(src.Path); err != nil {
			return fmt.Errorf("a folder source names a path inside the workspace: %w", err)
		}
		return nil
	}
	if strings.HasPrefix(src.URL, "-") {
		return fmt.Errorf("invalid url: %s", src.URL)
	}
	if err := checkRemote(src.URL); err != nil {
		return err
	}
	if src.Ref == "" {
		return errors.New("a source with a url needs the ref it follows")
	}
	if err := checkRef(src.Ref); err != nil {
		return err
	}
	if src.Path != "" {
		if err := insidePath(src.Path); err != nil {
			return fmt.Errorf("a remote source's path names a folder inside that repository: %w", err)
		}
	}
	return nil
}

// validateManifest checks a strictly decoded manifest. Paths are checked by
// checkManifestFiles, which needs the folder.
func validateManifest(m Manifest) error {
	if err := Name(m.ID); err != nil {
		return fmt.Errorf("id: %w", err)
	}
	switch m.Members {
	case ExtensionRequired, ExtensionDefaultOn, ExtensionOptIn:
	default:
		return fmt.Errorf("members is required, default-on, or opt-in, not %q", m.Members)
	}
	for key := range m.Settings {
		if err := settingKey(key); err != nil {
			return err
		}
	}
	seen := map[string]bool{}
	for _, hook := range m.Hooks {
		point, err := ParseHookPoint(hook.At)
		if err != nil {
			return err
		}
		if seen[point.String()] {
			return fmt.Errorf("%s appears twice; each point appears at most once", point)
		}
		seen[point.String()] = true
		if hook.Instructions == "" && hook.Script == nil {
			return fmt.Errorf("%s names neither instructions nor a script", point)
		}
		if hook.Script != nil {
			if err := validateScript(*hook.Script); err != nil {
				return fmt.Errorf("%s script: %w", point, err)
			}
		}
	}
	return validateManifestRoles(m)
}

var settingPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func settingKey(key string) error {
	if !settingPattern.MatchString(key) || len(key) > 100 {
		return fmt.Errorf("setting %q: use lowercase letters, digits and underscores", key)
	}
	return nil
}

// validateScript checks a script's shape, so a manifest written for a CLI that
// runs scripts is accepted here rather than refused.
func validateScript(script HookScript) error {
	if len(script.Run) == 0 {
		return errors.New("run names the program and its arguments")
	}
	for _, command := range [][]string{script.Run, script.RunWindows} {
		for _, arg := range command {
			if strings.TrimSpace(arg) == "" {
				return errors.New("an empty argument")
			}
			for part := range strings.SplitSeq(filepath.ToSlash(arg), "/") {
				if part == ".." {
					return fmt.Errorf("%s %w", arg, errEscapes)
				}
			}
		}
	}
	switch script.OnFailure {
	case "stop", "report":
	default:
		return fmt.Errorf("on_failure is stop or report, not %q", script.OnFailure)
	}
	switch script.Consent {
	case "", "none":
		if script.Action != "" {
			return errors.New("action is named only by a script with standing or each-time consent")
		}
	case "standing", "each-time":
		if err := Text(script.Action); err != nil {
			return fmt.Errorf("consent %s needs the action it covers: %w", script.Consent, err)
		}
	default:
		return fmt.Errorf("consent is none, standing, or each-time, not %q", script.Consent)
	}
	if script.Timeout != "" {
		if duration, err := time.ParseDuration(script.Timeout); err != nil || duration <= 0 {
			return fmt.Errorf("timeout is a duration such as 10m, not %q", script.Timeout)
		}
	}
	return nil
}

// cliSatisfies reports whether the running CLI meets a manifest's requires_cli,
// written ">=X.Y.Z". A development build, and a store that was never told its
// version, run the newest code there is.
func cliSatisfies(running, requirement string) (bool, error) {
	wanted, ok := strings.CutPrefix(strings.TrimSpace(requirement), ">=")
	if !ok {
		return false, fmt.Errorf("requires_cli is written \">=X.Y.Z\", not %q", requirement)
	}
	want, err := releaseNumbers(strings.TrimSpace(wanted))
	if err != nil {
		return false, fmt.Errorf("requires_cli: %w", err)
	}
	if running == "" || strings.HasSuffix(running, "-dev") {
		return true, nil
	}
	have, err := releaseNumbers(running)
	if err != nil {
		return true, nil
	}
	return slices.Compare(have, want) >= 0, nil
}

func releaseNumbers(version string) ([]int, error) {
	core, _, _ := strings.Cut(strings.TrimPrefix(version, "v"), "-")
	core, _, _ = strings.Cut(core, "+")
	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%q is not X.Y.Z", version)
	}
	numbers := make([]int, 3)
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("%q is not X.Y.Z", version)
		}
		numbers[i] = n
	}
	return numbers, nil
}

// Why an extension's manifest could not be used.
const (
	extensionUsable      = ""
	extensionUnavailable = "unavailable" // the folder, or the cached commit, is not there
	extensionUncached    = "uncached"    // a remote whose commit this machine has not fetched
	extensionInvalid     = "invalid"     // the source or manifest is wrong as written
	extensionNewerCLI    = "newer-cli"   // the manifest needs a newer CLI
	extensionEscapes     = "escapes"     // a path the manifest names leaves its folder
)

// loadedExtension is one entry with what this machine can read of it.
type loadedExtension struct {
	Entry  ExtensionEntry
	Index  int
	Folder string // workspace-relative slash path of the extension folder
	// Manifest is set only when it decoded strictly and validated.
	Manifest *Manifest
	Header   manifestHeader
	Fault    string // one of the extension* constants; empty when usable
	Problem  string
}

// extensionFolder is the workspace-relative folder an entry's manifest sits in.
func extensionFolder(entry ExtensionEntry) (string, error) {
	if err := checkSource(entry.Source); err != nil {
		return "", err
	}
	if !entry.Source.Remote() {
		return filepath.ToSlash(filepath.Clean(filepath.FromSlash(entry.Source.Path))), nil
	}
	if !commitPattern.MatchString(entry.Source.Commit) {
		return "", errors.New("records no resolved commit")
	}
	folder := extensionCache + "/" + entry.ID + "/" + entry.Source.Commit
	if entry.Source.Path != "" {
		folder += "/" + filepath.ToSlash(filepath.Clean(filepath.FromSlash(entry.Source.Path)))
	}
	return folder, nil
}

// cacheCheckout is the folder a remote entry's commit is fetched into.
func cacheCheckout(id, commit string) string { return extensionCache + "/" + id + "/" + commit }

// loadExtension reads one entry's manifest without touching the network.
func (s *Store) loadExtension(entry ExtensionEntry, index int) loadedExtension {
	loaded := loadedExtension{Entry: entry, Index: index}
	fail := func(fault, problem string) loadedExtension {
		loaded.Fault, loaded.Problem = fault, problem
		return loaded
	}
	if err := Name(entry.ID); err != nil {
		return fail(extensionInvalid, err.Error())
	}
	folder, err := extensionFolder(entry)
	if err != nil {
		if entry.Source.Remote() && entry.Source.Commit == "" {
			return fail(extensionInvalid, "the source records no resolved commit")
		}
		return fail(extensionInvalid, "source: "+err.Error())
	}
	loaded.Folder = folder
	if entry.Source.Remote() {
		checkout, err := s.Path(cacheCheckout(entry.ID, entry.Source.Commit))
		if err != nil {
			return fail(extensionInvalid, err.Error())
		}
		if info, err := os.Stat(checkout); err != nil || !info.IsDir() {
			return fail(extensionUncached, "commit "+entry.Source.Commit+" is not fetched on this machine")
		}
	}
	absolute, err := s.Path(folder)
	if err != nil {
		return fail(extensionInvalid, err.Error())
	}
	if info, err := os.Stat(absolute); err != nil || !info.IsDir() {
		return fail(extensionUnavailable, "no extension folder at "+folder)
	}
	data, err := readWithin(absolute, manifestName)
	if err != nil {
		if errors.Is(err, errEscapes) {
			return fail(extensionEscapes, manifestName+" "+errEscapes.Error())
		}
		return fail(extensionUnavailable, "no readable "+manifestName+" in "+folder)
	}
	// What a manifest this CLI cannot decode still says: which CLI it needs,
	// and whether a member may switch it off.
	_ = yaml.Unmarshal(data, &loaded.Header)
	if loaded.Header.RequiresCLI != "" {
		ok, err := cliSatisfies(s.CLIVersion, loaded.Header.RequiresCLI)
		if err != nil {
			return fail(extensionInvalid, folder+"/"+manifestName+": "+err.Error())
		}
		if !ok {
			return fail(extensionNewerCLI, fmt.Sprintf("needs a newer CLI: it requires %s, and this CLI is %s", loaded.Header.RequiresCLI, s.CLIVersion))
		}
	}
	var manifest Manifest
	if err := Decode(data, &manifest); err != nil {
		return fail(extensionInvalid, folder+"/"+manifestName+": "+err.Error())
	}
	if err := validateManifest(manifest); err != nil {
		if errors.Is(err, errEscapes) {
			return fail(extensionEscapes, folder+"/"+manifestName+": "+err.Error())
		}
		return fail(extensionInvalid, folder+"/"+manifestName+": "+err.Error())
	}
	for _, hook := range manifest.Hooks {
		if hook.Instructions == "" {
			continue
		}
		if err := insidePath(hook.Instructions); err != nil {
			return fail(extensionEscapes, fmt.Sprintf("%s: %s instructions %s", folder+"/"+manifestName, hook.At, err))
		}
		if _, err := readWithin(absolute, hook.Instructions); err != nil {
			if errors.Is(err, errEscapes) {
				return fail(extensionEscapes, fmt.Sprintf("%s: %s instructions %s", folder+"/"+manifestName, hook.At, err))
			}
			return fail(extensionInvalid, fmt.Sprintf("%s: %s names %s, which cannot be read", folder+"/"+manifestName, hook.At, hook.Instructions))
		}
	}
	for _, file := range manifestRoleFiles(manifest) {
		if err := insidePath(file.Path); err != nil {
			return fail(extensionEscapes, fmt.Sprintf("%s: %s instructions %s", folder+"/"+manifestName, file.Owner, err))
		}
		if _, err := readWithin(absolute, file.Path); err != nil {
			if errors.Is(err, errEscapes) {
				return fail(extensionEscapes, fmt.Sprintf("%s: %s instructions %s", folder+"/"+manifestName, file.Owner, err))
			}
			return fail(extensionInvalid, fmt.Sprintf("%s: %s names %s, which cannot be read", folder+"/"+manifestName, file.Owner, file.Path))
		}
	}
	loaded.Manifest = &manifest
	return loaded
}

// membersDefault is the manifest's default, or what could be read of it.
func (l loadedExtension) membersDefault() string {
	if l.Manifest != nil {
		return l.Manifest.Members
	}
	return l.Header.Members
}

// memberState says whether the extension is on for member, and why when not.
func (l loadedExtension) memberState(member string) (bool, string) {
	choice := l.Entry.Members[member]
	switch l.membersDefault() {
	case ExtensionRequired:
		return true, ""
	case ExtensionOptIn:
		if choice == "on" {
			return true, ""
		}
		if member == "" {
			return false, "opt-in, and no member is selected on this machine"
		}
		return false, "opt-in, and not switched on for member " + member
	default:
		if choice == "off" {
			return false, "off for member " + member
		}
		return true, ""
	}
}

// required reports whether members may not switch the extension off, as far as
// can be read.
func (l loadedExtension) required() bool { return l.membersDefault() == ExtensionRequired }

func (s *Store) loadExtensions(file ExtensionsFile) []loadedExtension {
	loaded := make([]loadedExtension, 0, len(file.Extensions))
	for index, entry := range file.Extensions {
		loaded = append(loaded, s.loadExtension(entry, index))
	}
	return loaded
}

// HookQuery is one hook point to resolve and what it fires for.
type HookQuery struct {
	Point  HookPoint
	Member string
	// Repositories is the record's or worktree's repositories. Scoped false
	// means the point carries no structured repository list, as an intent's
	// does, so only extensions limited to no repositories apply.
	Repositories []string
	Scoped       bool
	// Fetch fetches a remote commit this machine has not cached.
	Fetch bool
	// ListScoped lists a repository-limited extension's instructions, marked
	// with its repositories, where the point carries no repository list. It
	// serves a planner brief, written before the plans name repositories.
	ListScoped bool
	// Run lets allowed scripts run. Without it a script hook is skipped with
	// NotRun as its reason: a brief only lists instructions, and a CLI-fired
	// stage's scripts run only when its command fires them.
	Run    bool
	NotRun string
	// Confirm is the extension whose each-time script the person said yes to.
	// Only that extension's script at this point runs, with confirmed true; the
	// rest of the point already ran on the call that asked.
	Confirm string
	// Script is what a script receives beyond the point and its extension.
	Script ScriptContext
}

// ScriptContext is the record and worktrees a hook point fires for, and the
// folder its scripts run in.
type ScriptContext struct {
	Record    *ScriptRecord
	Records   []ScriptRecord
	Worktrees map[string]string
	// Dir is the worktree for a repository-scoped stage, and the workspace
	// root when empty.
	Dir string
}

// HookInstruction is one instruction file to follow at a point, as a
// workspace-relative path the agent opens itself.
type HookInstruction struct {
	ID       string         `json:"id"`
	File     string         `json:"file"`
	Settings map[string]any `json:"settings,omitempty"`
	// Repositories is set only in a listing for a point with no repository
	// list: these instructions apply to work touching those repositories.
	Repositories []string `json:"repositories,omitempty"`
}

// HookSkip is an extension or hook that did not apply, and why.
type HookSkip struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

// HookConfirmation is an each-time script waiting for the person's yes, with
// the action it would take in the manifest's words and the call that runs it.
type HookConfirmation struct {
	ID      string `json:"id"`
	Action  string `json:"action"`
	Confirm string `json:"confirm"`
}

// HookResult is what one hook point returns.
type HookResult struct {
	Point        string            `json:"point"`
	Instructions []HookInstruction `json:"instructions"`
	// Ran is each script that ran, and what came of it.
	Ran []HookRun `json:"ran,omitempty"`
	// ConfirmationRequired is each each-time script the person must say yes
	// to before it runs.
	ConfirmationRequired []HookConfirmation `json:"confirmation_required,omitempty"`
	Skipped              []HookSkip         `json:"skipped,omitempty"`
	// Stopped means the stage must not go on until a person decides: a
	// required extension could not be read or its script is not allowed on
	// this machine, or a script whose on_failure is stop failed.
	Stopped bool   `json:"stopped"`
	Resolve string `json:"resolve,omitempty"`
}

// empty reports a point where nothing applied and nothing was skipped.
func (r HookResult) empty() bool {
	return len(r.Instructions) == 0 && len(r.Ran) == 0 && len(r.ConfirmationRequired) == 0 && len(r.Skipped) == 0 && !r.Stopped
}

// scriptsNotRun is the reason a script hook gives where nothing may run it.
const scriptsNotRun = "script not run: only the coordinator runs scripts, through `extension hooks`"

// ResolveHooks resolves one hook point. listed is false when the workspace has
// no extensions.yaml, and the result is then empty.
func (s *Store) ResolveHooks(ctx context.Context, query HookQuery) (result HookResult, listed bool, err error) {
	result = HookResult{Point: query.Point.String(), Instructions: []HookInstruction{}}
	file, present, err := s.Extensions()
	if err != nil || !present {
		return result, present, err
	}
	var resolve []string
	confirmed := false
	for index, entry := range file.Extensions {
		if query.Confirm != "" && entry.ID != query.Confirm {
			continue
		}
		loaded := s.loadExtension(entry, index)
		if loaded.Fault == extensionUncached && query.Fetch {
			if _, fetchErr := s.fetchExtension(ctx, entry); fetchErr == nil {
				loaded = s.loadExtension(entry, index)
			}
		}
		if loaded.Fault != extensionUsable {
			reason := loaded.Fault + ": " + loaded.Problem
			// What cannot be read cannot be shown to be optional. A member's
			// own off still holds unless the manifest said it was required.
			result.Skipped = append(result.Skipped, HookSkip{entry.ID, reason})
			optional := loaded.Header.Members == ExtensionDefaultOn || loaded.Header.Members == ExtensionOptIn
			if !loaded.required() && (optional || loaded.Entry.Members[query.Member] == "off") {
				continue
			}
			result.Stopped = true
			resolve = append(resolve, extensionRemedy(loaded))
			continue
		}
		on, why := loaded.memberState(query.Member)
		if !on {
			result.Skipped = append(result.Skipped, HookSkip{entry.ID, why})
			continue
		}
		var hook *ManifestHook
		for i := range loaded.Manifest.Hooks {
			if point, _ := ParseHookPoint(loaded.Manifest.Hooks[i].At); point == query.Point {
				hook = &loaded.Manifest.Hooks[i]
			}
		}
		if hook == nil {
			continue
		}
		marked := []string(nil)
		if len(entry.Repositories) > 0 {
			switch {
			case !query.Scoped && query.ListScoped:
				marked = entry.Repositories
			case !query.Scoped:
				result.Skipped = append(result.Skipped, HookSkip{entry.ID, "limited to repositories " + strings.Join(entry.Repositories, ", ") + ", and " + query.Point.String() + " carries no repository list"})
				continue
			case !overlaps(entry.Repositories, query.Repositories):
				result.Skipped = append(result.Skipped, HookSkip{entry.ID, "outside repositories: limited to " + strings.Join(entry.Repositories, ", ")})
				continue
			}
		}
		if hook.Instructions != "" && query.Confirm == "" {
			result.Instructions = append(result.Instructions, HookInstruction{
				ID:           entry.ID,
				File:         loaded.Folder + "/" + filepath.ToSlash(filepath.Clean(filepath.FromSlash(hook.Instructions))),
				Settings:     entry.Settings,
				Repositories: marked,
			})
		}
		if hook.Script != nil {
			if remedy := s.runHookScript(ctx, query, loaded, *hook, &result); remedy != "" {
				resolve = append(resolve, remedy)
			}
			confirmed = confirmed || query.Confirm == entry.ID && hook.Script.Consent == "each-time"
		}
	}
	if query.Confirm != "" && !confirmed {
		return HookResult{}, true, fmt.Errorf("%s has no each-time script at %s that is on for this member; --confirm runs only a script that asked for the person's yes", query.Confirm, query.Point)
	}
	if result.Stopped {
		result.Resolve = "take this to the person before the stage goes on: " + strings.Join(unique(resolve), "; ")
	}
	return result, true, nil
}

// runHookScript is the run step between resolving a hook and returning it:
// it checks the allow, then consent, launches the script, and folds the outcome
// into result. It returns what resolves a stop it caused, or "".
func (s *Store) runHookScript(ctx context.Context, query HookQuery, loaded loadedExtension, hook ManifestHook, result *HookResult) string {
	id, script := loaded.Entry.ID, *hook.Script
	skip := func(reason string) { result.Skipped = append(result.Skipped, HookSkip{id, reason}) }
	if query.Confirm != "" && script.Consent != "each-time" {
		return ""
	}
	if !query.Run {
		reason := query.NotRun
		if reason == "" {
			reason = scriptsNotRun
		}
		skip(reason)
		return ""
	}
	if _, why := scriptCommand(script, scriptGOOS); why != "" {
		skip(why)
		return ""
	}
	version, err := s.extensionVersion(loaded)
	state := TrustNotAllowed
	if err == nil {
		state, err = s.trustState(query.Member, id, version)
	}
	if err != nil {
		skip("script not run: its version could not be read: " + err.Error())
		if loaded.required() {
			result.Stopped = true
			return "the script of " + id + " could not be checked against this machine's allow: " + err.Error()
		}
		return ""
	}
	if state != TrustAllowed {
		skip(notAllowedReason(state, id))
		if loaded.required() {
			// A required extension's work is part of the stage, so the stage
			// does not go on without it; allowing it is still the person's call.
			result.Stopped = true
			return id + " is required and its script is " + state + " on this machine: ask the person whether to run `extension allow --id " + id + "`"
		}
		return ""
	}
	confirmed := false
	if script.Consent == "each-time" {
		if _, fired := CLIFiredStages[query.Point.Stage]; fired {
			skip("each-time script not run: " + query.Point.String() + " is fired by `" + CLIFiredStages[query.Point.Stage] + "`, where no yes can be given; its author can move it to a stage the agent fires")
			return ""
		}
		if query.Confirm != id {
			confirm := "extension hooks --stage " + query.Point.Stage + " --when " + query.Point.When
			if query.Script.Record != nil {
				if query.Script.Record.Kind == "plan" {
					confirm += " --plan " + query.Script.Record.ID
				} else {
					confirm += " --intent " + query.Script.Record.ID
				}
			}
			if query.Scoped && len(query.Repositories) == 1 {
				confirm += " --repo " + query.Repositories[0]
			}
			result.ConfirmationRequired = append(result.ConfirmationRequired, HookConfirmation{ID: id, Action: script.Action, Confirm: confirm + " --confirm " + id})
			return ""
		}
		confirmed = true
	}
	folder, err := s.Path(loaded.Folder)
	if err != nil {
		skip("script not run: " + err.Error())
		return ""
	}
	dir := query.Script.Dir
	if dir == "" {
		dir = s.Root
	}
	worktrees := query.Script.Worktrees
	if worktrees == nil {
		worktrees = map[string]string{}
	}
	settings := loaded.Entry.Settings
	if settings == nil {
		settings = map[string]any{}
	}
	run := runScript(ctx, folder, script, dir, ScriptInput{
		Point: query.Point.String(), WorkspaceRoot: s.Root, Member: query.Member,
		Extension: ScriptExtension{ID: id, Settings: settings},
		Record:    query.Script.Record, Records: query.Script.Records,
		Worktrees: worktrees, Confirmed: confirmed,
	})
	result.Ran = append(result.Ran, run)
	if run.Status == ScriptFailed && script.OnFailure == "stop" {
		result.Stopped = true
		return id + " failed at " + query.Point.String() + " and its on_failure is stop; it is not re-run until the person decides"
	}
	return ""
}

// extensionRemedy names what makes an unreadable extension readable.
func extensionRemedy(loaded loadedExtension) string {
	id := loaded.Entry.ID
	switch loaded.Fault {
	case extensionUncached:
		return "fetch it with `extension sync`"
	case extensionNewerCLI:
		return "upgrade the CLI through `.agents/skills/cc-cli/SKILL.md` to use " + id
	case extensionEscapes:
		return "the author of " + id + " must keep every path inside its folder"
	default:
		return "fix the source of " + id + " or remove it with `extension remove --id " + id + "`; `check` names what is wrong"
	}
}

func overlaps(left, right []string) bool {
	for _, value := range left {
		if slices.Contains(right, value) {
			return true
		}
	}
	return false
}

// StageHook resolves one point of a stage a CLI command fires, for the
// repositories and record it fires for, running the scripts allowed there. The
// command calls it outside the workspace lock: the before point before taking
// the lock, the after point once it is released, so a slow script never blocks
// another command and a script may call the CLI itself. It is nil when the
// workspace lists no extensions or nothing applied, so a command's result keeps
// its shape there. A point it could not resolve is returned stopped, naming
// what to do.
func (s *Store) StageHook(ctx context.Context, point HookPoint, repositories []string, scoped bool, script ScriptContext) *HookResult {
	member, _ := s.ActiveMember()
	result, listed, err := s.ResolveHooks(ctx, HookQuery{Point: point, Member: member, Repositories: repositories, Scoped: scoped, Run: true, Script: script})
	if err != nil {
		result = HookResult{Point: point.String(), Instructions: []HookInstruction{}, Stopped: true,
			Resolve: "take this to the person: " + err.Error() + "; `check` names what is wrong"}
	} else if !listed || result.empty() {
		return nil
	}
	return &result
}

// ScriptRecordFor is a record as a script receives it, or nil when id names
// none.
func (s *Store) ScriptRecordFor(id string) *ScriptRecord {
	record, err := s.FindRecord(id)
	if err != nil {
		return nil
	}
	kind := "intent"
	if strings.HasPrefix(record.ID, "p") {
		kind = "plan"
	}
	path := record.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(s.Root, filepath.FromSlash(path))
	}
	repositories := record.Repositories
	if repositories == nil {
		repositories = []string{}
	}
	return &ScriptRecord{ID: record.ID, Kind: kind, Path: path, Repositories: repositories}
}

// PlanWorktrees is each repository's worktree prepared for plan on this
// machine, as absolute paths, read from the local association.
func (s *Store) PlanWorktrees(plan string) map[string]string {
	worktrees := map[string]string{}
	local, err := s.associations()
	if err != nil {
		return worktrees
	}
	for _, assoc := range local.Worktrees {
		if assoc.Plan != plan || assoc.Path == "" {
			continue
		}
		path := assoc.Path
		if !filepath.IsAbs(path) {
			path = filepath.Join(s.Root, path)
		}
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			worktrees[assoc.Repository] = path
		}
	}
	return worktrees
}

// PreparedWorktree is what `worktree prepare` returns: the worktree, and the
// extension blocks its stage fired, which are absent where none applied.
type PreparedWorktree struct {
	Worktree
	Extensions []HookResult `json:"extensions,omitempty"`
	// Next names the implement hook commands to run in this worktree, where
	// an extension applies at either implement point. The CLI never sees
	// implementation, so this is where it says so.
	Next []string `json:"next,omitempty"`
}

// ImplementHooks names the `extension hooks` commands to run before the first
// change and after the last commit in a worktree for plan and repo, or is
// empty when no extension on for the active member applies there. A skip
// because the extension is off for the member, opt-in and not switched on, or
// outside its repositories, does not count: nothing would follow from running
// the hook.
func (s *Store) ImplementHooks(ctx context.Context, plan, repo string) []string {
	member, _ := s.ActiveMember()
	applies := false
	for _, when := range []string{"before", "after"} {
		result, listed, err := s.ResolveHooks(ctx, HookQuery{Point: HookPoint{when, "implement"}, Member: member, Repositories: []string{repo}, Scoped: true})
		if err != nil || !listed {
			return nil
		}
		if len(result.Instructions) > 0 || result.Stopped {
			applies = true
		}
		for _, skip := range result.Skipped {
			if !strings.HasPrefix(skip.Reason, "off for member") && !strings.HasPrefix(skip.Reason, "opt-in, and") && !strings.HasPrefix(skip.Reason, "outside repositories") {
				applies = true
			}
		}
	}
	if !applies {
		return nil
	}
	target := "--repo " + repo
	if plan != "" {
		target = "--plan " + plan + " " + target
	}
	return []string{
		"sh .context-circuit/cli extension hooks --stage implement --when before " + target + "   # before the first change in this worktree",
		"sh .context-circuit/cli extension hooks --stage implement --when after " + target + "   # after the last commit in this worktree",
	}
}
