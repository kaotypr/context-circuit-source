# Named-plan review with optional host question prompts

Status: draft  
Repository: context-circuit-source  
Source: repository-evidence

## Review summary

`Review plan <id>` already routes to read-only `review-plan` and a Review
Card. That card lists risks and up to three human decisions as prose, so the
human has to parse them and type the next request. There is no dedicated
discovery skill, unnamed review does not clarify the target the way unnamed
approval does, and hosts that already expose an optional structured question
prompt are unused.

This plan keeps the existing route, context set, and read-only rule. It makes
named-plan review a specific request and skill, then offers the same 0–3
decisions through the host native question UI when that primitive is present.

## What approval authorizes

Approval authorizes bounded implementation of the tasks and acceptance
criteria in `plan.yaml`. It does not authorize execution, lease acquisition,
worktree creation, Git changes, delivery, publication, deployment, merge, or
cleanup.

## Scope and non-goals

Paths and behavior boundaries are linked to `plan.yaml`. Review stays
read-only. Question-prompt answers are presentation and conversation, not a
new gate. Host-local Cursor commands and policy files stay out of scope.

## Proposed solution

Keep `wrapper/runtime/engine.sh` as the only evaluator. Tighten Stage A so
named review phrasing selects `plan-review`, and so a review request without a
usable plan id becomes `clarify-target` instead of inspecting an arbitrary
bundle.

Add a thin `.agents/skills/cc-review-plan/SKILL.md` adapter for host
discovery. The human trigger remains natural language, especially
`Review plan <id>`, `Walk me through plan <id>`, and questions about that
plan’s risks or open decisions. `cc-plan` continues to draft and prepare
bundles; it does not keep a second review procedure.

The review procedure itself stays in the existing owners: load the
`plan-review` context set, inspect the named `plan.yaml`, `PLAN.md`, task
headers, declared Product Knowledge, dependencies, archive sidecar, and
delivery boundary, then emit the Review Card from `docs/plan-review.md`. The
card still reports outcome, summary, approval scope, non-effects, tasks,
acceptance mapping, risks/assumptions, contradictions, and next action.

When the card has one to three focused human decisions and the current host
exposes an optional native question prompt, present those decisions there
after the card. Known host primitives:

| Host | Optional primitive | Fallback |
| --- | --- | --- |
| Cursor Agent | `AskQuestion` | Review Card text |
| Claude Code | `AskUserQuestion` | Review Card text |
| Codex CLI | `request_user_input` when listed | Review Card text |

The prompt is optional. Codex often exposes it only in some modes or behind a
host setting; a missing, denied, or failed prompt is not `host-blocked` and
must not retry as if it were a required child. Do not record question
transcripts, answers, or tool payloads in `host_evidence`.

Answers never write `plan.yaml` or task status. A next-action label such as
`Revise this plan` or `Approve this plan` may continue only into that route’s
existing first card in the same session. Selecting `Approve this plan` still
requires the current session-bound confirmation; it cannot become
`approve-plan`.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| IPR-001 | Named-plan review request, clarify-target, and dedicated skill | — |
| IPR-002 | Optional host question prompt on top of the Review Card | IPR-001 |
| IPR-003 | Routing, host, release, and semantic coverage | IPR-001, IPR-002 |

## Acceptance criteria

- IPR-AC-01: Named review stays `review-plan` and read-only.
- IPR-AC-02: Missing plan id clarifies the target instead of guessing.
- IPR-AC-03: `cc-review-plan` is discoverable; humans still speak naturally.
- IPR-AC-04: Risks and 0–3 decisions appear on the card, and in the host
  question UI when that primitive exists.
- IPR-AC-05: Question answers cannot skip a human gate.
- IPR-AC-06: Existing semantic and release boundaries remain valid.

## Verification

The verification IDs and canonical commands are owned by `plan.yaml`. An
independent verifier must inspect the route decision, Review Card contract,
skill/release allowlists, and complete acceptance output. Live host question
UIs are optional and must not be required by offline CI.

## Risks, assumptions, and open decisions

- Native question prompts differ by host and are sometimes mode-gated; the
  product cannot require them.
- Putting next-action labels in the question UI is convenient, but a chosen
  label must not be treated as gate confirmation.
- Durable review artifacts are out of scope; answers stay in the conversation
  unless a later explicit revise request changes the bundle.
- Adding `question_prompt_capability` to `host_evidence` is deferred so
  transcripts and tool payloads cannot leak into records.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Prepared from the shipped `review-plan` route and `plan-review` context set,
the Review Card in `docs/plan-review.md`, host-capability and skill-discovery
contracts, and the observed optional question primitives on Cursor Agent,
Claude Code, and Codex CLI.
