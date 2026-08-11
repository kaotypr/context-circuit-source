import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { prepareActivityLifecycle } from "./lib/activity-lifecycle.js";
import type { ActivityCapability, ActivityEvent } from "./lib/types.js";

const events: ActivityEvent[] = ["task.starting", "task.review-ready", "task.completed", "task.blocked", "task.cancelled"];
const capabilities: ActivityCapability[] = ["read-tasks", "update-status", "create-tasks", "assign-task", "timers"];
const { values } = parseArgs({ options: {
  "run-id": { type: "string" },
  event: { type: "string" },
  available: { type: "string", multiple: true, default: [] },
} });
if (!values["run-id"] || !events.includes(values.event as ActivityEvent) || values.available.some((item) => !capabilities.includes(item as ActivityCapability))) {
  throw new Error("Usage: prepare-lifecycle --run-id <id> --event <semantic-event> [--available <capability>]");
}
const workspaceRoot = resolve(process.env.KAO_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
console.log(JSON.stringify(await prepareActivityLifecycle({
  workspaceRoot,
  runId: values["run-id"],
  event: values.event as ActivityEvent,
  availableCapabilities: values.available as ActivityCapability[],
}), null, 2));
