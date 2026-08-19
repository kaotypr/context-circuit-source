# Context Circuit OKF profile

Status: implementation profile over Open Knowledge Format v0.2

This profile adopts the portable knowledge envelope defined by the pinned OKF
v0.2 specification while adding Context Circuit requirements for explicit
bundle participation, authority, acceptance, lifecycle, and generation
validation. Base OKF conformance and this stricter profile are separate
results. Passing one does not imply passing the other.

The upstream specification is preserved at `sources/okf-v0.2/SPEC.md`; its
revision and checksum are recorded in `sources/okf-v0.2/source.yaml`.

## Base OKF contract

For a declared bundle, OKF v0.2 requires:

- UTF-8 Markdown with an initial YAML front-matter block;
- a line-1 `---` opening delimiter and a standalone `---` closing delimiter;
- a non-empty `type` on every non-reserved concept;
- reserved `index.md` and `log.md` structures;
- tolerance for unknown types, unknown front-matter keys, broken links, and
  missing optional indexes.

OKF's only universally required concept key is `type`. `title`,
`description`, `resource`, `tags`, provenance, trust, and lifecycle families
are optional in base OKF. A bare `verified` mapping is normalized by readers
to a one-item list. Missing `verified` means unverified; it is not a failure.

Base OKF does not prescribe a schema registry, a command runtime, a central
authority, or a fixed taxonomy. It is deliberately permissive.

## Context Circuit profile

The profile applies only to knowledge-oriented Markdown in a declared bundle:

- Idea Briefs and PRDs;
- Workspace, Project, Architecture, Convention, and Decision Knowledge;
- Domain, Workflow, and Role Knowledge;
- design notes and references;
- optional plan rationale when explicitly published as knowledge.

The profile does not apply to `plan.yaml`, task contracts, `workspace.yaml`,
provenance registries, runtime YAML, normative instructions, skills, or
arbitrary raw sources. Those artifacts remain linkable from a concept but keep
their own authority and schemas. In particular, the presence of Markdown in a
directory never promotes that directory to an OKF bundle.

Each profiled concept requires:

```yaml
type: Project Knowledge
title: Project context
description: One sentence suitable for an index or retrieval preview.
status: draft
```

Profile requirements:

- `type`, `title`, `description`, and explicit `status` are required;
- `description` is one sentence suitable for indexes and retrieval;
- `status` is one of `draft`, `stable`, or `deprecated`;
- `tags` are optional and should use stable lowercase vocabulary;
- a declared bundle has unique, stable Concept IDs;
- unknown OKF and Context Circuit extension fields survive rewrites;
- claim footnotes may refer only to source IDs present in the concept's
  `sources` list when claim-level attribution is used;
- acceptance and authority extensions must agree with the owning human gate;
- `verified` remains advisory and never changes acceptance or lifecycle
  authority.

Unknown `type` values remain valid. The profile's initial type vocabulary is
descriptive rather than a registry:

`Idea Brief`, `Product Requirements`, `Workspace Knowledge`, `Project
Knowledge`, `Architecture`, `Convention`, `Decision Log`, `Domain`,
`Workflow`, `Role`, `Plan Rationale`, `Design Note`, and `Reference`.

Consumers must treat an unknown value as a generic concept rather than
rejecting it.

## Bundle and filename rules

The generator or accepted route must explicitly declare a bundle root. Within
that root:

- lowercase `index.md` is a reserved directory listing;
- lowercase `log.md` is a reserved newest-first update log;
- every other `.md`, including `README.md`, is a concept;
- the bundle-root `index.md` may carry `okf_version: "0.2"`; other index files
  have no front matter;
- Concept ID is the bundle-relative path without `.md`;
- standard Markdown links express relationships and may be broken without
  failing base OKF.

The repository's uppercase `context/INDEX.md` is a routing/instruction
artifact unless a separate explicit compatibility contract declares it as
part of a bundle. It is not the lowercase OKF reserved index by spelling
similarity. An index-like `README.md` is likewise an explicit migration input,
not an automatic reserved file.

## Compatibility and normalization

### `type` and legacy `kind`

New writers emit only `type`. A reader may temporarily accept legacy `kind`
with this deterministic order:

1. read `type` when present;
2. when `type` is absent, copy `kind` into the normalized consumer `type` and
   mark the input as legacy;
3. when both are present with unequal values, fail the hard artifact-schema
   and profile checks;
4. when neither is present, fail base OKF for a non-reserved concept.

The alias is only for OKF concept metadata. Runtime records may legitimately
use a `kind` field under their runtime schema and must not be normalized into
OKF.

### `verified`

Readers produce one consumer shape:

```yaml
verified:
  - by: process:document-check
    at: 2026-08-18T10:00:00Z
```

A list is retained as a list. A mapping becomes a one-item list. Writers may
choose the canonical list form, but a compatibility reader must accept both.
Every entry keeps the OKF actor and timestamp contract. This normalization is
not a human verification or acceptance action.

### Unknown extensions

Consumers and rewriters preserve unknown top-level and nested extension
fields. This includes fields beginning with `x-` and future OKF or Context
Circuit fields. Schema validation checks known fields without deleting
unknown data. A round trip must retain the original unknown value and
structure.

## Lifecycle, acceptance, and trust

The profile maps knowledge meaning to OKF status:

| Context Circuit meaning | OKF status |
| --- | --- |
| Proposed, incomplete, or awaiting acceptance | `draft` |
| Accepted and current | `stable` |
| Superseded or historical | `deprecated` |

Use a Context Circuit extension when a gate must be machine-visible:

```yaml
acceptance:
  state: pending
  gate: product-knowledge-acceptance
```

The human gate, not the schema, changes acceptance. The plan lifecycle is
never translated to this table: `plan.yaml` remains exactly `draft`,
`approved`, or `done`, task status remains its projection, and runtime
completion remains evidence pending the `status-change` gate.

`generated` records who or what produced the current content. `verified`
records advisory confirmation against sources or a resource. Neither grants
authority, satisfies a plan or Product Knowledge gate, or authorizes
publication, merge, deployment, or cleanup.

Attested Computation is explicitly deferred. OKF's attestation concepts do not
describe ordinary document generation, plan verification, or Context Circuit
runtime completion evidence.

## Profile validation

Generation validates the initial front matter before presenting a document as
ready:

1. require the exact opening delimiter on line 1;
2. use the first later standalone `---` as the closing delimiter;
3. parse only that block as one YAML document, leaving later Markdown
   horizontal rules in the body;
4. validate base OKF, the profile, and the selected artifact schema as
   separate hard checks;
5. repair and rerun the complete sequence after any failure.

The result schema at `schemas/documents/validation-result-v1.yaml` keeps
`base_okf`, `context_circuit_profile`, `artifact_schema`, and advisory
diagnostics distinct. Broken links, missing optional indexes, freshness,
missing recommended metadata, and copy-quality findings are advisory codes;
they must not be reported as base OKF failures.

The YAML parser and schema validator are a portable deterministic capability
provided by the host. This repository only supplies inspectable schema data,
fixtures, and the contract. It does not add or silently require a Node,
Python, Ruby, package-manager, or other command runtime. If a host cannot
provide the approved capability, a front-matter artifact remains not-ready
and the host reports the capability gap; it must not substitute grep-like
parsing.

## Schema references

- `schemas/documents/okf-concept-v1.yaml` — profiled concept envelope and
  compatibility rules.
- `schemas/documents/okf-index-v1.yaml` — reserved index structure.
- `schemas/documents/okf-log-v1.yaml` — reserved log structure.
- `schemas/documents/front-matter-extraction-v1.yaml` — deterministic
  delimiter extraction result.
- `schemas/documents/validation-result-v1.yaml` — independent hard and
  advisory result categories.

## References

- `sources/okf-v0.2/SPEC.md`, sections 3–5, 8–9, and 11–12.
- `sources/document-system-improvements/details/01-artifact-model.md`.
- `sources/document-system-improvements/details/05-okf-profile.md`.
- `docs/document-system.md` for the complete artifact-family matrix and
  lifecycle authority boundaries.
