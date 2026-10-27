package scenario

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

const (
	nameOrientReply       = "orient-reply"
	nameOrientClean       = "orient-clean-trees"
	nameJoinMembers       = "join-members"
	nameJoinLocal         = "join-local-files"
	nameIntentRecord      = "intent-record"
	nameIntentNoPlan      = "intent-no-plan"
	nameIntentLangBody    = "intent-lang-body"
	namePlanRecord        = "plan-record"
	namePlanLinked        = "plan-linked"
	nameExecuteWorktree   = "execute-worktree"
	nameExecuteTests      = "execute-tests-pass"
	nameExecuteRecords    = "execute-records-untouched"
	nameCompletePlan      = "complete-plan-done"
	nameCompleteNote      = "complete-note"
	nameDirectReadme      = "direct-readme"
	nameDirectNoRecords   = "direct-no-records"
	nameBorrowedUntouched = "borrowed-untouched"
	nameBorrowedReply     = "borrowed-reply"

	codeCanary = "CODE-CANARY-41b2"

	// notEnglishRatio is the English function-word share at which prose
	// still counts as English. Ordinary English sits well above it.
	notEnglishRatio = 0.12
)

// draftVoid matches a reply that connects a draft invoice with voiding.
// The two words may appear in either order, up to 80 characters apart.
var draftVoid = regexp.MustCompile(`(?is)draft.{0,80}void|void.{0,80}draft`)

// issuedCredit matches a reply that connects an issued invoice with a credit.
// The two words may appear in either order, up to 80 characters apart.
var issuedCredit = regexp.MustCompile(`(?is)issued.{0,80}credit|credit.{0,80}issued`)

// approvalAsk matches a reply that asks for approval. It looks for the stem
// "approv", so "approve" and "approval" both count.
var approvalAsk = regexp.MustCompile(`(?i)approv`)

// planPresented matches a reply that talks about the plan.
var planPresented = regexp.MustCompile(`(?i)\bplan\b`)

// borrowedUpstream, borrowedHandbook, and borrowedReview match the S8 reply.
// The three words are independent and may appear in any order.
var (
	borrowedUpstream = regexp.MustCompile(`(?i)\bupstream\b`)
	borrowedHandbook = regexp.MustCompile(`(?i)\bhandbook\b`)
	borrowedReview   = regexp.MustCompile(`(?i)\breview\b`)
)

// conventionalSubject is type(scope): subject, the form the api fixture asks for.
var conventionalSubject = regexp.MustCompile(`^[a-z]+\([^)]+\): \S`)

// numberedList matches a numbered open-question item.
var numberedList = regexp.MustCompile(`(?m)^\d+\.\s+\S`)

// mdLink matches a markdown link target.
var mdLink = regexp.MustCompile(`\]\(([^)\s]+)\)`)

// csvWord matches a CSV mention.
var csvWord = regexp.MustCompile(`(?i)\bcsv\b`)

// csvDenied matches a sentence that says CSV is unsupported.
var csvDenied = regexp.MustCompile(`(?i)\bunsupported\b|\bnot supported\b`)

// csvAffirmed matches a sentence that records a CSV export.
var csvAffirmed = regexp.MustCompile(`(?i)\bcsv export\b|\bas csv\b|\bto csv\b`)

// recordIDMention matches the fixture record ids inside a knowledge note.
var recordIDMention = regexp.MustCompile(`\b(?:i100|p1000)\b`)

// reviewedDate matches a catalog entry's reviewed date at the end of a line.
var reviewedDate = regexp.MustCompile(`reviewed\s+(\d{4}-\d{2}-\d{2})\s*$`)

var intentHeadings = []string{
	"Goal",
	"Non-goals",
	"Constraints",
	"Success criteria",
	"Repository scope",
	"Open questions",
}

// planHeadings are required in every plan. planDetailsHeading is required
// only when the plan has supporting files: the record template says to
// remove the section when there are none.
var planHeadings = []string{
	"Outcome and repository scope",
	"Approach",
	"Tasks and order",
	"Risks and checks",
}

const planDetailsHeading = "Details"

// placeholderSentences are the sentences record create writes and a filled
// record must replace.
var placeholderSentences = []string{
	"Describe the intended outcome.",
	"Describe what is out of scope.",
	"Record relevant constraints.",
	"Describe an observable result.",
	"List the repositories likely to be involved.",
	"Describe the requested result and each repository's responsibility.",
	"Ground the approach in the selected repositories.",
	"Describe implementation steps and dependencies.",
	"Record useful risks and ordinary test, lint, or build commands.",
	"Link each supporting file here and assign it in required_files. Remove this section when there are no supporting files.",
}

// englishStopwords are function words. A high share of them means the
// prose is English. The comparison is whole words, lowercased.
var englishStopwords = map[string]bool{
	"a": true, "an": true, "the": true, "and": true, "or": true, "of": true,
	"to": true, "in": true, "for": true, "is": true, "are": true, "was": true,
	"were": true, "be": true, "that": true, "this": true, "with": true,
	"from": true, "on": true, "by": true, "it": true, "as": true, "at": true,
	"not": true, "but": true, "its": true, "their": true, "they": true,
	"can": true,
}

var wordPattern = regexp.MustCompile(`[A-Za-z]+`)

func init() {
	Registry[nameOrientReply] = checkOrientReply
	Registry[nameOrientClean] = checkOrientClean
	Registry[nameJoinMembers] = checkJoinMembers
	Registry[nameJoinLocal] = checkJoinLocal
	Registry[nameIntentRecord] = checkIntentRecord
	Registry[nameIntentNoPlan] = checkIntentNoPlan
	Registry[nameIntentLangBody] = checkIntentLangBody
	Registry[namePlanRecord] = checkPlanRecord
	Registry[namePlanLinked] = checkPlanLinked
	Registry[nameExecuteWorktree] = checkExecuteWorktree
	Registry[nameExecuteTests] = checkExecuteTests
	Registry[nameExecuteRecords] = checkExecuteRecords
	Registry[nameCompletePlan] = checkCompletePlan
	Registry[nameCompleteNote] = checkCompleteNote
	Registry[nameDirectReadme] = checkDirectReadme
	Registry[nameDirectNoRecords] = checkDirectNoRecords
	Registry[nameBorrowedUntouched] = checkBorrowedUntouched
	Registry[nameBorrowedReply] = checkBorrowedReply

	SoftRegistry["code-canary-seen"] = softCodeCanary
	SoftRegistry["read-api-code"] = softReadAPI
	SoftRegistry["subagents"] = softSubagents
	SoftRegistry["commit-subjects-conventional"] = softCommitSubjects
}

func softCodeCanary(ctx CheckContext) any {
	return strings.Contains(ctx.Stream, codeCanary)
}

func softReadAPI(ctx CheckContext) any {
	blob := ctx.Stream + "\n" + strings.Join(ctx.Parsed.Reads, "\n")
	if strings.Contains(blob, "internal/invoice") || strings.Contains(blob, "repositories/api") {
		return true
	}
	if !ctx.Parsed.ReadsAvailable {
		return "inconclusive"
	}
	return false
}

func softSubagents(ctx CheckContext) any {
	return ctx.Parsed.Subagents
}

func softCommitSubjects(ctx CheckContext) any {
	if ctx.Exec == nil {
		return "inconclusive"
	}
	repos := unionKeys(
		unionMap(ctx.Before.Heads, ctx.After.Heads),
		unionMap(ctx.Before.Refs, ctx.After.Refs),
	)
	for _, repo := range repos {
		dir := repoDir(ctx.World, repo)
		if dir == "" {
			return false
		}
		argv := []string{"git", "log", "--all", "--format=%s%x00"}
		if objects := refObjects(ctx.Before.Refs[repo]); len(objects) > 0 {
			argv = append(argv, "--not")
			argv = append(argv, objects...)
		}
		out, err := ctx.Exec(dir, argv...)
		if err != nil {
			return false
		}
		for _, subject := range commitMessages(out) {
			subject = strings.TrimSpace(strings.Split(subject, "\n")[0])
			if subject == "" {
				continue
			}
			if !conventionalSubject.MatchString(subject) {
				return false
			}
		}
	}
	return true
}

func checkOrientReply(ctx CheckContext) CheckResult {
	reply := ctx.Parsed.FinalReply
	if !strings.Contains(reply, "api") || !strings.Contains(reply, ".context-circuit/cli") {
		return failCheck(nameOrientReply, "reply is missing api or .context-circuit/cli")
	}
	if !draftVoid.MatchString(reply) || !issuedCredit.MatchString(reply) {
		return failCheck(nameOrientReply, "reply does not connect draft with void and issued with credit")
	}
	if strings.Contains(ctx.Stream, codeCanary) || strings.Contains(reply, codeCanary) {
		return failCheck(nameOrientReply, "api source canary appeared")
	}
	blob := ctx.Stream + "\n" + strings.Join(ctx.Parsed.Reads, "\n")
	if strings.Contains(reply, "domains/invoices.md") || strings.Contains(blob, "context/INDEX.md") || strings.Contains(blob, "context find") {
		return passCheck(nameOrientReply)
	}
	return failCheck(nameOrientReply, "no sign the catalog or the invoices note was used")
}

func checkOrientClean(ctx CheckContext) CheckResult {
	for _, repo := range []string{"ws", "api"} {
		if ctx.After.Heads[repo] == "" {
			return failCheck(nameOrientClean, repo+" HEAD is missing")
		}
		status, ok := ctx.After.Status[repo]
		if !ok {
			return failCheck(nameOrientClean, repo+" status is missing")
		}
		if len(statusLines(status)) != 0 {
			return failCheck(nameOrientClean, "workspace or api has uncommitted changes")
		}
	}
	return passCheck(nameOrientClean)
}

func checkJoinMembers(ctx CheckContext) CheckResult {
	body, err := readRel(ctx, "members.yaml")
	if err != nil {
		return failCheck(nameJoinMembers, "members.yaml: "+err.Error())
	}
	var doc struct {
		Members map[string]struct {
			Band  *int     `yaml:"band"`
			Roles []string `yaml:"roles"`
		} `yaml:"members"`
	}
	if err := yaml.Unmarshal([]byte(body), &doc); err != nil {
		return failCheck(nameJoinMembers, "members.yaml is not YAML")
	}
	dana, ok := doc.Members["dana"]
	if !ok || dana.Band == nil || *dana.Band != 2 {
		return failCheck(nameJoinMembers, "dana is missing or is not on band 2")
	}
	found := false
	for _, role := range dana.Roles {
		if role == "backend developer" {
			found = true
		}
	}
	if !found {
		return failCheck(nameJoinMembers, "dana has no backend developer role")
	}
	return passCheck(nameJoinMembers)
}

func checkJoinLocal(ctx CheckContext) CheckResult {
	body, err := readRel(ctx, "member.local.yaml")
	if err != nil {
		return failCheck(nameJoinLocal, "member.local.yaml: "+err.Error())
	}
	var id struct {
		Member string `yaml:"member"`
	}
	if yaml.Unmarshal([]byte(body), &id) != nil || id.Member != "dana" {
		return failCheck(nameJoinLocal, "member.local.yaml does not name dana")
	}
	raw, err := readRel(ctx, "repositories.local.yaml")
	if err != nil {
		return failCheck(nameJoinLocal, "repositories.local.yaml: "+err.Error())
	}
	var bind struct {
		Workspace *struct {
			Path string `yaml:"path"`
		} `yaml:"workspace"`
		Bindings map[string]struct {
			Path string `yaml:"path"`
		} `yaml:"bindings"`
	}
	if yaml.Unmarshal([]byte(raw), &bind) != nil || bind.Workspace == nil || strings.TrimSpace(bind.Workspace.Path) == "" {
		return failCheck(nameJoinLocal, "workspace binding is missing")
	}
	if ctx.World.WS == "" || !pathsEqual(resolveBinding(ctx.World.WS, bind.Workspace.Path), ctx.World.WS) {
		return failCheck(nameJoinLocal, "workspace binding is not the active workspace")
	}
	api, ok := bind.Bindings["api"]
	if !ok || strings.TrimSpace(api.Path) == "" {
		return failCheck(nameJoinLocal, "api binding is missing")
	}
	apiWant := filepath.Join(ctx.World.WS, "repositories", "api")
	if !pathsEqual(resolveBinding(ctx.World.WS, api.Path), apiWant) {
		return failCheck(nameJoinLocal, "api binding is not repositories/api")
	}
	if missing := missingRoleSet(ctx.World.WS, ctx.Host); missing != "" {
		return failCheck(nameJoinLocal, missing)
	}
	return passCheck(nameJoinLocal)
}

func checkIntentRecord(ctx CheckContext) CheckResult {
	added := addedMatch(ctx.Before.Files, ctx.After.Files, "intent/i100-*.md")
	if len(added) != 1 {
		return failCheck(nameIntentRecord, "expected one new intent/i100-*.md")
	}
	text, err := readRel(ctx, added[0])
	if err != nil {
		return failCheck(nameIntentRecord, err.Error())
	}
	front, body, ok := splitFront(text)
	if !ok {
		return failCheck(nameIntentRecord, "intent has no frontmatter")
	}
	author, _ := frontScalar(front, "created_by")
	created, hasCreated := frontScalar(front, "created_at")
	if author != "maya" || !hasCreated || strings.TrimSpace(created) == "" {
		return failCheck(nameIntentRecord, "frontmatter needs created_by maya and created_at")
	}
	if frontHas(front, "approved_at") {
		return failCheck(nameIntentRecord, "intent is already approved")
	}
	if !headingsPresent(body, intentHeadings) || hasPlaceholder(body) || !openQuestionsOK(body) {
		return failCheck(nameIntentRecord, "intent sections are incomplete or still the template")
	}
	if !approvalAsk.MatchString(ctx.Parsed.FinalReply) {
		return failCheck(nameIntentRecord, "reply does not ask for approval")
	}
	return passCheck(nameIntentRecord)
}

func checkIntentNoPlan(ctx CheckContext) CheckResult {
	if problem := requireRecordedCheckout(nameIntentNoPlan, ctx, "ws", "api"); problem != "" {
		return failCheck(nameIntentNoPlan, problem)
	}
	for _, file := range ctx.After.Files {
		if strings.HasPrefix(file, "plans/p") {
			return failCheck(nameIntentNoPlan, "a plan file exists")
		}
	}
	if ctx.Before.Worktrees["api"] != ctx.After.Worktrees["api"] {
		return failCheck(nameIntentNoPlan, "api worktrees changed")
	}
	if ctx.Before.Heads["api"] != ctx.After.Heads["api"] || len(statusLines(ctx.After.Status["api"])) != 0 {
		return failCheck(nameIntentNoPlan, "api checkout changed")
	}
	return passCheck(nameIntentNoPlan)
}

func checkIntentLangBody(ctx CheckContext) CheckResult {
	files := matchFiles(ctx.After.Files, "intent/i100-*.md")
	if len(files) != 1 {
		return failCheck(nameIntentLangBody, "expected one intent/i100-*.md")
	}
	text, err := readRel(ctx, files[0])
	if err != nil {
		return failCheck(nameIntentLangBody, err.Error())
	}
	_, body, ok := splitFront(text)
	if !ok {
		return failCheck(nameIntentLangBody, "intent has no frontmatter")
	}
	allowed := map[string]bool{}
	for _, title := range intentHeadings {
		allowed[title] = true
	}
	seen := map[string]bool{}
	for _, line := range strings.Split(body, "\n") {
		title, isHeading := headingTitle(line)
		if !isHeading {
			continue
		}
		if !allowed[title] {
			return failCheck(nameIntentLangBody, "heading is not English: "+title)
		}
		seen[title] = true
	}
	for _, title := range intentHeadings {
		if !seen[title] {
			return failCheck(nameIntentLangBody, "missing English heading "+title)
		}
	}
	words := proseWords(body)
	if len(words) < 20 {
		return failCheck(nameIntentLangBody, "body is too short to judge")
	}
	if englishRatio(words) >= notEnglishRatio {
		return failCheck(nameIntentLangBody, "body still reads as English")
	}
	return passCheck(nameIntentLangBody)
}

func checkPlanRecord(ctx CheckContext) CheckResult {
	added := addedMatch(ctx.Before.Files, ctx.After.Files, "plans/p1000-*/plan.md")
	if len(added) != 1 {
		return failCheck(namePlanRecord, "expected one new plans/p1000-*/plan.md")
	}
	text, err := readRel(ctx, added[0])
	if err != nil {
		return failCheck(namePlanRecord, err.Error())
	}
	front, body, ok := splitFront(text)
	if !ok {
		return failCheck(namePlanRecord, "plan has no frontmatter")
	}
	intent, _ := frontScalar(front, "intent")
	if intent != "i100" || !containsValue(frontList(front, "repositories"), "api") {
		return failCheck(namePlanRecord, "plan frontmatter needs intent i100 and repository api")
	}
	if !headingsPresent(body, planHeadings) || hasPlaceholder(body) {
		return failCheck(namePlanRecord, "plan sections are incomplete or still the template")
	}
	if !hasHeading(body, planDetailsHeading) && planHasSupportingFiles(front, body) {
		return failCheck(namePlanRecord, "plan has supporting files but no Details section")
	}
	if ok, missing := linksAssigned(front, body); !ok {
		return failCheck(namePlanRecord, "linked detail is not in required_files: "+missing)
	}
	if !planPresented.MatchString(ctx.Parsed.FinalReply) {
		return failCheck(namePlanRecord, "reply does not present the plan")
	}
	return passCheck(namePlanRecord)
}

func checkPlanLinked(ctx CheckContext) CheckResult {
	if problem := requireRecordedCheckout(namePlanLinked, ctx, "ws", "api"); problem != "" {
		return failCheck(namePlanLinked, problem)
	}
	files := matchFiles(ctx.After.Files, "intent/i100-*.md")
	if len(files) != 1 {
		return failCheck(namePlanLinked, "expected one intent/i100-*.md")
	}
	text, err := readRel(ctx, files[0])
	if err != nil {
		return failCheck(namePlanLinked, err.Error())
	}
	front, _, ok := splitFront(text)
	if !ok || !containsValue(frontList(front, "plans"), "p1000") {
		return failCheck(namePlanLinked, "i100 does not list p1000")
	}
	if ctx.Before.Worktrees["api"] != ctx.After.Worktrees["api"] {
		return failCheck(namePlanLinked, "a worktree appeared")
	}
	if ctx.Before.Heads["api"] != ctx.After.Heads["api"] || len(statusLines(ctx.After.Status["api"])) != 0 {
		return failCheck(namePlanLinked, "api checkout changed")
	}
	return passCheck(namePlanLinked)
}

func checkExecuteWorktree(ctx CheckContext) CheckResult {
	if ctx.Exec == nil {
		return inconclusiveCheck(nameExecuteWorktree, "command runner is absent")
	}
	plan, ok := findPlanWorktree(ctx)
	if !ok {
		return failCheck(nameExecuteWorktree, "no worktree on cc/p1000/*")
	}
	if ctx.Before.Heads["api"] != ctx.After.Heads["api"] || len(statusLines(ctx.After.Status["api"])) != 0 {
		return failCheck(nameExecuteWorktree, "bound checkout moved or is dirty")
	}
	if !boundOnMain(ctx) {
		return failCheck(nameExecuteWorktree, "bound checkout is not on main")
	}
	branch := strings.TrimPrefix(plan.branch, "refs/heads/")
	count, err := ctx.Exec(repoDir(ctx.World, "api"), "git", "rev-list", "--count", "main.."+branch)
	if err != nil {
		return failCheck(nameExecuteWorktree, err.Error())
	}
	n, _ := strconv.Atoi(strings.TrimSpace(count))
	if n < 1 {
		return failCheck(nameExecuteWorktree, "worktree branch is not ahead of main")
	}
	status, err := ctx.Exec(plan.path, "git", "status", "--porcelain")
	if err != nil {
		return failCheck(nameExecuteWorktree, err.Error())
	}
	if strings.TrimRight(status, "\n") != "" {
		return failCheck(nameExecuteWorktree, "worktree is dirty")
	}
	if !treeHasCSV(plan.path) {
		return failCheck(nameExecuteWorktree, "internal/invoice has no CSV export")
	}
	return passCheck(nameExecuteWorktree)
}

func checkExecuteTests(ctx CheckContext) CheckResult {
	if ctx.Exec == nil {
		return inconclusiveCheck(nameExecuteTests, "command runner is absent")
	}
	plan, ok := findPlanWorktree(ctx)
	if !ok {
		return failCheck(nameExecuteTests, "no worktree on cc/p1000/*")
	}
	if _, err := ctx.Exec(plan.path, "go", "test", "./..."); err != nil {
		detail := err.Error()
		if len(detail) > 500 {
			detail = detail[:500]
		}
		return failCheck(nameExecuteTests, detail)
	}
	return passCheck(nameExecuteTests)
}

func checkExecuteRecords(ctx CheckContext) CheckResult {
	if ctx.Before.Heads["ws"] != ctx.After.Heads["ws"] {
		return failCheck(nameExecuteRecords, "workspace HEAD changed")
	}
	for _, line := range statusLines(ctx.After.Status["ws"]) {
		rel := porcelainPath(line)
		if strings.HasPrefix(rel, "intent/") || strings.HasPrefix(rel, "plans/") || strings.HasPrefix(rel, ".context-circuit/ids/") {
			return failCheck(nameExecuteRecords, "record changed: "+rel)
		}
	}
	files := matchFiles(ctx.After.Files, "plans/p1000-*/plan.md")
	if len(files) != 1 {
		return failCheck(nameExecuteRecords, "expected one plan.md")
	}
	text, err := readRel(ctx, files[0])
	if err != nil {
		return failCheck(nameExecuteRecords, err.Error())
	}
	front, _, ok := splitFront(text)
	if !ok || frontHas(front, "completed_at") {
		return failCheck(nameExecuteRecords, "plan has completed_at")
	}
	return passCheck(nameExecuteRecords)
}

func checkCompletePlan(ctx CheckContext) CheckResult {
	files := matchFiles(ctx.After.Files, "plans/p1000-*/plan.md")
	if len(files) != 1 {
		return failCheck(nameCompletePlan, "expected one plan.md")
	}
	text, err := readRel(ctx, files[0])
	if err != nil {
		return failCheck(nameCompletePlan, err.Error())
	}
	front, body, ok := splitFront(text)
	if !ok || !frontHas(front, "completed_at") || !strings.Contains(body, "## Completion") {
		return failCheck(nameCompletePlan, "plan is missing completed_at or a completion note")
	}
	return passCheck(nameCompletePlan)
}

func checkCompleteNote(ctx CheckContext) CheckResult {
	note, err := readRel(ctx, "context/domains/invoices.md")
	if err != nil {
		return failCheck(nameCompleteNote, err.Error())
	}
	if !csvExportRecorded(note) {
		return failCheck(nameCompleteNote, "invoices note does not record that CSV export exists")
	}
	index, err := readRel(ctx, "context/INDEX.md")
	if err != nil {
		return failCheck(nameCompleteNote, err.Error())
	}
	found := false
	for _, line := range strings.Split(index, "\n") {
		if !strings.Contains(line, "domains/invoices.md") {
			continue
		}
		date := reviewedDate.FindStringSubmatch(strings.TrimRight(line, " \t"))
		got := ""
		if date != nil {
			got = date[1]
		}
		if !nearUTCDay(got, time.Now()) {
			return failCheck(nameCompleteNote, "reviewed date is "+got+", want the run's UTC date")
		}
		found = true
	}
	if !found {
		return failCheck(nameCompleteNote, "invoices catalog entry is missing")
	}
	var named string
	err = filepath.WalkDir(filepath.Join(ctx.World.WS, "context"), func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return err
		}
		body, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if recordIDMention.Match(body) {
			rel, _ := filepath.Rel(ctx.World.WS, p)
			named = filepath.ToSlash(rel)
		}
		return nil
	})
	if err != nil {
		return failCheck(nameCompleteNote, err.Error())
	}
	if named != "" {
		return failCheck(nameCompleteNote, named+" names a record id")
	}
	return passCheck(nameCompleteNote)
}

func checkDirectReadme(ctx CheckContext) CheckResult {
	if ctx.Before.Heads["api"] != ctx.After.Heads["api"] {
		return failCheck(nameDirectReadme, "api HEAD changed")
	}
	lines := statusLines(ctx.After.Status["api"])
	if len(lines) != 1 || lines[0] != " M README.md" {
		return failCheck(nameDirectReadme, "api status is not exactly one README edit")
	}
	body, err := os.ReadFile(filepath.Join(ctx.World.API, "README.md"))
	if err != nil {
		body, err = os.ReadFile(filepath.Join(ctx.World.WS, "repositories", "api", "README.md"))
	}
	if err != nil {
		return failCheck(nameDirectReadme, err.Error())
	}
	text := string(body)
	if strings.Contains(text, "recieve") || !strings.Contains(text, "receive") {
		return failCheck(nameDirectReadme, "README still has the typo or lost receive")
	}
	return passCheck(nameDirectReadme)
}

func checkDirectNoRecords(ctx CheckContext) CheckResult {
	for _, line := range statusLines(ctx.After.Status["ws"]) {
		rel := porcelainPath(line)
		if strings.HasPrefix(rel, "intent/") || strings.HasPrefix(rel, "plans/") || strings.HasPrefix(rel, ".context-circuit/ids/") {
			return failCheck(nameDirectNoRecords, "record path changed: "+rel)
		}
	}
	if addedUnder(ctx, "intent/") != nil || addedUnder(ctx, "plans/") != nil || addedUnder(ctx, ".context-circuit/ids/") != nil {
		return failCheck(nameDirectNoRecords, "a record file appeared")
	}
	if ctx.Before.Worktrees["api"] != ctx.After.Worktrees["api"] {
		return failCheck(nameDirectNoRecords, "a worktree appeared")
	}
	return passCheck(nameDirectNoRecords)
}

func checkBorrowedUntouched(ctx CheckContext) CheckResult {
	if ctx.Before.Heads["handbook"] == "" || ctx.Before.Heads["handbook"] != ctx.After.Heads["handbook"] {
		return failCheck(nameBorrowedUntouched, "handbook HEAD changed")
	}
	if len(statusLines(ctx.After.Status["handbook"])) != 0 {
		return failCheck(nameBorrowedUntouched, "handbook checkout is dirty")
	}
	if ctx.Before.Refs["handbook"] != ctx.After.Refs["handbook"] {
		return failCheck(nameBorrowedUntouched, "handbook branches changed")
	}
	if ctx.Before.RemoteRefs["handbook.git"] != ctx.After.RemoteRefs["handbook.git"] {
		return failCheck(nameBorrowedUntouched, "handbook remote changed")
	}
	for _, line := range statusLines(ctx.After.Status["ws"]) {
		rel := porcelainPath(line)
		if rel == "context" || strings.HasPrefix(rel, "context/") {
			return failCheck(nameBorrowedUntouched, "context changed: "+rel)
		}
	}
	return passCheck(nameBorrowedUntouched)
}

func checkBorrowedReply(ctx CheckContext) CheckResult {
	reply := ctx.Parsed.FinalReply
	if borrowedUpstream.MatchString(reply) && borrowedHandbook.MatchString(reply) && borrowedReview.MatchString(reply) {
		return passCheck(nameBorrowedReply)
	}
	return failCheck(nameBorrowedReply, "reply does not send the fix upstream to the handbook's review")
}

type wtEntry struct {
	path   string
	branch string
}

func findPlanWorktree(ctx CheckContext) (wtEntry, bool) {
	for _, tree := range parseWorktrees(ctx.After.Worktrees["api"]) {
		name := strings.TrimPrefix(tree.branch, "refs/heads/")
		ok, err := path.Match("cc/p1000/*", name)
		if err == nil && ok {
			return tree, true
		}
	}
	return wtEntry{}, false
}

func boundOnMain(ctx CheckContext) bool {
	api := ctx.World.API
	for _, tree := range parseWorktrees(ctx.After.Worktrees["api"]) {
		if pathsEqual(tree.path, api) {
			return tree.branch == "refs/heads/main"
		}
	}
	return false
}

func parseWorktrees(text string) []wtEntry {
	var out []wtEntry
	var cur wtEntry
	flush := func() {
		if cur.path != "" {
			out = append(out, cur)
		}
		cur = wtEntry{}
	}
	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.TrimSpace(line) == "":
			flush()
		case strings.HasPrefix(line, "worktree "):
			flush()
			cur.path = strings.TrimPrefix(line, "worktree ")
		case strings.HasPrefix(line, "branch "):
			cur.branch = strings.TrimPrefix(line, "branch ")
		}
	}
	flush()
	return out
}

func treeHasCSV(root string) bool {
	dir := filepath.Join(root, "internal", "invoice")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	fset := token.NewFileSet()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			continue
		}
		if goFileUsesCSV(file) {
			return true
		}
	}
	return false
}

// goFileUsesCSV reports whether a function in the file calls the
// encoding/csv package. A comment or a string that says CSV does not count.
func goFileUsesCSV(file *ast.File) bool {
	pkg := csvImportName(file)
	if pkg == "" {
		return false
	}
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return true
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if pkg == "." {
				if usesCSVDot(n) {
					found = true
				}
				return true
			}
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if ok && id.Name == pkg {
				found = true
			}
			return true
		})
		return !found
	})
	return found
}

func csvImportName(file *ast.File) string {
	for _, imp := range file.Imports {
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil || path != "encoding/csv" {
			continue
		}
		if imp.Name == nil {
			return "csv"
		}
		if imp.Name.Name == "_" {
			return ""
		}
		return imp.Name.Name
	}
	return ""
}

func usesCSVDot(n ast.Node) bool {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}
	id, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	switch id.Name {
	case "NewReader", "NewWriter":
		return true
	default:
		return false
	}
}

// requireRecordedCheckout fails when a HEAD is absent or a status was
// not recorded. An empty status string is a clean tree and is present.
func requireRecordedCheckout(name string, ctx CheckContext, repos ...string) string {
	for _, repo := range repos {
		if ctx.Before.Heads[repo] == "" || ctx.After.Heads[repo] == "" {
			return repo + " HEAD is missing"
		}
		if _, ok := ctx.After.Status[repo]; !ok {
			return repo + " status is missing"
		}
	}
	return ""
}

func resolveBinding(base, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !filepath.IsAbs(raw) {
		raw = filepath.Join(base, raw)
	}
	return raw
}

// roleHosts are the coding hosts whose role files agent setup writes.
// claude and claude-code share a directory: the scenario host is claude,
// and the CLI names that host claude-code.
var roleHosts = []string{"claude", "codex", "cursor"}

var roleFileNames = []string{"cc-explorer", "cc-planner", "cc-worker", "cc-reviewer"}

func roleLocation(host string) (dir, ext string, ok bool) {
	switch host {
	case "claude", "claude-code":
		return ".claude/agents", ".md", true
	case "codex":
		return ".codex/agents", ".toml", true
	case "cursor":
		return ".cursor/agents", ".md", true
	default:
		return "", "", false
	}
}

func roleSetComplete(root, host string) bool {
	dir, ext, ok := roleLocation(host)
	if !ok {
		return false
	}
	for _, name := range roleFileNames {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(dir), name+ext)); err != nil {
			return false
		}
	}
	return true
}

// missingRoleSet names the host whose role files are incomplete.
// A known host requires that host's four files. An empty or unknown
// host requires one host's complete set.
func missingRoleSet(root, host string) string {
	if _, _, ok := roleLocation(host); ok {
		if roleSetComplete(root, host) {
			return ""
		}
		return "missing role definitions for " + host
	}
	for _, candidate := range roleHosts {
		if roleSetComplete(root, candidate) {
			return ""
		}
	}
	return "missing role definitions"
}

func csvExportRecorded(note string) bool {
	for _, sentence := range splitSentences(note) {
		if !csvWord.MatchString(sentence) || csvDenied.MatchString(sentence) {
			continue
		}
		if csvAffirmed.MatchString(sentence) {
			return true
		}
	}
	return false
}

func splitSentences(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		for _, sentence := range strings.Split(line, ".") {
			sentence = strings.TrimSpace(sentence)
			if sentence != "" {
				out = append(out, sentence)
			}
		}
	}
	return out
}

// nearUTCDay accepts the UTC calendar date of now and the day on either
// side, so a note written just before midnight still matches.
func nearUTCDay(got string, now time.Time) bool {
	if got == "" {
		return false
	}
	now = now.UTC()
	for _, delta := range []int{-1, 0, 1} {
		if got == now.AddDate(0, 0, delta).Format("2006-01-02") {
			return true
		}
	}
	return false
}

func pathsEqual(a, b string) bool {
	ar, aerr := filepath.EvalSymlinks(a)
	br, berr := filepath.EvalSymlinks(b)
	if aerr != nil {
		ar = a
	}
	if berr != nil {
		br = b
	}
	return filepath.Clean(ar) == filepath.Clean(br)
}

func readRel(ctx CheckContext, rel string) (string, error) {
	body, err := os.ReadFile(filepath.Join(ctx.World.WS, filepath.FromSlash(rel)))
	return string(body), err
}

func splitFront(text string) (front, body string, ok bool) {
	const open = "---\n"
	if !strings.HasPrefix(text, open) {
		return "", "", false
	}
	rest := text[len(open):]
	i := strings.Index(rest, "\n---\n")
	if i < 0 {
		return "", "", false
	}
	return rest[:i], rest[i+len("\n---\n"):], true
}

func frontScalar(front, key string) (string, bool) {
	prefix := key + ":"
	for _, line := range strings.Split(front, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, prefix) {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(trim, prefix)), `"'`), true
		}
	}
	return "", false
}

func frontHas(front, key string) bool {
	_, ok := frontScalar(front, key)
	return ok
}

func frontList(front, key string) []string {
	var out []string
	in := false
	for _, line := range strings.Split(front, "\n") {
		trim := strings.TrimSpace(line)
		if !in {
			if !strings.HasPrefix(trim, key+":") {
				continue
			}
			rest := strings.TrimSpace(strings.TrimPrefix(trim, key+":"))
			if rest == "" {
				in = true
				continue
			}
			rest = strings.Trim(rest, "[]")
			for _, part := range strings.Split(rest, ",") {
				part = strings.Trim(strings.TrimSpace(part), `"'`)
				if part != "" {
					out = append(out, part)
				}
			}
			return out
		}
		if trim == "" {
			continue
		}
		if strings.HasPrefix(trim, "- ") {
			out = append(out, strings.Trim(strings.TrimSpace(strings.TrimPrefix(trim, "- ")), `"'`))
			continue
		}
		break
	}
	return out
}

func containsValue(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func headingsPresent(body string, titles []string) bool {
	for _, title := range titles {
		if !hasHeading(body, title) {
			return false
		}
	}
	return true
}

func hasHeading(body, title string) bool {
	for _, line := range strings.Split(body, "\n") {
		if strings.TrimSpace(line) == "## "+title {
			return true
		}
	}
	return false
}

func headingTitle(line string) (string, bool) {
	trim := strings.TrimSpace(line)
	if !strings.HasPrefix(trim, "## ") || strings.HasPrefix(trim, "###") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(trim, "## ")), true
}

func section(body, title string) (string, bool) {
	lines := strings.Split(body, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "## "+title {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return "", false
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "## ") {
			end = i
			break
		}
	}
	return strings.Join(lines[start:end], "\n"), true
}

func openQuestionsOK(body string) bool {
	text, ok := section(body, "Open questions")
	if !ok {
		return false
	}
	text = strings.TrimSpace(text)
	return text == "None known." || numberedList.MatchString(text)
}

func hasPlaceholder(body string) bool {
	for _, sentence := range placeholderSentences {
		if strings.Contains(body, sentence) {
			return true
		}
	}
	return false
}

func linksAssigned(front, body string) (bool, string) {
	sec, ok := section(body, "Details")
	if !ok {
		return true, ""
	}
	assigned := assignedFiles(front)
	for _, match := range mdLink.FindAllStringSubmatch(sec, -1) {
		target := match[1]
		if strings.Contains(target, "://") {
			continue
		}
		if !assigned[target] {
			return false, target
		}
	}
	return true, ""
}

// planHasSupportingFiles reports whether the plan names a supporting file:
// an entry in required_files, or a relative link anywhere in the body to a
// file inside the plan folder. A link with a scheme, an absolute path, an
// anchor, or a path leaving the folder is not one.
func planHasSupportingFiles(front, body string) bool {
	if len(assignedFiles(front)) > 0 {
		return true
	}
	for _, match := range mdLink.FindAllStringSubmatch(body, -1) {
		if inPlanFolder(match[1]) {
			return true
		}
	}
	return false
}

func inPlanFolder(target string) bool {
	if i := strings.IndexAny(target, "#?"); i >= 0 {
		target = target[:i]
	}
	if target == "" || strings.Contains(target, "://") || strings.HasPrefix(target, "/") || strings.HasPrefix(target, "mailto:") {
		return false
	}
	clean := path.Clean(target)
	return clean != "." && clean != ".." && !strings.HasPrefix(clean, "../")
}

func assignedFiles(front string) map[string]bool {
	out := map[string]bool{}
	var doc struct {
		RequiredFiles *struct {
			Shared       []string            `yaml:"shared"`
			ByRepository map[string][]string `yaml:"by_repository"`
		} `yaml:"required_files"`
	}
	if yaml.Unmarshal([]byte(front), &doc) != nil || doc.RequiredFiles == nil {
		return out
	}
	for _, file := range doc.RequiredFiles.Shared {
		out[file] = true
	}
	for _, files := range doc.RequiredFiles.ByRepository {
		for _, file := range files {
			out[file] = true
		}
	}
	return out
}

func proseWords(body string) []string {
	var b strings.Builder
	for _, line := range strings.Split(body, "\n") {
		if _, ok := headingTitle(line); ok || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		b.WriteString(line)
		b.WriteByte(' ')
	}
	return wordPattern.FindAllString(b.String(), -1)
}

func englishRatio(words []string) float64 {
	if len(words) == 0 {
		return 0
	}
	n := 0
	for _, word := range words {
		if englishStopwords[strings.ToLower(word)] {
			n++
		}
	}
	return float64(n) / float64(len(words))
}

func matchFiles(files []string, pattern string) []string {
	var out []string
	for _, file := range files {
		ok, err := path.Match(pattern, file)
		if err == nil && ok {
			out = append(out, file)
		}
	}
	return out
}

func addedMatch(before, after []string, pattern string) []string {
	have := map[string]bool{}
	for _, file := range before {
		have[file] = true
	}
	var out []string
	for _, file := range matchFiles(after, pattern) {
		if !have[file] {
			out = append(out, file)
		}
	}
	return out
}

func addedUnder(ctx CheckContext, prefix string) []string {
	have := map[string]bool{}
	for _, file := range ctx.Before.Files {
		have[file] = true
	}
	var out []string
	for _, file := range ctx.After.Files {
		if strings.HasPrefix(file, prefix) && !have[file] {
			out = append(out, file)
		}
	}
	return out
}

func statusLines(status string) []string {
	status = strings.TrimRight(status, "\n")
	if status == "" {
		return nil
	}
	return strings.Split(status, "\n")
}

func porcelainPath(line string) string {
	if len(line) < 4 {
		return strings.TrimSpace(line)
	}
	return line[3:]
}
