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
`context/INDEX.md`, `docs/product-knowledge.md`, and the relevant template.
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

The approved product-direction source may establish the model boundaries, but
it does not authorize reading unrelated raw sources. Keep raw source text in
`sources/`; summarize evidence in the generated document and provenance
record.

## Generate Domain Knowledge

Use `docs/templates/domain-context.md` and write the canonical page to
`context/domains/<domain>/README.md`. Create a workflow page below
`context/domains/<domain>/workflows/` only when the evidence establishes a
distinct workflow. Update `context/domains/README.md` with a routing link, not
duplicated facts.

Include only supported categories. The page may describe summary, scope,
behavior, workflows, interfaces, data, constraints, edge cases,
implementation references, and verification. Do not fill unsupported sections
with guesses or create empty category files.

## Generate Role Knowledge

Use `docs/templates/role-context.md` and write the canonical page to
`context/roles/<role>.md`. Add a simple metadata list such as:

```yaml
domains:
  - ../domains/product/README.md
  - ../domains/planning/README.md
```

The role body explains the role's definition, outcomes, product surfaces,
cross-domain story, handoffs, role-specific behavior, and limitations. Link to
the exact domain or workflow page that owns behavior instead of copying it.
Do not invent `contributes_to`, `acts_on`, `consumes`, `approves`, or any other
relationship matrix. Update `context/roles/README.md` with a short routing
link.

## Metadata, provenance, and acceptance

Every generated page records, when known:

- `kind` and `status` (`proposed`, `accepted`, or `needs-review`);
- the selected `sources` and `source_revisions`;
- `generated_at`, `review_date`, and `freshness`;
- `assumptions`, `unknowns`, and `contradictions`; and
- an `acceptance` state, accepted date, and human accepter when accepted.

Record the source path, reason, revision or freshness signal, and resulting
document in `context/sources.yaml`. Do not copy raw source text into a domain
or role page.

Generated material is `status: proposed` with a pending acceptance state until
the user accepts it. If evidence is incomplete, keep the unknown visible. If
sources or an existing accepted page contradict one another, preserve the
accepted decision, create a visible proposed refresh or contradiction note,
and stop for human review. Never silently overwrite accepted context.

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
