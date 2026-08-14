import { readFile } from "node:fs/promises"
import { dirname, resolve } from "node:path"
import { parseArgs } from "node:util"
import { fileURLToPath } from "node:url"
import { publishPlan } from "./lib/plan-publication.js"
import type { PublicationReference } from "./lib/types.js"

const { values } = parseArgs({ options: { plan: { type: "string" }, provider: { type: "string", default: "manual" }, references: { type: "string" } } })
if (!values.plan) throw new Error("Usage: cc publish-plan --plan <reference> [--provider <name>] [--references <json-file>]")
const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
const references = values.references ? JSON.parse(await readFile(resolve(process.cwd(), values.references), "utf8")) as PublicationReference[] : []
console.log(JSON.stringify(await publishPlan({ workspaceRoot: root, plan: values.plan, provider: values.provider ?? "manual", references }), null, 2))
