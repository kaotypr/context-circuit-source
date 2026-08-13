import { lstat, readFile, readdir, realpath } from "node:fs/promises";
import { dirname, join, relative, resolve } from "node:path";
import { parse as parseYaml, stringify as stringifyYaml } from "yaml";
import { git } from "./git.js";
import { assertInside, writeTextTransaction } from "./io.js";
import { bootstrapWorkspace, initializeWorkspace, reconcileIgnoredClones, type BootstrapWorkspaceSummary, type InitializationSummary } from "./initialize-workspace.js";
import type { WorkspaceBootstrapRequest, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";
import { renderWorkspaceContext } from "./workspace-context.js";
import { reconcileWorkspaceReadme } from "./workspace-readme.js";
import { cloneReferenceError, contextReferenceError, remoteReferenceError } from "./safe-reference.js";

export type WorkspaceConfigurationState = "fresh" | "existing";

export interface ConfigureWorkspaceSummary {
  route: "bootstrap" | "reconfigure" | "inspect-fresh" | "inspect-existing";
  state: WorkspaceConfigurationState;
  message: string;
  result: BootstrapWorkspaceSummary | InitializationSummary | null;
}

async function exists(path: string): Promise<boolean> {
  try { await lstat(path); return true; } catch (error) {
    if ((error as NodeJS.ErrnoException).code === "ENOENT") return false;
    throw error;
  }
}

async function hasHead(root: string): Promise<boolean> {
  try { await git(root, ["rev-parse", "--verify", "HEAD"]); return true; } catch { return false; }
}

const transactionResidue = /\.\d+\.\d+\.[0-9a-f]+\.(?:stage|backup)$/;

/**
 * A process can be interrupted between the sibling renames used by
 * writeTextTransaction. Those siblings may be the only recoverable copy of a
 * managed file, so a later configuration must never silently delete or replace
 * them.
 */
export async function interruptedConfigurationArtifacts(workspaceRoot: string): Promise<string[]> {
  const root = resolve(workspaceRoot);
  const directories = [root, join(root, "context"), join(root, "agents")];
  const artifacts: string[] = [];
  for (const directory of directories) {
    let entries;
    try { entries = await readdir(directory, { withFileTypes: true }); } catch (error) {
      if ((error as NodeJS.ErrnoException).code === "ENOENT") continue;
      throw error;
    }
    for (const entry of entries) {
      if (!entry.isFile() || !transactionResidue.test(entry.name)) continue;
      const original = entry.name.replace(transactionResidue, "");
      const managed = directory === root
        ? ["workspace.yaml", "README.md", ".gitignore"].includes(original)
        : directory === join(root, "context")
          ? original === "SOURCES.md"
          : /^[a-z][a-z0-9-]*\.md$/.test(original);
      if (managed) artifacts.push(relative(root, join(directory, entry.name)).replaceAll("\\", "/"));
    }
  }
  return artifacts.sort();
}

async function assertNoInterruptedConfiguration(workspaceRoot: string): Promise<void> {
  const artifacts = await interruptedConfigurationArtifacts(workspaceRoot);
  if (artifacts.length === 0) return;
  throw new Error(
    `Interrupted workspace configuration artifacts were found:\n- ${artifacts.join("\n- ")}\n` +
    "Configuration will not delete or overwrite them. Inspect each target, .stage, and .backup sibling; restore exactly one authoritative target manually; preserve uncertain copies; then rerun configure-workspace.",
  );
}

export async function detectWorkspaceConfigurationState(workspaceRoot: string): Promise<WorkspaceConfigurationState> {
  const root = resolve(workspaceRoot);
  if (!await exists(join(root, ".git")) || !await hasHead(root)) return "fresh";
  return "existing";
}

export function workspaceCredentialErrors(request: WorkspaceBootstrapRequest): string[] {
  const remotes: Array<[string, string | undefined]> = [
    ["configuration.workspace.remote", request.configuration.workspace.remote],
    ...Object.entries(request.configuration.repositories).map(([name, repository]) => [`configuration.repositories.${name}.remote`, repository.remote] as [string, string | undefined]),
  ];
  const errors: string[] = [];
  for (const [path, value] of remotes) {
    const error = value ? remoteReferenceError(value) : null;
    if (error) errors.push(`${path} ${error}`);
  }
  for (const [index, repository] of request.repositories.entries()) {
    const error = repository.url ? cloneReferenceError(repository.url) : null;
    if (error) errors.push(`repositories.${index}.url ${error}`);
  }
  for (const [index, source] of (request.context.sources ?? []).entries()) {
    const error = contextReferenceError(source.reference);
    if (error) errors.push(`context.sources.${index}.reference ${error}`);
  }
  return errors;
}

function assertSourceConsistency(request: WorkspaceBootstrapRequest): void {
  const configured = request.configuration.context?.authoritative_sources ?? [];
  const requested = request.context.sources ?? [];
  if (JSON.stringify(configured) !== JSON.stringify(requested)) {
    throw new Error("configuration.context.authoritative_sources must exactly match context.sources");
  }
  const names = new Set(Object.keys(request.configuration.repositories));
  for (const [index, source] of requested.entries()) {
    if (source.repository && !names.has(source.repository)) throw new Error(`context.sources.${index}.repository is not configured: ${source.repository}`);
  }
  for (const [name, repository] of Object.entries(request.configuration.repositories)) {
    const action = request.repositories.find((candidate) => candidate.name === name);
    if (repository.remote && action?.url && repository.remote !== action.url) throw new Error(`Configured remote and source URL differ for repository ${name}`);
  }
}

async function validateRequest(request: WorkspaceBootstrapRequest): Promise<void> {
  const contractErrors = await validateContract("workspace-bootstrap-request", request);
  const errors = [
    ...contractErrors.map((error) => `${error.instancePath || "/"} ${error.message}`),
    ...workspaceSemanticErrors(request.configuration),
    ...workspaceCredentialErrors(request),
  ];
  if (errors.length > 0) throw new Error(`Invalid workspace configuration request:\n- ${errors.join("\n- ")}`);
  assertSourceConsistency(request);
}

async function readRegularInside(root: string, path: string, label: string): Promise<string> {
  const candidate = assertInside(root, path);
  const info = await lstat(candidate);
  if (!info.isFile() || info.isSymbolicLink()) throw new Error(`${label} must be a regular non-symlink file`);
  assertInside(await realpath(root), await realpath(candidate));
  return readFile(candidate, "utf8");
}

async function assertExactGitRoot(path: string, label: string): Promise<void> {
  let info;
  try { info = await lstat(path); } catch (error) {
    if ((error as NodeJS.ErrnoException).code === "ENOENT") throw new Error(`${label} path does not exist`);
    throw error;
  }
  if (!info.isDirectory() || info.isSymbolicLink()) throw new Error(`${label} must be a real directory`);
  const top = await git(path, ["rev-parse", "--show-toplevel"]);
  if (await realpath(top) !== await realpath(path)) throw new Error(`${label} is not an exact Git root`);
}

async function assertBranch(path: string, name: string, branch: string): Promise<void> {
  for (const ref of [`refs/heads/${branch}`, `refs/remotes/origin/${branch}`]) {
    try { await git(path, ["rev-parse", "--verify", ref]); return; } catch { /* next */ }
  }
  throw new Error(`Repository ${name} has no local or origin default branch named ${branch}`);
}

async function preflightExistingRepositories(root: string, request: WorkspaceBootstrapRequest): Promise<void> {
  let submodules = "";
  try { submodules = await readFile(join(root, ".gitmodules"), "utf8"); } catch (error) { if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error; }
  for (const [name, repository] of Object.entries(request.configuration.repositories)) {
    const path = assertInside(root, resolve(root, repository.path));
    await assertExactGitRoot(path, `Repository ${name}`);
    await assertBranch(path, name, repository.default_branch);
    const relativePath = relative(root, path).replaceAll("\\", "/");
    const tracked = await git(root, ["ls-files", "--stage", "--", relativePath]);
    const registered = new RegExp(`^\\s*path\\s*=\\s*${relativePath.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}\\s*$`, "m").test(submodules);
    if (repository.mode === "submodule" && (!registered || !tracked.startsWith("160000 "))) throw new Error(`Repository ${name} is not a tracked submodule: ${relativePath}`);
    if (repository.mode === "ignored-clone" && (registered || tracked)) throw new Error(`Repository ${name} is tracked but configured as an ignored clone`);
  }
}

async function exactBootstrapRerun(root: string, request: WorkspaceBootstrapRequest): Promise<boolean> {
  if (request.authorize_reviewable_changes || !request.wrapper.authorize_initial_commit) return false;
  let installed: WorkspaceConfig;
  try { installed = parseYaml(await readFile(join(root, "workspace.yaml"), "utf8")) as WorkspaceConfig; } catch { return false; }
  if (JSON.stringify(installed) !== JSON.stringify(request.configuration)) return false;
  const readme = await readRegularInside(root, join(root, "README.md"), "README.md");
  const sources = await readRegularInside(root, join(root, "context", "SOURCES.md"), "context/SOURCES.md");
  const readmeConfig: WorkspaceConfig = request.configuration.workspace.purpose
    ? request.configuration
    : { ...request.configuration, workspace: { ...request.configuration.workspace, purpose: request.context.project_summary } };
  if (reconcileWorkspaceReadme(readme, readmeConfig) !== readme) return false;
  if (renderWorkspaceContext(request.context)["context/SOURCES.md"] !== sources) return false;
  try {
    await preflightExistingRepositories(root, {
      ...request,
      repositories: request.repositories.map((repository) => ({ name: repository.name, source: "existing", authorize_initial_commit: false })),
    });
  } catch { return false; }
  return true;
}

async function reconfigureWorkspace(workspaceRoot: string, request: WorkspaceBootstrapRequest, transactionOptions: { failRenameAt?: number } = {}): Promise<InitializationSummary> {
  const root = resolve(workspaceRoot);
  const changes = await git(root, ["status", "--porcelain=v1", "--untracked-files=all"]);
  if (changes) throw new Error(`Wrapper must be clean before reconfiguration; refusing to overwrite existing work:\n${changes}`);
  if (request.authorize_reviewable_changes !== true) throw new Error("Existing wrapper reconfiguration requires explicit authorize_reviewable_changes: true");
  if (request.wrapper.initialize_git || request.wrapper.authorize_initial_commit) {
    throw new Error("An existing wrapper must not initialize Git or authorize an initial commit");
  }
  const configuredNames = Object.keys(request.configuration.repositories).sort();
  const actionNames = request.repositories.map((repository) => repository.name).sort();
  if (new Set(actionNames).size !== actionNames.length || configuredNames.join("\n") !== actionNames.join("\n")) {
    throw new Error("Reconfiguration repository actions must match configured repositories exactly");
  }
  for (const repository of request.repositories) {
    if (repository.source !== "existing") throw new Error(`Existing wrapper reconfiguration accepts only inspected existing repository paths: ${repository.name}`);
    if (repository.authorize_initial_commit || repository.commit_message || repository.author_name || repository.author_email) {
      throw new Error(`Existing repository ${repository.name} must not include initial-commit authorization or metadata`);
    }
  }
  const top = await git(root, ["rev-parse", "--show-toplevel"]);
  if (await realpath(top) !== await realpath(root)) throw new Error("Workspace root is not the wrapper Git root");
  const current = await readRegularInside(root, join(root, "README.md"), "README.md");
  await readRegularInside(root, join(root, "context", "SOURCES.md"), "context/SOURCES.md");
  const sources = renderWorkspaceContext(request.context)["context/SOURCES.md"]!;
  const readme = reconcileWorkspaceReadme(current, request.configuration);
  const gitignorePath = join(root, ".gitignore");
  const gitignore = reconcileIgnoredClones(await readRegularInside(root, gitignorePath, ".gitignore"), request.configuration);
  await preflightExistingRepositories(root, request);
  const agentWrites: Array<[string, string]> = [];
  const agentsDirectory = join(root, "agents");
  const agentsInfo = await lstat(agentsDirectory);
  if (!agentsInfo.isDirectory() || agentsInfo.isSymbolicLink()) throw new Error("agents must be a real directory");
  for (const [name, repository] of Object.entries(request.configuration.repositories)) {
    const path = join(root, "agents", `${repository.agent}.md`);
    if (await exists(path)) await readRegularInside(root, path, `agents/${repository.agent}.md`);
    else agentWrites.push([path, `# ${repository.agent}\n\nFollow \`repository-worker.md\`. This repository owns the ${repository.role} role. Read ${name}'s repository-local instructions before work.\n`]);
  }
  const workspace = stringifyYaml(request.configuration);
  // All validation and output computation above is read-only. The transaction stages every sibling before changing any target.
  await writeTextTransaction([
    { path: join(root, "workspace.yaml"), value: workspace },
    { path: join(root, "context/SOURCES.md"), value: sources },
    ...agentWrites.map(([path, value]) => ({ path, value })),
    { path: join(root, "README.md"), value: readme },
    { path: gitignorePath, value: gitignore },
  ], transactionOptions);
  return initializeWorkspace({ workspaceRoot: root });
}

export async function configureWorkspace(options: { workspaceRoot: string; request?: WorkspaceBootstrapRequest; checkOnly?: boolean; transactionOptions?: { failRenameAt?: number } }): Promise<ConfigureWorkspaceSummary> {
  const workspaceRoot = resolve(options.workspaceRoot);
  await assertNoInterruptedConfiguration(workspaceRoot);
  const state = await detectWorkspaceConfigurationState(workspaceRoot);
  if (!options.request) {
    if (state === "fresh") return { route: "inspect-fresh", state, message: "Fresh wrapper detected; collect a configuration request, then run the internal bootstrap phase with exact initial-commit authorization.", result: null };
    const result = await initializeWorkspace({ workspaceRoot, apply: false });
    return { route: "inspect-existing", state, message: "Existing wrapper detected; configuration changes will remain reviewable and uncommitted.", result };
  }
  await validateRequest(options.request);
  if (options.checkOnly) return { route: state === "fresh" ? "inspect-fresh" : "inspect-existing", state, message: "Configuration request is valid; no files or Git state changed.", result: state === "existing" ? await initializeWorkspace({ workspaceRoot, apply: false }) : null };
  if (state === "fresh") {
    if (options.request.authorize_reviewable_changes) throw new Error("Fresh bootstrap must not authorize existing-wrapper reconfiguration");
    const result = await bootstrapWorkspace({ workspaceRoot, request: options.request });
    return { route: "bootstrap", state, message: "Fresh wrapper configured through the explicit bootstrap phase.", result };
  }
  if (await exactBootstrapRerun(workspaceRoot, options.request)) {
    const result = await initializeWorkspace({ workspaceRoot, apply: false });
    return { route: "inspect-existing", state, message: "Exact completed bootstrap request detected; configuration is already current and no files or commits changed.", result };
  }
  const result = await reconfigureWorkspace(workspaceRoot, options.request, options.transactionOptions);
  return { route: "reconfigure", state, message: "Existing wrapper configuration was updated as reviewable, uncommitted changes.", result };
}
