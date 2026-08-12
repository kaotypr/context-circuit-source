# Overview

## Summary

Add a lean, generated, and maintained Product Knowledge capability that gives humans and AI agents task-relevant business, role, domain, and workflow context without loading or modeling the entire application.

## Source

idea: Human-approved Product Knowledge design from the Context Circuit RFC discussion on 2026-08-12

## Affected repositories

- context-circuit (the Context Circuit framework source itself; Product Knowledge is developed in-place by maintainer commits, not run-task worktree fan-out).

## Assumptions

- Canonical Product Knowledge is stored as human-reviewable Markdown in the initialized Context Circuit wrapper and versioned with Git.
- PROJECT.md remains the compact product map, each business role has one file under context/roles, and each domain owns a summary plus detailed workflow pages.
- Role pages describe cross-domain goals and end-to-end role stories, while workflow pages remain authoritative for exact current behavior, variations, and business rules.
- Product Knowledge describes current effective behavior; plans describe proposed behavior, and implementation observations cannot silently redefine business truth.
- The existing gather-context workflow remains read-only toward canonical context and may emit runtime candidates, gaps, contradictions, and compact context packages.
- The first implementation remains provider-neutral and accepts host-resolved source references rather than embedding connectors or credentials.
- Existing initialized wrappers remain valid and may adopt Product Knowledge incrementally without documenting the complete product first.

## Open questions

Both original open questions were resolved with the human before approval:

- Execution ownership: the Context Circuit framework source repository itself owns execution. It is registered in `workspace.yaml` as `context-circuit` (path `.`, mode `submodule`, agent `repository-worker`, default branch `main`), and every PKNOW work item targets that key. Implementation proceeds as direct maintainer commits on this branch, matching the workflow-continuity precedent, rather than run-task worktree fan-out. The shipped neutral template still ships `repositories: {}` because the template build neutralizes the repositories node.
- Effective-state confirmation authority: a dedicated, workspace-configured confirming role — distinct from the per-page `owners` and from the merge-gate human — declares business workflows effective during Product Knowledge synchronization (PKNOW-060).

- None outstanding.
