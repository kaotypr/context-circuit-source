# Fixed rules and overridable defaults

Which workspace rules an extension can change, who decides that it may, and
where that decision is enforced.

## Two kinds of rule

| Kind | Rules | Changed by |
| --- | --- | --- |
| Fixed | approving an intent; authorizing push, pull requests, merges, deployment, deletion, and commits in a bound checkout; starting planning and execution; the plan record left alone during execution; no attribution; no secrets; borrowed knowledge read-only | nothing |
| Overridable default | `repair-loop`, `verify-during-execution`, `risk-review`, `delegation` | an extension that declares the override, once a person allows it for the workspace |

The list of defaults is closed. A manifest naming anything else as an override
is invalid, so an extension cannot claim approval or push authorization by
listing it. The fixed rules are what a workspace exists to keep for a person;
the defaults are judgments about cost and thoroughness that a team may
reasonably set differently.

## A person allows, the CLI enforces

An extension declares the overrides it wants. Nothing applies until a person
runs `extension override --id ID --allow NAME`, which writes the allowance
under the entry in `extensions.yaml`, refuses a default the manifest does not
declare, and may lower a repair loop's rounds but never raise them. The
allowance is per workspace; whether it applies to a stage still depends on the
extension being on for the member and its repositories covering the plan.

The instruction text only describes the exception. What makes it hold is the
CLI: an allowed override reaches the coordinator only as an `overrides` block
from `extension hooks`, and `agent dispatch` accepts the three dispatches an
override unlocks — a read-only verifier during execution, a worker whose task
carries findings, a reviewer nobody asked for — only with `--override NAME`
naming an allowed one. Those are exactly the dispatches refused without one, so
a workspace with no allowance behaves as it did before overrides existed.

## The coordinator runs the loop

A repair loop is implement, verify, fix, until clean or the round limit. The
coordinator runs every round: the verifier never dispatches the worker and the
worker never dispatches the verifier, and a role's own brief still forbids
rerunning a check until it passes. A round is one verification, so a limit of
two allows one fix, and `--round` lets the CLI refuse a round past the limit.
Fixes are committed on the worktree branch; the plan record is not edited and
the rounds go in the report. The verifier must not be the session that wrote
the code, so a host that cannot spawn an agent runs no loop.

`check` reports, each as needing a person, an override outside the list, an
allowance the manifest does not declare or on an extension that cannot be read,
rounds out of range, a verifier that is not read-only, and two extensions
delegating one stage to different roles.

Owner:

- `context-circuit-source@internal/workspace/extension_overrides.go` — the
  closed list, manifest validation, allowances, hook blocks, the dispatch
  check, status, and `check` findings.
- `context-circuit-source@internal/workspace/agents.go` — `--override` in
  `agent dispatch` and the brief's override line.
- `context-circuit-source@product/docs/customization.md` — the reference an
  agent follows.
