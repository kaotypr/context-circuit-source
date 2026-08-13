import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { configureWorkspace } from "./lib/configure-workspace.js";
import type { WorkspaceBootstrapRequest } from "./lib/types.js";
import { readJsonRegularInside } from "./lib/io.js";

const { values } = parseArgs({ options: { request: { type: "string" }, "check-only": { type: "boolean", default: false } } });
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const request = values.request ? await readJsonRegularInside<WorkspaceBootstrapRequest>(workspaceRoot, resolve(workspaceRoot, values.request), "Workspace configuration request") : undefined;
console.log(JSON.stringify(await configureWorkspace({ workspaceRoot, ...(request ? { request } : {}), checkOnly: values["check-only"] }), null, 2));
