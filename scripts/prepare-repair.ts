import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { prepareRepair } from "./lib/review-lifecycle.js";

const { values } = parseArgs({
  options: {
    "run-id": { type: "string" },
    repository: { type: "string" },
  },
});
if (!values["run-id"] || !values.repository) throw new Error("Usage: prepare-repair --run-id <id> --repository <name>");

const workspaceRoot = resolve(process.env.KAO_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
console.log(JSON.stringify(await prepareRepair({ workspaceRoot, runId: values["run-id"], repository: values.repository }), null, 2));
