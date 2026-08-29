# Conventions

Status: uninitialized

Record accepted language, testing, naming, review, and delivery conventions.
A convention becomes Product Knowledge only when it is demonstrated consistently
and accepted, not from a single observation.

Owned vs. consumed: a wrapper/gateway repository's own behavior stays in
`context/domains/`, tagged by `repositories`. Knowledge about an external service
the workspace consumes but does not own lives under `context/references/`, one
sub-directory per service, with no `repositories` tag.
