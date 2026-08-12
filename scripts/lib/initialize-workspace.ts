import { createHash } from "node:crypto";
import { access, lstat, mkdir, readFile, realpath } from "node:fs/promises";
import { dirname, join, relative, resolve } from "node:path";
import { stringify as stringifyYaml } from "yaml";
import { assertInside, writeTextAtomic } from "./io.js";
import { git } from "./git.js";
import type { BootstrapGitCommit, WorkspaceBootstrapRequest, WorkspaceConfig } from "./types.js";
import { readData, requiredWorkspaceDocuments, validateContract, workspaceDocumentErrors, workspaceSemanticErrors } from "./validation.js";
import { renderWorkspaceContext } from "./workspace-context.js";
import { renderProductKnowledgeBaseline, validateProductKnowledgeTree } from "./product-knowledge.js";
import { reconcileWorkspaceReadme } from "./workspace-readme.js";
import { cloneReferenceError } from "./safe-reference.js";

const ignoredStart = "# context-circuit:ignored-clones:start";
const ignoredEnd = "# context-circuit:ignored-clones:end";
declare const __CC_TEMPLATE_INVENTORY__: string[] | undefined;

export interface InitializationRepositorySummary {
  name: string;
  path: string;
  mode: "ignored-clone" | "submodule";
  role: string;
  agent: string;
  default_branch: string;
  current_branch: string;
  clean: boolean;
  remote: string | null;
  instructions: string | null;
}

export interface InitializationSummary {
  workspace: string;
  mode: "solo" | "team";
  default_branch: string;
  activity_provider: string;
  wrapper_change_policy: "pull-request" | "direct-commit";
  repositories: InitializationRepositorySummary[];
  required_documents: string[];
  wrapper_changes: string[];
  gitignore_changed: boolean;
  applied: boolean;
  warnings: string[];
}

export interface InitializeWorkspaceOptions {
  workspaceRoot: string;
  apply?: boolean;
  allowUnbornWrapper?: boolean;
}

export interface BootstrapWorkspaceOptions {
  workspaceRoot: string;
  request: WorkspaceBootstrapRequest;
}

export interface BootstrapWorkspaceSummary extends InitializationSummary {
  status: "initialized";
  bootstrap_actions: string[];
  wrapper_initial_commit: string | null;
}

function normalizedRepositoryPath(path: string): string {
  return `${path.replace(/^\.\//, "").replace(/\/$/, "")}/`;
}

function parseSubmodulePaths(raw: string): Set<string> {
  const paths = new Set<string>();
  for (const match of raw.matchAll(/^\s*path\s*=\s*(.+?)\s*$/gm)) paths.add(match[1]!.replace(/\/$/, ""));
  return paths;
}

function dangerousRepositoryIgnore(line: string): boolean {
  const normalized = line.trim().replace(/^\//, "");
  return ["repositories", "repositories/", "repositories/*", "repositories/**"].includes(normalized);
}

export function reconcileIgnoredClones(current: string, config: WorkspaceConfig): string {
  const lines = current.replace(/\r\n/g, "\n").split("\n");
  const start = lines.indexOf(ignoredStart);
  const end = lines.indexOf(ignoredEnd);
  if ((start === -1) !== (end === -1) || (start !== -1 && end < start)) {
    throw new Error("Malformed managed ignored-clone block in .gitignore");
  }
  const withoutManaged = start === -1 ? [...lines] : [...lines.slice(0, start), ...lines.slice(end + 1)];
  if (withoutManaged.some(dangerousRepositoryIgnore)) {
    throw new Error("Blanket repositories/ ignore conflicts with submodule support; use exact ignored-clone paths");
  }
  const configuredPaths = new Set(Object.values(config.repositories).map((repository) => normalizedRepositoryPath(repository.path)));
  const retained = withoutManaged.filter((line) => !configuredPaths.has(normalizedRepositoryPath(line.trim())));
  while (retained.length > 0 && retained.at(-1) === "") retained.pop();
  const ignored = Object.values(config.repositories)
    .filter((repository) => repository.mode === "ignored-clone")
    .map((repository) => normalizedRepositoryPath(repository.path))
    .sort();
  const managed = [ignoredStart, ...ignored, ignoredEnd];
  return `${[...retained, ...(retained.length > 0 ? [""] : []), ...managed].join("\n")}\n`;
}

async function repositoryRemote(path: string): Promise<string | null> {
  try {
    return await git(path, ["remote", "get-url", "origin"]);
  } catch {
    return null;
  }
}

async function assertDefaultBranch(path: string, repository: string, branch: string): Promise<void> {
  for (const ref of [`refs/heads/${branch}`, `refs/remotes/origin/${branch}`]) {
    try {
      await git(path, ["rev-parse", "--verify", ref]);
      return;
    } catch {
      // Try the next conventional local or origin reference.
    }
  }
  throw new Error(`Repository ${repository} has no local or origin default branch named ${branch}`);
}

async function pathExists(path: string): Promise<boolean> {
  try {
    await lstat(path);
    return true;
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code === "ENOENT") return false;
    throw error;
  }
}

async function assertSafeRepositoryPath(workspaceRoot: string, path: string, name: string): Promise<void> {
  if (path === workspaceRoot) throw new Error(`Repository ${name} path cannot be the wrapper root`);
  let ancestor = dirname(path);
  while (!await pathExists(ancestor)) {
    const parent = dirname(ancestor);
    if (parent === ancestor) throw new Error(`Cannot resolve repository parent for ${name}`);
    ancestor = parent;
  }
  const info = await lstat(ancestor);
  if (info.isSymbolicLink()) throw new Error(`Repository ${name} parent cannot be a symbolic link`);
  assertInside(await realpath(workspaceRoot), await realpath(ancestor));
}

async function hasHead(path: string): Promise<boolean> {
  try {
    await git(path, ["rev-parse", "--verify", "HEAD"]);
    return true;
  } catch {
    return false;
  }
}

function safeRemote(value: string, repository: string): string {
  const remote = value.trim();
  const error = cloneReferenceError(remote);
  if (error) throw new Error(`Repository ${repository} clone URL ${error}`);
  return remote;
}

async function assertExpectedUnbornTemplate(root: string): Promise<void> {
  const status = (await git(root, ["status", "--porcelain=v1", "--untracked-files=all"])).split("\n").filter(Boolean);
  const allowed = new Set<string>([...requiredWorkspaceDocuments, ".gitignore", "template-manifest.json"]);
  const templateDirectories = [".agents/", ".codex/", ".claude/", "agents/", "context/", "contributions/", "docs/"];
  const trustedInventory = typeof __CC_TEMPLATE_INVENTORY__ === "undefined" ? null : __CC_TEMPLATE_INVENTORY__;
  if (trustedInventory) {
    const manifestPath = join(root, "template-manifest.json");
    const manifest = JSON.parse(await readFile(manifestPath, "utf8")) as Record<string, unknown>;
    const inventory = manifest.file_inventory;
    if (manifest.name !== "context-circuit" || manifest.version !== "0.2.1" || manifest.node !== ">=22" || manifest.command !== "node .agents/bin/cc.mjs"
      || !Array.isArray(inventory) || JSON.stringify(inventory) !== JSON.stringify(trustedInventory)) {
      throw new Error("Extracted template manifest or inventory has been modified");
    }
    const expectedBundle = createHash("sha256").update(await readFile(join(root, ".agents", "bin", "cc.mjs"))).digest("hex");
    if (manifest.bundle_sha256 !== expectedBundle) throw new Error("Extracted template manifest bundle digest has been modified");
    allowed.clear();
    for (const path of trustedInventory) allowed.add(path);
    templateDirectories.length = 0;
  }
  const unexpected = status.filter((line) => {
    if (!line.startsWith("?? ")) return true;
    const path = line.slice(3);
    return !allowed.has(path) && !templateDirectories.some((prefix) => path.startsWith(prefix));
  });
  if (unexpected.length > 0) throw new Error(`Unborn wrapper contains authored or unexpected changes; refusing bootstrap:\n${unexpected.join("\n")}`);
  const config = await readData(join(root, "workspace.yaml")) as WorkspaceConfig;
  if (config.workspace.name !== "uninitialized-workspace" || Object.keys(config.repositories).length !== 0) {
    throw new Error("Unborn wrapper is not the neutral extracted-template baseline");
  }
}

function commitArgs(commit: BootstrapGitCommit): string[] {
  if (Boolean(commit.author_name) !== Boolean(commit.author_email)) throw new Error("Commit author_name and author_email must be supplied together");
  const args: string[] = [];
  if (commit.author_name && commit.author_email) args.push("-c", `user.name=${commit.author_name}`, "-c", `user.email=${commit.author_email}`);
  return [...args, "commit", "--allow-empty", "-m", commit.commit_message.trim()];
}

function agentDocument(name: string, role: string): string {
  const title = name.split("-").map((part) => `${part.charAt(0).toUpperCase()}${part.slice(1)}`).join(" ");
  return `# ${title} worker\n\nFollow \`repository-worker.md\`. This repository owns the ${role} role. Read its repository-local instructions and preserve its established architecture, conventions, and verification commands.\n`;
}

async function assertExactGitRoot(path: string, name: string): Promise<void> {
  const topLevel = await git(path, ["rev-parse", "--show-toplevel"]);
  if (await realpath(topLevel) !== await realpath(path)) throw new Error(`Repository path is not a Git root: ${name}`);
}

export async function bootstrapWorkspace(options: BootstrapWorkspaceOptions): Promise<BootstrapWorkspaceSummary> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const requestErrors = await validateContract("workspace-bootstrap-request", options.request);
  if (requestErrors.length > 0) throw new Error(`Invalid workspace-bootstrap-request: ${requestErrors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  const config = options.request.configuration;
  const semanticErrors = workspaceSemanticErrors(config);
  if (semanticErrors.length > 0) throw new Error(`Invalid workspace configuration: ${semanticErrors.join("; ")}`);
  if (config.workspace.name === "uninitialized-workspace") throw new Error("Bootstrap requires a human-selected workspace name");
  const configuredNames = Object.keys(config.repositories).sort();
  if (configuredNames.length === 0) throw new Error("Bootstrap requires at least one configured repository");
  const actionsByName = new Map(options.request.repositories.map((action) => [action.name, action]));
  if (actionsByName.size !== options.request.repositories.length || configuredNames.join("\n") !== [...actionsByName.keys()].sort().join("\n")) {
    throw new Error("Bootstrap repository actions must match configured repositories exactly");
  }

  const wrapperGitExists = await pathExists(join(workspaceRoot, ".git"));
  if (wrapperGitExists === options.request.wrapper.initialize_git) {
    throw new Error(wrapperGitExists ? "Wrapper is already a Git repository; initialize_git must be false" : "Wrapper is not a Git repository; initialize_git must be true");
  }
  if (wrapperGitExists) {
    await assertExactGitRoot(workspaceRoot, "wrapper");
    if (!await hasHead(workspaceRoot)) await assertExpectedUnbornTemplate(workspaceRoot);
    else {
      const changes = await git(workspaceRoot, ["status", "--porcelain=v1", "--untracked-files=normal"]);
      if (changes) throw new Error(`Wrapper has existing changes; refusing bootstrap:\n${changes}`);
    }
  }
  const wrapperHadHead = wrapperGitExists && await hasHead(workspaceRoot);
  if (!wrapperHadHead && !options.request.wrapper.authorize_initial_commit) throw new Error("A new or unborn wrapper requires explicit initial-commit authorization");
  if (wrapperHadHead && options.request.wrapper.authorize_initial_commit) throw new Error("An existing wrapper must not authorize another initial commit");
  const gitignorePath = join(workspaceRoot, ".gitignore");
  const currentGitignore = await readFile(gitignorePath, "utf8");
  const nextGitignore = reconcileIgnoredClones(currentGitignore, config);
  const readmePath = join(workspaceRoot, "README.md");
  const readmeInfo = await lstat(readmePath);
  if (!readmeInfo.isFile() || readmeInfo.isSymbolicLink()) throw new Error("README.md must be a regular non-symlink file");
  assertInside(await realpath(workspaceRoot), await realpath(readmePath));
  const currentReadme = await readFile(readmePath, "utf8");
  const sourcesPath = join(workspaceRoot, "context", "SOURCES.md");
  const sourcesInfo = await lstat(sourcesPath);
  if (!sourcesInfo.isFile() || sourcesInfo.isSymbolicLink()) throw new Error("context/SOURCES.md must be a regular non-symlink file");
  assertInside(await realpath(workspaceRoot), await realpath(sourcesPath));
  const readmeConfig: WorkspaceConfig = config.workspace.purpose
    ? config
    : { ...config, workspace: { ...config.workspace, purpose: options.request.context.project_summary } };
  const nextReadme = reconcileWorkspaceReadme(currentReadme, readmeConfig);
  for (const agent of new Set(Object.values(config.repositories).map((repository) => repository.agent))) {
    const agentPath = join(workspaceRoot, "agents", `${agent}.md`);
    if (!await pathExists(agentPath)) continue;
    const info = await lstat(agentPath);
    if (!info.isFile() || info.isSymbolicLink()) throw new Error(`Domain agent path must be a regular file: agents/${agent}.md`);
  }

  for (const name of configuredNames) {
    const repository = config.repositories[name]!;
    const action = actionsByName.get(name)!;
    const path = assertInside(workspaceRoot, resolve(workspaceRoot, repository.path));
    await assertSafeRepositoryPath(workspaceRoot, path, name);
    const exists = await pathExists(path);
    if (action.source === "existing") {
      if (repository.mode !== "ignored-clone") throw new Error(`Existing repository ${name} must use ignored-clone mode`);
      if (!exists) throw new Error(`Existing repository path is not accessible: ${repository.path}`);
      if ((await lstat(path)).isSymbolicLink()) throw new Error(`Existing repository ${name} cannot be a symbolic link`);
      assertInside(await realpath(workspaceRoot), await realpath(path));
      await assertExactGitRoot(path, name);
    } else {
      if (exists) throw new Error(`Bootstrap refuses to replace existing path for ${name}: ${repository.path}`);
    }
    if (action.source === "submodule" && repository.mode !== "submodule") throw new Error(`Submodule action requires submodule mode for ${name}`);
    if ((action.source === "new" || action.source === "clone") && repository.mode !== "ignored-clone") throw new Error(`${action.source} action requires ignored-clone mode for ${name}`);
    if ((action.source === "clone" || action.source === "submodule") && !action.url) throw new Error(`${action.source} action requires a URL for ${name}`);
    if ((action.source === "new" || action.source === "existing") && action.url) throw new Error(`${action.source} repository ${name} must not include a clone URL`);
    if (action.url) safeRemote(action.url, name);
    if (action.source === "new" && !action.authorize_initial_commit) throw new Error(`New repository ${name} requires explicit initial-commit authorization`);
    if (action.source === "new" && !action.commit_message) throw new Error(`New repository ${name} requires an initial commit message`);
    if (action.source !== "new" && action.authorize_initial_commit) throw new Error(`${action.source} repository ${name} must not authorize an initial commit`);
    if (action.source !== "new" && (action.commit_message || action.author_name || action.author_email)) throw new Error(`${action.source} repository ${name} must not include artificial commit metadata`);
  }

  const bootstrapActions: string[] = [];
  if (!wrapperGitExists) {
    await git(workspaceRoot, ["init", "--initial-branch", config.workspace.default_branch]);
    bootstrapActions.push(`initialized wrapper Git repository on ${config.workspace.default_branch}`);
  } else if (!wrapperHadHead) {
    const current = await git(workspaceRoot, ["symbolic-ref", "--short", "HEAD"]);
    if (current !== config.workspace.default_branch) throw new Error(`Unborn wrapper branch is ${current}, expected ${config.workspace.default_branch}`);
  }

  await writeTextAtomic(join(workspaceRoot, "workspace.yaml"), stringifyYaml(config));
  for (const [path, contents] of Object.entries(renderWorkspaceContext(options.request.context))) {
    await writeTextAtomic(join(workspaceRoot, path), contents);
  }
  if (options.request.context.product_knowledge) {
    for (const [path, contents] of Object.entries(renderProductKnowledgeBaseline(options.request.context.product_knowledge))) {
      const full = assertInside(workspaceRoot, resolve(workspaceRoot, path));
      await mkdir(dirname(full), { recursive: true });
      await writeTextAtomic(full, contents);
    }
  }
  await writeTextAtomic(readmePath, nextReadme);
  await mkdir(join(workspaceRoot, "agents"), { recursive: true });
  for (const [name, repository] of Object.entries(config.repositories)) {
    const agentPath = join(workspaceRoot, "agents", `${repository.agent}.md`);
    if (!await pathExists(agentPath)) await writeTextAtomic(agentPath, agentDocument(name, repository.role));
  }
  await writeTextAtomic(gitignorePath, nextGitignore);

  for (const name of configuredNames) {
    const repository = config.repositories[name]!;
    const action = actionsByName.get(name)!;
    const path = assertInside(workspaceRoot, resolve(workspaceRoot, repository.path));
    if (action.source === "new") {
      await mkdir(dirname(path), { recursive: true });
      await mkdir(path);
      await git(path, ["init", "--initial-branch", repository.default_branch]);
      await git(path, commitArgs(action as BootstrapGitCommit));
      bootstrapActions.push(`created ${name} with an empty base commit`);
    } else if (action.source === "clone") {
      await mkdir(dirname(path), { recursive: true });
      await git(workspaceRoot, ["clone", "--branch", repository.default_branch, "--single-branch", "--", safeRemote(action.url!, name), path]);
      bootstrapActions.push(`cloned ${name} into ${repository.path}`);
    } else if (action.source === "submodule") {
      await mkdir(dirname(path), { recursive: true });
      await git(workspaceRoot, ["-c", "protocol.file.allow=always", "submodule", "add", "-b", repository.default_branch, "--", safeRemote(action.url!, name), repository.path]);
      bootstrapActions.push(`registered ${name} as a submodule`);
    } else {
      bootstrapActions.push(`registered existing repository ${name}`);
    }
  }

  await initializeWorkspace({ workspaceRoot, allowUnbornWrapper: !wrapperHadHead });
  let wrapperInitialCommit: string | null = null;
  if (!wrapperHadHead) {
    await git(workspaceRoot, ["add", "-A"]);
    await git(workspaceRoot, commitArgs(options.request.wrapper));
    wrapperInitialCommit = await git(workspaceRoot, ["rev-parse", "HEAD"]);
    bootstrapActions.push("created configured wrapper initial commit");
  }
  const summary = await initializeWorkspace({ workspaceRoot });
  return { ...summary, status: "initialized", bootstrap_actions: bootstrapActions, wrapper_initial_commit: wrapperInitialCommit };
}

export async function initializeWorkspace(options: InitializeWorkspaceOptions): Promise<InitializationSummary> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const configPath = join(workspaceRoot, "workspace.yaml");
  const config = await readData(configPath) as WorkspaceConfig;
  const contractErrors = await validateContract("workspace", config);
  const productKnowledge = await validateProductKnowledgeTree(join(workspaceRoot, "context"));
  const errors = [
    ...contractErrors.map((error) => `${error.instancePath || "/"} ${error.message}`),
    ...workspaceSemanticErrors(config),
    ...await workspaceDocumentErrors(workspaceRoot, config),
    ...productKnowledge.errors.map((error) => `product-knowledge ${error}`),
  ];
  if (errors.length > 0) throw new Error(`Workspace initialization validation failed:\n- ${errors.join("\n- ")}`);

  await git(workspaceRoot, ["rev-parse", "--is-inside-work-tree"]);
  const wrapperTopLevel = await git(workspaceRoot, ["rev-parse", "--show-toplevel"]);
  if (await realpath(wrapperTopLevel) !== await realpath(workspaceRoot)) {
    throw new Error(`Workspace root is not the wrapper Git root: ${workspaceRoot}`);
  }
  if (options.allowUnbornWrapper && !await hasHead(workspaceRoot)) {
    const current = await git(workspaceRoot, ["symbolic-ref", "--short", "HEAD"]);
    if (current !== config.workspace.default_branch) throw new Error(`Wrapper branch is ${current}, expected ${config.workspace.default_branch}`);
  } else {
    await assertDefaultBranch(workspaceRoot, "wrapper", config.workspace.default_branch);
  }
  let submodulePaths = new Set<string>();
  try {
    submodulePaths = parseSubmodulePaths(await readFile(join(workspaceRoot, ".gitmodules"), "utf8"));
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
  }

  const repositories: InitializationRepositorySummary[] = [];
  const warnings: string[] = [];
  for (const [name, repository] of Object.entries(config.repositories)) {
    const path = assertInside(workspaceRoot, resolve(workspaceRoot, repository.path));
    try {
      await access(path);
    } catch {
      throw new Error(`Repository ${name} path is not accessible: ${repository.path}`);
    }
    assertInside(await realpath(workspaceRoot), await realpath(path));
    const topLevel = await git(path, ["rev-parse", "--show-toplevel"]);
    if (await realpath(topLevel) !== await realpath(path)) throw new Error(`Repository path is not a Git root: ${repository.path}`);
    const relativePath = relative(workspaceRoot, path).replaceAll("\\", "/");
    const trackedEntry = await git(workspaceRoot, ["ls-files", "--stage", "--", relativePath]);
    if (repository.mode === "submodule" && !submodulePaths.has(relativePath)) {
      throw new Error(`Repository ${name} is configured as a submodule but is not registered in .gitmodules: ${relativePath}`);
    }
    if (repository.mode === "submodule" && !trackedEntry.startsWith("160000 ")) {
      throw new Error(`Repository ${name} is configured as a submodule but the wrapper does not track a gitlink: ${relativePath}`);
    }
    if (repository.mode === "ignored-clone" && submodulePaths.has(relativePath)) {
      throw new Error(`Repository ${name} is registered as a submodule but configured as an ignored clone`);
    }
    if (repository.mode === "ignored-clone" && trackedEntry) {
      throw new Error(`Repository ${name} is tracked by the wrapper but configured as an ignored clone`);
    }
    await assertDefaultBranch(path, name, repository.default_branch);
    const instructionsPath = join(path, "AGENTS.md");
    let instructions: string | null = null;
    try {
      await access(instructionsPath);
      instructions = instructionsPath;
    } catch {
      warnings.push(`Repository ${name} has no repository-local AGENTS.md`);
    }
    const remote = await repositoryRemote(path);
    if (!remote) warnings.push(`Repository ${name} has no origin remote`);
    const currentBranch = await git(path, ["rev-parse", "--abbrev-ref", "HEAD"]);
    const repositoryChanges = (await git(path, ["status", "--porcelain=v1", "--untracked-files=normal"])).split("\n").filter(Boolean);
    repositories.push({
      name,
      path: relativePath,
      mode: repository.mode,
      role: repository.role,
      agent: repository.agent,
      default_branch: repository.default_branch,
      current_branch: currentBranch,
      clean: repositoryChanges.length === 0,
      remote,
      instructions,
    });
  }
  if (config.activity.provider === "none") warnings.push("No activity provider is configured; planless work remains available");

  const gitignorePath = join(workspaceRoot, ".gitignore");
  let currentGitignore = "";
  try {
    currentGitignore = await readFile(gitignorePath, "utf8");
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
  }
  const nextGitignore = reconcileIgnoredClones(currentGitignore, config);
  const gitignoreChanged = nextGitignore !== currentGitignore;
  if (options.apply !== false && gitignoreChanged) await writeTextAtomic(gitignorePath, nextGitignore);
  const wrapperChanges = (await git(workspaceRoot, ["status", "--porcelain=v1", "--untracked-files=all"])).split("\n").filter(Boolean);

  return {
    workspace: config.workspace.name,
    mode: config.workspace.mode,
    default_branch: config.workspace.default_branch,
    activity_provider: config.activity.provider,
    wrapper_change_policy: config.workflow.wrapper_change_policy,
    repositories,
    required_documents: [...new Set([
      ...requiredWorkspaceDocuments,
      ...new Set(Object.values(config.repositories).map((repository) => `agents/${repository.agent}.md`)),
    ])],
    wrapper_changes: wrapperChanges,
    gitignore_changed: gitignoreChanged,
    applied: options.apply !== false,
    warnings,
  };
}
