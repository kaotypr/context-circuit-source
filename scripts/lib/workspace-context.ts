import type { WorkspaceBootstrapContext } from "./types.js";

export const neutralWorkspaceContext: WorkspaceBootstrapContext = {
  project_summary: "This workspace has not been initialized. Run the initialize-workspace skill to record the product purpose and authoritative sources.",
  architecture: [],
  conventions: [],
  decisions: [],
};

function listDocument(title: string, values: string[], empty: string): string {
  const body = values.length > 0 ? values.map((value) => `- ${value.trim()}`).join("\n") : empty;
  return `# ${title}\n\n${body}\n`;
}

export function renderWorkspaceContext(context: WorkspaceBootstrapContext): Record<string, string> {
  return {
    "context/PROJECT.md": `# Project\n\n${context.project_summary.trim()}\n`,
    "context/ARCHITECTURE.md": listDocument("Architecture", context.architecture, "No project architecture has been recorded yet."),
    "context/CONVENTIONS.md": listDocument("Conventions", context.conventions, "No project-specific conventions have been recorded yet."),
    "context/DECISIONS.md": listDocument("Decisions", context.decisions, "No project decisions have been recorded yet."),
  };
}
