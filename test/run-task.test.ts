import assert from "node:assert/strict";
import { access, mkdir, readFile, symlink, writeFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { preparePlanlessTask } from "../scripts/lib/run-task.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace, taskOptions } from "./helpers.js";

test("planless preparation writes valid scoped evidence and creates a worktree", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    testScope: ["src/App.test.tsx"],
    now: new Date("2026-08-11T08:30:00.000Z"),
    discriminator: "abcd1234",
  });
  assert.equal(prepared.workId, "ADHOC-20260811-001");
  assert.match(prepared.runId, /^20260811T083000Z-abcd[0-9a-f]{4}$/);
  await access(join(prepared.worktree, ".git"));

  const brief = JSON.parse(await readFile(prepared.taskBrief, "utf8"));
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  const worker = JSON.parse(await readFile(prepared.workerInput, "utf8"));
  const verifier = JSON.parse(await readFile(prepared.verifierInput, "utf8"));
  assert.deepEqual(await validateContract("task-brief", brief), []);
  assert.deepEqual(await validateContract("runtime-manifest", manifest), []);
  assert.equal(manifest.status, "prepared");
  assert.equal(brief.plan.approval_state, "not-applicable");
  assert.equal(brief.activity.duplicate_effort_warning, true);
  assert.deepEqual(brief.implementation_scope, ["src/App.tsx"]);
  assert.deepEqual(brief.test_expectation, {
    policy: "required",
    paths: ["src/App.test.tsx"],
    rationale: "The worker must add or update tests in the declared test scope.",
  });
  assert.deepEqual(worker.allowed_scope, ["src/App.tsx", "src/App.test.tsx"]);
  assert.deepEqual(worker.test_expectation, brief.test_expectation);
  await access(worker.result_contract);
  assert.equal(verifier.read_only, true);
  assert.deepEqual(verifier.test_expectation, brief.test_expectation);
  await access(verifier.result_contract);
});

test("dirty repository protection preserves changes and creates no run", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const dirtyPath = join(workspace.repository, "unrecorded.txt");
  await writeFile(dirtyPath, "preserve me\n", "utf8");
  await assert.rejects(
    preparePlanlessTask({ workspaceRoot: workspace.root, ...taskOptions }),
    /unresolved local changes/,
  );
  assert.equal(await readFile(dirtyPath, "utf8"), "preserve me\n");
  await assert.rejects(access(join(workspace.root, ".runtime")));
});

test("runtime evidence refuses a symlinked runtime root", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const target = join(workspace.root, "evidence-target");
  await mkdir(target);
  await symlink(target, join(workspace.root, ".runtime"), "dir");
  await assert.rejects(
    preparePlanlessTask({ workspaceRoot: workspace.root, ...taskOptions }),
    /Runtime path must be a real directory/,
  );
  await assert.rejects(access(join(target, "id-state.json")));
});

test("existing-coverage policy rejects a missing test path", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await assert.rejects(
    preparePlanlessTask({
      workspaceRoot: workspace.root,
      ...taskOptions,
      testScope: ["src/missing.test.tsx"],
      testPolicy: "existing-coverage",
    }),
    /Existing-coverage test path does not exist/,
  );
});

test("separate preparations use isolated worktrees and branches", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const first = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T09:00:00.000Z"),
    discriminator: "11111111",
  });
  const second = await preparePlanlessTask({
    workspaceRoot: workspace.root,
    ...taskOptions,
    now: new Date("2026-08-11T09:00:01.000Z"),
    discriminator: "22222222",
  });
  assert.notEqual(first.branch, second.branch);
  assert.notEqual(first.worktree, second.worktree);
  assert.equal(first.workId, "ADHOC-20260811-001");
  assert.equal(second.workId, "ADHOC-20260811-002");
  await writeFile(join(first.worktree, "isolation.txt"), "first only\n", "utf8");
  await assert.rejects(access(join(second.worktree, "isolation.txt")));
});
