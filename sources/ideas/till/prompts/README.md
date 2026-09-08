# Till prompts — every case

Plain things you **type** to the coordinator while building Till. Each case
matches one expected conversation from Context Circuit's whole-surface catalog
(39 cases). One extra case covers a fuller write-up of a single change
(intent detail), which that catalog does not list on its own.

Play them in a **product workspace**. Speak like this pack: goals and decisions,
not internal names.

## How to play

1. **Spine first** (the numbered path below). That is a real Till build.
2. Then **branches** when the situation in *When* is true. Do not force a
   refusal case if nothing is being refused.
3. Copy the **Say** lines. Wait for the reply before the next line.
4. Stop if the coordinator invents a project folder you did not ask for, starts
   building before you approved, calls try-out work "independently checked", or
   says something shipped when you only accepted the result.

Name sources by path, for example
`sources/ideas/till/product-brief.md`.

## Spine (build Till)

| Play | File | Case |
| --- | --- | --- |
| 1 | [01-orientation.md](01-orientation.md) | onboarding → orient → init both folders → get up to speed |
| 2 | [02-intent-gate.md](02-intent-gate.md) | design the product, then approve one change (do not build yet) |
| 3 | [04-execute-and-verify.md](04-execute-and-verify.md) | build it; checked is not finished |
| 4 | [05-candidate-delivery.md](05-candidate-delivery.md) | accept, open onto main, Standard finishes |
| 5 | [06-knowledge.md](06-knowledge.md) | fold what you shipped into what the project knows |
| ★ | [08-whole-flow.md](08-whole-flow.md) | the same spine as **one** conversation (invoice list) |

## Coverage

Every catalog case has a prompt. The spine cases also appear as isolated
branches so you can replay a single beat.

### Orientation and knowledge foundation

| Case | File |
| --- | --- |
| onboarding-what-is-this | 01 |
| orient-new-project | 01 |
| connect-existing-repo | 01 |
| clone-or-init-new-repo | 01 |
| build-product-knowledge | 01 |
| query-product-knowledge | 01 |
| accept-or-defer-context-proposal | 01 |

### Intent gate

| Case | File |
| --- | --- |
| author-system-design-spawns-intents | 02 |
| review-intent-approve-derive | 02 |
| feasibility-question-then-approve | 02 |
| refuse-before-gate1 | 02 |
| tier-fails-upward-refuse-explore | 02 |
| scope-reach-feasibility-question | 02 |
| approve-and-build-one-turn | 02 |
| intent-detail (fuller write-up of one change) | 02 |

### Explore

| Case | File |
| --- | --- |
| direct-collaboration-explore | 03 |
| explore-promote-to-standard | 03 |

### Execute and verify

| Case | File |
| --- | --- |
| execute-standard-verify-not-complete | 04 |
| repository-grounding | 04 |
| verifier-host-blocked | 04 |
| run-stack-single-repo | 04 |
| run-stack-multi-repo | 04 |
| execution-tiering-hidden | 04 |
| three-failure-stop | 04 |
| interrupted-recovery | 04 |
| lease-ownership-conflict | 04 |

### Candidate, completion, delivery

| Case | File |
| --- | --- |
| stale-candidate-refuse-complete | 05 |
| critical-repair-then-complete | 05 |
| delivery-boundary-block-no-remote | 05 |
| standard-inferred-completion | 05 |
| change-set-one-verification | 05 |
| change-set-base-unbuildable | 05 |

### Knowledge loop

| Case | File |
| --- | --- |
| knowledge-debt-blocks-next-plan | 06 |
| reconcile-and-proceed | 06 |

### Organization and publish

| Case | File |
| --- | --- |
| archive-plan | 07 |
| restore-archived-plan | 07 |
| archive-restore-intent | 07 |
| publish-plan-to-external | 07 |
| publish-open-questions-thread | 07 |

### Whole flow

| Case | File |
| --- | --- |
| standard-feature-whole-flow | 08 |
