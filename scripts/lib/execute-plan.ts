import { randomBytes } from "node:crypto";
import { access, lstat, readFile, rm } from "node:fs/promises";
import { join, relative, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { assertCleanRepository, git } from "./git.js";
import { generateRunId, slugify } from "./ids.js";
import { assertInside, ensurePrivateDirectory, writeJsonAtomic } from "./io.js";
import { resolveRootPlanDirectory, setPlanState, validatePlanDirectory } from "./plans.js";
import type { ActivityCapability, PlanExecutionRequest, PlanIndex, PlanRuntimeRevision, PlanWorkItem, RuntimeManifest, RuntimeRepository, WorkspaceConfig } from "./types.js";
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

export interface ResumePlanOptions {
  workspaceRoot: string;
  runId: string;
  request: PlanExecutionRequest;
  reason: string;
  now?: Date;
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

function stableJson(value: unknown): string {
  if (Array.isArray(value)) return `[${value.map(stableJson).sort().join(",")}]`;
  if (value && typeof value === "object") return `{${Object.entries(value as Record<string, unknown>).sort(([a], [b]) => a.localeCompare(b)).map(([key, entry]) => `${JSON.stringify(key)}:${stableJson(entry)}`).join(",")}}`;
  return JSON.stringify(value);
}

function taskExecutionContract(item: PlanWorkItem): Record<string, unknown> {
  return {
    repository: item.repository,
    scope: item.scope,
    test_scope: item.test_scope,
    test_policy: item.test_policy,
    test_rationale: item.test_rationale ?? "The approved plan task contract is authoritative.",
    verification_commands: item.verification_commands,
    acceptance_criteria: item.acceptance_criteria,
    description: item.description ?? "",
    parent: item.parent,
    depends_on: item.depends_on,
  };
}

function inputExecutionContract(input: Record<string, unknown>): Record<string, unknown> {
  const expectation = (input.test_expectation ?? {}) as Record<string, unknown>;
  return {
    repository: input.repository,
    scope: input.implementation_scope ?? input.allowed_scope,
    test_scope: expectation.paths ?? [],
    test_policy: expectation.policy,
    test_rationale: expectation.rationale,
    verification_commands: input.verification_commands ?? [],
    acceptance_criteria: input.acceptance_criteria ?? [],
    description: input.description ?? "",
    parent: input.parent ?? null,
    depends_on: input.depends_on ?? [],
  };
}

function dependentClosure(items: PlanWorkItem[], roots: Set<string>): Set<string> {
  const closure = new Set(roots);
  let changed = true;
  while (changed) {
    changed = false;
    for (const item of items) {
      if (!closure.has(item.work_id) && item.depends_on.some((dependency) => closure.has(dependency))) {
        closure.add(item.work_id);
        changed = true;
      }
    }
  }
  return closure;
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
      description: item.description ?? "",
      parent: item.parent,
      depends_on: item.depends_on,
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

/** Resume one unmerged plan runtime after a newly approved material revision. */
export async function resumeExecutePlan(options: ResumePlanOptions): Promise<PreparedPlan> {
  if (!options.reason.trim()) throw new Error("An approved plan runtime revision requires a reason");
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const runRoot = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId));
  const manifestPath = assertInside(runRoot, join(runRoot, "manifest.json"));
  const manifest = JSON.parse(await readFile(manifestPath, "utf8")) as RuntimeManifest;
  const manifestErrors = contractMessages(await validateContract("runtime-manifest", manifest));
  if (manifestErrors.length > 0) throw new Error(`Invalid runtime-manifest: ${manifestErrors.join("; ")}`);
  if (manifest.source_kind !== "plan" || !manifest.task_graph || !manifest.plan_id || !manifest.plan_reference || !manifest.plan_version || !manifest.approved_digest) {
    throw new Error("Only an approved numbered-plan runtime can be resumed");
  }
  if (["closing", "closed"].includes(manifest.status)) throw new Error("A closed plan runtime cannot be revised");
  if (manifest.repositories.some((repository) => repository.active_task_id || repository.status === "running" || repository.status === "verifying")) {
    throw new Error("Pause active workers and verifiers before resuming an approved plan revision");
  }
  if (options.request.contract_version !== 1 || options.request.source.kind !== "plan" || options.request.source.reference !== manifest.plan_reference) {
    throw new Error("Plan revision must target the exact runtime plan reference");
  }
  const planDirectory = await resolveRootPlanDirectory(workspaceRoot, manifest.plan_reference);
  const validation = await validatePlanDirectory(planDirectory);
  if (!validation.index || !validation.work_breakdown || validation.errors.length > 0) throw new Error(`Revised plan validation failed: ${validation.errors.join("; ") || "plan metadata is unavailable"}`);
  const index = validation.index;
  if (index.status !== "approved") throw new Error("The revised plan must be explicitly approved before runtime resume");
  if (index.plan_id !== manifest.plan_id || index.plan_reference !== manifest.plan_reference) throw new Error("Revised plan identity does not match the existing runtime");
  if (index.plan_version <= manifest.plan_version) throw new Error(`Plan revision must be newer than runtime version ${manifest.plan_version}`);
  if (options.request.source.plan_version !== index.plan_version) throw new Error("Plan revision version is stale");
  if (options.request.source.approved_digest !== index.approved_digest) throw new Error("Plan revision approval digest is stale or does not match the approved plan material");
  const items = validation.work_breakdown.items;
  const graphIssues = graphErrors(items);
  if (graphIssues.length > 0) throw new Error(`Invalid revised task graph: ${graphIssues.join("; ")}`);

  const oldTasks = new Map(manifest.task_graph.map((task) => [(task.task_id ?? task.work_id), task]));
  const oldInputContracts = new Map<string, Record<string, unknown>>();
  for (const [taskId, task] of oldTasks) {
    try {
      const raw = JSON.parse(await readFile(assertInside(runtimeRoot, task.verifier_input), "utf8")) as Record<string, unknown>;
      oldInputContracts.set(taskId, inputExecutionContract(raw));
    } catch {
      oldInputContracts.set(taskId, {});
    }
  }
  const changed = new Set<string>();
  const added = new Set<string>();
  for (const item of items) {
    const prior = oldTasks.get(item.work_id);
    if (!prior) { added.add(item.work_id); changed.add(item.work_id); continue; }
    if (stableJson(oldInputContracts.get(item.work_id) ?? {}) !== stableJson(taskExecutionContract(item))) changed.add(item.work_id);
  }
  const removed = [...oldTasks.keys()].filter((taskId) => !items.some((item) => item.work_id === taskId));
  const invalidated = dependentClosure(items, changed);
  const preserved = items.filter((item) => {
    const prior = oldTasks.get(item.work_id);
    return Boolean(prior && prior.outcome === "passed" && prior.worker_result && prior.verifier_result && !invalidated.has(item.work_id));
  }).map((item) => item.work_id);
  for (const taskId of preserved) invalidated.delete(taskId);

  const now = (options.now ?? new Date()).toISOString();
  const planBriefPath = join(runtimeRoot, "plans", `${options.runId}-revision-${index.plan_version}.json`);
  const revisionPath = join(runRoot, "revisions", `revision-${String(index.plan_version).padStart(4, "0")}.json`);
  const worktreeRoot = join(runtimeRoot, "worktrees", options.runId);
  const currentHeads = new Map<string, string>();
  for (const repository of manifest.repositories) {
    await assertCleanRepository(repository.worktree);
    const branch = await git(repository.worktree, ["branch", "--show-current"]);
    if (branch !== repository.branch) throw new Error(`Runtime repository ${repository.name} is no longer on its prepared branch`);
    currentHeads.set(repository.name, await git(repository.worktree, ["rev-parse", "HEAD"]));
  }
  const planBrief = taskBrief(index, items, options.runId, now);
  const taskGraph: NonNullable<RuntimeManifest["task_graph"]> = [];
  const planWorkItems: NonNullable<RuntimeManifest["plan_work_items"]> = [];
  const repositoryByName = new Map(manifest.repositories.map((repository) => [repository.name, repository]));
  const firstReadyByRepository = new Set<string>();
  const passed = new Set(preserved);
  for (const item of items) {
    const repository = repositoryByName.get(item.repository);
    if (!repository) throw new Error(`Revised task targets an unprepared repository: ${item.repository}`);
    const prior = oldTasks.get(item.work_id);
    if (preserved.includes(item.work_id) && prior) {
      const preservedTask = { ...prior, plan_version: index.plan_version, plan_revision: index.plan_version, approved_digest: index.approved_digest! };
      taskGraph.push(preservedTask);
      const preservedSummary: NonNullable<RuntimeManifest["plan_work_items"]>[number] = { work_id: item.work_id, task_id: item.work_id, plan_reference: index.plan_reference!, repository: item.repository, plan_revision: index.plan_version, attempt: prior.attempt ?? 0, start_commit: prior.start_commit ?? repository.base_commit, ready: true, blocked_by: [], status: "passed", depends_on: item.depends_on, outcome: "passed", task_input: prior.worker_input, verifier_input: prior.verifier_input };
      if (prior.worker_result) preservedSummary.worker_result = prior.worker_result;
      if (prior.verifier_result) preservedSummary.verifier_result = prior.verifier_result;
      planWorkItems.push(preservedSummary);
      repository.task_inputs ??= [];
      repository.verifier_inputs ??= [];
      repository.task_inputs.push(prior.worker_input);
      repository.verifier_inputs.push(prior.verifier_input);
      continue;
    }
    const workerInput = join(runRoot, `${item.work_id}-revision-${index.plan_version}-worker-input.json`);
    const verifierInput = join(runRoot, `${item.work_id}-revision-${index.plan_version}-verifier-input.json`);
    const baseCommit = repository.base_commit;
    const ready = !firstReadyByRepository.has(item.repository) && item.depends_on.every((dependency) => passed.has(dependency));
    if (ready) firstReadyByRepository.add(item.repository);
    const startCommit = ready ? currentHeads.get(item.repository)! : null;
    const common = {
      contract_version: 2,
      plan_reference: index.plan_reference!, plan_id: index.plan_id, plan_version: index.plan_version, plan_revision: index.plan_version,
      approved_digest: index.approved_digest!, task_id: item.work_id, repository: item.repository, run_id: options.runId,
      worktree: repository.worktree, branch: repository.branch, base_commit: baseCommit, start_commit: startCommit,
      attempt: 0, ready, blocked_by: ready ? [] : item.depends_on, allowed_scope: taskScope(item), implementation_scope: item.scope,
      description: item.description ?? "", parent: item.parent, depends_on: item.depends_on,
      test_expectation: { policy: item.test_policy, paths: item.test_scope, rationale: item.test_rationale ?? "The approved plan task contract is authoritative." },
      result_contract: join(workspaceRoot, ".agents", "contracts", "worker-result.schema.json"),
      result_path: join(runtimeRoot, "results", `${options.runId}-${item.work_id}-revision-${index.plan_version}-worker.json`),
    };
    const verifierResultPath = join(runtimeRoot, "results", `${options.runId}-${item.work_id}-revision-${index.plan_version}-verifier.json`);
    await writeJsonAtomic(workerInput, { ...common, role: "repository-worker", task_brief: planBriefPath });
    await writeJsonAtomic(verifierInput, { ...common, role: "verifier", read_only: true, task_brief: planBriefPath, worker_result: common.result_path, result_path: verifierResultPath, acceptance_criteria: item.acceptance_criteria, verification_commands: item.verification_commands, result_contract: join(workspaceRoot, ".agents", "contracts", "verifier-result.schema.json") });
    repository.task_inputs ??= [];
    repository.verifier_inputs ??= [];
    repository.task_inputs.push(workerInput);
    repository.verifier_inputs.push(verifierInput);
    const task = { work_id: item.work_id, task_id: item.work_id, plan_id: index.plan_id, plan_reference: index.plan_reference!, plan_version: index.plan_version, approved_digest: index.approved_digest!, repository: item.repository, plan_revision: index.plan_version, attempt: 0, start_commit: startCommit, ready, blocked_by: ready ? [] : item.depends_on, status: ready ? "prepared" as const : "waiting" as const, depends_on: item.depends_on, outcome: "pending" as const, worker_input: workerInput, verifier_input: verifierInput };
    taskGraph.push(task);
    planWorkItems.push({ work_id: item.work_id, task_id: item.work_id, plan_reference: index.plan_reference!, repository: item.repository, plan_revision: index.plan_version, attempt: 0, start_commit: startCommit, ready, blocked_by: ready ? [] : item.depends_on, status: task.status, depends_on: item.depends_on, outcome: "pending", task_input: workerInput, verifier_input: verifierInput });
  }
  const planVerifierInput = join(runRoot, `plan-verifier-input-revision-${index.plan_version}.json`);
  const planVerifierResult = join(runtimeRoot, "results", `${options.runId}-plan-verifier-attempt-0-revision-${index.plan_version}.json`);
  await writeJsonAtomic(planVerifierInput, {
    contract_version: 1, role: "plan-verifier", read_only: true, plan_reference: index.plan_reference!, plan_id: index.plan_id,
    plan_revision: index.plan_version, plan_version: index.plan_version, run_id: options.runId, approved_digest: index.approved_digest!,
    task_id: `PLAN-${index.plan_id}`, repository: "plan", attempt: 0, ready: taskGraph.every((task) => task.outcome === "passed"),
    worktrees: manifest.repositories.map((repository) => ({ name: repository.name, worktree: repository.worktree, branch: repository.branch, base_commit: repository.base_commit })),
    tasks: taskGraph.map((task) => ({ task_id: task.task_id, repository: task.repository, worker_result: task.worker_input, verifier_result: task.verifier_input })),
    acceptance_criteria: ["Every approved plan task is independently verified.", "The cumulative repository worktrees pass holistic verification."],
    verification_commands: [...new Set(items.flatMap((item) => item.verification_commands))], instruction_paths: [join(workspaceRoot, "AGENTS.md"), join(workspaceRoot, "agents", "verifier.md")],
    result_contract: join(workspaceRoot, ".agents", "contracts", "plan-verifier-result.schema.json"), result_path: planVerifierResult,
  });
  const revision: PlanRuntimeRevision = {
    contract_version: 1, kind: "approved-plan-runtime-revision", plan_reference: index.plan_reference!, plan_id: index.plan_id, run_id: options.runId,
    prior_plan_version: manifest.plan_version, prior_approved_digest: manifest.approved_digest, plan_version: index.plan_version, approved_digest: index.approved_digest!, plan_revision: index.plan_version,
    reason: options.reason.trim(), changed_task_ids: [...changed].sort(), added_task_ids: [...added].sort(), removed_task_ids: removed.sort(), invalidated_task_ids: [...invalidated].sort(), preserved_task_ids: preserved.sort(), prior_manifest: manifestPath, created_at: now,
  };
  const revisionErrors = contractMessages(await validateContract("plan-runtime-revision", revision));
  if (revisionErrors.length > 0) throw new Error(`Invalid plan-runtime-revision: ${revisionErrors.join("; ")}`);
  await writeJsonAtomic(planBriefPath, planBrief);
  await writeJsonAtomic(revisionPath, revision);
  const { plan_verifier_result: _stalePlanVerifier, ...manifestWithoutStaleFinal } = manifest;
  const nextManifest: RuntimeManifest = {
    ...manifestWithoutStaleFinal,
    task_brief: planBriefPath,
    plan_version: index.plan_version,
    plan_revision: index.plan_version,
    approved_digest: index.approved_digest!,
    plan_work_items: planWorkItems,
    task_graph: taskGraph,
    plan_verifier_input: planVerifierInput,
    plan_verifier_status: "pending",
    plan_revisions: [...(manifest.plan_revisions ?? []), revisionPath],
    status: taskGraph.every((task) => task.outcome === "passed") ? "verifying" : "prepared",
    updated_at: now,
    evidence: [...new Set([...manifest.evidence, planBriefPath, revisionPath, planVerifierInput, ...taskGraph.flatMap((task) => [task.worker_input, task.verifier_input])])],
    repositories: manifest.repositories.map((repository) => {
      const next = { ...repository, status: (taskGraph.filter((task) => task.repository === repository.name).every((task) => task.outcome === "passed") ? "passed" : taskGraph.some((task) => task.repository === repository.name && task.status === "prepared") ? "prepared" : "waiting") as NonNullable<RuntimeRepository["status"]> };
      delete next.review_preparation;
      delete next.review_publication;
      delete next.review_state;
      delete next.merge_confirmation;
      delete next.closeout_record;
      return next;
    }),
  };
  const nextErrors = contractMessages(await validateContract("runtime-manifest", nextManifest));
  if (nextErrors.length > 0) throw new Error(`Revised runtime-manifest is invalid: ${nextErrors.join("; ")}`);
  await writeJsonAtomic(manifestPath, nextManifest);
  await setPlanState(planDirectory, { kind: "lifecycle", status: "in-progress", reason: `Approved plan revision ${index.plan_version} resumed in the existing runtime.`, actor: "engine", evidence: revisionPath }, options.now ?? new Date());
  return {
    planId: index.plan_id, planReference: index.plan_reference!, planVersion: index.plan_version, approvedDigest: index.approved_digest!, runId: options.runId,
    manifest: manifestPath, planBrief: planBriefPath, repositories: manifest.repositories.map((repository) => ({ name: repository.name, branch: repository.branch, worktree: repository.worktree, ready: repository.status === "prepared", blockedBy: [], taskInputs: repository.task_inputs ?? [], verifierInputs: repository.verifier_inputs ?? [] })), preparationStatus: "prepared",
  };
}
