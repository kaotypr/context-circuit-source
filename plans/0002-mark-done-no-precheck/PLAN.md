# Flip plan status on ask, with no check

Plan ID: 0002-mark-done-no-precheck
Intent: i016-mark-done-no-precheck
Status: done

## Original request and coverage

- What you asked: when you mark a plan done, or several plans done, each
  named plan's status becomes done. There is no check — no look at work or
  evidence, and no refusal because the plan was never built, failed a check,
  or was not "ready".
- Covered below: drop the unreadiness floor from the completion contract
  (UGD-001), ungate the mark-done verb (UGD-002), stop teaching a look in
  the agent instructions and Product Knowledge (UGD-003), and invert the
  tests that currently prove a refusal (UGD-004).
- Still explicit: a plan becomes done only when you ask. Checking the work,
  accepting the result, and delivering do not mark it done. Knowledge update
  after mark-done stays as it is. Explore stays planless.

## Objective and desired behavior

Asking to mark a plan done changes its status to done. Asking for several
plans does that for each named plan, on the same path. The agent does not
look first, and nothing else refuses for unreadiness.

## Constraints and non-goals

- Change the owning completion contract, not a skill-only exception that
  leaves a refusal in place.
- Same path for one named plan and for several named plans in one ask.
- Do not mark a plan done without an explicit ask.
- Do not let checking the work, accepting the result, or delivering mark a
  plan done.
- Do not change in-place Product Knowledge reconcile after mark-done.
- Explore stays planless.
- Do not deliver, publish, merge, or commit the implementation as part of
  this plan.

## Product Knowledge grounding

- completion (context/domains/completion/README.md) — owns mark-done; today
  it still gates the flip on a verified execution and completion-ready.
  Rewrite that so asking is enough; keep explicit mark-done and
  reconcile-after.
- assurance (context/domains/assurance/README.md) — the verifier floor stays
  at approval and at the independent check; it must stop being taught as a
  mark-done gate.
- verification (context/domains/verification/README.md) — checking the work
  still does not change plan status.
- delivery (context/domains/delivery/README.md) — delivering still does not
  mark a plan done.

Grounding summary: the unreadiness gate is layered across the contract, the
mark-done verb, the complete skill, and tests. Removing any one layer leaves
a refusal. Checking the work, accepting the result, and delivering already
do not mark done and must stay that way.

## Repositories and source evidence

- context-circuit-source — owns the completion contract, mark-done verb,
  complete skill, and completion domain. Evidence:
  `intent/i016-mark-done-no-precheck/trace/context-circuit-source.yaml` at
  HEAD `a128cdd`. Fifty unreadiness / refusal sites in the live product
  surface; none require a second repository.

## Tasks

1. UGD-001 (context-circuit-source, paths:
   `wrapper/contracts/invariants.yaml`,
   `wrapper/contracts/schemas/completion.yaml`,
   `wrapper/contracts/schemas/plan.yaml`; depends on: none) — drop the
   unreadiness floor from the completion contract; keep explicit mark-done.
   Acceptance UGD-AC-001, UGD-AC-002. Verification UGD-VT-001, UGD-VT-002.

2. UGD-002 (context-circuit-source, paths: `wrapper/runtime/engine.sh`;
   depends on: UGD-001) — ungate the mark-done verb so a never-built, failed,
   unverified, stale, or unaccepted draft becomes done when asked. Keep
   Explore planless and already-done as a status-machine refusal. Write a
   completion record only when the files exist. Acceptance UGD-AC-003,
   UGD-AC-004. Verification UGD-VT-003, UGD-VT-004.

3. UGD-003 (context-circuit-source, paths:
   `.agents/skills/cc-complete/SKILL.md`,
   `.agents/skills/cc-execute/SKILL.md`, `agents/coordinator.md`,
   `context/domains/completion/README.md`,
   `context/domains/assurance/README.md`; depends on: UGD-002) — stop
   teaching a look or refusal before the flip; loop the same verb for
   several named plans; keep reconcile-after. Acceptance UGD-AC-005,
   UGD-AC-006. Verification UGD-VT-005, UGD-VT-006.

4. UGD-004 (context-circuit-source, paths: `test/completion/`,
   `test/runtime/test-runtime.sh`, `test/assurance/test-tier.sh`,
   `test/candidate/test-candidate.sh`,
   `agent-harness/conversations/plots/stale-candidate-refuse-complete.md`,
   `agent-harness/scenarios/17-stale-candidate-requires-reverification/case.yaml`;
   depends on: UGD-003) — invert unreadiness tests; keep the tests that
   prove check / accept / deliver are not done. Acceptance UGD-AC-007,
   UGD-AC-008. Verification UGD-VT-007, UGD-VT-008, UGD-VT-009.

## Acceptance criteria

- UGD-AC-001 — The owning completion contract no longer requires a readiness
  or evidence look before mark-done can succeed.
- UGD-AC-002 — Verification, candidate acceptance, and delivery still do not
  mark a plan done in the contract.
- UGD-AC-003 — Asking to mark a never-built, failed, unverified, stale, or
  unaccepted plan done changes that plan's status to done.
- UGD-AC-004 — Delivery, candidate acceptance, change-set complete, and
  inferred completion still do not set plan status done. Explore stays
  planless.
- UGD-AC-005 — Agent instructions and Product Knowledge no longer refuse or
  inspect before flipping status on a mark-done ask, including for several
  named plans.
- UGD-AC-006 — Agent instructions still say verification, candidate
  acceptance, and delivery do not mark a plan done, Explore is planless, and
  knowledge update after mark-done stays as it is.
- UGD-AC-007 — Semantic tests prove mark-done succeeds with no execution,
  failed execution, missing acceptance, and stale candidate, and that two
  named plans each become done on the same path.
- UGD-AC-008 — Semantic tests still prove verification, candidate
  acceptance, delivery, inferred completion, and change-set complete do not
  mark a plan done.

These prove the intent outcomes AC-NO-CHECK and AC-STILL-EXPLICIT.

## Verification

- UGD-VT-001 — no "completion-ready refuses to complete" / "enforces the
  tier floor" coupling left on the completion invariants.
- UGD-VT-002 — plan.yaml no longer says draft → done is inferred from
  candidate acceptance plus delivery.
- UGD-VT-003 — the mark-done verb does not call completion-ready or fail the
  unreadiness codes. Explore-planless may still appear.
- UGD-VT-004 — delivery-record and change-set-complete do not set plan
  status; inferred completion still refuses.
- UGD-VT-005 — cc-complete has no unreadiness refusal and no
  completion-ready precheck. It still says delivery / acceptance /
  verification do not mark done.
- UGD-VT-006 — completion and assurance domain pages no longer teach
  completion-ready as a mark-done gate.
- UGD-VT-007 — `sh test/completion/test-completion.sh` passes, including
  never-built, failed, unverified, and two-plan cases.
- UGD-VT-008 — `sh test/completion/test-inferred.sh` passes: stale still
  becomes done; accept / deliver / infer / change-set still do not.
- UGD-VT-009 — `sh test/acceptance.sh` passes.

## Assumptions, open questions, risks

Assumptions (plan-level, from the look at the code):

- Missing execution / acceptor / candidate / revision: still flip status;
  write a full completion record only when those files exist.
- completion-ready remains a read-only eligibility query, disconnected from
  mark-done.
- Several named plans: loop the existing one-plan verb. No new engine verb.
- Already-done is a status-machine refusal, not unreadiness. Explore-tier
  executions stay refused.

Open questions: none that change the approved decision. The unreadiness
question is already answered: yes, asking still flips it.

Risks:

- Changing only the complete skill leaves the mark-done verb still refusing.
- Dropping the first gate but leaving finalize's missing-acceptor /
  missing-revision refusals still blocks never-built plans.
- Required completion-record fields become a hidden unreadiness gate if the
  status flip still waits on a full record.
- Letting check / accept / deliver mark done by accident.
- Leaving completion-ready described as "refuses to complete" in the
  contract or assurance page restores the gate in writing.
- The stale-candidate conversation plot currently requires a refusal; leaving
  it grades a correct ungated flip as failure.
- Collapsing "accept the result" into "mark done" once unreadiness is gone.

## Expected commits and delivery notes

- One commit on context-circuit-source covering the contract, mark-done
  verb, instructions, Product Knowledge, and tests.
- Delivery (pull request, merge) is a separate explicit action. This plan
  does not deliver.

## Expected Product Knowledge impact

- Reassess completion and assurance at completion: mark-done is ask-then-
  flip, with no unreadiness look; the verifier floor is not a complete
  gate; reconcile-after-mark-done is unchanged.
