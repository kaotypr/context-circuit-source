import assert from "node:assert/strict";
import { mkdir, readFile, stat, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import test from "node:test";
import { git } from "../scripts/lib/git.js";
import { recordResult } from "../scripts/lib/record-result.js";
import { preparePlanlessTask, type PreparedTask } from "../scripts/lib/run-task.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace, taskOptions } from "./helpers.js";

async function writeJson(path: string, value: unknown): Promise<void> {
  await mkdir(dirname(path), { recursive: true });
  await writeFile(path, `${JSON.stringify(value, null, 2)}\n`, "utf8");
}

async function implementFixture(prepared: PreparedTask, includeTest = false): Promise<{ commit: string; resultPath: string }> {
  const appPath = join(prepared.worktree, "src", "App.tsx");
  const source = await readFile(appPath, "utf8");
  await writeFile(appPath, source.replace("Increment\n", "Reset\n"), "utf8");
  const changedFiles = ["src/App.tsx"];
  if (includeTest) {
    const testPath = join(prepared.worktree, "src", "App.test.tsx");
    const testSource = await readFile(testPath, "utf8");
    await writeFile(testPath, `${testSource}\n// Required test-scope lifecycle fixture.\n`, "utf8");
    changedFiles.push("src/App.test.tsx");
  }
  await git(prepared.worktree, ["add", ...changedFiles]);
  await git(prepared.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", "test: implement task"]);
  const commit = await git(prepared.worktree, ["rev-parse", "HEAD"]);
  const workerInput = JSON.parse(await readFile(prepared.workerInput, "utf8"));
  await writeJson(workerInput.result_path, {
    contract_version: 1,
    work_id: prepared.workId,
    run_id: prepared.runId,
    repository: "frontend",
    status: "completed",
    summary: "Implemented the scoped fixture change.",
    branch: prepared.branch,
    worktree: prepared.worktree,
    commits: [commit],
    changed_files: changedFiles,
    checks: [{ command: "npm test", status: "not-run", evidence: "Lifecycle test fixture" }],
    risks: [],
  });
  return { commit, resultPath: workerInput.result_path };
}

async function writePassingVerifier(prepared: PreparedTask): Promise<string> {
  const verifierInput = JSON.parse(await readFile(prepared.verifierInput, "utf8"));
  await writeJson(verifierInput.result_path, {
    contract_version: 1,
    work_id: prepared.workId,
    run_id: prepared.runId,
    repository: "frontend",
    status: "pass",
    summary: "The requested behavior is verified.",
    acceptance: [{ criterion: taskOptions.acceptanceCriteria[0], status: "passed", evidence: "Inspected the scoped change." }],
    checks: ["git diff --check"],
    findings: [],
    verified_at: "2026-08-11T12:03:00.000Z",
  });
  return verifierInput.result_path;
}

test("records an idempotent prepared-to-running-to-verifying-to-passed lifecycle", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T12:00:00.000Z"),
    discriminator: "aabbccdd",
  });

  let manifest = await recordResult({
    workspaceRoot: workspace.root,
    runId: prepared.runId,
    repository: "frontend",
    stage: "worker-started",
    now: new Date("2026-08-11T12:01:00.000Z"),
  });
  assert.equal(manifest.status, "running");
  manifest = await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-started" });
  assert.equal(manifest.execution_events?.length, 1);

  const worker = await implementFixture(prepared);
  manifest = await recordResult({
    workspaceRoot: workspace.root,
    runId: prepared.runId,
    repository: "frontend",
    stage: "worker-result",
    now: new Date("2026-08-11T12:02:00.000Z"),
  });
  assert.equal(manifest.status, "verifying");
  assert.ok(manifest.evidence.includes(worker.resultPath));
  assert.equal((await stat(worker.resultPath)).mode & 0o777, 0o600);

  const verifierPath = await writePassingVerifier(prepared);
  manifest = await recordResult({
    workspaceRoot: workspace.root,
    runId: prepared.runId,
    repository: "frontend",
    stage: "verifier-result",
    now: new Date("2026-08-11T12:03:00.000Z"),
  });
  assert.equal(manifest.status, "passed");
  assert.ok(manifest.evidence.includes(verifierPath));
  assert.equal((await stat(verifierPath)).mode & 0o777, 0o600);
  assert.deepEqual(manifest.execution_events?.map((event) => event.stage), ["worker-started", "worker-result", "verifier-result"]);
  assert.deepEqual(await validateContract("runtime-manifest", manifest), []);

  const replayed = await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  assert.equal(replayed.status, "passed");
  assert.equal(replayed.execution_events?.length, 3);
});

test("rejects an identity mismatch without advancing the manifest", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T13:00:00.000Z"),
    discriminator: "11223344",
  });
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-started" });
  const worker = await implementFixture(prepared);
  const result = JSON.parse(await readFile(worker.resultPath, "utf8"));
  result.run_id = "20260811T130000Z-deadbeef";
  await writeJson(worker.resultPath, result);

  await assert.rejects(
    recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" }),
    /worker result run_id mismatch/,
  );
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.status, "running");
  assert.equal(manifest.execution_events.length, 1);
});

test("infers a missing start event when recording an already completed worker", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T14:00:00.000Z"),
    discriminator: "55667788",
  });
  await implementFixture(prepared);
  const manifest = await recordResult({
    workspaceRoot: workspace.root,
    runId: prepared.runId,
    repository: "frontend",
    stage: "worker-result",
  });
  assert.equal(manifest.status, "verifying");
  assert.equal(manifest.execution_events?.[0]?.stage, "worker-started");
  assert.equal(manifest.execution_events?.[0]?.inferred, true);
});

test("enforces required test scope before accepting a completed worker", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    testScope: ["src/App.test.tsx"],
    now: new Date("2026-08-11T15:00:00.000Z"),
    discriminator: "99aabbcc",
  });
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-started" });
  await implementFixture(prepared);
  await assert.rejects(
    recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" }),
    /Required test policy needs a changed file in test scope/,
  );
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.status, "running");
});

test("accepts a completed worker that changes required implementation and test scopes", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    testScope: ["src/App.test.tsx"],
    now: new Date("2026-08-11T16:00:00.000Z"),
    discriminator: "ddeeff00",
  });
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-started" });
  await implementFixture(prepared, true);
  const manifest = await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  assert.equal(manifest.status, "verifying");
});
