import assert from "node:assert/strict";
import { access, readFile, writeFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { prepareActivityLifecycle, recordActivityLifecycleAction } from "../scripts/lib/activity-lifecycle.js";
import { preparePlanlessTask, resumePlanlessTask } from "../scripts/lib/run-task.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace, taskOptions } from "./helpers.js";

async function configureActivity(root: string, lifecycle: string): Promise<void> {
  const path = join(root, "workspace.yaml");
  const raw = await readFile(path, "utf8");
  await writeFile(path, raw
    .replace("provider: none", "provider: fake-session")
    .replace("required_capabilities: []", "required_capabilities: [update-status]")
    .replace("optional_capabilities: []", `optional_capabilities: [timers]\n  lifecycle:\n${lifecycle}`), "utf8");
}

async function preparedRun(root: string, discriminator: string) {
  return preparePlanlessTask({ workspaceRoot: root, ...taskOptions, discriminator, now: new Date("2026-08-12T02:00:00Z") });
}

test("lifecycle preparation classifies pending, optional, and manual actions deterministically", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparedRun(workspace.root, "1afe0001");
  await configureActivity(workspace.root, `    task.completed:
      - id: set-completed
        capability: update-status
        policy: required
        description: Set the external task status to Completed.
      - id: stop-timer
        capability: timers
        policy: optional
        description: Stop the external task timer.
      - id: add-merge-note
        capability: update-status
        policy: manual
        description: Add the merge reference to the external task.
`);

  const record = await prepareActivityLifecycle({
    workspaceRoot: workspace.root,
    runId: prepared.runId,
    event: "task.completed",
    availableCapabilities: ["update-status"],
    now: new Date("2026-08-12T02:10:00Z"),
  });

  assert.equal(record.status, "pending");
  assert.deepEqual(record.actions.map((action) => action.status), ["pending", "skipped", "manual"]);
  assert.deepEqual(record.manual_fallbacks, ["Stop the external task timer.", "Add the merge reference to the external task."]);
  assert.deepEqual(await validateContract("activity-lifecycle-record", record), []);
  const repeated = await prepareActivityLifecycle({ workspaceRoot: workspace.root, runId: prepared.runId, event: "task.completed" });
  assert.deepEqual(repeated, record);
  const manifest = JSON.parse(await readFile(prepared.manifest, "utf8"));
  assert.equal(manifest.lifecycle_events.find((event: { event: string }) => event.event === "task.completed").status, "pending");
  assert.deepEqual(await validateContract("runtime-manifest", manifest), []);
});

test("recording confirmed host and manual results completes the event idempotently", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparedRun(workspace.root, "1afe0002");
  await configureActivity(workspace.root, `    task.review-ready:
      - id: set-review
        capability: update-status
        policy: required
        description: Set the external task status to Review.
      - id: attach-pr
        capability: update-status
        policy: manual
        description: Attach the pull-request reference manually.
`);
  await prepareActivityLifecycle({ workspaceRoot: workspace.root, runId: prepared.runId, event: "task.review-ready", availableCapabilities: ["update-status"] });
  const first = await recordActivityLifecycleAction({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.review-ready", actionId: "set-review",
    status: "completed", evidence: "Status transition confirmed by the session tool.", externalReference: "TASK-42",
    now: new Date("2026-08-12T02:11:00Z"),
  });
  assert.equal(first.status, "manual");
  const completed = await recordActivityLifecycleAction({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.review-ready", actionId: "attach-pr",
    status: "completed", evidence: "Human attached frontend#42.", externalReference: "frontend#42",
    now: new Date("2026-08-12T02:12:00Z"),
  });
  assert.equal(completed.status, "completed");
  const repeated = await recordActivityLifecycleAction({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.review-ready", actionId: "attach-pr",
    status: "completed", evidence: "Human attached frontend#42.", externalReference: "frontend#42",
  });
  assert.deepEqual(repeated, completed);
});

test("missing required capability becomes a manual gate and stops later actions", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparedRun(workspace.root, "1afe0003");
  await configureActivity(workspace.root, `    task.blocked:
      - id: set-blocked
        capability: update-status
        policy: required
        description: Set the external task status to Blocked.
      - id: pause-timer
        capability: timers
        policy: optional
        description: Pause the external task timer.
`);
  const missing = await prepareActivityLifecycle({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.blocked", availableCapabilities: ["timers"],
  });
  assert.equal(missing.status, "manual");
  assert.deepEqual(missing.actions.map((action) => action.status), ["manual", "skipped"]);
  assert.match(missing.actions[1]!.evidence ?? "", /required lifecycle action failed/);
});

test("optional capability absence degrades cleanly with an exact manual fallback", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparedRun(workspace.root, "1afe0004");
  await configureActivity(workspace.root, `    task.cancelled:
      - id: stop-timer
        capability: timers
        policy: optional
        description: Stop the timer and record the cancellation reason.
`);
  const record = await prepareActivityLifecycle({ workspaceRoot: workspace.root, runId: prepared.runId, event: "task.cancelled" });
  assert.equal(record.status, "completed");
  assert.equal(record.actions[0]?.status, "skipped");
  assert.deepEqual(record.manual_fallbacks, ["Stop the timer and record the cancellation reason."]);
  assert.match(record.warnings.join("\n"), /Optional capability is unavailable/);
});

test("lifecycle evidence rejects credential-bearing references", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparedRun(workspace.root, "1afe0005");
  await configureActivity(workspace.root, `    task.completed:
      - id: complete
        capability: update-status
        policy: required
        description: Complete the external task.
`);
  await prepareActivityLifecycle({ workspaceRoot: workspace.root, runId: prepared.runId, event: "task.completed", availableCapabilities: ["update-status"] });
  await assert.rejects(recordActivityLifecycleAction({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.completed", actionId: "complete",
    status: "completed", evidence: "Completed.", externalReference: "https://user:secret@example.invalid/task",
  }), /credential/);
});

test("configured starting hooks complete before worktree creation and resume revalidates the base", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await configureActivity(workspace.root, `    task.starting:
      - id: refresh-and-claim
        capability: update-status
        policy: required
        description: Refresh ownership and set the external task to In Progress.
`);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root, ...taskOptions, discriminator: "1afe0006",
    now: new Date("2026-08-12T02:00:00Z"), availableCapabilities: ["update-status"],
  });
  assert.equal(prepared.preparationStatus, "awaiting-activity");
  await assert.rejects(access(prepared.worktree));
  await recordActivityLifecycleAction({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.starting", actionId: "refresh-and-claim",
    status: "completed", evidence: "Ownership refreshed and claim confirmed.", externalReference: "TASK-99",
  });
  const resumed = await resumePlanlessTask({ workspaceRoot: workspace.root, runId: prepared.runId, now: new Date("2026-08-12T02:05:00Z") });
  assert.equal(resumed.preparationStatus, "prepared");
  await access(resumed.worktree);
  const manifest = JSON.parse(await readFile(resumed.manifest, "utf8"));
  assert.equal(manifest.status, "prepared");
  assert.equal(manifest.lifecycle_events.find((event: { event: string }) => event.event === "task.starting").status, "completed");
});

test("unavailable required starting capability waits for manual completion before worktree creation", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await configureActivity(workspace.root, `    task.starting:
      - id: claim-task
        capability: update-status
        policy: required
        description: Claim the external task before work starts.
`);
  const prepared = await preparePlanlessTask({
    workspaceRoot: workspace.root, ...taskOptions, discriminator: "1afe0007", now: new Date("2026-08-12T02:00:00Z"),
  });
  assert.equal(prepared.preparationStatus, "awaiting-activity");
  await assert.rejects(access(prepared.worktree));
  const branches = await import("../scripts/lib/git.js").then(({ git }) => git(workspace.repository, ["branch", "--list", prepared.branch]));
  assert.equal(branches, "");
  await recordActivityLifecycleAction({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.starting", actionId: "claim-task",
    status: "completed", evidence: "Human confirmed the task claim manually.", externalReference: "TASK-100",
  });
  const resumed = await resumePlanlessTask({ workspaceRoot: workspace.root, runId: prepared.runId });
  assert.equal(resumed.preparationStatus, "prepared");
});

test("confirmed failure of a required action stops remaining pending actions", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const prepared = await preparedRun(workspace.root, "1afe0008");
  await configureActivity(workspace.root, `    task.blocked:
      - id: set-blocked
        capability: update-status
        policy: required
        description: Set the external task status to Blocked.
      - id: pause-timer
        capability: timers
        policy: optional
        description: Pause the external task timer.
`);
  await prepareActivityLifecycle({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.blocked", availableCapabilities: ["update-status", "timers"],
  });
  const failed = await recordActivityLifecycleAction({
    workspaceRoot: workspace.root, runId: prepared.runId, event: "task.blocked", actionId: "set-blocked",
    status: "failed", evidence: "The session tool rejected the status transition.",
  });
  assert.equal(failed.status, "failed");
  assert.deepEqual(failed.actions.map((action) => action.status), ["failed", "skipped"]);
});
