import { cp, mkdir, mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { git } from "../scripts/lib/git.js";

export const projectRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");

export async function createTestWorkspace(): Promise<{ root: string; repository: string; cleanup: () => Promise<void> }> {
  const root = await mkdtemp(join(tmpdir(), "kao-delivery-test-"));
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
template_version: 0.1.0
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

export const taskOptions = {
  request: "Add a reset button",
  repository: "frontend",
  acceptanceCriteria: ["Reset returns the count to zero"],
  scope: ["src/App.tsx"],
  verificationCommands: ["npm test"],
};
