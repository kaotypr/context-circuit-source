import { createHash } from "node:crypto";
import { chmod, cp, mkdir, readFile, readdir, rm, writeFile } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { build } from "esbuild";
import { create as createTar } from "tar";
import { neutralWorkspaceContext, renderWorkspaceContext } from "./lib/workspace-context.js";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const packageJson = JSON.parse(await readFile(join(root, "package.json"), "utf8"));
const version = packageJson.version;
if (typeof version !== "string" || !/^\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$/.test(version)) {
  throw new Error(`Source package version is invalid: ${String(version)}`);
}
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
    if (entry.name === ".DS_Store" || entry.name === "__MACOSX" || entry.name.startsWith("._")) await rm(path, { recursive: true, force: true });
    else if (entry.isDirectory()) await removeLocalMetadata(path);
  }
}

const distributionRoot = join(root, ".dist");
const releaseStem = `context-circuit-${version}`;
const destination = join(distributionRoot, releaseStem);
await rm(distributionRoot, { recursive: true, force: true });
await mkdir(destination, { recursive: true });
for (const path of ["README.md", "AGENTS.md", "CLAUDE.md", "WORKFLOW.md", "workspace.yaml", ".gitignore", "agents", "context", "contributions", ".agents", ".codex", ".claude"]) {
  await cp(join(root, path), join(destination, path), { recursive: true });
}
await mkdir(join(destination, "docs"), { recursive: true });
for (const path of ["getting-started.md", "using-the-wrapper.md", "configuration.md", "command-reference.md"]) {
  await cp(join(root, "docs", path), join(destination, "docs", path));
}
for (const [path, contents] of Object.entries(renderWorkspaceContext(neutralWorkspaceContext))) {
  await writeFile(join(destination, path), contents, "utf8");
}
for (const path of [
  "PLAN.md", "package.json", "package-lock.json", "tsconfig.json", "node_modules", "fixtures", "scripts", "test", ".dist",
  "docs/phase-0-proof.md", "docs/phase-0-host-results.md", "docs/create-plan-host-results.md", "docs/finish-work-host-results.md",
  "docs/initialization-host-results.md", "docs/repair-host-results.md", "docs/whats-next-host-results.md", "docs/cross-repository-host-results.md",
]) await rm(join(destination, path), { recursive: true, force: true });
await removeLocalMetadata(destination);
const bundleSha256 = createHash("sha256").update(await readFile(join(destination, ".agents", "bin", "cc.mjs"))).digest("hex");
await writeFile(join(destination, "template-manifest.json"), `${JSON.stringify({
  name: "context-circuit",
  version,
  node: ">=22",
  command: "node .agents/bin/cc.mjs",
  bundle_sha256: bundleSha256,
  excluded_maintainer_inputs: ["PLAN.md", "package.json", "package-lock.json", "tsconfig.json", "node_modules/", "fixtures/", "scripts/", "test/"],
}, null, 2)}\n`, "utf8");
const archive = join(distributionRoot, `${releaseStem}.tar.gz`);
await createTar({
  cwd: distributionRoot,
  file: archive,
  gzip: true,
  noMtime: true,
  portable: true,
  filter: (path) => !path.split("/").some((part) => part === ".DS_Store" || part === "__MACOSX" || part.startsWith("._")),
}, [releaseStem]);
console.log(JSON.stringify({ version, destination, archive, binary }, null, 2));
