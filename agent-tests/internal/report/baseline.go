package report

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// ExportOptions selects finished run directories and the baseline
// directory to write. No host is started.
type ExportOptions struct {
	// Runs are run directories, each holding runs.jsonl and runs/.
	Runs []string
	// Out must not exist, or must be an empty directory.
	Out string
	// DenyFile is an optional local file with one term per line. The terms
	// are matched case-insensitively and are never printed.
	DenyFile string
	// RedactDeny replaces every case-insensitive occurrence of a deny-listed
	// term in file contents with redactedTerm instead of refusing the
	// export. A file name that contains a term still refuses it.
	RedactDeny bool
	Home       string
	User       string
	Stdout     io.Writer
}

// redactedTerm replaces a deny-listed term when RedactDeny is set.
const redactedTerm = "[private]"

// absolutePathPatterns are path prefixes that must not appear in a
// committed baseline once the home prefix has been redacted.
var absolutePathPatterns = []string{"/Users/", "/home/", "/private/", "/tmp/", "/var/folders/"}

// hostIsolation describes how each host is isolated during a run. It is
// recorded in hosts.yaml so a reader knows what the baseline measured.
var hostIsolation = map[string]string{
	"claude": "empty CLAUDE_CONFIG_DIR and HOME under the run root; --strict-mcp-config with no MCP config; --dangerously-skip-permissions",
	"codex":  "empty CODEX_HOME and HOME under the run root; --ignore-user-config; --sandbox workspace-write",
	"cursor": "empty HOME and XDG_CONFIG_HOME under the run root; --force --trust",
}

// exportFiles are the per-run files copied beside the produced files.
var exportFiles = []string{"reads.txt", "reply.md", "check.json"}

// exportDirs are artifact directories copied from artifacts/ of a run.
var exportDirs = []string{"files", "design", "git"}

type exportRun struct {
	run   record.Run
	dir   string
	id    string
	order int
}

type pairRef struct {
	Scenario string `json:"scenario"`
	Host     string `json:"host"`
}

type harnessRef struct {
	Scenario string `json:"scenario"`
	Host     string `json:"host"`
	RunID    string `json:"run_id"`
	Run      int    `json:"run"`
	Date     string `json:"date"`
}

type baselineSummary struct {
	SrcCommit     string       `json:"src_commit"`
	Rows          []summaryRow `json:"rows"`
	Missing       []pairRef    `json:"missing"`
	HarnessErrors []harnessRef `json:"harness_errors"`
	// RedactDeny is true when the export replaced deny-listed terms.
	RedactDeny bool `json:"redact_deny,omitempty"`
	// Redactions counts, per exported file, how many deny-listed terms were
	// replaced. It never names a term.
	Redactions map[string]int `json:"redactions,omitempty"`
}

// refusal is an export stopped by the scrub check. Its message names the
// file and the kind of match, never a deny-listed term.
type refusal struct {
	msg string
}

func (r *refusal) Error() string { return r.msg }

// ExportBaseline merges run directories into one baseline directory. Per
// scenario and host it keeps the newest run that is not a harness error;
// a pair with only harness errors is listed as missing. Everything is
// written to a temporary directory beside Out and renamed into place only
// when every file passed the scrub check.
func ExportBaseline(opts ExportOptions) error {
	out := opts.Stdout
	if out == nil {
		out = os.Stdout
	}
	if len(opts.Runs) == 0 {
		return fmt.Errorf("export-baseline: at least one --run is required")
	}
	if strings.TrimSpace(opts.Out) == "" {
		return fmt.Errorf("export-baseline: --out is required")
	}
	replaceEmpty, err := checkExportTarget(opts.Out)
	if err != nil {
		return err
	}
	deny, err := loadDenyTerms(opts.DenyFile)
	if err != nil {
		return err
	}
	chosen, missing, harness, err := chooseRuns(opts.Runs)
	if err != nil {
		return err
	}
	if len(chosen) == 0 {
		return fmt.Errorf("export-baseline: no run that is not a harness error was found")
	}
	srcCommit, err := oneSourceCommit(chosen)
	if err != nil {
		return err
	}
	hostsYAML, err := renderHostsYAML(chosen)
	if err != nil {
		return err
	}

	parent := filepath.Dir(filepath.Clean(opts.Out))
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(parent, ".export-baseline-")
	if err != nil {
		return err
	}
	done := false
	defer func() {
		if !done {
			_ = os.RemoveAll(tmp)
		}
	}()
	w := &exportWriter{root: tmp, home: opts.Home, user: opts.User, deny: deny}
	if opts.RedactDeny {
		w.redact = denyPattern(deny)
		w.counts = map[string]int{}
	}

	var log bytes.Buffer
	for _, item := range chosen {
		b, err := json.Marshal(item.run)
		if err != nil {
			return err
		}
		log.Write(b)
		log.WriteByte('\n')
	}
	if err := w.write("runs.jsonl", log.Bytes()); err != nil {
		return err
	}
	if err := w.write("hosts.yaml", []byte(hostsYAML)); err != nil {
		return err
	}
	for _, item := range chosen {
		if err := w.copyRun(item); err != nil {
			return err
		}
	}
	if err := w.write("README.md", []byte(renderBaselineReadme(srcCommit, chosen, missing))); err != nil {
		return err
	}
	// summary.json is written last so it can count the replacements made in
	// every other file.
	summary := baselineSummary{
		SrcCommit:     srcCommit,
		Rows:          summaryRows(chosen),
		Missing:       missing,
		HarnessErrors: harness,
	}
	if opts.RedactDeny {
		summary.RedactDeny = true
		summary.Redactions = w.counts
	}
	sb, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return err
	}
	if err := w.write("summary.json", append(sb, '\n')); err != nil {
		return err
	}
	if replaceEmpty {
		if err := os.Remove(opts.Out); err != nil {
			return err
		}
	}
	if err := os.Rename(tmp, opts.Out); err != nil {
		return err
	}
	done = true
	fmt.Fprintf(out, "exported %d run(s); %d scenario and host pair(s) missing\n", len(chosen), len(missing))
	if opts.RedactDeny {
		total := 0
		for _, n := range w.counts {
			total += n
		}
		fmt.Fprintf(out, "replaced %d deny-listed occurrence(s) in %d file(s) with %s\n", total, len(w.counts), redactedTerm)
	}
	return nil
}

// checkExportTarget refuses an existing file or a non-empty directory. It
// reports whether an empty directory is there to be replaced.
func checkExportTarget(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("export-baseline: --out exists and is not a directory")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	if len(entries) > 0 {
		return false, fmt.Errorf("export-baseline: --out already exists and is not empty; export never overwrites")
	}
	return true, nil
}

// loadDenyTerms reads one term per line. Blank lines and lines starting
// with # are skipped. The error never quotes the file or its contents.
func loadDenyTerms(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("export-baseline: the deny file could not be read")
	}
	var terms []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		terms = append(terms, strings.ToLower(line))
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("export-baseline: the deny file could not be read")
	}
	return terms, nil
}

// chooseRuns keeps, per scenario and host, the newest branch-side run that
// is not a harness error. Base-side runs measure another ref and are
// skipped.
func chooseRuns(dirs []string) ([]exportRun, []pairRef, []harnessRef, error) {
	groups := map[string][]exportRun{}
	var harness []harnessRef
	for i, dir := range dirs {
		runs, err := LoadRuns(filepath.Join(dir, "runs.jsonl"))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("export-baseline: run directory %d: %w", i+1, err)
		}
		id := filepath.Base(filepath.Clean(dir))
		for _, run := range runs {
			if run.Side == "base" || run.Scenario == "" || run.Host == "" {
				continue
			}
			item := exportRun{run: run, dir: dir, id: id, order: i}
			if run.Verdict == "harness-error" {
				harness = append(harness, harnessRef{Scenario: run.Scenario, Host: run.Host, RunID: id, Run: run.Run, Date: run.Date})
			}
			key := run.Scenario + "\x00" + run.Host
			groups[key] = append(groups[key], item)
		}
	}
	var keys []string
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var chosen []exportRun
	var missing []pairRef
	for _, key := range keys {
		var best *exportRun
		for i := range groups[key] {
			item := groups[key][i]
			if item.run.Verdict == "harness-error" {
				continue
			}
			if best == nil || newerRun(item, *best) {
				best = &groups[key][i]
			}
		}
		if best == nil {
			scenario, host, _ := strings.Cut(key, "\x00")
			missing = append(missing, pairRef{Scenario: scenario, Host: host})
			continue
		}
		chosen = append(chosen, *best)
	}
	sort.SliceStable(harness, func(i, j int) bool {
		a, b := harness[i], harness[j]
		if a.Scenario != b.Scenario {
			return a.Scenario < b.Scenario
		}
		if a.Host != b.Host {
			return a.Host < b.Host
		}
		if a.RunID != b.RunID {
			return a.RunID < b.RunID
		}
		return a.Run < b.Run
	})
	if missing == nil {
		missing = []pairRef{}
	}
	if harness == nil {
		harness = []harnessRef{}
	}
	return chosen, missing, harness, nil
}

// newerRun orders by date, then by run id (a UTC timestamp), then by the
// position of the run directory on the command line, then by run number.
func newerRun(a, b exportRun) bool {
	if a.run.Date != b.run.Date {
		return a.run.Date > b.run.Date
	}
	if a.id != b.id {
		return a.id > b.id
	}
	if a.order != b.order {
		return a.order > b.order
	}
	return a.run.Run > b.run.Run
}

func oneSourceCommit(chosen []exportRun) (string, error) {
	seen := map[string]bool{}
	var commits []string
	for _, item := range chosen {
		if !seen[item.run.SrcCommit] {
			seen[item.run.SrcCommit] = true
			commits = append(commits, item.run.SrcCommit)
		}
	}
	if len(commits) != 1 {
		sort.Strings(commits)
		return "", fmt.Errorf("export-baseline: the chosen runs come from %d source commits (%s); a baseline is one commit", len(commits), strings.Join(commits, ", "))
	}
	return commits[0], nil
}

// renderHostsYAML writes one entry per host. A field that differs between
// the chosen runs of one host refuses the export.
func renderHostsYAML(chosen []exportRun) (string, error) {
	byHost := map[string][]exportRun{}
	for _, item := range chosen {
		byHost[item.run.Host] = append(byHost[item.run.Host], item)
	}
	var names []string
	for name := range byHost {
		names = append(names, name)
	}
	sort.Strings(names)
	var b strings.Builder
	b.WriteString("# Hosts as recorded by the exported runs.\n")
	b.WriteString("# host_version_source is \"command\" when the host printed it for --version,\n")
	b.WriteString("# and \"config\" when the configured value was recorded instead.\n")
	b.WriteString("# model is the pinned id; host_models lists the names the host reported\n")
	b.WriteString("# for it when they differ.\n")
	b.WriteString("hosts:\n")
	for _, name := range names {
		runs := byHost[name]
		fields := []struct {
			key string
			get func(record.Run) string
		}{
			{"host_version", func(r record.Run) string { return r.HostVersion }},
			{"host_version_source", func(r record.Run) string { return r.HostVersionSource }},
			{"model", func(r record.Run) string { return r.Model }},
			{"effort", func(r record.Run) string { return r.Effort }},
			{"src_commit", func(r record.Run) string { return r.SrcCommit }},
			{"cli_version", func(r record.Run) string { return r.CLIVersion }},
			{"fixture_digest", func(r record.Run) string { return r.FixtureDigest }},
		}
		fmt.Fprintf(&b, "  %s:\n", name)
		for _, field := range fields {
			value, err := agreed(name, field.key, runs, field.get)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "    %s: %s\n", field.key, strconv.Quote(value))
		}
		// The host's own model names are informational; runs may differ.
		reported := map[string]bool{}
		for _, item := range runs {
			if item.run.HostModel != "" {
				reported[item.run.HostModel] = true
			}
		}
		if len(reported) > 0 {
			b.WriteString("    host_models:\n")
			for _, value := range sortedKeys(reported) {
				fmt.Fprintf(&b, "      - %s\n", strconv.Quote(value))
			}
		}
		isolation := hostIsolation[name]
		if isolation == "" {
			isolation = "not described"
		}
		fmt.Fprintf(&b, "    isolation: %s\n", strconv.Quote(isolation))
		dates := map[string]bool{}
		versions := map[string]string{}
		for _, item := range runs {
			dates[item.run.Date] = true
			versions[item.run.Scenario] = item.run.ScenarioVersion
		}
		b.WriteString("    dates:\n")
		for _, date := range sortedKeys(dates) {
			fmt.Fprintf(&b, "      - %s\n", strconv.Quote(date))
		}
		b.WriteString("    scenario_versions:\n")
		var scenarios []string
		for scenario := range versions {
			scenarios = append(scenarios, scenario)
		}
		sort.Strings(scenarios)
		for _, scenario := range scenarios {
			fmt.Fprintf(&b, "      %s: %s\n", strconv.Quote(scenario), strconv.Quote(versions[scenario]))
		}
	}
	return b.String(), nil
}

func agreed(host, key string, runs []exportRun, get func(record.Run) string) (string, error) {
	seen := map[string]bool{}
	for _, item := range runs {
		seen[get(item.run)] = true
	}
	if len(seen) != 1 {
		return "", fmt.Errorf("export-baseline: host %s: the chosen runs disagree on %s (%s)", host, key, strings.Join(sortedKeys(seen), ", "))
	}
	return get(runs[0].run), nil
}

func sortedKeys(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for key := range set {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func summaryRows(chosen []exportRun) []summaryRow {
	rows := make([]summaryRow, 0, len(chosen))
	for _, item := range chosen {
		run := item.run
		failed := run.FailedChecks
		if failed == nil {
			failed = []string{}
		}
		rows = append(rows, summaryRow{
			Scenario:        run.Scenario,
			ScenarioVersion: run.ScenarioVersion,
			Host:            run.Host,
			HostVersion:     run.HostVersion,
			Model:           run.Model,
			HostModel:       run.HostModel,
			Effort:          run.Effort,
			Verdict:         run.Verdict,
			InputTotal:      run.Tokens.InputTotal,
			Output:          run.Tokens.Output,
			CostUSD:         run.CostUSD,
			CostBasis:       run.CostBasis,
			Runs:            1,
			FailedChecks:    failed,
		})
	}
	return rows
}

func renderBaselineReadme(srcCommit string, chosen []exportRun, missing []pairRef) string {
	cli := ""
	if len(chosen) > 0 {
		cli = chosen[0].run.CLIVersion
	}
	var b strings.Builder
	b.WriteString("# Agent-test baseline\n\n")
	fmt.Fprintf(&b, "Exported with `agent-tests export-baseline` from runs of source commit `%s`, CLI version `%s`.\n\n", srcCommit, cli)
	b.WriteString("Host versions, model IDs, isolation and dates are in `hosts.yaml`. Each run's produced files, reads list, final reply and `check.json` are under `artifacts/<scenario>/<host>/`. No transcripts are kept.\n\n")
	b.WriteString("Compare a branch against it with `report --design --baseline` and this directory.\n\n")
	b.WriteString("## Summary\n\n")
	b.WriteString("| Scenario | Host | Verdict | Total input | Output | Cost (USD) | Failed checks |\n")
	b.WriteString("|---|---|---|---|---|---|---|\n")
	for _, item := range chosen {
		run := item.run
		failed := "—"
		if len(run.FailedChecks) > 0 {
			failed = strings.Join(run.FailedChecks, ", ")
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s | %s | %s | %s |\n",
			run.Scenario, run.Host, run.Verdict, baselineCount(run.Tokens.InputTotal), baselineCount(run.Tokens.Output), baselineCost(run.CostUSD), failed)
	}
	b.WriteString("\n## Missing\n\n")
	if len(missing) == 0 {
		b.WriteString("Every scenario and host in the exported runs has a run that is not a harness error.\n")
	} else {
		b.WriteString("These pairs had only harness errors and need a rerun:\n\n")
		for _, pair := range missing {
			fmt.Fprintf(&b, "- `%s` · %s\n", pair.Scenario, pair.Host)
		}
	}
	b.WriteString("\n## Host notes\n\n")
	b.WriteString("<!-- Edit by hand: what about a host's setup a reader needs to know. -->\n\n")
	b.WriteString("- Cursor ran with a plugin assigned to its account, whose notices appear in some replies (redacted where they named private terms).\n")
	b.WriteString("\n## Notable behaviour\n\n")
	b.WriteString("<!-- Written by hand after reading the runs. -->\n\n")
	b.WriteString("_Not written yet._\n")
	return b.String()
}

func baselineCount(v *int64) string {
	if v == nil {
		return "n/a"
	}
	return strconv.FormatInt(*v, 10)
}

func baselineCost(v *float64) string {
	if v == nil {
		return "n/a"
	}
	return strconv.FormatFloat(*v, 'f', 4, 64)
}

type exportWriter struct {
	root string
	home string
	user string
	deny []string
	// redact is set when deny-listed terms in contents are replaced rather
	// than refused; counts holds the replacements per relative path.
	redact *regexp.Regexp
	counts map[string]int
}

// denyPattern matches any deny-listed term case-insensitively in one pass,
// longest term first, so a replacement is never matched again.
func denyPattern(terms []string) *regexp.Regexp {
	if len(terms) == 0 {
		return nil
	}
	sorted := append([]string{}, terms...)
	sort.SliceStable(sorted, func(i, j int) bool { return len(sorted[i]) > len(sorted[j]) })
	quoted := make([]string, len(sorted))
	for i, term := range sorted {
		quoted[i] = regexp.QuoteMeta(term)
	}
	return regexp.MustCompile("(?i)(?:" + strings.Join(quoted, "|") + ")")
}

// write redacts body, refuses it when an absolute path or a deny-listed
// term remains, and writes it under the export root.
func (w *exportWriter) write(rel string, body []byte) error {
	rel = filepath.ToSlash(rel)
	if utf8.Valid(body) {
		body = []byte(record.Redact(string(body), w.home, w.user))
	}
	text := string(body)
	for _, pattern := range absolutePathPatterns {
		if strings.Contains(text, pattern) {
			return &refusal{fmt.Sprintf("export-baseline: refused: %s contains an absolute path (%s); nothing was written", rel, pattern)}
		}
	}
	lowerName := strings.ToLower(rel)
	for _, term := range w.deny {
		if strings.Contains(lowerName, term) {
			return &refusal{fmt.Sprintf("export-baseline: refused: a file name under %s contains a deny-listed term; nothing was written", runPrefix(rel))}
		}
	}
	if w.redact != nil {
		if n := len(w.redact.FindAllIndex(body, -1)); n > 0 {
			body = w.redact.ReplaceAllLiteral(body, []byte(redactedTerm))
			w.counts[rel] += n
		}
	}
	lowerText := strings.ToLower(string(body))
	if w.redact != nil {
		// The replacement text itself is not a match.
		lowerText = strings.ReplaceAll(lowerText, redactedTerm, "\x00")
	}
	for _, term := range w.deny {
		if strings.Contains(lowerText, term) {
			return &refusal{fmt.Sprintf("export-baseline: refused: %s contains a deny-listed term; nothing was written", rel)}
		}
	}
	dest := filepath.Join(w.root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dest, body, 0o644)
}

// runPrefix is the artifacts/<scenario>/<host> part of a path, so a
// refusal about a file name does not print the name.
func runPrefix(rel string) string {
	parts := strings.Split(rel, "/")
	if len(parts) > 3 {
		parts = parts[:3]
	}
	return strings.Join(parts, "/")
}

func (w *exportWriter) copyRun(item exportRun) error {
	side := item.run.Side
	if side == "" {
		side = "branch"
	}
	src := filepath.Join(item.dir, "runs", item.run.Scenario, item.run.Host, side, fmt.Sprintf("r%d", runNumber(item.run)))
	dest := filepath.ToSlash(filepath.Join("artifacts", item.run.Scenario, item.run.Host))
	for _, name := range exportFiles {
		body, ok, err := readRegular(filepath.Join(src, name))
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		if err := w.write(dest+"/"+name, body); err != nil {
			return err
		}
	}
	for _, name := range exportDirs {
		if err := w.copyTree(filepath.Join(src, "artifacts", name), dest+"/"+name); err != nil {
			return err
		}
	}
	return nil
}

func (w *exportWriter) copyTree(src, dest string) error {
	if !isDir(src) {
		return nil
	}
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !safeRel(rel) {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return w.write(dest+"/"+rel, body)
	})
}

func readRegular(path string) ([]byte, bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if !info.Mode().IsRegular() {
		return nil, false, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}
	return b, true, nil
}
