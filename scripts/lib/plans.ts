import { createHash, randomUUID } from "node:crypto";
import { lstat, mkdir, readdir, readFile, realpath, rename, rm } from "node:fs/promises";
import { basename, join, resolve } from "node:path";
import { parse as parseYaml, stringify as stringifyYaml } from "yaml";
import { assertInside, writeTextAtomic, writeTextExclusive } from "./io.js";
import type { PlanDraftRequest, PlanIndex, PlanWorkBreakdown, PlanWorkItem, WorkspaceConfig } from "./types.js";
import { readData, validateContract, workspaceSemanticErrors } from "./validation.js";

const documents = [
  "0001-overview.md",
  "0010-requirements.md",
  "0020-solution.md",
  "0040-delivery.md",
  "0050-verification.md",
  "0070-risks.md",
  "0080-work-breakdown.md",
] as const;

const tableHeader = "| Work ID | Title | Parent | Depends on | Area | External reference |";
const tableSeparator = "| --- | --- | --- | --- | --- | --- |";

export interface PlanValidationResult {
  index: PlanIndex | null;
  work_breakdown: PlanWorkBreakdown | null;
  errors: string[];
}

export interface PlanCreationSummary {
  plan_id: string;
  status: "draft";
  plan_version: 1;
  directory: string;
  index: string;
  documents: string[];
  work_ids: string[];
  approval_required: true;
}

export type PlanStateTransition =
  | { kind: "approve"; approved_by: string }
  | { kind: "material-revision"; reason: string }
  | { kind: "non-material-repair" };

function contractMessages(errors: Awaited<ReturnType<typeof validateContract>>): string[] {
  return errors.map((error) => `${error.instancePath || "/"} ${error.message}`);
}

function markdownList(values: string[], empty: string): string {
  return values.length > 0 ? values.map((value) => `- ${value}`).join("\n") : `- ${empty}`;
}

function assertMarkdownCell(value: string, field: string): void {
  if (value.includes("|") || /[\r\n]/.test(value)) throw new Error(`${field} cannot contain a table delimiter or newline`);
}

function allocateWorkItems(request: PlanDraftRequest): PlanWorkItem[] {
  const ids = new Map<string, string>();
  request.work_items.forEach((item, index) => {
    const suffix = index === 0 ? 1 : index * 10;
    ids.set(item.key, `${request.work_prefix}-${String(suffix).padStart(3, "0")}`);
  });
  return request.work_items.map((item) => ({
    work_id: ids.get(item.key)!,
    title: item.title,
    parent: item.parent ? ids.get(item.parent) ?? null : null,
    depends_on: (item.depends_on ?? []).map((key) => ids.get(key) ?? key),
    area: item.area,
    repository: item.repository,
    scope: item.scope,
    test_scope: item.test_scope,
    test_policy: item.test_policy,
    ...(item.test_rationale ? { test_rationale: item.test_rationale } : {}),
    verification_commands: item.verification_commands,
    acceptance_criteria: item.acceptance_criteria,
    external_reference: null,
  }));
}

function cycleErrors(items: Array<{ work_id: string; depends_on: string[] }>): string[] {
  const dependencies = new Map(items.map((item) => [item.work_id, item.depends_on]));
  const errors: string[] = [];
  const visited = new Set<string>();
  const active = new Set<string>();
  const visit = (id: string): void => {
    if (active.has(id)) {
      errors.push(`work dependency cycle includes ${id}`);
      return;
    }
    if (visited.has(id)) return;
    active.add(id);
    for (const dependency of dependencies.get(id) ?? []) visit(dependency);
    active.delete(id);
    visited.add(id);
  };
  for (const id of dependencies.keys()) visit(id);
  return [...new Set(errors)];
}

function parentCycleErrors(items: Array<{ work_id: string; parent: string | null }>): string[] {
  return cycleErrors(items.map((item) => ({ work_id: item.work_id, depends_on: item.parent ? [item.parent] : [] })))
    .map((error) => error.replace("work dependency cycle", "work parent cycle"));
}

export function planDraftSemanticErrors(request: PlanDraftRequest, config?: WorkspaceConfig): string[] {
  const errors: string[] = [];
  const keys = new Set<string>();
  for (const item of request.work_items) {
    if (keys.has(item.key)) errors.push(`duplicate work item key: ${item.key}`);
    keys.add(item.key);
    for (const [field, value] of [["title", item.title], ["area", item.area]] as const) {
      try {
        assertMarkdownCell(value, `work item ${item.key} ${field}`);
      } catch (error) {
        errors.push((error as Error).message);
      }
    }
  }
  for (const item of request.work_items) {
    if (item.parent && !keys.has(item.parent)) errors.push(`work item ${item.key} has unknown parent: ${item.parent}`);
    if (item.parent === item.key) errors.push(`work item ${item.key} cannot be its own parent`);
    for (const dependency of item.depends_on ?? []) {
      if (!keys.has(dependency)) errors.push(`work item ${item.key} has unknown dependency: ${dependency}`);
      if (dependency === item.key) errors.push(`work item ${item.key} cannot depend on itself`);
    }
  }
  const keyedDependencies = request.work_items.map((item) => ({ work_id: item.key, depends_on: item.depends_on ?? [] }));
  errors.push(...cycleErrors(keyedDependencies));
  errors.push(...parentCycleErrors(request.work_items.map((item) => ({ work_id: item.key, parent: item.parent ?? null }))));
  if (config) {
    for (const repository of request.affected_repositories) {
      if (!config.repositories[repository]) errors.push(`affected repository is not registered: ${repository}`);
    }
    for (const item of request.work_items) {
      if (!config.repositories[item.repository]) errors.push(`work item ${item.key} repository is not registered: ${item.repository}`);
      if (!request.affected_repositories.includes(item.repository)) errors.push(`work item ${item.key} repository is not affected: ${item.repository}`);
    }
  }
  return [...new Set(errors)];
}

export function planWorkBreakdownSemanticErrors(breakdown: PlanWorkBreakdown): string[] {
  const errors: string[] = [];
  const ids = new Set<string>();
  for (const item of breakdown.items) {
    if (ids.has(item.work_id)) errors.push(`duplicate work ID: ${item.work_id}`);
    ids.add(item.work_id);
    if (!item.work_id.startsWith(`${breakdown.work_prefix}-`)) {
      errors.push(`work ID does not use ${breakdown.work_prefix} prefix: ${item.work_id}`);
    }
  }
  for (const item of breakdown.items) {
    if (item.parent && !ids.has(item.parent)) errors.push(`${item.work_id} has unknown parent: ${item.parent}`);
    if (item.parent === item.work_id) errors.push(`${item.work_id} cannot be its own parent`);
    for (const dependency of item.depends_on) {
      if (!ids.has(dependency)) errors.push(`${item.work_id} has unknown dependency: ${dependency}`);
      if (dependency === item.work_id) errors.push(`${item.work_id} cannot depend on itself`);
    }
  }
  errors.push(...cycleErrors(breakdown.items));
  errors.push(...parentCycleErrors(breakdown.items));
  return [...new Set(errors)];
}

function materialDigest(files: Map<string, string>, names: readonly string[]): string {
  const hash = createHash("sha256");
  for (const name of names) hash.update(`${name}\0${files.get(name) ?? ""}\0`);
  return `sha256:${hash.digest("hex")}`;
}

export function parsePlanIndex(raw: string): PlanIndex {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---(?:\r?\n|$)/);
  if (!match) throw new Error("Plan README must begin with YAML frontmatter");
  return parseYaml(match[1]!) as PlanIndex;
}

export function parseWorkBreakdown(raw: string, index: PlanIndex): PlanWorkBreakdown {
  const lines = raw.replace(/\r\n/g, "\n").split("\n");
  const header = lines.indexOf(tableHeader);
  if (header === -1 || lines[header + 1] !== tableSeparator) {
    throw new Error("Work breakdown must contain the canonical six-column table and must not add live status columns");
  }
  const summaries: Array<Pick<PlanWorkItem, "work_id" | "title" | "parent" | "depends_on" | "area" | "external_reference">> = [];
  for (const line of lines.slice(header + 2)) {
    if (!line.startsWith("|")) break;
    const cells = line.slice(1, -1).split("|").map((cell) => cell.trim());
    if (cells.length !== 6) throw new Error(`Invalid work breakdown row: ${line}`);
    const [workId, title, parent, dependencies, area, external] = cells as [string, string, string, string, string, string];
    summaries.push({
      work_id: workId,
      title,
      parent: parent === "—" ? null : parent,
      depends_on: dependencies === "—" ? [] : dependencies.split(",").map((value) => value.trim()),
      area,
      external_reference: external === "—" ? null : external,
    });
  }
  const executionMatch = raw.match(/## Execution contracts\r?\n\r?\n```json\r?\n([\s\S]*?)\r?\n```/);
  if (!executionMatch) throw new Error("Work breakdown must contain the canonical execution contracts JSON block");
  const execution = JSON.parse(executionMatch[1]!) as { contract_version: number; items: Array<Omit<PlanWorkItem, "title" | "parent" | "depends_on" | "area" | "external_reference">> };
  if (execution.contract_version !== 1 || !Array.isArray(execution.items)) throw new Error("Invalid work execution contracts block");
  const executionById = new Map(execution.items.map((item) => [item.work_id, item]));
  const items: PlanWorkItem[] = summaries.map((summary) => {
    const details = executionById.get(summary.work_id);
    if (!details) throw new Error(`Missing execution contract for ${summary.work_id}`);
    return { ...summary, ...details };
  });
  for (const workId of executionById.keys()) if (!summaries.some((item) => item.work_id === workId)) throw new Error(`Execution contract references unknown work ID: ${workId}`);
  return { contract_version: 1, plan_id: index.plan_id, work_prefix: index.work_prefix, items };
}

async function regularFile(path: string): Promise<boolean> {
  try {
    const info = await lstat(path);
    return info.isFile() && !info.isSymbolicLink();
  } catch {
    return false;
  }
}

export async function validatePlanDirectory(planDirectory: string, expectedPlanId = basename(planDirectory)): Promise<PlanValidationResult> {
  const directory = resolve(planDirectory);
  const errors: string[] = [];
  let index: PlanIndex | null = null;
  let breakdown: PlanWorkBreakdown | null = null;
  try {
    const info = await lstat(directory);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error("Plan path must be a real directory");
    if (!await regularFile(join(directory, "README.md"))) throw new Error("Plan README must be a real file");
    index = parsePlanIndex(await readFile(join(directory, "README.md"), "utf8"));
    const indexErrors = contractMessages(await validateContract("plan-index", index));
    errors.push(...indexErrors);
    if (indexErrors.length > 0) return { index, work_breakdown: null, errors: [...new Set(errors)] };
    if (index.plan_id !== expectedPlanId) errors.push(`plan_id must match directory name: ${expectedPlanId}`);
    if (Date.parse(index.updated_at) < Date.parse(index.created_at)) errors.push("updated_at cannot be earlier than created_at");
    const sorted = [...index.documents].sort();
    if (JSON.stringify(sorted) !== JSON.stringify(index.documents)) errors.push("numbered plan documents must be listed in ascending order");
    if (!index.documents.includes(index.work_breakdown)) errors.push("work_breakdown must be listed in documents");
    const actualNumbered = (await readdir(directory)).filter((name) => /^[0-9]{4}-.+\.md$/.test(name)).sort();
    for (const document of actualNumbered) {
      if (!index.documents.includes(document)) errors.push(`numbered plan document is not listed in the index: ${document}`);
    }
    for (const document of index.documents) {
      if (!await regularFile(join(directory, document))) errors.push(`plan document is missing or unsafe: ${document}`);
    }
    if (errors.length === 0) {
      const material = new Map<string, string>();
      for (const document of index.documents) material.set(document, await readFile(join(directory, document), "utf8"));
      const digest = materialDigest(material, index.documents);
      if (index.material_digest !== digest) errors.push("material_digest does not match the numbered plan documents");
      if (index.status === "approved" && index.approved_digest !== digest) errors.push("approved_digest does not match the approved plan material");
      breakdown = parseWorkBreakdown(material.get(index.work_breakdown)!, index);
      errors.push(...contractMessages(await validateContract("plan-work-breakdown", breakdown)));
      errors.push(...planWorkBreakdownSemanticErrors(breakdown));
    }
  } catch (error) {
    errors.push((error as Error).message);
  }
  return { index, work_breakdown: breakdown, errors: [...new Set(errors)] };
}

async function actualMaterialDigest(directory: string, index: PlanIndex): Promise<string> {
  const material = new Map<string, string>();
  for (const document of index.documents) material.set(document, await readFile(join(directory, document), "utf8"));
  return materialDigest(material, index.documents);
}

export async function setPlanState(planDirectory: string, transition: PlanStateTransition, now = new Date()): Promise<PlanIndex> {
  const directory = resolve(planDirectory);
  const validation = await validatePlanDirectory(directory);
  const allowedStaleDigestErrors = new Set([
    "material_digest does not match the numbered plan documents",
    "approved_digest does not match the approved plan material",
  ]);
  const blocking = validation.errors.filter((error) => !allowedStaleDigestErrors.has(error));
  if (!validation.index || blocking.length > 0) throw new Error(`Plan state transition validation failed:\n- ${blocking.join("\n- ")}`);
  const index = validation.index;
  const digest = await actualMaterialDigest(directory, index);
  if (transition.kind === "approve") {
    if (index.status !== "draft") throw new Error("Only a draft plan can be approved");
    if (!transition.approved_by.trim()) throw new Error("Approval requires a non-empty approver");
    index.status = "approved";
    index.approved_at = now.toISOString();
    index.approved_by = transition.approved_by.trim();
    index.material_digest = digest;
    index.approved_digest = digest;
  } else if (transition.kind === "material-revision") {
    if (index.status !== "approved") throw new Error("Material revision transition requires an approved plan");
    if (!transition.reason.trim()) throw new Error("Material revision requires a reason");
    index.status = "draft";
    index.plan_version += 1;
    index.approved_at = null;
    index.approved_by = null;
    index.approved_digest = null;
    index.material_digest = digest;
    index.revision_reason = transition.reason.trim();
  } else {
    if (index.status !== "approved") throw new Error("Non-material repair transition requires an approved plan");
    index.material_digest = digest;
    index.approved_digest = digest;
  }
  index.updated_at = now.toISOString();
  const readmePath = join(directory, "README.md");
  const raw = await readFile(readmePath, "utf8");
  const match = raw.match(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/);
  if (!match) throw new Error("Plan README must begin with YAML frontmatter");
  await writeTextAtomic(readmePath, raw.replace(match[0], `---\n${stringifyYaml(index).trimEnd()}\n---\n`));
  const after = await validatePlanDirectory(directory);
  if (after.errors.length > 0) throw new Error(`Plan state transition produced invalid metadata:\n- ${after.errors.join("\n- ")}`);
  return index;
}

function renderDocument(title: string, sections: Array<[string, string]>): string {
  return `# ${title}\n\n${sections.map(([heading, body]) => `## ${heading}\n\n${body}`).join("\n\n")}\n`;
}

function renderPlan(request: PlanDraftRequest, createdAt: string): { index: PlanIndex; breakdown: PlanWorkBreakdown; files: Map<string, string> } {
  const workItems = allocateWorkItems(request);
  const breakdown: PlanWorkBreakdown = { contract_version: 1, plan_id: request.plan_id, work_prefix: request.work_prefix, items: workItems };
  const files = new Map<string, string>();
  files.set("0001-overview.md", renderDocument("Overview", [
    ["Summary", request.summary],
    ["Source", `${request.source.kind}: ${request.source.reference}`],
    ["Affected repositories", markdownList(request.affected_repositories, "None identified.")],
    ["Assumptions", markdownList(request.assumptions, "None recorded.")],
    ["Open questions", markdownList(request.open_questions, "None recorded.")],
  ]));
  files.set("0010-requirements.md", renderDocument("Requirements", [["Requirements and acceptance criteria", markdownList(request.requirements, "None recorded.")]]));
  files.set("0020-solution.md", renderDocument("Solution", [["Proposed solution", markdownList(request.solution, "None recorded.")]]));
  files.set("0040-delivery.md", renderDocument("Delivery", [["Delivery order", markdownList(request.delivery, "None recorded.")]]));
  files.set("0050-verification.md", renderDocument("Verification", [["Verification strategy", markdownList(request.verification, "None recorded.")]]));
  files.set("0070-risks.md", renderDocument("Risks", [["Risks and mitigations", markdownList(request.risks, "None recorded.")]]));
  const rows = workItems.map((item) => `| ${item.work_id} | ${item.title} | ${item.parent ?? "—"} | ${item.depends_on.join(", ") || "—"} | ${item.area} | — |`).join("\n");
  const execution = {
    contract_version: 1,
    items: workItems.map(({ work_id, repository, scope, test_scope, test_policy, test_rationale, verification_commands, acceptance_criteria }) => ({
      work_id, repository, scope, test_scope, test_policy, ...(test_rationale ? { test_rationale } : {}), verification_commands, acceptance_criteria,
    })),
  };
  files.set("0080-work-breakdown.md", `# Work breakdown\n\n${tableHeader}\n${tableSeparator}\n${rows}\n\n## Execution contracts\n\n\`\`\`json\n${JSON.stringify(execution, null, 2)}\n\`\`\`\n\nLive task status does not belong in this plan. Add confirmed external references only after an explicit publication action.\n`);
  const index: PlanIndex = {
    contract_version: 1,
    plan_id: request.plan_id,
    title: request.title,
    status: "draft",
    plan_version: 1,
    approved_at: null,
    approved_by: null,
    revision_reason: "Initial draft",
    source: request.source,
    work_prefix: request.work_prefix,
    documents: [...documents],
    work_breakdown: "0080-work-breakdown.md",
    material_digest: materialDigest(files, documents),
    approved_digest: null,
    created_at: createdAt,
    updated_at: createdAt,
  };
  const links = documents.map((document) => `- [${document.replace(/^[0-9]{4}-|\.md$/g, "").replaceAll("-", " ")}](./${document})`).join("\n");
  files.set("README.md", `---\n${stringifyYaml(index).trimEnd()}\n---\n\n# ${request.title}\n\n${request.summary}\n\n## Plan documents\n\n${links}\n\n## Approval gate\n\nHuman approval must explicitly cover scope, solution, delivery order, risks, and acceptance criteria before the metadata status changes to \`approved\`. The machine-readable frontmatter status is authoritative; approval updates metadata without rewriting this prose.\n`);
  return { index, breakdown, files };
}

export async function createPlanDraft(workspaceRootInput: string, request: PlanDraftRequest, now = new Date()): Promise<PlanCreationSummary> {
  const workspaceRoot = resolve(workspaceRootInput);
  const contractErrors = contractMessages(await validateContract("plan-draft-request", request));
  const config = await readData(join(workspaceRoot, "workspace.yaml")) as WorkspaceConfig;
  const workspaceErrors = contractMessages(await validateContract("workspace", config));
  const errors = [...contractErrors, ...workspaceErrors];
  if (contractErrors.length === 0 && workspaceErrors.length === 0) {
    errors.push(...workspaceSemanticErrors(config), ...planDraftSemanticErrors(request, config));
  }
  if (errors.length > 0) throw new Error(`Invalid plan draft request:\n- ${errors.join("\n- ")}`);
  const realWorkspace = await realpath(workspaceRoot);
  const contextRoot = assertInside(workspaceRoot, join(workspaceRoot, "context"));
  try {
    const info = await lstat(contextRoot);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`Context root must be a real directory: ${contextRoot}`);
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
    await mkdir(contextRoot, { mode: 0o755 });
  }
  if (await realpath(contextRoot) !== join(realWorkspace, "context")) throw new Error(`Context root must not traverse symbolic links: ${contextRoot}`);
  const plansRoot = assertInside(contextRoot, join(contextRoot, "plans"));
  try {
    const info = await lstat(plansRoot);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`Plan root must be a real directory: ${plansRoot}`);
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
    await mkdir(plansRoot, { mode: 0o755 });
  }
  const realPlansRoot = await realpath(plansRoot);
  assertInside(realWorkspace, realPlansRoot);
  if (realPlansRoot !== join(realWorkspace, "context", "plans")) {
    throw new Error(`Plan root must not traverse symbolic links: ${plansRoot}`);
  }
  const destination = assertInside(realPlansRoot, join(realPlansRoot, request.plan_id));
  try {
    await lstat(destination);
    throw new Error(`Plan already exists; refusing to overwrite: ${destination}`);
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
  }
  const temporary = join(realPlansRoot, `.${request.plan_id}.${randomUUID()}.tmp`);
  const rendered = renderPlan(request, now.toISOString());
  try {
    await mkdir(temporary, { mode: 0o755 });
    for (const [name, contents] of rendered.files) await writeTextExclusive(join(temporary, name), contents);
    const validation = await validatePlanDirectory(temporary, request.plan_id);
    if (validation.errors.length > 0) throw new Error(`Generated plan failed validation:\n- ${validation.errors.join("\n- ")}`);
    await rename(temporary, destination);
  } catch (error) {
    await rm(temporary, { recursive: true, force: true });
    throw error;
  }
  return {
    plan_id: request.plan_id,
    status: "draft",
    plan_version: 1,
    directory: destination,
    index: join(destination, "README.md"),
    documents: [...documents],
    work_ids: rendered.breakdown.items.map((item) => item.work_id),
    approval_required: true,
  };
}
