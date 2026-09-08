# Plan 0031 — Hide ok punchline in engine; trigger from coordinator

**Intent:** i022-ok-easter-egg  
**Repository:** context-circuit-source  
**Tier:** Standard  
**Status:** done

## Objective

When a user's entire message is the two letters ok (case-insensitive, optional
surrounding whitespace only), the coordinator invokes a hidden engine action,
prefixes its stdout as the first line of the reply, and continues the immediately
pending coordinator ask without granting new authority. The punchline string lives
only in `engine.sh`; ordinary instructions, skills, adapters, Product Knowledge,
engine help, and `engine-and-seam.md` do not contain it.

## Grounding (HEAD 1aceb907)

`engine.sh` dispatches verbs through `cc_main`; help prints only a generic Usage
line. Public verbs are listed in `engine-and-seam.md`. `coordinator.md` owns
conversational routing but has no lone-ok trigger today. INV-RUNTIME-01 keeps
routing out of the engine. `plan-allocate-id` fails on archived duplicate prefixes;
this plan uses the next in-band id **0031**.

## Decisions

**One plan, three tasks** — engine verb, coordinator trigger, and proof share one
execution and verification boundary; neither half is meaningfully verifiable alone
against AC-LONE-OK.

- Hidden verb **`ok-easter-egg`** prints a fixed stdout line; not listed in help
  or `engine-and-seam.md`.
- Coordinator matches **entire message** trimmed and case-folded to exactly `ok`.
- Proof: `test/easter-egg/test-ok-easter-egg.sh` plus a conversation harness plot
  for behavioral lone-ok and idle-ok cases.

## Tasks

### OKE-001 — Hidden engine invoke

Add `cc_ok_easter_egg` and register `ok-easter-egg` in `cc_main` only.

**Done when:** `sh test/easter-egg/test-ok-easter-egg.sh --case engine` passes.

### OKE-002 — Coordinator trigger

Extend `coordinator.md` to detect lone ok, invoke `ok-easter-egg`, prefix the
reply, continue pending work or wait when idle — without embedding the punchline.

**Done when:** `sh test/easter-egg/test-ok-easter-egg.sh --case coordinator` passes.

### OKE-003 — Proof suite and harness

Add grep/engine tests, criteria-map entries, acceptance wiring, and a conversation
plot for behavioral verification.

**Done when:** `sh test/easter-egg/test-ok-easter-egg.sh` and `sh test/acceptance.sh`
pass.

## Risks

- Punchline leaking into a skill or adapter during copy-paste — mitigated by grep
  suite over shipped surfaces.
- Coordinator treating ok as a hidden approval token — mitigated by explicit
  non-authorization wording and harness negative cases.
- Accidental listing in `engine-and-seam.md` — mitigated by omission check in OKE-003.

## Verification

Independent verifier runs task verification commands and full semantic acceptance.
No delivery, merge, or publication in this plan.
