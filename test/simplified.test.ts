import assert from "node:assert/strict"
import { access, readFile, writeFile } from "node:fs/promises"
import { join } from "node:path"
import { test } from "node:test"
import { archivePlan, createPlan, listPlans, resolvePlanDirectory, setPlanStatus, setTaskStatus, unarchivePlan, validatePlanDirectory } from "../scripts/lib/plans.js"
import { publishPlan } from "../scripts/lib/plan-publication.js"
import { importProductKnowledge, refreshProductKnowledge, staleProductKnowledgeSources } from "../scripts/lib/product-knowledge.js"
import { preparePlanExecution } from "../scripts/lib/run-task.js"
import { recommendWhatsNext } from "../scripts/lib/whats-next.js"
import { createTestWorkspace, sha256 } from "./helpers.js"

test("provides an Idea Brief entry point for empty workspaces", async () => {
  const skill = await readFile(join(process.cwd(), ".agents", "skills", "cc-idea-brief", "SKILL.md"), "utf8")
  assert.match(skill, /context\/IDEA-BRIEF\.md/)
  assert.match(skill, /cc-import-product-knowledge/)
})

test("creates and parses the YAML plan and task Markdown format", async () => {
  const workspace = await createTestWorkspace()
  try {
    const result = await createPlan(workspace.root, {
      id: "checkout",
      title: "Checkout improvements",
      source: { kind: "prd", reference: "docs/checkout.md" },
      repository: "app",
      work_prefix: "CHK",
      requirements: ["A shopper can complete checkout"],
      acceptance_criteria: ["A successful checkout creates an order"],
      tasks: [{ id: "CHK-0001", title: "Implement checkout", repository: "app", implementation_scope: ["README.md"], acceptance_criteria: ["The task is reviewable"] }],
    })
    assert.match(result.directory, /plans\/app-plans\/\d{4}-checkout-improvements$/)
    const validation = await validatePlanDirectory(result.directory)
    assert.deepEqual(validation.errors, [])
    assert.equal(validation.plan?.status, "draft")
    assert.equal(validation.tasks[0]?.status, "draft")
    assert.equal(validation.tasks[0]?.id, "CHK-0001")
    assert.match(await readFile(join(result.directory, "plan.yaml"), "utf8"), /status: draft/)
    assert.match(await readFile(join(result.directory, "tasks", "CHK-0001.md"), "utf8"), /^---\n/)
  } finally {
    await workspace.cleanup()
  }
})

test("supports exactly three statuses and plan-level whats-next", async () => {
  const workspace = await createTestWorkspace()
  try {
    const result = await createPlan(workspace.root, {
      id: "dependencies",
      title: "Dependency example",
      repository: "app",
      tasks: [
        { id: "DEP-0001", title: "First", repository: "app" },
        { id: "DEP-0002", title: "Second", repository: "app", dependencies: ["DEP-0001"] },
      ],
    })
    await setPlanStatus(workspace.root, result.directory, "approved")
    let next = await recommendWhatsNext(workspace.root)
    assert.equal(next.recommendation.plan_id, "dependencies")
    assert.deepEqual(next.recommendation.remaining_tasks, ["DEP-0001", "DEP-0002"])
    assert.equal(next.no_state_changed, true)
    await setTaskStatus(workspace.root, result.directory, "DEP-0001", "done")
    next = await recommendWhatsNext(workspace.root)
    assert.equal(next.recommendation.plan_id, "dependencies")
    assert.deepEqual(next.recommendation.remaining_tasks, ["DEP-0002"])
    await setTaskStatus(workspace.root, result.directory, "DEP-0002", "done")
    await setPlanStatus(workspace.root, result.directory, "done")
    assert.equal((await validatePlanDirectory(result.directory)).plan?.status, "done")
  } finally {
    await workspace.cleanup()
  }
})

test("detects stale Product Knowledge and proposes a refresh", async () => {
  const workspace = await createTestWorkspace()
  try {
    const source = join(workspace.root, "requirements.md")
    const first = "# Requirements\n\nInitial behavior.\n"
    await writeFile(source, first, "utf8")
    await importProductKnowledge(workspace.root, { source: "requirements.md", source_id: "requirements", product_knowledge: ["context/PROJECT.md"] })
    assert.deepEqual(await staleProductKnowledgeSources(workspace.root), [])
    await writeFile(source, `${first}\nChanged behavior.\n`, "utf8")
    const stale = await staleProductKnowledgeSources(workspace.root)
    assert.equal(stale[0]?.id, "requirements")
    const refresh = await refreshProductKnowledge(workspace.root)
    assert.match(refresh.proposal, /requirements/)
    assert.deepEqual(refresh.affected_pages, ["context/PROJECT.md"])
  } finally {
    await workspace.cleanup()
  }
})

test("publication stores URLs without status changes or activity records", async () => {
  const workspace = await createTestWorkspace()
  try {
    const result = await createPlan(workspace.root, { id: "publish", title: "Publish example", repository: "app", tasks: [{ id: "PUB-0001", title: "Task", repository: "app" }] })
    await setPlanStatus(workspace.root, result.directory, "approved")
    const published = await publishPlan({ workspaceRoot: workspace.root, plan: result.directory, provider: "github", references: [{ id: "publish", url: "https://example.test/plans/1" }, { id: "PUB-0001", url: "https://example.test/issues/2" }] })
    assert.equal(published.unchanged_status, "approved")
    assert.deepEqual(published.activity_records, [])
    assert.equal((await validatePlanDirectory(result.directory)).plan?.status, "approved")
    assert.equal(await access(join(workspace.root, ".runtime", "activity")).then(() => true).catch(() => false), false)
  } finally {
    await workspace.cleanup()
  }
})

test("publication requires explicit plan approval before execution", async () => {
  const workspace = await createTestWorkspace()
  try {
    const result = await createPlan(workspace.root, { id: "publish-draft", title: "Draft publication", repository: "app", tasks: [{ id: "PUB-DRAFT", title: "Task", repository: "app" }] })
    await assert.rejects(() => publishPlan({ workspaceRoot: workspace.root, plan: result.directory, provider: "github", references: [] }), /publication is offered after explicit approval and before execution/)
  } finally {
    await workspace.cleanup()
  }
})

test("archives any plan status, ignores it, and explicitly unarchives it", async () => {
  const workspace = await createTestWorkspace()
  try {
    const archivedByStatus = []
    for (const status of ["draft", "approved", "done"] as const) {
      const result = await createPlan(workspace.root, { id: `archive-${status}`, title: `Archive ${status}`, repository: "app", tasks: [{ id: `ARC-${status}`, title: "Task", repository: "app" }] })
      if (status !== "draft") await setPlanStatus(workspace.root, result.directory, status)
      const archived = await archivePlan(workspace.root, result.directory)
      archivedByStatus.push({ ...archived, status })
      assert.equal((await validatePlanDirectory(archived.destination)).plan?.status, status)
    }
    const next = await recommendWhatsNext(workspace.root)
    assert.equal(next.recommendation.action, "review")
    const archived = archivedByStatus[1]!
    assert.match(archived.destination, /archives\/plans\/app-plans\/\d{4}-archive-approved$/)
    const restored = await unarchivePlan(workspace.root, `archives/plans/app-plans/${archived.destination.split('/').at(-1)}`)
    assert.equal(restored.plan_id, "archive-approved")
    assert.equal((await validatePlanDirectory(await resolvePlanDirectory(workspace.root, restored.destination))).plan?.status, "approved")
    assert.equal((await listPlans(workspace.root, true)).filter((item) => item.plan?.id === "archive-approved").length, 1)
  } finally {
    await workspace.cleanup()
  }
})

test("runs an approved plan continuously without changing status or writing result contracts", async () => {
  const workspace = await createTestWorkspace()
  try {
    const result = await createPlan(workspace.root, { id: "execute", title: "Execution example", repository: "app", tasks: [
      { id: "RUN-0001", title: "Run first task", repository: "app" },
      { id: "RUN-0002", title: "Run second task", repository: "app" },
    ] })
    await setPlanStatus(workspace.root, result.directory, "approved")
    const prepared = await preparePlanExecution({ workspaceRoot: workspace.root, plan: result.directory, sessionId: "sess-runner", rootSessionId: "sess-root" })
    assert.equal(prepared.status_changed, false)
    assert.equal(prepared.plan_domain, "app")
    assert.equal(prepared.session_id, "sess-runner")
    assert.equal(prepared.lease_status, "active")
    assert.deepEqual(prepared.tasks.map((task) => task.id), ["RUN-0001", "RUN-0002"])
    assert.equal((await validatePlanDirectory(result.directory)).tasks[0]?.status, "draft")
    assert.match(prepared.prompt, /No result JSON contract is required/)
    assert.match(prepared.prompt, /without pausing for per-task human review/)
    const second = await preparePlanExecution({ workspaceRoot: workspace.root, plan: result.directory, sessionId: "sess-runner", rootSessionId: "sess-root" })
    assert.equal(second.worktree, prepared.worktree)
    assert.equal(second.branch, prepared.branch)
    assert.match(prepared.worktree, /worktrees\/app\/execute$/)
    assert.doesNotMatch(prepared.worktree, /RUN-0001|RUN-0002/)
    await workspace.cleanup()
  } catch (error) {
    await workspace.cleanup()
    throw error
  }
})
