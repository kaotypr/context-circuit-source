import assert from "node:assert/strict";
import { access, readFile } from "node:fs/promises";
import { spawnSync } from "node:child_process";
import { join } from "node:path";
import test from "node:test";
import { projectRoot } from "./helpers.js";

const build = spawnSync(process.execPath, ["--import", "tsx", "scripts/build-template.ts"], { cwd: projectRoot, encoding: "utf8" });
if (build.status !== 0) throw new Error(build.stderr || build.stdout);

test("bundled command validates without TypeScript tooling or installed wrapper dependencies", async () => {
  const binary = join(projectRoot, ".agents", "bin", "kao.mjs");
  await access(binary);
  const result = spawnSync(process.execPath, [binary, "validate"], { cwd: projectRoot, encoding: "utf8" });
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /Valid workspace/);
});

test("distributable excludes maintainer-only inputs and records its bundle digest", async () => {
  const destination = join(projectRoot, ".dist", "kao-delivery-workspace-0.1.0");
  const manifest = JSON.parse(await readFile(join(destination, "template-manifest.json"), "utf8"));
  assert.equal(manifest.version, "0.1.0");
  assert.match(manifest.bundle_sha256, /^[a-f0-9]{64}$/);
  for (const path of ["PLAN.md", "package.json", "node_modules", "fixtures", "scripts", "test"]) {
    await assert.rejects(access(join(destination, path)));
  }
  await access(join(destination, ".agents", "bin", "kao.mjs"));
});
