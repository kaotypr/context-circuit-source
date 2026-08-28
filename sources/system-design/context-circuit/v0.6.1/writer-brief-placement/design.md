# writer-brief-placement — design

## Capability

Keep the shipped `context-circuit-template` workspace root clean of internal
runtime machinery by relocating the `writer-brief.md` template out of the
user-facing root into an unobtrusive, machinery-only location, without changing
what the template is or how repository-grounding uses it.

## Problem

When a user instantiates `context-circuit-template`, the workspace root contains
`writer-brief.md` alongside the recognizable entry docs (`README.md`,
`AGENTS.md`, `CLAUDE.md`, `WORKFLOW.md`). Those four are legitimately
user-facing. `writer-brief.md` is not: it is an **internal runtime template**
(a fixed brief skeleton the runtime fills by deterministic slot substitution
during execution) that a user should never open or edit. Sitting undifferentiated
at the bare root, it reads as noise and invites the reasonable question "why is
this here?" — exactly the feedback that prompted this scope.

It is **not a stray file**. It is placed at root deliberately by release
assembly, and it is load-bearing:

- `scripts/release-artifact.sh` copies `wrapper/adapters/writer-brief.md` to the
  staged tree **root**, promoting it alongside the other adapter entry docs.
- `scripts/release-manifest.txt` marks `writer-brief.md` **required** at root.
- The shipped artifact **excludes `wrapper/adapters/`**, so at runtime the engine
  lookup in `cc_writer_brief_assemble` resolves the **root** copy — it is the
  only home the template has in the shipped template. (The
  `wrapper/adapters/writer-brief.md` fallback exists only for this source
  checkout, where `wrapper/adapters/` is present.)

So the placement is intentional; the wart is that a *runtime-only* artifact was
promoted to the same visible tier as *user-facing* entry docs.

## Principle

The workspace root is user-facing surface. A file promoted there must be
something a user is meant to read or edit. Runtime-only machinery belongs with
the other shipped machinery (the template already ships agent machinery under
`.agents/` and the runtime under `wrapper/runtime/`), not at the bare root.

The `writer-brief.md` **source of truth is unchanged**: it remains
`wrapper/adapters/writer-brief.md`, the owner named in
`wrapper/contracts/invariants.yaml` (INV-GROUND-01/03). Only the *promoted,
shipped runtime copy's location* moves.

## Fixed decisions

1. **Relocate the promoted copy** out of the user-facing root into an
   unobtrusive machinery location that already ships. The concrete candidate is
   `.agents/writer-brief.md` — grouping it with the shipped `.agents/skills/`
   machinery, out of the user's default view. (Final path is an implementation
   detail; `.runtime/` is **not** eligible — it is excluded from the artifact as
   runtime evidence.)
2. **The engine lookup follows the file, with the source fallback preserved.**
   `cc_writer_brief_assemble` prefers the new shipped location, then falls back
   to `wrapper/adapters/writer-brief.md` for the source checkout. No behavior
   change beyond the path.
3. **The source of truth stays** `wrapper/adapters/writer-brief.md`; the owner
   map and INV-GROUND rules are untouched.
4. **No user-visible behavior change** to repository-grounding: the assembled
   brief, its preflight, and grounding execution are identical; only where the
   unfilled template lives changes.

## Shape of the change (surfaces that move)

Each is owned elsewhere and changed through its owner's normal action:

- **Assembly** (`scripts/release-artifact.sh`): copy `wrapper/adapters/writer-brief.md`
  to the new location in the staged tree instead of the root.
- **Manifest** (`scripts/release-manifest.txt`): change `required writer-brief.md`
  to the new path (still required; the leak-assertion on `wrapper/adapters`
  stays).
- **Engine lookup** (`wrapper/runtime/engine.sh`, `cc_writer_brief_assemble`):
  point the primary lookup at the new location; keep the `wrapper/adapters/`
  fallback for source checkouts.
- **Prose references** that name the template location (e.g.
  `agents/coordinator.md`, the `cc-execute` / `cc-run-stack` skills): reword any
  that imply a root path so they match the new location. The
  `writer-brief-assemble` command interface is unchanged.

No core contract bump. The invariants owner map value and INV-GROUND semantics
are unchanged.

## Constraints and edge cases

- The source-checkout path (this repository) has no promoted root copy; the
  engine's `wrapper/adapters/` fallback must keep working there unchanged.
- The template-harness live scenarios (`11-repo-grounding`, the run-stack cases)
  assert grounding *behavior* (assembled-brief content, manifest recorded), not
  the template's on-disk location, so they are unaffected — but the assertion
  set should be checked so none inadvertently pins the root path.

## How it feeds the rest — unchanged

This is source. It reaches Product Knowledge and a plan through the normal path:
the coordinator gathers context from this named design source, proposes the
relocation as context/contract changes through the existing proposal path, a
human accepts, and a plan grounds in the result. No separate design-acceptance
gate.
