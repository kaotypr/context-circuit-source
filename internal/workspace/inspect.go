package workspace

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)

type Orientation struct {
	Workspace    Config              `json:"workspace"`
	Members      Members             `json:"members"`
	ActiveMember string              `json:"active_member,omitempty"`
	Repositories map[string]Snapshot `json:"repositories"`
	Issues       []string            `json:"issues"`
}

func (s *Store) Status(ctx context.Context) (Orientation, error) {
	cfg, err := s.Config()
	if err != nil {
		return Orientation{}, err
	}
	members, err := s.Members()
	if err != nil {
		return Orientation{}, err
	}
	result := Orientation{Workspace: cfg, Members: members, Repositories: map[string]Snapshot{}, Issues: []string{}}
	result.ActiveMember, err = s.ActiveMember()
	if err != nil {
		result.Issues = append(result.Issues, err.Error())
	}
	for id := range cfg.Repositories {
		_, path, err := s.Repository(ctx, id)
		if err != nil {
			result.Issues = append(result.Issues, id+": "+err.Error())
			continue
		}
		snapshot, err := Inspect(ctx, path)
		if err != nil {
			result.Issues = append(result.Issues, id+": "+err.Error())
			continue
		}
		result.Repositories[id] = snapshot
	}
	sort.Strings(result.Issues)
	return result, nil
}

// Check is an explicitly invoked diagnostic, not an execution admission gate.
func (s *Store) Check(ctx context.Context) ([]string, error) {
	state, err := s.Status(ctx)
	if err != nil {
		return nil, err
	}
	issues := state.Issues
	records, err := s.ListRecords(false)
	if err != nil {
		return append(issues, err.Error()), nil
	}
	var ledger Ledger
	if err := s.YAML(".context-circuit/ids.yaml", &ledger); err != nil {
		issues = append(issues, err.Error())
	}
	reserved := map[string]bool{}
	for _, id := range append(append([]string{}, ledger.Intents...), ledger.Plans...) {
		if reserved[id] || !recordPattern.MatchString(id) {
			issues = append(issues, "invalid/duplicate reserved ID: "+id)
		}
		reserved[id] = true
	}
	paths, err := s.RecordPaths(true)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, path := range paths {
		id := strings.SplitN(strings.TrimPrefix(path[strings.LastIndex(path, "/")+1:], "/"), "-", 2)[0]
		if seen[id] {
			issues = append(issues, "duplicate record ID: "+id)
		}
		seen[id] = true
		if !reserved[id] {
			issues = append(issues, "record ID missing from permanent ledger: "+id)
		}
	}
	for _, record := range records {
		if _, ok := state.Members.Members[record.CreatedBy]; !ok {
			issues = append(issues, record.ID+": unknown created_by member")
		}
		if strings.HasPrefix(record.ID, "p") {
			parent, err := s.FindRecord(record.Intent)
			if err != nil || !strings.HasPrefix(record.Intent, "i") {
				issues = append(issues, record.ID+": missing intent reference")
			} else if !slices.Contains(parent.Plans, record.ID) {
				issues = append(issues, record.ID+": not linked from intent "+parent.ID)
			}
			if len(record.Repositories) == 0 {
				issues = append(issues, record.ID+": no repositories")
			}
			for _, repo := range record.Repositories {
				if _, ok := state.Workspace.Repositories[repo]; !ok {
					issues = append(issues, record.ID+": unknown repository "+repo)
				}
			}
		} else {
			for _, id := range record.Plans {
				plan, err := s.FindRecord(id)
				if err != nil || plan.Intent != record.ID {
					issues = append(issues, record.ID+": broken plan link "+id)
				}
			}
		}
	}
	var visit func(string, map[string]bool) error
	visit = func(id string, trail map[string]bool) error {
		if trail[id] {
			return fmt.Errorf("dependency cycle at %s", id)
		}
		if !strings.HasPrefix(id, "p") {
			return fmt.Errorf("invalid plan dependency: %s", id)
		}
		record, err := s.FindRecord(id)
		if err != nil {
			return err
		}
		trail[id] = true
		defer delete(trail, id)
		for _, dep := range record.DependsOn {
			if err := visit(dep, trail); err != nil {
				return err
			}
		}
		return nil
	}
	for _, record := range records {
		if strings.HasPrefix(record.ID, "p") {
			if err := visit(record.ID, map[string]bool{}); err != nil {
				issues = append(issues, err.Error())
			}
		}
	}
	for _, relation := range state.Workspace.Relationships {
		for _, id := range []string{relation.From, relation.To} {
			if _, ok := state.Workspace.Repositories[id]; !ok {
				issues = append(issues, "relationship references unknown repository: "+id)
			}
		}
	}
	local, err := s.associations()
	if err != nil {
		issues = append(issues, err.Error())
	} else {
		for _, assoc := range local.Worktrees {
			_, tree, err := s.selectedTree(ctx, assoc.Repository, assoc.Path)
			if err != nil || tree.Prunable || tree.Branch != assoc.Branch {
				issues = append(issues, "worktree association needs attention: "+assoc.Path)
			}
			if assoc.Plan != "" {
				if _, err := s.FindRecord(assoc.Plan); err != nil {
					issues = append(issues, "worktree has missing plan: "+assoc.Plan)
				}
			}
		}
	}
	issues = unique(issues)
	sort.Strings(issues)
	return issues, nil
}

func (s *Store) FindContext(query string) ([]string, error) {
	data, err := s.Read("context/INDEX.md")
	if os.IsNotExist(err) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" && strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			result = append(result, line)
		}
	}
	return result, nil
}
