package report

import (
	"fmt"
	"strings"
)

// Render writes the Agent tests section of a pull request body.
func Render(m Model) string {
	var b strings.Builder
	b.WriteString("## Agent tests\n\n")
	fmt.Fprintf(&b, "**Change kind:** %s\n", m.ChangeKind)
	fmt.Fprintf(&b, "**Surface diff:** %s\n", surfaceLine(m))
	fmt.Fprintf(&b, "**Scenarios:** %s\n", scenariosLine(m.Scenarios))
	fmt.Fprintf(&b, "**Before:** %s\n\n", m.BeforeSource)
	b.WriteString(renderHosts(m.Hosts))
	b.WriteString("\n")
	b.WriteString(scenarioTable(m.Rows))
	b.WriteString("\n")
	fmt.Fprintf(&b, "%s\n", costLine(m))
	fmt.Fprintf(&b, "**Harness errors (excluded, rerun):** %d · **Inconclusive:** %d\n\n", m.Harness, m.Inconclusive)
	b.WriteString("### Warnings (explain each below; the reviewer decides)\n")
	if len(m.Warnings) == 0 {
		b.WriteString("- [ ] none\n")
	} else {
		b.WriteString("- [ ] " + strings.Join(m.Warnings, " · ") + "\n")
	}
	b.WriteString("\n")
	b.WriteString(gateBlock(m.Gates))
	b.WriteString("\n")
	b.WriteString(notable)
	b.WriteString("\n")
	b.WriteString(evidence)
	b.WriteString("\n")
	b.WriteString(checklist)
	return b.String()
}

func surfaceLine(m Model) string {
	if !m.HasSurface {
		return "unknown"
	}
	sha := shortSHA(m.BaseSHA)
	if sha == "" {
		sha = "unknown"
	}
	return fmt.Sprintf("`%d` %s and `%d` %s differ from base `%s`",
		m.FileCount, countWord(m.FileCount, "workspace file", "workspace files"),
		m.OutputCount, countWord(m.OutputCount, "CLI output", "CLI outputs"),
		sha)
}

func scenariosLine(scenarios []string) string {
	return "required by the trigger map: " + codeList(scenarios) + " · asked for by the reviewer: — · one-off: —"
}

func codeList(items []string) string {
	if len(items) == 0 {
		return "—"
	}
	parts := make([]string, len(items))
	for i, item := range items {
		parts[i] = "`" + item + "`"
	}
	return strings.Join(parts, ", ")
}

func renderHosts(hosts []HostInfo) string {
	var b strings.Builder
	b.WriteString("| Host | Host version | Model | Effort | Isolation |\n")
	b.WriteString("|---|---|---|---|---|\n")
	if len(hosts) == 0 {
		b.WriteString("| — | — | — | — | — |\n")
		return b.String()
	}
	for _, host := range hosts {
		fmt.Fprintf(&b, "| %s | %s | `%s` | %s | %s |\n", host.Host, host.Version, host.Model, host.Effort, host.Isolation)
	}
	return b.String()
}

func scenarioTable(rows []Row) string {
	var b strings.Builder
	b.WriteString("| Scenario | Host | Branch | Before | Total input tokens (branch / before / change) | Output tokens | Cost per run (branch / before) | Cost basis |\n")
	b.WriteString("|---|---|---|---|---|---|---|---|\n")
	if len(rows) == 0 {
		b.WriteString("| — | — | — | — | — | — | — | — |\n")
		return b.String()
	}
	for _, row := range rows {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			row.Scenario, row.Host, branchCell(row), beforeCell(row), inputCell(row), outputCell(row), costCell(row), basisCell(row))
	}
	return b.String()
}

func branchCell(row Row) string {
	if row.Measurement {
		return "—"
	}
	return verdictCell(row.Branch)
}

func beforeCell(row Row) string {
	if row.Measurement {
		return "—"
	}
	return verdictCell(row.Before)
}

func verdictCell(s sideStats) string {
	if s.Count == 0 {
		return "—"
	}
	if s.Summary {
		return verdictLabel(s) + " (summary)"
	}
	if len(s.Verdicts) == 1 {
		var verdict string
		for name := range s.Verdicts {
			verdict = name
		}
		if s.Count > 1 {
			return fmt.Sprintf("%s (%d runs)", verdict, s.Count)
		}
		return verdict
	}
	var names []string
	for name := range s.Verdicts {
		names = append(names, name)
	}
	sortStrings(names)
	parts := make([]string, len(names))
	for i, name := range names {
		parts[i] = fmt.Sprintf("%s %d", name, s.Verdicts[name])
	}
	return strings.Join(parts, ", ") + fmt.Sprintf(" (%d runs)", s.Count)
}

func verdictLabel(s sideStats) string {
	if len(s.Verdicts) == 1 {
		for name := range s.Verdicts {
			return name
		}
	}
	var names []string
	for name := range s.Verdicts {
		names = append(names, name)
	}
	sortStrings(names)
	if len(names) == 0 {
		return "—"
	}
	parts := make([]string, len(names))
	for i, name := range names {
		parts[i] = fmt.Sprintf("%s %d", name, s.Verdicts[name])
	}
	return strings.Join(parts, ", ")
}

func inputCell(row Row) string {
	var branch, before *float64
	var branchN, beforeN int
	if row.Measurement {
		branch, before = row.DeltaBranch, row.DeltaBefore
		branchN, beforeN = row.Branch.Count, row.Before.Count
		return "ws\u2212empty: " + valueCell(branch, before, branchN, beforeN, true, false, false)
	}
	return valueCell(row.Branch.inputMean(), row.Before.inputMean(), row.Branch.Count, row.Before.Count, true, row.Branch.Summary, row.Before.Summary)
}

func outputCell(row Row) string {
	if row.Measurement {
		return "—"
	}
	return valueCell(row.Branch.outputMean(), row.Before.outputMean(), row.Branch.Count, row.Before.Count, false, row.Branch.Summary, row.Before.Summary)
}

func costCell(row Row) string {
	if row.Measurement {
		return "—"
	}
	left := "n/a"
	right := "n/a"
	if v := row.Branch.costMean(); v != nil {
		left = formatUSD(*v)
	}
	if v := row.Before.costMean(); v != nil {
		right = formatUSD(*v)
	}
	if row.Branch.Count == 0 {
		left = "—"
	}
	if row.Before.Count == 0 {
		right = "—"
	}
	return left + " / " + right
}

func basisCell(row Row) string {
	left := row.Branch.basis()
	right := row.Before.basis()
	if row.Measurement {
		if left == "" && right == "" {
			return "—"
		}
	}
	if left == "" {
		left = "—"
	}
	if right == "" {
		right = "—"
	}
	if left == right {
		return left
	}
	return left + " / " + right
}

func valueCell(branch, before *float64, branchN, beforeN int, withPct, branchSummary, beforeSummary bool) string {
	left := sideCount(branch, branchN)
	right := sideCount(before, beforeN)
	cell := left + " / " + right
	if !withPct {
		return cell
	}
	if branch != nil && before != nil && branchN > 0 && beforeN > 0 {
		if pct, ok := formatPct(*before, *branch); ok {
			cell += " / " + pct
		}
	}
	if note := meanNote(branchN, beforeN, branchSummary, beforeSummary); note != "" {
		cell += " " + note
	}
	return cell
}

func sideCount(v *float64, n int) string {
	if n == 0 {
		return "—"
	}
	if v == nil {
		return "n/a"
	}
	return formatCount(*v)
}

func meanNote(branchN, beforeN int, branchSummary, beforeSummary bool) string {
	if branchSummary || beforeSummary {
		switch {
		case branchN > 1 && !branchSummary:
			return fmt.Sprintf("(mean of %d branch runs; summary)", branchN)
		case beforeN > 1 && !beforeSummary:
			return fmt.Sprintf("(summary; mean of %d before runs)", beforeN)
		default:
			return "(summary)"
		}
	}
	switch {
	case branchN > 1 && beforeN > 1 && branchN == beforeN:
		return fmt.Sprintf("(mean of %d runs)", branchN)
	case branchN > 1 && beforeN > 1:
		return fmt.Sprintf("(mean of %d branch runs and %d before runs)", branchN, beforeN)
	case branchN > 1:
		return fmt.Sprintf("(mean of %d branch runs)", branchN)
	case beforeN > 1:
		return fmt.Sprintf("(mean of %d before runs)", beforeN)
	default:
		return ""
	}
}

func costLine(m Model) string {
	cost := "n/a"
	extra := ""
	if m.TotalCost != nil {
		cost = formatUSD(*m.TotalCost)
		if m.CostOmitted {
			extra = "; runs without a cost are omitted"
		}
	}
	hosts := "none"
	if len(m.NoTokenHosts) > 0 {
		hosts = strings.Join(m.NoTokenHosts, ", ")
	}
	return fmt.Sprintf("**Total cost:** %s (list-price equivalent%s) · **Hosts with no token report:** %s", cost, extra, hosts)
}

func gateBlock(gates []Gate) string {
	var b strings.Builder
	b.WriteString("### Gate checks\n")
	if len(gates) != 7 {
		b.WriteString("- [ ] U1–U7 were not graded\n")
		return b.String()
	}
	fmt.Fprintf(&b, "- %s %s · %s %s · %s %s\n", box(gates[0].Pass), gates[0].Label, box(gates[1].Pass), gates[1].Label, box(gates[2].Pass), gates[2].Label)
	fmt.Fprintf(&b, "- %s %s · %s %s · %s %s · %s %s\n", box(gates[3].Pass), gates[3].Label, box(gates[4].Pass), gates[4].Label, box(gates[5].Pass), gates[5].Label, box(gates[6].Pass), gates[6].Label)
	return b.String()
}

func box(pass bool) string {
	if pass {
		return "[x]"
	}
	return "[ ]"
}

func sortStrings(items []string) {
	for i := 1; i < len(items); i++ {
		j := i
		for j > 0 && items[j] < items[j-1] {
			items[j], items[j-1] = items[j-1], items[j]
			j--
		}
	}
}

const notable = `### Notable behaviour
<!-- What the agent did differently that the design case does not already cover, and any failure.
     Cite the run id; quote at most a few transcript lines. -->
`

const evidence = `### Evidence
- Report and before/after: pasted above
- Cleaned zip of all runs: <!-- attach agent-tests/out/<run-id>/bundle.zip -->
  It holds every run's JSON line, stream and produced files. The runner cleans it: home and run-root
  paths become ` + "`$HOME`/`$RUN`" + `, the username is removed, and auth files are dropped. If it finds a
  key-shaped string, it refuses to build the zip.
`

const checklist = `### Reviewer checklist
Design case
- [ ] Present and before the evidence; light form only where no behaviour change is intended
- [ ] Larger change: links the issue where the reviewer agreed the design case
- [ ] The problem is shown with evidence, not asserted
- [ ] Before/after use the artifact the trigger map names, from the same scenario, fixture, host and model
- [ ] Every "better" claim is visible in the after artifact or a check, with its run count
- [ ] Costs give a size, or "none", for tokens, results on each host, existing workspaces and customized templates
- [ ] At least one alternative, including leaving it as it is

Agent tests
- [ ] Surface diff and scenarios match the trigger map; host-specific files ran on that host
- [ ] "Before" is a matching baseline, or base runs on the same machine
- [ ] Every warning is explained, and I accept the explanation (or asked for changes)
- [ ] Opened the zip for any failed or flagged run, and for the before/after pair
- [ ] More runs, scenarios or hosts asked for? (write which, or "none"; ask for more runs when a result is borderline)
- [ ] Hosts the contributor lacked: run by the reviewer (write which, or "none")
`
