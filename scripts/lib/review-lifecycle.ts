import { readFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { assertCleanRepository, git } from "./git.js";
import { assertInside, ensurePrivateDirectory, readJsonRegularInside, withExclusiveFile, writeJsonAtomic } from "./io.js";
import type { ExecutionEvent, MergeConfirmationRecord, ReviewCommand, ReviewPreparation, ReviewPublicationRecord, RuntimeManifest, RuntimeRepository, TaskBrief, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";

interface ResultInput {
  result_path: string;
}

interface WorkerResult {
  work_id: string;
  run_id: string;
  plan_reference?: string;
  plan_id?: string;
  plan_version?: number;
  plan_revision?: number;
  approved_digest?: string;
  task_id?: string;
  attempt?: number;
  repository: string;
  status: "completed" | "blocked" | "failed";
  branch: string;
  worktree: string;
  start_commit?: string;
  commits: string[];
  changed_files: string[];
}

interface VerifierResult {
  work_id: string;
  run_id: string;
  plan_reference?: string;
  plan_id?: string;
  plan_version?: number;
  plan_revision?: number;
  approved_digest?: string;
  task_id?: string;
  attempt?: number;
  repository: string;
  branch?: string;
  worktree?: string;
  start_commit?: string;
  status: "pass" | "fail" | "blocked";
  summary: string;
  acceptance: Array<{ criterion: string; status: "passed" | "failed" | "blocked"; evidence: string }>;
  checks: string[];
  findings: Array<{ severity: string; description: string; evidence: string }>;
}

interface PlanVerifierResult {
  contract_version: 1;
  plan_reference: string;
  plan_id: string;
  plan_version: number;
  approved_digest: string;
  run_id: string;
  task_id: string;
  repository: "plan";
  plan_revision: number;
  attempt: number;
  status: "pass" | "fail" | "blocked";
  summary: string;
  tasks: Array<{ task_id: string; repository: string; status: "passed"; head_commit: string; evidence: string }>;
  acceptance: Array<{ criterion: string; status: "passed" | "failed" | "blocked"; evidence: string }>;
  checks: string[];
  findings: Array<{ severity: string; description: string; evidence: string }>;
  verified_at: string;
}

interface ApprovedPlanBrief {
  kind: "approved-plan-execution";
  plan_id: string;
  plan_reference: string;
  plan_version: number;
  approved_digest: string;
  run_id: string;
  requested_outcome: string;
}

export interface PrepareLifecycleOptions {
  workspaceRoot: string;
  runId: string;
  repository: string;
  taskId?: string;
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

async function preparePlanTaskRepair(options: PrepareLifecycleOptions, runtimeRoot: string, manifestPath: string, manifest: RuntimeManifest, repository: RuntimeRepository, config: WorkspaceConfig): Promise<RepairPreparation> {
  if (!options.taskId) throw new Error("Cumulative plan repair requires a task ID");
  const task = manifest.task_graph?.find((candidate) => (candidate.task_id ?? candidate.work_id) === options.taskId);
  if (!task || task.repository !== repository.name) throw new Error(`Plan repair task does not belong to ${repository.name}: ${options.taskId}`);
  const attempt = task.attempt ?? 0;
  const maximumAttempts = config.workflow.maximum_repair_attempts;
  const lastEvent = manifest.execution_events?.filter((event) => event.repository === repository.name && event.idempotency_key.includes(`:${options.taskId}:`)).at(-1);
  if (task.status === "prepared" && lastEvent?.stage === "repair-prepared" && lastEvent.attempt === attempt) {
    return { status: "prepared", attempt, maximum_attempts: maximumAttempts, manifest: manifestPath, worker_input: task.worker_input, verifier_input: task.verifier_input };
  }
  if (task.status !== "failed" && task.status !== "blocked") throw new Error(`Plan repair preparation requires failed or blocked task status, received ${task.status ?? "pending"}`);
  const occurredAt = (options.now ?? new Date()).toISOString();
  if (attempt >= maximumAttempts) {
    const priorStatus = task.status;
    task.status = "blocked";
    task.outcome = "blocked";
    manifest.execution_events ??= [];
    repository.active_task_id = options.taskId;
    manifest.execution_events.push({ stage: "repair-exhausted", repository: repository.name, from_status: priorStatus ?? "failed", to_status: "blocked", inferred: false, attempt, idempotency_key: `${manifest.run_id}:execution:${repository.name}:${options.taskId}:repair-exhausted:attempt-${attempt}`, occurred_at: occurredAt });
    manifest.status = "blocked";
    manifest.updated_at = occurredAt;
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    return { status: "exhausted", attempt, maximum_attempts: maximumAttempts, manifest: manifestPath };
  }
  const priorWorkerInput = await readJsonRegularInside<Record<string, unknown>>(runtimeRoot, task.worker_input, "Plan worker input");
  const priorVerifierInput = await readJsonRegularInside<Record<string, unknown>>(runtimeRoot, task.verifier_input, "Plan verifier input");
  const nextAttempt = attempt + 1;
  const runRoot = join(runtimeRoot, "runs", options.runId);
  const workerInputPath = join(runRoot, `${options.taskId}-repair-${nextAttempt}-worker-input.json`);
  const verifierInputPath = join(runRoot, `${options.taskId}-repair-${nextAttempt}-verifier-input.json`);
  const workerResultPath = join(runtimeRoot, "results", `${options.runId}-${options.taskId}-repair-${nextAttempt}-worker.json`);
  const verifierResultPath = join(runtimeRoot, "results", `${options.runId}-${options.taskId}-repair-${nextAttempt}-verifier.json`);
  const startCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
  const workerInput = { ...priorWorkerInput, role: "repair-worker", attempt: nextAttempt, start_commit: startCommit, findings: ["Repair the failed or blocked plan task using the prior independent evidence."], previous_worker_result: priorWorkerInput.result_path, previous_verifier_result: priorVerifierInput.result_path, result_path: workerResultPath, ready: true, blocked_by: [] };
  const verifierInput = { ...priorVerifierInput, role: "verifier", attempt: nextAttempt, worker_result: workerResultPath, result_path: verifierResultPath, ready: true, blocked_by: [] };
  await writeJsonAtomic(workerInputPath, workerInput);
  await writeJsonAtomic(verifierInputPath, verifierInput);
  task.worker_input = workerInputPath;
  task.verifier_input = verifierInputPath;
  task.attempt = nextAttempt;
  task.start_commit = startCommit;
  task.status = "prepared";
  task.outcome = "pending";
  task.ready = true;
  task.blocked_by = [];
  const summary = manifest.plan_work_items?.find((candidate) => candidate.work_id === task.work_id);
  if (summary) {
    summary.task_input = workerInputPath;
    summary.verifier_input = verifierInputPath;
    summary.attempt = nextAttempt;
    summary.start_commit = startCommit;
    summary.status = "prepared";
    summary.outcome = "pending";
    summary.ready = true;
    summary.blocked_by = [];
  }
  repository.worker_input = workerInputPath;
  repository.verifier_input = verifierInputPath;
  repository.repair_attempts = Math.max(repository.repair_attempts ?? 0, nextAttempt);
  repository.active_task_id = options.taskId;
  manifest.execution_events ??= [];
  manifest.execution_events.push({ stage: "repair-prepared", repository: repository.name, from_status: "failed", to_status: "running", inferred: false, attempt: nextAttempt, idempotency_key: `${manifest.run_id}:execution:${repository.name}:${options.taskId}:repair-prepared:attempt-${nextAttempt}`, occurred_at: occurredAt });
  manifest.evidence.push(workerInputPath, verifierInputPath);
  repository.status = "prepared";
  manifest.status = "prepared";
  manifest.updated_at = occurredAt;
  await assertValid("runtime-manifest", manifest);
  await writeJsonAtomic(manifestPath, manifest);
  return { status: "prepared", attempt: nextAttempt, maximum_attempts: maximumAttempts, manifest: manifestPath, worker_input: workerInputPath, verifier_input: verifierInputPath };
}

export interface RecordReviewPublicationOptions {
  workspaceRoot: string;
  runId: string;
  repository: string;
  status: "published" | "failed";
  tool: "gh" | "glab" | "manual";
  pullRequest?: string;
  evidence: string;
  authorized?: boolean;
  now?: Date;
}

export interface ConfirmMergeOptions {
  workspaceRoot: string;
  runId: string;
  repository: string;
  mergeCommit: string;
  evidence: string;
  author: string;
  now?: Date;
}

async function assertValid(name: "workspace" | "runtime-manifest" | "task-brief" | "worker-result" | "verifier-result" | "plan-verifier-result" | "review-preparation" | "review-publication-record" | "merge-confirmation-record", value: unknown): Promise<void> {
  const errors = await validateContract(name as Parameters<typeof validateContract>[0], value);
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

async function assertCurrentWorker(runtimeRoot: string, manifest: RuntimeManifest, repository: RuntimeRepository, workerInput: ResultInput): Promise<WorkerResult> {
  const worker = await readJsonRegularInside<WorkerResult>(runtimeRoot, workerInput.result_path, "Worker result");
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
    const manifest = await readJsonRegularInside<RuntimeManifest>(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository(manifest, options.repository);
    if (manifest.task_graph) {
      if (!options.taskId) throw new Error("Cumulative plan repair requires a task ID; holistic verifier failures require an approved plan revision");
      return preparePlanTaskRepair(options, runtimeRoot, manifestPath, manifest, repository, config);
    }
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
    const taskBrief = await readJsonRegularInside<TaskBrief>(runtimeRoot, taskBriefPath, "Task brief");
    await assertValid("task-brief", taskBrief);
    const taskTarget = taskBrief.repositories.find((candidate) => candidate.name === repository.name);
    if (!taskTarget) throw new Error(`Task brief does not include repository ${repository.name}`);
    const priorWorkerInput = await readJsonRegularInside<ResultInput>(runtimeRoot, repository.worker_input, "Worker input");
    const priorVerifierInput = await readJsonRegularInside<ResultInput>(runtimeRoot, repository.verifier_input, "Verifier input");
    assertInside(runtimeRoot, priorWorkerInput.result_path);
    assertInside(runtimeRoot, priorVerifierInput.result_path);
    await assertCurrentWorker(runtimeRoot, manifest, repository, priorWorkerInput);
    const verifier = await readJsonRegularInside<VerifierResult>(runtimeRoot, priorVerifierInput.result_path, "Verifier result");
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

function shellQuote(value: string): string {
  return `'${value.replaceAll("'", `'"'"'`)}'`;
}

function command(description: string, cwd: string, argv: string[]): ReviewCommand {
  return { description, cwd, argv, shell: `cd -- ${shellQuote(cwd)} && ${argv.map(shellQuote).join(" ")}` };
}

async function assertBranchName(repository: string, branch: string, label: string): Promise<void> {
  try {
    await git(repository, ["check-ref-format", "--branch", branch]);
  } catch {
    throw new Error(`${label} is not a valid Git branch name`);
  }
}

function planReviewBody(brief: ApprovedPlanBrief, manifest: RuntimeManifest, repository: RuntimeRepository, verifier: PlanVerifierResult): string {
  const tasks = (manifest.task_graph ?? []).filter((task) => task.repository === repository.name);
  const taskEvidence = tasks.map((task) => `- ${task.task_id ?? task.work_id}: worker=${task.worker_result}; verifier=${task.verifier_result}`).join("\n");
  const acceptance = verifier.acceptance.map((item) => `- [x] ${item.criterion}`).join("\n");
  const checks = verifier.checks.length > 0 ? verifier.checks.map((check) => `- ${check}`).join("\n") : "- No cumulative command was configured; holistic evidence is recorded.";
  return `## Summary\n\n${brief.requested_outcome}\n\n## Approved plan\n\n- Plan: ${brief.plan_reference}\n- Version: ${brief.plan_version}\n- Approved digest: ${brief.approved_digest}\n\n## Tasks for ${repository.name}\n\n${taskEvidence}\n\n## Cumulative acceptance\n\n${acceptance}\n\n## Holistic verification\n\n${checks}\n\nFinal plan verifier result: ${manifest.plan_verifier_result}. No push or pull request was performed.\n`;
}

async function loadApprovedPlanForReview(runtimeRoot: string, manifest: RuntimeManifest, repository: RuntimeRepository): Promise<{ brief: ApprovedPlanBrief; verifier: PlanVerifierResult; verifierInput: Record<string, unknown> }> {
  if (manifest.source_kind !== "plan" || (manifest.contract_version !== 2 && manifest.contract_version !== 3) || !manifest.task_graph || !manifest.plan_reference || !manifest.plan_id || !manifest.plan_version || !manifest.plan_revision || !manifest.approved_digest || !manifest.plan_verifier_input || !manifest.plan_verifier_result || manifest.plan_verifier_status !== "passed") {
    throw new Error("Plan review requires a complete approved-plan runtime and a passing holistic verifier");
  }
  const brief = await readJsonRegularInside<ApprovedPlanBrief>(runtimeRoot, manifest.task_brief, "Approved plan brief");
  if (brief.kind !== "approved-plan-execution" || brief.plan_id !== manifest.plan_id || brief.plan_reference !== manifest.plan_reference || brief.plan_version !== manifest.plan_version || brief.approved_digest !== manifest.approved_digest || brief.run_id !== manifest.run_id) throw new Error("Approved plan brief identity does not match the runtime manifest");
  const verifierInput = await readJsonRegularInside<Record<string, unknown>>(runtimeRoot, manifest.plan_verifier_input, "Plan verifier input");
  const verifier = await readJsonRegularInside<PlanVerifierResult>(runtimeRoot, manifest.plan_verifier_result, "Plan verifier result");
  await assertValid("plan-verifier-result", verifier);
  if (verifier.plan_reference !== manifest.plan_reference || verifier.plan_id !== manifest.plan_id || verifier.plan_version !== manifest.plan_version || verifier.approved_digest !== manifest.approved_digest || verifier.run_id !== manifest.run_id || verifier.task_id !== `PLAN-${manifest.plan_id}` || verifier.repository !== "plan" || verifier.plan_revision !== manifest.plan_revision || verifier.attempt !== 0 || verifier.status !== "pass") throw new Error("Plan verifier result identity does not match the approved runtime");
  const tasks = manifest.task_graph;
  if (!tasks.every((task) => task.outcome === "passed" && task.worker_result && task.verifier_result)) throw new Error("Plan review requires every task to pass with worker and verifier evidence");
  const expected = tasks.map((task) => `${task.task_id ?? task.work_id}:${task.repository}`).sort();
  const actual = verifier.tasks.map((task) => `${task.task_id}:${task.repository}`).sort();
  if (expected.join("\n") !== actual.join("\n")) throw new Error("Plan verifier result does not cover the complete task graph");
  for (const task of tasks) {
    const taskRepository = manifest.repositories.find((candidate) => candidate.name === task.repository);
    if (!taskRepository) throw new Error(`Plan task repository is missing from the runtime: ${task.repository}`);
    const worker = await readJsonRegularInside<WorkerResult>(runtimeRoot, task.worker_result!, "Plan worker result");
    const taskVerifier = await readJsonRegularInside<VerifierResult>(runtimeRoot, task.verifier_result!, "Plan task verifier result");
    await assertValid("worker-result", worker);
    await assertValid("verifier-result", taskVerifier);
    if (worker.run_id !== manifest.run_id || worker.plan_id !== manifest.plan_id || worker.plan_reference !== manifest.plan_reference || worker.plan_version !== manifest.plan_version || worker.plan_revision !== manifest.plan_revision || worker.approved_digest !== manifest.approved_digest || worker.task_id !== (task.task_id ?? task.work_id) || worker.attempt !== (task.attempt ?? 0) || worker.repository !== task.repository || worker.branch !== taskRepository.branch || resolve(worker.worktree) !== resolve(taskRepository.worktree) || worker.start_commit !== task.start_commit || worker.status !== "completed") throw new Error(`Plan worker evidence identity is invalid for ${task.task_id ?? task.work_id}`);
    if (taskVerifier.run_id !== manifest.run_id || taskVerifier.plan_id !== manifest.plan_id || taskVerifier.plan_reference !== manifest.plan_reference || taskVerifier.plan_version !== manifest.plan_version || taskVerifier.plan_revision !== manifest.plan_revision || taskVerifier.approved_digest !== manifest.approved_digest || taskVerifier.task_id !== (task.task_id ?? task.work_id) || taskVerifier.attempt !== (task.attempt ?? 0) || taskVerifier.repository !== task.repository || taskVerifier.branch !== taskRepository.branch || resolve(taskVerifier.worktree ?? "") !== resolve(taskRepository.worktree) || taskVerifier.start_commit !== task.start_commit || taskVerifier.status !== "pass") throw new Error(`Plan task verifier evidence identity is invalid for ${task.task_id ?? task.work_id}`);
    const head = await git(taskRepository.worktree, ["rev-parse", "HEAD"]);
    if (worker.commits.at(-1) !== head || taskVerifier.acceptance.some((item) => item.status !== "passed")) throw new Error(`Plan task evidence is stale or incomplete for ${task.task_id ?? task.work_id}`);
  }
  if (!(manifest.task_graph ?? []).some((task) => task.repository === repository.name)) throw new Error(`Plan has no task for repository ${repository.name}`);
  return { brief, verifier, verifierInput };
}

async function preparePlanReview(options: PrepareLifecycleOptions, workspaceRoot: string, runtimeRoot: string, manifestPath: string, manifest: RuntimeManifest, repository: RuntimeRepository, config: WorkspaceConfig): Promise<ReviewPreparation> {
  const data = await loadApprovedPlanForReview(runtimeRoot, manifest, repository);
  const headCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
  const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${headCommit}`])).split("\n").filter(Boolean);
  const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${headCommit}`])).split("\n").filter(Boolean);
  if (commits.length === 0 || changedFiles.length === 0) throw new Error("Plan review preparation requires committed cumulative changes");
  const workerResult = manifest.plan_verifier_input!;
  const verifierResult = manifest.plan_verifier_result!;
  if (repository.review_preparation) {
    const existing = await readJsonRegularInside<ReviewPreparation>(runtimeRoot, repository.review_preparation, "Review preparation");
    await assertValid("review-preparation", existing);
    if (existing.contract_version === 2 && existing.head_commit === headCommit && existing.worker_result === workerResult && existing.verifier_result === verifierResult) return existing;
  }
  const remote = (await git(repository.worktree, ["remote"])).split("\n").filter(Boolean).includes("origin") ? "origin" : null;
  const baseBranch = config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch;
  await assertBranchName(repository.worktree, repository.branch, "Recorded source branch");
  await assertBranchName(repository.worktree, baseBranch, "Configured default branch");
  const baseRepository = assertInside(workspaceRoot, join(workspaceRoot, repository.base_path));
  const testCommands = Array.isArray(data.verifierInput.verification_commands) ? data.verifierInput.verification_commands.filter((value): value is string => typeof value === "string") : [];
  const commands = {
    diff: command("Inspect the exact cumulative base-to-head diff", repository.worktree, ["git", "diff", "--stat", `${repository.base_commit}...${headCommit}`]),
    commits: command("Inspect the exact cumulative commit list", repository.worktree, ["git", "log", "--oneline", `${repository.base_commit}..${headCommit}`]),
    show: command("Inspect the exact verified cumulative head commit", repository.worktree, ["git", "show", "--stat", "--oneline", headCommit]),
    tests: testCommands.map((value) => command(`Run recorded holistic verification: ${value}`, repository.worktree, ["sh", "-lc", value])),
    switch_target: command("Switch the base repository to the configured target branch", baseRepository, ["git", "switch", baseBranch]),
    merge: command("Human-only merge of the exact verified cumulative head", baseRepository, ["git", "merge", "--no-ff", headCommit]),
  };
  const preparedAt = (options.now ?? new Date()).toISOString();
  const preparationPath = join(runtimeRoot, "runs", options.runId, `${repository.name}-draft-pr.json`);
  const confirmArgv = ["node", ".agents/bin/cc.mjs", "confirm-merge", "--run-id", manifest.run_id, "--repository", repository.name, "--merge-commit", "<full-merge-commit>", "--author", "<author-slug>", "--evidence", "<single-line-human-merge-evidence>"];
  const preparation: ReviewPreparation = {
    contract_version: 2, work_id: manifest.work_id, run_id: manifest.run_id, repository: repository.name,
    status: remote ? "ready-for-publication" : "ready-for-local-review", remote, base_branch: baseBranch, head_branch: repository.branch,
    base_commit: repository.base_commit, head_commit: headCommit, commits, changed_files: changedFiles, title: `${manifest.work_id}: ${data.brief.requested_outcome}`,
    body: planReviewBody(data.brief, manifest, repository, data.verifier), worker_result: workerResult, verifier_result: verifierResult,
    blockers: [], prepared_at: preparedAt, commands, merge_handoff: { status: "merge-confirmation-required", confirmation_argv: confirmArgv, confirmation_shell: confirmArgv.map(shellQuote).join(" ") },
  };
  await assertValid("review-preparation", preparation);
  await writeJsonAtomic(preparationPath, preparation);
  repository.review_preparation = preparationPath;
  repository.review_state = preparation.status as "ready-for-local-review" | "ready-for-publication";
  addEvidence(manifest, preparationPath);
  const eventKey = `${manifest.run_id}:execution:${repository.name}:review-prepared:attempt-${repository.repair_attempts ?? 0}`;
  if (!manifest.execution_events?.some((event) => event.idempotency_key === eventKey)) addExecutionEvent(manifest, repository.name, "review-prepared", "passed", "passed", repository.repair_attempts ?? 0, preparedAt, preparationPath);
  await assertValid("runtime-manifest", manifest);
  await writeJsonAtomic(manifestPath, manifest);
  return preparation;
}

export async function prepareReview(options: PrepareLifecycleOptions): Promise<ReviewPreparation> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;
  const config = await loadWorkspace(workspaceRoot);

  return withExclusiveFile(lockPath, async () => {
    const manifest = await readJsonRegularInside<RuntimeManifest>(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository(manifest, options.repository);
    if (manifest.status !== "passed") throw new Error(`Draft review preparation requires passed status, received ${manifest.status}`);
    if (manifest.task_graph) return preparePlanReview(options, workspaceRoot, runtimeRoot, manifestPath, manifest, repository, config);
    const taskBrief = await readJsonRegularInside<TaskBrief>(runtimeRoot, manifest.task_brief, "Task brief");
    await assertValid("task-brief", taskBrief);
    const workerInput = await readJsonRegularInside<ResultInput>(runtimeRoot, repository.worker_input, "Worker input");
    const verifierInput = await readJsonRegularInside<ResultInput>(runtimeRoot, repository.verifier_input, "Verifier input");
    assertInside(runtimeRoot, workerInput.result_path);
    assertInside(runtimeRoot, verifierInput.result_path);
    const worker = await assertCurrentWorker(runtimeRoot, manifest, repository, workerInput);
    const verifier = await readJsonRegularInside<VerifierResult>(runtimeRoot, verifierInput.result_path, "Verifier result");
    await assertValid("verifier-result", verifier);
    assertVerifier(manifest, repository, taskBrief, verifier);
    if (verifier.status !== "pass") throw new Error(`Draft review preparation requires a passing verifier result, received ${verifier.status}`);
    const headCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
    const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${headCommit}`])).split("\n").filter(Boolean);
    const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${headCommit}`])).split("\n").filter(Boolean);
    if (commits.length === 0 || changedFiles.length === 0) throw new Error("Draft review preparation requires committed changes");
    if (repository.review_preparation) {
      const existing = await readJsonRegularInside<ReviewPreparation>(runtimeRoot, repository.review_preparation, "Review preparation");
      await assertValid("review-preparation", existing);
      if (existing.contract_version === 2 && existing.head_commit === headCommit && existing.worker_result === workerInput.result_path && existing.verifier_result === verifierInput.result_path) {
        return existing;
      }
    }
    const remotes = (await git(repository.worktree, ["remote"])).split("\n").filter(Boolean);
    const remote = remotes.includes("origin") ? "origin" : null;
    const baseBranch = config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch;
    await assertBranchName(repository.worktree, repository.branch, "Recorded source branch");
    await assertBranchName(repository.worktree, baseBranch, "Configured default branch");
    const taskTarget = taskBrief.repositories.find((candidate) => candidate.name === repository.name);
    const testArgv = taskTarget?.verification_commands ?? taskBrief.verification_commands;
    const commands = {
      diff: command("Inspect the exact base-to-head diff", repository.worktree, ["git", "diff", "--stat", `${repository.base_commit}...${headCommit}`]),
      commits: command("Inspect the exact commit list", repository.worktree, ["git", "log", "--oneline", `${repository.base_commit}..${headCommit}`]),
      show: command("Inspect the exact verified head commit", repository.worktree, ["git", "show", "--stat", "--oneline", headCommit]),
      tests: testArgv.map((value) => command(`Run recorded verification: ${value}`, repository.worktree, ["sh", "-lc", value])),
      switch_target: command("Switch the base repository to the configured target branch", assertInside(workspaceRoot, join(workspaceRoot, repository.base_path)), ["git", "switch", baseBranch]),
      merge: command("Human-only merge of the exact verified head", assertInside(workspaceRoot, join(workspaceRoot, repository.base_path)), ["git", "merge", "--no-ff", headCommit]),
    };
    const confirmArgv = ["node", ".agents/bin/cc.mjs", "confirm-merge", "--run-id", manifest.run_id, "--repository", repository.name, "--merge-commit", "<full-merge-commit>", "--author", "<author-slug>", "--evidence", "<single-line-human-merge-evidence>"];
    const preparedAt = (options.now ?? new Date()).toISOString();
    const preparationPath = join(runtimeRoot, "runs", options.runId, `${repository.name}-draft-pr.json`);
    const preparation: ReviewPreparation = {
      contract_version: 2,
      work_id: manifest.work_id,
      run_id: manifest.run_id,
      repository: repository.name,
      status: remote ? "ready-for-publication" : "ready-for-local-review",
      remote,
      base_branch: baseBranch,
      head_branch: repository.branch,
      base_commit: repository.base_commit,
      head_commit: headCommit,
      commits,
      changed_files: changedFiles,
      title: `${manifest.work_id}: ${taskBrief.requested_outcome}`,
      body: reviewBody(taskBrief, verifier),
      worker_result: workerInput.result_path,
      verifier_result: verifierInput.result_path,
      blockers: [],
      prepared_at: preparedAt,
      commands,
      merge_handoff: { status: "merge-confirmation-required", confirmation_argv: confirmArgv, confirmation_shell: confirmArgv.map(shellQuote).join(" ") },
    };
    if (!sameMembers(worker.changed_files, changedFiles)) {
      throw new Error("Current Git diff does not match the recorded worker result");
    }
    await assertValid("review-preparation", preparation);
    await writeJsonAtomic(preparationPath, preparation);
    repository.review_preparation = preparationPath;
    repository.review_state = preparation.status as "ready-for-local-review" | "ready-for-publication";
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
    const manifest = await readJsonRegularInside<RuntimeManifest>(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository(manifest, options.repository);
    if (!repository.review_preparation) throw new Error("Prepare the draft pull-request handoff before recording publication");
    const preparation = await readJsonRegularInside<ReviewPreparation>(runtimeRoot, repository.review_preparation, "Review preparation");
    await assertValid("review-preparation", preparation);
    if (preparation.work_id !== manifest.work_id || preparation.run_id !== manifest.run_id || preparation.repository !== repository.name) throw new Error("Review preparation identity does not match the active run");
    if (preparation.contract_version !== 2 || preparation.status !== "ready-for-publication" || preparation.remote !== "origin") {
      throw new Error(`Remote publication requires a ready-for-publication handoff with origin; current state is ${preparation.status}`);
    }
    if (!options.authorized) throw new Error("Remote publication recording requires explicit authorization confirmation");
    await assertCleanRepository(repository.worktree);
    const headCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
    if (headCommit !== preparation.head_commit) throw new Error("Worktree HEAD changed after review preparation");
    const evidence = safeReviewEvidence(options.evidence, "Publication evidence");
    const pullRequest = options.status === "published" ? safeReviewEvidence(options.pullRequest ?? "", "Pull-request reference") : null;
    if (options.status === "failed" && options.pullRequest) throw new Error("Failed publication cannot record a pull-request reference");
    const recordPath = join(runtimeRoot, "runs", options.runId, `${repository.name}-review-publication.json`);
    const record: ReviewPublicationRecord = {
      contract_version: options.status === "published" ? 2 : 1,
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
      ...(options.status === "published" ? { review_state: "published-for-review" as const } : {}),
    };
    await assertValid("review-publication-record", record);
    if (repository.review_publication) {
      const existing = await readJsonRegularInside<ReviewPublicationRecord>(runtimeRoot, repository.review_publication, "Review publication record");
      await assertValid("review-publication-record", existing);
      const comparable = (value: ReviewPublicationRecord) => JSON.stringify({ ...value, recorded_at: null });
      if (comparable(existing) !== comparable(record)) throw new Error("Review publication was already recorded with different confirmed evidence");
      return existing;
    }
    await writeJsonAtomic(recordPath, record);
    repository.review_publication = recordPath;
    repository.review_state = options.status === "published" ? "published-for-review" : "ready-for-publication";
    addEvidence(manifest, recordPath);
    manifest.updated_at = record.recorded_at;
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    return record;
  });
}

async function isAncestor(repository: string, ancestor: string, descendant: string): Promise<boolean> {
  try { await git(repository, ["merge-base", "--is-ancestor", ancestor, descendant]); return true; } catch { return false; }
}

export async function confirmMerge(options: ConfirmMergeOptions): Promise<MergeConfirmationRecord> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const manifestPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, "manifest.json"));
  if (!/^[a-f0-9]{40,64}$/.test(options.mergeCommit)) throw new Error("Merge commit must be a full lowercase Git object ID");
  const author = safeReviewEvidence(options.author, "Author");
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(author)) throw new Error("Author must be a lowercase slug");
  const evidence = safeReviewEvidence(options.evidence, "Merge evidence");
  const config = await loadWorkspace(workspaceRoot);
  return withExclusiveFile(`${manifestPath}.lock`, async () => {
    const manifest = await readJsonRegularInside<RuntimeManifest>(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository(manifest, options.repository);
    if (!repository.review_preparation) throw new Error("Prepare review before confirming a merge");
    const preparation = await readJsonRegularInside<ReviewPreparation>(runtimeRoot, repository.review_preparation, "Review preparation");
    await assertValid("review-preparation", preparation);
    if (preparation.work_id !== manifest.work_id || preparation.run_id !== manifest.run_id || preparation.repository !== repository.name) throw new Error("Review preparation identity does not match the active run");
    if (preparation.contract_version !== 2 || preparation.head_commit !== await git(repository.worktree, ["rev-parse", "HEAD"])) throw new Error("Review preparation does not match the current verified head");
    const baseBranch = config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch;
    if (preparation.base_branch !== baseBranch) throw new Error("Review preparation target differs from the configured default branch");
    const baseRepository = assertInside(workspaceRoot, join(workspaceRoot, repository.base_path));
    await git(baseRepository, ["cat-file", "-e", `${options.mergeCommit}^{commit}`]);
    const targetRefs = [`refs/heads/${baseBranch}`, `refs/remotes/origin/${baseBranch}`];
    let targetRef: string | null = null;
    let targetCommit: string | null = null;
    for (const ref of targetRefs) {
      try {
        const commit = await git(baseRepository, ["rev-parse", "--verify", `${ref}^{commit}`]);
        if (await isAncestor(baseRepository, options.mergeCommit, commit)) { targetRef = ref; targetCommit = commit; break; }
      } catch { /* Missing exact configured target ref is not merge evidence. */ }
    }
    if (!targetRef || !targetCommit) throw new Error(`Reported merge commit is not reachable from the configured default target ${baseBranch}`);
    if (!await isAncestor(baseRepository, preparation.head_commit, options.mergeCommit)) throw new Error("Verified review head is not reachable from the reported merge commit");
    if (!await isAncestor(baseRepository, preparation.base_commit, options.mergeCommit)) throw new Error("Recorded base is not reachable from the reported merge commit");
    const finishArgv = ["node", ".agents/bin/cc.mjs", "finish-work", "--run-id", manifest.run_id, "--repository", repository.name, "--outcome", "merged", "--author", author, "--merge-commit", options.mergeCommit];
    const publication = repository.review_publication ? await readJsonRegularInside<ReviewPublicationRecord>(runtimeRoot, repository.review_publication, "Review publication record") : null;
    if (publication?.status === "published" && publication.pull_request) finishArgv.push("--pull-request", publication.pull_request);
    const record: MergeConfirmationRecord = {
      contract_version: 1, work_id: manifest.work_id, run_id: manifest.run_id, repository: repository.name,
      status: "closeout-ready", base_branch: baseBranch, target_ref: targetRef, target_commit: targetCommit,
      head_commit: preparation.head_commit, merge_commit: options.mergeCommit, evidence,
      finish_work_argv: finishArgv, finish_work_shell: finishArgv.map(shellQuote).join(" "),
      idempotency_key: `${manifest.run_id}:merge-confirmation:${repository.name}`, confirmed_at: (options.now ?? new Date()).toISOString(),
    };
    await assertValid("merge-confirmation-record", record);
    const recordPath = join(runtimeRoot, "runs", options.runId, `${repository.name}-merge-confirmation.json`);
    if (repository.merge_confirmation) {
      const existing = await readJsonRegularInside<MergeConfirmationRecord>(runtimeRoot, repository.merge_confirmation, "Merge confirmation record");
      await assertValid("merge-confirmation-record", existing);
      const comparable = (value: MergeConfirmationRecord) => JSON.stringify({ ...value, confirmed_at: null });
      if (comparable(existing) !== comparable(record)) throw new Error("Merge was already confirmed with different evidence");
      return existing;
    }
    await writeJsonAtomic(recordPath, record);
    repository.merge_confirmation = recordPath;
    repository.review_state = "closeout-ready";
    addEvidence(manifest, recordPath);
    manifest.updated_at = record.confirmed_at;
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    return record;
  });
}
