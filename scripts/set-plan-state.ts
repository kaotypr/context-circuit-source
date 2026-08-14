import { dirname, resolve } from "node:path"
import { parseArgs } from "node:util"
import { fileURLToPath } from "node:url"
import { setPlanStatus } from "./lib/plans.js"
import type { PlanStatus } from "./lib/types.js"

const { values } = parseArgs({ options: { plan: { type: "string" }, status: { type: "string" } } })
if (!values.plan || !values.status || !["draft", "approved", "done"].includes(values.status)) throw new Error("Usage: cc set-plan-state --plan <reference> --status <draft|approved|done>")
const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
console.log(JSON.stringify(await setPlanStatus(root, values.plan, values.status as PlanStatus), null, 2))
