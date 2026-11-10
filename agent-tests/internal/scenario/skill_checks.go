package scenario

import (
	"regexp"
	"strings"
)

// Checks for the extension-skill scenario: an extension ships a QA skill that
// declares the named stage qa-pass/qa, and a second extension hooks after it.
// The agent uses the installed skill, runs the stage's points around its work,
// follows the second extension's instruction, and changes nothing outside the
// plan's worktree.

const (
	nameSkillUsed     = "extension-skill-used"
	nameSkillHooksRun = "extension-skill-hooks-run"
	nameSkillFollowed = "extension-skill-hook-followed"
	nameSkillScoped   = "extension-skill-scoped"

	// installedSkill is the QA skill as agent setup installs it.
	installedSkill = "x-qa-pass-qa"
)

var (
	// stagePoint matches an `extension hooks` call at the named stage, with
	// its flags in either order, and captures before or after.
	stagePoint = regexp.MustCompile(`extension hooks[^\n"]{0,200}?(?:--stage[ =]qa-pass/qa[^\n"]{0,200}?--when[ =](before|after)|--when[ =](before|after)[^\n"]{0,200}?--stage[ =]qa-pass/qa)`)
	// commandsRun matches a reply listing the commands the QA pass ran, as
	// the qa-report extension's instruction asks.
	commandsRun = regexp.MustCompile(`(?i)commands?\s+(that\s+)?(were\s+)?run|commands?\s+ran`)
	goTest      = regexp.MustCompile(`go test`)
)

func init() {
	Registry[nameSkillUsed] = checkSkillUsed
	Registry[nameSkillHooksRun] = checkSkillHooksRun
	Registry[nameSkillFollowed] = checkSkillFollowed
	Registry[nameSkillScoped] = checkSkillScoped
}

// checkSkillUsed passes when the stream names the installed skill: the host
// invoked it, or the agent read its SKILL.md.
func checkSkillUsed(ctx CheckContext) CheckResult {
	if strings.Contains(ctx.Stream, installedSkill) || strings.Contains(strings.Join(ctx.Parsed.Reads, "\n"), installedSkill) {
		return passCheck(nameSkillUsed)
	}
	return failCheck(nameSkillUsed, "the stream never names "+installedSkill)
}

// checkSkillHooksRun passes when the stream shows the before point of the
// named stage and, later, its after point.
func checkSkillHooksRun(ctx CheckContext) CheckResult {
	before, after := -1, -1
	for _, match := range stagePoint.FindAllStringSubmatchIndex(ctx.Stream, -1) {
		when := ""
		for _, group := range []int{2, 4} {
			if match[group] >= 0 {
				when = ctx.Stream[match[group]:match[group+1]]
			}
		}
		switch {
		case when == "before" && before < 0:
			before = match[0]
		case when == "after":
			after = match[0]
		}
	}
	switch {
	case before < 0:
		return failCheck(nameSkillHooksRun, "no `extension hooks --stage qa-pass/qa --when before` call")
	case after < 0:
		return failCheck(nameSkillHooksRun, "no `extension hooks --stage qa-pass/qa --when after` call")
	case after < before:
		return failCheck(nameSkillHooksRun, "the after point ran before the before point")
	}
	return passCheck(nameSkillHooksRun)
}

// checkSkillFollowed passes when the reply lists the commands the QA pass
// ran, which only the qa-report extension's instruction asks for.
func checkSkillFollowed(ctx CheckContext) CheckResult {
	if !commandsRun.MatchString(ctx.Parsed.FinalReply) || !goTest.MatchString(ctx.Parsed.FinalReply) {
		return failCheck(nameSkillFollowed, "the reply does not list the commands run, as the after qa-pass/qa instruction asks")
	}
	return passCheck(nameSkillFollowed)
}

// checkSkillScoped passes when the workspace and the bound api checkout are
// as they were: the QA pass works in the plan's worktree only.
func checkSkillScoped(ctx CheckContext) CheckResult {
	for _, repo := range []string{"ws", "api"} {
		if ctx.Before.Heads[repo] != ctx.After.Heads[repo] {
			return failCheck(nameSkillScoped, repo+" HEAD changed")
		}
		if ctx.Before.Status[repo] != ctx.After.Status[repo] {
			return failCheck(nameSkillScoped, repo+" changed: "+strings.TrimSpace(ctx.After.Status[repo]))
		}
	}
	return passCheck(nameSkillScoped)
}
