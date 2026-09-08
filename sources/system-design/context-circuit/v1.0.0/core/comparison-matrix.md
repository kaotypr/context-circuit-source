# Comparison matrix — v1.0 vs. current Context Circuit

Legend: 🟩 kept · 🟧 reworked · 🟦 new · 🟥 removed. The baseline is the current
product contract under `wrapper/` (v0.7), read line-by-line.

## The ideas v1.0 adopts and declines

Nine design ideas v1.0 weighed, and where each stands in current Context Circuit
versus v1.0. ✓ = present · ◐ = partial/inherited · ✗ = absent/declined.

| Idea | CC today | v1.0 | |
| --- | :--: | :--: | --- |
| Gate outcomes, not every transition | ✗ | ✓ | 🟧 |
| Candidate/tree‑bound evidence staleness | ✗ | ✓ | 🟦 |
| Consequence/risk‑based process tiering | ✗ | ✓ | 🟦 |
| Verifier gets evidence, not prose | ◐ | ✓ | 🟩 |
| No fake multi‑repo atomicity | ◐ | ✓ | 🟩 |
| Event‑sourced / replayable ledger = truth | ✗ | ✗ declined | — |
| Deterministic non‑model orchestrator | ◐ | ◐ | 🟩 |
| Fencing tokens on leases | ✗ | ✗ declined | — |
| Deterministic context/token compiler | ✗ | ✗ declined | — |

v1.0 adopts the top five and declines the bottom three. The one deliberate ◐
(deterministic non‑model orchestrator) keeps the engine model‑blind but leaves
routing to the conversational coordinator — a CC strength it declines to trade away
(see `preserved-core.md`, `roles-and-spawning.md`). The three declines
(event‑sourcing, fencing tokens, a context/token compiler) each solve a problem this
single‑host, human‑in‑the‑loop product does not have today; each is a defensible
future addition, not part of the core evolution.

## Core philosophy

| Aspect | Context Circuit today | v1.0 | |
| --- | --- | --- | --- |
| What is gated | every plan **state transition** | the **claim of correct** (intent) + **delivery** | 🟧 |
| Human gates | 4 (approve, execute, complete, deliver) | 2 (intent, delivery) | 🟧 |
| Unit of trust | the execution/attempt | the **candidate** | 🟧 |
| Definition of "correct" | one human glance in plan approval | explicit intent gate + **post-approval tracing** | 🟦 |
| Ceremony | fixed full lifecycle (+ separate pairing mode) | **one consequence ladder** Explore→Critical | 🟧 |

## Mechanisms

| Capability | Today | v1.0 | |
| --- | --- | --- | --- |
| Intent / bigger-picture object | none (sources → proposals → plan) | first-class `intent/` + frozen contract | 🟦 |
| Feasibility/risk checked against real code, before any plan or code change | no | **the tracer**, spawned automatically right after Gate 1 | 🟦 |
| Plan approval | explicit human, every plan | automatic once the intent is approved | 🟧 |
| Scope safety | task scope per plan (INV-PLAN-02) | settled at **delivery** (Gate 2); the **feasibility check** surfaces a required change beyond a bound scope | 🟧 |
| Acceptance criteria precision | vague until plan, or fixed upfront by guesswork | **outcome-level** in the intent; the **executable** check is a tracing output, carried in the plan | 🟦 |
| Evidence freshness | per-attempt, conversational | **candidate digest**, auto-void on change | 🟦 |
| One-PR / stacked plans | verifier per plan | one candidate → one verification | 🟧 |
| Human acceptance record | not first-class | **candidate-bound acceptance record** | 🟦 |
| Verification | always independent (INV-VERIFY-01) | tiered: Explore supervised / Std+Crit independent | 🟧 |
| Direct collaboration | separate mode (INV-PAIR-01) | **Explore tier**, promotable | 🟧 |
| Completion | explicit human flip | inferred (Explore/Std), explicit (Critical) | 🟧 |
| Knowledge reconciliation | generated, acted-on manually, forgettable | **debt marker blocks next grounding** | 🟧 |

## Preserved unchanged (the mechanics + spine)

| Kept | Rule | |
| --- | --- | --- |
| Writer ≠ checker; never self-verify | INV-VERIFY-02 | 🟩 |
| Model-blind, opaque runtime | INV-RUNTIME-01/02 | 🟩 |
| Worktree isolation from base tip | INV-EXEC-03 | 🟩 |
| Commit-before-verify, new-commit repairs | INV-EXEC-04 | 🟩 |
| Three-failure counter, waived≠passed | INV-REPAIR-01 | 🟩 |
| Preserve-on-failure | INV-PRESERVE-01 | 🟩 |
| Path leases + region overlap + descendant exemption | INV-CONCURRENCY-01 | 🟩 |
| Base selection + integration merge + stale rebuild | INV-CONCURRENCY-02 | 🟩 |
| Delivery separateness + drift guard (Gate 2) | INV-DELIVER-01/02 | 🟩 |
| Portable identity vs. local binding; base branch vs default | INV-REPO-01..04 | 🟩 |
| Credentials out; passive sources | INV-SEC-01/02 | 🟩 |
| Repository grounding + brief assembly | INV-GROUND-01/02/03 | 🟩 |
| Host neutrality; per-role model/effort tiering | INV-HOST-01 | 🟩 |
| Conventional Commits; no AI attribution | INV-COMMIT-01 | 🟩 |
| Publication orthogonal, export-only, self-contained | INV-EXTERNAL-01/02/03 | 🟩 |
| Retrieval-first Product Knowledge | INV-KNOWLEDGE-01 | 🟩 |
| Skills at `.agents/skills/<name>/SKILL.md` | INV-SKILL-01 | 🟩 |
| Archive/restore as status-blind moves | INV-ARCHIVE-01/02 | 🟩 |
| The engine-action CLI seam; `cc_digest` + atomic write | (runtime) | 🟩 |

## Removed

| Removed | Replaced by | |
| --- | --- | --- |
| Per-plan approval gate as a human step | intent gate; scope settled at delivery (Gate 2) | 🟥 |
| Unconditional verifier per execution | tiered verification | 🟥 |
| Explicit mark-done at low tier | inferred completion | 🟥 |
| `cc-pair` as a **separate** mode | Explore tier of the one ladder | 🟥 |
| Forgettable manual reconciliation | debt-blocked grounding | 🟥 |

## Deliberately NOT adopted (from the convergent-idea cluster)

| Idea | Why declined | |
| --- | --- | --- |
| Fencing tokens on leases | cheap later; single-host CC doesn't hit the race | 🟥 |
| Full event-sourced engine rewrite | append-only records capture most value; schema-evolution cost not worth it now | 🟥 |
| Deterministic context/token compiler | invisible present quality, not needed for the 8 pains; deferred | 🟥 |

## Summary line

v1.0 keeps everything the code read shows is load-bearing and correct, spends its
new energy on the two things Context Circuit under-invests in — **proving what
"correct" means** and **making evidence un-fakeable** — and moves the human's
attention from four transition gates to the two decisions only a human can make.
