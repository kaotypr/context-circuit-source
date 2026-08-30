# Context Circuit v0.7.0

Source design for Context Circuit v0.7.0 — a **delta on v0.6 / v0.6.1**. Read the
v0.6 design first ([../v0.6/](../v0.6/)) and the v0.5 base it builds on
([../v0.5/](../v0.5/)); v0.7.0 changes only what it names. This README is the
version index; each scope owns its own design.

## Scopes

- [direct-collaboration/](./direct-collaboration/) — **`cc-pair`, a direct
  interactive collaboration mode** for working on a bound repo with the agent, with
  three actors (user, coordinator, worker) and **no verifier, no lease, no
  execution records** — the human is the live oracle. It is **orthogonal to the
  plan lifecycle** (not a plan, not an execution), entered anytime a repo is bound,
  by intent/`/cc-pair`, or offered after a plan/stack execution completes. It
  answers three product questions — *no new role* (the worker is the existing role
  driven interactively), *one new skill* (`cc-pair`), and *one boundary invariant*
  (`INV-PAIR-01`: outside the lifecycle, own branch/worktree, output
  human-supervised and never "verified", never auto-completes or delivers). It adds
  **no plan-contract change**. It supersedes the earlier `ui-refinement` framing.
  Start at [direct-collaboration/design.md](./direct-collaboration/design.md).
- [execution-latency/](./execution-latency/) — making actions **finish sooner
  without changing what they mean** by applying two inference-layer levers whose
  safety the contract already guarantees — overlapping provably-independent
  run-stack plans (no new invariant; INV-CONCURRENCY-01/02) and per-role model
  tiering the host adapter applies on the child spawn (a coordinator/host concern;
  INV-HOST-01, INV-RUNTIME-01) — while the deterministic runtime stays thin. The
  only new record is the per-attempt `(model, effort)` evidence; an early
  "measure the layers" cut (phase timing, coordinator wall-clock) was built and
  trimmed. Start at [execution-latency/design.md](./execution-latency/design.md).
- [runtime-opacity/](./runtime-opacity/) — hardening the **invoke-not-read**
  boundary: the coordinator invokes `wrapper/runtime/engine.sh` as an opaque tool
  and must never read its implementation. A live v0.6 run showed a real coordinator
  read the engine source anyway, tripping the maintainer access gate
  non-deterministically. The fix strengthens the INV-RUNTIME-01 corollary (say
  "must not read," not "need not"; state it once and reference it) and closes any
  skill information gap that tempts the read — **no new invariant**. Start at
  [runtime-opacity/design.md](./runtime-opacity/design.md).
- [publication-intent/](./publication-intent/) — a **structured local home for the
  provider fields a publication pushes** (dates, time estimate) and a
  **consult-before-publish** preview. v0.6 records store external *identity* but not
  the field *values* sent, so reviewing or changing a date/estimate forces a live
  provider read. The fix adds a user-owned `intent/` layer (canonical
  `estimate_minutes`, `"2h 30m"` human I/O), a last-published `fields:` snapshot on
  the record, and a `cc-publish` preview that diffs and edits before pushing — **no
  new invariant, no new skill**; one INV-EXTERNAL-02 wording clarification permits a
  display-only drift read. Start at
  [publication-intent/design.md](./publication-intent/design.md).

## Layout convention

`sources/system-design/<product>/<version>/<scope>/`. Every folder has a
`README.md` landing/index; each scope's normative design is `design.md`, with
detail split into files or sub-folders as it grows. Version groupings use 3-number
semver going forward (v0.6.1 `design-layout-grouping`).

## Authority

v0.7.0 adds no owner it does not name and duplicates no rule.
**direct-collaboration** introduces one new capability (`cc-pair`), one new
invariant (`INV-PAIR-01`) for the interactive mode's boundary, and a light
session-pointer schema — and **no plan-contract change** (contract delta in
[direct-collaboration/skill-and-contract.md](./direct-collaboration/skill-and-contract.md)).
**execution-latency** introduces **no new skill and no new invariant** — its
safety is existing INV-CONCURRENCY-01/02 (overlap) and INV-HOST-01 /
INV-RUNTIME-01 (tiering); its only contract surface is additive schema fields and
coordinator/host policy (fan-out width, the tier ladder, and the host-adapter
spawn-application mechanism). **runtime-opacity** introduces **no new invariant** — it
strengthens the coordinator-side corollary of INV-RUNTIME-01 (invoke the engine,
never read it) as a wording and single-ownership cleanup across
`wrapper/adapters/AGENTS.md` and the invoking skills. **publication-intent**
introduces **no new skill and no new invariant** — a mode of the existing
`cc-publish`, one additive `intent/` file schema and a `fields:` snapshot on
`publication-record.yaml` (owners under `wrapper/contracts/schemas/`), and a wording
clarification of INV-EXTERNAL-02 (display-only drift read ≠ inbound flow). The
canonical owners stay under `wrapper/` (`wrapper/contracts/invariants.yaml`).
