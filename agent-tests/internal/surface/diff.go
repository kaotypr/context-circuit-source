package surface

import (
	"strings"
)

// Unified returns a unified diff of two texts. Equal texts return an
// empty string. A very large pair is summarized instead of diffed.
func Unified(name, before, after string) string {
	if before == after {
		return ""
	}
	a := splitLines(before)
	b := splitLines(after)
	if len(a)*len(b) > 1_500_000 {
		return "--- " + name + "\n+++ " + name + "\n(diff omitted: " + itoa(len(a)) + " / " + itoa(len(b)) + " lines)\n"
	}
	script := editScript(a, b)
	return formatDiff(name, script)
}

type edit struct {
	kind byte // ' ' '=' kept, '-' deleted, '+' inserted
	line string
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func editScript(a, b []string) []edit {
	n, m := len(a), len(b)
	dp := make([]int, (n+1)*(m+1))
	at := func(i, j int) *int { return &dp[i*(m+1)+j] }
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				*at(i, j) = *at(i-1, j-1) + 1
				continue
			}
			left := *at(i-1, j)
			up := *at(i, j-1)
			if left >= up {
				*at(i, j) = left
			} else {
				*at(i, j) = up
			}
		}
	}
	var rev []edit
	i, j := n, m
	for i > 0 || j > 0 {
		switch {
		case i > 0 && j > 0 && a[i-1] == b[j-1]:
			rev = append(rev, edit{' ', a[i-1]})
			i--
			j--
		case i > 0 && (j == 0 || *at(i-1, j) >= *at(i, j-1)):
			rev = append(rev, edit{'-', a[i-1]})
			i--
		default:
			rev = append(rev, edit{'+', b[j-1]})
			j--
		}
	}
	for l, r := 0, len(rev)-1; l < r; l, r = l+1, r-1 {
		rev[l], rev[r] = rev[r], rev[l]
	}
	return rev
}

func formatDiff(name string, script []edit) string {
	var bld strings.Builder
	bld.WriteString("--- " + name + "\n+++ " + name + "\n")
	hunks := splitHunks(script, 3)
	oldLine, newLine := 1, 1
	pos := 0
	for _, h := range hunks {
		for pos < h.start {
			if script[pos].kind != '+' {
				oldLine++
			}
			if script[pos].kind != '-' {
				newLine++
			}
			pos++
		}
		oldStart, newStart := oldLine, newLine
		oldCount, newCount := 0, 0
		var body strings.Builder
		for _, op := range h.edits {
			body.WriteByte(op.kind)
			body.WriteString(op.line)
			body.WriteByte('\n')
			if op.kind != '+' {
				oldCount++
				oldLine++
			}
			if op.kind != '-' {
				newCount++
				newLine++
			}
			pos++
		}
		bld.WriteString("@@ -" + itoa(oldStart) + "," + itoa(oldCount) + " +" + itoa(newStart) + "," + itoa(newCount) + " @@\n")
		bld.WriteString(body.String())
	}
	return bld.String()
}

type hunk struct {
	start int
	edits []edit
}

func splitHunks(script []edit, context int) []hunk {
	var hunks []hunk
	var cur []edit
	start := 0
	gapStart := 0
	gap := 0
	for i, op := range script {
		if op.kind == ' ' {
			if gap == 0 {
				gapStart = i
			}
			gap++
			if len(cur) > 0 && gap > context*2 {
				cur = append(cur, script[gapStart:gapStart+context]...)
				hunks = append(hunks, hunk{start, cur})
				cur = nil
				gap = context
				gapStart = i - context + 1
			}
			continue
		}
		if len(cur) == 0 {
			lead := gapStart
			if gap > context {
				lead = i - context
			}
			if gap == 0 {
				lead = i
			}
			start = lead
			cur = append(cur, script[lead:i]...)
		} else if gap > 0 {
			cur = append(cur, script[gapStart:i]...)
		}
		gap = 0
		cur = append(cur, op)
	}
	if len(cur) > 0 {
		if gap > context {
			cur = append(cur, script[gapStart:gapStart+context]...)
		} else if gap > 0 {
			cur = append(cur, script[gapStart:]...)
		}
		hunks = append(hunks, hunk{start, cur})
	}
	return hunks
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [16]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
