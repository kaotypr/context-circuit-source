import assert from "node:assert/strict";
import { access, mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { spawnSync } from "node:child_process";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import test from "node:test";
import { configureWorkspace } from "../scripts/lib/configure-workspace.js";
import { git } from "../scripts/lib/git.js";
import type { WorkspaceBootstrapRequest } from "../scripts/lib/types.js";
import { requiredWorkspaceDocuments } from "../scripts/lib/validation.js";
import { projectRoot } from "./helpers.js";

async function neutralWrapper(): Promise<{ root: string; cleanup: () => Promise<void> }> {
  const root = await mkdtemp(join(tmpdir(), "context-circuit-configure-"));
  for (const path of requiredWorkspaceDocuments) {
    const target = join(root, path);
    await mkdir(dirname(target), { recursive: true });
    await writeFile(target, path === "README.md" ? "# Context Circuit\n\nFramework introduction.\n" : `# Fixture: ${path}\n`, "utf8");
  }
  await writeFile(join(root, ".gitignore"), ".runtime/\n\n# context-circuit:ignored-clones:start\n# context-circuit:ignored-clones:end\n", "utf8");
  await writeFile(join(root, "workspace.yaml"), "version: 1\ntemplate_version: 0.2.1\nworkspace:\n  name: uninitialized-workspace\n  mode: team\n  default_branch: main\nrepositories: {}\nactivity:\n  provider: none\n  access: auto\n  required_capabilities: []\n  optional_capabilities: []\nworkflow:\n  human_gates: [plan-approval, task-selection, merge]\n  maximum_repair_attempts: 2\n  wrapper_change_policy: pull-request\n", "utf8");
  return { root, cleanup: async () => rm(root, { recursive: true, force: true }) };
}

function request(existing = false): WorkspaceBootstrapRequest {
  const sources = [{ kind: "prd" as const, reference: "docs/product-brief.md", purpose: "Authoritative product intent", repository: "app" }];
  return {
    contract_version: 1,
    ...(existing ? { authorize_reviewable_changes: true } : {}),
    configuration: {
      version: 1,
      template_version: "0.2.1",
      workspace: { name: "sample-workspace", mode: "team", default_branch: "main", purpose: "Coordinates delivery of the sample product.", remote: "git@example.invalid:sample/wrapper.git" },
      repositories: { app: { path: "repositories/app", mode: "ignored-clone", role: "product application", agent: "app", default_branch: "main", remote: "git@example.invalid:sample/app.git" } },
      activity: { provider: "none", access: "auto", required_capabilities: [], optional_capabilities: [] },
      workflow: { human_gates: ["plan-approval", "task-selection", "merge"], maximum_repair_attempts: 2, wrapper_change_policy: "pull-request", review_mode: "remote" },
      context: { authoritative_sources: sources },
    },
    context: { project_summary: "Coordinates delivery of the sample product.", architecture: [], conventions: [], decisions: [], sources },
    wrapper: existing
      ? { initialize_git: false, authorize_initial_commit: false, commit_message: "unused" }
      : { initialize_git: true, authorize_initial_commit: true, commit_message: "chore: configure workspace", author_name: "Test", author_email: "test@example.invalid" },
    repositories: [{ name: "app", source: existing ? "existing" : "new", authorize_initial_commit: !existing, ...(!existing ? { commit_message: "chore: initialize app", author_name: "Test", author_email: "test@example.invalid" } : {}) }],
  };
}

test("fresh configuration requires exact commit authorization and generates a project-first wrapper", async (t) => {
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  const unauthorized = request();
  unauthorized.wrapper.authorize_initial_commit = false;
  await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: unauthorized }), /initial-commit authorization/);
  await assert.rejects(access(join(workspace.root, ".git")));

  const summary = await configureWorkspace({ workspaceRoot: workspace.root, request: request() });
  assert.equal(summary.route, "bootstrap");
  const readme = await readFile(join(workspace.root, "README.md"), "utf8");
  assert.ok(readme.indexOf("# sample-workspace") < readme.indexOf("## Context Circuit framework"));
  assert.match(readme, /product application/);
  assert.match(await readFile(join(workspace.root, "context/SOURCES.md"), "utf8"), /docs\/product-brief\.md/);
  assert.equal(await git(workspace.root, ["status", "--porcelain=v1"]), "");
});

test("team reconfiguration preserves authored README and context, active runs, HEAD, and rerun output", async (t) => {
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  await configureWorkspace({ workspaceRoot: workspace.root, request: request() });
  await writeFile(join(workspace.root, "README.md"), `${await readFile(join(workspace.root, "README.md"), "utf8")}\n## Team notes\n\nKeep this authored section.\n`, "utf8");
  await writeFile(join(workspace.root, "context/ARCHITECTURE.md"), "# Architecture\n\nCurated architecture stays authored.\n", "utf8");
  await git(workspace.root, ["add", "README.md", "context/ARCHITECTURE.md"]);
  await git(workspace.root, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "docs: add authored context"]);
  const head = await git(workspace.root, ["rev-parse", "HEAD"]);
  await mkdir(join(workspace.root, ".runtime", "runs", "active"), { recursive: true });
  await writeFile(join(workspace.root, ".runtime", "runs", "active", "manifest.json"), "captured-configuration\n", "utf8");

  const changed = request(true);
  changed.configuration.workspace.purpose = "Updated purpose for future work.";
  const result = await configureWorkspace({ workspaceRoot: workspace.root, request: changed });
  assert.equal(result.route, "reconfigure");
  assert.equal(await git(workspace.root, ["rev-parse", "HEAD"]), head);
  assert.match(await readFile(join(workspace.root, "README.md"), "utf8"), /Keep this authored section/);
  assert.equal(await readFile(join(workspace.root, "context/ARCHITECTURE.md"), "utf8"), "# Architecture\n\nCurated architecture stays authored.\n");
  assert.equal(await readFile(join(workspace.root, ".runtime", "runs", "active", "manifest.json"), "utf8"), "captured-configuration\n");
  const firstStatus = await git(workspace.root, ["status", "--porcelain=v1", "--untracked-files=all"]);
  await configureWorkspace({ workspaceRoot: workspace.root, request: changed });
  assert.equal(await git(workspace.root, ["status", "--porcelain=v1", "--untracked-files=all"]), firstStatus);
});

test("existing configuration requires reviewed authorization and rejects credential-bearing source references without mutation", async (t) => {
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  await configureWorkspace({ workspaceRoot: workspace.root, request: request() });
  const before = await git(workspace.root, ["status", "--porcelain=v1"]);
  const unreviewed = request(true);
  delete unreviewed.authorize_reviewable_changes;
  await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: unreviewed }), /authorize_reviewable_changes/);
  const unsafe = request(true);
  unsafe.context.sources![0]!.reference = "https://user:secret@example.invalid/prd";
  unsafe.configuration.context!.authoritative_sources[0]!.reference = unsafe.context.sources![0]!.reference;
  await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: unsafe }), /credentials/);
  assert.equal(await git(workspace.root, ["status", "--porcelain=v1"]), before);
});

test("legacy initialize invocation reports state-aware compatibility routing", async (t) => {
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  const fresh = spawnSync(process.execPath, ["--import", "tsx", join(projectRoot, "scripts/initialize-workspace.ts"), "--check-only"], { cwd: projectRoot, encoding: "utf8", env: { ...process.env, CONTEXT_CIRCUIT_WORKSPACE_ROOT: workspace.root } });
  assert.equal(fresh.status, 0, fresh.stderr);
  assert.match(fresh.stdout, /inspect-fresh/);
  await configureWorkspace({ workspaceRoot: workspace.root, request: request() });
  const existing = spawnSync(process.execPath, ["--import", "tsx", join(projectRoot, "scripts/initialize-workspace.ts"), "--check-only"], { cwd: projectRoot, encoding: "utf8", env: { ...process.env, CONTEXT_CIRCUIT_WORKSPACE_ROOT: workspace.root } });
  assert.equal(existing.status, 0, existing.stderr);
  assert.match(existing.stdout, /inspect-existing/);
});
