# Risks

## Risks and mitigations

- Product Knowledge may become too large for efficient agent context; mitigate with a short product map, role and domain routing, bounded context packages, and a large-application retrieval fixture.
- Role stories and workflow pages may duplicate or contradict each other; mitigate with explicit ownership rules, link validation, and workflow authority for detailed behavior.
- Generated observations from code may be mistaken for intended business policy; mitigate by keeping discovery read-only and requiring human confirmation before canonical synchronization.
- Knowledge can become stale after product changes; mitigate by requiring impact declarations in plans and planless tasks and reviewing knowledge deltas during verification and closeout.
- Updating current knowledge immediately after merge may document behavior that is not released or enabled; mitigate with explicit effectiveness confirmation.
- Mandatory documentation could block adoption for mature products; mitigate with optional incremental initialization and just-in-time workflow discovery.
- Product Knowledge may contain sensitive operational information; mitigate with credential and secret rejection, source references rather than copied restricted content, and existing repository access controls.
- Schema and workflow expansion may break existing wrapper archives; mitigate with compatibility fixtures, additive defaults, migration guidance, and release-template validation.
- Host adapters may implement different discovery or review behavior; mitigate by keeping contracts and lifecycle rules canonical under .agents and verifying Codex and Claude parity.
