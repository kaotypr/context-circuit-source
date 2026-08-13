import assert from "node:assert/strict";
import { access, mkdir, readFile, writeFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { prepareExecutePlan } from "../scripts/lib/execute-plan.js";
import { generateRunId } from "../scripts/lib/ids.js";
import { git } from "../scripts/lib/git.js";
import { generatePlanBatch, setPlanState } from "../scripts/lib/plans.js";
import type { PlanGenerationRequest } from "../scripts/lib/types.js";
import { createTestWorkspace } from "./helpers.js";

function request(): PlanGenerationRequest {
  return {
    contract_version: 2,
    source: { kind: "prd", reference: "docs/execute.md" },
    plans: [{
      plan_id: "complete-runtime", title: "Complete runtime", repository: "frontend", work_prefix: "RUN",
      summary: "Prepare a complete approved plan.", assumptions: [], open_questions: [], requirements: ["All tasks are selected."], solution: ["Use one cumulative repository worktree."], delivery: ["Prepare before starting workers."], verification: ["Verify the prepared graph."], risks: [],
      work_items: [
        { key: "first", title: "First task", area: "foundation", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["The first task is present."] },
        { key: "second", title: "Second task", area: "foundation", repository: "frontend", scope: ["src/App.test.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["The second task is present."], depends_on: ["first"] },
      ],
    }],
  };
}

test("execute-plan prepares the complete approved graph on one cumulative repository worktree", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const generated = await generatePlanBatch(workspace.root, request(), new Date("2026-08-14T09:00:00Z"));
  const plan = generated.plans[0]!;
  const approved = await setPlanState(plan.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-14T09:01:00Z"));
  const prepared = await prepareExecutePlan({
    workspaceRoot: workspace.root,
    request: { contract_version: 1, source: { kind: "plan", reference: approved.plan_reference!, plan_version: approved.plan_version, approved_digest: approved.approved_digest! } },
    now: new Date("2026-08-14T09:02:00Z"),
    discriminator: "12345678",
  });
  assert.equal(prepared.repositories.length, 1);
  assert.equal(prepared.repositories[0]!.taskInputs.length, 2);
  assert.equal(prepared.repositories[0]!.taskInputs[0]!.includes("RUN-001-worker-input.json"), true);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.contract_version, 2);
  assert.equal(manifest.plan_work_items.length, 2);
  assert.deepEqual(manifest.plan_work_items.map((item: { outcome: string }) => item.outcome), ["pending", "pending"]);
  assert.equal(manifest.repositories[0].branch, prepared.repositories[0]!.branch);
  await access(join(prepared.repositories[0]!.worktree, ".git"));
});

test("execute-plan rejects draft or stale plans before runtime or worktree creation", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const generated = await generatePlanBatch(workspace.root, request());
  const plan = generated.plans[0]!;
  const draftRequest = { contract_version: 1 as const, source: { kind: "plan" as const, reference: plan.plan_reference, plan_version: 1, approved_digest: `sha256:${"0".repeat(64)}` } };
  await assert.rejects(prepareExecutePlan({ workspaceRoot: workspace.root, request: draftRequest }), /explicit approval/);
  await assert.rejects(access(join(workspace.root, ".runtime")));
  const approved = await setPlanState(plan.directory, { kind: "approve", approved_by: "owner" });
  await assert.rejects(prepareExecutePlan({ workspaceRoot: workspace.root, request: { ...draftRequest, source: { ...draftRequest.source, approved_digest: `sha256:${"0".repeat(64)}` } } }), /stale/);
  assert.equal(approved.status, "approved");
  await assert.rejects(access(join(workspace.root, ".runtime")));
});

test("execute-plan refuses dirty bases without mutating the plan or runtime", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const generated = await generatePlanBatch(workspace.root, request());
  const plan = generated.plans[0]!;
  const approved = await setPlanState(plan.directory, { kind: "approve", approved_by: "owner" });
  await writeFile(join(workspace.repository, "unrecorded.txt"), "keep me\n", "utf8");
  await assert.rejects(prepareExecutePlan({ workspaceRoot: workspace.root, request: { contract_version: 1, source: { kind: "plan", reference: approved.plan_reference!, plan_version: 1, approved_digest: approved.approved_digest! } } }), /unresolved local changes/);
  assert.equal(await readFile(join(workspace.repository, "unrecorded.txt"), "utf8"), "keep me\n");
  await assert.rejects(access(join(workspace.root, ".runtime")));
});

test("execute-plan rolls back partial runtime output after a pre-worktree failure", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const generated = await generatePlanBatch(workspace.root, request());
  const plan = generated.plans[0]!;
  const approved = await setPlanState(plan.directory, { kind: "approve", approved_by: "owner" });
  const now = new Date("2026-08-14T09:02:00Z");
  const discriminator = "87654321";
  const runId = generateRunId(`execute-plan:${approved.plan_reference!}`, now, discriminator);
  const branch = `plan/complete-runtime-${runId.slice(-8)}`;
  const branchLock = join(workspace.repository, ".git", "refs", "heads", "plan", `complete-runtime-${runId.slice(-8)}.lock`);
  await mkdir(join(workspace.repository, ".git", "refs", "heads", "plan"), { recursive: true });
  await writeFile(branchLock, "", "utf8");

  await assert.rejects(prepareExecutePlan({
    workspaceRoot: workspace.root,
    request: { contract_version: 1, source: { kind: "plan", reference: approved.plan_reference!, plan_version: approved.plan_version, approved_digest: approved.approved_digest! } },
    now,
    discriminator,
  }));

  await assert.rejects(access(join(workspace.root, ".runtime", "plans", `${runId}.json`)));
  await assert.rejects(access(join(workspace.root, ".runtime", "runs", runId)));
  await assert.rejects(access(join(workspace.root, ".runtime", "worktrees", runId)));
  await assert.rejects(git(workspace.repository, ["show-ref", "--verify", "--quiet", `refs/heads/${branch}`]));
  assert.equal(await readFile(branchLock, "utf8"), "");
});
