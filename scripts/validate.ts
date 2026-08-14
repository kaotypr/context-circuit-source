import { dirname, relative, resolve } from "node:path"
import { fileURLToPath } from "node:url"
import { parse as parseYaml } from "yaml"
import { readFile } from "node:fs/promises"
import { workspaceDocumentErrors, workspaceSemanticErrors } from "./lib/validation.js"
import { listPlans } from "./lib/plans.js"
import { validateProductKnowledgeTree } from "./lib/product-knowledge.js"
import type { WorkspaceConfig } from "./lib/types.js"

const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
const config = parseYaml(await readFile(resolve(root, "workspace.yaml"), "utf8")) as WorkspaceConfig
const errors = [...workspaceSemanticErrors(config), ...(await workspaceDocumentErrors(root, config))]
const knowledge = await validateProductKnowledgeTree(resolve(root, "context"))
errors.push(...knowledge.errors.map((error) => `product-knowledge ${error}`))
const plans = await listPlans(root, true)
const planIds = new Set<string>()
const taskIds = new Set<string>()
for (const plan of plans) {
  errors.push(...plan.errors.map((error) => `${plan.directory}: ${error}`))
  if (plan.plan) {
    if (planIds.has(plan.plan.id)) errors.push(`duplicate plan id: ${plan.plan.id}`)
    planIds.add(plan.plan.id)
    for (const task of plan.tasks) {
      if (taskIds.has(task.id)) errors.push(`duplicate task id: ${task.id}`)
      taskIds.add(task.id)
    }
  }
}
for (const plan of plans) {
  if (!plan.plan) continue
  for (const dependency of plan.plan.dependencies ?? []) if (!plans.some((candidate) => candidate.plan?.id === dependency)) errors.push(`${plan.plan.id}: plan dependency does not resolve: ${dependency}`)
  for (const connection of plan.plan.connections ?? []) if (!plans.some((candidate) => candidate.plan?.id === connection.target) && !plan.tasks.some((task) => task.id === connection.target)) errors.push(`${plan.plan.id}: connection target does not resolve: ${connection.target}`)
}
const planEdges = new Map(plans.filter((item) => item.plan && !item.archived).map((item) => [item.plan!.id, item.plan!.dependencies ?? []]))
const visited = new Set<string>()
const active = new Set<string>()
function visitPlan(id: string, path: string[]): void {
  if (active.has(id)) {
    errors.push(`plan dependency cycle: ${[...path, id].join(" -> ")}`)
    return
  }
  if (visited.has(id)) return
  active.add(id)
  for (const dependency of planEdges.get(id) ?? []) visitPlan(dependency, [...path, id])
  active.delete(id)
  visited.add(id)
}
for (const id of planEdges.keys()) visitPlan(id, [])
const activePaths = new Set<string>()
const archivedPaths = new Set<string>()
for (const item of plans) {
  const planRoot = item.archived ? resolve(root, "archives", "plans") : resolve(root, "plans")
  const path = relative(planRoot, item.directory).replaceAll("\\", "/")
  const target = item.archived ? archivedPaths : activePaths
  if (target.has(path)) errors.push(`duplicate plan path: ${path}`)
  target.add(path)
  if (item.archived && activePaths.has(path)) errors.push(`active and archived plan paths collide: ${path}`)
  if (!item.archived && archivedPaths.has(path)) errors.push(`active and archived plan paths collide: ${path}`)
}
if (errors.length > 0) {
  console.error("Invalid Context Circuit workspace:")
  for (const error of [...new Set(errors)]) console.error(`- ${error}`)
  process.exitCode = 1
} else {
  console.log(`Valid Context Circuit workspace (${plans.length} plans; ${knowledge.pages} Product Knowledge pages)`)
}
