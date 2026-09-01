# Glossary and vocabulary

vNext adds concepts, so it must also hold a concept budget: **any new user-facing term
should retire one.** This file lists the new vocabulary, the internal→human
translations the coordinator uses, and the net change in how many concepts a human
carries.

## New terms

| Term | Meaning | Replaces / relates to |
| --- | --- | --- |
| **Intent** | The first-class decision for one change: goal, non-goals, constraints, acceptance criteria, scope, tier. Approved once, frozen. | new front door; the thing you approve instead of the plan |
| **Contract** | The machine record of an intent (`contract.yaml`); its `contract_digest` is frozen at approval. | — |
| **Acceptance criterion** | One executable-or-manual condition that defines "correct." | sharpens today's implicit plan acceptance |
| **Spec adversary** | An independent role that attacks the criteria *before* code. | new; applies "writer ≠ checker" to the spec |
| **Scope envelope** | The repositories and path regions an intent authorizes; plans must stay inside it. | new; what lets plan approval collapse safely |
| **Candidate** | The exact proposed result — a digest over the commit map, bases, and contract digest. Evidence binds to it. | replaces "the attempt/latest" as the unit of trust |
| **Tier** | Consequence level: **Explore / Standard / Critical**. Decides how much assurance, incl. whether a verifier spawns. | replaces the plan-vs-pairing binary |
| **Human acceptance** | A first-class record: this human accepted this candidate. | replaces the "mark done" flip |
| **Reconciliation debt** | A marker that delivered work has not yet updated Product Knowledge; it blocks the next grounding. | makes "gather context" un-forgettable |
| **Promote** | Turn an Explore session into a candidate-bearing change by attaching an intent and raising the tier. | replaces the pairing→plan cliff |
| **Change set** | The set of plans delivered as one pull request; the unit a shared candidate is computed over. | new; the delivery grouping |

## Internal → human translations (what the coordinator says)

The human never hears the left column.

| Internal | Said to the human |
| --- | --- |
| worktree / branch `cc/<plan>/<repo>` | "a separate working copy" |
| `intent-envelope-check` EXCEEDS | "this needs to go outside what you approved" |
| `candidate` void | "the code changed, so the earlier check no longer applies" |
| tier = critical | "higher-risk, so it gets an independent check and an explicit sign-off" |
| tier = explore | "quick and human-supervised — no independent check" |
| `BASE_UNBUILDABLE` | "these changes don't combine cleanly yet" |
| knowledge debt pending | "there's merged work I haven't folded into what the project knows yet" |
| host-blocked | "this environment can't run the required check safely" |

## Words with one meaning (kept disjoint)

As today, these never blur into each other:

- **Check** = observe evidence. It does **not** mean accept or deliver.
- **Accept** = a human confirms a candidate is right. It does not deliver.
- **Deliver** = the irreversible act (PR / merge / push / deploy). Gate 2.
- **Approve** = the human agrees to an **intent** (Gate 1). It is not delivery.
- **Verified** = passed an *independent* check (Standard/Critical only). Explore output
  is **human-supervised**, never "verified."

## The concept budget (net change)

vNext adds terms — but it also **removes** things the human juggles today, so the net
count should not rise:

| Added (must learn) | Removed (no longer juggled) |
| --- | --- |
| intent, tier, candidate | per-plan approval as a step |
| spec adversary | the "mark done" flip |
| reconciliation debt | remembering to "gather context" |
| promote (replaces a cliff) | "pairing vs plan" as two separate modes |

The design's rule of thumb: if a change adds a
user-facing noun, it must delete one. The measure of success is not how few concepts
the spec has, but how few the *human* has to hold to get a change shipped.
