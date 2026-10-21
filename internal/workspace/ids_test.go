package workspace

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestIDFilePathValidatesTheMember(t *testing.T) {
	if got, err := IDFilePath("maya"); err != nil || got != ".context-circuit/ids/maya.yaml" {
		t.Fatalf("got %q, %v", got, err)
	}
	for _, bad := range []string{"", "Maya", "../maya", "maya.yaml", "a/b"} {
		if _, err := IDFilePath(bad); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
	// Every clone must be able to check out every ID file, Windows included.
	for _, device := range []string{"con", "prn", "aux", "nul", "com1", "com9", "lpt1", "lpt9"} {
		if _, err := IDFilePath(device); err == nil || !strings.Contains(err.Error(), "Windows device name") {
			t.Errorf("accepted device name %q: %v", device, err)
		}
	}
	for _, fine := range []string{"con-team", "com10", "console", "lpt"} {
		if _, err := IDFilePath(fine); err != nil {
			t.Errorf("refused %q: %v", fine, err)
		}
	}
}

func TestAddingAMemberNamedLikeADeviceIsRefused(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{})
	if err := store.AddMember("nul", "Nul", 2, "", "", nil); err == nil || !strings.Contains(err.Error(), "Windows device name") {
		t.Fatalf("added a member whose ID file no Windows clone could hold: %v", err)
	}
}

func TestIDFilesListsBothLayoutsInOrder(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{})
	for _, rel := range []string{IDFilesDir + "/zed.yaml", IDFilesDir + "/alex.yaml", IDFilesDir + "/notes.txt", LegacyIDFile} {
		if err := store.WriteYAML(rel, Ledger{Intents: []string{}, Plans: []string{}}, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(store.Root, IDFilesDir, "nested.yaml"), 0o755); err != nil {
		t.Fatal(err)
	}
	files, err := store.IDFiles()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{LegacyIDFile, IDFilesDir + "/alex.yaml", IDFilesDir + "/reserved.yaml", IDFilesDir + "/zed.yaml"}
	if !slices.Equal(files, want) {
		t.Fatalf("got %v, want %v", files, want)
	}
}

func TestIsIDFile(t *testing.T) {
	for rel, want := range map[string]bool{
		".context-circuit/ids.yaml":          true,
		".context-circuit/ids/maya.yaml":     true,
		".context-circuit/ids/a/b.yaml":      false,
		".context-circuit/ids/maya.txt":      false,
		".context-circuit/role-tiering.yaml": false,
		"members.yaml":                       false,
	} {
		if IsIDFile(rel) != want {
			t.Errorf("IsIDFile(%q) = %v", rel, !want)
		}
	}
}

func TestIDReservationsNameTheFileAtFault(t *testing.T) {
	for name, content := range map[string]string{
		"missing list":     "intents: []\n",
		"malformed ID":     "intents:\n  - x1\nplans: []\n",
		"wrong prefix":     "intents:\n  - p1000\nplans: []\n",
		"duplicate":        "intents:\n  - i100\n  - i100\nplans: []\n",
		"conflict markers": "intents:\n<<<<<<< HEAD\n  - i100\n=======\n  - i200\n>>>>>>> other\nplans: []\n",
	} {
		store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{})
		if err := store.Write(IDFilesDir+"/alex.yaml", []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := store.idReservations(); err == nil || !strings.Contains(err.Error(), "alex.yaml") {
			t.Errorf("%s: the error does not name the file: %v", name, err)
		}
	}
}

// legacy builds a schema 2 workspace holding the shared ID file and the given
// records, each "id:author" (an empty author is a record written before every
// record carried one).
func legacy(t *testing.T, reserved Ledger, records ...string) *Store {
	t.Helper()
	store := banded(t, map[string]*int{"maya": BandOf(1), "alex": BandOf(2)}, Ledger{})
	if err := os.Remove(filepath.Join(store.Root, IDFilesDir, "reserved.yaml")); err != nil {
		t.Fatal(err)
	}
	if err := store.Write("workspace.yaml", []byte("# kept\nversion: 2 # schema\nname: Acme\npurpose: Billing\nrepositories: {}\nrelationships: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := store.WriteYAML(LegacyIDFile, reserved, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, spec := range records {
		id, author, _ := strings.Cut(spec, ":")
		header := "---\nid: " + id + "\ncreated_at: 2026-09-15T10:53:00Z\n"
		if author != "" {
			header += "created_by: " + author + "\n"
		}
		rel := "intent/" + id + "-r.md"
		if strings.HasPrefix(id, "p") {
			rel = "plans/" + id + "-r/plan.md"
		}
		if strings.HasSuffix(id, "9") {
			// Archived records are split like any other.
			rel = strings.Replace(rel, "/", "/archive/", 1)
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(store.Root, rel)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(store.Root, rel), []byte(header+"---\n\n# R\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return store
}

func memberFile(t *testing.T, store *Store, member string) Ledger {
	t.Helper()
	rel, _ := IDFilePath(member)
	var l Ledger
	if err := store.YAML(rel, &l); err != nil {
		t.Fatal(err)
	}
	return l
}

func TestSplitMovesEveryReservationToItsOwner(t *testing.T) {
	store := legacy(t,
		Ledger{Intents: []string{"i001", "i100", "i200", "i209", "i050"}, Plans: []string{"p1000", "p0007"}},
		"i001:", "i100:maya", "i200:alex", "i209:alex", "p1000:maya", "p0007:ghost")
	split, err := store.SplitIDLedger("alex")
	if err != nil {
		t.Fatal(err)
	}
	// i001 has no author and i050 no record: both go to the member splitting.
	// i209 is archived. ghost is not on the roster and still gets a file.
	// Numbers outside the owner's band are accepted.
	for member, want := range map[string]string{"maya": "{[i100] [p1000]}", "alex": "{[i001 i200 i209 i050] []}", "ghost": "{[] [p0007]}"} {
		if got := fmtLedger(memberFile(t, store, member)); got != want {
			t.Errorf("%s: got %s want %s", member, got, want)
		}
	}
	if _, err := os.Stat(filepath.Join(store.Root, LegacyIDFile)); !os.IsNotExist(err) {
		t.Fatal("the shared ID file survived the split")
	}
	data, _ := store.Read("workspace.yaml")
	if !strings.Contains(string(data), "version: 3 # schema") || !strings.Contains(string(data), "# kept") {
		t.Fatalf("the schema was not moved to 3 with comments kept:\n%s", data)
	}
	if split.FromSchema != 2 || split.ToSchema != 3 || !split.LegacyRemoved || !strings.Contains(split.Notice(), "person's decision") {
		t.Fatalf("unexpected report: %+v\n%s", split, split.Notice())
	}
	if pending, err := store.IDSplitPending(); err != nil || pending {
		t.Fatalf("a finished split is still pending: %v %v", pending, err)
	}
}

// A split that finished is reported as finished even when listing the files
// it changed fails: the notice falls back to naming none rather than the split
// reading as undone.
func TestAFinishedSplitSurvivesAFailedChangeListing(t *testing.T) {
	store := legacy(t, Ledger{Intents: []string{"i100"}, Plans: []string{}}, "i100:maya")
	original := changedSince
	changedSince = func(*Store, SharedSnapshot) ([]string, error) { return nil, os.ErrPermission }
	t.Cleanup(func() { changedSince = original })
	split, err := store.SplitIDLedger("maya")
	if err != nil || split == nil {
		t.Fatalf("a finished split was reported as failed: %v", err)
	}
	if len(split.Files) != 0 || !strings.Contains(split.Notice(), "Commit and push the changed workspace files") {
		t.Fatalf("the notice did not fall back: %+v\n%s", split, split.Notice())
	}
	if pending, err := store.IDSplitPending(); err != nil || pending {
		t.Fatalf("the split did not finish: %v %v", pending, err)
	}
}

func fmtLedger(l Ledger) string {
	return "{[" + strings.Join(l.Intents, " ") + "] [" + strings.Join(l.Plans, " ") + "]}"
}

// A teammate on an older CLI reserved numbers before pulling the split, and
// merging brought the shared file back. Only the new numbers move.
func TestResplitSkipsWhatAMemberFileHolds(t *testing.T) {
	store := legacy(t, Ledger{Intents: []string{"i100"}, Plans: []string{}}, "i100:maya")
	if _, err := store.SplitIDLedger("maya"); err != nil {
		t.Fatal(err)
	}
	if err := store.WriteYAML(LegacyIDFile, Ledger{Intents: []string{"i100", "i200"}, Plans: []string{}}, 0o644); err != nil {
		t.Fatal(err)
	}
	if pending, _ := store.IDSplitPending(); !pending {
		t.Fatal("a returning shared file is not pending")
	}
	split, err := store.SplitIDLedger("alex")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(split.Skipped, []string{"i100"}) || split.FromSchema != 3 || split.ToSchema != 3 || strings.Contains(split.Notice(), "schema") {
		t.Fatalf("unexpected re-split: %+v\n%s", split, split.Notice())
	}
	if got := fmtLedger(memberFile(t, store, "maya")); got != "{[i100] []}" {
		t.Fatalf("maya: %s", got)
	}
	if got := fmtLedger(memberFile(t, store, "alex")); got != "{[i200] []}" {
		t.Fatalf("alex: %s", got)
	}
}

func TestSplitRefusesWhatItCannotReadWhole(t *testing.T) {
	for name, content := range map[string]string{
		"conflict":  "intents:\n<<<<<<< HEAD\n  - i100\n=======\n  - i200\n>>>>>>> other\nplans: []\n",
		"duplicate": "intents:\n  - i100\n  - i100\nplans: []\n",
	} {
		store := legacy(t, Ledger{Intents: []string{}, Plans: []string{}})
		if err := store.Write(LegacyIDFile, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := store.SplitIDLedger("maya"); err == nil || !strings.Contains(err.Error(), LegacyIDFile) {
			t.Errorf("%s: split did not refuse naming the file: %v", name, err)
		}
		if _, err := os.Stat(filepath.Join(store.Root, LegacyIDFile)); err != nil {
			t.Errorf("%s: a refused split removed the shared file", name)
		}
		if data, _ := store.Read("workspace.yaml"); !strings.Contains(string(data), "version: 2") {
			t.Errorf("%s: a refused split moved the schema", name)
		}
	}
}

// Interrupted after writing the member files, the workspace is still schema 2
// with the shared file present. Running again converges without duplicates.
func TestAnInterruptedSplitConverges(t *testing.T) {
	store := legacy(t, Ledger{Intents: []string{"i100", "i200"}, Plans: []string{}}, "i100:maya", "i200:alex")
	if err := store.reserve("maya", []string{"i100"}, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SplitIDLedger("maya"); err != nil {
		t.Fatal(err)
	}
	if got := fmtLedger(memberFile(t, store, "maya")); got != "{[i100] []}" {
		t.Fatalf("maya: %s", got)
	}
	if got := fmtLedger(memberFile(t, store, "alex")); got != "{[i200] []}" {
		t.Fatalf("alex: %s", got)
	}
	// A schema 2 workspace with no shared file only has its version moved.
	if err := store.Update("workspace.yaml", []string{"version"}, 2, 0o644); err != nil {
		t.Fatal(err)
	}
	split, err := store.SplitIDLedger("maya")
	if err != nil || split.ToSchema != 3 || len(split.Moved) != 0 || split.LegacyRemoved {
		t.Fatalf("unexpected: %+v %v", split, err)
	}
}

func TestSchemaRange(t *testing.T) {
	store := banded(t, map[string]*int{"maya": BandOf(1)}, Ledger{})
	for version, want := range map[int]string{2: "", 3: "", 4: "requires a newer CLI", 1: "v1 migration"} {
		if err := store.Update("workspace.yaml", []string{"version"}, version, 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := store.Config()
		if want == "" && err != nil || want != "" && (err == nil || !strings.Contains(err.Error(), want)) {
			t.Errorf("schema %d: got %v, want %q", version, err, want)
		}
	}
}
