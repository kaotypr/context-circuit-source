package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// banded builds a minimal workspace whose roster holds the given bands. A zero
// band means the member is unbanded.
func banded(t *testing.T, bands map[string]int, reserved Ledger) *Store {
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
	if err := store.WriteYAML("workspace.yaml", Config{Version: 2, Name: "Acme", Purpose: "Billing", Repositories: map[string]Repository{}, Relationships: []Relationship{}}, 0o644); err != nil {
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
	if err := store.WriteYAML(".context-circuit/ids.yaml", reserved, 0o644); err != nil {
		t.Fatal(err)
	}
	return store
}

// A band exists so two clones that cannot see each other still never choose the
// same number. That only holds if each member draws from its own block and an
// unbanded member stays out of every block.
func TestBandedAllocationUsesDisjointBlocks(t *testing.T) {
	store := banded(t, map[string]int{"maya": 1, "alex": 2, "sam": 0}, Ledger{})
	for _, tc := range []struct{ author, kind, want string }{
		{"maya", "intent", "i100"},
		{"maya", "plan", "p1000"},
		{"maya", "plan", "p1001"},
		{"alex", "intent", "i200"},
		{"alex", "plan", "p2000"},
		{"maya", "intent", "i101"},
	} {
		got, err := store.allocate(tc.kind, tc.author)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Fatalf("%s %s: got %s want %s", tc.author, tc.kind, got, tc.want)
		}
	}
}

// Without any band the workspace keeps allocating exactly as it did before
// bands existed, so a solo workspace pays nothing for the feature.
func TestUnbandedAllocationIsUnchanged(t *testing.T) {
	store := banded(t, map[string]int{"maya": 0}, Ledger{})
	for _, want := range []string{"i001", "i002", "i003"} {
		got, err := store.allocate("intent", "maya")
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	}
	for _, want := range []string{"p0001", "p0002"} {
		got, err := store.allocate("plan", "maya")
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	}
}

// An unbanded member walking into someone's band would defeat the whole point,
// so allocation steps over every declared block.
func TestUnbandedAllocationSkipsDeclaredBands(t *testing.T) {
	var full []string
	for n := 1; n <= 99; n++ {
		full = append(full, fmt.Sprintf("i%03d", n))
	}
	store := banded(t, map[string]int{"maya": 1, "sam": 0}, Ledger{Intents: full})
	got, err := store.allocate("intent", "sam")
	if err != nil {
		t.Fatal(err)
	}
	// 100-199 belongs to maya's band 1, so the next unbanded number is 200.
	if got != "i200" {
		t.Fatalf("unbanded allocation entered a band: %s", got)
	}
}

// A band decides which number comes next; it never makes a reserved number
// reusable, including one reserved inside the band itself.
func TestBandedAllocationSkipsReservedNumbers(t *testing.T) {
	store := banded(t, map[string]int{"maya": 1}, Ledger{Plans: []string{"p1000", "p1001"}})
	got, err := store.allocate("plan", "maya")
	if err != nil {
		t.Fatal(err)
	}
	if got != "p1002" {
		t.Fatalf("got %s want p1002", got)
	}
}

func TestExhaustedBandReportsTheMember(t *testing.T) {
	var full []string
	for n := 100; n <= 199; n++ {
		full = append(full, fmt.Sprintf("i%03d", n))
	}
	store := banded(t, map[string]int{"maya": 1}, Ledger{Intents: full})
	_, err := store.allocate("intent", "maya")
	if err == nil || !strings.Contains(err.Error(), "band 1") || !strings.Contains(err.Error(), "maya") {
		t.Fatalf("expected an exhaustion error naming the member and band, got %v", err)
	}
}

// Two members sharing a band share a range, which is the one thing a band must
// never allow. Refuse to read the roster at all until it is resolved.
func TestDuplicateBandsAreRefused(t *testing.T) {
	store := banded(t, map[string]int{"maya": 1, "alex": 1}, Ledger{})
	if _, err := store.Members(); err == nil || !strings.Contains(err.Error(), "band 1") {
		t.Fatalf("expected a duplicate-band error, got %v", err)
	}
}

func TestBandAssignmentGuardsCollisionAndClearing(t *testing.T) {
	store := banded(t, map[string]int{"maya": 1}, Ledger{})
	if err := store.AddMember("alex", "Alex", 1); err == nil {
		t.Fatal("added a member holding an occupied band")
	}
	if err := store.AddMember("alex", "Alex", 2); err != nil {
		t.Fatal(err)
	}
	if err := store.SetMemberBand("alex", 1); err == nil {
		t.Fatal("moved a member onto an occupied band")
	}
	if err := store.SetMemberBand("alex", 3); err != nil {
		t.Fatal(err)
	}
	members, err := store.Members()
	if err != nil {
		t.Fatal(err)
	}
	if members.Members["alex"].Band != 3 || members.Members["alex"].Name != "Alex" {
		t.Fatalf("band change lost data: %+v", members.Members["alex"])
	}
	// Clearing returns the member to the unbanded range without touching names.
	if err := store.SetMemberBand("alex", 0); err != nil {
		t.Fatal(err)
	}
	got, err := store.allocate("plan", "alex")
	if err != nil {
		t.Fatal(err)
	}
	if got != "p0001" {
		t.Fatalf("cleared band did not return to the unbanded range: %s", got)
	}
}

func TestBandRangeIsValidated(t *testing.T) {
	store := banded(t, map[string]int{"maya": 0}, Ledger{})
	for _, band := range []int{-1, 10000} {
		if err := store.SetMemberBand("maya", band); err == nil {
			t.Fatalf("accepted band %d", band)
		}
	}
}
