import { readFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { createPlanDraft, generatePlanBatch, migrateCurrentPlans } from "./lib/plans.js";
import type { PlanDraftRequest, PlanGenerationRequest } from "./lib/types.js";

const { values } = parseArgs({
  options: { input: { type: "string" }, migrate: { type: "boolean", default: false } },
});
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
if (values.migrate) {
  console.log(JSON.stringify(await migrateCurrentPlans(workspaceRoot), null, 2));
  process.exit(0);
}
if (!values.input) throw new Error("Usage: cc create-plan --input <plan-generation-request.json> | cc create-plan --migrate");
const inputPath = resolve(process.cwd(), values.input);
const request = JSON.parse(await readFile(inputPath, "utf8")) as PlanDraftRequest | PlanGenerationRequest;
if (request.contract_version === 2) console.log(JSON.stringify(await generatePlanBatch(workspaceRoot, request), null, 2));
else console.log(JSON.stringify(await createPlanDraft(workspaceRoot, request), null, 2));
