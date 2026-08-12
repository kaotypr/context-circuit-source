import assert from "node:assert/strict";
import { access, mkdir, mkdtemp, readFile, rm, symlink, writeFile } from "node:fs/promises";
import { spawnSync } from "node:child_process";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import test from "node:test";
import { configureWorkspace } from "../scripts/lib/configure-workspace.js";
import { git } from "../scripts/lib/git.js";
import type { WorkspaceBootstrapRequest } from "../scripts/lib/types.js";
import { requiredWorkspaceDocuments } from "../scripts/lib/validation.js";
import { projectRoot } from "./helpers.js";
import { reconcileWorkspaceReadme } from "../scripts/lib/workspace-readme.js";

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
  assert.notEqual(firstStatus, "");
  await git(workspace.root, ["add", "-A"]);
  await git(workspace.root, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "chore: review configuration"]);
  await configureWorkspace({ workspaceRoot: workspace.root, request: changed });
  assert.equal(await git(workspace.root, ["status", "--porcelain=v1", "--untracked-files=all"]), "");
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

test("dirty existing wrappers reject tracked and untracked work without changing bytes", async (t) => {
  for (const kind of ["tracked", "untracked"] as const) {
    const workspace = await neutralWrapper();
    t.after(workspace.cleanup);
    await configureWorkspace({ workspaceRoot: workspace.root, request: request() });
    const readme = await readFile(join(workspace.root, "README.md"), "utf8");
    if (kind === "tracked") await writeFile(join(workspace.root, "README.md"), `${readme}authored dirty bytes\n`, "utf8");
    else await writeFile(join(workspace.root, "authored.txt"), "untracked authored bytes\n", "utf8");
    const beforeReadme = await readFile(join(workspace.root, "README.md"));
    const beforeConfig = await readFile(join(workspace.root, "workspace.yaml"));
    await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: request(true) }), /must be clean/);
    assert.deepEqual(await readFile(join(workspace.root, "README.md")), beforeReadme);
    assert.deepEqual(await readFile(join(workspace.root, "workspace.yaml")), beforeConfig);
  }
});

test("repository preflight failure leaves every managed file byte-identical", async (t) => {
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  await configureWorkspace({ workspaceRoot: workspace.root, request: request() });
  const paths = ["README.md", "workspace.yaml", "context/SOURCES.md", ".gitignore"];
  const before = new Map(await Promise.all(paths.map(async (path) => [path, await readFile(join(workspace.root, path))] as const)));
  const invalid = request(true);
  invalid.configuration.repositories.app!.path = "repositories/missing";
  await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: invalid }), /Repository app/);
  for (const path of paths) assert.deepEqual(await readFile(join(workspace.root, path)), before.get(path));
});

test("README reconciliation rejects duplicate and malformed managed blocks", () => {
  const config = request().configuration;
  const block = reconcileWorkspaceReadme("", config);
  assert.throws(() => reconcileWorkspaceReadme(`${block}\n${block}`, config), /at most one/);
  assert.throws(() => reconcileWorkspaceReadme("<!-- context-circuit:workspace:start -->\n", config), /Malformed/);
  assert.throws(() => reconcileWorkspaceReadme("<!-- context-circuit:workspace:end -->\n", config), /Malformed/);
});

test("README and source register symlinks reject before configuration writes", async (t) => {
  for (const relativePath of ["README.md", "context/SOURCES.md"]) {
    const workspace = await neutralWrapper();
    const outside = join(await mkdtemp(join(tmpdir(), "context-circuit-outside-")), "target.md");
    t.after(workspace.cleanup);
    t.after(async () => rm(dirname(outside), { recursive: true, force: true }));
    await configureWorkspace({ workspaceRoot: workspace.root, request: request() });
    await writeFile(outside, "outside bytes\n", "utf8");
    await rm(join(workspace.root, relativePath));
    await symlink(outside, join(workspace.root, relativePath));
    await git(workspace.root, ["add", relativePath]);
    await git(workspace.root, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", `test: symlink ${relativePath}`]);
    const configBefore = await readFile(join(workspace.root, "workspace.yaml"));
    await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: request(true) }), /regular non-symlink/);
    assert.equal(await readFile(outside, "utf8"), "outside bytes\n");
    assert.deepEqual(await readFile(join(workspace.root, "workspace.yaml")), configBefore);
  }
});

test("context references accept stable forms and reject traversal, absolute, backslash, and userinfo", async (t) => {
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  const unsafe = ["../secret.md", "docs/%2e%2e/secret.md", "/tmp/secret.md", "docs\\secret.md", "https://oauth-token@example.invalid/prd", "https://oauth-token%40example.invalid/prd"];
  for (const reference of unsafe) {
    const candidate = request();
    candidate.context.sources![0]!.reference = reference;
    candidate.configuration.context!.authoritative_sources[0]!.reference = reference;
    await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: candidate }), /reference|userinfo|absolute|traversal|backslashes/);
    await assert.rejects(access(join(workspace.root, ".git")));
  }
  const safe = request();
  safe.context.sources![0]!.reference = "github:sample/app#42";
  safe.configuration.context!.authoritative_sources[0]!.reference = "github:sample/app#42";
  assert.equal((await configureWorkspace({ workspaceRoot: workspace.root, request: safe })).route, "bootstrap");
});

test("wrapper and repository remotes reject plain and percent-encoded URL userinfo", async (t) => {
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  for (const remote of ["https://oauth-token@example.invalid/repo.git", "https://oauth-token%40example.invalid/repo.git", "https://user%3Asecret@example.invalid/repo.git"]) {
    const candidate = request();
    candidate.configuration.workspace.remote = remote;
    await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: candidate }), /credentials|userinfo|invalid|valid URL/);
    await assert.rejects(access(join(workspace.root, ".git")));
    candidate.configuration.workspace.remote = "git@example.invalid:sample/wrapper.git";
    candidate.configuration.repositories.app!.remote = remote;
    await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: candidate }), /credentials|userinfo|invalid|valid URL/);
    await assert.rejects(access(join(workspace.root, ".git")));
  }
});

test("exact fresh request rerun is a no-op while material changes require review authorization", async (t) => {
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  const original = request();
  await configureWorkspace({ workspaceRoot: workspace.root, request: original });
  const head = await git(workspace.root, ["rev-parse", "HEAD"]);
  const rerun = await configureWorkspace({ workspaceRoot: workspace.root, request: original });
  assert.equal(rerun.route, "inspect-existing");
  assert.equal(await git(workspace.root, ["rev-parse", "HEAD"]), head);
  assert.equal(await git(workspace.root, ["status", "--porcelain=v1"]), "");
  const changed = request();
  changed.configuration.workspace.purpose = "Materially changed purpose.";
  await assert.rejects(configureWorkspace({ workspaceRoot: workspace.root, request: changed }), /authorize_reviewable_changes/);
});

test("unborn extracted-template Git bootstrap accepts its baseline and rejects extra authored files", async (t) => {
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  await git(workspace.root, ["init", "--initial-branch=main"]);
  const unborn = request();
  unborn.wrapper.initialize_git = false;
  const result = await configureWorkspace({ workspaceRoot: workspace.root, request: unborn });
  assert.equal(result.route, "bootstrap");
  assert.equal(await git(workspace.root, ["rev-list", "--count", "HEAD"]), "1");

  const dirty = await neutralWrapper();
  t.after(dirty.cleanup);
  await git(dirty.root, ["init", "--initial-branch=main"]);
  await writeFile(join(dirty.root, "AUTHORED.md"), "do not overwrite\n", "utf8");
  const dirtyRequest = request();
  dirtyRequest.wrapper.initialize_git = false;
  await assert.rejects(configureWorkspace({ workspaceRoot: dirty.root, request: dirtyRequest }), /authored or unexpected/);
  assert.equal(await readFile(join(dirty.root, "AUTHORED.md"), "utf8"), "do not overwrite\n");
  await assert.rejects(git(dirty.root, ["rev-parse", "HEAD"]));
});

test("built bundle enforces request-path safety and exact rerun idempotence", async (t) => {
  const build = spawnSync(process.execPath, ["--import", "tsx", "scripts/build-template.ts"], { cwd: projectRoot, encoding: "utf8" });
  assert.equal(build.status, 0, build.stderr || build.stdout);
  const workspace = await neutralWrapper();
  t.after(workspace.cleanup);
  const binary = join(projectRoot, ".agents", "bin", "cc.mjs");
  await mkdir(join(workspace.root, ".runtime", "bootstrap"), { recursive: true });
  const requestPath = join(workspace.root, ".runtime", "bootstrap", "request.json");
  await writeFile(requestPath, `${JSON.stringify(request())}\n`, "utf8");
  const first = spawnSync(process.execPath, [binary, "configure-workspace", "--request", ".runtime/bootstrap/request.json"], { cwd: workspace.root, encoding: "utf8" });
  assert.equal(first.status, 0, first.stderr || first.stdout);
  const head = await git(workspace.root, ["rev-parse", "HEAD"]);
  const second = spawnSync(process.execPath, [binary, "configure-workspace", "--request", ".runtime/bootstrap/request.json"], { cwd: workspace.root, encoding: "utf8" });
  assert.equal(second.status, 0, second.stderr || second.stdout);
  assert.match(second.stdout, /inspect-existing/);
  assert.equal(await git(workspace.root, ["rev-parse", "HEAD"]), head);
  assert.equal(await git(workspace.root, ["status", "--porcelain=v1"]), "");

  const outsideDirectory = await mkdtemp(join(tmpdir(), "context-circuit-request-outside-"));
  t.after(async () => rm(outsideDirectory, { recursive: true, force: true }));
  const outsideRequest = join(outsideDirectory, "request.json");
  await writeFile(outsideRequest, `${JSON.stringify(request(true))}\n`, "utf8");
  const outside = spawnSync(process.execPath, [binary, "configure-workspace", "--request", outsideRequest], { cwd: workspace.root, encoding: "utf8" });
  assert.notEqual(outside.status, 0);
  assert.match(outside.stderr, /outside/);

  await rm(requestPath);
  await symlink(outsideRequest, requestPath);
  const linked = spawnSync(process.execPath, [binary, "configure-workspace", "--request", ".runtime/bootstrap/request.json"], { cwd: workspace.root, encoding: "utf8" });
  assert.notEqual(linked.status, 0);
  assert.match(linked.stderr, /regular file/);
});
