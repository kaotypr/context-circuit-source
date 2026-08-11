import { readFile } from "node:fs/promises";
import { join, relative, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { assertInside, ensurePrivateDirectory, withExclusiveFile, writeJsonAtomic } from "./io.js";
import type {
  ActivityActionResult,
  ActivityCapability,
  ActivityEvent,
  ActivityLifecycleRecord,
  RuntimeManifest,
  WorkspaceConfig,
} from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";

export interface PrepareLifecycleOptions {
  workspaceRoot: string;
  runId: string;
  event: ActivityEvent;
  availableCapabilities?: ActivityCapability[];
  now?: Date;
}

export interface RecordLifecycleActionOptions {
  workspaceRoot: string;
  runId: string;
  event: ActivityEvent;
  actionId: string;
  status: "completed" | "failed";
  evidence: string;
  externalReference?: string;
  now?: Date;
}

async function readJson<T>(path: string): Promise<T> {
  return JSON.parse(await readFile(path, "utf8")) as T;
}

async function assertValid(name: "workspace" | "runtime-manifest" | "activity-lifecycle-record", value: unknown): Promise<void> {
  const errors = await validateContract(name, value);
  if (errors.length > 0) throw new Error(`Invalid ${name}: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}

async function loadConfig(workspaceRoot: string): Promise<WorkspaceConfig> {
  const config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  await assertValid("workspace", config);
  const errors = workspaceSemanticErrors(config);
  if (errors.length > 0) throw new Error(`Invalid workspace: ${errors.join("; ")}`);
  return config;
}

function safeEvidence(value: string, field: string): string {
  const clean = value.trim();
  if (!clean || /[\r\n]/.test(clean)) throw new Error(`${field} must be a non-empty single line`);
  if (/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----|https?:\/\/[^\s/@:]+:[^\s/@]+@/i.test(clean)) {
    throw new Error(`${field} appears to contain a credential or private key`);
  }
  return clean;
}

function overallStatus(actions: ActivityActionResult[]): ActivityLifecycleRecord["status"] {
  if (actions.some((action) => action.policy === "required" && action.status === "failed")) return "failed";
  if (actions.some((action) => action.status === "pending")) return "pending";
  if (actions.some((action) => action.status === "manual")) return "manual";
  return actions.length === 0 ? "skipped" : "completed";
}

function syncManifestEvent(manifest: RuntimeManifest, record: ActivityLifecycleRecord, recordPath: string, workspaceRoot: string): void {
  const existing = manifest.lifecycle_events.find((item) => item.event === record.event);
  const occurredAt = record.updated_at;
  const value = {
    event: record.event,
    status: record.status,
    idempotency_key: existing?.idempotency_key ?? `${manifest.run_id}:lifecycle:${record.event}:${record.provider}`,
    occurred_at: occurredAt,
    record: relative(workspaceRoot, recordPath).replaceAll("\\", "/"),
    actions: record.actions,
  };
  if (existing) Object.assign(existing, value);
  else manifest.lifecycle_events.push(value);
  manifest.updated_at = occurredAt;
}

function paths(workspaceRoot: string, runId: string, event: ActivityEvent): { runtimeRoot: string; manifestPath: string; recordPath: string } {
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime"));
  const runRoot = assertInside(runtimeRoot, join(runtimeRoot, "runs", runId));
  return {
    runtimeRoot,
    manifestPath: join(runRoot, "manifest.json"),
    recordPath: join(runRoot, "activity", `${event}.json`),
  };
}

export async function prepareActivityLifecycle(options: PrepareLifecycleOptions): Promise<ActivityLifecycleRecord> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const config = await loadConfig(workspaceRoot);
  const { manifestPath, recordPath } = paths(workspaceRoot, options.runId, options.event);
  await ensurePrivateDirectory(join(workspaceRoot, ".runtime"));
  const lockPath = `${recordPath}.lock`;
  return withExclusiveFile(lockPath, async () => {
    try {
      const existing = await readJson<ActivityLifecycleRecord>(recordPath);
      await assertValid("activity-lifecycle-record", existing);
      return existing;
    } catch (error) {
      if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
    }
    const manifest = await readJson<RuntimeManifest>(manifestPath);
    await assertValid("runtime-manifest", manifest);
    if (manifest.run_id !== options.runId) throw new Error("Manifest run ID does not match lifecycle request");
    const available = new Set(options.availableCapabilities ?? []);
    const configured = config.activity.lifecycle?.[options.event] ?? [];
    const actions: ActivityActionResult[] = [];
    const warnings: string[] = [];
    const manualFallbacks: string[] = [];
    let stopped = false;
    for (const action of configured) {
      let status: ActivityActionResult["status"];
      let evidence: string | null = null;
      if (stopped) {
        status = "skipped";
        evidence = "Not attempted after a required lifecycle action failed.";
      } else if (action.policy === "manual") {
        status = "manual";
        manualFallbacks.push(action.description);
      } else if (available.has(action.capability)) {
        status = "pending";
      } else if (action.policy === "required") {
        status = "manual";
        evidence = `Required capability is unavailable: ${action.capability}; manual completion is required.`;
        manualFallbacks.push(action.description);
        stopped = true;
      } else {
        status = "skipped";
        evidence = `Optional capability is unavailable: ${action.capability}.`;
        warnings.push(`${action.id}: ${evidence}`);
        manualFallbacks.push(action.description);
      }
      actions.push({
        ...action,
        status,
        idempotency_key: `${manifest.run_id}:${options.event}:${action.id}`,
        evidence,
        external_reference: null,
      });
    }
    if (config.activity.provider === "none") warnings.push("No activity provider is configured; the semantic event is recorded without an external write.");
    else if (configured.length === 0) warnings.push(`No lifecycle actions are configured for ${options.event}.`);
    const now = (options.now ?? new Date()).toISOString();
    const record: ActivityLifecycleRecord = {
      contract_version: 1,
      work_id: manifest.work_id,
      run_id: manifest.run_id,
      provider: config.activity.provider,
      event: options.event,
      status: overallStatus(actions),
      actions,
      warnings,
      manual_fallbacks: manualFallbacks,
      prepared_at: now,
      updated_at: now,
    };
    await assertValid("activity-lifecycle-record", record);
    syncManifestEvent(manifest, record, recordPath, workspaceRoot);
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(recordPath, record);
    await writeJsonAtomic(manifestPath, manifest);
    return record;
  });
}

export async function recordActivityLifecycleAction(options: RecordLifecycleActionOptions): Promise<ActivityLifecycleRecord> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const { manifestPath, recordPath } = paths(workspaceRoot, options.runId, options.event);
  return withExclusiveFile(`${recordPath}.lock`, async () => {
    const record = await readJson<ActivityLifecycleRecord>(recordPath);
    await assertValid("activity-lifecycle-record", record);
    const action = record.actions.find((item) => item.id === options.actionId);
    if (!action) throw new Error(`Lifecycle event has no action named ${options.actionId}`);
    const evidence = safeEvidence(options.evidence, "Evidence");
    const externalReference = options.externalReference ? safeEvidence(options.externalReference, "External reference") : null;
    if (action.status === "completed" || action.status === "failed") {
      if (action.status === options.status && action.evidence === evidence && action.external_reference === externalReference) return record;
      throw new Error(`Lifecycle action ${action.id} already has a different terminal result`);
    }
    if (action.status === "skipped") throw new Error(`Lifecycle action ${action.id} was skipped and cannot receive an external result`);
    action.status = options.status;
    action.evidence = evidence;
    action.external_reference = externalReference;
    if (options.status === "failed" && action.policy === "required") {
      const index = record.actions.indexOf(action);
      for (const remaining of record.actions.slice(index + 1)) {
        if (remaining.status === "pending" || remaining.status === "manual") {
          remaining.status = "skipped";
          remaining.evidence = "Not attempted after a required lifecycle action failed.";
        }
      }
      if (!record.manual_fallbacks.includes(action.description)) record.manual_fallbacks.push(action.description);
    } else if (options.status === "failed") {
      record.warnings.push(`${action.id}: ${evidence}`);
    }
    record.status = overallStatus(record.actions);
    record.updated_at = (options.now ?? new Date()).toISOString();
    const manifest = await readJson<RuntimeManifest>(manifestPath);
    syncManifestEvent(manifest, record, recordPath, workspaceRoot);
    await assertValid("activity-lifecycle-record", record);
    await assertValid("runtime-manifest", manifest);
    await writeJsonAtomic(recordPath, record);
    await writeJsonAtomic(manifestPath, manifest);
    return record;
  });
}
