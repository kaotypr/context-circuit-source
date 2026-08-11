import { resolve } from "node:path";
import { parseArgs } from "node:util";
import { validatePlanDirectory } from "./lib/plans.js";

const { positionals } = parseArgs({ allowPositionals: true });
if (!positionals[0]) throw new Error("Usage: cc validate-plan context/plans/<plan-id>");
const planDirectory = resolve(process.cwd(), positionals[0]);
const result = await validatePlanDirectory(planDirectory);
if (result.errors.length > 0) {
  console.error(`Invalid plan ${planDirectory}:`);
  for (const error of result.errors) console.error(`- ${error}`);
  process.exitCode = 1;
} else {
  console.log(`Valid ${result.index?.status} plan ${result.index?.plan_id} version ${result.index?.plan_version}`);
}
