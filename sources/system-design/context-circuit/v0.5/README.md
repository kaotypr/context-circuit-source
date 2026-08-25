# Context Circuit v0.5 design set

This folder contains the complete v0.5 source design. The normative overview is
[context-circuit-v0.5-design.md](./context-circuit-v0.5-design.md). The other
files are detailed design modules. They may clarify the overview but must not
introduce a competing lifecycle, router, authority, or policy owner.

## Reading order

1. [context-circuit-v0.5-design.md](./context-circuit-v0.5-design.md) — product
   definition and fixed decisions.
2. [01-agent-behavior.md](./01-agent-behavior.md) — agent behavior for
   orientation, context work, planning, execution, verification, repair, and
   completion.
3. [02-human-interaction.md](./02-human-interaction.md) — conversational
   requests, approval, refusals, status reporting, and human decisions.
4. [03-workspace-and-repositories.md](./03-workspace-and-repositories.md) —
   workspace structure, repository registration, bindings, branches, and
   worktrees.
5. [04-context-lifecycle.md](./04-context-lifecycle.md) — how Product
   Knowledge is gathered, proposed, accepted, refreshed, and used.
6. [05-planning-and-execution.md](./05-planning-and-execution.md) — plan
   structure, execution brief, worker/verifier prompts, commits, and repair.
7. [06-runtime-engine.md](./06-runtime-engine.md) — runtime responsibilities,
   records, host adapters, safety checks, and the engine refactor.
8. [07-examples-and-acceptance.md](./07-examples-and-acceptance.md) — concrete
   files, conversations, traces, and acceptance scenarios.
9. [08-terminology.md](./08-terminology.md) — controlled vocabulary and
   deprecated product terms.
10. [09-source-and-template.md](./09-source-and-template.md) — source
    repository, product template, and instantiated workspace boundaries.
11. [10-template-publication.md](./10-template-publication.md) — publishing the
    template repository via GitHub Actions: its independent semantic version,
    the `release/` ledger, one-commit-per-bump history, tags, and releases.

## Authority

The overview owns product decisions. The plan file owns plan intent and human
status. Product Knowledge owns accepted project facts. The runtime owns
execution evidence and repository safety. The host adapter owns provider-native
child creation. The worker owns implementation commits. The verifier owns
independent verification results.

A detail module must point back to the owner when it describes a rule. It must
not copy a rule into a second canonical form.

## Source boundary

These documents are maintainer design material. They are not copied into an
instantiated workspace as Product Knowledge and are not normal agent context.
An implementation session should receive only the relevant contracts and detail
needed for its bounded change.
