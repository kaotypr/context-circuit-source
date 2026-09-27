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

// Deletion describes removing a member's records, or every record in a solo
// workspace, together with the ledger reservations they held. Without
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
// archived, and releases their numbers so allocation can hand them out again. A
// solo workspace has no roster to select by, so there it removes every record
// and empties the ledger. It refuses rather than leave another member's record
// pointing at something gone, or this machine's worktree prepared for a plan
// whose number is about to become reusable.
func (s *Store) DeleteRecords(member string, all, confirm bool) (Deletion, error) {
	members, err := s.Members()
	if err != nil {
		return Deletion{}, err
	}
	solo := len(members.Members) == 0
	switch {
	case all && member != "":
		return Deletion{}, errors.New("pass --member or --all, not both")
	case solo && !all:
		return Deletion{}, errors.New("this workspace lists no members, so its records have no author to select by: pass --all to delete every intent and plan")
	case !solo && all:
		return Deletion{}, errors.New("--all is for a solo workspace; this one lists members, so name whose records to delete with --member ID")
	case !all:
		if err := Name(member); err != nil {
			return Deletion{}, err
		}
	}
	var ledger Ledger
	if err := s.YAML(".context-circuit/ids.yaml", &ledger); err != nil {
		return Deletion{}, err
	}
	if ledger.Intents == nil || ledger.Plans == nil {
		return Deletion{}, errors.New("ID ledger needs intents and plans lists; restore it from workspace history")
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
	if len(result.Records) == 0 && (!all || len(ledger.Intents)+len(ledger.Plans) == 0) {
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
			blockers = append(blockers, fmt.Sprintf("worktree %s in %s was prepared for %s; remove it with `worktree remove --repo %s --path %s` first, or keep the plan", assoc.Path, assoc.Repository, assoc.Plan, assoc.Repository, assoc.Path))
		}
	}
	if len(blockers) > 0 {
		sort.Strings(blockers)
		return Deletion{}, fmt.Errorf("nothing was deleted, because these would be left pointing at deleted records: %s. Resolve each first: its author edits the link out or deletes the record", strings.Join(blockers, "; "))
	}
	// A solo deletion resets the ledger outright, including numbers an
	// interrupted create reserved and never wrote a file for.
	keep := func(ids []string) []string {
		kept := []string{}
		for _, id := range ids {
			if !all && !selected[id] {
				kept = append(kept, id)
			}
		}
		return kept
	}
	for _, id := range append(append([]string{}, ledger.Intents...), ledger.Plans...) {
		if all || selected[id] {
			result.Released = append(result.Released, id)
		}
	}
	sort.Strings(result.Released)
	if !confirm {
		result.ConfirmationRequired = "nothing was deleted: show the person every record above and the IDs it releases, and run again with --confirm only after they agree to this list"
		return result, nil
	}
	// The ledger goes first. Interrupted after it, the files still stand and
	// allocation still counts their filenames as taken, so no number is handed
	// out twice; running the same deletion again finishes the job. Interrupted
	// the other way round, the numbers would stay reserved with nothing left to
	// select them by.
	intents, plans := keep(ledger.Intents), keep(ledger.Plans)
	if !slices.Equal(intents, ledger.Intents) {
		if err := s.Update(".context-circuit/ids.yaml", []string{"intents"}, intents, 0644); err != nil {
			return Deletion{}, err
		}
	}
	if !slices.Equal(plans, ledger.Plans) {
		if err := s.Update(".context-circuit/ids.yaml", []string{"plans"}, plans, 0644); err != nil {
			return Deletion{}, err
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
