import type { WorkspaceConfig } from "./types.js";

const managedStart = "<!-- context-circuit:workspace:start -->";
const managedEnd = "<!-- context-circuit:workspace:end -->";

const canonicalFrameworkTitle = "# Context Circuit\n";

function repositoryRows(config: WorkspaceConfig): string {
  const repositories = Object.entries(config.repositories).sort(([left], [right]) => left.localeCompare(right));
  if (repositories.length === 0) return "No product repositories are registered yet.";
  return [
    "| Repository | Role | Path | Base branch | Remote |",
    "| --- | --- | --- | --- | --- |",
    ...repositories.map(([name, repository]) => `| ${name} | ${repository.role} | \`${repository.path}\` | \`${repository.default_branch}\` | ${repository.remote ?? "Not recorded"} |`),
  ].join("\n");
}

function sourceLinks(config: WorkspaceConfig): string {
  const sources = config.context?.authoritative_sources ?? [];
  if (sources.length === 0) return "- [Context source register](context/SOURCES.md) — no authoritative sources recorded; unknowns remain explicit.";
  return [
    "- [Context source register](context/SOURCES.md)",
    ...sources.map((source) => `- ${source.kind}: ${source.reference} — ${source.purpose}`),
  ].join("\n");
}

export function renderManagedWorkspaceReadme(config: WorkspaceConfig): string {
  const purpose = config.workspace.purpose ?? "Project purpose has not been recorded yet.";
  return `${managedStart}
# ${config.workspace.name}

${purpose}

## Product repositories

${repositoryRows(config)}

## Common actions

- Configure this wrapper: \`$configure-workspace\` (Codex) or \`/configure-workspace\` (Claude Code).
- Choose reviewed work: \`$whats-next\`.
- Run explicitly selected work: \`$run-task\`.
- Create an optional reviewed plan: \`$create-plan\`.
- Curate completed-work learning: \`$sync-context\`.

## Project context

- [Project summary](context/PROJECT.md)
- [Architecture](context/ARCHITECTURE.md)
- [Conventions](context/CONVENTIONS.md)
- [Decisions](context/DECISIONS.md)
${sourceLinks(config)}
- [Approved and draft plans](context/plans/)

Workspace mode: **${config.workspace.mode}**. Review mode: **${config.workflow.review_mode ?? (config.workflow.wrapper_change_policy === "pull-request" ? "remote" : "local")}**.
${managedEnd}`;
}

export function reconcileWorkspaceReadme(current: string, config: WorkspaceConfig): string {
  const managed = renderManagedWorkspaceReadme(config);
  const start = current.indexOf(managedStart);
  const end = current.indexOf(managedEnd);
  if ((start === -1) !== (end === -1) || (start !== -1 && end < start)) throw new Error("Malformed managed workspace block in README.md");
  if (start !== -1) {
    const after = end + managedEnd.length;
    return `${current.slice(0, start)}${managed}${current.slice(after)}`.replace(/\s*$/, "\n");
  }
  const existing = current.trim();
  if (!existing) return `${managed}\n`;
  if (existing.startsWith(canonicalFrameworkTitle)) {
    return `${managed}\n\n## Context Circuit framework\n\n${existing.slice(canonicalFrameworkTitle.length).trimStart()}\n`;
  }
  return `${managed}\n\n${existing}\n`;
}

