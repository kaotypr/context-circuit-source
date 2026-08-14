import { access, mkdir, readFile, realpath } from "node:fs/promises"
import { join, resolve } from "node:path"
import { parse as parseYaml } from "yaml"
import { assertCleanRepository, git } from "./git.js"
import { assertInside, writeTextAtomic } from "./io.js"
import { resolvePlanDirectory, validatePlanDirectory } from "./plans.js"
import type { PlanTask, PlanYaml, PreparedPlanExecution, RepositoryConfig, WorkspaceConfig } from "./types.js"

function slugify(value: string): string {
  return value.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "").slice(0, 48) || "plan"
}

async function exists(path: string): Promise<boolean> {
  try {
    await access(path)
    return true
  } catch {
    return false
  }
}

async function loadWorkspace(root: string): Promise<WorkspaceConfig> {
  const value = parseYaml(await readFile(join(root, "workspace.yaml"), "utf8")) as WorkspaceConfig
  if (!value?.repositories) throw new Error("workspace.yaml has no repositories")
  return value
}

function planDomain(plan: PlanYaml, tasks: PlanTask[]): string {
  const domains = new Set([...(plan.repositories ?? []), ...tasks.map((task) => task.repository)])
  if (domains.size !== 1) throw new Error(`Plan ${plan.id} spans multiple repository domains (${[...domains].join(", ")}); split it into one plan per domain before execution`)
  return [...domains][0]!
}

function promptFor(plan: PlanYaml, tasks: PlanTask[], repositoryPath: string): string {
  const lines = [
    `# Plan ${plan.id}: ${plan.title}`,
    "",
    `Repository domain: ${plan.repositories[0]} (${repositoryPath})`,
    "",
    "Work under the assigned plan worktree only. Continue through every unfinished task in dependency order without pausing for per-task human review. Publication, if desired, must happen before this run; do not publish during execution. Human status remains authoritative: do not mark tasks or the plan done, merge, or create lifecycle records.",
    "",
    "## Plan objective",
    "",
    plan.title,
    "",
    "## Tasks",
    "",
  ]
  for (const task of tasks) {
    lines.push(
      `### ${task.id}: ${task.title}`,
      "",
      `Status: ${task.status}`,
      ...(task.dependencies?.length ? [`Dependencies: ${task.dependencies.join(", ")}`] : ["Dependencies: none"]),
      "",
      task.description ?? task.title,
      "",
      "Implementation scope:",
      ...(task.implementation_scope?.length ? task.implementation_scope.map((item) => `- ${item}`) : ["- Use the task description and repository instructions to determine the smallest implementation scope."]),
      "",
      "Test scope and expectations:",
      ...(task.test_scope?.length ? task.test_scope.map((item) => `- ${item}`) : ["- Run the relevant repository checks and report what was run."]),
      ...(task.test_expectations?.length ? task.test_expectations.map((item) => `- ${item}`) : []),
      "",
      "Verification commands:",
      ...(task.verification_commands?.length ? task.verification_commands.map((item) => `- ${item}`) : ["- Human review decides which additional checks are appropriate."]),
      "",
      "Acceptance criteria:",
      ...(task.acceptance_criteria?.length ? task.acceptance_criteria.map((item) => `- ${item}`) : ["- Explain how the requested outcome was handled."]),
      "",
    )
  }
  lines.push(
    "## Agent loop",
    "",
    "1. Read repository instructions, Product Knowledge references, and the plan documents.",
    "2. Select the next unfinished task whose dependencies are satisfied.",
    "3. Implement and test it in this same plan worktree.",
    "4. Continue to the next task without asking for human review between tasks.",
    "5. Stop only when all unfinished tasks are implemented, or explain the blocker and any independent work that remains.",
    "",
    "## Final handoff",
    "",
    "Return one plan-level human-readable summary with changed files, tests run, questions, blockers, and limitations. A human reviews the plan once after this run. No result JSON contract is required.",
    "",
  )
  return lines.join("\n")
}

interface WorktreeRecord {
  path: string
  branch?: string
}

async function listWorktrees(repositoryRoot: string): Promise<WorktreeRecord[]> {
  const output = await git(repositoryRoot, ["worktree", "list", "--porcelain"])
  const records: WorktreeRecord[] = []
  let current: WorktreeRecord | null = null
  for (const line of output.split("\n")) {
    if (line.startsWith("worktree ")) {
      if (current) records.push(current)
      current = { path: line.slice("worktree ".length) }
    } else if (line.startsWith("branch ") && current) {
      current.branch = line.slice("branch ".length).replace(/^refs\/heads\//, "")
    }
  }
  if (current) records.push(current)
  return records
}

async function canonicalPath(path: string): Promise<string> {
  return realpath(path).catch(() => resolve(path))
}

async function planWorktree(root: string, repository: string, repositoryConfig: RepositoryConfig, plan: PlanYaml): Promise<{ worktree: string; branch: string; baseCommit: string }> {
  const repositoryRoot = assertInside(root, resolve(root, repositoryConfig.path))
  const domain = slugify(repository)
  const planKey = slugify(plan.id)
  const branch = `agent/${domain}/plan-${planKey}`
  const worktree = assertInside(root, join(root, ".runtime", "worktrees", domain, planKey))
  const records = await listWorktrees(repositoryRoot)
  const worktreePath = await canonicalPath(worktree)
  let existing: WorktreeRecord | undefined
  for (const record of records) {
    if (await canonicalPath(record.path) === worktreePath) {
      existing = record
      break
    }
  }
  const branchOwner = records.find((record) => record.branch === branch)

  if (existing) {
    if (existing.branch !== branch) throw new Error(`Plan worktree path is already assigned to another branch: ${worktree}`)
    const baseCommit = await git(worktree, ["merge-base", "HEAD", `refs/heads/${repositoryConfig.default_branch}`]).catch(() => git(worktree, ["rev-parse", "HEAD"]))
    return { worktree, branch, baseCommit }
  }
  if (branchOwner) throw new Error(`Plan branch ${branch} is already attached to another worktree: ${branchOwner.path}`)
  if (await exists(worktree)) throw new Error(`Plan worktree path already exists but is not registered by Git: ${worktree}`)

  await assertCleanRepository(repositoryRoot)
  const baseCommit = await git(repositoryRoot, ["rev-parse", "HEAD"])
  await mkdir(resolve(worktree, ".."), { recursive: true })
  await git(repositoryRoot, ["worktree", "add", "-b", branch, worktree, baseCommit])
  return { worktree, branch, baseCommit }
}

async function createPlanWorktree(root: string, repository: string, repositoryConfig: RepositoryConfig, plan: PlanYaml, tasks: PlanTask[]): Promise<PreparedPlanExecution> {
  const { worktree, branch, baseCommit } = await planWorktree(root, repository, repositoryConfig, plan)
  const stamp = new Date().toISOString().replace(/[-:.TZ]/g, "").slice(0, 14)
  const prompt = promptFor(plan, tasks, repositoryConfig.path)
  const promptFile = assertInside(root, join(root, ".runtime", "plans", `${slugify(plan.id)}-${stamp}.md`))
  await mkdir(resolve(promptFile, ".."), { recursive: true })
  await writeTextAtomic(promptFile, prompt)
  const shellPath = (path: string): string => `'${path.replaceAll("'", "'\\''")}'`
  const verificationCommands = [...new Set(tasks.flatMap((task) => task.verification_commands ?? []))]
  return {
    plan,
    plan_domain: repository,
    tasks,
    worktree,
    branch,
    base_commit: baseCommit,
    prompt,
    prompt_file: promptFile,
    review_commands: [
      `git -C ${shellPath(worktree)} diff ${baseCommit}...HEAD`,
      `git -C ${shellPath(worktree)} status --short`,
      ...verificationCommands,
    ],
    status_changed: false,
  }
}

export async function preparePlanExecution(options: { workspaceRoot: string; plan: string }): Promise<PreparedPlanExecution> {
  const root = resolve(options.workspaceRoot)
  const directory = await resolvePlanDirectory(root, options.plan)
  const validation = await validatePlanDirectory(directory)
  if (!validation.plan || validation.errors.length > 0) throw new Error(`Cannot execute invalid plan:\n- ${validation.errors.join("\n- ")}`)
  const plan = validation.plan
  if (plan.status !== "approved") throw new Error(`Plan ${plan.id} is ${plan.status}; only explicitly approved plans are executable`)
  const tasks = validation.tasks.filter((task) => task.status !== "done")
  if (tasks.length === 0) throw new Error(`Plan ${plan.id} has no unfinished tasks`)
  const domain = planDomain(plan, validation.tasks)
  const config = await loadWorkspace(root)
  const repository = config.repositories[domain]
  if (!repository) throw new Error(`Plan repository is not registered: ${domain}`)
  return createPlanWorktree(root, domain, repository, plan, tasks)
}
