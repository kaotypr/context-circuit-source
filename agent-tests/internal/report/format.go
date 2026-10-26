package report

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func formatCount(n float64) string {
	rounded := int64(math.Round(n))
	sign := ""
	if rounded < 0 {
		sign = "-"
		rounded = -rounded
	}
	switch {
	case rounded >= 1_000_000:
		return sign + trimOne(float64(rounded)/1_000_000) + "M"
	case rounded >= 1000:
		return sign + trimOne(float64(rounded)/1000) + "k"
	default:
		return sign + strconv.FormatInt(rounded, 10)
	}
}

func trimOne(v float64) string {
	return strings.TrimSuffix(strconv.FormatFloat(v, 'f', 1, 64), ".0")
}

func formatUSD(v float64) string {
	cents := math.Round(v * 100)
	return fmt.Sprintf("$%.2f", cents/100)
}

func formatPct(before, after float64) (string, bool) {
	if before == 0 {
		return "", false
	}
	rounded := math.Round((after - before) / before * 100)
	if rounded < 0 {
		return fmt.Sprintf("%.0f%%", rounded), true
	}
	return fmt.Sprintf("+%.0f%%", rounded), true
}

// formatWarnPct keeps one decimal so a warning never rounds a 15.1% rise
// down to the same digits as the +15% threshold.
func formatWarnPct(before, after float64) (string, bool) {
	if before == 0 {
		return "", false
	}
	pct := (after - before) / before * 100
	return fmt.Sprintf("%+.1f%%", pct), true
}

// overTokenWarn reports whether after is more than 15% above before.
// An equal 15% is not a warning. A missing or zero before has no percentage.
func overTokenWarn(before, after float64) bool {
	if before <= 0 || after <= before {
		return false
	}
	return after*100 > before*115
}

func shortSHA(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

func countWord(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
