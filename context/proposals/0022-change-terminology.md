# Proposal 0022 — change: TERMINOLOGY.md (v0.6 terms)

- id: 0022-change-terminology
- target_context_unit: context/TERMINOLOGY.md
- operation: change
- confidence: high
- status: review-needed
- related_plan: —
- affected_commits: [74eb510, e3e95cd, 048eb67]
- evidence_refs:
  - sources/system-design/context-circuit/v0.6/
  - wrapper/contracts/schemas/{lease,grounding-manifest,plan,execution}.yaml
  - wrapper/adapters/writer-brief.md
  - docs/terminology.md (canonical glossary; add the same terms there on acceptance)

## Statement

Add the v0.6 terms to `context/TERMINOLOGY.md` (and the canonical
`docs/terminology.md`). Proposed glossary entries:

- **Plan stack** — a named set of approved plans executed in one run.
- **Plan dependency** — inter-plan ordering (`plan_dependencies`), distinct from a
  task's intra-plan `depends_on`; declaring it requires plan `schema_version: 2`.
- **Path lease** — a `(repository, path-region)` reservation extending the
  one-writer lock; overlap = equal / path-prefix ancestor / repository-wide `.`;
  descendants exempt; held until delivery.
- **Execution base** — the tip a plan's worktree is built on: **anchor tip** (no
  same-repo predecessor), a **stacked** predecessor branch (one), or a
  runtime-authored **integration base** (a `--no-ff` merge of ≥2 predecessor
  branches), recorded as `base_commit` with `based_on`; the kept ref lives at
  `refs/cc-base/<plan>/<repo>`.
- **Drift guard** — rebase-onto-current-anchor-tip + re-verify before a pull
  request when the recorded base has diverged.
- **Repository grounding** — the writer reading and honoring the target
  repository's own agent guidance, discovered live from the worktree.
- **Grounding manifest** — the discovered data (files, skills+descriptions,
  environment) recorded per execution.
- **Worktree hardening** — deterministic toolchain detection/preparation reported
  as `environment: ready | no-toolchain`.
- **Writer brief** — the fixed template (`writer-brief.md`), filled from the
  grounding manifest and the plan and delivered to the worker; the coordinator
  authors only the task focus.
- **System design (authoring)** — a structured source describing the shape of a
  change, authored via the `cc-system-design` skill; a source, not a lifecycle
  stage (no status/gate).

## Change

Add the entries above to `context/TERMINOLOGY.md`; mirror them in
`docs/terminology.md` (the canonical glossary) so the internal→user-facing
projection stays complete. Keep user-facing translations plain (e.g. an execution
base and integration base are internal mechanics the coordinator never exposes to
a lay user).
