# Role Knowledge

`context/roles/<role>.md` is the canonical home for one business or project
role's cross-domain perspective. It links to canonical Domain Knowledge and
explains outcomes, surfaces, journeys, handoffs, role-specific behavior, and
limitations without becoming a second copy of domain facts.

Generated role pages use [the role template](../../docs/templates/role-context.md).
Their metadata has a simple `domains` list containing links to relevant domain
pages. It is not a relationship matrix and does not describe technical access
control.

## Selecting role context

For a role request, read this index, the requested role page, and only the
domain pages and workflow pages linked from that role. An unrelated session
should not scan all role, domain, or source files. If a role page is missing,
the agent should report the gap and gather only the evidence needed to create
it.

Role pages can be `proposed`, `accepted`, or `needs-review`. Human acceptance
is required before proposed material becomes accepted Product Knowledge. A
contradictory refresh remains visible for human review.

## Role entries

- [Product-source maintainer](maintainer.md) — ships the wrapper and seed; does not treat this checkout as a customer workspace. Status: proposed.
