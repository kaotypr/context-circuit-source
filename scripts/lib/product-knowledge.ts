import { createHash } from "node:crypto";
import { access, mkdir, readdir, readFile, rm } from "node:fs/promises";
import { dirname, join, relative, resolve } from "node:path";
import { parse as parseYaml, stringify as stringifyYaml } from "yaml";
import { writeTextAtomic } from "./io.js";
import { validateContract } from "./validation.js";
import type { SchemaName } from "./validation.js";
import type { ProductKnowledgeBaselineSpec, ProductKnowledgeImpact, TaskContextPackage, WorkspaceConfig } from "./types.js";

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

export interface ProductKnowledgeSyncUpdate {
  target: string;
  content: string;
  source_contribution?: string;
}

export interface ProductKnowledgeSyncInput {
  workspaceRoot: string;
  confirming_role: string;
  confirmed_effective: boolean;
  updates: ProductKnowledgeSyncUpdate[];
  synced_at: string;
}

export interface ProductKnowledgeSyncRecord {
  contract_version: 1;
  synced_at: string;
  confirming_role: string;
  updated_pages: string[];
  source_contributions: string[];
}

function currentBehaviorPageKind(contextRelative: string): "role" | "domain" | "workflow" | null {
  const normalized = contextRelative.replace(/\\/g, "/");
  if (/^roles\/(?!README\.md$)[^/]+\.md$/.test(normalized)) return "role";
  if (/^domains\/[^/]+\/README\.md$/.test(normalized)) return "domain";
  if (/^domains\/[^/]+\/workflows\/[^/]+\.md$/.test(normalized)) return "workflow";
  return null;
}

/**
 * Synchronize canonical Product Knowledge only after a dedicated, workspace
 * configured confirming role declares the behavior effective. This is the single
 * gate that writes current-behavior pages: without an explicit effectiveness
 * confirmation from the configured role, nothing is written and the behavior
 * remains proposed. Only referenced role, domain, and workflow pages may change,
 * and the resulting tree must validate or every write is rolled back.
 */
export async function synchronizeProductKnowledge(input: ProductKnowledgeSyncInput): Promise<ProductKnowledgeSyncRecord> {
  const root = resolve(input.workspaceRoot);
  const config = parseYaml(await readFile(join(root, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  const confirmingRole = config.product_knowledge?.confirming_role;
  if (!confirmingRole) throw new Error("Workspace has no configured Product Knowledge confirming role; synchronization is not permitted");
  if (!input.confirmed_effective) throw new Error("Product Knowledge synchronization requires an explicit effectiveness confirmation; the behavior remains proposed");
  if (input.confirming_role !== confirmingRole) throw new Error(`Confirming role '${input.confirming_role}' is not the configured Product Knowledge confirming role '${confirmingRole}'`);
  if (input.updates.length === 0) throw new Error("Synchronization requires at least one referenced current-behavior page");

  const contextDir = join(root, "context");
  const targets = new Map<string, { absolute: string; previous: string | null }>();
  for (const update of input.updates) {
    const contextRelative = update.target.replace(/^context\//, "");
    if (!currentBehaviorPageKind(contextRelative)) throw new Error(`Synchronization target is not a current-behavior page: ${update.target}`);
    const absolute = resolve(contextDir, contextRelative);
    if (relative(contextDir, absolute).startsWith("..")) throw new Error(`Synchronization target escapes context: ${update.target}`);
    if (targets.has(absolute)) throw new Error(`Duplicate synchronization target: ${update.target}`);
    let previous: string | null = null;
    try {
      previous = await readFile(absolute, "utf8");
    } catch {
      previous = null;
    }
    targets.set(absolute, { absolute, previous });
  }

  const applied: string[] = [];
  try {
    for (const update of input.updates) {
      const contextRelative = update.target.replace(/^context\//, "");
      const absolute = resolve(contextDir, contextRelative);
      await mkdir(dirname(absolute), { recursive: true });
      await writeTextAtomic(absolute, update.content);
      applied.push(absolute);
    }
    const validation = await validateProductKnowledgeTree(contextDir);
    if (validation.errors.length > 0) throw new Error(`Synchronized Product Knowledge is invalid: ${validation.errors.join("; ")}`);
  } catch (error) {
    for (const absolute of applied) {
      const snapshot = targets.get(absolute)!;
      if (snapshot.previous === null) await rm(absolute, { force: true });
      else await writeTextAtomic(absolute, snapshot.previous);
    }
    throw error;
  }

  const record: ProductKnowledgeSyncRecord = {
    contract_version: 1,
    synced_at: input.synced_at,
    confirming_role: input.confirming_role,
    updated_pages: input.updates.map((update) => update.target),
    source_contributions: [...new Set(input.updates.map((update) => update.source_contribution).filter((value): value is string => Boolean(value)))],
  };
  const errors = await validateContract("product-knowledge-sync-record", record);
  if (errors.length > 0) throw new Error(`Invalid Product Knowledge sync record: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  return record;
}

/**
 * Resolve a compact, immutable Product Knowledge package for a task. Only the
 * referenced pages that actually exist are included (the bounded business
 * baseline, not the whole tree), and their content is pinned by a digest so a
 * worker and an independent verifier share the exact same baseline. Missing
 * references are proposed pages and are intentionally excluded from the snapshot.
 */
export async function buildTaskContextPackage(input: {
  workspaceRoot: string;
  references: string[];
  impact: ProductKnowledgeImpact;
  proposed_change?: string | null;
  revision?: string;
}): Promise<TaskContextPackage> {
  const root = resolve(input.workspaceRoot);
  const contents = new Map<string, string>();
  for (const reference of input.references) {
    const target = resolve(root, reference);
    if (relative(root, target).startsWith("..")) throw new Error(`Product Knowledge reference escapes the workspace: ${reference}`);
    try {
      contents.set(reference, await readFile(target, "utf8"));
    } catch {
      // A referenced page that does not exist yet is proposed, not part of the
      // current baseline snapshot.
    }
  }
  const contextPaths = [...contents.keys()].sort();
  const hash = createHash("sha256");
  for (const path of contextPaths) hash.update(`${path}\0${contents.get(path)}\0`);
  const contentDigest = `sha256:${hash.digest("hex")}`;
  return {
    contract_version: 1,
    revision: input.revision?.trim() || contentDigest,
    content_digest: contentDigest,
    context_paths: contextPaths,
    impact: input.impact,
    proposed_change: input.proposed_change ?? null,
  };
}

function slugify(name: string): string {
  const slug = name.trim().toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
  if (!slug) throw new Error(`Cannot derive a Product Knowledge slug from '${name}'`);
  return slug;
}

function frontmatterBlock(data: Record<string, unknown>): string {
  return `---\n${stringifyYaml(data).trimEnd()}\n---\n`;
}

const placeholder = "Not documented yet.";

/**
 * Render a minimal, reviewable Product Knowledge baseline as workspace-relative
 * files under `context/`. The baseline records what a human already knows (product
 * purpose, major roles, major domains, and their known workflows) and keeps every
 * unknown explicit rather than inventing detail. The output is a valid tree under
 * the PKNOW-001 contracts, so coverage can grow one page at a time.
 */
export function renderProductKnowledgeBaseline(spec: ProductKnowledgeBaselineSpec): Record<string, string> {
  const roleSlugs = new Map<string, string>();
  for (const role of spec.roles) {
    const slug = slugify(role);
    if ([...roleSlugs.values()].includes(slug)) throw new Error(`Duplicate role slug in baseline: ${slug}`);
    roleSlugs.set(role, slug);
  }
  const domainSlugs = new Map<string, string>();
  for (const domain of spec.domains) {
    const slug = slugify(domain.name);
    if ([...domainSlugs.values()].includes(slug)) throw new Error(`Duplicate domain slug in baseline: ${slug}`);
    domainSlugs.set(domain.name, slug);
  }
  const gaps = spec.unknowns.length > 0 ? spec.unknowns : ["No unknowns recorded yet."];
  const owners = ["Unassigned — record the owner."];
  const files: Record<string, string> = {};

  const roleList = spec.roles.length > 0
    ? spec.roles.map((role) => `- [${role}](roles/${roleSlugs.get(role)}.md)`).join("\n")
    : "- None documented yet.";
  const domainList = spec.domains.length > 0
    ? spec.domains.map((domain) => `- [${domain.name}](domains/${domainSlugs.get(domain.name)}/README.md)`).join("\n")
    : "- None documented yet.";
  files["context/PROJECT.md"] = `${frontmatterBlock({
    kind: "product-map",
    title: spec.title,
    roles: spec.roles.map((role) => `roles/${roleSlugs.get(role)}.md`),
    domains: spec.domains.map((domain) => `domains/${domainSlugs.get(domain.name)}/README.md`),
    sources: spec.sources,
    review_date: spec.review_date,
    known_gaps: gaps,
  })}\n# ${spec.title}\n\n${spec.purpose}\n\n## Roles\n\n${roleList}\n\n## Domains\n\n${domainList}\n`;

  files["context/GLOSSARY.md"] = "# Glossary\n\nDefine business terms here as they are confirmed.\n";
  const roleIndexList = spec.roles.length > 0
    ? spec.roles.map((role) => `- [${role}](${roleSlugs.get(role)}.md)`).join("\n")
    : "- None documented yet.";
  files["context/roles/README.md"] = `# Roles\n\n${roleIndexList}\n`;

  const relevantDomains = spec.domains.map((domain) => `../domains/${domainSlugs.get(domain.name)}/README.md`);
  const relatedWorkflows = spec.domains
    .filter((domain) => (domain.workflows ?? []).length > 0)
    .map((domain) => `../domains/${domainSlugs.get(domain.name)}/workflows/${slugify(domain.workflows![0]!)}.md`);
  const relatedList = relatedWorkflows.length > 0
    ? relatedWorkflows.map((reference) => `- [Workflow](${reference})`).join("\n")
    : "None documented yet.";
  for (const role of spec.roles) {
    files[`context/roles/${roleSlugs.get(role)}.md`] = `${frontmatterBlock({
      kind: "role",
      title: role,
      owners,
      sources: spec.sources,
      review_date: spec.review_date,
      relevant_domains: relevantDomains,
      related_workflows: relatedWorkflows,
      known_gaps: gaps,
    })}\n# ${role}\n\n## Role definition\n\n${placeholder}\n\n## Primary outcomes\n\n${placeholder}\n\n## Product surfaces\n\n${placeholder}\n\n## End-to-end role story\n\n${placeholder}\n\n## Related workflows\n\n${relatedList}\n\n## Role-specific behavior\n\n${placeholder}\n\n## Limitations\n\n${placeholder}\n`;
  }

  for (const domain of spec.domains) {
    const domainSlug = domainSlugs.get(domain.name)!;
    const workflows = domain.workflows ?? [];
    const workflowSlugs = workflows.map((workflow) => slugify(workflow));
    const workflowList = workflows.length > 0
      ? workflows.map((workflow, index) => `- [${workflow}](workflows/${workflowSlugs[index]}.md)`).join("\n")
      : "None documented yet.";
    files[`context/domains/${domainSlug}/README.md`] = `${frontmatterBlock({
      kind: "domain",
      title: domain.name,
      owners,
      sources: spec.sources,
      review_date: spec.review_date,
      workflows: workflowSlugs.map((slug) => `workflows/${slug}.md`),
      known_gaps: gaps,
    })}\n# ${domain.name}\n\n## Summary\n\n${placeholder}\n\n## Workflows\n\n${workflowList}\n`;
    workflows.forEach((workflow, index) => {
      files[`context/domains/${domainSlug}/workflows/${workflowSlugs[index]}.md`] = `${frontmatterBlock({
        kind: "workflow",
        title: workflow,
        owners,
        sources: spec.sources,
        review_date: spec.review_date,
        implementation_ownership: "Unassigned — record the implementing repository or team.",
        known_gaps: gaps,
      })}\n# ${workflow}\n\n## Outcome\n\n${placeholder}\n\n## Actors\n\n${placeholder}\n\n## Entry points\n\n${placeholder}\n\n## Current flow\n\n${placeholder}\n\n## Variations\n\n${placeholder}\n\n## Business rules\n\n${placeholder}\n`;
    });
  }

  return files;
}
