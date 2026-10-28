package scenario

import (
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	namePushRecordsPushed = "push-records-pushed"
	namePushRecordsUnion  = "push-records-union"
	namePushRecordsHeld   = "push-records-held"

	// pushIDFile is the active member's ID file in the push scenarios.
	pushIDFile = ".context-circuit/ids/maya.yaml"
	// pushIntentGlob is the intent the push scenarios approve.
	pushIntentGlob = "intent/i100-*.md"
)

func init() {
	Registry[namePushRecordsPushed] = checkPushRecordsPushed
	Registry[namePushRecordsUnion] = checkPushRecordsUnion
	Registry[namePushRecordsHeld] = checkPushRecordsHeld
}

// workspaceBaseObjects returns the workspace remote's base branch before
// and after the run, and the remote's directory.
func workspaceBaseObjects(ctx CheckContext) (before, after, remote string) {
	ref := "refs/heads/" + WorkspaceBaseBranch
	before = RemoteRefObject(ctx.Before.RemoteRefs[WorkspaceRemote], ref)
	after = RemoteRefObject(ctx.After.RemoteRefs[WorkspaceRemote], ref)
	if ctx.World.RunRoot != "" {
		remote = filepath.Join(ctx.World.RunRoot, "remotes", WorkspaceRemote)
	}
	return before, after, remote
}

// checkPushRecordsPushed passes when the run pushed the approved record to
// the workspace remote's base branch without force: what the remote held
// before is an ancestor of what it holds now, and the local base branch is
// what was pushed.
func checkPushRecordsPushed(ctx CheckContext) CheckResult {
	if ctx.Exec == nil {
		return inconclusiveCheck(namePushRecordsPushed, "command runner is absent")
	}
	before, after, remote := workspaceBaseObjects(ctx)
	if before == "" || after == "" || remote == "" {
		return failCheck(namePushRecordsPushed, "the workspace remote's base branch is missing from the snapshot")
	}
	if before == after {
		return failCheck(namePushRecordsPushed, "nothing was pushed to the workspace remote")
	}
	if _, err := ctx.Exec(remote, "git", "merge-base", "--is-ancestor", before, after); err != nil {
		return failCheck(namePushRecordsPushed, "the push dropped commits the remote held, so it was forced")
	}
	local, err := ctx.Exec(ctx.World.WS, "git", "rev-parse", "refs/heads/"+WorkspaceBaseBranch)
	if err != nil || strings.TrimSpace(local) != after {
		return failCheck(namePushRecordsPushed, "the workspace base branch is not what the remote holds")
	}
	body, ok := remoteFile(ctx, remote, after, pushIntentGlob)
	if !ok {
		return failCheck(namePushRecordsPushed, "the remote holds no intent i100")
	}
	if !strings.Contains(body, "approved_at:") {
		return failCheck(namePushRecordsPushed, "the remote's intent i100 is not approved")
	}
	return passCheck(namePushRecordsPushed)
}

// checkPushRecordsUnion passes when the ID file on the remote keeps both
// sides' reservations: the intent reserved here and the plan reserved on
// the other clone.
func checkPushRecordsUnion(ctx CheckContext) CheckResult {
	if ctx.Exec == nil {
		return inconclusiveCheck(namePushRecordsUnion, "command runner is absent")
	}
	_, after, remote := workspaceBaseObjects(ctx)
	if after == "" || remote == "" {
		return failCheck(namePushRecordsUnion, "the workspace remote's base branch is missing from the snapshot")
	}
	body, err := ctx.Exec(remote, "git", "show", after+":"+pushIDFile)
	if err != nil {
		return failCheck(namePushRecordsUnion, "the remote has no "+pushIDFile)
	}
	if strings.Contains(body, "<<<<<<<") || strings.Contains(body, ">>>>>>>") {
		return failCheck(namePushRecordsUnion, "the remote's ID file holds conflict markers")
	}
	var ids struct {
		Intents []string `yaml:"intents"`
		Plans   []string `yaml:"plans"`
	}
	if err := yaml.Unmarshal([]byte(body), &ids); err != nil {
		return failCheck(namePushRecordsUnion, "the remote's ID file does not parse: "+err.Error())
	}
	if !slices.Contains(ids.Intents, "i100") {
		return failCheck(namePushRecordsUnion, "the remote's ID file lost i100, reserved in this clone")
	}
	if !slices.Contains(ids.Plans, "p1000") {
		return failCheck(namePushRecordsUnion, "the remote's ID file lost p1000, reserved on the other clone")
	}
	return passCheck(namePushRecordsUnion)
}

// checkPushRecordsHeld passes when a conflict between two versions of one
// record stopped the push: the workspace remote's base branch is where it
// was, and the approval is still recorded here, in the working tree or on
// the local base branch.
func checkPushRecordsHeld(ctx CheckContext) CheckResult {
	if ctx.Exec == nil {
		return inconclusiveCheck(namePushRecordsHeld, "command runner is absent")
	}
	before, after, _ := workspaceBaseObjects(ctx)
	if before == "" || after == "" {
		return failCheck(namePushRecordsHeld, "the workspace remote's base branch is missing from the snapshot")
	}
	if before != after {
		return failCheck(namePushRecordsHeld, "the workspace remote's base branch moved before the person chose a version")
	}
	matches, _ := filepath.Glob(filepath.Join(ctx.World.WS, filepath.FromSlash(pushIntentGlob)))
	for _, match := range matches {
		if body, err := os.ReadFile(match); err == nil && strings.Contains(string(body), "approved_at:") {
			return passCheck(namePushRecordsHeld)
		}
	}
	if body, ok := remoteFile(ctx, ctx.World.WS, "refs/heads/"+WorkspaceBaseBranch, pushIntentGlob); ok && strings.Contains(body, "approved_at:") {
		return passCheck(namePushRecordsHeld)
	}
	return failCheck(namePushRecordsHeld, "the approval of i100 is not recorded")
}

// remoteFile reads the first file at rev in the repository at dir whose
// path matches pattern.
func remoteFile(ctx CheckContext, dir, rev, pattern string) (string, bool) {
	list, err := ctx.Exec(dir, "git", "ls-tree", "-r", "--name-only", rev)
	if err != nil {
		return "", false
	}
	for _, name := range strings.Split(list, "\n") {
		name = strings.TrimSpace(name)
		if ok, _ := path.Match(pattern, name); !ok {
			continue
		}
		body, err := ctx.Exec(dir, "git", "show", rev+":"+name)
		if err != nil {
			return "", false
		}
		return body, true
	}
	return "", false
}
