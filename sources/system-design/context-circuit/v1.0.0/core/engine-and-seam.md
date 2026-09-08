# Engine verbs and the stable seam

All v1.0 changes are expressed as new or changed verbs *behind the existing CLI
seam* (`sh wrapper/runtime/engine.sh <action> <args>`, invoke-not-read). The seam
is the reason evolution is safe: skills that call stable verbs keep working.

## New verbs

| Verb | Purpose | Mechanism |
| --- | --- | --- |
| `intent-validate <dir>` | structural/schema check only: does `intent/<id>/` have the required files, does `contract.yaml` parse and match its schema — not a feasibility or risk check | M1 |
| `intent-allocate-id <slug>` | next `i<NNN>-slug` in the current member's intent band (active+archived; never wrap) | M1 |
| `member-roster-validate` | validate committed `members.yaml` (non-overlapping bands) | M1 |
| `member-identity-read` | read gitignored `member.local.yaml`; fail closed if missing | M1 |
| `member-band-resolve` | resolve local identity to roster bands; never accept a numeric range | M1 |
| `intent-approve <id>` | draft→approved; freeze `contract_digest`; the one upstream gate | M1 |
| `intent-archive <id>` | move `intent/<id>/` to `intent/archive/<id>/`; status-blind (INV-ARCHIVE-01/02) | M1 |
| `intent-restore <id>` | move an archived intent back to the active area | M1 |
| `candidate-digest <plan> <exec>` | deterministic digest over commit map + bases + `contract_digest` | M2 |
| `candidate-current <plan>` | report the current candidate id | M2 |
| `human-acceptance-record <candidate> <file>` | first-class acceptance bound to a candidate | M2 |
| `knowledge-debt <root>` | list delivered candidates with unresolved reconciliation | M4 |
| `knowledge-reconciled <candidate>` | clear the debt marker (accept or explicit defer) | M4 |

## Changed verbs

| Verb | Change |
| --- | --- |
| `plan-allocate-id <slug>` | next `NNNN-slug` in the current member's plan band (active+archived; never wrap); a stack reserves consecutive in-band ids |
| `plan-validate` | require `intent:` and a structural/schema check; there is no scope-containment verb — feasibility is a coordinator judgment run once on the tracer's findings (`intent-feasibility.md`), and scope-safety is settled at delivery (Gate 2) |
| `plan-approve` | narrowed: automatic once the intent is approved, not a human gate (INV-APPROVE-01 reworked) |
| `verifier-result-record` | bind the result to the current candidate (INV-CANDIDATE-01); keep the exact read-only tip-check |
| `completion-ready` | read the intent `tier`; enforce the verifier floor (Critical needs a candidate-bound `passed`); support inferred completion at lower tiers |
| `plan-complete` | at Critical, explicit; otherwise a projection of accepted + delivered; emit the reconciliation-debt marker |
| `context-impact-record` | extended to carry/clear the reconciliation-debt marker keyed to the candidate |
| `execution-begin` / run-stack loop | compute the candidate over the delivered tip set; spawn the verifier only at Standard/Critical |

## Preserved verbs (untouched)

`workspace-*`, `repository-*`, `worktree-prepare`, `base-prepare`, `lease-*`,
`plan-archive`/`plan-restore`, `plan-index-*`, `attempt-*`, `worker-commit-record`,
`worker-handoff-record`, `verifier-prepare`, `attempt-evidence-record`,
`repair-allowed`, `delivery-*`, `pair-*`, `recovery-inspect`, `lock-*`, and the two
helpers `cc_digest` + atomic write. The pairing verbs (`pair-begin`, `pair-inspect`,
`pair-close`, `pair-delivery-targets`) are unchanged; the *skill* reframes them as
the Explore tier with a promote step.

## Coupling map — what ports unchanged

From the skill/engine coupling read:

| Skill | Coupling | Under v1.0 |
| --- | --- | --- |
| `cc-publish` | none (host/MCP only) | **unchanged** (depends only on `cc_digest` + atomic write, preserved) |
| `cc-system-design` | none | **unchanged** — and now feeds the new intent object naturally |
| `cc-workspace` | workspace + member identity verbs | **additive** — `member-roster-validate`, `member-identity-read`, `member-band-resolve` |
| `cc-archive` | 2 simple verbs | **unchanged** (`plan-archive`, `plan-restore`) |
| `cc-pair` | 3 `pair-*` verbs | **reframed** as Explore tier; verbs unchanged, adds promote |
| `cc-execute` / `cc-run-stack` | deep (≈11 verbs + `.runtime/executions/`) | **additive changes** — record candidate, tier-gate the verifier |
| *new* `cc-intent` | — | authors `INTENT.md` + `contract.yaml` (the plain human decision) and takes the single upstream approval, `intent-approve` (freezes `contract_digest`); does not read the codebase and runs no validation step — there is none any more |
| *new* `cc-trace` | — | spawned on `intent-approve`; fans out one read-only tracer child per repository, in parallel; collects manifests; runs the feasibility check on the tracer's findings (surfacing an infeasible intent or a required change beyond a bound scope); hands the coordinator the grounded basis for planning |
| `cc-plan` | plan authoring | creates plan(s) from the trace manifest (not a blind read of `context/`); requires parent intent; sets `plan.intent`; knowledge-debt preflight |

Roughly two-thirds of the periphery ports unchanged; only the execution drivers
and `cc-plan` take real (additive) work, plus two new skills.

## The later typed-language port (optional, separate track)

Because the seam is stable, the *mechanics* (leases, base selection, verifier
enforcement, records) could later be re-implemented in a typed language behind the
identical verb surface and emit format — a **behavior-preserving port first**, then
policy — killing the two accidental-complexity hotspots identified in the read:
the hand-rolled YAML-in-awk readers and the string-emit/re-parse return convention.
This is justified only if testability and removing the YAML fragility are worth the
regression risk; it is **not** required to deliver any of M1–M4. v1.0's policy
changes and the mechanics port are independent tracks that share the seam.
