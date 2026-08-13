import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { readJsonRegularInside } from "./lib/io.js";
import { prepareImportContext } from "./lib/import-context.js";
import type { ImportContextRequest } from "./lib/types.js";

const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const { values } = parseArgs({ options: { request: { type: "string" } } });
if (!values.request) throw new Error("Usage: import-context --request <import-context-request.json>");

const request = await readJsonRegularInside<ImportContextRequest>(
  workspaceRoot,
  resolve(workspaceRoot, values.request),
  "Import context request",
);
const prepared = await prepareImportContext({ workspaceRoot, request });

console.log(JSON.stringify({
  ...prepared,
  handoff: {
    context_sync_request: {
      contribution: prepared.contribution,
      guidance: "Curate the discovered evidence into a context-sync-request JSON; preserve source citations and explicit unknowns.",
      required_command: ["node", ".agents/bin/cc.mjs", "sync-context", "--request", "<context-sync-request.json>"],
    },
    context_review: {
      guidance: "After sync-context returns a sync_id and the curated changes are committed, prepare the review handoff.",
      required_command: ["node", ".agents/bin/cc.mjs", "prepare-context-review", "--sync-id", "<sync-id>"],
    },
  },
}, null, 2));
