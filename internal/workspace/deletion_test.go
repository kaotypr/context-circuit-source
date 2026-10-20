package workspace

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// deletable writes an intent and its plan folder, both by maya and both
// reserved in her ID file, as a confirmed deletion finds them.
func deletable(t *testing.T) *Store {
	t.Helper()
	store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{})
	files := map[string]string{
		"intent/i100-billing.md":           "---\nid: i100\ncreated_by: maya\ncreated_at: 2026-09-15T10:53:00Z\nplans:\n- p1000\n---\n\n# Billing\n",
		"plans/p1000-billing-api/plan.md":  "---\nid: p1000\ncreated_by: maya\ncreated_at: 2026-09-15T10:54:00Z\nintent: i100\nrepositories:\n- api\nrequired_files:\n  shared: [notes.md]\n---\n\n# Billing API\n",
		"plans/p1000-billing-api/notes.md": "# Notes\n",
		IDFilesDir + "/maya.yaml":          "intents:\n  - i100\nplans:\n  - p1000\n",
	}
	for rel, content := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(store.Root, rel)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(store.Root, rel), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return store
}

// failOn makes removal of one path fail until the returned function restores it.
func failOn(t *testing.T, store *Store, rel string) {
	t.Helper()
	target := filepath.Join(store.Root, filepath.FromSlash(rel))
	removePath = func(path string) error {
		if filepath.Clean(path) == target {
			return errors.New("locked")
		}
		return os.RemoveAll(path)
	}
	t.Cleanup(func() { removePath = os.RemoveAll })
}

func mayaReserved(t *testing.T, store *Store) []string {
	t.Helper()
	ledger, err := store.readIDFile(IDFilesDir + "/maya.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return append(slices.Clone(ledger.Intents), ledger.Plans...)
}

func present(store *Store, rel string) bool {
	_, err := os.Lstat(filepath.Join(store.Root, filepath.FromSlash(rel)))
	return err == nil
}

// Interrupted inside a plan folder, the plan still holds plan.md and every
// number is still reserved, so the call the error names finishes the job.
func TestAnInterruptedDeletionByIDIsFinishedByRunningItAgain(t *testing.T) {
	store := deletable(t)
	ctx := context.Background()
	failOn(t, store, "plans/p1000-billing-api/plan.md")
	_, err := store.DeleteRecords(ctx, "", false, []string{"i100", "p1000"}, true)
	if err == nil || !strings.Contains(err.Error(), "`record delete --id i100 --id p1000 --confirm` to finish") {
		t.Fatalf("the interrupted deletion did not name the call that finishes it: %v", err)
	}
	if present(store, "plans/p1000-billing-api/notes.md") || !present(store, "plans/p1000-billing-api/plan.md") || !present(store, "intent/i100-billing.md") {
		t.Fatal("plan.md must go last, after the plan's other files and before its intent")
	}
	if got := mayaReserved(t, store); !slices.Equal(got, []string{"i100", "p1000"}) {
		t.Fatalf("an interrupted deletion released numbers: %v", got)
	}

	removePath = os.RemoveAll
	done, err := store.DeleteRecords(ctx, "", false, []string{"i100", "p1000"}, true)
	if err != nil || !done.Deleted {
		t.Fatalf("running the deletion again did not finish it: %+v %v", done, err)
	}
	if present(store, "plans/p1000-billing-api") || present(store, "intent/i100-billing.md") || len(mayaReserved(t, store)) != 0 {
		t.Fatal("the re-run left files or reservations behind")
	}
}

// A --member deletion interrupted after its plan is gone cannot reach the
// plan's number by member again, so the error names the --id call instead.
func TestAnInterruptedDeletionByMemberNamesTheIDsLeft(t *testing.T) {
	store := deletable(t)
	ctx := context.Background()
	failOn(t, store, "intent/i100-billing.md")
	_, err := store.DeleteRecords(ctx, "maya", false, nil, true)
	if err == nil || !strings.Contains(err.Error(), "`record delete --id i100 --id p1000 --confirm`") {
		t.Fatalf("the interrupted deletion did not name the IDs left: %v", err)
	}
	if present(store, "plans/p1000-billing-api") || !present(store, "intent/i100-billing.md") {
		t.Fatal("the plan must go before the intent it points to")
	}
	if got := mayaReserved(t, store); !slices.Equal(got, []string{"i100", "p1000"}) {
		t.Fatalf("an interrupted deletion released numbers: %v", got)
	}

	removePath = os.RemoveAll
	done, err := store.DeleteRecords(ctx, "", false, []string{"i100", "p1000"}, true)
	if err != nil || !done.Deleted {
		t.Fatalf("the named call did not finish the deletion: %+v %v", done, err)
	}
	if r := done.Records; len(r) != 2 || !slices.ContainsFunc(r, func(d DeletedRecord) bool { return d.ID == "p1000" && d.Status == reservationOnly }) {
		t.Fatalf("the removed plan's number was not reached as a reservation: %+v", r)
	}
	if present(store, "intent/i100-billing.md") || len(mayaReserved(t, store)) != 0 {
		t.Fatal("the named call left files or reservations behind")
	}
}

// Interrupted between ID files, the error names exactly the IDs still reserved
// and the call that releases them, and that call finishes the job.
func TestADeletionInterruptedInTheIDFilesNamesTheIDsStillReserved(t *testing.T) {
	store := deletable(t)
	ctx := context.Background()
	original := writeIDList
	writeIDList = func(s *Store, rel, key string, ids []string) error {
		if key == "plans" {
			return errors.New("disk full")
		}
		return original(s, rel, key, ids)
	}
	t.Cleanup(func() { writeIDList = original })
	_, err := store.DeleteRecords(ctx, "", false, []string{"i100", "p1000"}, true)
	if err == nil || !strings.Contains(err.Error(), "still reserved: p1000.") || !strings.Contains(err.Error(), "`record delete --id p1000 --confirm` to finish") {
		t.Fatalf("the interrupted deletion did not name the IDs still reserved: %v", err)
	}
	if got := mayaReserved(t, store); !slices.Equal(got, []string{"p1000"}) {
		t.Fatalf("want only p1000 still reserved, got %v", got)
	}

	writeIDList = original
	done, err := store.DeleteRecords(ctx, "", false, []string{"p1000"}, true)
	if err != nil || !done.Deleted || len(mayaReserved(t, store)) != 0 {
		t.Fatalf("the named call did not finish the deletion: %+v %v", done, err)
	}
}

// After --all, the unfinished hint re-runs --all, which also releases numbers
// no record file holds; and a plan folder whose plan.md is already gone is
// named for removal by hand, because no re-run reaches it.
func TestAnInterruptedDeletionOfEverythingSaysToRunItAgain(t *testing.T) {
	store := deletable(t)
	failOn(t, store, "plans/p1000-billing-api")
	_, err := store.DeleteRecords(context.Background(), "", true, nil, true)
	if err == nil || !strings.Contains(err.Error(), "`record delete --all --confirm` to finish") || strings.Contains(err.Error(), "--id") {
		t.Fatalf("the interrupted --all deletion did not say to run --all again: %v", err)
	}
	if !strings.Contains(err.Error(), "remove the folder plans/p1000-billing-api/ by hand") {
		t.Fatalf("the folder left without plan.md was not named: %v", err)
	}
}

// Unlinking a plan from its intent keeps a CRLF intent CRLF.
func TestUnlinkingAPlanKeepsCRLF(t *testing.T) {
	store := deletable(t)
	intent := filepath.Join(store.Root, "intent", "i100-billing.md")
	data, err := os.ReadFile(intent)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(intent, []byte(strings.ReplaceAll(string(data), "\n", "\r\n")), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DeleteRecords(context.Background(), "", false, []string{"p1000"}, true); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(intent)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(data), "\n") != strings.Count(string(data), "\r\n") || strings.Contains(string(data), "p1000") {
		t.Fatalf("the intent lost its CRLF line endings or kept the plan:\n%q", data)
	}
}
