import { chmod, lstat, mkdir, open, readFile, realpath, rename } from "node:fs/promises"
import { dirname, resolve, sep } from "node:path"

export function assertInside(root: string, candidate: string): string {
  const resolvedRoot = resolve(root)
  const resolvedCandidate = resolve(candidate)
  if (resolvedCandidate !== resolvedRoot && !resolvedCandidate.startsWith(`${resolvedRoot}${sep}`)) throw new Error(`Refusing path outside ${resolvedRoot}: ${resolvedCandidate}`)
  return resolvedCandidate
}

export async function readJsonRegularInside<T>(root: string, candidate: string, label: string): Promise<T> {
  const path = assertInside(root, candidate)
  const info = await lstat(path)
  if (!info.isFile() || info.isSymbolicLink()) throw new Error(`${label} must be a regular file: ${path}`)
  assertInside(await realpath(root), await realpath(path))
  return JSON.parse(await readFile(path, "utf8")) as T
}

export async function writeTextAtomic(path: string, value: string, mode = 0o644): Promise<void> {
  await mkdir(dirname(path), { recursive: true, mode: 0o755 })
  const temporary = `${path}.${process.pid}.tmp`
  const handle = await open(temporary, "wx", mode)
  try {
    await handle.writeFile(value, "utf8")
    await handle.sync()
  } finally {
    await handle.close()
  }
  await rename(temporary, path)
  await chmod(path, mode)
}

export async function writeTextExclusive(path: string, value: string, mode = 0o644): Promise<void> {
  await mkdir(dirname(path), { recursive: true, mode: 0o755 })
  const handle = await open(path, "wx", mode)
  try {
    await handle.writeFile(value, "utf8")
    await handle.sync()
  } finally {
    await handle.close()
  }
  await chmod(path, mode)
}
