package workspace

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// finished reports whether a record's work is done: a completed plan, or an
// approved intent with at least one plan, every one of them completed. An
// approved intent without plans still waits for planning, so it is open work.
// byID holds every record, archived ones included, because an intent's plans
// may already be archived.
func finished(record Record, byID map[string]Record) bool {
	if strings.HasPrefix(record.ID, "p") {
		return record.CompletedAt != ""
	}
	if record.ApprovedAt == "" {
		return false
	}
	plans := intentPlans(record, byID)
	if len(plans) == 0 {
		return false
	}
	for _, id := range plans {
		plan, ok := byID[id]
		if !ok || plan.CompletedAt == "" {
			return false
		}
	}
	return true
}

// intentPlans lists, sorted, the plans an intent lists and the plans naming
// it as their intent: either link alone makes a plan part of its work.
func intentPlans(intent Record, byID map[string]Record) []string {
	plans := slices.Clone(intent.Plans)
	for id, record := range byID {
		if record.Intent == intent.ID && !slices.Contains(plans, id) {
			plans = append(plans, id)
		}
	}
	slices.Sort(plans)
	return plans
}

// RecordSelection is whose records a listing shows. With neither Member nor
// All it is the active member's; a record with no author, written by an
// earlier version, is listed only with All.
type RecordSelection struct {
	Member   string
	All      bool
	Archived bool
}

// RecordListing is a listing and, when it could not show the active member's
// records because no member is selected, why it shows everyone's instead.
type RecordListing struct {
	Records []Record
	Notice  string
}

// SelectRecords lists the records a selection names. Reading is allowed
// before setup, so a machine with no member selected is shown every record
// and told why, rather than refused.
func (s *Store) SelectRecords(selection RecordSelection) (RecordListing, error) {
	if selection.All && selection.Member != "" {
		return RecordListing{}, errors.New("pass --member or --all, not both")
	}
	member, all, notice := selection.Member, selection.All, ""
	if member != "" {
		if err := Name(member); err != nil {
			return RecordListing{}, err
		}
	} else if !all {
		active, err := s.ActiveMember()
		if err != nil {
			all = true
			notice = fmt.Sprintf("notice: no member is selected on this machine, so this lists every member's records (%v); select one with `member use --id ID` to list your own\n", err)
		}
		member = active
	}
	records, err := s.ListRecords(selection.Archived)
	if err != nil {
		return RecordListing{}, err
	}
	if all {
		return RecordListing{records, notice}, nil
	}
	selected := []Record{}
	for _, record := range records {
		if record.CreatedBy == member {
			selected = append(selected, record)
		}
	}
	return RecordListing{selected, notice}, nil
}

// WorkItem is one record of a member's open work.
type WorkItem struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// MyWork is the active member's own open work, which status opens with: their
// intents waiting for approval, their approved intents with no plans yet,
// their plans not completed, and how many of their finished records are not
// archived. Nobody else's records are counted.
type MyWork struct {
	Member               string     `json:"member"`
	AwaitingApproval     []WorkItem `json:"awaiting_approval"`
	ApprovedWithoutPlans []WorkItem `json:"approved_without_plans"`
	OpenPlans            []WorkItem `json:"open_plans"`
	UnarchivedFinished   int        `json:"unarchived_finished"`
}

// allRecords reads every record, archived ones included, in path order, and
// which of them are archived.
func (s *Store) allRecords() ([]Record, map[string]bool, error) {
	paths, err := s.RecordPaths(true)
	if err != nil {
		return nil, nil, err
	}
	records := make([]Record, 0, len(paths))
	archived := map[string]bool{}
	for _, path := range paths {
		record, err := s.readRecord(path)
		if err != nil {
			return nil, nil, err
		}
		record.Content = ""
		records = append(records, record)
		if _, isArchived, _ := placement.Classify(path); isArchived {
			archived[record.ID] = true
		}
	}
	return records, archived, nil
}

func indexByID(records []Record) map[string]Record {
	byID := make(map[string]Record, len(records))
	for _, record := range records {
		byID[record.ID] = record
	}
	return byID
}

// unarchivedFinished lists the member's finished records that are not archived.
func unarchivedFinished(member string, records []Record, archived map[string]bool, byID map[string]Record) []Record {
	var out []Record
	for _, record := range records {
		if record.CreatedBy == member && !archived[record.ID] && finished(record, byID) {
			out = append(out, record)
		}
	}
	return out
}

// myWork reads the member's open work from every active record.
func (s *Store) myWork(member string) (*MyWork, error) {
	records, archived, err := s.allRecords()
	if err != nil {
		return nil, err
	}
	byID := indexByID(records)
	work := &MyWork{Member: member, AwaitingApproval: []WorkItem{}, ApprovedWithoutPlans: []WorkItem{}, OpenPlans: []WorkItem{}}
	for _, record := range records {
		if record.CreatedBy != member || archived[record.ID] {
			continue
		}
		item := WorkItem{record.ID, record.Path}
		switch {
		case strings.HasPrefix(record.ID, "p"):
			if record.CompletedAt == "" {
				work.OpenPlans = append(work.OpenPlans, item)
			}
		case record.ApprovedAt == "":
			work.AwaitingApproval = append(work.AwaitingApproval, item)
		case len(intentPlans(record, byID)) == 0:
			work.ApprovedWithoutPlans = append(work.ApprovedWithoutPlans, item)
		}
	}
	work.UnarchivedFinished = len(unarchivedFinished(member, records, archived, byID))
	return work, nil
}

// UnarchivedFinishedCount is how many of the member's finished records are not
// archived yet, which completion offers to archive.
func (s *Store) UnarchivedFinishedCount(member string) (int, error) {
	records, archived, err := s.allRecords()
	if err != nil {
		return 0, err
	}
	return len(unarchivedFinished(member, records, archived, indexByID(records))), nil
}
