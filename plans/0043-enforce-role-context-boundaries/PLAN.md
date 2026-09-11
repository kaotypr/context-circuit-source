# Plan 0043 — Enforce role-specific context boundaries

**Intent:** i031-enforce-role-context-boundaries  
**Repository:** agent-harness  
**Tier:** Standard  
**Status:** draft

## Objective

Make strict Agent Harness results prove that every created role stayed within
its declared normal or explicitly expanded context boundary and hard budget.
Available, provided, active-read, and charged context remain separate, and a
host that cannot prove a required measurement reports an inconclusive context
result instead of passing on an assumed zero.

## Grounding

The clean target at `93ebd8b99b30eb667b02e47e15299e1aa4f6b85d`
already resolves declarative scenarios, builds independent run worlds, launches
Codex, Claude Code, and Cursor, preserves root/child topology and native events,
captures a host usage object, redacts streams, and aggregates host outcomes.
It has no role envelope or expansion contract, context-event vocabulary,
per-role ledger, boundary/budget grader, contributor report, or honest rule for
missing context measurement.

## Decisions

- Use one plan: scenario authority, native event translation, ledger, grading,
  and end-to-end reporting form one compatibility and verification seam.
- Freeze normal envelopes, hard budgets, and exact expansion authority in the
  resolved scenario inputs so reproduction cannot change what was authorized.
- Build ledgers only from enforceable views or positively attributed native
  evidence. Missing role, resource, or measurement proof stays unavailable or
  unattributed.
- Grade per role before host and matrix aggregation. Observed boundary or budget
  violations fail; missing proof is inconclusive and never erases a known fail.
- Rank contributors from resource events. Host token counts are supplemental and
  enforceable only when trustworthy role-level attribution exists.

## Tasks

| id | title | depends on |
|---|---|---|
| CB-001 | Define and freeze role context contracts | — |
| CB-002 | Build attributed role context ledgers | CB-001 |
| CB-003 | Enforce boundaries and budgets in strict grading | CB-002 |
| CB-004 | Prove strict context behavior end to end | CB-003 |

Canonical paths, changes, acceptance statements, and runnable checks are in
`plan.yaml`.

## Risks

- Host protocols expose different resource and token events. Translation must
  preserve native evidence and report missing attribution, not synthesize it.
- Process-wide usage can look like role usage. Only trustworthy role-linked
  telemetry may satisfy or breach a role token ceiling.
- Resource identity normalization can accidentally broaden authority. Exact
  named surfaces and explicit adjacency tests bound matching.
- Background index traffic may resemble an active read. It remains a distinct
  dimension unless native evidence proves delivery to the role session.
- A strict boundary pass needs either enforced unavailability or positive read
  evidence; replay fixtures remain parser tests and cannot prove a live host.

## Verification

Run the focused contract, evidence, driver, grading, lifecycle, matrix, and
end-to-end tests, then the full unit and compile suites. The opt-in live command
checks all configured native hosts and may honestly return inconclusive when a
host cannot supply the required role context proof.
