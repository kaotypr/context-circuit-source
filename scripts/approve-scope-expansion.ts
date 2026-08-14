import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { approveScopeExpansion } from "./lib/scope-approval.js";

const { values } = parseArgs({ options: {
  "run-id": { type: "string" },
  repository: { type: "string" },
  "task-id": { type: "string" },
  "approved-by": { type: "string" },
  reason: { type: "string" },
} });

if (!values["run-id"] || !values.repository || !values["approved-by"] || !values.reason) {
  throw new Error("Usage: approve-scope-expansion --run-id <id> --repository <name> [--task-id <id>] --approved-by <identifier> --reason <text>");
}

const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
console.log(JSON.stringify(await approveScopeExpansion({
  workspaceRoot,
  runId: values["run-id"],
  repository: values.repository,
  ...(values["task-id"] ? { taskId: values["task-id"] } : {}),
  approvedBy: values["approved-by"],
  reason: values.reason,
}), null, 2));
