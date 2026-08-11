import { createHash } from "node:crypto";
import { chmod, cp, mkdir, readFile, readdir, rm, writeFile } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { build } from "esbuild";
import { neutralWorkspaceContext, renderWorkspaceContext } from "./lib/workspace-context.js";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const binary = join(root, ".agents", "bin", "cc.mjs");
await mkdir(dirname(binary), { recursive: true });
await build({
  entryPoints: [join(root, "scripts", "cc.ts")],
  outfile: binary,
  bundle: true,
  platform: "node",
  format: "esm",
  target: "node22",
  banner: { js: "#!/usr/bin/env node\nimport { createRequire as __ccCreateRequire } from 'node:module'; const require = __ccCreateRequire(import.meta.url);" },
  legalComments: "none",
});
await chmod(binary, 0o755);

async function removeLocalMetadata(directory: string): Promise<void> {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const path = join(directory, entry.name);
    if (entry.name === ".DS_Store" || entry.name === "__MACOSX") await rm(path, { recursive: true, force: true });
    else if (entry.isDirectory()) await removeLocalMetadata(path);
  }
}

const distributionRoot = join(root, ".dist");
const destination = join(distributionRoot, "context-circuit-0.2.0");
await rm(distributionRoot, { recursive: true, force: true });
await mkdir(destination, { recursive: true });
for (const path of ["README.md", "AGENTS.md", "CLAUDE.md", "WORKFLOW.md", "workspace.yaml", ".gitignore", "agents", "context", "contributions", ".agents", ".codex", ".claude", "docs"]) {
  await cp(join(root, path), join(destination, path), { recursive: true });
}
for (const [path, contents] of Object.entries(renderWorkspaceContext(neutralWorkspaceContext))) {
  await writeFile(join(destination, path), contents, "utf8");
}
for (const path of [
  "PLAN.md", "package.json", "package-lock.json", "tsconfig.json", "node_modules", "fixtures", "scripts", "test", ".dist",
  "docs/phase-0-proof.md", "docs/phase-0-host-results.md", "docs/create-plan-host-results.md", "docs/finish-work-host-results.md",
  "docs/initialization-host-results.md", "docs/repair-host-results.md", "docs/whats-next-host-results.md",
  "docs/cross-repository-host-results.md",
]) await rm(join(destination, path), { recursive: true, force: true });
await removeLocalMetadata(destination);
const bundleSha256 = createHash("sha256").update(await readFile(join(destination, ".agents", "bin", "cc.mjs"))).digest("hex");
await writeFile(join(destination, ".template-version"), "0.2.0\n", "utf8");
await writeFile(join(destination, "template-manifest.json"), `${JSON.stringify({
  name: "context-circuit",
  version: "0.2.0",
  node: ">=22",
  command: "node .agents/bin/cc.mjs",
  bundle_sha256: bundleSha256,
  excluded_maintainer_inputs: ["PLAN.md", "package.json", "package-lock.json", "tsconfig.json", "node_modules/", "fixtures/", "scripts/", "test/"],
}, null, 2)}\n`, "utf8");
console.log(JSON.stringify({ version: "0.2.0", destination, binary }, null, 2));
