import { access, lstat, mkdir, readFile, readdir, realpath } from "node:fs/promises";
import { basename, join, relative, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { git } from "./git.js";
import { assertInside, ensurePrivateDirectory, readJsonRegularInside, withExclusiveFile, writeJsonAtomic, writeTextExclusive } from "./io.js";
import type { CloseoutRecord, ExecutionEvent, MergeConfirmationRecord, RuntimeManifest, RuntimeRepository, TaskBrief, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";

interface ResultInput {
  result_path: string;
}

interface VerifierResult {
  work_id: string;
  run_id: string;
  repository: string;
  status: "pass" | "fail" | "blocked";
  summary: string;
  checks: string[];
  acceptance: Array<{ criterion: string; status: string; evidence: string }>;
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

export interface FinishWorkOptions {
  workspaceRoot: string;
  runId: string;
  repository: string;
  outcome: "merged" | "abandoned";
  author: string;
  reason?: string;
  mergeCommit?: string;
  pullRequests?: string[];
  cleanup?: boolean;
  now?: Date;
}

const contributionHeadings = [
  "## Outcome",
  "## Affected repositories",
  "## Pull requests and commits",
  "## Verification",
  "## Decisions and deviations",
  "## Remaining risks and follow-up",
  "## Candidate durable learnings",
] as const;

async function assertValid(name: "workspace" | "runtime-manifest" | "task-brief" | "worker-result" | "verifier-result" | "merge-confirmation-record" | "closeout-record", value: unknown): Promise<void> {
  const errors = await validateContract(name, value);
  if (errors.length > 0) throw new Error(`Invalid ${name}: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}

function findRepository(manifest: RuntimeManifest, name: string): RuntimeRepository {
  const repository = manifest.repositories.find((candidate) => candidate.name === name);
  if (!repository) throw new Error(`Run ${manifest.run_id} has no repository named ${name}`);
  return repository;
}

function compactTimestamp(date: Date): string {
  return date.toISOString().replace(/[-:]/g, "").replace(/\.\d{3}Z$/, "Z");
}

function safeToken(value: string, label: string): string {
  const normalized = value.toLowerCase().trim().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "");
  if (!normalized || !/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(normalized)) throw new Error(`${label} must contain letters or numbers`);
  return normalized;
}

function taskSlug(brief: TaskBrief): string {
  return safeToken(brief.requested_outcome, "Task outcome").slice(0, 48).replace(/-$/, "") || "work";
}

function list(items: string[], empty: string): string {
  return items.length > 0 ? items.map((item) => `- ${item}`).join("\n") : `- ${empty}`;
}

function contributionDocument(
  manifest: RuntimeManifest,
  repository: RuntimeRepository,
  brief: TaskBrief,
  record: Pick<CloseoutRecord, "outcome" | "author" | "reason" | "pull_requests" | "commits" | "changed_files" | "verification" | "head_commit">,
): string {
  const outcome = record.outcome === "merged" ? "Merged after human review." : `Deliberately abandoned by the human.${record.reason ? ` ${record.reason}` : ""}`;
  const changed = record.changed_files.length > 0 ? ` Changed files: ${record.changed_files.join(", ")}.` : " No product files changed.";
  return `# ${manifest.work_id}: ${brief.requested_outcome}\n\n` +
    `- Run: \`${manifest.run_id}\`\n` +
    `- Task source: direct request\n` +
    `- Plan: none\n` +
    `- Author: \`${record.author}\`\n\n` +
    `## Outcome\n\n${outcome}\n\n` +
    `## Affected repositories\n\n- \`${repository.name}\` on branch \`${repository.branch}\`.${changed}\n\n` +
    `## Pull requests and commits\n\n${list(record.pull_requests.map((item) => `Pull request: ${item}`), "No pull-request reference was recorded.")}\n${list(record.commits.map((item) => `Commit: \`${item}\``), `No commits beyond base \`${repository.base_commit}\`.`)}\n- Recorded head: \`${record.head_commit}\`\n\n` +
    `## Verification\n\n${list(record.verification, "No verifier evidence was available.")}\n\n` +
    `## Decisions and deviations\n\n- ${record.outcome === "merged" ? "No closeout deviation was recorded." : "The run was deliberately abandoned instead of merged."}\n\n` +
    `## Remaining risks and follow-up\n\n- ${record.reason ?? "No closeout-specific follow-up was recorded."}\n\n` +
    `## Candidate durable learnings\n\n- Review this contribution during the next context synchronization; no canonical-context change is asserted automatically.\n`;
}

export function contributionDocumentErrors(path: string, content: string, runId?: string): string[] {
  const errors: string[] = [];
  if (!/^\d{8}T\d{6}Z-[a-z0-9]+(?:-[a-z0-9]+)*-[a-z0-9]+(?:-[a-z0-9]+)*\.md$/.test(basename(path))) {
    errors.push("contribution filename must be <UTC timestamp>-<author>-<slug>.md");
  }
  for (const heading of contributionHeadings) if (!content.includes(`${heading}\n`)) errors.push(`contribution is missing ${heading}`);
  if (runId && !content.includes(`- Run: \`${runId}\``)) errors.push("contribution does not reference the expected run");
  if (/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----|https?:\/\/[^\s/@:]+:[^\s/@]+@/i.test(content)) {
    errors.push("contribution appears to contain a credential or private key");
  }
  return errors;
}

async function loadWorkspace(workspaceRoot: string): Promise<WorkspaceConfig> {
  const config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  await assertValid("workspace", config);
  const errors = workspaceSemanticErrors(config);
  if (errors.length > 0) throw new Error(`Invalid workspace: ${errors.join("; ")}`);
  return config;
}

async function optionalVerifier(runtimeRoot: string, manifest: RuntimeManifest, repository: RuntimeRepository): Promise<VerifierResult | null> {
  try {
    const input = await readJsonRegularInside<ResultInput>(runtimeRoot, repository.verifier_input, "Verifier input");
    const resultPath = assertInside(runtimeRoot, input.result_path);
    const result = await readJsonRegularInside<VerifierResult>(runtimeRoot, resultPath, "Verifier result");
    await assertValid("verifier-result", result);
    if (result.work_id !== manifest.work_id || result.run_id !== manifest.run_id || result.repository !== repository.name) {
      throw new Error("Verifier result identity does not match the closeout run");
    }
    return result;
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code === "ENOENT") return null;
    throw error;
  }
}

async function assertCurrentWorker(runtimeRoot: string, manifest: RuntimeManifest, repository: RuntimeRepository, headCommit: string, commits: string[], changedFiles: string[]): Promise<void> {
  const input = await readJsonRegularInside<ResultInput>(runtimeRoot, repository.worker_input, "Worker input");
  const worker = await readJsonRegularInside<WorkerResult>(runtimeRoot, input.result_path, "Worker result");
  await assertValid("worker-result", worker);
  if (worker.work_id !== manifest.work_id || worker.run_id !== manifest.run_id || worker.repository !== repository.name) {
    throw new Error("Worker result identity does not match the closeout run");
  }
  if (worker.status !== "completed" || worker.branch !== repository.branch || resolve(worker.worktree) !== resolve(repository.worktree)) {
    throw new Error("Closeout requires the completed worker recorded for this branch and worktree");
  }
  if (worker.commits.at(-1) !== headCommit || worker.commits.join("\n") !== commits.join("\n")) {
    throw new Error("Worktree commits changed after the recorded worker result");
  }
  if (worker.changed_files.slice().sort().join("\n") !== changedFiles.slice().sort().join("\n")) {
    throw new Error("Worktree changed-file set differs from the recorded worker result");
  }
}

async function findExistingContribution(root: string, runId: string, repository: string): Promise<string | null> {
  try {
    for (const entry of await readdir(root, { withFileTypes: true })) {
      if (!entry.isFile() || !entry.name.endsWith(".md")) continue;
      const path = join(root, entry.name);
      const content = await readFile(path, "utf8");
      if (content.includes(`- Run: \`${runId}\``) && content.includes(`- \`${repository}\` on branch`)) return path;
    }
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
  }
  return null;
}

async function ensureContributionRoot(workspaceRoot: string, path: string): Promise<void> {
  await mkdir(path, { recursive: true, mode: 0o755 });
  const info = await lstat(path);
  if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`Contribution path must be a real directory: ${path}`);
  assertInside(await realpath(workspaceRoot), await realpath(path));
}

function addExecutionEvent(manifest: RuntimeManifest, repository: RuntimeRepository, stage: "closeout-prepared" | "closeout-cleaned", from: ExecutionEvent["from_status"], to: ExecutionEvent["to_status"], occurredAt: string, resultPath: string): void {
  const key = `${manifest.run_id}:execution:${repository.name}:${stage}`;
  if (manifest.execution_events?.some((event) => event.idempotency_key === key)) return;
  manifest.execution_events ??= [];
  manifest.execution_events.push({
    stage,
    repository: repository.name,
    from_status: from,
    to_status: to,
    inferred: false,
    attempt: repository.repair_attempts ?? 0,
    result_path: resultPath,
    idempotency_key: key,
    occurred_at: occurredAt,
  });
  repository.status = to;
  const statuses = manifest.repositories.map((candidate) => candidate.status ?? manifest.status);
  if (statuses.every((status) => status === "closed")) manifest.status = "closed";
  else if (statuses.some((status) => status === "closing" || status === "closed")) manifest.status = "closing";
  else if (statuses.every((status) => status === "passed")) manifest.status = "passed";
  else if (statuses.includes("failed")) manifest.status = "failed";
  else if (statuses.includes("blocked")) manifest.status = "blocked";
  manifest.updated_at = occurredAt;
}

function addLifecycleEvent(manifest: RuntimeManifest, outcome: CloseoutRecord["outcome"], occurredAt: string): void {
  const event = outcome === "merged" ? "task.completed" : "task.cancelled";
  const key = `${manifest.run_id}:lifecycle:${event}:activity-none`;
  if (manifest.lifecycle_events.some((item) => item.idempotency_key === key)) return;
  manifest.lifecycle_events.push({ event, status: "skipped", idempotency_key: key, occurred_at: occurredAt });
}

function assertCloseoutLifecycleReady(manifest: RuntimeManifest, config: WorkspaceConfig, outcome: CloseoutRecord["outcome"]): void {
  if (config.activity.provider === "none") return;
  const event = outcome === "merged" ? "task.completed" : "task.cancelled";
  const lifecycle = manifest.lifecycle_events.find((item) => item.event === event);
  if (!lifecycle) {
    throw new Error(`Prepare configured activity hooks before closeout: node .agents/bin/cc.mjs prepare-lifecycle --run-id ${manifest.run_id} --event ${event}`);
  }
  if (lifecycle.status !== "completed" && lifecycle.status !== "skipped") {
    throw new Error(`Configured activity hook ${event} is ${lifecycle.status}; complete required or manual actions before closeout`);
  }
}

async function isAncestor(repository: string, ancestor: string, descendant: string): Promise<boolean> {
  try {
    await git(repository, ["merge-base", "--is-ancestor", ancestor, descendant]);
    return true;
  } catch {
    return false;
  }
}

async function assertVerifiedMergeConfirmation(
  workspaceRoot: string, runtimeRoot: string, config: WorkspaceConfig, manifest: RuntimeManifest, repository: RuntimeRepository, requestedMergeCommit?: string,
): Promise<MergeConfirmationRecord> {
  if (!repository.merge_confirmation) throw new Error(`Merged closeout is not ready: record and verify the human merge with node .agents/bin/cc.mjs confirm-merge --run-id ${manifest.run_id} --repository ${repository.name} --merge-commit <full-sha> --author <slug> --evidence <single-line-evidence>`);
  const record = await readJsonRegularInside<MergeConfirmationRecord>(runtimeRoot, repository.merge_confirmation, "Merge confirmation record");
  await assertValid("merge-confirmation-record", record);
  if (record.work_id !== manifest.work_id || record.run_id !== manifest.run_id || record.repository !== repository.name || record.head_commit !== await git(repository.worktree, ["rev-parse", "HEAD"])) {
    throw new Error("Merge confirmation identity or verified head does not match the active run");
  }
  if (requestedMergeCommit && requestedMergeCommit !== record.merge_commit) throw new Error("Requested merge commit differs from the verified merge confirmation");
  const defaultBranch = config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch;
  if (record.base_branch !== defaultBranch || ![`refs/heads/${defaultBranch}`, `refs/remotes/origin/${defaultBranch}`].includes(record.target_ref)) {
    throw new Error("Merge confirmation does not target the configured default branch");
  }
  const baseRepository = assertInside(workspaceRoot, join(workspaceRoot, repository.base_path));
  const currentTarget = await git(baseRepository, ["rev-parse", "--verify", `${record.target_ref}^{commit}`]);
  if (currentTarget !== record.target_commit || !await isAncestor(baseRepository, record.merge_commit, currentTarget) || !await isAncestor(baseRepository, record.head_commit, record.merge_commit)) {
    throw new Error("Merge confirmation no longer proves the exact reviewed head is reachable from the recorded default target");
  }
  return record;
}

async function verifiedDefaultRefs(repository: string, branch: string): Promise<string[]> {
  const refs: string[] = [];
  for (const ref of [`refs/heads/${branch}`, `refs/remotes/origin/${branch}`]) {
    try {
      await git(repository, ["rev-parse", "--verify", ref]);
      refs.push(ref);
    } catch {
      // Missing local or origin default refs are represented by cleanup blockers.
    }
  }
  return refs;
}

async function cleanupBlockers(workspaceRoot: string, config: WorkspaceConfig, repository: RuntimeRepository, record: CloseoutRecord): Promise<string[]> {
  const blockers: string[] = [];
  const contributionPath = assertInside(workspaceRoot, join(workspaceRoot, record.contribution));
  try {
    await git(workspaceRoot, ["ls-files", "--error-unmatch", "--", record.contribution]);
    if (await git(workspaceRoot, ["status", "--porcelain=v1", "--", record.contribution])) {
      blockers.push("Contribution has uncommitted wrapper changes; commit it through the configured wrapper workflow before cleanup.");
    }
  } catch {
    blockers.push("Contribution is not durably tracked by wrapper Git; commit it before cleanup.");
  }
  try {
    await access(contributionPath);
  } catch {
    blockers.push("Contribution file is missing; runtime cleanup would discard the only closeout record.");
  }

  const baseRepository = assertInside(workspaceRoot, join(workspaceRoot, repository.base_path));
  try {
    if (await git(baseRepository, ["status", "--porcelain=v1", "--untracked-files=normal"])) blockers.push("Base repository is dirty.");
    if (await git(repository.worktree, ["status", "--porcelain=v1", "--untracked-files=normal"])) blockers.push("Run worktree has uncommitted changes.");
    if (await git(repository.worktree, ["branch", "--show-current"]) !== repository.branch) blockers.push("Run worktree is on an unexpected branch.");
    if (await git(repository.worktree, ["rev-parse", "HEAD"]) !== record.head_commit) blockers.push("Run worktree HEAD changed after closeout preparation.");
  } catch (error) {
    blockers.push(`Run worktree is unavailable: ${(error as Error).message}`);
    return blockers;
  }

  const defaultBranch = config.repositories[repository.name]?.default_branch ?? config.workspace.default_branch;
  const defaultRefs = await verifiedDefaultRefs(baseRepository, defaultBranch);
  const headOnDefault = (await Promise.all(defaultRefs.map((ref) => isAncestor(baseRepository, record.head_commit, ref)))).some(Boolean);
  if (record.outcome === "merged") {
    let mergeEvidence = headOnDefault;
    if (!mergeEvidence && record.merge_commit) {
      const mergeOnDefault = (await Promise.all(defaultRefs.map((ref) => isAncestor(baseRepository, record.merge_commit!, ref)))).some(Boolean);
      mergeEvidence = mergeOnDefault && await isAncestor(baseRepository, record.base_commit, record.merge_commit);
    }
    if (!mergeEvidence) blockers.push("Merged outcome is not reachable from the configured default branch; fetch the merge or provide a verified merge commit.");
  } else if (record.head_commit !== record.base_commit && !headOnDefault) {
    const remoteRefs = await git(baseRepository, ["for-each-ref", "--format=%(refname)", "--contains", record.head_commit, "refs/remotes/"]);
    if (!remoteRefs) blockers.push("Abandoned branch contains commits that are neither merged nor preserved by a remote ref.");
  }
  return blockers;
}

async function closePreparedRun(workspaceRoot: string, manifestPath: string, manifest: RuntimeManifest, repository: RuntimeRepository, recordPath: string, record: CloseoutRecord, config: WorkspaceConfig, occurredAt: string): Promise<CloseoutRecord> {
  const blockers = await cleanupBlockers(workspaceRoot, config, repository, record);
  if (blockers.length > 0) {
    const blocked: CloseoutRecord = { ...record, status: "blocked", cleanup: { ...record.cleanup, requested: true }, blockers, updated_at: occurredAt };
    await assertValid("closeout-record", blocked);
    await writeJsonAtomic(recordPath, blocked);
    return blocked;
  }
  const baseRepository = assertInside(workspaceRoot, join(workspaceRoot, repository.base_path));
  await git(baseRepository, ["worktree", "remove", repository.worktree]);
  const closed: CloseoutRecord = {
    ...record,
    status: "closed",
    cleanup: { requested: true, worktree_removed: true, branch_preserved: true, runtime_evidence_preserved: true },
    blockers: [],
    updated_at: occurredAt,
  };
  await assertValid("closeout-record", closed);
  await writeJsonAtomic(recordPath, closed);
  addExecutionEvent(manifest, repository, "closeout-cleaned", "closing", "closed", occurredAt, recordPath);
  await assertValid("runtime-manifest", manifest);
  await writeJsonAtomic(manifestPath, manifest);
  return closed;
}

export async function finishWork(options: FinishWorkOptions): Promise<CloseoutRecord> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, "manifest.json"));
  const lockPath = `${manifestPath}.lock`;
  const author = safeToken(options.author, "Author");
  const invocationTime = options.now ?? new Date();
  if (options.mergeCommit && !/^[a-f0-9]{40,64}$/.test(options.mergeCommit)) throw new Error("--merge-commit must be a full lowercase Git object ID");
  if (options.pullRequests?.some((reference) => !reference.trim() || /[\r\n]/.test(reference))) throw new Error("Pull-request references must be non-empty single lines");
  const config = await loadWorkspace(workspaceRoot);
  const wrapperTopLevel = await git(workspaceRoot, ["rev-parse", "--show-toplevel"]);
  if (await realpath(wrapperTopLevel) !== await realpath(workspaceRoot)) throw new Error("Workspace root must be the wrapper Git root before closeout");

  return withExclusiveFile(lockPath, async () => {
    const manifest = await readJsonRegularInside<RuntimeManifest>(runtimeRoot, manifestPath, "Runtime manifest");
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match the requested run");
    assertCloseoutLifecycleReady(manifest, config, options.outcome);
    const repository = findRepository(manifest, options.repository);
    const recordPath = assertInside(runtimeRoot, join(runtimeRoot, "runs", options.runId, `${repository.name}-closeout.json`));
    if (repository.closeout_record) {
      const existing = await readJsonRegularInside<CloseoutRecord>(runtimeRoot, repository.closeout_record, "Closeout record");
      await assertValid("closeout-record", existing);
      if (existing.outcome !== options.outcome || existing.author !== author) throw new Error("Closeout was already prepared with different human intent");
      if (existing.outcome === "merged") await assertVerifiedMergeConfirmation(workspaceRoot, runtimeRoot, config, manifest, repository, options.mergeCommit);
      if (existing.status === "closed" || !options.cleanup) return existing;
      return closePreparedRun(workspaceRoot, manifestPath, manifest, repository, recordPath, existing, config, invocationTime.toISOString());
    }
    const repositoryStatus = repository.status ?? manifest.status;
    if (!["passed", "failed", "blocked", "cancelled"].includes(repositoryStatus)) throw new Error(`Closeout preparation requires a terminal repository outcome, received ${repositoryStatus}`);
    if (options.outcome === "merged" && repositoryStatus !== "passed") throw new Error(`Merged closeout requires a passed repository, received ${repositoryStatus}`);
    if (options.outcome === "abandoned" && !options.reason?.trim()) throw new Error("Deliberate abandonment requires --reason");
    const mergeConfirmation = options.outcome === "merged"
      ? await assertVerifiedMergeConfirmation(workspaceRoot, runtimeRoot, config, manifest, repository, options.mergeCommit)
      : null;

    const brief = await readJsonRegularInside<TaskBrief>(runtimeRoot, manifest.task_brief, "Task brief");
    await assertValid("task-brief", brief);
    const headCommit = await git(repository.worktree, ["rev-parse", "HEAD"]);
    if (await git(repository.worktree, ["branch", "--show-current"]) !== repository.branch) throw new Error("Run worktree is on an unexpected branch");
    const commits = (await git(repository.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${headCommit}`])).split("\n").filter(Boolean);
    const changedFiles = (await git(repository.worktree, ["diff", "--name-only", `${repository.base_commit}...${headCommit}`])).split("\n").filter(Boolean);
    if (options.outcome === "merged") await assertCurrentWorker(runtimeRoot, manifest, repository, headCommit, commits, changedFiles);
    const verifier = await optionalVerifier(runtimeRoot, manifest, repository);
    if (options.outcome === "merged" && verifier?.status !== "pass") throw new Error("Merged closeout requires the recorded passing verifier result");
    const verification = verifier ? [verifier.summary, ...verifier.checks, ...verifier.acceptance.map((item) => `${item.criterion}: ${item.status} — ${item.evidence}`)] : [];
    const preparedAt = invocationTime.toISOString();
    const contributionsRoot = assertInside(workspaceRoot, join(workspaceRoot, "contributions", "general"));
    await ensureContributionRoot(workspaceRoot, contributionsRoot);
    const existingContribution = await findExistingContribution(contributionsRoot, options.runId, repository.name);
    const contributionPath = existingContribution ?? join(contributionsRoot, `${compactTimestamp(invocationTime)}-${author}-${taskSlug(brief)}-${repository.name}.md`);
    const contributionRelative = relative(workspaceRoot, contributionPath).replaceAll("\\", "/");
    const record: CloseoutRecord = {
      contract_version: 1,
      work_id: manifest.work_id,
      run_id: manifest.run_id,
      repository: repository.name,
      outcome: options.outcome,
      status: "prepared",
      author,
      reason: options.reason?.trim() || null,
      contribution: contributionRelative,
      branch: repository.branch,
      base_commit: repository.base_commit,
      head_commit: headCommit,
      merge_commit: mergeConfirmation?.merge_commit ?? options.mergeCommit ?? null,
      pull_requests: [...new Set(options.pullRequests ?? [])],
      commits,
      changed_files: changedFiles,
      verification,
      cleanup: { requested: Boolean(options.cleanup), worktree_removed: false, branch_preserved: true, runtime_evidence_preserved: true },
      blockers: [],
      prepared_at: preparedAt,
      updated_at: preparedAt,
    };
    const document = contributionDocument(manifest, repository, brief, record);
    const documentErrors = contributionDocumentErrors(contributionPath, document, options.runId);
    if (documentErrors.length > 0) throw new Error(`Invalid contribution: ${documentErrors.join("; ")}`);
    if (existingContribution && await readFile(existingContribution, "utf8") !== document) {
      throw new Error("An append-only contribution already exists for this run with different closeout content");
    }
    if (!existingContribution) await writeTextExclusive(contributionPath, document);
    await assertValid("closeout-record", record);
    await writeJsonAtomic(recordPath, record);
    repository.closeout_record = recordPath;
    repository.contribution = contributionRelative;
    if (!manifest.evidence.includes(recordPath)) manifest.evidence.push(recordPath);
    const fromStatus = repositoryStatus as ExecutionEvent["from_status"];
    addExecutionEvent(manifest, repository, "closeout-prepared", fromStatus, "closing", preparedAt, recordPath);
    if (config.activity.provider === "none") addLifecycleEvent(manifest, options.outcome, preparedAt);
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(manifestPath, manifest);
    if (!options.cleanup) return record;
    return closePreparedRun(workspaceRoot, manifestPath, manifest, repository, recordPath, record, config, preparedAt);
  });
}
