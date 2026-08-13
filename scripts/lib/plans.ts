import { createHash, randomUUID } from "node:crypto";
import { cp, lstat, mkdir, readdir, readFile, realpath, rename, rm } from "node:fs/promises";
import { basename, dirname, join, relative, resolve } from "node:path";
import { Ajv2020 } from "ajv/dist/2020.js";
import { parse as parseYaml, stringify as stringifyYaml } from "yaml";
import { assertInside, writeTextAtomic, writeTextExclusive } from "./io.js";
import type { PlanDraftRequest, PlanGenerationDefinition, PlanGenerationRequest, PlanIndex, PlanLifecycleStatus, PlanTaskContract, PlanTrack, PlanWorkBreakdown, PlanWorkItem, ProductKnowledgePlanDeclaration, WorkspaceConfig } from "./types.js";
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

const tableHeader = "| Work ID | Title | Parent | Depends on | Repository | Area | External reference |";
const tableSeparator = "| --- | --- | --- | --- | --- | --- | --- |";
const legacyTableHeader = "| Work ID | Title | Parent | Depends on | Area | External reference |";
const legacyTableSeparator = "| --- | --- | --- | --- | --- | --- |";

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
  | { kind: "non-material-repair" }
  | { kind: "lifecycle"; status: Exclude<PlanLifecycleStatus, "draft" | "approved" | "archived">; reason: string; actor?: string; evidence?: string }
  | { kind: "archive"; actor: string; evidence: string };

function contractMessages(errors: Awaited<ReturnType<typeof validateContract>>): string[] {
  return errors.map((error) => `${error.instancePath || "/"} ${error.message}`);
}

function markdownList(values: string[], empty: string): string {
  return values.length > 0 ? values.map((value) => `- ${value}`).join("\n") : `- ${empty}`;
}

function productImpactBody(declaration: ProductKnowledgePlanDeclaration): string {
  const references = markdownList(declaration.references, "None referenced.");
  const proposed = declaration.proposed_change ?? "No product behavior change is proposed.";
  return `- Impact: ${declaration.impact}\n\nReferenced Product Knowledge:\n\n${references}\n\nProposed change:\n\n${proposed}`;
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
    description: item.description?.trim() || item.title.trim(),
    ...(item.subtasks ? { subtasks: item.subtasks.map((key) => ids.get(key) ?? key) } : {}),
    ...(item.connections ? {
      connections: normalizeConnectionList(item.connections).map((connection) => ({
        ...connection,
        target: ids.get(connection.target) ?? connection.target,
      })),
    } : {}),
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
  const productKnowledge = request.product_knowledge;
  if (productKnowledge) {
    const requiresChange = ["behavior-change", "new-workflow", "retired-workflow"];
    if (requiresChange.includes(productKnowledge.impact) && !productKnowledge.proposed_change?.trim()) {
      errors.push(`product knowledge impact '${productKnowledge.impact}' requires a proposed_change summary`);
    }
    if (productKnowledge.impact === "none" && productKnowledge.proposed_change) {
      errors.push("product knowledge impact 'none' must not include a proposed_change");
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

export function planWorkBreakdownSemanticErrors(breakdown: PlanWorkBreakdown, config?: WorkspaceConfig): string[] {
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
    for (const subtask of item.subtasks ?? []) {
      if (!ids.has(subtask)) errors.push(`${item.work_id} has unknown subtask: ${subtask}`);
      if (subtask === item.work_id) errors.push(`${item.work_id} cannot list itself as a subtask`);
      const child = breakdown.items.find((candidate) => candidate.work_id === subtask);
      if (child && child.parent !== item.work_id) errors.push(`${item.work_id} lists ${subtask} as a subtask but its parent is ${child.parent ?? "none"}`);
    }
    if (config && !config.repositories[item.repository]) errors.push(`${item.work_id} repository is not registered: ${item.repository}`);
  }
  errors.push(...cycleErrors(breakdown.items.map((item) => ({
    work_id: item.work_id,
    depends_on: [...item.depends_on, ...(item.connections ?? []).filter((connection) => connection.type === "depends-on" && ids.has(connection.target)).map((connection) => connection.target)],
  }))));
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

export function parseWorkBreakdown(raw: string, index: PlanIndex, config?: WorkspaceConfig): PlanWorkBreakdown {
  const lines = raw.replace(/\r\n/g, "\n").split("\n");
  const currentHeader = lines.indexOf(tableHeader);
  const legacyHeader = lines.indexOf(legacyTableHeader);
  const legacy = currentHeader === -1 && legacyHeader !== -1;
  const header = currentHeader === -1 ? legacyHeader : currentHeader;
  if (header === -1 || lines[header + 1] !== (legacy ? legacyTableSeparator : tableSeparator)) {
    throw new Error("Work breakdown must contain the canonical seven-column table and must not add live status columns");
  }
  const summaries: Array<Pick<PlanWorkItem, "work_id" | "title" | "parent" | "depends_on" | "area" | "external_reference"> & { repository: string | null }> = [];
  for (const line of lines.slice(header + 2)) {
    if (!line.startsWith("|")) break;
    const cells = line.slice(1, -1).split("|").map((cell) => cell.trim());
    if (cells.length !== (legacy ? 6 : 7)) throw new Error(`Invalid work breakdown row: ${line}`);
    const [workId, title, parent, dependencies] = cells;
    const repository = legacy ? null : cells[4]!;
    const area = cells[legacy ? 4 : 5]!;
    const external = cells[legacy ? 5 : 6]!;
    summaries.push({
      work_id: workId!,
      title: title!,
      parent: parent === "—" ? null : parent!,
      depends_on: dependencies === "—" ? [] : dependencies!.split(",").map((value) => value.trim()),
      repository,
      area,
      external_reference: external === "—" ? null : external,
    });
  }
  const executionMatch = raw.match(/## Execution contracts\r?\n\r?\n```json\r?\n([\s\S]*?)\r?\n```/);
  if (!executionMatch) throw new Error("Work breakdown must contain the canonical execution contracts JSON block");
  const execution = JSON.parse(executionMatch[1]!) as { contract_version: number; items: Array<Partial<Omit<PlanWorkItem, "title" | "parent" | "depends_on" | "area" | "external_reference">> & { work_id: string }> };
  if (![1, 2].includes(execution.contract_version) || !Array.isArray(execution.items)) throw new Error("Invalid work execution contracts block");
  if (!legacy && execution.contract_version !== 2) throw new Error("The canonical seven-column work breakdown requires execution contract version 2");
  if (legacy && execution.contract_version !== 1) throw new Error("The legacy six-column work breakdown requires execution contract version 1");
  const executionById = new Map(execution.items.map((item) => [item.work_id, item]));
  const items: PlanWorkItem[] = summaries.map((summary) => {
    const details = executionById.get(summary.work_id);
    if (!details) throw new Error(`Missing execution contract for ${summary.work_id}`);
    let repository = details.repository ?? summary.repository;
    if (summary.repository && details.repository && summary.repository !== details.repository) {
      throw new Error(`Repository mismatch for ${summary.work_id}: table has ${summary.repository}, execution contract has ${details.repository}`);
    }
    if (!repository && legacy && execution.contract_version === 1) {
      if (!config) throw new Error(`Legacy work item ${summary.work_id} has no repository; validate it inside a configured workspace or migrate the plan`);
      repository = config.repositories[summary.area] ? summary.area : null;
      if (!repository) throw new Error(`Legacy work item ${summary.work_id} has no repository and area '${summary.area}' is not an exact registered repository key; add an explicit repository through a material plan revision`);
    }
    if (!repository) throw new Error(`Work item ${summary.work_id} has no explicit repository`);
    return { ...summary, ...details, repository } as PlanWorkItem;
  });
  for (const workId of executionById.keys()) if (!summaries.some((item) => item.work_id === workId)) throw new Error(`Execution contract references unknown work ID: ${workId}`);
  return { contract_version: 2, plan_id: index.plan_id, work_prefix: index.work_prefix, items };
}

async function regularFile(path: string): Promise<boolean> {
  try {
    const info = await lstat(path);
    return info.isFile() && !info.isSymbolicLink();
  } catch {
    return false;
  }
}

async function workspaceRootForPlan(directory: string): Promise<string> {
  let current = resolve(directory);
  while (true) {
    if (await regularFile(join(current, "workspace.yaml"))) return current;
    const parent = dirname(current);
    if (parent === current) break;
    current = parent;
  }
  throw new Error(`Unable to locate workspace.yaml for plan: ${directory}`);
}

function matchingPlanIndexes(indexes: PlanIndex[], reference: string): PlanIndex[] {
  const normalized = planDependencyKey(reference);
  return indexes.filter((candidate) => candidate.plan_id === reference || planDependencyKey(candidate.plan_reference ?? "") === normalized);
}

function planReferenceErrors(current: PlanIndex, indexes: PlanIndex[]): string[] {
  const errors: string[] = [];
  const all = [...indexes.filter((candidate) => candidate.plan_id !== current.plan_id), current];
  const dependencies = new Map<string, string[]>();
  for (const candidate of all) {
    const edges: string[] = [];
    for (const reference of candidate.depends_on_plans ?? []) {
      const matches = matchingPlanIndexes(all, reference);
      if (matches.length === 0) {
        if (candidate.plan_id === current.plan_id) errors.push(`${candidate.plan_id} references unknown plan: ${reference}`);
        continue;
      }
      if (matches.length > 1) {
        if (candidate.plan_id === current.plan_id) errors.push(`Plan reference is ambiguous: ${reference}`);
        continue;
      }
      if (matches[0]!.plan_id === candidate.plan_id) errors.push(`plan ${candidate.plan_id} cannot depend on itself`);
      edges.push(matches[0]!.plan_id);
    }
    for (const connection of candidate.connections ?? []) {
      const matches = matchingPlanIndexes(all, connection.target);
      if (matches.length === 0) {
        if (candidate.plan_id === current.plan_id) errors.push(`${candidate.plan_id} has an unresolved connection target: ${connection.target}`);
        continue;
      }
      if (matches.length > 1) {
        if (candidate.plan_id === current.plan_id) errors.push(`Plan connection target is ambiguous: ${connection.target}`);
        continue;
      }
      if (connection.type === "depends-on") {
        if (matches[0]!.plan_id === candidate.plan_id) errors.push(`plan ${candidate.plan_id} cannot depend on itself`);
        edges.push(matches[0]!.plan_id);
      }
    }
    dependencies.set(candidate.plan_id, edges);
  }
  errors.push(...cycleErrors([...dependencies.entries()].map(([work_id, depends_on]) => ({ work_id, depends_on }))).map((error) => error.replace("work dependency cycle", "plan dependency cycle")));
  return [...new Set(errors)];
}

async function validateRootTaskFiles(
  workspaceRoot: string,
  directory: string,
  index: PlanIndex,
  breakdown: PlanWorkBreakdown,
  material: Map<string, string>,
  knownPlans: PlanIndex[],
  resolveExternalPlanReferences: boolean,
  errors: string[],
): Promise<void> {
  const taskDirectory = join(directory, "tasks");
  if (!await regularDirectory(taskDirectory)) {
    errors.push("tasks/ must be a real directory");
    return;
  }
  const entries = await readdir(taskDirectory, { withFileTypes: true });
  for (const entry of entries) {
    if (entry.name === "README.md") continue;
    if (entry.isDirectory() || entry.isSymbolicLink()) errors.push(`tasks/ contains an unsafe or unexpected entry: ${entry.name}`);
  }
  const expected = new Map(breakdown.items.map((item) => [item.work_id, item]));
  for (const entry of entries) {
    if (!entry.name.endsWith(".md") || entry.name === "README.md") continue;
    const taskId = entry.name.slice(0, -3);
    if (!expected.has(taskId)) errors.push(`task file does not resolve to a task contract: ${entry.name}`);
  }
  for (const item of breakdown.items) {
    const name = `tasks/${item.work_id}.md`;
    const raw = material.get(name);
    if (raw === undefined) {
      errors.push(`task file is missing: ${name}`);
      continue;
    }
    const parsed = parseTaskFrontmatter(raw);
    errors.push(...parsed.errors.map((error) => `${name}: ${error}`));
    const frontmatter = parsed.value as Partial<PlanTaskContract> | null;
    if (!frontmatter || typeof frontmatter !== "object") continue;
    errors.push(...(await validatePlanTaskContract(workspaceRoot, frontmatter)).map((error) => `${name}: ${error}`));
    if (frontmatter.task_id !== item.work_id) errors.push(`${name}: task_id does not match its stable filename`);
    if (frontmatter.plan_id !== index.plan_id) errors.push(`${name}: plan_id does not match the containing plan`);
    if (frontmatter.repository !== item.repository) errors.push(`${name}: repository does not match the task graph`);
    if (frontmatter.parent_task !== item.parent) errors.push(`${name}: parent_task does not match the task graph`);
    if (!sameStringList(frontmatter.depends_on, item.depends_on)) errors.push(`${name}: depends_on does not match the task graph`);
    if (!sameStringList(frontmatter.subtasks, item.subtasks ?? [])) errors.push(`${name}: subtasks does not match the task graph`);
    if (JSON.stringify(frontmatter.connections ?? []) !== JSON.stringify(item.connections ?? [])) errors.push(`${name}: connections do not match the task graph`);
    const markdown = markdownSections(parsed.body);
    errors.push(...markdown.errors.map((error) => `${name}: ${error}`));
    if (markdown.title !== item.title) errors.push(`${name}: Markdown title does not match the task graph`);
    errors.push(...authoredLiveStatusErrors(raw).map((error) => `${name}: ${error}`));
    for (const connection of item.connections ?? []) {
      if (expected.has(connection.target)) continue;
      if (!resolveExternalPlanReferences) continue;
      const matches = matchingPlanIndexes(knownPlans, connection.target);
      if (matches.length === 0) errors.push(`${name}: connection target does not resolve exactly: ${connection.target}`);
      else if (matches.length > 1) errors.push(`${name}: connection target is ambiguous: ${connection.target}`);
    }
  }
}

async function validateRootPlanDirectory(
  directory: string,
  index: PlanIndex,
  config: WorkspaceConfig,
  errors: string[],
  resolveExternalPlanReferences = true,
  suppliedPlanIndexes?: PlanIndex[],
): Promise<PlanWorkBreakdown | null> {
  const expectedFolder = index.plan_reference?.replace(/^.*\//, "").replace(/@v[0-9]+$/, "");
  const archivedFolder = expectedFolder?.replace(/^\d{3,}-/, "");
  if (expectedFolder && basename(directory) !== expectedFolder && basename(directory) !== archivedFolder) errors.push(`plan folder does not match plan reference: ${basename(directory)}`);
  if (!index.repository_collection || !index.track || index.plan_number === undefined || !index.task_index) {
    errors.push("root plan metadata is incomplete");
    return null;
  }
  const primaryRepository = index.repository_collection.replace(/-plans$/, "");
  if (!config.repositories[primaryRepository] || !index.affected_repositories?.includes(primaryRepository)) {
    errors.push("repository_collection must be derived from an exact affected repository key");
  }
  const expectedDocuments = [...ROOT_PLAN_DOCUMENTS];
  if (JSON.stringify([...index.documents].sort()) !== JSON.stringify([...expectedDocuments].sort())) errors.push("root plan documents must use the fixed unnumbered inventory");
  const material = new Map<string, string>();
  for (const document of expectedDocuments) {
    if (!await regularFile(join(directory, document))) errors.push(`plan document is missing or unsafe: ${document}`);
    else {
      const raw = await readFile(join(directory, document), "utf8");
      material.set(document, raw);
      errors.push(...authoredLiveStatusErrors(raw).map((error) => `${document}: ${error}`));
    }
  }
  const taskIndexPath = join(directory, "tasks", "README.md");
  if (index.task_index !== "tasks/README.md" || index.work_breakdown !== "tasks/README.md") errors.push("root plan task index must be tasks/README.md");
  if (!await regularFile(taskIndexPath)) errors.push("tasks/README.md is missing or unsafe");
  else {
    const raw = await readFile(taskIndexPath, "utf8");
    material.set("tasks/README.md", raw);
    errors.push(...authoredLiveStatusErrors(raw).map((error) => `tasks/README.md: ${error}`));
  }
  if (await regularDirectory(join(directory, "tasks"))) {
    for (const entry of await readdir(join(directory, "tasks"))) {
      if (entry === "README.md" || !entry.endsWith(".md")) continue;
      if (!await regularFile(join(directory, "tasks", entry))) errors.push(`task file is missing or unsafe: ${entry}`);
      else material.set(`tasks/${entry}`, await readFile(join(directory, "tasks", entry), "utf8"));
    }
  }
  const digest = rootMaterialDigest(material);
  if (index.material_digest !== digest) errors.push("material_digest does not match root plan documents and task files");
  if (index.status === "approved" && index.approved_digest !== digest) errors.push("approved_digest does not match the approved root plan material");
  const rawTasks = material.get("tasks/README.md");
  const executionMatch = rawTasks?.match(/```json\r?\n([\s\S]*?)\r?\n```/);
  if (!executionMatch) {
    errors.push("tasks/README.md must contain the complete task graph JSON");
    return null;
  }
  try {
    const breakdown = JSON.parse(executionMatch[1]!) as PlanWorkBreakdown;
    const contractErrors = contractMessages(await validateContract("plan-work-breakdown", breakdown));
    errors.push(...contractErrors, ...planWorkBreakdownSemanticErrors(breakdown, config));
    if (breakdown.plan_id !== index.plan_id) errors.push("task graph plan_id does not match the plan README");
    if (breakdown.work_prefix !== index.work_prefix) errors.push("task graph work_prefix does not match the plan README");
    const knownPlans = suppliedPlanIndexes ?? (resolveExternalPlanReferences ? [
      ...(await readRootIndexesIfPresent(join(await workspaceRootForPlan(directory), "plans"))),
      ...(await readRootIndexesIfPresent(join(await workspaceRootForPlan(directory), "archived", "plans"))),
    ] : [index]);
    if (resolveExternalPlanReferences) errors.push(...planReferenceErrors(index, knownPlans));
    await validateRootTaskFiles(await workspaceRootForPlan(directory), directory, index, breakdown, material, knownPlans, resolveExternalPlanReferences, errors);
    return breakdown;
  } catch (error) {
    errors.push(`Invalid root task graph: ${(error as Error).message}`);
    return null;
  }
}

export async function validatePlanDirectory(planDirectory: string, expectedPlanId = basename(planDirectory)): Promise<PlanValidationResult> {
  const directory = resolve(planDirectory);
  const errors: string[] = [];
  let index: PlanIndex | null = null;
  let breakdown: PlanWorkBreakdown | null = null;
  try {
    const workspaceRoot = await workspaceRootForPlan(directory);
    const config = await readData(join(workspaceRoot, "workspace.yaml")) as WorkspaceConfig;
    const workspaceErrors = contractMessages(await validateContract("workspace", config));
    workspaceErrors.push(...workspaceSemanticErrors(config));
    if (workspaceErrors.length > 0) throw new Error(`Invalid workspace configuration: ${workspaceErrors.join("; ")}`);
    const info = await lstat(directory);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error("Plan path must be a real directory");
    if (!await regularFile(join(directory, "README.md"))) throw new Error("Plan README must be a real file");
    index = parsePlanIndex(await readFile(join(directory, "README.md"), "utf8"));
    const indexErrors = contractMessages(await validateContract("plan-index", index));
    errors.push(...indexErrors);
    if (indexErrors.length > 0) return { index, work_breakdown: null, errors: [...new Set(errors)] };
    if (index.contract_version === 2) {
      breakdown = await validateRootPlanDirectory(directory, index, config, errors);
      return { index, work_breakdown: breakdown, errors: [...new Set(errors)] };
    }
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
      breakdown = parseWorkBreakdown(material.get(index.work_breakdown)!, index, config);
      errors.push(...contractMessages(await validateContract("plan-work-breakdown", breakdown)));
      errors.push(...planWorkBreakdownSemanticErrors(breakdown, config));
    }
  } catch (error) {
    errors.push((error as Error).message);
  }
  return { index, work_breakdown: breakdown, errors: [...new Set(errors)] };
}

async function actualMaterialDigest(directory: string, index: PlanIndex): Promise<string> {
  if (index.contract_version === 2) {
    const material = new Map<string, string>();
    for (const document of index.documents) material.set(document, await readFile(join(directory, document), "utf8"));
    const taskIndex = index.task_index ?? index.work_breakdown;
    material.set(taskIndex, await readFile(join(directory, taskIndex), "utf8"));
    for (const entry of await readdir(join(directory, "tasks"))) {
      if (entry === "README.md" || !entry.endsWith(".md")) continue;
      material.set(`tasks/${entry}`, await readFile(join(directory, "tasks", entry), "utf8"));
    }
    return rootMaterialDigest(material);
  }
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
    "material_digest does not match root plan documents and task files",
    "approved_digest does not match the approved root plan material",
  ]);
  const blocking = validation.errors.filter((error) => !allowedStaleDigestErrors.has(error));
  if (!validation.index || blocking.length > 0) throw new Error(`Plan state transition validation failed:\n- ${blocking.join("\n- ")}`);
  const index = validation.index;
  if (transition.kind === "lifecycle" && index.contract_version === 2 && index.status === transition.status && index.status_reason === transition.reason.trim() && index.status_actor === (transition.actor?.trim() || "engine") && index.status_evidence === (transition.evidence?.trim() || null)) {
    return index;
  }
  const digest = await actualMaterialDigest(directory, index);
  if (transition.kind === "approve") {
    if (index.status !== "draft") throw new Error("Only a draft plan can be approved");
    if (!transition.approved_by.trim()) throw new Error("Approval requires a non-empty approver");
    index.status = "approved";
    index.approved_at = now.toISOString();
    index.approved_by = transition.approved_by.trim();
    index.material_digest = digest;
    index.approved_digest = digest;
    if (index.contract_version === 2) {
      index.status_updated_at = now.toISOString();
      index.status_reason = "Human approval recorded";
      index.status_actor = transition.approved_by.trim();
      index.status_evidence = index.approved_digest;
    }
  } else if (transition.kind === "material-revision") {
    if (!["approved", "in-progress", "review-ready", "merge-pending"].includes(index.status)) throw new Error("Material revision transition requires an approved or unmerged active plan");
    if (!transition.reason.trim()) throw new Error("Material revision requires a reason");
    index.status = "draft";
    index.plan_version += 1;
    index.approved_at = null;
    index.approved_by = null;
    index.approved_digest = null;
    index.material_digest = digest;
    index.revision_reason = transition.reason.trim();
    if (index.contract_version === 2) {
      index.status_updated_at = now.toISOString();
      index.status_reason = transition.reason.trim();
      index.status_actor = "engine";
      index.status_evidence = index.material_digest;
    }
  } else if (transition.kind === "non-material-repair") {
    if (index.status !== "approved") throw new Error("Non-material repair transition requires an approved plan");
    index.material_digest = digest;
    index.approved_digest = digest;
  } else if (transition.kind === "lifecycle") {
    if (index.contract_version !== 2) throw new Error("Lifecycle plan status transitions require a root plan");
    if (!transition.reason.trim()) throw new Error("Lifecycle status transition requires a reason");
    index.status = transition.status;
    index.status_updated_at = now.toISOString();
    index.status_reason = transition.reason.trim();
    index.status_actor = transition.actor?.trim() || "engine";
    index.status_evidence = transition.evidence?.trim() || null;
  } else {
    if (index.contract_version !== 2) throw new Error("Archive transitions require a root plan");
    if (index.status === "archived") return index;
    if (index.status === "in-progress" || index.status === "merge-pending") throw new Error("Active plans must be deliberately stopped before archiving");
    index.status = "archived";
    index.archived_at = now.toISOString();
    index.status_updated_at = now.toISOString();
    index.status_reason = "Explicit archive action";
    index.status_actor = transition.actor.trim();
    index.status_evidence = transition.evidence.trim();
  }
  index.updated_at = now.toISOString();
  const readmePath = join(directory, "README.md");
  const raw = await readFile(readmePath, "utf8");
  const match = raw.match(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/);
  if (!match) throw new Error("Plan README must begin with YAML frontmatter");
  await writeTextAtomic(readmePath, raw.replace(match[0], `---\n${stringifyYaml(index).trimEnd()}\n---\n`));
  const after = await validatePlanDirectory(directory);
  if (after.errors.length > 0) throw new Error(`Plan state transition produced invalid metadata:\n- ${after.errors.join("\n- ")}`);
  if (index.contract_version === 2) await refreshRootRegistry(await workspaceRootForPlan(directory));
  return index;
}

function renderDocument(title: string, sections: Array<[string, string]>): string {
  return `# ${title}\n\n${sections.map(([heading, body]) => `## ${heading}\n\n${body}`).join("\n\n")}\n`;
}

function renderPlan(request: PlanDraftRequest, createdAt: string): { index: PlanIndex; breakdown: PlanWorkBreakdown; files: Map<string, string> } {
  const workItems = allocateWorkItems(request);
  const breakdown: PlanWorkBreakdown = { contract_version: 2, plan_id: request.plan_id, work_prefix: request.work_prefix, items: workItems };
  const files = new Map<string, string>();
  const overviewSections: Array<[string, string]> = [
    ["Summary", request.summary],
    ["Source", `${request.source.kind}: ${request.source.reference}`],
    ["Affected repositories", markdownList(request.affected_repositories, "None identified.")],
  ];
  if (request.product_knowledge) overviewSections.push(["Product impact", productImpactBody(request.product_knowledge)]);
  overviewSections.push(
    ["Assumptions", markdownList(request.assumptions, "None recorded.")],
    ["Open questions", markdownList(request.open_questions, "None recorded.")],
  );
  files.set("0001-overview.md", renderDocument("Overview", overviewSections));
  files.set("0010-requirements.md", renderDocument("Requirements", [["Requirements and acceptance criteria", markdownList(request.requirements, "None recorded.")]]));
  files.set("0020-solution.md", renderDocument("Solution", [["Proposed solution", markdownList(request.solution, "None recorded.")]]));
  files.set("0040-delivery.md", renderDocument("Delivery", [["Delivery order", markdownList(request.delivery, "None recorded.")]]));
  files.set("0050-verification.md", renderDocument("Verification", [["Verification strategy", markdownList(request.verification, "None recorded.")]]));
  files.set("0070-risks.md", renderDocument("Risks", [["Risks and mitigations", markdownList(request.risks, "None recorded.")]]));
  const rows = workItems.map((item) => `| ${item.work_id} | ${item.title} | ${item.parent ?? "—"} | ${item.depends_on.join(", ") || "—"} | ${item.repository} | ${item.area} | — |`).join("\n");
  const execution = {
    contract_version: 2,
    items: workItems.map(({ work_id, repository, scope, test_scope, test_policy, test_rationale, verification_commands, acceptance_criteria, description, subtasks, connections }) => ({
      work_id, repository, scope, test_scope, test_policy, ...(test_rationale ? { test_rationale } : {}), verification_commands, acceptance_criteria, description,
      ...(subtasks ? { subtasks } : {}), ...(connections ? { connections } : {}),
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
    ...(request.product_knowledge ? { product_knowledge: request.product_knowledge } : {}),
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

// Root-plan support is intentionally kept below the legacy helpers. Existing
// authored plans and work IDs can still be inspected and migrated, while all
// newly generated plans use the repository-scoped registry described by the
// plan-generation contract.
export const ROOT_PLAN_DOCUMENTS = [
  "overview.md",
  "requirements.md",
  "acceptance-criteria.md",
  "solution.md",
  "delivery.md",
  "verification.md",
  "risks.md",
] as const;

export interface RootPlanSummary {
  plan_id: string;
  plan_number: number;
  track: PlanTrack;
  repository: string;
  repository_collection: string;
  plan_reference: string;
  status: "draft";
  plan_version: 1;
  directory: string;
  index: string;
  documents: string[];
  task_files: string[];
  work_ids: string[];
  approval_required: true;
}

export interface PlanGenerationSummary {
  batch_id: string;
  root: string;
  plans: RootPlanSummary[];
  migrated_from_legacy: string[];
  atomic: true;
}

export interface PlanMigrationSummary {
  root: string;
  migrated: RootPlanSummary[];
  removed_legacy_root: string | null;
}

function planSlug(value: string): string {
  const normalized = value.trim().toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
  if (!normalized) throw new Error("Plan slug must contain at least one ASCII letter or digit");
  return normalized;
}

function rootPlanPath(collection: string, track: PlanTrack, number: number, slug: string): string {
  return join(collection, track === "bau" ? "__BAU__" : "", `${String(number).padStart(3, "0")}-${slug}`);
}

function rootPlanReference(collection: string, track: PlanTrack, number: number, slug: string, version = 1): string {
  const path = rootPlanPath(collection, track, number, slug).replaceAll("\\", "/");
  return `plans/${path}@v${version}`;
}

function rootMaterialDigest(files: Map<string, string>): string {
  const names = [...files.keys()].sort();
  return materialDigest(files, names);
}

const TASK_BODY_SECTIONS = [
  "Description",
  "Scope",
  "Test expectation",
  "Verification commands",
  "Acceptance criteria",
] as const;

function authoredLiveStatusErrors(raw: string): string[] {
  const errors: string[] = [];
  if (/^status\s*:/im.test(raw)) errors.push("authored Markdown must not contain a live status field");
  if (/^#{2,6}\s+status\b/im.test(raw)) errors.push("authored Markdown must not contain a live Status section");
  if (/\|\s*status\s*\|/i.test(raw)) errors.push("authored Markdown must not contain a live Status column");
  return errors;
}

function markdownSections(raw: string): { title: string | null; sections: Map<string, string>; errors: string[] } {
  const errors: string[] = [];
  const frontmatter = raw.match(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/);
  const body = frontmatter ? raw.slice(frontmatter[0].length) : raw;
  const title = body.match(/^#\s+([^\r\n]+)\r?$/m)?.[1]?.trim() ?? null;
  const headings = [...body.matchAll(/^##\s+([^\r\n]+)\r?$/gm)];
  const sections = new Map<string, string>();
  for (const [index, heading] of headings.entries()) {
    const name = heading[1]!.trim();
    if (sections.has(name)) errors.push(`duplicate Markdown section: ${name}`);
    const start = (heading.index ?? 0) + heading[0].length;
    const end = headings[index + 1]?.index ?? body.length;
    sections.set(name, body.slice(start, end).trim());
  }
  for (const section of TASK_BODY_SECTIONS) {
    if (!sections.has(section)) errors.push(`task Markdown is missing the ${section} section`);
    else if (!sections.get(section)) errors.push(`task Markdown section is empty: ${section}`);
  }
  return { title, sections, errors };
}

function parseTaskFrontmatter(raw: string): { value: unknown; body: string; errors: string[] } {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---(?:\r?\n|$)/);
  if (!match) return { value: null, body: raw, errors: ["task Markdown must begin with YAML frontmatter"] };
  try {
    return { value: parseYaml(match[1]!), body: raw.slice(match[0].length), errors: [] };
  } catch (error) {
    return { value: null, body: raw.slice(match[0].length), errors: [`invalid task frontmatter: ${(error as Error).message}`] };
  }
}

async function validatePlanTaskContract(workspaceRoot: string, value: unknown): Promise<string[]> {
  try {
    const connectionSchema = await readData(join(workspaceRoot, ".agents", "contracts", "plan-connection.schema.json")) as Record<string, unknown>;
    const taskSchema = await readData(join(workspaceRoot, ".agents", "contracts", "plan-task.schema.json")) as Record<string, unknown>;
    const ajv = new Ajv2020({ allErrors: true, strict: false });
    ajv.addSchema(connectionSchema);
    const validate = ajv.compile(taskSchema);
    return validate(value) ? [] : (validate.errors ?? []).map((error) => `${error.instancePath || "/"} ${error.message ?? "is invalid"}`);
  } catch (error) {
    return [`unable to validate task frontmatter: ${(error as Error).message}`];
  }
}

function sameStringList(left: string[] | undefined, right: string[]): boolean {
  return JSON.stringify(left ?? []) === JSON.stringify(right);
}

function normalizeConnectionList(value: unknown): Array<{ type: "depends-on" | "integrates-with" | "blocks" | "related" | "supersedes"; target: string; description?: string }> {
  if (!Array.isArray(value)) return [];
  return value.map((entry) => {
    const connection = entry as { type?: string; target?: string; description?: string };
    if (!connection || typeof connection.type !== "string" || typeof connection.target !== "string" || !connection.target.trim()) {
      throw new Error("Every plan or task connection must have a type and target");
    }
    const types = ["depends-on", "integrates-with", "blocks", "related", "supersedes"] as const;
    if (!types.includes(connection.type as typeof types[number])) throw new Error(`Unknown connection type: ${connection.type}`);
    return {
      type: connection.type as typeof types[number],
      target: connection.target.trim(),
      ...(connection.description?.trim() ? { description: connection.description.trim() } : {}),
    };
  });
}

function rootWorkItems(definition: PlanGenerationDefinition): PlanWorkItem[] {
  const ids = new Map<string, string>();
  const used = new Set<string>();
  definition.work_items.forEach((item, index) => {
    const candidate = item.work_id?.trim() || `${definition.work_prefix}-${String(index === 0 ? 1 : index * 10).padStart(3, "0")}`;
    if (!candidate.startsWith(`${definition.work_prefix}-`)) throw new Error(`Work ID ${candidate} does not use ${definition.work_prefix} prefix`);
    if (used.has(candidate)) throw new Error(`duplicate work ID: ${candidate}`);
    used.add(candidate);
    ids.set(item.key, candidate);
    ids.set(candidate, candidate);
  });
  const items = definition.work_items.map((item) => {
    const workId = ids.get(item.key)!;
    const parent = item.parent ? ids.get(item.parent) ?? item.parent : null;
    const dependsOn = (item.depends_on ?? []).map((dependency) => ids.get(dependency) ?? dependency);
    const subtasks = (item.subtasks ?? []).map((subtask) => ids.get(subtask) ?? subtask);
    return {
      work_id: workId,
      title: item.title.trim(),
      parent,
      depends_on: dependsOn,
      area: item.area.trim(),
      repository: item.repository,
      scope: [...item.scope],
      test_scope: [...item.test_scope],
      test_policy: item.test_policy,
      ...(item.test_rationale?.trim() ? { test_rationale: item.test_rationale.trim() } : {}),
      verification_commands: [...item.verification_commands],
      acceptance_criteria: [...item.acceptance_criteria],
      external_reference: null,
      description: item.description?.trim() || item.title.trim(),
      ...(subtasks.length > 0 ? { subtasks } : {}),
      ...(item.connections ? {
        connections: normalizeConnectionList(item.connections).map((connection) => ({
          ...connection,
          target: ids.get(connection.target) ?? connection.target,
        })),
      } : {}),
    };
  });
  const byId = new Map(items.map((item) => [item.work_id, item]));
  for (const item of items) {
    if (!item.parent || !byId.has(item.parent)) continue;
    const parent = byId.get(item.parent)!;
    if (!parent.subtasks?.includes(item.work_id)) parent.subtasks = [...(parent.subtasks ?? []), item.work_id];
  }
  return items;
}

function rootPlanSemanticErrors(request: PlanGenerationRequest, config: WorkspaceConfig): string[] {
  const errors: string[] = [];
  const ids = new Set<string>();
  const numbers = new Set<string>();
  const definitions = request.plans;
  for (const definition of definitions) {
    if (ids.has(definition.plan_id)) errors.push(`duplicate plan ID: ${definition.plan_id}`);
    ids.add(definition.plan_id);
    const collection = `${definition.repository}-plans`;
    if (definition.repository_collection && definition.repository_collection !== collection) {
      errors.push(`repository collection must be the exact registered collection ${collection}, not ${definition.repository_collection}`);
    }
    if (!config.repositories[definition.repository]) errors.push(`repository is not registered: ${definition.repository}`);
    const affected = [...new Set([definition.repository, ...(definition.affected_repositories ?? [])])];
    for (const repository of affected) {
      if (!config.repositories[repository]) errors.push(`repository is not registered: ${repository}`);
    }
    const track = definition.track ?? "epic";
    if (definition.plan_number !== undefined) {
      const numberKey = `${collection}/${track}/${definition.plan_number}`;
      if (numbers.has(numberKey)) errors.push(`duplicate plan number in ${collection} ${track}: ${definition.plan_number}`);
      numbers.add(numberKey);
    }
    try {
      const items = rootWorkItems(definition);
      const known = new Set(items.map((item) => item.work_id));
      for (const item of items) {
        if (item.parent && !known.has(item.parent)) errors.push(`${item.work_id} has unknown parent: ${item.parent}`);
        for (const dependency of item.depends_on) if (!known.has(dependency)) errors.push(`${item.work_id} has unknown dependency: ${dependency}`);
        for (const subtask of item.subtasks ?? []) if (!known.has(subtask)) errors.push(`${item.work_id} has unknown subtask: ${subtask}`);
      }
      errors.push(...planWorkBreakdownSemanticErrors({ contract_version: 2, plan_id: definition.plan_id, work_prefix: definition.work_prefix, items }, config));
    } catch (error) {
      errors.push((error as Error).message);
    }
  }
  const batchIds = new Set(definitions.map((definition) => definition.plan_id));
  for (const definition of definitions) {
    for (const dependency of definition.depends_on_plans ?? []) {
      if (dependency === definition.plan_id) errors.push(`plan ${definition.plan_id} cannot depend on itself`);
      if (!batchIds.has(dependency) && !dependency.startsWith("plans/")) {
        // Existing plan IDs are resolved against the registry during the atomic preflight.
        continue;
      }
    }
    try { normalizeConnectionList(definition.connections); } catch (error) { errors.push((error as Error).message); }
  }
  errors.push(...cycleErrors(definitions.map((definition) => ({ work_id: definition.plan_id, depends_on: definition.depends_on_plans ?? [] }))));
  return [...new Set(errors)];
}

function renderTaskFile(planId: string, item: PlanWorkItem): string {
  const frontmatter = {
    task_id: item.work_id,
    plan_id: planId,
    repository: item.repository,
    parent_task: item.parent,
    depends_on: item.depends_on,
    ...(item.subtasks && item.subtasks.length > 0 ? { subtasks: item.subtasks } : {}),
    ...(item.connections ? { connections: item.connections } : { connections: [] }),
  };
  return `---\n${stringifyYaml(frontmatter).trimEnd()}\n---\n\n# ${item.title}\n\n## Description\n\n${item.description ?? item.title}\n\n## Scope\n\n${markdownList(item.scope, "No implementation scope recorded.")}\n\n## Test expectation\n\n- Policy: ${item.test_policy}\n- Paths: ${item.test_scope.length > 0 ? item.test_scope.join(", ") : "None"}\n${item.test_rationale ? `- Rationale: ${item.test_rationale}\n` : ""}\n## Verification commands\n\n${markdownList(item.verification_commands, "No command recorded.")}\n\n## Acceptance criteria\n\n${markdownList(item.acceptance_criteria, "None recorded.")}\n`;
}

function renderRootPlan(definition: PlanGenerationDefinition, number: number, createdAt: string, collection: string): { index: PlanIndex; files: Map<string, string>; taskFiles: string[]; folder: string } {
  const track = definition.track ?? "epic";
  const slug = planSlug(definition.slug ?? definition.plan_id);
  const folder = rootPlanPath(collection, track, number, slug);
  const planReference = rootPlanReference(collection, track, number, slug, 1);
  const source = definition.source ?? { kind: "document", reference: "generated request" };
  const items = rootWorkItems(definition);
  const taskFiles = items.map((item) => `tasks/${item.work_id}.md`);
  const files = new Map<string, string>();
  files.set("overview.md", renderDocument("Overview", [
    ["Summary", definition.summary],
    ["Source", `${source.kind}: ${source.reference}`],
    ["Affected repositories", markdownList([...new Set([definition.repository, ...(definition.affected_repositories ?? [])])], "None identified.")],
    ["Assumptions", markdownList(definition.assumptions, "None recorded.")],
    ["Open questions", markdownList(definition.open_questions, "None recorded.")],
  ]));
  files.set("requirements.md", renderDocument("Requirements", [["Requirements and acceptance criteria", markdownList(definition.requirements, "None recorded.")]]));
  files.set("acceptance-criteria.md", renderDocument("Acceptance criteria", [["Plan acceptance", markdownList([...new Set(definition.work_items.flatMap((item) => item.acceptance_criteria))], "None recorded.")]]));
  files.set("solution.md", renderDocument("Solution", [["Proposed solution", markdownList(definition.solution, "None recorded.")]]));
  files.set("delivery.md", renderDocument("Delivery", [["Delivery order", markdownList(definition.delivery, "None recorded.")]]));
  files.set("verification.md", renderDocument("Verification", [["Verification strategy", markdownList(definition.verification, "None recorded.")]]));
  files.set("risks.md", renderDocument("Risks", [["Risks and mitigations", markdownList(definition.risks, "None recorded.")]]));
  const execution = { contract_version: 2, plan_id: definition.plan_id, work_prefix: definition.work_prefix, items };
  const taskLinks = items.map((item) => `- [${item.work_id}](./${item.work_id}.md)`).join("\n");
  files.set("tasks/README.md", `# Tasks\n\nThe task files below are the complete immutable task graph for this plan.\n\n## Task files\n\n${taskLinks}\n\n## Task contracts\n\n\`\`\`json\n${JSON.stringify(execution, null, 2)}\n\`\`\`\n`);
  for (const item of items) files.set(`tasks/${item.work_id}.md`, renderTaskFile(definition.plan_id, item));
  const digest = rootMaterialDigest(files);
  const affected = [...new Set([definition.repository, ...(definition.affected_repositories ?? [])])];
  const index: PlanIndex = {
    contract_version: 2,
    plan_id: definition.plan_id,
    title: definition.title,
    repository_collection: collection,
    track,
    plan_number: number,
    plan_reference: planReference,
    status: "draft",
    status_updated_at: createdAt,
    status_reason: "Generated plan draft",
    status_actor: "engine",
    status_evidence: planReference,
    archived_at: null,
    plan_version: 1,
    approved_at: null,
    approved_by: null,
    revision_reason: "Initial draft",
    source,
    source_reference: source.reference,
    work_prefix: definition.work_prefix,
    documents: [...ROOT_PLAN_DOCUMENTS],
    work_breakdown: "tasks/README.md",
    task_index: "tasks/README.md",
    material_digest: digest,
    approved_digest: null,
    affected_repositories: affected,
    depends_on_plans: [...(definition.depends_on_plans ?? [])],
    connections: normalizeConnectionList(definition.connections),
    created_at: createdAt,
    updated_at: createdAt,
    ...(definition.product_knowledge ? { product_knowledge: definition.product_knowledge } : {}),
  };
  const links = [...ROOT_PLAN_DOCUMENTS, "tasks/README.md"].map((document) => `- [${document}](./${document})`).join("\n");
  files.set("README.md", `---\n${stringifyYaml(index).trimEnd()}\n---\n\n# ${definition.title}\n\n${definition.summary}\n\n## Plan documents\n\n${links}\n\n## Plan identity\n\n- Collection: ${collection}\n- Track: ${track}\n- Stable number: ${number}\n- Plan reference: ${planReference}\n\nHuman approval is required before execution. Live task status belongs in runtime evidence, not this plan.\n`);
  return { index, files, taskFiles, folder };
}

async function regularDirectory(path: string): Promise<boolean> {
  try {
    const info = await lstat(path);
    return info.isDirectory() && !info.isSymbolicLink();
  } catch { return false; }
}

async function ensureRootPlansRoot(workspaceRoot: string): Promise<string> {
  const root = assertInside(workspaceRoot, join(workspaceRoot, "plans"));
  try {
    const info = await lstat(root);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`Plan root must be a real directory: ${root}`);
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
    await mkdir(root, { recursive: true, mode: 0o755 });
  }
  return root;
}

async function readRootIndexes(root: string): Promise<PlanIndex[]> {
  const indexes: PlanIndex[] = [];
  for (const collectionEntry of await readdir(root, { withFileTypes: true })) {
    if (!collectionEntry.isDirectory() || collectionEntry.name.startsWith(".")) continue;
    const collection = join(root, collectionEntry.name);
    const pending = [collection];
    while (pending.length > 0) {
      const directory = pending.pop()!;
      for (const entry of await readdir(directory, { withFileTypes: true })) {
        if (!entry.isDirectory() || entry.name.startsWith(".")) continue;
        const candidate = join(directory, entry.name);
        if (entry.name === "__BAU__") { pending.push(candidate); continue; }
        try {
          const raw = await readFile(join(candidate, "README.md"), "utf8");
          const index = parsePlanIndex(raw);
          if (index.contract_version === 2) indexes.push(index);
        } catch { /* non-plan directories do not define registry entries */ }
      }
    }
  }
  return indexes;
}

async function readRootIndexesIfPresent(root: string): Promise<PlanIndex[]> {
  return await regularDirectory(root) ? readRootIndexes(root) : [];
}

function planDependencyKey(value: string): string {
  return value.replace(/^plans\//, "").replace(/@v[0-9]+$/, "");
}

async function assertRootPlanReferences(root: string, definitions: PlanGenerationDefinition[], rendered: Array<{ index: PlanIndex; files: Map<string, string> }>): Promise<void> {
  const existing = await regularDirectory(root) ? await readRootIndexes(root) : [];
  const all = [...existing, ...rendered.map((item) => item.index)];
  const graph = new Map<string, string[]>();
  for (const index of all) graph.set(index.plan_id, []);
  for (const [position, definition] of definitions.entries()) {
    const plan = rendered[position]!.index;
    const refs = [...(plan.depends_on_plans ?? []), ...(plan.connections ?? []).map((connection) => connection.target)];
    for (const reference of refs) {
      const matches = matchingPlanIndexes(all, reference);
      if (matches.length === 0) throw new Error(`Plan ${definition.plan_id} references unknown plan: ${reference}`);
      if (matches.length > 1) throw new Error(`Plan reference is ambiguous: ${reference}`);
      if ((plan.depends_on_plans ?? []).includes(reference) || plan.connections?.some((connection) => connection.type === "depends-on" && connection.target === reference)) {
        graph.get(plan.plan_id)!.push(matches[0]!.plan_id);
      }
    }
    const taskGraph = JSON.parse(rendered[position]!.files.get("tasks/README.md")!.match(/```json\n([\s\S]*?)\n```/)![1]!) as PlanWorkBreakdown;
    const taskIds = new Set(taskGraph.items.map((item) => item.work_id));
    for (const item of taskGraph.items) for (const connection of item.connections ?? []) {
      if (taskIds.has(connection.target)) continue;
      const matches = matchingPlanIndexes(all, connection.target);
      if (matches.length === 0) throw new Error(`Task ${item.work_id} references unknown plan or task: ${connection.target}`);
      if (matches.length > 1) throw new Error(`Task connection target is ambiguous: ${connection.target}`);
    }
    if (definition.plan_number !== undefined && rendered[position]?.index.plan_number !== definition.plan_number) throw new Error(`Plan number allocation changed during validation for ${definition.plan_id}`);
  }
  const cycle = cycleErrors([...graph.entries()].map(([work_id, depends_on]) => ({ work_id, depends_on })));
  if (cycle.length > 0) throw new Error(cycle.join("; ").replaceAll("work dependency cycle", "plan dependency cycle"));
}

async function validateRenderedRootPlan(directory: string, expected: PlanIndex): Promise<void> {
  const readme = await readFile(join(directory, "README.md"), "utf8");
  const index = parsePlanIndex(readme);
  const errors = contractMessages(await validateContract("plan-index", index));
  if (errors.length > 0) throw new Error(`Generated root plan index is invalid: ${errors.join("; ")}`);
  if (index.plan_id !== expected.plan_id || index.plan_reference !== expected.plan_reference) throw new Error("Generated plan identity changed during validation");
  const workspaceRoot = await workspaceRootForPlan(directory);
  const config = await readData(join(workspaceRoot, "workspace.yaml")) as WorkspaceConfig;
  const generatedErrors: string[] = [];
  await validateRootPlanDirectory(directory, index, config, generatedErrors, false, [index]);
  if (generatedErrors.length > 0) throw new Error(`Generated root plan is invalid: ${generatedErrors.join("; ")}`);
}

function collectionReadme(indexes: PlanIndex[], collection: string): string {
  const rows = indexes.filter((index) => index.repository_collection === collection).sort((a, b) => (a.plan_number ?? 0) - (b.plan_number ?? 0)).map((index) => `| ${String(index.plan_number).padStart(3, "0")} | ${index.plan_id} | ${index.track} | ${index.status} | ${index.plan_reference} | ${(index.depends_on_plans ?? []).join(", ") || "—"} |`).join("\n");
  return `# ${collection}\n\nPeer plans are ordered by stable number; explicit dependencies remain authoritative.\n\n| Number | Plan | Track | Status | Reference | Depends on |\n| --- | --- | --- | --- | --- | --- |\n${rows || "| — | No plans | — | — | — | — |"}\n`;
}

function rootReadme(indexes: PlanIndex[]): string {
  const collections = [...new Set(indexes.map((index) => index.repository_collection).filter((value): value is string => Boolean(value)))].sort();
  const rows = indexes.sort((a, b) => `${a.repository_collection}/${a.plan_number}`.localeCompare(`${b.repository_collection}/${b.plan_number}`)).map((index) => `| ${index.repository_collection} | ${String(index.plan_number).padStart(3, "0")} | ${index.plan_id} | ${index.track} | ${index.status} | ${index.plan_reference} |`).join("\n");
  return `# Plans roadmap\n\nThe root registry groups complete peer plans by exact registered repository collection. A conceptual area is never a repository identity.\n\nCollections: ${collections.join(", ") || "None"}\n\n| Collection | Number | Plan | Track | Status | Reference |\n| --- | --- | --- | --- | --- | --- |\n${rows || "| — | — | No plans | — | — | — |"}\n`;
}

async function writeStagedRootPlan(stageRoot: string, rendered: { index: PlanIndex; files: Map<string, string>; folder: string }): Promise<void> {
  const directory = join(stageRoot, rendered.folder);
  await mkdir(directory, { recursive: true, mode: 0o755 });
  for (const [name, contents] of rendered.files) await writeTextExclusive(join(directory, name), contents);
  await validateRenderedRootPlan(directory, rendered.index);
}

function rootExistingIndexByIdentity(indexes: PlanIndex[], planId: string, collection: string, track: PlanTrack, number: number, slug: string): string | null {
  for (const index of indexes) {
    if (index.plan_id === planId) return `plan ID collision: ${planId}`;
    if (index.repository_collection === collection && index.track === track && index.plan_number === number) return `plan number collision in ${collection} ${track}: ${number}`;
    if (index.repository_collection === collection && index.plan_number === number && index.track === track && index.plan_reference?.includes(`/${String(number).padStart(3, "0")}-${slug}`)) return `plan folder collision: ${slug}`;
  }
  return null;
}

async function atomicInstallPlans(root: string, stageRoot: string, oldRoot: string | null): Promise<void> {
  const backup = `${root}.backup-${randomUUID()}`;
  const existed = await regularDirectory(root);
  try {
    if (existed) await rename(root, backup);
    await rename(stageRoot, root);
    if (existed) await rm(backup, { recursive: true, force: true });
  } catch (error) {
    try { if (await regularDirectory(root)) await rm(root, { recursive: true, force: true }); } catch { /* preserve original failure */ }
    try { if (existed && await regularDirectory(backup)) await rename(backup, root); } catch { /* surface original failure */ }
    throw error;
  }
  if (oldRoot && await regularDirectory(oldRoot)) await rm(oldRoot, { recursive: true, force: true });
}

async function normalizedGenerationRequest(request: PlanGenerationRequest): Promise<PlanGenerationRequest> {
  return {
    contract_version: 2,
    source: { kind: request.source.kind, reference: request.source.reference.trim() },
    plans: request.plans.map((definition) => ({
      ...definition,
      repository: definition.repository.trim(),
      ...(definition.repository_collection ? { repository_collection: definition.repository_collection.trim() } : {}),
      track: definition.track ?? "epic",
      slug: definition.slug?.trim() || planSlug(definition.plan_id),
      ...(definition.source ? { source: { kind: definition.source.kind, reference: definition.source.reference.trim() } } : {}),
    })),
  };
}

export async function generatePlanBatch(workspaceRootInput: string, input: PlanGenerationRequest, now = new Date()): Promise<PlanGenerationSummary> {
  const workspaceRoot = resolve(workspaceRootInput);
  const contractErrors = contractMessages(await validateContract("plan-generation-request", input));
  const config = await readData(join(workspaceRoot, "workspace.yaml")) as WorkspaceConfig;
  const workspaceErrors = contractMessages(await validateContract("workspace", config));
  const errors = [...contractErrors, ...workspaceErrors];
  if (contractErrors.length === 0 && workspaceErrors.length === 0) {
    const normalized = await normalizedGenerationRequest(input);
    errors.push(...workspaceSemanticErrors(config), ...rootPlanSemanticErrors(normalized, config));
    if (errors.length > 0) throw new Error(`Invalid plan generation request:\n- ${errors.join("\n- ")}`);
    input = normalized;
  }
  if (errors.length > 0) throw new Error(`Invalid plan generation request:\n- ${errors.join("\n- ")}`);
  const root = assertInside(workspaceRoot, join(workspaceRoot, "plans"));
  let rootExists = false;
  try {
    const info = await lstat(root);
    if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`Plan root must be a real directory: ${root}`);
    rootExists = true;
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
  }
  const existing = rootExists ? await readRootIndexes(root) : [];
  const stageRoot = join(dirname(root), `.plans-${randomUUID()}.tmp`);
  const migratedFromLegacy: string[] = [];
  const legacyRoot = join(workspaceRoot, "context", "plans");
  try {
    if (rootExists) await cp(root, stageRoot, { recursive: true });
    else await mkdir(stageRoot, { recursive: true, mode: 0o755 });
    const rendered: Array<{ index: PlanIndex; files: Map<string, string>; folder: string; taskFiles: string[] }> = [];
    const nextNumbers = new Map<string, number>();
    for (const index of existing) nextNumbers.set(`${index.repository_collection}/${index.track}`, Math.max(nextNumbers.get(`${index.repository_collection}/${index.track}`) ?? 0, index.plan_number ?? 0));
    for (const definition of input.plans) {
      const collection = `${definition.repository}-plans`;
      const track = definition.track ?? "epic";
      const priorNumber = nextNumbers.get(`${collection}/${track}`) ?? 0;
      const number = definition.plan_number ?? (priorNumber + 1);
      if (definition.plan_number !== undefined && number !== priorNumber + 1) {
        throw new Error(`Plan numbers must be allocated monotonically in ${collection} ${track}; expected ${priorNumber + 1}, received ${number}`);
      }
      nextNumbers.set(`${collection}/${track}`, number);
      const slug = planSlug(definition.slug ?? definition.plan_id);
      const collision = rootExistingIndexByIdentity([...existing, ...rendered.map((item) => item.index)], definition.plan_id, collection, track, number, slug);
      if (collision) throw new Error(collision);
      const result = renderRootPlan(definition, number, now.toISOString(), collection);
      rendered.push({ ...result, taskFiles: result.taskFiles });
    }
    await assertRootPlanReferences(root, input.plans, rendered);
    const allIndexes = [...existing, ...rendered.map((item) => item.index)];
    const collections = [...new Set(allIndexes.map((index) => index.repository_collection).filter((value): value is string => Boolean(value)))];
    await writeTextAtomic(join(stageRoot, "README.md"), rootReadme(allIndexes));
    for (const collection of collections) {
      const collectionDirectory = join(stageRoot, collection);
      await mkdir(collectionDirectory, { recursive: true, mode: 0o755 });
      await writeTextAtomic(join(collectionDirectory, "README.md"), collectionReadme(allIndexes, collection));
      const bau = allIndexes.some((index) => index.repository_collection === collection && index.track === "bau");
      if (bau) await mkdir(join(collectionDirectory, "__BAU__"), { recursive: true, mode: 0o755 });
    }
    for (const item of rendered) await writeStagedRootPlan(stageRoot, item);
    const stageIndexes = await readRootIndexes(stageRoot);
    if (stageIndexes.length !== allIndexes.length) throw new Error("Atomic plan generation produced an incomplete registry");
    for (const index of stageIndexes) {
      const directory = join(stageRoot, relative(root, resolve(root, index.plan_reference!.replace(/^plans\//, "").replace(/@v[0-9]+$/, ""))));
      if (!await regularDirectory(directory)) {
        // Existing indexes are already validated by their own generation history;
        // newly written indexes were fully validated before installation.
        continue;
      }
    }
    await atomicInstallPlans(root, stageRoot, null);
    return { batch_id: randomUUID(), root, plans: rendered.map((item) => ({ plan_id: item.index.plan_id, plan_number: item.index.plan_number!, track: item.index.track!, repository: item.index.repository_collection!.replace(/-plans$/, ""), repository_collection: item.index.repository_collection!, plan_reference: item.index.plan_reference!, status: "draft", plan_version: 1, directory: join(root, item.folder), index: join(root, item.folder, "README.md"), documents: [...ROOT_PLAN_DOCUMENTS], task_files: item.taskFiles, work_ids: (JSON.parse(item.files.get("tasks/README.md")!.match(/```json\n([\s\S]*?)\n```/)![1]!) as PlanWorkBreakdown).items.map((work) => work.work_id), approval_required: true })), migrated_from_legacy: migratedFromLegacy, atomic: true };
  } catch (error) {
    await rm(stageRoot, { recursive: true, force: true });
    throw error;
  }
}

export async function migrateCurrentPlans(workspaceRootInput: string, now = new Date()): Promise<PlanMigrationSummary> {
  const workspaceRoot = resolve(workspaceRootInput);
  const legacyRoot = join(workspaceRoot, "context", "plans");
  if (!await regularDirectory(legacyRoot)) return { root: join(workspaceRoot, "plans"), migrated: [], removed_legacy_root: null };
  const config = await readData(join(workspaceRoot, "workspace.yaml")) as WorkspaceConfig;
  const legacyEntries = (await readdir(legacyRoot, { withFileTypes: true })).filter((entry) => entry.isDirectory() && !entry.name.startsWith("."));
  if (legacyEntries.length === 0) return { root: join(workspaceRoot, "plans"), migrated: [], removed_legacy_root: null };
  const definitions: PlanGenerationDefinition[] = [];
  const statuses: Array<"draft" | "approved"> = [];
  for (const entry of legacyEntries) {
    const directory = join(legacyRoot, entry.name);
    const validation = await validatePlanDirectory(directory, entry.name);
    if (!validation.index || !validation.work_breakdown || validation.errors.length > 0) throw new Error(`Cannot migrate current plan ${entry.name}: ${validation.errors.join("; ")}`);
    const repositories = [...new Set(validation.work_breakdown.items.map((item) => item.repository))];
    if (repositories.length !== 1) throw new Error(`Cannot migrate ${entry.name}: plan identity is ambiguous across repositories`);
    const repository = repositories[0]!;
    if (!config.repositories[repository]) throw new Error(`Cannot migrate ${entry.name}: repository is not registered: ${repository}`);
    const byId = new Map(validation.work_breakdown.items.map((item) => [item.work_id, item]));
    const items = validation.work_breakdown.items.map((item) => ({
      key: item.work_id.toLowerCase(), work_id: item.work_id, title: item.title, area: item.area, repository: item.repository,
      scope: item.scope, test_scope: item.test_scope, test_policy: item.test_policy, ...(item.test_rationale ? { test_rationale: item.test_rationale } : {}), verification_commands: item.verification_commands, acceptance_criteria: item.acceptance_criteria,
      ...(item.parent ? { parent: item.parent.toLowerCase() } : {}), ...(item.depends_on.length > 0 ? { depends_on: item.depends_on.map((dependency) => dependency.toLowerCase()) } : {}),
      ...(item.description ? { description: item.description } : {}), ...(item.subtasks ? { subtasks: item.subtasks.map((subtask) => subtask.toLowerCase()) } : {}), ...(item.connections ? { connections: item.connections } : {}),
    }));
    const source = validation.index.source;
    definitions.push({
      plan_id: validation.index.plan_id,
      title: validation.index.title,
      repository,
      track: "epic",
      plan_number: definitions.length + 1,
      slug: entry.name,
      source,
      work_prefix: validation.index.work_prefix,
      summary: `Migrated plan ${validation.index.title}.`,
      affected_repositories: [repository], assumptions: [], open_questions: [], requirements: ["Preserve the approved plan material and stable work IDs."], solution: ["Use the migrated plan as the new root-plan source."], delivery: ["Review the migrated plan before execution."], verification: ["Validate every migrated task contract."], risks: [], work_items: items,
    });
    statuses.push(validation.index.status === "approved" ? "approved" : "draft");
    void byId;
  }
  const root = await ensureRootPlansRoot(workspaceRoot);
  const existing = await readRootIndexes(root);
  if (existing.length > 0) throw new Error("Cannot migrate current plans into a non-empty root registry without an explicit collision-free revision");
  const request: PlanGenerationRequest = { contract_version: 2, source: { kind: "document", reference: "context/plans" }, plans: definitions };
  const generated = await generatePlanBatch(workspaceRoot, request, now);
  if (statuses.some((status) => status === "approved")) {
    for (const [index, status] of statuses.entries()) if (status === "approved") {
      const summary = generated.plans[index]!;
      await setPlanState(summary.directory, { kind: "approve", approved_by: "migration" }, now);
    }
  }
  // The generation transaction has already installed the new root. Removing
  // the old root is the final migration action and is never attempted when
  // validation or collision preflight fails.
  if (await regularDirectory(legacyRoot)) await rm(legacyRoot, { recursive: true, force: true });
  return { root, migrated: generated.plans, removed_legacy_root: legacyRoot };
}

export async function resolveRootPlanDirectory(workspaceRootInput: string, reference: string): Promise<string> {
  const workspaceRoot = resolve(workspaceRootInput);
  const root = assertInside(workspaceRoot, join(workspaceRoot, "plans"));
  const trimmed = reference.trim();
  if (!trimmed || trimmed.startsWith("context/plans") || trimmed.includes("\\") || trimmed.split("/").includes("..")) throw new Error(`Legacy or invalid plan reference is not executable: ${reference}`);
  const withoutVersion = trimmed.replace(/@v[0-9]+$/, "");
  if (withoutVersion.startsWith("plans/")) {
    const candidate = assertInside(root, join(workspaceRoot, withoutVersion));
    if (await regularDirectory(candidate)) return candidate;
    throw new Error(`Plan reference does not resolve: ${reference}`);
  }
  const indexes = await readRootIndexes(root);
  const exact = indexes.filter((index) => index.plan_reference === trimmed || planDependencyKey(index.plan_reference ?? "") === withoutVersion || index.plan_id === trimmed);
  if (exact.length === 0) throw new Error(`Plan reference does not resolve: ${reference}`);
  if (exact.length > 1) throw new Error(`Plan reference is ambiguous: ${reference}`);
  const index = exact[0]!;
  const path = index.plan_reference!.replace(/^plans\//, "").replace(/@v[0-9]+$/, "");
  return assertInside(root, join(workspaceRoot, "plans", path));
}

async function refreshRootRegistry(workspaceRootInput: string): Promise<void> {
  const workspaceRoot = resolve(workspaceRootInput);
  const root = await ensureRootPlansRoot(workspaceRoot);
  const indexes = await readRootIndexes(root);
  await writeTextAtomic(join(root, "README.md"), rootReadme(indexes));
  for (const collection of [...new Set(indexes.map((index) => index.repository_collection).filter((value): value is string => Boolean(value)))]) {
    const directory = join(root, collection);
    await mkdir(directory, { recursive: true, mode: 0o755 });
    await writeTextAtomic(join(directory, "README.md"), collectionReadme(indexes, collection));
  }
}

export async function archivePlan(workspaceRootInput: string, reference: string, actor: string, evidence: string, now = new Date()): Promise<{ source: string; destination: string; plan_id: string }> {
  if (!actor.trim() || !evidence.trim()) throw new Error("Archiving requires an actor and credential-free evidence");
  const workspaceRoot = resolve(workspaceRootInput);
  const source = await resolveRootPlanDirectory(workspaceRoot, reference);
  const validation = await validatePlanDirectory(source);
  if (!validation.index || validation.errors.length > 0 || validation.index.contract_version !== 2) throw new Error(`Only a valid root plan can be archived: ${validation.errors.join("; ")}`);
  const index = validation.index;
  if (index.status === "archived") return { source, destination: source, plan_id: index.plan_id };
  if (index.status === "in-progress" || index.status === "merge-pending") throw new Error(`Plan ${index.plan_id} is active and cannot be archived`);
  const raw = await readFile(join(source, "README.md"), "utf8");
  const match = raw.match(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/);
  if (!match) throw new Error("Plan README must begin with YAML frontmatter");
  index.status = "archived";
  index.archived_at = now.toISOString();
  index.status_updated_at = now.toISOString();
  index.status_reason = "Explicit archive action";
  index.status_actor = actor.trim();
  index.status_evidence = evidence.trim();
  const slug = basename(source).replace(/^[0-9]{3,}-/, "");
  const destination = assertInside(workspaceRoot, join(workspaceRoot, "archived", "plans", index.repository_collection!, index.track === "bau" ? "__BAU__" : "", slug));
  if (await lstat(destination).then(() => true).catch(() => false)) throw new Error(`Archive destination collision: ${destination}`);
  const temporary = `${source}.archive-${randomUUID()}`;
  await writeTextAtomic(join(source, "README.md"), raw.replace(match[0], `---\n${stringifyYaml(index).trimEnd()}\n---\n`));
  try {
    await mkdir(dirname(destination), { recursive: true, mode: 0o755 });
    await rename(source, temporary);
    await rename(temporary, destination);
  } catch (error) {
    if (await regularDirectory(temporary)) await rename(temporary, source);
    throw error;
  }
  await refreshRootRegistry(workspaceRoot);
  return { source, destination, plan_id: index.plan_id };
}
