import { dirname, resolve } from "node:path"
import { parseArgs } from "node:util"
import { fileURLToPath } from "node:url"
import { setTaskStatus } from "./lib/plans.js"
import type { TaskStatus } from "./lib/types.js"

const { values } = parseArgs({ options: { plan: { type: "string" }, task: { type: "string" }, status: { type: "string" } } })
if (!values.plan || !values.task || !values.status || !["draft", "approved", "done"].includes(values.status)) throw new Error("Usage: cc set-task-state --plan <reference> --task <id> --status <draft|approved|done>")
const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
console.log(JSON.stringify(await setTaskStatus(root, values.plan, values.task, values.status as TaskStatus), null, 2))
