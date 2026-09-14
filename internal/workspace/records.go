package workspace

import (
	"bytes"
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
	Intent       string   `yaml:"intent,omitempty" json:"intent,omitempty"`
	Repositories []string `yaml:"repositories,omitempty" json:"repositories,omitempty"`
	DependsOn    []string `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	Completed    ISODate  `yaml:"completed,omitempty" json:"completed,omitempty"`
	Plans        []string `yaml:"plans,omitempty" json:"plans,omitempty"`
	Path         string   `yaml:"-" json:"path"`
	Content      string   `yaml:"-" json:"content,omitempty"`
}

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

func (s *Store) allocate(kind string) (string, error) {
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
	maxID := 0
	for id := range seen {
		if strings.HasPrefix(id, prefix) {
			n, err := strconv.Atoi(id[1:])
			if err != nil || n >= 999999999 {
				return "", errors.New("numeric ID is outside the supported range")
			}
			if n > maxID {
				maxID = n
			}
		}
	}
	id := fmt.Sprintf("%s%0*d", prefix, width, maxID+1)
	values := ledger.Intents
	if kind == "plan" {
		values = ledger.Plans
	}
	// Reserve before creating a record: an interrupted write consumes a number,
	// and removing or archiving its file will never make that number reusable.
	return id, s.Update(".context-circuit/ids.yaml", []string{key}, append(values, id), 0644)
}

func (s *Store) CreateRecord(kind, slug, title, intent string, repos, dependencies []string) (Record, error) {
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
	id, err := s.allocate(kind)
	if err != nil {
		return Record{}, err
	}
	record := Record{ID: id, CreatedBy: author}
	body := "## Goal\n\nDescribe the intended outcome.\n\n## Non-goals\n\nDescribe what is out of scope.\n\n## Constraints\n\nRecord relevant constraints.\n\n## Success criteria\n\n- Describe an observable result.\n\n## Repository scope\n\nList the repositories likely to be involved.\n\nApproval: Pending\n"
	folder := "intent"
	if kind == "plan" {
		folder = "plans"
		record.Intent, record.Repositories, record.DependsOn = intent, unique(repos), unique(dependencies)
		body = "## Approach\n\nGround the approach in the selected repositories.\n\n## Tasks and order\n\n- Describe implementation steps and dependencies.\n\n## Risks and checks\n\nRecord useful risks and ordinary test, lint, or build commands.\n\n## Progress and result\n\nRecord changes, observed checks, remaining work, and delivery references here.\n"
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

// ISODate is a calendar date recorded as ISO 8601 YYYY-MM-DD. Marshaling it as a
// raw scalar keeps the unquoted form: the YAML encoder quotes a plain Go string
// that looks like a date, which would make dates read differently depending on
// which field they sit in.
type ISODate string

func (d ISODate) MarshalYAML() ([]byte, error) { return []byte(string(d)), nil }

// Every date written anywhere in a workspace record uses this one format.
func isoDate(moment time.Time) string { return moment.UTC().Format(time.DateOnly) }

func (s *Store) Note(id, text, kind string) error {
	if strings.TrimSpace(text) == "" || strings.ContainsRune(text, 0) {
		return errors.New("provide a nonempty note")
	}
	record, err := s.FindRecord(id)
	if err != nil {
		return err
	}
	stamp := isoDate(time.Now())
	next := strings.TrimRight(record.Content, "\n")
	switch kind {
	case "approve":
		if !strings.HasPrefix(id, "i") {
			return errors.New("approval is recorded on an intent")
		}
		if err := Text(text); err != nil {
			return err
		}
		line := "Approval: Approved on " + stamp + " — " + text
		rx := regexp.MustCompile(`(?m)^Approval:.*$`)
		if len(rx.FindAllStringIndex(next, -1)) > 1 {
			return errors.New("multiple approval lines; identify the intended approval before updating")
		}
		if rx.MatchString(next) {
			next = rx.ReplaceAllStringFunc(next, func(string) string { return line })
		} else {
			next += "\n\n" + line
		}
	case "complete":
		if !strings.HasPrefix(id, "p") {
			return errors.New("completion is recorded on a plan")
		}
		next += "\n\n## Completion — " + stamp + "\n\n" + text
	case "note":
		next += "\n\n## Update — " + stamp + "\n\n" + text
	default:
		return errors.New("unknown note operation")
	}
	data := []byte(next + "\n")
	// Completion also records a machine-readable date so dependency ordering can
	// tell a finished plan from an unfinished one without reading prose. The
	// human-readable section stays; one write keeps both consistent.
	if kind == "complete" {
		header, body, err := splitRecord(data)
		if err != nil {
			return err
		}
		updated, err := Edit(header, []string{"completed"}, ISODate(stamp))
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
	seen := map[string]bool{}
	entries := []string{}
	for _, repository := range record.Repositories {
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
