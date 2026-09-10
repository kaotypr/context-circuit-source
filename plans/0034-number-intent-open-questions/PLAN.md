# Plan 0034 — Number open questions on the intent page

**Intent:** i025-number-intent-open-questions  
**Repository:** context-circuit-source  
**Tier:** Standard  
**Status:** done

## Objective

Open questions on the human-facing intent use the same numbered 1, 2, 3 list style
as **The plans**: bold question, italic answer beneath, stable numbers when
answered, next unused number for questions added later, and an unnumbered empty-state
line when there are none. New intents, the template, the filled example, the
`cc-intent` skill, the intent domain page, and the Gate 1 approval ask teach and
use this form. Already-written intents are left unchanged. No machine ids enter
`contract.yaml` or `contract_digest`.

## Grounding (HEAD 05e384ed)

| Surface | Current state |
|---------|---------------|
| `.context-circuit/docs/templates/intent.md` | Empty state line present; HTML comment shows unnumbered `**question**` / `_Answer:_` |
| `.context-circuit/docs/templates/intent.example.md` | One answered question, not numbered |
| `.agents/skills/cc-intent/SKILL.md` | Phase-aware Open questions; no numbering convention |
| `context/domains/intent/README.md` | Bold/italic shape documented; no numbering |
| Release manifest | Already requires both intent templates |

**The plans** already demonstrates the target list shape (`1. **title**` /
`_After this:_`). Historical intents such as `i026-compiled-runtime` keep
unnumbered questions and must not be rewritten.

`plan-allocate-id` fails `PLAN_PREFIX_COLLISION` on archived duplicate numeric
prefixes; this plan uses the next in-band id **0034**.

## Decisions

**One plan, two tasks** — template/skill/domain updates and acceptance proof share
one execution and verification boundary.

- Canonical template paths are `.context-circuit/docs/templates/` (not legacy
  `docs/templates/`).
- Numbering policy lives only in template, `cc-intent` skill, and intent domain
  page (INV-SKILL-01 / one owner per rule).
- `template/` thin routes stay unchanged; they already import the skill owner.
- `intent.example.md` shows one numbered answered question; skill text carries
  stable-number and empty-state rules.

## Tasks

### NQ-001 — Teach numbered open questions on intent authoring surfaces

**Repository:** context-circuit-source  
**Paths:** `.context-circuit/docs/templates/intent.md`,
`.context-circuit/docs/templates/intent.example.md`, `.agents/skills/cc-intent/`,
`context/domains/intent/README.md`

Update the blank template comment, number the example question as `1.`, extend
`cc-intent` Open questions and approval guidance (including answer-by-number in
the approval ask), and mirror the convention on the intent domain page.

**Acceptance:** NQ-AC-NUMBERED, NQ-AC-STABLE, NQ-AC-EMPTY, NQ-AC-STYLE,
NQ-AC-AUTHORING, NQ-AC-EXISTING-UNTOUCHED

**Verification:**

- `NQ-VT-001` — numbered guidance present in template, example, skill, domain
- `NQ-VT-002` — stable/empty/untouched rules present; `i026` remains unnumbered

### NQ-002 — Prove numbered open questions in the acceptance suite

**Repository:** context-circuit-source  
**Paths:** `test/intent/`, `test/acceptance.sh`  
**Depends on:** NQ-001

Add `test/intent/test-open-questions-numbering.sh` and wire it into
`test/acceptance.sh`.

**Acceptance:** NQ-AC-NUMBERED, NQ-AC-STABLE, NQ-AC-EMPTY, NQ-AC-AUTHORING

**Verification:**

- `NQ-VT-003` — `sh test/intent/test-open-questions-numbering.sh`
- `NQ-VT-004` — `sh test/acceptance.sh`

## Risks

- **Drift into other skills** — numbering policy copied into `cc-trace` or
  coordinator files would violate one-owner-per-rule; NQ-VT-001 scopes checks to
  the three owning surfaces.
- **Accidental mass renumbering** — worker might “fix” historical intents; NQ-VT-002
  and explicit non-goal wording guard against that.

## Expected Product Knowledge impact

Review `context/domains/intent/README.md` on completion (`knowledge_impact: intent`).
