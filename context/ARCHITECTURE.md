# Architecture

Context Circuit has four layers: a compact entry spine, one two-stage router,
human-facing plan artifacts, and private resumable runtime evidence. The source
repository separates `wrapper/` (shipped immutable contracts) from `template/`
(blank mutable seed). An instantiated workspace keeps wrapper, identity,
context, sources, plans, runtime, and registered repositories distinct.

The router and invariant catalogs are canonical. Skills are thin natural-
language adapters. Runtime records are filesystem evidence, not a database or
scheduler. Writers are isolated by atomic plan leases and worktrees; a separate
verifier is read-only for implementation state.
