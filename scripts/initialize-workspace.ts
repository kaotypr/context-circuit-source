import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { initializeWorkspace } from "./lib/initialize-workspace.js";

const { values } = parseArgs({
  options: {
    "check-only": { type: "boolean", default: false },
  },
});

const workspaceRoot = resolve(process.env.KAO_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const summary = await initializeWorkspace({
  workspaceRoot,
  apply: !values["check-only"],
});
console.log(JSON.stringify(summary, null, 2));
