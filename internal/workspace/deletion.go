package workspace

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// Deletion describes removing a member's records, every record in a workspace
// of one member, or the records a person named by ID, together with the ID
// reservations they held. Without confirmation it is a preview that deletes
// nothing, because the person asked has to see exactly which records and
// numbers go before they agree to it.
type Deletion struct {
	Member   string          `json:"member,omitempty"`
	All      bool            `json:"all,omitempty"`
	IDs      []string        `json:"ids,omitempty"`
	Records  []DeletedRecord `json:"records"`
	Released []string        `json:"released_ids"`
	Deleted  bool            `json:"deleted"`

	// Unlinked names each plan taken out of its own intent's plans: list.
	// Deleting a plan by ID would otherwise almost never be possible, because
	// nearly every plan made from an intent is listed there.
	Unlinked []string `json:"unlinked,omitempty"`

	// Warnings name what Git still holds for a deleted plan number: branches
	// and worktree directories a reused number would run into. They never
	// refuse, because those often live on a remote and only the person can say
	// whether the work on them still matters.
	Warnings []string `json:"warnings,omitempty"`

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
	Path      string `json:"path,omitempty"`

	// Status is draft or approved for an intent, open or completed for a plan,
	// and "reservation only" for a reserved number no record file holds, so
	// the person agrees to what they actually see.
	Status string `json:"status"`
}

const reservationOnly = "reservation only"

func deletedStatus(record Record) string {
	switch {
	case strings.HasPrefix(record.ID, "i") && record.ApprovedAt == "":
		return "draft"
	case strings.HasPrefix(record.ID, "i"):
		return "approved"
	case record.CompletedAt == "":
		return "open"
	}
	return "completed"
}

// DeleteRecords removes every intent and plan one member created, active and
// archived, or the records named by ID, and releases their numbers so
// allocation can hand them out again. While the roster lists exactly one member
// it can instead remove every record and empty every ID file. Records written
// without an author by an earlier version are reached by --all or by their ID;
// nothing selects them by member. It refuses rather than leave a record outside
// the selection pointing at something gone, or this machine's worktree prepared
// for a plan whose number is about to become reusable.
func (s *Store) DeleteRecords(ctx context.Context, member string, all bool, ids []string, confirm bool) (Deletion, error) {
	members, err := s.Members()
	if err != nil {
		return Deletion{}, err
	}
	byID := len(ids) > 0
	switch {
	case byID && (all || member != ""):
		return Deletion{}, errors.New("--id names the records itself; pass it without --member or --all")
	case all && member != "":
		return Deletion{}, errors.New("pass --member or --all, not both")
	case all && len(members.Members) != 1:
		return Deletion{}, fmt.Errorf("--all is for a workspace whose roster lists exactly one member; this one lists %d, so name whose records to delete with --member ID, or the records with --id ID", len(members.Members))
	case byID:
		for _, id := range ids {
			if !recordPattern.MatchString(id) {
				return Deletion{}, fmt.Errorf("nothing was deleted: %q is not an intent or plan ID, such as i100 or p1000", id)
			}
		}
		ids = slices.Clone(ids)
		slices.Sort(ids)
		ids = slices.Compact(ids)
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
	reserved := map[string]bool{}
	for _, rel := range idFiles {
		ledger, err := s.readIDFile(rel)
		if err != nil {
			return Deletion{}, err
		}
		ledgers[rel] = ledger
		for _, id := range append(slices.Clone(ledger.Intents), ledger.Plans...) {
			reserved[id] = true
		}
	}
	paths, err := s.RecordPaths(true)
	if err != nil {
		return Deletion{}, err
	}
	records := []Record{}
	byRecord := map[string]Record{}
	for _, path := range paths {
		record, err := s.readRecord(path)
		if err != nil {
			return Deletion{}, err
		}
		records = append(records, record)
		byRecord[record.ID] = record
	}
	result := Deletion{Member: member, All: all, IDs: ids, Records: []DeletedRecord{}, Released: []string{}}
	selected := map[string]bool{}
	for _, record := range records {
		if all || byID && slices.Contains(ids, record.ID) || !byID && record.CreatedBy == member {
			selected[record.ID] = true
			path := record.Path
			if folderPlan(record) {
				path = filepath.ToSlash(filepath.Dir(record.Path)) + "/"
			}
			result.Records = append(result.Records, DeletedRecord{record.ID, record.CreatedBy, path, deletedStatus(record)})
		}
	}
	if byID {
		// A number an interrupted create reserved and never wrote a file for is
		// otherwise only reachable with --all, which a roster of two or more
		// cannot use.
		var unknown []string
		for _, id := range ids {
			switch {
			case selected[id]:
			case reserved[id]:
				selected[id] = true
				result.Records = append(result.Records, DeletedRecord{ID: id, Status: reservationOnly})
			default:
				unknown = append(unknown, id)
			}
		}
		if len(unknown) > 0 {
			return Deletion{}, fmt.Errorf("nothing was deleted: no record or reserved ID matches %s", strings.Join(unknown, ", "))
		}
	}
	if len(result.Records) == 0 && (!all || len(reserved) == 0) {
		if all {
			return Deletion{}, errors.New("there are no intents, plans, or reserved IDs to delete")
		}
		return Deletion{}, fmt.Errorf("no intents or plans were created by member %s", member)
	}
	// Nothing stops one member planning under another's intent, and a person
	// naming records by ID names only some of them. Deleting the half that was
	// selected would leave a record pointing at a number that may soon mean
	// something else. A plan listed by its own intent is the one link that is
	// edited out instead, and only for records named by ID.
	var blockers []string
	unlink := map[string][]string{}
	addIDs := map[string]bool{}
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
			addIDs[record.ID] = true
		}
		for _, dep := range record.DependsOn {
			if selected[dep] {
				blockers = append(blockers, fmt.Sprintf("%s (%s) depends on %s", record.ID, owner, dep))
				addIDs[record.ID] = true
			}
		}
		for _, plan := range record.Plans {
			if !selected[plan] {
				continue
			}
			if byID && byRecord[plan].Intent == record.ID {
				unlink[record.ID] = append(unlink[record.ID], plan)
				continue
			}
			blockers = append(blockers, fmt.Sprintf("%s (%s) links plan %s", record.ID, owner, plan))
			addIDs[record.ID] = true
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
		message := fmt.Sprintf("nothing was deleted, because these would be left pointing at deleted records: %s. Resolve each first: its author edits the link out or deletes the record", strings.Join(blockers, "; "))
		if byID && len(addIDs) > 0 {
			var more []string
			for _, id := range slices.Sorted(maps.Keys(addIDs)) {
				more = append(more, "--id "+id)
			}
			message += fmt.Sprintf(", or delete them together by adding %s to this call", strings.Join(more, " "))
		}
		return Deletion{}, errors.New(message)
	}
	unlinkFrom := slices.Sorted(maps.Keys(unlink))
	for _, id := range unlinkFrom {
		for _, plan := range unlink[id] {
			result.Unlinked = append(result.Unlinked, fmt.Sprintf("%s from intent %s", plan, id))
		}
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
	if byID {
		var plans []string
		for _, id := range ids {
			if strings.HasPrefix(id, "p") {
				plans = append(plans, id)
			}
		}
		result.Warnings = s.leftovers(ctx, plans)
	}
	if !confirm {
		result.ConfirmationRequired = "nothing was deleted: show the person every record above, its author and status, and the IDs it releases, and run again with --confirm only after they agree to this list"
		return result, nil
	}
	// Each plan leaves its intent's list first. Interrupted after that, the
	// plan still stands and still names its intent, so the same call finishes
	// the job; released first, an interruption would leave the intent listing
	// a number nothing selects any more.
	for _, id := range unlinkFrom {
		intent := byRecord[id]
		remaining := []string{}
		for _, plan := range intent.Plans {
			if !slices.Contains(unlink[id], plan) {
				remaining = append(remaining, plan)
			}
		}
		if err := s.editRecord(intent, "plans", remaining); err != nil {
			return Deletion{}, fmt.Errorf("nothing was released: could not remove %s from intent %s: %w", strings.Join(unlink[id], ", "), id, err)
		}
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
		if record.Path == "" {
			continue
		}
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

// leftovers names what each bound checkout on this machine still holds for the
// given plan numbers: local and remote-tracking cc/<plan>/* branches, worktrees
// on such a branch, and the default .worktrees/<plan> directory. A plan
// created later under a released number would meet a branch that already
// exists, and --reuse would start it from the old plan's commits. Nothing here
// refuses: a repository that cannot be read is reported as not inspected.
func (s *Store) leftovers(ctx context.Context, plans []string) []string {
	if len(plans) == 0 {
		return nil
	}
	bindings, err := s.Bindings()
	if err != nil {
		return []string{fmt.Sprintf("bound checkouts were not inspected for leftover branches or worktrees: %v", err)}
	}
	repos := slices.Sorted(maps.Keys(bindings.Bindings))
	var warnings []string
	var listed []string
	for _, repo := range repos {
		checkout, err := s.Repository(ctx, repo)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s was not inspected for leftover branches or worktrees: %v", repo, err))
			continue
		}
		refs, err := Git(ctx, checkout.Path, "for-each-ref", "--format=%(refname)", "refs/heads/cc", "refs/remotes")
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s was not inspected for leftover branches: %v", repo, err))
		}
		for _, ref := range strings.Fields(refs) {
			for _, plan := range plans {
				switch {
				case strings.HasPrefix(ref, "refs/heads/cc/"+plan+"/"):
					warnings = append(warnings, fmt.Sprintf("branch %s in %s is left from an earlier %s; a plan given this number again meets it, and --reuse would start from its commits", strings.TrimPrefix(ref, "refs/heads/"), repo, plan))
				case strings.HasPrefix(ref, "refs/remotes/") && strings.Contains(ref, "/cc/"+plan+"/"):
					warnings = append(warnings, fmt.Sprintf("remote-tracking branch %s in %s is left from an earlier %s", strings.TrimPrefix(ref, "refs/remotes/"), repo, plan))
				}
			}
		}
		trees, err := s.Worktrees(ctx, repo)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s was not inspected for leftover worktrees: %v", repo, err))
		}
		for _, tree := range trees {
			for _, plan := range plans {
				if !tree.Main && strings.HasPrefix(tree.Branch, "cc/"+plan+"/") {
					listed = append(listed, tree.Path)
					warnings = append(warnings, fmt.Sprintf("worktree directory %s in %s is on branch %s", tree.Path, repo, tree.Branch))
				}
			}
		}
	}
	for _, plan := range plans {
		dir, err := s.LocalPath(filepath.Join(".worktrees", plan))
		if err != nil {
			continue
		}
		if resolved, e := filepath.EvalSymlinks(dir); e == nil {
			dir = resolved
		}
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		covered := slices.ContainsFunc(listed, func(path string) bool {
			return strings.HasPrefix(filepath.Clean(path)+string(filepath.Separator), filepath.Clean(dir)+string(filepath.Separator))
		})
		if !covered {
			warnings = append(warnings, fmt.Sprintf("worktree directory %s is left from an earlier %s", dir, plan))
		}
	}
	return warnings
}
