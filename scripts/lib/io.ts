import { chmod, lstat, mkdir, open, readFile, realpath, rename, unlink } from "node:fs/promises";
import { dirname, resolve, sep } from "node:path";

export function assertInside(root: string, candidate: string): string {
  const resolvedRoot = resolve(root);
  const resolvedCandidate = resolve(candidate);
  if (resolvedCandidate !== resolvedRoot && !resolvedCandidate.startsWith(`${resolvedRoot}${sep}`)) {
    throw new Error(`Refusing path outside ${resolvedRoot}: ${resolvedCandidate}`);
  }
  return resolvedCandidate;
}

export async function readJsonRegularInside<T>(root: string, candidate: string, label: string): Promise<T> {
  const path = assertInside(root, candidate);
  const info = await lstat(path);
  if (!info.isFile() || info.isSymbolicLink()) throw new Error(`${label} must be a regular file: ${path}`);
  const realRoot = await realpath(root);
  const realPath = await realpath(path);
  assertInside(realRoot, realPath);
  return JSON.parse(await readFile(path, "utf8")) as T;
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

export async function writeTextAtomic(path: string, value: string, mode = 0o644): Promise<void> {
  const temporary = `${path}.${process.pid}.tmp`;
  const handle = await open(temporary, "wx", mode);
  try {
    await handle.writeFile(value, "utf8");
    await handle.sync();
  } finally {
    await handle.close();
  }
  await rename(temporary, path);
  await chmod(path, mode);
}

export interface TextTransactionEntry { path: string; value: string; mode?: number }

export async function writeTextTransaction(entries: TextTransactionEntry[], options: { failRenameAt?: number } = {}): Promise<void> {
  const unique = new Set(entries.map((entry) => resolve(entry.path)));
  if (unique.size !== entries.length) throw new Error("Text transaction targets must be unique");
  const nonce = `${process.pid}.${Date.now()}.${Math.random().toString(16).slice(2)}`;
  const staged: Array<{ target: string; temporary: string; backup: string; existed: boolean }> = [];
  const backedUp: typeof staged = [];
  const installed: typeof staged = [];
  let renameCount = 0;
  const transactionRename = async (from: string, to: string): Promise<void> => {
    renameCount += 1;
    if (options.failRenameAt === renameCount) throw new Error(`Injected transaction rename failure at ${renameCount}`);
    await rename(from, to);
  };
  try {
    // Stage every sibling first. Any permissions, space, or open failure occurs before a target changes.
    for (const entry of entries) {
      const target = resolve(entry.path);
      const temporary = `${target}.${nonce}.stage`;
      const backup = `${target}.${nonce}.backup`;
      let existed = false;
      let mode = entry.mode ?? 0o644;
      try {
        const info = await lstat(target);
        if (!info.isFile() || info.isSymbolicLink()) throw new Error(`Transaction target must be a regular file: ${target}`);
        existed = true;
        mode = info.mode & 0o777;
      } catch (error) {
        if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
      }
      const handle = await open(temporary, "wx", mode);
      try { await handle.writeFile(entry.value, "utf8"); await handle.sync(); } finally { await handle.close(); }
      await chmod(temporary, mode);
      staged.push({ target, temporary, backup, existed });
    }
    for (const item of staged) {
      if (!item.existed) continue;
      await transactionRename(item.target, item.backup);
      backedUp.push(item);
    }
    for (const item of staged) {
      await transactionRename(item.temporary, item.target);
      installed.push(item);
    }
    for (const item of backedUp) await unlink(item.backup);
  } catch (error) {
    for (const item of installed.reverse()) {
      try { await unlink(item.target); } catch (cleanupError) { if ((cleanupError as NodeJS.ErrnoException).code !== "ENOENT") throw cleanupError; }
    }
    for (const item of backedUp.reverse()) {
      try { await rename(item.backup, item.target); } catch (rollbackError) {
        throw new AggregateError([error, rollbackError], `Text transaction failed and rollback could not restore ${item.target}`);
      }
    }
    throw error;
  } finally {
    for (const item of staged) {
      for (const path of [item.temporary, item.backup]) {
        try { await unlink(path); } catch (cleanupError) { if ((cleanupError as NodeJS.ErrnoException).code !== "ENOENT") throw cleanupError; }
      }
    }
  }
}

export async function writeTextExclusive(path: string, value: string, mode = 0o644): Promise<void> {
  await mkdir(dirname(path), { recursive: true, mode: 0o755 });
  const handle = await open(path, "wx", mode);
  try {
    await handle.writeFile(value, "utf8");
    await handle.sync();
  } finally {
    await handle.close();
  }
  await chmod(path, mode);
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
