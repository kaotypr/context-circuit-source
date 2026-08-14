import { resolve } from "node:path"
import { parseArgs } from "node:util"
import { validatePlanDirectory } from "./lib/plans.js"

const { positionals } = parseArgs({ allowPositionals: true })
if (!positionals[0]) throw new Error("Usage: cc validate-plan plans/<repository-key>-plans/<number>-<slug>")
const directory = resolve(process.cwd(), positionals[0])
const result = await validatePlanDirectory(directory)
if (result.errors.length > 0) {
  console.error(`Invalid plan ${directory}:`)
  for (const error of result.errors) console.error(`- ${error}`)
  process.exitCode = 1
} else {
  console.log(`Valid ${result.plan!.status} plan ${result.plan!.id} with ${result.tasks.length} task(s)`)
}
