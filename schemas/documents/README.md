# Document schema catalog

This directory contains versioned, inspectable data contracts for Context
Circuit document families. The schemas describe inputs and deterministic
validation results; they do not install a parser, define a command-line
interface, or authorize a migration.

## Schema notation

Each YAML schema contains:

- `schema_version: 1` — version of the schema document;
- `id` and `version` — stable contract identity and artifact-contract version;
- `applies_to` — the authority and representation selected by the generator;
- `required` and `properties` — the fields and value shapes to check;
- `additional_properties: preserve` — unknown fields survive reads and
  rewrites;
- `compatibility` — explicit reader and writer behavior where historical
  inputs exist.

The notation is intentionally data-only. A host supplies the deterministic
YAML parser and schema-validation capability described in
`docs/host-capabilities.md`; no repository runtime is implied.

## Catalog

| Schema | Applies to | OKF status |
| --- | --- | --- |
| `okf-concept-v1.yaml` | Non-reserved Markdown concepts in declared bundles | Base OKF plus the Context Circuit profile |
| `okf-index-v1.yaml` | Reserved `index.md` files | Base OKF reserved structure |
| `okf-log-v1.yaml` | Reserved `log.md` files | Base OKF reserved structure |
| `front-matter-extraction-v1.yaml` | Initial front-matter extraction for any applicable Markdown artifact | Pre-schema hard check |
| `plan-v1.yaml` | Canonical `plan.yaml` | Non-OKF plan authority |
| `task-v1.yaml` | Task Markdown with YAML front matter | Non-OKF task authority |
| `instruction-skill-v1.yaml` | Existing `.agents/skills/*/SKILL.md` and other declared instruction front matter | Non-OKF normative-instruction authority |
| `runtime-record-v1.yaml` | Agent-only runtime YAML records | Non-OKF runtime authority |
| `validation-result-v1.yaml` | Generation-time validation result | Non-OKF diagnostic evidence |

## Selection rules

1. Select a schema from the artifact-family matrix in
   `docs/document-system.md`, not from a filename alone.
2. Apply `front-matter-extraction-v1.yaml` before parsing or schema validation
   for Markdown with front matter.
3. Apply base OKF and `okf-concept-v1.yaml` only after an explicit bundle
   declaration.
4. Select `instruction-skill-v1.yaml` for existing skill metadata such as
   `name` and `description`; it validates the front matter without promoting
   skills or instructions to OKF concepts.
5. Keep plan, task, runtime, instruction, and source artifacts on their
   non-OKF schemas and authorities.
6. Use `validation-result-v1.yaml` to report each applicable hard category and
   advisory diagnostics independently.

Schema data is additive and inspectable. Historical artifacts remain readable
through the compatibility rules; this catalog does not require bulk migration.
