import assert from "node:assert/strict";
import { readFile, readdir, writeFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { git } from "../scripts/lib/git.js";
import { createPlanDraft, setPlanState } from "../scripts/lib/plans.js";
import type { FakeActivitySource, PlanDraftRequest, WorkCandidate } from "../scripts/lib/types.js";
import { recommendWhatsNext } from "../scripts/lib/whats-next.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace } from "./helpers.js";

function planRequest(workItems: PlanDraftRequest["work_items"] = [
  { key: "reset", title: "Add reset behavior", area: "frontend" },
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
    { key: "contract", title: "Define reset contract", area: "frontend" },
    { key: "ui", title: "Implement reset UI", area: "frontend", depends_on: ["contract"] },
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
