import { readFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { preparePlanPublication } from "./lib/plan-publication.js";
import type { PlanPublicationDiscovery } from "./lib/types.js";

const { values } = parseArgs({ options: { plan: { type: "string" }, discovery: { type: "string" } } });
if (!values.plan || !values.discovery) throw new Error("Usage: prepare-plan-publication --plan <plan-id> --discovery <json>");
const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const discovery = JSON.parse(await readFile(resolve(process.cwd(), values.discovery), "utf8")) as PlanPublicationDiscovery;
console.log(JSON.stringify(await preparePlanPublication({ workspaceRoot, planId: values.plan, discovery }), null, 2));
