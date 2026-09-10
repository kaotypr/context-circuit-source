# Plan 0041 — Tell the truth about a ready worktree, and prove it end to end

**Intent:** i027-runnable-worktrees
**Repository:** context-circuit-source
**Tier:** Standard
**Status:** draft
**Depends on:** [0039-worktree-ignored-overlay](../0039-worktree-ignored-overlay/PLAN.md),
[0040-worktree-toolchain-provision](../0040-worktree-toolchain-provision/PLAN.md)

## Objective

Once the runtime really prepares a worktree, every place that *describes* a
prepared worktree has to say so. The manifest schema stops describing
`environment: ready` as a detection. The shipped worker brief says what the
runtime actually guarantees. The three owning Product Knowledge pages describe
create-time overlay, lockfile correction, and fail-closed setup instead of
"full dependency provisioning is a later phase". One dated decision and an
index pointer record it. And coverage closes the loop: the shipped-template
harness proves a *released* workspace behaves the same way.

## Grounding (HEAD 7d72f468)

| Surface | Current state |
|---------|---------------|
| `grounding-manifest.yaml` | `environment` is `[ready, no-toolchain]`, described as "a toolchain was detected and prepared". The enum's description is the only statement of what preparation means, and it mentions neither overlay nor lockfile correction |
| `worker-brief.md` (line 16) | Carries the `@@ENVIRONMENT@@` slot, filled by `cc_worker_brief_assemble` (engine.sh 1788-1791) with either "Ready: dependencies provisioned, commit hooks handled…" — today a forward-looking claim — or the no-toolchain sentence |
| `context/domains/repository-grounding/README.md` | The owning page. Behavior section says hardening reports `ready \| no-toolchain` and that "full dependency provisioning is a later phase". Its greenfield constraint stays true |
| `context/domains/plan-execution/README.md` (40, 60, 65, 84, 104, 113, 120) | Describes isolated worktrees and branches, nothing about what is *in* them beyond tracked files |
| `context/domains/direct-collaboration/README.md` (75-91) | Describes the Explore worktree, its cleanup, and that closing preserves it — all still true; says nothing about preparation or fail-closed setup |
| `context/DECISIONS.md` | Dated append-only list: `## YYYY-MM-DD — title`, then Decision / Rationale / Consequence |
| `context/INDEX.md` | Short paragraphs pointing at dated decisions and their domain pages |
| `agent-harness/test-template-runtime.sh` (137-139, 242-245) | Drives a *staged shipped template* through a real execution with worktrees at `.runtime/worktrees/<plan>/api` — where released behavior gets proved |
| INV-GROUND-03 | Already says workarounds are eliminated by hardening and never documented by hand. **No invariant rewrite is needed**; the risk is the brief acquiring a workaround sentence |

## Decisions

- **Separate plan** because it is only verifiable once both behaviors exist, and
  its failure surface is documentation drift and shipped-template behavior, not
  runtime correctness. Folding it into 0040 means either publishing a knowledge
  page ahead of the code, or holding a runtime change hostage to a knowledge
  review.
- **Smallest honest change to the brief.** Its ready sentence already claims
  dependencies are provisioned; after 0040 that is true. Re-read it, adjust only
  what 0040 made inaccurate, and add nothing that tells a worker how to cope
  with missing dependencies.
- **No new invariant, no new domain page.** One dated decision plus one index
  paragraph, in the house pattern of the 2026-09-08 and 2026-09-10 entries.
- **Prove the released artifact.** A user runs the shipped template, not this
  checkout, so the harness carries the create-time assertion.

## Tasks

| id | title | depends on |
|----|-------|------------|
| HR-001 | Manifest schema + worker brief describe real preparation | — |
| HR-002 | Update the three owning Product Knowledge pages in place | HR-001 |
| HR-003 | Record the dated decision; point the index at it | HR-002 |
| HR-004 | Prove a released workspace creates runnable worktrees | HR-003 |

Full changes, acceptance, and runnable done-checks are in `plan.yaml`.

## Risks

- **Second owner creep.** The easiest mistake here is restating the enum or an
  invariant on a domain page. HR-002 forbids it explicitly and HR-VT-004/005
  check for the retired wording.
- **Pinned strings.** `test/contracts/test-contracts.sh` pins exact text across
  shipped surfaces; changing the brief or schema wording can break assertions
  that look unrelated. HR-004 reconciles the whole suite rather than one file.
- **Documenting `.env` cloning.** The decision entry names it plainly. That is
  intentional — a reader must know secrets reach every worktree — but the entry
  must not turn into an operational guide for handling them.
- **Harness drift.** `test-template-runtime.sh` stages a shipped tree; if the
  release manifest ever stops carrying a file the new assertion needs, the
  harness fails for a reason that looks like this plan's fault.

## Verification

`sh agent-harness/test-template-runtime.sh`,
`sh test/grounding/test-grounding.sh`, `sh test/pairing/test-pairing.sh`,
`sh test/contracts/test-contracts.sh`, `sh test/knowledge/test-durable-context.sh`,
and finally `sh test/acceptance.sh`.
