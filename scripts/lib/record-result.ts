import { chmod, readFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { assertCleanRepository, git } from "./git.js";
import { assertInside, ensurePrivateDirectory, withExclusiveFile, writeJsonAtomic } from "./io.js";
import type { ExecutionEvent, RuntimeManifest, RuntimeRepository, TaskBrief, TestExpectation } from "./types.js";
import { validateContract } from "./validation.js";

export type RecordStage = "worker-started" | "worker-result" | "verifier-result";

interface WorkerInput {
  task_brief: string;
  repository: string;
  worktree: string;
  branch: string;
  base_commit: string;
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
  worker_result: string;
  acceptance_criteria: string[];
  test_expectation?: TestExpectation;
  result_path: string;
}

interface WorkerResult {
  work_id: string;
  run_id: string;
  repository: string;
  status: "completed" | "blocked" | "failed";
  branch: string;
  worktree: string;
  commits: string[];
  changed_files: string[];
  checks: Array<{ status: "passed" | "failed" | "not-run" }>;
}

interface VerifierResult {
  work_id: string;
  run_id: string;
  repository: string;
  status: "pass" | "fail" | "blocked";
  acceptance: Array<{ criterion: string; status: string; evidence: string }>;
}

export interface RecordResultOptions {
  workspaceRoot: string;
  runId: string;
  repository: string;
  stage: RecordStage;
  now?: Date;
}

function assertIdentifier(value: string, label: string, pattern: RegExp): void {
  if (!pattern.test(value)) throw new Error(`Invalid ${label}: ${value}`);
}

async function readJson<T>(path: string): Promise<T> {
  return JSON.parse(await readFile(path, "utf8")) as T;
}

async function assertValid(name: "runtime-manifest" | "task-brief" | "worker-result" | "verifier-result", value: unknown): Promise<void> {
  const errors = await validateContract(name, value);
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

export async function recordResult(options: RecordResultOptions): Promise<RuntimeManifest> {
  assertIdentifier(options.runId, "run ID", /^[0-9]{8}T[0-9]{6}Z-[a-f0-9]{8}$/);
  assertIdentifier(options.repository, "repository", /^[a-z][a-z0-9-]*$/);
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;

  return withExclusiveFile(lockPath, async () => {
    const manifest = await readJson<RuntimeManifest>(manifestPath);
    await assertValid("runtime-manifest", manifest);
    assertEqual(manifest.run_id, options.runId, "manifest run_id");
    const repository = findRepository(manifest, options.repository);
    const attempt = repository.repair_attempts ?? 0;
    assertInside(runtimeRoot, repository.worktree);
    const taskBriefPath = assertInside(runtimeRoot, manifest.task_brief);
    const workerInputPath = assertInside(runtimeRoot, repository.worker_input);
    const verifierInputPath = assertInside(runtimeRoot, repository.verifier_input);
    const brief = await readJson<TaskBrief>(taskBriefPath);
    await assertValid("task-brief", brief);
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

    const existing = manifest.execution_events?.find((event) => event.idempotency_key === eventKey(options.runId, options.repository, options.stage, attempt));
    const occurredAt = (options.now ?? new Date()).toISOString();
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
      for (const item of manifest.plan_work_items ?? []) {
        if (item.repository === options.repository) item.outcome = target;
      }
      if (target === "passed") await unlockDependents(runtimeRoot, manifest);
    }

    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    return manifest;
  });
}
