package workspace

import (
	"context"
	"errors"
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

// A Finding is one thing wrong and what discharges it. The issue alone was what
// this diagnostic reported for its whole life, and it left the caller to derive
// the remedy in the one situation where the workspace is already inconsistent
// and deriving is most expensive. Every other command here hands back the next
// move — `record approve` returns planning_recommended, `worktree prepare` names
// the branch to pass — and this was the exception.
//
// Resolve is an instruction, not a category: it names the command where one
// exists, the edit where an edit is the whole of it, and it opens with
// `needs a person:` where neither is true and the next step is somebody's
// judgment. That last case is the one worth marking, because an agent that
// cannot tell it apart improvises against a broken workspace.
type Finding struct {
	Issue   string `json:"issue"`
	Resolve string `json:"resolve"`
}

// needsAPerson opens a resolution that no command and no edit discharges.
const needsAPerson = "needs a person: "

func found(issue, resolve string) Finding { return Finding{issue, resolve} }

// wrap gives a shared resolution to findings produced in bulk, where the issue
// text differs and the remedy does not.
func wrap(issues []string, resolve string) []Finding {
	out := make([]Finding, 0, len(issues))
	for _, issue := range issues {
		out = append(out, Finding{issue, resolve})
	}
	return out
}

func uniqueFindings(values []Finding) []Finding {
	out, seen := []Finding{}, map[string]bool{}
	for _, value := range values {
		if !seen[value.Issue] {
			seen[value.Issue] = true
			out = append(out, value)
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].Issue < out[b].Issue })
	return out
}

type Orientation struct {
	Workspace Config `json:"workspace"`
	// Self is the Git state of the workspace's own repository, present only
	// once it has been described. Shared records are read from whatever commit
	// this checkout is sitting on, so a stale or dirty workspace is part of
	// what orientation has to report rather than something to go looking for.
	Self         *Snapshot `json:"workspace_repository,omitempty"`
	Members      Members   `json:"members"`
	ActiveMember string    `json:"active_member,omitempty"`
	// IdentityRequired names the outstanding act when this machine has selected
	// no member of the roster, including when the roster lists nobody or cannot
	// be read. It is said beside the data, the way approval says
	// planning_recommended, because an empty active_member alone does not tell
	// an agent what to do next.
	IdentityRequired string `json:"identity_required,omitempty"`
	// BandRequired names the outstanding act when a member holds no band, which
	// only a roster written before every member needed one can have.
	// BandSuggestions offers a band for each such member, and NextBand the one
	// for a member added next; the person confirms each or picks another.
	BandRequired    string              `json:"band_required,omitempty"`
	BandSuggestions map[string]int      `json:"band_suggestions,omitempty"`
	NextBand        int                 `json:"next_band,omitempty"`
	Repositories    map[string]Snapshot `json:"repositories"`
	// Knowledge is every borrowed repository and what this machine currently
	// sees of it. It is reported beside the repositories, and separately from
	// them, because reading it is part of orienting and working in it is not.
	Knowledge []KnowledgeState `json:"knowledge_repositories,omitempty"`
	Issues    []Finding        `json:"issues"`
}

func (s *Store) Status(ctx context.Context) (Orientation, error) {
	cfg, err := s.Config()
	if err != nil {
		return Orientation{}, err
	}
	result := Orientation{Workspace: cfg, Members: Members{Members: map[string]Member{}}, Repositories: map[string]Snapshot{}, Issues: []Finding{}}
	// An unreadable roster is reported, not fatal: two clones that both changed
	// it are exactly when a person most needs to see why they are being asked
	// to repair it.
	members, err := s.Members()
	if err != nil {
		result.Issues = append(result.Issues, found("members.yaml: "+err.Error(),
			needsAPerson+"repair members.yaml by hand so each member appears once, holds one band, and shares it with nobody; when two clones clashed, choose which of their edits stands"))
		result.IdentityRequired = identityRequired
	} else {
		result.Members = members
		result.ActiveMember, err = s.activeMember(members)
		if err != nil {
			result.Issues = append(result.Issues, found(err.Error(),
				"select this machine's member with `member use --id ID`; add them to the roster first with `member add --id ID --name NAME --band N` if nobody has"))
			if errors.Is(err, ErrIdentityRequired) {
				result.IdentityRequired = identityRequired
			}
		}
		if suggestions, err := s.SuggestBands(); err != nil {
			result.Issues = append(result.Issues, found(err.Error(),
				needsAPerson+"bands cannot be suggested until the ID ledger and record files can be read; `check` names what is wrong"))
		} else {
			result.NextBand = suggestions.Next
			if len(suggestions.Members) > 0 {
				result.BandSuggestions = suggestions.Members
			}
		}
		if unbanded := unbandedMembers(members, result.BandSuggestions); len(unbanded) > 0 {
			result.Issues = append(result.Issues, unbanded...)
			result.BandRequired = "decline any request that changes the workspace, its repositories, or its records until every member holds a band: offer each member in band_suggestions that band, let the person confirm it or pick another unused one, record it with `member band --id ID --band N`, and continue the request they made"
		}
		if cfg.Version < CurrentSchema {
			result.Issues = append(result.Issues, schema2Roles(members)...)
		}
		result.Issues = append(result.Issues, deviceNameMembers(members)...)
	}
	if checkout, described, err := s.WorkspaceCheckout(ctx); err != nil {
		result.Issues = append(result.Issues, found("workspace repository: "+err.Error(),
			"bind this machine's checkout of the workspace with `workspace connect --base BRANCH`"))
	} else if described {
		snapshot, err := Inspect(ctx, checkout.Path)
		if err != nil {
			result.Issues = append(result.Issues, found("workspace repository: "+err.Error(),
				needsAPerson+"the recorded workspace checkout cannot be read as a Git repository; repair or re-point it, then run `workspace connect` again"))
		} else {
			snapshot.BaseBranch = checkout.BaseBranch
			result.Self = &snapshot
		}
	}
	for id := range cfg.Repositories {
		checkout, err := s.Repository(ctx, id)
		if err != nil {
			result.Issues = append(result.Issues, found(id+": "+err.Error(),
				"obtain it with `repo clone --id "+id+" --path repositories/"+id+" --base BRANCH`, or bind a checkout you already have with `repo connect`"))
			continue
		}
		snapshot, err := Inspect(ctx, checkout.Path)
		if err != nil {
			result.Issues = append(result.Issues, found(id+": "+err.Error(),
				needsAPerson+"the bound path cannot be read as a Git repository; restore it, or re-point the binding with `repo connect`"))
			continue
		}
		snapshot.BaseBranch = checkout.BaseBranch
		result.Repositories[id] = snapshot
	}
	if knowledge, err := s.KnowledgeStates(ctx); err != nil {
		result.Issues = append(result.Issues, found(err.Error(),
			needsAPerson+"the knowledge records could not be read; inspect workspace.yaml and repositories.local.yaml"))
	} else {
		result.Knowledge = knowledge
	}
	result.Issues = uniqueFindings(result.Issues)
	return result, nil
}

// workspaceRepositoryIssues reports what the workspace's own Git state costs
// when it is left undescribed. Shared records travel through that repository,
// so a workspace nobody can locate is one a second machine cannot be told how
// to obtain, and an undescribed one is invisible to status.
func (s *Store) workspaceRepositoryIssues(ctx context.Context, state Orientation) []Finding {
	if state.Workspace.WorkspaceRepository != nil {
		return nil
	}
	if _, err := rootOf(ctx, s.Root); err != nil {
		// A workspace that is not under version control has nothing to
		// describe. Sharing it is a choice, not an omission.
		return nil
	}
	return []Finding{{"the workspace is a Git checkout but describes no workspace repository, so status cannot report its branch and a second machine is not told where to clone it",
		"`workspace connect --base BRANCH`, which reads the URL from this checkout's origin"}}
}

// identityRequired is the instruction status gives while this machine has no
// member of a readable roster selected.
const identityRequired = "decline any request that changes the workspace, its repositories, or its records until the person says which member they are: select them with `member use --id ID`, or add them with `member add --id ID --name NAME --band N` and then select them, and continue the request they made"

// unbandedMembers reports each member holding no band, with the band it would
// be offered. Every member holds one, and changes are refused until they do.
func unbandedMembers(members Members, suggestions map[string]int) []Finding {
	var out []Finding
	for _, id := range members.IDs() {
		if members.Members[id].Band != nil {
			continue
		}
		resolve := "give it a band nobody holds with `member band --id " + id + " --band N`"
		if band, ok := suggestions[id]; ok {
			resolve = fmt.Sprintf("give it a band with `member band --id %s --band %d`, or another band nobody holds that the person picks", id, band)
		}
		out = append(out, found("member "+id+" holds no band", resolve))
	}
	return out
}

// schema2Roles reports roles in a schema 2 roster, which only a hand edit can
// write. A 2.3.0 CLI still reads that workspace and stops at the unknown field
// instead of asking its user to upgrade, which is why roles wait for schema 3.
func schema2Roles(members Members) []Finding {
	var out []Finding
	for _, id := range members.IDs() {
		if len(members.Members[id].Roles) == 0 {
			continue
		}
		out = append(out, found("member "+id+" records roles in a schema 2 workspace, which a 2.3.0 CLI cannot read",
			"remove the roles from members.yaml by hand, or finish setup so the next changing command moves the workspace to schema 3"))
	}
	return out
}

// deviceNameMembers reports roster IDs that are Windows device names. A 2.3.0
// roster could hold one, and no clone on Windows can check out the ID file it
// would name, so the split refuses it and every change waits on a rename.
func deviceNameMembers(members Members) []Finding {
	var out []Finding
	for _, id := range members.IDs() {
		if windowsDevice.MatchString(id) {
			out = append(out, found("member "+id+" is a Windows device name, which cannot name an ID file",
				needsAPerson+renameMemberByHand))
		}
	}
	return out
}

// Check is an explicitly invoked diagnostic, not an execution admission gate.
func recordKind(prefix string) string {
	if prefix == "i" {
		return "an intent"
	}
	return "a plan"
}

func (s *Store) Check(ctx context.Context) ([]Finding, error) {
	state, err := s.Status(ctx)
	if err != nil {
		return nil, err
	}
	issues := state.Issues
	issues = append(issues, s.workspaceRepositoryIssues(ctx, state)...)
	for _, baseRel := range []string{"plans", "plans/archive"} {
		base, err := s.Path(baseRel)
		if err != nil {
			continue
		}
		entries, err := os.ReadDir(base)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() || !recordFilename.MatchString(entry.Name()+".md") {
				continue
			}
			rel := filepath.ToSlash(filepath.Join(baseRel, entry.Name(), "plan.md"))
			if _, err := s.Read(rel); err != nil {
				// An interrupted record delete can leave a plan folder emptied of
				// everything; there is nothing in it to restore.
				if left, err := os.ReadDir(filepath.Join(base, entry.Name())); err == nil && len(left) == 0 {
					issues = append(issues, found("empty plan folder: "+filepath.ToSlash(filepath.Join(baseRel, entry.Name())),
						"remove the empty folder, which an interrupted `record delete` leaves behind; restore its plan.md instead only if the plan was meant to stay"))
					continue
				}
				issues = append(issues, found("missing or invalid plan entry: "+rel, "restore a regular plan.md in this plan folder"))
			}
		}
	}
	records, err := s.ListRecords(false)
	if err != nil {
		return append(issues, found(err.Error(),
			needsAPerson+"the records could not be listed; inspect intent/ and plans/ for an unreadable file")), nil
	}
	issues, reserved, err := s.checkIDFiles(issues)
	if err != nil {
		return nil, err
	}
	paths, err := s.RecordPaths(true)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, path := range paths {
		name := filepath.Base(path)
		if name == "plan.md" {
			name = filepath.Base(filepath.Dir(path))
		}
		id := strings.SplitN(name, "-", 2)[0]
		if seen[id] {
			issues = append(issues, found("duplicate record ID: "+id,
				needsAPerson+"two clones allocated the same number; renumber one record and every reference to it before either is shared, and keep both reservations"))
		}
		seen[id] = true
		if !reserved[id] {
			owner := "<member>"
			if record, err := s.readRecord(path); err == nil && record.CreatedBy != "" {
				owner = record.CreatedBy
			} else if active, err := s.ActiveMember(); err == nil {
				owner = active
			}
			issues = append(issues, found("record ID reserved in no ID file: "+id,
				needsAPerson+"add the ID to "+IDFilesDir+"/"+owner+".yaml, the file of the record's created_by or, when it has none, of the local member, so nothing allocates it again"))
		}
	}
	// An unreadable roster is already reported; reading every author as unknown
	// against it would bury that finding under ones that are not true.
	_, rosterErr := s.Members()
	for _, record := range records {
		// A record with no author was written by an earlier version while the
		// workspace had no roster, and stays valid.
		if _, ok := state.Members.Members[record.CreatedBy]; !ok && record.CreatedBy != "" && rosterErr == nil {
			issues = append(issues, found(record.ID+": unknown created_by member",
				"add them with `member add --id ID --name NAME --band N`, or correct the record's created_by to a member the roster carries"))
		}
		// A roster member with such an ID is reported by status; an author the
		// roster does not carry still stops the split.
		if _, ok := state.Members.Members[record.CreatedBy]; !ok && windowsDevice.MatchString(record.CreatedBy) {
			issues = append(issues, found(record.ID+": created_by "+record.CreatedBy+" is a Windows device name, which cannot name an ID file",
				needsAPerson+"correct the record's created_by to a member the roster carries whose ID is not a Windows device name"))
		}
		// An instant that no longer parses is reported rather than read as an event
		// that never happened. Each gate belongs to one kind of record: nothing
		// approves a plan, and an intent is never the thing that completes.
		if record.CreatedAt == "" {
			issues = append(issues, found(record.ID+": created_at is missing; a record written before 2.0.0-rc.4 is not readable by this CLI",
				needsAPerson+"add the instant this record was created to its frontmatter; a record from an earlier version is not migrated for you"))
		} else if _, err := time.Parse(TimeLayout, string(record.CreatedAt)); err != nil {
			issues = append(issues, found(record.ID+": created_at is not a canonical ISO 8601 UTC timestamp: "+string(record.CreatedAt),
				needsAPerson+"rewrite it in frontmatter as 2026-09-15T10:53:00Z; only a person knows which instant was meant"))
		}
		for _, gate := range []struct {
			field, prefix string
			value         ISOTime
		}{{"approved_at", "i", record.ApprovedAt}, {"completed_at", "p", record.CompletedAt}} {
			if value := string(gate.value); value == "" {
				continue
			} else if !strings.HasPrefix(record.ID, gate.prefix) {
				issues = append(issues, found(record.ID+": "+gate.field+" belongs to "+recordKind(gate.prefix),
					needsAPerson+"remove the field from this record; a gate stamped on the wrong kind of record claims a decision that was never made"))
			} else if _, err := time.Parse(TimeLayout, value); err != nil {
				issues = append(issues, found(record.ID+": "+gate.field+" is not a canonical ISO 8601 UTC timestamp: "+value,
					needsAPerson+"rewrite it in frontmatter as 2026-09-15T10:53:00Z, or remove it if the decision it claims was never taken"))
			}
		}
		if strings.HasPrefix(record.ID, "p") {
			if !folderPlan(record) && record.RequiredFiles != nil {
				issues = append(issues, found(record.ID+": legacy plan cannot assign required_files", "remove required_files from the single-file plan"))
			}
			if record.Intent != "" {
				parent, err := s.FindRecord(record.Intent)
				if err != nil || !strings.HasPrefix(record.Intent, "i") {
					issues = append(issues, found(record.ID+": invalid intent reference",
						needsAPerson+"correct the intent ID and its backlink, or remove the field for a standalone plan"))
				} else if !slices.Contains(parent.Plans, record.ID) {
					issues = append(issues, found(record.ID+": not linked from intent "+parent.ID,
						needsAPerson+"add this plan to that intent's plans; the link is written at creation and nothing relinks it afterwards"))
				}
			}
			if folderPlan(record) {
				full, err := s.readRecord(record.Path)
				if err == nil {
					issues = append(issues, s.planFileIssues(full)...)
				}
			}
			if len(record.Repositories) == 0 {
				issues = append(issues, found(record.ID+": no repositories",
					needsAPerson+"name the repositories this plan changes in its frontmatter; a plan with none cannot be prepared or ordered"))
			}
			for _, repo := range record.Repositories {
				if _, ok := state.Workspace.Repositories[repo]; !ok {
					issues = append(issues, found(record.ID+": unknown repository "+repo,
						"register it with `repo connect --id "+repo+" --path PATH --base BRANCH`, or correct the plan to name a repository the workspace describes"))
				}
			}
		} else {
			for _, id := range record.Plans {
				plan, err := s.FindRecord(id)
				if err != nil || plan.Intent != record.ID {
					issues = append(issues, found(record.ID+": broken plan link "+id,
						needsAPerson+"the linked plan is not there; correct the link, or restore the record if it was moved"))
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
				issues = append(issues, found(err.Error(),
					"break the cycle with `record dependencies --id PLAN --depends-on ID ...`, naming only the plans that truly precede it; omit --depends-on to clear them"))
			}
		}
	}
	for _, relation := range state.Workspace.Relationships {
		for _, id := range []string{relation.From, relation.To} {
			if _, ok := state.Workspace.Repositories[id]; !ok {
				issues = append(issues, found("relationship references unknown repository: "+id,
					"register it with `repo connect --id "+id+" --path PATH --base BRANCH`, or remove the relationship from workspace.yaml"))
			}
		}
	}
	local, err := s.associations()
	if err != nil {
		issues = append(issues, found(err.Error(),
			needsAPerson+"the plan-to-worktree associations cannot be read; they are a convenience for resuming work, so removing .context-circuit/local/worktrees.yaml is safe if it is beyond repair"))
	} else {
		for _, assoc := range local.Worktrees {
			_, tree, err := s.selectedTree(ctx, assoc.Repository, assoc.Path)
			if errors.Is(err, errNotWorktree) {
				issues = append(issues, found("worktree association needs attention: "+assoc.Path,
					"Git no longer lists this worktree: `worktree prune --repo "+assoc.Repository+"` drops the associations of worktrees that are gone and touches no worktree, branch, or file; use `worktree repair --repo "+assoc.Repository+" --path PATH` instead if the checkout only moved"))
			} else if err != nil || tree.Prunable || tree.Branch != assoc.Branch {
				issues = append(issues, found("worktree association needs attention: "+assoc.Path,
					"`worktree repair --repo "+assoc.Repository+" --path "+assoc.Path+"` where the checkout moved; Git's own inventory is authoritative, and this association is only a convenience"))
			}
			if assoc.Plan != "" {
				if plan, err := s.FindRecord(assoc.Plan); err != nil {
					issues = append(issues, found("worktree has missing plan: "+assoc.Plan,
						needsAPerson+"the plan this worktree was prepared for is gone; the worktree and its branch are untouched, and removing one is a separate authorized request"))
				} else if assoc.PlanCreatedAt != "" && assoc.PlanCreatedAt != plan.CreatedAt {
					// The number was released and handed out again: the worktree
					// holds the earlier plan's work, which this one must not meet.
					issues = append(issues, found("worktree "+assoc.Path+" was prepared for an earlier "+assoc.Plan+", created at "+string(assoc.PlanCreatedAt)+", not the "+assoc.Plan+" created at "+string(plan.CreatedAt),
						needsAPerson+"remove it with `worktree remove --repo "+assoc.Repository+" --path \""+assoc.Path+"\"` (dirty work is kept unless discard is explicitly authorized) and delete branch "+assoc.Branch+" once its work no longer matters, or run `worktree prune --repo "+assoc.Repository+"` if Git no longer lists it; removing either is a separate authorized request"))
				}
			}
		}
	}
	knowledge, err := s.knowledgeIssues()
	if err != nil {
		issues = append(issues, found(err.Error(),
			needsAPerson+"the context tree could not be walked; inspect context/ for an unreadable file or a symlink"))
	} else {
		issues = append(issues, knowledge...)
	}
	issues = append(issues, borrowedIssues(state.Knowledge)...)
	issues = append(issues, s.reviewIssues(ctx)...)
	return uniqueFindings(issues), nil
}

// Context notes describe the project, not the workspace machinery that produced
// them. A note naming a record or a raw evidence file rots once that ephemeral
// file is archived, and a recorded evidence path becomes a standing instruction
// to read material that must stay passive. Repository paths carry a logical
// repository ID and remain the durable anchor, exact or patterned.
var knowledgeRecord = regexp.MustCompile(`(^|[^A-Za-z0-9-])(i[0-9]{3,}|p[0-9]{4,})([^A-Za-z0-9-]|$)`)
var knowledgeMachinery = regexp.MustCompile(`(^|[^A-Za-z0-9@/._-])((intent|plans|sources)/|\.context-circuit/)`)

func (s *Store) knowledgeIssues() ([]Finding, error) {
	base, err := s.Path("context")
	if err != nil {
		return nil, err
	}
	var issues []Finding
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
			issues = append(issues, wrap(noteReadability(relative, data),
				"reshape the note as `.agents/skills/cc-knowledge/SKILL.md` describes; the finding names the line and what it is carrying")...)
		}
		for index, line := range strings.Split(string(data), "\n") {
			for _, pattern := range []*regexp.Regexp{knowledgeMachinery, knowledgeRecord} {
				if match := pattern.FindStringSubmatch(line); match != nil {
					issues = append(issues, found(fmt.Sprintf("%s:%d: knowledge note names workspace machinery (%s)", relative, index+1, match[2]),
						"remove the reference; a note describes the project and never a record or a file of raw evidence, and which record produced it belongs in that plan"))
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
func (s *Store) catalogIssues(notes map[string]string) ([]Finding, error) {
	data, err := s.Read("context/INDEX.md")
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var issues []Finding
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
			issues = append(issues, found(fmt.Sprintf("context/INDEX.md:%d: heading has no line saying what belongs under it (%s)", headingLine+1, heading),
				"write one line under the heading saying what belongs there, before its first entry, naming the neighbouring concern where two are easily confused"))
		}
		if catalogRepositories.MatchString(line) && !catalogReviewed.MatchString(strings.TrimRight(line, " \t")) {
			issues = append(issues, found(fmt.Sprintf("context/INDEX.md:%d: catalog entry names repositories without a `reviewed YYYY-MM-DD` date", index+1),
				"end the entry with `· reviewed YYYY-MM-DD`, the date the note was last confirmed against the code"))
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
				issues = append(issues, found(fmt.Sprintf("context/INDEX.md:%d: catalog entry links outside the catalog (%s)", index+1, target),
					"point the entry at a note inside context/, or drop the entry; the catalog lists this workspace's own notes and nothing else"))
				continue
			}
			referenced[target] = true
			if _, ok := notes[target]; !ok && target != "INDEX.md" {
				issues = append(issues, found(fmt.Sprintf("context/INDEX.md:%d: catalog entry links to a missing note (%s)", index+1, target),
					"write the note, or correct the link; an entry naming a note that is not there is a confident miss"))
			}
			// One directory, one heading: notes split across two groups leave a
			// reader guessing which of them a new note belongs under.
			if concern, _, nested := strings.Cut(target, "/"); nested && heading != "" {
				if other, seen := concerns[concern]; seen && other != heading {
					if key := concern + "/"; !flagged[key] {
						flagged[key] = true
						issues = append(issues, found(fmt.Sprintf("context/INDEX.md:%d: notes in %s/ are catalogued under two headings (%s and %s)", index+1, concern, other, heading),
							"catalogue that directory's notes under one heading, so the directory listing and this catalog tell one story"))
					}
				} else if !seen {
					concerns[concern] = heading
				}
			}
		}
	}
	for inside, relative := range notes {
		if !referenced[inside] {
			issues = append(issues, found(relative+": note is not listed in the catalog",
				"add its entry to context/INDEX.md as one unwrapped line, or remove the note; an uncatalogued note is reachable only by someone who knows its filename"))
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
