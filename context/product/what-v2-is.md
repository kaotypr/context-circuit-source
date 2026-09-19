# What Context Circuit v2 is

Context Circuit is a shared workspace station for AI-assisted development across
one or more Git repositories, for solo developers and teams. A workspace holds
the project's identity, its logical repositories and their relationships, its
members, its durable knowledge, and the approved intents and derived plans that
change it.

The reason the product exists is the knowledge circuit: accepted understanding
grounds the next change, and what that change durably altered returns to the
knowledge so the change after it starts from truth instead of rediscovery.
Everything else is machinery serving that loop — see
[durable notes and the retrieval catalog](../domains/knowledge-notes.md).

## The v2 thesis

v2 is a rewrite, not an evolution of v1. Its governing claim:

> A system's own records prove only that the system wrote them. The branch, the
> diff, the passing check, and what a person actually said are what can be
> trusted.

Three consequences shape the whole product:

1. **Git is authoritative; workspace notes are notes.** A progress note records
   what a session believed. On resume the real branch and diff are read first
   and the note second.
2. **Mechanism splits from judgment.** Work that must happen identically every
   time lives in a compiled CLI; everything interpretive stays with the agent,
   in the conversation, where a person can see it. See
   [the CLI and agent seam](../architecture/cli-and-agent-seam.md).
3. **Nothing claims what has not been established.** A returned subagent
   specification says it launched nothing; a written host role file is not a
   loaded one; derived execution order says it ran, merged, and reserved
   nothing; the diagnostic is labelled a diagnostic.

## Two decisions belong to a person

Approving the intended outcome, and authorizing any outward action. Consent
cannot come from a command, an editable note, or another agent. Everything
between those two points is ordinary work done in the open — there is no second
approval between a plan and its execution. See
[authorization boundaries](../domains/authorization.md).

Whether that path applies at all is also theirs. A change whose outcome is
already the request can be made [directly](../domains/direct-changes.md), with no
intent, no plan, and no worktree, and with every outward gate still in force.

## What v2 deliberately does not have

Consequence tiers, frozen criteria digests, candidate identities, path leases,
execution or verification or host-evidence records, automatic repair loops,
risk-triggered review, an external publication surface, mandatory delegation, or
a shipped model catalog. Each was judged to cost more in ceremony and false
confidence than it returned, with the concern now carried by an artifact a
person can inspect directly.

v2 also does not migrate a v1 workspace. Initialization is for fresh workspaces;
a v1 workspace keeps working with v1 and is never rewritten in place.
