import { mkdir, readFile, rm, stat } from "node:fs/promises"
import { join, resolve } from "node:path"
import { parse as parseYaml, stringify as stringifyYaml } from "yaml"
import { writeTextAtomic } from "./io.js"
import type { PlanLease, SessionRecord } from "./types.js"

const identifierPattern = /^[a-z0-9][a-z0-9._-]{0,63}$/
const leaseStaleAfterSeconds = 1800

function assertIdentifier(value: string, label: string): string {
  if (!identifierPattern.test(value)) throw new Error(label + " is not a safe runtime identifier: " + value)
  return value
}

function runtimePath(root: string, ...parts: string[]): string {
  return resolve(root, ".runtime", ...parts)
}

export function sessionRuntimeDirectory(root: string, sessionId: string): string {
  return runtimePath(root, "sessions", assertIdentifier(sessionId, "session_id"))
}

export function planRuntimeDirectory(root: string, planId: string): string {
  return runtimePath(root, "plans", assertIdentifier(planId, "plan_id"))
}

export function planWorktreePath(root: string, repository: string, planId: string): string {
  return runtimePath(root, "worktrees", assertIdentifier(repository, "repository"), assertIdentifier(planId, "plan_id"))
}

export async function readSessionRecord(root: string, sessionId: string): Promise<SessionRecord> {
  const path = join(sessionRuntimeDirectory(root, sessionId), "session.yaml")
  return parseYaml(await readFile(path, "utf8")) as SessionRecord
}

export async function writeSessionRecord(root: string, record: SessionRecord): Promise<void> {
  const directory = sessionRuntimeDirectory(root, record.session_id)
  await mkdir(directory, { recursive: true })
  await writeTextAtomic(join(directory, "session.yaml"), stringifyYaml(record))
}

export async function acquirePlanLease(options: {
  workspaceRoot: string
  plan: string
  sessionId: string
  rootSessionId?: string
  worktree: string
}): Promise<PlanLease> {
  const planDirectory = planRuntimeDirectory(options.workspaceRoot, options.plan)
  const lockDirectory = join(planDirectory, "lease.lock")
  const ownerPath = join(lockDirectory, "owner.yaml")
  const leasePath = join(planDirectory, "lease.yaml")
  const now = new Date().toISOString()
  const lease: PlanLease = {
    schema_version: 1,
    plan: options.plan,
    session_id: assertIdentifier(options.sessionId, "session_id"),
    root_session_id: assertIdentifier(options.rootSessionId ?? options.sessionId, "root_session_id"),
    worktree: options.worktree,
    status: "active",
    acquired_at: now,
    heartbeat_at: now,
    released_at: null,
    stale_after_seconds: leaseStaleAfterSeconds,
  }
  await mkdir(planDirectory, { recursive: true })
  try {
    await mkdir(lockDirectory)
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "EEXIST") throw error
    const existing = await readLease(leasePath)
    if (!existing) throw new Error("Plan lease exists but is incomplete; refusing to infer ownership: " + options.plan)
    if (existing.session_id !== lease.session_id) {
      if (leaseIsStale(existing)) {
        throw new Error("Plan lease is stale and requires explicit human takeover for " + options.plan + ": owned by " + existing.session_id)
      }
      throw new Error("Plan lease contention for " + options.plan + ": owned by " + existing.session_id)
    }
    if (existing.status !== "active") {
      throw new Error("Plan lease is not active for session " + existing.session_id + ": " + existing.status)
    }
    if (existing.worktree !== options.worktree) {
      throw new Error("Plan lease worktree mismatch for " + options.plan + ": owned path is " + existing.worktree)
    }
    const refreshed = { ...existing, heartbeat_at: now }
    await writeTextAtomic(leasePath, stringifyYaml(refreshed))
    return refreshed
  }
  try {
    await writeTextAtomic(ownerPath, stringifyYaml({
      schema_version: 1,
      session_id: lease.session_id,
      root_session_id: lease.root_session_id,
      plan: lease.plan,
      acquired_at: lease.acquired_at,
    }))
    await writeTextAtomic(leasePath, stringifyYaml(lease))
    return lease
  } catch (error) {
    await rm(lockDirectory, { recursive: true, force: true })
    throw error
  }
}

export async function readLease(path: string): Promise<PlanLease | null> {
  try {
    return parseYaml(await readFile(path, "utf8")) as PlanLease
  } catch {
    return null
  }
}

export function leaseIsStale(lease: PlanLease, now = new Date()): boolean {
  if (lease.status !== "active") return false
  const heartbeat = Date.parse(lease.heartbeat_at)
  if (!Number.isFinite(heartbeat)) return true
  return now.getTime() - heartbeat > lease.stale_after_seconds * 1000
}

export async function releasePlanLease(options: {
  workspaceRoot: string
  plan: string
  sessionId: string
}): Promise<PlanLease> {
  const planDirectory = planRuntimeDirectory(options.workspaceRoot, options.plan)
  const lockDirectory = join(planDirectory, "lease.lock")
  const leasePath = join(planDirectory, "lease.yaml")
  const existing = await readLease(leasePath)
  if (!existing) throw new Error("Cannot release missing plan lease: " + options.plan)
  if (existing.session_id !== options.sessionId) throw new Error("Cannot release a plan lease owned by " + existing.session_id)
  const released = { ...existing, status: "released" as const, released_at: new Date().toISOString(), heartbeat_at: new Date().toISOString() }
  await writeTextAtomic(leasePath, stringifyYaml(released))
  await rm(lockDirectory, { recursive: true, force: false })
  return released
}

export async function runtimePathExists(path: string): Promise<boolean> {
  try {
    await stat(path)
    return true
  } catch {
    return false
  }
}
