import { createHash, randomBytes } from "node:crypto";
import { readFile } from "node:fs/promises";
import { join } from "node:path";
import { ensurePrivateDirectory, writeJsonAtomic } from "./io.js";

function utcStamp(now: Date): { day: string; instant: string } {
  const iso = now.toISOString();
  return {
    day: iso.slice(0, 10).replaceAll("-", ""),
    instant: iso.slice(0, 19).replaceAll("-", "").replaceAll(":", ""),
  };
}

export async function generateIds(
  runtimeRoot: string,
  request: string,
  now = new Date(),
  discriminator = randomBytes(4).toString("hex"),
): Promise<{ workId: string; runId: string }> {
  if (!/^[a-f0-9]{8}$/.test(discriminator)) {
    throw new Error("Run discriminator must contain exactly eight lowercase hexadecimal characters");
  }
  await ensurePrivateDirectory(runtimeRoot);
  const { day, instant } = utcStamp(now);
  const statePath = join(runtimeRoot, "id-state.json");
  let state: { day: string; next: number } = { day, next: 1 };
  try {
    state = JSON.parse(await readFile(statePath, "utf8")) as typeof state;
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
  }
  const sequence = state.day === day ? state.next : 1;
  await writeJsonAtomic(statePath, { day, next: sequence + 1 });
  const requestFingerprint = createHash("sha256").update(request).digest("hex").slice(0, 4);
  return {
    workId: `ADHOC-${day}-${String(sequence).padStart(3, "0")}`,
    runId: `${instant}Z-${discriminator.slice(0, 4)}${requestFingerprint}`,
  };
}

export function slugify(value: string): string {
  return value
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "")
    .slice(0, 32) || "task";
}
