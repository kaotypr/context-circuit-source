import assert from "node:assert/strict";
import { mkdtemp, mkdir, readFile, readdir, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { syncReleaseArtifact, validateGeneratedRelease } from "../scripts/release-artifact.js";
import { projectRoot } from "./helpers.js";

test("generated release metadata, folder, archive, and tag agree", async () => {
  const release = await validateGeneratedRelease(join(projectRoot, ".dist"), "v0.2.0");
  assert.equal(release.version, "0.2.0");
  await assert.rejects(validateGeneratedRelease(join(projectRoot, ".dist"), "v0.2.1"), /does not match/);
});

async function releaseRepository(t: test.TestContext): Promise<string> {
  const root = await mkdtemp(join(tmpdir(), "context-circuit-release-sync-"));
  t.after(() => rm(root, { recursive: true, force: true }));
  await mkdir(join(root, "artifact"));
  await writeFile(join(root, "artifact", "README.md"), "placeholder\n");
  await writeFile(join(root, "package.json"), `${JSON.stringify({ name: "context-circuit", version: "0.0.1" }, null, 2)}\n`);
  await writeFile(join(root, "package-lock.json"), `${JSON.stringify({ name: "context-circuit", version: "0.0.1", lockfileVersion: 3, packages: { "": { name: "context-circuit", version: "0.0.1" } } }, null, 2)}\n`);
  return root;
}

test("release synchronization is idempotent", async (t) => {
  const root = await releaseRepository(t);
  const release = await validateGeneratedRelease(join(projectRoot, ".dist"));
  assert.deepEqual(await syncReleaseArtifact(root, release), { changed: true });
  assert.deepEqual(await syncReleaseArtifact(root, release), { changed: false });
  assert.equal(JSON.parse(await readFile(join(root, "package.json"), "utf8")).version, release.version);
  assert.equal((await readdir(join(root, "artifact"))).length, 1);
});

test("release synchronization rejects a conflicting same-version artifact", async (t) => {
  const root = await releaseRepository(t);
  const release = await validateGeneratedRelease(join(projectRoot, ".dist"));
  await rm(join(root, "artifact", "README.md"));
  await writeFile(join(root, "artifact", `context-circuit-${release.version}.tar.gz`), "different bytes");
  await assert.rejects(syncReleaseArtifact(root, release), /checksum differs/);
});
