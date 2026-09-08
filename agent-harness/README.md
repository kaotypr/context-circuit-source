# Cc template-runtime laboratory

This source-only test area checks how the assembled `context-circuit-template`
works as a universal project workspace. It is deliberately separate from the
source repository's own ignored `.runtime/`, which stores execution state for
plans about the Cc source project.

Each scenario should exercise the built template through realistic v1.0.0
project prompts and small fixture repositories. A scenario may define:

- the project prompt or conversation;
- one or more fixture repositories and their starting revisions;
- starting Product Knowledge and source evidence;
- expected plan details and repository mappings;
- intent authoring and Gate 1 approval, automatic plan derivation, execution,
  verification, repair, and failure-limit outcomes;
- candidate-bound acceptance, separate Gate 2 delivery, tier-specific completion,
  host-blocked outcomes, and preserved runtime evidence.

The harness must initialize an isolated temporary workspace from the assembled
template. It must not import the source repository's Product Knowledge, plans,
`.runtime/`, credentials, or implementation state. Generated workspaces,
worktrees, commits, and runtime records are disposable test artifacts.

The v1.0.0 scenario set covers safe orientation, repository connection on the
human's selected base branch, the intent front door, post-approval tracing with
feasibility questions from the real code, refusal before Gate 1, an out-of-scope
reach surfaced as a feasibility question, schema-3 plan
derivation, Standard and Critical assurance, independent verification, worker
repair, candidate staleness, host blocking, run-stack ordering, repository
grounding, Explore/direct collaboration, archive/restore, Standard inferred
completion, reconciliation debt, change-set integration, and the separate
delivery boundary. Legacy branch migration remains intentionally covered only
by the semantic runtime suites.

This directory belongs to `context-circuit-source` and is excluded from
`context-circuit-template` assembly.
