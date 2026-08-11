import { access, readFile, realpath } from "node:fs/promises";
import { join, relative, resolve } from "node:path";
import { assertInside, writeTextAtomic } from "./io.js";
import { git } from "./git.js";
import type { WorkspaceConfig } from "./types.js";
import { readData, requiredWorkspaceDocuments, validateContract, workspaceDocumentErrors, workspaceSemanticErrors } from "./validation.js";

const ignoredStart = "# kao-delivery-workspace:ignored-clones:start";
const ignoredEnd = "# kao-delivery-workspace:ignored-clones:end";

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

export async function initializeWorkspace(options: InitializeWorkspaceOptions): Promise<InitializationSummary> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const configPath = join(workspaceRoot, "workspace.yaml");
  const config = await readData(configPath) as WorkspaceConfig;
  const contractErrors = await validateContract("workspace", config);
  const errors = [
    ...contractErrors.map((error) => `${error.instancePath || "/"} ${error.message}`),
    ...workspaceSemanticErrors(config),
    ...await workspaceDocumentErrors(workspaceRoot, config),
  ];
  if (errors.length > 0) throw new Error(`Workspace initialization validation failed:\n- ${errors.join("\n- ")}`);

  await git(workspaceRoot, ["rev-parse", "--is-inside-work-tree"]);
  const wrapperTopLevel = await git(workspaceRoot, ["rev-parse", "--show-toplevel"]);
  if (await realpath(wrapperTopLevel) !== await realpath(workspaceRoot)) {
    throw new Error(`Workspace root is not the wrapper Git root: ${workspaceRoot}`);
  }
  await assertDefaultBranch(workspaceRoot, "wrapper", config.workspace.default_branch);
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
