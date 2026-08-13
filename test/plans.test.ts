import assert from "node:assert/strict";
import { access, readFile, writeFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { archivePlan, createPlanDraft, generatePlanBatch, migrateCurrentPlans, parsePlanIndex, setPlanState, validatePlanDirectory } from "../scripts/lib/plans.js";
import type { PlanDraftRequest, PlanGenerationRequest } from "../scripts/lib/types.js";
import { createTestWorkspace } from "./helpers.js";

const request: PlanDraftRequest = {
  contract_version: 1,
  plan_id: "billing-v2",
  title: "Billing v2",
  source: { kind: "prd", reference: "docs/billing-v2.md" },
  work_prefix: "BILLING",
  summary: "Make retry behavior explicit and visible.",
  affected_repositories: ["frontend"],
  assumptions: [],
  open_questions: ["Which retry states are customer-visible?"],
  requirements: ["Retry state has explicit acceptance criteria."],
  solution: ["Define the contract before implementing user-interface state."],
  delivery: ["Approve the contract, then implement repository work."],
  verification: ["Verify every acceptance criterion independently."],
  risks: ["Unresolved states could create inconsistent behavior."],
  work_items: [
    { key: "retry-contract", title: "Establish retry contract", area: "architecture", repository: "frontend", scope: ["src/contract.ts"], test_scope: [], test_policy: "verifier-only", verification_commands: ["npm test"], acceptance_criteria: ["The retry contract is explicit."] },
    { key: "retry-ui", title: "Display retry state", area: "frontend", repository: "frontend", scope: ["src/App.tsx"], test_scope: ["src/App.test.tsx"], test_policy: "required", verification_commands: ["npm test"], acceptance_criteria: ["The retry state is displayed."], parent: "retry-contract", depends_on: ["retry-contract"] },
  ],
};

test("create-plan writes a validated numbered draft with stable work IDs", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, request, new Date("2026-08-11T08:30:00Z"));
  assert.equal(created.status, "draft");
  assert.deepEqual(created.work_ids, ["BILLING-001", "BILLING-010"]);
  assert.deepEqual(created.documents, [
    "0001-overview.md", "0010-requirements.md", "0020-solution.md", "0040-delivery.md",
    "0050-verification.md", "0070-risks.md", "0080-work-breakdown.md",
  ]);
  const validation = await validatePlanDirectory(created.directory);
  assert.deepEqual(validation.errors, []);
  assert.equal(validation.index?.approved_at, null);
  assert.equal(validation.work_breakdown?.contract_version, 2);
  assert.equal(validation.work_breakdown?.items[0]?.repository, "frontend");
  assert.equal(validation.work_breakdown?.items[0]?.area, "architecture");
  assert.equal(validation.work_breakdown?.items[1]?.parent, "BILLING-001");
  assert.deepEqual(validation.work_breakdown?.items[1]?.depends_on, ["BILLING-001"]);
  assert.doesNotMatch(await readFile(join(created.directory, "0080-work-breakdown.md"), "utf8"), /\| Status \|/);
});

test("create-plan refuses to overwrite an existing plan", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  await createPlanDraft(workspace.root, request, new Date("2026-08-11T08:30:00Z"));
  await assert.rejects(createPlanDraft(workspace.root, request), /refusing to overwrite/);
});

test("approved plan metadata requires an approver and timestamp", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, request, new Date("2026-08-11T08:30:00Z"));
  const indexPath = join(created.directory, "README.md");
  await setPlanState(created.directory, { kind: "approve", approved_by: "kao" }, new Date("2026-08-11T09:00:00Z"));
  assert.deepEqual((await validatePlanDirectory(created.directory)).errors, []);
  assert.equal(parsePlanIndex(await readFile(indexPath, "utf8")).status, "approved");
  assert.doesNotMatch(await readFile(indexPath, "utf8"), /This plan is a draft/);
  await writeFile(join(created.directory, "0010-requirements.md"), "# Requirements\n\nMaterially changed.\n", "utf8");
  assert.match((await validatePlanDirectory(created.directory)).errors.join("\n"), /material_digest|approved_digest/);
  const revised = await setPlanState(created.directory, { kind: "material-revision", reason: "Acceptance criteria changed" }, new Date("2026-08-11T10:00:00Z"));
  assert.equal(revised.status, "draft");
  assert.equal(revised.plan_version, 2);
  assert.equal(revised.approved_at, null);
  assert.deepEqual((await validatePlanDirectory(created.directory)).errors, []);
});

test("plan validation rejects a live-status table shape", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, request, new Date("2026-08-11T08:30:00Z"));
  const breakdownPath = join(created.directory, "0080-work-breakdown.md");
  const raw = await readFile(breakdownPath, "utf8");
  await writeFile(breakdownPath, raw.replace("| Work ID | Title |", "| Work ID | Status | Title |"), "utf8");
  assert.match((await validatePlanDirectory(created.directory)).errors.join("\n"), /canonical seven-column table/);
});

test("plan approval rejects an unknown explicit repository after draft creation", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, request, new Date("2026-08-11T08:30:00Z"));
  const breakdownPath = join(created.directory, "0080-work-breakdown.md");
  const raw = await readFile(breakdownPath, "utf8");
  await writeFile(breakdownPath, raw.replaceAll("| frontend |", "| missing |").replaceAll('"repository": "frontend"', '"repository": "missing"'), "utf8");
  assert.match((await validatePlanDirectory(created.directory)).errors.join("\n"), /repository is not registered: missing/);
  await assert.rejects(setPlanState(created.directory, { kind: "approve", approved_by: "kao" }), /repository is not registered: missing/);
});

test("legacy area-only plans normalize only an exact registered repository key without changing approved material", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const legacyRequest: PlanDraftRequest = { ...request, plan_id: "legacy-plan", work_items: [{ ...request.work_items[0]!, area: "frontend" }] };
  const created = await createPlanDraft(workspace.root, legacyRequest, new Date("2026-08-11T08:30:00Z"));
  const breakdownPath = join(created.directory, "0080-work-breakdown.md");
  const current = await readFile(breakdownPath, "utf8");
  const legacy = current
    .replace("| Work ID | Title | Parent | Depends on | Repository | Area | External reference |", "| Work ID | Title | Parent | Depends on | Area | External reference |")
    .replace("| --- | --- | --- | --- | --- | --- | --- |", "| --- | --- | --- | --- | --- | --- |")
    .replace("| BILLING-001 | Establish retry contract | — | — | frontend | frontend | — |", "| BILLING-001 | Establish retry contract | — | — | frontend | — |")
    .replace('"contract_version": 2', '"contract_version": 1')
    .replace(/\n      "repository": "frontend",/, "");
  await writeFile(breakdownPath, legacy, "utf8");
  const approved = await setPlanState(created.directory, { kind: "approve", approved_by: "kao" }, new Date("2026-08-11T09:00:00Z"));
  const approvedMaterial = await readFile(breakdownPath, "utf8");
  const validation = await validatePlanDirectory(created.directory);
  assert.deepEqual(validation.errors, []);
  assert.equal(validation.work_breakdown?.items[0]?.repository, "frontend");
  assert.equal(validation.index?.approved_digest, approved.approved_digest);
  assert.equal(await readFile(breakdownPath, "utf8"), approvedMaterial);
});

test("legacy area-only plans do not guess repositories from descriptive text", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, { ...request, plan_id: "legacy-description", work_items: [request.work_items[0]!] }, new Date("2026-08-11T08:30:00Z"));
  const breakdownPath = join(created.directory, "0080-work-breakdown.md");
  const current = await readFile(breakdownPath, "utf8");
  const legacy = current
    .replace("| Work ID | Title | Parent | Depends on | Repository | Area | External reference |", "| Work ID | Title | Parent | Depends on | Area | External reference |")
    .replace("| --- | --- | --- | --- | --- | --- | --- |", "| --- | --- | --- | --- | --- | --- |")
    .replace("| BILLING-001 | Establish retry contract | — | — | frontend | architecture | — |", "| BILLING-001 | Establish retry contract | — | — | application foundation | — |")
    .replace('"contract_version": 2', '"contract_version": 1')
    .replace(/\n      "repository": "frontend",/, "");
  await writeFile(breakdownPath, legacy, "utf8");
  await assert.rejects(setPlanState(created.directory, { kind: "approve", approved_by: "kao" }), /not an exact registered repository key/);
});

test("plan draft rejects unknown repositories and dependency cycles before writing", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const invalid: PlanDraftRequest = {
    ...request,
    plan_id: "invalid-plan",
    affected_repositories: ["backend"],
    work_items: [
      { key: "one", title: "One", area: "frontend", repository: "frontend", scope: ["src/one.ts"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["One passes."], depends_on: ["two"] },
      { key: "two", title: "Two", area: "frontend", repository: "frontend", scope: ["src/two.ts"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["Two passes."], depends_on: ["one"] },
    ],
  };
  await assert.rejects(createPlanDraft(workspace.root, invalid), /not registered|dependency cycle/);
  await assert.rejects(createPlanDraft(workspace.root, {} as PlanDraftRequest), /Invalid plan draft request/);
});

const productKnowledgeRequest: PlanDraftRequest = {
  ...request,
  plan_id: "checkout-retry",
  work_prefix: "CKO",
  product_knowledge: {
    impact: "behavior-change",
    references: ["context/domains/checkout/workflows/place-order.md", "context/roles/shopper.md"],
    proposed_change: "Auto-retry a declined card once before failing the order.",
  },
};

test("create-plan records Product Knowledge references and impact", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, productKnowledgeRequest, new Date("2026-08-11T08:30:00Z"));
  const validation = await validatePlanDirectory(created.directory);
  assert.deepEqual(validation.errors, []);
  assert.equal(validation.index?.product_knowledge?.impact, "behavior-change");
  assert.deepEqual(validation.index?.product_knowledge?.references, [
    "context/domains/checkout/workflows/place-order.md",
    "context/roles/shopper.md",
  ]);
  const overview = await readFile(join(created.directory, "0001-overview.md"), "utf8");
  assert.match(overview, /## Product impact/);
  assert.match(overview, /Auto-retry a declined card once/);
});

test("Product Knowledge impact survives approval", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const created = await createPlanDraft(workspace.root, productKnowledgeRequest, new Date("2026-08-11T08:30:00Z"));
  await setPlanState(created.directory, { kind: "approve", approved_by: "maintainer" }, new Date("2026-08-11T09:00:00Z"));
  const index = parsePlanIndex(await readFile(join(created.directory, "README.md"), "utf8"));
  assert.equal(index.status, "approved");
  assert.equal(index.product_knowledge?.impact, "behavior-change");
});

test("a behavior-changing plan requires a proposed change summary", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const invalid: PlanDraftRequest = {
    ...productKnowledgeRequest,
    plan_id: "checkout-retry-invalid",
    product_knowledge: { impact: "new-workflow", references: [] },
  };
  await assert.rejects(createPlanDraft(workspace.root, invalid), /requires a proposed_change/);
});

test("an impact of none must not include a proposed change", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const invalid: PlanDraftRequest = {
    ...productKnowledgeRequest,
    plan_id: "checkout-none-invalid",
    product_knowledge: { impact: "none", references: [], proposed_change: "Should not be here." },
  };
  await assert.rejects(createPlanDraft(workspace.root, invalid), /must not include a proposed_change/);
});

const generationRequest: PlanGenerationRequest = {
  contract_version: 2,
  source: { kind: "prd", reference: "docs/runtime.md" },
  plans: [
    {
      plan_id: "registry-foundation", title: "Registry foundation", repository: "frontend", work_prefix: "REG",
      summary: "Create the registry foundation.", assumptions: [], open_questions: [], requirements: ["The registry is exact."],
      solution: ["Use repository-scoped collections."], delivery: ["Generate the foundation first."], verification: ["Validate the generated inventory."], risks: [],
      work_items: [{ key: "registry", title: "Create registry", area: "planning", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["The registry is created."] }],
    },
    {
      plan_id: "registry-ui", title: "Registry UI", repository: "frontend", work_prefix: "REGUI", plan_number: 2, depends_on_plans: ["registry-foundation"],
      summary: "Expose the registry.", assumptions: [], open_questions: [], requirements: ["The registry is visible."], solution: ["Render the registry."], delivery: ["Build on the foundation."], verification: ["Verify the registry."], risks: [],
      work_items: [{ key: "ui", title: "Render registry", area: "planning", repository: "frontend", scope: ["src/App.tsx"], test_scope: [], test_policy: "verifier-only", verification_commands: [], acceptance_criteria: ["The registry is visible."] }],
    },
  ],
};

test("root generation creates ordered peer plans in one exact repository collection", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const generated = await generatePlanBatch(workspace.root, generationRequest, new Date("2026-08-14T08:30:00Z"));
  assert.equal(generated.atomic, true);
  assert.deepEqual(generated.plans.map((plan) => plan.plan_number), [1, 2]);
  assert.equal(generated.plans.every((plan) => plan.repository_collection === "frontend-plans"), true);
  assert.equal(generated.plans[0]?.directory.endsWith("plans/frontend-plans/001-registry-foundation"), true);
  assert.equal(generated.plans[1]?.directory.endsWith("plans/frontend-plans/002-registry-ui"), true);
  for (const plan of generated.plans) {
    const validation = await validatePlanDirectory(plan.directory);
    assert.deepEqual(validation.errors, []);
    assert.equal(validation.index?.contract_version, 2);
    assert.equal(validation.work_breakdown?.items.length, 1);
  }
  assert.match(await readFile(join(workspace.root, "plans", "README.md"), "utf8"), /frontend-plans/);
  assert.match(await readFile(join(generated.plans[0]!.directory, "tasks", "REG-001.md"), "utf8"), /task_id: REG-001/);
});

test("root generation refuses a conceptual collection and leaves no partial plan output", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const invalid = structuredClone(generationRequest);
  invalid.plans[0]!.repository_collection = "foundation-plans";
  await assert.rejects(generatePlanBatch(workspace.root, invalid), /exact registered collection/);
  await assert.rejects(access(join(workspace.root, "plans")));
});

test("BAU plans use an independent collection sequence and explicit archive moves", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const request = structuredClone(generationRequest);
  request.plans = [structuredClone(request.plans[0]!)];
  request.plans[0]!.plan_id = "maintenance-refresh";
  request.plans[0]!.track = "bau";
  const generated = await generatePlanBatch(workspace.root, request);
  const plan = generated.plans[0]!;
  await setPlanState(plan.directory, { kind: "approve", approved_by: "owner" });
  const archived = await archivePlan(workspace.root, plan.plan_reference, "owner", "Explicitly archived for test coverage.");
  assert.equal(archived.destination.endsWith("archived/plans/frontend-plans/__BAU__/maintenance-refresh"), true);
  assert.equal((await validatePlanDirectory(archived.destination)).index?.status, "archived");
  await assert.rejects(access(plan.directory));
});

test("current context plans migrate once to the root collection while preserving work IDs", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const legacy = await createPlanDraft(workspace.root, request, new Date("2026-08-14T10:00:00Z"));
  await setPlanState(legacy.directory, { kind: "approve", approved_by: "owner" }, new Date("2026-08-14T10:00:01Z"));
  const migrated = await migrateCurrentPlans(workspace.root, new Date("2026-08-14T10:01:00Z"));
  assert.equal(migrated.migrated.length, 1);
  assert.equal(migrated.migrated[0]!.repository_collection, "frontend-plans");
  assert.deepEqual(migrated.migrated[0]!.work_ids, ["BILLING-001", "BILLING-010"]);
  assert.equal((await validatePlanDirectory(migrated.migrated[0]!.directory)).index?.status, "approved");
  await assert.rejects(access(join(workspace.root, "context", "plans", "billing-v2")));
});
