package scenario

import (
	"fmt"
	"regexp"
	"strings"
)

// Checks for the override scenarios: a repair loop a person allowed runs
// verifier, fixing worker, and verifier again in the coordinator, stops at its
// round limit, and does not run at all where nobody allowed it.

const (
	nameLoopSequence     = "override-loop-sequence"
	nameLoopClean        = "override-loop-clean"
	nameNotAllowedHeld   = "override-not-allowed-held"
	nameNotAllowedSaid   = "override-not-allowed-said"
	nameLimitHeld        = "override-limit-held"
	nameLimitReported    = "override-limit-reported"
	overrideVerifierRole = "x-csv-check-csv-verifier"
)

var (
	// cleanAfterTwo matches a reply saying the loop ended clean in round two.
	cleanAfterTwo = regexp.MustCompile(`(?i)\bclean\b[^.\n]{0,40}\b(2|two)\s+rounds?\b|\b(2|two)\s+rounds?\b[^.\n]{0,40}\bclean\b`)
	// notAllowedWord matches a reply saying the override is not allowed here.
	notAllowedWord = regexp.MustCompile(`(?i)\b(not|n't|never)\s+(been\s+|yet\s+)?allowed\b|\bno allowance\b|\bnot permitted\b|\bextension override\b`)
	// limitWord matches a reply saying the loop stopped at its limit.
	limitWord = regexp.MustCompile(`(?i)\blimit\b|\b(2|two|both|all)\s+rounds?\b|\blast round\b|\bstopp?(ed|ing)?\b`)
	// remainingFinding is the finding the strict verifier always reports.
	remainingFinding = regexp.MustCompile(`(?i)trailing row|count of invoices|invoice count|row count`)
	roleWorker       = regexp.MustCompile(`--role[ =']+worker\b`)
)

func init() {
	Registry[nameLoopSequence] = checkLoopSequence
	Registry[nameLoopClean] = checkLoopClean
	Registry[nameNotAllowedHeld] = checkNotAllowedHeld
	Registry[nameNotAllowedSaid] = checkNotAllowedSaid
	Registry[nameLimitHeld] = checkLimitHeld
	Registry[nameLimitReported] = checkLimitReported
}

// dispatchStep is one `agent dispatch` call the stream shows, by the role it
// names and whether it names the repair-loop override.
type dispatchStep struct {
	verifier, worker, override bool
}

// dispatchSteps reads every `agent dispatch` call for the fixture plan from
// the stream, in order. A call runs from one occurrence to the end of its
// command: the next line break that is not a shell continuation, the end of
// the JSON string holding it, or the next occurrence. Only a call naming
// p1000 counts, which leaves out the examples in a doc or skill the agent
// read. A call repeated verbatim, as a host echoing it does, counts once.
func dispatchSteps(stream string) []dispatchStep {
	var steps []dispatchStep
	last := ""
	for _, line := range strings.Split(stream, "\n") {
		line = strings.ReplaceAll(line, `\\\n`, " ")
		parts := strings.Split(line, "agent dispatch")
		for _, part := range parts[1:] {
			if end := strings.Index(part, `\n`); end >= 0 {
				part = part[:end]
			}
			if end := unescapedQuote(part); end >= 0 {
				part = part[:end]
			}
			if !strings.Contains(part, "p1000") || part == last {
				continue
			}
			last = part
			steps = append(steps, dispatchStep{
				verifier: strings.Contains(part, overrideVerifierRole),
				worker:   roleWorker.MatchString(part),
				override: strings.Contains(part, "--override") && strings.Contains(part, "repair-loop"),
			})
		}
	}
	return steps
}

// unescapedQuote is the index of the first double quote not escaped by a
// backslash, or -1.
func unescapedQuote(text string) int {
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '\\':
			i++
		case '"':
			return i
		}
	}
	return -1
}

// loopSteps keeps the dispatches the repair loop made, as "V" for a
// verification and "W" for a fix.
func loopSteps(stream string) string {
	var b strings.Builder
	for _, step := range dispatchSteps(stream) {
		switch {
		case step.verifier && step.override:
			b.WriteByte('V')
		case step.worker && step.override:
			b.WriteByte('W')
		}
	}
	return b.String()
}

// checkLoopSequence passes when the coordinator verified, had the finding
// fixed, and verified again, each under the override, in that order.
func checkLoopSequence(ctx CheckContext) CheckResult {
	steps := loopSteps(ctx.Stream)
	if !strings.Contains(steps, "VWV") {
		return failCheck(nameLoopSequence, fmt.Sprintf("want verifier, worker, verifier under --override repair-loop; the stream shows %q", steps))
	}
	return passCheck(nameLoopSequence)
}

func checkLoopClean(ctx CheckContext) CheckResult {
	if !cleanAfterTwo.MatchString(ctx.Parsed.FinalReply) {
		return failCheck(nameLoopClean, "the reply does not say the loop was clean after 2 rounds")
	}
	return passCheck(nameLoopClean)
}

// checkNotAllowedHeld passes when no verifier ran during execution: no
// dispatch named it, under an override or not.
func checkNotAllowedHeld(ctx CheckContext) CheckResult {
	for _, step := range dispatchSteps(ctx.Stream) {
		if step.verifier || step.override {
			return failCheck(nameNotAllowedHeld, "the verifier was dispatched although no override was allowed")
		}
	}
	if strings.Contains(ctx.Stream, `"cc-`+overrideVerifierRole+`"`) {
		return failCheck(nameNotAllowedHeld, "the verifier was spawned although no override was allowed")
	}
	return passCheck(nameNotAllowedHeld)
}

func checkNotAllowedSaid(ctx CheckContext) CheckResult {
	if !notAllowedWord.MatchString(ctx.Parsed.FinalReply) {
		return failCheck(nameNotAllowedSaid, "the reply does not say the override is not allowed in this workspace")
	}
	return passCheck(nameNotAllowedSaid)
}

// checkLimitHeld passes when the loop ran its two verifications, one fix
// between them, and nothing after the last.
func checkLimitHeld(ctx CheckContext) CheckResult {
	steps := loopSteps(ctx.Stream)
	if strings.Count(steps, "V") != 2 {
		return failCheck(nameLimitHeld, fmt.Sprintf("want two verifications; the stream shows %q", steps))
	}
	if !strings.HasSuffix(steps, "V") {
		return failCheck(nameLimitHeld, fmt.Sprintf("a fix followed the last round; the stream shows %q", steps))
	}
	return passCheck(nameLimitHeld)
}

func checkLimitReported(ctx CheckContext) CheckResult {
	reply := ctx.Parsed.FinalReply
	if !remainingFinding.MatchString(reply) {
		return failCheck(nameLimitReported, "the reply does not give the remaining finding")
	}
	if !limitWord.MatchString(reply) {
		return failCheck(nameLimitReported, "the reply does not say the loop stopped at its round limit")
	}
	return passCheck(nameLimitReported)
}
