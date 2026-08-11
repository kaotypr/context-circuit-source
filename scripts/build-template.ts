import { createHash } from "node:crypto";
import { chmod, cp, mkdir, readFile, rm, writeFile } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { build } from "esbuild";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const binary = join(root, ".agents", "bin", "kao.mjs");
await mkdir(dirname(binary), { recursive: true });
await build({
  entryPoints: [join(root, "scripts", "kao.ts")],
  outfile: binary,
  bundle: true,
  platform: "node",
  format: "esm",
  target: "node22",
  banner: { js: "#!/usr/bin/env node\nimport { createRequire as __kaoCreateRequire } from 'node:module'; const require = __kaoCreateRequire(import.meta.url);" },
  legalComments: "none",
});
await chmod(binary, 0o755);

const destination = join(root, ".dist", "kao-delivery-workspace-0.1.0");
await rm(destination, { recursive: true, force: true });
await mkdir(destination, { recursive: true });
for (const path of ["README.md", "AGENTS.md", "CLAUDE.md", "WORKFLOW.md", "workspace.yaml", ".gitignore", "agents", "context", "contributions", ".agents", ".codex", ".claude", "docs"]) {
  await cp(join(root, path), join(destination, path), { recursive: true });
}
for (const path of [
  "PLAN.md", "package.json", "package-lock.json", "tsconfig.json", "node_modules", "fixtures", "scripts", "test", ".dist",
  "docs/phase-0-proof.md", "docs/phase-0-host-results.md", "docs/create-plan-host-results.md", "docs/finish-work-host-results.md",
  "docs/initialization-host-results.md", "docs/repair-host-results.md", "docs/whats-next-host-results.md",
  "docs/cross-repository-host-results.md",
  "context/plans/.DS_Store",
]) await rm(join(destination, path), { recursive: true, force: true });
const bundleSha256 = createHash("sha256").update(await readFile(join(destination, ".agents", "bin", "kao.mjs"))).digest("hex");
await writeFile(join(destination, ".template-version"), "0.1.0\n", "utf8");
await writeFile(join(destination, "template-manifest.json"), `${JSON.stringify({
  name: "kao-delivery-workspace",
  version: "0.1.0",
  node: ">=22",
  command: "node .agents/bin/kao.mjs",
  bundle_sha256: bundleSha256,
  excluded_maintainer_inputs: ["PLAN.md", "package.json", "package-lock.json", "tsconfig.json", "node_modules/", "fixtures/", "scripts/", "test/"],
}, null, 2)}\n`, "utf8");
console.log(JSON.stringify({ version: "0.1.0", destination, binary }, null, 2));
