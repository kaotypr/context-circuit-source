import { access, readFile } from "node:fs/promises";
import { dirname, join, relative, resolve } from "node:path";
import { parse as parseYaml } from "yaml";
import { validateContract } from "./validation.js";

export interface OnboardingPackManifest {
  contract_version: 1;
  generated_at: string;
  revision: string;
  roles: string[];
  included_paths: string[];
  known_gaps: string[];
  generated_view: true;
}

export interface OnboardingPack {
  manifest: OnboardingPackManifest;
  markdown: string;
}

interface IncludedPage {
  path: string; // workspace-relative
  content: string;
  frontmatter: Record<string, unknown> | null;
}

function parseFrontmatter(raw: string): Record<string, unknown> | null {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---/);
  if (!match) return null;
  try {
    const parsed = parseYaml(match[1]!);
    return parsed && typeof parsed === "object" && !Array.isArray(parsed) ? (parsed as Record<string, unknown>) : null;
  } catch {
    return null;
  }
}

async function readIfPresent(root: string, workspaceRelative: string): Promise<IncludedPage | null> {
  const absolute = resolve(root, workspaceRelative);
  if (relative(root, absolute).startsWith("..")) return null;
  try {
    await access(absolute);
  } catch {
    return null;
  }
  const content = await readFile(absolute, "utf8");
  return { path: workspaceRelative, content, frontmatter: parseFrontmatter(content) };
}

function references(frontmatter: Record<string, unknown> | null, field: string): string[] {
  const value = frontmatter?.[field];
  return Array.isArray(value) ? value.filter((entry): entry is string => typeof entry === "string") : [];
}

/**
 * Generate a role- or domain-oriented onboarding pack from canonical Product
 * Knowledge. The pack is a generated view — reproducible from a single revision,
 * explicit about the source paths it includes, and clearly not an independent
 * source of truth. It selects the product map, the requested role pages, the
 * domains and workflows those roles reference, architecture context, and the
 * known gaps of every included page.
 */
export async function generateOnboardingPack(input: {
  workspaceRoot: string;
  roles: string[];
  revision: string;
  generated_at: string;
}): Promise<OnboardingPack> {
  const root = resolve(input.workspaceRoot);
  if (input.roles.length === 0) throw new Error("An onboarding pack requires at least one role");
  const roleSlugs = input.roles.map((role) => role.replace(/^context\/roles\//, "").replace(/\.md$/, ""));

  const included: IncludedPage[] = [];
  const seen = new Set<string>();
  const add = async (workspaceRelative: string): Promise<IncludedPage | null> => {
    if (seen.has(workspaceRelative)) return included.find((page) => page.path === workspaceRelative) ?? null;
    const page = await readIfPresent(root, workspaceRelative);
    if (!page) return null;
    seen.add(workspaceRelative);
    included.push(page);
    return page;
  };

  await add("context/PROJECT.md");
  await add("context/ARCHITECTURE.md");

  for (const slug of roleSlugs) {
    const rolePath = `context/roles/${slug}.md`;
    const rolePage = await add(rolePath);
    if (!rolePage) throw new Error(`Onboarding role page does not exist: ${rolePath}`);
    const roleDir = dirname(join(root, rolePath));
    for (const field of ["relevant_domains", "related_workflows"] as const) {
      for (const reference of references(rolePage.frontmatter, field)) {
        const resolved = resolve(roleDir, reference);
        const workspaceRelative = relative(root, resolved);
        if (!workspaceRelative.startsWith("..")) await add(workspaceRelative);
      }
    }
  }

  const knownGaps = [...new Set(included.flatMap((page) => references(page.frontmatter, "known_gaps")))];
  const includedPaths = included.map((page) => page.path);

  const manifest: OnboardingPackManifest = {
    contract_version: 1,
    generated_at: input.generated_at,
    revision: input.revision,
    roles: [...new Set(roleSlugs)],
    included_paths: includedPaths,
    known_gaps: knownGaps,
    generated_view: true,
  };
  const errors = await validateContract("onboarding-pack", manifest);
  if (errors.length > 0) throw new Error(`Invalid onboarding pack manifest: ${errors.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ")}`);

  const header = [
    "# Onboarding pack (generated view)",
    "",
    "> This document is a generated view assembled from Product Knowledge. It is not",
    "> an independent source of truth; the linked canonical pages remain authoritative.",
    "",
    `- Revision: \`${input.revision}\``,
    `- Generated at: ${input.generated_at}`,
    `- Roles: ${manifest.roles.join(", ")}`,
    `- Included pages: ${includedPaths.map((path) => `\`${path}\``).join(", ")}`,
    "",
    "## Known gaps",
    "",
    knownGaps.length > 0 ? knownGaps.map((gap) => `- ${gap}`).join("\n") : "- None recorded.",
  ].join("\n");
  const body = included.map((page) => `## Source: \`${page.path}\`\n\n${page.content.trim()}`).join("\n\n---\n\n");
  const markdown = `${header}\n\n---\n\n${body}\n`;

  return { manifest, markdown };
}
