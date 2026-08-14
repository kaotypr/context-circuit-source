import { dirname, resolve } from "node:path"
import { parseArgs } from "node:util"
import { fileURLToPath } from "node:url"
import { preparePlanExecution } from "./lib/run-task.js"

const { values } = parseArgs({ options: { plan: { type: "string" }, session: { type: "string" }, "root-session": { type: "string" } } })
const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
if (values.plan) {
  console.log(JSON.stringify(await preparePlanExecution({ workspaceRoot: root, plan: values.plan, sessionId: values.session, rootSessionId: values["root-session"] }), null, 2))
} else {
  throw new Error("Usage: cc run-task --plan <reference>")
}
