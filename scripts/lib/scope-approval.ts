import { join, resolve } from "node:path";
import { assertInside, ensurePrivateDirectory, readJsonRegularInside, withExclusiveFile, writeJsonAtomic } from "./io.js";
import type { RuntimeManifest, RuntimeRepository, ScopeApproval, ScopeAuthorization } from "./types.js";
import { validateContract } from "./validation.js";

export interface ApproveScopeExpansionOptions {
  workspaceRoot: string;
  runId: string;
  repository: string;
  taskId?: string;
  approvedBy: string;
  reason: string;
  now?: Date;
}

export interface ApprovedScopeExpansion {
  approval: ScopeApproval;
  approval_path: string;
  worker_input: string;
  verifier_input: string;
}

async function assertValid(name: "runtime-manifest" | "scope-approval", value: unknown): Promise<void> {
  const errors = await validateContract(name, value);
  if (errors.length > 0) throw new Error(`Invalid ${name}: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}

function safeSingleLine(value: string, label: string): string {
  const clean = value.trim();
  if (!clean || /[\r\n]/.test(clean)) throw new Error(`${label} must be a non-empty single line`);
  if (/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----|https?:\/\/[^\s/@:]+:[^\s/@]+@/i.test(clean)) {
    throw new Error(`${label} appears to contain a credential or private key`);
  }
  return clean;
}

function findRepository(manifest: RuntimeManifest, name: string): RuntimeRepository {
  const repository = manifest.repositories.find((candidate) => candidate.name === name);
  if (!repository) throw new Error(`Run ${manifest.run_id} has no repository named ${name}`);
  return repository;
}

function sameAuthorization(left: ScopeAuthorization | undefined, right: ScopeAuthorization): boolean {
  return left?.mode === right.mode && left.approval_path === right.approval_path;
}

function taskIdentity(manifest: RuntimeManifest, repository: RuntimeRepository, requestedTaskId?: string): { taskId: string; workerInput: string; verifierInput: string; ready: boolean; status: string } {
  if (!manifest.task_graph) {
    if (requestedTaskId && requestedTaskId !== manifest.work_id) throw new Error(`Direct task ID must match ${manifest.work_id}`);
    return {
      taskId: manifest.work_id,
      workerInput: repository.worker_input,
      verifierInput: repository.verifier_input,
      ready: repository.status !== "waiting",
      status: repository.status ?? manifest.status,
    };
  }
  if (!requestedTaskId) throw new Error("Plan scope expansion requires --task-id");
  const task = manifest.task_graph.find((candidate) => (candidate.task_id ?? candidate.work_id) === requestedTaskId);
  if (!task) throw new Error(`Run ${manifest.run_id} has no plan task named ${requestedTaskId}`);
  if (task.repository !== repository.name) throw new Error(`Plan task ${requestedTaskId} does not belong to ${repository.name}`);
  return {
    taskId: requestedTaskId,
    workerInput: task.worker_input,
    verifierInput: task.verifier_input,
    ready: task.ready === true,
    status: task.status ?? "waiting",
  };
}

export async function approveScopeExpansion(options: ApproveScopeExpansionOptions): Promise<ApprovedScopeExpansion> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const runRoot = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId));
  const manifestPath = assertInside(runRoot, join(runRoot, "manifest.json"));
  const approvedBy = safeSingleLine(options.approvedBy, "Approver");
  const reason = safeSingleLine(options.reason, "Reason");

  return withExclusiveFile(`${manifestPath}.lock`, async () => {
    const manifest = await readJsonRegularInside<RuntimeManifest>(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    const repository = findRepository(manifest, options.repository);
    const identity = taskIdentity(manifest, repository, options.taskId);
    if (!identity.ready) throw new Error(`Task ${identity.taskId} is not ready for scope expansion approval`);
    if (!["prepared", "running", "blocked"].includes(identity.status)) {
      throw new Error(`Task ${identity.taskId} cannot receive scope expansion approval from status ${identity.status}`);
    }

    const approvalDirectory = assertInside(runRoot, join(runRoot, "scope-approvals"));
    await ensurePrivateDirectory(approvalDirectory);
    const approvalPath = assertInside(approvalDirectory, join(approvalDirectory, `${repository.name}-${identity.taskId}.json`));
    const approvedAt = (options.now ?? new Date()).toISOString();
    const approval: ScopeApproval = {
      contract_version: 1,
      scope_approval_id: `${manifest.run_id}:scope-expansion:${repository.name}:${identity.taskId}`,
      work_id: manifest.work_id,
      run_id: manifest.run_id,
      task_id: identity.taskId,
      repository: repository.name,
      mode: "task-level-expansion",
      reason,
      approved_by: approvedBy,
      approved_at: approvedAt,
    };
    let approvalToWrite = approval;
    try {
      const existing = await readJsonRegularInside<ScopeApproval>(runtimeRoot, approvalPath, "Scope approval");
      await assertValid("scope-approval", existing);
      if (JSON.stringify({ ...existing, approved_at: null }) !== JSON.stringify({ ...approval, approved_at: null })) {
        throw new Error("Scope expansion was already approved with different evidence");
      }
      approvalToWrite = existing;
    } catch (error) {
      if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
    }
    await assertValid("scope-approval", approval);
    await writeJsonAtomic(approvalPath, approvalToWrite);

    const authorization: ScopeAuthorization = { mode: "task-level-expansion", approval_path: approvalPath };
    for (const inputPath of [identity.workerInput, identity.verifierInput]) {
      const absolute = assertInside(runtimeRoot, inputPath);
      const input = await readJsonRegularInside<Record<string, unknown>>(runtimeRoot, absolute, "Task input");
      const current = input.scope_authorization as ScopeAuthorization | undefined;
      if (current && !sameAuthorization(current, authorization)) throw new Error("Task input already has a different scope authorization");
      input.scope_authorization = authorization;
      await writeJsonAtomic(absolute, input);
    }

    manifest.scope_approvals ??= [];
    if (!manifest.scope_approvals.includes(approvalPath)) manifest.scope_approvals.push(approvalPath);
    if (!manifest.evidence.includes(approvalPath)) manifest.evidence.push(approvalPath);
    manifest.updated_at = approvalToWrite.approved_at;
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    return { approval: approvalToWrite, approval_path: approvalPath, worker_input: identity.workerInput, verifier_input: identity.verifierInput };
  });
}
