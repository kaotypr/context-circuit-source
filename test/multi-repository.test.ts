import assert from "node:assert/strict";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import test from "node:test";
import { git } from "../scripts/lib/git.js";
import { recordResult } from "../scripts/lib/record-result.js";
import { prepareContractFirstTask } from "../scripts/lib/run-task.js";
import type { RunTaskRequest } from "../scripts/lib/types.js";
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
