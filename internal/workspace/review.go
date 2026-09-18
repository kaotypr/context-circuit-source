package workspace

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// A catalog entry carries the date its note was last confirmed against the code,
// and until now nothing ever read that date back. `check` enforced that it was
// present; no command compared it to anything. So it was write-only: knowledge
// went stale invisibly, and the one loop that brought a note back for judgment —
// completion, which returns the entries scoped to a plan's repositories — only
// ever fires for work done through this workspace. A merge by somebody else, a
// commit from before the workspace existed, a hotfix pushed directly: none of
// those reach it, and the note sits there reading confidently and wrong.
//
// The date is compared against the code the note itself points at, rather than
// against the calendar. A note whose anchors nobody has touched is not stale,
// however old its date is; a note with forty commits under its anchors since it
// was last confirmed is stale with evidence, and the finding can say how many.
// Calendar age is the fallback, and only for a note that anchors to nothing and
// only where a workspace asks for it, because a birthday on its own reports
// every note eventually and teaches people to scroll past findings.

var catalogReviewedDate = regexp.MustCompile(`reviewed\s+(\d{4}-\d{2}-\d{2})\s*$`)

// reviewAnchor is the anchor shape as a catalog consumer needs it: the logical
// repository ID and the path under it. noteAnchor finds them in prose; this
// splits what it found.
func splitAnchor(anchor string) (string, string, bool) {
	id, rest, ok := strings.Cut(anchor, "@")
	if !ok || id == "" || rest == "" {
		return "", "", false
	}
	// A patterned anchor names a shape rather than a path — `web@src/features/
	// <feature>/` — and Git cannot match the placeholder. The directory above it
	// is the most precise thing that is still true, so ask about that instead of
	// asking about nothing.
	if cut := strings.IndexAny(rest, "<>*?["); cut >= 0 {
		rest = rest[:cut]
		if slash := strings.LastIndex(rest, "/"); slash >= 0 {
			rest = rest[:slash+1]
		} else {
			return "", "", false
		}
	}
	if rest == "" {
		return "", "", false
	}
	return id, rest, true
}

// noteAnchors collects every anchor a note carries, grouped by repository. The
// Owner block is where they belong, and the readability pass reports one written
// anywhere else; for this purpose all of them count, because an anchor in the
// wrong place still says which code the note describes.
func noteAnchors(body []byte) map[string][]string {
	anchors := map[string][]string{}
	for _, match := range noteAnchor.FindAllStringSubmatch(string(body), -1) {
		id, target, ok := splitAnchor(match[2])
		if !ok {
			continue
		}
		if !contains(anchors[id], target) {
			anchors[id] = append(anchors[id], target)
		}
	}
	return anchors
}

func contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

// movedSince counts commits touching these paths after the day a note was last
// confirmed. The day itself is excluded: a note confirmed on the same day the
// code changed was confirmed against that change, and reporting it would fire on
// every note the moment it is written.
func movedSince(ctx context.Context, checkout Checkout, reviewed time.Time, paths []string) int {
	revision := "HEAD"
	if checkout.BaseBranch != "" {
		if _, err := Git(ctx, checkout.Path, "rev-parse", "--verify", "--end-of-options", "refs/heads/"+checkout.BaseBranch); err == nil {
			revision = "refs/heads/" + checkout.BaseBranch
		}
	}
	args := []string{"log", "--format=%H", "--since=" + reviewed.AddDate(0, 0, 1).Format(DateLayout), "--end-of-options", revision, "--"}
	out, err := Git(ctx, checkout.Path, append(args, paths...)...)
	if err != nil || strings.TrimSpace(out) == "" {
		return 0
	}
	return len(strings.Split(strings.TrimSpace(out), "\n"))
}

// reviewIssues reports notes whose code has moved since they were last
// confirmed. It never fails the diagnostic: a repository this machine has not
// obtained, a note that cannot be read, or a date that does not parse is left to
// the findings that already cover those, rather than reported twice.
func (s *Store) reviewIssues(ctx context.Context) []Finding {
	data, err := s.Read("context/INDEX.md")
	if err != nil {
		return nil
	}
	cfg, err := s.Config()
	if err != nil {
		return nil
	}
	checkouts, missing := map[string]Checkout{}, map[string]bool{}
	var issues []Finding
	fenced := false
	for index, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			fenced = !fenced
			continue
		}
		if fenced || !catalogEntry.MatchString(line) {
			continue
		}
		stamp := catalogReviewedDate.FindStringSubmatch(strings.TrimRight(line, " \t"))
		if stamp == nil {
			continue
		}
		reviewed, err := time.Parse(DateLayout, stamp[1])
		if err != nil {
			continue
		}
		target := entryTarget(line)
		if target == "" {
			continue
		}
		body, err := s.Read("context/" + target)
		if err != nil {
			continue
		}
		anchors := noteAnchors(body)
		if len(anchors) == 0 {
			if finding := unanchoredReview(cfg, index+1, target, reviewed); finding != nil {
				issues = append(issues, *finding)
			}
			continue
		}
		ids := make([]string, 0, len(anchors))
		for id := range anchors {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		var evidence []string
		for _, id := range ids {
			if missing[id] {
				continue
			}
			checkout, ok := checkouts[id]
			if !ok {
				checkout, err = s.Repository(ctx, id)
				if err != nil {
					missing[id] = true
					continue
				}
				checkouts[id] = checkout
			}
			paths := anchors[id]
			sort.Strings(paths)
			if moved := movedSince(ctx, checkout, reviewed, paths); moved > 0 {
				evidence = append(evidence, fmt.Sprintf("%s@%s has %s since", id, paths[0], plural(moved, "commit")))
			}
		}
		if len(evidence) == 0 {
			continue
		}
		issues = append(issues, found(
			fmt.Sprintf("context/INDEX.md:%d: %s was last confirmed %s, and %s", index+1, target, stamp[1], strings.Join(evidence, "; ")),
			needsAPerson+"read the note against the code it anchors to, then edit the note and its catalog entry together and move the reviewed date — or move the date alone if nothing it says changed"))
	}
	return issues
}

// unanchoredReview falls back to calendar age for a note that anchors to no
// code — an external contract in references/, or a note about the product
// rather than the system. It reports nothing unless a workspace asks for it,
// because age alone is not evidence of anything and a finding that fires on
// every note eventually is one nobody reads.
func unanchoredReview(cfg Config, line int, target string, reviewed time.Time) *Finding {
	if cfg.KnowledgeReviewDays <= 0 {
		return nil
	}
	age := int(time.Since(reviewed).Hours() / 24)
	if age <= cfg.KnowledgeReviewDays {
		return nil
	}
	finding := found(
		fmt.Sprintf("context/INDEX.md:%d: %s anchors to no code and was last confirmed %s, %s ago", line, target, reviewed.Format(DateLayout), plural(age, "day")),
		needsAPerson+"confirm it against whatever it does describe, then move its reviewed date; nothing in this workspace will ever bring it back on its own")
	return &finding
}

func plural(count int, noun string) string {
	if count == 1 {
		return "1 " + noun
	}
	return strconv.Itoa(count) + " " + noun + "s"
}

// entryTarget is the note a catalog entry links to, relative to context/.
func entryTarget(line string) string {
	for _, match := range catalogLink.FindAllStringSubmatch(line, -1) {
		target := match[1]
		if strings.HasPrefix(target, "#") || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
			continue
		}
		target, _, _ = strings.Cut(target, "#")
		target = path.Clean(strings.TrimPrefix(target, "./"))
		if target == "" || target == "." || strings.HasPrefix(target, "../") || strings.HasPrefix(target, "/") {
			continue
		}
		return target
	}
	return ""
}
