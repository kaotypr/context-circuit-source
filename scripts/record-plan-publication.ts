import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { recordPlanPublication } from "./lib/plan-publication.js";

const { values } = parseArgs({ options: { plan: { type: "string" }, work: { type: "string" }, status: { type: "string" }, evidence: { type: "string" }, reference: { type: "string" } } });
if (!values.plan || !values.work || !["created", "failed"].includes(values.status ?? "") || !values.evidence) throw new Error("Usage: record-plan-publication --plan <id> --work <id> --status <created|failed> --evidence <text> [--reference <external-ref>]");
const workspaceRoot = resolve(process.env.KAO_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
console.log(JSON.stringify(await recordPlanPublication({ workspaceRoot, planId: values.plan, workId: values.work, status: values.status as "created" | "failed", evidence: values.evidence, ...(values.reference ? { externalReference: values.reference } : {}) }), null, 2));
