import { mkdir } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { parseArgs } from "node:util";
import { assertInside, writeTextAtomic } from "./lib/io.js";
import { generateOnboardingPack } from "./lib/product-knowledge-onboarding.js";

const { values } = parseArgs({
  options: {
    roles: { type: "string" },
    revision: { type: "string" },
    "generated-at": { type: "string" },
    out: { type: "string" },
  },
});

const workspaceRoot = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const roles = (values.roles ?? "").split(",").map((role) => role.trim()).filter(Boolean);
if (roles.length === 0) throw new Error("Usage: onboarding-pack --roles <role[,role...]> [--revision <rev>] [--generated-at <iso>] [--out <path>]");

const pack = await generateOnboardingPack({
  workspaceRoot,
  roles,
  revision: values.revision ?? "working-tree",
  generated_at: values["generated-at"] ?? new Date().toISOString(),
});

if (values.out) {
  const out = assertInside(workspaceRoot, resolve(workspaceRoot, values.out));
  await mkdir(dirname(out), { recursive: true });
  await writeTextAtomic(out, pack.markdown);
  console.log(JSON.stringify({ ...pack.manifest, out: values.out }, null, 2));
} else {
  console.log(JSON.stringify(pack.manifest, null, 2));
}
