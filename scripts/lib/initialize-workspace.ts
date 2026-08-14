import { access, lstat, mkdir, readFile } from "node:fs/promises"
import { dirname, join, relative, resolve } from "node:path"
import { parse as parseYaml, stringify as stringifyYaml } from "yaml"
import { git, assertCleanRepository } from "./git.js"
import { assertInside, writeTextAtomic } from "./io.js"
import type { BootstrapWorkspaceSummary, InitializationRepositorySummary, InitializationSummary, WorkspaceBootstrapRequest, WorkspaceConfig } from "./types.js"
import { workspaceDocumentErrors, workspaceSemanticErrors } from "./validation.js"
import { renderWorkspaceContext } from "./workspace-context.js"
import { renderProductKnowledgeBaseline, validateProductKnowledgeTree } from "./product-knowledge.js"
import { reconcileWorkspaceReadme } from "./workspace-readme.js"
import { cloneReferenceError, remoteReferenceError } from "./safe-reference.js"

const ignoredStart = "# context-circuit:ignored-clones:start"
const ignoredEnd = "# context-circuit:ignored-clones:end"

function pathExists(path: string): Promise<boolean> {
  return access(path).then(() => true).catch(() => false)
}

function normalizedRepositoryPath(path: string): string {
  return `${path.replace(/^\.\//, "").replace(/\/$/, "")}/`
}

function parseSubmodulePaths(raw: string): Set<string> {
  return new Set([...raw.matchAll(/^\s*path\s*=\s*(.+?)\s*$/gm)].map((match) => match[1]!.replace(/\/$/, "")))
}

export function reconcileIgnoredClones(current: string, config: WorkspaceConfig): string {
  const lines = current.replace(/\r\n/g, "\n").split("\n")
  const start = lines.indexOf(ignoredStart)
  const end = lines.indexOf(ignoredEnd)
  if ((start === -1) !== (end === -1) || (start !== -1 && end < start)) throw new Error("Malformed managed ignored-clone block in .gitignore")
  const without = start === -1 ? lines : [...lines.slice(0, start), ...lines.slice(end + 1)]
  if (without.some((line) => ["repositories", "repositories/", "repositories/*", "repositories/**"].includes(line.trim().replace(/^\//, "")))) throw new Error("Blanket repositories/ ignore conflicts with submodule support")
  const configured = new Set(Object.values(config.repositories).map((repo) => normalizedRepositoryPath(repo.path)))
  const retained = without.filter((line) => !configured.has(normalizedRepositoryPath(line.trim())))
  while (retained.at(-1) === "") retained.pop()
  const ignored = Object.values(config.repositories).filter((repo) => repo.mode === "ignored-clone").map((repo) => normalizedRepositoryPath(repo.path)).sort()
  return [...retained, ignoredStart, ...ignored, ignoredEnd, ""].join("\n")
}

async function loadConfig(root: string): Promise<WorkspaceConfig> {
  const config = parseYaml(await readFile(join(root, "workspace.yaml"), "utf8")) as WorkspaceConfig
  const errors = workspaceSemanticErrors(config)
  if (errors.length > 0) throw new Error(`Invalid workspace configuration:\n- ${errors.join("\n- ")}`)
  return config
}

async function repositoryRemote(path: string): Promise<string | null> {
  try { return await git(path, ["remote", "get-url", "origin"]) } catch { return null }
}

async function currentBranch(path: string): Promise<string> {
  try { return await git(path, ["branch", "--show-current"]) } catch { return "unknown" }
}

export async function initializeWorkspace(options: { workspaceRoot: string; apply?: boolean; allowUnbornWrapper?: boolean }): Promise<InitializationSummary> {
  const root = resolve(options.workspaceRoot)
  const config = await loadConfig(root)
  const documentErrors = await workspaceDocumentErrors(root, config)
  const knowledge = await validateProductKnowledgeTree(join(root, "context"))
  const errors = [...documentErrors, ...knowledge.errors.map((error) => `product-knowledge ${error}`)]
  if (errors.length > 0) throw new Error(`Workspace validation failed:\n- ${errors.join("\n- ")}`)
  await git(root, ["rev-parse", "--show-toplevel"])
  const repositories: InitializationRepositorySummary[] = []
  const warnings: string[] = []
  for (const [name, repository] of Object.entries(config.repositories)) {
    const path = assertInside(root, resolve(root, repository.path))
    if (!await pathExists(path)) throw new Error(`Repository ${name} path is not accessible: ${repository.path}`)
    await git(path, ["rev-parse", "--show-toplevel"])
    const changes = (await git(path, ["status", "--porcelain=v1", "--untracked-files=normal"])).split("\n").filter(Boolean)
    const instructions = await pathExists(join(path, "AGENTS.md")) ? join(path, "AGENTS.md") : null
    if (!instructions) warnings.push(`Repository ${name} has no repository-local AGENTS.md`)
    repositories.push({ name, path: relative(root, path) || ".", mode: repository.mode, role: repository.role, agent: repository.agent, default_branch: repository.default_branch, current_branch: await currentBranch(path), clean: changes.length === 0, remote: await repositoryRemote(path), instructions })
  }
  const gitignorePath = join(root, ".gitignore")
  const currentGitignore = await readFile(gitignorePath, "utf8").catch(() => "")
  const nextGitignore = reconcileIgnoredClones(currentGitignore, config)
  if (options.apply !== false && nextGitignore !== currentGitignore) await writeTextAtomic(gitignorePath, nextGitignore)
  return { workspace: config.workspace.name, mode: config.workspace.mode, default_branch: config.workspace.default_branch, repositories, required_documents: [...new Set([...documentErrors.map((error) => error.replace(/^required workspace document is missing: /, "")), ...Object.values(config.repositories).map((repo) => `agents/${repo.agent}.md`)])], wrapper_changes: (await git(root, ["status", "--porcelain=v1", "--untracked-files=all"])).split("\n").filter(Boolean), gitignore_changed: nextGitignore !== currentGitignore, applied: options.apply !== false, warnings }
}

function safeReference(value: string, label: string): void {
  const error = value.includes("://") ? remoteReferenceError(value) : cloneReferenceError(value)
  if (error) throw new Error(`${label} ${error}`)
}

export async function bootstrapWorkspace(options: { workspaceRoot: string; request: WorkspaceBootstrapRequest }): Promise<BootstrapWorkspaceSummary> {
  const root = resolve(options.workspaceRoot)
  const request = options.request
  const config = request.configuration
  const errors = workspaceSemanticErrors(config)
  if (errors.length > 0) throw new Error(`Invalid workspace configuration:\n- ${errors.join("\n- ")}`)
  const hasGit = await pathExists(join(root, ".git"))
  if (request.wrapper.initialize_git && hasGit) throw new Error("Wrapper is already a Git repository; initialize_git must be false")
  if (!request.wrapper.initialize_git && !hasGit) throw new Error("Wrapper is not a Git repository; initialize_git must be true")
  if (!hasGit) await git(root, ["init", "--initial-branch", config.workspace.default_branch])
  else await assertCleanRepository(root)
  const currentReadme = await readFile(join(root, "README.md"), "utf8").catch(() => "")
  const sourceContext = request.context.sources ?? []
  const readmeConfig = config.workspace.purpose ? config : { ...config, workspace: { ...config.workspace, purpose: request.context.project_summary } }
  await writeTextAtomic(join(root, "workspace.yaml"), stringifyYaml(config))
  for (const [path, contents] of Object.entries(renderWorkspaceContext(request.context))) await writeTextAtomic(join(root, path), contents)
  if (request.context.product_knowledge) for (const [path, contents] of Object.entries(renderProductKnowledgeBaseline(request.context.product_knowledge))) await writeTextAtomic(join(root, path), contents)
  await writeTextAtomic(join(root, "README.md"), reconcileWorkspaceReadme(currentReadme, readmeConfig))
  await writeTextAtomic(join(root, ".gitignore"), reconcileIgnoredClones(await readFile(join(root, ".gitignore"), "utf8").catch(() => ".runtime/\n.dist/\nnode_modules/\n"), config))
  await mkdir(join(root, "agents"), { recursive: true })
  for (const [name, repository] of Object.entries(config.repositories)) {
    const path = assertInside(root, resolve(root, repository.path))
    const action = request.repositories.find((candidate) => candidate.name === name)
    if (!action) throw new Error(`Missing setup action for repository: ${name}`)
    if (action.url) safeReference(action.url, `Repository ${name} URL`)
    const agentPath = join(root, "agents", `${repository.agent}.md`)
    if (!await pathExists(agentPath)) await writeTextAtomic(agentPath, `# ${repository.agent}\n\nRead the repository's local instructions and implement the approved plan continuously in its assigned plan worktree.\n`)
    if (action.source === "existing") continue
    if (await pathExists(path)) throw new Error(`Refusing to replace existing repository path: ${repository.path}`)
    await mkdir(dirname(path), { recursive: true })
    if (action.source === "new") {
      await mkdir(path)
      await git(path, ["init", "--initial-branch", repository.default_branch])
      await writeTextAtomic(join(path, ".gitkeep"), "")
      await git(path, ["add", "."])
      await git(path, ["commit", "-m", action.commit_message ?? `Initialize ${name}`])
    } else if (action.source === "clone" && action.url) {
      await git(root, ["clone", "--branch", repository.default_branch, "--single-branch", "--", action.url, path])
    } else if (action.source === "submodule" && action.url) {
      await git(root, ["-c", "protocol.file.allow=always", "submodule", "add", "-b", repository.default_branch, "--", action.url, repository.path])
    }
  }
  const initial = !await git(root, ["rev-parse", "--verify", "HEAD"]).then(() => true).catch(() => false)
  let commit: string | null = null
  if (initial && request.wrapper.authorize_initial_commit !== false) {
    await git(root, ["add", "-A"])
    await git(root, ["commit", "-m", request.wrapper.commit_message ?? "Configure Context Circuit workspace"])
    commit = await git(root, ["rev-parse", "HEAD"])
  }
  const summary = await initializeWorkspace({ workspaceRoot: root })
  return { ...summary, status: "initialized", bootstrap_actions: ["configured workspace and Product Knowledge context"], wrapper_initial_commit: commit }
}
