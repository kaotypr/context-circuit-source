import assert from "node:assert/strict";
import { mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { prepareImportContext } from "../scripts/lib/import-context.js";
import { git } from "../scripts/lib/git.js";
import { validateContract } from "../scripts/lib/validation.js";

async function fixture(options: { instructions?: boolean; repositoryContext?: boolean } = {}): Promise<{ root: string; repository: string; cleanup: () => Promise<void> }> {
  const root = await mkdtemp(join(tmpdir(), "context-import-"));
  const repository = join(root, "repositories", "product");
  await mkdir(join(repository, "docs"), { recursive: true });
  await mkdir(join(root, "contributions"), { recursive: true });
  await writeFile(join(root, ".gitignore"), ".runtime/\nrepositories/\n", "utf8");
  await writeFile(join(root, "workspace.yaml"), `version: 1
template_version: 0.2.1
workspace: { name: import-test, mode: team, default_branch: main }
repositories:
  product: { path: repositories/product, mode: ignored-clone, role: product, agent: repository-worker, default_branch: main }
activity: { provider: none, access: auto, required_capabilities: [], optional_capabilities: [] }
workflow: { human_gates: [plan-approval, task-selection, merge], maximum_repair_attempts: 2, wrapper_change_policy: pull-request }
`, "utf8");
  if (options.instructions !== false) await writeFile(join(repository, "AGENTS.md"), "# Product instructions\n", "utf8");
  await writeFile(join(repository, "docs", "architecture.md"), "# Architecture\n\nA source-cited design.\n", "utf8");
  await writeFile(join(repository, "package.json"), "{\"name\":\"product\"}\n", "utf8");
  if (options.repositoryContext) {
    await mkdir(join(repository, "context"), { recursive: true });
    await writeFile(join(repository, "context", "DECISIONS.md"), "# Decisions\n\nUse deterministic discovery.\n", "utf8");
  }
  await git(repository, ["init", "--initial-branch=main"]);
  await git(repository, ["add", "."]);
  await git(repository, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "fixture"]);
  await git(root, ["init", "--initial-branch=main"]);
  await git(root, ["add", "."]);
  await git(root, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "wrapper"]);
  return { root, repository, cleanup: async () => rm(root, { recursive: true, force: true }) };
}

const request = { contract_version: 1 as const, repository: "product", authorize_contribution_write: true as const };

test("discovers in fixed order and emits a source-cited contribution plus bounded manifest", async (t) => {
  const workspace = await fixture();
  t.after(workspace.cleanup);
  const result = await prepareImportContext({ workspaceRoot: workspace.root, request, now: new Date("2026-08-13T12:00:00Z") });
  assert.equal(result.contribution, "contributions/import-context/product/20260813T120000Z-import-context-product.md");
  assert.deepEqual(result.evidence.slice(0, 2).map((item) => [item.kind, item.path]), [["root-instruction", "AGENTS.md"], ["documentation", "docs/architecture.md"]]);
  assert.ok(result.evidence.length <= 50);
  const contribution = await readFile(join(workspace.root, result.contribution), "utf8");
  for (const evidence of result.evidence) assert.ok(contribution.includes(`\`${evidence.path}\``));
  const manifest = JSON.parse(await readFile(result.manifest, "utf8"));
  assert.equal(manifest.source_commit, await git(workspace.repository, ["rev-parse", "HEAD"]));
  assert.equal(await git(workspace.repository, ["status", "--porcelain"]), "");
  assert.equal((await validateContract("context-sync-request", { contract_version: 1, author: "curator", contributions: [result.contribution], proposals: [{ source_contribution: result.contribution, classification: "durable-wrapper", summary: "Curate evidence", proposed_change: "Cited update", target: "context/PROJECT.md" }] })).length, 0);
});

test("includes optional repository context last as high-trust evidence", async (t) => {
  const workspace = await fixture({ repositoryContext: true });
  t.after(workspace.cleanup);
  const result = await prepareImportContext({ workspaceRoot: workspace.root, request, now: new Date("2026-08-13T12:01:00Z") });
  const context = result.evidence.find((item) => item.path === "context/DECISIONS.md");
  assert.equal(context?.kind, "repository-context");
  assert.equal(context?.trust, "high");
  assert.equal(result.evidence.at(-1)?.path, "context/DECISIONS.md");
});

test("refuses invalid and unregistered import requests with distinct errors", async (t) => {
  const workspace = await fixture();
  t.after(workspace.cleanup);
  await assert.rejects(prepareImportContext({ workspaceRoot: workspace.root, request: { contract_version: 1, repository: "Product" } }), /Invalid import request/);
  await assert.rejects(prepareImportContext({ workspaceRoot: workspace.root, request: { ...request, repository: "missing" } }), /Unregistered repository: missing/);
});

test("refuses missing required discovery inputs", async (t) => {
  const workspace = await fixture({ instructions: false });
  t.after(workspace.cleanup);
  await assert.rejects(prepareImportContext({ workspaceRoot: workspace.root, request }), /Missing or unreadable required discovery inputs/);
});

test("refuses dirty workspace state without altering unrecorded work", async (t) => {
  const workspace = await fixture();
  t.after(workspace.cleanup);
  const unrecorded = join(workspace.root, "keep.txt");
  await writeFile(unrecorded, "preserve\n", "utf8");
  await assert.rejects(prepareImportContext({ workspaceRoot: workspace.root, request }), /Dirty workspace Git state/);
  assert.equal(await readFile(unrecorded, "utf8"), "preserve\n");
});

test("refuses dirty source state without altering unrecorded work", async (t) => {
  const workspace = await fixture();
  t.after(workspace.cleanup);
  const unrecorded = join(workspace.repository, "keep.txt");
  await writeFile(unrecorded, "preserve\n", "utf8");
  await assert.rejects(prepareImportContext({ workspaceRoot: workspace.root, request }), /Dirty source Git state/);
  assert.equal(await readFile(unrecorded, "utf8"), "preserve\n");
});
