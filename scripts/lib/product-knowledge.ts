import { createHash } from "node:crypto"
import { access, mkdir, readdir, readFile } from "node:fs/promises"
import { dirname, join, relative, resolve } from "node:path"
import { parse as parseYaml, stringify as stringifyYaml } from "yaml"
import { writeTextAtomic } from "./io.js"
import { parseFrontmatter } from "./validation.js"
import type { ProductKnowledgeBaselineSpec, ProductKnowledgeImportRequest, ProductKnowledgeImportResult, ProductKnowledgeRefreshResult, ProductKnowledgeValidationResult, SourceRecord, SourceRegistry, WorkspaceConfig } from "./types.js"

function isRecord(value: unknown): value is Record<string, unknown> {
  return Boolean(value) && typeof value === "object" && !Array.isArray(value)
}

function safeSlug(value: string): string {
  return value.trim().toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "") || "source"
}

function digest(content: string): string {
  return `sha256:${createHash("sha256").update(content).digest("hex")}`
}

async function exists(path: string): Promise<boolean> {
  try {
    await access(path)
    return true
  } catch {
    return false
  }
}

function markdownFiles(directory: string): Promise<string[]> {
  return readdir(directory, { withFileTypes: true }).then((entries) => entries.filter((entry) => entry.isFile() && entry.name.endsWith(".md")).map((entry) => join(directory, entry.name)).sort())
}

function linksInMarkdown(raw: string): string[] {
  return [...raw.matchAll(/\]\(([^)#]+)(?:#[^)]+)?\)/g)].map((match) => match[1]!).filter((link) => !link.startsWith("http"))
}

async function resolveLinks(root: string, page: string, errors: string[]): Promise<void> {
  for (const link of linksInMarkdown(await readFile(page, "utf8"))) {
    const target = resolve(dirname(page), link)
    if (relative(root, target).startsWith("..") || !await exists(target)) errors.push(`${relative(root, page)} references missing page: ${link}`)
  }
}

export async function validateProductKnowledgeTree(contextDirInput: string): Promise<ProductKnowledgeValidationResult> {
  const root = resolve(contextDirInput)
  const pages: string[] = []
  if (await exists(join(root, "PROJECT.md"))) pages.push(join(root, "PROJECT.md"))
  for (const directory of [join(root, "roles"), join(root, "domains")]) {
    if (!await exists(directory)) continue
    const walk = async (current: string): Promise<void> => {
      for (const entry of await readdir(current, { withFileTypes: true })) {
        if (entry.isSymbolicLink()) continue
        const path = join(current, entry.name)
        if (entry.isDirectory()) await walk(path)
        else if (entry.isFile() && entry.name.endsWith(".md") && entry.name !== "README.md") pages.push(path)
        else if (entry.isFile() && entry.name === "README.md" && current !== root) pages.push(path)
      }
    }
    await walk(directory)
  }
  const errors: string[] = []
  for (const page of pages) {
    const raw = await readFile(page, "utf8")
    const parsed = parseFrontmatter(raw)
    if (raw.trimStart().startsWith("---")) errors.push(...parsed.errors.map((error) => `${relative(root, page)}: ${error}`))
    if (parsed.value) {
      const kind = parsed.value.kind
      if (!['product-map', 'role', 'domain', 'workflow'].includes(String(kind))) errors.push(`${relative(root, page)}: kind must be product-map, role, domain, or workflow`)
      if (typeof parsed.value.title !== 'string' || !parsed.value.title.trim()) errors.push(`${relative(root, page)}: title is required`)
    }
    await resolveLinks(root, page, errors)
  }
  return { present: pages.length > 0, pages: pages.length, errors: [...new Set(errors)] }
}

function defaultSourcesPath(root: string): string {
  return join(root, "context", "sources.yaml")
}

export async function readSourceRegistry(workspaceRootInput: string): Promise<SourceRegistry> {
  const root = resolve(workspaceRootInput)
  const path = defaultSourcesPath(root)
  if (!await exists(path)) return { sources: [] }
  const parsed = parseYaml(await readFile(path, "utf8"))
  if (!isRecord(parsed) || !Array.isArray(parsed.sources)) throw new Error("context/sources.yaml must contain a sources list")
  const sources = parsed.sources.filter(isRecord).map((source) => ({
    id: String(source.id ?? ""),
    kind: String(source.kind ?? "other"),
    location: String(source.location ?? source.reference ?? ""),
    revision: String(source.revision ?? ""),
    product_knowledge: Array.isArray(source.product_knowledge) ? source.product_knowledge.filter((item): item is string => typeof item === "string") : [],
    ...(typeof source.imported_at === "string" ? { imported_at: source.imported_at } : {}),
  }))
  return { sources }
}

async function sourceContent(root: string, location: string): Promise<string | null> {
  if (location.startsWith("http://") || location.startsWith("https://")) return null
  const path = resolve(root, location)
  if (relative(root, path).startsWith("..") || !await exists(path)) return null
  return readFile(path, "utf8")
}

export async function staleProductKnowledgeSources(workspaceRootInput: string): Promise<SourceRecord[]> {
  const root = resolve(workspaceRootInput)
  const registry = await readSourceRegistry(root)
  const stale: SourceRecord[] = []
  for (const source of registry.sources) {
    if (source.revision === "unrecorded") continue
    const content = await sourceContent(root, source.location)
    if (content !== null && digest(content) !== source.revision) stale.push(source)
    else if (content === null && !source.location.startsWith("http")) stale.push(source)
  }
  return stale
}

async function writeRegistry(root: string, registry: SourceRegistry): Promise<void> {
  await mkdir(join(root, "context"), { recursive: true })
  await writeTextAtomic(defaultSourcesPath(root), stringifyYaml(registry))
}

function sourceRecord(root: string, request: ProductKnowledgeImportRequest, content: string): SourceRecord {
  const location = request.source.replaceAll("\\", "/")
  const id = request.source_id?.trim() || safeSlug(location.split('/').at(-1)?.replace(/\.[^.]+$/, '') || location)
  return {
    id,
    kind: request.kind?.trim() || "document",
    location: relative(root, resolve(root, location)).replaceAll("\\", "/"),
    revision: digest(content),
    product_knowledge: request.product_knowledge?.length ? [...request.product_knowledge] : ["context/PROJECT.md"],
    imported_at: new Date().toISOString(),
  }
}

export async function importProductKnowledge(workspaceRootInput: string, request: ProductKnowledgeImportRequest): Promise<ProductKnowledgeImportResult> {
  const root = resolve(workspaceRootInput)
  if (!request.source.trim()) throw new Error("Product Knowledge import requires a source path")
  const sourcePath = resolve(root, request.source)
  if (relative(root, sourcePath).startsWith("..")) throw new Error("Product Knowledge source must stay inside the workspace")
  const content = await readFile(sourcePath, "utf8")
  const source = sourceRecord(root, request, content)
  const registry = await readSourceRegistry(root)
  const prior = registry.sources.find((candidate) => candidate.id === source.id)
  registry.sources = [...registry.sources.filter((candidate) => candidate.id !== source.id), source].sort((left, right) => left.id.localeCompare(right.id))
  await writeRegistry(root, registry)
  const updatedPages: string[] = []
  const proposalLines = [`# Product Knowledge import proposal`, ``, `Source: \`${source.location}\``, `Revision: \`${source.revision}\``, ``]
  for (const page of source.product_knowledge) {
    const target = resolve(root, page)
    if (relative(root, target).startsWith("..")) throw new Error(`Product Knowledge target escapes workspace: ${page}`)
    if (!await exists(target)) {
      await mkdir(dirname(target), { recursive: true })
      const title = request.title?.trim() || source.id
      const purpose = request.purpose?.trim() || content.trim().split(/\r?\n/).find((line) => line.trim()) || `Imported from ${source.location}`
      await writeTextAtomic(target, `---\nkind: product-map\ntitle: ${title}\nsources:\n  - ${source.location}\nreview_date: ${new Date().toISOString().slice(0, 10)}\n---\n\n# ${title}\n\n${purpose}\n`)
      updatedPages.push(page)
    } else {
      proposalLines.push(`## ${page}`, ``, `The page already exists. Review the source and propose a source-cited Markdown change; no existing page was overwritten.`, ``)
    }
  }
  return { source, updated_pages: updatedPages, ...(proposalLines.length > 5 ? { proposal: proposalLines.join("\n") } : {}) }
}

export async function refreshProductKnowledge(workspaceRootInput: string): Promise<ProductKnowledgeRefreshResult> {
  const root = resolve(workspaceRootInput)
  const stale = await staleProductKnowledgeSources(root)
  const affectedPages = [...new Set(stale.flatMap((source) => source.product_knowledge))].sort()
  const lines = ['# Product Knowledge refresh proposal', '', stale.length ? 'The following sources changed; review and accept, revise, or reject the proposed updates.' : 'No Product Knowledge source changes were detected.', '']
  for (const source of stale) {
    const content = await sourceContent(root, source.location)
    lines.push(`## ${source.id}`, '', `- Location: \`${source.location}\``, `- Recorded revision: \`${source.revision}\``, `- Current revision: \`${content === null ? 'unavailable' : digest(content)}\``, `- Affected Product Knowledge: ${source.product_knowledge.map((page) => `\`${page}\``).join(', ')}`, '')
  }
  return { stale, proposal: lines.join('\n'), affected_pages: affectedPages }
}

export function renderProductKnowledgeBaseline(spec: ProductKnowledgeBaselineSpec): Record<string, string> {
  const sources = spec.sources.length ? spec.sources.map((source) => `- ${source}`).join("\n") : "- None recorded."
  const roles = spec.roles.length ? spec.roles.map((role) => `- ${role}`).join("\n") : "- None recorded."
  const domains = spec.domains.length ? spec.domains.map((domain) => `- ${domain.name}${domain.workflows?.length ? ` — ${domain.workflows.join(', ')}` : ''}`).join("\n") : "- None recorded."
  const unknowns = spec.unknowns.length ? spec.unknowns.map((unknown) => `- ${unknown}`).join("\n") : "- None recorded."
  return {
    "context/PROJECT.md": `---\nkind: product-map\ntitle: ${spec.title}\nsources:\n${spec.sources.map((source) => `  - ${source}`).join("\n")}\nreview_date: ${spec.review_date}\n---\n\n# ${spec.title}\n\n${spec.purpose}\n\n## Sources\n\n${sources}\n\n## Roles\n\n${roles}\n\n## Domains\n\n${domains}\n\n## Known gaps\n\n${unknowns}\n`,
    "context/sources.yaml": stringifyYaml({ sources: spec.sources.map((source) => ({ id: safeSlug(source), kind: "document", location: source, revision: "unrecorded", product_knowledge: ["context/PROJECT.md"] })) }),
  }
}
