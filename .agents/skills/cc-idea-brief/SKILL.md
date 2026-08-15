---
name: cc-idea-brief
description: Capture an uncertain idea or intent in the smallest useful human-reviewable brief.
---

# Capture an Idea Brief

Use this skill when the user is exploring an idea, clarifying intent, or asks
for a lightweight brief. An Idea Brief is sufficient on its own; do not force a
PRD or plan when the user does not need one.

## Read

Start from the user's conversation and the smallest relevant Product Knowledge.
If the request names or requires source material, identify the specific files,
state why each will be read, and read only those files under `sources/`. Do not
scan or ingest the source inbox by default.

Record the files and reasons in the brief's `Provenance` section. Keep raw
source text in `sources/`; summarize relevant evidence instead of copying it.

## Output

Draft the Idea Brief at a user/team-selected path under `sources/` using the
contract in `docs/idea-brief.md`. Do not impose a subdirectory or filename
convention. Record the exact chosen path in the brief's `Provenance` section.
Include the user's intent, desired outcome, audience or users, constraints,
assumptions, open questions, relevant Product Knowledge, and source
provenance. Mark it `status: draft` until accepted.

## Gate and uncertainty

Human gate: ask the user to accept, revise, or leave the brief as a draft.
Separate observed facts, the user's stated intent, assumptions, and unresolved
questions. Ask the user to accept, revise, or leave the brief as a draft before
calling it accepted. Never silently turn an assumption into a decision.

## Next action

After acceptance, offer the optional next action that fits the request: refresh
context, create a PRD, create a plan, or begin a small implementation. The user
may stop after the Idea Brief.
