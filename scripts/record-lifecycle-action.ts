import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { recordActivityLifecycleAction } from "./lib/activity-lifecycle.js";
import type { ActivityEvent } from "./lib/types.js";

const events: ActivityEvent[] = ["task.starting", "task.review-ready", "task.completed", "task.blocked", "task.cancelled"];
const { values } = parseArgs({ options: {
  "run-id": { type: "string" }, event: { type: "string" }, action: { type: "string" },
  status: { type: "string" }, evidence: { type: "string" }, reference: { type: "string" },
} });
if (!values["run-id"] || !events.includes(values.event as ActivityEvent) || !values.action || !["completed", "failed"].includes(values.status ?? "") || !values.evidence) {
  throw new Error("Usage: record-lifecycle-action --run-id <id> --event <event> --action <id> --status <completed|failed> --evidence <text> [--reference <ref>]");
}
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
console.log(JSON.stringify(await recordActivityLifecycleAction({
  workspaceRoot, runId: values["run-id"], event: values.event as ActivityEvent, actionId: values.action,
  status: values.status as "completed" | "failed", evidence: values.evidence,
  ...(values.reference ? { externalReference: values.reference } : {}),
}), null, 2));
