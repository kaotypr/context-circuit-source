import assert from "node:assert/strict";
import { readFile, writeFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { createPlanDraft, parsePlanIndex, setPlanState, validatePlanDirectory } from "../scripts/lib/plans.js";
import type { PlanDraftRequest } from "../scripts/lib/types.js";
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
