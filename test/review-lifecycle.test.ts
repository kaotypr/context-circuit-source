import assert from "node:assert/strict";
import { mkdir, readFile, symlink, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import test from "node:test";
import { git } from "../scripts/lib/git.js";
import { recordResult } from "../scripts/lib/record-result.js";
import { confirmMerge, prepareRepair, prepareReview, recordReviewPublication } from "../scripts/lib/review-lifecycle.js";
import { preparePlanlessTask, type PreparedTask } from "../scripts/lib/run-task.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace, taskOptions } from "./helpers.js";

async function writeJson(path: string, value: unknown): Promise<void> {
  await mkdir(dirname(path), { recursive: true });
  await writeFile(path, `${JSON.stringify(value, null, 2)}\n`, "utf8");
}

async function currentInputs(prepared: PreparedTask): Promise<{ worker: Record<string, unknown>; verifier: Record<string, unknown> }> {
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  return {
    worker: JSON.parse(await readFile(manifest.repositories[0].worker_input, "utf8")),
    verifier: JSON.parse(await readFile(manifest.repositories[0].verifier_input, "utf8")),
  };
}

async function commitWorkerResult(prepared: PreparedTask, marker: string): Promise<void> {
  const appPath = join(prepared.worktree, "src", "App.tsx");
  await writeFile(appPath, `${await readFile(appPath, "utf8")}\n// ${marker}\n`, "utf8");
  await git(prepared.worktree, ["add", "src/App.tsx"]);
  await git(prepared.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", `test: ${marker}`]);
  const inputs = await currentInputs(prepared);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const repository = manifest.repositories[0];
  const head = await git(prepared.worktree, ["rev-parse", "HEAD"]);
  const commits = (await git(prepared.worktree, ["rev-list", "--reverse", `${repository.base_commit}..${head}`])).split("\n").filter(Boolean);
  const changedFiles = (await git(prepared.worktree, ["diff", "--name-only", `${repository.base_commit}...${head}`])).split("\n").filter(Boolean);
  await writeJson(inputs.worker.result_path as string, {
    contract_version: 1,
    work_id: prepared.workId,
    run_id: prepared.runId,
    repository: "frontend",
    status: "completed",
    summary: `Completed ${marker}.`,
    branch: prepared.branch,
    worktree: prepared.worktree,
    commits,
    changed_files: changedFiles,
    checks: [{ command: "git diff --check", status: "passed", evidence: "Passed" }],
    risks: [],
  });
}

async function writeVerifier(prepared: PreparedTask, status: "pass" | "fail"): Promise<void> {
  const inputs = await currentInputs(prepared);
  await writeJson(inputs.verifier.result_path as string, {
    contract_version: 1,
    work_id: prepared.workId,
    run_id: prepared.runId,
    repository: "frontend",
    status,
    summary: status === "pass" ? "Independent verification passed." : "The implementation needs repair.",
    acceptance: [{
      criterion: taskOptions.acceptanceCriteria[0],
      status: status === "pass" ? "passed" : "failed",
      evidence: status === "pass" ? "Behavior inspected." : "Reset behavior is incomplete.",
    }],
    checks: ["git diff --check"],
    findings: status === "pass" ? [] : [{ severity: "high", description: "Reset remains incomplete", evidence: "Verifier inspection" }],
    verified_at: "2026-08-11T18:03:00.000Z",
  });
}

async function failedRun(workspaceRoot: string): Promise<PreparedTask> {
  const prepared = await preparePlanlessTask({
    workspaceRoot,
    ...taskOptions,
    now: new Date("2026-08-11T18:00:00.000Z"),
    discriminator: "a1b2c3d4",
  });
  await commitWorkerResult(prepared, "initial implementation");
  await recordResult({ workspaceRoot, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "fail");
  await recordResult({ workspaceRoot, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  return prepared;
}

test("prepares a fresh bounded repair and blocks after the configured limit", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const configPath = join(workspace.root, "workspace.yaml");
  await writeFile(configPath, (await readFile(configPath, "utf8")).replace("maximum_repair_attempts: 2", "maximum_repair_attempts: 1"), "utf8");
  const prepared = await failedRun(workspace.root);

  const repair = await prepareRepair({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", now: new Date("2026-08-11T18:04:00.000Z") });
  assert.equal(repair.status, "prepared");
  assert.equal(repair.attempt, 1);
  const repairInput = JSON.parse(await readFile(repair.worker_input!, "utf8"));
  assert.equal(repairInput.role, "repair-worker");
  assert.match(repairInput.findings.join("\n"), /Reset remains incomplete/);
  assert.deepEqual(await prepareRepair({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" }), repair);

  await commitWorkerResult(prepared, "repair attempt one");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "fail");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  const exhausted = await prepareRepair({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", now: new Date("2026-08-11T18:08:00.000Z") });
  assert.equal(exhausted.status, "exhausted");
  assert.equal(exhausted.attempt, 1);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.status, "blocked");
  assert.deepEqual(manifest.execution_events.slice(-4).map((event: { stage: string }) => event.stage), ["repair-prepared", "worker-result", "verifier-result", "repair-exhausted"]);
});

test("prepares an idempotent draft-PR handoff without pushing or exposing remote credentials", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await git(workspace.repository, ["remote", "add", "origin", "https://user:secret@example.invalid/project.git"]);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T19:00:00.000Z"),
    discriminator: "d4c3b2a1",
  });
  await commitWorkerResult(prepared, "passing implementation");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "pass");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });

  const review = await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", now: new Date("2026-08-11T19:04:00.000Z") });
  assert.equal(review.status, "ready-for-publication");
  assert.equal(review.remote, "origin");
  assert.equal(review.blockers.length, 0);
  assert.doesNotMatch(JSON.stringify(review), /secret/);
  assert.deepEqual(await validateContract("review-preparation", review), []);
  assert.deepEqual(await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" }), review);
  await assert.rejects(recordReviewPublication({
    workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", status: "published", tool: "gh",
    pullRequest: "https://example.invalid/project/pull/42", evidence: "Unconfirmed publication claim",
  }), /explicit authorization/);
  const publication = await recordReviewPublication({
    workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", status: "published", tool: "gh",
    pullRequest: "https://example.invalid/project/pull/42", evidence: "gh confirmed draft pull request 42", authorized: true, now: new Date("2026-08-11T19:05:00.000Z"),
  });
  assert.equal(publication.status, "published");
  assert.deepEqual(await validateContract("review-publication-record", publication), []);
  assert.deepEqual(await recordReviewPublication({
    workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", status: "published", tool: "gh",
    pullRequest: "https://example.invalid/project/pull/42", evidence: "gh confirmed draft pull request 42", authorized: true,
  }), publication);
  await assert.rejects(recordReviewPublication({
    workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", status: "published", tool: "gh",
    pullRequest: "https://example.invalid/project/pull/43", evidence: "conflict", authorized: true,
  }), /different confirmed evidence/);
  assert.equal(await git(prepared.worktree, ["status", "--porcelain=v1"]), "");
  assert.equal(await git(prepared.worktree, ["branch", "--show-current"]), prepared.branch);
});

test("records a local-review-ready handoff when origin is unavailable", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T20:00:00.000Z"),
    discriminator: "1234abcd",
  });
  await commitWorkerResult(prepared, "passing implementation");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "pass");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  const review = await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  assert.equal(review.status, "ready-for-local-review");
  assert.deepEqual(review.blockers, []);
  assert.equal(review.commands?.merge.argv.at(-1), review.head_commit);
});

test("rejects runtime evidence changed after verifier recording", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T21:00:00.000Z"),
    discriminator: "bad0cafe",
  });
  await commitWorkerResult(prepared, "passing implementation");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "pass");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  const inputs = await currentInputs(prepared);
  const verifierPath = inputs.verifier.result_path as string;
  const verifier = JSON.parse(await readFile(verifierPath, "utf8"));
  verifier.run_id = "20260811T210000Z-deadbeef";
  await writeJson(verifierPath, verifier);

  await assert.rejects(
    prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" }),
    /Verifier result identity does not match/,
  );
});

test("emits shell-safe exact commands and rejects forged branch evidence without Git mutation", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const configPath = join(workspace.root, "workspace.yaml");
  const prepared = await preparePlanlessTask({ workspaceRoot: workspace.root, ...taskOptions, now: new Date("2026-08-11T21:10:00.000Z"), discriminator: "5afe5afe" });
  await commitWorkerResult(prepared, "shell-safe review");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "pass");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  await writeFile(configPath, (await readFile(configPath, "utf8")).replaceAll("default_branch: main", "default_branch: main;touch-pwned"), "utf8");
  const review = await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  assert.match(review.commands!.switch_target.shell, /'main;touch-pwned'/);
  assert.equal(review.commands!.merge.argv.at(-1), review.head_commit);
  const beforeRefs = await git(workspace.repository, ["show-ref"]);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  manifest.repositories[0].branch = "forged;touch-pwned";
  await writeJson(prepared.manifest, manifest);
  await assert.rejects(prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" }), /branch or worktree does not match|valid Git branch/);
  assert.equal(await git(workspace.repository, ["show-ref"]), beforeRefs);
});

test("verifies the exact reviewed head on the configured default target before closeout", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({ workspaceRoot: workspace.root, ...taskOptions, now: new Date("2026-08-11T21:20:00.000Z"), discriminator: "c105e000" });
  await commitWorkerResult(prepared, "merge confirmation");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "pass");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  const review = await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  const beforeRefs = await git(workspace.repository, ["show-ref"]);
  await assert.rejects(confirmMerge({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", mergeCommit: review.head_commit, author: "kao", evidence: "Forged merge report" }), /not reachable from the configured default target/);
  assert.equal(await git(workspace.repository, ["show-ref"]), beforeRefs);
  await git(workspace.repository, ["-c", "user.name=Merger", "-c", "user.email=merger@example.invalid", "merge", "--no-ff", "-m", "merge: reviewed head", prepared.branch]);
  const mergeCommit = await git(workspace.repository, ["rev-parse", "HEAD"]);
  const confirmation = await confirmMerge({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", mergeCommit, author: "kao", evidence: "Human confirmed merge after review." });
  assert.equal(confirmation.status, "closeout-ready");
  assert.equal(confirmation.target_ref, "refs/heads/main");
  assert.equal(confirmation.head_commit, review.head_commit);
  assert.match(confirmation.finish_work_shell, /'finish-work'/);
  assert.deepEqual(await validateContract("merge-confirmation-record", confirmation), []);
  const configPath = join(workspace.root, "workspace.yaml");
  await writeFile(configPath, (await readFile(configPath, "utf8")).replaceAll("default_branch: main", "default_branch: trunk"), "utf8");
  await assert.rejects(confirmMerge({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", mergeCommit, author: "kao", evidence: "Wrong target" }), /target differs from the configured default branch/);
});

test("review test commands preserve exact shell semantics in a quoted worktree cwd", async (t) => {
  const workspace = await createTestWorkspace("context circuit 'quoted path' ");
  t.after(workspace.cleanup);
  const recorded = `printf '%s\\n' "a b" && test -f src/App.tsx`;
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root, ...taskOptions, verificationCommands: [recorded],
    now: new Date("2026-08-11T21:30:00.000Z"), discriminator: "c0ffee00",
  });
  await commitWorkerResult(prepared, "quoted worktree command");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "pass");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  const review = await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  const testCommand = review.commands!.tests[0]!;
  assert.equal(testCommand.cwd, prepared.worktree);
  assert.deepEqual(testCommand.argv, ["sh", "-lc", recorded]);
  assert.match(testCommand.shell, /^cd -- '/);
  assert.match(testCommand.shell, /'"'"'/);
  const { execFile } = await import("node:child_process");
  const output = await new Promise<string>((resolvePromise, reject) => execFile("sh", ["-lc", testCommand.shell], { encoding: "utf8" }, (error, stdout) => error ? reject(error) : resolvePromise(stdout)));
  assert.equal(output, "a b\n");
});

test("publication rejects a symlinked review preparation before manifest mutation", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await git(workspace.repository, ["remote", "add", "origin", "https://example.invalid/project.git"]);
  const prepared = await preparePlanlessTask({ workspaceRoot: workspace.root, ...taskOptions, now: new Date("2026-08-11T21:40:00.000Z"), discriminator: "1eadf11e" });
  await commitWorkerResult(prepared, "symlink publication guard");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "pass");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const originalPreparation = manifest.repositories[0].review_preparation as string;
  const symlinkPath = join(dirname(originalPreparation), "frontend-symlinked-review.json");
  await symlink(originalPreparation, symlinkPath);
  manifest.repositories[0].review_preparation = symlinkPath;
  await writeJson(prepared.manifest, manifest);
  const before = await readFile(prepared.manifest, "utf8");
  await assert.rejects(recordReviewPublication({
    workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", status: "published", tool: "gh",
    pullRequest: "https://example.invalid/project/pull/9", evidence: "Forged through symlink", authorized: true,
  }), /Review preparation must be a regular file/);
  assert.equal(await readFile(prepared.manifest, "utf8"), before);
  assert.equal(await git(prepared.worktree, ["status", "--porcelain=v1"]), "");
});

test("merge confirmation rejects a symlinked publication record before manifest mutation", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await git(workspace.repository, ["remote", "add", "origin", "https://example.invalid/project.git"]);
  const prepared = await preparePlanlessTask({ workspaceRoot: workspace.root, ...taskOptions, now: new Date("2026-08-11T21:50:00.000Z"), discriminator: "f11e1ead" });
  await commitWorkerResult(prepared, "publication symlink guard");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  await writeVerifier(prepared, "pass");
  await recordResult({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  await prepareReview({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend" });
  await recordReviewPublication({
    workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", status: "published", tool: "gh",
    pullRequest: "https://example.invalid/project/pull/10", evidence: "Confirmed publication", authorized: true,
  });
  await git(workspace.repository, ["-c", "user.name=Merger", "-c", "user.email=merger@example.invalid", "merge", "--no-ff", "-m", "merge: publication symlink guard", prepared.branch]);
  const mergeCommit = await git(workspace.repository, ["rev-parse", "HEAD"]);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const originalPublication = manifest.repositories[0].review_publication as string;
  const symlinkPath = join(dirname(originalPublication), "frontend-symlinked-publication.json");
  await symlink(originalPublication, symlinkPath);
  manifest.repositories[0].review_publication = symlinkPath;
  await writeJson(prepared.manifest, manifest);
  const before = await readFile(prepared.manifest, "utf8");
  await assert.rejects(confirmMerge({
    workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", mergeCommit, author: "kao", evidence: "Human confirmed merge",
  }), /Review publication record must be a regular file/);
  assert.equal(await readFile(prepared.manifest, "utf8"), before);
});
