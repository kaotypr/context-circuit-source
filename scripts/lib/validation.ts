import { lstat, readFile } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { Ajv2020, type ErrorObject } from "ajv/dist/2020.js";
import { parse as parseYaml } from "yaml";
import type { WorkspaceConfig } from "./types.js";
import { contextReferenceError, remoteReferenceError } from "./safe-reference.js";

export const schemaNames = ["workspace", "workspace-bootstrap-request", "workspace-configure-request", "task-brief", "worker-result", "verifier-result", "plan-verifier-result", "runtime-manifest", "plan-runtime-revision", "run-task-request", "review-preparation", "review-publication-record", "merge-confirmation-record", "closeout-record", "context-sync-request", "context-sync-record", "plan-index", "plan-work-breakdown", "plan-task", "plan-connection", "plan-draft-request", "plan-generation-request", "work-candidate", "fake-activity-source", "whats-next-result", "activity-lifecycle-record", "plan-publication-discovery", "plan-publication-record", "product-knowledge-project", "product-knowledge-role", "product-knowledge-workflow", "product-knowledge-domain", "product-knowledge-candidate", "task-context-package", "product-knowledge-sync-record", "onboarding-pack"] as const;
export type SchemaName = (typeof schemaNames)[number];

export const requiredWorkspaceDocuments = [
  "README.md",
  "AGENTS.md",
  "CLAUDE.md",
  "WORKFLOW.md",
  "workspace.yaml",
  "context/PROJECT.md",
  "context/ARCHITECTURE.md",
  "context/CONVENTIONS.md",
  "context/DECISIONS.md",
  "context/SOURCES.md",
  "agents/coordinator.md",
  "agents/repository-worker.md",
  "agents/verifier.md",
  ".agents/bin/cc.mjs",
  ".agents/contracts/workspace.schema.json",
  ".agents/contracts/workspace-bootstrap-request.schema.json",
  ".agents/contracts/workspace-configure-request.schema.json",
  ".agents/contracts/review-preparation.schema.json",
  ".agents/contracts/review-publication-record.schema.json",
  ".agents/contracts/merge-confirmation-record.schema.json",
  ".agents/contracts/closeout-record.schema.json",
  ".agents/contracts/run-task-request.schema.json",
  ".agents/contracts/plan-verifier-result.schema.json",
  ".agents/contracts/plan-runtime-revision.schema.json",
  ".agents/contracts/context-sync-request.schema.json",
  ".agents/contracts/context-sync-record.schema.json",
  ".agents/contracts/plan-index.schema.json",
  ".agents/contracts/plan-work-breakdown.schema.json",
  ".agents/contracts/plan-draft-request.schema.json",
  ".agents/contracts/plan-generation-request.schema.json",
  ".agents/contracts/plan-task.schema.json",
  ".agents/contracts/plan-connection.schema.json",
  ".agents/contracts/work-candidate.schema.json",
  ".agents/contracts/fake-activity-source.schema.json",
  ".agents/contracts/whats-next-result.schema.json",
  ".agents/contracts/activity-lifecycle-record.schema.json",
  ".agents/contracts/plan-publication-discovery.schema.json",
  ".agents/contracts/plan-publication-record.schema.json",
  ".agents/contracts/product-knowledge-project.schema.json",
  ".agents/contracts/product-knowledge-role.schema.json",
  ".agents/contracts/product-knowledge-workflow.schema.json",
  ".agents/contracts/product-knowledge-domain.schema.json",
  ".agents/contracts/product-knowledge-candidate.schema.json",
  ".agents/contracts/task-context-package.schema.json",
  ".agents/contracts/product-knowledge-sync-record.schema.json",
  ".agents/contracts/onboarding-pack.schema.json",
  ".agents/skills/cc-initialize-workspace/SKILL.md",
  ".agents/skills/cc-configure-workspace/SKILL.md",
  ".agents/skills/cc-gather-context/SKILL.md",
  ".agents/skills/cc-run-task/SKILL.md",
  ".agents/skills/cc-finish-work/SKILL.md",
  ".agents/skills/cc-create-plan/SKILL.md",
  ".agents/skills/cc-whats-next/SKILL.md",
  ".agents/skills/cc-publish-plan-tasks/SKILL.md",
  ".agents/skills/cc-sync-context/SKILL.md",
  ".codex/skills/cc-initialize-workspace/SKILL.md",
  ".codex/skills/cc-configure-workspace/SKILL.md",
  ".codex/skills/cc-gather-context/SKILL.md",
  ".codex/skills/cc-run-task/SKILL.md",
  ".codex/skills/cc-finish-work/SKILL.md",
  ".codex/skills/cc-create-plan/SKILL.md",
  ".codex/skills/cc-whats-next/SKILL.md",
  ".codex/skills/cc-publish-plan-tasks/SKILL.md",
  ".codex/skills/cc-sync-context/SKILL.md",
  ".claude/commands/cc-initialize-workspace.md",
  ".claude/commands/cc-configure-workspace.md",
  ".claude/commands/cc-gather-context.md",
  ".claude/commands/cc-run-task.md",
  ".claude/commands/cc-finish-work.md",
  ".claude/commands/cc-create-plan.md",
  ".claude/commands/cc-whats-next.md",
  ".claude/commands/cc-publish-plan-tasks.md",
  ".claude/commands/cc-sync-context.md",
  "docs/getting-started.md",
  "docs/using-the-wrapper.md",
  "docs/configuration.md",
  "docs/command-reference.md",
] as const;

const projectRoot = resolve(dirname(fileURLToPath(import.meta.url)), "../..");

export async function readData(path: string): Promise<unknown> {
  const raw = await readFile(path, "utf8");
  return path.endsWith(".yaml") || path.endsWith(".yml") ? parseYaml(raw) : JSON.parse(raw);
}

export async function validateContract(name: SchemaName, value: unknown): Promise<ErrorObject[]> {
  const schema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", `${name}.schema.json`), "utf8"));
  const ajv = new Ajv2020({ allErrors: true, strict: false });
  ajv.addFormat("email", { type: "string", validate: (value: string) => /^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(value) });
  ajv.addFormat("date-time", {
    type: "string",
    validate: (value: string) => /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/.test(value) && !Number.isNaN(Date.parse(value)),
  });
  if (name === "fake-activity-source") {
    const candidateSchema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", "work-candidate.schema.json"), "utf8"));
    ajv.addSchema(candidateSchema);
  }
  if (name === "runtime-manifest") {
    const lifecycleSchema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", "activity-lifecycle-record.schema.json"), "utf8"));
    ajv.addSchema(lifecycleSchema);
  }
  if (name === "workspace-bootstrap-request" || name === "workspace-configure-request") {
    const workspaceSchema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", "workspace.schema.json"), "utf8"));
    ajv.addSchema(workspaceSchema);
    if (name === "workspace-configure-request") {
      const bootstrapSchema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", "workspace-bootstrap-request.schema.json"), "utf8"));
      ajv.addSchema(bootstrapSchema);
    }
  }
  const validate = ajv.compile(schema);
  return validate(value) ? [] : [...(validate.errors ?? [])];
}

export function workspaceSemanticErrors(config: WorkspaceConfig): string[] {
  const errors: string[] = [];
  const paths = new Map<string, string>();
  const remotes: Array<[string, string | undefined]> = [
    ["workspace.remote", config.workspace.remote],
    ...Object.entries(config.repositories).map(([name, repository]) => [`repositories.${name}.remote`, repository.remote] as [string, string | undefined]),
  ];
  for (const [path, value] of remotes) {
    const error = value ? remoteReferenceError(value) : null;
    if (error) errors.push(`${path} ${error}`);
  }
  for (const [index, source] of (config.context?.authoritative_sources ?? []).entries()) {
    const error = contextReferenceError(source.reference);
    if (error) errors.push(`context.authoritative_sources.${index}.reference ${error}`);
  }
  for (const [name, repository] of Object.entries(config.repositories)) {
    const normalized = repository.path.replace(/^\.\//, "").replace(/\/$/, "");
    const prior = paths.get(normalized);
    if (prior) errors.push(`repositories.${name}.path duplicates repositories.${prior}.path`);
    paths.set(normalized, name);
  }
  for (const [index, source] of (config.context?.authoritative_sources ?? []).entries()) {
    if (source.repository && !config.repositories[source.repository]) errors.push(`context.authoritative_sources.${index}.repository is not configured: ${source.repository}`);
  }
  const required = new Set(config.activity.required_capabilities);
  const declared = new Set([...config.activity.required_capabilities, ...config.activity.optional_capabilities]);
  for (const capability of config.activity.optional_capabilities) {
    if (required.has(capability)) errors.push(`activity capability is both required and optional: ${capability}`);
  }
  const lifecycle = config.activity.lifecycle ?? {};
  if (config.activity.provider === "none" && Object.values(lifecycle).some((actions) => (actions?.length ?? 0) > 0)) {
    errors.push("activity.lifecycle cannot configure external actions when provider is none");
  }
  for (const [event, actions] of Object.entries(lifecycle)) {
    const ids = new Set<string>();
    for (const action of actions ?? []) {
      if (ids.has(action.id)) errors.push(`activity.lifecycle.${event} has duplicate action id: ${action.id}`);
      ids.add(action.id);
      if (!declared.has(action.capability)) errors.push(`activity.lifecycle.${event}.${action.id} uses undeclared capability: ${action.capability}`);
      if (action.policy === "required" && !required.has(action.capability)) {
        errors.push(`required lifecycle action ${event}.${action.id} must use a required capability`);
      }
    }
  }
  if (config.workspace.mode === "team" && config.workflow.wrapper_change_policy !== "pull-request") {
    errors.push("team mode requires workflow.wrapper_change_policy: pull-request");
  }
  return errors;
}

export async function workspaceDocumentErrors(workspaceRoot: string, config: WorkspaceConfig): Promise<string[]> {
  const required = [
    ...requiredWorkspaceDocuments,
    ...new Set(Object.values(config.repositories).map((repository) => `agents/${repository.agent}.md`)),
  ];
  const errors: string[] = [];
  for (const path of required) {
    try {
      const info = await lstat(resolve(workspaceRoot, path));
      if (!info.isFile() || info.isSymbolicLink()) throw new Error("not a regular file");
    } catch {
      errors.push(`required workspace document is missing: ${path}`);
    }
  }
  return errors;
}
