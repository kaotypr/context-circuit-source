import { lstat, readdir, readFile } from "node:fs/promises";
import { join, relative, resolve, sep } from "node:path";
import type { ErrorObject } from "ajv";
import { contributionDocumentErrors } from "./finish-work.js";
import { git } from "./git.js";
import { validatePlanDirectory } from "./plans.js";
import type {
  CandidateState,
  CloseoutRecord,
  FakeActivitySource,
  NextAction,
  RuntimeManifest,
  TaskBrief,
  WhatsNextResult,
  WorkCandidate,
  WorkspaceConfig,
} from "./types.js";
import { readData, validateContract, workspaceSemanticErrors } from "./validation.js";

interface AssessedCandidate {
  candidate: WorkCandidate;
  blockers: string[];
  rank: number;
}

interface StateObservation {
  state: CandidateState;
  source_reference: string;
  precedence: number;
  plan_reference?: string;
}

interface ProjectedState {
  state: CandidateState;
  observations: StateObservation[];
  contradiction: boolean;
}

function contractMessages(errors: ErrorObject[]): string[] {
  return errors.map((error) => `${error.instancePath || "/"} ${error.message}`);
}

async function isDirectory(path: string): Promise<boolean> {
  try {
    const info = await lstat(path);
    return info.isDirectory() && !info.isSymbolicLink();
  } catch {
    return false;
  }
}

function reference(workspaceRoot: string, path: string): string {
  const raw = relative(workspaceRoot, path);
  return raw && !raw.startsWith(`..${sep}`) && raw !== ".." ? raw.replaceAll("\\", "/") : path;
}

function inside(root: string, path: string): boolean {
  const resolvedRoot = resolve(root);
  const resolvedPath = resolve(path);
  return resolvedPath === resolvedRoot || resolvedPath.startsWith(`${resolvedRoot}${sep}`);
}

function acceptanceIsSufficient(raw: string): boolean {
  return /^- (?!None recorded\.$).+/m.test(raw);
}

function rank(candidate: WorkCandidate): number {
  if (candidate.kind === "reconciliation") return 0;
  if (candidate.urgent) return 1;
  if (candidate.state === "in-progress" || candidate.state === "closeout") return 2;
  if (["review", "verification-failure", "ci-failure"].includes(candidate.kind) || candidate.state === "failed" || candidate.state === "review") return 3;
  if (candidate.kind !== "plan-work-item") return 4;
  return 5;
}

function blockers(candidate: WorkCandidate, currentUser: string): string[] {
  const values: string[] = [];
  if (candidate.kind === "reconciliation") values.push("read-only sources report contradictory work states");
  if (candidate.owner && candidate.owner !== currentUser) values.push(`owned by another active contributor: ${candidate.owner}`);
  for (const dependency of candidate.dependencies) {
    if (dependency.state !== "completed") values.push(`dependency ${dependency.reference} is ${dependency.state}`);
  }
  if (candidate.plan_approval_state !== "approved" && candidate.plan_approval_state !== "not-applicable") {
    values.push(`governing plan is ${candidate.plan_approval_state}`);
  }
  if (!candidate.scope_sufficient) values.push("scope is insufficient");
  if (!candidate.acceptance_sufficient) values.push("acceptance criteria are insufficient");
  if (candidate.repositories.length === 0) values.push("no affected repository is resolved");
  if (!candidate.access_available) values.push("required repository access is unavailable");
  if (candidate.contract_blocked) values.push("an unresolved contract decision blocks implementation");
  return [...new Set(values)];
}

function actionKind(candidate: WorkCandidate, candidateBlockers: string[]): NextAction["action"] {
  if (candidate.kind === "reconciliation") return "reconcile";
  if (candidateBlockers.length > 0) return "enable";
  if (candidate.state === "review") return "review";
  if (candidate.state === "closeout") return "closeout";
  return "execute";
}

function actionFor(candidate: WorkCandidate, candidateBlockers: string[], config: WorkspaceConfig): NextAction {
  const action = actionKind(candidate, candidateBlockers);
  const stateSources = candidate.state_sources?.length
    ? candidate.state_sources.map((item) => `state ${item.state}: ${item.source_reference}`)
    : [`state ${candidate.state}: ${candidate.source_reference}`];
  const evidence = [
    ...stateSources,
    `plan approval: ${candidate.plan_approval_state}`,
    `plan approval source: ${candidate.plan_reference ?? candidate.source_reference}`,
    candidate.dependencies.length === 0 ? "dependencies: none" : `dependencies: ${candidate.dependencies.map((item) => `${item.reference}=${item.state}`).join(", ")}`,
    `dependency source: ${candidate.plan_reference ?? candidate.source_reference}`,
    `scope sufficient: ${candidate.scope_sufficient}`,
    `scope source: ${candidate.source_reference}`,
    `acceptance sufficient: ${candidate.acceptance_sufficient}`,
    `acceptance source: ${candidate.source_reference}`,
    `repository access available: ${candidate.access_available}`,
    "repository access source: workspace.yaml#repositories",
  ];
  const sequence = candidate.repositories.flatMap((repository) => {
    const agent = config.repositories[repository]?.agent ?? repository;
    return [`repository worker (${agent})`, "independent verifier"];
  });
  const title = action === "execute" ? candidate.title : actionTitle(candidate, candidateBlockers, action);
  return {
    action,
    candidate_id: candidate.candidate_id,
    title,
    why: action === "execute"
      ? `${rankingReason(candidate)}; all readiness checks passed.`
      : action === "review"
        ? "Verified implementation evidence is ready for human review or merge preparation; implementation must not be duplicated."
        : action === "closeout"
          ? "Implementation has advanced beyond execution and the remaining work is human-gated closeout or cleanup."
          : action === "reconcile"
            ? "Configured activity and local outcome evidence disagree; reconcile the cited sources without mutating them or starting duplicate implementation."
            : `No candidate is currently executable; this is the smallest visible action that addresses the first blocker for ${candidate.title}.`,
    readiness_evidence: evidence,
    source_references: [...new Set([candidate.source_reference, ...(candidate.plan_reference ? [candidate.plan_reference] : []), ...(candidate.state_sources ?? []).map((item) => item.source_reference)])],
    repositories: candidate.repositories,
    agent_sequence: action === "execute" ? [...new Set(sequence)] : [],
    blockers: candidateBlockers,
    risks: candidate.risks,
  };
}

function rankingReason(candidate: WorkCandidate): string {
  if (candidate.urgent) return "It is explicitly urgent";
  if (candidate.state === "in-progress") return "It is actionable work already in progress";
  if (["review", "verification-failure", "ci-failure"].includes(candidate.kind) || candidate.state === "failed") return "It addresses review, verification, or CI feedback";
  if (candidate.kind !== "plan-work-item") return "It is the highest-priority approved ready source candidate";
  return "It is the next dependency-ready item in an approved plan";
}

function actionTitle(candidate: WorkCandidate, candidateBlockers: string[], action: NextAction["action"]): string {
  if (action === "review") return `Review or prepare merge for ${candidate.title}`;
  if (action === "closeout") return `Complete closeout or cleanup for ${candidate.title}`;
  if (action === "reconcile") return `Reconcile contradictory state for ${candidate.title}`;
  const first = candidateBlockers[0] ?? "readiness is not established";
  if (first.startsWith("governing plan is draft")) return `Approve the governing plan for ${candidate.title}`;
  if (first.startsWith("dependency ")) return `Resolve or confirm ${first.replace(" is ", " as ")}`;
  if (first.startsWith("owned by another")) return `Confirm ownership before starting ${candidate.title}`;
  if (first.includes("contract")) return `Resolve the blocking contract for ${candidate.title}`;
  if (first.includes("repository")) return `Register or restore repository access for ${candidate.title}`;
  return `Clarify ${first} for ${candidate.title}`;
}

async function repositoryAccess(workspaceRoot: string, config: WorkspaceConfig, repositories: string[]): Promise<boolean> {
  if (repositories.length === 0) return false;
  for (const name of repositories) {
    const repository = config.repositories[name];
    if (!repository || !await isDirectory(resolve(workspaceRoot, repository.path))) return false;
  }
  return true;
}

function runtimeState(manifest: RuntimeManifest, workId: string): CandidateState {
  const item = manifest.plan_work_items?.find((candidate) => candidate.work_id === workId);
  if (manifest.status === "cancelled") return "cancelled";
  if (manifest.status === "closing") return "closeout";
  if (manifest.status === "closed") return "completed";
  if (manifest.status === "passed" || item?.outcome === "passed") return "review";
  if (manifest.status === "failed" || manifest.status === "blocked" || item?.outcome === "failed" || item?.outcome === "blocked") return "failed";
  return "in-progress";
}

async function validatedJson<T>(name: "runtime-manifest" | "task-brief" | "closeout-record", path: string): Promise<T> {
  const info = await lstat(path);
  if (!info.isFile() || info.isSymbolicLink()) throw new Error(`${name} is not a regular file`);
  const value = JSON.parse(await readFile(path, "utf8")) as T;
  const errors = contractMessages(await validateContract(name, value));
  if (errors.length > 0) throw new Error(errors.join("; "));
  return value;
}

async function discoverRuntimeObservations(workspaceRoot: string): Promise<{ observations: Map<string, StateObservation[]>; warnings: string[] }> {
  const observations = new Map<string, StateObservation[]>();
  const warnings: string[] = [];
  const runsRoot = join(workspaceRoot, ".runtime", "runs");
  if (!await isDirectory(runsRoot)) return { observations, warnings };
  const entries = (await readdir(runsRoot, { withFileTypes: true })).filter((entry) => entry.isDirectory() && !entry.isSymbolicLink()).sort((a, b) => a.name.localeCompare(b.name));
  for (const entry of entries) {
    const manifestPath = join(runsRoot, entry.name, "manifest.json");
    try {
      const info = await lstat(manifestPath);
      if (!info.isFile() || info.isSymbolicLink()) throw new Error("manifest is not a regular file");
      const manifest = await validatedJson<RuntimeManifest>("runtime-manifest", manifestPath);
      if (manifest.source_kind !== "plan" || !manifest.plan_work_items) continue;
      if (!inside(workspaceRoot, manifest.task_brief)) throw new Error("task brief escapes the workspace");
      const brief = await validatedJson<TaskBrief>("task-brief", manifest.task_brief);
      if (brief.source.kind !== "plan" || brief.plan.approval_state !== "approved" || brief.work_id !== manifest.work_id || !brief.plan.work_ids.includes(manifest.work_id)) throw new Error("plan task brief identity does not match manifest");
      for (const item of manifest.plan_work_items) {
        let state = runtimeState(manifest, item.work_id);
        const sources = [reference(workspaceRoot, manifestPath)];
        const repository = manifest.repositories.find((candidate) => candidate.name === item.repository);
        if (repository?.closeout_record) {
          if (!inside(workspaceRoot, repository.closeout_record)) throw new Error("closeout record escapes the workspace");
          const closeout = await validatedJson<CloseoutRecord>("closeout-record", repository.closeout_record);
          if (closeout.run_id !== manifest.run_id || closeout.work_id !== item.work_id || closeout.repository !== item.repository) throw new Error("closeout identity does not match manifest");
          state = closeout.status === "closed" ? (closeout.outcome === "merged" ? "completed" : "cancelled") : "closeout";
          sources.push(reference(workspaceRoot, repository.closeout_record));
        }
        const values = observations.get(item.work_id) ?? [];
        for (const source_reference of sources) values.push({ state, source_reference, precedence: state === "completed" || state === "cancelled" ? 30 : 20, plan_reference: brief.source.reference.replace(/\/README\.md$/, "").replace(/\/$/, "") });
        observations.set(item.work_id, values);
      }
    } catch (error) {
      const code = (error as NodeJS.ErrnoException).code;
      if (code !== "ENOENT") warnings.push(`Ignored malformed runtime evidence ${entry.name}: ${(error as Error).message}`);
    }
  }
  return { observations, warnings };
}

async function discoverDurableContributions(workspaceRoot: string): Promise<{ observations: Map<string, StateObservation[]>; warnings: string[] }> {
  const observations = new Map<string, StateObservation[]>();
  const warnings: string[] = [];
  const root = join(workspaceRoot, "contributions");
  if (!await isDirectory(root)) return { observations, warnings };
  const groups = (await readdir(root, { withFileTypes: true })).filter((entry) => entry.isDirectory() && !entry.isSymbolicLink()).sort((a, b) => a.name.localeCompare(b.name));
  for (const group of groups) {
    const directory = join(root, group.name);
    const files = (await readdir(directory, { withFileTypes: true })).filter((entry) => entry.isFile() && !entry.isSymbolicLink() && entry.name.endsWith(".md")).sort((a, b) => a.name.localeCompare(b.name));
    for (const file of files) {
      const path = join(directory, file.name);
      const source_reference = reference(workspaceRoot, path);
      try {
        await git(workspaceRoot, ["ls-files", "--error-unmatch", "--", source_reference]);
      } catch {
        continue;
      }
      const content = await readFile(path, "utf8");
      const work = content.match(/^# ([A-Z][A-Z0-9]{1,15}-\d{3,}):/m)?.[1];
      const run = content.match(/^- Run: `([^`]+)`$/m)?.[1];
      const merged = content.includes("Merged after human review.");
      const abandoned = content.includes("Deliberately abandoned by the human.");
      const documentErrors = contributionDocumentErrors(path, content, run);
      if (!work || !run || merged === abandoned || documentErrors.length > 0) {
        warnings.push(`Ignored unrecognized durable contribution ${source_reference}`);
        continue;
      }
      const values = observations.get(work) ?? [];
      values.push({ state: merged ? "completed" : "cancelled", source_reference, precedence: 40 });
      observations.set(work, values);
    }
  }
  return { observations, warnings };
}

function project(observations: StateObservation[]): ProjectedState | null {
  if (observations.length === 0) return null;
  const ordered = observations.slice().sort((a, b) => b.precedence - a.precedence || a.source_reference.localeCompare(b.source_reference));
  return { state: ordered[0]!.state, observations: ordered, contradiction: new Set(ordered.map((item) => item.state)).size > 1 };
}

async function discoverPlanCandidates(
  workspaceRoot: string,
  config: WorkspaceConfig,
  activityFacts: Map<string, WorkCandidate[]>,
  localObservations: Map<string, StateObservation[]>,
): Promise<{ candidates: WorkCandidate[]; warnings: string[]; matchedFacts: Set<string> }> {
  const plansRoot = join(workspaceRoot, "context", "plans");
  if (!await isDirectory(plansRoot)) return { candidates: [], warnings: [], matchedFacts: new Set() };
  const entries = (await readdir(plansRoot, { withFileTypes: true })).filter((entry) => entry.isDirectory() && !entry.isSymbolicLink()).sort((a, b) => a.name.localeCompare(b.name));
  const candidates: WorkCandidate[] = [];
  const warnings: string[] = [];
  const matchedFacts = new Set<string>();
  const validations = new Map<string, Awaited<ReturnType<typeof validatePlanDirectory>>>();
  const workIdCounts = new Map<string, number>();
  for (const entry of entries) {
    const planDirectory = join(plansRoot, entry.name);
    const validation = await validatePlanDirectory(planDirectory);
    if (!validation.index || !validation.work_breakdown || validation.errors.length > 0) {
      warnings.push(`Skipped invalid plan ${entry.name}: ${validation.errors.join("; ") || "missing parsed plan material"}`);
      continue;
    }
    validations.set(entry.name, validation);
    for (const item of validation.work_breakdown.items) workIdCounts.set(item.work_id, (workIdCounts.get(item.work_id) ?? 0) + 1);
  }
  for (const entry of entries) {
    const planDirectory = join(plansRoot, entry.name);
    const validation = validations.get(entry.name);
    if (!validation?.index || !validation.work_breakdown) continue;
    const planDirectoryReference = relative(workspaceRoot, planDirectory).replaceAll("\\", "/");
    const projectedByWork = new Map<string, ProjectedState>();
    for (const item of validation.work_breakdown.items) {
      const relevantLocal = (localObservations.get(item.work_id) ?? []).filter((observation) => {
        if (observation.plan_reference) return observation.plan_reference === planDirectoryReference;
        if (observation.precedence === 40 && workIdCounts.get(item.work_id)! > 1) return false;
        return true;
      });
      const activity = activityFacts.get(item.work_id) ?? [];
      const value = project([...relevantLocal, ...activity.map((fact) => ({ state: fact.state, source_reference: fact.source_reference, precedence: 10 }))]);
      if (value) projectedByWork.set(item.work_id, value);
      if (workIdCounts.get(item.work_id)! > 1 && (localObservations.get(item.work_id) ?? []).some((observation) => !observation.plan_reference && observation.precedence === 40)) {
        warnings.push(`Ignored ambiguous durable contribution for ${item.work_id}: multiple plans use that work ID`);
      }
    }
    const requirementPath = join(planDirectory, "0010-requirements.md");
    const acceptanceSufficient = acceptanceIsSufficient(await readFile(requirementPath, "utf8"));
    for (const item of validation.work_breakdown.items) {
      const facts = activityFacts.get(item.work_id) ?? [];
      for (const fact of facts) matchedFacts.add(fact.candidate_id);
      const fact = facts[0];
      const projection = projectedByWork.get(item.work_id);
      const repositories = fact?.repositories.length ? fact.repositories : (config.repositories[item.repository] ? [item.repository] : []);
      const dependencies = item.depends_on.map((dependency) => ({
        reference: dependency,
        state: projectedByWork.get(dependency)?.state === "completed" ? "completed" as const : projectedByWork.has(dependency) ? "pending" as const : "unknown" as const,
      }));
      const planReference = relative(workspaceRoot, join(planDirectory, "README.md"));
      const state = projection?.state ?? "ready";
      const contradiction = Boolean(projection?.contradiction) || (Boolean(projection) && validation.index.status !== "approved");
      const stateSources = projection?.observations.map(({ state: observed, source_reference }) => ({ state: observed, source_reference })) ?? [{ state: "ready" as const, source_reference: `${relative(workspaceRoot, join(planDirectory, validation.index.work_breakdown))}#${item.work_id}` }];
      candidates.push({
        contract_version: 1,
        candidate_id: `plan:${validation.index.plan_id}:${item.work_id}`,
        kind: contradiction ? "reconciliation" : state === "review" ? "review" : "plan-work-item",
        work_id: item.work_id,
        title: item.title,
        state,
        urgent: fact?.urgent ?? false,
        priority: fact?.priority ?? 0,
        owner: fact?.owner ?? null,
        plan_reference: planReference,
        plan_approval_state: validation.index.status,
        dependencies,
        scope_sufficient: fact?.scope_sufficient ?? (item.title.trim().length > 0 && repositories.length > 0),
        acceptance_sufficient: fact?.acceptance_sufficient ?? acceptanceSufficient,
        repositories,
        access_available: (fact?.access_available ?? true) && await repositoryAccess(workspaceRoot, config, repositories),
        contract_blocked: fact?.contract_blocked ?? false,
        source_reference: `${relative(workspaceRoot, join(planDirectory, validation.index.work_breakdown))}#${item.work_id}`,
        state_sources: stateSources,
        risks: [...new Set([...(fact?.risks ?? []), ...(contradiction ? ["Starting implementation before reconciliation could duplicate or overwrite completed work."] : [])])],
      });
    }
  }
  return { candidates, warnings, matchedFacts };
}

export async function recommendWhatsNext(workspaceRootInput: string, activity: FakeActivitySource | null = null, now = new Date()): Promise<WhatsNextResult> {
  const workspaceRoot = resolve(workspaceRootInput);
  const config = await readData(join(workspaceRoot, "workspace.yaml")) as WorkspaceConfig;
  const workspaceErrors = contractMessages(await validateContract("workspace", config));
  workspaceErrors.push(...workspaceSemanticErrors(config));
  if (workspaceErrors.length > 0) throw new Error(`Invalid workspace configuration:\n- ${workspaceErrors.join("\n- ")}`);
  if (activity) {
    const activityErrors = contractMessages(await validateContract("fake-activity-source", activity));
    if (activityErrors.length > 0) throw new Error(`Invalid fake activity source:\n- ${activityErrors.join("\n- ")}`);
  }
  const currentUser = activity?.current_user ?? "local-user";
  const facts = new Map<string, WorkCandidate[]>();
  for (const candidate of activity?.candidates ?? []) if (candidate.work_id) facts.set(candidate.work_id, [...(facts.get(candidate.work_id) ?? []), candidate]);
  const runtime = await discoverRuntimeObservations(workspaceRoot);
  const durable = await discoverDurableContributions(workspaceRoot);
  const localObservations = new Map(runtime.observations);
  for (const [workId, observations] of durable.observations) localObservations.set(workId, [...(localObservations.get(workId) ?? []), ...observations]);
  const discovered = await discoverPlanCandidates(workspaceRoot, config, facts, localObservations);
  const external = (activity?.candidates ?? []).filter((candidate) => !discovered.matchedFacts.has(candidate.candidate_id));
  const hydratedExternal: WorkCandidate[] = [];
  for (const candidate of external) hydratedExternal.push({ ...candidate, state_sources: candidate.state_sources ?? [{ state: candidate.state, source_reference: candidate.source_reference }], access_available: candidate.access_available && await repositoryAccess(workspaceRoot, config, candidate.repositories) });
  const candidates = [...discovered.candidates, ...hydratedExternal];
  const duplicateIds = candidates.filter((candidate, index) => candidates.findIndex((value) => value.candidate_id === candidate.candidate_id) !== index);
  if (duplicateIds.length > 0) throw new Error(`Duplicate candidate ID: ${duplicateIds[0]!.candidate_id}`);
  const excluded = candidates.filter((candidate) => candidate.state === "completed" || candidate.state === "cancelled");
  const assessed: AssessedCandidate[] = candidates.filter((candidate) => candidate.state !== "completed" && candidate.state !== "cancelled").map((candidate) => ({ candidate, blockers: blockers(candidate, currentUser), rank: rank(candidate) }));
  const ordered = assessed.slice().sort((left, right) => left.rank - right.rank || right.candidate.priority - left.candidate.priority || left.candidate.candidate_id.localeCompare(right.candidate.candidate_id));
  const executable = ordered.filter((item) => item.blockers.length === 0);
  const blocked = ordered.filter((item) => item.blockers.length > 0);
  const reconciliation = blocked.filter((item) => item.candidate.kind === "reconciliation");
  let recommendation: NextAction;
  let alternatives: NextAction[];
  if (reconciliation.length > 0) {
    recommendation = actionFor(reconciliation[0]!.candidate, reconciliation[0]!.blockers, config);
    alternatives = [...reconciliation.slice(1), ...executable].slice(0, 2).map((item) => actionFor(item.candidate, item.blockers, config));
  } else if (executable.length > 0) {
    recommendation = actionFor(executable[0]!.candidate, [], config);
    alternatives = executable.slice(1, 3).map((item) => actionFor(item.candidate, [], config));
  } else if (blocked.length > 0) {
    recommendation = actionFor(blocked[0]!.candidate, blocked[0]!.blockers, config);
    alternatives = blocked.slice(1, 3).map((item) => actionFor(item.candidate, item.blockers, config));
  } else {
    recommendation = { action: "enable", candidate_id: null, title: "Create or approve a scoped work source", why: "No executable or blocked candidate was found in the configured read-only sources.", readiness_evidence: ["approved plan candidates none: context/plans", "activity candidates none: workspace.yaml#activity", "active runtime candidates none: .runtime/runs", "durable outcome candidates none: contributions"], source_references: ["context/plans", "workspace.yaml#activity", ".runtime/runs", "contributions"], repositories: [], agent_sequence: [], blockers: ["no available candidate provides sufficient scope and acceptance criteria"], risks: [] };
    alternatives = [];
  }
  const result: WhatsNextResult = { contract_version: 1, generated_at: now.toISOString(), recommendation, alternatives, considered: { total: candidates.length, executable: executable.length, blocked: blocked.length, excluded: excluded.length }, warnings: [...new Set([...runtime.warnings, ...durable.warnings, ...discovered.warnings])], no_state_changed: true };
  const resultErrors = contractMessages(await validateContract("whats-next-result", result));
  if (resultErrors.length > 0) throw new Error(`Generated invalid whats-next result:\n- ${resultErrors.join("\n- ")}`);
  return result;
}
