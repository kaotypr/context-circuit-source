package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// Deletion describes removing a member's records, or every record in a
// workspace of one member, together with the ID reservations they held. Without
// confirmation it is a preview that deletes nothing, because the person asked
// has to see exactly which records and numbers go before they agree to it.
type Deletion struct {
	Member   string          `json:"member,omitempty"`
	All      bool            `json:"all,omitempty"`
	Records  []DeletedRecord `json:"records"`
	Released []string        `json:"released_ids"`
	Deleted  bool            `json:"deleted"`

	// ConfirmationRequired is set on a preview and names what the caller still
	// owes: the person's agreement to this exact list.
	ConfirmationRequired string `json:"confirmation_required,omitempty"`

	// SyncRequired is set once files are gone. Deletion happens in this
	// checkout alone, and a clone that has not pulled it still holds the
	// records and may hand their released numbers out again.
	SyncRequired string `json:"sync_required,omitempty"`
}

type DeletedRecord struct {
	ID        string `json:"id"`
	CreatedBy string `json:"created_by,omitempty"`
	Path      string `json:"path"`
}

// DeleteRecords removes every intent and plan one member created, active and
// archived, and releases their numbers so allocation can hand them out again.
// While the roster lists exactly one member it can instead remove every record
// and empty every ID file, which is how records written without an author by an
// earlier version are removed: nothing selects them by member. It refuses
// rather than leave another member's record pointing at something gone, or
// this machine's worktree prepared for a plan whose number is about to become
// reusable.
func (s *Store) DeleteRecords(member string, all, confirm bool) (Deletion, error) {
	members, err := s.Members()
	if err != nil {
		return Deletion{}, err
	}
	switch {
	case all && member != "":
		return Deletion{}, errors.New("pass --member or --all, not both")
	case all && len(members.Members) != 1:
		return Deletion{}, fmt.Errorf("--all is for a workspace whose roster lists exactly one member; this one lists %d, so name whose records to delete with --member ID", len(members.Members))
	case !all:
		if err := Name(member); err != nil {
			return Deletion{}, err
		}
	}
	// Every ID file is read, the legacy shared one included: a preview runs
	// before any split and writes nothing, so it sees whichever layout is there.
	idFiles, err := s.IDFiles()
	if err != nil {
		return Deletion{}, err
	}
	ledgers := map[string]Ledger{}
	reservedCount := 0
	for _, rel := range idFiles {
		ledger, err := s.readIDFile(rel)
		if err != nil {
			return Deletion{}, err
		}
		ledgers[rel] = ledger
		reservedCount += len(ledger.Intents) + len(ledger.Plans)
	}
	paths, err := s.RecordPaths(true)
	if err != nil {
		return Deletion{}, err
	}
	records := []Record{}
	for _, path := range paths {
		record, err := s.readRecord(path)
		if err != nil {
			return Deletion{}, err
		}
		records = append(records, record)
	}
	result := Deletion{Member: member, All: all, Records: []DeletedRecord{}, Released: []string{}}
	selected := map[string]bool{}
	for _, record := range records {
		if all || record.CreatedBy == member {
			selected[record.ID] = true
			path := record.Path
			if folderPlan(record) {
				path = filepath.ToSlash(filepath.Dir(record.Path)) + "/"
			}
			result.Records = append(result.Records, DeletedRecord{record.ID, record.CreatedBy, path})
		}
	}
	if len(result.Records) == 0 && (!all || reservedCount == 0) {
		if all {
			return Deletion{}, errors.New("there are no intents, plans, or reserved IDs to delete")
		}
		return Deletion{}, fmt.Errorf("no intents or plans were created by member %s", member)
	}
	// Selection is by author, and nothing stops one member planning under
	// another's intent. Deleting the half that matched would leave the other
	// member's record pointing at a number that may soon mean something else.
	var blockers []string
	for _, record := range records {
		if selected[record.ID] {
			continue
		}
		owner := record.CreatedBy
		if owner == "" {
			owner = "no recorded author"
		}
		if selected[record.Intent] {
			blockers = append(blockers, fmt.Sprintf("%s (%s) belongs to intent %s", record.ID, owner, record.Intent))
		}
		for _, dep := range record.DependsOn {
			if selected[dep] {
				blockers = append(blockers, fmt.Sprintf("%s (%s) depends on %s", record.ID, owner, dep))
			}
		}
		for _, plan := range record.Plans {
			if selected[plan] {
				blockers = append(blockers, fmt.Sprintf("%s (%s) links plan %s", record.ID, owner, plan))
			}
		}
	}
	local, err := s.associations()
	if err != nil {
		return Deletion{}, err
	}
	for _, assoc := range local.Worktrees {
		if selected[assoc.Plan] {
			blockers = append(blockers, fmt.Sprintf("worktree %s in %s was prepared for %s; remove it with `worktree remove --repo %s --path %s` first (if Git no longer lists it, `worktree prune --repo %s` drops the leftover association), or keep the plan", assoc.Path, assoc.Repository, assoc.Plan, assoc.Repository, assoc.Path, assoc.Repository))
		}
	}
	if len(blockers) > 0 {
		sort.Strings(blockers)
		return Deletion{}, fmt.Errorf("nothing was deleted, because these would be left pointing at deleted records: %s. Resolve each first: its author edits the link out or deletes the record", strings.Join(blockers, "; "))
	}
	// Deleting everything empties every ID file outright, including numbers an
	// interrupted create reserved and never wrote a file for. The files stay,
	// so each member keeps the file their next record is reserved in.
	keep := func(ids []string) []string {
		kept := []string{}
		for _, id := range ids {
			if !all && !selected[id] {
				kept = append(kept, id)
			}
		}
		return kept
	}
	released := map[string]bool{}
	for _, rel := range idFiles {
		for _, id := range append(append([]string{}, ledgers[rel].Intents...), ledgers[rel].Plans...) {
			if (all || selected[id]) && !released[id] {
				released[id] = true
				result.Released = append(result.Released, id)
			}
		}
	}
	sort.Strings(result.Released)
	if !confirm {
		result.ConfirmationRequired = "nothing was deleted: show the person every record above and the IDs it releases, and run again with --confirm only after they agree to this list"
		return result, nil
	}
	// The ID files go first, each released number from every file holding it.
	// Interrupted after them, the record files still stand and allocation still
	// counts their filenames as taken, so no number is handed out twice; running
	// the same deletion again finishes the job. Interrupted the other way round,
	// the numbers would stay reserved with nothing left to select them by.
	for _, rel := range idFiles {
		ledger := ledgers[rel]
		intents, plans := keep(ledger.Intents), keep(ledger.Plans)
		if !slices.Equal(intents, ledger.Intents) {
			if err := s.Update(rel, []string{"intents"}, intents, 0644); err != nil {
				return Deletion{}, err
			}
		}
		if !slices.Equal(plans, ledger.Plans) {
			if err := s.Update(rel, []string{"plans"}, plans, 0644); err != nil {
				return Deletion{}, err
			}
		}
	}
	for _, record := range result.Records {
		p, err := s.Path(strings.TrimSuffix(record.Path, "/"))
		if err != nil {
			return Deletion{}, err
		}
		if strings.HasSuffix(record.Path, "/") {
			err = os.RemoveAll(p)
		} else {
			err = os.Remove(p)
		}
		if err != nil {
			return Deletion{}, fmt.Errorf("released the IDs but could not remove %s; run the same deletion again to finish: %w", record.Path, err)
		}
	}
	result.Deleted = true
	result.SyncRequired = "the deletion exists only in this checkout: until it is committed and every clone has pulled it, another clone still holds these records and may allocate their numbers again. Committing and pushing are the person's decision"
	return result, nil
}
