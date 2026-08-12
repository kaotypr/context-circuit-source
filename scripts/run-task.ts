import { dirname, resolve } from "node:path";
import { readFile } from "node:fs/promises";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { prepareContractFirstTask, preparePlanTask, preparePlanlessTask, resumePlanlessTask } from "./lib/run-task.js";
import type { ActivityCapability, PlanRunTaskRequest, RunTaskRequest, StructuredRunTaskRequest, TestExpectationPolicy } from "./lib/types.js";

const testPolicies: TestExpectationPolicy[] = ["required", "existing-coverage", "verifier-only", "not-required"];
const activityCapabilities: ActivityCapability[] = ["read-tasks", "update-status", "create-tasks", "assign-task", "timers"];
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));

const { values } = parseArgs({
  options: {
    request: { type: "string" },
    repository: { type: "string" },
    acceptance: { type: "string", multiple: true, default: [] },
    scope: { type: "string", multiple: true, default: [] },
    "test-scope": { type: "string", multiple: true, default: [] },
    "test-policy": { type: "string" },
    "test-rationale": { type: "string" },
    verify: { type: "string", multiple: true, default: [] },
    available: { type: "string", multiple: true, default: [] },
    "resume-run": { type: "string" },
    "request-file": { type: "string" },
  },
});

if (values["resume-run"]) {
  console.log(JSON.stringify(await resumePlanlessTask({ workspaceRoot, runId: values["resume-run"] }), null, 2));
  process.exit(0);
}

if (values.available.some((capability) => !activityCapabilities.includes(capability as ActivityCapability))) {
  throw new Error(`Unknown activity capability; expected one of: ${activityCapabilities.join(", ")}`);
}

if (values["request-file"]) {
  const request = JSON.parse(await readFile(resolve(values["request-file"]), "utf8")) as StructuredRunTaskRequest;
  const prepared = "source" in request && request.source?.kind === "plan"
    ? await preparePlanTask({ workspaceRoot, request: request as PlanRunTaskRequest, availableCapabilities: values.available as ActivityCapability[] })
    : await prepareContractFirstTask({ workspaceRoot, request: request as RunTaskRequest, availableCapabilities: values.available as ActivityCapability[] });
  console.log(JSON.stringify(prepared, null, 2));
  process.exit(0);
}

if (!values.request || !values.repository || (values["test-policy"] && !testPolicies.includes(values["test-policy"] as TestExpectationPolicy))) {
  throw new Error("Usage: run-task --request <text> --repository <name> --acceptance <criterion> --scope <path> [...] | run-task --request-file <json>");
}
const prepared = await preparePlanlessTask({
  workspaceRoot,
  request: values.request,
  repository: values.repository,
  acceptanceCriteria: values.acceptance,
  scope: values.scope,
  testScope: values["test-scope"],
  ...(values["test-policy"] ? { testPolicy: values["test-policy"] as TestExpectationPolicy } : {}),
  ...(values["test-rationale"] ? { testRationale: values["test-rationale"] } : {}),
  verificationCommands: values.verify,
  availableCapabilities: values.available as ActivityCapability[],
});

console.log(JSON.stringify(prepared, null, 2));
if (prepared.preparationStatus === "prepared") console.warn("Warning: exclusive ownership is guaranteed only when a configured starting action confirmed it.");
else console.warn(`Worktree not created: activity preflight is ${prepared.preparationStatus}. Complete the recorded actions, then rerun with --resume-run ${prepared.runId}.`);
