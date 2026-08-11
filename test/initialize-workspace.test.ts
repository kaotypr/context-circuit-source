import assert from "node:assert/strict";
import { access, mkdir, readFile, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import test from "node:test";
import { stringify as stringifyYaml } from "yaml";
import { git } from "../scripts/lib/git.js";
import { initializeWorkspace, reconcileIgnoredClones } from "../scripts/lib/initialize-workspace.js";
import type { WorkspaceConfig } from "../scripts/lib/types.js";
import { requiredWorkspaceDocuments } from "../scripts/lib/validation.js";
import { createTestWorkspace } from "./helpers.js";

async function ensureFile(path: string, contents = "# Initialization fixture\n"): Promise<void> {
  try {
    await access(path);
  } catch {
    await mkdir(dirname(path), { recursive: true });
    await writeFile(path, contents, "utf8");
  }
}

async function prepareWrapper(): Promise<Awaited<ReturnType<typeof createTestWorkspace>>> {
  const workspace = await createTestWorkspace();
  for (const path of requiredWorkspaceDocuments) await ensureFile(join(workspace.root, path));
  await ensureFile(join(workspace.root, ".gitignore"), "custom.log\nrepositories/frontend/\n");
  await git(workspace.root, ["init", "--initial-branch=main"]);
  await git(workspace.root, ["add", "."]);
  await git(workspace.root, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "wrapper fixture"]);
  return workspace;
}

test("initialization reconciles an exact ignored-clone block idempotently", async (t) => {
  const workspace = await prepareWrapper();
  t.after(workspace.cleanup);

  const checked = await initializeWorkspace({ workspaceRoot: workspace.root, apply: false });
  assert.equal(checked.gitignore_changed, true);
  assert.equal(checked.applied, false);
  assert.equal(checked.repositories[0]?.mode, "ignored-clone");
  assert.equal(checked.repositories[0]?.current_branch, "main");
  assert.equal(checked.repositories[0]?.clean, true);
  assert.ok(checked.repositories[0]?.instructions?.endsWith("/AGENTS.md"));
  assert.ok(checked.required_documents.includes("agents/frontend.md"));
  assert.deepEqual(checked.wrapper_changes, []);
  assert.equal(await readFile(join(workspace.root, ".gitignore"), "utf8"), "custom.log\nrepositories/frontend/\n");

  const applied = await initializeWorkspace({ workspaceRoot: workspace.root });
  assert.equal(applied.gitignore_changed, true);
  assert.ok(applied.wrapper_changes.some((change) => change.endsWith(".gitignore")));
  const ignored = await readFile(join(workspace.root, ".gitignore"), "utf8");
  assert.equal(ignored, "custom.log\n\n# kao-delivery-workspace:ignored-clones:start\nrepositories/frontend/\n# kao-delivery-workspace:ignored-clones:end\n");
  const rerun = await initializeWorkspace({ workspaceRoot: workspace.root });
  assert.equal(rerun.gitignore_changed, false);
});

test("initialization rejects a blanket repositories ignore", async (t) => {
  const workspace = await prepareWrapper();
  t.after(workspace.cleanup);
  await writeFile(join(workspace.root, ".gitignore"), "repositories/\n", "utf8");
  await assert.rejects(initializeWorkspace({ workspaceRoot: workspace.root }), /Blanket repositories\/ ignore/);
});

test("initialization reports an inaccessible configured repository", async (t) => {
  const workspace = await prepareWrapper();
  t.after(workspace.cleanup);
  const raw = await readFile(join(workspace.root, "workspace.yaml"), "utf8");
  await writeFile(join(workspace.root, "workspace.yaml"), raw.replace("repositories/frontend", "repositories/missing"), "utf8");
  await assert.rejects(initializeWorkspace({ workspaceRoot: workspace.root }), /Repository frontend path is not accessible/);
});

test("submodule mode requires both .gitmodules registration and a tracked gitlink", async (t) => {
  const workspace = await prepareWrapper();
  t.after(workspace.cleanup);
  const config = (await import("yaml")).parse(await readFile(join(workspace.root, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  config.repositories.frontend!.mode = "submodule";
  await writeFile(join(workspace.root, "workspace.yaml"), stringifyYaml(config), "utf8");
  await writeFile(join(workspace.root, ".gitmodules"), '[submodule "frontend"]\n\tpath = repositories/frontend\n\turl = ../frontend.git\n', "utf8");

  await assert.rejects(initializeWorkspace({ workspaceRoot: workspace.root }), /does not track a gitlink/);
  await git(workspace.root, ["add", "-f", "repositories/frontend"]);
  const summary = await initializeWorkspace({ workspaceRoot: workspace.root });
  assert.equal(summary.repositories[0]?.mode, "submodule");
  assert.doesNotMatch(await readFile(join(workspace.root, ".gitignore"), "utf8"), /repositories\/frontend\//);
});

test("solo mode permits a direct-commit wrapper policy", async (t) => {
  const workspace = await prepareWrapper();
  t.after(workspace.cleanup);
  const config = (await import("yaml")).parse(await readFile(join(workspace.root, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  config.workspace.mode = "solo";
  config.workflow.wrapper_change_policy = "direct-commit";
  await writeFile(join(workspace.root, "workspace.yaml"), stringifyYaml(config), "utf8");

  const summary = await initializeWorkspace({ workspaceRoot: workspace.root, apply: false });
  assert.equal(summary.mode, "solo");
  assert.equal(summary.wrapper_change_policy, "direct-commit");
});

test("ignore reconciliation refuses malformed managed markers", () => {
  const config = {
    repositories: { frontend: { path: "repositories/frontend", mode: "ignored-clone" } },
  } as unknown as WorkspaceConfig;
  assert.throws(() => reconcileIgnoredClones("# kao-delivery-workspace:ignored-clones:start\n", config), /Malformed managed/);
});
