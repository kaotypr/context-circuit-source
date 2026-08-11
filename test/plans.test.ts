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
    { key: "retry-contract", title: "Establish retry contract", area: "architecture" },
    { key: "retry-ui", title: "Display retry state", area: "frontend", parent: "retry-contract", depends_on: ["retry-contract"] },
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
  assert.match((await validatePlanDirectory(created.directory)).errors.join("\n"), /canonical six-column table/);
});

test("plan draft rejects unknown repositories and dependency cycles before writing", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const invalid: PlanDraftRequest = {
    ...request,
    plan_id: "invalid-plan",
    affected_repositories: ["backend"],
    work_items: [
      { key: "one", title: "One", area: "frontend", depends_on: ["two"] },
      { key: "two", title: "Two", area: "frontend", depends_on: ["one"] },
    ],
  };
  await assert.rejects(createPlanDraft(workspace.root, invalid), /not registered|dependency cycle/);
  await assert.rejects(createPlanDraft(workspace.root, {} as PlanDraftRequest), /Invalid plan draft request/);
});
