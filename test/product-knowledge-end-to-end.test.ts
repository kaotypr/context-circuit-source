import assert from "node:assert/strict";
import { access, appendFile, cp, mkdtemp, readFile, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { createPlanDraft, setPlanState, validatePlanDirectory } from "../scripts/lib/plans.js";
import { buildTaskContextPackage, synchronizeProductKnowledge, validateProductKnowledgeTree } from "../scripts/lib/product-knowledge.js";
import { discoverProductKnowledge, writeProductKnowledgeDiscovery } from "../scripts/lib/product-knowledge-discovery.js";
import { generateOnboardingPack } from "../scripts/lib/product-knowledge-onboarding.js";
import { resolveProductKnowledgeCloseout } from "../scripts/lib/finish-work.js";
import type { PlanDraftRequest } from "../scripts/lib/types.js";
import { createTestWorkspace, projectRoot } from "./helpers.js";

const largeFixture = join(projectRoot, "fixtures", "product-knowledge-large");
const workflowTarget = "context/domains/checkout/workflows/place-order.md";

const effectiveWorkflow = `---
kind: workflow
title: Place Order
owners:
  - Checkout team
sources:
  - Commerce platform PRD
review_date: 2026-03-01
implementation_ownership: checkout-service
known_gaps:
  - None outstanding.
---

# Place Order

## Outcome

A confirmed, paid order.

## Actors

Shopper and checkout-service.

## Entry points

The checkout screen.

## Current flow

1. Validate the cart.
2. Authorize payment, retrying once after a declined card.
3. Confirm the order.

## Variations

Guests enter an address inline.

## Business rules

An order is confirmed only after payment authorization succeeds.
`;

test("Product Knowledge composes end to end across the delivery journey", async (t) => {
  const workspace = await createTestWorkspace();
  t.after(workspace.cleanup);
  const root = workspace.root;
  await cp(largeFixture, join(root, "context"), { recursive: true });
  await appendFile(join(root, "workspace.yaml"), "product_knowledge:\n  confirming_role: product-owner\n", "utf8");

  // (a) Bounded, valid tree at scale.
  const tree = await validateProductKnowledgeTree(join(root, "context"));
  assert.equal(tree.present, true);
  assert.deepEqual(tree.errors, []);
  assert.equal(tree.pages, 20);

  // (b) Read-only discovery resolves only relevant pages and never mutates canonical context.
  const before = await readFile(join(root, "context", workflowTarget.replace("context/", "")), "utf8");
  const discovery = await discoverProductKnowledge({
    contextDir: join(root, "context"),
    revision: "rev-1",
    topic: "declined card retry",
    select: ["domains/checkout/workflows/place-order.md"],
    observations: [{ topic: "retry", summary: "Orders should retry once after a declined card.", target: "domains/checkout/workflows/place-order.md", conflicts_with_current: true, sources: ["Checkout PRD addendum"] }],
    generated_at: "2026-03-01T09:00:00.000Z",
  });
  assert.deepEqual(discovery.resolved_context, ["PROJECT.md", "domains/checkout/workflows/place-order.md"]);
  assert.equal(discovery.findings[0]?.classification, "contradiction");
  const discoveryPath = await writeProductKnowledgeDiscovery(root, discovery);
  assert.match(discoveryPath, /\.runtime\/product-knowledge\//);
  assert.equal(await readFile(join(root, "context", workflowTarget.replace("context/", "")), "utf8"), before);

  // (c) A plan declares Product Knowledge impact and is approved.
  const draft: PlanDraftRequest = {
    contract_version: 1, plan_id: "checkout-retry", title: "Checkout retry", source: { kind: "prd", reference: "docs/retry.md" },
    work_prefix: "CKO", summary: "Retry a declined card once.", affected_repositories: ["frontend"], assumptions: [], open_questions: [],
    requirements: ["Retry once after a declined card."], solution: ["Add a single retry."], delivery: ["Implement then verify."], verification: ["Independent verification."], risks: [],
    work_items: [{ key: "retry", title: "Add retry", area: "checkout", repository: "frontend", scope: ["src/App.tsx"], test_scope: ["src/App.test.tsx"], test_policy: "required", verification_commands: ["npm test"], acceptance_criteria: ["A declined card retries once."] }],
    product_knowledge: { impact: "behavior-change", references: [workflowTarget], proposed_change: "Retry a declined card once before failing." },
  };
  const created = await createPlanDraft(root, draft, new Date("2026-03-01T08:00:00Z"));
  assert.deepEqual((await validatePlanDirectory(created.directory)).errors, []);
  const approved = await setPlanState(created.directory, { kind: "approve", approved_by: "maintainer" }, new Date("2026-03-01T08:30:00Z"));
  assert.equal(approved.product_knowledge?.impact, "behavior-change");

  // (d) A bounded, immutable task context package resolves only the referenced page.
  const pkg = await buildTaskContextPackage({ workspaceRoot: root, references: approved.product_knowledge!.references, impact: approved.product_knowledge!.impact, proposed_change: approved.product_knowledge!.proposed_change ?? null });
  assert.deepEqual(pkg.context_paths, [workflowTarget]);
  assert.match(pkg.content_digest, /^sha256:[a-f0-9]{64}$/);

  // (e) Unexpected impact blocks silent synchronization.
  const closeout = resolveProductKnowledgeCloseout([discovery.findings[0]?.classification === "contradiction" ? "contradicts-current" : undefined]);
  assert.equal(closeout.synchronization, "pending-review");

  // (f) Effective synchronization is gated on the confirming role and updates only the referenced page.
  await assert.rejects(synchronizeProductKnowledge({ workspaceRoot: root, confirming_role: "product-owner", confirmed_effective: false, updates: [{ target: workflowTarget, content: effectiveWorkflow }], synced_at: "2026-03-01T10:00:00.000Z" }), /remains proposed/);
  const record = await synchronizeProductKnowledge({ workspaceRoot: root, confirming_role: "product-owner", confirmed_effective: true, updates: [{ target: workflowTarget, content: effectiveWorkflow, source_contribution: "contributions/general/20260301T000000Z-maintainer-retry.md" }], synced_at: "2026-03-01T10:00:00.000Z" });
  assert.deepEqual(record.updated_pages, [workflowTarget]);
  assert.match(await readFile(join(root, "context", workflowTarget.replace("context/", "")), "utf8"), /retrying once after a declined card/);
  assert.deepEqual((await validateProductKnowledgeTree(join(root, "context"))).errors, []);

  // (g) A reproducible onboarding pack is a generated view.
  const first = await generateOnboardingPack({ workspaceRoot: root, roles: ["shopper"], revision: "rev-2", generated_at: "2026-03-01T11:00:00.000Z" });
  const second = await generateOnboardingPack({ workspaceRoot: root, roles: ["shopper"], revision: "rev-2", generated_at: "2026-03-01T11:00:00.000Z" });
  assert.equal(first.manifest.generated_view, true);
  assert.equal(first.markdown, second.markdown);
  assert.ok(first.manifest.included_paths.includes("context/roles/shopper.md"));
});

test("an existing wrapper without Product Knowledge stays valid", async () => {
  const root = await mkdtemp(join(tmpdir(), "pk-e2e-absent-"));
  try {
    const result = await validateProductKnowledgeTree(join(root, "context"));
    assert.deepEqual(result, { present: false, pages: 0, errors: [] });
  } finally {
    await rm(root, { recursive: true, force: true });
  }
});

test("host adapters remain thin delegates to the canonical Product Knowledge behavior", async () => {
  for (const skill of ["gather-context", "sync-context", "create-plan", "initialize-workspace"]) {
    const codex = await readFile(join(projectRoot, ".codex", "skills", skill, "SKILL.md"), "utf8");
    assert.match(codex, new RegExp(`\\.agents/skills/${skill}/SKILL\\.md`));
    const claude = await readFile(join(projectRoot, ".claude", "commands", `${skill}.md`), "utf8");
    assert.match(claude, new RegExp(`\\.agents/skills/${skill}/SKILL\\.md`));
  }
  // Every Product Knowledge contract is a required workspace document.
  const { requiredWorkspaceDocuments } = await import("../scripts/lib/validation.js");
  const required = requiredWorkspaceDocuments as readonly string[];
  for (const contract of ["product-knowledge-project", "product-knowledge-role", "product-knowledge-workflow", "product-knowledge-domain", "product-knowledge-candidate", "task-context-package", "product-knowledge-sync-record", "onboarding-pack"]) {
    assert.ok(required.includes(`.agents/contracts/${contract}.schema.json`), contract);
    await access(join(projectRoot, ".agents", "contracts", `${contract}.schema.json`));
  }
});
