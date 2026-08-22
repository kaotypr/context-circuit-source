---
name: cc-entry
description: Enter or resume a Context Circuit workspace through the single two-stage router.
---

Use when a human asks to start, resume, orient, or explain the workspace.
Read the shipped entry adapters and Tier 0 only. Source the host-neutral
primitives from `wrapper/runtime/engine.sh`, select one probe, then load only
that context set before emitting one normalized route decision.

The router owner is `wrapper/contracts/routes.yaml`; cite invariant IDs rather
than restating lifecycle or ownership policy. A read-only orientation does not
create runtime state. Never infer a gate, lease, approval, or execution from a
generic start/resume request.

Before any write-like action, inspect `workspace.yaml` identity state. If the
workspace is `instantiated-workspace` with `identity.status: uninitialized`,
or has no accepted identity, route the request to initialization and present
the `identity-acceptance` gate. This includes generic requests such as “help me
build this”; they must not create files, plans, tasks, runtime state, or
repositories. Read-only orientation and explicitly selected read-only evidence
may continue without mutation. A `product-source` maintainer checkout is
already identified and follows its separate maintainer-source rules.

After identity acceptance, a generic build, implement, create, add, or change
request with no named approved plan routes to `draft-plan`. Drafting a plan is
the only safe next action; it does not authorize implementation. Approval and
the separate execution request remain required before any writer or
implementation files can be created.

For Codex CLI, Claude Code, and Cursor Agent CLI, record the normalized host
identifier and observed capability in `host_evidence` when entering or
resuming. The host adapter is not a second router: use the same probe and
reason codes. If a required writer or verifier child cannot be created, return
`block-missing-child-primitive` with `host-blocked`; never self-verify.
