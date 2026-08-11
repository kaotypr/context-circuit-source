import { readFile } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { Ajv2020, type ErrorObject } from "ajv/dist/2020.js";
import { parse as parseYaml } from "yaml";
import type { WorkspaceConfig } from "./types.js";

export const schemaNames = ["workspace", "task-brief", "worker-result", "verifier-result", "runtime-manifest"] as const;
export type SchemaName = (typeof schemaNames)[number];

const projectRoot = resolve(dirname(fileURLToPath(import.meta.url)), "../..");

export async function readData(path: string): Promise<unknown> {
  const raw = await readFile(path, "utf8");
  return path.endsWith(".yaml") || path.endsWith(".yml") ? parseYaml(raw) : JSON.parse(raw);
}

export async function validateContract(name: SchemaName, value: unknown): Promise<ErrorObject[]> {
  const schema = JSON.parse(await readFile(join(projectRoot, ".agents", "contracts", `${name}.schema.json`), "utf8"));
  const ajv = new Ajv2020({ allErrors: true, strict: false });
  ajv.addFormat("date-time", {
    type: "string",
    validate: (value: string) => /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/.test(value) && !Number.isNaN(Date.parse(value)),
  });
  const validate = ajv.compile(schema);
  return validate(value) ? [] : [...(validate.errors ?? [])];
}

export function workspaceSemanticErrors(config: WorkspaceConfig): string[] {
  const errors: string[] = [];
  const paths = new Map<string, string>();
  for (const [name, repository] of Object.entries(config.repositories)) {
    const normalized = repository.path.replace(/\/$/, "");
    const prior = paths.get(normalized);
    if (prior) errors.push(`repositories.${name}.path duplicates repositories.${prior}.path`);
    paths.set(normalized, name);
  }
  const required = new Set(config.activity.required_capabilities);
  for (const capability of config.activity.optional_capabilities) {
    if (required.has(capability)) errors.push(`activity capability is both required and optional: ${capability}`);
  }
  if (config.workspace.mode === "team" && config.workflow.wrapper_change_policy !== "pull-request") {
    errors.push("team mode requires workflow.wrapper_change_policy: pull-request");
  }
  return errors;
}
