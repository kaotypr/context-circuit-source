# Repository grounding — worktree hardening

## The point

The execution-environment workarounds that clutter hand-written writer prompts —
symlinked `node_modules`, `--config.verify-deps-before-run=false`, `--no-verify`
because a commit hook hangs — are **not the repository's knowledge and not the
user's to capture.** They are symptoms of Context Circuit's own worktree setup
being incomplete. So the fix is not to document them; it is to **prepare a
worktree where they are unnecessary.**

## What hardening does

During worktree preparation the runtime, deterministically and host-neutral:

1. **Detects the toolchain** from the worktree (e.g. a lockfile → its package
   manager; a module manifest → its ecosystem). A runtime constant maps known
   lockfiles/manifests to a provisioning strategy.
2. **Provisions dependencies** so the writer has a working environment without a
   per-worktree install — reusing the bound checkout's already-installed
   dependencies where the ecosystem allows, and configuring the tool so its own
   preflight does not abort on that arrangement.
3. **Handles commit hooks** for the isolated execution context so a commit does
   not hang or require a manual bypass — the runtime owns the commit step, so it
   decides how hooks run there.
4. **Reports the result** as the manifest's `environment` field (`ready` or
   `no-toolchain`), which the brief's Environment section renders.

## Before and after

| Hand-written prompt (today) | Hardened worktree (v0.6) |
| --- | --- |
| "symlink `node_modules` from the main checkout" | runtime provisions deps |
| "run every command with `--config.verify-deps-before-run=false`" | runtime configures the tool; flag unneeded |
| "lefthook hangs; run biome manually then `git commit --no-verify`" | runtime handles hooks for the exec context |
| writer follows a workaround list | writer runs the repo's normal commands |

## What hardening does NOT remove

Genuinely **repo-specific** operational detail — *how do I test / build / lint
this project* — is not a CC quirk. It lives in the repo's own `AGENTS.md` /
manifest scripts and reaches the writer through **discovery** (the grounding
directive), not through hardening. Hardening removes CC-induced friction; the
repo's own commands stay the repo's to state.

## Boundaries

- Deterministic and host-neutral (INV-RUNTIME-01): detection is a table lookup,
  provisioning and hook handling are mechanical. No intelligence, no network
  policy beyond what the ecosystem's provisioning already implies.
- The isolation model is unchanged: each execution still gets its own worktree and
  branch; hardening only makes that worktree *ready*, it does not share writers or
  weaken bounds.
- **Nothing is captured manually.** If a residual toolchain hint must persist for
  speed, it is a runtime-generated cache, never a user-authored file.

## Phasing

Hardening is the **heaviest runtime piece** of this scope (toolchain detection,
provisioning, hook handling across ecosystems). It may land as its **own
implementation phase** after discovery, the brief template, and the worker
contract — those deliver the primary value (the writer reads the repo's guidance)
even before every ecosystem is hardened. The profiles the manifests reveal across
repos tell you which quirks are CC-induced (harden) versus repo-inherent (stay in
the repo's own docs).
