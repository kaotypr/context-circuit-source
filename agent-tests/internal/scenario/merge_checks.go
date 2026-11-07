package scenario

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	nameAfterMergeOffered = "after-merge-offered"
	nameFFSynced          = "ff-synced"

	// apiRemote is the bare remote the world gives the api repository.
	apiRemote = "api.git"
	// ffDirtyFile is the uncommitted file the ff scenario leaves in the
	// workspace, which must keep the workspace from moving.
	ffDirtyFile = "notes.txt"
)

func init() {
	Registry[nameAfterMergeOffered] = checkAfterMergeOffered
	Registry[nameFFSynced] = checkFFSynced
}

var (
	// offerSync matches a reply offering to bring the base branch up to date.
	offerSync = regexp.MustCompile(`(?i)fast[- ]?forward|repo sync|/ff\b|\bsync\b|up to date|bring .* main`)
	// offerComplete matches a reply offering to mark the plan completed.
	offerComplete = regexp.MustCompile(`(?i)\bcomplet`)
)

// checkAfterMergeOffered passes when a person saying their request merged
// was offered the sync and the completion and nothing was done: the plan
// is not completed, api's base branch has not moved, and the reply offers
// both.
func checkAfterMergeOffered(ctx CheckContext) CheckResult {
	files := matchFiles(ctx.After.Files, "plans/p1000-*/plan.md")
	if len(files) != 1 {
		return failCheck(nameAfterMergeOffered, "expected one plan.md")
	}
	text, err := readRel(ctx, files[0])
	if err != nil {
		return failCheck(nameAfterMergeOffered, err.Error())
	}
	if front, _, ok := splitFront(text); !ok || frontHas(front, "completed_at") {
		return failCheck(nameAfterMergeOffered, "the plan was completed before the person answered the offer")
	}
	before := parseRefMap(ctx.Before.Refs["api"])["refs/heads/main"]
	after := parseRefMap(ctx.After.Refs["api"])["refs/heads/main"]
	if before == "" || before != after {
		return failCheck(nameAfterMergeOffered, "api's base branch moved before the person answered the offer")
	}
	reply := ctx.Parsed.FinalReply
	if !offerSync.MatchString(reply) || !offerComplete.MatchString(reply) {
		return failCheck(nameAfterMergeOffered, "the reply does not offer both the sync and the completion")
	}
	return passCheck(nameAfterMergeOffered)
}

// checkFFSynced passes when /ff fast-forwarded api's base branch to what
// its remote holds, with HEAD on it, and left the dirty workspace exactly
// where it was with its uncommitted file in place.
func checkFFSynced(ctx CheckContext) CheckResult {
	if ctx.Exec == nil {
		return inconclusiveCheck(nameFFSynced, "command runner is absent")
	}
	remote := RemoteRefObject(ctx.After.RemoteRefs[apiRemote], "refs/heads/main")
	before := parseRefMap(ctx.Before.Refs["api"])["refs/heads/main"]
	after := parseRefMap(ctx.After.Refs["api"])["refs/heads/main"]
	switch {
	case remote == "" || before == "" || after == "":
		return failCheck(nameFFSynced, "api's base branch or its remote is missing from the snapshot")
	case after != remote:
		return failCheck(nameFFSynced, "api's base branch is not at what its remote holds")
	case ctx.After.Heads["api"] != after:
		return failCheck(nameFFSynced, "api's HEAD is not on its base branch")
	}
	if _, err := ctx.Exec(ctx.World.API, "git", "merge-base", "--is-ancestor", before, after); err != nil {
		return failCheck(nameFFSynced, "api's base branch was not fast-forwarded: its earlier tip is not an ancestor")
	}
	wsBefore := parseRefMap(ctx.Before.Refs["ws"])["refs/heads/"+WorkspaceBaseBranch]
	wsAfter := parseRefMap(ctx.After.Refs["ws"])["refs/heads/"+WorkspaceBaseBranch]
	if wsBefore == "" || wsBefore != wsAfter || ctx.Before.Heads["ws"] != ctx.After.Heads["ws"] {
		return failCheck(nameFFSynced, "the workspace moved although it holds an uncommitted file")
	}
	body, err := os.ReadFile(filepath.Join(ctx.World.WS, ffDirtyFile))
	if err != nil || !strings.Contains(string(body), "draft") {
		return failCheck(nameFFSynced, "the workspace's uncommitted "+ffDirtyFile+" was removed or changed")
	}
	return passCheck(nameFFSynced)
}
