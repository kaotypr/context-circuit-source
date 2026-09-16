package workspace

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

type Ledger struct {
	Intents []string `yaml:"intents" json:"intents"`
	Plans   []string `yaml:"plans" json:"plans"`
}
type Record struct {
	ID           string   `yaml:"id" json:"id"`
	CreatedBy    string   `yaml:"created_by" json:"created_by"`
	CreatedAt    ISOTime  `yaml:"created_at" json:"created_at"`
	ApprovedAt   ISOTime  `yaml:"approved_at,omitempty" json:"approved_at,omitempty"`
	Intent       string   `yaml:"intent,omitempty" json:"intent,omitempty"`
	Repositories []string `yaml:"repositories,omitempty" json:"repositories,omitempty"`
	DependsOn    []string `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	CompletedAt  ISOTime  `yaml:"completed_at,omitempty" json:"completed_at,omitempty"`
	Plans        []string `yaml:"plans,omitempty" json:"plans,omitempty"`
	Path         string   `yaml:"-" json:"path"`
	Content      string   `yaml:"-" json:"content,omitempty"`

	// Set only when a record is created in a workspace behind its remote, where
	// the number just reserved may already be taken in another clone. Never
	// written to the record: it describes this allocation, not this record.
	SyncRequired string `yaml:"-" json:"sync_required,omitempty"`

	// Plans already holding unmerged work in this plan's repositories. Said
	// here because this is the last moment declaring a dependency is free:
	// after this, the omission is only visible as a worktree missing work.
	UnmergedPlans string `yaml:"-" json:"unmerged_plans,omitempty"`
}

var legacyRecordKey = regexp.MustCompile(`(?m)^completed:`)
var recordPattern = regexp.MustCompile(`^(i[0-9]{3,}|p[0-9]{4,})$`)
var recordFilename = regexp.MustCompile(`^(i[0-9]{3,}|p[0-9]{4,})-[a-z][a-z0-9-]*\.md$`)

func splitRecord(data []byte) ([]byte, []byte, error) {
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(data, []byte("---\n")) {
		return nil, nil, errors.New("record needs YAML frontmatter")
	}
	end := bytes.Index(data[4:], []byte("\n---\n"))
	if end < 0 {
		return nil, nil, errors.New("record frontmatter is not closed")
	}
	return data[4 : 4+end], data[4+end+5:], nil
}

func (s *Store) RecordPaths(archived bool) ([]string, error) {
	var paths []string
	for _, folder := range []string{"intent", "plans"} {
		base, err := s.Path(folder)
		if err != nil {
			return nil, err
		}
		err = filepath.WalkDir(base, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("record directories cannot contain symlinks: %s", path)
			}
			if entry.IsDir() {
				if path != base && !archived {
					return filepath.SkipDir
				}
				return nil
			}
			if recordFilename.MatchString(entry.Name()) {
				rel, _ := filepath.Rel(s.Root, path)
				paths = append(paths, filepath.ToSlash(rel))
			}
			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func (s *Store) FindRecord(id string) (Record, error) {
	if !recordPattern.MatchString(id) {
		return Record{}, fmt.Errorf("invalid record ID: %s", id)
	}
	paths, err := s.RecordPaths(true)
	if err != nil {
		return Record{}, err
	}
	var matches []string
	for _, path := range paths {
		if strings.HasPrefix(filepath.Base(path), id+"-") {
			matches = append(matches, path)
		}
	}
	if len(matches) != 1 {
		return Record{}, fmt.Errorf("expected one record for %s; found %d", id, len(matches))
	}
	return s.readRecord(matches[0])
}

func (s *Store) readRecord(path string) (Record, error) {
	var record Record
	data, err := s.Read(path)
	if err != nil {
		return record, err
	}
	header, _, err := splitRecord(data)
	if err != nil {
		return record, fmt.Errorf("%s: %w", path, err)
	}
	if err := Decode(header, &record); err != nil {
		// A record from 2.0.0-rc.3 or earlier fails here on a field this
		// candidate renamed. The decoder can only say the key is unknown, which
		// reads as a corrupt file rather than a candidate boundary, so the one
		// renamed key that shipped is named along with what replaced it.
		if legacyRecordKey.Match(header) {
			return record, fmt.Errorf("%s: written by 2.0.0-rc.3 or earlier, which this candidate cannot read: `completed` is now `completed_at`, approval is now the `approved_at` instant, and every record carries `created_at`. Records are not converted; start a fresh workspace, or rewrite this frontmatter by hand: %w", path, err)
		}
		return record, fmt.Errorf("%s: %w", path, err)
	}
	if !recordPattern.MatchString(record.ID) || !strings.HasPrefix(filepath.Base(path), record.ID+"-") {
		return record, fmt.Errorf("record ID and filename do not match: %s", path)
	}
	if err := Name(record.CreatedBy); err != nil {
		return record, err
	}
	if strings.HasPrefix(record.ID, "i") && !strings.HasPrefix(path, "intent/") || strings.HasPrefix(record.ID, "p") && !strings.HasPrefix(path, "plans/") {
		return record, fmt.Errorf("record is in the wrong directory: %s", path)
	}
	record.Path, record.Content = path, string(data)
	return record, nil
}

// bandWidth is the size of one member's allocation block per record kind. The
// blocks are N*width..N*width+width-1, so distinct bands can never overlap and
// band 1 still lands on the kind's minimum digit width (i100, p1000).
func bandWidth(kind string) int {
	if kind == "plan" {
		return 1000
	}
	return 100
}

// allocate reserves the next free number for a record kind. An author holding a
// band draws from that band's block alone; an author without one draws from the
// numbers no band has claimed. Either way a reserved number is never reused, so
// a band changes which number comes next and nothing about the ledger's rules.
func (s *Store) allocate(kind, author string) (string, error) {
	var ledger Ledger
	if err := s.YAML(".context-circuit/ids.yaml", &ledger); err != nil {
		return "", err
	}
	if ledger.Intents == nil || ledger.Plans == nil {
		return "", errors.New("ID ledger needs intents and plans lists; restore it from workspace history")
	}
	seen := map[string]bool{}
	for _, group := range []struct {
		prefix string
		values []string
	}{{"i", ledger.Intents}, {"p", ledger.Plans}} {
		for _, id := range group.values {
			if !recordPattern.MatchString(id) || !strings.HasPrefix(id, group.prefix) || seen[id] {
				return "", errors.New("invalid or duplicate reserved ID; resolve the workspace merge first")
			}
			seen[id] = true
		}
	}
	paths, err := s.RecordPaths(true)
	if err != nil {
		return "", err
	}
	fileIDs := map[string]bool{}
	for _, path := range paths {
		id := strings.SplitN(filepath.Base(path), "-", 2)[0]
		if fileIDs[id] {
			return "", fmt.Errorf("duplicate record ID: %s", id)
		}
		fileIDs[id], seen[id] = true, true
	}
	prefix, width, key := "i", 3, "intents"
	if kind == "plan" {
		prefix, width, key = "p", 4, "plans"
	}
	taken := map[int]bool{}
	for id := range seen {
		if strings.HasPrefix(id, prefix) {
			n, err := strconv.Atoi(id[1:])
			if err != nil || n >= 999999999 {
				return "", errors.New("numeric ID is outside the supported range")
			}
			taken[n] = true
		}
	}
	members, err := s.Members()
	if err != nil {
		return "", err
	}
	block := bandWidth(kind)
	band := members.Members[author].Band
	low, high := 1, 999999998
	if band != 0 {
		low, high = band*block, band*block+block-1
	}
	number := 0
	for n := low; n <= high; n++ {
		// An unbanded author must not walk into a band another member is
		// holding, or the band would stop preventing anything.
		if band == 0 && members.claimed(n, block) {
			continue
		}
		if !taken[n] {
			number = n
			break
		}
	}
	if number == 0 {
		if band != 0 {
			return "", fmt.Errorf("member %s has used every %s number in band %d (%d-%d); assign a further band", author, kind, band, low, high)
		}
		return "", errors.New("no unbanded numeric ID remains; assign the author a band")
	}
	id := fmt.Sprintf("%s%0*d", prefix, width, number)
	values := ledger.Intents
	if kind == "plan" {
		values = ledger.Plans
	}
	// Reserve before creating a record: an interrupted write consumes a number,
	// and removing or archiving its file will never make that number reusable.
	return id, s.Update(".context-circuit/ids.yaml", []string{key}, append(values, id), 0644)
}

func (s *Store) CreateRecord(ctx context.Context, kind, slug, title, intent string, repos, dependencies []string) (Record, error) {
	if kind != "intent" && kind != "plan" {
		return Record{}, errors.New("record kind must be intent or plan")
	}
	if err := Name(slug); err != nil {
		return Record{}, err
	}
	if err := Text(title); err != nil {
		return Record{}, err
	}
	cfg, err := s.Config()
	if err != nil {
		return Record{}, err
	}
	author, err := s.ActiveMember()
	if err != nil {
		return Record{}, err
	}
	var parent Record
	if kind == "plan" {
		if !strings.HasPrefix(intent, "i") {
			return Record{}, errors.New("plan needs an intent ID")
		}
		parent, err = s.FindRecord(intent)
		if err != nil {
			return Record{}, err
		}
		if len(repos) == 0 {
			return Record{}, errors.New("plan needs at least one repository")
		}
		for _, repo := range repos {
			if _, ok := cfg.Repositories[repo]; !ok {
				return Record{}, fmt.Errorf("unknown plan repository: %s", repo)
			}
		}
		for _, dep := range dependencies {
			if !strings.HasPrefix(dep, "p") {
				return Record{}, errors.New("plan dependencies must be plan IDs")
			}
			if _, err := s.FindRecord(dep); err != nil {
				return Record{}, err
			}
		}
	}
	id, err := s.allocate(kind, author)
	if err != nil {
		return Record{}, err
	}
	record := Record{ID: id, CreatedBy: author, CreatedAt: ISOTime(isoTime(time.Now()))}
	body := "## Goal\n\nDescribe the intended outcome.\n\n## Non-goals\n\nDescribe what is out of scope.\n\n## Constraints\n\nRecord relevant constraints.\n\n## Success criteria\n\n- Describe an observable result.\n\n## Repository scope\n\nList the repositories likely to be involved.\n\n## Open questions\n\nNone known.\n"
	folder := "intent"
	if kind == "plan" {
		folder = "plans"
		record.Intent, record.Repositories, record.DependsOn = intent, unique(repos), unique(dependencies)
		body = "## Approach\n\nGround the approach in the selected repositories.\n\n## Tasks and order\n\n- Describe implementation steps and dependencies.\n\n## Risks and checks\n\nRecord useful risks and ordinary test, lint, or build commands.\n"
	}
	header, err := yaml.Marshal(record)
	if err != nil {
		return Record{}, err
	}
	record.Path = fmt.Sprintf("%s/%s-%s.md", folder, id, slug)
	data := []byte("---\n" + string(header) + "---\n\n# " + title + "\n\n" + body)
	p, err := s.Path(record.Path)
	if err != nil {
		return Record{}, err
	}
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return Record{}, err
	}
	_, writeErr := f.Write(data)
	closeErr := f.Close()
	if writeErr != nil {
		return Record{}, writeErr
	}
	if closeErr != nil {
		return Record{}, closeErr
	}
	if kind == "plan" {
		if err := s.editRecord(parent, "plans", append(parent.Plans, id)); err != nil {
			return record, fmt.Errorf("created %s but could not link it from %s: %w", record.Path, parent.Path, err)
		}
	}
	record.Content = string(data)
	// The refusal that catches a dependent plan started from the wrong base is
	// armed by `depends_on`, which the coordinator has to remember to pass. This
	// is where remembering is still cheap: the plan is being written, the
	// repositories are known, and nothing has been prepared from anything yet.
	if kind == "plan" {
		notes := []string{}
		for _, repoID := range record.Repositories {
			checkout, err := s.Repository(ctx, repoID)
			if err != nil {
				continue
			}
			base := "refs/heads/" + checkout.BaseBranch
			if _, err := Git(ctx, checkout.Path, "rev-parse", "--verify", "--end-of-options", base+"^{commit}"); err != nil {
				continue
			}
			remoteBase := remoteCounterpart(ctx, checkout.Path, checkout.BaseBranch)
			found := s.unmergedSiblings(ctx, checkout.Path, repoID, id, base, remoteBase)
			for _, item := range found {
				if slices.Contains(record.DependsOn, item.Plan) {
					continue
				}
				notes = append(notes, fmt.Sprintf("%s in %s (%s, %s commit(s))", item.Plan, repoID, item.Branch, item.Count))
			}
		}
		if len(notes) > 0 {
			record.UnmergedPlans = fmt.Sprintf("work not in the base branch already exists: %s. Record it with `record dependencies --id %s --depends-on ID` if this plan builds on it, so preparation starts from it rather than the base.",
				strings.Join(notes, ", "), id)
		}
	}
	// Bands keep two members apart; nothing keeps one workspace apart from its
	// own unsynchronized clone, and a number is reserved before its file exists.
	// The ledger that would have shown the collision is the thing not pulled, so
	// say it where the number is handed out rather than leaving it to `check`
	// after both records are written and referenced.
	if branch, err := Git(ctx, s.Root, "symbolic-ref", "--quiet", "--short", "HEAD"); err == nil {
		if n, remote := behindRemote(ctx, s.Root, branch); n > 0 {
			record.SyncRequired = fmt.Sprintf("this workspace is %d commit(s) behind %s, so %s may already be reserved in another clone. Pull the workspace and run `check` before publishing this ID or branching from it.",
				n, remote, id)
		}
	}
	return record, nil
}

func unique(values []string) []string {
	out := []string{}
	for _, value := range values {
		if !slices.Contains(out, value) {
			out = append(out, value)
		}
	}
	return out
}

func (s *Store) editRecord(record Record, key string, value any) error {
	header, body, err := splitRecord([]byte(record.Content))
	if err != nil {
		return err
	}
	updated, err := Edit(header, []string{key}, value)
	if err != nil {
		return err
	}
	return s.replaceRecord(record, []byte("---\n"+strings.TrimRight(string(updated), "\n")+"\n---\n"+string(body)))
}

func (s *Store) replaceRecord(record Record, next []byte) error {
	current, err := s.Read(record.Path)
	if err != nil {
		return err
	}
	if string(current) != record.Content {
		return errors.New("record changed during edit; retry")
	}
	return s.Write(record.Path, next, 0644)
}

// ISOTime is an instant recorded as a canonical ISO 8601 UTC timestamp,
// 2026-09-15T10:53:00Z. A gate records when it happened rather than a word
// saying that it did, so the two cannot disagree and a second decision on the
// same day is still distinct. Marshaling it as a raw scalar keeps the unquoted
// form: the YAML encoder quotes a plain Go string that looks like a timestamp,
// which would make instants read differently depending on which field they sit
// in.
type ISOTime string

// TimeLayout is the one format every instant in a workspace record uses.
const TimeLayout = time.RFC3339

func (t ISOTime) MarshalYAML() ([]byte, error) { return []byte(string(t)), nil }

func isoTime(moment time.Time) string { return moment.UTC().Format(TimeLayout) }

func (s *Store) Note(id, text, kind string) error {
	if strings.TrimSpace(text) == "" || strings.ContainsRune(text, 0) {
		return errors.New("provide a nonempty note")
	}
	record, err := s.FindRecord(id)
	if err != nil {
		return err
	}
	stamp := isoTime(time.Now())
	next := strings.TrimRight(record.Content, "\n")
	switch kind {
	case "approve":
		if !strings.HasPrefix(id, "i") {
			return errors.New("approval is recorded on an intent")
		}
		if err := Text(text); err != nil {
			return err
		}
		// Renewed approval is a second decision, not a corrected first one, so
		// each one keeps its own date and words rather than overwriting them.
		next += "\n\n## Approval — " + stamp + "\n\n" + text
	case "complete":
		if !strings.HasPrefix(id, "p") {
			return errors.New("completion is recorded on a plan")
		}
		next += "\n\n## Completion — " + stamp + "\n\n" + text
	default:
		return errors.New("unknown note operation")
	}
	data := []byte(next + "\n")
	// Both gates also record a machine-readable frontmatter value: the status a
	// person set, and the completion date dependency ordering reads to tell a
	// finished plan from an unfinished one without reading prose. The
	// human-readable section stays; one write keeps both consistent.
	field, value := "", any(nil)
	switch kind {
	case "approve":
		field, value = "approved_at", ISOTime(stamp)
	case "complete":
		field, value = "completed_at", ISOTime(stamp)
	}
	if field != "" {
		header, body, err := splitRecord(data)
		if err != nil {
			return err
		}
		updated, err := Edit(header, []string{field}, value)
		if err != nil {
			return err
		}
		data = []byte("---\n" + strings.TrimRight(string(updated), "\n") + "\n---\n" + string(body))
	}
	return s.replaceRecord(record, data)
}

// A catalog entry is a linked list item. Prose in the index may legitimately
// discuss a repository scope without cataloguing a note, and the shipped index
// explains the entry shape using an example, so only entry-shaped lines are
// offered as reconcile candidates.
var catalogEntry = regexp.MustCompile(`^\s*[-*]\s+\[[^\]]+\]\(`)

// KnowledgeCandidates lists the catalog entries scoped to a plan's repositories.
// Completion reconciles durable knowledge, so Go locates the entries worth
// considering and the agent decides which of them actually changed meaning.
func (s *Store) KnowledgeCandidates(id string) ([]string, error) {
	record, err := s.FindRecord(id)
	if err != nil {
		return nil, err
	}
	return s.CatalogEntriesFor(record.Repositories)
}

// CatalogEntriesFor lists the entries claiming any of these repositories. A
// change made without a plan reconciles the same knowledge, so the lookup takes
// repositories rather than a record.
func (s *Store) CatalogEntriesFor(repositories []string) ([]string, error) {
	seen := map[string]bool{}
	entries := []string{}
	for _, repository := range repositories {
		lines, err := s.FindContext("{" + repository + "}")
		if err != nil {
			return nil, err
		}
		for _, line := range lines {
			if catalogEntry.MatchString(line) && !seen[line] {
				seen[line] = true
				entries = append(entries, line)
			}
		}
	}
	sort.Strings(entries)
	return entries, nil
}

func (s *Store) SetDependencies(id string, dependencies []string) error {
	record, err := s.FindRecord(id)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(id, "p") {
		return errors.New("dependencies belong to a plan")
	}
	dependencies = unique(dependencies)
	var visit func(string, map[string]bool) error
	visit = func(current string, trail map[string]bool) error {
		if !strings.HasPrefix(current, "p") || current == id || trail[current] {
			return fmt.Errorf("invalid or cyclic dependency at %s", current)
		}
		dep, err := s.FindRecord(current)
		if err != nil {
			return err
		}
		trail[current] = true
		defer delete(trail, current)
		for _, next := range dep.DependsOn {
			if err := visit(next, trail); err != nil {
				return err
			}
		}
		return nil
	}
	for _, dep := range dependencies {
		if err := visit(dep, map[string]bool{}); err != nil {
			return err
		}
	}
	return s.editRecord(record, "depends_on", dependencies)
}

func (s *Store) ListRecords(archived bool) ([]Record, error) {
	paths, err := s.RecordPaths(archived)
	if err != nil {
		return nil, err
	}
	result := []Record{}
	for _, path := range paths {
		r, err := s.readRecord(path)
		if err != nil {
			return nil, err
		}
		r.Content = ""
		result = append(result, r)
	}
	return result, nil
}
