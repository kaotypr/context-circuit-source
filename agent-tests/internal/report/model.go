package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// Input is the evidence Build turns into a report model.
type Input struct {
	Branch       []record.Run
	Before       []record.Run
	BaselineName string
	Surface      Surface
	HasSurface   bool
}

// Model is the Agent tests section, before it is rendered.
type Model struct {
	ChangeKind   string
	HasSurface   bool
	FileCount    int
	OutputCount  int
	BaseSHA      string
	Scenarios    []string
	BeforeSource string
	Hosts        []HostInfo
	Rows         []Row
	TotalCost    *float64
	CostOmitted  bool
	NoTokenHosts []string
	Harness      int
	Inconclusive int
	Warnings     []string
	Gates        []Gate
}

// HostInfo is one row of the host table.
type HostInfo struct {
	Host      string
	Version   string
	Model     string
	Effort    string
	Isolation string
}

// Row is one scenario and host, possibly averaged.
type Row struct {
	Scenario    string
	Host        string
	Branch      sideStats
	Before      sideStats
	Measurement bool
	DeltaBranch *float64
	DeltaBefore *float64
}

// Gate is one U-check box.
type Gate struct {
	Label string
	Pass  bool
}

type sideStats struct {
	Count         int
	Verdicts      map[string]int
	InputSum      float64
	InputN        int
	InputMissing  bool
	OutputSum     float64
	OutputN       int
	OutputMissing bool
	CostSum       float64
	CostN         int
	CostMissing   bool
	Summary       bool
	Basis         map[string]bool
	Failed        []string
}

type overheadSides struct {
	Empty sideStats
	World sideStats
}

// Build aggregates runs. Several runs of one scenario and host use the mean.
func Build(in Input) Model {
	branchRest, branchOH := extractOverhead(in.Branch)
	beforeRest, beforeOH := extractOverhead(in.Before)
	branchStats := statsFor(branchRest)
	beforeStats := statsFor(beforeRest)

	keys := map[string]bool{}
	for key := range branchStats {
		keys[key] = true
	}
	for key := range beforeStats {
		keys[key] = true
	}
	var ordered []string
	for key := range keys {
		ordered = append(ordered, key)
	}
	sort.Strings(ordered)

	var rows []Row
	for _, key := range ordered {
		scenario, host, _ := strings.Cut(key, "\x00")
		row := Row{Scenario: scenario, Host: host}
		if st := branchStats[key]; st != nil {
			row.Branch = *st
		}
		if st := beforeStats[key]; st != nil {
			row.Before = *st
		}
		rows = append(rows, row)
	}
	rows = append(rows, measurementRows(branchOH, beforeOH)...)

	m := Model{
		ChangeKind:   changeKind(in.Surface, in.HasSurface),
		HasSurface:   in.HasSurface,
		FileCount:    in.Surface.Files,
		OutputCount:  in.Surface.Outputs,
		BaseSHA:      in.Surface.BaseSHA,
		Scenarios:    uniqueSorted(in.Surface.Scenarios),
		BeforeSource: beforeSource(in),
		Hosts:        hostTable(append(append([]record.Run{}, in.Branch...), in.Before...)),
		Rows:         rows,
		NoTokenHosts: hostsMissingTokens(append(append([]record.Run{}, in.Branch...), in.Before...)),
		Gates:        gateBoxes(in.Branch),
	}
	m.TotalCost, m.CostOmitted = sumCost(in.Branch)
	m.Harness, m.Inconclusive = countOutcomes(in.Branch)
	m.Warnings = warnings(rows)
	return m
}

func changeKind(s Surface, has bool) string {
	if !has || len(s.Kinds) == 0 {
		return "unknown"
	}
	return strings.Join(uniqueSorted(s.Kinds), ", ")
}

func beforeSource(in Input) string {
	if in.BaselineName != "" {
		return "committed baseline `" + in.BaselineName + "`"
	}
	if len(in.Before) == 0 {
		return "unknown"
	}
	for _, run := range in.Before {
		if run.SrcCommit != "" {
			return "base `" + shortSHA(run.SrcCommit) + "` run alternately on this machine"
		}
	}
	return "base run alternately on this machine"
}

func extractOverhead(runs []record.Run) ([]record.Run, map[string]overheadSides) {
	empty := map[string][]record.Run{}
	world := map[string][]record.Run{}
	var rest []record.Run
	for _, run := range runs {
		if run.Verdict == "harness-error" {
			rest = append(rest, run)
			continue
		}
		switch run.Scenario {
		case "overhead-empty":
			empty[run.Host] = append(empty[run.Host], run)
		case "overhead-world":
			world[run.Host] = append(world[run.Host], run)
		default:
			rest = append(rest, run)
		}
	}
	pairs := map[string]overheadSides{}
	hosts := map[string]bool{}
	for host := range empty {
		hosts[host] = true
	}
	for host := range world {
		hosts[host] = true
	}
	for host := range hosts {
		if len(empty[host]) > 0 && len(world[host]) > 0 {
			pairs[host] = overheadSides{Empty: addAll(empty[host]), World: addAll(world[host])}
			continue
		}
		rest = append(rest, empty[host]...)
		rest = append(rest, world[host]...)
	}
	return rest, pairs
}

func addAll(runs []record.Run) sideStats {
	var st sideStats
	for _, run := range runs {
		st.add(run)
	}
	return st
}

func measurementRows(branch, before map[string]overheadSides) []Row {
	hosts := map[string]bool{}
	for host := range branch {
		hosts[host] = true
	}
	for host := range before {
		hosts[host] = true
	}
	var names []string
	for host := range hosts {
		names = append(names, host)
	}
	sort.Strings(names)
	var rows []Row
	for _, host := range names {
		b := branch[host]
		a := before[host]
		row := Row{
			Scenario:    "overhead (S0)",
			Host:        host,
			Measurement: true,
			Branch:      b.World,
			Before:      a.World,
			DeltaBranch: deltaMean(b.World, b.Empty),
			DeltaBefore: deltaMean(a.World, a.Empty),
		}
		rows = append(rows, row)
	}
	return rows
}

func deltaMean(world, empty sideStats) *float64 {
	w := world.inputMean()
	e := empty.inputMean()
	if w == nil || e == nil {
		return nil
	}
	d := *w - *e
	return &d
}

func statsFor(runs []record.Run) map[string]*sideStats {
	out := map[string]*sideStats{}
	for _, run := range runs {
		if run.Verdict == "harness-error" {
			continue
		}
		key := run.Scenario + "\x00" + run.Host
		st := out[key]
		if st == nil {
			st = &sideStats{}
			out[key] = st
		}
		st.add(run)
	}
	return out
}

func (s *sideStats) add(run record.Run) {
	if s.Verdicts == nil {
		s.Verdicts = map[string]int{}
	}
	if s.Basis == nil {
		s.Basis = map[string]bool{}
	}
	s.Count++
	if run.Summary {
		s.Summary = true
	}
	if run.Verdict != "" {
		s.Verdicts[run.Verdict]++
	}
	if run.Tokens.InputTotal == nil {
		s.InputMissing = true
	} else {
		s.InputSum += float64(*run.Tokens.InputTotal)
		s.InputN++
	}
	if run.Tokens.Output == nil {
		s.OutputMissing = true
	} else {
		s.OutputSum += float64(*run.Tokens.Output)
		s.OutputN++
	}
	if run.CostUSD == nil {
		s.CostMissing = true
	} else {
		s.CostSum += *run.CostUSD
		s.CostN++
	}
	if run.CostBasis != "" {
		s.Basis[run.CostBasis] = true
	}
	for _, check := range run.FailedChecks {
		if !contains(s.Failed, check) {
			s.Failed = append(s.Failed, check)
		}
	}
	sort.Strings(s.Failed)
}

func (s sideStats) inputMean() *float64 {
	return s.mean(s.InputSum, s.InputN, s.InputMissing)
}

func (s sideStats) outputMean() *float64 {
	return s.mean(s.OutputSum, s.OutputN, s.OutputMissing)
}

func (s sideStats) costMean() *float64 {
	return s.mean(s.CostSum, s.CostN, s.CostMissing)
}

func (s sideStats) mean(sum float64, n int, missing bool) *float64 {
	if s.Count == 0 || missing || n == 0 {
		return nil
	}
	v := sum / float64(n)
	return &v
}

func (s sideStats) basis() string {
	if s.Count == 0 {
		return ""
	}
	if len(s.Basis) == 0 {
		return "n/a"
	}
	var keys []string
	for key := range s.Basis {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

func (s sideStats) worst() (string, bool) {
	bestRank := 99
	best := ""
	for verdict := range s.Verdicts {
		rank := verdictRank(verdict)
		if rank < bestRank {
			bestRank = rank
			best = verdict
		}
	}
	if best == "" {
		return "", false
	}
	return best, true
}

func verdictRank(verdict string) int {
	switch verdict {
	case "pass":
		return 3
	case "inconclusive":
		return 2
	case "fail":
		return 1
	default:
		return 0
	}
}

func hostTable(runs []record.Run) []HostInfo {
	type acc struct {
		versions map[string]bool
		models   map[string]bool
		efforts  map[string]bool
	}
	by := map[string]*acc{}
	var names []string
	for _, run := range runs {
		if run.Host == "" {
			continue
		}
		item := by[run.Host]
		if item == nil {
			item = &acc{map[string]bool{}, map[string]bool{}, map[string]bool{}}
			by[run.Host] = item
			names = append(names, run.Host)
		}
		item.versions[blank(run.HostVersion, "unknown")] = true
		item.models[blank(run.Model, "unknown")] = true
		effort := run.Effort
		if effort == "" {
			effort = "default"
		}
		item.efforts[effort] = true
	}
	sort.Strings(names)
	var out []HostInfo
	for _, name := range names {
		item := by[name]
		out = append(out, HostInfo{
			Host:      name,
			Version:   joinKeys(item.versions),
			Model:     joinKeys(item.models),
			Effort:    joinKeys(item.efforts),
			Isolation: isolationOf(name),
		})
	}
	return out
}

func isolationOf(host string) string {
	switch host {
	case "claude", "claude-code":
		return "empty config dir, no MCP"
	case "codex":
		return "empty CODEX_HOME, no user config"
	case "cursor":
		return "empty home"
	default:
		return "unknown"
	}
}

func hostsMissingTokens(runs []record.Run) []string {
	missing := map[string]bool{}
	for _, run := range runs {
		if run.Host == "" || run.Verdict == "harness-error" {
			continue
		}
		if run.Tokens.InputTotal == nil {
			missing[run.Host] = true
		}
	}
	var out []string
	for host := range missing {
		out = append(out, host)
	}
	sort.Strings(out)
	return out
}

func sumCost(runs []record.Run) (*float64, bool) {
	var sum float64
	var n int
	omitted := false
	for _, run := range runs {
		if run.CostUSD == nil {
			omitted = true
			continue
		}
		sum += *run.CostUSD
		n++
	}
	if n == 0 {
		return nil, omitted
	}
	return &sum, omitted
}

func countOutcomes(runs []record.Run) (harness, inconclusive int) {
	for _, run := range runs {
		switch run.Verdict {
		case "harness-error":
			harness++
		case "inconclusive":
			inconclusive++
		}
	}
	return harness, inconclusive
}

func gateBoxes(branch []record.Run) []Gate {
	labels := []struct {
		check string
		label string
	}{
		{"launcher-only", "U1 launcher only"},
		{"nothing-pushed", "U2 nothing pushed"},
		{"no-stray-commits", "U3 no commit outside allowed branch"},
		{"no-attribution", "U4 no attribution"},
		{"check-clean", "U5 check clean at finish"},
		{"sources-passive", "U6 `sources/` untouched"},
		{"no-outside-writes", "U7 no writes outside workspace"},
	}
	graded := 0
	failed := map[string]bool{}
	for _, run := range branch {
		if run.Verdict == "harness-error" {
			continue
		}
		graded++
		for _, check := range run.FailedChecks {
			failed[check] = true
		}
	}
	if graded == 0 {
		return nil
	}
	out := make([]Gate, len(labels))
	for i, item := range labels {
		out[i] = Gate{Label: item.label, Pass: graded > 0 && !failed[item.check]}
	}
	return out
}

func warnings(rows []Row) []string {
	type item struct {
		key  string
		text string
	}
	var items []item
	for _, row := range rows {
		key := row.Scenario + "\x00" + row.Host
		if row.Measurement {
			if text, ok := tokenWarning(row.Scenario, row.Host, row.DeltaBefore, row.DeltaBranch); ok {
				items = append(items, item{key: key + "\x00t", text: text})
			}
			continue
		}
		if row.Branch.Count > 0 && (row.Branch.Verdicts["fail"] > 0 || len(row.Branch.Failed) > 0) {
			items = append(items, item{key: key + "\x00f", text: failWarning(row)})
		}
		if lower, ok := lowerWarning(row); ok {
			items = append(items, item{key: key + "\x00l", text: lower})
		}
		if text, ok := tokenWarning(row.Scenario, row.Host, row.Before.inputMean(), row.Branch.inputMean()); ok {
			items = append(items, item{key: key + "\x00t", text: text})
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].key < items[j].key })
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = item.text
	}
	return out
}

func failWarning(row Row) string {
	if len(row.Branch.Failed) == 0 {
		return fmt.Sprintf("`%s` on %s fails on the branch", row.Scenario, row.Host)
	}
	return fmt.Sprintf("`%s` on %s fails on the branch (%s)", row.Scenario, row.Host, strings.Join(row.Branch.Failed, ", "))
}

func lowerWarning(row Row) (string, bool) {
	branch, bok := row.Branch.worst()
	before, aok := row.Before.worst()
	if !bok || !aok || verdictRank(branch) >= verdictRank(before) {
		return "", false
	}
	return fmt.Sprintf("`%s` on %s: a lower result than before (%s, before %s)", row.Scenario, row.Host, branch, before), true
}

func tokenWarning(scenario, host string, before, branch *float64) (string, bool) {
	if before == nil || branch == nil || !overTokenWarn(*before, *branch) {
		return "", false
	}
	pct, ok := formatWarnPct(*before, *branch)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("`%s` on %s: total input tokens above +15%% (%s)", scenario, host, pct), true
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
	return out
}

func joinKeys(set map[string]bool) string {
	var keys []string
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

func blank(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func contains(list []string, item string) bool {
	for _, have := range list {
		if have == item {
			return true
		}
	}
	return false
}
