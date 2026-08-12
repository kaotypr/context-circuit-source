import assert from "node:assert/strict";
import { cp, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { parse as parseYaml } from "yaml";
import { validateProductKnowledgeTree } from "../scripts/lib/product-knowledge.js";
import { projectRoot } from "./helpers.js";

const smallFixture = join(projectRoot, "fixtures", "product-knowledge-small");
const largeFixture = join(projectRoot, "fixtures", "product-knowledge-large");
const template = join(projectRoot, ".agents", "templates", "product-knowledge");

async function withMutableTree(source: string, mutate: (root: string) => Promise<void>): Promise<Awaited<ReturnType<typeof validateProductKnowledgeTree>>> {
  const root = await mkdtemp(join(tmpdir(), "pk-tree-"));
  try {
    await cp(source, root, { recursive: true });
    await mutate(root);
    return await validateProductKnowledgeTree(root);
  } finally {
    await rm(root, { recursive: true, force: true });
  }
}

function frontmatter(raw: string): Record<string, unknown> {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---/);
  return match ? (parseYaml(match[1]!) as Record<string, unknown>) : {};
}

test("the small fixture is a valid Product Knowledge tree", async () => {
  const result = await validateProductKnowledgeTree(smallFixture);
  assert.equal(result.present, true);
  assert.deepEqual(result.errors, []);
  assert.ok(result.pages >= 4);
});

test("the large multi-domain fixture validates", async () => {
  const result = await validateProductKnowledgeTree(largeFixture);
  assert.equal(result.present, true);
  assert.deepEqual(result.errors, []);
  // 1 product map + 4 roles + 5 domains + 10 workflows.
  assert.equal(result.pages, 20);
});

test("the shipped template is itself a valid Product Knowledge tree", async () => {
  const result = await validateProductKnowledgeTree(template);
  assert.equal(result.present, true);
  assert.deepEqual(result.errors, []);
});

test("an agent can select a domain from the product map without loading the whole tree", async () => {
  const opened: string[] = [];
  const open = async (relative: string): Promise<string> => {
    opened.push(relative);
    return readFile(join(largeFixture, relative), "utf8");
  };
  // Read only the compact product map, then route to exactly one domain.
  const map = frontmatter(await open("PROJECT.md"));
  const domains = map.domains as string[];
  assert.ok(domains.length >= 5);
  const chosen = domains.find((d) => d.includes("checkout"))!.replace(/^\.\//, "");
  const domain = frontmatter(await open(chosen));
  for (const workflow of domain.workflows as string[]) {
    await open(join("domains", "checkout", workflow));
  }
  // The bounded selection never touched other domains or any role page.
  assert.ok(!opened.some((path) => path.includes("catalog") || path.includes("payments") || path.includes("roles")));
  assert.ok(opened.length < 6);
});

test("a directory with no roles or domains is treated as absent and stays valid", async () => {
  const root = await mkdtemp(join(tmpdir(), "pk-empty-"));
  try {
    await writeFile(join(root, "PROJECT.md"), "# Project\n\nNot initialized.\n");
    const result = await validateProductKnowledgeTree(root);
    assert.deepEqual(result, { present: false, pages: 0, errors: [] });
  } finally {
    await rm(root, { recursive: true, force: true });
  }
});

test("a missing required section is reported", async () => {
  const result = await withMutableTree(smallFixture, async (root) => {
    const path = join(root, "domains", "checkout", "workflows", "place-order.md");
    const raw = await readFile(path, "utf8");
    await writeFile(path, raw.replace(/## Business rules[\s\S]*$/, ""));
  });
  assert.ok(result.errors.some((error) => error.includes("place-order.md") && error.includes("Business rules")));
});

test("missing owners metadata is reported", async () => {
  const result = await withMutableTree(smallFixture, async (root) => {
    const path = join(root, "roles", "shopper.md");
    const raw = await readFile(path, "utf8");
    await writeFile(path, raw.replace(/owners:\n  - Storefront product team\n/, "owners: []\n"));
  });
  assert.ok(result.errors.some((error) => error.includes("shopper.md") && error.includes("owners")));
});

test("a broken relative reference is reported", async () => {
  const result = await withMutableTree(smallFixture, async (root) => {
    const path = join(root, "domains", "checkout", "README.md");
    const raw = await readFile(path, "utf8");
    await writeFile(path, raw.replace("workflows/place-order.md", "workflows/missing.md"));
  });
  assert.ok(result.errors.some((error) => error.includes("missing.md") && error.includes("does not resolve")));
});

test("a role page that claims workflow authority is reported", async () => {
  const result = await withMutableTree(smallFixture, async (root) => {
    const path = join(root, "roles", "shopper.md");
    const raw = await readFile(path, "utf8");
    await writeFile(path, `${raw}\n## Business rules\n\nRules that belong on a workflow page.\n`);
  });
  assert.ok(result.errors.some((error) => error.includes("shopper.md") && error.includes("Business rules") && error.includes("workflow")));
});

test("a page whose frontmatter kind mismatches its location is reported", async () => {
  const result = await withMutableTree(smallFixture, async (root) => {
    const path = join(root, "roles", "shopper.md");
    const raw = await readFile(path, "utf8");
    await writeFile(path, raw.replace("kind: role", "kind: workflow"));
  });
  assert.ok(result.errors.some((error) => error.includes("shopper.md") && error.includes("kind")));
});
