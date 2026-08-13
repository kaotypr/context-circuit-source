import { randomBytes } from "node:crypto";
import { access, lstat, readFile, rm } from "node:fs/promises";
import { join, relative, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { assertCleanRepository, git } from "./git.js";
import { generateRunId, slugify } from "./ids.js";
import { assertInside, ensurePrivateDirectory, writeJsonAtomic } from "./io.js";
import { resolveRootPlanDirectory, setPlanState, validatePlanDirectory } from "./plans.js";
import type { ActivityCapability, PlanExecutionRequest, PlanIndex, PlanWorkItem, RuntimeManifest, RuntimeRepository, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";

export interface PreparedPlanRepository {
  name: string;
  branch: string;
  worktree: string;
  ready: boolean;
  blockedBy: string[];
  taskInputs: string[];
  verifierInputs: string[];
}

export interface PreparedPlan {
  planId: string;
  planReference: string;
  planVersion: number;
  approvedDigest: string;
  runId: string;
  manifest: string;
  planBrief: string;
  repositories: PreparedPlanRepository[];
  preparationStatus: "prepared";
}

export interface PreparePlanOptions {
  workspaceRoot: string;
  request: PlanExecutionRequest;
  now?: Date;
  discriminator?: string;
  availableCapabilities?: ActivityCapability[];
}

interface RepositoryBase {
  name: string;
  path: string;
  commit: string;
  defaultBranch: string;
  agent: string;
}

function contractMessages(errors: Awaited<ReturnType<typeof validateContract>>): string[] {
  return errors.map((error) => `${error.instancePath || "/"} ${error.message}`);
}

async function isAbsent(path: string): Promise<boolean> {
  try { await lstat(path); return false; } catch (error) { return (error as NodeJS.ErrnoException).code === "ENOENT"; }
}

async function assertNoSymlinkDirectory(path: string, label: string): Promise<void> {
  try {
    const info = await lstat(path);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`${label} must be a real directory: ${path}`);
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
  }
}

function graphErrors(items: PlanWorkItem[]): string[] {
  const byId = new Map(items.map((item) => [item.work_id, item]));
  const errors: string[] = [];
  const active = new Set<string>();
  const visited = new Set<string>();
  const visit = (id: string): void => {
    if (active.has(id)) { errors.push(`task dependency cycle includes ${id}`); return; }
    if (visited.has(id)) return;
    const item = byId.get(id);
    if (!item) { errors.push(`invalid task reference: ${id}`); return; }
    active.add(id);
    for (const dependency of item.depends_on) visit(dependency);
    active.delete(id);
    visited.add(id);
  };
  for (const item of items) {
    for (const dependency of item.depends_on) if (!byId.has(dependency)) errors.push(`${item.work_id} references unknown task ${dependency}`);
    visit(item.work_id);
  }
  return [...new Set(errors)];
}

async function assertPlanDependencies(workspaceRoot: string, index: PlanIndex): Promise<void> {
  for (const dependency of index.depends_on_plans ?? []) {
    const directory = await resolveRootPlanDirectory(workspaceRoot, dependency);
    const validation = await validatePlanDirectory(directory);
    if (!validation.index || validation.errors.length > 0) throw new Error(`Plan dependency is invalid: ${dependency}`);
    if (validation.index.status !== "completed") throw new Error(`Plan ${index.plan_id} is dependency-blocked by ${dependency}`);
  }
}

async function prepareRepositoryBases(workspaceRoot: string, config: WorkspaceConfig, repositories: string[]): Promise<Map<string, RepositoryBase>> {
  const result = new Map<string, RepositoryBase>();
  for (const name of repositories) {
    const registered = config.repositories[name];
    if (!registered) throw new Error(`repository is not registered: ${name}`);
    const path = assertInside(workspaceRoot, join(workspaceRoot, registered.path));
    await access(path);
    await assertCleanRepository(path);
    const commit = await git(path, ["rev-parse", registered.default_branch]);
    result.set(name, { name, path, commit, defaultBranch: registered.default_branch, agent: registered.agent });
    const instructionPaths = [
      join(workspaceRoot, "AGENTS.md"),
      join(workspaceRoot, "agents", `${registered.agent}.md`),
      join(workspaceRoot, "agents", "repository-worker.md"),
      join(workspaceRoot, "agents", "verifier.md"),
    ];
    await Promise.all(instructionPaths.map((instructionPath) => access(instructionPath)));
  }
  return result;
}

function taskScope(item: PlanWorkItem): string[] {
  return [...new Set([...item.scope, ...(item.test_policy === "required" ? item.test_scope : [])])];
}

function taskBrief(index: PlanIndex, items: PlanWorkItem[], runId: string, createdAt: string): Record<string, unknown> {
  return {
    contract_version: 2,
    kind: "approved-plan-execution",
    plan_id: index.plan_id,
    plan_reference: index.plan_reference!,
    plan_version: index.plan_version,
    approved_digest: index.approved_digest,
    run_id: runId,
    requested_outcome: index.title,
    repositories: index.affected_repositories,
    task_ids: items.map((item) => item.work_id),
    acceptance_criteria: items.flatMap((item) => item.acceptance_criteria),
    created_at: createdAt,
  };
}

async function assertBranchAvailable(base: RepositoryBase, branch: string): Promise<void> {
  try {
    await git(base.path, ["show-ref", "--verify", "--quiet", `refs/heads/${branch}`]);
    throw new Error(`Branch collision in ${base.name}: ${branch}`);
  } catch (error) {
    const message = (error as Error).message;
    if (!message.includes("failed in")) throw error;
  }
}

async function removeCreatedWorktrees(created: Array<{ base: RepositoryBase; worktree: string; branch: string }>): Promise<void> {
  for (const item of [...created].reverse()) {
    try { await git(item.base.path, ["worktree", "remove", "--force", item.worktree]); } catch { /* only remove output created by this attempt */ }
    try { await git(item.base.path, ["branch", "-D", item.branch]); } catch { /* branch may not have been created */ }
  }
}

export async function prepareExecutePlan(options: PreparePlanOptions): Promise<PreparedPlan> {
  const request = options.request;
  if (request.contract_version !== 1 || request.source?.kind !== "plan" || !request.source.reference || !Number.isInteger(request.source.plan_version) || !/^sha256:[a-f0-9]{64}$/.test(request.source.approved_digest)) {
    throw new Error("Invalid execute-plan request: expected one plan reference, plan version, and sha256 approval digest");
  }
  const workspaceRoot = resolve(options.workspaceRoot);
  const config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  const workspaceErrors = contractMessages(await validateContract("workspace", config));
  workspaceErrors.push(...workspaceSemanticErrors(config));
  if (workspaceErrors.length > 0) throw new Error(`Invalid workspace: ${workspaceErrors.join("; ")}`);
  const planDirectory = await resolveRootPlanDirectory(workspaceRoot, options.request.source.reference);
  const validation = await validatePlanDirectory(planDirectory);
  if (!validation.index || !validation.work_breakdown || validation.errors.length > 0) throw new Error(`Plan validation failed: ${validation.errors.join("; ") || "plan metadata is unavailable"}`);
  const index = validation.index;
  if (index.contract_version !== 2) throw new Error("execute-plan requires a root numbered plan; migrate context/plans first");
  if (index.status !== "approved") throw new Error(`Plan ${index.plan_id} is ${index.status}; explicit approval is required`);
  if (options.request.source.plan_version !== index.plan_version) throw new Error(`Plan version is stale: requested ${options.request.source.plan_version}, current ${index.plan_version}`);
  if (options.request.source.approved_digest !== index.approved_digest) throw new Error("Plan approval digest is stale or does not match the approved plan material");
  await assertPlanDependencies(workspaceRoot, index);
  const items = validation.work_breakdown.items;
  const graphIssues = graphErrors(items);
  if (graphIssues.length > 0) throw new Error(`Invalid approved task graph: ${graphIssues.join("; ")}`);
  const affectedRepositories = [...new Set(index.affected_repositories ?? items.map((item) => item.repository))];
  if (affectedRepositories.length === 0) throw new Error("Approved plan has no affected repositories");
  for (const item of items) {
    if (!config.repositories[item.repository]) throw new Error(`repository is not registered: ${item.repository}`);
    if (!affectedRepositories.includes(item.repository)) throw new Error(`Task ${item.work_id} targets a repository outside the plan affected_repositories: ${item.repository}`);
  }
  const bases = await prepareRepositoryBases(workspaceRoot, config, affectedRepositories);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  await assertNoSymlinkDirectory(runtimeRoot, "Runtime path");
  const now = options.now ?? new Date();
  const discriminator = options.discriminator ?? randomBytes(4).toString("hex");
  const runId = generateRunId(`execute-plan:${index.plan_reference}`, now, discriminator);
  const runRoot = join(runtimeRoot, "runs", runId);
  const worktreeRoot = join(runtimeRoot, "worktrees", runId);
  if (!(await isAbsent(runRoot)) || !(await isAbsent(worktreeRoot))) throw new Error(`Runtime collision: ${runId}`);
  const branches = new Map<string, string>();
  for (const repository of affectedRepositories) {
    const base = bases.get(repository)!;
    const branch = `plan/${slugify(index.plan_id)}-${runId.slice(-8)}`;
    await assertBranchAvailable(base, branch);
    branches.set(repository, branch);
  }
  const createdAt = now.toISOString();
  const planBriefPath = join(runtimeRoot, "plans", `${runId}.json`);
  const manifestPath = join(runRoot, "manifest.json");
  if (!(await isAbsent(planBriefPath))) throw new Error(`Runtime collision: ${runId}`);
  const planBrief = taskBrief(index, items, runId, createdAt);
  const runtimeRepositories: RuntimeRepository[] = [];
  const preparedRepositories: PreparedPlanRepository[] = affectedRepositories.map((name) => ({ name, branch: branches.get(name)!, worktree: assertInside(runtimeRoot, join(worktreeRoot, name)), ready: false, blockedBy: [], taskInputs: [], verifierInputs: [] }));
  const taskGraph: NonNullable<RuntimeManifest["task_graph"]> = [];
  const planWorkItems: NonNullable<RuntimeManifest["plan_work_items"]> = [];
  const firstReadyByRepository = new Set<string>();
  for (const item of items) {
    const repository = preparedRepositories.find((entry) => entry.name === item.repository)!;
    const workerInput = join(runRoot, `${item.work_id}-worker-input.json`);
    const verifierInput = join(runRoot, `${item.work_id}-verifier-input.json`);
    const verifierResultPath = join(runtimeRoot, "results", `${runId}-${item.work_id}-verifier.json`);
    const base = bases.get(item.repository)!;
    const ready = item.depends_on.length === 0 && !firstReadyByRepository.has(item.repository);
    if (ready) firstReadyByRepository.add(item.repository);
    const startCommit = ready ? base.commit : null;
    const attempt = 0;
    const common = {
      contract_version: 2,
      plan_reference: index.plan_reference!,
      plan_id: index.plan_id,
      plan_version: index.plan_version,
      plan_revision: index.plan_version,
      approved_digest: index.approved_digest,
      task_id: item.work_id,
      repository: item.repository,
      run_id: runId,
      worktree: repository.worktree,
      branch: repository.branch,
      base_commit: base.commit,
      start_commit: startCommit,
      attempt,
      ready,
      blocked_by: ready ? [] : item.depends_on,
      allowed_scope: taskScope(item),
      implementation_scope: item.scope,
      test_expectation: { policy: item.test_policy, paths: item.test_scope, rationale: item.test_rationale ?? "The approved plan task contract is authoritative." },
      instruction_paths: [join(workspaceRoot, "AGENTS.md"), join(workspaceRoot, "agents", `${base.agent}.md`), join(workspaceRoot, "agents", "repository-worker.md")],
      result_contract: join(workspaceRoot, ".agents", "contracts", "worker-result.schema.json"),
      result_path: join(runtimeRoot, "results", `${runId}-${item.work_id}-worker.json`),
    };
    await writeJsonAtomic(workerInput, { ...common, role: "repository-worker" });
    await writeJsonAtomic(verifierInput, { ...common, role: "verifier", read_only: true, worker_result: common.result_path, result_path: verifierResultPath, acceptance_criteria: item.acceptance_criteria, verification_commands: item.verification_commands, instruction_paths: [join(workspaceRoot, "AGENTS.md"), join(workspaceRoot, "agents", "verifier.md")] });
    repository.taskInputs.push(workerInput);
    repository.verifierInputs.push(verifierInput);
    taskGraph.push({ work_id: item.work_id, task_id: item.work_id, plan_id: index.plan_id, plan_reference: index.plan_reference!, plan_version: index.plan_version, approved_digest: index.approved_digest!, repository: item.repository, plan_revision: index.plan_version, attempt, start_commit: startCommit, ready, blocked_by: ready ? [] : item.depends_on, status: ready ? "prepared" : "waiting", depends_on: item.depends_on, outcome: "pending", worker_input: workerInput, verifier_input: verifierInput });
    planWorkItems.push({ work_id: item.work_id, task_id: item.work_id, plan_reference: index.plan_reference!, repository: item.repository, plan_revision: index.plan_version, attempt, start_commit: startCommit, ready, blocked_by: ready ? [] : item.depends_on, status: ready ? "prepared" : "waiting", depends_on: item.depends_on, outcome: "pending", task_input: workerInput, verifier_input: verifierInput });
    if (ready) repository.ready = true;
  }
  for (const repository of preparedRepositories) {
    const base = bases.get(repository.name)!;
    runtimeRepositories.push({ name: repository.name, base_path: relative(workspaceRoot, base.path), base_commit: base.commit, branch: repository.branch, worktree: repository.worktree, worker_input: repository.taskInputs[0] ?? join(runRoot, `${repository.name}-worker-input.json`), verifier_input: repository.verifierInputs[0] ?? join(runRoot, `${repository.name}-verifier-input.json`), task_inputs: repository.taskInputs, verifier_inputs: repository.verifierInputs, lock_path: join(runRoot, `${repository.name}.task.lock`), status: repository.ready ? "prepared" : "waiting", repair_attempts: 0 });
  }
  const planVerifierInput = join(runRoot, "plan-verifier-input.json");
  const planVerifierResult = join(runtimeRoot, "results", `${runId}-plan-verifier-attempt-0.json`);
  await writeJsonAtomic(planVerifierInput, {
    contract_version: 1,
    role: "plan-verifier",
    read_only: true,
    plan_reference: index.plan_reference!,
    plan_id: index.plan_id,
    plan_revision: index.plan_version,
    plan_version: index.plan_version,
    run_id: runId,
    approved_digest: index.approved_digest!,
    task_id: `PLAN-${index.plan_id}`,
    repository: "plan",
    attempt: 0,
    ready: false,
    worktrees: runtimeRepositories.map((repository) => ({ name: repository.name, worktree: repository.worktree, branch: repository.branch, base_commit: repository.base_commit })),
    tasks: taskGraph.map((task) => ({ task_id: task.task_id, repository: task.repository, worker_result: task.worker_input, verifier_result: task.verifier_input })),
    acceptance_criteria: ["Every approved plan task is independently verified.", "The cumulative repository worktrees pass holistic verification."],
    verification_commands: [...new Set(items.flatMap((item) => item.verification_commands))],
    instruction_paths: [join(workspaceRoot, "AGENTS.md"), join(workspaceRoot, "agents", "verifier.md")],
    result_contract: join(workspaceRoot, ".agents", "contracts", "plan-verifier-result.schema.json"),
    result_path: planVerifierResult,
  });
  const manifest: RuntimeManifest = {
    contract_version: 2,
    work_id: `PLAN-${index.plan_id}`,
    run_id: runId,
    source_kind: "plan",
    status: "preparing",
    created_at: createdAt,
    updated_at: createdAt,
    task_brief: planBriefPath,
    repositories: runtimeRepositories,
    plan_reference: index.plan_reference!,
    plan_id: index.plan_id,
    plan_version: index.plan_version,
    plan_revision: index.plan_version,
    approved_digest: index.approved_digest!,
    plan_work_items: planWorkItems,
    task_graph: taskGraph,
    plan_verifier_input: planVerifierInput,
    plan_verifier_status: "pending",
    evidence: [planBriefPath, manifestPath, planVerifierInput, ...taskGraph.flatMap((task) => [task.worker_input, task.verifier_input])],
    warnings: config.activity.provider === "none" ? ["No activity tool is configured; this run cannot guarantee exclusive ownership."] : [],
    execution_events: [],
    lifecycle_events: [{ event: "task.starting", status: config.activity.provider === "none" ? "skipped" : "pending", idempotency_key: `${runId}:task.starting`, occurred_at: createdAt }],
  };
  const manifestErrors = contractMessages(await validateContract("runtime-manifest", manifest));
  if (manifestErrors.length > 0) throw new Error(`Generated runtime-manifest is invalid: ${manifestErrors.join("; ")}`);
  await ensurePrivateDirectory(runtimeRoot);
  const createdWorktrees: Array<{ base: RepositoryBase; worktree: string; branch: string }> = [];
  let ownsPlanBrief = false;
  let ownsRunRoot = false;
  let ownsWorktreeRoot = false;
  try {
    ownsPlanBrief = true;
    await writeJsonAtomic(planBriefPath, planBrief);
    ownsRunRoot = true;
    await writeJsonAtomic(manifestPath, manifest);
    ownsWorktreeRoot = true;
    await ensurePrivateDirectory(worktreeRoot);
    for (const repository of preparedRepositories) {
      const base = bases.get(repository.name)!;
      await git(base.path, ["worktree", "add", "-b", repository.branch, repository.worktree, base.commit]);
      createdWorktrees.push({ base, worktree: repository.worktree, branch: repository.branch });
    }
    manifest.status = "prepared";
    manifest.updated_at = new Date().toISOString();
    await writeJsonAtomic(manifestPath, manifest);
    await setPlanState(planDirectory, { kind: "lifecycle", status: "in-progress", reason: "Approved plan runtime prepared.", actor: "engine", evidence: manifestPath }, new Date());
  } catch (error) {
    await removeCreatedWorktrees(createdWorktrees);
    if (ownsWorktreeRoot) try { await rm(worktreeRoot, { recursive: true, force: true }); } catch { /* preserve original error */ }
    if (ownsRunRoot) try { await rm(runRoot, { recursive: true, force: true }); } catch { /* preserve original error */ }
    if (ownsPlanBrief) try { await rm(planBriefPath, { force: true }); } catch { /* preserve original error */ }
    throw error;
  }
  return { planId: index.plan_id, planReference: index.plan_reference!, planVersion: index.plan_version, approvedDigest: index.approved_digest!, runId, manifest: manifestPath, planBrief: planBriefPath, repositories: preparedRepositories, preparationStatus: "prepared" };
}
