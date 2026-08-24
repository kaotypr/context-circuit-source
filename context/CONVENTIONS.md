# Conventions

Use the smallest route-selected context. Name exact source files before reading
them and record provenance. Keep one normative owner per rule and use invariant
IDs in prose and tests. Treat `plan.yaml` as canonical status and `PLAN.md` as
the human explanation.

Preserve dirty state and runtime evidence. Never store credentials or infer
authorization from eligibility, tests, Git state, or provider status. Writers
and verifiers stay isolated, and human-gated actions remain separate.

A policy change must not live only in a context page: an update that changes
repository identity, branch policy, a security boundary, or execution behavior
also updates the owning workspace or contract file through that owner's normal
action. Weigh evidence by scope: repository evidence shows what the code does,
not what the product should do; a plan is intent, not implementation; a verifier
confirms behavior, not a domain decision.
