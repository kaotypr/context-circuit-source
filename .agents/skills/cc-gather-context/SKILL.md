---
name: cc-gather-context
description: Create or refresh request-scoped Domain or Role Knowledge from selected evidence.
---

# Gather Domain or Role Knowledge

Use this skill when the user asks to create, refresh, or select detailed
Product Knowledge for a bounded domain or a business/project role. It supports
two complementary views:

- Domain Knowledge is canonical for one bounded project area and its exact
  behavior, workflows, interfaces, data, constraints, edge cases, and verified
  implementation references.
- Role Knowledge is a cross-domain perspective for a business or project role.
  It links to canonical domains and does not duplicate their facts.

Business/project roles are not agent execution roles. Coordinator, worker, and
verifier instructions remain under `agents/`; these role pages are not
technical access-control policies.

## Read and select evidence

Start with the user's request, `AGENTS.md`, `WORKFLOW.md`, `workspace.yaml`,
`context/index.md`, `docs/product-knowledge.md`, and the relevant template.
Read the smallest existing context set needed for the requested domain or
role. Do not scan all of `context/domains/`, `context/roles/`, or `sources/`.

For a source-based request:

1. Identify the requested knowledge kind and bounded slug.
2. Propose the exact raw source files and repository documents needed, with a
   one-line reason for each.
3. Read only those selected files. A source path not selected for the request
   must not be opened as a side effect.
4. Separate observed evidence, current implementation, accepted decisions,
   assumptions, unknowns, and contradictions.

Accepted project identity in `PROJECT.md`, or a human-reviewed Idea Brief or
PRD, may establish model boundaries. It does not authorize reading unrelated
raw sources. Do not refresh project identity from a separate direction
summary. Keep raw source text in `sources/`; summarize evidence in the
generated document and provenance record.

## Generate Product Knowledge

When a request explicitly asks for a Product Knowledge refresh, use
`docs/templates/product-context.md` and update the canonical project artifact
selected by the workspace contract. If that artifact is `context/PROJECT.md`,
it remains part of the top-level routing and metadata layer rather than an
implicit OKF bundle. Emit `type: Product Knowledge`, a title, a description,
and `status: draft` with a separate pending acceptance gate. Preserve the
workspace's project-identity boundary and do not rewrite an uninitialized page
without the user's project decision.

## Generate Domain Knowledge

Use `docs/templates/domain-context.md` and write the canonical domain concept
to `context/domains/<domain>/README.md`. Create a workflow concept below
`context/domains/<domain>/workflows/` only when the evidence establishes a
distinct workflow, using `docs/templates/workflow-context.md`. Update
`context/domains/index.md` with a routing link, not duplicated facts.

Include only supported categories. The page may describe summary, scope,
behavior, workflows, interfaces, data, constraints, edge cases,
implementation references, and verification. Do not fill unsupported sections
with guesses or create empty category files.

## Generate Role Knowledge

Use `docs/templates/role-context.md` and write the canonical page to
`context/roles/<role>.md`. Update `context/roles/index.md` with a short
routing link. Add a simple metadata list such as:

```yaml
domains:
  - ../domains/product/README.md
  - ../domains/planning/README.md
```
The role body explains the role's definition, outcomes, product surfaces,
cross-domain story, handoffs, role-specific behavior, and limitations. Link to
the exact domain or workflow page that owns behavior instead of copying it.
Do not invent `contributes_to`, `acts_on`, `consumes`, `approves`, or any other
relationship matrix.

## Metadata, provenance, and acceptance

Every generated Product, Domain, Workflow, or Role page records, when known:

- `type`, `title`, `description`, and OKF `status` (`draft`, `stable`, or
  `deprecated`);
- the selected `sources` and `source_revisions`;
- `generated_at`, `review_date`, and `freshness`;
- `assumptions`, `unknowns`, and `contradictions`; and
- an `acceptance` state, accepted date, and human accepter when accepted.

Record the source path, reason, revision or freshness signal, and resulting
document in `context/sources.yaml`. Do not copy raw source text into a domain
or role page.

Generated material is `status: draft` with a pending acceptance state until the
user accepts it. If evidence is incomplete, keep the unknown visible. If
sources or an existing accepted page contradict one another, preserve the
accepted decision, create a visible proposed refresh or contradiction note,
and stop for human review. Never silently overwrite accepted context.

New output omits `kind`, unsupported sections, and generation comments.
Compatibility readers may normalize a legacy `kind` only when `type` is absent;
conflicting `kind` and `type` values fail validation. Unknown top-level and
nested extension fields survive refreshes. A bare `verified` mapping is
normalized to a one-item list and remains advisory.

## Bundle declaration and generation validation

`context/domains/` and `context/roles/` are the current explicit OKF bundle
roots. Their lowercase root `index.md` files are reserved indexes. The
top-level `context/` directory is a routing and metadata layer, not an
implicit bundle, and the entire `sources/` inbox is never a bundle by default.
A user-selected source directory is a bundle only when the request or accepted
contract declares it.

Before presenting any front-matter-bearing page as ready, use the host-provided
`deterministic-yaml-schema-validation` capability to extract and parse the
initial block, validate the selected artifact schema, apply OKF/profile checks
only for a declared bundle, and rerun every applicable hard check after repair.
Keep extraction, YAML, base OKF, profile, artifact schema, and advisory results
independent. An unavailable capability leaves the page not-ready; no Node,
Python, Ruby, package manager, or repository runtime may be added.

## Refresh and context selection

When the canonical page already exists, update that page rather than creating a
duplicate path. Preserve accepted decisions and provenance that remains valid;
add changed evidence and a visible review state. A role refresh loads the role
page and only the domain/workflow pages linked by its `domains` and workflow
links. A later planning session should select the smallest relevant domain and
role set from the indexes rather than loading the entire context tree.

## Output

Report the requested kind and slug, selected evidence and reasons, generated
paths, provenance recorded, assumptions, unknowns, contradictions, acceptance
state, and the next human review action. Do not call the document accepted just
because generation or tests succeeded.
