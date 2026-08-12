# Overview

## Summary

Add a lean, generated, and maintained Product Knowledge capability that gives humans and AI agents task-relevant business, role, domain, and workflow context without loading or modeling the entire application.

## Source

idea: Human-approved Product Knowledge design from the Context Circuit RFC discussion on 2026-08-12

## Affected repositories

- None identified.

## Assumptions

- Canonical Product Knowledge is stored as human-reviewable Markdown in the initialized Context Circuit wrapper and versioned with Git.
- PROJECT.md remains the compact product map, each business role has one file under context/roles, and each domain owns a summary plus detailed workflow pages.
- Role pages describe cross-domain goals and end-to-end role stories, while workflow pages remain authoritative for exact current behavior, variations, and business rules.
- Product Knowledge describes current effective behavior; plans describe proposed behavior, and implementation observations cannot silently redefine business truth.
- The existing gather-context workflow remains read-only toward canonical context and may emit runtime candidates, gaps, contradictions, and compact context packages.
- The first implementation remains provider-neutral and accepts host-resolved source references rather than embedding connectors or credentials.
- Existing initialized wrappers remain valid and may adopt Product Knowledge incrementally without documenting the complete product first.

## Open questions

- This neutral source checkout has no registered product repository, so which configured wrapper or repository registration will own execution of the approved implementation plan?
- Which workspace-specific human roles may confirm business workflows and declare implemented behavior effective during Product Knowledge synchronization?
