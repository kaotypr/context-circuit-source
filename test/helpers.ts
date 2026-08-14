import { createHash } from "node:crypto"
import { mkdir, mkdtemp, rm, writeFile } from "node:fs/promises"
import { tmpdir } from "node:os"
import { join } from "node:path"
import { git } from "../scripts/lib/git.js"

export function sha256(value: string): string {
  return `sha256:${createHash("sha256").update(value).digest("hex")}`
}

export async function createTestWorkspace(prefix = "context-circuit-test-"): Promise<{ root: string; repository: string; cleanup: () => Promise<void> }> {
  const root = await mkdtemp(join(tmpdir(), prefix))
  const repository = join(root, "repositories", "app")
  await mkdir(repository, { recursive: true })
  await writeFile(join(repository, "README.md"), "# App\n", "utf8")
  await mkdir(join(root, "context"), { recursive: true })
  await writeFile(join(root, "context", "PROJECT.md"), "# Product\n", "utf8")
  await writeFile(join(root, "context", "SOURCES.md"), "# Sources\n", "utf8")
  await writeFile(join(root, "workspace.yaml"), `version: 1
workspace:
  name: test
  mode: solo
  default_branch: main
repositories:
  app:
    path: repositories/app
    mode: ignored-clone
    role: application
    agent: repository-worker
    default_branch: main
workflow:
  human_gates: [plan-approval, task-selection, status-change]
  wrapper_change_policy: direct-commit
`, "utf8")
  await git(repository, ["init", "--initial-branch=main"])
  await git(repository, ["add", "."])
  await git(repository, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "fixture"])
  return { root, repository, cleanup: () => rm(root, { recursive: true, force: true }) }
}
