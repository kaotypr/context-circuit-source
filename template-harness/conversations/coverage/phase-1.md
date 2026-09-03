# Phase 1 coverage cross-check — orientation & context foundation

How the six phase-1 plots cover the existing `scenarios/` cases in this phase, and
the harness capability gaps they surface. Descriptive only: every rule cited is
owned by `wrapper/contracts/invariants.yaml` or `test/acceptance/criteria-map.yaml`.

## Existing cases → plots

| existing case | plot | relationship |
| --- | --- | --- |
| `01-new-project-simple-idea` | `orient-new-project` | **evolves** — same empty seed and `orient` budget; expectations now sourced from the plot |
| `02-connect-existing-repo` | `connect-existing-repo` | **evolves** — same `my-notes` seed (master default / develop working branch) and budgets; expectations now sourced from the plot |

No existing case maps to `clone-or-init-new-repo`, `query-product-knowledge`, or
`accept-or-defer-context-proposal` — these are the three real coverage gaps this
phase fills. `build-product-knowledge` already has its worked plot + case.

## AC / invariant mapping drift (surfaced, not hidden)

The library is meant to catch exactly this kind of drift between a case's declared
coverage and what its dialogue actually demonstrates.

- **`01` historically tags `[AC-01, AC-04, AC-05]`; the plot narrows to `[AC-01]`.**
  Case 01 is pure orientation: nothing is planned or executed. AC-04 (a plan
  preserves request detail) and AC-05 (a plan outside its scope envelope is
  re-gated) require a plan to exist, so this dialogue does not demonstrate them —
  they are demonstrated by the intent/plan plots (`refuse-before-gate1` (04),
  `scope-envelope-regate` (16), `tier-fails-upward-refuse-explore` (20)). The
  plot keeps `AC-01` (register zero repositories — here, zero) plus the fail-closed
  invariants the dialogue really exercises (`INV-REPO-04`, `INV-PLAN-04`,
  `INV-EXEC-01`, `INV-SEC-02`). **Recommendation:** when case 01 is regenerated,
  its AC list narrows to `[AC-01]`; confirm no acceptance-coverage report depends
  on 01 carrying AC-04/AC-05 (they remain covered by 04/16/20).
- **`02` mapping is unchanged.** `[AC-01, AC-26, AC-27]` + `INV-REPO-01/02/04`,
  `INV-SEC-01` all fit connect-in-place on a user-selected branch; the plot keeps
  them verbatim.

## Harness capability gaps surfaced this phase

These are gaps in the *generated case* layer (the executable harness), not in the
plots. Each is flagged inline in the affected `generated/*.case.yaml`. None blocks
authoring; they block eventual regeneration/execution and are listed here for the
generator/harness workstream.

1. **id-agnostic / path-shape post-conditions** (same class already noted in
   `plots/standard-feature-whole-flow.md`). Fresh, non-id-seeded cases need
   assertions the id-keyed seeded cases don't:
   - `repository_under_project_area: true` — a new repo landed under the project
     area (`repositories/`), the positive inverse of case 02's `no_repository_clone`
     (`clone-or-init-new-repo`). Owned rule: `INV-REPO-03`.
   - `initial_commit_present: true` — a freshly initialized repo has the base
     commit execution later requires (`clone-or-init-new-repo`).
   - `no_remote_created: true` — init stayed local; no remote/push
     (`clone-or-init-new-repo`). Owned rule: `INV-DELIVER-01`.
   - `knowledge_units_accepted: N` / `knowledge_units_pending: ">=1"` — per-unit
     accept/defer outcome without fixed unit ids (`accept-or-defer-context-proposal`).
     Owned rule: `INV-KNOWLEDGE-02`.

2. **Seed overlays the current runner does not provide.** The runner seeds
   `repositories` and `sources`; two phase-1 plots need more:
   - `setup.knowledge` — pre-accepted Product Knowledge to query against
     (`query-product-knowledge`). Needed so retrieval-first recall
     (`INV-KNOWLEDGE-01`, AC-02) has something to recall.
   - `setup.proposals` — a pending proposal awaiting a decision
     (`accept-or-defer-context-proposal`). Needed so the consent gate
     (`INV-KNOWLEDGE-02`) has a real pending unit to accept/defer.

3. **Access-policy `forbidden` on a seeded repo's own files.**
   `query-product-knowledge` forbids reading `notes/**` during the `query` phase to
   enforce "recall recorded knowledge, don't re-derive from code." The file-access
   trace already supports `forbidden` globs, so this is expected to work as-is —
   noted here only because it forbids files that *exist* in the seed (a stronger
   assertion than the usual "don't touch archive/engine/sources").
