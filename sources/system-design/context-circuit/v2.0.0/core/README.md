# Context Circuit v2.0 — mechanism in Go, judgment in the agent

**Status: source material only.** This is design material, not a plan, not
Product Knowledge, not a release, and not permission to change anything. It
describes the v2.0.0 core as built against the 2.0.0-rc.1 candidate and argues
why it is shaped this way. Where it disagrees with `product/AGENTS.md.in`,
`product/docs/`, or `internal/`, those own the answer and this document is stale.

The **core** scope of the Context Circuit v2.0.0 design. For the version index,
see [../README.md](../README.md).

## What this is

v2 is a **rewrite**, not an evolution. v1's `wrapper/runtime/engine.sh`, its
invariant catalog, its contract schemas, and its evidence records are gone
rather than ported. In their place: a ~2,500-line Go executable that does
bookkeeping and Git, a single shared instruction file that does everything else,
and two skills — one to install the CLI, one to dispatch subagents.

This scope covers the whole coordination model, because in v2 the model is
small enough to hold in one scope and its pieces only make sense together.

## The thesis in one sentence

> A system's own records prove only that the system wrote them — so v2 keeps the
> branch, the diff, the passing check, and the human's own words as the things
> that can be trusted, moves everything mechanical into a model-blind
> executable, leaves everything interpretive with the agent in the open, and
> says plainly when something has not been established.

## Where it came from

Four pressures, all observed on the v1 line rather than theorized:

1. **The runtime was the product's ceiling.** A ~2,000-line shell engine meant
   Windows needed WSL, contributors needed to reason about shell quoting to
   change a record shape, and every new capability grew the engine.
2. **Ceremony outpaced safety.** Consequence tiers, candidate digests,
   verification records, and path leases produced artifacts that *looked* like
   proof of correctness while proving only that the machinery had run. The
   friction was paid every increment; the safety was mostly notional.
3. **Records drifted from Git.** Plan notes recorded what a session believed at
   the time. Trusting them on resume was the single most reliable way to be
   wrong, and no amount of additional record-keeping fixed it.
4. **Knowledge rotted in specific, repeatable ways.** Notes cited plan IDs and
   `sources/` paths that were later archived; catalog entries and notes drifted
   apart; reconciliation was remembered only when someone remembered it.

v2 answers (1) with Go, (2) by deletion, (3) by making Git authoritative and
saying so in the instruction, and (4) with three enforced knowledge rules and a
completion step that hands the agent the candidate entries.

## Reading order

**The argument**
- [`design.md`](./design.md) — the normative overview. A reviewer can stop here.
- [`retired-machinery.md`](./retired-machinery.md) — what v1 had, why each piece
  went, and what protects that concern now. This is what separates a rewrite
  from a regression.

**The seam**
- [`executable-and-agent.md`](./executable-and-agent.md) — the central division:
  what the Go executable owns, what the agent owns, why the boundary falls where
  it does, and the CLI surface as the interface contract between them.
- [`records-and-ids.md`](./records-and-ids.md) — intents, plans, global ID
  allocation, locking, dates, archival, and the multi-clone honesty limit.

**The circuit**
- [`knowledge-circuit.md`](./knowledge-circuit.md) — `context/` notes, the
  catalog, the glossary, the durable-only boundary, and reconciliation on
  explicit completion.

**Doing the work**
- [`worktrees-and-reuse.md`](./worktrees-and-reuse.md) — worktree mechanics,
  the preservation rules, and copy-on-write environment reuse.
- [`stacked-execution.md`](./stacked-execution.md) — order derivation, waves vs.
  the linear chain, integration merges, and the stop conditions.
- [`roles-and-dispatch.md`](./roles-and-dispatch.md) — the four roles, role
  tiering without a model catalog, native host mapping, the two axes of
  parallelism, and independent review.

**The boundaries**
- [`authority-and-delivery.md`](./authority-and-delivery.md) — the two human
  decisions, scoped run authorization, delivery, completion, and cleanup as
  separate acts.
- [`versioning-and-distribution.md`](./versioning-and-distribution.md) — two
  products, two release lines, per-workspace CLI pinning, side-by-side installs,
  and release assembly.

**Vocabulary**
- [`glossary.md`](./glossary.md) — v2 terms, and the v1 terms that no longer
  exist.
