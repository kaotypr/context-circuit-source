# Mechanism 3 — Tiered assurance, and cc-pair folded in

Fixes pain 5 ("do I even need a verifier?") and pain 6's over-spawning. Resolves
the cc-pair question by folding direct collaboration into the bottom of one ladder
instead of a separate mode.

## The problem it solves

Today verification is **unconditional**: one independent verifier per execution,
always (INV-VERIFY-01). For low-risk solo work that feels heavy, and for stacked
plans it multiplies. Separately, `cc-pair` exists as a *disjoint* low-ceremony mode
with no verifier at all — which means Context Circuit already has a "no verifier"
answer, but it quarantined that answer into an island with a cliff between it and
the plan world.

vNext replaces both the always-on verifier and the separate pairing mode with **one
continuous consequence ladder**.

## The ladder

`tier` is declared on the intent's `contract.yaml` (M1) and may be raised by the
human. It composes assurance from parts, rather than switching between two worlds:

| Tier | Independent verifier | Spec adversary | Completion | Typical use |
| --- | --- | --- | --- | --- |
| **Explore** | none — human-supervised | skipped/light | inferred from acceptance | small, live, reversible, solo |
| **Standard** | required, candidate-bound | on the criteria | inferred from acceptance + delivery | ordinary feature work |
| **Critical** | required + explicit re-check | full battery | **explicit human** completion | many repos, security, migration, production, breaking change |

Rules:

- **"Never self-verify" stays absolute** (INV-VERIFY-02). Explore is *human-
  supervised*, which is honest — its output is never labeled "verified." It is not
  a weaker verifier; it is no independent verifier plus a live human oracle,
  exactly what `cc-pair` is today.
- **Standard and Critical require an independent read-only verifier bound to the
  current candidate** (M2). The verifier's structural read-only enforcement is the
  unchanged kept mechanic.
- **The tier decides whether a verifier spawns at all** (pain 5) and, for stacked
  work, how many (pain 6 — combined with M2's one-candidate rule).
- **The engine stays model-blind.** Tier is *declared* by the coordinator and the
  engine only enforces the floor: `completion-ready` refuses to complete a Critical
  candidate without an independent `passed` verdict bound to the current candidate,
  and refuses to treat an Explore result as "verified." The runtime never chooses a
  tier — consistent with INV-RUNTIME-01 and the way `complexity` is already an
  additive, non-gating hint today.

## How the tier is chosen (and the second crown jewel)

Tier is set at intent time from **transparent risk signals** — repository count,
reversibility, security/privacy surface, money, data migration, production
availability, novelty — and the human may override upward. The signal set is
declared and auditable ("why is this Critical?" fits on one line).

Skipping the verifier at Explore is safe **only** if tiering is deterministic and
**fails upward** (when unsure, tier higher). This is the second crown-jewel check
(the first is M1's envelope drift). See `risks-and-open-questions.md`.

## cc-pair becomes the Explore tier

Direct collaboration is not deleted — it is **promoted from a separate mode to the
bottom rung of the ladder**. The separate-mode design is what gave cc-pair its one
structural flaw (a disjoint island with a cliff between it and the plan world);
making it the Explore tier removes that flaw while keeping everything cc-pair does
well.

What is preserved from today's `cc-pair` (INV-PAIR-01) unchanged:

- User + coordinator (interprets, never writes) + one worker, live human as the
  acceptance oracle.
- Git isolation: a fresh `cc-pair/<session>` branch and worktree from a chosen
  base; the active checkout is never touched — **topological safety, not policy**.
- Output is human-supervised, never "verified"; worker commits only on explicit
  request; delivery blocks on anchor drift rather than silently rebasing.
- The light resumable pointer as the only runtime state.

What changes — the cliff becomes a ramp:

- **Recordless by default, but promotable.** An Explore session works fluidly with
  no candidate and no ceremony (the whole point of pairing). If it turns out to be
  real work, the human **promotes it in place**: attach an intent + criteria (the
  adversary can now challenge them), raise the tier to Standard/Critical, and spawn
  the independent verifier — **without stopping and restarting as a separate
  plan**. Promotion is the moment the work enters the trust system and gets a
  candidate.
- **Its output stops being invisible.** Once promoted, the work has a candidate, so
  it can feed reconciliation into Product Knowledge (M4) and deliver through the
  normal path — closing the gap where today's pairing work never teaches the
  knowledge base anything.

So the answer to "verifier or not?" stops being a *world you switch into* and
becomes a *tier you slide*. `cc-pair`'s good instinct (not everything needs a plan)
is kept; its one structural flaw (a disjoint island) is removed.

## Skills and engine

- **Changed:** `completion-ready` reads the intent `tier` and enforces the
  verifier floor accordingly; `cc-execute`/`cc-run-stack` spawn the independent
  verifier only at Standard/Critical.
- **`cc-pair` reframed** as the Explore tier with an explicit **promote** step
  (attach intent, raise tier) rather than a standalone mode. INV-PAIR-01's content
  is preserved but relocated into the tier ladder (see `invariant-deltas.md`).

## What stays the same

The verifier's read-only enforcement, the three-failure counter, preserve-on-
failure, and the pairing git-isolation mechanics (`cc_pair_*`) are all unchanged.
Tiering is a *policy* over when those mechanics run — not a change to the mechanics.
