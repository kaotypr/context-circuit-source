import assert from "node:assert/strict";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import test from "node:test";
import { git } from "../scripts/lib/git.js";
import { recordResult } from "../scripts/lib/record-result.js";
import { prepareExecutePlan } from "../scripts/lib/execute-plan.js";
import { prepareContractFirstTask } from "../scripts/lib/run-task.js";
import { generatePlanBatch, setPlanState } from "../scripts/lib/plans.js";
import type { PlanGenerationRequest, RunTaskRequest } from "../scripts/lib/types.js";
import { createTwoRepositoryTestWorkspace } from "./helpers.js";

const request: RunTaskRequest = {
  contract_version: 1,
  request: "Add a shared reset-counter contract and consume it in the frontend",
  acceptance_criteria: ["The frontend and API agree that reset returns count zero"],
  shared_contract: { repository: "backend", paths: ["src/contract.ts"] },
  repositories: [
    {
      name: "backend",
      depends_on: [],
      scope: ["src/contract.ts", "src/api.ts"],
      test_scope: [],
      test_policy: "verifier-only",
      acceptance_criteria: ["The API contract defines a reset state with count zero"],
      verification_commands: ["npm test"],
    },
    {
      name: "frontend",
      depends_on: ["backend"],
      scope: ["src/App.tsx"],
      test_scope: ["src/App.test.tsx"],
      test_policy: "existing-coverage",
      acceptance_criteria: ["The frontend consumes the verified reset contract"],
      verification_commands: ["npm test"],
    },
  ],
};

async function writeJson(path: string, value: unknown): Promise<void> {
  await mkdir(dirname(path), { recursive: true });
  await writeFile(path, `${JSON.stringify(value, null, 2)}\n`, "utf8");
}

test("prepares separate worktrees and locks dependent repository inputs", async (t) => {
  const workspace = await createTwoRepositoryTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await prepareContractFirstTask({ workspaceRoot: workspace.root, request, now: new Date("2026-08-12T01:00:00Z"), discriminator: "1234abcd" });
  assert.equal(prepared.repositories.length, 2);
  const backend = prepared.repositories.find((repository) => repository.name === "backend")!;
  const frontend = prepared.repositories.find((repository) => repository.name === "frontend")!;
  assert.notEqual(backend.worktree, frontend.worktree);
  assert.equal(backend.ready, true);
  assert.equal(frontend.ready, false);
  const frontendInput = JSON.parse(await readFile(frontend.workerInput, "utf8"));
  assert.equal(frontendInput.ready, false);
  assert.deepEqual(frontendInput.blocked_by, ["backend"]);
  assert.deepEqual(frontendInput.allowed_scope, ["src/App.tsx"]);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.deepEqual(manifest.repositories.map((repository: { name: string; status: string }) => [repository.name, repository.status]), [["backend", "prepared"], ["frontend", "waiting"]]);
  await assert.rejects(recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-started" }), /blocked by: backend/);
});

test("independent contract verification unlocks the dependent worker", async (t) => {
  const workspace = await createTwoRepositoryTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await prepareContractFirstTask({ workspaceRoot: workspace.root, request, now: new Date("2026-08-12T02:00:00Z"), discriminator: "abcdef12" });
  const backend = prepared.repositories.find((repository) => repository.name === "backend")!;
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "backend", stage: "worker-started" });
  const contractPath = join(backend.worktree, "src", "contract.ts");
  await writeFile(contractPath, `${await readFile(contractPath, "utf8")}\n// Contract established for the cross-repository proof.\n`, "utf8");
  await git(backend.worktree, ["add", "src/contract.ts"]);
  await git(backend.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", "feat: establish reset contract"]);
  const commit = await git(backend.worktree, ["rev-parse", "HEAD"]);
  const workerInput = JSON.parse(await readFile(backend.workerInput, "utf8"));
  await writeJson(workerInput.result_path, {
    contract_version: 1, work_id: prepared.workId, run_id: prepared.runId, repository: "backend", status: "completed",
    summary: "Established the shared contract.", branch: backend.branch, worktree: backend.worktree, commits: [commit], changed_files: ["src/contract.ts"],
    checks: [{ command: "npm test", status: "not-run", evidence: "Verifier supplies independent evidence." }], risks: [],
  });
  let manifest = await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "backend", stage: "worker-result" });
  assert.equal(manifest.status, "verifying");
  const verifierInput = JSON.parse(await readFile(backend.verifierInput, "utf8"));
  await writeJson(verifierInput.result_path, {
    contract_version: 1, work_id: prepared.workId, run_id: prepared.runId, repository: "backend", status: "pass", summary: "Contract is stable.",
    acceptance: [{ criterion: request.repositories[0]!.acceptance_criteria[0], status: "passed", evidence: "Contract declares a zero count reset response." }],
    checks: ["npm test"], findings: [], verified_at: "2026-08-12T02:03:00Z",
  });
  manifest = await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "backend", stage: "verifier-result" });
  assert.equal(manifest.status, "prepared");
  assert.equal(manifest.repositories.find((repository) => repository.name === "backend")?.status, "passed");
  assert.equal(manifest.repositories.find((repository) => repository.name === "frontend")?.status, "prepared");
  const frontend = prepared.repositories.find((repository) => repository.name === "frontend")!;
  const unlocked = JSON.parse(await readFile(frontend.workerInput, "utf8"));
  assert.equal(unlocked.ready, true);
  assert.deepEqual(unlocked.blocked_by, []);
  assert.equal(unlocked.shared_contract.approval, "verified");
  const unlockedVerifier = JSON.parse(await readFile(frontend.verifierInput, "utf8"));
  assert.equal(unlockedVerifier.shared_contract.approval, "verified");
  manifest = await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-started" });
  assert.equal(manifest.status, "running");
});

test("rejects an ambiguous dependency graph before creating runtime state", async (t) => {
  const workspace = await createTwoRepositoryTestWorkspace();
  t.after(workspace.cleanup);
  const ambiguous: RunTaskRequest = { ...request, repositories: request.repositories.map((repository) => repository.name === "frontend" ? { ...repository, depends_on: [] } : repository) };
  await assert.rejects(prepareContractFirstTask({ workspaceRoot: workspace.root, request: ambiguous }), /must depend on the shared contract repository/);
});

test("protects every base repository before creating cross-repository runtime state", async (t) => {
  const workspace = await createTwoRepositoryTestWorkspace();
  t.after(workspace.cleanup);
  const unrecorded = join(workspace.frontend, "preserve-me.txt");
  await writeFile(unrecorded, "do not discard\n", "utf8");
  await assert.rejects(prepareContractFirstTask({ workspaceRoot: workspace.root, request }), /unresolved local changes/);
  assert.equal(await readFile(unrecorded, "utf8"), "do not discard\n");
});

test("executes cumulative tasks across repositories with independent worktrees and verified unlocking", async (t) => {
  const workspace = await createTwoRepositoryTestWorkspace();
  t.after(workspace.cleanup);
  const planRequest: PlanGenerationRequest = {
    contract_version: 2,
    source: { kind: "prd", reference: "docs/multi-repository-runtime.md" },
    plans: [{
      plan_id: "multi-runtime", title: "Multi-repository runtime", repository: "backend", affected_repositories: ["backend", "frontend"], work_prefix: "MUL",
      summary: "Run a backend task before the dependent frontend task.", assumptions: [], open_questions: [], requirements: ["Only verified backend work unlocks frontend work."], solution: ["Use one cumulative worktree per repository."], delivery: ["Verify backend before preparing frontend."], verification: ["The runtime records repository-scoped branches and worktrees."], risks: [],
      work_items: [
        { key: "backend", title: "Backend contract", area: "runtime", repository: "backend", scope: ["src/contract.ts"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Backend contract is updated."] },
        { key: "frontend", title: "Frontend consumer", area: "runtime", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Frontend consumes the verified contract."], depends_on: ["MUL-001"] },
      ],
    }],
  };
  const generated = await generatePlanBatch(workspace.root, planRequest, new Date("2026-08-14T13:00:00Z"));
  const approved = await setPlanState(generated.plans[0]!.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-14T13:01:00Z"));
  const prepared = await prepareExecutePlan({ workspaceRoot: workspace.root, request: { contract_version: 1, source: { kind: "plan", reference: approved.plan_reference!, plan_version: approved.plan_version, approved_digest: approved.approved_digest! } }, now: new Date("2026-08-14T13:02:00Z"), discriminator: "5566ddee" });
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.repositories.length, 2);
  assert.notEqual(manifest.repositories[0].worktree, manifest.repositories[1].worktree);
  assert.equal(manifest.task_graph.find((task: any) => task.task_id === "MUL-001").ready, true);
  assert.equal(manifest.task_graph.find((task: any) => task.task_id === "MUL-010").ready, false);
  const backend = prepared.repositories.find((repository) => repository.name === "backend")!;
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "backend", taskId: "MUL-001", stage: "worker-started" });
  await writeFile(join(backend.worktree, "src/contract.ts"), `${await readFile(join(backend.worktree, "src/contract.ts"), "utf8")}\n// verified backend contract\n`, "utf8");
  await git(backend.worktree, ["add", "src/contract.ts"]);
  await git(backend.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", "test: backend contract"]);
  const workerInput = JSON.parse(await readFile(join(workspace.root, ".runtime", "runs", prepared.runId, "MUL-001-worker-input.json"), "utf8"));
  const head = await git(backend.worktree, ["rev-parse", "HEAD"]);
  await writeJson(workerInput.result_path, { contract_version: 2, plan_reference: workerInput.plan_reference, plan_id: workerInput.plan_id, plan_version: workerInput.plan_version, approved_digest: workerInput.approved_digest, task_id: "MUL-001", repository: "backend", plan_revision: workerInput.plan_revision, attempt: 0, run_id: prepared.runId, status: "completed", summary: "backend", branch: workerInput.branch, worktree: workerInput.worktree, start_commit: workerInput.start_commit, commits: [head], changed_files: ["src/contract.ts"], checks: [], risks: [] });
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "backend", taskId: "MUL-001", stage: "worker-result" });
  const verifierInput = JSON.parse(await readFile(join(workspace.root, ".runtime", "runs", prepared.runId, "MUL-001-verifier-input.json"), "utf8"));
  await writeJson(verifierInput.result_path, { contract_version: 2, plan_reference: verifierInput.plan_reference, plan_id: verifierInput.plan_id, plan_version: verifierInput.plan_version, approved_digest: verifierInput.approved_digest, task_id: "MUL-001", repository: "backend", plan_revision: verifierInput.plan_revision, attempt: 0, run_id: prepared.runId, status: "pass", summary: "backend passed", branch: verifierInput.branch, worktree: verifierInput.worktree, start_commit: verifierInput.start_commit, acceptance: [{ criterion: verifierInput.acceptance_criteria[0], status: "passed", evidence: "verified" }], checks: [], findings: [], verified_at: "2026-08-14T13:04:00Z" });
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "backend", taskId: "MUL-001", stage: "verifier-result" });
  const unlocked = JSON.parse(await readFile(join(workspace.root, ".runtime", "runs", prepared.runId, "MUL-010-worker-input.json"), "utf8"));
  assert.equal(unlocked.ready, true);
  assert.equal(unlocked.start_commit, await git(workspace.frontend, ["rev-parse", "main"]));
});
