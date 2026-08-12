import { lstat, mkdir, readFile, realpath } from "node:fs/promises";
import { join, resolve } from "node:path";
import { stringify as stringifyYaml } from "yaml";
import { git } from "./git.js";
import { writeTextAtomic } from "./io.js";
import { bootstrapWorkspace, initializeWorkspace, type BootstrapWorkspaceSummary, type InitializationSummary } from "./initialize-workspace.js";
import type { WorkspaceBootstrapRequest, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";
import { renderWorkspaceContext } from "./workspace-context.js";
import { reconcileWorkspaceReadme } from "./workspace-readme.js";

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

export async function detectWorkspaceConfigurationState(workspaceRoot: string): Promise<WorkspaceConfigurationState> {
  const root = resolve(workspaceRoot);
  if (!await exists(join(root, ".git")) || !await hasHead(root)) return "fresh";
  return "existing";
}

function containsCredential(value: string): boolean {
  return /https?:\/\/[^\s/@:]+:[^\s/@]+@/i.test(value)
    || /(?:token|password|passwd|secret|api[_-]?key)\s*[=:]/i.test(value)
    || /-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/i.test(value);
}

export function workspaceCredentialErrors(request: WorkspaceBootstrapRequest): string[] {
  const values: Array<[string, string | undefined]> = [
    ["configuration.workspace.remote", request.configuration.workspace.remote],
    ...Object.entries(request.configuration.repositories).map(([name, repository]) => [`configuration.repositories.${name}.remote`, repository.remote] as [string, string | undefined]),
    ...request.repositories.map((repository) => [`repositories.${repository.name}.url`, repository.url] as [string, string | undefined]),
    ...(request.context.sources ?? []).map((source, index) => [`context.sources.${index}.reference`, source.reference] as [string, string | undefined]),
  ];
  return values.filter(([, value]) => value && containsCredential(value)).map(([path]) => `${path} appears to contain credentials`);
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

async function reconfigureWorkspace(workspaceRoot: string, request: WorkspaceBootstrapRequest): Promise<InitializationSummary> {
  const root = resolve(workspaceRoot);
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
  const current = await readFile(join(root, "README.md"), "utf8");
  const sources = renderWorkspaceContext(request.context)["context/SOURCES.md"]!;
  await writeTextAtomic(join(root, "workspace.yaml"), stringifyYaml(request.configuration));
  await writeTextAtomic(join(root, "context/SOURCES.md"), sources);
  await mkdir(join(root, "agents"), { recursive: true });
  for (const [name, repository] of Object.entries(request.configuration.repositories)) {
    const path = join(root, "agents", `${repository.agent}.md`);
    if (!await exists(path)) await writeTextAtomic(path, `# ${repository.agent}\n\nFollow \`repository-worker.md\`. This repository owns the ${repository.role} role. Read ${name}'s repository-local instructions before work.\n`);
  }
  await writeTextAtomic(join(root, "README.md"), reconcileWorkspaceReadme(current, request.configuration));
  return initializeWorkspace({ workspaceRoot: root });
}

export async function configureWorkspace(options: { workspaceRoot: string; request?: WorkspaceBootstrapRequest; checkOnly?: boolean }): Promise<ConfigureWorkspaceSummary> {
  const workspaceRoot = resolve(options.workspaceRoot);
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
  const result = await reconfigureWorkspace(workspaceRoot, options.request);
  return { route: "reconfigure", state, message: "Existing wrapper configuration was updated as reviewable, uncommitted changes.", result };
}
