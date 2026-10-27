package workspace

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// shelfPlacement puts every record one folder deeper, under a "shelf" folder
// in each root. It exists to show that callers follow the placement seam
// rather than a layout of their own.
type shelfPlacement struct{ flatPlacement }

func (shelfPlacement) NewRecordPath(kind, id, slug string) string {
	if kind == "plan" {
		return "plans/shelf/" + id + "-" + slug + "/plan.md"
	}
	return "intent/shelf/" + id + "-" + slug + ".md"
}

func (shelfPlacement) Classify(rel string) (string, bool, bool) {
	root, rest, ok := strings.Cut(rel, "/")
	shelved, ok2 := strings.CutPrefix(rest, "shelf/")
	if !ok || !ok2 {
		return "", false, false
	}
	return flatPlacement{}.Classify(root + "/" + shelved)
}

func (p shelfPlacement) RecordPaths(s *Store, archived bool) ([]string, error) {
	var paths []string
	for _, root := range p.Roots() {
		err := filepath.WalkDir(filepath.Join(s.Root, root), func(path string, entry os.DirEntry, err error) error {
			if err != nil || entry.IsDir() {
				return err
			}
			rel, _ := filepath.Rel(s.Root, path)
			rel = filepath.ToSlash(rel)
			if _, isArchived, ok := p.Classify(rel); ok && (archived || !isArchived) {
				paths = append(paths, rel)
			}
			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	return paths, nil
}

func (shelfPlacement) GitPathMatcher(id string) func(string) bool {
	return func(rel string) bool {
		root, rest, _ := strings.Cut(rel, "/")
		return (root == "intent" || root == "plans") && strings.HasPrefix(rest, "shelf/"+id+"-")
	}
}

func TestCallersFollowThePlacementSeam(t *testing.T) {
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(t.TempDir(), "gitconfig"))
	previous := placement
	placement = shelfPlacement{}
	t.Cleanup(func() { placement = previous })

	root := t.TempDir()
	s := &Store{Root: root}
	// No seed files are exported, so agent setup has no role settings to read;
	// the records this test needs are written before that step.
	if err := s.Init(map[string][]byte{}, nil, "Seam", "Show the placement seam", "maya", "Maya", 1, nil); err != nil && !strings.HasPrefix(err.Error(), "workspace initialized") {
		t.Fatal(err)
	}
	cfg, err := s.Config()
	if err != nil {
		t.Fatal(err)
	}
	cfg.Repositories["api"] = Repository{DefaultBranch: "main"}
	if err := s.WriteYAML("workspace.yaml", cfg, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{"intent/shelf", "plans/shelf"} {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(dir)), 0755); err != nil {
			t.Fatal(err)
		}
	}
	testGit(t, root, "init", "-b", "main")

	ctx := context.Background()
	intent, err := s.CreateRecord(ctx, "intent", "billing", "Billing", "", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := "intent/shelf/" + intent.ID + "-billing.md"; intent.Path != want {
		t.Fatalf("intent placed at %s, want %s", intent.Path, want)
	}
	plan, err := s.CreateRecord(ctx, "plan", "api", "API", intent.ID, []string{"api"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := "plans/shelf/" + plan.ID + "-api/plan.md"; plan.Path != want {
		t.Fatalf("plan placed at %s, want %s", plan.Path, want)
	}

	paths, err := s.RecordPaths(true)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{intent.Path, plan.Path}; !reflect.DeepEqual(paths, want) {
		t.Fatalf("found %v, want %v", paths, want)
	}
	for _, record := range []Record{intent, plan} {
		found, err := s.FindRecord(record.ID)
		if err != nil {
			t.Fatal(err)
		}
		if found.Path != record.Path {
			t.Fatalf("found %s at %s, want %s", record.ID, found.Path, record.Path)
		}
	}

	files, unchecked := s.UncommittedRecords(ctx, plan.ID)
	if unchecked != "" {
		t.Fatalf("git unavailable: %s", unchecked)
	}
	var listed []string
	for _, file := range files {
		if !IsIDFile(file.Path) {
			listed = append(listed, file.Path)
		}
	}
	if want := []string{plan.Path}; !reflect.DeepEqual(listed, want) {
		t.Fatalf("uncommitted records %v, want %v", listed, want)
	}
}

func TestFlatPlacementPaths(t *testing.T) {
	p := flatPlacement{}
	for _, tc := range []struct {
		rel      string
		id       string
		archived bool
		ok       bool
	}{
		{"intent/i001-billing.md", "i001", false, true},
		{"intent/archive/i001-billing.md", "i001", true, true},
		{"plans/p0001-api/plan.md", "p0001", false, true},
		{"plans/archive/p0001-api/plan.md", "p0001", true, true},
		{"plans/p0001-api.md", "p0001", false, true},
		{"plans/p0001-api/notes.md", "", false, false},
		{"intent/i001-billing/plan.md", "", false, false},
		{"intent/sub/i001-billing.md", "", false, false},
		{"context/i001-billing.md", "", false, false},
	} {
		id, archived, ok := p.Classify(tc.rel)
		if id != tc.id || archived != tc.archived || ok != tc.ok {
			t.Errorf("Classify(%s) = %q %v %v, want %q %v %v", tc.rel, id, archived, ok, tc.id, tc.archived, tc.ok)
		}
	}
	for rel, want := range map[string]string{
		"intent/i001-billing.md":         "intent/archive/i001-billing.md",
		"plans/p0001-api/plan.md":        "plans/archive/p0001-api/plan.md",
		"intent/archive/i001-billing.md": "intent/archive/i001-billing.md",
	} {
		if got := p.ArchivedPath(rel); got != want {
			t.Errorf("ArchivedPath(%s) = %s, want %s", rel, got, want)
		}
	}
	if got := p.NewRecordPath("intent", "i001", "billing"); got != "intent/i001-billing.md" {
		t.Errorf("intent placed at %s", got)
	}
	if got := p.NewRecordPath("plan", "p0001", "api"); got != "plans/p0001-api/plan.md" {
		t.Errorf("plan placed at %s", got)
	}
}
