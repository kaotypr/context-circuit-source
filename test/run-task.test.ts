import assert from "node:assert/strict";
import { access, mkdir, readFile, symlink, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import test from "node:test";
import { createPlanDraft, setPlanState } from "../scripts/lib/plans.js";
import { git } from "../scripts/lib/git.js";
import { recordResult } from "../scripts/lib/record-result.js";
import { preparePlanTask, preparePlanlessTask } from "../scripts/lib/run-task.js";
import type { PlanDraftRequest, PlanRunTaskRequest } from "../scripts/lib/types.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace, taskOptions } from "./helpers.js";

const planDraft: PlanDraftRequest = {
  contract_version: 1,
  plan_id: "reset-flow",
  title: "Reset flow",
  source: { kind: "prd", reference: "docs/reset.md" },
  work_prefix: "RESET",
  summary: "Deliver reset behavior.",
  affected_repositories: ["frontend"],
  assumptions: [], open_questions: [],
  requirements: ["Reset returns the count to zero."],
  solution: ["Implement and verify reset behavior."],
  delivery: ["Establish the contract before the UI."],
  verification: ["Run the focused tests."], risks: [],
  work_items: [
    { key: "contract", title: "Define reset behavior", area: "reset contract", repository: "frontend", scope: ["src/App.tsx"], test_scope: ["src/App.test.tsx"], test_policy: "required", verification_commands: ["npm test"], acceptance_criteria: ["Reset returns the count to zero."] },
    { key: "ui", title: "Implement reset UI", area: "reset interface", repository: "frontend", scope: ["src/App.tsx"], test_scope: ["src/App.test.tsx"], test_policy: "required", verification_commands: ["npm test"], acceptance_criteria: ["The reset control is visible."], depends_on: ["contract"] },
  ],
};

function planRequest(planVersion: number, approvedDigest: string): PlanRunTaskRequest {
  return {
    contract_version: 1,
    source: { kind: "plan", reference: "context/plans/reset-flow", plan_version: planVersion, approved_digest: approvedDigest },
    work_ids: ["RESET-001"],
  };
}

test("planless preparation writes valid scoped evidence and creates a worktree", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    testScope: ["src/App.test.tsx"],
    now: new Date("2026-08-11T08:30:00.000Z"),
    discriminator: "abcd1234",
  });
  assert.equal(prepared.workId, "ADHOC-20260811-001");
  assert.match(prepared.runId, /^20260811T083000Z-abcd[0-9a-f]{4}$/);
  await access(join(prepared.worktree, ".git"));

  const brief = JSON.parse(await readFile(prepared.taskBrief, "utf8"));
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const worker = JSON.parse(await readFile(prepared.workerInput, "utf8"));
  const verifier = JSON.parse(await readFile(prepared.verifierInput, "utf8"));
  assert.deepEqual(await validateContract("task-brief", brief), []);
  assert.deepEqual(await validateContract("runtime-manifest", manifest), []);
  assert.equal(manifest.status, "prepared");
  assert.equal(brief.plan.approval_state, "not-applicable");
  assert.equal(brief.activity.duplicate_effort_warning, true);
  assert.deepEqual(brief.implementation_scope, ["src/App.tsx"]);
  assert.deepEqual(brief.test_expectation, {
    policy: "required",
    paths: ["src/App.test.tsx"],
    rationale: "The worker must add or update tests in the declared test scope.",
  });
  assert.deepEqual(worker.allowed_scope, ["src/App.tsx", "src/App.test.tsx"]);
  assert.deepEqual(worker.test_expectation, brief.test_expectation);
  await access(worker.result_contract);
  assert.equal(verifier.read_only, true);
  assert.deepEqual(verifier.test_expectation, brief.test_expectation);
  await access(verifier.result_contract);
});

test("dirty repository protection preserves changes and creates no run", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const dirtyPath = join(workspace.repository, "unrecorded.txt");
  await writeFile(dirtyPath, "preserve me\n", "utf8");
  await assert.rejects(
    preparePlanlessTask({ workspaceRoot: workspace.root, ...taskOptions }),
    /unresolved local changes/,
  );
  assert.equal(await readFile(dirtyPath, "utf8"), "preserve me\n");
  await assert.rejects(access(join(workspace.root, ".runtime")));
});

test("runtime evidence refuses a symlinked runtime root", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const target = join(workspace.root, "evidence-target");
  await mkdir(target);
  await symlink(target, join(workspace.root, ".runtime"), "dir");
  await assert.rejects(
    preparePlanlessTask({ workspaceRoot: workspace.root, ...taskOptions }),
    /Runtime path must be a real directory/,
  );
  await assert.rejects(access(join(target, "id-state.json")));
});

test("existing-coverage policy rejects a missing test path", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await assert.rejects(
    preparePlanlessTask({
      workspaceRoot: workspace.root,
      ...taskOptions,
      testScope: ["src/missing.test.tsx"],
      testPolicy: "existing-coverage",
    }),
    /Existing-coverage test path does not exist/,
  );
});

test("separate preparations use isolated worktrees and branches", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const first = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T09:00:00.000Z"),
    discriminator: "11111111",
  });
  const second = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T09:00:01.000Z"),
    discriminator: "22222222",
  });
  assert.notEqual(first.branch, second.branch);
  assert.notEqual(first.worktree, second.worktree);
  assert.equal(first.workId, "ADHOC-20260811-001");
  assert.equal(second.workId, "ADHOC-20260811-002");
  await writeFile(join(first.worktree, "isolation.txt"), "first only\n", "utf8");
  await assert.rejects(access(join(second.worktree, "isolation.txt")));
});

test("approved plan selection preserves stable work IDs and approval evidence", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, planDraft, new Date("2026-08-11T07:00:00Z"));
  const approved = await setPlanState(created.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-11T07:30:00Z"));
  const prepared = await preparePlanTask({
    workspaceRoot: workspace.root,
    request: planRequest(approved.plan_version, approved.approved_digest!),
    now: new Date("2026-08-11T08:30:00Z"),
    discriminator: "abcd1234",
  });
  assert.equal(prepared.workId, "RESET-001");
  const brief = JSON.parse(await readFile(prepared.taskBrief, "utf8"));
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(brief.source.kind, "plan");
  assert.equal(brief.source.reference, "context/plans/reset-flow");
  assert.deepEqual(brief.plan, {
    reference: "context/plans/reset-flow", approval_state: "approved", plan_version: 1,
    approved_digest: approved.approved_digest, work_ids: ["RESET-001"],
  });
  assert.equal(brief.requested_outcome, "Define reset behavior");
  assert.deepEqual(brief.scope, ["src/App.tsx", "src/App.test.tsx"]);
  assert.deepEqual(brief.acceptance_criteria, ["Reset returns the count to zero."]);
  assert.deepEqual(manifest.plan_work_items, [
    { work_id: "RESET-001", repository: "frontend", depends_on: [], outcome: "pending" },
  ]);
  assert.deepEqual(await validateContract("task-brief", brief), []);
  assert.deepEqual(await validateContract("runtime-manifest", manifest), []);
  await access(join(prepared.worktree, ".git"));

  const worker = JSON.parse(await readFile(prepared.workerInput, "utf8"));
  for (const path of ["src/App.tsx", "src/App.test.tsx"]) {
    const file = join(prepared.worktree, path);
    await writeFile(file, `${await readFile(file, "utf8")}\n// Approved reset-flow evidence.\n`, "utf8");
  }
  await git(prepared.worktree, ["add", "src/App.tsx", "src/App.test.tsx"]);
  await git(prepared.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", "feat: implement approved reset flow"]);
  const commit = await git(prepared.worktree, ["rev-parse", "HEAD"]);
  await mkdir(dirname(worker.result_path), { recursive: true });
  await writeFile(worker.result_path, `${JSON.stringify({
    contract_version: 1, work_id: prepared.workId, run_id: prepared.runId, repository: "frontend", status: "completed",
    summary: "Implemented the selected plan work.", branch: prepared.branch, worktree: prepared.worktree, commits: [commit],
    changed_files: ["src/App.test.tsx", "src/App.tsx"], checks: [{ command: "npm test", status: "passed", evidence: "Focused reset tests pass." }], risks: [],
  }, null, 2)}\n`, "utf8");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  const verifier = JSON.parse(await readFile(prepared.verifierInput, "utf8"));
  await writeFile(verifier.result_path, `${JSON.stringify({
    contract_version: 1, work_id: prepared.workId, run_id: prepared.runId, repository: "frontend", status: "pass",
    summary: "Verified the selected plan work.", acceptance: [{ criterion: "Reset returns the count to zero.", status: "passed", evidence: "The focused reset behavior was independently verified." }],
    checks: ["npm test"], findings: [], verified_at: "2026-08-11T08:35:00Z",
  }, null, 2)}\n`, "utf8");
  const untamperedManifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const tamperedManifest = structuredClone(untamperedManifest);
  const oversizedManifest = structuredClone(untamperedManifest);
  oversizedManifest.plan_work_items.push({ ...oversizedManifest.plan_work_items[0], work_id: "RESET-999" });
  assert.notDeepEqual(await validateContract("runtime-manifest", oversizedManifest), []);
  tamperedManifest.plan_work_items[0].work_id = "RESET-999";
  assert.deepEqual(await validateContract("runtime-manifest", tamperedManifest), []);
  await writeFile(prepared.manifest, `${JSON.stringify(tamperedManifest, null, 2)}\n`, "utf8");
  await assert.rejects(
    recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" }),
    /plan work item work_id mismatch/,
  );
  assert.equal(JSON.parse(await readFile(prepared.manifest, "utf8")).plan_work_items[0].outcome, "pending");
  await writeFile(prepared.manifest, `${JSON.stringify(untamperedManifest, null, 2)}\n`, "utf8");
  const passed = await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  assert.deepEqual(passed.plan_work_items?.map((item) => [item.work_id, item.outcome]), [["RESET-001", "passed"]]);
});

test("plan execution rejects caller-defined scope, acceptance, and dependency evidence", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, planDraft, new Date("2026-08-11T07:00:00Z"));
  const approved = await setPlanState(created.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-11T07:30:00Z"));
  const unsafe = {
    ...planRequest(approved.plan_version, approved.approved_digest!),
    request: "Caller override",
    dependency_evidence: [{ work_id: "RESET-001", evidence: "trust me" }],
    acceptance_criteria: ["Unrelated acceptance"],
    repositories: [{ name: "frontend", depends_on: [], scope: ["secrets.txt"], test_scope: [], test_policy: "not-required", verification_commands: [], acceptance_criteria: ["Anything"] }],
  };
  await assert.rejects(preparePlanTask({ workspaceRoot: workspace.root, request: unsafe as PlanRunTaskRequest }), /Invalid run-task-request/);
  await assert.rejects(access(join(workspace.root, ".runtime")));
});

test("task brief contract rejects inconsistent plan and direct source metadata", async () => {
  const base = {
    contract_version: 1, work_id: "RESET-001", run_id: "20260811T083000Z-abcd1234",
    requested_outcome: "Reset", scope: ["src/App.tsx"], acceptance_criteria: ["Reset works"], repositories: [{ name: "frontend", dependency_order: 0 }],
    activity: { reference: null, claim_status: "not-applicable", duplicate_effort_warning: true }, assumptions: [], risks: [], verification_commands: [],
    authorization: { kind: "confirmed-selection", evidence: "Selected." }, created_at: "2026-08-11T08:30:00Z",
  };
  assert.notDeepEqual(await validateContract("task-brief", { ...base, source: { kind: "plan", reference: "context/plans/reset-flow" }, plan: { reference: null, approval_state: "not-applicable" } }), []);
  assert.notDeepEqual(await validateContract("task-brief", { ...base, work_id: "ADHOC-20260811-001", source: { kind: "direct-request" }, plan: { reference: "context/plans/reset-flow", approval_state: "approved", plan_version: 1, approved_digest: `sha256:${"0".repeat(64)}`, work_ids: ["RESET-001"] } }), []);
  assert.notDeepEqual(await validateContract("task-brief", { ...base, work_id: "ADHOC-20260811-001", source: { kind: "direct-request" }, plan: { reference: null, approval_state: "not-applicable", plan_version: 1 } }), []);
  for (const kind of ["issue", "pull-request", "activity-task"] as const) {
    assert.deepEqual(await validateContract("task-brief", {
      ...base,
      source: { kind, reference: `${kind}:123` },
      plan: { reference: "context/plans/reset-flow", approval_state: "unknown" },
    }), [], kind);
  }
});

test("draft, stale, unknown, and dependency-blocked plans fail before runtime creation", async (t) => {
  await t.test("draft", async () => {
    const workspace = await createTestWorkspace();
    try {
      await createPlanDraft(workspace.root, planDraft, new Date("2026-08-11T07:00:00Z"));
      await assert.rejects(preparePlanTask({ workspaceRoot: workspace.root, request: planRequest(1, `sha256:${"0".repeat(64)}`) }), /is draft/);
      await assert.rejects(access(join(workspace.root, ".runtime")));
    } finally { await workspace.cleanup(); }
  });
  for (const scenario of ["stale", "unknown", "blocked"] as const) await t.test(scenario, async () => {
    const workspace = await createTestWorkspace();
    try {
      const created = await createPlanDraft(workspace.root, planDraft, new Date("2026-08-11T07:00:00Z"));
      const approved = await setPlanState(created.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-11T07:30:00Z"));
      const request = planRequest(approved.plan_version, approved.approved_digest!);
      if (scenario === "stale") request.source.approved_digest = `sha256:${"0".repeat(64)}`;
      if (scenario === "unknown") request.work_ids = ["RESET-999"];
      if (scenario === "blocked") request.work_ids = ["RESET-010"];
      await assert.rejects(preparePlanTask({ workspaceRoot: workspace.root, request }), scenario === "stale" ? /digest is stale/ : scenario === "unknown" ? /Unknown plan work ID/ : /dependency-blocked/);
      await assert.rejects(access(join(workspace.root, ".runtime")));
    } finally { await workspace.cleanup(); }
  });
});
