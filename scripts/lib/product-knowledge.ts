import { access, readdir, readFile } from "node:fs/promises";
import { dirname, join, relative, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { validateContract } from "./validation.js";
import type { SchemaName } from "./validation.js";

export type ProductKnowledgeKind = "role" | "workflow" | "domain" | "product-map";

export interface ProductKnowledgeValidationResult {
  present: boolean;
  pages: number;
  errors: string[];
}

const schemaByKind: Record<ProductKnowledgeKind, SchemaName> = {
  role: "product-knowledge-role",
  workflow: "product-knowledge-workflow",
  domain: "product-knowledge-domain",
  "product-map": "product-knowledge-project",
};

// Sections each page kind must own. Authority never overlaps: workflow pages own
// exact behavior, role pages own the cross-domain story, and domain summaries own
// business-area routing.
const requiredSections: Record<ProductKnowledgeKind, string[]> = {
  role: ["Role definition", "Primary outcomes", "Product surfaces", "End-to-end role story", "Related workflows", "Role-specific behavior", "Limitations"],
  workflow: ["Outcome", "Actors", "Entry points", "Current flow", "Variations", "Business rules"],
  domain: ["Summary", "Workflows"],
  "product-map": ["Roles", "Domains"],
};

// Sections that would usurp another page's authority. Only workflow pages may own
// exact current flow and business rules.
const forbiddenSections: Record<ProductKnowledgeKind, string[]> = {
  role: ["Current flow", "Business rules"],
  workflow: [],
  domain: ["Current flow", "Business rules"],
  "product-map": ["Current flow", "Business rules"],
};

// Frontmatter reference fields that must resolve to a real page inside the tree.
const referenceFields: Record<ProductKnowledgeKind, string[]> = {
  role: ["relevant_domains", "related_workflows"],
  workflow: [],
  domain: ["workflows"],
  "product-map": ["roles", "domains"],
};

interface ParsedPage {
  frontmatter: Record<string, unknown> | null;
  body: string;
}

function parsePage(raw: string): ParsedPage {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---(?:\r?\n([\s\S]*))?$/);
  if (!match) return { frontmatter: null, body: raw };
  let frontmatter: Record<string, unknown> | null = null;
  try {
    const parsed = parseYaml(match[1]!);
    frontmatter = parsed && typeof parsed === "object" && !Array.isArray(parsed) ? (parsed as Record<string, unknown>) : null;
  } catch {
    frontmatter = null;
  }
  return { frontmatter, body: match[2] ?? "" };
}

function sectionTitles(body: string): Set<string> {
  const titles = new Set<string>();
  for (const line of body.replace(/\r\n/g, "\n").split("\n")) {
    const heading = line.match(/^##\s+(.+?)\s*$/);
    if (heading) titles.add(heading[1]!);
  }
  return titles;
}

function relativeLinkTargets(body: string): string[] {
  const targets: string[] = [];
  for (const match of body.matchAll(/\]\(([^)]+)\)/g)) {
    const target = match[1]!.trim().split("#")[0]!.split(/\s+/)[0]!;
    if (!target) continue;
    if (/^[a-z][a-z0-9+.-]*:/i.test(target) || target.startsWith("//") || target.startsWith("/")) continue;
    targets.push(target);
  }
  return targets;
}

async function isFile(path: string): Promise<boolean> {
  try {
    await access(path);
    return true;
  } catch {
    return false;
  }
}

async function directoryExists(path: string): Promise<boolean> {
  try {
    const entries = await readdir(path, { withFileTypes: true });
    return Array.isArray(entries);
  } catch {
    return false;
  }
}

async function resolveReference(root: string, pagePath: string, reference: string, label: string, errors: string[]): Promise<void> {
  const target = resolve(dirname(pagePath), reference);
  const within = relative(root, target);
  if (within.startsWith("..")) {
    errors.push(`${relative(root, pagePath)}: ${label} escapes the Product Knowledge tree: ${reference}`);
    return;
  }
  if (!(await isFile(target))) errors.push(`${relative(root, pagePath)}: ${label} does not resolve: ${reference}`);
}

async function validatePage(root: string, pagePath: string, kind: ProductKnowledgeKind, errors: string[]): Promise<void> {
  const rel = relative(root, pagePath);
  const { frontmatter, body } = parsePage(await readFile(pagePath, "utf8"));
  if (!frontmatter) {
    errors.push(`${rel}: missing YAML frontmatter for ${kind} page`);
    return;
  }
  if (frontmatter.kind !== kind) {
    errors.push(`${rel}: frontmatter kind must be '${kind}' but is '${String(frontmatter.kind)}'`);
    return;
  }
  for (const error of await validateContract(schemaByKind[kind], frontmatter)) {
    errors.push(`${rel}: metadata ${error.instancePath || "/"} ${error.message}`);
  }
  const titles = sectionTitles(body);
  for (const section of requiredSections[kind]) {
    if (!titles.has(section)) errors.push(`${rel}: missing required section '## ${section}'`);
  }
  for (const section of forbiddenSections[kind]) {
    if (titles.has(section)) errors.push(`${rel}: section '## ${section}' belongs to workflow pages, not ${kind} pages`);
  }
  for (const field of referenceFields[kind]) {
    const references = frontmatter[field];
    if (!Array.isArray(references)) continue;
    for (const reference of references) {
      if (typeof reference === "string") await resolveReference(root, pagePath, reference, `metadata ${field}`, errors);
    }
  }
  for (const link of relativeLinkTargets(body)) {
    await resolveReference(root, pagePath, link, "relative link", errors);
  }
}

async function markdownFiles(directory: string): Promise<string[]> {
  const entries = await readdir(directory, { withFileTypes: true });
  return entries.filter((entry) => entry.isFile() && entry.name.endsWith(".md")).map((entry) => join(directory, entry.name));
}

/**
 * Validate an optional Product Knowledge tree rooted at a context directory.
 *
 * The tree is discovered from `roles/` and `domains/`. When neither exists the
 * result is `present: false` with no errors, so an existing wrapper that has not
 * adopted Product Knowledge stays valid. Every page that does exist is validated
 * against its metadata contract, required and forbidden sections, and resolvable
 * references, so adoption can grow one page at a time.
 */
export async function validateProductKnowledgeTree(contextDir: string): Promise<ProductKnowledgeValidationResult> {
  const root = resolve(contextDir);
  const rolesDir = join(root, "roles");
  const domainsDir = join(root, "domains");
  const hasRoles = await directoryExists(rolesDir);
  const hasDomains = await directoryExists(domainsDir);
  if (!hasRoles && !hasDomains) return { present: false, pages: 0, errors: [] };

  const errors: string[] = [];
  let pages = 0;

  const projectPath = join(root, "PROJECT.md");
  if (await isFile(projectPath)) {
    const { frontmatter } = parsePage(await readFile(projectPath, "utf8"));
    if (frontmatter && frontmatter.kind === "product-map") {
      pages += 1;
      await validatePage(root, projectPath, "product-map", errors);
    }
  }

  if (hasRoles) {
    for (const file of await markdownFiles(rolesDir)) {
      if (file.endsWith("README.md")) {
        for (const link of relativeLinkTargets(await readFile(file, "utf8"))) await resolveReference(root, file, link, "relative link", errors);
        continue;
      }
      pages += 1;
      await validatePage(root, file, "role", errors);
    }
  }

  if (hasDomains) {
    const domainEntries = (await readdir(domainsDir, { withFileTypes: true })).filter((entry) => entry.isDirectory());
    for (const entry of domainEntries) {
      const domainDir = join(domainsDir, entry.name);
      const readmePath = join(domainDir, "README.md");
      if (!(await isFile(readmePath))) {
        errors.push(`domains/${entry.name}: missing README.md domain summary`);
      } else {
        pages += 1;
        await validatePage(root, readmePath, "domain", errors);
      }
      const workflowsDir = join(domainDir, "workflows");
      if (await directoryExists(workflowsDir)) {
        for (const file of await markdownFiles(workflowsDir)) {
          pages += 1;
          await validatePage(root, file, "workflow", errors);
        }
      }
    }
  }

  return { present: true, pages, errors: [...new Set(errors)] };
}
