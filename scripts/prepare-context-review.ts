import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { prepareContextReview } from "./lib/context-sync.js";

const workspaceRoot = resolve(process.env.KAO_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const { values } = parseArgs({ options: { "sync-id": { type: "string" } } });
if (!values["sync-id"]) throw new Error("Usage: prepare-context-review --sync-id <id>");
console.log(JSON.stringify(await prepareContextReview({ workspaceRoot, syncId: values["sync-id"] }), null, 2));
