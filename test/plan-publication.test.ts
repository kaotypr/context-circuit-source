import assert from "node:assert/strict";
import { readFile, writeFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { createPlanDraft, setPlanState, validatePlanDirectory } from "../scripts/lib/plans.js";
import { preparePlanPublication, recordPlanPublication } from "../scripts/lib/plan-publication.js";
import type { PlanDraftRequest, PlanPublicationDiscovery } from "../scripts/lib/types.js";
import { validateContract } from "../scripts/lib/validation.js";
import { createTestWorkspace } from "./helpers.js";

const request: PlanDraftRequest = {
  contract_version: 1, plan_id: "publish-proof", title: "Publication proof",
  source: { kind: "idea", reference: "publication test" }, work_prefix: "PUB",
  summary: "Publish stable work without duplicates.", affected_repositories: ["frontend"],
  assumptions: [], open_questions: [], requirements: ["Each work ID maps to at most one confirmed external task."],
  solution: ["Discover mappings before explicit creation."], delivery: ["Publish parents and dependencies first."],
  verification: ["Validate confirmed mappings."], risks: ["Partial external success must remain recoverable."],
  work_items: [
    { key: "contract", title: "Define publication contract", area: "publication architecture", repository: "frontend", scope: ["src/contract.ts"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Contract is defined."] },
    { key: "implementation", title: "Implement publication", area: "delivery workflow", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Publication is implemented."], parent: "contract", depends_on: ["contract"] },
    { key: "docs", title: "Document publication", area: "operator guidance", repository: "frontend", scope: ["README.md"], test_scope: [], test_policy: "not-required", verification_commands: [], acceptance_criteria: ["Publication is documented."], depends_on: ["implementation"] },
  ],
};

async function configure(root: string): Promise<void> {
  const path = join(root, "workspace.yaml");
  await writeFile(path, (await readFile(path, "utf8"))
    .replace("provider: none", "provider: fake-session")
    .replace("required_capabilities: []", "required_capabilities: [create-tasks]"), "utf8");
}

function discovery(mappings: PlanPublicationDiscovery["mappings"] = []): PlanPublicationDiscovery {
  return { contract_version: 1, provider: "fake-session", destination: "Fake project / Delivery", mappings };
}

test("approved plan publication discovers mappings, orders dependencies, and records confirmed references", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await configure(workspace.root);
  const created = await createPlanDraft(workspace.root, request, new Date("2026-08-12T03:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-12T03:05:00Z"));
  const prepared = await preparePlanPublication({
    workspaceRoot: workspace.root, planId: request.plan_id,
    discovery: discovery([{ work_id: "PUB-001", external_reference: "TASK-1", evidence: "Confirmed by read-only discovery." }]),
    now: new Date("2026-08-12T03:10:00Z"),
  });
  assert.deepEqual(prepared.items.map((item) => item.work_id), ["PUB-001", "PUB-010", "PUB-020"]);
  assert.deepEqual(prepared.items.map((item) => item.status), ["existing", "proposed", "proposed"]);
  assert.equal(prepared.contract_version, 2);
  assert.equal(prepared.items[0]?.action, "skip-existing");
  assert.deepEqual(prepared.items.map((item) => item.repository), ["frontend", "frontend", "frontend"]);
  assert.deepEqual(await validateContract("plan-publication-record", prepared), []);

  const first = await recordPlanPublication({ workspaceRoot: workspace.root, planId: request.plan_id, workId: "PUB-010", status: "created", evidence: "Create response confirmed.", externalReference: "TASK-2", now: new Date("2026-08-12T03:11:00Z") });
  assert.equal(first.status, "in-progress");
  const completed = await recordPlanPublication({ workspaceRoot: workspace.root, planId: request.plan_id, workId: "PUB-020", status: "created", evidence: "Create response confirmed.", externalReference: "TASK-3", now: new Date("2026-08-12T03:12:00Z") });
  assert.equal(completed.status, "completed");
  const plan = await validatePlanDirectory(created.directory);
  assert.deepEqual(plan.errors, []);
  assert.equal(plan.index?.status, "approved");
  assert.deepEqual(plan.work_breakdown?.items.map((item) => item.external_reference), [null, "TASK-2", "TASK-3"]);
  const retried = await recordPlanPublication({ workspaceRoot: workspace.root, planId: request.plan_id, workId: "PUB-020", status: "created", evidence: "Create response confirmed.", externalReference: "TASK-3" });
  assert.deepEqual(retried, completed);
});

test("partial publication stops with confirmed successes and retryable failure evidence", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await configure(workspace.root);
  const created = await createPlanDraft(workspace.root, { ...request, plan_id: "partial-proof" }, new Date("2026-08-12T03:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-12T03:05:00Z"));
  await preparePlanPublication({ workspaceRoot: workspace.root, planId: "partial-proof", discovery: discovery(), now: new Date("2026-08-12T03:06:00Z") });
  await recordPlanPublication({ workspaceRoot: workspace.root, planId: "partial-proof", workId: "PUB-001", status: "created", evidence: "Confirmed.", externalReference: "TASK-1", now: new Date("2026-08-12T03:07:00Z") });
  const partial = await recordPlanPublication({ workspaceRoot: workspace.root, planId: "partial-proof", workId: "PUB-010", status: "failed", evidence: "Provider rejected the request.", now: new Date("2026-08-12T03:08:00Z") });
  assert.equal(partial.status, "partial");
  assert.deepEqual(partial.items.map((item) => item.status), ["created", "failed", "proposed"]);
  assert.equal((await validatePlanDirectory(created.directory)).work_breakdown?.items[0]?.external_reference, "TASK-1");
});

test("publication rejects draft plans and conflicting discovered mappings", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await configure(workspace.root);
  const created = await createPlanDraft(workspace.root, { ...request, plan_id: "draft-proof" }, new Date("2026-08-12T03:00:00Z"));
  await assert.rejects(preparePlanPublication({ workspaceRoot: workspace.root, planId: "draft-proof", discovery: discovery() }), /Only an approved plan/);
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-12T03:05:00Z"));
  await assert.rejects(preparePlanPublication({ workspaceRoot: workspace.root, planId: "draft-proof", discovery: discovery([
    { work_id: "UNKNOWN-001", external_reference: "TASK-X", evidence: "Discovery." },
  ]) }), /unknown work ID/i);
});

test("publication revalidates plan repositories against current workspace configuration", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await configure(workspace.root);
  const created = await createPlanDraft(workspace.root, { ...request, plan_id: "removed-repository" }, new Date("2026-08-12T03:00:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "reviewer" }, new Date("2026-08-12T03:05:00Z"));
  const workspacePath = join(workspace.root, "workspace.yaml");
  await writeFile(workspacePath, (await readFile(workspacePath, "utf8")).replace("  frontend:\n", "  replacement:\n"), "utf8");
  await assert.rejects(preparePlanPublication({ workspaceRoot: workspace.root, planId: "removed-repository", discovery: discovery() }), /repository is not registered: frontend/);
});
