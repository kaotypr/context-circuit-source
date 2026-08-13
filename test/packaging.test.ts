import assert from "node:assert/strict";
import { access, mkdtemp, readFile, readdir, rm } from "node:fs/promises";
import { spawnSync } from "node:child_process";
import { tmpdir } from "node:os";
import { join, relative } from "node:path";
import test from "node:test";
import { extract as extractTar, list as listTar } from "tar";
import { projectRoot } from "./helpers.js";

const hostWorkflows = [
  "cc-configure-workspace",
  "cc-create-plan",
  "cc-finish-work",
  "cc-gather-context",
  "cc-import-context",
  "cc-initialize-workspace",
  "cc-publish-plan-tasks",
  "cc-run-task",
  "cc-sync-context",
  "cc-whats-next",
] as const;

const displayLabels = new Map([
  ["cc-configure-workspace", "CC Configure Workspace"],
  ["cc-create-plan", "CC Create Plan"],
  ["cc-finish-work", "CC Finish Work"],
  ["cc-gather-context", "CC Gather Context"],
  ["cc-import-context", "CC Import Context"],
  ["cc-initialize-workspace", "CC Initialize Workspace"],
  ["cc-publish-plan-tasks", "CC Publish Plan Tasks"],
  ["cc-run-task", "CC Run Task"],
  ["cc-sync-context", "CC Sync Context"],
  ["cc-whats-next", "CC What's Next"],
]);

const importContextAssets = [
  ".agents/contracts/import-context-request.schema.json",
  ".agents/skills/cc-import-context/SKILL.md",
  ".agents/skills/cc-import-context/agents/openai.yaml",
  ".codex/skills/cc-import-context/SKILL.md",
  ".claude/commands/cc-import-context.md",
  "docs/command-reference.md",
  "docs/context-sync.md",
] as const;

const build = spawnSync(process.execPath, ["--import", "tsx", "scripts/build-template.ts"], { cwd: projectRoot, encoding: "utf8" });
if (build.status !== 0) throw new Error(build.stderr || build.stdout);

test("bundled command validates without TypeScript tooling or installed wrapper dependencies", async () => {
  const binary = join(projectRoot, ".agents", "bin", "cc.mjs");
  await access(binary);
  const result = spawnSync(process.execPath, [binary, "validate"], { cwd: projectRoot, encoding: "utf8" });
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /Valid workspace/);
});

test("distributable and release archive contain only wrapper inputs", async () => {
  const destination = join(projectRoot, ".dist", "context-circuit-0.2.1");
  const archive = join(projectRoot, ".dist", "context-circuit-0.2.1.tar.gz");
  const manifest = JSON.parse(await readFile(join(destination, "template-manifest.json"), "utf8"));
  assert.equal(manifest.version, "0.2.1");
  assert.match(manifest.bundle_sha256, /^[a-f0-9]{64}$/);
  assert.ok(Array.isArray(manifest.file_inventory));
  assert.ok(manifest.file_inventory.includes(".agents/bin/cc.mjs"));
  assert.ok(manifest.file_inventory.includes("template-manifest.json"));
  assert.ok(manifest.file_inventory.includes(".agents/skills/cc-run-task/SKILL.md"));
  for (const path of importContextAssets) {
    await access(join(projectRoot, path));
    assert.ok(manifest.file_inventory.includes(path), `${path} is missing from the trusted source inventory`);
    await access(join(destination, path));
  }
  assert.equal(manifest.file_inventory.some((path: string) => path.includes("/w-") || path.includes("/configure-workspace/") || path.includes("/run-task/")), false);
  for (const path of ["PLAN.md", "package.json", "node_modules", "fixtures", "scripts", "test"]) {
    await assert.rejects(access(join(destination, path)));
  }
  await access(join(destination, ".agents", "bin", "cc.mjs"));
  const canonicalSkills = (await readdir(join(projectRoot, ".agents", "skills"), { withFileTypes: true }))
    .filter((entry) => entry.isDirectory()).map((entry) => entry.name).sort();
  assert.deepEqual(canonicalSkills, hostWorkflows);
  assert.ok(canonicalSkills.every((skill) => skill.startsWith("cc-")));
  assert.equal(canonicalSkills.some((skill) => skill.startsWith("w-")), false);
  for (const skill of canonicalSkills) {
    await access(join(projectRoot, ".agents", "skills", skill, "SKILL.md"));
    await access(join(projectRoot, ".agents", "skills", skill, "agents", "openai.yaml"));
    await access(join(projectRoot, ".codex", "skills", skill, "SKILL.md"));
    await access(join(projectRoot, ".claude", "commands", `${skill}.md`));
    await access(join(destination, ".agents", "skills", skill, "SKILL.md"));
    await access(join(destination, ".agents", "skills", skill, "agents", "openai.yaml"));
    await access(join(destination, ".codex", "skills", skill, "SKILL.md"));
    await access(join(destination, ".claude", "commands", `${skill}.md`));
    assert.match(await readFile(join(projectRoot, ".agents", "skills", skill, "agents", "openai.yaml"), "utf8"), new RegExp(`display_name: "${displayLabels.get(skill)}"`));
    for (const path of [
      join(projectRoot, ".agents", "skills", skill, "SKILL.md"),
      join(projectRoot, ".codex", "skills", skill, "SKILL.md"),
      join(projectRoot, ".claude", "commands", `${skill}.md`),
    ]) assert.doesNotMatch(await readFile(path, "utf8"), /\bw-(?:configure-workspace|create-plan|finish-work|gather-context|initialize-workspace|publish-plan-tasks|run-task|sync-context|whats-next)\b/);
  }
  for (const path of [
    ".agents/skills/cc-configure-workspace/SKILL.md",
    ".agents/skills/cc-gather-context/SKILL.md",
    ".codex/skills/cc-configure-workspace/SKILL.md",
    ".codex/skills/cc-gather-context/SKILL.md",
    ".claude/commands/cc-configure-workspace.md",
    ".claude/commands/cc-gather-context.md",
  ]) await access(join(destination, path));
  for (const prefix of ["w-", ""]) {
    for (const skill of hostWorkflows.map((workflow) => workflow.slice(3))) {
      await assert.rejects(access(join(projectRoot, ".agents", "skills", `${prefix}${skill}`)));
      await assert.rejects(access(join(projectRoot, ".codex", "skills", `${prefix}${skill}`)));
      await assert.rejects(access(join(projectRoot, ".claude", "commands", `${prefix}${skill}.md`)));
      await assert.rejects(access(join(destination, ".agents", "skills", `${prefix}${skill}`)));
      await assert.rejects(access(join(destination, ".codex", "skills", `${prefix}${skill}`)));
      await assert.rejects(access(join(destination, ".claude", "commands", `${prefix}${skill}.md`)));
    }
  }
  const bundled = await readFile(join(projectRoot, ".agents", "bin", "cc.mjs"), "utf8");
  assert.match(bundled, /\.agents\/skills\/cc-run-task\/SKILL\.md/);
  assert.doesNotMatch(bundled, /\.agents\/skills\/w-run-task\/SKILL\.md/);
  assert.match(bundled, /run-task/);
  assert.match(bundled, /whats-next/);
  assert.match(bundled, /Usage: import-context --request <import-context-request\.json>/);
  assert.match(bundled, /sync-context.*--request/);
  for (const path of importContextAssets) assert.ok(bundled.includes(path), `${path} is missing from the standalone bundle inventory`);
  const generatedReadme = await readFile(join(destination, "README.md"), "utf8");
  assert.match(generatedReadme, /\$cc-run-task/);
  assert.doesNotMatch(generatedReadme, /\$w-run-task|\/w-run-task/);
  await assert.rejects(access(join(destination, ".agents", "bin", "kao.mjs")));
  await assert.rejects(access(join(destination, ".template-version")));
  await access(join(destination, "context", "plans", ".gitkeep"));
  await assert.rejects(access(join(destination, "repositories")));
  await assert.rejects(access(join(destination, "agents", "frontend.md")));
  const distributedDocs = (await readdir(join(destination, "docs"))).sort();
  assert.deepEqual(distributedDocs, ["command-reference.md", "configuration.md", "context-sync.md", "getting-started.md", "product-knowledge.md", "using-the-wrapper.md"]);
  assert.match(await readFile(join(destination, "workspace.yaml"), "utf8"), /repositories: \{\}/);
  assert.doesNotMatch(await readFile(join(destination, ".gitignore"), "utf8"), /repositories\/frontend/);
  assert.match(await readFile(join(destination, "context", "PROJECT.md"), "utf8"), /has not been initialized/);
  assert.doesNotMatch(await readFile(join(destination, "context", "PROJECT.md"), "utf8"), /Context Circuit 0\.2\.0/);
  assert.doesNotMatch(await readFile(join(destination, "WORKFLOW.md"), "utf8"), /template is in team mode/i);

  const pending = [destination];
  const actualInventory: string[] = [];
  while (pending.length > 0) {
    const directory = pending.pop()!;
    for (const entry of await readdir(directory, { withFileTypes: true })) {
      assert.notEqual(entry.name, ".DS_Store");
      assert.notEqual(entry.name, "__MACOSX");
      assert.equal(entry.name.startsWith("._"), false);
      if (entry.isDirectory()) pending.push(join(directory, entry.name));
      else if (entry.isFile()) actualInventory.push(relative(destination, join(directory, entry.name)).replaceAll("\\", "/"));
    }
  }
  assert.deepEqual(actualInventory.sort(), manifest.file_inventory);

  await access(archive);
  const archiveEntries: string[] = [];
  await listTar({ file: archive, onReadEntry: (entry) => archiveEntries.push(entry.path) });
  assert.ok(archiveEntries.includes("context-circuit-0.2.1/README.md"));
  assert.ok(archiveEntries.includes("context-circuit-0.2.1/.agents/bin/cc.mjs"));
  for (const path of importContextAssets) assert.ok(archiveEntries.includes(`context-circuit-0.2.1/${path}`), `${path} is missing from the release archive`);
  assert.equal(archiveEntries.some((path) => path.endsWith("/.template-version")), false);
  assert.equal(archiveEntries.some((path) => path.split("/").some((part) => part === ".DS_Store" || part === "__MACOSX" || part.startsWith("._"))), false);
});

test("release archive validates after extraction without wrapper dependencies", async (t) => {
  const extractionRoot = await mkdtemp(join(tmpdir(), "context-circuit-release-"));
  t.after(async () => rm(extractionRoot, { recursive: true, force: true }));
  await extractTar({ cwd: extractionRoot, file: join(projectRoot, ".dist", "context-circuit-0.2.1.tar.gz") });
  const wrapperRoot = join(extractionRoot, "context-circuit-0.2.1");
  const result = spawnSync(process.execPath, [join(wrapperRoot, ".agents", "bin", "cc.mjs"), "validate", "--check-paths", "--check-documents"], {
    cwd: wrapperRoot,
    encoding: "utf8",
  });
  assert.equal(result.status, 0, result.stderr || result.stdout);
});
