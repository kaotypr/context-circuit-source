import assert from "node:assert/strict";
import { access, mkdir, readFile, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import test from "node:test";
import { finishWork } from "../scripts/lib/finish-work.js";
import { git } from "../scripts/lib/git.js";
import { createPlanDraft, setPlanState } from "../scripts/lib/plans.js";
import { recordResult } from "../scripts/lib/record-result.js";
import { confirmMerge, prepareReview, recordReviewPublication } from "../scripts/lib/review-lifecycle.js";
import { preparePlanTask, type PreparedTask } from "../scripts/lib/run-task.js";
import type { PlanDraftRequest } from "../scripts/lib/types.js";
import { recommendWhatsNext } from "../scripts/lib/whats-next.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace } from "./helpers.js";

const plan: PlanDraftRequest = {
  contract_version: 1,
  plan_id: "journey",
  title: "Complete journey",
  source: { kind: "prd", reference: "docs/journey.md" },
  work_prefix: "JOURNEY",
  summary: "Exercise the complete plan-linked workflow.",
  affected_repositories: ["frontend"],
  assumptions: [], open_questions: [],
  requirements: ["The selected behavior is delivered."],
  solution: ["Make one scoped product change."],
  delivery: ["Review and merge through the human gate."],
  verification: ["Verify the acceptance criterion independently."],
  risks: [],
  work_items: [{
    key: "behavior", title: "Deliver selected behavior", area: "application behavior", repository: "frontend",
    scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: ["npm test"],
    acceptance_criteria: ["The selected behavior is delivered."],
  }],
};

async function writeJson(path: string, value: unknown): Promise<void> {
  await mkdir(dirname(path), { recursive: true });
  await writeFile(path, `${JSON.stringify(value, null, 2)}\n`, "utf8");
}

async function initializeWrapper(root: string): Promise<void> {
  await writeFile(join(root, ".gitignore"), ".runtime/\nrepositories/frontend/\n", "utf8");
  await git(root, ["init", "--initial-branch=main"]);
  await git(root, ["add", "."]);
  await git(root, ["-c", "user.name=Wrapper", "-c", "user.email=wrapper@example.invalid", "commit", "-m", "test: configure wrapper"]);
}

async function preparePassedPlanRun(root: string, discriminator: string): Promise<PreparedTask> {
  const created = await createPlanDraft(root, plan, new Date("2026-08-12T03:00:00Z"));
  const approved = await setPlanState(created.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-12T03:01:00Z"));
  await git(root, ["add", "context/plans/journey"]);
  await git(root, ["-c", "user.name=Owner", "-c", "user.email=owner@example.invalid", "commit", "-m", "plan: approve journey"]);
  const prepared = await preparePlanTask({
    workspaceRoot: root,
    request: { contract_version: 1, source: { kind: "plan", reference: "context/plans/journey", plan_version: approved.plan_version, approved_digest: approved.approved_digest! }, work_ids: ["JOURNEY-001"] },
    discriminator, now: new Date("2026-08-12T03:02:00Z"),
  });
  const app = join(prepared.worktree, "src", "App.tsx");
  await writeFile(app, `${await readFile(app, "utf8")}\n// end-to-end journey evidence\n`, "utf8");
  await git(prepared.worktree, ["add", "src/App.tsx"]);
  await git(prepared.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", "feat: deliver journey behavior"]);
  const head = await git(prepared.worktree, ["rev-parse", "HEAD"]);
  const worker = JSON.parse(await readFile(prepared.workerInput, "utf8"));
  await writeJson(worker.result_path, {
    contract_version: 1, work_id: prepared.workId, run_id: prepared.runId, repository: "frontend", status: "completed",
    summary: "Delivered the selected plan item.", branch: prepared.branch, worktree: prepared.worktree,
    commits: [head], changed_files: ["src/App.tsx"], checks: [{ command: "npm test", status: "passed", evidence: "Fixture suite passed." }], risks: [],
  });
  await recordResult({ workspaceRoot: root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  const verifier = JSON.parse(await readFile(prepared.verifierInput, "utf8"));
  await writeJson(verifier.result_path, {
    contract_version: 1, work_id: prepared.workId, run_id: prepared.runId, repository: "frontend", status: "pass",
    summary: "Independent end-to-end verification passed.",
    acceptance: [{ criterion: "The selected behavior is delivered.", status: "passed", evidence: "Committed diff inspected." }],
    checks: ["npm test passed"], findings: [], verified_at: "2026-08-12T03:04:00Z",
  });
  await recordResult({ workspaceRoot: root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  return prepared;
}

test("configured approved-plan journey reaches verified merge, safe closeout, and completed exclusion", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await initializeWrapper(workspace.root);
  const prepared = await preparePassedPlanRun(workspace.root, "e2e00001");
  const review = await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  assert.equal(review.status, "ready-for-local-review");
  await git(workspace.repository, ["-c", "user.name=Human", "-c", "user.email=human@example.invalid", "merge", "--no-ff", "-m", "merge: journey", prepared.branch]);
  const mergeCommit = await git(workspace.repository, ["rev-parse", "HEAD"]);
  const confirmation = await confirmMerge({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", mergeCommit, author: "owner", evidence: "Human verified and merged the reviewed head." });
  assert.equal(confirmation.status, "closeout-ready");
  const closeout = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "owner", mergeCommit, now: new Date("2026-08-12T03:10:00Z") });
  assert.equal(closeout.status, "prepared");
  const contribution = await readFile(join(workspace.root, closeout.contribution), "utf8");
  assert.match(contribution, /Task source: plan/);
  assert.match(contribution, /Plan: `context\/plans\/journey`/);
  await git(workspace.root, ["add", closeout.contribution]);
  await git(workspace.root, ["-c", "user.name=Owner", "-c", "user.email=owner@example.invalid", "commit", "-m", "docs: record journey outcome"]);
  const closed = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "owner", mergeCommit, cleanup: true });
  assert.equal(closed.status, "closed");
  await access(prepared.manifest);
  assert.equal(await git(workspace.repository, ["show-ref", "--verify", `refs/heads/${prepared.branch}`]).then(() => true), true);
  const next = await recommendWhatsNext(workspace.root, null, new Date("2026-08-12T03:20:00Z"));
  assert.equal(next.considered.excluded, 1);
  assert.equal(next.recommendation.candidate_id, null);
});

test("remote publication fixture records confirmed evidence without pushing", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await initializeWrapper(workspace.root);
  await git(workspace.repository, ["remote", "add", "origin", "https://example.invalid/product.git"]);
  const prepared = await preparePassedPlanRun(workspace.root, "e2e00002");
  const refsBefore = await git(workspace.repository, ["for-each-ref", "--format=%(refname)", "refs/remotes/"]);
  const review = await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  assert.equal(review.status, "ready-for-publication");
  const publication = await recordReviewPublication({
    workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", status: "published", tool: "manual",
    pullRequest: "https://example.invalid/product/pull/7", evidence: "Fixture confirmation; no network operation was performed.", authorized: true,
  });
  assert.deepEqual(await validateContract("review-publication-record", publication), []);
  assert.equal(await git(workspace.repository, ["for-each-ref", "--format=%(refname)", "refs/remotes/"]), refsBefore);
  assert.equal(await git(prepared.worktree, ["rev-parse", "HEAD"]), publication.head_commit);
});
