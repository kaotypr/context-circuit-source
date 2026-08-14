import { chmod, readFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { assertCleanRepository, git } from "./git.js";
import { assertInside, ensurePrivateDirectory, withExclusiveFile, writeJsonAtomic } from "./io.js";
import type { ExecutionEvent, RuntimeManifest, RuntimeRepository, TaskBrief, TestExpectation } from "./types.js";
import { validateContract } from "./validation.js";

export type RecordStage = "worker-started" | "worker-result" | "verifier-result" | "plan-verifier-result";

interface WorkerInput {
  task_brief: string;
  repository: string;
  worktree: string;
  branch: string;
  base_commit: string;
  plan_id?: string;
  plan_reference?: string;
  plan_version?: number;
  plan_revision?: number;
  approved_digest?: string;
  task_id?: string;
  run_id?: string;
  attempt?: number;
  start_commit?: string | null;
  ready?: boolean;
  blocked_by?: string[];
  allowed_scope: string[];
  implementation_scope?: string[];
  test_expectation?: TestExpectation;
  result_path: string;
}

interface VerifierInput {
  task_brief: string;
  repository: string;
  worktree: string;
  branch: string;
  base_commit: string;
  plan_id?: string;
  plan_reference?: string;
  plan_version?: number;
  plan_revision?: number;
  approved_digest?: string;
  task_id?: string;
  run_id?: string;
  attempt?: number;
  start_commit?: string | null;
  worker_result: string;
  acceptance_criteria: string[];
  test_expectation?: TestExpectation;
  result_path: string;
}

interface PlanTaskInput extends WorkerInput {
  run_id: string;
  plan_id: string;
  plan_reference: string;
  plan_version: number;
  plan_revision: number;
  approved_digest: string;
  task_id: string;
  attempt: number;
  start_commit: string | null;
  worker_result?: string;
  acceptance_criteria: string[];
}

interface PlanVerifierInput {
  contract_version: 1;
  role: "plan-verifier";
  plan_reference: string;
  plan_id: string;
  plan_version: number;
  approved_digest: string;
  plan_revision: number;
  run_id: string;
  task_id: string;
  repository: "plan";
  attempt: number;
  ready: boolean;
  result_path: string;
  worktrees: Array<{ name: string; worktree: string; branch: string; base_commit: string }>;
  tasks: Array<{ task_id: string; repository: string; worker_result: string; verifier_result: string }>;
  acceptance_criteria: string[];
}

interface WorkerResult {
  contract_version: number;
  work_id?: string;
  run_id?: string;
  plan_reference?: string;
  plan_id?: string;
  plan_version?: number;
  task_id?: string;
  plan_revision?: number;
  approved_digest?: string;
  attempt?: number;
  repository: string;
  status: "completed" | "blocked" | "failed";
  branch: string;
  worktree: string;
  start_commit?: string;
  commits: string[];
  changed_files: string[];
  checks: Array<{ status: "passed" | "failed" | "not-run" }>;
}

interface VerifierResult {
  contract_version: number;
  work_id?: string;
  run_id?: string;
  plan_reference?: string;
  plan_id?: string;
  plan_version?: number;
  task_id?: string;
  plan_revision?: number;
  approved_digest?: string;
  attempt?: number;
  repository: string;
  branch?: string;
  worktree?: string;
  start_commit?: string;
  status: "pass" | "fail" | "blocked";
  acceptance: Array<{ criterion: string; status: string; evidence: string }>;
}

export interface RecordResultOptions {
  workspaceRoot: string;
  runId: string;
  repository?: string;
  taskId?: string;
  stage: RecordStage;
  now?: Date;
}

function assertIdentifier(value: string, label: string, pattern: RegExp): void {
  if (!pattern.test(value)) throw new Error(`Invalid ${label}: ${value}`);
}

async function readJson<T>(path: string): Promise<T> {
  return JSON.parse(await readFile(path, "utf8")) as T;
}

async function assertValid(name: "runtime-manifest" | "task-brief" | "worker-result" | "verifier-result" | "plan-verifier-result", value: unknown): Promise<void> {
  const errors = await validateContract(name as Parameters<typeof validateContract>[0], value);
  if (errors.length > 0) {
    throw new Error(`Invalid ${name}: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  }
}

function assertEqual(actual: unknown, expected: unknown, label: string): void {
  if (actual !== expected) throw new Error(`${label} mismatch: expected ${String(expected)}, received ${String(actual)}`);
}

function inAllowedScope(path: string, scopes: string[]): boolean {
  return scopes.some((scope) => {
    const normalized = scope.replace(/\/$/, "");
    return path === normalized || path.startsWith(`${normalized}/`);
  });
}

function sameMembers(left: string[], right: string[]): boolean {
  return [...left].sort().join("\n") === [...right].sort().join("\n");
}

function sameTestExpectation(left: TestExpectation, right: TestExpectation): boolean {
  return left.policy === right.policy && left.rationale === right.rationale && sameMembers(left.paths, right.paths);
}

function findRepository(manifest: RuntimeManifest, name: string): RuntimeRepository {
  const repository = manifest.repositories.find((candidate) => candidate.name === name);
  if (!repository) throw new Error(`Run ${manifest.run_id} has no repository named ${name}`);
  return repository;
}

function assertPlanWorkItemAssociation(manifest: RuntimeManifest, brief: TaskBrief, repository: string): void {
  assertEqual(manifest.source_kind, brief.source.kind, "manifest source_kind");
  const planLinked = brief.source.kind === "plan";
  if (!planLinked) {
    if (manifest.plan_work_items !== undefined) throw new Error("Non-plan run must not contain plan work items");
    return;
  }
  if (brief.plan.approval_state !== "approved") throw new Error("Plan-linked task brief must contain approved plan metadata");
  if (!manifest.plan_work_items || manifest.plan_work_items.length !== 1) {
    throw new Error("Plan-linked run must contain exactly one plan work item");
  }
  const item = manifest.plan_work_items[0]!;
  assertEqual(item.work_id, manifest.work_id, "plan work item work_id");
  if (brief.plan.work_ids.length !== 1 || brief.plan.work_ids[0] !== item.work_id) {
    throw new Error("Plan work item identity does not match task brief work IDs");
  }
  if (manifest.repositories.length !== 1 || manifest.repositories[0]!.name !== item.repository ||
      brief.repositories.length !== 1 || brief.repositories[0]!.name !== item.repository || item.repository !== repository) {
    throw new Error("Plan work item repository does not match task brief and recorded repository");
  }
}

function eventKey(runId: string, repository: string, stage: RecordStage, attempt: number): string {
  const suffix = attempt === 0 ? "" : `:attempt-${attempt}`;
  return `${runId}:execution:${repository}:${stage}${suffix}`;
}

function appendEvent(
  manifest: RuntimeManifest,
  stage: RecordStage,
  repository: string,
  from: ExecutionEvent["from_status"],
  to: ExecutionEvent["to_status"],
  occurredAt: string,
  inferred: boolean,
  attempt: number,
  resultPath?: string,
): void {
  const event: ExecutionEvent = {
    stage,
    repository,
    from_status: from,
    to_status: to,
    inferred,
    attempt,
    idempotency_key: eventKey(manifest.run_id, repository, stage, attempt),
    occurred_at: occurredAt,
  };
  if (resultPath) event.result_path = resultPath;
  manifest.execution_events ??= [];
  manifest.execution_events.push(event);
  const runtimeRepository = findRepository(manifest, repository);
  runtimeRepository.status = to;
  refreshManifestStatus(manifest);
  manifest.updated_at = occurredAt;
  if (resultPath && !manifest.evidence.includes(resultPath)) manifest.evidence.push(resultPath);
}

function refreshManifestStatus(manifest: RuntimeManifest): void {
  const statuses = manifest.repositories.map((repository) => repository.status ?? manifest.status);
  if (statuses.every((status) => status === "passed")) manifest.status = "passed";
  else if (statuses.includes("failed")) manifest.status = "failed";
  else if (statuses.includes("blocked")) manifest.status = "blocked";
  else if (statuses.includes("verifying")) manifest.status = "verifying";
  else if (statuses.includes("running")) manifest.status = "running";
  else manifest.status = "prepared";
}

async function unlockDependents(runtimeRoot: string, manifest: RuntimeManifest): Promise<void> {
  for (const candidate of manifest.repositories) {
    if (candidate.status !== "waiting") continue;
    const dependencies = candidate.depends_on ?? [];
    if (!dependencies.every((name) => findRepository(manifest, name).status === "passed")) continue;
    for (const [path, worker] of [[candidate.worker_input, true], [candidate.verifier_input, false]] as const) {
      const inputPath = assertInside(runtimeRoot, path);
      const input = await readJson<Record<string, unknown>>(inputPath);
      if (worker) {
        input.ready = true;
        input.blocked_by = [];
      }
      if (input.shared_contract && typeof input.shared_contract === "object") {
        (input.shared_contract as Record<string, unknown>).approval = "verified";
      }
      await writeJsonAtomic(inputPath, input);
    }
    candidate.status = "prepared";
  }
  refreshManifestStatus(manifest);
}

async function assertWorktree(repository: RuntimeRepository): Promise<string> {
  await assertCleanRepository(repository.worktree);
  const branch = await git(repository.worktree, ["branch", "--show-current"]);
  assertEqual(branch, repository.branch, "worktree branch");
  return git(repository.worktree, ["rev-parse", "HEAD"]);
}

function assertTaskIdentity(manifest: RuntimeManifest, brief: TaskBrief, repository: string): void {
  assertEqual(brief.work_id, manifest.work_id, "task brief work_id");
  assertEqual(brief.run_id, manifest.run_id, "task brief run_id");
  if (!brief.repositories.some((candidate) => candidate.name === repository)) {
    throw new Error(`Task brief does not include repository ${repository}`);
  }
}

async function validateWorkerResult(
  manifest: RuntimeManifest,
  repository: RuntimeRepository,
  input: WorkerInput,
  testExpectation: TestExpectation,
): Promise<"verifying" | "failed" | "blocked"> {
  const result = await readJson<WorkerResult>(input.result_path);
  await assertValid("worker-result", result);
  assertEqual(result.work_id, manifest.work_id, "worker result work_id");
  assertEqual(result.run_id, manifest.run_id, "worker result run_id");
  assertEqual(result.repository, repository.name, "worker result repository");
  assertEqual(result.branch, repository.branch, "worker result branch");
  assertEqual(resolve(result.worktree), resolve(repository.worktree), "worker result worktree");

  if (result.status === "completed") {
    if (result.commits.length === 0) throw new Error("Completed worker result must record at least one commit");
    if (result.checks.some((check) => check.status === "failed")) throw new Error("Completed worker result cannot contain a failed check");
    const head = await assertWorktree(repository);
    assertEqual(result.commits.at(-1), head, "worker result final commit");
    const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${head}`])).split("\n").filter(Boolean);
    if (commits.join("\n") !== result.commits.join("\n")) throw new Error("worker result commits does not match the ordered base-to-head Git history");
    const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${head}`])).split("\n").filter(Boolean);
    if (!sameMembers(changedFiles, result.changed_files)) throw new Error("worker result changed_files does not match the base-to-head Git diff");
    const outsideScope = changedFiles.filter((path) => !inAllowedScope(path, input.allowed_scope));
    if (outsideScope.length > 0) throw new Error(`Worker changed files outside allowed scope: ${outsideScope.join(", ")}`);
    if (testExpectation.policy === "required" && !changedFiles.some((path) => inAllowedScope(path, testExpectation.paths))) {
      throw new Error(`Required test policy needs a changed file in test scope: ${testExpectation.paths.join(", ")}`);
    }
    return "verifying";
  }
  return result.status;
}

async function validateVerifierResult(
  manifest: RuntimeManifest,
  repository: RuntimeRepository,
  input: VerifierInput,
): Promise<"passed" | "failed" | "blocked"> {
  const worker = await readJson<WorkerResult>(input.worker_result);
  await assertValid("worker-result", worker);
  const result = await readJson<VerifierResult>(input.result_path);
  await assertValid("verifier-result", result);
  assertEqual(result.work_id, manifest.work_id, "verifier result work_id");
  assertEqual(result.run_id, manifest.run_id, "verifier result run_id");
  assertEqual(result.repository, repository.name, "verifier result repository");
  if (!sameMembers(result.acceptance.map((item) => item.criterion), input.acceptance_criteria)) {
    throw new Error("verifier result acceptance criteria do not match verifier input");
  }
  const acceptanceStatuses = result.acceptance.map((item) => item.status);
  if (result.status === "pass" && acceptanceStatuses.some((status) => status !== "passed")) {
    throw new Error("Passing verifier result requires every acceptance criterion to pass");
  }
  if (result.status === "fail" && !acceptanceStatuses.includes("failed")) {
    throw new Error("Failing verifier result must identify a failed acceptance criterion");
  }
  if (result.status === "blocked" && !acceptanceStatuses.includes("blocked")) {
    throw new Error("Blocked verifier result must identify a blocked acceptance criterion");
  }
  const head = await assertWorktree(repository);
  assertEqual(worker.commits.at(-1), head, "verified worker commit");
  return result.status === "pass" ? "passed" : result.status === "fail" ? "failed" : "blocked";
}

type PlanTaskRecord = NonNullable<RuntimeManifest["task_graph"]>[number];

function findPlanTask(manifest: RuntimeManifest, taskId: string): PlanTaskRecord {
  const task = manifest.task_graph?.find((candidate) => (candidate.task_id ?? candidate.work_id) === taskId);
  if (!task) throw new Error(`Run ${manifest.run_id} has no plan task named ${taskId}`);
  return task;
}

function planTaskAttempt(task: PlanTaskRecord): number {
  return task.attempt ?? 0;
}

function planEventKey(manifest: RuntimeManifest, taskId: string, repository: string, stage: RecordStage, attempt: number): string {
  return `${manifest.run_id}:execution:${repository}:${taskId}:${stage}:attempt-${attempt}`;
}

function setPlanTask(manifest: RuntimeManifest, task: PlanTaskRecord, status: NonNullable<PlanTaskRecord["status"]>, outcome?: NonNullable<PlanTaskRecord["outcome"]>): void {
  task.status = status;
  if (outcome) task.outcome = outcome;
  const summary = manifest.plan_work_items?.find((candidate) => candidate.work_id === task.work_id);
  if (summary) {
    summary.status = status;
    if (outcome) summary.outcome = outcome;
    if (task.attempt !== undefined) summary.attempt = task.attempt;
    if (task.start_commit !== undefined) summary.start_commit = task.start_commit;
    if (task.ready !== undefined) summary.ready = task.ready;
    if (task.blocked_by !== undefined) summary.blocked_by = task.blocked_by;
    if (task.worker_result !== undefined) summary.worker_result = task.worker_result;
    if (task.verifier_result !== undefined) summary.verifier_result = task.verifier_result;
  }
}

function refreshPlanRepositories(manifest: RuntimeManifest): void {
  const tasks = manifest.task_graph ?? [];
  for (const repository of manifest.repositories) {
    const repositoryTasks = tasks.filter((task) => task.repository === repository.name);
    const active = repositoryTasks.find((task) => task.status === "running" || task.status === "verifying");
    const failed = repositoryTasks.find((task) => task.outcome === "failed");
    const blocked = repositoryTasks.find((task) => task.outcome === "blocked");
    const prepared = repositoryTasks.find((task) => task.status === "prepared");
    const allPassed = repositoryTasks.length > 0 && repositoryTasks.every((task) => task.outcome === "passed");
    if (active?.task_id) repository.active_task_id = active.task_id;
    else if (repository.active_task_id) {
      const lockedTask = repositoryTasks.find((task) => (task.task_id ?? task.work_id) === repository.active_task_id);
      if (!lockedTask || (lockedTask.outcome !== "failed" && lockedTask.outcome !== "blocked")) delete repository.active_task_id;
    }
    if (active) repository.status = active.status === "verifying" ? "verifying" : "running";
    else if (blocked) repository.status = "blocked";
    else if (failed) repository.status = "failed";
    else if (allPassed) repository.status = "passed";
    else if (prepared) repository.status = "prepared";
    else repository.status = "waiting";
  }
}

function refreshPlanManifestStatus(manifest: RuntimeManifest): void {
  const tasks = manifest.task_graph ?? [];
  refreshPlanRepositories(manifest);
  if (tasks.some((task) => task.outcome === "blocked")) manifest.status = "blocked";
  else if (tasks.some((task) => task.outcome === "failed")) manifest.status = "failed";
  else if (manifest.plan_verifier_status === "blocked") manifest.status = "blocked";
  else if (manifest.plan_verifier_status === "failed") manifest.status = "failed";
  else if (tasks.length > 0 && tasks.every((task) => task.outcome === "passed")) {
    manifest.status = manifest.plan_verifier_status === "passed" ? "passed" : "verifying";
  } else if (tasks.some((task) => task.status === "verifying")) manifest.status = "verifying";
  else if (tasks.some((task) => task.status === "running")) manifest.status = "running";
  else manifest.status = "prepared";
}

function assertPlanRepositoryLock(repository: RuntimeRepository, taskId: string): void {
  if (repository.active_task_id && repository.active_task_id !== taskId) {
    throw new Error(`Repository ${repository.name} is locked by task ${repository.active_task_id}`);
  }
}

function assertPlanInputIdentity(manifest: RuntimeManifest, task: PlanTaskRecord, repository: RuntimeRepository, input: PlanTaskInput): void {
  assertEqual(input.run_id, manifest.run_id, "plan task input run_id");
  assertEqual(input.plan_id, manifest.plan_id, "plan task input plan_id");
  assertEqual(input.plan_reference, manifest.plan_reference, "plan task input plan_reference");
  assertEqual(input.plan_version, manifest.plan_version, "plan task input plan_version");
  assertEqual(input.plan_revision, manifest.plan_revision ?? manifest.plan_version, "plan task input plan_revision");
  assertEqual(input.approved_digest, manifest.approved_digest, "plan task input approved_digest");
  assertEqual(input.task_id, task.task_id ?? task.work_id, "plan task input task_id");
  assertEqual(input.repository, repository.name, "plan task input repository");
  assertEqual(resolve(input.worktree), resolve(repository.worktree), "plan task input worktree");
  assertEqual(input.branch, repository.branch, "plan task input branch");
  assertEqual(input.base_commit, repository.base_commit, "plan task input base_commit");
  assertEqual(input.attempt, planTaskAttempt(task), "plan task input attempt");
  assertEqual(input.start_commit, task.start_commit, "plan task input start_commit");
}

async function validatePlanWorkerResult(manifest: RuntimeManifest, task: PlanTaskRecord, repository: RuntimeRepository, input: PlanTaskInput): Promise<WorkerResult> {
  const result = await readJson<WorkerResult>(input.result_path);
  await assertValid("worker-result", result);
  assertEqual(result.contract_version, 2, "plan worker result contract_version");
  assertEqual(result.run_id, manifest.run_id, "plan worker result run_id");
  assertEqual(result.plan_id, manifest.plan_id, "plan worker result plan_id");
  assertEqual(result.plan_reference, manifest.plan_reference, "plan worker result plan_reference");
  assertEqual(result.plan_version, manifest.plan_version, "plan worker result plan_version");
  assertEqual(result.plan_revision, manifest.plan_revision ?? manifest.plan_version, "plan worker result plan_revision");
  assertEqual(result.approved_digest, manifest.approved_digest, "plan worker result approved_digest");
  assertEqual(result.task_id, task.task_id ?? task.work_id, "plan worker result task_id");
  assertEqual(result.attempt, planTaskAttempt(task), "plan worker result attempt");
  assertEqual(result.repository, repository.name, "plan worker result repository");
  assertEqual(result.branch, repository.branch, "plan worker result branch");
  assertEqual(resolve(result.worktree), resolve(repository.worktree), "plan worker result worktree");
  assertEqual(result.start_commit, task.start_commit, "plan worker result start_commit");
  const head = await assertWorktree(repository);
  if (task.start_commit === null || task.start_commit === undefined) throw new Error("Plan task has no immutable start commit");
  if (result.status === "completed") {
    if (result.commits.length === 0) throw new Error("Completed plan worker result must record at least one commit");
    if (result.checks.some((check) => check.status === "failed")) throw new Error("Completed plan worker result cannot contain a failed check");
    assertEqual(result.commits.at(-1), head, "plan worker result final commit");
    const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${task.start_commit}..${head}`])).split("\n").filter(Boolean);
    if (commits.join("\n") !== result.commits.join("\n")) throw new Error("Plan worker result commits does not match the task-local start-to-head history");
    const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${task.start_commit}...${head}`])).split("\n").filter(Boolean);
    if (!sameMembers(changedFiles, result.changed_files)) throw new Error("Plan worker result changed_files does not match the task-local start-to-head diff");
    const outsideScope = changedFiles.filter((path) => !inAllowedScope(path, input.allowed_scope));
    if (outsideScope.length > 0) throw new Error(`Plan worker changed files outside allowed scope: ${outsideScope.join(", ")}`);
    const expectation = input.test_expectation;
    if (expectation?.policy === "required" && !changedFiles.some((path) => inAllowedScope(path, expectation.paths))) {
      throw new Error(`Required test policy needs a changed file in test scope: ${expectation.paths.join(", ")}`);
    }
  }
  return result;
}

async function unlockPlanDependents(runtimeRoot: string, manifest: RuntimeManifest): Promise<void> {
  const tasks = manifest.task_graph ?? [];
  const passed = new Set(tasks.filter((task) => task.outcome === "passed").map((task) => task.task_id ?? task.work_id));
  for (const repository of manifest.repositories) {
    if (repository.active_task_id) continue;
    const candidate = tasks.find((task) => task.repository === repository.name && task.outcome === "pending" && (task.status === "waiting" || task.status === undefined) && task.depends_on.every((dependency) => passed.has(dependency)));
    if (!candidate) continue;
    const startCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
    candidate.start_commit = candidate.start_commit ?? startCommit;
    candidate.ready = true;
    candidate.blocked_by = [];
    setPlanTask(manifest, candidate, "prepared");
    const workerPath = assertInside(runtimeRoot, candidate.worker_input);
    const verifierPath = assertInside(runtimeRoot, candidate.verifier_input);
    for (const path of [workerPath, verifierPath]) {
      const input = await readJson<Record<string, unknown>>(path);
      input.ready = true;
      input.blocked_by = [];
      input.start_commit = candidate.start_commit;
      await writeJsonAtomic(path, input);
    }
  }
  if (tasks.every((task) => task.outcome === "passed")) {
    manifest.plan_verifier_status = "pending";
    if (manifest.plan_verifier_input) {
      const inputPath = assertInside(runtimeRoot, manifest.plan_verifier_input);
      const input = await readJson<Record<string, unknown>>(inputPath);
      input.ready = true;
      await writeJsonAtomic(inputPath, input);
    }
  }
  refreshPlanManifestStatus(manifest);
}

function assertPlanVerifierShape(value: unknown): asserts value is { contract_version: 1; plan_reference: string; plan_id: string; plan_version: number; approved_digest: string; run_id: string; task_id: string; repository: "plan"; plan_revision: number; attempt: number; status: "pass" | "fail" | "blocked"; summary: string; tasks: Array<{ task_id: string; repository: string; status: "passed"; head_commit: string; evidence: string }>; acceptance: Array<{ criterion: string; status: string; evidence: string }>; checks: string[]; findings: Array<{ severity: string; description: string; evidence: string }>; verified_at: string } {
  if (!value || typeof value !== "object") throw new Error("Plan verifier result must be an object");
  const result = value as Record<string, unknown>;
  for (const key of ["plan_reference", "plan_id", "run_id", "summary", "verified_at"]) if (typeof result[key] !== "string" || !result[key]) throw new Error(`Plan verifier result is missing ${key}`);
  if (result.contract_version !== 1 || !["pass", "fail", "blocked"].includes(result.status as string) || !Number.isInteger(result.plan_revision) || !Number.isInteger(result.attempt)) throw new Error("Plan verifier result identity or status is invalid");
  if (!Array.isArray(result.tasks) || !Array.isArray(result.checks) || !Array.isArray(result.findings)) throw new Error("Plan verifier result tasks, checks, and findings are required");
  for (const task of result.tasks) {
    if (!task || typeof task !== "object" || typeof (task as Record<string, unknown>).task_id !== "string" || typeof (task as Record<string, unknown>).repository !== "string" || (task as Record<string, unknown>).status !== "passed" || typeof (task as Record<string, unknown>).head_commit !== "string" || typeof (task as Record<string, unknown>).evidence !== "string") throw new Error("Plan verifier task evidence is invalid");
  }
}

async function recordCumulativePlanResult(options: RecordResultOptions, runtimeRoot: string, manifest: RuntimeManifest, occurredAt: string): Promise<RuntimeManifest> {
  if (options.stage === "plan-verifier-result") {
    if (!manifest.plan_verifier_input) throw new Error("Cumulative plan is missing its final verifier input");
    const inputPath = assertInside(runtimeRoot, manifest.plan_verifier_input);
    const input = await readJson<PlanVerifierInput>(inputPath);
    assertEqual(input.plan_reference, manifest.plan_reference, "plan verifier input plan_reference");
    assertEqual(input.plan_id, manifest.plan_id, "plan verifier input plan_id");
    assertEqual(input.plan_version, manifest.plan_version, "plan verifier input plan_version");
    assertEqual(input.approved_digest, manifest.approved_digest, "plan verifier input approved_digest");
    assertEqual(input.plan_revision, manifest.plan_revision ?? manifest.plan_version, "plan verifier input plan_revision");
    assertEqual(input.run_id, manifest.run_id, "plan verifier input run_id");
    assertEqual(input.task_id, `PLAN-${manifest.plan_id}`, "plan verifier input task_id");
    assertEqual(input.repository, "plan", "plan verifier input repository");
    assertEqual(input.attempt, 0, "plan verifier input attempt");
    if (!input.ready) throw new Error("Plan final verifier is still locked");
    const tasks = manifest.task_graph ?? [];
    if (!tasks.every((task) => task.outcome === "passed")) throw new Error("Plan final verifier requires every task to pass independently");
    assertInside(runtimeRoot, input.result_path);
    const result = await readJson<Record<string, unknown>>(input.result_path);
    await assertValid("plan-verifier-result", result);
    assertPlanVerifierShape(result);
    assertEqual(result.plan_reference, manifest.plan_reference, "plan verifier result plan_reference");
    assertEqual(result.plan_id, manifest.plan_id, "plan verifier result plan_id");
    assertEqual(result.plan_version, manifest.plan_version, "plan verifier result plan_version");
    assertEqual(result.approved_digest, manifest.approved_digest, "plan verifier result approved_digest");
    assertEqual(result.run_id, manifest.run_id, "plan verifier result run_id");
    assertEqual(result.task_id, `PLAN-${manifest.plan_id}`, "plan verifier result task_id");
    assertEqual(result.repository, "plan", "plan verifier result repository");
    assertEqual(result.plan_revision, manifest.plan_revision ?? manifest.plan_version, "plan verifier result plan_revision");
    assertEqual(result.attempt, 0, "plan verifier result attempt");
    const expected = tasks.map((task) => `${task.task_id ?? task.work_id}:${task.repository}`).sort();
    const actual = result.tasks.map((task) => `${task.task_id}:${task.repository}`).sort();
    if (expected.join("\n") !== actual.join("\n")) throw new Error("Plan verifier result tasks do not match the complete plan task graph");
    for (const task of result.tasks) {
      const repository = findRepository(manifest, task.repository);
      const head = await assertWorktree(repository);
      assertEqual(task.head_commit, head, `plan verifier head for ${task.task_id}`);
    }
    const acceptance = result.acceptance as Array<{ criterion: string; status: string }>;
    if (input.acceptance_criteria.length !== acceptance.length || !sameMembers(acceptance.map((item) => item.criterion), input.acceptance_criteria)) {
      throw new Error("Plan verifier acceptance criteria do not match verifier input");
    }
    if (result.status === "pass" && acceptance.some((item) => item.status !== "passed")) throw new Error("Passing plan verifier result requires every plan acceptance criterion to pass");
    if (result.status === "fail" && !acceptance.some((item) => item.status === "failed")) throw new Error("Failing plan verifier result must identify a failed plan acceptance criterion");
    if (result.status === "blocked" && !acceptance.some((item) => item.status === "blocked")) throw new Error("Blocked plan verifier result must identify a blocked plan acceptance criterion");
    const existing = manifest.execution_events?.find((event) => event.idempotency_key === planEventKey(manifest, input.task_id, "plan", options.stage, 0));
    if (existing) return manifest;
    await chmod(input.result_path, 0o600);
    manifest.plan_verifier_status = result.status === "pass" ? "passed" : result.status === "fail" ? "failed" : "blocked";
    manifest.plan_verifier_result = input.result_path;
    manifest.evidence.push(input.result_path);
    manifest.execution_events ??= [];
    manifest.execution_events.push({ stage: "plan-verifier-result", repository: "plan", from_status: "verifying", to_status: manifest.plan_verifier_status === "passed" ? "passed" : manifest.plan_verifier_status === "failed" ? "failed" : "blocked", inferred: false, attempt: 0, result_path: input.result_path, idempotency_key: planEventKey(manifest, input.task_id, "plan", options.stage, 0), occurred_at: occurredAt });
    refreshPlanManifestStatus(manifest);
    manifest.updated_at = occurredAt;
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(assertInside(runtimeRoot, join(runtimeRoot, "runs", manifest.run_id, "manifest.json")), manifest);
    return manifest;
  }

  const taskId = options.taskId;
  if (!taskId) throw new Error("Cumulative plan result recording requires --task-id");
  if (!options.repository) throw new Error("Cumulative plan task result recording requires --repository");
  const task = findPlanTask(manifest, taskId);
  const repository = findRepository(manifest, options.repository);
  if (task.repository !== repository.name) throw new Error(`Plan task ${taskId} belongs to ${task.repository}, not ${repository.name}`);
  const workerPath = assertInside(runtimeRoot, task.worker_input);
  const verifierPath = assertInside(runtimeRoot, task.verifier_input);
  const workerInput = await readJson<PlanTaskInput>(workerPath);
  const verifierInput = await readJson<PlanTaskInput>(verifierPath);
  assertPlanInputIdentity(manifest, task, repository, workerInput);
  assertPlanInputIdentity(manifest, task, repository, verifierInput);
  assertEqual(resolve(verifierInput.worker_result!), resolve(workerInput.result_path), "plan verifier input worker_result");
  const attempt = planTaskAttempt(task);
  const existing = manifest.execution_events?.find((event) => event.idempotency_key === planEventKey(manifest, taskId, repository.name, options.stage, attempt));
  if (existing) return manifest;

  if (options.stage === "worker-started") {
    if (task.status !== "prepared" || !task.ready) throw new Error(`Plan task ${taskId} is not ready: ${(task.blocked_by ?? task.depends_on).join(", ") || "repository lock"}`);
    assertPlanRepositoryLock(repository, taskId);
    const head = await assertWorktree(repository);
    assertEqual(head, task.start_commit, "plan task start HEAD");
    repository.active_task_id = taskId;
    setPlanTask(manifest, task, "running");
    manifest.execution_events ??= [];
    manifest.execution_events.push({ stage: options.stage, repository: repository.name, from_status: "prepared", to_status: "running", inferred: false, attempt, idempotency_key: planEventKey(manifest, taskId, repository.name, options.stage, attempt), occurred_at: occurredAt });
  } else if (options.stage === "worker-result") {
    assertPlanRepositoryLock(repository, taskId);
    const result = await validatePlanWorkerResult(manifest, task, repository, workerInput);
    await chmod(workerInput.result_path, 0o600);
    if (task.status === "prepared") {
      repository.active_task_id = taskId;
      setPlanTask(manifest, task, "running");
      manifest.execution_events ??= [];
      manifest.execution_events.push({ stage: "worker-started", repository: repository.name, from_status: "prepared", to_status: "running", inferred: true, attempt, idempotency_key: planEventKey(manifest, taskId, repository.name, "worker-started", attempt), occurred_at: occurredAt });
    }
    if (task.status !== "running") throw new Error(`Plan worker-result requires running task status, received ${task.status ?? "pending"}`);
    task.worker_result = workerInput.result_path;
    const target = result.status === "completed" ? "verifying" : result.status;
    setPlanTask(manifest, task, target, target === "verifying" ? undefined : target);
    manifest.execution_events ??= [];
    manifest.execution_events.push({ stage: options.stage, repository: repository.name, from_status: "running", to_status: target, inferred: false, attempt, result_path: workerInput.result_path, idempotency_key: planEventKey(manifest, taskId, repository.name, options.stage, attempt), occurred_at: occurredAt });
  } else {
    if (task.status !== "verifying") throw new Error(`Plan verifier-result requires verifying task status, received ${task.status ?? "pending"}`);
    if (repository.active_task_id !== taskId) throw new Error(`Repository ${repository.name} is not locked by task ${taskId}`);
    const worker = await validatePlanWorkerResult(manifest, task, repository, workerInput);
    const result = await readJson<VerifierResult>(verifierInput.result_path);
    await assertValid("verifier-result", result);
    assertEqual(result.contract_version, 2, "plan verifier result contract_version");
    assertEqual(result.run_id, manifest.run_id, "plan verifier result run_id");
    assertEqual(result.plan_id, manifest.plan_id, "plan verifier result plan_id");
    assertEqual(result.plan_reference, manifest.plan_reference, "plan verifier result plan_reference");
    assertEqual(result.plan_version, manifest.plan_version, "plan verifier result plan_version");
    assertEqual(result.plan_revision, manifest.plan_revision ?? manifest.plan_version, "plan verifier result plan_revision");
    assertEqual(result.approved_digest, manifest.approved_digest, "plan verifier result approved_digest");
    assertEqual(result.task_id, taskId, "plan verifier result task_id");
    assertEqual(result.attempt, attempt, "plan verifier result attempt");
    assertEqual(result.repository, repository.name, "plan verifier result repository");
    assertEqual(result.branch, repository.branch, "plan verifier result branch");
    if (!result.worktree) throw new Error("Plan verifier result is missing worktree");
    assertEqual(resolve(result.worktree), resolve(repository.worktree), "plan verifier result worktree");
    assertEqual(result.start_commit, task.start_commit, "plan verifier result start_commit");
    if (!sameMembers(result.acceptance.map((item) => item.criterion), verifierInput.acceptance_criteria)) throw new Error("Plan verifier acceptance criteria do not match verifier input");
    const acceptanceStatuses = result.acceptance.map((item) => item.status);
    if (result.status === "pass" && acceptanceStatuses.some((status) => status !== "passed")) throw new Error("Passing plan verifier result requires every acceptance criterion to pass");
    if (result.status === "fail" && !acceptanceStatuses.includes("failed")) throw new Error("Failing plan verifier result must identify a failed acceptance criterion");
    if (result.status === "blocked" && !acceptanceStatuses.includes("blocked")) throw new Error("Blocked plan verifier result must identify a blocked acceptance criterion");
    assertEqual(worker.commits.at(-1), await git(repository.worktree, ["rev-parse", "HEAD"]), "verified plan task commit");
    await chmod(verifierInput.result_path, 0o600);
    const target = result.status === "pass" ? "passed" : result.status === "fail" ? "failed" : "blocked";
    task.verifier_result = verifierInput.result_path;
    setPlanTask(manifest, task, target, target);
    if (target === "passed") delete repository.active_task_id;
    manifest.execution_events ??= [];
    manifest.execution_events.push({ stage: options.stage, repository: repository.name, from_status: "verifying", to_status: target, inferred: false, attempt, result_path: verifierInput.result_path, idempotency_key: planEventKey(manifest, taskId, repository.name, options.stage, attempt), occurred_at: occurredAt });
    if (target === "passed") await unlockPlanDependents(runtimeRoot, manifest);
  }
  refreshPlanManifestStatus(manifest);
  manifest.updated_at = occurredAt;
  await assertValid("runtime-manifest", manifest);
  await writeJsonAtomic(assertInside(runtimeRoot, join(runtimeRoot, "runs", manifest.run_id, "manifest.json")), manifest);
  return manifest;
}

export async function recordResult(options: RecordResultOptions): Promise<RuntimeManifest> {
  assertIdentifier(options.runId, "run ID", /^[0-9]{8}T[0-9]{6}Z-[a-f0-9]{8}$/);
  if (options.repository) assertIdentifier(options.repository, "repository", /^[a-z][a-z0-9-]*$/);
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;

  return withExclusiveFile(lockPath, async () => {
    const manifest = await readJson<RuntimeManifest>(manifestPath);
    await assertValid("runtime-manifest", manifest);
    assertEqual(manifest.run_id, options.runId, "manifest run_id");
    const occurredAt = (options.now ?? new Date()).toISOString();
    if (manifest.source_kind === "plan" && manifest.task_graph && manifest.plan_verifier_input) return recordCumulativePlanResult(options, runtimeRoot, manifest, occurredAt);
    if (!options.repository) throw new Error("Legacy result recording requires --repository");
    const requestedRepository = options.repository;
    const repository = findRepository(manifest, requestedRepository);
    const attempt = repository.repair_attempts ?? 0;
    assertInside(runtimeRoot, repository.worktree);
    const taskBriefPath = assertInside(runtimeRoot, manifest.task_brief);
    const workerInputPath = assertInside(runtimeRoot, repository.worker_input);
    const verifierInputPath = assertInside(runtimeRoot, repository.verifier_input);
    const brief = await readJson<TaskBrief>(taskBriefPath);
    await assertValid("task-brief", brief);
    assertPlanWorkItemAssociation(manifest, brief, options.repository);
    assertTaskIdentity(manifest, brief, options.repository);
    const target = brief.repositories.find((candidate) => candidate.name === options.repository)!;
    const targetScope = target.scope ?? brief.scope;
    const implementationScope = target.implementation_scope ?? brief.implementation_scope ?? targetScope;
    const testExpectation: TestExpectation = target.test_expectation ?? brief.test_expectation ?? {
      policy: "verifier-only",
      paths: [],
      rationale: "Legacy task brief has no authorized test edit scope; verifier evidence is required.",
    };
    const workerInput = await readJson<WorkerInput>(workerInputPath);
    const verifierInput = await readJson<VerifierInput>(verifierInputPath);
    assertInside(runtimeRoot, workerInput.result_path);
    assertInside(runtimeRoot, verifierInput.worker_result);
    assertInside(runtimeRoot, verifierInput.result_path);
    assertEqual(workerInput.repository, repository.name, "worker input repository");
    assertEqual(resolve(workerInput.task_brief), resolve(taskBriefPath), "worker input task_brief");
    assertEqual(resolve(workerInput.worktree), resolve(repository.worktree), "worker input worktree");
    assertEqual(workerInput.branch, repository.branch, "worker input branch");
    assertEqual(workerInput.base_commit, repository.base_commit, "worker input base_commit");
    if (!sameMembers(workerInput.allowed_scope, targetScope)) throw new Error("worker input allowed_scope does not match task brief repository scope");
    if ((target.implementation_scope || brief.implementation_scope) && (!workerInput.implementation_scope || !sameMembers(workerInput.implementation_scope, implementationScope))) {
      throw new Error("worker input implementation_scope does not match task brief");
    }
    if ((target.test_expectation || brief.test_expectation) && (!workerInput.test_expectation || !sameTestExpectation(workerInput.test_expectation, testExpectation))) {
      throw new Error("worker input test_expectation does not match task brief");
    }
    assertEqual(verifierInput.repository, repository.name, "verifier input repository");
    assertEqual(resolve(verifierInput.task_brief), resolve(taskBriefPath), "verifier input task_brief");
    assertEqual(resolve(verifierInput.worktree), resolve(repository.worktree), "verifier input worktree");
    assertEqual(verifierInput.branch, repository.branch, "verifier input branch");
    assertEqual(verifierInput.base_commit, repository.base_commit, "verifier input base_commit");
    assertEqual(resolve(verifierInput.worker_result), resolve(workerInput.result_path), "verifier input worker_result");
    if (!sameMembers(verifierInput.acceptance_criteria, target.acceptance_criteria ?? brief.acceptance_criteria)) throw new Error("verifier input acceptance_criteria does not match task brief repository criteria");
    if ((target.test_expectation || brief.test_expectation) && (!verifierInput.test_expectation || !sameTestExpectation(verifierInput.test_expectation, testExpectation))) {
      throw new Error("verifier input test_expectation does not match task brief");
    }

    const existing = manifest.execution_events?.find((event) => event.idempotency_key === eventKey(options.runId, requestedRepository, options.stage, attempt));
    const currentStatus = repository.status ?? manifest.status;

    if (options.stage === "worker-started") {
      if (existing) return manifest;
      if (currentStatus === "waiting") throw new Error(`worker-started for ${repository.name} is blocked by: ${(repository.depends_on ?? []).join(", ")}`);
      if (currentStatus !== "prepared") throw new Error(`worker-started requires prepared repository status, received ${currentStatus}`);
      const head = await assertWorktree(repository);
      assertEqual(head, repository.base_commit, "worker start HEAD");
      appendEvent(manifest, options.stage, options.repository, "prepared", "running", occurredAt, false, attempt);
    } else if (options.stage === "worker-result") {
      const target = await validateWorkerResult(manifest, repository, workerInput, testExpectation);
      await chmod(workerInput.result_path, 0o600);
      if (existing) return manifest;
      if (currentStatus === "prepared") {
        appendEvent(manifest, "worker-started", options.repository, "prepared", "running", occurredAt, true, attempt);
      }
      if ((repository.status ?? manifest.status) !== "running") throw new Error(`worker-result requires running repository status, received ${repository.status ?? manifest.status}`);
      appendEvent(manifest, options.stage, options.repository, "running", target, occurredAt, false, attempt, workerInput.result_path);
    } else {
      const target = await validateVerifierResult(manifest, repository, verifierInput);
      await chmod(verifierInput.result_path, 0o600);
      if (existing) return manifest;
      if (currentStatus !== "verifying") throw new Error(`verifier-result requires verifying repository status, received ${currentStatus}`);
      appendEvent(manifest, options.stage, options.repository, "verifying", target, occurredAt, false, attempt, verifierInput.result_path);
      const item = manifest.plan_work_items?.find((candidate) => candidate.work_id === manifest.work_id);
      if (manifest.plan_work_items && (!item || item.repository !== options.repository)) {
        throw new Error("Plan work item identity does not match the verified manifest work and repository");
      }
      if (item) item.outcome = target;
      if (target === "passed") await unlockDependents(runtimeRoot, manifest);
    }

    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    return manifest;
  });
}
