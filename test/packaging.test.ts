import assert from "node:assert/strict";
import { access, mkdtemp, readFile, readdir, rm } from "node:fs/promises";
import { spawnSync } from "node:child_process";
import { tmpdir } from "node:os";
import { join, relative } from "node:path";
import test from "node:test";
import { extract as extractTar, list as listTar } from "tar";
import { projectRoot } from "./helpers.js";

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
  for (const path of ["PLAN.md", "package.json", "node_modules", "fixtures", "scripts", "test"]) {
    await assert.rejects(access(join(destination, path)));
  }
  await access(join(destination, ".agents", "bin", "cc.mjs"));
  const canonicalSkills = (await readdir(join(projectRoot, ".agents", "skills"), { withFileTypes: true }))
    .filter((entry) => entry.isDirectory()).map((entry) => entry.name).sort();
  const expectedHostSkills = [
    "w-configure-workspace",
    "w-create-plan",
    "w-finish-work",
    "w-gather-context",
    "w-initialize-workspace",
    "w-publish-plan-tasks",
    "w-run-task",
    "w-sync-context",
    "w-whats-next",
  ];
  assert.deepEqual(canonicalSkills, expectedHostSkills);
  assert.ok(canonicalSkills.every((skill) => skill.startsWith("w-")));
  for (const skill of canonicalSkills) {
    await access(join(projectRoot, ".agents", "skills", skill, "SKILL.md"));
    await access(join(projectRoot, ".agents", "skills", skill, "agents", "openai.yaml"));
    await access(join(projectRoot, ".codex", "skills", skill, "SKILL.md"));
    await access(join(projectRoot, ".claude", "commands", `${skill}.md`));
    await access(join(destination, ".agents", "skills", skill, "SKILL.md"));
    await access(join(destination, ".agents", "skills", skill, "agents", "openai.yaml"));
    await access(join(destination, ".codex", "skills", skill, "SKILL.md"));
    await access(join(destination, ".claude", "commands", `${skill}.md`));
  }
  for (const path of [
    ".agents/skills/w-configure-workspace/SKILL.md",
    ".agents/skills/w-gather-context/SKILL.md",
    ".codex/skills/w-configure-workspace/SKILL.md",
    ".codex/skills/w-gather-context/SKILL.md",
    ".claude/commands/w-configure-workspace.md",
    ".claude/commands/w-gather-context.md",
  ]) await access(join(destination, path));
  for (const skill of [
    "configure-workspace",
    "create-plan",
    "finish-work",
    "gather-context",
    "initialize-workspace",
    "publish-plan-tasks",
    "run-task",
    "sync-context",
    "whats-next",
  ]) {
    await assert.rejects(access(join(projectRoot, ".agents", "skills", skill)));
    await assert.rejects(access(join(projectRoot, ".codex", "skills", skill)));
    await assert.rejects(access(join(projectRoot, ".claude", "commands", `${skill}.md`)));
    await assert.rejects(access(join(destination, ".agents", "skills", skill)));
    await assert.rejects(access(join(destination, ".codex", "skills", skill)));
    await assert.rejects(access(join(destination, ".claude", "commands", `${skill}.md`)));
  }
  await assert.rejects(access(join(destination, ".agents", "bin", "kao.mjs")));
  await assert.rejects(access(join(destination, ".template-version")));
  await access(join(destination, "context", "plans", ".gitkeep"));
  await assert.rejects(access(join(destination, "repositories")));
  await assert.rejects(access(join(destination, "agents", "frontend.md")));
  const distributedDocs = (await readdir(join(destination, "docs"))).sort();
  assert.deepEqual(distributedDocs, ["command-reference.md", "configuration.md", "getting-started.md", "product-knowledge.md", "using-the-wrapper.md"]);
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
