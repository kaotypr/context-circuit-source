package workspace

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
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
func recordKind(prefix string) string {
	if prefix == "i" {
		return "an intent"
	}
	return "a plan"
}

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
		// An instant that no longer parses is reported rather than read as an event
		// that never happened. Each gate belongs to one kind of record: nothing
		// approves a plan, and an intent is never the thing that completes.
		if record.CreatedAt == "" {
			issues = append(issues, record.ID+": created_at is missing; a record written before 2.0.0-rc.4 is not readable by this candidate")
		} else if _, err := time.Parse(TimeLayout, string(record.CreatedAt)); err != nil {
			issues = append(issues, record.ID+": created_at is not a canonical ISO 8601 UTC timestamp: "+string(record.CreatedAt))
		}
		for _, gate := range []struct {
			field, prefix string
			value         ISOTime
		}{{"approved_at", "i", record.ApprovedAt}, {"completed_at", "p", record.CompletedAt}} {
			if value := string(gate.value); value == "" {
				continue
			} else if !strings.HasPrefix(record.ID, gate.prefix) {
				issues = append(issues, record.ID+": "+gate.field+" belongs to "+recordKind(gate.prefix))
			} else if _, err := time.Parse(TimeLayout, value); err != nil {
				issues = append(issues, record.ID+": "+gate.field+" is not a canonical ISO 8601 UTC timestamp: "+value)
			}
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
	knowledge, err := s.knowledgeIssues()
	if err != nil {
		issues = append(issues, err.Error())
	} else {
		issues = append(issues, knowledge...)
	}
	issues = unique(issues)
	sort.Strings(issues)
	return issues, nil
}

// Context notes describe the project, not the workspace machinery that produced
// them. A note naming a record or a raw evidence file rots once that ephemeral
// file is archived, and a recorded evidence path becomes a standing instruction
// to read material that must stay passive. Repository paths carry a logical
// repository ID and remain the durable anchor, exact or patterned.
var knowledgeRecord = regexp.MustCompile(`(^|[^A-Za-z0-9-])(i[0-9]{3,}|p[0-9]{4,})([^A-Za-z0-9-]|$)`)
var knowledgeMachinery = regexp.MustCompile(`(^|[^A-Za-z0-9@/._-])((intent|plans|sources)/|\.context-circuit/)`)

func (s *Store) knowledgeIssues() ([]string, error) {
	base, err := s.Path("context")
	if err != nil {
		return nil, err
	}
	var issues []string
	notes := map[string]string{}
	err = filepath.WalkDir(base, func(filename string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("context directories cannot contain symlinks: %s", filename)
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}
		relative, err := filepath.Rel(s.Root, filename)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		// The catalog is not a note, and a README is navigation rather than
		// knowledge; everything else is reachable only through the catalog.
		if entry.Name() != "INDEX.md" && entry.Name() != "README.md" {
			inside, err := filepath.Rel(base, filename)
			if err != nil {
				return err
			}
			notes[filepath.ToSlash(inside)] = relative
		}
		data, err := s.Read(relative)
		if err != nil {
			return err
		}
		for index, line := range strings.Split(string(data), "\n") {
			for _, pattern := range []*regexp.Regexp{knowledgeMachinery, knowledgeRecord} {
				if match := pattern.FindStringSubmatch(line); match != nil {
					issues = append(issues, fmt.Sprintf("%s:%d: knowledge note names workspace machinery (%s)", relative, index+1, match[2]))
				}
			}
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	catalog, err := s.catalogIssues(notes)
	if err != nil {
		return nil, err
	}
	return append(issues, catalog...), nil
}

var catalogLink = regexp.MustCompile(`\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)

// The catalog is the only way into a note, so an entry naming a note that is not
// there is a confident miss, and a note no entry names is reachable only by
// someone who already knows its filename. An absent catalog is not an error: the
// agent then falls back to filenames and search terms.
func (s *Store) catalogIssues(notes map[string]string) ([]string, error) {
	data, err := s.Read("context/INDEX.md")
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var issues []string
	referenced, fenced := map[string]bool{}, false
	for index, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			fenced = !fenced
			continue
		}
		// A fenced example teaches the entry shape; it catalogs nothing.
		if fenced || !catalogEntry.MatchString(line) {
			continue
		}
		for _, match := range catalogLink.FindAllStringSubmatch(line, -1) {
			target := match[1]
			if strings.HasPrefix(target, "#") || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
				continue
			}
			target, _, _ = strings.Cut(target, "#")
			target = path.Clean(strings.TrimPrefix(target, "./"))
			if target == "" || target == "." {
				continue
			}
			if strings.HasPrefix(target, "../") || strings.HasPrefix(target, "/") {
				issues = append(issues, fmt.Sprintf("context/INDEX.md:%d: catalog entry links outside the catalog (%s)", index+1, target))
				continue
			}
			referenced[target] = true
			if _, ok := notes[target]; !ok && target != "INDEX.md" {
				issues = append(issues, fmt.Sprintf("context/INDEX.md:%d: catalog entry links to a missing note (%s)", index+1, target))
			}
		}
	}
	for inside, relative := range notes {
		if !referenced[inside] {
			issues = append(issues, relative+": note is not listed in the catalog")
		}
	}
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
