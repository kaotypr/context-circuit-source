# Nest Context Circuit system files under `.context-circuit`

Plan ID: 0003-nest-system-home
Intent: i017-nest-system-home
Status: draft

## Original request and coverage

- What you asked: Context Circuit's own system files live under
  `.context-circuit/`, so the workspace root is the project's, not the
  product's. Three product folders move there as they are: `wrapper/`,
  `agents/`, and `docs/`. Hosts still find the front door. `.agents/` stays
  at the workspace root. New and upgraded workspaces match. Project files
  stay where they are.
- Covered below: move the three folders and retarget the runtime
  (NEST-001), keep the host front door as thin pointers (NEST-002), ship
  the nested home in the template and in upgrades (NEST-003), then
  retarget leftover path names and prove the old root homes are gone
  (NEST-004).
- Unchanged: gates, tracer, independent check, one-owner-per-rule,
  invoke-not-read, host-blocked. No implicit deliver, publish, merge, or
  commit.

## Objective and desired behavior

In the source checkout and in an instantiated workspace, the shipped
wrapper, role files, and product docs live under `.context-circuit/`. The
runtime still runs from that nested home. Hosts still enter through root
`AGENTS.md` / `WORKFLOW.md` and `.agents/`. A new workspace is nested from
the start; an upgrade reaches the same layout without moving the user's
project files. No leftover product `wrapper/`, `agents/`, or `docs/` at
the workspace root.

## Constraints and non-goals

- Nested home is `.context-circuit/wrapper`, `.context-circuit/agents`,
  and `.context-circuit/docs`. Nest, do not flatten.
- Host-required discovery stays at the surfaces hosts already look at.
  `.agents/` stays at the workspace root.
- The blank template seed remains `template/` at the source root; its
  contents follow the nested layout.
- An upgrade may replace template-owned files and must preserve
  workspace-owned files. When the layout change alters the meaning of a
  record, report migration-needed and preserve the old state.
- Do not change what the product does (gates, tracer, verifier floor,
  self-verify).
- Do not move workspace-owned files into `.context-circuit`.
- Do not rename the wrapper or drop the wrapper/template split.
- Do not add a new host, router, or runtime scheduler.
- Preserve unrelated and dirty work. Run the semantic acceptance suite
  after meaningful phases.
- Do not deliver, publish, merge, or commit the implementation as part of
  this plan.

## Product Knowledge grounding

- host-adapters (context/domains/host-adapters/README.md) — hosts enter
  through root `AGENTS.md` / `WORKFLOW.md` and thin adapters. After the
  nest, those files stay at root as pointers into `.context-circuit`;
  they must not become a second copy of the product.
- source-release-and-upgrade (context/domains/source-release-and-upgrade/README.md)
  — assembler, shipped/never-ship sets, and upgrade preservation. The
  shipped layer's shape becomes nested; `template/` stays the source-root
  seed identity; workspace-owned files stay unmoved.
- repository-grounding (context/domains/repository-grounding/README.md) —
  worker-brief lives beside the runtime. The lookup must follow the
  nested wrapper or grounding and pairing fail.
- tracing (context/domains/tracing/README.md) — schema and role paths
  move with the product home; `intent/` stays at the workspace root.

Grounding summary: the product home is hardcoded across the engine's
worker-brief lookup, the release assembler, skills that invoke the
engine, every test `eng()` helper, and ~30 harness plots. Moving the
folders without retargeting those cites leaves the product unusable.
Host discovery and skill discovery stay at root; only the three product
folders nest.

## Repositories and source evidence

- context-circuit-source — owns the shipped layer, role files, product
  docs, assembler, skills, tests, and Product Knowledge. Evidence:
  `intent/i017-nest-system-home/trace/context-circuit-source.yaml` at HEAD
  `50cb84d`. Completeness proof found 718 old-home path references to
  retarget or prove KEEP. One repository; no second repository required.

## Tasks

1. NEST-001 (context-circuit-source, paths:
   `wrapper/`, `agents/`, `docs/`, `.context-circuit/`; depends on: none)
   — move `wrapper/`, `agents/`, and `docs/` under `.context-circuit/` as
   folders. Retarget the engine's worker-brief lookup, shipped/upgrade
   sets, owner map, and schema `owner:` fields. Add a product-home helper
   for engine-internal lookups and a workspace-validate probe that the
   nested wrapper exists. Acceptance NEST-AC-001, NEST-AC-002.
   Verification NEST-VT-001, NEST-VT-002.

2. NEST-002 (context-circuit-source, paths:
   `.context-circuit/wrapper/adapters/`, `AGENTS.md`, `WORKFLOW.md`,
   `CLAUDE.md`, `CURSOR.md`, `.cursor/rules/role-tiering-spawn.mdc`;
   depends on: NEST-001) — keep
   workspace-root host discovery as thin pointers into
   `.context-circuit`. `.agents/` stays at root. Source-checkout
   `AGENTS.md` stays maintainer-specific and is not the shipped pointer.
   Acceptance NEST-AC-003. Verification NEST-VT-003, NEST-VT-004.

3. NEST-003 (context-circuit-source, paths: `scripts/release-artifact.sh`,
   `scripts/release-manifest.txt`, `scripts/build-dist.sh`,
   `.context-circuit/wrapper/migrations/README.md`, `test/release/`,
   `agent-harness/test-template-runtime.sh`; depends on: NEST-001) —
   assemble and upgrade onto the nested home. New workspaces are nested
   from the start. An upgrade moves template-owned product folders under
   `.context-circuit/`, retargets root pointers, and preserves
   workspace-owned files. Acceptance NEST-AC-004, NEST-AC-005.
   Verification NEST-VT-005, NEST-VT-006.

4. NEST-004 (context-circuit-source, paths: `.agents/skills/`, `test/`,
   `agent-harness/`, `context/`, `README.md`; depends on: NEST-001,
   NEST-002, NEST-003) — retarget every remaining product-home cite in
   skills, tests, harness plots, and Product Knowledge. Prove no leftover
   root `wrapper/`, `agents/`, or `docs/`, and that workspace-owned trees
   stayed put. Acceptance NEST-AC-006, NEST-AC-007, NEST-AC-008.
   Verification NEST-VT-007, NEST-VT-008, NEST-VT-009.

## Acceptance criteria

- NEST-AC-001 — In the source checkout, `.context-circuit/wrapper`,
  `.context-circuit/agents`, and `.context-circuit/docs` exist, and there
  is no leftover product `wrapper/`, `agents/`, or `docs/` at the
  workspace root.
- NEST-AC-002 — The runtime finds the nested home: worker-brief assemble
  and engine invoke use `.context-circuit/wrapper/...`. A workspace
  missing that nested wrapper fails validate loudly.
- NEST-AC-003 — Root `AGENTS.md` / `WORKFLOW.md` / `CLAUDE.md` /
  `CURSOR.md` exist and point into `.context-circuit`. They do not carry
  a second copy of the product. `.agents/` remains at the workspace root.
- NEST-AC-004 — The assembled template artifact uses
  `.context-circuit/{wrapper,agents,docs}`, not root `wrapper/`. Root
  adapters and `.agents/` still ship.
- NEST-AC-005 — Upgrade guidance moves template-owned product folders
  under `.context-circuit/`, preserves workspace-owned files, and reports
  migration-needed when a record's meaning still depends on a top-level
  `wrapper/` path.
- NEST-AC-006 — Skills, tests, harness plots, and Product Knowledge that
  named the old homes now name `.context-circuit/...`. The leftover-path
  completeness check has zero actionable matches.
- NEST-AC-007 — Workspace-owned trees (`context/`, `plans/`, `intent/`,
  `template/`, identity, runtime evidence) stay at the workspace root and
  are not under `.context-circuit/`.
- NEST-AC-008 — Lifecycle, one-owner-per-rule, invoke-not-read, and
  host-blocked fail-closed are unchanged.

These prove the intent outcomes AC-NESTED-HOME, AC-RUNTIME-FINDS-HOME,
AC-HOST-FRONT-DOOR, AC-TEMPLATE-AND-UPGRADE, AC-NO-LEFTOVER-ROOT, and
AC-WORKSPACE-OWNED-STAY.

## Verification

- NEST-VT-001 — nested product folders exist; root product trees are
  absent.
- NEST-VT-002 — engine worker-brief lookup and workspace-validate probe
  the nested wrapper path, not `wrapper/runtime/` at the workspace root.
- NEST-VT-003 — root host surfaces exist; `AGENTS.md` or `WORKFLOW.md`
  points into `.context-circuit`; `.agents/` remains at root.
- NEST-VT-004 — shipped adapters under `.context-circuit/wrapper/adapters`
  still own the product text; root files are pointers, not a second copy.
- NEST-VT-005 — `release-artifact.sh` stages `.context-circuit/{wrapper,agents,docs}`
  and copies adapters from the nested adapters directory.
- NEST-VT-006 — `sh test/release/test-release.sh` and
  `sh agent-harness/test-template-runtime.sh` pass with the nested layout.
- NEST-VT-007 — leftover old-home grep (same command as the completeness
  proof, excluding archives, sources, and this intent) has zero matches.
- NEST-VT-008 — workspace-owned trees remain at root and are absent under
  `.context-circuit/`.
- NEST-VT-009 — `sh test/acceptance.sh` passes.

## Assumptions, open questions, risks

Assumptions (plan-level, from the look at the code):

- Engine-internal lookups use a single product-home helper
  (`.context-circuit/wrapper` from the workspace root). Skills, tests, and
  harness invoke `.context-circuit/wrapper/runtime/engine.sh` explicitly —
  the same pattern as today's `sh wrapper/runtime/engine.sh`, retargeted.
  No second runtime and no path-rewriting layer.
- No new upgrade engine verb. Document the layout move in
  `wrapper/migrations/README.md` and `upgrade_boundary`. The assembler
  emits the nested home for new workspaces. An existing workspace is
  upgraded by replacing template-owned product trees and root pointers
  while preserving workspace-owned files. If a record's meaning still
  depends on top-level `wrapper/`, report migration-needed and preserve
  the old state. `repository-binding-migrate` stays the precedent for
  workspace-owned file migration; this layout move is template-owned.
- Source-checkout `AGENTS.md` stays maintainer-specific (source safety),
  retargeted to nested paths. Instantiated workspaces and the assembled
  artifact get the thin product pointer copied from
  `.context-circuit/wrapper/adapters/AGENTS.md`. Do not fold maintainer
  rules into the shipped pointer, and do not duplicate product policy at
  the source root.
- `cc_workspace_validate` gains a nested-home probe so a workspace
  missing `.context-circuit/wrapper` fails closed. Other verbs keep their
  current shape.
- `template/` stays at the source root as seed identity. Assembly still
  copies product trees from the source checkout; the seed is not required
  to contain a nested product home of its own.
- Source-only `docs/release.md` moves with `docs/` in this checkout and
  remains excluded from the release-manifest ship set.
- `.agents/` stays at the workspace root (already answered on the
  intent). `docs/` moves (already answered).

Open questions: none that change the approved decision.

Risks:

- Moving the folders without retargeting `test/lib/assert.sh` leaves
  every suite sourcing a missing engine.
- Skills that still say `sh wrapper/runtime/engine.sh` fail closed after
  the move — coordinators cannot allocate intents, execute, pair, or
  trace.
- `release-artifact.sh` still tar-staging top-level `wrapper agents docs`
  ships the old layout even if the source checkout is nested.
- Updating only the owner map, without schema `owner:` fields and skill
  citations, leaves one-owner-per-rule inconsistent.
- Making root `AGENTS.md` a full product adapter instead of a thin
  pointer duplicates policy outside the nested home.
- Accidentally moving `.agents/`, `context/`, `plans/`, `intent/`, or
  `.runtime/` into `.context-circuit/`.
- Folding maintainer source `AGENTS.md` into the shipped pointer leaks
  maintainer-only rules into instantiated workspaces, or drops them from
  source.
- ~30 harness scenarios forbid `wrapper/runtime/engine.sh` as the
  invoke-not-read boundary; those paths must update together or grading
  rejects correct coordinator behavior.
- Upgrade docs that only describe a new workspace leave existing
  workspaces with root product folders.

## Expected commits and delivery notes

- One commit on context-circuit-source covering the nested home, host
  pointers, assembler and upgrade notes, skills, tests, harness, and
  Product Knowledge.
- Delivery (pull request, merge) is a separate explicit action. This plan
  does not deliver.

## Expected Product Knowledge impact

- Reassess host-adapters, source-release-and-upgrade, repository-grounding,
  tracing, and architecture at completion: product home is
  `.context-circuit/{wrapper,agents,docs}`; root host files are pointers;
  `.agents/` stays at root; `template/` stays the source-root seed;
  upgrade preserves workspace-owned files. Domain implementation
  references that named the old homes are retargeted in NEST-004.
