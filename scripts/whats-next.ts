import { readFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { recommendWhatsNext } from "./lib/whats-next.js";
import type { FakeActivitySource } from "./lib/types.js";

const { values } = parseArgs({ options: { "activity-fixture": { type: "string" } } });
const workspaceRoot = resolve(process.env.KAO_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const activity = values["activity-fixture"]
  ? JSON.parse(await readFile(resolve(process.cwd(), values["activity-fixture"]), "utf8")) as FakeActivitySource
  : null;
console.log(JSON.stringify(await recommendWhatsNext(workspaceRoot, activity), null, 2));
