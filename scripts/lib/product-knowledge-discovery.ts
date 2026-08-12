import { access, mkdir } from "node:fs/promises";
import { join, relative, resolve } from "node:path";
import { assertInside, writeTextAtomic } from "./io.js";
import { validateContract } from "./validation.js";
import { containsSecret, contextReferenceError } from "./safe-reference.js";

export type ProductKnowledgeFinding = "candidate" | "contradiction" | "gap";

export interface ProductKnowledgeObservation {
  topic: string;
  summary: string;
  target?: string | null;
  conflicts_with_current?: boolean;
  sources: string[];
}

export interface DiscoverProductKnowledgeInput {
  contextDir: string;
  revision: string;
  topic: string;
  select?: string[];
  observations: ProductKnowledgeObservation[];
  generated_at: string;
}

export interface ProductKnowledgeDiscoveryResult {
  contract_version: 1;
  generated_at: string;
  topic: string;
  revision: string;
  resolved_context: string[];
  findings: Array<{ classification: ProductKnowledgeFinding; topic: string; summary: string; target: string | null; sources: string[] }>;
}

async function pageExists(root: string, reference: string): Promise<boolean> {
  const target = resolve(root, reference);
  if (relative(root, target).startsWith("..")) return false;
  try {
    await access(target);
    return true;
  } catch {
    return false;
  }
}

function assertCredentialFree(label: string, value: string): void {
  if (containsSecret(value)) throw new Error(`${label} appears to contain credentials and cannot enter Product Knowledge discovery`);
}

/**
 * Reconcile host-gathered observations against the canonical Product Knowledge
 * tree without mutating it. Discovery resolves only the selected, relevant pages
 * (plus the product map) and classifies every observation as a candidate addition,
 * a contradiction of current behavior, or a gap. Credentials and secrets are
 * rejected before anything is recorded.
 */
export async function discoverProductKnowledge(input: DiscoverProductKnowledgeInput): Promise<ProductKnowledgeDiscoveryResult> {
  const root = resolve(input.contextDir);
  if (!input.topic.trim()) throw new Error("Discovery requires a non-empty topic");
  if (!input.revision.trim()) throw new Error("Discovery requires a baseline revision or digest");

  const resolved = new Set<string>();
  if (await pageExists(root, "PROJECT.md")) resolved.add("PROJECT.md");
  for (const reference of input.select ?? []) {
    const error = contextReferenceError(reference);
    if (error) throw new Error(`Selected context path ${reference} ${error}`);
    if (await pageExists(root, reference)) resolved.add(reference);
  }

  const findings: ProductKnowledgeDiscoveryResult["findings"] = [];
  for (const observation of input.observations) {
    if (!observation.topic.trim()) throw new Error("Every observation requires a non-empty topic");
    if (!observation.summary.trim()) throw new Error("Every observation requires a non-empty summary");
    assertCredentialFree(`Observation '${observation.topic}' summary`, observation.summary);
    for (const source of observation.sources) {
      if (!source.trim() || /[\r\n]/.test(source)) throw new Error(`Observation '${observation.topic}' has an empty or multi-line source reference`);
      assertCredentialFree(`Observation '${observation.topic}' source`, source);
    }
    const target = observation.target ?? null;
    const targetResolves = Boolean(target) && (await pageExists(root, target!));
    let classification: ProductKnowledgeFinding;
    if (observation.conflicts_with_current) {
      classification = targetResolves ? "contradiction" : "gap";
    } else {
      classification = targetResolves ? "candidate" : "gap";
    }
    findings.push({ classification, topic: observation.topic, summary: observation.summary, target, sources: [...observation.sources] });
  }

  return {
    contract_version: 1,
    generated_at: input.generated_at,
    topic: input.topic,
    revision: input.revision,
    resolved_context: [...resolved].sort(),
    findings,
  };
}

/**
 * Persist a discovery result as ignored runtime evidence. The artifact is written
 * under `.runtime/product-knowledge/`; canonical context is never touched.
 */
export async function writeProductKnowledgeDiscovery(workspaceRoot: string, result: ProductKnowledgeDiscoveryResult): Promise<string> {
  const errors = await validateContract("product-knowledge-candidate", result);
  if (errors.length > 0) throw new Error(`Invalid Product Knowledge discovery result: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);
  const root = resolve(workspaceRoot);
  const directory = assertInside(root, join(root, ".runtime", "product-knowledge"));
  await mkdir(directory, { recursive: true });
  const stamp = result.generated_at.replace(/[:.]/g, "-");
  const slug = result.topic.trim().toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "").slice(0, 48) || "discovery";
  const path = assertInside(directory, join(directory, `${stamp}-${slug}.json`));
  await writeTextAtomic(path, `${JSON.stringify(result, null, 2)}\n`);
  return path;
}
