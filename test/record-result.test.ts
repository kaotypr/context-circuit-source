import assert from "node:assert/strict";
import { mkdir, readFile, stat, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import test from "node:test";
import { git } from "../scripts/lib/git.js";
import { recordResult } from "../scripts/lib/record-result.js";
import { prepareExecutePlan } from "../scripts/lib/execute-plan.js";
import { preparePlanlessTask, type PreparedTask } from "../scripts/lib/run-task.js";
import { generatePlanBatch, setPlanState } from "../scripts/lib/plans.js";
import { validateContract } from "../scripts/lib/validation.js";
import type { PlanGenerationRequest } from "../scripts/lib/types.js";
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

test("executes a cumulative plan in dependency order and gates review on holistic verification", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const request: PlanGenerationRequest = {
    contract_version: 2,
    source: { kind: "prd", reference: "docs/runtime.md" },
    plans: [{
      plan_id: "cumulative-runtime", title: "Cumulative runtime", repository: "frontend", work_prefix: "CUM",
      summary: "Run a dependency-ordered cumulative plan.", assumptions: [], open_questions: [], requirements: ["Tasks run in order."], solution: ["Use one worktree."], delivery: ["Workers commit each task."], verification: ["The holistic verifier checks the cumulative worktree."], risks: [],
      work_items: [
        { key: "one", title: "First cumulative change", area: "runtime", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["First change is present."] },
        { key: "two", title: "Second cumulative change", area: "runtime", repository: "frontend", scope: ["src/App.test.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Second change is present."], depends_on: ["one"] },
        { key: "three", title: "Third cumulative change", area: "runtime", repository: "frontend", scope: ["package.json"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Third change is present."], depends_on: ["two"] },
      ],
    }],
  };
  const generated = await generatePlanBatch(workspace.root, request, new Date("2026-08-14T10:00:00Z"));
  const approved = await setPlanState(generated.plans[0]!.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-14T10:01:00Z"));
  const prepared = await prepareExecutePlan({
    workspaceRoot: workspace.root,
    request: { contract_version: 1, source: { kind: "plan", reference: approved.plan_reference!, plan_version: approved.plan_version, approved_digest: approved.approved_digest! } },
    now: new Date("2026-08-14T10:02:00Z"), discriminator: "cafebabe",
  });
  const readManifest = async () => JSON.parse(await readFile(prepared.manifest, "utf8")) as any;
  const complete = async (taskId: string, marker: string): Promise<void> => {
    let manifest = await readManifest();
    const task = manifest.task_graph.find((candidate: any) => candidate.task_id === taskId);
    await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId, stage: "worker-started" });
    const input = JSON.parse(await readFile(task.worker_input, "utf8"));
    const target = join(prepared.repositories[0]!.worktree, taskId === "CUM-001" ? "src/App.tsx" : taskId === "CUM-010" ? "src/App.test.tsx" : "package.json");
    await writeFile(target, `${await readFile(target, "utf8")}\n// ${marker}\n`, "utf8");
    const relativeTarget = taskId === "CUM-001" ? "src/App.tsx" : taskId === "CUM-010" ? "src/App.test.tsx" : "package.json";
    await git(prepared.repositories[0]!.worktree, ["add", relativeTarget]);
    await git(prepared.repositories[0]!.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", `test: ${marker}`]);
    const head = await git(prepared.repositories[0]!.worktree, ["rev-parse", "HEAD"]);
    const commits = (await git(prepared.repositories[0]!.worktree, ["rev-list", "--reverse", `${input.start_commit}..${head}`])).split("\n").filter(Boolean);
    const changedFiles = (await git(prepared.repositories[0]!.worktree, ["diff", "--name-only", `${input.start_commit}...${head}`])).split("\n").filter(Boolean);
    await writeJson(input.result_path, { contract_version: 2, plan_reference: input.plan_reference, plan_id: input.plan_id, plan_version: input.plan_version, approved_digest: input.approved_digest, task_id: taskId, repository: "frontend", plan_revision: input.plan_revision, attempt: input.attempt, run_id: prepared.runId, status: "completed", summary: marker, branch: input.branch, worktree: input.worktree, start_commit: input.start_commit, commits, changed_files: changedFiles, checks: [{ command: "git diff --check", status: "passed", evidence: "clean" }], risks: [] });
    manifest = await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId, stage: "worker-result" });
    const verifier = JSON.parse(await readFile(task.verifier_input, "utf8"));
    await writeJson(verifier.result_path, { contract_version: 2, plan_reference: verifier.plan_reference, plan_id: verifier.plan_id, plan_version: verifier.plan_version, approved_digest: verifier.approved_digest, task_id: taskId, repository: "frontend", plan_revision: verifier.plan_revision, attempt: verifier.attempt, run_id: prepared.runId, status: "pass", summary: "passed", branch: verifier.branch, worktree: verifier.worktree, start_commit: verifier.start_commit, acceptance: [{ criterion: verifier.acceptance_criteria[0], status: "passed", evidence: "verified" }], checks: ["git diff --check"], findings: [], verified_at: "2026-08-14T10:03:00Z" });
    await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId, stage: "verifier-result" });
  };
  await complete("CUM-001", "one");
  let manifest = await readManifest();
  assert.equal(manifest.task_graph.find((task: any) => task.task_id === "CUM-010").ready, true);
  assert.equal(manifest.task_graph.find((task: any) => task.task_id === "CUM-020").ready, false);
  await complete("CUM-010", "two");
  await complete("CUM-020", "three");
  manifest = await readManifest();
  assert.equal(manifest.status, "verifying");
  assert.equal(manifest.plan_verifier_status, "pending");
  const finalInput = JSON.parse(await readFile(manifest.plan_verifier_input, "utf8"));
  const finalHead = await git(prepared.repositories[0]!.worktree, ["rev-parse", "HEAD"]);
  const tasks = manifest.task_graph.map((task: any) => ({ task_id: task.task_id, repository: task.repository, status: "passed", head_commit: finalHead, evidence: "all task verifiers passed" }));
  await writeJson(finalInput.result_path, { contract_version: 1, plan_reference: finalInput.plan_reference, plan_id: finalInput.plan_id, plan_version: finalInput.plan_version, approved_digest: finalInput.approved_digest, run_id: prepared.runId, plan_revision: finalInput.plan_revision, task_id: finalInput.task_id, repository: "plan", attempt: 0, status: "pass", summary: "holistic pass", tasks, acceptance: finalInput.acceptance_criteria.map((criterion: string) => ({ criterion, status: "passed", evidence: "cumulative verification" })), checks: ["cumulative verification"], findings: [], verified_at: "2026-08-14T10:04:00Z" });
  manifest = await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, stage: "plan-verifier-result" });
  assert.equal(manifest.status, "passed");
  assert.equal(manifest.plan_verifier_status, "passed");
});

test("serializes ready tasks in one repository until the current task verifier passes", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const request: PlanGenerationRequest = {
    contract_version: 2,
    source: { kind: "prd", reference: "docs/serialization.md" },
    plans: [{
      plan_id: "serialization", title: "Serialized task runtime", repository: "frontend", work_prefix: "SER",
      summary: "Serialize tasks sharing one cumulative worktree.", assumptions: [], open_questions: [], requirements: ["One task owns the repository at a time."], solution: ["Keep a runtime repository lock."], delivery: ["Pass the current task before unlocking the next."], verification: ["The next task cannot start during worker or verifier execution."], risks: [],
      work_items: [
        { key: "first", title: "First serialized task", area: "runtime", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["First task passes."] },
        { key: "second", title: "Second serialized task", area: "runtime", repository: "frontend", scope: ["src/App.test.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Second task passes."] },
      ],
    }],
  };
  const generated = await generatePlanBatch(workspace.root, request, new Date("2026-08-14T11:00:00Z"));
  const approved = await setPlanState(generated.plans[0]!.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-14T11:01:00Z"));
  const prepared = await prepareExecutePlan({ workspaceRoot: workspace.root, request: { contract_version: 1, source: { kind: "plan", reference: approved.plan_reference!, plan_version: approved.plan_version, approved_digest: approved.approved_digest! } }, now: new Date("2026-08-14T11:02:00Z"), discriminator: "1122aabb" });
  const readManifest = async () => JSON.parse(await readFile(prepared.manifest, "utf8")) as any;
  let manifest = await readManifest();
  const first = manifest.task_graph[0];
  const second = manifest.task_graph[1];
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId: first.task_id, stage: "worker-started" });
  assert.equal((await readManifest()).repositories[0].active_task_id, first.task_id);
  await assert.rejects(recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId: second.task_id, stage: "worker-started" }), /not ready/);
  await writeFile(join(prepared.repositories[0]!.worktree, "src/App.tsx"), `${await readFile(join(prepared.repositories[0]!.worktree, "src/App.tsx"), "utf8")}\n// serialized first\n`, "utf8");
  await git(prepared.repositories[0]!.worktree, ["add", "src/App.tsx"]);
  await git(prepared.repositories[0]!.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", "test: serialized first"]);
  const firstInput = JSON.parse(await readFile(first.worker_input, "utf8"));
  const firstHead = await git(prepared.repositories[0]!.worktree, ["rev-parse", "HEAD"]);
  await writeJson(firstInput.result_path, { contract_version: 2, plan_reference: firstInput.plan_reference, plan_id: firstInput.plan_id, plan_version: firstInput.plan_version, approved_digest: firstInput.approved_digest, task_id: first.task_id, repository: "frontend", plan_revision: firstInput.plan_revision, attempt: 0, run_id: prepared.runId, status: "completed", summary: "first", branch: firstInput.branch, worktree: firstInput.worktree, start_commit: firstInput.start_commit, commits: [firstHead], changed_files: ["src/App.tsx"], checks: [], risks: [] });
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId: first.task_id, stage: "worker-result" });
  await assert.rejects(recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId: second.task_id, stage: "worker-started" }), /not ready/);
  const firstVerifier = JSON.parse(await readFile(first.verifier_input, "utf8"));
  await writeJson(firstVerifier.result_path, { contract_version: 2, plan_reference: firstVerifier.plan_reference, plan_id: firstVerifier.plan_id, plan_version: firstVerifier.plan_version, approved_digest: firstVerifier.approved_digest, task_id: first.task_id, repository: "frontend", plan_revision: firstVerifier.plan_revision, attempt: 0, run_id: prepared.runId, status: "pass", summary: "first passed", branch: firstVerifier.branch, worktree: firstVerifier.worktree, start_commit: firstVerifier.start_commit, acceptance: [{ criterion: firstVerifier.acceptance_criteria[0], status: "passed", evidence: "verified" }], checks: [], findings: [], verified_at: "2026-08-14T11:04:00Z" });
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId: first.task_id, stage: "verifier-result" });
  manifest = await readManifest();
  assert.equal(manifest.repositories[0].active_task_id, undefined);
  assert.equal(manifest.task_graph[1].ready, true);
  assert.equal(manifest.task_graph[1].start_commit, firstHead);
});

test("rejects a plan worker result whose immutable plan identity is forged", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const request: PlanGenerationRequest = {
    contract_version: 2, source: { kind: "prd", reference: "docs/identity.md" }, plans: [{
      plan_id: "identity", title: "Identity runtime", repository: "frontend", work_prefix: "IDN", summary: "Pin plan task evidence.", assumptions: [], open_questions: [], requirements: ["Identity is immutable."], solution: ["Validate every result field."], delivery: ["Reject forged results."], verification: ["A wrong plan ID cannot advance the manifest."], risks: [],
      work_items: [{ key: "only", title: "Identity task", area: "runtime", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Identity is checked."] }],
    }],
  };
  const generated = await generatePlanBatch(workspace.root, request, new Date("2026-08-14T12:00:00Z"));
  const approved = await setPlanState(generated.plans[0]!.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-14T12:01:00Z"));
  const prepared = await prepareExecutePlan({ workspaceRoot: workspace.root, request: { contract_version: 1, source: { kind: "plan", reference: approved.plan_reference!, plan_version: approved.plan_version, approved_digest: approved.approved_digest! } }, now: new Date("2026-08-14T12:02:00Z"), discriminator: "3344ccdd" });
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const task = manifest.task_graph[0];
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId: task.task_id, stage: "worker-started" });
  await writeFile(join(prepared.repositories[0]!.worktree, "src/App.tsx"), `${await readFile(join(prepared.repositories[0]!.worktree, "src/App.tsx"), "utf8")}\n// identity\n`, "utf8");
  await git(prepared.repositories[0]!.worktree, ["add", "src/App.tsx"]);
  await git(prepared.repositories[0]!.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", "test: identity"]);
  const input = JSON.parse(await readFile(task.worker_input, "utf8"));
  const head = await git(prepared.repositories[0]!.worktree, ["rev-parse", "HEAD"]);
  await writeJson(input.result_path, { contract_version: 2, plan_reference: input.plan_reference, plan_id: "forged", plan_version: input.plan_version, approved_digest: input.approved_digest, task_id: task.task_id, repository: "frontend", plan_revision: input.plan_revision, attempt: 0, run_id: prepared.runId, status: "completed", summary: "forged", branch: input.branch, worktree: input.worktree, start_commit: input.start_commit, commits: [head], changed_files: ["src/App.tsx"], checks: [], risks: [] });
  await assert.rejects(recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", taskId: task.task_id, stage: "worker-result" }), /plan worker result plan_id mismatch/);
  assert.equal(JSON.parse(await readFile(prepared.manifest, "utf8")).status, "running");
});
