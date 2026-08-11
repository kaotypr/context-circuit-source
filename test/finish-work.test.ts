import assert from "node:assert/strict";
import { access, mkdir, readFile, readdir, symlink, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import test from "node:test";
import { prepareActivityLifecycle, recordActivityLifecycleAction } from "../scripts/lib/activity-lifecycle.js";
import { contributionDocumentErrors, finishWork } from "../scripts/lib/finish-work.js";
import { git } from "../scripts/lib/git.js";
import { recordResult } from "../scripts/lib/record-result.js";
import { preparePlanlessTask, type PreparedTask } from "../scripts/lib/run-task.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace, taskOptions } from "./helpers.js";

async function writeJson(path: string, value: unknown): Promise<void> {
  await mkdir(dirname(path), { recursive: true });
  await writeFile(path, `${JSON.stringify(value, null, 2)}\n`, "utf8");
}

async function initializeWrapper(root: string): Promise<void> {
  await writeFile(join(root, ".gitignore"), ".runtime/\nrepositories/frontend/\n", "utf8");
  await git(root, ["init", "--initial-branch=main"]);
  await git(root, ["add", ".gitignore", "workspace.yaml", "AGENTS.md", "agents", ".agents"]);
  await git(root, ["-c", "user.name=Wrapper", "-c", "user.email=wrapper@example.invalid", "commit", "-m", "test: initialize wrapper"]);
}

async function passRun(root: string, discriminator: string): Promise<PreparedTask> {
  const prepared = await preparePlanlessTask({
    workspaceRoot: root,
    ...taskOptions,
    now: new Date("2026-08-11T22:00:00.000Z"),
    discriminator,
  });
  const appPath = join(prepared.worktree, "src", "App.tsx");
  await writeFile(appPath, `${await readFile(appPath, "utf8")}\n// closeout fixture\n`, "utf8");
  await git(prepared.worktree, ["add", "src/App.tsx"]);
  await git(prepared.worktree, ["-c", "user.name=Worker", "-c", "user.email=worker@example.invalid", "commit", "-m", "test: closeout fixture"]);
  const head = await git(prepared.worktree, ["rev-parse", "HEAD"]);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const workerInput = JSON.parse(await readFile(prepared.workerInput, "utf8"));
  await writeJson(workerInput.result_path, {
    contract_version: 1,
    work_id: prepared.workId,
    run_id: prepared.runId,
    repository: "frontend",
    status: "completed",
    summary: "Implemented the closeout fixture.",
    branch: prepared.branch,
    worktree: prepared.worktree,
    commits: [head],
    changed_files: ["src/App.tsx"],
    checks: [{ command: "npm test", status: "passed", evidence: "Fixture evidence" }],
    risks: [],
  });
  await recordResult({ workspaceRoot: root, runId: prepared.runId, repository: "frontend", stage: "worker-result" });
  const verifierInput = JSON.parse(await readFile(manifest.repositories[0].verifier_input, "utf8"));
  await writeJson(verifierInput.result_path, {
    contract_version: 1,
    work_id: prepared.workId,
    run_id: prepared.runId,
    repository: "frontend",
    status: "pass",
    summary: "Independent verification passed.",
    acceptance: [{ criterion: taskOptions.acceptanceCriteria[0], status: "passed", evidence: "Behavior inspected." }],
    checks: ["npm test passed"],
    findings: [],
    verified_at: "2026-08-11T22:03:00.000Z",
  });
  await recordResult({ workspaceRoot: root, runId: prepared.runId, repository: "frontend", stage: "verifier-result" });
  return prepared;
}

async function commitContribution(root: string, contribution: string): Promise<void> {
  await git(root, ["add", contribution]);
  await git(root, ["-c", "user.name=Closer", "-c", "user.email=closer@example.invalid", "commit", "-m", "docs: record contribution"]);
}

test("prepares one append-only contribution and safely removes only a merged clean worktree", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await initializeWrapper(workspace.root);
  const prepared = await passRun(workspace.root, "c10c10c1");
  await git(workspace.repository, ["merge", "--no-ff", "-m", "merge: closeout fixture", prepared.branch]);

  const closeout = await finishWork({
    workspaceRoot: workspace.root,
    runId: prepared.runId,
    repository: "frontend",
    outcome: "merged",
    author: "Kao Typr",
    pullRequests: ["frontend#42"],
    now: new Date("2026-08-11T22:10:00.000Z"),
  });
  assert.equal(closeout.status, "prepared");
  assert.equal(closeout.author, "kao-typr");
  assert.deepEqual(await validateContract("closeout-record", closeout), []);
  const contribution = await readFile(join(workspace.root, closeout.contribution), "utf8");
  assert.deepEqual(contributionDocumentErrors(closeout.contribution, contribution, prepared.runId), []);
  assert.match(contribution, /frontend#42/);
  assert.doesNotMatch(contribution, /raw conversation/i);
  const repeated = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "kao-typr" });
  assert.deepEqual(repeated, closeout);
  assert.equal((await readdir(join(workspace.root, "contributions", "general"))).filter((name) => name.endsWith(".md")).length, 1);

  const blocked = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "kao-typr", cleanup: true });
  assert.equal(blocked.status, "blocked");
  assert.match(blocked.blockers.join("\n"), /durably tracked/i);
  await access(prepared.worktree);

  await commitContribution(workspace.root, closeout.contribution);
  const closed = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "kao-typr", cleanup: true, now: new Date("2026-08-11T22:20:00.000Z") });
  assert.equal(closed.status, "closed");
  assert.equal(closed.cleanup.worktree_removed, true);
  await assert.rejects(access(prepared.worktree));
  assert.equal(await git(workspace.repository, ["show-ref", "--verify", `refs/heads/${prepared.branch}`]).then(() => true), true);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.status, "closed");
  assert.deepEqual(manifest.execution_events.slice(-2).map((event: { stage: string }) => event.stage), ["closeout-prepared", "closeout-cleaned"]);
  assert.equal(manifest.lifecycle_events.at(-1).event, "task.completed");
  assert.equal(manifest.lifecycle_events.at(-1).status, "skipped");
  assert.deepEqual(await validateContract("runtime-manifest", manifest), []);
});

test("preserves an abandoned worktree when its commits are unpushed", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await initializeWrapper(workspace.root);
  const prepared = await passRun(workspace.root, "abadd00d");
  const closeout = await finishWork({
    workspaceRoot: workspace.root,
    runId: prepared.runId,
    repository: "frontend",
    outcome: "abandoned",
    author: "kao",
    reason: "The experiment is no longer required.",
    now: new Date("2026-08-11T23:10:00.000Z"),
  });
  await commitContribution(workspace.root, closeout.contribution);
  const blocked = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "abandoned", author: "kao", cleanup: true });
  assert.equal(blocked.status, "blocked");
  assert.match(blocked.blockers.join("\n"), /neither merged nor preserved by a remote ref/i);
  await access(prepared.worktree);
  assert.equal(await git(prepared.worktree, ["rev-parse", "HEAD"]), closeout.head_commit);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.status, "closing");
  assert.equal(manifest.lifecycle_events.at(-1).event, "task.cancelled");
});

test("refuses to remove a merged worktree with uncommitted changes", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await initializeWrapper(workspace.root);
  const prepared = await passRun(workspace.root, "d17d17d1");
  await git(workspace.repository, ["merge", "--no-ff", "-m", "merge: dirty closeout fixture", prepared.branch]);
  const closeout = await finishWork({
    workspaceRoot: workspace.root,
    runId: prepared.runId,
    repository: "frontend",
    outcome: "merged",
    author: "kao",
    now: new Date("2026-08-12T00:10:00.000Z"),
  });
  await commitContribution(workspace.root, closeout.contribution);
  await writeFile(join(prepared.worktree, "unrecorded.txt"), "preserve me\n", "utf8");
  const blocked = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "kao", cleanup: true });
  assert.equal(blocked.status, "blocked");
  assert.match(blocked.blockers.join("\n"), /uncommitted changes/i);
  assert.equal(await readFile(join(prepared.worktree, "unrecorded.txt"), "utf8"), "preserve me\n");
});

test("contribution validation rejects incomplete or credential-bearing documents", () => {
  const errors = contributionDocumentErrors(
    "20260811T230000Z-kao-test.md",
    "# Incomplete\n\nhttps://user:secret@example.invalid/repo\n",
    "run-1",
  );
  assert.ok(errors.some((error) => error.includes("missing ## Outcome")));
  assert.ok(errors.some((error) => error.includes("expected run")));
  assert.ok(errors.some((error) => error.includes("credential")));
});

test("stops before closeout mutation when configured activity hooks are unavailable", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await initializeWrapper(workspace.root);
  const prepared = await passRun(workspace.root, "ac710000");
  const configPath = join(workspace.root, "workspace.yaml");
  await writeFile(configPath, (await readFile(configPath, "utf8"))
    .replace("provider: none", "provider: example")
    .replace("required_capabilities: []", "required_capabilities: [update-status]")
    .replace("optional_capabilities: []", `optional_capabilities: []
  lifecycle:
    task.completed:
      - id: complete-task
        capability: update-status
        policy: required
        description: Set the external task to Completed.
`), "utf8");
  await assert.rejects(
    finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "kao" }),
    /Prepare configured activity hooks before closeout/,
  );
  await assert.rejects(access(join(workspace.root, "contributions", "general")));
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.status, "passed");

  const hook = await prepareActivityLifecycle({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.completed", availableCapabilities: ["update-status"],
  });
  assert.equal(hook.status, "pending");
  await recordActivityLifecycleAction({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.completed", actionId: "complete-task",
    status: "completed", evidence: "External task completion confirmed.", externalReference: "TASK-42",
  });
  const closeout = await finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "kao" });
  assert.equal(closeout.status, "prepared");
  const after = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(after.lifecycle_events.find((event: { event: string }) => event.event === "task.completed").status, "completed");
});

test("refuses a symlinked contribution directory", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await initializeWrapper(workspace.root);
  const prepared = await passRun(workspace.root, "5afe0000");
  await mkdir(join(workspace.root, "contributions"));
  await symlink(workspace.repository, join(workspace.root, "contributions", "general"), "dir");
  await assert.rejects(
    finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "kao" }),
    /Contribution path must be a real directory/,
  );
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.status, "passed");
});

test("rejects tampered verifier evidence before writing a contribution", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await initializeWrapper(workspace.root);
  const prepared = await passRun(workspace.root, "e71de0ce");
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const verifierInput = JSON.parse(await readFile(manifest.repositories[0].verifier_input, "utf8"));
  const verifier = JSON.parse(await readFile(verifierInput.result_path, "utf8"));
  verifier.run_id = "20260812T010000Z-deadbeef";
  await writeJson(verifierInput.result_path, verifier);
  await assert.rejects(
    finishWork({ workspaceRoot: workspace.root, runId: prepared.runId, repository: "frontend", outcome: "merged", author: "kao" }),
    /Verifier result identity does not match/,
  );
  await assert.rejects(access(join(workspace.root, "contributions", "general")));
});
