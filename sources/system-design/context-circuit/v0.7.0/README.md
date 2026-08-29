# Context Circuit v0.7.0

Source design for Context Circuit v0.7.0 — a **delta on v0.6 / v0.6.1**. Read the
v0.6 design first ([../v0.6/](../v0.6/)) and the v0.5 base it builds on
([../v0.5/](../v0.5/)); v0.7.0 changes only what it names. This README is the
version index; each scope owns its own design.

## Scopes

- [ui-refinement/](./ui-refinement/) — an **interactive, human-gated refinement
  mode** for taste-driven work whose acceptance criterion is subjective ("refine
  the UI until it looks good"). Its driving instance is frontend/UI layout, but
  the mechanism is general. It answers three product questions at once — *no new
  role* (roles are authority-shaped, domain competence is loaded context), *one
  new skill* (`cc-refine`), and *a distinct execution mode bracketed by the normal
  spine* (one plan, one lease, one independent verify + baseline freeze at the
  end). Start at [ui-refinement/design.md](./ui-refinement/design.md).
- [execution-latency/](./execution-latency/) — making actions **finish sooner
  without changing what they mean**. Measure first (two clocks: deterministic vs
  inference), then apply two inference-layer levers whose safety the contract
  already guarantees — overlapping provably-independent run-stack plans (no new
  invariant; INV-CONCURRENCY-01/02) and model tiering (a coordinator/host concern;
  INV-HOST-01, INV-RUNTIME-01) — while the deterministic runtime stays thin. Start
  at [execution-latency/design.md](./execution-latency/design.md).
- [runtime-opacity/](./runtime-opacity/) — hardening the **invoke-not-read**
  boundary: the coordinator invokes `wrapper/runtime/engine.sh` as an opaque tool
  and must never read its implementation. A live v0.6 run showed a real coordinator
  read the engine source anyway, tripping the maintainer access gate
  non-deterministically. The fix strengthens the INV-RUNTIME-01 corollary (say
  "must not read," not "need not"; state it once and reference it) and closes any
  skill information gap that tempts the read — **no new invariant**. Start at
  [runtime-opacity/design.md](./runtime-opacity/design.md).
- [publication-fields/](./publication-fields/) — a **structured local home for the
  provider fields a publication pushes** (dates, time estimate) and a
  **consult-before-publish** preview. v0.6 records store external *identity* but not
  the field *values* sent, so reviewing or changing a date/estimate forces a live
  provider read. The fix adds a user-owned `desired/` intent layer (canonical
  `estimate_minutes`, `"2h 30m"` human I/O), a last-published `fields:` snapshot on
  the record, and a `cc-publish` preview that diffs and edits before pushing — **no
  new invariant, no new skill**; one INV-EXTERNAL-02 wording clarification permits a
  display-only drift read. Start at
  [publication-fields/design.md](./publication-fields/design.md).

## Layout convention

`sources/system-design/<product>/<version>/<scope>/`. Every folder has a
`README.md` landing/index; each scope's normative design is `design.md`, with
detail split into files or sub-folders as it grows. Version groupings use 3-number
semver going forward (v0.6.1 `design-layout-grouping`).

## Authority

v0.7.0 adds no owner it does not name and duplicates no rule. **ui-refinement**
introduces one new capability (`cc-refine`), one new plan acceptance *kind*
(human-gated visual acceptance), and provisional invariants for the interactive
mode (contract delta in
[ui-refinement/skill-and-schema.md](./ui-refinement/skill-and-schema.md)).
**execution-latency** introduces **no new skill and no new invariant** — its
safety is existing INV-CONCURRENCY-01/02 (overlap) and INV-HOST-01 /
INV-RUNTIME-01 (tiering); its only contract surface is additive schema fields and
coordinator policy. **runtime-opacity** introduces **no new invariant** — it
strengthens the coordinator-side corollary of INV-RUNTIME-01 (invoke the engine,
never read it) as a wording and single-ownership cleanup across
`wrapper/adapters/AGENTS.md` and the invoking skills. **publication-fields**
introduces **no new skill and no new invariant** — a mode of the existing
`cc-publish`, one additive `desired/` file schema and a `fields:` snapshot on
`publication-record.yaml` (owners under `wrapper/contracts/schemas/`), and a wording
clarification of INV-EXTERNAL-02 (display-only drift read ≠ inbound flow). The
canonical owners stay under `wrapper/` (`wrapper/contracts/invariants.yaml`).
