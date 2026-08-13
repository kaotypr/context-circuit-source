import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { access, cp, lstat, mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join, resolve, sep } from "node:path";
import test from "node:test";
import { configureWorkspace } from "../scripts/lib/configure-workspace.js";
import { finishWork } from "../scripts/lib/finish-work.js";
import { git } from "../scripts/lib/git.js";
import { createPlanDraft, setPlanState } from "../scripts/lib/plans.js";
import { recordResult } from "../scripts/lib/record-result.js";
import { confirmMerge, prepareReview, recordReviewPublication } from "../scripts/lib/review-lifecycle.js";
import { preparePlanTask, type PreparedTask } from "../scripts/lib/run-task.js";
import type { PlanDraftRequest, ReviewCommand, TaskBrief, TestExpectation, WorkspaceBootstrapRequest } from "../scripts/lib/types.js";
import { recommendWhatsNext } from "../scripts/lib/whats-next.js";
import { requiredWorkspaceDocuments, validateContract } from "../scripts/lib/validation.js";
import { projectRoot } from "./helpers.js";

const workerIdentity = "worker-e2e";
const verifierIdentity = "verifier-e2e";
const acceptanceCriterion = "The selected behavior is delivered.";

const plan: PlanDraftRequest = {
  contract_version: 1,
  plan_id: "journey",
  title: "Complete journey",
  source: { kind: "prd", reference: "docs/journey.md" },
  work_prefix: "JOURNEY",
  summary: "Exercise the complete plan-linked workflow.",
  affected_repositories: ["frontend"],
  assumptions: [], open_questions: [],
  requirements: [acceptanceCriterion],
  solution: ["Make one scoped product change."],
  delivery: ["Review and merge through the human gate."],
  verification: ["Verify the acceptance criterion independently."],
  risks: [],
  work_items: [{
    key: "behavior", title: "Deliver selected behavior", area: "application behavior", repository: "frontend",
    scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: ["npm test"],
    acceptance_criteria: [acceptanceCriterion],
  }],
};

interface ConfiguredFixture {
  root: string;
  repository: string;
  cleanup: () => Promise<void>;
}

interface CommandEvidence {
  command: string;
  cwd: string;
  status: number;
  stdout: string;
  stderr: string;
}

interface ScopedVerifierInput {
  contract_version: 1;
  role: "verifier";
  read_only: true;
  task_brief: string;
  repository: string;
  worktree: string;
  branch: string;
  base_commit: string;
  worker_result: string;
  acceptance_criteria: string[];
  test_expectation: TestExpectation;
  verification_commands: string[];
  instruction_paths: string[];
  result_contract: string;
  result_path: string;
}

async function writeJson(path: string, value: unknown): Promise<void> {
  await mkdir(dirname(path), { recursive: true });
  await writeFile(path, `${JSON.stringify(value, null, 2)}\n`, "utf8");
}

function execute(cwd: string, command: string, args: string[]): CommandEvidence {
  const env = { ...process.env };
  delete env.NODE_TEST_CONTEXT;
  const result = spawnSync(command, args, { cwd, encoding: "utf8", env });
  assert.equal(result.error, undefined);
  return { command: [command, ...args].join(" "), cwd, status: result.status ?? -1, stdout: result.stdout, stderr: result.stderr };
}

function executeReviewCommand(value: ReviewCommand): CommandEvidence {
  return execute(value.cwd, value.argv[0]!, value.argv.slice(1));
}

function concise(value: string): string {
  return value.trim().replace(/\s+/g, " ").slice(0, 300);
}

function assertInside(root: string, candidate: string, label: string): string {
  const resolvedRoot = resolve(root);
  const resolved = resolve(candidate);
  assert.ok(resolved === resolvedRoot || resolved.startsWith(`${resolvedRoot}${sep}`), `${label} must stay inside the fixture workspace`);
  return resolved;
}

async function readRegularJson<T>(root: string, candidate: string, label: string): Promise<T> {
  const path = assertInside(root, candidate, label);
  const info = await lstat(path);
  assert.equal(info.isFile() && !info.isSymbolicLink(), true, `${label} must be a regular file`);
  return JSON.parse(await readFile(path, "utf8")) as T;
}

function bootstrapRequest(): WorkspaceBootstrapRequest {
  return {
    contract_version: 1,
    configuration: {
      version: 1,
      template_version: "0.2.1",
      workspace: { name: "e2e-workspace", mode: "team", default_branch: "main", purpose: "Exercises the complete Context Circuit journey." },
      repositories: { frontend: { path: "repositories/frontend", mode: "ignored-clone", role: "product application", agent: "frontend", default_branch: "main" } },
      activity: { provider: "none", access: "auto", required_capabilities: [], optional_capabilities: [] },
      workflow: { human_gates: ["plan-approval", "task-selection", "merge"], maximum_repair_attempts: 2, wrapper_change_policy: "pull-request", review_mode: "local" },
      context: { authoritative_sources: [{ kind: "repository-documentation", reference: "repositories/frontend/README.md", purpose: "Product fixture instructions", repository: "frontend" }] },
    },
    context: {
      project_summary: "Exercises the complete Context Circuit journey.",
      architecture: ["One configured frontend repository owns product behavior."],
      conventions: ["Run npm test for product verification."],
      decisions: ["Keep merge and closeout human-gated."],
      sources: [{ kind: "repository-documentation", reference: "repositories/frontend/README.md", purpose: "Product fixture instructions", repository: "frontend" }],
    },
    wrapper: { initialize_git: true, authorize_initial_commit: true, commit_message: "chore: configure e2e workspace", author_name: "Configurator", author_email: "configurator@example.invalid" },
    repositories: [{ name: "frontend", source: "existing", authorize_initial_commit: false }],
  };
}

async function configuredFixture(): Promise<ConfiguredFixture> {
  const root = await mkdtemp(join(tmpdir(), "context-circuit-e2e-"));
  for (const path of requiredWorkspaceDocuments) {
    const target = join(root, path);
    await mkdir(dirname(target), { recursive: true });
    await cp(join(projectRoot, path), target);
  }
  await cp(join(projectRoot, ".gitignore"), join(root, ".gitignore"));
  const repository = join(root, "repositories", "frontend");
  await mkdir(join(repository, "src"), { recursive: true });
  await mkdir(join(repository, "test"), { recursive: true });
  await writeFile(join(repository, "README.md"), "# Product fixture\n", "utf8");
  await writeFile(join(repository, "package.json"), `${JSON.stringify({ private: true, scripts: { test: "node --test test/behavior.test.mjs" } }, null, 2)}\n`, "utf8");
  await writeFile(join(repository, "src", "App.tsx"), "export const selectedBehavior = 'pending';\n", "utf8");
  await writeFile(join(repository, "test", "behavior.test.mjs"), `import assert from "node:assert/strict";\nimport { readFile } from "node:fs/promises";\nimport test from "node:test";\ntest("selected behavior is delivered", async () => { assert.match(await readFile(new URL("../src/App.tsx", import.meta.url), "utf8"), /selectedBehavior = 'delivered'/); console.log("behavior-check: delivered"); });\n`, "utf8");
  await git(repository, ["init", "--initial-branch=main"]);
  await git(repository, ["add", "."]);
  await git(repository, ["-c", "user.name=Product", "-c", "user.email=product@example.invalid", "commit", "-m", "test: product fixture"]);
  const productHead = await git(repository, ["rev-parse", "HEAD"]);

  const configured = await configureWorkspace({ workspaceRoot: root, request: bootstrapRequest() });
  assert.equal(configured.route, "bootstrap");
  assert.equal(await git(root, ["rev-list", "--count", "HEAD"]), "1");
  assert.equal(await git(repository, ["rev-parse", "HEAD"]), productHead);
  assert.equal(await git(root, ["status", "--porcelain=v1"]), "");
  assert.equal(await git(repository, ["status", "--porcelain=v1"]), "");
  const readme = await readFile(join(root, "README.md"), "utf8");
  assert.match(readme, /e2e-workspace/);
  assert.match(readme, /repositories\/frontend/);
  const sources = await readFile(join(root, "context", "SOURCES.md"), "utf8");
  assert.match(sources, /repositories\/frontend\/README\.md/);
  assert.match(sources, /Product fixture instructions/);
  return { root, repository, cleanup: () => rm(root, { recursive: true, force: true }) };
}

async function runWorker(root: string, prepared: PreparedTask): Promise<CommandEvidence> {
  const app = join(prepared.worktree, "src", "App.tsx");
  await writeFile(app, (await readFile(app, "utf8")).replace("'pending'", "'delivered'"), "utf8");
  await git(prepared.worktree, ["add", "src/App.tsx"]);
  await git(prepared.worktree, ["-c", `user.name=${workerIdentity}`, "-c", "user.email=worker@example.invalid", "commit", "-m", "feat: deliver journey behavior"]);
  const head = await git(prepared.worktree, ["rev-parse", "HEAD"]);
  const testRun = execute(prepared.worktree, "npm", ["test"]);
  assert.equal(testRun.status, 0, testRun.stderr);
  assert.match(`${testRun.stdout}\n${testRun.stderr}`, /behavior-check: delivered/);
  const worker = JSON.parse(await readFile(prepared.workerInput, "utf8"));
  await writeJson(worker.result_path, {
    contract_version: 1, work_id: prepared.workId, run_id: prepared.runId, repository: "frontend", status: "completed",
    summary: `${workerIdentity} delivered the selected plan item.`, branch: prepared.branch, worktree: prepared.worktree,
    commits: [head], changed_files: ["src/App.tsx"],
    checks: [{ command: "npm test", status: "passed", evidence: `exit ${testRun.status}; ${concise(testRun.stdout)}` }], risks: [],
  });
  const manifest = await recordResult({ workspaceRoot: root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  assert.equal(manifest.status, "verifying");
  return testRun;
}

async function runIndependentVerifier(root: string, verifierInputPath: string): Promise<{ testRuns: CommandEvidence[]; resultPath: string }> {
  assert.notEqual(workerIdentity, verifierIdentity);

  // Validate the emitted verifier boundary and its task contract before running
  // any repository inspection or configured command.
  const verifier = await readRegularJson<ScopedVerifierInput>(root, verifierInputPath, "Verifier input");
  assert.equal(verifier.contract_version, 1);
  assert.equal(verifier.role, "verifier");
  assert.equal(verifier.read_only, true);
  assert.match(verifier.repository, /^[a-z][a-z0-9-]*$/);
  assert.match(verifier.base_commit, /^[a-f0-9]{40,64}$/);
  assert.ok(verifier.acceptance_criteria.length > 0);
  assert.ok(verifier.verification_commands.length > 0);
  assert.ok(verifier.verification_commands.every((command) => command.trim() && !/[\r\n]/.test(command)));
  assertInside(root, verifier.worktree, "Verifier worktree");
  assertInside(join(root, ".runtime"), verifier.result_path, "Verifier result path");
  const brief = await readRegularJson<TaskBrief>(root, verifier.task_brief, "Verifier task brief");
  assert.deepEqual(await validateContract("task-brief", brief), []);
  assert.equal(brief.run_id.length > 0, true);
  const target = brief.repositories.find((candidate) => candidate.name === verifier.repository);
  assert.ok(target, `Verifier repository ${verifier.repository} must exist in the task brief`);
  assert.deepEqual(verifier.acceptance_criteria, target.acceptance_criteria);
  assert.deepEqual(verifier.verification_commands, target.verification_commands);
  assert.deepEqual(verifier.test_expectation, target.test_expectation);

  const head = await git(verifier.worktree, ["rev-parse", "HEAD"]);
  assert.equal(await git(verifier.worktree, ["branch", "--show-current"]), verifier.branch);
  const inspected = execute(verifier.worktree, "git", ["diff", "--check", `${verifier.base_commit}...${head}`]);
  assert.equal(inspected.status, 0, inspected.stderr);
  const changedFiles = (await git(verifier.worktree, ["diff", "--name-only", `${verifier.base_commit}...${head}`])).split("\n").filter(Boolean);
  assert.ok(changedFiles.length > 0);
  assert.ok(changedFiles.every((path) => brief.scope.includes(path)), `Verifier diff escaped task scope: ${changedFiles.join(", ")}`);
  const diff = execute(verifier.worktree, "git", ["diff", "--unified=0", `${verifier.base_commit}...${head}`, "--", ...changedFiles]);
  assert.equal(diff.status, 0, diff.stderr);
  assert.match(diff.stdout, /selectedBehavior = 'delivered'/);
  const testRuns = verifier.verification_commands.map((command) => execute(verifier.worktree, "sh", ["-lc", command]));
  for (const testRun of testRuns) {
    assert.equal(testRun.status, 0, `${testRun.command}: ${testRun.stderr}`);
    assert.match(`${testRun.stdout}\n${testRun.stderr}`, /behavior-check: delivered/);
  }
  await writeJson(verifier.result_path, {
    contract_version: 1, work_id: brief.work_id, run_id: brief.run_id, repository: verifier.repository, status: "pass",
    summary: `${verifierIdentity} independently inspected the committed diff and reran the configured command.`,
    acceptance: verifier.acceptance_criteria.map((criterion) => ({ criterion, status: "passed", evidence: `Exact ${verifier.base_commit}...${head} diff stayed in scope and configured checks passed.` })),
    checks: [`git diff --check exit ${inspected.status}`, ...testRuns.map((run) => `${run.command} exit ${run.status}: ${concise(run.stdout)}`)], findings: [], verified_at: "2026-08-12T03:04:00Z",
  });
  const manifest = await recordResult({ workspaceRoot: root, runId: brief.run_id, repository: verifier.repository, stage: "verifier-result" });
  assert.equal(manifest.status, "passed");
  assert.deepEqual(manifest.execution_events?.slice(-2).map((event) => event.stage), ["worker-result", "verifier-result"]);
  assert.deepEqual(await validateContract("runtime-manifest", manifest), []);
  return { testRuns, resultPath: verifier.result_path };
}

async function preparePassedPlanRun(root: string, discriminator: string): Promise<PreparedTask> {
  const created = await createPlanDraft(root, plan, new Date("2026-08-12T03:00:00Z"));
  const approved = await setPlanState(created.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-12T03:01:00Z"));
  await git(root, ["add", "context/plans/journey"]);
  await git(root, ["-c", "user.name=Owner", "-c", "user.email=owner@example.invalid", "commit", "-m", "plan: approve journey"]);
  const prepared = await preparePlanTask({
    workspaceRoot: root,
    request: { contract_version: 1, source: { kind: "plan", reference: "context/plans/journey", plan_version: approved.plan_version, approved_digest: approved.approved_digest! }, work_ids: ["JOURNEY-001"] },
    discriminator, now: new Date("2026-08-12T03:02:00Z"),
  });
  const workerTest = await runWorker(root, prepared);
  const emittedVerifier = await readRegularJson<ScopedVerifierInput>(root, prepared.verifierInput, "Verifier input");
  const tamperedPath = join(root, ".runtime", "runs", prepared.runId, "tampered-verifier-input.json");
  await writeJson(tamperedPath, { ...emittedVerifier, repository: "backend" });
  await assert.rejects(runIndependentVerifier(root, tamperedPath), /must exist in the task brief/);
  await assert.rejects(access(emittedVerifier.result_path));
  const verifier = await runIndependentVerifier(root, prepared.verifierInput);
  assert.equal(workerTest.command, verifier.testRuns[0]?.command.replace(/^sh -lc /, ""));
  assert.equal(workerTest.cwd, verifier.testRuns[0]?.cwd);
  const workerResult = JSON.parse(await readFile(JSON.parse(await readFile(prepared.workerInput, "utf8")).result_path, "utf8"));
  const verifierResult = JSON.parse(await readFile(verifier.resultPath, "utf8"));
  assert.match(workerResult.summary, new RegExp(`^${workerIdentity}`));
  assert.match(verifierResult.summary, new RegExp(`^${verifierIdentity}`));
  return prepared;
}

test("configured approved-plan journey reaches verified merge, safe closeout, and completed exclusion", async (t) => {
  const workspace = await configuredFixture();
  t.after(workspace.cleanup);
  const prepared = await preparePassedPlanRun(workspace.root, "e2e00001");
  const review = await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  assert.equal(review.status, "ready-for-local-review");
  assert.equal(review.commands?.diff.cwd, prepared.worktree);
  assert.equal(review.commands?.show.cwd, prepared.worktree);
  assert.equal(review.commands?.tests[0]?.cwd, prepared.worktree);
  const reviewDiff = executeReviewCommand(review.commands!.diff);
  const reviewShow = executeReviewCommand(review.commands!.show);
  const reviewTest = executeReviewCommand(review.commands!.tests[0]!);
  assert.equal(reviewDiff.status, 0, reviewDiff.stderr);
  assert.match(reviewDiff.stdout, /src\/App\.tsx/);
  assert.equal(reviewShow.status, 0, reviewShow.stderr);
  assert.match(reviewShow.stdout, new RegExp(review.head_commit.slice(0, 7)));
  assert.equal(reviewTest.status, 0, reviewTest.stderr);
  assert.match(`${reviewTest.stdout}\n${reviewTest.stderr}`, /behavior-check: delivered/);
  await git(workspace.repository, ["-c", "user.name=Human", "-c", "user.email=human@example.invalid", "merge", "--no-ff", "-m", "merge: journey", prepared.branch]);
  const mergeCommit = await git(workspace.repository, ["rev-parse", "HEAD"]);
  const confirmation = await confirmMerge({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", mergeCommit, author: "owner", evidence: "Human verified and merged the reviewed head." });
  assert.equal(confirmation.status, "closeout-ready");
  const closeout = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "owner", mergeCommit, now: new Date("2026-08-12T03:10:00Z") });
  assert.equal(closeout.status, "prepared");
  const contribution = await readFile(join(workspace.root, closeout.contribution), "utf8");
  assert.match(contribution, /Task source: plan/);
  assert.match(contribution, /Plan: `context\/plans\/journey`/);
  await git(workspace.root, ["add", closeout.contribution]);
  await git(workspace.root, ["-c", "user.name=Owner", "-c", "user.email=owner@example.invalid", "commit", "-m", "docs: record journey outcome"]);
  const closed = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "owner", mergeCommit, cleanup: true });
  assert.equal(closed.status, "closed");
  await access(prepared.manifest);
  assert.equal(await git(workspace.repository, ["show-ref", "--verify", `refs/heads/${prepared.branch}`]).then(() => true), true);
  const next = await recommendWhatsNext(workspace.root, null, new Date("2026-08-12T03:20:00Z"));
  assert.equal(next.considered.excluded, 1);
  assert.equal(next.recommendation.candidate_id, null);
});

test("remote publication fixture records confirmed evidence without pushing", async (t) => {
  const workspace = await configuredFixture();
  t.after(workspace.cleanup);
  await git(workspace.repository, ["remote", "add", "origin", "https://example.invalid/product.git"]);
  const prepared = await preparePassedPlanRun(workspace.root, "e2e00002");
  const refsBefore = await git(workspace.repository, ["for-each-ref", "--format=%(refname)", "refs/remotes/"]);
  const review = await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  assert.equal(review.status, "ready-for-publication");
  const publication = await recordReviewPublication({
    workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", status: "published", tool: "manual",
    pullRequest: "https://example.invalid/product/pull/7", evidence: "Fixture confirmation; no network operation was performed.", authorized: true,
  });
  assert.deepEqual(await validateContract("review-publication-record", publication), []);
  assert.equal(await git(workspace.repository, ["for-each-ref", "--format=%(refname)", "refs/remotes/"]), refsBefore);
  assert.equal(await git(prepared.worktree, ["rev-parse", "HEAD"]), publication.head_commit);
});
