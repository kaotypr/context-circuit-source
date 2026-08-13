# Product Knowledge delivery record

Durable, append-only record of the Product Knowledge work delivered on
`codex/product-knowledge`. This is maintainer evidence for the Context Circuit
source repository; it is intentionally excluded from the shipped wrapper template.

## Outcome

Context Circuit now carries an optional, incremental Product Knowledge capability:
a lean business view (product map, roles, domains, and workflow pages) that gives
humans and agents task-relevant business context without loading the whole
application. It is absent-safe, contract-validated, and threaded through the whole
delivery journey — planning, execution, evidence, effective synchronization, and
onboarding — while preserving every existing safety property.

## Delivered work items

Implemented in delivery order as direct maintainer commits under the approved plan
`context/plans/product-knowledge` (`approved_by: maintainer`):

1. **PKNOW-001 — Contracts, layout, templates, validation.** Product-knowledge
   project/role/workflow/domain page contracts, a deterministic tree validator
   (required sections, resolvable references, non-overlapping authority), neutral
   templates, and small and large fixtures proving bounded selection.
2. **PKNOW-010 — Optional initialization baseline.** A bootstrap request may create
   a minimal reviewed baseline that preserves explicit unknowns and validates;
   omitting it keeps an existing wrapper valid.
3. **PKNOW-020 — Read-only discovery.** Source-backed candidate, contradiction, and
   gap classification against the canonical tree, credential-rejecting, persisted
   only as ignored runtime evidence.
4. **PKNOW-030 — Plan references and impact.** Plans declare referenced pages, a
   product impact, and a proposed change as approved material; additive and
   backward-compatible.
5. **PKNOW-040 — Bounded task context package.** An immutable, digest-pinned
   package of only the referenced pages carried to worker and verifier through the
   shared task brief, for planned and planless tasks.
6. **PKNOW-050 — Impact through evidence and closeout.** Worker and verifier report
   observed impact; closeout records it and withholds synchronization for human
   review when impact is unexpectedly broader or contradictory.
7. **PKNOW-060 — Effective synchronization.** Canonical current-behavior pages
   change only after a dedicated, workspace-configured confirming role declares the
   behavior effective; writes are bounded to referenced pages and roll back if the
   tree becomes invalid.
8. **PKNOW-070 — Onboarding packs.** Reproducible, revision-stamped, source-listed
   generated views that are explicitly not a source of truth.
9. **PKNOW-080 — Adapters, templates, guides.** Thin Codex and Claude adapters
   inherit the canonical behavior; the shipped template carries the Product
   Knowledge guide and every advertised contract.
10. **PKNOW-090 — End-to-end proof.** A single journey test composes the whole
    lifecycle and asserts bounded context, backward-compatible absence, read-only
    discovery, plan impact, blocked silent synchronization, effective-only updates,
    reproducible packs, and host-adapter parity.

## Enabling decision

The neutral source checkout registers no product repository. Per the human
decision, the Context Circuit framework source registers itself as the
`context-circuit` repository (`path: .`, submodule mode, `repository-worker`
agent), and this feature was implemented in-place as maintainer commits — matching
the workflow-continuity precedent — rather than run-task worktree fan-out. The
template build neutralizes `workspace.yaml` so the shipped template still ships
`repositories: {}`.

## Verification

- `node --import tsx --test --test-concurrency=1 test/*.test.ts` — full suite green.
- `npm run typecheck`, `npm run validate -- --check-paths --check-documents`,
  `npm run build`, and `git diff --check` — clean.
- Packaging and release-artifact tests confirm a version- and digest-matched
  archive that ships every Product Knowledge contract, template, and guide while
  keeping the template neutral.

## Remaining risks and follow-up

- Effective synchronization is a library gate validated against the workspace's own
  configured confirming role; wiring it into the full reviewed-worktree
  `sync-context` handoff (branch, review, human merge) remains a natural extension.
- The neutral template currently ships the maintainer plan under
  `context/plans/product-knowledge/`. Stripping in-progress maintainer plans from
  the shipped template (as `PLAN.md` already is) is an optional follow-up.
