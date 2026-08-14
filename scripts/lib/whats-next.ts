import { resolve } from "node:path"
import { listPlans, planReference } from "./plans.js"
import { staleProductKnowledgeSources } from "./product-knowledge.js"
import type { NextPlanRecommendation, PlanListItem, WhatsNextResult } from "./types.js"

function planDone(plans: PlanListItem[], reference: string): boolean {
  return plans.some((candidate) => !candidate.archived && candidate.plan?.id === reference && candidate.plan.status === "done")
}

function planReady(planResult: PlanListItem, plans: PlanListItem[]): boolean {
  if (!planResult.plan || planResult.plan.status !== "approved") return false
  if ((planResult.plan.dependencies ?? []).some((dependency) => !planDone(plans, dependency))) return false
  return planResult.tasks.some((task) => task.status !== "done")
}

function recommendation(planResult: PlanListItem, root: string): NextPlanRecommendation {
  const plan = planResult.plan!
  const remainingTasks = planResult.tasks.filter((task) => task.status !== "done").map((task) => task.id)
  return {
    action: "execute",
    plan_id: plan.id,
    title: plan.title,
    why: `Plan ${plan.id} is approved, its plan dependencies are done, and it has unfinished tasks. Execute the plan continuously in its isolated domain worktree, then request one human review.`,
    plan_reference: planReference(root, planResult.directory),
    product_knowledge_references: [...(plan.product_knowledge?.references ?? []), ...planResult.tasks.flatMap((task) => task.product_knowledge ?? [])].filter((value, index, values) => values.indexOf(value) === index),
    repository: plan.repositories[0],
    remaining_tasks: remainingTasks,
  }
}

export async function recommendWhatsNext(workspaceRootInput: string): Promise<WhatsNextResult> {
  const root = resolve(workspaceRootInput)
  const active = await listPlans(root, false)
  const warnings = active.flatMap((item) => item.errors.map((error) => `${item.directory}: ${error}`))
  const stale = await staleProductKnowledgeSources(root)
  if (stale.length > 0) warnings.push(`Product Knowledge sources changed: ${stale.map((source) => source.id).join(', ')}. Review refresh-product-knowledge before relying on affected pages.`)
  const approved = active.filter((item) => item.plan?.status === "approved")
  const ready: PlanListItem[] = []
  let tasks = 0
  for (const item of approved) {
    tasks += item.tasks.length
    if (planReady(item, active)) ready.push(item)
  }
  ready.sort((left, right) => `${left.plan!.number}-${left.plan!.id}`.localeCompare(`${right.plan!.number}-${right.plan!.id}`))
  if (ready.length > 0) {
    const first = recommendation(ready[0]!, root)
    return {
      generated_at: new Date().toISOString(),
      recommendation: first,
      alternatives: ready.slice(1, 6).map((plan) => recommendation(plan, root)),
      considered: { plans: active.length, tasks, executable: ready.length, blocked: Math.max(0, tasks - ready.length) },
      warnings,
      no_state_changed: true,
    }
  }
  const draft = active.filter((item) => item.plan?.status === "draft").sort((left, right) => (left.plan?.number ?? 0) - (right.plan?.number ?? 0))[0]
  const enabling: NextPlanRecommendation = draft?.plan
    ? {
        action: "review",
        plan_id: draft.plan.id,
        title: `Review draft plan ${draft.plan.id}`,
        why: "No approved plan is ready for continuous execution. Review and explicitly approve a draft plan, or create the next plan.",
        plan_reference: planReference(root, draft.directory),
        product_knowledge_references: draft.plan.product_knowledge?.references ?? [],
      }
    : {
        action: "review",
        title: "Create or approve a plan",
        why: "No approved plan is ready for continuous execution. Create a draft plan from the relevant source, or explicitly approve an existing draft.",
        product_knowledge_references: [],
      }
  return {
    generated_at: new Date().toISOString(),
    recommendation: enabling,
    alternatives: [],
    considered: { plans: active.length, tasks, executable: 0, blocked: tasks },
    warnings,
    no_state_changed: true,
  }
}
