import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { configureWorkspace } from "./lib/configure-workspace.js";
import type { WorkspaceBootstrapRequest } from "./lib/types.js";
import { readData } from "./lib/validation.js";

const { values } = parseArgs({ options: { request: { type: "string" }, "check-only": { type: "boolean", default: false } } });
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const request = values.request ? await readData(resolve(values.request)) as WorkspaceBootstrapRequest : undefined;
console.log(JSON.stringify(await configureWorkspace({ workspaceRoot, ...(request ? { request } : {}), checkOnly: values["check-only"] }), null, 2));
