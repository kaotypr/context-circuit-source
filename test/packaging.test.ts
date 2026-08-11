import assert from "node:assert/strict";
import { access, readFile, readdir } from "node:fs/promises";
import { spawnSync } from "node:child_process";
import { join } from "node:path";
import test from "node:test";
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

test("distributable excludes maintainer-only inputs and records its bundle digest", async () => {
  const destination = join(projectRoot, ".dist", "context-circuit-0.2.0");
  const manifest = JSON.parse(await readFile(join(destination, "template-manifest.json"), "utf8"));
  assert.equal(manifest.version, "0.2.0");
  assert.match(manifest.bundle_sha256, /^[a-f0-9]{64}$/);
  for (const path of ["PLAN.md", "package.json", "node_modules", "fixtures", "scripts", "test"]) {
    await assert.rejects(access(join(destination, path)));
  }
  await access(join(destination, ".agents", "bin", "cc.mjs"));
  await assert.rejects(access(join(destination, ".agents", "bin", "kao.mjs")));
  await access(join(destination, "context", "plans", ".gitkeep"));
  await assert.rejects(access(join(destination, "repositories")));
  await assert.rejects(access(join(destination, "agents", "frontend.md")));
  assert.match(await readFile(join(destination, "workspace.yaml"), "utf8"), /repositories: \{\}/);
  assert.doesNotMatch(await readFile(join(destination, ".gitignore"), "utf8"), /repositories\/frontend/);
  assert.match(await readFile(join(destination, "context", "PROJECT.md"), "utf8"), /has not been initialized/);
  assert.doesNotMatch(await readFile(join(destination, "context", "PROJECT.md"), "utf8"), /Context Circuit 0\.2\.0/);
  assert.doesNotMatch(await readFile(join(destination, "WORKFLOW.md"), "utf8"), /template is in team mode/i);

  const pending = [destination];
  while (pending.length > 0) {
    const directory = pending.pop()!;
    for (const entry of await readdir(directory, { withFileTypes: true })) {
      assert.notEqual(entry.name, ".DS_Store");
      assert.notEqual(entry.name, "__MACOSX");
      if (entry.isDirectory()) pending.push(join(directory, entry.name));
    }
  }
});
