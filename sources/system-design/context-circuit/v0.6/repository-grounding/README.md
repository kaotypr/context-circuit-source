# Repository grounding — design

The v0.6 scope that grounds the **writer** in the target repository's own agent
guidance and hands it a ready-to-work worktree — so repository knowledge is
discovered and delivered, never hand-authored into each writer prompt. v0.6 is a
**delta on v0.5** — see [../README.md](../README.md) for the version index.

## Reading order

1. [design.md](./design.md) — the overview: the problem, the three kinds of
   knowledge and who owns each, principles, and fixed decisions.
2. [discovery-and-manifest.md](./discovery-and-manifest.md) — the deterministic
   runtime discovery step and the grounding manifest it emits.
3. [writer-brief-template.md](./writer-brief-template.md) — the brief template
   (contract), the grounding directive, precedence, and deliver-vs-author.
4. [worktree-hardening.md](./worktree-hardening.md) — preparing a ready worktree
   so execution-environment workarounds disappear.
5. [schema-and-contracts.md](./schema-and-contracts.md) — the contract deltas:
   engine functions, the brief template, invariants, manifest, handoff field.
6. [examples.md](./examples.md) — worked briefs (rich / doc-less / empty repo)
   and the acceptance criteria.

## Authority

This scope adds no owner and duplicates no rule; it specifies intent and points
to the canonical owners under `wrapper/`. Contract deltas are summarized in
[schema-and-contracts.md](./schema-and-contracts.md).
