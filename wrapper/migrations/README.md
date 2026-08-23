# Wrapper migrations

The manifest version is compared with the live session and context receipt.
`compatible` updates are wrapper-only; `migration-needed` changes require a
current human migration confirmation; `blocked` state remains read-only.
Records without a version are `legacy-unknown` and stay readable.

The v1-to-v2 plan migration in `v1-to-v2.yaml` is additive. It preserves plan
intent, status, context, runtime, dirty work, repository identity, local
bindings, and registered repositories. Rollback restores only wrapper-owned
files, including the shipped root CLAUDE.md adapter, and never overwrites
`repositories.local.yaml` or repository contents. Optional host evidence
fields are additive; unknown or contradictory host state remains read-only.

Identity-region insertion is additive. Upgrade inserts a missing bounded
identity region into legacy `context/WORKSPACE.md`, `context/PROJECT.md`, and
`context/INDEX.md` from current `workspace.yaml` metadata. It does not rewrite
authored Product Knowledge outside that region and does not change accepted
identity. Rollback removes the inserted region and still leaves accepted
identity and authored Product Knowledge unchanged. Effect identifiers remain
descriptive and never authorize a route or gate.

Canonical packet fixtures remain wrapper-owned: the context-receipt,
delegation, child-start, plan, and task schemas plus the plan, task, and PRD
templates and gate card are installed with the wrapper. A compatible receipt is retained
as runtime evidence, and a legacy child-start record is retained as well;
migration-needed or unknown evidence is classified for a current human decision and is never
rewritten merely to pass upgrade. Legacy-readable evidence never gains a
commit marker or current ownership through compatibility handling.
