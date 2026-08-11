import { cp, mkdir, rm, writeFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { git } from "./lib/git.js";

const root = resolve(process.cwd());
const proofRoot = join(root, ".runtime", "cross-host-proof", "20260812-release");
await rm(proofRoot, { recursive: true, force: true });

for (const host of ["codex", "claude"]) {
  const wrapper = join(proofRoot, host);
  await cp(join(root, ".dist", "kao-delivery-workspace-0.1.0"), wrapper, { recursive: true });
  await mkdir(join(wrapper, "repositories"), { recursive: true });
  await cp(join(root, "fixtures", "typescript-api"), join(wrapper, "repositories", "backend"), { recursive: true });
  await cp(join(root, "fixtures", "react-app"), join(wrapper, "repositories", "frontend"), { recursive: true });
  await writeFile(join(wrapper, "agents", "backend.md"), "# Backend agent\n\nOwn the application API and its shared response contracts.\n", "utf8");
  await writeFile(join(wrapper, "workspace.yaml"), `version: 1
template_version: 0.1.0
workspace: { name: cross-host-${host}, mode: team, default_branch: main }
repositories:
  backend: { path: repositories/backend, mode: ignored-clone, role: application-api, agent: backend, default_branch: main }
  frontend: { path: repositories/frontend, mode: ignored-clone, role: web-application, agent: frontend, default_branch: main }
activity: { provider: none, access: auto, required_capabilities: [], optional_capabilities: [] }
workflow: { human_gates: [plan-approval, task-selection, merge], maximum_repair_attempts: 2, wrapper_change_policy: pull-request }
`, "utf8");
  await writeFile(join(wrapper, ".gitignore"), `.runtime/
repositories/backend/
repositories/frontend/
`, "utf8");
  const request = {
    contract_version: 1,
    request: "Version the reset response contract, then expose the verified version in the frontend",
    acceptance_criteria: ["The API returns contract version 1 and the frontend displays Contract v1"],
    shared_contract: { repository: "backend", paths: ["src/contract.ts"] },
    repositories: [
      { name: "backend", depends_on: [], scope: ["src/contract.ts", "src/api.ts"], test_scope: ["src/api.test.ts"], test_policy: "required", verification_commands: ["npm test", "npm run typecheck"], acceptance_criteria: ["ResetCounterResponse requires contractVersion 1 and resetCounter returns it"] },
      { name: "frontend", depends_on: ["backend"], scope: ["src/App.tsx"], test_scope: ["src/App.test.tsx"], test_policy: "required", verification_commands: ["npm test", "npm run build"], acceptance_criteria: ["The counter UI displays Contract v1 based on the verified shared contract"] },
    ],
  };
  await mkdir(join(wrapper, ".runtime", "proof"), { recursive: true });
  await writeFile(join(wrapper, ".runtime", "proof", "request.json"), `${JSON.stringify(request, null, 2)}\n`, "utf8");
  for (const repository of ["backend", "frontend"]) {
    const path = join(wrapper, "repositories", repository);
    await git(path, ["init", "--initial-branch=main"]);
    await git(path, ["add", "."]);
    await git(path, ["-c", "user.name=Proof", "-c", "user.email=proof@example.invalid", "commit", "-m", "fixture"]);
  }
  await git(wrapper, ["init", "--initial-branch=main"]);
  await git(wrapper, ["add", "."]);
  await git(wrapper, ["-c", "user.name=Proof", "-c", "user.email=proof@example.invalid", "commit", "-m", "template"]);
}
console.log(proofRoot);
