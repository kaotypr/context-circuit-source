import { createHash } from "node:crypto";
import { access, lstat, readFile, realpath } from "node:fs/promises";
import { join, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { contributionDocumentErrors } from "./finish-work.js";
import { assertCleanRepository, git } from "./git.js";
import { assertInside, ensurePrivateDirectory, withExclusiveFile, writeJsonAtomic } from "./io.js";
import type { ContextSyncRecord, ContextSyncRequest, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";

export interface PrepareContextSyncOptions { workspaceRoot: string; request: ContextSyncRequest; now?: Date }
export interface PrepareContextReviewOptions { workspaceRoot: string; syncId: string; now?: Date }

async function assertValid(name: "workspace" | "context-sync-request" | "context-sync-record", value: unknown): Promise<void> {
  const errors = await validateContract(name, value);
  if (errors.length > 0) throw new Error(`Invalid ${name}: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}

async function readJson<T>(path: string): Promise<T> {
  return JSON.parse(await readFile(path, "utf8")) as T;
}

function compactTimestamp(date: Date): string {
  return date.toISOString().replace(/[-:]/g, "").replace(/\.\d{3}Z$/, "Z");
}

async function loadWorkspace(workspaceRoot: string): Promise<WorkspaceConfig> {
  const config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  await assertValid("workspace", config);
  const errors = workspaceSemanticErrors(config);
  if (errors.length > 0) throw new Error(`Invalid workspace: ${errors.join("; ")}`);
  return config;
}

async function validateRequestSemantics(workspaceRoot: string, config: WorkspaceConfig, request: ContextSyncRequest): Promise<string[]> {
  const listed = new Set(request.contributions);
  for (const contribution of request.contributions) {
    const path = assertInside(workspaceRoot, join(workspaceRoot, contribution));
    const info = await lstat(path);
    if (!info.isFile() || info.isSymbolicLink()) throw new Error(`Contribution must be a regular file: ${contribution}`);
    assertInside(await realpath(workspaceRoot), await realpath(path));
    const errors = contributionDocumentErrors(path, await readFile(path, "utf8"));
    if (errors.length > 0) throw new Error(`Invalid contribution ${contribution}: ${errors.join("; ")}`);
  }
  const durableTargets: string[] = [];
  for (const proposal of request.proposals) {
    if (!listed.has(proposal.source_contribution)) throw new Error(`Proposal source is not listed: ${proposal.source_contribution}`);
    if (proposal.classification === "durable-wrapper") {
      if (!proposal.target || !proposal.proposed_change?.trim() || proposal.target_repository) throw new Error("durable-wrapper proposals require target and proposed_change only");
      if ((await readFile(join(workspaceRoot, proposal.target), "utf8")).includes(proposal.source_contribution)) {
        throw new Error(`Contribution is already cited by ${proposal.target}: ${proposal.source_contribution}`);
      }
      durableTargets.push(proposal.target);
    } else if (proposal.classification === "repository-local") {
      if (!proposal.target_repository || !config.repositories[proposal.target_repository] || proposal.target || proposal.proposed_change) {
        throw new Error("repository-local proposals require one registered target_repository and no wrapper target/change");
      }
    } else if (proposal.target || proposal.target_repository || proposal.proposed_change) {
      throw new Error(`${proposal.classification} proposals cannot mutate wrapper or repository context`);
    }
  }
  if (durableTargets.length === 0) throw new Error("Context synchronization has no durable wrapper proposal; report classifications without creating a worktree");
  return [...new Set(durableTargets)].sort();
}

export async function prepareContextSync(options: PrepareContextSyncOptions): Promise<ContextSyncRecord> {
  const workspaceRoot = resolve(options.workspaceRoot);
  await assertValid("context-sync-request", options.request);
  const config = await loadWorkspace(workspaceRoot);
  const allowedPaths = await validateRequestSemantics(workspaceRoot, config, options.request);
  if (await realpath(await git(workspaceRoot, ["rev-parse", "--show-toplevel"])) !== await realpath(workspaceRoot)) {
    throw new Error("Workspace root must be the wrapper Git root");
  }
  await assertCleanRepository(workspaceRoot);
  const baseBranch = await git(workspaceRoot, ["branch", "--show-current"]);
  if (!baseBranch) throw new Error("Context synchronization requires an attached wrapper branch");
  const baseCommit = await git(workspaceRoot, ["rev-parse", "HEAD"]);
  const now = options.now ?? new Date();
  const digest = createHash("sha256").update(JSON.stringify(options.request)).update(now.toISOString()).digest("hex").slice(0, 8);
  const syncId = `${compactTimestamp(now)}-${digest}`;
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const syncRoot = assertInside(runtimeRoot, join(runtimeRoot, "context-sync", syncId));
  await ensurePrivateDirectory(syncRoot);
  const requestPath = join(syncRoot, "request.json");
  const recordPath = join(syncRoot, "record.json");
  const worktree = assertInside(runtimeRoot, join(runtimeRoot, "worktrees", "context-sync", syncId, "wrapper"));
  const branch = `agent/context-sync-${syncId.toLowerCase()}`;
  await writeJsonAtomic(requestPath, options.request);
  const preparedAt = now.toISOString();
  const record: ContextSyncRecord = {
    contract_version: 1,
    sync_id: syncId,
    status: "prepared",
    wrapper_mode: config.workflow.wrapper_change_policy,
    base_branch: baseBranch,
    base_commit: baseCommit,
    branch,
    worktree,
    request: requestPath,
    allowed_wrapper_paths: allowedPaths,
    repository_follow_ups: options.request.proposals.filter((proposal) => proposal.classification === "repository-local").map((proposal) => ({ repository: proposal.target_repository!, summary: proposal.summary, source_contribution: proposal.source_contribution })),
    future_tasks: options.request.proposals.filter((proposal) => proposal.classification === "future-task").map((proposal) => ({ summary: proposal.summary, source_contribution: proposal.source_contribution })),
    retained_one_offs: options.request.proposals.filter((proposal) => proposal.classification === "one-off").map((proposal) => proposal.summary),
    changed_files: [],
    commits: [],
    remote: null,
    blockers: [],
    prepared_at: preparedAt,
    updated_at: preparedAt,
  };
  try {
    await ensurePrivateDirectory(join(runtimeRoot, "worktrees"));
    await ensurePrivateDirectory(join(runtimeRoot, "worktrees", "context-sync"));
    await ensurePrivateDirectory(join(runtimeRoot, "worktrees", "context-sync", syncId));
    await git(workspaceRoot, ["worktree", "add", "-b", branch, worktree, baseCommit]);
  } catch (error) {
    record.status = "blocked";
    record.blockers = [`Wrapper worktree preparation failed: ${(error as Error).message}`];
    await writeJsonAtomic(recordPath, record);
    throw error;
  }
  await assertValid("context-sync-record", record);
  await writeJsonAtomic(recordPath, record);
  return record;
}

export async function prepareContextReview(options: PrepareContextReviewOptions): Promise<ContextSyncRecord> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const recordPath = assertInside(runtimeRoot, join(runtimeRoot, "context-sync", options.syncId, "record.json"));
  return withExclusiveFile(`${recordPath}.lock`, async () => {
    const record = await readJson<ContextSyncRecord>(recordPath);
    await assertValid("context-sync-record", record);
    if (record.sync_id !== options.syncId) throw new Error("Context sync ID mismatch");
    const worktree = assertInside(runtimeRoot, record.worktree);
    const config = await loadWorkspace(workspaceRoot);
    const request = await readJson<ContextSyncRequest>(assertInside(runtimeRoot, record.request));
    await assertValid("context-sync-request", request);
    const recomputedPaths = await validateRequestSemantics(workspaceRoot, config, request);
    if (recomputedPaths.join("\n") !== record.allowed_wrapper_paths.slice().sort().join("\n")) throw new Error("Context sync allowed paths do not match the validated request");
    const registrations = await git(workspaceRoot, ["worktree", "list", "--porcelain"]);
    if (!registrations.split("\n").includes(`worktree ${await realpath(worktree)}`)) throw new Error("Context sync worktree is not registered by the wrapper repository");
    await assertCleanRepository(worktree);
    if (await git(worktree, ["branch", "--show-current"]) !== record.branch) throw new Error("Wrapper worktree branch changed");
    const head = await git(worktree, ["rev-parse", "HEAD"]);
    await git(worktree, ["merge-base", "--is-ancestor", record.base_commit, head]);
    const commits = (await git(record.worktree, ["rev-list", "--reverse", `${record.base_commit}..${head}`])).split("\n").filter(Boolean);
    const changedFiles = (await git(record.worktree, ["diff", "--name-only", `${record.base_commit}...${head}`])).split("\n").filter(Boolean);
    const outsideScope = changedFiles.filter((path) => !record.allowed_wrapper_paths.includes(path));
    if (outsideScope.length > 0) throw new Error(`Context sync changed files outside approved wrapper scope: ${outsideScope.join(", ")}`);
    if (commits.length === 0 || changedFiles.length === 0) throw new Error("Context review requires committed canonical-context changes");
    for (const proposal of request.proposals.filter((candidate) => candidate.classification === "durable-wrapper")) {
      const content = await readFile(join(record.worktree, proposal.target!), "utf8");
      if (!content.includes(proposal.source_contribution)) throw new Error(`Canonical update must cite source contribution: ${proposal.source_contribution}`);
    }
    const remotes = (await git(record.worktree, ["remote"])).split("\n").filter(Boolean);
    const remote = remotes.includes("origin") ? "origin" : null;
    const blockers = record.wrapper_mode === "pull-request" && !remote ? ["Wrapper has no origin remote; configure one before pushing or opening the required review."] : [];
    record.status = blockers.length > 0 ? "blocked" : "review-ready";
    record.commits = commits;
    record.changed_files = changedFiles;
    record.remote = remote;
    record.blockers = blockers;
    record.updated_at = (options.now ?? new Date()).toISOString();
    await assertValid("context-sync-record", record);
    await writeJsonAtomic(recordPath, record);
    return record;
  });
}
