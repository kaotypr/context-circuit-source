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

	// BlockedBy, Refused, and AddIDs are set only on a refusal. The records
	// asked for are still listed above, so the person sees what they asked for
	// and what stands in the way in one answer rather than a bare error.
	BlockedBy []BlockingRecord `json:"blocked_by,omitempty"`
	Refused   string           `json:"refused,omitempty"`
	AddIDs    []string         `json:"add_ids,omitempty"`
}

// BlockingRecord is a record outside the selection that still points at one
// inside it, and how.
type BlockingRecord struct {
	ID        string `json:"id"`
	CreatedBy string `json:"created_by,omitempty"`
	Path      string `json:"path,omitempty"`
	Status    string `json:"status"`
	Link      string `json:"link"`
}

// DeletionRefused is a deletion that would leave something pointing at a
// deleted record. It carries the preview, so a caller can show it beside the
// reason; nothing was deleted.
type DeletionRefused struct {
	Preview Deletion
}

func (e *DeletionRefused) Error() string { return e.Preview.Refused }

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

// deletionPath names a record the way deletion removes it: a folder plan as its
// folder, anything else as its file.
func deletionPath(record Record) string {
	if folderPlan(record) {
		return filepath.ToSlash(filepath.Dir(record.Path)) + "/"
	}
	return record.Path
}

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
			result.Records = append(result.Records, DeletedRecord{record.ID, record.CreatedBy, deletionPath(record), deletedStatus(record)})
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
	var blocking []BlockingRecord
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
		block := func(link string) {
			blockers = append(blockers, fmt.Sprintf("%s (%s) %s", record.ID, owner, link))
			blocking = append(blocking, BlockingRecord{record.ID, record.CreatedBy, deletionPath(record), deletedStatus(record), link})
			addIDs[record.ID] = true
		}
		if selected[record.Intent] {
			block("belongs to intent " + record.Intent)
		}
		for _, dep := range record.DependsOn {
			if selected[dep] {
				block("depends on " + dep)
			}
		}
		for _, plan := range record.Plans {
			if !selected[plan] {
				continue
			}
			// A plan with no file left, as an interrupted deletion or create
			// leaves one, points at nothing, so its entry goes the same way.
			if byID && (byRecord[plan].Intent == record.ID || byRecord[plan].ID == "") {
				unlink[record.ID] = append(unlink[record.ID], plan)
				continue
			}
			block("links plan " + plan)
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
				result.AddIDs = append(result.AddIDs, id)
			}
			message += fmt.Sprintf(", or delete them together by adding %s to this call", strings.Join(more, " "))
		}
		sort.SliceStable(blocking, func(i, j int) bool { return blocking[i].ID < blocking[j].ID })
		result.BlockedBy = blocking
		result.Refused = message
		return result, &DeletionRefused{Preview: result}
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
	var plans []string
	for _, id := range result.Released {
		if strings.HasPrefix(id, "p") {
			plans = append(plans, id)
		}
	}
	for _, record := range result.Records {
		if strings.HasPrefix(record.ID, "p") && !slices.Contains(plans, record.ID) {
			plans = append(plans, record.ID)
		}
	}
	slices.Sort(plans)
	result.Warnings = s.leftovers(ctx, plans)
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
	// The record files go next, and every number stays reserved until they are
	// all gone: interrupted here, each selected ID still holds a file or a
	// reservation, so --id reaches it and allocation never hands it out. Plans
	// go before intents, so a record goes before anything it points to.
	for _, plansFirst := range []bool{true, false} {
		for _, record := range result.Records {
			if record.Path == "" || strings.HasPrefix(record.ID, "p") != plansFirst {
				continue
			}
			if err := s.removeRecord(record.Path); err != nil {
				return Deletion{}, s.unfinished(record.Path, result.Records, reserved, all, err)
			}
		}
	}
	// The ID files go last, each released number from every file holding it.
	for _, rel := range idFiles {
		ledger := ledgers[rel]
		intents, plans := keep(ledger.Intents), keep(ledger.Plans)
		if !slices.Equal(intents, ledger.Intents) {
			if err := writeIDList(s, rel, "intents", intents); err != nil {
				return Deletion{}, s.unreleased(rel, idFiles, ledgers, result.Released, all, err)
			}
		}
		if !slices.Equal(plans, ledger.Plans) {
			if err := writeIDList(s, rel, "plans", plans); err != nil {
				return Deletion{}, s.unreleased(rel, idFiles, ledgers, result.Released, all, err)
			}
		}
	}
	result.Deleted = true
	result.SyncRequired = "the deletion exists only in this checkout: until it is committed and every clone has pulled it, another clone still holds these records and may allocate their numbers again. Committing and pushing are the person's decision. Every other machine should pull it, then run `check`, which names a worktree prepared there for a deleted plan, and remove that worktree and its cc/pNNNN branches (`worktree remove`, then `worktree prune --repo ID`) before creating plans, so a plan given a released number does not meet the old plan's work"
	return result, nil
}

// removePath removes one file or directory tree. Tests replace it to interrupt
// a deletion part-way through.
var removePath = os.RemoveAll

// removeRecord removes one record file, or a plan folder with its plan.md last:
// interrupted part-way, the folder still holds plan.md and is still a record
// the same deletion selects.
func (s *Store) removeRecord(rel string) error {
	p, err := s.Path(strings.TrimSuffix(rel, "/"))
	if err != nil {
		return err
	}
	if !strings.HasSuffix(rel, "/") {
		return removePath(p)
	}
	entries, err := os.ReadDir(p)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, entry := range entries {
		if entry.Name() != "plan.md" {
			if err := removePath(filepath.Join(p, entry.Name())); err != nil {
				return err
			}
		}
	}
	if err := removePath(filepath.Join(p, "plan.md")); err != nil {
		return err
	}
	return removePath(p)
}

// writeIDList rewrites one list in one ID file. Tests replace it to interrupt a
// deletion between ID files.
var writeIDList = func(s *Store, rel, key string, ids []string) error {
	return s.Update(rel, []string{key}, ids, 0644)
}

// finishing names the call that finishes an interrupted deletion. --all is
// re-run as itself, because it also releases numbers no record file held; any
// other deletion is finished by exactly the IDs still left, which --id reaches
// on any roster.
func finishing(ids []string, all bool) string {
	if all {
		return "Run `record delete --all --confirm` to finish"
	}
	return "Run `record delete --id " + strings.Join(ids, " --id ") + " --confirm` to finish"
}

// unreleased explains a deletion stopped while rewriting ID files. Every record
// file is already gone, so what is left is reservations: the error names the
// IDs still reserved, read back from the files, and the call that releases
// them, or says nothing is left when none is.
func (s *Store) unreleased(rel string, idFiles []string, ledgers map[string]Ledger, released []string, all bool, cause error) error {
	var left []string
	for _, file := range idFiles {
		ledger, err := s.readIDFile(file)
		if err != nil {
			ledger = ledgers[file]
		}
		for _, id := range append(slices.Clone(ledger.Intents), ledger.Plans...) {
			if slices.Contains(released, id) && !slices.Contains(left, id) {
				left = append(left, id)
			}
		}
	}
	slices.Sort(left)
	if len(left) == 0 {
		return fmt.Errorf("could not rewrite %s, but every record file is gone and no released ID is still reserved, so nothing is left to finish: %w", rel, cause)
	}
	return fmt.Errorf("could not rewrite %s: every record file is gone, and these IDs are still reserved: %s. %s: %w",
		rel, strings.Join(left, ", "), finishing(left, all), cause)
}

// unfinished explains a deletion stopped while removing files. Nothing was
// released yet, so every selected ID still holding a file or a reservation is
// reachable by --id on any roster; --member could not reach a removed record's
// number again, so the call that finishes it is the --id one, or --all again.
func (s *Store) unfinished(path string, records []DeletedRecord, reserved map[string]bool, all bool, cause error) error {
	var left []string
	for _, record := range records {
		held := reserved[record.ID]
		if !held && record.Path != "" {
			if p, err := s.Path(strings.TrimSuffix(record.Path, "/")); err == nil {
				if strings.HasSuffix(record.Path, "/") {
					p = filepath.Join(p, "plan.md")
				}
				_, err := os.Lstat(p)
				held = err == nil
			}
		}
		if held && !slices.Contains(left, record.ID) {
			left = append(left, record.ID)
		}
	}
	slices.Sort(left)
	if len(left) == 0 {
		return fmt.Errorf("could not remove %s; every record is gone and no ID was released, so remove it by hand: %w", path, cause)
	}
	// Only plan.md makes a folder a record, so once it is gone no re-run
	// reaches what is left of the folder.
	folder := ""
	if strings.HasSuffix(path, "/") {
		if p, err := s.Path(strings.TrimSuffix(path, "/")); err == nil {
			if _, err := os.Lstat(filepath.Join(p, "plan.md")); os.IsNotExist(err) {
				folder = fmt.Sprintf(", then remove the folder %s by hand: it no longer holds plan.md, so no re-run reaches it", path)
			}
		}
	}
	return fmt.Errorf("could not remove %s, so no ID was released and these stay reserved: %s. %s%s: %w",
		path, strings.Join(left, ", "), finishing(left, all), folder, cause)
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
