import { readFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { assertInside, ensurePrivateDirectory, withExclusiveFile, writeJsonAtomic, writeTextAtomic } from "./io.js";
import { setPlanState, validatePlanDirectory } from "./plans.js";
import type { PlanPublicationDiscovery, PlanPublicationItem, PlanPublicationRecord, PlanWorkItem, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";

export interface PreparePublicationOptions { workspaceRoot: string; planId: string; discovery: PlanPublicationDiscovery; now?: Date }
export interface RecordPublicationOptions {
  workspaceRoot: string; planId: string; workId: string; status: "created" | "failed";
  evidence: string; externalReference?: string; now?: Date;
}

async function assertValid(name: "workspace" | "plan-publication-discovery" | "plan-publication-record", value: unknown): Promise<void> {
  const errors = await validateContract(name, value);
  if (errors.length) throw new Error(`Invalid ${name}: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
}

function safeLine(value: string, field: string): string {
  const clean = value.trim();
  if (!clean || /[|\r\n]/.test(clean)) throw new Error(`${field} must be a non-empty single line without table delimiters`);
  if (/https?:\/\/[^\s/@:]+:[^\s/@]+@|-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/i.test(clean)) throw new Error(`${field} appears to contain a credential`);
  return clean;
}

function ordered(items: PlanWorkItem[]): PlanWorkItem[] {
  const remaining = [...items];
  const emitted = new Set<string>();
  const result: PlanWorkItem[] = [];
  while (remaining.length) {
    const index = remaining.findIndex((item) => [...item.depends_on, ...(item.parent ? [item.parent] : [])].every((id) => emitted.has(id)));
    if (index < 0) throw new Error("Plan publication order cannot resolve dependencies and parents");
    const [item] = remaining.splice(index, 1);
    result.push(item!);
    emitted.add(item!.work_id);
  }
  return result;
}

function status(items: PlanPublicationItem[]): PlanPublicationRecord["status"] {
  const done = items.filter((item) => item.status === "created" || item.status === "existing").length;
  const failed = items.filter((item) => item.status === "failed").length;
  if (done === items.length) return "completed";
  if (failed && done) return "partial";
  if (failed) return "failed";
  return done ? "in-progress" : "proposed";
}

function recordPath(workspaceRoot: string, planId: string): string {
  return assertInside(workspaceRoot, join(workspaceRoot, ".runtime", "publications", `${planId}.json`));
}

export async function preparePlanPublication(options: PreparePublicationOptions): Promise<PlanPublicationRecord> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  await assertValid("workspace", config);
  const semantic = workspaceSemanticErrors(config);
  if (semantic.length) throw new Error(`Invalid workspace: ${semantic.join("; ")}`);
  await assertValid("plan-publication-discovery", options.discovery);
  if (config.activity.provider === "none" || config.activity.provider !== options.discovery.provider) throw new Error("Publication discovery provider must match the configured non-none activity provider");
  if (![...config.activity.required_capabilities, ...config.activity.optional_capabilities].includes("create-tasks")) throw new Error("Configured activity provider does not declare create-tasks capability");
  const planDirectory = assertInside(workspaceRoot, join(workspaceRoot, "context", "plans", options.planId));
  const plan = await validatePlanDirectory(planDirectory);
  if (plan.errors.length || !plan.index || !plan.work_breakdown) throw new Error(`Plan is invalid: ${plan.errors.join("; ")}`);
  if (plan.index.status !== "approved" || !plan.index.approved_digest) throw new Error("Only an approved plan can be published");
  const index = plan.index;
  const breakdown = plan.work_breakdown;
  const discovered = new Map<string, { reference: string; evidence: string }>();
  for (const mapping of options.discovery.mappings) {
    if (discovered.has(mapping.work_id)) throw new Error(`Duplicate discovered mapping: ${mapping.work_id}`);
    discovered.set(mapping.work_id, { reference: safeLine(mapping.external_reference, "External reference"), evidence: safeLine(mapping.evidence, "Evidence") });
  }
  const path = recordPath(workspaceRoot, options.planId);
  await ensurePrivateDirectory(join(workspaceRoot, ".runtime", "publications"));
  return withExclusiveFile(`${path}.lock`, async () => {
    try {
      const existing = JSON.parse(await readFile(path, "utf8")) as PlanPublicationRecord;
      await assertValid("plan-publication-record", existing);
      if (existing.plan_version !== index.plan_version || existing.provider !== options.discovery.provider || existing.destination !== options.discovery.destination) throw new Error("Existing publication record conflicts with this request");
      return existing;
    } catch (error) {
      if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
    }
    const items = ordered(breakdown.items).map<PlanPublicationItem>((item) => {
      const known = item.external_reference ? { reference: item.external_reference, evidence: "Confirmed mapping already stored in the approved plan." } : discovered.get(item.work_id);
      if (item.external_reference && discovered.get(item.work_id)?.reference !== undefined && discovered.get(item.work_id)!.reference !== item.external_reference) throw new Error(`Conflicting external mapping for ${item.work_id}`);
      return { ...item, action: known ? "skip-existing" : "create", status: known ? "existing" : "proposed", external_reference: known?.reference ?? null, evidence: known?.evidence ?? null, idempotency_key: `${index.plan_id}:v${index.plan_version}:${item.work_id}` };
    });
    for (const workId of discovered.keys()) if (!items.some((item) => item.work_id === workId)) throw new Error(`Discovered mapping references unknown work ID: ${workId}`);
    const now = (options.now ?? new Date()).toISOString();
    const record: PlanPublicationRecord = { contract_version: 1, plan_id: index.plan_id, plan_version: index.plan_version, approved_digest: index.approved_digest!, provider: options.discovery.provider, destination: safeLine(options.discovery.destination, "Destination"), status: status(items), items, warnings: [], prepared_at: now, updated_at: now };
    await assertValid("plan-publication-record", record);
    await writeJsonAtomic(path, record);
    return record;
  });
}

async function writeMapping(planDirectory: string, breakdownName: string, workId: string, reference: string, now: Date): Promise<string> {
  const path = join(planDirectory, breakdownName);
  const raw = await readFile(path, "utf8");
  let found = false;
  const updated = raw.split("\n").map((line) => {
    if (!line.startsWith(`| ${workId} |`)) return line;
    const cells = line.slice(1, -1).split("|").map((cell) => cell.trim());
    if (cells.length !== 6) throw new Error(`Invalid work-breakdown row for ${workId}`);
    if (cells[5] !== "—" && cells[5] !== reference) throw new Error(`Plan already maps ${workId} to a different external reference`);
    cells[5] = reference;
    found = true;
    return `| ${cells.join(" | ")} |`;
  }).join("\n");
  if (!found) throw new Error(`Plan has no work item ${workId}`);
  await writeTextAtomic(path, updated);
  return (await setPlanState(planDirectory, { kind: "non-material-repair" }, now)).approved_digest!;
}

export async function recordPlanPublication(options: RecordPublicationOptions): Promise<PlanPublicationRecord> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const path = recordPath(workspaceRoot, options.planId);
  return withExclusiveFile(`${path}.lock`, async () => {
    const record = JSON.parse(await readFile(path, "utf8")) as PlanPublicationRecord;
    await assertValid("plan-publication-record", record);
    const item = record.items.find((candidate) => candidate.work_id === options.workId);
    if (!item) throw new Error(`Publication has no work item ${options.workId}`);
    const evidence = safeLine(options.evidence, "Evidence");
    const reference = options.externalReference ? safeLine(options.externalReference, "External reference") : null;
    if (options.status === "created" && !reference) throw new Error("Created publication result requires a confirmed external reference");
    if (item.status === "created" || item.status === "existing") {
      if (item.external_reference === reference && item.evidence === evidence) return record;
      throw new Error(`${item.work_id} already has a different confirmed mapping`);
    }
    if (options.status === "failed") {
      item.status = "failed"; item.evidence = evidence; item.external_reference = null;
    } else {
      const planDirectory = assertInside(workspaceRoot, join(workspaceRoot, "context", "plans", record.plan_id));
      const validation = await validatePlanDirectory(planDirectory);
      if (validation.errors.length || !validation.index || validation.index.status !== "approved" || validation.index.plan_version !== record.plan_version) throw new Error("Approved plan changed during publication");
      record.approved_digest = await writeMapping(planDirectory, validation.index.work_breakdown, item.work_id, reference!, options.now ?? new Date());
      item.status = "created"; item.external_reference = reference; item.evidence = evidence;
    }
    record.status = status(record.items);
    record.updated_at = (options.now ?? new Date()).toISOString();
    await assertValid("plan-publication-record", record);
    await writeJsonAtomic(path, record);
    return record;
  });
}
