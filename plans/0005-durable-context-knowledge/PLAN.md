# Make context hold only durable product knowledge

Plan ID: 0005-durable-context-knowledge
Intent: i019-durable-context-knowledge
Status: draft

## Original request and coverage

- What you asked: when an agent creates or updates Product Knowledge, those
  files hold only durable knowledge about the product. They must not name a
  plan, an intent file, or a sources file. Existing live context is rewritten
  now. No live context file names a sources path. `DECISIONS.md` records the
  knowledge change, not files.
- Covered below: put the rule on the knowledge owners (DCK-001), teach every
  write path (DCK-002), rewrite live context including `DECISIONS.md`
  (DCK-003), align the blank seed (DCK-004), and prove it with checks
  (DCK-005).
- Still allowed: explaining the product concepts plan, intent, and sources.
  When knowledge is written (gathering context, or reconciling after
  mark-done) does not change.

## Objective and desired behavior

Live context stays useful after plans, intents, and sources are archived or
deleted. An agent that writes those files records what is true about the
product, not the ephemeral artifact the change came from.

## Constraints and non-goals

- Put the rule on the knowledge contract, not only in a skill.
- Linking between durable context units stays allowed.
- Invariant IDs and shipped contract or adapter paths may still be named.
- Do not change when knowledge is written.
- Do not rewrite sources files, plan files, or intent files themselves.
- Do not add a third human gate for knowledge updates.
- Do not deliver, publish, merge, or commit the implementation as part of
  this plan.

## Product Knowledge grounding

- completion (context/domains/completion/README.md) — mark-done still
  reconciles live context when knowledge was affected; that write must
  follow the durable-only rule.
- conventions (context/CONVENTIONS.md) — today it tells agents to record
  source-file provenance in context; that conflicts with the rule.
- sources (context/SOURCES.md) — today it records sources paths in live
  context; it must become a concept page.
- decisions (context/DECISIONS.md) — today it cites proposals, file edits,
  and particular artifacts; it must record the knowledge change only.
- index (context/INDEX.md) — retrieval catalog; today it points at
  `context/sources.yaml` as provenance rules.

Grounding summary: the rule is missing from the owners, the live pages
already name ephemeral files, and the blank seed teaches the old habit. All
three have to change together or the next write undoes the rewrite.

## Repositories and source evidence

- context-circuit-source — owns the knowledge contract, writing paths, live
  `context/`, and `template/context/`. Evidence: the look at HEAD `4245a80`.
  Twenty-six live context files, sixty-eight path-citation lines. No second
  repository.

## Tasks

1. DCK-001 (context-circuit-source, paths:
   `.context-circuit/wrapper/contracts/invariants.yaml`,
   `.context-circuit/docs/product-knowledge.md`,
   `.context-circuit/wrapper/contracts/schemas/context-index.yaml`,
   `context/CONVENTIONS.md`; depends on: none) — add the durable-only rule
   to the knowledge owners. Acceptance DCK-AC-001, DCK-AC-002. Verification
   DCK-VT-001, DCK-VT-002.

2. DCK-002 (context-circuit-source, paths:
   `.context-circuit/wrapper/adapters/WORKFLOW.md`,
   `.context-circuit/agents/coordinator.md`,
   `.agents/skills/cc-complete/SKILL.md`,
   `.context-circuit/docs/templates/domain-context.md`; depends on:
   DCK-001) — gather and mark-done writes follow the rule. Acceptance
   DCK-AC-003. Verification DCK-VT-003.

3. DCK-003 (context-circuit-source, paths: `context/`; depends on:
   DCK-001) — rewrite live context, including `DECISIONS.md` and a
   comment-free empty `context/sources.yaml`. Acceptance DCK-AC-004,
   DCK-AC-005, DCK-AC-006. Verification DCK-VT-004, DCK-VT-005, DCK-VT-006.

4. DCK-004 (context-circuit-source, paths: `template/context/SOURCES.md`,
   `template/context/sources.yaml`, `template/context/DECISIONS.md`,
   `template/context/INDEX.md`, `workspace.yaml`,
   `scripts/release-manifest.txt`; depends on: DCK-001, DCK-002) — align
   the blank seed. Acceptance DCK-AC-007. Verification DCK-VT-007.

5. DCK-005 (context-circuit-source, paths: `test/contracts/test-contracts.sh`,
   `test/`; depends on: DCK-002, DCK-003, DCK-004) — prove it; stop tests
   from requiring ephemeral paths in live context. Acceptance DCK-AC-008.
   Verification DCK-VT-008.

## Acceptance criteria

- DCK-AC-001 — The knowledge invariants and product-knowledge doc state
  that live context is durable-only and never names a plan, an intent file,
  or a sources file.
- DCK-AC-002 — CONVENTIONS.md no longer requires recording a sources path
  in context. Index provenance does not mean a sources path.
- DCK-AC-003 — Every path that generates or updates context files states
  the durable-only rule.
- DCK-AC-004 — No live context file names a plan, an intent file, or a
  sources file, including `context/sources.yaml`.
- DCK-AC-005 — DECISIONS.md records durable product-knowledge decisions,
  not file edits and not the artifacts the change came from.
- DCK-AC-006 — Product Knowledge still explains the product concepts plan,
  intent, and sources.
- DCK-AC-007 — The blank context seed does not name those artifacts, and it
  tells agents that DECISIONS records the knowledge change, not files.
- DCK-AC-008 — Checks fail if live context or the blank seed names a banned
  path, or if DECISIONS.md cites proposals instead of the knowledge change.

These prove the intent outcomes AC-NO-EPHEMERAL-CITE, AC-REWRITE-EXISTING,
AC-DECISIONS-KNOWLEDGE, AC-CONCEPTS-REMAIN, and AC-ALL-WRITES.

## Verification

- DCK-VT-001 — knowledge invariants and product-knowledge.md state durable.
- DCK-VT-002 — CONVENTIONS.md no longer says “record provenance”.
- DCK-VT-003 — coordinator, WORKFLOW.md, and cc-complete state durable;
  the domain template no longer keeps raw text in sources/.
- DCK-VT-004 — no live context file still contains a banned path.
- DCK-VT-005 — DECISIONS.md has no “Accepted from proposal”.
- DCK-VT-006 — intent, plan-organization, and SOURCES concept pages remain.
- DCK-VT-007 — the blank seed has no banned path and DECISIONS.md states
  durable.
- DCK-VT-008 — `sh test/acceptance.sh` passes.

## Assumptions, open questions, risks

Assumptions (plan-level, from the look at the code):

- Generic product-layout notation (`plans/<id>/`, `intent/<id>/INTENT.md`)
  stays as concept. Particular artifact ids and every sources path go.
- `context/sources.yaml` stays as a comment-free empty catalog so workspace
  identity and the release seed still have a file to point at. It does not
  record sources paths. INDEX and SOURCES.md stop treating it as a
  provenance-path catalog.
- Index provenance, if the field remains, is not a sources path.
- Template seed is rewritten; it is shipped product surface, not only this
  checkout’s live pages.
- Historical DECISIONS.md entries keep what became true and drop proposal
  citations, file-edit lists, and particular artifacts.

Open questions: none that change the approved decision. Existing pages are
rewritten now. No live context file names a sources path.

Risks:

- Changing only the complete skill leaves gather-context still recording
  sources paths.
- Rewriting live context without changing the owners lets the next
  mark-done reintroduce citations.
- A blunt `sources/` search can also hit legitimate concept prose; checks
  must ban paths and particular ids, not the words plan, intent, or sources.
- Stripping DECISIONS.md file lists can drop the durable decision if the
  rewrite is careless.
- `test/contracts/test-contracts.sh` currently requires ephemeral paths in
  a domain page; leaving that assertion fails a correct rewrite.

## Expected commits and delivery notes

- One commit on context-circuit-source covering the contract, writing
  paths, live context, blank seed, and tests.
- Delivery (pull request, merge) is a separate explicit action. This plan
  does not deliver.

## Expected Product Knowledge impact

- Reassess conventions, sources, decisions, index, and completion at
  completion: live context is durable-only; DECISIONS.md records the
  knowledge change; gather and mark-done writes follow that rule.
