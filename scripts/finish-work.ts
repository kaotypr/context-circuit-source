import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { finishWork } from "./lib/finish-work.js";

const { values } = parseArgs({
  options: {
    "run-id": { type: "string" },
    repository: { type: "string" },
    outcome: { type: "string" },
    author: { type: "string" },
    reason: { type: "string" },
    "merge-commit": { type: "string" },
    "pull-request": { type: "string", multiple: true },
    cleanup: { type: "boolean", default: false },
  },
});

if (!values["run-id"] || !values.repository || !values.outcome || !values.author) {
  throw new Error("Required: --run-id <id> --repository <name> --outcome <merged|abandoned> --author <slug>");
}
if (values.outcome !== "merged" && values.outcome !== "abandoned") throw new Error("--outcome must be merged or abandoned");

const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const result = await finishWork({
  workspaceRoot,
  runId: values["run-id"],
  repository: values.repository,
  outcome: values.outcome,
  author: values.author,
  cleanup: values.cleanup,
  ...(values.reason ? { reason: values.reason } : {}),
  ...(values["merge-commit"] ? { mergeCommit: values["merge-commit"] } : {}),
  ...(values["pull-request"] ? { pullRequests: values["pull-request"] } : {}),
});
console.log(JSON.stringify(result, null, 2));
