# Domain Knowledge

`context/domains/<domain>/` is the canonical home for knowledge about one
bounded project area. The domain page owns exact behavior, workflows, rules,
interfaces, data, constraints, edge cases, implementation references, and
verification when those facts are supported by evidence.

Generated domain pages use [the domain template](../../.context-circuit/docs/templates/domain-context.md)
and normally live at `context/domains/<domain>/README.md`. Detailed workflow
pages may live below that domain at `workflows/<workflow>.md`.

## Selecting domain context

For a domain request, first read this index and then only the named domain page
and the workflow pages it links to. Do not recursively read every domain or the
raw source inbox. If the request spans domains, select the smallest set that
covers the stated outcome and record the selection in the plan or handoff.

Domain pages can be `proposed`, `accepted`, or `needs-review`. Proposed pages
are useful evidence but are not accepted Product Knowledge until a human gate
is satisfied. A refresh of an accepted page must preserve accepted decisions;
conflicting evidence is surfaced for review rather than silently replacing it.

## Domain entries

- [Workspace orientation and repository binding](repository-binding/README.md) — orientation, portable identity, host-local bindings, fail-closed resolution. Status: accepted.
- [Intent and Gate 1](intent/README.md) — the first-class intent front door: goal, outcome criteria, coarse optional scope, provisional tier; approval (Gate 1) freezes the contract and confirms understanding. Status: accepted.
- [Tracing and feasibility](tracing/README.md) — post-approval read-only tracer (one child per repository) reports a trace manifest; the coordinator runs the feasibility check before writing plans. Status: accepted.
- [Planning and plan review](plan-review/README.md) — grounded plan authoring and non-executing review. Status: accepted.
- [Plan authorization](plan-authorization/README.md) — a plan is authorized by its approved intent (scope-free: approved intent + unchanged criteria); no separate plan gate. Status: accepted.
- [Plan execution](plan-execution/README.md) — one-worker loop, isolated worktrees, repair, preserved runtime records. Status: accepted.
- [Direct collaboration](direct-collaboration/README.md) — `cc-pair` as the Explore tier of the assurance ladder: live human-supervised work in one repository with a coordinator and one worker, no verifier, promotable in place. Status: accepted.
- [Run-stack](run-stack/README.md) — executing a set of intent-authorized plans in one run: inter-plan dependencies, path leases, execution bases, scheduling, failure containment, drift guard. Status: accepted.
- [Repository grounding](repository-grounding/README.md) — the worker honors the target repository's own agent guidance, discovered live and delivered via a generated brief. Status: accepted.
- [Assurance and consequence tiering](assurance/README.md) — the Explore/Standard/Critical ladder and its verifier floor; the one safety-critical automated check, model-blind and fail-upward. Status: accepted.
- [Verification](verification/README.md) — independent read-only check (Standard/Critical), candidate-bound; sole authority for `verified`. Status: accepted.
- [Completion](completion/README.md) — human-gated completion to `done` (`draft → done`) and knowledge reconciliation. Status: accepted.
- [Plan organization](plan-organization/README.md) — status-agnostic archive and restore. Status: accepted.
- [Delivery](delivery/README.md) — separate pull-request/merge/push actions; source/target branch model. Status: accepted.
- [External surface](external-surface/README.md) — manually-triggered, config-driven publishing of Context Circuit data (a plan and its tasks, or a plan's open-question thread) to external systems, orthogonal to the core workflow; no runtime/network code, no new authority. Status: accepted.
- [Host adapters](host-adapters/README.md) — Codex, Claude Code, and Cursor Agent CLI as evidence-only transports. Status: accepted.
- [Source release and upgrade](source-release-and-upgrade/README.md) — assembly, ships-vs-never-ships, upgrade preservation + migration-needed. Status: accepted.
- [System-design authoring](system-design-authoring/README.md) — the `cc-system-design` skill for authoring a system design as structured source material. Status: accepted.
