package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// banded builds a minimal workspace whose roster holds the given bands. A nil
// band is a member with no band at all, which only an upgraded roster has. The
// reserved numbers sit in a member file nobody on the roster allocates into, so
// every allocation shows a number in another member's file is still taken.
func banded(t *testing.T, bands map[string]*int, reserved Ledger) *Store {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{".context-circuit", "intent", "plans"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	store := &Store{Root: root}
	members := Members{Members: map[string]Member{}}
	for id, band := range bands {
		members.Members[id] = Member{Name: strings.ToUpper(id), Band: band}
	}
	if err := store.WriteYAML("workspace.yaml", Config{Version: CurrentSchema, Name: "Acme", Purpose: "Billing", Repositories: map[string]Repository{}, Relationships: []Relationship{}}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := store.WriteYAML("members.yaml", members, 0o644); err != nil {
		t.Fatal(err)
	}
	if reserved.Intents == nil {
		reserved.Intents = []string{}
	}
	if reserved.Plans == nil {
		reserved.Plans = []string{}
	}
	if err := store.WriteYAML(IDFilesDir+"/reserved.yaml", reserved, 0o644); err != nil {
		t.Fatal(err)
	}
	return store
}

// A number one clone moved into its member's file and another clone into its
// own is reserved twice. That is check's to report, and allocation treats it
// as taken rather than refusing.
func TestANumberReservedInTwoFilesIsTaken(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1), "alex": BandOf(2)}, Ledger{Intents: []string{"i100"}})
	if err := store.WriteYAML(IDFilesDir+"/alex.yaml", Ledger{Intents: []string{"i100"}, Plans: []string{}}, 0o644); err != nil {
		t.Fatal(err)
	}
	allocates(t, store, "maya", "intent", "i101")
}

// A member's record creation writes their own file and no other, which is what
// keeps two members' clones from conflicting.
func TestAllocationWritesOnlyTheAuthorsFile(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1), "alex": BandOf(2)}, Ledger{Plans: []string{"p1000"}})
	before, err := store.Read(IDFilesDir + "/reserved.yaml")
	if err != nil {
		t.Fatal(err)
	}
	allocates(t, store, "maya", "intent", "i100")
	allocates(t, store, "maya", "plan", "p1001")
	allocates(t, store, "alex", "intent", "i200")
	if after, _ := store.Read(IDFilesDir + "/reserved.yaml"); string(after) != string(before) {
		t.Fatalf("allocating changed another member's file:\n%s", after)
	}
	var maya, alex Ledger
	if err := store.YAML(IDFilesDir+"/maya.yaml", &maya); err != nil {
		t.Fatal(err)
	}
	if err := store.YAML(IDFilesDir+"/alex.yaml", &alex); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(maya) != "{[i100] [p1001]}" || fmt.Sprint(alex) != "{[i200] []}" {
		t.Fatalf("unexpected member files: maya %v, alex %v", maya, alex)
	}
	// The header a new file opens with survives every later append.
	if data, _ := store.Read(IDFilesDir + "/maya.yaml"); !strings.HasPrefix(string(data), idFileHeader("maya")) {
		t.Fatalf("maya's file lost its header:\n%s", data)
	}
}

// A duplicate inside one file is a merge nobody resolved, and allocation still
// refuses it.
func TestADuplicateInsideOneFileIsRefused(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{Intents: []string{"i100", "i100"}})
	if _, err := store.allocate("intent", "maya"); err == nil || !strings.Contains(err.Error(), "reserved twice") {
		t.Fatalf("allocated through a duplicate reservation: %v", err)
	}
}

func allocates(t *testing.T, store *Store, author, kind string, want ...string) {
	t.Helper()
	for _, w := range want {
		got, err := store.allocate(kind, author)
		if err != nil {
			t.Fatal(err)
		}
		if got != w {
			t.Fatalf("%s %s: got %s want %s", author, kind, got, w)
		}
	}
}

// A band exists so two clones that cannot see each other still never choose the
// same number. That only holds if each member draws from its own block.
func TestBandedAllocationUsesDisjointBlocks(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1), "alex": BandOf(2), "arif": BandOf(5)}, Ledger{})
	allocates(t, store, "maya", "intent", "i100")
	allocates(t, store, "maya", "plan", "p1000", "p1001")
	allocates(t, store, "alex", "intent", "i200")
	allocates(t, store, "alex", "plan", "p2000")
	allocates(t, store, "maya", "intent", "i101")
	allocates(t, store, "arif", "intent", "i500")
	allocates(t, store, "arif", "plan", "p5000")
}

// Band 0 is a real band whose range starts at 1: there is never an i000 or a
// p0000.
func TestBandZeroStartsAtOne(t *testing.T) {
	store := banded(t, map[string]*int{"lutfi": BandOf(0), "maya": BandOf(1)}, Ledger{})
	allocates(t, store, "lutfi", "intent", "i001", "i002")
	allocates(t, store, "lutfi", "plan", "p0001", "p0002")
}

// Old records sitting inside a band, from before every member held one, keep
// their numbers, and allocation steps over them.
func TestBandZeroSkipsOldTakenNumbers(t *testing.T) {
	store := banded(t, map[string]*int{"lutfi": BandOf(0)}, Ledger{Intents: []string{"i001", "i002"}, Plans: []string{"p0001"}})
	if err := os.WriteFile(filepath.Join(store.Root, "intent", "i003-old.md"), []byte("---\nid: i003\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	allocates(t, store, "lutfi", "intent", "i004")
	allocates(t, store, "lutfi", "plan", "p0002")
}

// A band decides which number comes next; it never makes a reserved number
// reusable, including one reserved inside the band itself.
func TestBandedAllocationSkipsReservedNumbers(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{Plans: []string{"p1000", "p1001"}})
	allocates(t, store, "maya", "plan", "p1002")
}

func TestExhaustedBandReportsTheMember(t *testing.T) {
	var full []string
	for n := 100; n <= 199; n++ {
		full = append(full, fmt.Sprintf("i%03d", n))
	}
	store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{Intents: full})
	_, err := store.allocate("intent", "maya")
	if err == nil || !strings.Contains(err.Error(), "band 1") || !strings.Contains(err.Error(), "maya") || !strings.Contains(err.Error(), "100-199") || !strings.Contains(err.Error(), "member band") {
		t.Fatalf("expected an exhaustion error naming the member, band, range, and remedy, got %v", err)
	}
}

func TestAuthorWithoutABandIsRefused(t *testing.T) {
	store := banded(t, map[string]*int{"maya": nil}, Ledger{})
	if _, err := store.allocate("intent", "maya"); err == nil || !strings.Contains(err.Error(), "member band") {
		t.Fatalf("allocated for a member holding no band: %v", err)
	}
}

// Two members sharing a band share a range, which is the one thing a band must
// never allow, band 0 included. Refuse to read the roster until it is resolved.
func TestDuplicateBandsAreRefused(t *testing.T) {
	for _, band := range []int{0, 1} {
		store := banded(t, map[string]*int{"maya": BandOf(band), "alex": BandOf(band)}, Ledger{})
		want := fmt.Sprintf("members alex and maya both hold band %d", band)
		if _, err := store.Members(); err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("expected %q, got %v", want, err)
		}
	}
}

// goccy omits a nil pointer and must write a non-nil zero, or band 0 would be
// read back as no band at all.
func TestBandZeroSurvivesWriteAndReread(t *testing.T) {
	store := banded(t, map[string]*int{"lutfi": BandOf(0), "maya": nil}, Ledger{})
	data, err := store.Read("members.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "band: 0") {
		t.Fatalf("band 0 was not written:\n%s", data)
	}
	members, err := store.Members()
	if err != nil {
		t.Fatal(err)
	}
	if band := members.Members["lutfi"].Band; band == nil || *band != 0 {
		t.Fatalf("band 0 did not survive a re-read: %v", band)
	}
	if members.Members["maya"].Band != nil {
		t.Fatal("a member with no band was read back holding one")
	}
}

// A member given two bands by two clones is not something to read the last
// value of: goccy refuses a duplicate key, and the roster is reported.
func TestDuplicateBandKeyIsRefused(t *testing.T) {
	store := banded(t, map[string]*int{}, Ledger{})
	if err := store.Write("members.yaml", []byte("members:\n  maya:\n    name: Maya\n    band: 1\n    band: 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Members(); err == nil {
		t.Fatal("read a member holding two bands")
	}
	if err := store.Write("members.yaml", []byte("members:\n  maya:\n    name: Maya\n    band: 1\n  maya:\n    name: Maya\n    band: 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Members(); err == nil {
		t.Fatal("read a member listed twice")
	}
}

func TestConflictMarkersAreNamed(t *testing.T) {
	store := banded(t, map[string]*int{}, Ledger{})
	roster := "members:\n  maya:\n    name: Maya\n<<<<<<< HEAD\n    band: 1\n=======\n    band: 2\n>>>>>>> other\n"
	if err := store.Write("members.yaml", []byte(roster), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Members(); err == nil || !strings.Contains(err.Error(), "conflict markers") {
		t.Fatalf("expected the conflict to be named, got %v", err)
	}
}

func TestBandAssignmentGuardsCollision(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1), "lutfi": BandOf(0)}, Ledger{})
	for _, band := range []int{0, 1} {
		if err := store.AddMember("alex", "Alex", band, "", "", nil); err == nil {
			t.Fatalf("added a member holding occupied band %d", band)
		}
	}
	if err := store.AddMember("alex", "Alex", 2, "", "", nil); err != nil {
		t.Fatal(err)
	}
	if err := store.SetMemberBand("alex", 1); err == nil {
		t.Fatal("moved a member onto an occupied band")
	}
	if err := store.SetMemberBand("alex", 0); err == nil {
		t.Fatal("moved a member onto occupied band 0")
	}
	if err := store.SetMemberBand("alex", 3); err != nil {
		t.Fatal(err)
	}
	members, err := store.Members()
	if err != nil {
		t.Fatal(err)
	}
	if band := members.Members["alex"].Band; band == nil || *band != 3 || members.Members["alex"].Name != "Alex" {
		t.Fatalf("band change lost data: %+v", members.Members["alex"])
	}
}

func TestSettingBandZero(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{})
	if err := store.SetMemberBand("maya", 0); err != nil {
		t.Fatal(err)
	}
	allocates(t, store, "maya", "intent", "i001")
}

func TestBandRangeIsValidated(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{})
	for _, band := range []int{-1, 10000} {
		if err := store.SetMemberBand("maya", band); err == nil {
			t.Fatalf("accepted band %d", band)
		}
	}
}

// Suggestions start empty: never band 0, never a held band, never one offered
// twice, and never one whose range holds a taken number.
func TestBandSuggestions(t *testing.T) {
	store := banded(t, map[string]*int{"ana": nil, "ben": nil, "cy": BandOf(2), "dee": BandOf(0)}, Ledger{Intents: []string{"i001", "i120"}, Plans: []string{"p3500"}})
	suggestions, err := store.SuggestBands()
	if err != nil {
		t.Fatal(err)
	}
	// i120 blocks band 1, cy holds 2, p3500 blocks 3.
	if suggestions.Members["ana"] != 4 || suggestions.Members["ben"] != 5 || len(suggestions.Members) != 2 || suggestions.Next != 6 {
		t.Fatalf("unexpected suggestions: %+v", suggestions)
	}
}

func TestSetupNamesWhatIsOutstanding(t *testing.T) {
	store := banded(t, map[string]*int{}, Ledger{})
	if err := store.RequireSetup(); err == nil || !strings.Contains(err.Error(), "member add") {
		t.Fatalf("an empty roster passed setup: %v", err)
	}
	store = banded(t, map[string]*int{"maya": BandOf(1), "ana": nil}, Ledger{})
	if err := store.RequireSetup(); err == nil || !strings.Contains(err.Error(), "member use") {
		t.Fatalf("no selected member passed setup: %v", err)
	}
	if err := store.UseMember("maya"); err != nil {
		t.Fatal(err)
	}
	if err := store.RequireSetup(); err == nil || !strings.Contains(err.Error(), "ana band 2") {
		t.Fatalf("a member without a band passed setup: %v", err)
	}
	if err := store.SetMemberBand("ana", 2); err != nil {
		t.Fatal(err)
	}
	if err := store.RequireSetup(); err != nil {
		t.Fatal(err)
	}
}
