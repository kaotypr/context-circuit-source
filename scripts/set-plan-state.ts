import { dirname, join, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { assertInside } from "./lib/io.js";
import { setPlanState, type PlanStateTransition } from "./lib/plans.js";

const { values } = parseArgs({
  options: {
    plan: { type: "string" },
    "approve-by": { type: "string" },
    "material-revision": { type: "string" },
    "non-material-repair": { type: "boolean", default: false },
  },
});
if (!values.plan) throw new Error("Usage: kao set-plan-state --plan context/plans/<plan-id> (--approve-by <name> | --material-revision <reason> | --non-material-repair)");
const requested = [Boolean(values["approve-by"]), Boolean(values["material-revision"]), values["non-material-repair"]].filter(Boolean).length;
if (requested !== 1) throw new Error("Choose exactly one plan state transition");
let transition: PlanStateTransition;
if (values["approve-by"]) transition = { kind: "approve", approved_by: values["approve-by"] };
else if (values["material-revision"]) transition = { kind: "material-revision", reason: values["material-revision"] };
else transition = { kind: "non-material-repair" };
const workspaceRoot = resolve(process.env.KAO_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const planDirectory = assertInside(join(workspaceRoot, "context", "plans"), resolve(process.cwd(), values.plan));
console.log(JSON.stringify(await setPlanState(planDirectory, transition), null, 2));
