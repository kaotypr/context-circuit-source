import assert from "node:assert/strict";
import { mkdir, readFile, readdir, symlink, writeFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { git } from "../scripts/lib/git.js";
import { createPlanDraft, setPlanState } from "../scripts/lib/plans.js";
import { preparePlanTask } from "../scripts/lib/run-task.js";
import type { FakeActivitySource, PlanDraftRequest, PlanRunTaskRequest, WorkCandidate } from "../scripts/lib/types.js";
import { recommendWhatsNext } from "../scripts/lib/whats-next.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace } from "./helpers.js";

function planRequest(workItems: PlanDraftRequest["work_items"] = [
  { key: "reset", title: "Add reset behavior", area: "application foundation", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only" as const, verification_commands: [], acceptance_criteria: ["Reset behavior works."] },
]): PlanDraftRequest {
  return {
    contract_version: 1,
    plan_id: "counter-reset",
    title: "Counter reset",
    source: { kind: "idea", reference: "human request" },
    work_prefix: "RESET",
    summary: "Add an explicit reset path.",
    affected_repositories: ["frontend"],
    assumptions: [],
    open_questions: [],
    requirements: ["Reset returns the displayed count to zero."],
    solution: ["Add reset behavior to the existing counter."],
    delivery: ["Implement the dependency-ready work items in order."],
    verification: ["Run the repository test suite."],
    risks: [],
    work_items: workItems,
  };
}

function candidate(overrides: Partial<WorkCandidate> = {}): WorkCandidate {
  return {
    contract_version: 1,
    candidate_id: "activity:ready",
    kind: "activity-task",
    work_id: null,
    title: "Ready activity task",
    state: "ready",
    urgent: false,
    priority: 50,
    owner: null,
    plan_reference: null,
    plan_approval_state: "not-applicable",
    dependencies: [],
    scope_sufficient: true,
    acceptance_sufficient: true,
    repositories: ["frontend"],
    access_available: true,
    contract_blocked: false,
    source_reference: "fake-activity://ready",
    risks: [],
    ...overrides,
  };
}

async function preparePlanRun(workspace: { root: string }, now = new Date("2026-08-11T09:30:00Z")) {
  const created = await createPlanDraft(workspace.root, planRequest(), new Date("2026-08-11T08:00:00Z"));
  const approved = await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-11T09:00:00Z"));
  const request: PlanRunTaskRequest = {
    contract_version: 1,
    source: { kind: "plan", reference: "context/plans/counter-reset", plan_version: approved.plan_version, approved_digest: approved.approved_digest! },
    work_ids: ["RESET-001"],
  };
  return preparePlanTask({ workspaceRoot: workspace.root, request, now, discriminator: "abcdef12" });
}

test("whats-next recommends a dependency-ready item from an approved plan", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, planRequest(), new Date("2026-08-11T08:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-11T09:00:00Z"));

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.action, "execute");
  assert.equal(result.recommendation.candidate_id, "plan:counter-reset:RESET-001");
  assert.match(result.recommendation.why, /dependency-ready item in an approved plan/);
  assert.deepEqual(result.recommendation.repositories, ["frontend"]);
  assert.equal(result.no_state_changed, true);
  assert.deepEqual(await validateContract("whats-next-result", result), []);
});

test("whats-next skips plans whose explicit repository is no longer registered", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, planRequest(), new Date("2026-08-11T08:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-11T09:00:00Z"));
  const workspacePath = join(workspace.root, "workspace.yaml");
  await writeFile(workspacePath, (await readFile(workspacePath, "utf8")).replace("  frontend:\n", "  replacement:\n"), "utf8");
  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));
  assert.equal(result.recommendation.action, "enable");
  assert.match(result.warnings.join("\n"), /repository is not registered: frontend/);
});

test("whats-next ranks urgent work before in-progress and priority-ready work", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const activity: FakeActivitySource = {
    contract_version: 1,
    current_user: "kao",
    candidates: [
      candidate({ candidate_id: "activity:priority", priority: 100 }),
      candidate({ candidate_id: "activity:progress", title: "Continue active work", state: "in-progress", priority: 1, owner: "kao" }),
      candidate({ candidate_id: "activity:urgent", title: "Resolve production incident", urgent: true, priority: 0 }),
    ],
  };

  const result = await recommendWhatsNext(workspace.root, activity, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.candidate_id, "activity:urgent");
  assert.deepEqual(result.alternatives.map((item) => item.candidate_id), ["activity:progress", "activity:priority"]);
});

test("fake activity facts safely establish plan dependency completion", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, planRequest([
    { key: "contract", title: "Define reset contract", area: "frontend", repository: "frontend", scope: ["src/contract.ts"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Reset contract is defined."] },
    { key: "ui", title: "Implement reset UI", area: "frontend", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Reset UI is implemented."], depends_on: ["contract"] },
  ]), new Date("2026-08-11T08:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-11T09:00:00Z"));
  const activity: FakeActivitySource = {
    contract_version: 1,
    current_user: "kao",
    candidates: [
      candidate({ candidate_id: "activity:reset-contract", work_id: "RESET-001", state: "completed", source_reference: "fake-activity://RESET-001" }),
    ],
  };

  const result = await recommendWhatsNext(workspace.root, activity, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.candidate_id, "plan:counter-reset:RESET-010");
  assert.ok(result.recommendation.readiness_evidence.includes("dependencies: RESET-001=completed"));
  assert.equal(result.considered.excluded, 1);
});

test("draft plans yield a concrete approval action instead of invented implementation", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await createPlanDraft(workspace.root, planRequest(), new Date("2026-08-11T08:00:00Z"));

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.action, "enable");
  assert.match(result.recommendation.title, /Approve the governing plan/);
  assert.match(result.recommendation.blockers.join("\n"), /governing plan is draft/);
});

test("whats-next is read-only even when a repository contains unrecorded work", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, planRequest(), new Date("2026-08-11T08:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-11T09:00:00Z"));
  await writeFile(join(workspace.repository, "src", "App.tsx"), "// unrecorded user work\n", { flag: "a" });
  const statusBefore = await git(workspace.repository, ["status", "--porcelain=v1"]);
  const planBefore = await readFile(join(created.directory, "README.md"), "utf8");
  const pathsBefore = (await readdir(workspace.root, { recursive: true })).sort();

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.no_state_changed, true);
  assert.equal(await git(workspace.repository, ["status", "--porcelain=v1"]), statusBefore);
  assert.equal(await readFile(join(created.directory, "README.md"), "utf8"), planBefore);
  assert.deepEqual((await readdir(workspace.root, { recursive: true })).sort(), pathsBefore);
});

test("whats-next rejects malformed fake activity input before recommendation", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const malformed = { contract_version: 1, current_user: "kao", candidates: [{ title: "missing fields" }] } as unknown as FakeActivitySource;
  await assert.rejects(recommendWhatsNext(workspace.root, malformed), /Invalid fake activity source/);
});

test("passed plan-linked runtime is recommended for review instead of duplicate execution", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanRun(workspace);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  manifest.status = "passed";
  manifest.repositories[0].status = "passed";
  manifest.plan_work_items[0].outcome = "passed";
  await writeFile(prepared.manifest, `${JSON.stringify(manifest, null, 2)}\n`);

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.action, "review");
  assert.equal(result.recommendation.candidate_id, "plan:counter-reset:RESET-001");
  assert.match(result.recommendation.title, /Review or prepare merge/);
  assert.ok(result.recommendation.source_references.some((item) => item.endsWith("/manifest.json")));
  assert.ok(result.recommendation.source_references.some((item) => item.includes(".runtime/tasks/")));
  assert.equal(result.considered.executable, 1);
});

test("forged plan task brief run identity and approval digest cannot project runtime state", async (t) => {
  for (const forgery of ["run", "digest"] as const) {
    await t.test(forgery, async (t) => {
      const workspace = await createTestWorkspace();
      t.after(workspace.cleanup);
      const prepared = await preparePlanRun(workspace);
      const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
      manifest.status = "passed";
      manifest.repositories[0].status = "passed";
      manifest.plan_work_items[0].outcome = "passed";
      await writeFile(prepared.manifest, `${JSON.stringify(manifest, null, 2)}\n`);
      const brief = JSON.parse(await readFile(prepared.taskBrief, "utf8"));
      if (forgery === "run") brief.run_id = "20260811T093000Z-deadbeef";
      else brief.plan.approved_digest = `sha256:${"0".repeat(64)}`;
      await writeFile(prepared.taskBrief, `${JSON.stringify(brief, null, 2)}\n`);

      const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

      assert.equal(result.recommendation.action, "execute");
      assert.ok(result.warnings.some((item) => item.includes("Ignored malformed runtime evidence") && item.includes(forgery === "run" ? "run or work identity" : "version or digest")));
    });
  }
});

test("closed plan-linked runtime is excluded without an activity provider", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanRun(workspace);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  manifest.status = "closed";
  manifest.repositories[0].status = "closed";
  manifest.plan_work_items[0].outcome = "passed";
  await writeFile(prepared.manifest, `${JSON.stringify(manifest, null, 2)}\n`);

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.candidate_id, null);
  assert.equal(result.considered.excluded, 1);
  assert.equal(result.no_state_changed, true);
});

test("closing plan-linked runtime is recommended for closeout instead of execution", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanRun(workspace);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  manifest.status = "closing";
  manifest.repositories[0].status = "closing";
  manifest.plan_work_items[0].outcome = "passed";
  await writeFile(prepared.manifest, `${JSON.stringify(manifest, null, 2)}\n`);

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.action, "closeout");
  assert.match(result.recommendation.title, /Complete closeout or cleanup/);
  assert.equal(result.no_state_changed, true);
});

test("runtime and activity contradiction yields read-only reconciliation", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanRun(workspace);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  manifest.status = "passed";
  manifest.repositories[0].status = "passed";
  manifest.plan_work_items[0].outcome = "passed";
  await writeFile(prepared.manifest, `${JSON.stringify(manifest, null, 2)}\n`);
  const activity: FakeActivitySource = {
    contract_version: 1,
    current_user: "kao",
    candidates: [candidate({ candidate_id: "activity:stale-reset", work_id: "RESET-001", state: "ready", source_reference: "fake-activity://RESET-001" })],
  };

  const result = await recommendWhatsNext(workspace.root, activity, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.action, "reconcile");
  assert.match(result.recommendation.title, /Reconcile contradictory state/);
  assert.match(result.recommendation.blockers.join("\n"), /contradictory work states/);
  assert.ok(result.recommendation.source_references.includes("fake-activity://RESET-001"));
  assert.ok(result.recommendation.source_references.some((item) => item.endsWith("/manifest.json")));
  assert.equal(result.no_state_changed, true);
});

test("malformed and unrelated runtime evidence is ignored safely", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, planRequest(), new Date("2026-08-11T08:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-11T09:00:00Z"));
  const malformed = join(workspace.root, ".runtime", "runs", "malformed");
  const unrelated = join(workspace.root, ".runtime", "runs", "unrelated");
  await mkdir(malformed, { recursive: true });
  await mkdir(unrelated, { recursive: true });
  await writeFile(join(malformed, "manifest.json"), "{ definitely not json\n");
  await writeFile(join(unrelated, "note.txt"), "not runtime evidence\n");

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.action, "execute");
  assert.ok(result.warnings.some((item) => item.includes("Ignored malformed runtime evidence malformed")));
  assert.ok(!result.warnings.some((item) => item.includes("unrelated")));
});

test("tracked contribution with a nonexistent run is warned and cannot establish completion", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, planRequest(), new Date("2026-08-11T08:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-11T09:00:00Z"));
  const contributionDirectory = join(workspace.root, "contributions", "general");
  const contribution = join(contributionDirectory, "20260811T093000Z-reviewer-reset-frontend.md");
  await mkdir(contributionDirectory, { recursive: true });
  await writeFile(contribution, `# RESET-001: Add reset behavior

- Run: \`20260811T090000Z-abcdef12\`

## Outcome

Merged after human review.

## Affected repositories

- \`frontend\`.

## Pull requests and commits

- Recorded.

## Verification

- Passed.

## Decisions and deviations

- None.

## Remaining risks and follow-up

- None.

## Candidate durable learnings

- None.
`);
  await git(workspace.root, ["init", "--initial-branch=main"]);
  await git(workspace.root, ["add", "."]);
  await git(workspace.root, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "record completed work"]);

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.action, "execute");
  assert.equal(result.considered.excluded, 0);
  assert.ok(result.warnings.some((item) => item.includes("Ignored unassociated durable contribution") && item.includes("cited run and closeout relationship")));
});

test("tracked contribution only becomes terminal through a matching validated plan closeout", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanRun(workspace);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const contributionDirectory = join(workspace.root, "contributions", "general");
  const contributionRelative = "contributions/general/20260811T093000Z-reviewer-reset-frontend.md";
  const contribution = join(workspace.root, contributionRelative);
  const closeoutPath = join(workspace.root, ".runtime", "runs", prepared.runId, "frontend-closeout.json");
  await mkdir(contributionDirectory, { recursive: true });
  await writeFile(contribution, `# RESET-001: Add reset behavior

- Run: \`${prepared.runId}\`

## Outcome

Merged after human review.

## Affected repositories

- \`frontend\`.

## Pull requests and commits

- Recorded.

## Verification

- Passed.

## Decisions and deviations

- None.

## Remaining risks and follow-up

- None.

## Candidate durable learnings

- None.
`);
  await writeFile(closeoutPath, `${JSON.stringify({
    contract_version: 1, work_id: "RESET-001", run_id: prepared.runId, repository: "frontend", outcome: "merged", status: "prepared", author: "reviewer", reason: null,
    contribution: contributionRelative, branch: manifest.repositories[0].branch, base_commit: manifest.repositories[0].base_commit, head_commit: manifest.repositories[0].base_commit,
    merge_commit: null, pull_requests: [], commits: [], changed_files: [], verification: ["Passed"], cleanup: { requested: false, worktree_removed: false, branch_preserved: true, runtime_evidence_preserved: true },
    blockers: [], prepared_at: "2026-08-11T09:30:00Z", updated_at: "2026-08-11T09:30:00Z",
  }, null, 2)}\n`);
  manifest.status = "closing";
  manifest.repositories[0].status = "closing";
  manifest.repositories[0].closeout_record = closeoutPath;
  manifest.repositories[0].contribution = contributionRelative;
  manifest.plan_work_items[0].outcome = "passed";
  await writeFile(prepared.manifest, `${JSON.stringify(manifest, null, 2)}\n`);
  await git(workspace.root, ["init", "--initial-branch=main"]);
  await git(workspace.root, ["add", contributionRelative]);
  await git(workspace.root, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "record completed work"]);

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.candidate_id, null);
  assert.equal(result.considered.excluded, 1);
  assert.equal(result.warnings.length, 0);
});

test("symlinked runtime and contribution entries are warned and never followed", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, planRequest(), new Date("2026-08-11T08:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-11T09:00:00Z"));
  await mkdir(join(workspace.root, ".runtime", "runs"), { recursive: true });
  await mkdir(join(workspace.root, "contributions"), { recursive: true });
  await symlink(workspace.repository, join(workspace.root, ".runtime", "runs", "outside-run"));
  await symlink(workspace.repository, join(workspace.root, "contributions", "outside-group"));

  const result = await recommendWhatsNext(workspace.root, null, new Date("2026-08-11T10:00:00Z"));

  assert.equal(result.recommendation.action, "execute");
  assert.ok(result.warnings.includes("Ignored symlinked runtime evidence outside-run"));
  assert.ok(result.warnings.includes("Ignored symlinked contribution group outside-group"));
});
