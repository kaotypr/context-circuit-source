import { createHash } from "node:crypto";
import { copyFile, mkdir, readFile, readdir, rm, stat, writeFile } from "node:fs/promises";
import { basename, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { list as listTar } from "tar";

export type ReleaseArtifact = { archive: string; folder: string; version: string };

async function sha256(file: string): Promise<string> {
  return createHash("sha256").update(await readFile(file)).digest("hex");
}

function forbiddenPath(relative: string): boolean {
  const parts = relative.split("/").filter(Boolean);
  const sourceInputs = new Set(["PLAN.md", "package.json", "package-lock.json", "tsconfig.json", "node_modules", "fixtures", "scripts", "test", ".dist"]);
  if (parts.some((part) => part === ".DS_Store" || part === "__MACOSX" || part.startsWith("._"))) return true;
  if (parts.some((part) => part === "node_modules")) return true;
  if (parts[0] && sourceInputs.has(parts[0])) return true;
  return parts.some((part) => part === ".env" || /^\.env\./.test(part) || /(^|[._-])(credentials?|secrets?|tokens?)([._-]|$)/i.test(part) || /\.(pem|key|p12)$/i.test(part));
}

export async function validateGeneratedRelease(distributionRoot: string, tag = ""): Promise<ReleaseArtifact> {
  const entries = await readdir(distributionRoot, { withFileTypes: true });
  const folders = entries.filter((entry) => entry.isDirectory() && entry.name.startsWith("context-circuit-"));
  if (folders.length !== 1) throw new Error("Generated distribution must contain exactly one Context Circuit folder.");
  const folder = join(distributionRoot, folders[0]!.name);
  const manifest = JSON.parse(await readFile(join(folder, "template-manifest.json"), "utf8"));
  const version = manifest.version;
  if (manifest.name !== "context-circuit" || typeof version !== "string") throw new Error("Generated metadata is not Context Circuit release metadata.");
  const expectedStem = `context-circuit-${version}`;
  if (basename(folder) !== expectedStem) throw new Error(`Generated folder must be ${expectedStem}.`);
  const archive = join(distributionRoot, `${expectedStem}.tar.gz`);
  await stat(archive);
  if (tag && tag !== `v${version}`) throw new Error(`Release tag ${tag} does not match v${version}.`);

  const names: string[] = [];
  const obsoleteIdentity: string[] = [];
  await listTar({
    file: archive,
    preservePaths: true,
    strict: true,
    onReadEntry(entry) {
      const name = entry.path.replace(/\/$/, "");
      names.push(name);
      const top = name.split("/")[0];
      if (top !== expectedStem) throw new Error(`Archive entry is outside ${expectedStem}: ${entry.path}`);
      const relative = name.slice(expectedStem.length + 1);
      if (relative && forbiddenPath(relative)) throw new Error(`Forbidden release input in archive: ${relative}`);
      if (entry.type === "File") {
        const chunks: Buffer[] = [];
        entry.on("data", (chunk) => chunks.push(Buffer.from(chunk)));
        entry.on("end", () => {
          if (/kao[- ]delivery[- ]workspace/i.test(Buffer.concat(chunks).toString("utf8"))) obsoleteIdentity.push(relative);
        });
      }
    },
  });
  if (obsoleteIdentity.length > 0) throw new Error(`Obsolete product identity found in: ${obsoleteIdentity.join(", ")}`);
  for (const required of [`${expectedStem}/template-manifest.json`, `${expectedStem}/.agents/bin/cc.mjs`]) {
    if (!names.includes(required)) throw new Error(`Archive is missing ${required}.`);
  }
  return { archive, folder, version };
}

async function readJson(file: string): Promise<Record<string, any>> {
  return JSON.parse(await readFile(file, "utf8"));
}

export async function syncReleaseArtifact(releaseRoot: string, release: ReleaseArtifact): Promise<{ changed: boolean }> {
  const artifactRoot = join(releaseRoot, "artifact");
  const expectedName = `context-circuit-${release.version}.tar.gz`;
  const destination = join(artifactRoot, expectedName);
  await mkdir(artifactRoot, { recursive: true });
  const artifactEntries = await readdir(artifactRoot, { withFileTypes: true });
  const existing = artifactEntries.find((entry) => entry.isFile() && entry.name === expectedName);
  if (existing && await sha256(destination) !== await sha256(release.archive)) {
    throw new Error(`Refusing to replace ${expectedName}: the existing same-version checksum differs.`);
  }
  const packageFile = join(releaseRoot, "package.json");
  const lockFile = join(releaseRoot, "package-lock.json");
  const packageJson = await readJson(packageFile);
  const lockJson = await readJson(lockFile);
  const alreadyCurrent = artifactEntries.length === 1 && existing &&
    packageJson.version === release.version && lockJson.version === release.version && lockJson.packages?.[""]?.version === release.version;
  if (alreadyCurrent) return { changed: false };

  for (const entry of artifactEntries) await rm(join(artifactRoot, entry.name), { recursive: true, force: true });
  await copyFile(release.archive, destination);
  packageJson.version = release.version;
  lockJson.version = release.version;
  if (lockJson.packages?.[""]) lockJson.packages[""].version = release.version;
  await writeFile(packageFile, `${JSON.stringify(packageJson, null, 2)}\n`, "utf8");
  await writeFile(lockFile, `${JSON.stringify(lockJson, null, 2)}\n`, "utf8");
  return { changed: true };
}

if (process.argv[1] && fileURLToPath(import.meta.url) === resolve(process.argv[1])) {
  const [distributionRoot, releaseRoot, tag = ""] = process.argv.slice(2);
  if (!distributionRoot || !releaseRoot) throw new Error("Usage: release-artifact.ts <distribution-root> <release-root> [tag]");
  const release = await validateGeneratedRelease(distributionRoot, tag);
  const result = await syncReleaseArtifact(releaseRoot, release);
  process.stdout.write(`${JSON.stringify({ ...release, ...result })}\n`);
}
