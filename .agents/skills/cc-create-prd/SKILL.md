---
name: cc-create-prd
description: Create a source-grounded product requirements document with explicit uncertainty and human acceptance.
---

# Create a PRD

Use this skill when the user asks for a PRD, product requirements, or a durable
definition of a feature or project. Start from the user's stated intent; an
Idea Brief is useful but not mandatory when the request is already clear.

## Read

Read only the relevant Product Knowledge, accepted Idea Brief if one exists,
and repository evidence needed for the requested scope. For source-based work,
identify the selected files under `sources/`, explain why they are relevant,
and read only those files. Normal PRD creation must not ingest the source
inbox.

Record every source path, reason, and known revision in `Provenance`. Preserve
raw source material in `sources/` rather than copying it into the PRD or
accepted context.

## Output

Draft the PRD at a user/team-selected path under `sources/` using the
contract in `docs/prd.md`. Do not impose a subdirectory or filename
convention. Record the exact chosen path in the PRD's `Provenance` section.
Emit `type: Product Requirements`, a non-empty `title`, a one-sentence
`description`, and the OKF lifecycle `status: draft`. Include a separate
pending `acceptance` extension with the human acceptance gate. Include the
problem and outcome, users, requirements, non-goals, scenarios and acceptance
criteria, constraints, dependencies, assumptions, open questions, Product
Knowledge references, and provenance. Do not emit a new `kind` field,
unsupported sections, or generation comments.

The selected path is an OKF concept only if the request or an accepted
workspace contract explicitly declares its containing bundle. Never treat all
of `sources/` as a bundle. When no bundle is declared, validate the artifact
schema while reporting the OKF checks as `not-applicable`.

## Generation validation

Before presenting the PRD as ready, use the host-provided
`deterministic-yaml-schema-validation` capability to:

1. require `---` on line 1 and close at the first later standalone `---`;
2. parse only that initial block as one deterministic YAML document;
3. validate the PRD artifact schema and, only for a declared bundle, base OKF
   and the Context Circuit profile; and
4. repair hard failures and rerun the complete sequence from the beginning.

Keep front-matter extraction, YAML, base OKF, profile, artifact schema, and
advisories as independent result categories. Advisory findings never make
output blocked. If the capability is unavailable, leave the PRD not-ready.

## Gate and uncertainty

Human gate: ask the user to accept or revise the PRD before treating it as
accepted.
Distinguish requirements from evidence, assumptions, proposals, and open
questions. Ask the user to accept or revise the PRD. A changed or contradictory
source is surfaced for human review; it is not silently reconciled.

## Next action

After acceptance, offer the optional next action: refresh Product Knowledge,
draft an approved plan, or stop with the PRD as the durable artifact. Do not
silently approve a plan or treat the PRD as implementation authorization.
