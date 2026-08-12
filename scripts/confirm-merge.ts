import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { confirmMerge } from "./lib/review-lifecycle.js";

const { values } = parseArgs({ options: {
  "run-id": { type: "string" }, repository: { type: "string" }, "merge-commit": { type: "string" },
  evidence: { type: "string" }, author: { type: "string" },
} });
if (!values["run-id"] || !values.repository || !values["merge-commit"] || !values.evidence || !values.author) {
  throw new Error("Usage: cc confirm-merge --run-id <id> --repository <name> --merge-commit <full-sha> --author <slug> --evidence <single-line-evidence>");
}
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
console.log(JSON.stringify(await confirmMerge({
  workspaceRoot, runId: values["run-id"], repository: values.repository, mergeCommit: values["merge-commit"], evidence: values.evidence, author: values.author,
}), null, 2));
