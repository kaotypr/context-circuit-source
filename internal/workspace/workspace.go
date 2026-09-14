package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Repository struct {
	BaseBranch string `yaml:"base_branch" json:"base_branch"`
}
type Relationship struct {
	From        string `yaml:"from" json:"from"`
	To          string `yaml:"to" json:"to"`
	Description string `yaml:"description" json:"description"`
}
type Config struct {
	Version       int                   `yaml:"version" json:"version"`
	Name          string                `yaml:"name" json:"name"`
	Purpose       string                `yaml:"purpose" json:"purpose"`
	Repositories  map[string]Repository `yaml:"repositories" json:"repositories"`
	Relationships []Relationship        `yaml:"relationships" json:"relationships"`
}
type Member struct {
	Name string `yaml:"name" json:"name"`
}
type Members struct {
	Members map[string]Member `yaml:"members" json:"members"`
}
type Identity struct {
	Member string `yaml:"member" json:"member"`
}
type Binding struct {
	Path string `yaml:"path" json:"path"`
}
type Bindings struct {
	Bindings map[string]Binding `yaml:"bindings" json:"bindings"`
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
		return cfg, err
	}
	if cfg.Version != 2 || cfg.Name == "" {
		return cfg, errors.New("initialize a fresh v2 workspace first; v1 migration is not automatic")
	}
	if cfg.Repositories == nil || cfg.Relationships == nil {
		return cfg, errors.New("workspace.yaml needs repositories and relationships collections")
	}
	for id, repo := range cfg.Repositories {
		if err := Name(id); err != nil {
			return cfg, err
		}
		if err := Text(repo.BaseBranch); err != nil {
			return cfg, err
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
	for id, member := range members.Members {
		if err := Name(id); err != nil {
			return members, err
		}
		if err := Text(member.Name); err != nil {
			return members, err
		}
	}
	return members, nil
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
	for rel, seed := range files {
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
	cfg := Config{2, name, purpose, map[string]Repository{}, []Relationship{}}
	if err := s.Export(missing); err != nil {
		return err
	}
	if err := s.WriteYAML("workspace.yaml", cfg, 0644); err != nil {
		return err
	}
	if err := s.WriteYAML("members.yaml", Members{map[string]Member{member: {display}}}, 0644); err != nil {
		return err
	}
	if err := s.WriteYAML("repositories.local.yaml", Bindings{map[string]Binding{}}, 0600); err != nil {
		return err
	}
	return s.WriteYAML("member.local.yaml", Identity{member}, 0600)
}

func (s *Store) AddMember(id, name string) error {
	if _, err := s.Config(); err != nil {
		return err
	}
	if err := Name(id); err != nil {
		return err
	}
	if err := Text(name); err != nil {
		return err
	}
	members, err := s.Members()
	if err != nil {
		return err
	}
	if existing, ok := members.Members[id]; ok {
		if existing.Name == name {
			return nil
		}
		return fmt.Errorf("member %s already has another name", id)
	}
	return s.Update("members.yaml", []string{"members", id}, Member{name}, 0644)
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
