# Conventions

Use the smallest route-selected context. Name exact source files before reading
them and record provenance. Keep one normative owner per rule and use invariant
IDs in prose and tests. Treat `plan.yaml` as canonical status and `PLAN.md` as
the human explanation.

Preserve dirty state and runtime evidence. Never store credentials or infer
authorization from eligibility, tests, Git state, or provider status. Writers
and verifiers stay isolated, and human-gated actions remain separate.
