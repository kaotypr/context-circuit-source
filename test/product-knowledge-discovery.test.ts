import assert from "node:assert/strict";
import { access, mkdtemp, readdir, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { discoverProductKnowledge, writeProductKnowledgeDiscovery } from "../scripts/lib/product-knowledge-discovery.js";
import { validateProductKnowledgeTree } from "../scripts/lib/product-knowledge.js";
import { validateContract } from "../scripts/lib/validation.js";
import { projectRoot } from "./helpers.js";

const contextDir = join(projectRoot, "fixtures", "product-knowledge-small");
const generated_at = "2026-03-01T10:00:00.000Z";

test("discovery resolves only the selected relevant pages plus the product map", async () => {
  const result = await discoverProductKnowledge({
    contextDir,
    revision: "abc123",
    topic: "order payment retry",
    select: ["domains/checkout/workflows/place-order.md"],
    observations: [],
    generated_at,
  });
  assert.deepEqual(result.resolved_context, ["PROJECT.md", "domains/checkout/workflows/place-order.md"]);
  assert.ok(!result.resolved_context.includes("roles/shopper.md"));
});

test("discovery classifies contradictions, candidates, and gaps", async () => {
  const result = await discoverProductKnowledge({
    contextDir,
    revision: "abc123",
    topic: "checkout changes",
    select: [],
    observations: [
      { topic: "declined card retry", summary: "The order now auto-retries once after a declined card.", target: "domains/checkout/workflows/place-order.md", conflicts_with_current: true, sources: ["Checkout PRD addendum"] },
      { topic: "guest express checkout", summary: "A one-tap express path exists for returning shoppers.", target: "domains/checkout/workflows/place-order.md", sources: ["Checkout PRD addendum"] },
      { topic: "returns", summary: "Shoppers can request a return within 30 days.", target: null, sources: ["Returns policy"] },
      { topic: "loyalty program", summary: "Points accrue per order.", target: "domains/loyalty/README.md", conflicts_with_current: true, sources: ["Loyalty brief"] },
    ],
    generated_at,
  });
  const byTopic = Object.fromEntries(result.findings.map((finding) => [finding.topic, finding.classification]));
  assert.equal(byTopic["declined card retry"], "contradiction");
  assert.equal(byTopic["guest express checkout"], "candidate");
  assert.equal(byTopic["returns"], "gap");
  // A conflict against a page that does not exist is a gap, not a contradiction.
  assert.equal(byTopic["loyalty program"], "gap");
});

test("discovery rejects credentials in summaries and sources", async () => {
  await assert.rejects(discoverProductKnowledge({
    contextDir,
    revision: "abc123",
    topic: "leak",
    observations: [{ topic: "leak", summary: "The api_key=sk-secret is used to call the provider.", target: null, sources: ["notes"] }],
    generated_at,
  }), /credentials/);
  await assert.rejects(discoverProductKnowledge({
    contextDir,
    revision: "abc123",
    topic: "leak",
    observations: [{ topic: "leak", summary: "Uses the provider.", target: null, sources: ["token=abcdef1234"] }],
    generated_at,
  }), /credentials/);
});

test("discovery writes ignored runtime evidence and never mutates canonical context", async () => {
  const workspace = await mkdtemp(join(tmpdir(), "pk-discovery-"));
  try {
    const before = await validateProductKnowledgeTree(contextDir);
    const result = await discoverProductKnowledge({
      contextDir,
      revision: "abc123",
      topic: "returns gap",
      observations: [{ topic: "returns", summary: "Returns are not documented.", target: null, sources: ["Returns policy"] }],
      generated_at,
    });
    const path = await writeProductKnowledgeDiscovery(workspace, result);
    assert.match(path, /\.runtime\/product-knowledge\/.*returns-gap\.json$/);
    await access(path);
    assert.deepEqual(await validateContract("product-knowledge-candidate", result), []);
    // The canonical tree is unchanged and nothing was written under context/.
    const after = await validateProductKnowledgeTree(contextDir);
    assert.deepEqual(after, before);
    await assert.rejects(access(join(workspace, "context")));
    assert.deepEqual((await readdir(workspace)).sort(), [".runtime"]);
  } finally {
    await rm(workspace, { recursive: true, force: true });
  }
});
