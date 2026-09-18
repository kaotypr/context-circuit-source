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
	Workspace Config `json:"workspace"`
	// Self is the Git state of the workspace's own repository, present only
	// once it has been described. Shared records are read from whatever commit
	// this checkout is sitting on, so a stale or dirty workspace is part of
	// what orientation has to report rather than something to go looking for.
	Self         *Snapshot           `json:"workspace_repository,omitempty"`
	Members      Members             `json:"members"`
	ActiveMember string              `json:"active_member,omitempty"`
	Repositories map[string]Snapshot `json:"repositories"`
	// Knowledge is every borrowed repository and what this machine currently
	// sees of it. It is reported beside the repositories, and separately from
	// them, because reading it is part of orienting and working in it is not.
	Knowledge []KnowledgeState `json:"knowledge_repositories,omitempty"`
	Issues    []string         `json:"issues"`
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
	if checkout, described, err := s.WorkspaceCheckout(ctx); err != nil {
		result.Issues = append(result.Issues, "workspace repository: "+err.Error())
	} else if described {
		snapshot, err := Inspect(ctx, checkout.Path)
		if err != nil {
			result.Issues = append(result.Issues, "workspace repository: "+err.Error())
		} else {
			snapshot.BaseBranch = checkout.BaseBranch
			result.Self = &snapshot
		}
	}
	for id := range cfg.Repositories {
		checkout, err := s.Repository(ctx, id)
		if err != nil {
			result.Issues = append(result.Issues, id+": "+err.Error())
			continue
		}
		snapshot, err := Inspect(ctx, checkout.Path)
		if err != nil {
			result.Issues = append(result.Issues, id+": "+err.Error())
			continue
		}
		snapshot.BaseBranch = checkout.BaseBranch
		result.Repositories[id] = snapshot
	}
	if knowledge, err := s.KnowledgeStates(ctx); err != nil {
		result.Issues = append(result.Issues, err.Error())
	} else {
		result.Knowledge = knowledge
	}
	sort.Strings(result.Issues)
	return result, nil
}

// workspaceRepositoryIssues reports what the workspace's own Git state costs
// when it is left undescribed. Shared records travel through that repository,
// so a workspace nobody can locate is one a second machine cannot be told how
// to obtain, and an undescribed one is invisible to status.
func (s *Store) workspaceRepositoryIssues(ctx context.Context, state Orientation) []string {
	if state.Workspace.WorkspaceRepository != nil {
		return nil
	}
	if _, err := rootOf(ctx, s.Root); err != nil {
		// A workspace that is not under version control has nothing to
		// describe. Sharing it is a choice, not an omission.
		return nil
	}
	return []string{"the workspace is a Git checkout but describes no workspace repository; record it with workspace connect so status reports its branch and a second machine knows where to clone it"}
}

// unbandedMembers reports members allocating from the shared range in a
// workspace that has more than one. Bands are what let clones that cannot see
// each other allocate without colliding; a solo workspace needs none, so the
// finding waits until a second member makes collision possible.
func unbandedMembers(members Members) []string {
	if len(members.Members) < 2 {
		return nil
	}
	var unbanded []string
	for id, member := range members.Members {
		if member.Band == 0 {
			unbanded = append(unbanded, id)
		}
	}
	if len(unbanded) == 0 {
		return nil
	}
	sort.Strings(unbanded)
	return []string{fmt.Sprintf("members allocating from the shared range in a workspace of %d: %s; give each a distinct band with member band before they work apart", len(members.Members), strings.Join(unbanded, ", "))}
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
	issues = append(issues, s.workspaceRepositoryIssues(ctx, state)...)
	issues = append(issues, unbandedMembers(state.Members)...)
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
	issues = append(issues, borrowedIssues(state.Knowledge)...)
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
		// The glossary is a table of terms rather than a note, and is the one
		// place an anchor belongs outside an Owner block.
		if entry.Name() != "INDEX.md" && entry.Name() != "README.md" && entry.Name() != "glossary.md" {
			issues = append(issues, noteReadability(relative, data)...)
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

// Groups carry entries; the document's own title does not.
var catalogHeading = regexp.MustCompile(`^#{2,6}\s+\S`)

// An entry naming repositories claims to describe their code, so it carries the
// date that claim was last confirmed. The glossary and anything else describing
// the project rather than a repository names none and needs no date.
var catalogRepositories = regexp.MustCompile(`\{[a-z0-9][a-z0-9-]*(?:\s*,\s*[a-z0-9][a-z0-9-]*)*\}`)
var catalogReviewed = regexp.MustCompile(`reviewed\s+\d{4}-\d{2}-\d{2}\s*$`)

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
	heading, headingLine, scoped := "", 0, false
	concerns, flagged := map[string]string{}, map[string]bool{}
	for index, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			fenced = !fenced
			continue
		}
		if !fenced && catalogHeading.MatchString(trimmed) {
			heading, headingLine, scoped = strings.TrimLeft(trimmed, "# "), index, false
			continue
		}
		// A fenced example teaches the entry shape; it catalogs nothing.
		if fenced || !catalogEntry.MatchString(line) {
			// A group says what belongs under it before its first entry, so the
			// next note has somewhere obvious to go rather than a guess.
			if !fenced && trimmed != "" {
				scoped = true
			}
			continue
		}
		if heading != "" && !scoped && !flagged[heading] {
			flagged[heading] = true
			issues = append(issues, fmt.Sprintf("context/INDEX.md:%d: heading has no line saying what belongs under it (%s)", headingLine+1, heading))
		}
		if catalogRepositories.MatchString(line) && !catalogReviewed.MatchString(strings.TrimRight(line, " \t")) {
			issues = append(issues, fmt.Sprintf("context/INDEX.md:%d: catalog entry names repositories without a `reviewed YYYY-MM-DD` date", index+1))
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
			// One directory, one heading: notes split across two groups leave a
			// reader guessing which of them a new note belongs under.
			if concern, _, nested := strings.Cut(target, "/"); nested && heading != "" {
				if other, seen := concerns[concern]; seen && other != heading {
					if key := concern + "/"; !flagged[key] {
						flagged[key] = true
						issues = append(issues, fmt.Sprintf("context/INDEX.md:%d: notes in %s/ are catalogued under two headings (%s and %s)", index+1, concern, other, heading))
					}
				} else if !seen {
					concerns[concern] = heading
				}
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

// ContextMatch is one line of a knowledge index that matched. Borrowed says the
// line came from a repository this workspace reads but does not own, which is
// the difference between an entry to edit and one to raise upstream.
type ContextMatch struct {
	Source   string `json:"source"`
	Borrowed bool   `json:"borrowed,omitempty"`
	Entry    string `json:"entry"`
}

// ContextSearch carries what was found and what could not be read. A borrowed
// index this machine has not obtained would otherwise narrow the result
// silently, and a search that found nothing reads identically to one that never
// looked — which is the confident miss the catalog exists to prevent.
type ContextSearch struct {
	Matches    []ContextMatch `json:"matches"`
	Unsearched []string       `json:"unsearched,omitempty"`
}

func matchLines(source string, borrowed bool, data []byte, query string) []ContextMatch {
	var found []ContextMatch
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" && strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			found = append(found, ContextMatch{Source: source, Borrowed: borrowed, Entry: line})
		}
	}
	return found
}

// FindContext searches this workspace's catalog and every borrowed knowledge
// index. Matching stays plain case-insensitive substring over whole lines,
// because a borrowed repository writes its index its own way and parsing it as
// this product's catalog would find nothing in a file that is perfectly good.
func (s *Store) FindContext(ctx context.Context, query string) (ContextSearch, error) {
	result := ContextSearch{Matches: []ContextMatch{}}
	data, err := s.Read("context/INDEX.md")
	if err != nil && !os.IsNotExist(err) {
		return result, err
	}
	if err == nil {
		result.Matches = append(result.Matches, matchLines("context/INDEX.md", false, data, query)...)
	}
	cfg, err := s.Config()
	if err != nil {
		return result, err
	}
	ids := make([]string, 0, len(cfg.KnowledgeRepositories))
	for id := range cfg.KnowledgeRepositories {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		index, data, err := s.KnowledgeIndex(ctx, id)
		if err != nil {
			result.Unsearched = append(result.Unsearched, id+": "+err.Error())
			continue
		}
		result.Matches = append(result.Matches, matchLines(id+"/"+index, true, data, query)...)
	}
	return result, nil
}
