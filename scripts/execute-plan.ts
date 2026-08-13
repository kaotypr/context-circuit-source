import { dirname, resolve } from "node:path";
import { readFile } from "node:fs/promises";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { prepareExecutePlan } from "./lib/execute-plan.js";
import type { ActivityCapability, PlanExecutionRequest } from "./lib/types.js";

const { values } = parseArgs({
  options: {
    "request-file": { type: "string" },
    plan: { type: "string" },
    version: { type: "string" },
    "approved-digest": { type: "string" },
    available: { type: "string", multiple: true, default: [] },
  },
});
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const capabilities: ActivityCapability[] = ["read-tasks", "update-status", "create-tasks", "assign-task", "timers"];
if (values.available.some((capability) => !capabilities.includes(capability as ActivityCapability))) throw new Error(`Unknown activity capability; expected one of: ${capabilities.join(", ")}`);
let request: PlanExecutionRequest;
if (values["request-file"]) {
  request = JSON.parse(await readFile(resolve(process.cwd(), values["request-file"]), "utf8")) as PlanExecutionRequest;
} else if (values.plan && values.version && values["approved-digest"]) {
  request = { contract_version: 1, source: { kind: "plan", reference: values.plan, plan_version: Number(values.version), approved_digest: values["approved-digest"] } };
} else {
  throw new Error("Usage: execute-plan --request-file <json> | execute-plan --plan <reference> --version <number> --approved-digest <sha256:...>");
}
console.log(JSON.stringify(await prepareExecutePlan({ workspaceRoot, request, availableCapabilities: values.available as ActivityCapability[] }), null, 2));
