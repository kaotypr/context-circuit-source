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

Draft `contributions/prds/<slug>.md` using the contract in `docs/prd.md`. Include
the problem and outcome, users, requirements, non-goals, scenarios and
acceptance criteria, constraints, dependencies, assumptions, open questions,
Product Knowledge references, and provenance. Mark it `status: draft` until
the user accepts it.

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
