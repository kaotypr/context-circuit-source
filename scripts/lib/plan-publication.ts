import { readFile } from "node:fs/promises"
import { join, resolve } from "node:path"
import { parse as parseYaml, stringify as stringifyYaml } from "yaml"
import { writeTextAtomic } from "./io.js"
import { parseTaskFrontmatter, resolvePlanDirectory, validatePlanDirectory } from "./plans.js"
import type { PlanStatus, PublicationReference, PublicationResult } from "./types.js"

export interface PublicationOptions {
  workspaceRoot: string
  plan: string
  provider: string
  references?: PublicationReference[]
}

export async function publishPlan(options: PublicationOptions): Promise<PublicationResult> {
  const root = resolve(options.workspaceRoot)
  const directory = await resolvePlanDirectory(root, options.plan)
  const validation = await validatePlanDirectory(directory)
  if (!validation.plan || validation.errors.length > 0) throw new Error(`Cannot publish invalid plan:\n- ${validation.errors.join("\n- ")}`)
  const plan = validation.plan
  if (plan.status !== "approved") throw new Error(`Plan ${plan.id} is ${plan.status}; publication is offered after explicit approval and before execution`)
  const references = options.references ?? []
  const ids = new Set([plan.id, ...validation.tasks.map((task) => task.id)])
  for (const reference of references) {
    if (!ids.has(reference.id)) throw new Error(`Publication reference does not resolve to plan or task: ${reference.id}`)
    if (!/^https?:\/\//.test(reference.url)) throw new Error(`Publication reference must be an http(s) URL: ${reference.url}`)
  }
  const byId = new Map(references.map((reference) => [reference.id, reference.url]))
  if (byId.has(plan.id)) {
    plan.external_reference = byId.get(plan.id)!
    await writeTextAtomic(join(directory, "plan.yaml"), stringifyYaml(plan))
  }
  for (const task of validation.tasks) {
    const url = byId.get(task.id)
    if (!url) continue
    const path = join(directory, "tasks", `${task.id}.md`)
    const parsed = parseTaskFrontmatter(await readFile(path, "utf8"))
    if (!parsed.value || parsed.errors.length > 0) throw new Error(`Invalid task while publishing ${task.id}: ${parsed.errors.join('; ')}`)
    parsed.value.external_reference = url
    const metadata = { ...parsed.value }
    await writeTextAtomic(path, `---\n${stringifyYaml(metadata).trimEnd()}\n---\n\n${parsed.body.trimStart()}`)
  }
  const status: PlanStatus = plan.status
  return { plan_id: plan.id, provider: options.provider, published: references, unchanged_status: status, activity_records: [] }
}

export function publicationPayload(workspaceRoot: string, planId: string, provider: string, references: PublicationReference[] = []): string {
  return JSON.stringify({ provider, plan: planId, references, instructions: "Publish the current plan and selected tasks, preserve IDs, and return URLs. Do not update status or create lifecycle records." }, null, 2)
}
