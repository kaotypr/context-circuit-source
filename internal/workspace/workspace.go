package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Repository is the shared description of one logical repository: where it
// lives and which branch it defaults to. Where a machine keeps its checkout,
// and which branch that machine starts work from, belong to the local binding
// in repositories.local.yaml instead, because both differ per machine.
type Repository struct {
	URL           string `yaml:"url,omitempty" json:"url,omitempty"`
	DefaultBranch string `yaml:"default_branch" json:"default_branch"`
}
type Relationship struct {
	From        string `yaml:"from" json:"from"`
	To          string `yaml:"to" json:"to"`
	Description string `yaml:"description" json:"description"`
}

// CLIRegistry names a GitLab project mirroring CLI releases. It is shared, so
// one member records the mirror and everyone who clones the workspace installs
// from it without configuring their own environment. Unset, the installer reads
// the product's own GitHub releases.
//
// WorkspaceRepository describes the Git repository carrying this workspace
// itself. It sits beside `repositories` rather than inside it because every
// consumer of that map treats an entry as somewhere work happens: a plan names
// repositories, a relationship joins two, and a worktree is cut from one. The
// workspace is none of those — a worktree of it would duplicate the records the
// CLI is reading — so it is described separately and reported separately.
// KnowledgeRepository is a repository of knowledge this workspace reads but
// does not own: an organization's knowledge center, shared by many workspaces
// and changed through its own repository. It sits beside `repositories` for the
// same reason `workspace_repository` does — every consumer of that map treats an
// entry as somewhere work happens, and a plan, a relationship, or a worktree cut
// from a knowledge repository is never what the caller meant.
//
// Index names the file that maps its contents, the way `context/INDEX.md` maps
// this workspace's own notes. It is recorded because a borrowed repository
// chooses its own entry point and nothing here may rename it.
type KnowledgeRepository struct {
	URL           string `yaml:"url,omitempty" json:"url,omitempty"`
	DefaultBranch string `yaml:"default_branch" json:"default_branch"`
	Index         string `yaml:"index,omitempty" json:"index,omitempty"`
}
type Config struct {
	Version               int                            `yaml:"version" json:"version"`
	Name                  string                         `yaml:"name" json:"name"`
	Purpose               string                         `yaml:"purpose" json:"purpose"`
	CLIRegistry           string                         `yaml:"cli_registry,omitempty" json:"cli_registry,omitempty"`
	WorkspaceRepository   *Repository                    `yaml:"workspace_repository,omitempty" json:"workspace_repository,omitempty"`
	Repositories          map[string]Repository          `yaml:"repositories" json:"repositories"`
	KnowledgeRepositories map[string]KnowledgeRepository `yaml:"knowledge_repositories,omitempty" json:"knowledge_repositories,omitempty"`
	// KnowledgeReviewDays reports a note that anchors to no code and has gone
	// this long unconfirmed. It is off unless a workspace sets it, because age
	// alone is evidence of nothing: a note whose anchors nobody touched is not
	// stale however old it is, and that case is answered from the code instead.
	KnowledgeReviewDays int            `yaml:"knowledge_review_days,omitempty" json:"knowledge_review_days,omitempty"`
	Relationships       []Relationship `yaml:"relationships" json:"relationships"`
}

// Band is an optional allocation block index. Members holding distinct bands
// allocate from disjoint numeric ranges, so two clones that cannot see each
// other still never choose the same ID. Band 0 means unbanded: those members
// share the range outside every declared band, which is how a solo workspace
// and every workspace predating a band keeps allocating from i001 and p0001.
//
// Language names the language this member's intents and plans are written in,
// as a person would say it rather than as a tag, because the only reader is a
// dispatch brief quoting it into a sentence. A tag would need a table mapping
// it to a name, and that table is wrong for the first language it omits. Unset
// means English, which is what every record written before this field assumed.
//
// Tone is the register that language is written in, recorded the same way and
// for the same reason: it is quoted into the instruction that writes a record,
// never parsed. It exists because naming a language settles which words are
// used and nothing about how they are put together, and an agent composing in
// English and translating produces prose that is grammatical, stiff, and
// formal in a way nobody chose. A team says what it wants once —
// "semi-formal; keep technical terms in English" — instead of reading it back
// in every record. Unset leaves the general composition guidance to decide.
type Member struct {
	Name     string `yaml:"name" json:"name"`
	Band     int    `yaml:"band,omitempty" json:"band,omitempty"`
	Language string `yaml:"language,omitempty" json:"language,omitempty"`
	Tone     string `yaml:"tone,omitempty" json:"tone,omitempty"`
}
type Members struct {
	Members map[string]Member `yaml:"members" json:"members"`
}
type Identity struct {
	Member string `yaml:"member" json:"member"`
}

// Binding is this machine's checkout of a repository: its path, and the branch
// work starts from here. A binding written before the base moved here records
// no base and resolves to the shared default.
type Binding struct {
	Path       string `yaml:"path" json:"path"`
	BaseBranch string `yaml:"base_branch,omitempty" json:"base_branch,omitempty"`
}

// Workspace is this machine's binding for the workspace repository itself. It
// is separate from the map for the same reason the shared record is, and it
// carries a path because a workspace need not sit at its repository's root: a
// workspace kept in a subdirectory records the containing root as `..`.
// Knowledge holds this machine's checkout of each borrowed knowledge
// repository. A knowledge binding records a path and no base branch: work never
// starts here, so there is no branch to start it from. The tracked branch is the
// shared `default_branch`, because every member reads the same knowledge and a
// machine reading another branch of it is reading something nobody else has.
type Bindings struct {
	Workspace *Binding           `yaml:"workspace,omitempty" json:"workspace,omitempty"`
	Bindings  map[string]Binding `yaml:"bindings" json:"bindings"`
	Knowledge map[string]Binding `yaml:"knowledge,omitempty" json:"knowledge,omitempty"`
}

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)

func Name(value string) error {
	if !namePattern.MatchString(value) || len(value) > 100 {
		return fmt.Errorf("use a lowercase ID/slug with hyphens, up to 100 characters: %q", value)
	}
	return nil
}
func Text(value string) error {
	if strings.TrimSpace(value) == "" || strings.ContainsAny(value, "\r\n\x00") {
		return errors.New("provide nonempty single-line text")
	}
	return nil
}

func (s *Store) Config() (Config, error) {
	var cfg Config
	if err := s.YAML("workspace.yaml", &cfg); err != nil {
		// Candidates through 2.0.0-rc.6 recorded base_branch here, where one
		// value described every machine at once. The decoder can only report an
		// unknown key, which reads as a corrupt file, so name the move.
		if strings.Contains(err.Error(), "base_branch") {
			return cfg, fmt.Errorf("%w; base_branch now belongs to this machine's binding in repositories.local.yaml, and a repository here records url and default_branch", err)
		}
		return cfg, err
	}
	// Separate the two reasons a version does not match, as initialization
	// already does: a higher schema is this CLI being old, and answering it with
	// migration advice sends the reader to rewrite a workspace that is fine.
	if cfg.Version > 2 {
		return cfg, fmt.Errorf("workspace schema %d requires a newer CLI; this CLI supports schema 2", cfg.Version)
	}
	if cfg.Version != 2 || cfg.Name == "" {
		return cfg, errors.New("initialize a fresh v2 workspace first; v1 migration is not automatic")
	}
	if cfg.Repositories == nil || cfg.Relationships == nil {
		return cfg, errors.New("workspace.yaml needs repositories and relationships collections")
	}
	if cfg.CLIRegistry != "" {
		// Reject here what the installer would reject later, so a mirror that
		// cannot be resolved fails where it is recorded rather than at install.
		rest, https := strings.CutPrefix(cfg.CLIRegistry, "https://")
		host, project, split := strings.Cut(rest, "/")
		if !https || !split || host == "" || strings.Trim(project, "/") == "" {
			return cfg, errors.New("cli_registry must be an https project URL, such as https://gitlab.example.com/group/project")
		}
	}
	for id, repo := range cfg.Repositories {
		if err := Name(id); err != nil {
			return cfg, err
		}
		if err := Text(repo.DefaultBranch); err != nil {
			return cfg, err
		}
		if repo.URL != "" {
			if err := checkRemote(repo.URL); err != nil {
				return cfg, fmt.Errorf("%s: %w", id, err)
			}
		}
	}
	return cfg, nil
}

func (s *Store) Members() (Members, error) {
	var members Members
	if err := s.YAML("members.yaml", &members); err != nil {
		return members, err
	}
	if members.Members == nil {
		return members, errors.New("members.yaml needs a members mapping")
	}
	bands := map[int]string{}
	for id, member := range members.Members {
		if err := Name(id); err != nil {
			return members, err
		}
		if err := Text(member.Name); err != nil {
			return members, err
		}
		if err := Band(member.Band); err != nil {
			return members, fmt.Errorf("%s: %w", id, err)
		}
		// Two members sharing a band share a range, which defeats the only
		// thing a band is for. Refuse to allocate until it is resolved.
		if member.Band != 0 {
			if other, taken := bands[member.Band]; taken {
				return members, fmt.Errorf("members %s and %s both hold band %d", other, id, member.Band)
			}
			bands[member.Band] = id
		}
	}
	return members, nil
}

// claimed reports whether a number falls inside some member's band block, which
// is what keeps unbanded allocation out of banded territory.
func (m Members) claimed(number, block int) bool {
	for _, member := range m.Members {
		if member.Band != 0 && number >= member.Band*block && number < member.Band*block+block {
			return true
		}
	}
	return false
}

// Band validates an allocation block index. Zero is unbanded.
// Language is quoted into a brief, so it is one short line of prose naming a
// language, not free text carrying instructions of its own.
func Language(language string) error {
	if language == "" {
		return nil
	}
	if err := Text(language); err != nil {
		return err
	}
	if len(language) > 40 {
		return errors.New("name the language as a person would say it, such as English or Bahasa Indonesia")
	}
	return nil
}

// Tone is quoted into the instruction that writes a record, so it is one line
// of prose describing a register, not free text carrying instructions of its
// own. It is allowed more room than a language name because naming a register
// usually takes a clause and an example rather than a word.
func Tone(tone string) error {
	if tone == "" {
		return nil
	}
	if err := Text(tone); err != nil {
		return err
	}
	if len(tone) > 200 {
		return errors.New("describe the register in one line, such as 'semi-formal; keep technical terms in English'")
	}
	return nil
}

func Band(band int) error {
	if band < 0 || band > 9999 {
		return errors.New("band must be between 1 and 9999, or absent")
	}
	return nil
}

func (s *Store) Bindings() (Bindings, error) {
	bindings := Bindings{Bindings: map[string]Binding{}}
	err := s.YAML("repositories.local.yaml", &bindings)
	if os.IsNotExist(err) {
		return bindings, nil
	}
	if err == nil && bindings.Bindings == nil {
		err = errors.New("local bindings must be a mapping")
	}
	return bindings, err
}

// BaseBranch reports the branch work starts from for one repository on this
// machine: the binding's own base when it records one, and the repository's
// shared default otherwise.
func (b Bindings) BaseBranch(cfg Config, id string) string {
	if binding, ok := b.Bindings[id]; ok && binding.BaseBranch != "" {
		return binding.BaseBranch
	}
	return cfg.Repositories[id].DefaultBranch
}

// WorkspaceBase reports the branch workspace edits start from on this machine,
// falling back to the shared default the same way a repository's does.
func (b Bindings) WorkspaceBase(cfg Config) string {
	if b.Workspace != nil && b.Workspace.BaseBranch != "" {
		return b.Workspace.BaseBranch
	}
	if cfg.WorkspaceRepository != nil {
		return cfg.WorkspaceRepository.DefaultBranch
	}
	return ""
}

func (s *Store) ActiveMember() (string, error) {
	var local Identity
	if err := s.YAML("member.local.yaml", &local); err != nil {
		return "", fmt.Errorf("select a workspace member: %w", err)
	}
	members, err := s.Members()
	if err != nil {
		return "", err
	}
	if _, ok := members.Members[local.Member]; !ok {
		return "", errors.New("active member is absent from members.yaml")
	}
	return local.Member, nil
}

func (s *Store) Init(files map[string][]byte, name, purpose, member, display string) error {
	for _, value := range []string{name, purpose, display} {
		if err := Text(value); err != nil {
			return err
		}
	}
	if err := Name(member); err != nil {
		return err
	}
	for _, rel := range []string{"member.local.yaml", "repositories.local.yaml"} {
		p, err := s.Path(rel)
		if err != nil {
			return err
		}
		if _, err := os.Lstat(p); !os.IsNotExist(err) {
			return fmt.Errorf("existing local state: %s", rel)
		}
	}
	// A previously exported blank seed is also supported, but only identical seed
	// text files may already occupy generated destinations. Treat Git's CRLF/LF
	// checkout conversion as equivalent, but never overwrite content edits.
	missing := map[string][]byte{}
	// The template and CLI have independent versions. A blank compatible
	// template owns its own instructions; initialization changes its data only.
	var blank Config
	existingSeed := false
	if err := s.YAML("workspace.yaml", &blank); err == nil {
		// Separate the two reasons a template cannot be initialized. An
		// incompatible schema is a CLI problem; data is a workspace problem.
		if blank.Version != 2 {
			return fmt.Errorf("workspace schema %d requires a newer CLI; this CLI supports schema 2", blank.Version)
		}
		if blank.Name != "" || blank.Purpose != "" || blank.WorkspaceRepository != nil || len(blank.Repositories) != 0 || len(blank.Relationships) != 0 {
			return errors.New("existing workspace data; initialization stopped")
		}
		var members Members
		if err := s.YAML("members.yaml", &members); err != nil {
			return err
		}
		if len(members.Members) != 0 {
			return errors.New("existing members; initialization stopped")
		}
		existingSeed = true
	} else if !os.IsNotExist(err) {
		return err
	}
	for rel, seed := range files {
		if existingSeed {
			continue
		}
		current, err := s.Read(rel)
		if os.IsNotExist(err) {
			missing[rel] = seed
		} else if err != nil {
			return err
		} else if !bytes.Equal(bytes.ReplaceAll(current, []byte("\r\n"), []byte("\n")), bytes.ReplaceAll(seed, []byte("\r\n"), []byte("\n"))) {
			return fmt.Errorf("initialization would overwrite existing data: %s", rel)
		}
	}
	if paths, err := s.RecordPaths(true); err != nil {
		return err
	} else if len(paths) != 0 {
		return errors.New("existing intent or plan records; initialization stopped")
	}
	cfg := Config{Version: 2, Name: name, Purpose: purpose, Repositories: map[string]Repository{}, Relationships: []Relationship{}}
	if err := s.Export(missing); err != nil {
		return err
	}
	if err := s.WriteYAML("workspace.yaml", cfg, 0644); err != nil {
		return err
	}
	if err := s.WriteYAML("members.yaml", Members{map[string]Member{member: {Name: display}}}, 0644); err != nil {
		return err
	}
	if err := s.WriteYAML("repositories.local.yaml", Bindings{Bindings: map[string]Binding{}}, 0600); err != nil {
		return err
	}
	if err := s.WriteYAML("member.local.yaml", Identity{member}, 0600); err != nil {
		return err
	}
	// The README travels with the workspace, so it carries the workspace's name
	// and the versions it received rather than the product's newest release.
	if err := s.TitleReadme(name); err != nil {
		return err
	}
	// Write role definitions for every host now, from this CLI, rather than
	// leaving a setup step between a new workspace and its first delegation.
	// Seeding them from the template instead would ship files no inventory
	// records, which a later version reads as a customization and refuses to
	// replace. These carry the shipped inherit settings; `agent configure`
	// still needs `agent setup` afterwards.
	_, err := s.SetupAgents("")
	return err
}

func (s *Store) AddMember(id, name string, band int, language, tone string) error {
	if _, err := s.Config(); err != nil {
		return err
	}
	if err := Name(id); err != nil {
		return err
	}
	if err := Text(name); err != nil {
		return err
	}
	if err := Band(band); err != nil {
		return err
	}
	if err := Language(language); err != nil {
		return err
	}
	if err := Tone(tone); err != nil {
		return err
	}
	members, err := s.Members()
	if err != nil {
		return err
	}
	if existing, ok := members.Members[id]; ok {
		if existing.Name != name {
			return fmt.Errorf("member %s already has another name", id)
		}
		if language != "" && existing.Language != language {
			return fmt.Errorf("member %s already writes in %s; use member language to change it", id, existingLanguage(existing))
		}
		if tone != "" && existing.Tone != tone {
			return fmt.Errorf("member %s already records a tone; use member tone to change it", id)
		}
		if band == 0 || existing.Band == band {
			return nil
		}
		return fmt.Errorf("member %s already holds band %d; use member band to change it", id, existing.Band)
	}
	if err := s.freeBand(members, id, band); err != nil {
		return err
	}
	return s.Update("members.yaml", []string{"members", id}, Member{name, band, language, tone}, 0644)
}

// existingLanguage names what a member writes in, including the unset case, so
// a refusal says what is already recorded rather than leaving a blank.
func existingLanguage(member Member) string {
	if member.Language == "" {
		return "English"
	}
	return member.Language
}

// SetMemberLanguage records the language a member's intents and plans are
// written in. Knowledge is English whatever this says: a note outlives the
// member who wrote it and anchors to code, while an intent is approved by a
// person who has to understand it.
func (s *Store) SetMemberLanguage(id, language string) error {
	if err := Language(language); err != nil {
		return err
	}
	if language == "" {
		return errors.New("name the language, such as English or Bahasa Indonesia")
	}
	members, err := s.Members()
	if err != nil {
		return err
	}
	member, ok := members.Members[id]
	if !ok {
		return fmt.Errorf("unknown member: %s", id)
	}
	member.Language = language
	return s.Update("members.yaml", []string{"members", id}, member, 0644)
}

// SetMemberTone records the register a member's intents and plans are written
// in. It is set and cleared rather than only set, because a team that recorded
// a register and changed its mind wants the general guidance back, and there is
// no sentinel value for prose the way band 0 clears a band.
func (s *Store) SetMemberTone(id, tone string, clear bool) error {
	if err := Tone(tone); err != nil {
		return err
	}
	if tone == "" && !clear {
		return errors.New("describe the register, such as 'semi-formal; keep technical terms in English', or pass --clear")
	}
	if tone != "" && clear {
		return errors.New("pass --tone or --clear, not both")
	}
	members, err := s.Members()
	if err != nil {
		return err
	}
	member, ok := members.Members[id]
	if !ok {
		return fmt.Errorf("unknown member: %s", id)
	}
	member.Tone = tone
	return s.Update("members.yaml", []string{"members", id}, member, 0644)
}

// SetMemberBand assigns or clears an existing member's allocation band. Records
// already allocated keep their IDs and their reservations: a band decides which
// numbers come next, never which numbers were right.
func (s *Store) SetMemberBand(id string, band int) error {
	if err := Band(band); err != nil {
		return err
	}
	members, err := s.Members()
	if err != nil {
		return err
	}
	member, ok := members.Members[id]
	if !ok {
		return fmt.Errorf("unknown member: %s", id)
	}
	if err := s.freeBand(members, id, band); err != nil {
		return err
	}
	member.Band = band
	return s.Update("members.yaml", []string{"members", id}, member, 0644)
}

func (s *Store) freeBand(members Members, id string, band int) error {
	if band == 0 {
		return nil
	}
	for other, member := range members.Members {
		if other != id && member.Band == band {
			return fmt.Errorf("band %d already belongs to member %s", band, other)
		}
	}
	return nil
}

func (s *Store) UseMember(id string) error {
	members, err := s.Members()
	if err != nil {
		return err
	}
	if _, ok := members.Members[id]; !ok {
		return fmt.Errorf("unknown member: %s", id)
	}
	return s.WriteYAML("member.local.yaml", Identity{id}, 0600)
}

func (s *Store) Relate(from, to, description string) error {
	cfg, err := s.Config()
	if err != nil {
		return err
	}
	if from == to {
		return errors.New("a relationship needs two distinct repositories")
	}
	for _, id := range []string{from, to} {
		if _, ok := cfg.Repositories[id]; !ok {
			return fmt.Errorf("unknown repository: %s", id)
		}
	}
	if err := Text(description); err != nil {
		return err
	}
	relation := Relationship{from, to, description}
	for _, existing := range cfg.Relationships {
		if existing == relation {
			return nil
		}
	}
	return s.Update("workspace.yaml", []string{"relationships"}, append(cfg.Relationships, relation), 0644)
}
