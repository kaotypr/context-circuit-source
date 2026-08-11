import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { bootstrapWorkspace, initializeWorkspace } from "./lib/initialize-workspace.js";
import type { WorkspaceBootstrapRequest } from "./lib/types.js";
import { readData } from "./lib/validation.js";

const { values } = parseArgs({
  options: {
    "check-only": { type: "boolean", default: false },
    bootstrap: { type: "string" },
  },
});

const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
if (values.bootstrap && values["check-only"]) throw new Error("--bootstrap and --check-only cannot be combined");
const summary = values.bootstrap
  ? await bootstrapWorkspace({ workspaceRoot, request: await readData(resolve(values.bootstrap)) as WorkspaceBootstrapRequest })
  : await initializeWorkspace({ workspaceRoot, apply: !values["check-only"] });
console.log(JSON.stringify(summary, null, 2));
