import { readFile } from "node:fs/promises"
import { dirname, resolve } from "node:path"
import { parseArgs } from "node:util"
import { fileURLToPath } from "node:url"
import { createPlan } from "./lib/plans.js"
import type { PlanCreateRequest } from "./lib/types.js"

const { values } = parseArgs({ options: { input: { type: "string" } } })
if (!values.input) throw new Error("Usage: cc create-plan --input <plan-request.json>")
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
const input = JSON.parse(await readFile(resolve(process.cwd(), values.input), "utf8")) as PlanCreateRequest & { plans?: PlanCreateRequest[] }
const requests = input.plans?.length ? input.plans : [input]
const results = []
for (const request of requests) results.push(await createPlan(workspaceRoot, request))
console.log(JSON.stringify(results.length === 1 ? results[0] : results, null, 2))
