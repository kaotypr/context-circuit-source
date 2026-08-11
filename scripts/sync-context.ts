import { readFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { prepareContextSync } from "./lib/context-sync.js";
import type { ContextSyncRequest } from "./lib/types.js";

const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const { values } = parseArgs({ options: { request: { type: "string" } } });
if (!values.request) throw new Error("Usage: sync-context --request <context-sync-request.json>");
const request = JSON.parse(await readFile(resolve(values.request), "utf8")) as ContextSyncRequest;
console.log(JSON.stringify(await prepareContextSync({ workspaceRoot, request }), null, 2));
