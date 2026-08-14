import { lstat, readFile } from "node:fs/promises"
import { join, relative, resolve } from "node:path"
import { parse as parseYaml } from "yaml"
import type { RepositoryConfig, WorkspaceConfig } from "./types.js"
import { contextReferenceError, remoteReferenceError } from "./safe-reference.js"

export async function readData(path: string): Promise<unknown> {
  const raw = await readFile(path, "utf8")
  return path.endsWith(".yaml") || path.endsWith(".yml") ? parseYaml(raw) : JSON.parse(raw)
}

function record(value: unknown): value is Record<string, unknown> {
  return Boolean(value) && typeof value === "object" && !Array.isArray(value)
}

function nonEmptyString(value: unknown): value is string {
  return typeof value === "string" && value.trim().length > 0
}

function stringList(value: unknown): value is string[] {
  return Array.isArray(value) && value.every((item) => typeof item === "string" && item.trim().length > 0)
}

export function workspaceErrors(value: unknown): string[] {
  if (!record(value)) return ["workspace.yaml must contain a mapping"]
  const errors: string[] = []
  const workspace = value.workspace
  if (!record(workspace)) errors.push("workspace is required")
  else {
    if (!nonEmptyString(workspace.name)) errors.push("workspace.name is required")
    if (!['solo', 'team'].includes(String(workspace.mode))) errors.push("workspace.mode must be solo or team")
    if (!nonEmptyString(workspace.default_branch)) errors.push("workspace.default_branch is required")
  }
  if (!record(value.repositories)) errors.push("repositories is required")
  else {
    for (const [name, repo] of Object.entries(value.repositories)) {
      if (!record(repo)) {
        errors.push(`repositories.${name} must be a mapping`)
        continue
      }
      if (!nonEmptyString(repo.path)) errors.push(`repositories.${name}.path is required`)
      if (!['ignored-clone', 'submodule'].includes(String(repo.mode))) errors.push(`repositories.${name}.mode must be ignored-clone or submodule`)
      if (!nonEmptyString(repo.role)) errors.push(`repositories.${name}.role is required`)
      if (!nonEmptyString(repo.agent)) errors.push(`repositories.${name}.agent is required`)
      if (!nonEmptyString(repo.default_branch)) errors.push(`repositories.${name}.default_branch is required`)
      if (repo.remote !== undefined) {
        const remoteError = typeof repo.remote === 'string' ? remoteReferenceError(repo.remote) : 'must be a string'
        if (remoteError) errors.push(`repositories.${name}.remote ${remoteError}`)
      }
    }
  }
  const sources = record(value.context) ? value.context.sources_file : undefined
  if (sources !== undefined && !nonEmptyString(sources)) errors.push("context.sources_file must be a path")
  return errors
}

export function workspaceSemanticErrors(config: WorkspaceConfig): string[] {
  const errors = [...workspaceErrors(config)]
  const paths = new Map<string, string>()
  for (const [name, repository] of Object.entries(config.repositories ?? {})) {
    const normalized = repository.path.replace(/^\.\//, '').replace(/\/$/, '')
    const prior = paths.get(normalized)
    if (prior) errors.push(`repositories.${name}.path duplicates repositories.${prior}.path`)
    paths.set(normalized, name)
    if (repository.remote) {
      const remoteError = remoteReferenceError(repository.remote)
      if (remoteError) errors.push(`repositories.${name}.remote ${remoteError}`)
    }
  }
  for (const [index, source] of (config.context?.authoritative_sources ?? []).entries()) {
    if (source.repository && !config.repositories[source.repository]) errors.push(`context.authoritative_sources.${index}.repository is not configured: ${source.repository}`)
    const reference = source.location || source.reference
    if (reference) {
      const referenceError = contextReferenceError(reference)
      if (referenceError && !reference.includes('://')) errors.push(`context.authoritative_sources.${index}.location ${referenceError}`)
    }
  }
  return [...new Set(errors)]
}

export async function workspaceDocumentErrors(workspaceRoot: string, config: WorkspaceConfig): Promise<string[]> {
  const required = [
    'README.md', 'AGENTS.md', 'CLAUDE.md', 'WORKFLOW.md', 'workspace.yaml',
    'context/PROJECT.md', 'context/ARCHITECTURE.md', 'context/CONVENTIONS.md',
    'context/DECISIONS.md', 'context/SOURCES.md', 'context/sources.yaml', 'plans/README.md',
    '.agents/bin/cc.mjs',
    ...new Set(Object.values(config.repositories).map((repo) => `agents/${repo.agent}.md`)),
  ]
  const errors: string[] = []
  for (const path of required) {
    try {
      const info = await lstat(resolve(workspaceRoot, path))
      if (!info.isFile() || info.isSymbolicLink()) throw new Error('not a regular file')
    } catch {
      errors.push(`required workspace document is missing: ${path}`)
    }
  }
  return errors
}

export function parseFrontmatter(raw: string): { value: Record<string, unknown> | null; body: string; errors: string[] } {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---(?:\r?\n|$)/)
  if (!match) return { value: null, body: raw, errors: ['Markdown must begin with YAML frontmatter'] }
  try {
    const value = parseYaml(match[1]!)
    if (!record(value)) return { value: null, body: raw.slice(match[0].length), errors: ['frontmatter must contain a YAML mapping'] }
    return { value, body: raw.slice(match[0].length), errors: [] }
  } catch (error) {
    return { value: null, body: raw.slice(match[0].length), errors: [`invalid YAML frontmatter: ${(error as Error).message}`] }
  }
}

export function requiredString(value: Record<string, unknown>, key: string, errors: string[], prefix = ''): string | undefined {
  if (!nonEmptyString(value[key])) errors.push(`${prefix}${key} is required`)
  return typeof value[key] === 'string' ? value[key].trim() : undefined
}

export function optionalStringList(value: Record<string, unknown>, key: string, errors: string[], prefix = ''): string[] {
  if (value[key] === undefined) return []
  if (!stringList(value[key])) {
    errors.push(`${prefix}${key} must be a list of non-empty strings`)
    return []
  }
  return value[key].map((item) => item.trim())
}

export async function referencedLocalFilesExist(root: string, references: string[]): Promise<string[]> {
  const errors: string[] = []
  for (const reference of references) {
    if (reference.includes('*') || reference.includes('<') || reference.includes('>')) continue
    const candidate = resolve(root, reference)
    if (relative(root, candidate).startsWith('..')) {
      errors.push(`referenced path escapes workspace: ${reference}`)
      continue
    }
    try {
      const info = await lstat(candidate)
      if (info.isSymbolicLink()) errors.push(`referenced path is a symlink: ${reference}`)
    } catch {
      errors.push(`referenced local file does not exist: ${reference}`)
    }
  }
  return errors
}

export function repositoryForPath(config: WorkspaceConfig, name: string): RepositoryConfig | undefined {
  return config.repositories?.[name]
}
