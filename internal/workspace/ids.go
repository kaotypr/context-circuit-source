package workspace

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
)

// LegacyIDFile is the one shared reservation list a schema 2 workspace kept.
// Every member appended to the end of the same two lists, so two clones that
// each reserved a number conflicted in Git even though bands kept the numbers
// apart. It is only ever read and split now, never written.
const LegacyIDFile = ".context-circuit/ids.yaml"

// IDFilesDir holds one reservation file per member. A member's record creation
// writes only to that member's file, so once the split is shared, two members'
// clones change different files. The split itself writes every author's file,
// which is why one person makes it and pushes before the others change anything.
const IDFilesDir = ".context-circuit/ids"

// CurrentSchema is the workspace schema this CLI writes. Schema 3 is the
// per-member ID layout; schema 2 is still read, and moved forward by the split.
const CurrentSchema = 3

// Ledger is one ID file: the intent and plan numbers it reserves. Both lists
// are always present, so a file that lost one is reported rather than read as
// reserving nothing.
type Ledger struct {
	Intents []string `yaml:"intents" json:"intents"`
	Plans   []string `yaml:"plans" json:"plans"`
}

// windowsDevice matches the names Windows reserves for devices whatever the
// extension, so `con.yaml` cannot be checked out there. A member ID is a file
// name on every clone, so these are refused on every OS.
var windowsDevice = regexp.MustCompile(`^(con|prn|aux|nul|com[1-9]|lpt[1-9])$`)

// MemberID validates an ID a new member can be given: a lowercase name that is
// also a file name every clone can check out.
func MemberID(id string) error {
	if err := Name(id); err != nil {
		return err
	}
	if windowsDevice.MatchString(id) {
		return fmt.Errorf("member ID %q is a Windows device name, which no clone on Windows could check out as an ID file; choose another", id)
	}
	return nil
}

// renameMemberByHand is the remedy for a member ID that cannot name an ID file.
// No command renames a member, so it is an edit, and choosing the new ID is the
// person's.
const renameMemberByHand = "rename that member by hand to an ID that is not a Windows device name: its key in members.yaml, the created_by of every record they authored, and member.local.yaml on each machine that selected them; then run the command again"

// IDFilePath names a member's ID file relative to the workspace root.
func IDFilePath(member string) (string, error) {
	if err := MemberID(member); err != nil {
		return "", err
	}
	return IDFilesDir + "/" + member + ".yaml", nil
}

// IDFiles lists every ID file that exists: the legacy shared file if it is
// present, then every regular YAML file in IDFilesDir, sorted, slash-separated,
// and relative to the workspace root.
func (s *Store) IDFiles() ([]string, error) {
	files := []string{}
	legacy, err := s.Path(LegacyIDFile)
	if err != nil {
		return nil, err
	}
	if _, err := os.Lstat(legacy); err == nil {
		files = append(files, LegacyIDFile)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	dir, err := s.Path(IDFilesDir)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	for _, entry := range entries {
		if entry.Type().IsRegular() && strings.HasSuffix(entry.Name(), ".yaml") {
			files = append(files, IDFilesDir+"/"+entry.Name())
		}
	}
	sort.Strings(files)
	return files, nil
}

// IsIDFile reports whether a workspace-relative path, such as one Git status
// names, is an ID file. It covers a deleted file too, which IDFiles cannot.
func IsIDFile(rel string) bool {
	rel = filepath.ToSlash(rel)
	if rel == LegacyIDFile {
		return true
	}
	name, ok := strings.CutPrefix(rel, IDFilesDir+"/")
	return ok && strings.HasSuffix(name, ".yaml") && !strings.Contains(name, "/")
}

// readIDFile decodes one ID file strictly. A merge left unresolved is the
// likeliest reason one stops parsing, so it is named.
func (s *Store) readIDFile(rel string) (Ledger, error) {
	var ledger Ledger
	if err := s.YAML(rel, &ledger); err != nil {
		if data, readErr := s.Read(rel); readErr == nil && conflictMarker.Match(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))) {
			return ledger, fmt.Errorf("%s holds unresolved Git conflict markers", rel)
		}
		return ledger, err
	}
	if ledger.Intents == nil || ledger.Plans == nil {
		return ledger, fmt.Errorf("%s needs intents and plans lists; restore it from workspace history", rel)
	}
	return ledger, nil
}

// ledgerProblems lists what is wrong inside one readable ID file: a malformed
// ID, one under the wrong list, or one reserved twice in the same file.
func ledgerProblems(rel string, ledger Ledger) []string {
	var problems []string
	seen := map[string]bool{}
	for _, group := range []struct {
		prefix string
		values []string
	}{{"i", ledger.Intents}, {"p", ledger.Plans}} {
		for _, id := range group.values {
			switch {
			case !recordPattern.MatchString(id) || !strings.HasPrefix(id, group.prefix):
				problems = append(problems, fmt.Sprintf("%s: invalid reserved ID %q", rel, id))
			case seen[id]:
				problems = append(problems, fmt.Sprintf("%s: %s is reserved twice", rel, id))
			}
			seen[id] = true
		}
	}
	return problems
}

// idReservations reads every ID file, refusing at the first one that cannot be
// read or holds a malformed, misplaced, or repeated ID.
func (s *Store) idReservations() (map[string]Ledger, error) {
	files, err := s.IDFiles()
	if err != nil {
		return nil, err
	}
	result := map[string]Ledger{}
	for _, rel := range files {
		ledger, err := s.readIDFile(rel)
		if err != nil {
			return nil, err
		}
		if problems := ledgerProblems(rel, ledger); len(problems) > 0 {
			return nil, fmt.Errorf("%s; resolve the workspace merge first", problems[0])
		}
		result[rel] = ledger
	}
	return result, nil
}

// reservedIDs is the union of every ID file, each ID mapped to the files that
// reserve it. A number reserved in two files is taken, not an error: two clones
// that split before pulling each other leave exactly that, and `check` reports
// it for a person to tidy.
func (s *Store) reservedIDs() (map[string][]string, error) {
	files, err := s.idReservations()
	if err != nil {
		return nil, err
	}
	reserved := map[string][]string{}
	for rel, ledger := range files {
		for _, id := range append(append([]string{}, ledger.Intents...), ledger.Plans...) {
			reserved[id] = append(reserved[id], rel)
		}
	}
	for id := range reserved {
		sort.Strings(reserved[id])
	}
	return reserved, nil
}

// reserve appends IDs to one member's file and to no other, creating the file
// with both lists when the member has none yet.
func (s *Store) reserve(member string, intents, plans []string) error {
	rel, err := IDFilePath(member)
	if err != nil {
		return err
	}
	p, err := s.Path(rel)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(p); os.IsNotExist(err) {
		// Store.Write creates missing directories private to this user, which is
		// right for local state and wrong for a directory every clone shares.
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			return err
		}
		data, err := yaml.Marshal(Ledger{Intents: append([]string{}, intents...), Plans: append([]string{}, plans...)})
		if err != nil {
			return err
		}
		return s.Write(rel, append([]byte(idFileHeader(member)), data...), 0644)
	} else if err != nil {
		return err
	}
	ledger, err := s.readIDFile(rel)
	if err != nil {
		return err
	}
	if len(intents) > 0 {
		if err := s.Update(rel, []string{"intents"}, append(ledger.Intents, intents...), 0644); err != nil {
			return err
		}
	}
	if len(plans) > 0 {
		if err := s.Update(rel, []string{"plans"}, append(ledger.Plans, plans...), 0644); err != nil {
			return err
		}
	}
	return nil
}

// checkIDFiles reports what is wrong across the ID files and returns every ID
// any of them reserves. A file it cannot read reserves nothing it can see, so
// that file is reported rather than every record it would have covered being
// read as reserved nowhere. Numbers outside the owner's band, and files for
// members no longer on the roster, are valid and not reported.
func (s *Store) checkIDFiles(issues []Finding) ([]Finding, map[string]bool, error) {
	files, err := s.IDFiles()
	if err != nil {
		return issues, nil, err
	}
	reserved := map[string]bool{}
	holders := map[string][]string{}
	for _, rel := range files {
		if rel != LegacyIDFile {
			if _, err := IDFilePath(strings.TrimSuffix(strings.TrimPrefix(rel, IDFilesDir+"/"), ".yaml")); err != nil {
				issues = append(issues, found(rel+" is not named after a member ID",
					needsAPerson+"move its reservations into "+IDFilesDir+"/<member>.yaml for the member who holds them and remove it; its numbers still count as taken until then"))
			}
		}
		ledger, err := s.readIDFile(rel)
		if err != nil {
			issues = append(issues, found(err.Error(),
				needsAPerson+"repair "+rel+" by hand, preserving every reservation it already held"))
			continue
		}
		for _, problem := range ledgerProblems(rel, ledger) {
			issues = append(issues, found(problem,
				needsAPerson+"correct "+rel+" by hand without dropping a reservation; only `record delete` releases one"))
		}
		seen := map[string]bool{}
		for _, id := range append(append([]string{}, ledger.Intents...), ledger.Plans...) {
			reserved[id] = true
			// The next changing command splits a legacy file and skips what a
			// member file already holds, so an overlap with it is not a person's.
			if rel != LegacyIDFile && !seen[id] {
				seen[id] = true
				holders[id] = append(holders[id], rel)
			}
		}
	}
	ids := make([]string, 0, len(holders))
	for id := range holders {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if len(holders[id]) > 1 {
			issues = append(issues, found(id+" is reserved in more than one ID file: "+strings.Join(holders[id], ", "),
				needsAPerson+"two clones each moved it into their own member's file; keep it in the file of the member who created the record and remove the extra entry by hand"))
		}
	}
	return issues, reserved, nil
}

// idFileHeader opens a new member file. Besides telling a reader whose file it
// is, it keeps the file from reading as a copy of the shared file it replaces:
// Git pairs a deleted file with a similar new one as a rename, and two clones
// that each split before pulling would otherwise have renamed the shared file
// to two different member files, which Git reports as a conflict.
func idFileHeader(member string) string {
	return "# Intent and plan numbers reserved by member " + member + ". Only " + member + "'s new\n" +
		"# records add to this file, and only `record delete` takes a number out of it.\n"
}

// IDSplit reports what moving the legacy shared ID file into member files did.
type IDSplit struct {
	Moved         map[string][]string `json:"moved"`
	Skipped       []string            `json:"skipped"`
	LegacyRemoved bool                `json:"legacy_removed"`
	FromSchema    int                 `json:"from_schema"`
	ToSchema      int                 `json:"to_schema"`
	Commit        string              `json:"commit_recommended"`
}

// Notice is the one paragraph a person reads on the error stream when a
// command split the ID file before doing what it was asked.
func (split IDSplit) Notice() string {
	var parts []string
	files := make([]string, 0, len(split.Moved))
	for rel := range split.Moved {
		files = append(files, rel)
	}
	sort.Strings(files)
	for _, rel := range files {
		parts = append(parts, fmt.Sprintf("%s to %s", strings.Join(split.Moved[rel], ", "), rel))
	}
	var b strings.Builder
	b.WriteString("notice: moved this workspace's ID reservations into one file per member")
	if split.FromSchema != split.ToSchema {
		fmt.Fprintf(&b, " and moved its schema from %d to %d, which CLIs older than 2.4.0 refuse", split.FromSchema, split.ToSchema)
	}
	b.WriteString(".")
	if len(parts) > 0 {
		b.WriteString(" Moved " + strings.Join(parts, "; ") + ".")
	}
	if len(split.Skipped) > 0 {
		b.WriteString(" Already reserved in a member file: " + strings.Join(split.Skipped, ", ") + ".")
	}
	if split.LegacyRemoved {
		b.WriteString(" Removed " + LegacyIDFile + ".")
	}
	b.WriteString(" " + split.Commit + "\n")
	return b.String()
}

// IDSplitPending reports whether the workspace still needs SplitIDLedger: it
// holds the legacy shared file, or still records schema 2.
func (s *Store) IDSplitPending() (bool, error) {
	files, err := s.IDFiles()
	if err != nil {
		return false, err
	}
	if len(files) > 0 && files[0] == LegacyIDFile {
		return true, nil
	}
	cfg, err := s.Config()
	if err != nil {
		return false, err
	}
	return cfg.Version < CurrentSchema, nil
}

// SplitIDLedger moves every reservation in the legacy shared file into the file
// of the member whose record holds it, gives a reservation with no author or no
// record to the local member, moves the schema to 3, and removes the shared
// file. It is the same routine for the first split and for a shared file merged
// back later from a clone that had not pulled the split: an ID already reserved
// in a member file is skipped, never added twice. Interrupted anywhere, running
// it again finishes the job.
func (s *Store) SplitIDLedger(member string) (*IDSplit, error) {
	if err := MemberID(member); err != nil {
		return nil, fmt.Errorf("the member selected on this machine, %q, cannot have an ID file: %w; %s", member, err, renameMemberByHand)
	}
	cfg, err := s.Config()
	if err != nil {
		return nil, err
	}
	split := &IDSplit{Moved: map[string][]string{}, Skipped: []string{}, FromSchema: cfg.Version, ToSchema: cfg.Version}
	files, err := s.IDFiles()
	if err != nil {
		return nil, err
	}
	hasLegacy := len(files) > 0 && files[0] == LegacyIDFile
	if hasLegacy {
		legacy, err := s.readIDFile(LegacyIDFile)
		if err != nil {
			return nil, fmt.Errorf("%w; resolve the merge in %s first, keeping every entry, then run the command again", err, LegacyIDFile)
		}
		if problems := ledgerProblems(LegacyIDFile, legacy); len(problems) > 0 {
			return nil, fmt.Errorf("%s; correct %s by hand without dropping a reservation, then run the command again", strings.Join(problems, "; "), LegacyIDFile)
		}
		existing := map[string]bool{}
		for _, rel := range files[1:] {
			ledger, err := s.readIDFile(rel)
			if err != nil {
				return nil, err
			}
			for _, id := range append(append([]string{}, ledger.Intents...), ledger.Plans...) {
				existing[id] = true
			}
		}
		authors := map[string]string{}
		paths, err := s.RecordPaths(true)
		if err != nil {
			return nil, err
		}
		for _, path := range paths {
			record, err := s.readRecord(path)
			if err != nil {
				return nil, err
			}
			authors[record.ID] = record.CreatedBy
		}
		type pending struct{ intents, plans []string }
		owners := map[string]*pending{}
		for _, id := range append(append([]string{}, legacy.Intents...), legacy.Plans...) {
			if existing[id] {
				split.Skipped = append(split.Skipped, id)
				continue
			}
			owner := authors[id]
			if owner == "" {
				owner = member
			}
			if owners[owner] == nil {
				owners[owner] = &pending{}
			}
			if strings.HasPrefix(id, "i") {
				owners[owner].intents = append(owners[owner].intents, id)
			} else {
				owners[owner].plans = append(owners[owner].plans, id)
			}
		}
		names := make([]string, 0, len(owners))
		for owner := range owners {
			names = append(names, owner)
		}
		sort.Strings(names)
		// Every author is checked before anything is written, so a name that
		// cannot be a file stops the split whole and is named with its records.
		for _, owner := range names {
			if err := MemberID(owner); err != nil {
				held := append(append([]string{}, owners[owner].intents...), owners[owner].plans...)
				return nil, fmt.Errorf("%s name %q in created_by, which cannot name an ID file: %w; %s", strings.Join(held, ", "), owner, err, renameMemberByHand)
			}
		}
		for _, owner := range names {
			if err := s.reserve(owner, owners[owner].intents, owners[owner].plans); err != nil {
				return nil, err
			}
			rel, _ := IDFilePath(owner)
			split.Moved[rel] = append(append([]string{}, owners[owner].intents...), owners[owner].plans...)
		}
	}
	// The member running the split always has a file afterwards, so a workspace
	// whose shared file was empty still ends in the per-member layout.
	if err := s.reserve(member, nil, nil); err != nil {
		return nil, err
	}
	if cfg.Version < CurrentSchema {
		if err := s.Update("workspace.yaml", []string{"version"}, CurrentSchema, 0644); err != nil {
			return nil, err
		}
		split.ToSchema = CurrentSchema
	}
	if hasLegacy {
		p, err := s.Path(LegacyIDFile)
		if err != nil {
			return nil, err
		}
		if err := os.Remove(p); err != nil {
			return nil, err
		}
		split.LegacyRemoved = true
	}
	split.Commit = "Commit and push workspace.yaml and .context-circuit/ so every clone sees the same reservations; committing is the person's decision."
	return split, nil
}
