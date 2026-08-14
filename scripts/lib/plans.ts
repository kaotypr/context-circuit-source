import { lstat, mkdir, readdir, readFile, rename, rm } from "node:fs/promises"
import { basename, dirname, join, relative, resolve } from "node:path"
import { parse as parseYaml, stringify as stringifyYaml } from "yaml"
import { assertInside, writeTextAtomic, writeTextExclusive } from "./io.js"
import { readData, parseFrontmatter, requiredString, optionalStringList, referencedLocalFilesExist, workspaceErrors } from "./validation.js"
import type { PlanConnection, PlanCreateRequest, PlanCreationSummary, PlanListItem, PlanStatus, PlanTask, PlanTaskInput, PlanValidationResult, PlanYaml, TaskStatus, WorkspaceConfig } from "./types.js"

export const PLAN_DOCUMENTS = [
  "overview.md",
  "requirements.md",
  "acceptance-criteria.md",
  "solution.md",
  "delivery.md",
  "verification.md",
  "risks.md",
] as const

const statuses = new Set<PlanStatus>(["draft", "approved", "done"])

function isRecord(value: unknown): value is Record<string, unknown> {
  return Boolean(value) && typeof value === "object" && !Array.isArray(value)
}

function nonEmpty(value: unknown): value is string {
  return typeof value === "string" && value.trim().length > 0
}

function slugify(value: string): string {
  const slug = value.trim().toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "")
  if (!slug) throw new Error("Plan slug must contain at least one ASCII letter or digit")
  return slug
}

function numberFolder(number: number, slug: string): string {
  return `${String(number).padStart(4, "0")}-${slug}`
}

function collectionForPlan(plan: PlanYaml): string {
  const repository = plan.repositories[0]
  if (!repository) throw new Error(`Plan ${plan.id} must name at least one repository`)
  return `${repository}-plans`
}

function markdownList(values: string[] | undefined, empty = "None recorded."): string {
  return values && values.length > 0 ? values.map((value) => `- ${value}`).join("\n") : `- ${empty}`
}

function renderDocument(title: string, body: string): string {
  return `# ${title}\n\n${body.trim()}\n`
}

function normalizeConnections(value: unknown): PlanConnection[] | undefined {
  if (value === undefined) return undefined
  if (!Array.isArray(value)) throw new Error("connections must be a list")
  return value.map((entry) => {
    if (!isRecord(entry) || !nonEmpty(entry.type) || !nonEmpty(entry.target)) throw new Error("Every connection needs a type and target")
    const types = ["depends-on", "integrates-with", "blocks", "related", "supersedes"]
    if (!types.includes(entry.type)) throw new Error(`Unknown connection type: ${entry.type}`)
    return {
      type: entry.type as PlanConnection["type"],
      target: entry.target.trim(),
      ...(nonEmpty(entry.description) ? { description: entry.description.trim() } : {}),
    }
  })
}

function normalizeTaskInput(input: PlanTaskInput, planId: string, generatedId: string): PlanTask {
  const id = String(input.id ?? input.task_id ?? input.work_id ?? input.key ?? generatedId).trim()
  if (!id) throw new Error("Every task needs an id")
  const dependencies = input.dependencies ?? input.depends_on
  return {
    id,
    plan_id: planId,
    title: input.title.trim(),
    status: "draft",
    ...(input.description?.trim() ? { description: input.description.trim() } : {}),
    repository: input.repository.trim(),
    ...(input.area?.trim() ? { area: input.area.trim() } : {}),
    ...(input.parent_task?.trim() ? { parent_task: input.parent_task.trim() } : {}),
    ...(dependencies && dependencies.length > 0 ? { dependencies: [...dependencies] } : {}),
    ...(input.subtasks && input.subtasks.length > 0 ? { subtasks: [...input.subtasks] } : {}),
    ...(input.connections && input.connections.length > 0 ? { connections: normalizeConnections(input.connections) } : {}),
    ...(input.implementation_scope && input.implementation_scope.length > 0 ? { implementation_scope: [...input.implementation_scope] } : {}),
    ...(input.test_scope && input.test_scope.length > 0 ? { test_scope: [...input.test_scope] } : {}),
    ...(input.test_expectations && input.test_expectations.length > 0 ? { test_expectations: [...input.test_expectations] } : {}),
    ...(input.verification_commands && input.verification_commands.length > 0 ? { verification_commands: [...input.verification_commands] } : {}),
    ...(input.acceptance_criteria && input.acceptance_criteria.length > 0 ? { acceptance_criteria: [...input.acceptance_criteria] } : {}),
    ...(input.product_knowledge && input.product_knowledge.length > 0 ? { product_knowledge: [...input.product_knowledge] } : {}),
    ...(input.external_reference?.trim() ? { external_reference: input.external_reference.trim() } : {}),
  }
}

function taskFrontmatter(task: PlanTask): Record<string, unknown> {
  const { id, plan_id, title, status, ...optional } = task
  return { id, plan_id, title, status, ...optional }
}

function taskMarkdown(task: PlanTask): string {
  return [
    `# ${task.title}`,
    "",
    "## Description",
    "",
    task.description ?? task.title,
    "",
    "## Implementation scope",
    "",
    markdownList(task.implementation_scope),
    "",
    "## Test scope and expectations",
    "",
    markdownList(task.test_scope),
    task.test_expectations?.length ? `\n${markdownList(task.test_expectations)}` : "",
    "",
    "## Verification commands",
    "",
    markdownList(task.verification_commands, "No commands recorded; human review decides what to run."),
    "",
    "## Acceptance criteria",
    "",
    markdownList(task.acceptance_criteria, "No additional criteria recorded."),
    "",
  ].join("\n")
}

function planReadme(plan: PlanYaml, summary: string | undefined): string {
  return `# ${plan.title}\n\n${summary?.trim() || "This numbered plan is maintained as human-reviewed delivery intent."}\n\nPlan metadata lives in [plan.yaml](./plan.yaml). The plan status changes only through an explicit human request.\n\n## Documents\n\n${PLAN_DOCUMENTS.map((name) => `- [${name.replace(/\.md$/, "")}](${name})`).join("\n")}\n\n## Tasks\n\nSee [tasks/README.md](./tasks/README.md).\n`
}

async function exists(path: string): Promise<boolean> {
  try {
    await lstat(path)
    return true
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code === "ENOENT") return false
    throw error
  }
}

async function workspaceRootFor(path: string): Promise<string> {
  let current = resolve(path)
  while (true) {
    if (await exists(join(current, "workspace.yaml"))) return current
    const parent = dirname(current)
    if (parent === current) break
    current = parent
  }
  throw new Error(`Unable to locate workspace.yaml for plan: ${path}`)
}

async function loadWorkspace(root: string): Promise<WorkspaceConfig> {
  const value = await readData(join(root, "workspace.yaml"))
  const errors = workspaceErrors(value)
  if (errors.length > 0) throw new Error(`Invalid workspace.yaml:\n- ${errors.join("\n- ")}`)
  return value as WorkspaceConfig
}

function normalizePlan(value: unknown, directory?: string): PlanYaml {
  if (!isRecord(value)) throw new Error("plan.yaml must contain a mapping")
  const errors: string[] = []
  const id = requiredString(value, "id", errors) ?? ""
  const title = requiredString(value, "title", errors) ?? ""
  const status = value.status
  if (!statuses.has(status as PlanStatus)) errors.push("status must be draft, approved, or done")
  const number = typeof value.number === "number" ? value.number : Number(value.number)
  if (!Number.isInteger(number) || number < 1) errors.push("number must be a positive integer")
  const repositories = Array.isArray(value.repositories)
    ? value.repositories.filter((entry): entry is string => typeof entry === "string" && entry.trim().length > 0).map((entry) => entry.trim())
    : typeof value.repository === "string" ? [value.repository.trim()] : []
  if (repositories.length === 0) errors.push("repositories must contain at least one repository")
  if (errors.length > 0) throw new Error(`${directory ? `${directory}: ` : ""}${errors.join("; ")}`)
  return {
    id,
    number,
    title,
    status: status as PlanStatus,
    ...(nonEmpty(value.track) ? { track: value.track as PlanYaml["track"] } : {}),
    ...(isRecord(value.source) && nonEmpty(value.source.kind) && nonEmpty(value.source.reference) ? { source: { kind: value.source.kind, reference: value.source.reference } } : {}),
    repositories,
    ...(Array.isArray(value.dependencies) && value.dependencies.length > 0 ? { dependencies: value.dependencies.filter((entry): entry is string => typeof entry === "string") } : {}),
    ...(Array.isArray(value.depends_on_plans) && value.depends_on_plans.length > 0 ? { dependencies: value.depends_on_plans.filter((entry): entry is string => typeof entry === "string") } : {}),
    ...(normalizeConnections(value.connections) ? { connections: normalizeConnections(value.connections) } : {}),
    ...(isRecord(value.product_knowledge) ? { product_knowledge: value.product_knowledge as unknown as PlanYaml["product_knowledge"] } : {}),
    ...(nonEmpty(value.external_reference) ? { external_reference: value.external_reference.trim() } : {}),
    ...(nonEmpty(value.created_at) ? { created_at: value.created_at } : {}),
    ...(nonEmpty(value.updated_at) ? { updated_at: value.updated_at } : {}),
  }
}

export function parseTaskFrontmatter(raw: string): { value: PlanTask | null; body: string; errors: string[] } {
  const parsed = parseFrontmatter(raw)
  if (parsed.errors.length > 0 || !parsed.value) return { value: null, body: parsed.body, errors: parsed.errors }
  const value = parsed.value
  const errors: string[] = []
  const id = requiredString(value, "id", errors) ?? (typeof value.task_id === "string" ? value.task_id : undefined)
  const planId = requiredString(value, "plan_id", errors) ?? ""
  const title = requiredString(value, "title", errors) ?? ""
  const repository = requiredString(value, "repository", errors) ?? ""
  const status = value.status
  if (!statuses.has(status as TaskStatus)) errors.push("status must be draft, approved, or done")
  if (!id && typeof value.task_id !== "string") errors.push("id is required")
  const list = (key: string): string[] | undefined => {
    if (value[key] === undefined) return undefined
    if (!Array.isArray(value[key]) || !value[key].every((entry) => typeof entry === "string")) {
      errors.push(`${key} must be a list of strings`)
      return undefined
    }
    return (value[key] as string[]).map((entry) => entry.trim()).filter(Boolean)
  }
  const dependencies = list("dependencies") ?? list("depends_on")
  if (errors.length > 0) return { value: null, body: parsed.body, errors }
  const task: PlanTask = {
    id: id!,
    plan_id: planId,
    title,
    status: status as TaskStatus,
    repository,
    ...(nonEmpty(value.description) ? { description: value.description.trim() } : {}),
    ...(nonEmpty(value.area) ? { area: value.area.trim() } : {}),
    ...(nonEmpty(value.parent_task) ? { parent_task: value.parent_task.trim() } : {}),
    ...(dependencies?.length ? { dependencies } : {}),
    ...(list("subtasks")?.length ? { subtasks: list("subtasks") } : {}),
    ...(list("implementation_scope")?.length ? { implementation_scope: list("implementation_scope") } : {}),
    ...(list("test_scope")?.length ? { test_scope: list("test_scope") } : {}),
    ...(list("test_expectations")?.length ? { test_expectations: list("test_expectations") } : {}),
    ...(list("verification_commands")?.length ? { verification_commands: list("verification_commands") } : {}),
    ...(list("acceptance_criteria")?.length ? { acceptance_criteria: list("acceptance_criteria") } : {}),
    ...(list("product_knowledge")?.length ? { product_knowledge: list("product_knowledge") } : {}),
    ...(normalizeConnections(value.connections) ? { connections: normalizeConnections(value.connections) } : {}),
    ...(nonEmpty(value.external_reference) ? { external_reference: value.external_reference.trim() } : {}),
  }
  return { value: task, body: parsed.body, errors: [] }
}

async function readTaskFiles(directory: string): Promise<{ tasks: PlanTask[]; errors: string[] }> {
  const tasksDirectory = join(directory, "tasks")
  const errors: string[] = []
  if (!await exists(tasksDirectory)) return { tasks: [], errors: ["tasks/ directory is missing"] }
  const entries = await readdir(tasksDirectory, { withFileTypes: true })
  const tasks: PlanTask[] = []
  for (const entry of entries.sort((a, b) => a.name.localeCompare(b.name))) {
    if (entry.name === "README.md") continue
    if (entry.isDirectory() || entry.isSymbolicLink() || !entry.name.endsWith(".md")) {
      errors.push(`tasks/ contains an unexpected entry: ${entry.name}`)
      continue
    }
    const path = join(tasksDirectory, entry.name)
    const parsed = parseTaskFrontmatter(await readFile(path, "utf8"))
    errors.push(...parsed.errors.map((error) => `tasks/${entry.name}: ${error}`))
    if (parsed.value) {
      if (parsed.value.id !== entry.name.slice(0, -3)) errors.push(`tasks/${entry.name}: id must match its filename`)
      tasks.push(parsed.value)
    }
  }
  return { tasks, errors }
}

function cycleErrors(edges: Map<string, string[]>, label: string): string[] {
  const errors: string[] = []
  const visited = new Set<string>()
  const active = new Set<string>()
  const visit = (id: string, path: string[]): void => {
    if (active.has(id)) {
      errors.push(`${label} dependency cycle: ${[...path, id].join(" -> ")}`)
      return
    }
    if (visited.has(id)) return
    active.add(id)
    for (const dependency of edges.get(id) ?? []) visit(dependency, [...path, id])
    active.delete(id)
    visited.add(id)
  }
  for (const id of edges.keys()) visit(id, [])
  return errors
}

export async function validatePlanDirectory(directoryInput: string, expectedPath?: string): Promise<PlanValidationResult> {
  const directory = resolve(directoryInput)
  const errors: string[] = []
  let plan: PlanYaml | null = null
  try {
    plan = normalizePlan(parseYaml(await readFile(join(directory, "plan.yaml"), "utf8")), directory)
  } catch (error) {
    errors.push((error as Error).message)
  }
  if (expectedPath && plan && plan.id !== expectedPath) errors.push(`plan id does not match expected reference: ${expectedPath}`)
  for (const document of PLAN_DOCUMENTS) if (!await exists(join(directory, document))) errors.push(`missing plan document: ${document}`)
  const taskResult = await readTaskFiles(directory)
  errors.push(...taskResult.errors)
  const tasks = taskResult.tasks
  if (plan) {
    const ids = new Set<string>()
    const config = await loadWorkspace(await workspaceRootFor(directory)).catch(() => null)
    for (const repository of plan.repositories) if (config && !config.repositories[repository]) errors.push(`plan repository is not registered: ${repository}`)
    for (const task of tasks) {
      if (ids.has(task.id)) errors.push(`duplicate task id: ${task.id}`)
      ids.add(task.id)
      if (task.plan_id !== plan.id) errors.push(`${task.id}: plan_id does not match plan.yaml`) 
      if (config && !config.repositories[task.repository]) errors.push(`${task.id}: repository is not registered: ${task.repository}`)
      for (const dependency of task.dependencies ?? []) if (!tasks.some((candidate) => candidate.id === dependency)) errors.push(`${task.id}: dependency does not resolve: ${dependency}`)
      for (const parent of task.parent_task ? [task.parent_task] : []) if (!tasks.some((candidate) => candidate.id === parent)) errors.push(`${task.id}: parent_task does not resolve: ${parent}`)
      for (const subtask of task.subtasks ?? []) if (!tasks.some((candidate) => candidate.id === subtask)) errors.push(`${task.id}: subtask does not resolve: ${subtask}`)
      for (const connection of task.connections ?? []) if (!tasks.some((candidate) => candidate.id === connection.target)) errors.push(`${task.id}: connection target does not resolve: ${connection.target}`)
      const workspaceRoot = await workspaceRootFor(directory)
      const repositoryRoot = config?.repositories[task.repository] ? resolve(workspaceRoot, config.repositories[task.repository]!.path) : workspaceRoot
      for (const reference of [...(task.implementation_scope ?? []), ...(task.test_scope ?? [])]) errors.push(...await referencedLocalFilesExist(repositoryRoot, [reference]))
    }
    errors.push(...cycleErrors(new Map(tasks.map((task) => [task.id, task.dependencies ?? []])), "task"))
  }
  return { plan, tasks, errors: [...new Set(errors)], directory }
}

function plansRoot(root: string, archived: boolean): string {
  return archived ? join(root, "archives", "plans") : join(root, "plans")
}

function isWithin(parent: string, child: string): boolean {
  const path = relative(parent, child).replaceAll("\\", "/")
  return path === "" || (!path.startsWith("../") && path !== ".." && !path.startsWith("/"))
}

function isArchivedPlanDirectory(root: string, directory: string): boolean {
  return isWithin(plansRoot(root, true), directory)
}

function planReferenceForDirectory(root: string, directory: string): string {
  const archive = isArchivedPlanDirectory(root, directory)
  const base = plansRoot(root, archive)
  return `${archive ? "archives/plans" : "plans"}/${relative(base, directory).replaceAll("\\", "/")}`
}

async function planDirectories(root: string, archived: boolean): Promise<string[]> {
  const collectionRoot = plansRoot(root, archived)
  if (!await exists(collectionRoot)) return []
  const result: string[] = []
  for (const collection of await readdir(collectionRoot, { withFileTypes: true })) {
    if (!collection.isDirectory() || collection.isSymbolicLink() || !collection.name.endsWith("-plans")) continue
    const parent = join(collectionRoot, collection.name)
    for (const entry of await readdir(parent, { withFileTypes: true })) {
      if (!entry.isDirectory() || entry.isSymbolicLink() || !/^\d{4,}-/.test(entry.name)) continue
      result.push(join(parent, entry.name))
    }
  }
  return result.sort()
}

export async function listPlans(workspaceRootInput: string, includeArchived = true): Promise<PlanListItem[]> {
  const root = resolve(workspaceRootInput)
  const directories = [...await planDirectories(root, false), ...(includeArchived ? await planDirectories(root, true) : [])]
  const result: PlanListItem[] = []
  for (const directory of directories) {
    const validation = await validatePlanDirectory(directory)
    result.push({ ...validation, archived: isArchivedPlanDirectory(root, directory) })
  }
  return result
}

function planIdentity(plan: PlanYaml): string {
  return `${collectionForPlan(plan)}/${numberFolder(plan.number, slugify(plan.title))}`
}

export async function resolvePlanDirectory(workspaceRootInput: string, reference: string, includeArchived = false): Promise<string> {
  const root = resolve(workspaceRootInput)
  const trimmed = reference.trim().replaceAll('\\', '/')
  if (!trimmed || trimmed.includes('..')) throw new Error(`Invalid plan reference: ${reference}`)
  const active = await planDirectories(root, false)
  const archived = includeArchived ? await planDirectories(root, true) : []
  const direct = resolve(root, trimmed)
  if (await exists(join(direct, 'plan.yaml')) && (active.includes(direct) || archived.includes(direct))) return direct
  const matches: string[] = []
  for (const directory of [...active, ...archived]) {
    let candidate: PlanYaml | null = null
    try { candidate = normalizePlan(parseYaml(await readFile(join(directory, 'plan.yaml'), 'utf8'))) } catch { /* validation is intentionally deferred for location-based actions */ }
    if (candidate && (candidate.id === trimmed || planReferenceForDirectory(root, directory) === trimmed)) matches.push(directory)
  }
  if (matches.length === 0) throw new Error(`Plan reference does not resolve: ${reference}`)
  if (matches.length > 1) throw new Error(`Plan reference is ambiguous: ${reference}`)
  return matches[0]!
}

function taskPrefix(request: PlanCreateRequest, planId: string): string {
  return request.work_prefix?.trim() || planId.replace(/[^A-Za-z0-9]+/g, '').slice(0, 12).toUpperCase() || 'TASK'
}

function normalizeCreateRequest(request: PlanCreateRequest): { plan: PlanYaml; tasks: PlanTask[]; summary?: string; documents: Record<string, string> } {
  const planId = (request.id ?? request.plan_id ?? slugify(request.title)).trim()
  const repositories = [...new Set(request.repositories ?? (request.repository ? [request.repository] : []))].filter(Boolean)
  if (repositories.length === 0) throw new Error('A plan must name at least one repository')
  const number = request.number ?? 0
  if (number < 0 || !Number.isInteger(number)) throw new Error('Plan number must be a positive integer')
  const plan: PlanYaml = {
    id: planId,
    number,
    title: request.title.trim(),
    status: 'draft',
    track: request.track ?? 'epic',
    ...(request.source ? { source: request.source } : {}),
    repositories,
    ...((request.dependencies ?? request.depends_on_plans)?.length ? { dependencies: [...(request.dependencies ?? request.depends_on_plans)!] } : {}),
    ...(request.connections?.length ? { connections: normalizeConnections(request.connections) } : {}),
    ...(request.product_knowledge ? { product_knowledge: request.product_knowledge } : {}),
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  }
  const inputs = request.tasks ?? request.work_items ?? []
  const prefix = taskPrefix(request, planId)
  const tasks = inputs.map((input, index) => normalizeTaskInput(input, planId, `${prefix}-${String(index + 1).padStart(4, '0')}`))
  const requirements = request.requirements ?? []
  const acceptance = request.acceptance_criteria ?? []
  const documents: Record<string, string> = {
    'overview.md': renderDocument('Overview', [request.summary ?? request.title, '', 'Source', request.source ? `${request.source.kind}: ${request.source.reference}` : 'No source recorded.', '', 'Assumptions', markdownList(request.assumptions), '', 'Open questions', markdownList(request.open_questions)].join('\n')),
    'requirements.md': renderDocument('Requirements', markdownList(requirements)),
    'acceptance-criteria.md': renderDocument('Acceptance criteria', markdownList(acceptance)),
    'solution.md': renderDocument('Solution', markdownList(request.solution)),
    'delivery.md': renderDocument('Delivery', markdownList(request.delivery)),
    'verification.md': renderDocument('Verification', markdownList(request.verification)),
    'risks.md': renderDocument('Risks', markdownList(request.risks)),
  }
  return { plan, tasks, summary: request.summary, documents }
}

export async function createPlan(workspaceRootInput: string, request: PlanCreateRequest): Promise<PlanCreationSummary> {
  const root = resolve(workspaceRootInput)
  const config = await loadWorkspace(root)
  const normalized = normalizeCreateRequest(request)
  const errors: string[] = []
  for (const repository of normalized.plan.repositories) if (!config.repositories[repository]) errors.push(`repository is not registered: ${repository}`)
  for (const task of normalized.tasks) if (!config.repositories[task.repository]) errors.push(`${task.id}: repository is not registered: ${task.repository}`)
  const collection = collectionForPlan(normalized.plan)
  const collectionRoot = join(root, 'plans', collection)
  const archivedCollectionRoot = join(root, 'archives', 'plans', collection)
  await mkdir(collectionRoot, { recursive: true })
  const existing = [...await planDirectories(root, false), ...await planDirectories(root, true)]
  const collectionExisting = existing.filter((directory) => dirname(directory) === collectionRoot || dirname(directory) === archivedCollectionRoot)
  const numbers = (await Promise.all(collectionExisting.map(async (directory) => (await validatePlanDirectory(directory)).plan?.number ?? 0))).filter(Boolean)
  if (!normalized.plan.number) normalized.plan.number = Math.max(0, ...numbers) + 1
  const folder = numberFolder(normalized.plan.number, slugify(request.slug ?? request.title))
  const destination = assertInside(collectionRoot, join(collectionRoot, folder))
  if (await exists(destination)) errors.push(`plan path already exists: ${destination}`)
  if (collectionExisting.some((directory) => basename(directory) === folder)) errors.push(`active or archived plan path already exists: ${folder}`)
  if (errors.length > 0) throw new Error(`Cannot create plan:\n- ${errors.join('\n- ')}`)
  const temporary = join(collectionRoot, `.${folder}.tmp`)
  await rm(temporary, { recursive: true, force: true })
  try {
    await mkdir(join(temporary, 'tasks'), { recursive: true })
    await writeTextExclusive(join(temporary, 'plan.yaml'), stringifyYaml(normalized.plan))
    await writeTextExclusive(join(temporary, 'README.md'), planReadme(normalized.plan, normalized.summary))
    for (const document of PLAN_DOCUMENTS) await writeTextExclusive(join(temporary, document), normalized.documents[document]!)
    await writeTextExclusive(join(temporary, 'tasks', 'README.md'), `# Tasks\n\n${normalized.tasks.length ? normalized.tasks.map((task) => `- [${task.id}](${task.id}.md) — ${task.title}`).join('\n') : '- No tasks have been broken down yet.'}\n`)
    for (const task of normalized.tasks) {
      await writeTextExclusive(join(temporary, 'tasks', `${task.id}.md`), `---\n${stringifyYaml(taskFrontmatter(task)).trimEnd()}\n---\n\n${taskMarkdown(task)}`)
    }
    const validation = await validatePlanDirectory(temporary)
    if (validation.errors.length > 0) throw new Error(`Generated plan failed validation:\n- ${validation.errors.join('\n- ')}`)
    await rename(temporary, destination)
  } catch (error) {
    await rm(temporary, { recursive: true, force: true })
    throw error
  }
  return { plan_id: normalized.plan.id, number: normalized.plan.number, status: 'draft', directory: destination, plan: join(destination, 'plan.yaml'), tasks: normalized.tasks.map((task) => join(destination, 'tasks', `${task.id}.md`)) }
}

export async function setPlanStatus(workspaceRootInput: string, reference: string, status: PlanStatus): Promise<PlanYaml> {
  const root = resolve(workspaceRootInput)
  const directory = await resolvePlanDirectory(root, reference)
  const plan = normalizePlan(parseYaml(await readFile(join(directory, 'plan.yaml'), 'utf8')))
  plan.status = status
  plan.updated_at = new Date().toISOString()
  await writeTextAtomic(join(directory, 'plan.yaml'), stringifyYaml(plan))
  return plan
}

export async function setTaskStatus(workspaceRootInput: string, reference: string, taskId: string, status: TaskStatus): Promise<PlanTask> {
  const root = resolve(workspaceRootInput)
  const directory = await resolvePlanDirectory(root, reference)
  const taskPath = assertInside(join(directory, 'tasks'), join(directory, 'tasks', `${taskId}.md`))
  const parsed = parseTaskFrontmatter(await readFile(taskPath, 'utf8'))
  if (!parsed.value || parsed.errors.length > 0) throw new Error(`Invalid task ${taskId}: ${parsed.errors.join('; ')}`)
  parsed.value.status = status
  await writeTextAtomic(taskPath, `---\n${stringifyYaml(taskFrontmatter(parsed.value)).trimEnd()}\n---\n\n${parsed.body.trimStart()}`)
  return parsed.value
}

export async function archivePlan(workspaceRootInput: string, reference: string): Promise<{ plan_id: string; source: string; destination: string }> {
  const root = resolve(workspaceRootInput)
  const source = await resolvePlanDirectory(root, reference)
  if (isArchivedPlanDirectory(root, source)) throw new Error(`Plan is already archived: ${reference}`)
  const plan = normalizePlan(parseYaml(await readFile(join(source, 'plan.yaml'), 'utf8')))
  const destination = join(root, 'archives', 'plans', collectionForPlan(plan), basename(source))
  if (await exists(destination)) throw new Error(`Archive destination collision: ${destination}`)
  await mkdir(dirname(destination), { recursive: true })
  await rename(source, destination)
  return { plan_id: plan.id, source, destination }
}

export async function unarchivePlan(workspaceRootInput: string, reference: string): Promise<{ plan_id: string; source: string; destination: string }> {
  const root = resolve(workspaceRootInput)
  const source = await resolvePlanDirectory(root, reference, true)
  if (!isArchivedPlanDirectory(root, source)) throw new Error(`Plan is not archived: ${reference}`)
  const plan = normalizePlan(parseYaml(await readFile(join(source, 'plan.yaml'), 'utf8')))
  const destination = join(root, 'plans', collectionForPlan(plan), basename(source))
  if (await exists(destination)) throw new Error(`Unarchive destination collision: ${destination}`)
  await rename(source, destination)
  return { plan_id: plan.id, source, destination }
}

export function planReference(root: string, directory: string): string {
  return planReferenceForDirectory(root, directory)
}

export function taskFromPlan(planResult: PlanValidationResult, taskId: string): PlanTask {
  const task = planResult.tasks.find((candidate) => candidate.id === taskId)
  if (!task) throw new Error(`Task does not resolve: ${taskId}`)
  return task
}
