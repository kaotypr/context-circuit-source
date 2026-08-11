import { lstat, readdir, readFile } from "node:fs/promises";
import { join, relative, resolve } from "node:path";
import type { ErrorObject } from "ajv";
import { validatePlanDirectory } from "./plans.js";
import type {
  FakeActivitySource,
  NextAction,
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

function contractMessages(errors: ErrorObject[]): string[] {
  return errors.map((error) => `${error.instancePath || "/"} ${error.message}`);
}

async function isDirectory(path: string): Promise<boolean> {
  try {
    return (await lstat(path)).isDirectory();
  } catch {
    return false;
  }
}

function acceptanceIsSufficient(raw: string): boolean {
  return /^- (?!None recorded\.$).+/m.test(raw);
}

function rank(candidate: WorkCandidate): number {
  if (candidate.urgent) return 1;
  if (candidate.state === "in-progress") return 2;
  if (["review", "verification-failure", "ci-failure"].includes(candidate.kind) || candidate.state === "failed") return 3;
  if (candidate.kind !== "plan-work-item") return 4;
  return 5;
}

function blockers(candidate: WorkCandidate, currentUser: string): string[] {
  const values: string[] = [];
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

function actionFor(candidate: WorkCandidate, candidateBlockers: string[], config: WorkspaceConfig, action: "execute" | "enable"): NextAction {
  const evidence = [
    `state: ${candidate.state}`,
    `plan approval: ${candidate.plan_approval_state}`,
    candidate.dependencies.length === 0 ? "dependencies: none" : `dependencies: ${candidate.dependencies.map((item) => `${item.reference}=${item.state}`).join(", ")}`,
    `scope sufficient: ${candidate.scope_sufficient}`,
    `acceptance sufficient: ${candidate.acceptance_sufficient}`,
    `repository access available: ${candidate.access_available}`,
  ];
  const sequence = candidate.repositories.flatMap((repository) => {
    const agent = config.repositories[repository]?.agent ?? repository;
    return [`repository worker (${agent})`, "independent verifier"];
  });
  const title = action === "execute" ? candidate.title : enablingTitle(candidate, candidateBlockers);
  return {
    action,
    candidate_id: candidate.candidate_id,
    title,
    why: action === "execute"
      ? `${rankingReason(candidate)}; all readiness checks passed.`
      : `No candidate is currently executable; this is the smallest visible action that addresses the first blocker for ${candidate.title}.`,
    readiness_evidence: evidence,
    source_references: [...new Set([candidate.source_reference, ...(candidate.plan_reference ? [candidate.plan_reference] : [])])],
    repositories: candidate.repositories,
    agent_sequence: action === "execute" ? [...new Set(sequence)] : [],
    blockers: candidateBlockers,
    risks: candidate.risks,
  };
}

function rankingReason(candidate: WorkCandidate): string {
  if (candidate.urgent) return "It is explicitly urgent";
  if (candidate.state === "in-progress") return "It is actionable work already in progress";
  if (["review", "verification-failure", "ci-failure"].includes(candidate.kind) || candidate.state === "failed") {
    return "It addresses review, verification, or CI feedback";
  }
  if (candidate.kind !== "plan-work-item") return "It is the highest-priority approved ready source candidate";
  return "It is the next dependency-ready item in an approved plan";
}

function enablingTitle(candidate: WorkCandidate, candidateBlockers: string[]): string {
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

async function discoverPlanCandidates(
  workspaceRoot: string,
  config: WorkspaceConfig,
  facts: Map<string, WorkCandidate>,
): Promise<{ candidates: WorkCandidate[]; warnings: string[]; matchedFacts: Set<string> }> {
  const plansRoot = join(workspaceRoot, "context", "plans");
  if (!await isDirectory(plansRoot)) return { candidates: [], warnings: [], matchedFacts: new Set() };
  const entries = (await readdir(plansRoot, { withFileTypes: true }))
    .filter((entry) => entry.isDirectory() && !entry.isSymbolicLink())
    .sort((left, right) => left.name.localeCompare(right.name));
  const candidates: WorkCandidate[] = [];
  const warnings: string[] = [];
  const matchedFacts = new Set<string>();
  for (const entry of entries) {
    const planDirectory = join(plansRoot, entry.name);
    const validation = await validatePlanDirectory(planDirectory);
    if (!validation.index || !validation.work_breakdown || validation.errors.length > 0) {
      warnings.push(`Skipped invalid plan ${entry.name}: ${validation.errors.join("; ") || "missing parsed plan material"}`);
      continue;
    }
    const requirementPath = join(planDirectory, "0010-requirements.md");
    const acceptanceSufficient = acceptanceIsSufficient(await readFile(requirementPath, "utf8"));
    for (const item of validation.work_breakdown.items) {
      const fact = facts.get(item.work_id);
      if (fact) matchedFacts.add(fact.candidate_id);
      const repositories = fact?.repositories.length ? fact.repositories : (config.repositories[item.area] ? [item.area] : []);
      const dependencyFacts = new Map([...facts.values()].filter((value) => value.work_id).map((value) => [value.work_id!, value.state]));
      const dependencies = item.depends_on.map((reference) => ({
        reference,
        state: dependencyFacts.get(reference) === "completed"
          ? "completed" as const
          : dependencyFacts.has(reference) ? "pending" as const : "unknown" as const,
      }));
      const planReference = relative(workspaceRoot, join(planDirectory, "README.md"));
      candidates.push({
        contract_version: 1,
        candidate_id: `plan:${validation.index.plan_id}:${item.work_id}`,
        kind: "plan-work-item",
        work_id: item.work_id,
        title: item.title,
        state: fact?.state ?? "ready",
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
        risks: fact?.risks ?? [],
      });
    }
  }
  return { candidates, warnings, matchedFacts };
}

export async function recommendWhatsNext(
  workspaceRootInput: string,
  activity: FakeActivitySource | null = null,
  now = new Date(),
): Promise<WhatsNextResult> {
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
  const facts = new Map<string, WorkCandidate>();
  for (const candidate of activity?.candidates ?? []) {
    if (candidate.work_id) facts.set(candidate.work_id, candidate);
  }
  const discovered = await discoverPlanCandidates(workspaceRoot, config, facts);
  const external = (activity?.candidates ?? []).filter((candidate) => !discovered.matchedFacts.has(candidate.candidate_id));
  const hydratedExternal: WorkCandidate[] = [];
  for (const candidate of external) {
    hydratedExternal.push({
      ...candidate,
      access_available: candidate.access_available && await repositoryAccess(workspaceRoot, config, candidate.repositories),
    });
  }
  const candidates = [...discovered.candidates, ...hydratedExternal];
  const duplicateIds = candidates.filter((candidate, index) => candidates.findIndex((value) => value.candidate_id === candidate.candidate_id) !== index);
  if (duplicateIds.length > 0) throw new Error(`Duplicate candidate ID: ${duplicateIds[0]!.candidate_id}`);
  const excluded = candidates.filter((candidate) => candidate.state === "completed" || candidate.state === "cancelled");
  const assessed: AssessedCandidate[] = candidates
    .filter((candidate) => candidate.state !== "completed" && candidate.state !== "cancelled")
    .map((candidate) => ({ candidate, blockers: blockers(candidate, currentUser), rank: rank(candidate) }));
  const ordered = [...assessed].sort((left, right) =>
    left.rank - right.rank
    || right.candidate.priority - left.candidate.priority
    || left.candidate.candidate_id.localeCompare(right.candidate.candidate_id));
  const executable = ordered.filter((item) => item.blockers.length === 0);
  const blocked = ordered.filter((item) => item.blockers.length > 0);
  let recommendation: NextAction;
  let alternatives: NextAction[];
  if (executable.length > 0) {
    recommendation = actionFor(executable[0]!.candidate, [], config, "execute");
    alternatives = executable.slice(1, 3).map((item) => actionFor(item.candidate, [], config, "execute"));
  } else if (blocked.length > 0) {
    recommendation = actionFor(blocked[0]!.candidate, blocked[0]!.blockers, config, "enable");
    alternatives = blocked.slice(1, 3).map((item) => actionFor(item.candidate, item.blockers, config, "enable"));
  } else {
    recommendation = {
      action: "enable",
      candidate_id: null,
      title: "Create or approve a scoped work source",
      why: "No executable or blocked candidate was found in the configured read-only sources.",
      readiness_evidence: ["approved plan candidates: none", "activity candidates: none"],
      source_references: ["context/plans", "workspace.yaml#activity"],
      repositories: [],
      agent_sequence: [],
      blockers: ["no available candidate provides sufficient scope and acceptance criteria"],
      risks: [],
    };
    alternatives = [];
  }
  const result: WhatsNextResult = {
    contract_version: 1,
    generated_at: now.toISOString(),
    recommendation,
    alternatives,
    considered: { total: candidates.length, executable: executable.length, blocked: blocked.length, excluded: excluded.length },
    warnings: discovered.warnings,
    no_state_changed: true,
  };
  const resultErrors = contractMessages(await validateContract("whats-next-result", result));
  if (resultErrors.length > 0) throw new Error(`Generated invalid whats-next result:\n- ${resultErrors.join("\n- ")}`);
  return result;
}
