import { access } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";
import { readData, schemaNames, validateContract, workspaceDocumentErrors, workspaceSemanticErrors, type SchemaName } from "./lib/validation.js";
import type { WorkspaceConfig } from "./lib/types.js";

const { values, positionals } = parseArgs({
  options: {
    schema: { type: "string", default: "workspace" },
    "check-paths": { type: "boolean", default: false },
    "check-documents": { type: "boolean", default: false },
  },
  allowPositionals: true,
});

const schema = values.schema as SchemaName;
if (!schemaNames.includes(schema)) throw new Error(`Unknown schema: ${schema}`);
const workspaceRoot = resolve(process.env.KAO_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const path = resolve(workspaceRoot, positionals[0] ?? "workspace.yaml");
const value = await readData(path);
const contractErrors = await validateContract(schema, value);
const errors = contractErrors.map((error) => `${error.instancePath || "/"} ${error.message}`);

if (schema === "workspace" && errors.length === 0) {
  const config = value as WorkspaceConfig;
  errors.push(...workspaceSemanticErrors(config));
  if (values["check-documents"]) errors.push(...await workspaceDocumentErrors(workspaceRoot, config));
  if (values["check-paths"]) {
    for (const [name, repository] of Object.entries(config.repositories)) {
      try {
        await access(resolve(workspaceRoot, repository.path));
      } catch {
        errors.push(`repositories.${name}.path does not exist: ${repository.path}`);
      }
    }
  }
}

if (errors.length > 0) {
  console.error(`Invalid ${schema} document ${path}:`);
  for (const error of errors) console.error(`- ${error}`);
  process.exitCode = 1;
} else {
  console.log(`Valid ${schema}: ${path}`);
}
