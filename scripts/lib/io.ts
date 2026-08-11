import { chmod, lstat, mkdir, open, rename, unlink } from "node:fs/promises";
import { dirname, resolve, sep } from "node:path";

export function assertInside(root: string, candidate: string): string {
  const resolvedRoot = resolve(root);
  const resolvedCandidate = resolve(candidate);
  if (resolvedCandidate !== resolvedRoot && !resolvedCandidate.startsWith(`${resolvedRoot}${sep}`)) {
    throw new Error(`Refusing path outside ${resolvedRoot}: ${resolvedCandidate}`);
  }
  return resolvedCandidate;
}

export async function ensurePrivateDirectory(path: string): Promise<void> {
  try {
    await mkdir(path, { recursive: true, mode: 0o700 });
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "EEXIST") throw error;
  }
  const info = await lstat(path);
  if (!info.isDirectory() || info.isSymbolicLink()) {
    throw new Error(`Runtime path must be a real directory: ${path}`);
  }
  await chmod(path, 0o700);
}

export async function writeJsonAtomic(path: string, value: unknown): Promise<void> {
  await ensurePrivateDirectory(dirname(path));
  const temporary = `${path}.${process.pid}.tmp`;
  const handle = await open(temporary, "wx", 0o600);
  try {
    await handle.writeFile(`${JSON.stringify(value, null, 2)}\n`, "utf8");
    await handle.sync();
  } finally {
    await handle.close();
  }
  await rename(temporary, path);
  await chmod(path, 0o600);
}

export async function withExclusiveFile<T>(path: string, operation: () => Promise<T>): Promise<T> {
  await ensurePrivateDirectory(dirname(path));
  let handle;
  try {
    handle = await open(path, "wx", 0o600);
    await handle.writeFile(`${process.pid}\n`, "utf8");
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code === "EEXIST") {
      throw new Error(`Another runtime recorder holds the lock: ${path}`);
    }
    throw error;
  }
  try {
    return await operation();
  } finally {
    await handle.close();
    await unlink(path);
  }
}
