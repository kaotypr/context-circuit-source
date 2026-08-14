import { access, readFile } from "node:fs/promises"
import { join, resolve } from "node:path"
import { parse as parseYaml } from "yaml"
import { bootstrapWorkspace, initializeWorkspace } from "./initialize-workspace.js"
import { workspaceSemanticErrors } from "./validation.js"
import type { ConfigureWorkspaceSummary, WorkspaceBootstrapRequest, WorkspaceConfigurationState, WorkspaceConfig } from "./types.js"

async function exists(path: string): Promise<boolean> {
  try { await access(path); return true } catch { return false }
}

export async function detectWorkspaceConfigurationState(workspaceRootInput: string): Promise<WorkspaceConfigurationState> {
  const root = resolve(workspaceRootInput)
  return await exists(join(root, ".git")) ? "existing" : "fresh"
}

export async function configureWorkspace(options: { workspaceRoot: string; request?: WorkspaceBootstrapRequest; checkOnly?: boolean }): Promise<ConfigureWorkspaceSummary> {
  const root = resolve(options.workspaceRoot)
  const state = await detectWorkspaceConfigurationState(root)
  if (!options.request) {
    if (state === "fresh") return { route: "inspect-fresh", state, message: "Fresh wrapper detected; collect a setup request before writing workspace files.", result: null }
    return { route: "inspect-existing", state, message: "Existing wrapper detected; reporting current workspace state.", result: await initializeWorkspace({ workspaceRoot: root, apply: false }) }
  }
  const errors = workspaceSemanticErrors(options.request.configuration)
  if (errors.length > 0) throw new Error(`Invalid workspace configuration:\n- ${errors.join("\n- ")}`)
  if (options.checkOnly) return { route: state === "fresh" ? "inspect-fresh" : "inspect-existing", state, message: "Setup request is valid; no files changed.", result: state === "existing" ? await initializeWorkspace({ workspaceRoot: root, apply: false }) : null }
  if (state === "fresh") return { route: "bootstrap", state, message: "Configured a fresh Context Circuit workspace.", result: await bootstrapWorkspace({ workspaceRoot: root, request: options.request }) }
  const current = parseYaml(await readFile(join(root, "workspace.yaml"), "utf8")) as WorkspaceConfig
  if (JSON.stringify(current) === JSON.stringify(options.request.configuration)) return { route: "inspect-existing", state, message: "Workspace configuration is already current.", result: await initializeWorkspace({ workspaceRoot: root, apply: false }) }
  const summary = await bootstrapWorkspace({ workspaceRoot: root, request: { ...options.request, wrapper: { ...options.request.wrapper, initialize_git: false, authorize_initial_commit: false }, repositories: options.request.repositories.map((repository) => ({ name: repository.name, source: "existing" as const, authorize_initial_commit: false })) } })
  return { route: "reconfigure", state, message: "Updated workspace configuration for human review.", result: summary }
}
