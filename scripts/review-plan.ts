import { dirname, resolve } from "node:path"
import { fileURLToPath } from "node:url"
import { parseArgs } from "node:util"
import { git } from "./lib/git.js"
import { resolvePlanDirectory, validatePlanDirectory } from "./lib/plans.js"

const { values } = parseArgs({ options: { plan: { type: "string" }, worktree: { type: "string" }, base: { type: "string" } } })
if (!values.plan || !values.worktree) throw new Error("Usage: cc review-plan --plan <reference> --worktree <path> [--base <commit>]")
const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
const directory = await resolvePlanDirectory(root, values.plan)
const validation = await validatePlanDirectory(directory)
if (!validation.plan || validation.errors.length > 0) throw new Error(`Cannot review invalid plan:\n- ${validation.errors.join("\n- ")}`)
const worktree = resolve(values.worktree)
const base = values.base ?? await git(worktree, ["merge-base", "HEAD", "HEAD~1"]).catch(() => "HEAD")
console.log(JSON.stringify({
  plan: validation.plan,
  remaining_tasks: validation.tasks.filter((task) => task.status !== "done").map((task) => task.id),
  worktree,
  base,
  branch: await git(worktree, ["branch", "--show-current"]),
  changed_files: (await git(worktree, ["diff", "--name-only", `${base}...HEAD`])).split("\n").filter(Boolean),
  status: await git(worktree, ["status", "--short"]),
  commands: [`git -C ${JSON.stringify(worktree)} diff ${base}...HEAD`, `git -C ${JSON.stringify(worktree)} log --oneline ${base}..HEAD`],
  human_review_required: true,
  review_scope: "whole-plan",
}, null, 2))
