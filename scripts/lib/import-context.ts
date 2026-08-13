import { lstat, mkdir, readFile, readdir, realpath } from "node:fs/promises";
import { dirname, join, relative, resolve, sep } from "node:path";
import { fileURLToPath } from "node:url";
import { Ajv2020 } from "ajv/dist/2020.js";
import { parse as parseYaml } from "yaml";
import { assertCleanRepository, git } from "./git.js";
import { assertInside, ensurePrivateDirectory, writeJsonAtomic, writeTextExclusive } from "./io.js";
import type { ImportContextEvidence, ImportContextLimits, ImportContextManifest, ImportContextRequest, ImportContextResult, WorkspaceConfig } from "./types.js";
import { validateContract, workspaceSemanticErrors } from "./validation.js";

const defaults: Required<ImportContextLimits> = { max_files: 50, max_file_bytes: 128 * 1024, max_total_bytes: 1024 * 1024 };
const rootInstructions = ["README.md", "README", "AGENTS.md", "CLAUDE.md", "CONTRIBUTING.md", "DEVELOPING.md"];
const documentationRoots = ["docs", "doc", "documentation"];
const repositoryInstructions = new Set(["AGENTS.md", "CLAUDE.md", "CONTRIBUTING.md", "DEVELOPING.md"]);
const structuralFiles = new Set(["package.json", "pyproject.toml", "Cargo.toml", "go.mod", "pom.xml", "build.gradle", "Makefile", "Dockerfile", "docker-compose.yml", "docker-compose.yaml"]);
const excludedDiscoveryDirectories = new Set([".git", "node_modules", "vendor", "dist", "build", ".runtime"]);
const projectRoot = resolve(dirname(fileURLToPath(import.meta.url)), "../..");

export interface PrepareImportContextOptions { workspaceRoot: string; request: unknown; now?: Date }

function compactTimestamp(date: Date): string {
  return date.toISOString().replace(/[-:]/g, "").replace(/\.\d{3}Z$/, "Z");
}

function safeExcerpt(content: string, limit: number): string {
  const bounded = content.slice(0, limit).replace(/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----[\s\S]*?-----END (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/gi, "[redacted private key]").replace(/https?:\/\/[^\s/@:]+:[^\s/@]+@/g, "[redacted credential URL]");
  return bounded.replace(/\r\n/g, "\n").trimEnd();
}

async function validateRequest(value: unknown): Promise<ImportContextRequest> {
  const schema = JSON.parse(await readFile(join(projectRoot, ".agents/contracts/import-context-request.schema.json"), "utf8"));
  const validate = new Ajv2020({ allErrors: true, strict: false }).compile(schema);
  if (!validate(value)) {
    const detail = (validate.errors ?? []).map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ");
    throw new Error(`Invalid import request: ${detail}`);
  }
  return value as ImportContextRequest;
}

async function loadWorkspace(workspaceRoot: string): Promise<WorkspaceConfig> {
  let config: WorkspaceConfig;
  try {
    config = parseYaml(await readFile(join(workspaceRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  } catch (error) {
    throw new Error(`Missing or unreadable required discovery input workspace.yaml: ${(error as Error).message}`);
  }
  const contractErrors = await validateContract("workspace", config);
  const semanticErrors = workspaceSemanticErrors(config);
  if (contractErrors.length || semanticErrors.length) throw new Error(`Invalid workspace configuration: ${[...contractErrors.map((error) => `${error.instancePath || "/"} ${error.message}`), ...semanticErrors].join("; ")}`);
  return config;
}

async function regularFile(path: string, required = false): Promise<boolean> {
  try {
    const info = await lstat(path);
    if (!info.isFile() || info.isSymbolicLink()) {
      if (required) throw new Error("not a regular file");
      return false;
    }
    await readFile(path);
    return true;
  } catch (error) {
    if (required) throw new Error(`Missing or unreadable required discovery input ${path}: ${(error as Error).message}`);
    if ((error as NodeJS.ErrnoException).code === "ENOENT") return false;
    throw new Error(`Unreadable discovery input ${path}: ${(error as Error).message}`);
  }
}

async function walkFiles(root: string, directory: string): Promise<string[]> {
  const result: string[] = [];
  let entries;
  try {
    entries = await readdir(directory, { withFileTypes: true });
  } catch (error) {
    throw new Error(`Unreadable discovery input ${relative(root, directory)}: ${(error as Error).message}`);
  }
  for (const entry of entries.sort((a, b) => a.name.localeCompare(b.name))) {
    if (entry.isSymbolicLink()) continue;
    const path = join(directory, entry.name);
    if (entry.isDirectory() && !excludedDiscoveryDirectories.has(entry.name)) result.push(...await walkFiles(root, path));
    else if (entry.isFile()) result.push(path);
  }
  return result;
}

async function optionalDirectory(root: string, name: string): Promise<string[]> {
  const path = join(root, name);
  try {
    const info = await lstat(path);
    if (info.isSymbolicLink()) return [];
    if (!info.isDirectory()) throw new Error("not a directory");
    return walkFiles(root, path);
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code === "ENOENT") return [];
    throw new Error(`Unreadable discovery input ${name}: ${(error as Error).message}`);
  }
}

function normalizeCitation(root: string, path: string): string {
  const citation = relative(root, path).split(sep).join("/");
  if (!citation || citation.startsWith("../") || citation === "..") throw new Error(`Discovery path escapes source repository: ${path}`);
  return citation;
}

async function discover(root: string, limits: Required<ImportContextLimits>): Promise<ImportContextEvidence[]> {
  const candidates: Array<{ path: string; kind: ImportContextEvidence["kind"]; trust: ImportContextEvidence["trust"]; synthetic?: string }> = [];
  for (const name of rootInstructions) if (await regularFile(join(root, name))) candidates.push({ path: join(root, name), kind: "root-instruction", trust: "standard" });
  if (!candidates.some((item) => item.kind === "root-instruction")) throw new Error(`Missing or unreadable required discovery inputs: ${rootInstructions.join(", ")}`);

  for (const directory of documentationRoots) for (const path of await optionalDirectory(root, directory)) {
    if (/\.(?:md|mdx|txt|rst)$/i.test(path)) candidates.push({ path, kind: "documentation", trust: "standard" });
  }
  for (const path of await walkFiles(root, root)) {
    const citation = normalizeCitation(root, path);
    if (citation.includes("/") && repositoryInstructions.has(citation.split("/").at(-1)!) && !candidates.some((item) => item.path === path)) {
      candidates.push({ path, kind: "repository-instruction", trust: "standard" });
    }
  }
  const top = (await readdir(root, { withFileTypes: true })).filter((entry) => !entry.isSymbolicLink()).sort((a, b) => a.name.localeCompare(b.name));
  for (const entry of top.filter((item) => item.isDirectory() && item.name !== ".git")) candidates.push({ path: join(root, entry.name), kind: "structural-signal", trust: "standard", synthetic: `[top-level directory: ${entry.name}/]` });
  for (const entry of top.filter((item) => item.isFile() && (structuralFiles.has(item.name) || item.name.startsWith(".github")))) {
    if (!candidates.some((candidate) => candidate.path === join(root, entry.name))) candidates.push({ path: join(root, entry.name), kind: "structural-signal", trust: "standard" });
  }
  for (const path of await optionalDirectory(root, "context")) if (/\.md$/i.test(path)) candidates.push({ path, kind: "repository-context", trust: "high" });

  const evidence: ImportContextEvidence[] = [];
  let total = 0;
  for (const candidate of candidates) {
    if (evidence.length >= limits.max_files) break;
    const citation = normalizeCitation(root, candidate.path);
    let bytes = 0;
    let excerpt = candidate.synthetic ?? "";
    if (!candidate.synthetic) {
      const content = await readFile(candidate.path);
      bytes = content.byteLength;
      const available = Math.min(limits.max_file_bytes, limits.max_total_bytes - total);
      if (available <= 0) break;
      excerpt = safeExcerpt(content.toString("utf8", 0, available), available);
      total += Math.min(bytes, available);
    }
    evidence.push({ path: citation, kind: candidate.kind, trust: candidate.trust, bytes, excerpt });
  }
  return evidence;
}

function contributionDocument(repository: string, commit: string, generatedAt: string, evidence: ImportContextEvidence[]): string {
  const rendered = evidence.map((item) => `### \`${item.path}\`\n\n- Kind: \`${item.kind}\`\n- Trust: \`${item.trust}\`\n- Source commit: \`${commit}\`\n\n\`\`\`text\n${item.excerpt}\n\`\`\``).join("\n\n");
  return `# Import context evidence: ${repository}\n\n- Run: \`import-context-${generatedAt}\`\n\n## Outcome\n\nRead-only discovery captured bounded, source-cited repository evidence for human curation.\n\n## Affected repositories\n\n- \`${repository}\` was read at commit \`${commit}\`; the source repository was not modified.\n\n## Pull requests and commits\n\n- No source-repository commit or pull request was created.\n\n## Verification\n\n- Workspace and source Git state were clean before discovery.\n- Every evidence item below cites an exact repository-relative path.\n\n## Decisions and deviations\n\n- Evidence is recorded without inferring unsupported product facts.\n\n## Remaining risks and follow-up\n\n- Human curation must preserve citations and record unknowns explicitly.\n\n## Candidate durable learnings\n\n${rendered}\n`;
}

export async function prepareImportContext(options: PrepareImportContextOptions): Promise<ImportContextResult> {
  const workspaceRoot = resolve(options.workspaceRoot);
  const request = await validateRequest(options.request);
  const config = await loadWorkspace(workspaceRoot);
  const registration = config.repositories[request.repository];
  if (!registration) throw new Error(`Unregistered repository: ${request.repository}`);
  const sourceRoot = resolve(workspaceRoot, registration.path);
  if (await realpath(await git(workspaceRoot, ["rev-parse", "--show-toplevel"])) !== await realpath(workspaceRoot)) throw new Error("Workspace root must be the wrapper Git root");
  await assertCleanRepository(workspaceRoot).catch((error) => { throw new Error(`Dirty workspace Git state: ${(error as Error).message}`); });
  let sourceReal: string;
  try { sourceReal = await realpath(sourceRoot); } catch (error) { throw new Error(`Missing or unreadable required discovery input repository ${request.repository}: ${(error as Error).message}`); }
  assertInside(await realpath(workspaceRoot), sourceReal);
  if (sourceReal === await realpath(workspaceRoot)) throw new Error("Source repository must be separate from the wrapper repository");
  if (await realpath(await git(sourceReal, ["rev-parse", "--show-toplevel"])) !== sourceReal) throw new Error(`Registered repository path is not its Git root: ${request.repository}`);
  await assertCleanRepository(sourceReal).catch((error) => { throw new Error(`Dirty source Git state: ${(error as Error).message}`); });
  const sourceCommit = await git(sourceReal, ["rev-parse", "HEAD"]);
  const limits = { ...defaults, ...(request.limits ?? {}) };
  const evidence = await discover(sourceReal, limits);
  const now = options.now ?? new Date();
  const timestamp = compactTimestamp(now);
  const contributionRelative = `contributions/import-context/${request.repository}/${timestamp}-import-context-${request.repository}.md`;
  const contributionPath = assertInside(workspaceRoot, join(workspaceRoot, contributionRelative));
  await mkdir(dirname(contributionPath), { recursive: true });
  assertInside(await realpath(workspaceRoot), await realpath(dirname(contributionPath)));
  await writeTextExclusive(contributionPath, contributionDocument(request.repository, sourceCommit, timestamp, evidence));
  const runtimeRoot = assertInside(workspaceRoot, join(workspaceRoot, ".runtime", "import-context", `${timestamp}-${request.repository}`));
  await ensurePrivateDirectory(runtimeRoot);
  const manifestPath = join(runtimeRoot, "manifest.json");
  const manifest: ImportContextManifest = { contract_version: 1, repository: request.repository, source_root: sourceReal, source_commit: sourceCommit, evidence, limits, contribution: contributionRelative, generated_at: now.toISOString() };
  await writeJsonAtomic(manifestPath, manifest);
  return { contribution: contributionRelative, manifest: manifestPath, evidence };
}
