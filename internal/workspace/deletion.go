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
	ids, err = checkSelection(member, all, ids, len(members.Members))
	if err != nil {
		return Deletion{}, err
	}
	byID := len(ids) > 0
	idFiles, ledgers, reserved, err := s.readLedgers()
	if err != nil {
		return Deletion{}, err
	}
	records, byRecord, err := s.readAllRecords()
	if err != nil {
		return Deletion{}, err
	}
	result := Deletion{Member: member, All: all, IDs: ids, Records: []DeletedRecord{}, Released: []string{}}
	selected, chosen, err := selectRecords(records, reserved, member, all, ids)
	if err != nil {
		return Deletion{}, err
	}
	result.Records = chosen
	blocks := findBlockers(records, byRecord, selected, byID)
	local, err := s.associations()
	if err != nil {
		return Deletion{}, err
	}
	blocks.addWorktrees(local, selected)
	if len(blocks.messages) > 0 {
		return refuse(result, blocks, byID)
	}
	unlinkFrom := slices.Sorted(maps.Keys(blocks.unlink))
	for _, id := range unlinkFrom {
		for _, plan := range blocks.unlink[id] {
			result.Unlinked = append(result.Unlinked, fmt.Sprintf("%s from intent %s", plan, id))
		}
	}
	result.Released = releasedIDs(idFiles, ledgers, selected, all)
	result.Warnings = s.leftovers(ctx, leftoverPlans(result.Released, result.Records))
	if !confirm {
		result.ConfirmationRequired = "nothing was deleted: show the person every record above, its author and status, and the IDs it releases, and run again with --confirm only after they agree to this list"
		return result, nil
	}
	if err := s.unlinkPlans(unlinkFrom, byRecord, blocks.unlink); err != nil {
		return Deletion{}, err
	}
	if err := s.removeRecords(result.Records, reserved, all); err != nil {
		return Deletion{}, err
	}
	if err := s.releaseIDs(idFiles, ledgers, selected, result.Released, all); err != nil {
		return Deletion{}, err
	}
	result.Deleted = true
	result.SyncRequired = "the deletion exists only in this checkout: until it is committed and every clone has pulled it, another clone still holds these records and may allocate their numbers again. Committing and pushing are the person's decision. Every other machine should pull it, then run `check`, which names a worktree prepared there for a deleted plan, and remove that worktree and its cc/pNNNN branches (`worktree remove`, then `worktree prune --repo ID`) before creating plans, so a plan given a released number does not meet the old plan's work"
	return result, nil
}

// checkSelection validates how the caller named the records and returns the
// IDs sorted and without duplicates when they were named by ID.
func checkSelection(member string, all bool, ids []string, memberCount int) ([]string, error) {
	switch {
	case len(ids) > 0 && (all || member != ""):
		return nil, errors.New("--id names the records itself; pass it without --member or --all")
	case all && member != "":
		return nil, errors.New("pass --member or --all, not both")
	case all && memberCount != 1:
		return nil, fmt.Errorf("--all is for a workspace whose roster lists exactly one member; this one lists %d, so name whose records to delete with --member ID, or the records with --id ID", memberCount)
	case len(ids) > 0:
		for _, id := range ids {
			if !recordPattern.MatchString(id) {
				return nil, fmt.Errorf("nothing was deleted: %q is not an intent or plan ID, such as i100 or p1000", id)
			}
		}
		ids = slices.Clone(ids)
		slices.Sort(ids)
		return slices.Compact(ids), nil
	case !all:
		if err := Name(member); err != nil {
			return nil, err
		}
	}
	return ids, nil
}

// readLedgers reads every ID file, the legacy shared one included: a preview
// runs before any split and writes nothing, so it sees whichever layout is
// there. reserved holds every number any file reserves.
func (s *Store) readLedgers() (idFiles []string, ledgers map[string]Ledger, reserved map[string]bool, err error) {
	idFiles, err = s.IDFiles()
	if err != nil {
		return nil, nil, nil, err
	}
	ledgers = map[string]Ledger{}
	reserved = map[string]bool{}
	for _, rel := range idFiles {
		ledger, err := s.readIDFile(rel)
		if err != nil {
			return nil, nil, nil, err
		}
		ledgers[rel] = ledger
		for _, id := range append(slices.Clone(ledger.Intents), ledger.Plans...) {
			reserved[id] = true
		}
	}
	return idFiles, ledgers, reserved, nil
}

// readAllRecords reads every active and archived record, also keyed by ID.
func (s *Store) readAllRecords() ([]Record, map[string]Record, error) {
	paths, err := s.RecordPaths(true)
	if err != nil {
		return nil, nil, err
	}
	records := []Record{}
	byRecord := map[string]Record{}
	for _, path := range paths {
		record, err := s.readRecord(path)
		if err != nil {
			return nil, nil, err
		}
		records = append(records, record)
		byRecord[record.ID] = record
	}
	return records, byRecord, nil
}

// selectRecords returns the set of IDs the request selects and the records to
// list for them, in record order and then ID order for reservation-only IDs.
func selectRecords(records []Record, reserved map[string]bool, member string, all bool, ids []string) (map[string]bool, []DeletedRecord, error) {
	byID := len(ids) > 0
	named := map[string]bool{}
	for _, id := range ids {
		named[id] = true
	}
	selected := map[string]bool{}
	chosen := []DeletedRecord{}
	for _, record := range records {
		if all || byID && named[record.ID] || !byID && record.CreatedBy == member {
			selected[record.ID] = true
			chosen = append(chosen, DeletedRecord{record.ID, record.CreatedBy, deletionPath(record), deletedStatus(record)})
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
				chosen = append(chosen, DeletedRecord{ID: id, Status: reservationOnly})
			default:
				unknown = append(unknown, id)
			}
		}
		if len(unknown) > 0 {
			return nil, nil, fmt.Errorf("nothing was deleted: no record or reserved ID matches %s", strings.Join(unknown, ", "))
		}
	}
	if len(chosen) == 0 && (!all || len(reserved) == 0) {
		if all {
			return nil, nil, errors.New("there are no intents, plans, or reserved IDs to delete")
		}
		return nil, nil, fmt.Errorf("no intents or plans were created by member %s", member)
	}
	return selected, chosen, nil
}

// blockers collects what stands between the selection and its deletion.
type blockers struct {
	messages []string
	records  []BlockingRecord
	// unlink holds, per intent, the selected plans edited out of its list.
	unlink map[string][]string
	addIDs map[string]bool
}

func (b *blockers) block(record Record, link string) {
	owner := record.CreatedBy
	if owner == "" {
		owner = "no recorded author"
	}
	b.messages = append(b.messages, fmt.Sprintf("%s (%s) %s", record.ID, owner, link))
	b.records = append(b.records, BlockingRecord{record.ID, record.CreatedBy, deletionPath(record), deletedStatus(record), link})
	b.addIDs[record.ID] = true
}

// findBlockers finds the records outside the selection that point inside it.
// Nothing stops one member planning under another's intent, and a person
// naming records by ID names only some of them. Deleting the half that was
// selected would leave a record pointing at a number that may soon mean
// something else. A plan listed by its own intent is the one link that is
// edited out instead, and only for records named by ID.
func findBlockers(records []Record, byRecord map[string]Record, selected map[string]bool, byID bool) *blockers {
	b := &blockers{unlink: map[string][]string{}, addIDs: map[string]bool{}}
	for _, record := range records {
		if selected[record.ID] {
			continue
		}
		if selected[record.Intent] {
			b.block(record, "belongs to intent "+record.Intent)
		}
		for _, dep := range record.DependsOn {
			if selected[dep] {
				b.block(record, "depends on "+dep)
			}
		}
		for _, plan := range record.Plans {
			if !selected[plan] {
				continue
			}
			// A plan with no file left, as an interrupted deletion or create
			// leaves one, points at nothing, so its entry goes the same way.
			if byID && (byRecord[plan].Intent == record.ID || byRecord[plan].ID == "") {
				b.unlink[record.ID] = append(b.unlink[record.ID], plan)
				continue
			}
			b.block(record, "links plan "+plan)
		}
	}
	return b
}

// addWorktrees blocks on this machine's worktrees prepared for a selected plan.
func (b *blockers) addWorktrees(local Associations, selected map[string]bool) {
	for _, assoc := range local.Worktrees {
		if selected[assoc.Plan] {
			b.messages = append(b.messages, fmt.Sprintf("worktree %s in %s was prepared for %s; remove it with `worktree remove --repo %s --path %s` first (if Git no longer lists it, `worktree prune --repo %s` drops the leftover association), or keep the plan", assoc.Path, assoc.Repository, assoc.Plan, assoc.Repository, assoc.Path, assoc.Repository))
		}
	}
}

// refuse turns the blockers into the refusal that carries the preview.
func refuse(result Deletion, b *blockers, byID bool) (Deletion, error) {
	sort.Strings(b.messages)
	message := fmt.Sprintf("nothing was deleted, because these would be left pointing at deleted records: %s. Resolve each first: its author edits the link out or deletes the record", strings.Join(b.messages, "; "))
	if byID && len(b.addIDs) > 0 {
		add := slices.Sorted(maps.Keys(b.addIDs))
		more := make([]string, 0, len(add))
		for _, id := range add {
			more = append(more, "--id "+id)
			result.AddIDs = append(result.AddIDs, id)
		}
		message += fmt.Sprintf(", or delete them together by adding %s to this call", strings.Join(more, " "))
	}
	sort.SliceStable(b.records, func(i, j int) bool { return b.records[i].ID < b.records[j].ID })
	result.BlockedBy = b.records
	result.Refused = message
	return result, &DeletionRefused{Preview: result}
}

// releasedIDs lists, sorted, the numbers the deletion frees. Deleting
// everything frees every number in every ID file, including ones an
// interrupted create reserved and never wrote a file for.
func releasedIDs(idFiles []string, ledgers map[string]Ledger, selected map[string]bool, all bool) []string {
	released := []string{}
	seen := map[string]bool{}
	for _, rel := range idFiles {
		for _, id := range append(append([]string{}, ledgers[rel].Intents...), ledgers[rel].Plans...) {
			if (all || selected[id]) && !seen[id] {
				seen[id] = true
				released = append(released, id)
			}
		}
	}
	sort.Strings(released)
	return released
}

// leftoverPlans lists, sorted, the plan numbers whose Git leftovers are looked
// for: released plans and plans among the records.
func leftoverPlans(released []string, records []DeletedRecord) []string {
	var plans []string
	seen := map[string]bool{}
	for _, id := range released {
		if strings.HasPrefix(id, "p") {
			plans = append(plans, id)
			seen[id] = true
		}
	}
	for _, record := range records {
		if strings.HasPrefix(record.ID, "p") && !seen[record.ID] {
			plans = append(plans, record.ID)
			seen[record.ID] = true
		}
	}
	slices.Sort(plans)
	return plans
}

// unlinkPlans takes each plan out of its own intent's list first. Interrupted
// after that, the plan still stands and still names its intent, so the same
// call finishes the job; released first, an interruption would leave the
// intent listing a number nothing selects any more.
func (s *Store) unlinkPlans(unlinkFrom []string, byRecord map[string]Record, unlink map[string][]string) error {
	for _, id := range unlinkFrom {
		intent := byRecord[id]
		drop := map[string]bool{}
		for _, plan := range unlink[id] {
			drop[plan] = true
		}
		remaining := []string{}
		for _, plan := range intent.Plans {
			if !drop[plan] {
				remaining = append(remaining, plan)
			}
		}
		if err := s.editRecord(intent, "plans", remaining); err != nil {
			return fmt.Errorf("nothing was released: could not remove %s from intent %s: %w", strings.Join(unlink[id], ", "), id, err)
		}
	}
	return nil
}

// removeRecords removes the record files, and every number stays reserved
// until they are all gone: interrupted here, each selected ID still holds a
// file or a reservation, so --id reaches it and allocation never hands it out.
// Plans go before intents, so a record goes before anything it points to.
func (s *Store) removeRecords(records []DeletedRecord, reserved map[string]bool, all bool) error {
	for _, plansFirst := range []bool{true, false} {
		for _, record := range records {
			if record.Path == "" || strings.HasPrefix(record.ID, "p") != plansFirst {
				continue
			}
			if err := s.removeRecord(record.Path); err != nil {
				return s.unfinished(record.Path, records, reserved, all, err)
			}
		}
	}
	return nil
}

// releaseIDs rewrites the ID files last, dropping each released number from
// every file holding it. Deleting everything empties every ID file outright;
// the files stay, so each member keeps the file their next record is reserved
// in.
func (s *Store) releaseIDs(idFiles []string, ledgers map[string]Ledger, selected map[string]bool, released []string, all bool) error {
	keep := func(ids []string) []string {
		kept := []string{}
		for _, id := range ids {
			if !all && !selected[id] {
				kept = append(kept, id)
			}
		}
		return kept
	}
	for _, rel := range idFiles {
		ledger := ledgers[rel]
		intents, plans := keep(ledger.Intents), keep(ledger.Plans)
		if !slices.Equal(intents, ledger.Intents) {
			if err := writeIDList(s, rel, "intents", intents); err != nil {
				return s.unreleased(rel, idFiles, ledgers, released, all, err)
			}
		}
		if !slices.Equal(plans, ledger.Plans) {
			if err := writeIDList(s, rel, "plans", plans); err != nil {
				return s.unreleased(rel, idFiles, ledgers, released, all, err)
			}
		}
	}
	return nil
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
	releasedSet := map[string]bool{}
	for _, id := range released {
		releasedSet[id] = true
	}
	var left []string
	seen := map[string]bool{}
	for _, file := range idFiles {
		ledger, err := s.readIDFile(file)
		if err != nil {
			ledger = ledgers[file]
		}
		for _, id := range append(slices.Clone(ledger.Intents), ledger.Plans...) {
			if releasedSet[id] && !seen[id] {
				seen[id] = true
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
	seen := map[string]bool{}
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
		if held && !seen[record.ID] {
			seen[record.ID] = true
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
	var warnings []string
	var listed []string
	for _, repo := range slices.Sorted(maps.Keys(bindings.Bindings)) {
		checkout, err := s.Repository(ctx, repo)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s was not inspected for leftover branches or worktrees: %v", repo, err))
			continue
		}
		warnings = append(warnings, branchLeftovers(ctx, repo, checkout.Path, plans)...)
		treeWarnings, treePaths := s.worktreeLeftovers(ctx, repo, plans)
		warnings = append(warnings, treeWarnings...)
		listed = append(listed, treePaths...)
	}
	return append(warnings, s.directoryLeftovers(plans, listed)...)
}

// branchLeftovers names the cc/<plan>/* branches, local or remote-tracking,
// that one repository still holds for the plans.
func branchLeftovers(ctx context.Context, repo, path string, plans []string) []string {
	var warnings []string
	refs, err := Git(ctx, path, "for-each-ref", "--format=%(refname)", "refs/heads/cc", "refs/remotes")
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
	return warnings
}

// worktreeLeftovers names the worktrees of one repository on a cc/<plan>/*
// branch, and returns their paths so the default directory is not named twice.
func (s *Store) worktreeLeftovers(ctx context.Context, repo string, plans []string) (warnings, paths []string) {
	trees, err := s.Worktrees(ctx, repo)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("%s was not inspected for leftover worktrees: %v", repo, err))
	}
	for _, tree := range trees {
		for _, plan := range plans {
			if !tree.Main && strings.HasPrefix(tree.Branch, "cc/"+plan+"/") {
				paths = append(paths, tree.Path)
				warnings = append(warnings, fmt.Sprintf("worktree directory %s in %s is on branch %s", tree.Path, repo, tree.Branch))
			}
		}
	}
	return warnings, paths
}

// directoryLeftovers names each default .worktrees/<plan> directory that exists
// and is not inside a worktree already listed.
func (s *Store) directoryLeftovers(plans, listed []string) []string {
	var warnings []string
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
		// listed holds prefixes to match against, not exact names, so it stays
		// a scan; it is one entry per leftover worktree.
		covered := slices.ContainsFunc(listed, func(path string) bool {
			return strings.HasPrefix(filepath.Clean(path)+string(filepath.Separator), filepath.Clean(dir)+string(filepath.Separator))
		})
		if !covered {
			warnings = append(warnings, fmt.Sprintf("worktree directory %s is left from an earlier %s", dir, plan))
		}
	}
	return warnings
}
