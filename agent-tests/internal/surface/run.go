package surface

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/triggers"
)

// Options selects the two refs. Ref defaults to HEAD. Path is the PATH
// given to git and to the Go build.
type Options struct {
	Checkout     string
	BaseRef      string
	Ref          string
	OutDir       string
	TriggersFile string
	HostsFile    string
	ScenariosDir string
	Path         string
	Home         string
	Stdout       io.Writer
}

// Run builds both refs, compares the workspaces and the CLI script, and
// writes surface.json and surface.md. It exits successfully for the
// caller either way: an empty surface is a result, not an error.
func Run(ctx context.Context, opts Options) (Report, error) {
	if opts.Checkout == "" {
		return Report{}, fmt.Errorf("checkout is empty")
	}
	if opts.BaseRef == "" {
		return Report{}, fmt.Errorf("base ref is empty")
	}
	if opts.Ref == "" {
		opts.Ref = "HEAD"
	}
	if opts.OutDir == "" {
		return Report{}, fmt.Errorf("output directory is empty")
	}
	if opts.Path == "" {
		opts.Path = os.Getenv("PATH")
	}
	if opts.TriggersFile == "" {
		opts.TriggersFile = filepath.Join(opts.Checkout, "agent-tests", "triggers.yaml")
	}
	if opts.HostsFile == "" {
		opts.HostsFile = filepath.Join(opts.Checkout, "agent-tests", "hosts.yaml")
	}
	if opts.ScenariosDir == "" {
		opts.ScenariosDir = filepath.Join(opts.Checkout, "agent-tests", "scenarios")
	}
	out := opts.Stdout
	if out == nil {
		out = os.Stdout
	}
	table, err := triggers.Load(opts.TriggersFile)
	if err != nil {
		return Report{}, err
	}
	if err := validateTable(table, opts.HostsFile, opts.ScenariosDir); err != nil {
		return Report{}, err
	}
	baseSHA, err := revParse(ctx, opts.Checkout, opts.Path, opts.BaseRef)
	if err != nil {
		return Report{}, err
	}
	branchSHA, err := revParse(ctx, opts.Checkout, opts.Path, opts.Ref)
	if err != nil {
		return Report{}, err
	}
	id := time.Now().UTC().Format("20060102T150405Z") + "-" + short(baseSHA) + "-" + short(branchSHA)
	dir := filepath.Join(opts.OutDir, id, "surface")
	report := Report{
		ID:        id,
		Base:      baseSHA,
		Branch:    branchSHA,
		BaseRef:   opts.BaseRef,
		BranchRef: opts.Ref,
		Mapping:   mappingNote,
		Dir:       dir,
	}
	if baseSHA == branchSHA {
		report.Empty = true
		report.Message = emptyMessage
		report = finish(report)
		if err := writeReport(dir, report); err != nil {
			return Report{}, err
		}
		fmt.Fprint(out, Render(report))
		return report, nil
	}
	baseSide, err := captureSide(ctx, opts.Checkout, baseSHA, opts.Path, opts.Home)
	if err != nil {
		return Report{}, fmt.Errorf("base: %w", err)
	}
	branchSide, err := captureSide(ctx, opts.Checkout, branchSHA, opts.Path, opts.Home)
	if err != nil {
		return Report{}, fmt.Errorf("branch: %w", err)
	}
	report, err = compareSides(ctx, report, table, opts, baseSide, branchSide)
	if err != nil {
		return Report{}, err
	}
	if err := writeReport(dir, report); err != nil {
		return Report{}, err
	}
	if err := writeTree(filepath.Join(dir, "base"), baseSide.files); err != nil {
		return Report{}, err
	}
	if err := writeTree(filepath.Join(dir, "branch"), branchSide.files); err != nil {
		return Report{}, err
	}
	if err := writeOutputs(dir, baseSide, branchSide); err != nil {
		return Report{}, err
	}
	fmt.Fprint(out, Render(report))
	return report, nil
}

func compareSides(ctx context.Context, report Report, table triggers.Table, opts Options, baseSide, branchSide captured) (Report, error) {
	baseNorm := normalizeTree(baseSide.files, baseSide.norm)
	branchNorm := normalizeTree(branchSide.files, branchSide.norm)
	files := CompareTrees(baseNorm, branchNorm)
	outputs := pairOutputs(baseSide, branchSide)
	sizes := measureSizes(baseSide.files, branchSide.files)
	changedOut := changedOutputs(outputs)
	empty := len(files) == 0 && len(changedOut) == 0 && sizesFlat(sizes)
	report.Files = files
	report.Outputs = changedOut
	report.Sizes = sizes
	if empty {
		report.Empty = true
		report.Message = emptyMessage
		return finish(report), nil
	}
	report.Message = summaryMessage(files, outputs)
	gitPaths, err := nameStatus(ctx, opts.Checkout, opts.Path, report.Base, report.Branch)
	if err != nil {
		return Report{}, err
	}
	baseManifest := manifestAt(ctx, opts.Checkout, opts.Path, report.Base)
	branchManifest := manifestAt(ctx, opts.Checkout, opts.Path, report.Branch)
	selected := SelectSources(table, SourceInput{
		GitPaths:       gitPaths,
		Files:          files,
		ManifestBase:   baseManifest,
		ManifestBranch: branchManifest,
		OutputsChanged: len(changedOut) > 0,
	})
	report.Ignored = selected.Ignored
	headings := headingsFor(ctx, opts, report, selected.Paths)
	var scenarios, hosts, kinds, notes, noRow []string
	var larger bool
	for _, source := range selected.Paths {
		hit := table.Apply(source, headings)
		if len(hit.Rows) == 0 {
			noRow = append(noRow, source)
			continue
		}
		extra := helpScenariosFor(source, outputs)
		scenarios = append(scenarios, hit.Scenarios...)
		scenarios = append(scenarios, extra...)
		hosts = append(hosts, hit.Hosts...)
		kinds = append(kinds, hit.Kinds...)
		notes = append(notes, hit.Notes...)
		if hit.Larger {
			larger = true
		}
		report.Sources = append(report.Sources, SourceHit{
			Path:      source,
			Scenarios: uniqueSorted(append(append([]string{}, hit.Scenarios...), extra...)),
			Hosts:     hit.Hosts,
			Kinds:     hit.Kinds,
			Notes:     hit.Notes,
		})
	}
	if selected.EveryScenario {
		structure := table.StructureHit()
		scenarios = append(scenarios, structure.Scenarios...)
		hosts = append(hosts, structure.Hosts...)
		kinds = append(kinds, structure.Kinds...)
		notes = append(notes, structure.Notes...)
		report.OrientHosts = append(report.OrientHosts, structure.OrientHosts...)
		larger = true
	}
	hosts = append(hosts, roleHosts(files)...)
	report.Scenarios = uniqueSorted(scenarios)
	report.Hosts = uniqueSorted(hosts)
	report.Kinds = uniqueSorted(kinds)
	report.NoRow = uniqueSorted(noRow)
	report.Notes = uniqueSorted(notes)
	reasons := collectReasons(files, outputs, baseNorm, branchNorm, larger)
	report.LargerReasons = reasons
	report.Larger = len(reasons) > 0 || larger
	return finish(report), nil
}

func finish(report Report) Report {
	if report.Files == nil {
		report.Files = []FileChange{}
	}
	if report.Outputs == nil {
		report.Outputs = []OutputChange{}
	}
	if report.Scenarios == nil {
		report.Scenarios = []string{}
	}
	if report.Hosts == nil {
		report.Hosts = []string{}
	}
	if report.Kinds == nil {
		report.Kinds = []string{}
	}
	if report.LargerReasons == nil {
		report.LargerReasons = []Reason{}
	}
	if report.NoRow == nil {
		report.NoRow = []string{}
	}
	if report.Ignored == nil {
		report.Ignored = []string{}
	}
	if report.Sources == nil {
		report.Sources = []SourceHit{}
	}
	if report.Notes == nil {
		report.Notes = []string{}
	}
	if report.OrientHosts == nil {
		report.OrientHosts = []string{}
	}
	return report
}

func collectReasons(files []FileChange, outputs []OutputChange, base, branch map[string][]byte, rowLarger bool) []Reason {
	var reasons []Reason
	reasons = append(reasons, LayoutReasons(files)...)
	reasons = append(reasons, SkillReasons(base, branch)...)
	reasons = append(reasons, CLIReasons(outputs)...)
	var lines []string
	for _, file := range files {
		before, after := fileBytes(file, base, branch)
		lines = append(lines, ChangedLines(string(before), string(after))...)
	}
	for _, out := range outputs {
		if out.Changed {
			lines = append(lines, ChangedLines(out.Base, out.Branch)...)
		}
	}
	reasons = append(reasons, GateReasons(lines)...)
	reasons = append(reasons, MigrationReasons(base, branch)...)
	if rowLarger && !hasKind(reasons, "layout") {
		reasons = append(reasons, Reason{Kind: "layout", Exact: true, Detail: "trigger row marks a structure change as larger"})
	}
	if reasons == nil {
		return []Reason{}
	}
	return reasons
}

func fileBytes(file FileChange, base, branch map[string][]byte) (before, after []byte) {
	switch file.Status {
	case "moved":
		return base[file.From], branch[file.Path]
	case "removed":
		return base[file.Path], nil
	case "added":
		return nil, branch[file.Path]
	default:
		return base[file.Path], branch[file.Path]
	}
}

func hasKind(reasons []Reason, kind string) bool {
	for _, reason := range reasons {
		if reason.Kind == kind {
			return true
		}
	}
	return false
}

func headingsFor(ctx context.Context, opts Options, report Report, sources []string) []string {
	for _, source := range sources {
		if source != "product/AGENTS.md.in" {
			continue
		}
		before, _, err1 := gitShow(ctx, opts.Checkout, opts.Path, report.Base, "product/AGENTS.md.in")
		after, _, err2 := gitShow(ctx, opts.Checkout, opts.Path, report.Branch, "product/AGENTS.md.in")
		if err1 != nil || err2 != nil {
			return nil
		}
		return triggers.Sections(triggers.SectionDiff(before, after))
	}
	return nil
}

func manifestAt(ctx context.Context, checkout, pathEnv, rev string) map[string]string {
	text, ok, err := gitShow(ctx, checkout, pathEnv, rev, "scripts/release-manifest.txt")
	if err != nil || !ok {
		return WithReadme(nil)
	}
	parsed, err := ParseManifest(text)
	if err != nil {
		return WithReadme(nil)
	}
	return WithReadme(parsed)
}

func pairOutputs(baseSide, branchSide captured) []OutputChange {
	baseBy := map[string]capturedOutput{}
	for _, out := range baseSide.outputs {
		baseBy[out.name] = out
	}
	branchBy := map[string]capturedOutput{}
	for _, out := range branchSide.outputs {
		branchBy[out.name] = out
	}
	names := map[string]bool{}
	for name := range baseBy {
		names[name] = true
	}
	for name := range branchBy {
		names[name] = true
	}
	var ordered []string
	for name := range names {
		ordered = append(ordered, name)
	}
	sort.Strings(ordered)
	var out []OutputChange
	for _, name := range ordered {
		baseText := baseSide.norm.Text(baseBy[name].text)
		branchText := branchSide.norm.Text(branchBy[name].text)
		item := OutputChange{
			Name:         name,
			Changed:      baseText != branchText,
			Base:         baseText,
			Branch:       branchText,
			StdoutBase:   baseBy[name].stdout,
			StdoutBranch: branchBy[name].stdout,
		}
		if item.Changed {
			item.Diff = Unified(name, baseText, branchText)
		}
		out = append(out, item)
	}
	return out
}

func normalizeTree(raw map[string][]byte, n Norm) map[string][]byte {
	out := make(map[string][]byte, len(raw))
	for rel, body := range raw {
		out[rel] = n.Bytes(body)
	}
	return out
}

var alwaysLoaded = []string{
	"AGENTS.md",
	"CLAUDE.md",
	"CURSOR.md",
	".cursor/rules/context-circuit.mdc",
	"README.md",
}

func measureSizes(base, branch map[string][]byte) []SizeChange {
	paths := append([]string{}, alwaysLoaded...)
	seen := map[string]bool{}
	for _, path := range paths {
		seen[path] = true
	}
	for _, tree := range []map[string][]byte{base, branch} {
		for rel := range tree {
			if name, ok := skillName(rel); ok {
				full := ".agents/skills/" + name + "/SKILL.md"
				if !seen[full] {
					seen[full] = true
					paths = append(paths, full)
				}
			}
		}
	}
	sort.Strings(paths[len(alwaysLoaded):])
	var out []SizeChange
	for _, rel := range paths {
		item := SizeChange{Path: rel, BaseBytes: len(base[rel]), BranchBytes: len(branch[rel])}
		item.DeltaBytes = item.BranchBytes - item.BaseBytes
		item.BaseTokens = item.BaseBytes / 4
		item.BranchTokens = item.BranchBytes / 4
		item.DeltaTokens = item.BranchTokens - item.BaseTokens
		out = append(out, item)
	}
	return out
}

func sizesFlat(sizes []SizeChange) bool {
	for _, size := range sizes {
		if size.DeltaBytes != 0 || size.DeltaTokens != 0 {
			return false
		}
	}
	return true
}

func roleHosts(files []FileChange) []string {
	var out []string
	for _, file := range files {
		for _, rel := range []string{file.Path, file.From} {
			switch {
			case strings.HasPrefix(rel, ".claude/agents/"):
				// The CLI names this host claude-code. hosts.yaml names it claude.
				out = append(out, "claude")
			case strings.HasPrefix(rel, ".cursor/agents/"):
				out = append(out, "cursor")
			case strings.HasPrefix(rel, ".codex/agents/"):
				out = append(out, "codex")
			}
		}
	}
	return out
}

func validateTable(table triggers.Table, hostsFile, scenariosDir string) error {
	loaded, err := hosts.LoadHosts(hostsFile)
	if err != nil {
		return err
	}
	hostNames := map[string]bool{}
	for name := range loaded {
		hostNames[name] = true
	}
	entries, err := os.ReadDir(scenariosDir)
	if err != nil {
		return err
	}
	scenarios := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") {
			scenarios[strings.TrimSuffix(name, ".yaml")] = true
		}
	}
	if len(scenarios) == 0 {
		return fmt.Errorf("no scenarios found")
	}
	return table.Validate(scenarios, hostNames)
}

func writeReport(dir string, report Report) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	if err := os.WriteFile(filepath.Join(dir, "surface.json"), body, 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "surface.md"), []byte(Render(report)), 0o644)
}

func writeOutputs(dir string, baseSide, branchSide captured) error {
	for _, side := range []struct {
		name string
		cap  captured
	}{
		{"base", baseSide},
		{"branch", branchSide},
	} {
		outDir := filepath.Join(dir, "outputs", side.name)
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return err
		}
		for _, out := range side.cap.outputs {
			name := out.name + ".txt"
			if strings.Contains(name, "/") || strings.Contains(name, "..") {
				return fmt.Errorf("output name %q", out.name)
			}
			text := side.cap.norm.Text(out.text)
			if err := os.WriteFile(filepath.Join(outDir, name), []byte(text), 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

func uniqueSorted(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range in {
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	sort.Strings(out)
	if out == nil {
		return []string{}
	}
	return out
}
