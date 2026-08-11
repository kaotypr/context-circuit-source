import assert from "node:assert/strict";
import { access, mkdir, mkdtemp, readFile, rm, symlink, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import test from "node:test";
import { stringify as stringifyYaml } from "yaml";
import { git } from "../scripts/lib/git.js";
import { bootstrapWorkspace, initializeWorkspace, reconcileIgnoredClones } from "../scripts/lib/initialize-workspace.js";
import type { WorkspaceBootstrapRequest, WorkspaceConfig } from "../scripts/lib/types.js";
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

async function prepareNeutralWrapper(): Promise<{ root: string; cleanup: () => Promise<void> }> {
  const root = await mkdtemp(join(tmpdir(), "context-circuit-bootstrap-"));
  for (const path of requiredWorkspaceDocuments) await ensureFile(join(root, path));
  await ensureFile(join(root, ".gitignore"), ".DS_Store\n.runtime/\n\n# context-circuit:ignored-clones:start\n# context-circuit:ignored-clones:end\n");
  await ensureFile(join(root, "workspace.yaml"), "version: 1\ntemplate_version: 0.2.0\nworkspace:\n  name: uninitialized-workspace\n  mode: team\n  default_branch: main\nrepositories: {}\nactivity:\n  provider: none\n  access: auto\n  required_capabilities: []\n  optional_capabilities: []\nworkflow:\n  human_gates: [plan-approval, task-selection, merge]\n  maximum_repair_attempts: 2\n  wrapper_change_policy: pull-request\n");
  return { root, cleanup: async () => rm(root, { recursive: true, force: true }) };
}

function bootstrapRequest(source: "new" | "clone" | "existing" | "submodule", url?: string): WorkspaceBootstrapRequest {
  return {
    contract_version: 1,
    configuration: {
      version: 1,
      template_version: "0.2.0",
      workspace: { name: "example-product", mode: "team", default_branch: "main" },
      repositories: {
        app: { path: "repositories/app", mode: source === "submodule" ? "submodule" : "ignored-clone", role: "application", agent: "app", default_branch: "main" },
      },
      activity: { provider: "none", access: "auto", required_capabilities: [], optional_capabilities: [] },
      workflow: { human_gates: ["plan-approval", "task-selection", "merge"], maximum_repair_attempts: 2, wrapper_change_policy: "pull-request" },
    },
    context: {
      project_summary: "Example Product provides a small application used to verify workspace initialization.",
      architecture: ["The app repository owns the product implementation."],
      conventions: ["Keep product-specific instructions in the app repository."],
      decisions: ["Use an ignored repository clone for local product work."],
    },
    wrapper: { initialize_git: true, authorize_initial_commit: true, commit_message: "chore: initialize Context Circuit", author_name: "Test", author_email: "test@example.invalid" },
    repositories: [{
      name: "app",
      source,
      ...(url ? { url } : {}),
      authorize_initial_commit: source === "new",
      ...(source === "new" ? { commit_message: "chore: initialize repository", author_name: "Test", author_email: "test@example.invalid" } : {}),
    }],
  };
}

async function prepareCloneSource(): Promise<{ root: string; cleanup: () => Promise<void> }> {
  const root = await mkdtemp(join(tmpdir(), "context-circuit-source-"));
  await git(root, ["init", "--initial-branch=main"]);
  await writeFile(join(root, "README.md"), "# Existing product\n", "utf8");
  await git(root, ["add", "."]);
  await git(root, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "initial product"]);
  return { root, cleanup: async () => rm(root, { recursive: true, force: true }) };
}

test("bootstrap creates a neutral wrapper and new repository with authorized base commits", async (t) => {
  const workspace = await prepareNeutralWrapper();
  t.after(workspace.cleanup);

  const summary = await bootstrapWorkspace({ workspaceRoot: workspace.root, request: bootstrapRequest("new") });
  assert.equal(summary.status, "initialized");
  assert.ok(summary.wrapper_initial_commit);
  assert.deepEqual(summary.repositories.map((repository) => repository.name), ["app"]);
  assert.equal(await git(workspace.root, ["status", "--porcelain=v1"]), "");
  assert.equal(await git(join(workspace.root, "repositories", "app"), ["rev-list", "--count", "HEAD"]), "1");
  assert.equal(await git(workspace.root, ["ls-files", "repositories/app"]), "");
  assert.match(await readFile(join(workspace.root, ".gitignore"), "utf8"), /repositories\/app\//);
  assert.match(await readFile(join(workspace.root, "agents", "app.md"), "utf8"), /application role/);
  assert.match(await readFile(join(workspace.root, "context", "PROJECT.md"), "utf8"), /Example Product provides/);
  assert.match(await readFile(join(workspace.root, "context", "ARCHITECTURE.md"), "utf8"), /app repository owns/);
  assert.doesNotMatch(await readFile(join(workspace.root, "context", "PROJECT.md"), "utf8"), /Context Circuit 0\.2\.0/);
  await assert.rejects(access(join(workspace.root, "agents", "frontend.md")));
  assert.equal((await initializeWorkspace({ workspaceRoot: workspace.root, apply: false })).gitignore_changed, false);
});

test("bootstrap clones a configured repository without adding an artificial commit", async (t) => {
  const workspace = await prepareNeutralWrapper();
  const source = await prepareCloneSource();
  t.after(workspace.cleanup);
  t.after(source.cleanup);
  const sourceHead = await git(source.root, ["rev-parse", "HEAD"]);

  const summary = await bootstrapWorkspace({ workspaceRoot: workspace.root, request: bootstrapRequest("clone", source.root) });
  assert.equal(summary.repositories[0]?.remote, source.root);
  assert.equal(await git(join(workspace.root, "repositories", "app"), ["rev-parse", "HEAD"]), sourceHead);
});

test("bootstrap registers an existing local repository without changing it", async (t) => {
  const workspace = await prepareNeutralWrapper();
  const path = join(workspace.root, "repositories", "app");
  await mkdir(path, { recursive: true });
  await git(path, ["init", "--initial-branch=main"]);
  await git(path, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "--allow-empty", "-m", "existing"]);
  const before = await git(path, ["rev-parse", "HEAD"]);
  t.after(workspace.cleanup);

  await bootstrapWorkspace({ workspaceRoot: workspace.root, request: bootstrapRequest("existing") });
  assert.equal(await git(path, ["rev-parse", "HEAD"]), before);
  assert.equal(await git(path, ["rev-list", "--count", "HEAD"]), "1");
});

test("bootstrap registers a submodule and tracks its gitlink in the configured first commit", async (t) => {
  const workspace = await prepareNeutralWrapper();
  const source = await prepareCloneSource();
  t.after(workspace.cleanup);
  t.after(source.cleanup);

  await bootstrapWorkspace({ workspaceRoot: workspace.root, request: bootstrapRequest("submodule", source.root) });
  assert.match(await git(workspace.root, ["ls-files", "--stage", "repositories/app"]), /^160000 /);
  assert.doesNotMatch(await readFile(join(workspace.root, ".gitignore"), "utf8"), /repositories\/app\//);
});

test("bootstrap refuses an existing target before initializing the wrapper", async (t) => {
  const workspace = await prepareNeutralWrapper();
  await mkdir(join(workspace.root, "repositories", "app"), { recursive: true });
  t.after(workspace.cleanup);

  await assert.rejects(bootstrapWorkspace({ workspaceRoot: workspace.root, request: bootstrapRequest("new") }), /refuses to replace existing path/);
  await assert.rejects(access(join(workspace.root, ".git")));
});

test("bootstrap rejects credential-bearing clone URLs before initializing Git", async (t) => {
  const workspace = await prepareNeutralWrapper();
  t.after(workspace.cleanup);
  const request = bootstrapRequest("clone", "https://user:secret@example.invalid/app.git");

  await assert.rejects(bootstrapWorkspace({ workspaceRoot: workspace.root, request }), /appears to contain credentials/);
  await assert.rejects(access(join(workspace.root, ".git")));
});

test("bootstrap rejects repository parents that escape through a symlink", async (t) => {
  const workspace = await prepareNeutralWrapper();
  const outside = await mkdtemp(join(tmpdir(), "context-circuit-outside-"));
  await symlink(outside, join(workspace.root, "repositories"));
  t.after(workspace.cleanup);
  t.after(async () => rm(outside, { recursive: true, force: true }));

  await assert.rejects(bootstrapWorkspace({ workspaceRoot: workspace.root, request: bootstrapRequest("new") }), /parent cannot be a symbolic link/);
  await assert.rejects(access(join(workspace.root, ".git")));
});

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
  assert.equal(ignored, "custom.log\n\n# context-circuit:ignored-clones:start\nrepositories/frontend/\n# context-circuit:ignored-clones:end\n");
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
  assert.throws(() => reconcileIgnoredClones("# context-circuit:ignored-clones:start\n", config), /Malformed managed/);
});
