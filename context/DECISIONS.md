# Decisions

## 2026-08-14 — Deliberate simplification

- Context Circuit is a human-controlled context, planning, and agent-work tool.
- Small YAML parsers and Markdown lint checks define authored plans, tasks, and agent handoffs.
- Plans and tasks have only `draft`, `approved`, and `done`; status changes are explicit human actions.
- Independent verification, repair, activity lifecycle, publication events, merge confirmation, closeout records, and automatic status synchronization are not required workflow concepts.
- Plans archive by moving to `archives/plans/<repository-key>-plans/<number>-<slug>/`; archiving preserves status and content.
- Product Knowledge provenance is a small `context/sources.yaml` register. Source changes produce a refresh proposal and never cause silent canonical edits.
