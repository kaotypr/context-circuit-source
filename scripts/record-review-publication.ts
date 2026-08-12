import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { recordReviewPublication } from "./lib/review-lifecycle.js";

const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const { values } = parseArgs({ options: {
  "run-id": { type: "string" }, repository: { type: "string" }, status: { type: "string" }, tool: { type: "string" },
  "pull-request": { type: "string" }, evidence: { type: "string" }, authorized: { type: "boolean" },
} });
if (!values["run-id"] || !values.repository || !values.evidence || !["published", "failed"].includes(values.status ?? "") || !["gh", "glab", "manual"].includes(values.tool ?? "")) {
  throw new Error("Usage: cc record-review-publication --run-id <id> --repository <name> --status <published|failed> --tool <gh|glab|manual> [--pull-request <ref>] --evidence <text> --authorized");
}
console.log(JSON.stringify(await recordReviewPublication({
  workspaceRoot, runId: values["run-id"], repository: values.repository,
  status: values.status as "published" | "failed", tool: values.tool as "gh" | "glab" | "manual",
  ...(values["pull-request"] ? { pullRequest: values["pull-request"] } : {}), evidence: values.evidence,
  ...(values.authorized ? { authorized: true } : {}),
}), null, 2));
