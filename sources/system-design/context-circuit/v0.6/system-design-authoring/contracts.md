# System-design authoring — contracts

This scope has an intentionally tiny surface: **one shipped skill and its
allowlist wiring**. Nothing in the runtime changes.

## The one owner

| Owner | Change |
| --- | --- |
| `.agents/skills/cc-system-design/SKILL.md` **(new, shipped)** | the authoring skill: how to author and structure a system design under `sources/system-design/<product>/<version>/<scope>/` — the three-tier layout, the detail-per-file rubric, scope-by-concern separation, and embedded mermaid (see [authoring-rubric.md](./authoring-rubric.md)). Read-as-procedure packet, resolved by path; purely skill-invoked (`/cc-system-design`), no WORKFLOW action. It drafts structured source files and nothing else — it never approves, accepts, plans, executes, or writes Product Knowledge (INV-SKILL-01). |

## Allowlist wiring (the only non-skill edits)

Adding a shipped skill requires the four coordinated allowlist updates, exactly as
`cc-run-stack` did — no other files change:

- `scripts/release-manifest.txt` — add `required .agents/skills/cc-system-design/SKILL.md`.
- `scripts/release-artifact.sh` — add `cc-system-design` to the shipped-skill allowlist.
- `test/release/test-release.sh` — add `cc-system-design` to the allowlist and bump the shipped-skill count.
- `test/contracts/test-contracts.sh` — add `cc-system-design` to the required-skill list.

## What this scope does NOT change

- **No `wrapper/runtime/engine.sh` change** — no new function, no dispatch verb.
- **No `.runtime` record and no new schema** — a system design has no status,
  acceptance record, or lifecycle state; it is a source.
- **No new invariant** — the skill owns its own authoring rules; there is no
  runtime-enforced rule to index. (INV-SKILL-01 already governs how skills ship
  and are resolved.)
- **No `wrapper/adapters/WORKFLOW.md` action** — the skill is purely
  skill-invoked; there is no *design the system* conversational route.
- **No `plan.yaml` change, no `manifest.yaml` version change** — the plan schema
  `[1, 2]` and `runtime_version 0.6.0` were set by the run-stack scope.
- **No first-class `design/` area** — the design stays under `sources/`.
