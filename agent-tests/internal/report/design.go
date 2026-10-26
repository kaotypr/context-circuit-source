package report

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// DesignInput pairs runs and reads their produced files.
type DesignInput struct {
	RunDir      string
	BaselineDir string
	Branch      []record.Run
	Before      []record.Run
	Home        string
	User        string
	RunRoots    []string
}

// RenderDesign writes design.md for the paired runs.
func RenderDesign(in DesignInput) (string, error) {
	pairs := pairRuns(in.Branch, in.Before)
	var b strings.Builder
	b.WriteString("## Before / after\n\n")
	b.WriteString("Each pair is one scenario and one host. The run on a side is the passing run with the median total input tokens. The lowest run number wins a tie. When no run passed, the median of the remaining runs is used.\n")
	if len(pairs) == 0 {
		b.WriteString("\nNo base and branch pair was found for any scenario and host.\n")
		return b.String(), nil
	}
	for _, pair := range pairs {
		section, err := renderPair(in, pair)
		if err != nil {
			return "", err
		}
		b.WriteString("\n")
		b.WriteString(section)
	}
	if !strings.HasSuffix(b.String(), "\n") {
		b.WriteByte('\n')
	}
	return b.String(), nil
}

type paired struct {
	Scenario string
	Host     string
	Branch   record.Run
	Before   record.Run
	BranchN  int
	BeforeN  int
	How      string
}

func pairRuns(branch, before []record.Run) []paired {
	bGroups := groupRuns(branch)
	aGroups := groupRuns(before)
	keys := map[string]bool{}
	for key := range bGroups {
		if _, ok := aGroups[key]; ok {
			keys[key] = true
		}
	}
	var ordered []string
	for key := range keys {
		ordered = append(ordered, key)
	}
	sort.Strings(ordered)
	var out []paired
	for _, key := range ordered {
		scenario, host, _ := strings.Cut(key, "\x00")
		bRun, bHow, bok := pickRun(bGroups[key])
		aRun, _, aok := pickRun(aGroups[key])
		if !bok || !aok {
			continue
		}
		out = append(out, paired{
			Scenario: scenario,
			Host:     host,
			Branch:   bRun,
			Before:   aRun,
			BranchN:  len(bGroups[key]),
			BeforeN:  len(aGroups[key]),
			How:      bHow,
		})
	}
	return out
}

func groupRuns(runs []record.Run) map[string][]record.Run {
	out := map[string][]record.Run{}
	for _, run := range runs {
		if run.Verdict == "harness-error" || run.Scenario == "" || run.Host == "" {
			continue
		}
		key := run.Scenario + "\x00" + run.Host
		out[key] = append(out[key], run)
	}
	return out
}

// pickRun selects the passing run with the median total input tokens.
func pickRun(runs []record.Run) (record.Run, string, bool) {
	if len(runs) == 0 {
		return record.Run{}, "", false
	}
	passing := filterRuns(runs, func(run record.Run) bool { return run.Verdict == "pass" })
	pool := passing
	how := "the passing run with the median total input tokens"
	if len(pool) == 0 {
		pool = filterRuns(runs, func(run record.Run) bool { return run.Verdict != "harness-error" })
		how = "the median run; none passed"
	}
	if len(pool) == 0 {
		pool = append([]record.Run{}, runs...)
		how = "the median run; none passed"
	}
	pool = medianOrder(pool)
	return pool[(len(pool)-1)/2], how, true
}

func medianOrder(runs []record.Run) []record.Run {
	with := filterRuns(runs, func(run record.Run) bool { return run.Tokens.InputTotal != nil })
	use := runs
	if len(with) > 0 {
		use = with
	}
	sort.SliceStable(use, func(i, j int) bool {
		if len(with) > 0 {
			ti := *use[i].Tokens.InputTotal
			tj := *use[j].Tokens.InputTotal
			if ti != tj {
				return ti < tj
			}
		}
		if use[i].Run != use[j].Run {
			return use[i].Run < use[j].Run
		}
		return false
	})
	return use
}

func filterRuns(runs []record.Run, keep func(record.Run) bool) []record.Run {
	var out []record.Run
	for _, run := range runs {
		if keep(run) {
			out = append(out, run)
		}
	}
	return out
}

func renderPair(in DesignInput, pair paired) (string, error) {
	branchDir := locateRun(in.RunDir, pair.Branch, false)
	beforeDir := locateRun(in.baselineRoot(), pair.Before, in.BaselineDir != "")
	branchFiles, err := loadProduced(branchDir)
	if err != nil {
		return "", err
	}
	beforeFiles, err := loadProduced(beforeDir)
	if err != nil {
		return "", err
	}
	branchFiles = cleanFiles(branchFiles, in)
	beforeFiles = cleanFiles(beforeFiles, in)
	var b strings.Builder
	fmt.Fprintf(&b, "### `%s` · %s\n\n", pair.Scenario, pair.Host)
	fmt.Fprintf(&b, "Pair: branch r%d (%s), before r%d (%s). Chosen as %s.\n\n",
		runNumber(pair.Branch), runWord(pair.BranchN), runNumber(pair.Before), runWord(pair.BeforeN), pair.How)
	b.WriteString(summaryTable(beforeFiles, branchFiles))
	b.WriteString("\n")
	b.WriteString(fileDiffs(beforeFiles, branchFiles))
	b.WriteString("### Reads\n\n")
	b.WriteString(readsSection(beforeDir, branchDir, in))
	b.WriteString("\n### Reply excerpts\n\n")
	b.WriteString("First 12 lines. Key-shaped strings are replaced with `[redacted]`.\n\n")
	b.WriteString("#### Before\n\n")
	b.WriteString(replyFence(readReply(beforeDir), in))
	b.WriteString("\n#### After\n\n")
	b.WriteString(replyFence(readReply(branchDir), in))
	b.WriteString("\n### Tokens\n\n")
	b.WriteString(tokenSection(pair.Before, pair.Branch))
	return b.String(), nil
}

func (in DesignInput) baselineRoot() string {
	if in.BaselineDir != "" {
		return in.BaselineDir
	}
	return in.RunDir
}

func locateRun(root string, run record.Run, baseline bool) string {
	if baseline {
		alt := filepath.Join(root, "artifacts", run.Scenario, run.Host)
		if isDir(alt) {
			return alt
		}
	}
	side := run.Side
	if side == "" {
		if baseline {
			side = "base"
		} else {
			side = "branch"
		}
	}
	return filepath.Join(root, "runs", run.Scenario, run.Host, side, fmt.Sprintf("r%d", runNumber(run)))
}

func runWord(n int) string {
	if n == 1 {
		return "1 run"
	}
	return fmt.Sprintf("%d runs", n)
}

func runNumber(run record.Run) int {
	if run.Run < 1 {
		return 1
	}
	return run.Run
}

func isDir(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

func loadProduced(dir string) (map[string][]byte, error) {
	out := map[string][]byte{}
	if dir == "" || !isDir(dir) {
		return out, nil
	}
	roots := producedRoots(dir)
	for _, root := range roots {
		if !isDir(root.abs) {
			continue
		}
		err := filepath.WalkDir(root.abs, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.Type()&os.ModeSymlink != 0 {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if d.IsDir() {
				return nil
			}
			if !d.Type().IsRegular() {
				return nil
			}
			rel, err := filepath.Rel(root.abs, path)
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
			out[root.prefix+"/"+rel] = body
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

type producedRoot struct {
	abs    string
	prefix string
}

func producedRoots(dir string) []producedRoot {
	art := filepath.Join(dir, "artifacts")
	if isDir(art) {
		return []producedRoot{
			{filepath.Join(art, "files"), "files"},
			{filepath.Join(art, "design"), "design"},
		}
	}
	return []producedRoot{
		{filepath.Join(dir, "files"), "files"},
		{filepath.Join(dir, "design"), "design"},
	}
}

func cleanFiles(files map[string][]byte, in DesignInput) map[string][]byte {
	out := make(map[string][]byte, len(files))
	for name, body := range files {
		out[name] = []byte(publicText(string(body), in.Home, in.User, in.RunRoots))
	}
	return out
}

func summaryTable(before, after map[string][]byte) string {
	var b strings.Builder
	b.WriteString("| | Before | After |\n")
	b.WriteString("|---|---|---|\n")
	fmt.Fprintf(&b, "| Files produced | %s | %s |\n", fileList(before), fileList(after))
	fmt.Fprintf(&b, "| Headings | %s | %s |\n", headingList(before), headingList(after))
	fmt.Fprintf(&b, "| Lines | %d | %d |\n", lineTotal(before), lineTotal(after))
	return b.String()
}

func fileList(files map[string][]byte) string {
	if len(files) == 0 {
		return "—"
	}
	names := mapKeys(files)
	parts := make([]string, len(names))
	for i, name := range names {
		parts[i] = "`" + name + "`"
	}
	return strings.Join(parts, ", ")
}

func headingList(files map[string][]byte) string {
	heads := headingsOf(files)
	if len(heads) == 0 {
		return "—"
	}
	return strings.Join(heads, " · ")
}

func headingsOf(files map[string][]byte) []string {
	var out []string
	seen := map[string]bool{}
	for _, name := range mapKeys(files) {
		for _, line := range splitLines(string(files[name])) {
			text, ok := headingText(line)
			if !ok || seen[text] {
				continue
			}
			seen[text] = true
			out = append(out, text)
		}
	}
	return out
}

func headingText(line string) (string, bool) {
	i := 0
	for i < len(line) && line[i] == '#' {
		i++
	}
	if i == 0 || i > 6 || i >= len(line) || line[i] != ' ' {
		return "", false
	}
	text := strings.TrimSpace(line[i:])
	if text == "" {
		return "", false
	}
	return text, true
}

func lineTotal(files map[string][]byte) int {
	n := 0
	for _, body := range files {
		n += countLines(string(body))
	}
	return n
}

func fileDiffs(before, after map[string][]byte) string {
	names := map[string]bool{}
	for name := range before {
		names[name] = true
	}
	for name := range after {
		names[name] = true
	}
	var ordered []string
	for name := range names {
		ordered = append(ordered, name)
	}
	sort.Strings(ordered)
	var b strings.Builder
	for _, name := range ordered {
		diff := unified(string(before[name]), string(after[name]), name)
		if diff == "" {
			continue
		}
		fmt.Fprintf(&b, "<details><summary>%s (before → after)</summary>\n\n```diff\n%s```\n</details>\n\n", name, diff)
	}
	return b.String()
}

func readsSection(beforeDir, branchDir string, in DesignInput) string {
	before, beforeOK := readReads(beforeDir)
	after, afterOK := readReads(branchDir)
	before = publicText(before, in.Home, in.User, in.RunRoots)
	after = publicText(after, in.Home, in.User, in.RunRoots)
	if !beforeOK && !afterOK {
		return "reads list not available\n"
	}
	if !beforeOK {
		return "before: reads list not available\n\n```diff\n" + unified("", after, "reads.txt") + "```\n"
	}
	if !afterOK {
		return "after: reads list not available\n\n```diff\n" + unified(before, "", "reads.txt") + "```\n"
	}
	if before == after {
		return "reads list unchanged\n"
	}
	return "```diff\n" + unified(before, after, "reads.txt") + "```\n"
}

func readReads(dir string) (string, bool) {
	for _, rel := range []string{"reads.txt", filepath.Join("artifacts", "reads.txt")} {
		b, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			continue
		}
		text := string(b)
		if strings.TrimSpace(text) == "" || strings.TrimSpace(text) == "reads list not available" {
			return "", false
		}
		return text, true
	}
	return "", false
}

func readReply(dir string) string {
	for _, rel := range []string{"reply.md", "reply", filepath.Join("artifacts", "reply.md"), filepath.Join("artifacts", "reply")} {
		b, err := os.ReadFile(filepath.Join(dir, rel))
		if err == nil {
			return string(b)
		}
	}
	return ""
}

func replyFence(text string, in DesignInput) string {
	text = publicText(text, in.Home, in.User, in.RunRoots)
	text = firstLines(text, 12)
	if strings.TrimSpace(text) == "" {
		text = "(none)"
	}
	return "```\n" + text + "\n```\n"
}

func firstLines(s string, n int) string {
	lines := splitLines(s)
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}

func tokenSection(before, branch record.Run) string {
	return "Total input tokens: " + tokenArrow(before.Tokens.InputTotal, branch.Tokens.InputTotal) + "\n" +
		"Output tokens: " + tokenArrow(before.Tokens.Output, branch.Tokens.Output) + "\n"
}

func tokenArrow(before, after *int64) string {
	left := "n/a"
	right := "n/a"
	if before != nil {
		left = formatCount(float64(*before))
	}
	if after != nil {
		right = formatCount(float64(*after))
	}
	cell := left + " → " + right
	if before != nil && after != nil {
		if pct, ok := formatPct(float64(*before), float64(*after)); ok {
			cell += " (" + pct + ")"
		}
	}
	return cell
}

func mapKeys(files map[string][]byte) []string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func safeRel(rel string) bool {
	if rel == "" || strings.HasPrefix(rel, "/") || strings.Contains(rel, `\`) {
		return false
	}
	clean := pathClean(rel)
	return clean != ".." && !strings.HasPrefix(clean, "../")
}

func pathClean(rel string) string {
	return strings.TrimPrefix(filepath.ToSlash(filepath.Clean(rel)), "./")
}
