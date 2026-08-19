# Context Circuit document system

Status: implementation contract for generated documents

This document defines the shared document-family, authority, lifecycle, and
generation-readiness contract. It is intentionally separate from the
Open Knowledge Format (OKF) v0.2 specification and from any host's parser
implementation. The upstream OKF specification is the source for the
interoperability rules; this document defines where Context Circuit applies
those rules and where it does not.

## Contract vocabulary

- **Artifact family** identifies the purpose and audience of a file.
- **Authority** identifies which artifact owns a fact or lifecycle state.
- **Declared bundle** is a directory explicitly named by a generator, profile,
  or accepted workspace contract as an OKF bundle root. A directory is not an
  OKF bundle merely because it contains Markdown.
- **OKF concept** is a non-reserved Markdown document inside a declared bundle.
- **Artifact schema** is the versioned Context Circuit contract for one
  artifact family. Schema data is inspectable; its presence does not install a
  parser or authorize a migration.
- **Generation-ready** means that every applicable hard check has passed. An
  advisory finding never makes a hard check pass or fail.

The base OKF result, the Context Circuit profile result, the applicable
artifact-schema result, and advisory diagnostics are reported independently.
They must not be collapsed into one generic “document valid” result.

## Artifact-family matrix

| Artifact family | Audience | Authority and owner | Lifecycle | Format | OKF applicability |
| --- | --- | --- | --- | --- | --- |
| Idea Brief / PRD | Human reviewers and agents | User intent and accepted product artifact; the human acceptance gate owns acceptance | OKF `draft`, `stable`, or `deprecated`; Context Circuit acceptance is an extension and a gate | Markdown with YAML front matter | OKF concept only when inside a declared knowledge bundle |
| Product, Workspace, Project, Domain, Workflow, Role, Architecture, Convention, Decision, Design, or Reference Knowledge | Human reviewers and selective agent readers | Accepted Product Knowledge and its cited sources; the human acceptance gate owns acceptance | OKF lifecycle mapped from proposed/current/superseded meaning | Markdown with YAML front matter | OKF concept inside a declared knowledge bundle |
| Plan rationale | Human reviewers | The associated plan and cited evidence; `plan.yaml` remains lifecycle authority | Proposed or historical knowledge lifecycle when published as a concept | Markdown with YAML front matter | Conditional: OKF only when explicitly published into a declared knowledge bundle |
| `plan.yaml` | Human reviewers, coordinator, and writers | Canonical plan lifecycle record | Exactly `draft`, `approved`, or `done`; `cc-approve-plan` and `cc-finish-plan` remain the gates | Structured YAML | Not OKF |
| Task contract | Writer, verifier, and human reviewer | Task contract under the canonical plan; task status is a projection of plan status | Exactly `draft`, `ready`, or `done`; not a second approval or lease | Markdown with initial YAML front matter, validated by the task schema | Not OKF |
| Workspace configuration and provenance registries | Hosts and workspace coordinators | `workspace.yaml` or the named provenance registry | Configuration and registry semantics defined by their owning contracts | Structured YAML | Not OKF |
| Session, delegation, lease, prompt, evidence, handoff, completion, stack graph, and stack progress records | Agent-only runtime readers | `.runtime/` ownership and lease contracts | Runtime states such as `executing`, `blocked`, and `implemented`; they never redefine product or plan truth | Structured YAML, with optional handoff Markdown summary | Not OKF |
| `AGENTS.md`, `WORKFLOW.md`, skills, and other normative instructions | Agents and hosts | The normative instruction owner | Instruction revisions are not knowledge lifecycle | Markdown | Not OKF |
| Raw or arbitrary source files | The user or team as source owner | The passive `sources/` inbox and the source's own authority | Source-specific; reading is request-scoped | Original source format | Not OKF |
| `index.md` and `log.md` inside a declared bundle | Human and agent discovery | The declared bundle | Reserved OKF behavior, not concept lifecycle | Markdown with the OKF reserved structures | OKF reserved files, not concepts |
| `README.md` and uppercase `INDEX.md` | Existing humans, routes, or compatibility readers | The local contract that explicitly names them | No implicit OKF lifecycle | Existing Markdown contract | Compatibility inputs only; never silently treated as reserved OKF files |

The matrix does not move an artifact between authorities. In particular,
`verified` is trust metadata and remains advisory; it cannot approve a Product
Knowledge artifact, approve a plan, change task status, complete a plan, or
authorize publication.

## Declared OKF boundaries

A generator must identify the bundle root before applying OKF rules. The
declaration is explicit and stable for the generated artifact; it may be a
profile entry, a route-specific contract, or another accepted workspace
record. Discovery must not use “directory contains a Markdown file” as a
bundle test.

Within a declared bundle:

1. `index.md` and `log.md` use their OKF reserved structures at every level.
2. Every other `.md` file is a concept and must have a parseable initial
   front-matter block and a non-empty `type` for base OKF conformance.
3. The Concept ID is the bundle-relative path without `.md`.
4. `README.md` is a concept unless an explicit compatibility contract says it
   is being migrated to `index.md`.
5. A bundle-root `index.md` may declare `okf_version: "0.2"`; an index at
   another level has no front matter.

Outside a declared bundle, ordinary Markdown is governed by its local
artifact family. This keeps the routing `context/INDEX.md`, plans, runtime
records, normative instructions, and arbitrary sources outside OKF even when
they contain headings or links that resemble a knowledge document.

## Shared vocabulary and compatibility

New OKF writers emit `type`, never a new `kind` field. A reader may support
historical input using this exact compatibility rule:

1. If `type` is present, use it.
2. If `type` is absent and `kind` is present, normalize `kind` to the
   consumer's `type` value and report a legacy-compatibility diagnostic.
3. If both are present and their values differ, fail the applicable hard
   schema/profile check. Do not guess which value is authoritative.
4. If neither is present, a non-reserved OKF concept fails base OKF because
   its required type is missing.

The compatibility rule applies to OKF concept metadata only. It must not
translate runtime `kind` fields, such as a runtime record's record kind, into
OKF `type` fields.

Consumers normalize a bare `verified` mapping to a one-item list:

```yaml
verified:
  - by: human:kaotypr
    at: 2026-08-18T09:35:00Z
```

An existing list remains a list. Missing `verified` means unverified and is
not a base OKF failure. Rewriters must preserve every unknown OKF or Context
Circuit extension field, including nested unknown fields. A reader that cannot
understand an extension may carry it as opaque data; it must not silently drop
it.

## Lifecycle and human gates

Knowledge concepts use the OKF lifecycle mapping:

| Context Circuit meaning | OKF `status` |
| --- | --- |
| Proposed, incomplete, or awaiting acceptance | `draft` |
| Accepted and current | `stable` |
| Superseded or retained for history | `deprecated` |

`acceptance.state`, `acceptance.gate`, and human evidence make the Context
Circuit gate visible, but schema validation does not satisfy that gate.

Plans retain exactly `draft` → `approved` → `done` in `plan.yaml`. Included task
status remains the synchronized projection `draft` → `ready` → `done`.
Runtime `implemented`, `blocked`, and completion evidence remain runtime
observations. A document contract must never rewrite those lifecycle meanings
as OKF statuses.

## Generation-time validation

Every generated Markdown artifact with front matter must pass the complete
applicable validation sequence before it is presented, committed, or used as
ready output:

1. **Extract the initial block.** Line 1 must be exactly `---`. Starting at
   line 2, the first later line that is exactly `---` closes the front matter.
   The body begins on the following line.
2. **Ignore later horizontal rules.** After the first closing delimiter, every
   later standalone `---` belongs to the Markdown body. It must not reopen or
   corrupt front-matter extraction.
3. **Parse YAML.** The host-provided deterministic capability parses the
   extracted block as one YAML document and returns structured data or a
   stable parse error. Text matching or a shell-only YAML approximation is not
   a substitute.
4. **Validate the artifact schema.** Check required fields, value types,
   enumerations, nesting, compatibility aliases, and extension preservation
   against the versioned schema selected by the artifact family.
5. **Validate applicable OKF and profile rules.** Reserved filenames,
   non-empty `type`, bundle identity, profile-required fields, and Concept ID
   rules apply only when the artifact is an explicitly declared OKF bundle
   member.
6. **Repair and repeat.** If a hard check fails, repair the generated
   artifact within scope and rerun extraction, YAML parsing, schema checks,
   and applicable OKF/profile checks from the beginning. A prior partial pass
   is not reusable.

The portable capability boundary is documented in
`docs/host-capabilities.md` and recorded in the approved runtime decision
outside this repository. It is a host-provided deterministic YAML/schema
capability. The repository does not add or silently require Node, Python,
Ruby, a package manager, or another command runtime.

### Validation result contract

The versioned result schema at
`schemas/documents/validation-result-v1.yaml` uses these independent
categories:

| Result category | Values | Blocking behavior |
| --- | --- | --- |
| `front_matter_extraction` | `pass`, `fail`, `not-applicable` | `fail` blocks front-matter output |
| `yaml_parse` | `pass`, `fail`, `not-applicable` | `fail` blocks output |
| `base_okf` | `pass`, `fail`, `not-applicable` | `fail` blocks only an artifact declared as OKF |
| `context_circuit_profile` | `pass`, `fail`, `not-applicable` | `fail` blocks only an artifact declared for the profile |
| `artifact_schema` | `pass`, `fail`, `not-applicable` | `fail` blocks the selected artifact family |
| `advisories` | zero or more stable diagnostic codes | Never changes a hard result |

The top-level `status` is `ready` only when every applicable hard category is
`pass`. It is `blocked` when any applicable hard category is `fail`.
`not-applicable` is explicit; it is not an implicit pass. Advisory examples
include broken links, a missing optional `index.md`, freshness, missing
recommended metadata, and copy or retrieval-quality concerns. Such findings
must not be reported as base OKF failures because OKF v0.2 explicitly
tolerates them.

## Schema catalog

The inspectable, versioned data contracts are catalogued in
`schemas/documents/README.md`. They cover the OKF concept and reserved files,
front-matter extraction, plan and task artifacts, runtime records, and
validation results. The catalog defines compatibility and preservation
semantics but does not itself authorize a migration or provide a parser.

## References

- `sources/okf-v0.2/SPEC.md`, sections 3–5, 8–9, and 11–12 — base bundle,
  front matter, reserved files, trust normalization, and permissive
  conformance.
- `sources/document-system-improvements/details/01-artifact-model.md` —
  artifact audiences, authority, and OKF boundary.
- `sources/document-system-improvements/details/05-okf-profile.md` —
  Context Circuit profile, lifecycle mapping, compatibility, and validation.
- `docs/runtime-contract.md` — runtime ownership, lifecycle, and human-gate
  semantics that remain outside OKF.
