package report

import (
	"fmt"
	"strings"
)

type edit struct {
	op      byte
	line    string
	oldLine int
	newLine int
}

func unified(before, after, name string) string {
	if before == after {
		return ""
	}
	a := splitLines(before)
	b := splitLines(after)
	if len(a) > 800 || len(b) > 800 {
		return dumpReplace(a, b, name)
	}
	return formatHunks(lcsEdits(a, b), name)
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	parts := strings.Split(s, "\n")
	if strings.HasSuffix(s, "\n") && len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

func countLines(s string) int {
	return len(splitLines(s))
}

func lcsEdits(a, b []string) []edit {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
				continue
			}
			if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
				continue
			}
			dp[i][j] = dp[i][j+1]
		}
	}
	var edits []edit
	i, j := 0, 0
	oldLine, newLine := 1, 1
	for i < n && j < m {
		if a[i] == b[j] {
			edits = append(edits, edit{' ', a[i], oldLine, newLine})
			i++
			j++
			oldLine++
			newLine++
			continue
		}
		if dp[i+1][j] >= dp[i][j+1] {
			edits = append(edits, edit{'-', a[i], oldLine, 0})
			i++
			oldLine++
			continue
		}
		edits = append(edits, edit{'+', b[j], 0, newLine})
		j++
		newLine++
	}
	for i < n {
		edits = append(edits, edit{'-', a[i], oldLine, 0})
		i++
		oldLine++
	}
	for j < m {
		edits = append(edits, edit{'+', b[j], 0, newLine})
		j++
		newLine++
	}
	return edits
}

func formatHunks(edits []edit, name string) string {
	const ctx = 3
	var b strings.Builder
	fmt.Fprintf(&b, "--- before/%s\n+++ after/%s\n", name, name)
	n := len(edits)
	i := 0
	for i < n {
		for i < n && edits[i].op == ' ' {
			i++
		}
		if i >= n {
			break
		}
		start := i - ctx
		if start < 0 {
			start = 0
		}
		end := i
		for end < n {
			if edits[end].op != ' ' {
				end++
				continue
			}
			j := end
			for j < n && edits[j].op == ' ' {
				j++
			}
			if j < n && j-end <= ctx*2 {
				end = j
				continue
			}
			end += ctx
			if end > n {
				end = n
			}
			break
		}
		writeHunk(&b, edits[start:end])
		i = end
	}
	return b.String()
}

func writeHunk(b *strings.Builder, edits []edit) {
	oldStart, newStart := 0, 0
	oldN, newN := 0, 0
	for _, item := range edits {
		if item.op != '+' && oldStart == 0 {
			oldStart = item.oldLine
		}
		if item.op != '-' && newStart == 0 {
			newStart = item.newLine
		}
		if item.op != '+' {
			oldN++
		}
		if item.op != '-' {
			newN++
		}
	}
	if oldN == 0 {
		oldStart = 0
	}
	if newN == 0 {
		newStart = 0
	}
	fmt.Fprintf(b, "@@ -%d,%d +%d,%d @@\n", oldStart, oldN, newStart, newN)
	for _, item := range edits {
		fmt.Fprintf(b, "%c%s\n", item.op, item.line)
	}
}

func dumpReplace(before, after []string, name string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "--- before/%s\n+++ after/%s\n", name, name)
	fmt.Fprintf(&b, "@@ -1,%d +1,%d @@\n", len(before), len(after))
	for _, line := range before {
		fmt.Fprintf(&b, "-%s\n", line)
	}
	for _, line := range after {
		fmt.Fprintf(&b, "+%s\n", line)
	}
	return b.String()
}
