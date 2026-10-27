package triggers

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
)

func TestCatalogue(t *testing.T) {
	table := loadTable(t)
	root := repoRoot(t)
	scenarios := scenarioIDs(t, filepath.Join(root, "agent-tests", "scenarios"))
	hostNames := hostIDs(t, filepath.Join(root, "agent-tests", "hosts.yaml"))
	if err := table.Validate(scenarios, hostNames); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"archive", "borrowed", "complete", "direct", "execute", "intent", "intent-lang", "join", "orient", "plan",
		"push-records", "push-records-conflict",
	}
	var got []string
	for id := range scenarios {
		got = append(got, id)
	}
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Fatalf("scenarios %v", got)
	}
	hit := table.StructureHit()
	if !slices.Equal(hit.Scenarios, want) {
		t.Fatalf("structure scenarios %v", hit.Scenarios)
	}
	if len(hit.Hosts) != 0 {
		t.Fatalf("structure hosts %v, want any one host", hit.Hosts)
	}
	if !slices.Equal(hit.OrientHosts, []string{"claude", "codex", "cursor"}) {
		t.Fatalf("orient hosts %v", hit.OrientHosts)
	}
	if !hit.Larger {
		t.Fatal("structure rows are larger")
	}
}

func TestSamplePaths(t *testing.T) {
	table := loadTable(t)
	all := []string{"archive", "borrowed", "complete", "direct", "execute", "intent", "intent-lang", "join", "orient", "plan", "push-records", "push-records-conflict"}
	cases := []struct {
		path      string
		scenarios []string
		hosts     []string
	}{
		{"product/AGENTS.md.in", []string{"orient"}, nil},
		{"product/CLAUDE.md.in", []string{"orient"}, []string{"claude"}},
		{"product/CURSOR.md.in", []string{"orient"}, []string{"cursor"}},
		{"product/cursor.mdc.in", []string{"orient"}, []string{"cursor"}},
		{"product/skills/cc-workspace/SKILL.md", []string{"join"}, nil},
		{"product/skills/cc-intent/SKILL.md", []string{"intent"}, nil},
		{"product/skills/cc-plan/SKILL.md", []string{"plan"}, nil},
		{"product/skills/cc-plan/route.md", []string{"plan"}, []string{"claude", "cursor"}},
		{"product/skills/cc-dispatch/SKILL.md", []string{"execute", "plan"}, nil},
		{"product/skills/cc-worktree/SKILL.md", []string{"execute"}, nil},
		{"product/skills/cc-knowledge/SKILL.md", []string{"borrowed", "complete"}, nil},
		{"product/skills/cc-complete/SKILL.md", []string{"archive", "complete"}, nil},
		{"product/skills/cc-direct/SKILL.md", []string{"direct"}, nil},
		{"product/skills/cc-cli/SKILL.md", []string{"orient"}, nil},
		{"product/skills/cc-cli/scripts/install.sh", []string{"orient"}, nil},
		{"product/skills/cc-cli/scripts/install.ps1", []string{"orient"}, nil},
		{"product/docs/commands.md", []string{"intent"}, nil},
		{"product/docs/working.md", []string{"complete", "intent", "plan"}, nil},
		{"product/docs/workspace.md", []string{"borrowed", "join"}, nil},
		{"product/docs/worktrees.md", []string{"execute"}, nil},
		{"product/docs/agents.md", []string{"execute", "plan"}, nil},
		{"product/docs/how-it-works.md", []string{"orient"}, nil},
		{"product/launcher/cli.sh", []string{"join", "orient"}, nil},
		{"product/workspace-README.md", []string{"orient"}, nil},
		{"template/workspace.yaml", []string{"join", "orient"}, nil},
		{"internal/cli/topics.go", []string{"orient"}, nil},
		{"internal/cli/commands.go", []string{"borrowed", "complete", "execute", "intent", "join", "plan"}, nil},
		{"internal/cli/cli.go", []string{"complete", "intent", "plan", "push-records"}, nil},
		{"internal/workspace/readme.go", []string{"join"}, nil},
		{"internal/workspace/records.go", []string{"intent", "plan"}, nil},
		{"internal/workspace/agents.go", []string{"execute", "plan"}, nil},
		{"scripts/release-manifest.txt", all, nil},
		{"assets.go", all, nil},
	}
	for _, tc := range cases {
		hit := table.Apply(tc.path, nil)
		if !slices.Equal(hit.Scenarios, tc.scenarios) {
			t.Errorf("%s scenarios %v, want %v", tc.path, hit.Scenarios, tc.scenarios)
		}
		if !slices.Equal(hit.Hosts, tc.hosts) {
			t.Errorf("%s hosts %v, want %v", tc.path, hit.Hosts, tc.hosts)
		}
		if slices.Contains(hit.Hosts, "codex") && strings.HasSuffix(tc.path, "/route.md") {
			t.Errorf("%s includes codex", tc.path)
		}
	}
	commands := table.Apply("internal/cli/commands.go", nil)
	note := strings.Join(commands.Notes, " ")
	if !strings.Contains(note, "Upper bound") || !strings.Contains(note, "narrows this list to the command that changed") {
		t.Fatalf("commands note %q", note)
	}
	route := table.Apply("product/skills/cc-plan/route.md", nil)
	if slices.Contains(route.Hosts, "codex") || !slices.Equal(route.Hosts, []string{"claude", "cursor"}) {
		t.Fatalf("route hosts %v", route.Hosts)
	}
	if route.Kinds[0] != "skill-route" {
		t.Fatalf("route kind %v", route.Kinds)
	}
}

func TestAGENTSSections(t *testing.T) {
	table := loadTable(t)
	root := repoRoot(t)
	body, err := os.ReadFile(filepath.Join(root, "product", "AGENTS.md.in"))
	if err != nil {
		t.Fatal(err)
	}
	var headings []string
	for _, line := range strings.Split(string(body), "\n") {
		if strings.HasPrefix(line, "## ") && !strings.HasPrefix(line, "### ") {
			headings = append(headings, strings.TrimPrefix(line, "## "))
		}
	}
	var listed []string
	for _, sec := range table.Sections {
		listed = append(listed, sec.Heading)
	}
	if !slices.Equal(headings, listed) {
		t.Fatalf("headings\nfile: %v\nmap: %v", headings, listed)
	}
	lang := table.Apply("product/AGENTS.md.in", []string{"Language"})
	if !slices.Equal(lang.Scenarios, []string{"intent-lang", "orient"}) {
		t.Fatalf("language %v", lang.Scenarios)
	}
	review := table.Apply("product/AGENTS.md.in", []string{"Independent review on request"})
	if !slices.Equal(review.Scenarios, []string{"orient"}) {
		t.Fatalf("review %v", review.Scenarios)
	}
	file := table.Apply("product/AGENTS.md.in", nil)
	if !slices.Equal(file.Scenarios, []string{"orient"}) || len(file.Hosts) != 0 {
		t.Fatalf("file row %+v", file)
	}
}

func TestNoneList(t *testing.T) {
	table := loadTable(t)
	for _, path := range []string{
		"README.md",
		"product/README.md",
		"internal/cow/copy.go",
		"internal/cli/cli_test.go",
		"VERSION",
		"CLI_VERSION",
		"scripts/build.sh",
		"agent-tests/internal/triggers/triggers.go",
	} {
		if !table.NoScenarios(path) {
			t.Errorf("expected no scenarios for %s", path)
		}
	}
	for _, path := range []string{
		"template/intent/README.md",
		"product/skills/cc-intent/SKILL.md",
		"scripts/release-manifest.txt",
	} {
		if table.NoScenarios(path) {
			t.Errorf("expected a scenario path for %s", path)
		}
	}
	for _, path := range []string{
		"product/skills/cc-review/SKILL.md",
		"product/skills/cc-deliver/route.md",
		"product/skills/cc-stacked/SKILL.md",
		"product/skills/cc-system-design/SKILL.md",
	} {
		if len(table.Match(path)) != 0 || table.NoScenarios(path) {
			t.Errorf("%s should be no row, not none", path)
		}
	}
}

func TestRowsMatchRealFiles(t *testing.T) {
	table := loadTable(t)
	root := repoRoot(t)
	matched := map[string]bool{}
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		for _, row := range table.Match(filepath.ToSlash(rel)) {
			matched[row.Match] = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range table.Rows {
		if !matched[row.Match] {
			t.Errorf("no file matches %s", row.Match)
		}
	}
}

func TestManifestSourcesAreClassified(t *testing.T) {
	table := loadTable(t)
	root := repoRoot(t)
	body, err := os.ReadFile(filepath.Join(root, "scripts", "release-manifest.txt"))
	if err != nil {
		t.Fatal(err)
	}
	unlisted := []string{
		"product/skills/cc-review/",
		"product/skills/cc-deliver/",
		"product/skills/cc-stacked/",
		"product/skills/cc-system-design/",
	}
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			t.Fatalf("manifest line %q", line)
		}
		source := fields[0]
		if len(table.Match(source)) > 0 || table.NoScenarios(source) {
			continue
		}
		ok := false
		for _, prefix := range unlisted {
			if strings.HasPrefix(source, prefix) {
				ok = true
			}
		}
		if !ok {
			t.Errorf("unclassified manifest source %s", source)
		}
	}
}

func loadTable(t *testing.T) Table {
	t.Helper()
	table, err := Load(filepath.Join(repoRoot(t), "agent-tests", "triggers.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return table
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "product", "AGENTS.md.in")); err != nil {
		t.Fatal(err)
	}
	return root
}

func scenarioIDs(t *testing.T, dir string) map[string]bool {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]bool{}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") {
			out[strings.TrimSuffix(name, ".yaml")] = true
		}
	}
	return out
}

func hostIDs(t *testing.T, path string) map[string]bool {
	t.Helper()
	loaded, err := hosts.LoadHosts(path)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]bool{}
	for name := range loaded {
		out[name] = true
	}
	return out
}
