import assert from "node:assert/strict"
import { mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises"
import { tmpdir } from "node:os"
import { join } from "node:path"
import { test } from "node:test"
import { parse as parseYaml } from "yaml"
import { acquirePlanLease, planWorktreePath, releasePlanLease, writeSessionRecord } from "../scripts/lib/runtime.js"
import type { SessionRecord } from "../scripts/lib/types.js"

const root = process.cwd()

async function temporaryRuntime(): Promise<{ root: string; cleanup: () => Promise<void> }> {
  const runtimeRoot = await mkdtemp(join(tmpdir(), "context-circuit-runtime-"))
  return { root: runtimeRoot, cleanup: () => rm(runtimeRoot, { recursive: true, force: true }) }
}

test("fresh entry exposes the agent workflow and runtime contract", async () => {
  const workflow = await readFile(join(root, "docs", "agent-workspace-workflow.md"), "utf8")
  const runtime = await readFile(join(root, "docs", "runtime-contract.md"), "utf8")
  const skill = await readFile(join(root, ".agents", "skills", "cc-session-entry", "SKILL.md"), "utf8")
  assert.match(workflow, /Start or resume work in this workspace/)
  assert.match(workflow, /There is no single global/)
  assert.match(runtime, /exclusive-create directory/)
  assert.match(runtime, /handoff\.md/)
  assert.match(skill, /There is no global current-session file/)
  assert.match(skill, /Do not silently approve plans/)
})

test("root and child session records preserve explicit identity and scope", async () => {
  const fixture = await temporaryRuntime()
  try {
    const rootSession: SessionRecord = {
      schema_version: 1,
      session_id: "sess-root",
      parent_session_id: null,
      root_session_id: "sess-root",
      kind: "root",
      role: "coordinator",
      objective: "Coordinate the approved workflow",
      scope: "workspace",
      write_access: false,
      status: "executing",
      created_at: "2026-08-14T12:00:00Z",
      updated_at: "2026-08-14T12:00:00Z",
      next_action: "Delegate the first task",
      blockers: [],
    }
    const childSession: SessionRecord = {
      ...rootSession,
      session_id: "sess-child",
      parent_session_id: "sess-root",
      root_session_id: "sess-root",
      kind: "subagent",
      role: "implementer",
      objective: "Implement one bounded task",
      plan: "plans/context-circuit-plans/0001-agent-workspace-workflow",
      task: "AWF-0001",
      scope: { paths: ["AGENTS.md", "WORKFLOW.md"] },
      write_access: true,
      worktree: ".runtime/worktrees/context-circuit/agent-workspace-workflow",
    }
    await writeSessionRecord(fixture.root, rootSession)
    await writeSessionRecord(fixture.root, childSession)
    const stored = parseYaml(await readFile(join(fixture.root, ".runtime", "sessions", "sess-child", "session.yaml"), "utf8")) as SessionRecord
    assert.equal(stored.parent_session_id, "sess-root")
    assert.equal(stored.root_session_id, "sess-root")
    assert.equal(stored.task, "AWF-0001")
    assert.deepEqual(stored.scope, { paths: ["AGENTS.md", "WORKFLOW.md"] })
    assert.equal(stored.write_access, true)
    assert.equal(await readFile(join(fixture.root, ".runtime", "sessions", "sess-root", "session.yaml"), "utf8").then(Boolean), true)
  } finally {
    await fixture.cleanup()
  }
})

test("child delegation requires bounded scope, permissions, and stop conditions", async () => {
  const fixture = await temporaryRuntime()
  try {
    const directory = join(fixture.root, ".runtime", "sessions", "sess-child")
    await mkdir(directory, { recursive: true })
    const packet = {
      schema_version: 1,
      session_id: "sess-child",
      parent_session_id: "sess-root",
      root_session_id: "sess-root",
      role: "implementer",
      objective: "Implement AWF-0001",
      scope: { plan: "agent-workspace-workflow", task: "AWF-0001", paths: ["AGENTS.md"] },
      non_goals: ["Do not change plan status"],
      context_refs: ["AGENTS.md", "WORKFLOW.md"],
      repository: "context-circuit",
      worktree: ".runtime/worktrees/context-circuit/agent-workspace-workflow",
      permissions: { write_worktree: true, write_runtime_session: true, write_plan: false, write_activity: false },
      acceptance_criteria: ["The contract is internally consistent"],
      stop_conditions: ["Required change falls outside scope"],
      handoff_schema: "session-handoff-v1",
    }
    await writeFile(join(directory, "delegation.yaml"), "schema_version: 1\n" + Object.entries(packet).slice(1).map(([key, value]) => key + ": " + JSON.stringify(value)).join("\n") + "\n", "utf8")
    const raw = await readFile(join(directory, "delegation.yaml"), "utf8")
    const parsed = parseYaml(raw) as Record<string, any>
    assert.match(raw, /parent_session_id/)
    assert.match(raw, /write_plan/)
    assert.match(raw, /stop_conditions/)
    assert.match(raw, /handoff_schema/)
    assert.equal(parsed.scope.task, "AWF-0001")
    assert.equal(parsed.permissions.write_plan, false)
  } finally {
    await fixture.cleanup()
  }
})

test("same-plan lease contention is non-destructive while different plans proceed", async () => {
  const fixture = await temporaryRuntime()
  try {
    const attempts = await Promise.allSettled([
      acquirePlanLease({
        workspaceRoot: fixture.root,
        plan: "plan-alpha",
        sessionId: "sess-alpha",
        worktree: planWorktreePath(fixture.root, "app", "plan-alpha"),
      }),
      acquirePlanLease({
        workspaceRoot: fixture.root,
        plan: "plan-alpha",
        sessionId: "sess-contender",
        worktree: planWorktreePath(fixture.root, "app", "plan-alpha"),
      }),
    ])
    assert.equal(attempts.filter((attempt) => attempt.status === "fulfilled").length, 1)
    assert.equal(attempts.filter((attempt) => attempt.status === "rejected").length, 1)
    const first = attempts.find((attempt): attempt is PromiseFulfilledResult<Awaited<ReturnType<typeof acquirePlanLease>>> => attempt.status === "fulfilled")!.value
    const rejected = attempts.find((attempt): attempt is PromiseRejectedResult => attempt.status === "rejected")!.reason as Error
    assert.match(rejected.message, /Plan lease contention|lease exists but is incomplete/)
    assert.equal(first.status, "active")
    const [secondPlan, thirdPlan] = await Promise.all([
      acquirePlanLease({
        workspaceRoot: fixture.root,
        plan: "plan-beta",
        sessionId: "sess-beta",
        worktree: planWorktreePath(fixture.root, "app", "plan-beta"),
      }),
      acquirePlanLease({
        workspaceRoot: fixture.root,
        plan: "plan-gamma",
        sessionId: "sess-gamma",
        worktree: planWorktreePath(fixture.root, "app", "plan-gamma"),
      }),
    ])
    assert.equal(secondPlan.session_id, "sess-beta")
    assert.equal(thirdPlan.session_id, "sess-gamma")
    assert.notEqual(first.worktree, secondPlan.worktree)
    const firstLease = parseYaml(await readFile(join(fixture.root, ".runtime", "plans", "plan-alpha", "lease.yaml"), "utf8")) as Record<string, unknown>
    assert.ok(["sess-alpha", "sess-contender"].includes(String(firstLease.session_id)))
    await releasePlanLease({ workspaceRoot: fixture.root, plan: "plan-alpha", sessionId: first.session_id })
    await releasePlanLease({ workspaceRoot: fixture.root, plan: "plan-beta", sessionId: "sess-beta" })
    await releasePlanLease({ workspaceRoot: fixture.root, plan: "plan-gamma", sessionId: "sess-gamma" })
  } finally {
    await fixture.cleanup()
  }
})

test("invalid runtime identifiers are rejected before filesystem writes", async () => {
  const fixture = await temporaryRuntime()
  try {
    await assert.rejects(
      () => acquirePlanLease({
        workspaceRoot: fixture.root,
        plan: "../outside",
        sessionId: "sess-safe",
        worktree: ".runtime/worktrees/app/outside",
      }),
      /safe runtime identifier/
    )
  } finally {
    await fixture.cleanup()
  }
})

test("interrupted sessions retain handoff evidence for explicit recovery", async () => {
  const fixture = await temporaryRuntime()
  try {
    const session: SessionRecord = {
      schema_version: 1,
      session_id: "sess-interrupted",
      parent_session_id: "sess-root",
      root_session_id: "sess-root",
      kind: "subagent",
      role: "implementer",
      objective: "Implement a bounded task",
      scope: "assigned task",
      write_access: true,
      worktree: ".runtime/worktrees/app/plan-alpha",
      status: "blocked",
      created_at: "2026-08-14T12:00:00Z",
      updated_at: "2026-08-14T12:20:00Z",
      next_action: "Request explicit recovery",
      blockers: ["worker interrupted before verification"],
    }
    await writeSessionRecord(fixture.root, session)
    const sessionDirectory = join(fixture.root, ".runtime", "sessions", session.session_id)
    await writeFile(join(sessionDirectory, "handoff.md"), "# Session handoff\n\nstatus: blocked\n\nThe worker stopped before verification.\n", "utf8")
    assert.match(await readFile(join(sessionDirectory, "handoff.md"), "utf8"), /stopped before verification/)
    assert.equal((parseYaml(await readFile(join(sessionDirectory, "session.yaml"), "utf8")) as SessionRecord).status, "blocked")
  } finally {
    await fixture.cleanup()
  }
})

test("verification and completion remain explicit human-gated outcomes", async () => {
  const workflow = await readFile(join(root, "docs", "agent-workspace-workflow.md"), "utf8")
  const worker = await readFile(join(root, "agents", "repository-worker.md"), "utf8")
  const reviewer = await readFile(join(root, "agents", "reviewer.md"), "utf8")
  const cleanupGate = await readFile(join(root, "docs", "legacy-cleanup-gate.md"), "utf8")
  assert.match(workflow, /Verification fails/)
  assert.match(workflow, /Human approval is required/)
  assert.match(workflow, /source changes during execution/)
  assert.match(worker, /approve work/)
  assert.match(reviewer, /Remain read-only/)
  assert.match(reviewer, /Do not repair/)
  assert.match(cleanupGate, /separate,\s+human-reviewed cleanup plan/)
  assert.match(cleanupGate, /No file deletion/)
})

test("legacy command implementation remains available during migration", async () => {
  const command = await readFile(join(root, ".agents", "bin", "cc.mjs"), "utf8")
  const runner = await readFile(join(root, "scripts", "lib", "run-task.ts"), "utf8")
  assert.match(command, /run-task/)
  assert.match(runner, /acquirePlanLease/)
})
