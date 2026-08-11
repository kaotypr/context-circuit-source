import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { preparePlanlessTask } from "./lib/run-task.js";
import type { TestExpectationPolicy } from "./lib/types.js";

const testPolicies: TestExpectationPolicy[] = ["required", "existing-coverage", "verifier-only", "not-required"];
const workspaceRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");

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
  },
});

if (!values.request || !values.repository || (values["test-policy"] && !testPolicies.includes(values["test-policy"] as TestExpectationPolicy))) {
  throw new Error("Usage: run-task --request <text> --repository <name> --acceptance <criterion> --scope <path> [--test-scope <path>] [--test-policy <policy>] [--verify <command>]");
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
});

console.log(JSON.stringify(prepared, null, 2));
console.warn("Warning: no activity claim was attempted; duplicate effort is possible.");
