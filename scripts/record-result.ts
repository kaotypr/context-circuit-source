import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { recordResult, type RecordStage } from "./lib/record-result.js";

const stages: RecordStage[] = ["worker-started", "worker-result", "verifier-result"];
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const { values } = parseArgs({
  options: {
    "run-id": { type: "string" },
    repository: { type: "string" },
    stage: { type: "string" },
  },
});

if (!values["run-id"] || !values.repository || !values.stage || !stages.includes(values.stage as RecordStage)) {
  throw new Error("Usage: record-result --run-id <id> --repository <name> --stage <worker-started|worker-result|verifier-result>");
}

const manifest = await recordResult({
  workspaceRoot,
  runId: values["run-id"],
  repository: values.repository,
  stage: values.stage as RecordStage,
});
console.log(JSON.stringify({ runId: manifest.run_id, status: manifest.status, executionEvents: manifest.execution_events?.length ?? 0 }, null, 2));
