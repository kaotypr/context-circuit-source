import { readFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { createPlanDraft } from "./lib/plans.js";
import type { PlanDraftRequest } from "./lib/types.js";

const { values } = parseArgs({
  options: { input: { type: "string" } },
});
if (!values.input) throw new Error("Usage: cc create-plan --input <plan-draft-request.json>");
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const inputPath = resolve(process.cwd(), values.input);
const request = JSON.parse(await readFile(inputPath, "utf8")) as PlanDraftRequest;
console.log(JSON.stringify(await createPlanDraft(workspaceRoot, request), null, 2));
