import { dirname, resolve } from "node:path"
import { parseArgs } from "node:util"
import { fileURLToPath } from "node:url"
import { unarchivePlan } from "./lib/plans.js"

const { positionals } = parseArgs({ allowPositionals: true })
if (!positionals[0]) throw new Error("Usage: cc unarchive-plan <archived-plan-reference>")
const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
console.log(JSON.stringify(await unarchivePlan(root, positionals[0]), null, 2))
