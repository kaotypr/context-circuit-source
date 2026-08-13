import assert from "node:assert/strict";
import { cp, mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { generateOnboardingPack } from "../scripts/lib/product-knowledge-onboarding.js";
import { validateContract } from "../scripts/lib/validation.js";
import { projectRoot } from "./helpers.js";

const smallFixture = join(projectRoot, "fixtures", "product-knowledge-small");
const generated_at = "2026-03-01T10:00:00.000Z";

async function packWorkspace(): Promise<{ root: string; cleanup: () => Promise<void> }> {
  const root = await mkdtemp(join(tmpdir(), "pk-onboarding-"));
  await cp(smallFixture, join(root, "context"), { recursive: true });
  await writeFile(join(root, "context", "ARCHITECTURE.md"), "# Architecture\n\nThe storefront is a single web app.\n", "utf8");
  return { root, cleanup: async () => rm(root, { recursive: true, force: true }) };
}

test("an onboarding pack selects the product map, role, referenced domain and workflow, and architecture", async (t) => {
  const ws = await packWorkspace();
  t.after(ws.cleanup);
  const pack = await generateOnboardingPack({ workspaceRoot: ws.root, roles: ["shopper"], revision: "rev-1", generated_at });
  assert.deepEqual(await validateContract("onboarding-pack", pack.manifest), []);
  assert.equal(pack.manifest.generated_view, true);
  assert.deepEqual(pack.manifest.included_paths, [
    "context/PROJECT.md",
    "context/ARCHITECTURE.md",
    "context/roles/shopper.md",
    "context/domains/checkout/README.md",
    "context/domains/checkout/workflows/place-order.md",
  ]);
  assert.deepEqual(pack.manifest.roles, ["shopper"]);
  // Known gaps from the included pages are surfaced.
  assert.ok(pack.manifest.known_gaps.some((gap) => /Returns and refunds|Guest vs registered|Split payments|declined card/.test(gap)));
  // The markdown is clearly a generated view and cites its sources.
  assert.match(pack.markdown, /generated view/i);
  assert.match(pack.markdown, /not[\s\S]*an independent source of truth/i);
  assert.match(pack.markdown, /Source: `context\/domains\/checkout\/workflows\/place-order\.md`/);
});

test("an onboarding pack is reproducible from one revision", async (t) => {
  const ws = await packWorkspace();
  t.after(ws.cleanup);
  const first = await generateOnboardingPack({ workspaceRoot: ws.root, roles: ["shopper"], revision: "rev-1", generated_at });
  const second = await generateOnboardingPack({ workspaceRoot: ws.root, roles: ["shopper"], revision: "rev-1", generated_at });
  assert.equal(first.markdown, second.markdown);
  assert.deepEqual(first.manifest, second.manifest);
});

test("an onboarding pack requires each requested role to exist", async (t) => {
  const ws = await packWorkspace();
  t.after(ws.cleanup);
  await assert.rejects(generateOnboardingPack({ workspaceRoot: ws.root, roles: ["ghost"], revision: "rev-1", generated_at }), /role page does not exist/);
});
