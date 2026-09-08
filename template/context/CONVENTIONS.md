# Conventions

Status: uninitialized

Record accepted language, testing, naming, review, and delivery conventions.
A convention becomes Product Knowledge only when it is demonstrated consistently
and accepted, not from a single observation.

Live context files hold durable product knowledge only (INV-KNOWLEDGE-03). They
never name a particular plan, intent file, or sources file, and never name a
path into the sources tree. Path patterns that explain product structure (for
example `intent/<id>/`) remain allowed. Invariant IDs and shipped contract or
adapter paths remain nameable.

`DECISIONS.md` records what is now true about the product — decision, rationale,
and consequence — not edited paths or the ephemeral artifact behind the change.

Owned vs. consumed: a wrapper/gateway repository's own behavior stays in
`context/domains/`, tagged by `repositories`. Knowledge about an external service
the workspace consumes but does not own lives under `context/references/`, one
sub-directory per service, with no `repositories` tag.
