# Lock in mark-done without an agent pre-check

Plan ID: 0001-mark-done-no-precheck
Intent: i016-mark-done-no-precheck
Status: draft

## Original request and coverage

- Human request (via intent i016): when asked to mark a plan (or several
  named plans) done, the agent marks each one done without first inspecting
  whether that plan's work and evidence exist; readiness stays a runtime
  refusal, reported plainly when a plan is not ready.
- Covered below: the instruction clarification (MDP-001), the several-named-
  plans case, and verification proving no pre-check language exists and the
  runtime refusal is unchanged.

## Objective and desired behavior

A human asking to mark one or several named plans done sees each one become
done directly — the agent calls `plan-complete` with no separate look at the
plan's work or evidence first. If a named plan is not ready, the runtime
`completion-ready` gate refuses and that refusal is reported plainly. This
holds identically whether one plan or several are named in the same ask.

## Constraints and non-goals

- Do not weaken or change `cc_completion_ready` or its refusal codes.
- Do not change what triggers mark-done — still an explicit human ask.
- Do not make verification, candidate acceptance, or delivery mark a plan
  done.
- Do not change in-place Product Knowledge reconcile after mark-done.
- Non-goal: inventing a new runtime verb or gate; the existing gate already
  does the readiness job.

## Product Knowledge grounding

- completion (context/domains/completion/README.md) — the owning domain;
  already documents completion-ready as the sole gate and that
  verification/candidate-acceptance/delivery never mark a plan done. This
  plan preserves that and makes the no-pre-check rule explicit and proven,
  including for several named plans in one ask.

## Repositories and source evidence

- context-circuit-source — owns `.agents/skills/cc-complete/SKILL.md`,
  `agents/coordinator.md`, and the completion runtime verb/gate. Evidence:
  `intent/i016-mark-done-no-precheck/trace/context-circuit-source.yaml` —
  the tracer found no existing pre-check code in either the skill or the
  coordinator instructions, and confirmed `cc_completion_ready` /
  `cc_plan_complete` in `wrapper/runtime/engine.sh` already implement the
  runtime-only gate.

## Tasks

1. MDP-001 (context-circuit-source, paths:
   `.agents/skills/cc-complete/SKILL.md`, `agents/coordinator.md`,
   `test/acceptance.sh`, `test/completion/`; depends on: none) — make the
   no-pre-check rule explicit for one plan and for several named plans in
   one ask, and add verification proving it; acceptance MDP-AC-001,
   MDP-AC-002, MDP-AC-003; verification MDP-VT-001, MDP-VT-002, MDP-VT-003.

## Acceptance criteria

- MDP-AC-001 — Asking to mark a plan done, or several named plans done,
  marks each named plan done without an agent step that inspects whether
  that plan's work or evidence exist.
- MDP-AC-002 — If a named plan is not ready to be done, mark-done refuses
  and that is reported plainly, with no separate agent-side look at work or
  evidence.
- MDP-AC-003 — A plan still becomes done only on an explicit mark-done ask;
  verification, candidate acceptance, and delivery do not mark it done.

## Verification

- MDP-VT-001 — `grep -rniE 'inspect.*work|check.*work|verify.*work|evidence.*before.*done' .agents/skills/cc-complete/ agents/coordinator.md`
  finds no remaining pre-check language.
- MDP-VT-002 — `grep -n 'cc_completion_ready\|cc_plan_complete' wrapper/runtime/engine.sh`
  shows the gate and verb are present and unchanged in shape.
- MDP-VT-003 — `sh test/acceptance.sh` passes, including new/extended
  completion coverage for the several-named-plans case.

## Assumptions, open questions, risks

- Assumption: the "pre-check" the intent refers to is agent conversational
  behavior/instructions, not existing runtime code — confirmed by the
  tracer, which found no pre-check code path in the current repository
  (disposition: already-answered).
- Open question (plan-resolution, carried from trace): whether a hosted
  template or another host adapter's coordinator text independently
  describes an inspection step is out of this repository's scope; not
  addressed here.
- Risk: accidentally touching `cc_completion_ready` while editing nearby
  text — mitigated by MDP-VT-002 asserting the gate/verb remain present and
  by keeping all edits to instruction text, not the runtime function body.
- Risk: missing a second location describing a pre-check — mitigated by
  MDP-VT-001's grep across both the skill and coordinator instructions.

## Expected commits and delivery notes

One commit in context-circuit-source. Delivery (pull request, merge) is a
separate explicit action, not implied by this plan or its execution.

## Expected Product Knowledge impact

No durable change to `context/domains/completion/README.md` is currently
expected — it already documents completion-ready as the sole gate and
verification/candidate-acceptance/delivery never marking a plan done. Review
this unit at completion to confirm nothing drifted.
