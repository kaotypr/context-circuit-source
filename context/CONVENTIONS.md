# Conventions

Use the smallest route-selected context. Name exact source files before reading
them and record provenance. Keep one normative owner per rule and use invariant
IDs in prose and tests. Treat `plan.yaml` as canonical status and `PLAN.md` as
the human explanation.

Separate owned knowledge from consumed external services. Knowledge the
workspace owns is organized by concept and scoped by repository: a
wrapper/gateway repository's own behavior stays in `context/domains/`, tagged by
`repositories`. Knowledge about an **external service the workspace consumes but
does not own** — a third-party API a dependency gates, its auth model, error
semantics, rate limits — is not a repository and carries no `repositories` tag;
it lives under `context/references/`, one sub-directory per service. A wrapper
repo is not its service: keep the two apart.

Preserve dirty state and runtime evidence. Never store credentials or infer
authorization from eligibility, tests, Git state, or provider status. Writers
and verifiers stay isolated, and human-gated actions remain separate.

A policy change must not live only in a context page: an update that changes
repository identity, branch policy, a security boundary, or execution behavior
also updates the owning workspace or contract file through that owner's normal
action. Weigh evidence by scope: repository evidence shows what the code does,
not what the product should do; a plan is intent, not implementation; a verifier
confirms behavior, not a domain decision.
