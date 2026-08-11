import { cp, mkdir, mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { git } from "../scripts/lib/git.js";

export const projectRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");

export async function createTestWorkspace(): Promise<{ root: string; repository: string; cleanup: () => Promise<void> }> {
  const root = await mkdtemp(join(tmpdir(), "context-circuit-test-"));
  const repository = join(root, "repositories", "frontend");
  await mkdir(dirname(repository), { recursive: true });
  await cp(join(projectRoot, "fixtures", "react-app"), repository, { recursive: true });
  await mkdir(join(root, ".agents"), { recursive: true });
  await cp(join(projectRoot, ".agents", "contracts"), join(root, ".agents", "contracts"), { recursive: true });
  await mkdir(join(root, "agents"), { recursive: true });
  await writeFile(join(root, "AGENTS.md"), "# Test instructions\n", "utf8");
  await writeFile(join(root, "agents", "frontend.md"), "# Frontend\n", "utf8");
  await writeFile(join(root, "agents", "repository-worker.md"), "# Worker\n", "utf8");
  await writeFile(join(root, "agents", "verifier.md"), "# Verifier\n", "utf8");
  await writeFile(join(root, "workspace.yaml"), `version: 1
template_version: 0.2.0
workspace:
  name: test
  mode: team
  default_branch: main
repositories:
  frontend:
    path: repositories/frontend
    mode: ignored-clone
    role: web-application
    agent: frontend
    default_branch: main
activity:
  provider: none
  access: auto
  required_capabilities: []
  optional_capabilities: []
workflow:
  human_gates: [plan-approval, task-selection, merge]
  maximum_repair_attempts: 2
  wrapper_change_policy: pull-request
`, "utf8");
  await git(repository, ["init", "--initial-branch=main"]);
  await git(repository, ["add", "."]);
  await git(repository, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "fixture"]);
  return {
    root,
    repository,
    cleanup: async () => rm(root, { recursive: true, force: true }),
  };
}

export async function createTwoRepositoryTestWorkspace(): Promise<{ root: string; frontend: string; backend: string; cleanup: () => Promise<void> }> {
  const workspace = await createTestWorkspace();
  const backend = join(workspace.root, "repositories", "backend");
  await cp(join(projectRoot, "fixtures", "typescript-api"), backend, { recursive: true });
  await writeFile(join(workspace.root, "agents", "backend.md"), "# Backend\n", "utf8");
  await writeFile(join(workspace.root, "workspace.yaml"), `version: 1
template_version: 0.2.0
workspace:
  name: test
  mode: team
  default_branch: main
repositories:
  frontend:
    path: repositories/frontend
    mode: ignored-clone
    role: web-application
    agent: frontend
    default_branch: main
  backend:
    path: repositories/backend
    mode: ignored-clone
    role: application-api
    agent: backend
    default_branch: main
activity:
  provider: none
  access: auto
  required_capabilities: []
  optional_capabilities: []
workflow:
  human_gates: [plan-approval, task-selection, merge]
  maximum_repair_attempts: 2
  wrapper_change_policy: pull-request
`, "utf8");
  await git(backend, ["init", "--initial-branch=main"]);
  await git(backend, ["add", "."]);
  await git(backend, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "fixture"]);
  return { root: workspace.root, frontend: workspace.repository, backend, cleanup: workspace.cleanup };
}

export const taskOptions = {
  request: "Add a reset button",
  repository: "frontend",
  acceptanceCriteria: ["Reset returns the count to zero"],
  scope: ["src/App.tsx"],
  verificationCommands: ["npm test"],
};
