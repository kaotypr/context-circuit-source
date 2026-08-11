import { readFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { assertCleanRepository, git } from "./git.js";
import { assertInside, ensurePrivateDirectory, withExclusiveFile, writeJsonAtomic } from "./io.js";
import type { ExecutionEvent, ReviewPreparation, ReviewPublicationRecord, RuntimeManifest, RuntimeRepository, TaskBrief, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";

interface ResultInput {
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
}

interface VerifierResult {
  work_id: string;
  run_id: string;
  repository: string;
  status: "pass" | "fail" | "blocked";
  summary: string;
  acceptance: Array<{ criterion: string; status: "passed" | "failed" | "blocked"; evidence: string }>;
  checks: string[];
  findings: Array<{ severity: string; description: string; evidence: string }>;
}

export interface PrepareLifecycleOptions {
  workspaceRoot: string;
  runId: string;
  repository: string;
  now?: Date;
}

export interface RepairPreparation {
  status: "prepared" | "exhausted";
  attempt: number;
  maximum_attempts: number;
  manifest: string;
  worker_input?: string;
  verifier_input?: string;
}

export interface RecordReviewPublicationOptions {
  workspaceRoot: string;
  runId: string;
  repository: string;
  status: "published" | "failed";
  tool: "gh" | "glab" | "manual";
  pullRequest?: string;
  evidence: string;
  now?: Date;
}

async function readJson<T>(path: string): Promise<T> {
  return JSON.parse(await readFile(path, "utf8")) as T;
}

async function assertValid(name: "workspace" | "runtime-manifest" | "task-brief" | "worker-result" | "verifier-result" | "review-preparation" | "review-publication-record", value: unknown): Promise<void> {
  const errors = await validateContract(name, value);
  if (errors.length > 0) throw new Error(`Invalid ${name}: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}

function findRepository(manifest: RuntimeManifest, name: string): RuntimeRepository {
  const repository = manifest.repositories.find((candidate) => candidate.name === name);
  if (!repository) throw new Error(`Run ${manifest.run_id} has no repository named ${name}`);
  return repository;
}

function addEvidence(manifest: RuntimeManifest, ...paths: string[]): void {
  for (const path of paths) if (!manifest.evidence.includes(path)) manifest.evidence.push(path);
}

function addExecutionEvent(
  manifest: RuntimeManifest,
  repository: string,
  stage: ExecutionEvent["stage"],
  fromStatus: ExecutionEvent["from_status"],
  toStatus: ExecutionEvent["to_status"],
  attempt: number,
  occurredAt: string,
  resultPath?: string,
): void {
  const event: ExecutionEvent = {
    stage,
    repository,
    from_status: fromStatus,
    to_status: toStatus,
    inferred: false,
    attempt,
    idempotency_key: `${manifest.run_id}:execution:${repository}:${stage}:attempt-${attempt}`,
    occurred_at: occurredAt,
  };
  if (resultPath) event.result_path = resultPath;
  manifest.execution_events ??= [];
  manifest.execution_events.push(event);
  const runtimeRepository = findRepository(manifest, repository);
  runtimeRepository.status = toStatus;
  const statuses = manifest.repositories.map((candidate) => candidate.status ?? manifest.status);
  if (statuses.every((status) => status === "passed")) manifest.status = "passed";
  else if (statuses.includes("failed")) manifest.status = "failed";
  else if (statuses.includes("blocked")) manifest.status = "blocked";
  else if (statuses.includes("verifying")) manifest.status = "verifying";
  else if (statuses.includes("running")) manifest.status = "running";
  else manifest.status = "prepared";
  manifest.updated_at = occurredAt;
}

async function loadWorkspace(workspaceRoot: string): Promise<WorkspaceConfig> {
  const config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  await assertValid("workspace", config);
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace: ${semanticErrors.join("; ")}`);
  return config;
}

function sameMembers(left: string[], right: string[]): boolean {
  return left.slice().sort().join("\n") === right.slice().sort().join("\n");
}

async function assertCurrentWorker(manifest: RuntimeManifest, repository: RuntimeRepository, workerInput: ResultInput): Promise<WorkerResult> {
  const worker = await readJson<WorkerResult>(workerInput.result_path);
  await assertValid("worker-result", worker);
  if (worker.work_id !== manifest.work_id || worker.run_id !== manifest.run_id || worker.repository !== repository.name) {
    throw new Error("Worker result identity does not match the active run");
  }
  if (worker.status !== "completed") throw new Error(`Review lifecycle requires a completed worker, received ${worker.status}`);
  if (worker.branch !== repository.branch || resolve(worker.worktree) !== resolve(repository.worktree)) {
    throw new Error("Worker result branch or worktree does not match the runtime manifest");
  }
  await assertCleanRepository(repository.worktree);
  const branch = await git(repository.worktree, ["branch", "--show-current"]);
  if (branch !== repository.branch) throw new Error(`Worktree branch mismatch: expected ${repository.branch}, received ${branch}`);
  const head = await git(repository.worktree, ["rev-parse", "HEAD"]);
  if (worker.commits.at(-1) !== head) throw new Error("Current worktree HEAD does not match the recorded worker result");
  const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${head}`])).split("\n").filter(Boolean);
  const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${head}`])).split("\n").filter(Boolean);
  if (commits.join("\n") !== worker.commits.join("\n") || !sameMembers(changedFiles, worker.changed_files)) {
    throw new Error("Worker result no longer matches the current base-to-head Git history");
  }
  return worker;
}

function assertVerifier(manifest: RuntimeManifest, repository: RuntimeRepository, brief: TaskBrief, verifier: VerifierResult): void {
  if (verifier.work_id !== manifest.work_id || verifier.run_id !== manifest.run_id || verifier.repository !== repository.name) {
    throw new Error("Verifier result identity does not match the active run");
  }
  const acceptanceCriteria = brief.repositories.find((candidate) => candidate.name === repository.name)?.acceptance_criteria ?? brief.acceptance_criteria;
  if (!sameMembers(verifier.acceptance.map((item) => item.criterion), acceptanceCriteria)) {
    throw new Error("Verifier acceptance criteria do not match the task brief");
  }
  if (verifier.status === "pass" && verifier.acceptance.some((item) => item.status !== "passed")) {
    throw new Error("Passing verifier result no longer has complete passing acceptance evidence");
  }
  if (verifier.status === "fail" && !verifier.acceptance.some((item) => item.status === "failed")) {
    throw new Error("Failing verifier result no longer identifies failed acceptance evidence");
  }
}

function repairFindings(verifier: VerifierResult): string[] {
  const findings = verifier.findings.map((finding) => `[${finding.severity}] ${finding.description} Evidence: ${finding.evidence}`);
  for (const acceptance of verifier.acceptance) {
    if (acceptance.status !== "passed") findings.push(`Acceptance ${acceptance.status}: ${acceptance.criterion}. Evidence: ${acceptance.evidence}`);
  }
  return findings;
}

export async function prepareRepair(options: PrepareLifecycleOptions): Promise<RepairPreparation> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;
  const config = await loadWorkspace(workspaceRoot);

  return withExclusiveFile(lockPath, async () => {
    const manifest = await readJson<RuntimeManifest>(manifestPath);
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository(manifest, options.repository);
    const attempt = repository.repair_attempts ?? 0;
    const maximumAttempts = config.workflow.maximum_repair_attempts;
    const lastEvent = manifest.execution_events?.filter((event) => event.repository === repository.name).at(-1);
    const repositoryStatus = repository.status ?? manifest.status;
    if (repositoryStatus === "running" && lastEvent?.stage === "repair-prepared" && lastEvent.attempt === attempt) {
      return { status: "prepared", attempt, maximum_attempts: maximumAttempts, manifest: manifestPath, worker_input: repository.worker_input, verifier_input: repository.verifier_input };
    }
    if (repositoryStatus === "blocked" && lastEvent?.stage === "repair-exhausted") {
      return { status: "exhausted", attempt, maximum_attempts: maximumAttempts, manifest: manifestPath };
    }
    if (repositoryStatus !== "failed") throw new Error(`Repair preparation requires failed status, received ${repositoryStatus}`);

    const occurredAt = (options.now ?? new Date()).toISOString();
    if (attempt >= maximumAttempts) {
      addExecutionEvent(manifest, repository.name, "repair-exhausted", "failed", "blocked", attempt, occurredAt);
      manifest.warnings.push(`Maximum repair attempts exhausted for ${repository.name}: ${maximumAttempts}`);
      await assertValid("runtime-manifest", manifest);
      await writeJsonAtomic(manifestPath, manifest);
      return { status: "exhausted", attempt, maximum_attempts: maximumAttempts, manifest: manifestPath };
    }

    const taskBriefPath = assertInside(runtimeRoot, manifest.task_brief);
    const taskBrief = await readJson<TaskBrief>(taskBriefPath);
    await assertValid("task-brief", taskBrief);
    const taskTarget = taskBrief.repositories.find((candidate) => candidate.name === repository.name);
    if (!taskTarget) throw new Error(`Task brief does not include repository ${repository.name}`);
    const priorWorkerInput = await readJson<ResultInput>(assertInside(runtimeRoot, repository.worker_input));
    const priorVerifierInput = await readJson<ResultInput>(assertInside(runtimeRoot, repository.verifier_input));
    assertInside(runtimeRoot, priorWorkerInput.result_path);
    assertInside(runtimeRoot, priorVerifierInput.result_path);
    await assertCurrentWorker(manifest, repository, priorWorkerInput);
    const verifier = await readJson<VerifierResult>(priorVerifierInput.result_path);
    await assertValid("verifier-result", verifier);
    assertVerifier(manifest, repository, taskBrief, verifier);
    if (verifier.status !== "fail") throw new Error(`Repair preparation requires a failing verifier result, received ${verifier.status}`);
    const findings = repairFindings(verifier);
    if (findings.length === 0) throw new Error("Failing verifier result contains no actionable findings");

    const nextAttempt = attempt + 1;
    const runRoot = join(runtimeRoot, "runs", options.runId);
    const workerInputPath = join(runRoot, `${repository.name}-repair-${nextAttempt}-worker-input.json`);
    const verifierInputPath = join(runRoot, `${repository.name}-repair-${nextAttempt}-verifier-input.json`);
    const workerResultPath = join(runtimeRoot, "results", `${options.runId}-${repository.name}-repair-${nextAttempt}-worker.json`);
    const verifierResultPath = join(runtimeRoot, "results", `${options.runId}-${repository.name}-repair-${nextAttempt}-verifier.json`);
    const repositoryConfig = config.repositories[repository.name];
    if (!repositoryConfig) throw new Error(`Workspace has no repository named ${repository.name}`);
    const instructionPaths = [
      join(workspaceRoot, "AGENTS.md"),
      join(workspaceRoot, "agents", "repository-worker.md"),
      join(workspaceRoot, "agents", `${repositoryConfig.agent}.md`),
    ];
    const workerInput = {
      contract_version: 1,
      role: "repair-worker",
      attempt: nextAttempt,
      task_brief: taskBriefPath,
      repository: repository.name,
      worktree: repository.worktree,
      branch: repository.branch,
      base_commit: repository.base_commit,
      allowed_scope: taskTarget.scope ?? taskBrief.scope,
      implementation_scope: taskTarget.implementation_scope ?? taskBrief.implementation_scope,
      test_expectation: taskTarget.test_expectation ?? taskBrief.test_expectation,
      findings,
      previous_worker_result: priorWorkerInput.result_path,
      previous_verifier_result: priorVerifierInput.result_path,
      instruction_paths: instructionPaths,
      result_contract: join(workspaceRoot, ".agents", "contracts", "worker-result.schema.json"),
      result_path: workerResultPath,
    };
    const verifierInput = {
      contract_version: 1,
      role: "verifier",
      read_only: true,
      attempt: nextAttempt,
      task_brief: taskBriefPath,
      repository: repository.name,
      worktree: repository.worktree,
      branch: repository.branch,
      base_commit: repository.base_commit,
      worker_result: workerResultPath,
      acceptance_criteria: taskTarget.acceptance_criteria ?? taskBrief.acceptance_criteria,
      test_expectation: taskTarget.test_expectation ?? taskBrief.test_expectation,
      verification_commands: taskTarget.verification_commands ?? taskBrief.verification_commands,
      instruction_paths: [join(workspaceRoot, "AGENTS.md"), join(workspaceRoot, "agents", "verifier.md")],
      result_contract: join(workspaceRoot, ".agents", "contracts", "verifier-result.schema.json"),
      result_path: verifierResultPath,
    };
    await writeJsonAtomic(workerInputPath, workerInput);
    await writeJsonAtomic(verifierInputPath, verifierInput);
    repository.worker_input = workerInputPath;
    repository.verifier_input = verifierInputPath;
    repository.repair_attempts = nextAttempt;
    addEvidence(manifest, workerInputPath, verifierInputPath);
    addExecutionEvent(manifest, repository.name, "repair-prepared", "failed", "running", nextAttempt, occurredAt);
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    return { status: "prepared", attempt: nextAttempt, maximum_attempts: maximumAttempts, manifest: manifestPath, worker_input: workerInputPath, verifier_input: verifierInputPath };
  });
}

function reviewBody(brief: TaskBrief, verifier: VerifierResult): string {
  const acceptance = brief.acceptance_criteria.map((criterion) => `- [x] ${criterion}`).join("\n");
  const checks = verifier.checks.length > 0 ? verifier.checks.map((check) => `- ${check}`).join("\n") : "- No repository command was configured; verifier evidence is recorded.";
  return `## Summary\n\n${brief.requested_outcome}\n\n## Acceptance\n\n${acceptance}\n\n## Verification\n\n${checks}\n\nPrepared from run \`${brief.run_id}\`. No push or pull request was performed.\n`;
}

export async function prepareReview(options: PrepareLifecycleOptions): Promise<ReviewPreparation> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;
  const config = await loadWorkspace(workspaceRoot);

  return withExclusiveFile(lockPath, async () => {
    const manifest = await readJson<RuntimeManifest>(manifestPath);
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository(manifest, options.repository);
    if (manifest.status !== "passed") throw new Error(`Draft review preparation requires passed status, received ${manifest.status}`);
    const taskBrief = await readJson<TaskBrief>(assertInside(runtimeRoot, manifest.task_brief));
    await assertValid("task-brief", taskBrief);
    const workerInput = await readJson<ResultInput>(assertInside(runtimeRoot, repository.worker_input));
    const verifierInput = await readJson<ResultInput>(assertInside(runtimeRoot, repository.verifier_input));
    assertInside(runtimeRoot, workerInput.result_path);
    assertInside(runtimeRoot, verifierInput.result_path);
    const worker = await assertCurrentWorker(manifest, repository, workerInput);
    const verifier = await readJson<VerifierResult>(verifierInput.result_path);
    await assertValid("verifier-result", verifier);
    assertVerifier(manifest, repository, taskBrief, verifier);
    if (verifier.status !== "pass") throw new Error(`Draft review preparation requires a passing verifier result, received ${verifier.status}`);
    const headCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
    const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${headCommit}`])).split("\n").filter(Boolean);
    const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${headCommit}`])).split("\n").filter(Boolean);
    if (commits.length === 0 || changedFiles.length === 0) throw new Error("Draft review preparation requires committed changes");
    if (repository.review_preparation) {
      const existing = await readJson<ReviewPreparation>(assertInside(runtimeRoot, repository.review_preparation));
      await assertValid("review-preparation", existing);
      if (existing.head_commit === headCommit && existing.worker_result === workerInput.result_path && existing.verifier_result === verifierInput.result_path) {
        return existing;
      }
    }
    const remotes = (await git(repository.worktree, ["remote"])).split("\n").filter(Boolean);
    const remote = remotes.includes("origin") ? "origin" : null;
    const blockers = remote ? [] : ["Repository has no origin remote; configure one before pushing or opening a draft pull request."];
    const preparedAt = (options.now ?? new Date()).toISOString();
    const preparationPath = join(runtimeRoot, "runs", options.runId, `${repository.name}-draft-pr.json`);
    const preparation: ReviewPreparation = {
      contract_version: 1,
      work_id: manifest.work_id,
      run_id: manifest.run_id,
      repository: repository.name,
      status: blockers.length === 0 ? "ready" : "blocked",
      remote,
      base_branch: config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch,
      head_branch: repository.branch,
      base_commit: repository.base_commit,
      head_commit: headCommit,
      commits,
      changed_files: changedFiles,
      title: `${manifest.work_id}: ${taskBrief.requested_outcome}`,
      body: reviewBody(taskBrief, verifier),
      worker_result: workerInput.result_path,
      verifier_result: verifierInput.result_path,
      blockers,
      prepared_at: preparedAt,
    };
    if (!sameMembers(worker.changed_files, changedFiles)) {
      throw new Error("Current Git diff does not match the recorded worker result");
    }
    await assertValid("review-preparation", preparation);
    await writeJsonAtomic(preparationPath, preparation);
    repository.review_preparation = preparationPath;
    addEvidence(manifest, preparationPath);
    const eventKey = `${manifest.run_id}:execution:${repository.name}:review-prepared:attempt-${repository.repair_attempts ?? 0}`;
    if (!manifest.execution_events?.some((event) => event.idempotency_key === eventKey)) {
      addExecutionEvent(manifest, repository.name, "review-prepared", "passed", "passed", repository.repair_attempts ?? 0, preparedAt, preparationPath);
    }
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    return preparation;
  });
}

function safeReviewEvidence(value: string, label: string): string {
  const trimmed = value.trim();
  if (!trimmed || /[\r\n]/.test(trimmed)) throw new Error(`${label} must be a non-empty single line`);
  if (/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----|https?:\/\/[^\s/@:]+:[^\s/@]+@/i.test(trimmed)) throw new Error(`${label} appears to contain credentials`);
  return trimmed;
}

export async function recordReviewPublication(options: RecordReviewPublicationOptions): Promise<ReviewPublicationRecord> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const manifestPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, "manifest.json"));
  return withExclusiveFile(`${manifestPath}.lock`, async () => {
    const manifest = await readJson<RuntimeManifest>(manifestPath);
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository(manifest, options.repository);
    if (!repository.review_preparation) throw new Error("Prepare the draft pull-request handoff before recording publication");
    const preparation = await readJson<ReviewPreparation>(assertInside(runtimeRoot, repository.review_preparation));
    await assertValid("review-preparation", preparation);
    if (preparation.status !== "ready") throw new Error(`Draft pull-request handoff is ${preparation.status}`);
    await assertCleanRepository(repository.worktree);
    const headCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
    if (headCommit !== preparation.head_commit) throw new Error("Worktree HEAD changed after review preparation");
    const evidence = safeReviewEvidence(options.evidence, "Publication evidence");
    const pullRequest = options.status === "published" ? safeReviewEvidence(options.pullRequest ?? "", "Pull-request reference") : null;
    if (options.status === "failed" && options.pullRequest) throw new Error("Failed publication cannot record a pull-request reference");
    const recordPath = join(runtimeRoot, "runs", options.runId, `${repository.name}-review-publication.json`);
    const record: ReviewPublicationRecord = {
      contract_version: 1,
      work_id: manifest.work_id,
      run_id: manifest.run_id,
      repository: repository.name,
      status: options.status,
      tool: options.tool,
      pull_request: pullRequest,
      evidence,
      head_commit: headCommit,
      idempotency_key: `${manifest.run_id}:review-publication:${repository.name}`,
      recorded_at: (options.now ?? new Date()).toISOString(),
    };
    await assertValid("review-publication-record", record);
    if (repository.review_publication) {
      const existing = await readJson<ReviewPublicationRecord>(assertInside(runtimeRoot, repository.review_publication));
      await assertValid("review-publication-record", existing);
      const comparable = (value: ReviewPublicationRecord) => JSON.stringify({ ...value, recorded_at: null });
      if (comparable(existing) !== comparable(record)) throw new Error("Review publication was already recorded with different confirmed evidence");
      return existing;
    }
    await writeJsonAtomic(recordPath, record);
    repository.review_publication = recordPath;
    addEvidence(manifest, recordPath);
    manifest.updated_at = record.recorded_at;
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    return record;
  });
}
