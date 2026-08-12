import { randomBytes } from "node:crypto";
import { access, readFile } from "node:fs/promises";
import { join, relative, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { assertCleanRepository, git } from "./git.js";
import { generateIds, generateRunId, slugify } from "./ids.js";
import { assertInside, ensurePrivateDirectory, writeJsonAtomic } from "./io.js";
import type { PlanRunTaskRequest, RunTaskRequest, RuntimeManifest, RuntimeRepository, TaskBrief, TaskRepositoryTarget, TestExpectation, TestExpectationPolicy, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";
import { prepareActivityLifecycle } from "./activity-lifecycle.js";
import type { ActivityCapability } from "./types.js";
import { validatePlanDirectory } from "./plans.js";

export interface PrepareTaskOptions {
  workspaceRoot: string;
  request: string;
  repository: string;
  acceptanceCriteria: string[];
  scope: string[];
  testScope?: string[];
  testPolicy?: TestExpectationPolicy;
  testRationale?: string;
  verificationCommands?: string[];
  now?: Date;
  discriminator?: string;
  availableCapabilities?: ActivityCapability[];
}

export interface PreparedTask {
  workId: string;
  runId: string;
  branch: string;
  worktree: string;
  taskBrief: string;
  manifest: string;
  workerInput: string;
  verifierInput: string;
  repositories: PreparedTaskRepository[];
  preparationStatus: "prepared" | "awaiting-activity" | "blocked";
}

export interface PreparedTaskRepository {
  name: string;
  branch: string;
  worktree: string;
  workerInput: string;
  verifierInput: string;
  ready: boolean;
  blockedBy: string[];
}

export interface PrepareContractFirstTaskOptions {
  workspaceRoot: string;
  request: RunTaskRequest;
  now?: Date;
  discriminator?: string;
  availableCapabilities?: ActivityCapability[];
}

export interface PreparePlanTaskOptions {
  workspaceRoot: string;
  request: PlanRunTaskRequest;
  now?: Date;
  discriminator?: string;
  availableCapabilities?: ActivityCapability[];
}

export interface ResumeTaskOptions { workspaceRoot: string; runId: string; now?: Date }

async function assertValid(name: "workspace" | "task-brief" | "runtime-manifest", value: unknown): Promise<void> {
  const errors = await validateContract(name, value);
  if (errors.length > 0) {
    throw new Error(`Generated ${name} is invalid: ${errors.map((error) => `${error.instancePath} ${error.message}`).join("; ")}`);
  }
}

function repositoryExpectation(input: RunTaskRequest["repositories"][number]): TestExpectation {
  const paths = normalizeScope(input.test_scope, `test scope for ${input.name}`);
  if ((input.test_policy === "required" || input.test_policy === "existing-coverage") && paths.length === 0) {
    throw new Error(`${input.test_policy} test policy requires at least one test scope entry for ${input.name}`);
  }
  return { policy: input.test_policy, paths, rationale: input.test_rationale?.trim() || defaultTestRationale(input.test_policy) };
}

function normalizeContractFirstRequest(input: RunTaskRequest, workId: string, runId: string, createdAt: string): TaskBrief {
  const request = input.request.trim();
  const acceptanceCriteria = input.acceptance_criteria.map((item) => item.trim()).filter(Boolean);
  if (!request) throw new Error("A direct request is required");
  if (acceptanceCriteria.length === 0) throw new Error("At least one acceptance criterion is required");
  const names = input.repositories.map((repository) => repository.name);
  if (new Set(names).size !== names.length) throw new Error("Repository names must be unique");
  if (!names.includes(input.shared_contract.repository)) throw new Error("Shared contract repository must be included in repositories");
  const contractPaths = normalizeScope(input.shared_contract.paths, "shared contract paths");
  const byName = new Map(input.repositories.map((repository) => [repository.name, repository]));
  for (const repository of input.repositories) {
    if (repository.depends_on.includes(repository.name)) throw new Error(`${repository.name} cannot depend on itself`);
    for (const dependency of repository.depends_on) if (!byName.has(dependency)) throw new Error(`${repository.name} has unknown dependency ${dependency}`);
  }
  const contractOwner = byName.get(input.shared_contract.repository)!;
  if (contractOwner.depends_on.length > 0) throw new Error("Shared contract repository cannot depend on another repository");
  const visiting = new Set<string>();
  const orders = new Map<string, number>();
  const orderOf = (name: string): number => {
    const known = orders.get(name);
    if (known !== undefined) return known;
    if (visiting.has(name)) throw new Error(`Repository dependency cycle includes ${name}`);
    visiting.add(name);
    const repository = byName.get(name)!;
    const order = repository.depends_on.length === 0 ? 0 : Math.max(...repository.depends_on.map(orderOf)) + 1;
    visiting.delete(name);
    orders.set(name, order);
    return order;
  };
  const dependsOnContract = (name: string, seen = new Set<string>()): boolean => {
    if (name === input.shared_contract.repository) return true;
    if (seen.has(name)) return false;
    seen.add(name);
    return byName.get(name)!.depends_on.some((dependency) => dependsOnContract(dependency, seen));
  };
  for (const name of names) {
    orderOf(name);
    if (name !== input.shared_contract.repository && !dependsOnContract(name)) throw new Error(`${name} must depend on the shared contract repository`);
  }
  const targets: TaskRepositoryTarget[] = input.repositories.map((repository) => {
    const implementationScope = normalizeScope(repository.scope, `implementation scope for ${repository.name}`);
    const testExpectation = repositoryExpectation(repository);
    const scope = [...new Set([...implementationScope, ...(testExpectation.policy === "required" ? testExpectation.paths : [])])];
    const repositoryAcceptance = repository.acceptance_criteria.map((criterion) => criterion.trim()).filter(Boolean);
    if (repositoryAcceptance.length === 0) throw new Error(`At least one acceptance criterion is required for ${repository.name}`);
    return {
      name: repository.name,
      dependency_order: orders.get(repository.name)!,
      depends_on: repository.depends_on,
      scope,
      implementation_scope: implementationScope,
      test_expectation: testExpectation,
      verification_commands: repository.verification_commands.map((command) => command.trim()).filter(Boolean),
      acceptance_criteria: repositoryAcceptance,
    };
  }).sort((left, right) => left.dependency_order - right.dependency_order || left.name.localeCompare(right.name));
  const ownerScope = targets.find((target) => target.name === input.shared_contract.repository)!.scope!;
  for (const path of contractPaths) if (!ownerScope.some((scope) => path === scope || path.startsWith(`${scope}/`))) {
    throw new Error(`Shared contract path is outside ${input.shared_contract.repository} scope: ${path}`);
  }
  return {
    contract_version: 1,
    work_id: workId,
    run_id: runId,
    source: { kind: "direct-request" },
    requested_outcome: request,
    scope: [...new Set(targets.flatMap((target) => target.scope!))],
    acceptance_criteria: acceptanceCriteria,
    repositories: targets,
    shared_contract: { repository: input.shared_contract.repository, paths: contractPaths },
    plan: { reference: null, approval_state: "not-applicable" },
    activity: { reference: null, claim_status: "not-applicable", duplicate_effort_warning: true },
    assumptions: ["The direct request is authoritative for this planless run.", "Dependent repositories remain locked until their declared dependencies pass independent verification."],
    risks: ["No authoritative claim is available; duplicate effort is possible."],
    verification_commands: [...new Set(targets.flatMap((target) => target.verification_commands ?? []))],
    authorization: { kind: "explicit-user-request", evidence: "The human explicitly invoked run-task for this direct request." },
    created_at: createdAt,
  };
}

function normalizePlanRequest(input: PlanRunTaskRequest, runId: string, createdAt: string): TaskBrief {
  const request = input.request.trim();
  const acceptanceCriteria = input.acceptance_criteria.map((item) => item.trim()).filter(Boolean);
  if (!request) throw new Error("A plan execution outcome is required");
  if (acceptanceCriteria.length === 0) throw new Error("At least one acceptance criterion is required");
  const names = input.repositories.map((repository) => repository.name);
  if (new Set(names).size !== names.length) throw new Error("Repository names must be unique");
  const byName = new Map(input.repositories.map((repository) => [repository.name, repository]));
  const visiting = new Set<string>();
  const orders = new Map<string, number>();
  const orderOf = (name: string): number => {
    const known = orders.get(name);
    if (known !== undefined) return known;
    if (visiting.has(name)) throw new Error(`Repository dependency cycle includes ${name}`);
    visiting.add(name);
    const repository = byName.get(name)!;
    for (const dependency of repository.depends_on) if (!byName.has(dependency)) throw new Error(`${name} has unknown dependency ${dependency}`);
    const order = repository.depends_on.length === 0 ? 0 : Math.max(...repository.depends_on.map(orderOf)) + 1;
    visiting.delete(name);
    orders.set(name, order);
    return order;
  };
  const targets: TaskRepositoryTarget[] = input.repositories.map((repository) => {
    if (repository.depends_on.includes(repository.name)) throw new Error(`${repository.name} cannot depend on itself`);
    const implementationScope = normalizeScope(repository.scope, `implementation scope for ${repository.name}`);
    const testExpectation = repositoryExpectation(repository);
    const repositoryAcceptance = repository.acceptance_criteria.map((criterion) => criterion.trim()).filter(Boolean);
    if (repositoryAcceptance.length === 0) throw new Error(`At least one acceptance criterion is required for ${repository.name}`);
    return {
      name: repository.name,
      dependency_order: orderOf(repository.name),
      depends_on: repository.depends_on,
      scope: [...new Set([...implementationScope, ...(testExpectation.policy === "required" ? testExpectation.paths : [])])],
      implementation_scope: implementationScope,
      test_expectation: testExpectation,
      verification_commands: repository.verification_commands.map((command) => command.trim()).filter(Boolean),
      acceptance_criteria: repositoryAcceptance,
    };
  }).sort((left, right) => left.dependency_order - right.dependency_order || left.name.localeCompare(right.name));
  return {
    contract_version: 1,
    work_id: input.work_ids[0]!,
    run_id: runId,
    source: { kind: "plan", reference: input.source.reference },
    requested_outcome: request,
    scope: [...new Set(targets.flatMap((target) => target.scope!))],
    acceptance_criteria: acceptanceCriteria,
    repositories: targets,
    plan: {
      reference: input.source.reference,
      approval_state: "approved",
      plan_version: input.source.plan_version,
      approved_digest: input.source.approved_digest,
      work_ids: input.work_ids,
    },
    activity: { reference: null, claim_status: "not-applicable", duplicate_effort_warning: true },
    assumptions: ["The selected work IDs and approved plan material are authoritative for this run."],
    risks: ["No authoritative claim is available; duplicate effort is possible."],
    verification_commands: [...new Set(targets.flatMap((target) => target.verification_commands ?? []))],
    authorization: { kind: "confirmed-selection", evidence: `The human selected approved plan work: ${input.work_ids.join(", ")}.` },
    created_at: createdAt,
  };
}

export function normalizeDirectRequest(input: {
  request: string;
  repository: string;
  acceptanceCriteria: string[];
  scope: string[];
  testScope?: string[];
  testPolicy?: TestExpectationPolicy;
  testRationale?: string;
  verificationCommands: string[];
  workId: string;
  runId: string;
  createdAt: string;
}): TaskBrief {
  const request = input.request.trim();
  const acceptanceCriteria = input.acceptanceCriteria.map((item) => item.trim()).filter(Boolean);
  const implementationScope = normalizeScope(input.scope, "implementation scope");
  const testScope = normalizeScope(input.testScope ?? [], "test scope");
  const testPolicy = input.testPolicy ?? (testScope.length > 0 ? "required" : "verifier-only");
  if (!request) throw new Error("A direct request is required");
  if (acceptanceCriteria.length === 0) throw new Error("At least one acceptance criterion is required");
  if (implementationScope.length === 0) throw new Error("At least one implementation scope entry is required");
  if ((testPolicy === "required" || testPolicy === "existing-coverage") && testScope.length === 0) {
    throw new Error(`${testPolicy} test policy requires at least one test scope entry`);
  }
  const testExpectation: TestExpectation = {
    policy: testPolicy,
    paths: testScope,
    rationale: input.testRationale?.trim() || defaultTestRationale(testPolicy),
  };
  const scope = [...new Set([...implementationScope, ...(testPolicy === "required" ? testScope : [])])];
  return {
    contract_version: 1,
    work_id: input.workId,
    run_id: input.runId,
    source: { kind: "direct-request" },
    requested_outcome: request,
    scope,
    implementation_scope: implementationScope,
    test_expectation: testExpectation,
    acceptance_criteria: acceptanceCriteria,
    repositories: [{ name: input.repository, dependency_order: 0 }],
    plan: { reference: null, approval_state: "not-applicable" },
    activity: {
      reference: null,
      claim_status: "not-applicable",
      duplicate_effort_warning: true,
    },
    assumptions: ["The direct request is authoritative for this planless run."],
    risks: ["No authoritative claim is available; duplicate effort is possible."],
    verification_commands: input.verificationCommands,
    authorization: {
      kind: "explicit-user-request",
      evidence: "The human explicitly invoked run-task for this direct request.",
    },
    created_at: input.createdAt,
  };
}

function normalizeScope(values: string[], label: string): string[] {
  const normalized = values.map((item) => item.trim().replace(/\/$/, "")).filter(Boolean);
  for (const path of normalized) {
    if (path.startsWith("/") || path.includes("\\") || path.split("/").includes("..")) {
      throw new Error(`${label} entries must be repository-relative paths: ${path}`);
    }
  }
  return [...new Set(normalized)];
}

function defaultTestRationale(policy: TestExpectationPolicy): string {
  switch (policy) {
    case "required": return "The worker must add or update tests in the declared test scope.";
    case "existing-coverage": return "Declared existing tests are expected to cover the requested behavior.";
    case "not-required": return "No repository test change is required for this task.";
    default: return "No test edit scope is authorized; the verifier must supply independent acceptance evidence.";
  }
}

export async function preparePlanlessTask(options: PrepareTaskOptions): Promise<PreparedTask> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const configPath = join(workspaceRoot, "workspace.yaml");
  const config = parseYaml(await readFile(configPath, "utf8")) as WorkspaceConfig;
  await assertValid("workspace", config);
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace: ${semanticErrors.join("; ")}`);

  const repositoryConfig = config.repositories[options.repository];
  if (!repositoryConfig) throw new Error(`Unknown repository: ${options.repository}`);
  const repositoryPath = assertInside(workspaceRoot, join(workspaceRoot, repositoryConfig.path));
  await access(repositoryPath);
  await assertCleanRepository(repositoryPath);
  const baseCommit = await git(repositoryPath, ["rev-parse", repositoryConfig.default_branch]);

  const requiredInstructionPaths = [
    join(workspaceRoot, "AGENTS.md"),
    join(workspaceRoot, "agents", `${repositoryConfig.agent}.md`),
    join(workspaceRoot, "agents", "repository-worker.md"),
    join(workspaceRoot, "agents", "verifier.md"),
  ];
  await Promise.all(requiredInstructionPaths.map((path) => access(path)));

  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const now = options.now ?? new Date();
  const { workId, runId } = await generateIds(runtimeRoot, options.request, now, options.discriminator);
  const createdAt = now.toISOString();
  const branch = `agent/${workId.toLowerCase()}-${slugify(options.request)}-${runId.slice(-8)}`;
  const runRoot = assertInside(runtimeRoot, join(runtimeRoot, "runs", runId));
  const worktree = assertInside(runtimeRoot, join(runtimeRoot, "worktrees", runId, options.repository));
  const taskBriefPath = join(runtimeRoot, "tasks", `${runId}.json`);
  const manifestPath = join(runRoot, "manifest.json");
  const workerInputPath = join(runRoot, `${options.repository}-worker-input.json`);
  const verifierInputPath = join(runRoot, `${options.repository}-verifier-input.json`);
  const workerResultPath = join(runtimeRoot, "results", `${runId}-${options.repository}-worker.json`);
  const verifierResultPath = join(runtimeRoot, "results", `${runId}-${options.repository}-verifier.json`);

  const taskBrief = normalizeDirectRequest({
    request: options.request,
    repository: options.repository,
    acceptanceCriteria: options.acceptanceCriteria,
    scope: options.scope,
    ...(options.testScope ? { testScope: options.testScope } : {}),
    ...(options.testPolicy ? { testPolicy: options.testPolicy } : {}),
    ...(options.testRationale ? { testRationale: options.testRationale } : {}),
    verificationCommands: options.verificationCommands ?? [],
    workId,
    runId,
    createdAt,
  });
  await assertValid("task-brief", taskBrief);
  if (taskBrief.test_expectation?.policy === "existing-coverage") {
    for (const path of taskBrief.test_expectation.paths) {
      try {
        await access(assertInside(repositoryPath, join(repositoryPath, path)));
      } catch {
        throw new Error(`Existing-coverage test path does not exist: ${path}`);
      }
    }
  }
  await writeJsonAtomic(taskBriefPath, taskBrief);

  const runtimeRepository: RuntimeRepository = {
    name: options.repository,
    base_path: relative(workspaceRoot, repositoryPath),
    base_commit: baseCommit,
    branch,
    worktree,
    worker_input: workerInputPath,
    verifier_input: verifierInputPath,
    repair_attempts: 0,
  };
  const manifest: RuntimeManifest = {
    contract_version: 1,
    work_id: workId,
    run_id: runId,
    status: "preparing",
    created_at: createdAt,
    updated_at: createdAt,
    task_brief: taskBriefPath,
    repositories: [runtimeRepository],
    evidence: [taskBriefPath, manifestPath, workerInputPath, verifierInputPath],
    warnings: config.activity.provider === "none" ? ["No activity tool is configured; this run cannot guarantee exclusive ownership."] : [],
    execution_events: [],
    lifecycle_events: config.activity.provider === "none" ? [{
      event: "task.starting",
      status: "skipped",
      idempotency_key: `${runId}:task.starting:activity-none`,
      occurred_at: createdAt,
    }] : [],
  };
  await assertValid("runtime-manifest", manifest);
  await writeJsonAtomic(manifestPath, manifest);

  const instructionPaths = [
    ...requiredInstructionPaths.slice(0, 3),
  ];
  try {
    await access(join(repositoryPath, "AGENTS.md"));
    instructionPaths.push(join(worktree, "AGENTS.md"));
  } catch {
    // Repository-local instructions are optional when the repository has none.
  }
  const workerInput = {
    contract_version: 1,
    role: "repository-worker",
    task_brief: taskBriefPath,
    repository: options.repository,
    worktree,
    branch,
    base_commit: baseCommit,
    ready: true,
    blocked_by: [],
    allowed_scope: taskBrief.scope,
    implementation_scope: taskBrief.implementation_scope,
    test_expectation: taskBrief.test_expectation,
    instruction_paths: instructionPaths,
    result_contract: join(workspaceRoot, ".agents", "contracts", "worker-result.schema.json"),
    result_path: workerResultPath,
  };
  const verifierInput = {
    contract_version: 1,
    role: "verifier",
    read_only: true,
    task_brief: taskBriefPath,
    repository: options.repository,
    worktree,
    branch,
    base_commit: baseCommit,
    worker_result: workerResultPath,
    acceptance_criteria: taskBrief.acceptance_criteria,
    test_expectation: taskBrief.test_expectation,
    verification_commands: taskBrief.verification_commands,
    instruction_paths: [join(workspaceRoot, "AGENTS.md"), join(workspaceRoot, "agents", "verifier.md"), ...instructionPaths.slice(3)],
    result_contract: join(workspaceRoot, ".agents", "contracts", "verifier-result.schema.json"),
    result_path: verifierResultPath,
  };
  await writeJsonAtomic(workerInputPath, workerInput);
  await writeJsonAtomic(verifierInputPath, verifierInput);

  if (config.activity.provider !== "none") {
    const lifecycle = await prepareActivityLifecycle({
      workspaceRoot,
      runId,
      event: "task.starting",
      availableCapabilities: options.availableCapabilities ?? [],
      now,
    });
    if (lifecycle.status !== "completed" && lifecycle.status !== "skipped") {
      return {
        workId, runId, branch, worktree, taskBrief: taskBriefPath, manifest: manifestPath,
        workerInput: workerInputPath, verifierInput: verifierInputPath,
        repositories: [{ name: options.repository, branch, worktree, workerInput: workerInputPath, verifierInput: verifierInputPath, ready: true, blockedBy: [] }],
        preparationStatus: lifecycle.status === "failed" ? "blocked" : "awaiting-activity",
      };
    }
  }

  try {
    await ensurePrivateDirectory(join(runtimeRoot, "worktrees", runId));
    await git(repositoryPath, ["worktree", "add", "-b", branch, worktree, baseCommit]);
    runtimeRepository.status = "prepared";
    manifest.status = "prepared";
    manifest.updated_at = new Date().toISOString();
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
  } catch (error) {
    manifest.status = "blocked";
    manifest.updated_at = new Date().toISOString();
    manifest.evidence.push(`Preparation failed: ${(error as Error).message}`);
    await writeJsonAtomic(manifestPath, manifest);
    throw error;
  }

  return {
    workId, runId, branch, worktree, taskBrief: taskBriefPath, manifest: manifestPath,
    workerInput: workerInputPath, verifierInput: verifierInputPath,
    repositories: [{ name: options.repository, branch, worktree, workerInput: workerInputPath, verifierInput: verifierInputPath, ready: true, blockedBy: [] }],
    preparationStatus: "prepared",
  };
}

export async function prepareContractFirstTask(options: PrepareContractFirstTaskOptions): Promise<PreparedTask> {
  const requestErrors = await validateContract("run-task-request", options.request);
  if (requestErrors.length > 0) {
    throw new Error(`Invalid run-task-request: ${requestErrors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  }
  const workspaceRoot = resolve(options.workspaceRoot);
  const config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  await assertValid("workspace", config);
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace: ${semanticErrors.join("; ")}`);

  const repositoryBases = new Map<string, { path: string; commit: string }>();
  for (const target of options.request.repositories) {
    const registered = config.repositories[target.name];
    if (!registered) throw new Error(`Unknown repository: ${target.name}`);
    const path = assertInside(workspaceRoot, join(workspaceRoot, registered.path));
    await access(path);
    await assertCleanRepository(path);
    repositoryBases.set(target.name, { path, commit: await git(path, ["rev-parse", registered.default_branch]) });
    await Promise.all([
      access(join(workspaceRoot, "AGENTS.md")),
      access(join(workspaceRoot, "agents", `${registered.agent}.md`)),
      access(join(workspaceRoot, "agents", "repository-worker.md")),
      access(join(workspaceRoot, "agents", "verifier.md")),
    ]);
  }

  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const now = options.now ?? new Date();
  const { workId, runId } = await generateIds(runtimeRoot, options.request.request, now, options.discriminator);
  const createdAt = now.toISOString();
  const branch = `agent/${workId.toLowerCase()}-${slugify(options.request.request)}-${runId.slice(-8)}`;
  const runRoot = assertInside(runtimeRoot, join(runtimeRoot, "runs", runId));
  const taskBriefPath = join(runtimeRoot, "tasks", `${runId}.json`);
  const manifestPath = join(runRoot, "manifest.json");
  const taskBrief = normalizeContractFirstRequest(options.request, workId, runId, createdAt);
  await assertValid("task-brief", taskBrief);

  for (const target of taskBrief.repositories) {
    if (target.test_expectation?.policy !== "existing-coverage") continue;
    const base = repositoryBases.get(target.name)!;
    for (const path of target.test_expectation.paths) {
      try { await access(assertInside(base.path, join(base.path, path))); }
      catch { throw new Error(`Existing-coverage test path does not exist in ${target.name}: ${path}`); }
    }
  }
  await writeJsonAtomic(taskBriefPath, taskBrief);

  const runtimeRepositories: RuntimeRepository[] = [];
  const preparedRepositories: PreparedTaskRepository[] = [];
  for (const target of taskBrief.repositories) {
    const base = repositoryBases.get(target.name)!;
    const registered = config.repositories[target.name]!;
    const worktree = assertInside(runtimeRoot, join(runtimeRoot, "worktrees", runId, target.name));
    const workerInputPath = join(runRoot, `${target.name}-worker-input.json`);
    const verifierInputPath = join(runRoot, `${target.name}-verifier-input.json`);
    const workerResultPath = join(runtimeRoot, "results", `${runId}-${target.name}-worker.json`);
    const verifierResultPath = join(runtimeRoot, "results", `${runId}-${target.name}-verifier.json`);
    const blockedBy = target.depends_on ?? [];
    const ready = blockedBy.length === 0;
    const instructionPaths = [
      join(workspaceRoot, "AGENTS.md"),
      join(workspaceRoot, "agents", `${registered.agent}.md`),
      join(workspaceRoot, "agents", "repository-worker.md"),
    ];
    try { await access(join(base.path, "AGENTS.md")); instructionPaths.push(join(worktree, "AGENTS.md")); } catch { /* optional */ }
    const sharedContract = {
      repository: taskBrief.shared_contract!.repository,
      paths: taskBrief.shared_contract!.paths,
      worktree: assertInside(runtimeRoot, join(runtimeRoot, "worktrees", runId, taskBrief.shared_contract!.repository)),
      approval: target.name === taskBrief.shared_contract!.repository ? "must-pass-independent-verification" : "pending",
    };
    await writeJsonAtomic(workerInputPath, {
      contract_version: 1,
      role: "repository-worker",
      task_brief: taskBriefPath,
      repository: target.name,
      worktree,
      branch,
      base_commit: base.commit,
      ready,
      blocked_by: blockedBy,
      shared_contract: sharedContract,
      allowed_scope: target.scope,
      implementation_scope: target.implementation_scope,
      test_expectation: target.test_expectation,
      instruction_paths: instructionPaths,
      result_contract: join(workspaceRoot, ".agents", "contracts", "worker-result.schema.json"),
      result_path: workerResultPath,
    });
    await writeJsonAtomic(verifierInputPath, {
      contract_version: 1,
      role: "verifier",
      read_only: true,
      task_brief: taskBriefPath,
      repository: target.name,
      worktree,
      branch,
      base_commit: base.commit,
      worker_result: workerResultPath,
      acceptance_criteria: target.acceptance_criteria,
      test_expectation: target.test_expectation,
      verification_commands: target.verification_commands,
      shared_contract: sharedContract,
      instruction_paths: [join(workspaceRoot, "AGENTS.md"), join(workspaceRoot, "agents", "verifier.md")],
      result_contract: join(workspaceRoot, ".agents", "contracts", "verifier-result.schema.json"),
      result_path: verifierResultPath,
    });
    runtimeRepositories.push({
      name: target.name,
      base_path: relative(workspaceRoot, base.path),
      base_commit: base.commit,
      branch,
      worktree,
      worker_input: workerInputPath,
      verifier_input: verifierInputPath,
      status: ready ? "prepared" : "waiting",
      depends_on: blockedBy,
      repair_attempts: 0,
    });
    preparedRepositories.push({ name: target.name, branch, worktree, workerInput: workerInputPath, verifierInput: verifierInputPath, ready, blockedBy });
  }
  const manifest: RuntimeManifest = {
    contract_version: 1,
    work_id: workId,
    run_id: runId,
    status: "preparing",
    created_at: createdAt,
    updated_at: createdAt,
    task_brief: taskBriefPath,
    repositories: runtimeRepositories,
    evidence: [taskBriefPath, manifestPath, ...preparedRepositories.flatMap((repository) => [repository.workerInput, repository.verifierInput])],
    warnings: config.activity.provider === "none" ? ["No activity tool is configured; this run cannot guarantee exclusive ownership."] : [],
    execution_events: [],
    lifecycle_events: config.activity.provider === "none" ? [{ event: "task.starting", status: "skipped", idempotency_key: `${runId}:task.starting:activity-none`, occurred_at: createdAt }] : [],
  };
  await assertValid("runtime-manifest", manifest);
  await writeJsonAtomic(manifestPath, manifest);

  if (config.activity.provider !== "none") {
    const lifecycle = await prepareActivityLifecycle({ workspaceRoot, runId, event: "task.starting", availableCapabilities: options.availableCapabilities ?? [], now });
    if (lifecycle.status !== "completed" && lifecycle.status !== "skipped") {
      const primary = preparedRepositories.find((repository) => repository.name === taskBrief.shared_contract!.repository)!;
      return { workId, runId, branch: primary.branch, worktree: primary.worktree, taskBrief: taskBriefPath, manifest: manifestPath, workerInput: primary.workerInput, verifierInput: primary.verifierInput, repositories: preparedRepositories, preparationStatus: lifecycle.status === "failed" ? "blocked" : "awaiting-activity" };
    }
  }

  try {
    await ensurePrivateDirectory(join(runtimeRoot, "worktrees", runId));
    for (const repository of runtimeRepositories) {
      const base = repositoryBases.get(repository.name)!;
      await git(base.path, ["worktree", "add", "-b", repository.branch, repository.worktree, repository.base_commit]);
    }
    manifest.status = "prepared";
    manifest.updated_at = (options.now ?? new Date()).toISOString();
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
  } catch (error) {
    manifest.status = "blocked";
    manifest.updated_at = (options.now ?? new Date()).toISOString();
    manifest.evidence.push(`Preparation failed: ${(error as Error).message}`);
    await writeJsonAtomic(manifestPath, manifest);
    throw error;
  }
  const primary = preparedRepositories.find((repository) => repository.name === taskBrief.shared_contract!.repository)!;
  return { workId, runId, branch: primary.branch, worktree: primary.worktree, taskBrief: taskBriefPath, manifest: manifestPath, workerInput: primary.workerInput, verifierInput: primary.verifierInput, repositories: preparedRepositories, preparationStatus: "prepared" };
}

export async function preparePlanTask(options: PreparePlanTaskOptions): Promise<PreparedTask> {
  const requestErrors = await validateContract("run-task-request", options.request);
  if (requestErrors.length > 0) {
    throw new Error(`Invalid run-task-request: ${requestErrors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  }
  const workspaceRoot = resolve(options.workspaceRoot);
  const planDirectory = assertInside(join(workspaceRoot, "context", "plans"), resolve(workspaceRoot, options.request.source.reference));
  const validation = await validatePlanDirectory(planDirectory);
  if (!validation.index || !validation.work_breakdown || validation.errors.length > 0) {
    throw new Error(`Plan validation failed: ${validation.errors.join("; ") || "plan metadata is unavailable"}`);
  }
  const index = validation.index;
  const breakdown = validation.work_breakdown;
  if (index.status !== "approved") throw new Error(`Plan ${index.plan_id} is draft; explicit approval is required`);
  if (options.request.source.plan_version !== index.plan_version) {
    throw new Error(`Plan version is stale: requested ${options.request.source.plan_version}, current ${index.plan_version}`);
  }
  if (options.request.source.approved_digest !== index.approved_digest) {
    throw new Error("Plan approval digest is stale or does not match the approved plan material");
  }
  const selected = new Set(options.request.work_ids);
  const positions = new Map(options.request.work_ids.map((workId, position) => [workId, position]));
  const items = new Map(breakdown.items.map((item) => [item.work_id, item]));
  const evidence = new Map<string, string>();
  for (const entry of options.request.dependency_evidence) {
    if (evidence.has(entry.work_id)) throw new Error(`Duplicate dependency evidence for ${entry.work_id}`);
    evidence.set(entry.work_id, entry.evidence.trim());
  }
  for (const workId of options.request.work_ids) {
    const item = items.get(workId);
    if (!item) throw new Error(`Unknown plan work ID: ${workId}`);
    for (const dependency of item.depends_on) {
      if (selected.has(dependency)) {
        if (positions.get(dependency)! >= positions.get(workId)!) throw new Error(`${workId} must be selected after dependency ${dependency}`);
      } else if (!evidence.get(dependency)) {
        throw new Error(`${workId} is dependency-blocked by ${dependency}; confirmed completion evidence is required`);
      }
    }
  }
  for (const workId of evidence.keys()) {
    if (!items.has(workId)) throw new Error(`Dependency evidence references unknown plan work ID: ${workId}`);
    if (selected.has(workId)) throw new Error(`Dependency evidence must not pre-complete selected work ID: ${workId}`);
  }

  const config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  await assertValid("workspace", config);
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace: ${semanticErrors.join("; ")}`);
  const requestedRepositories = new Set(options.request.repositories.map((repository) => repository.name));
  const requestedByName = new Map(options.request.repositories.map((repository) => [repository.name, repository]));
  const repositoryDependsOn = (name: string, dependency: string, seen = new Set<string>()): boolean => {
    if (name === dependency) return true;
    if (seen.has(name)) return false;
    seen.add(name);
    return (requestedByName.get(name)?.depends_on ?? []).some((candidate) => repositoryDependsOn(candidate, dependency, seen));
  };
  for (const workId of options.request.work_ids) {
    const item = items.get(workId)!;
    const repository = item.area;
    if (!config.repositories[repository]) throw new Error(`Plan work ${workId} area is not a registered repository: ${repository}`);
    if (!requestedRepositories.has(repository)) throw new Error(`Plan work ${workId} requires repository ${repository}`);
    for (const dependency of item.depends_on.filter((candidate) => selected.has(candidate))) {
      const dependencyRepository = items.get(dependency)!.area;
      if (!repositoryDependsOn(repository, dependencyRepository)) {
        throw new Error(`Repository ${repository} must depend on ${dependencyRepository} for selected plan dependency ${dependency}`);
      }
    }
  }
  for (const repository of requestedRepositories) {
    if (![...selected].some((workId) => items.get(workId)!.area === repository)) {
      throw new Error(`Repository ${repository} is not affected by the selected plan work`);
    }
  }

  const repositoryBases = new Map<string, { path: string; commit: string }>();
  for (const target of options.request.repositories) {
    const registered = config.repositories[target.name];
    if (!registered) throw new Error(`Unknown repository: ${target.name}`);
    const path = assertInside(workspaceRoot, join(workspaceRoot, registered.path));
    await access(path);
    await assertCleanRepository(path);
    repositoryBases.set(target.name, { path, commit: await git(path, ["rev-parse", registered.default_branch]) });
    await Promise.all([
      access(join(workspaceRoot, "AGENTS.md")),
      access(join(workspaceRoot, "agents", `${registered.agent}.md`)),
      access(join(workspaceRoot, "agents", "repository-worker.md")),
      access(join(workspaceRoot, "agents", "verifier.md")),
    ]);
  }

  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const now = options.now ?? new Date();
  const runId = generateRunId(options.request.request, now, options.discriminator ?? randomBytes(4).toString("hex"));
  const workId = options.request.work_ids[0]!;
  const createdAt = now.toISOString();
  const branch = `agent/${workId.toLowerCase()}-${slugify(options.request.request)}-${runId.slice(-8)}`;
  const runRoot = assertInside(runtimeRoot, join(runtimeRoot, "runs", runId));
  const taskBriefPath = join(runtimeRoot, "tasks", `${runId}.json`);
  const manifestPath = join(runRoot, "manifest.json");
  const taskBrief = normalizePlanRequest(options.request, runId, createdAt);
  await assertValid("task-brief", taskBrief);
  for (const target of taskBrief.repositories) {
    if (target.test_expectation?.policy !== "existing-coverage") continue;
    const base = repositoryBases.get(target.name)!;
    for (const path of target.test_expectation.paths) {
      try { await access(assertInside(base.path, join(base.path, path))); }
      catch { throw new Error(`Existing-coverage test path does not exist in ${target.name}: ${path}`); }
    }
  }
  await writeJsonAtomic(taskBriefPath, taskBrief);

  const runtimeRepositories: RuntimeRepository[] = [];
  const preparedRepositories: PreparedTaskRepository[] = [];
  for (const target of taskBrief.repositories) {
    const base = repositoryBases.get(target.name)!;
    const registered = config.repositories[target.name]!;
    const worktree = assertInside(runtimeRoot, join(runtimeRoot, "worktrees", runId, target.name));
    const workerInputPath = join(runRoot, `${target.name}-worker-input.json`);
    const verifierInputPath = join(runRoot, `${target.name}-verifier-input.json`);
    const workerResultPath = join(runtimeRoot, "results", `${runId}-${target.name}-worker.json`);
    const verifierResultPath = join(runtimeRoot, "results", `${runId}-${target.name}-verifier.json`);
    const blockedBy = target.depends_on ?? [];
    const ready = blockedBy.length === 0;
    const instructionPaths = [join(workspaceRoot, "AGENTS.md"), join(workspaceRoot, "agents", `${registered.agent}.md`), join(workspaceRoot, "agents", "repository-worker.md")];
    try { await access(join(base.path, "AGENTS.md")); instructionPaths.push(join(worktree, "AGENTS.md")); } catch { /* optional */ }
    await writeJsonAtomic(workerInputPath, {
      contract_version: 1, role: "repository-worker", task_brief: taskBriefPath, repository: target.name, worktree, branch,
      base_commit: base.commit, ready, blocked_by: blockedBy, allowed_scope: target.scope, implementation_scope: target.implementation_scope,
      test_expectation: target.test_expectation, instruction_paths: instructionPaths,
      result_contract: join(workspaceRoot, ".agents", "contracts", "worker-result.schema.json"), result_path: workerResultPath,
    });
    await writeJsonAtomic(verifierInputPath, {
      contract_version: 1, role: "verifier", read_only: true, task_brief: taskBriefPath, repository: target.name, worktree, branch,
      base_commit: base.commit, worker_result: workerResultPath, acceptance_criteria: target.acceptance_criteria,
      test_expectation: target.test_expectation, verification_commands: target.verification_commands,
      instruction_paths: [join(workspaceRoot, "AGENTS.md"), join(workspaceRoot, "agents", "verifier.md")],
      result_contract: join(workspaceRoot, ".agents", "contracts", "verifier-result.schema.json"), result_path: verifierResultPath,
    });
    runtimeRepositories.push({ name: target.name, base_path: relative(workspaceRoot, base.path), base_commit: base.commit, branch, worktree, worker_input: workerInputPath, verifier_input: verifierInputPath, status: ready ? "prepared" : "waiting", depends_on: blockedBy, repair_attempts: 0 });
    preparedRepositories.push({ name: target.name, branch, worktree, workerInput: workerInputPath, verifierInput: verifierInputPath, ready, blockedBy });
  }
  const manifest: RuntimeManifest = {
    contract_version: 1, work_id: workId, run_id: runId, status: "preparing", created_at: createdAt, updated_at: createdAt,
    task_brief: taskBriefPath, repositories: runtimeRepositories,
    plan_work_items: options.request.work_ids.map((selectedWorkId) => ({ work_id: selectedWorkId, repository: items.get(selectedWorkId)!.area, depends_on: items.get(selectedWorkId)!.depends_on, outcome: "pending" })),
    evidence: [taskBriefPath, manifestPath, ...options.request.dependency_evidence.map((entry) => `${entry.work_id}: ${entry.evidence}`), ...preparedRepositories.flatMap((repository) => [repository.workerInput, repository.verifierInput])],
    warnings: config.activity.provider === "none" ? ["No activity tool is configured; this run cannot guarantee exclusive ownership."] : [],
    execution_events: [], lifecycle_events: config.activity.provider === "none" ? [{ event: "task.starting", status: "skipped", idempotency_key: `${runId}:task.starting:activity-none`, occurred_at: createdAt }] : [],
  };
  await assertValid("runtime-manifest", manifest);
  await writeJsonAtomic(manifestPath, manifest);
  if (config.activity.provider !== "none") {
    const lifecycle = await prepareActivityLifecycle({ workspaceRoot, runId, event: "task.starting", availableCapabilities: options.availableCapabilities ?? [], now });
    if (lifecycle.status !== "completed" && lifecycle.status !== "skipped") {
      const primary = preparedRepositories[0]!;
      return { workId, runId, branch, worktree: primary.worktree, taskBrief: taskBriefPath, manifest: manifestPath, workerInput: primary.workerInput, verifierInput: primary.verifierInput, repositories: preparedRepositories, preparationStatus: lifecycle.status === "failed" ? "blocked" : "awaiting-activity" };
    }
  }
  try {
    await ensurePrivateDirectory(join(runtimeRoot, "worktrees", runId));
    for (const repository of runtimeRepositories) {
      await git(repositoryBases.get(repository.name)!.path, ["worktree", "add", "-b", repository.branch, repository.worktree, repository.base_commit]);
    }
    manifest.status = "prepared";
    manifest.updated_at = now.toISOString();
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
  } catch (error) {
    manifest.status = "blocked";
    manifest.updated_at = new Date().toISOString();
    manifest.evidence.push(`Preparation failed: ${(error as Error).message}`);
    await writeJsonAtomic(manifestPath, manifest);
    throw error;
  }
  const primary = preparedRepositories[0]!;
  return { workId, runId, branch, worktree: primary.worktree, taskBrief: taskBriefPath, manifest: manifestPath, workerInput: primary.workerInput, verifierInput: primary.verifierInput, repositories: preparedRepositories, preparationStatus: "prepared" };
}

export async function resumePlanlessTask(options: ResumeTaskOptions): Promise<PreparedTask> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const manifestPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, "manifest.json"));
  const manifest = JSON.parse(await readFile(manifestPath, "utf8")) as RuntimeManifest;
  await assertValid("runtime-manifest", manifest);
  if (manifest.run_id !== options.runId || manifest.status !== "preparing") throw new Error(`Run ${options.runId} is not awaiting preparation`);
  const lifecycle = manifest.lifecycle_events.find((event) => event.event === "task.starting");
  if (!lifecycle || (lifecycle.status !== "completed" && lifecycle.status !== "skipped")) {
    throw new Error(`task.starting lifecycle is ${lifecycle?.status ?? "missing"}; complete required or manual actions before resuming`);
  }
  const config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  await assertValid("workspace", config);
  const basePaths = new Map<string, string>();
  for (const repository of manifest.repositories) {
    const registered = config.repositories[repository.name];
    if (!registered) throw new Error(`Unknown repository: ${repository.name}`);
    const repositoryPath = assertInside(workspaceRoot, join(workspaceRoot, repository.base_path));
    await assertCleanRepository(repositoryPath);
    const currentBase = await git(repositoryPath, ["rev-parse", registered.default_branch]);
    if (currentBase !== repository.base_commit) throw new Error(`Repository base changed during activity preflight for ${repository.name}; prepare a fresh run`);
    basePaths.set(repository.name, repositoryPath);
  }
  try {
    await ensurePrivateDirectory(join(runtimeRoot, "worktrees", options.runId));
    for (const repository of manifest.repositories) {
      await git(basePaths.get(repository.name)!, ["worktree", "add", "-b", repository.branch, repository.worktree, repository.base_commit]);
    }
    manifest.status = "prepared";
    manifest.updated_at = (options.now ?? new Date()).toISOString();
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
  } catch (error) {
    manifest.status = "blocked";
    manifest.updated_at = (options.now ?? new Date()).toISOString();
    manifest.evidence.push(`Preparation failed: ${(error as Error).message}`);
    await writeJsonAtomic(manifestPath, manifest);
    throw error;
  }
  const brief = JSON.parse(await readFile(manifest.task_brief, "utf8")) as TaskBrief;
  const primary = manifest.repositories.find((repository) => repository.name === brief.shared_contract?.repository) ?? manifest.repositories[0]!;
  const repositories = manifest.repositories.map((repository) => ({
    name: repository.name,
    branch: repository.branch,
    worktree: repository.worktree,
    workerInput: repository.worker_input,
    verifierInput: repository.verifier_input,
    ready: (repository.depends_on?.length ?? 0) === 0,
    blockedBy: repository.depends_on ?? [],
  }));
  return {
    workId: manifest.work_id,
    runId: manifest.run_id,
    branch: primary.branch,
    worktree: primary.worktree,
    taskBrief: manifest.task_brief,
    manifest: manifestPath,
    workerInput: primary.worker_input,
    verifierInput: primary.verifier_input,
    repositories,
    preparationStatus: "prepared",
  };
}
