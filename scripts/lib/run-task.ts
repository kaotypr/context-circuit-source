import { access, readFile } from "node:fs/promises";
import { join, relative, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { assertCleanRepository, git } from "./git.js";
import { generateIds, slugify } from "./ids.js";
import { assertInside, ensurePrivateDirectory, writeJsonAtomic } from "./io.js";
import type { RuntimeManifest, TaskBrief, TestExpectation, TestExpectationPolicy, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";

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
}

async function assertValid(name: "workspace" | "task-brief" | "runtime-manifest", value: unknown): Promise<void> {
  const errors = await validateContract(name, value);
  if (errors.length > 0) {
    throw new Error(`Generated ${name} is invalid: ${errors.map((error) => `${error.instancePath} ${error.message}`).join("; ")}`);
  }
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

  const runtimeRepository = {
    name: options.repository,
    base_path: relative(workspaceRoot, repositoryPath),
    base_commit: baseCommit,
    branch,
    worktree,
    worker_input: workerInputPath,
    verifier_input: verifierInputPath,
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
    warnings: ["No activity tool is configured; this run cannot guarantee exclusive ownership."],
    execution_events: [],
    lifecycle_events: [{
      event: "task.starting",
      status: "skipped",
      idempotency_key: `${runId}:task.starting:activity-none`,
      occurred_at: createdAt,
    }],
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

  try {
    await ensurePrivateDirectory(join(runtimeRoot, "worktrees", runId));
    await git(repositoryPath, ["worktree", "add", "-b", branch, worktree, baseCommit]);
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

  return { workId, runId, branch, worktree, taskBrief: taskBriefPath, manifest: manifestPath, workerInput: workerInputPath, verifierInput: verifierInputPath };
}
