# Product Knowledge

Product Knowledge is concise, source-cited project context under `context/`.
It describes durable intent, users, domains, workflows, architecture,
conventions, decisions, and known gaps.

A typical tree contains:

- `PROJECT.md`;
- `ARCHITECTURE.md`;
- `CONVENTIONS.md`;
- `DECISIONS.md`;
- `roles/<role>.md` for business or project perspectives;
- `domains/<domain>/README.md` for canonical bounded-area knowledge;
- `domains/<domain>/workflows/` for exact workflow behavior;
- `SOURCES.md` and `sources.yaml`.

Product Knowledge is not an instruction layer and does not override
`AGENTS.md`, `WORKFLOW.md`, repository-local instructions, approved plans, or
human decisions.

When an agent receives a PRD or other source, it should:

1. Read and cite the source.
2. Separate requirements from current implementation.
3. Draft concise, useful context.
4. Present changes for human confirmation.
5. Record provenance and source revision.
6. Refresh only through a human-reviewed proposal when the source changes.

## Domain Knowledge

Domain Knowledge is canonical for one bounded project area. It may describe
purpose, scope, behavior, workflows, interfaces, data, constraints, edge
cases, implementation references, and verification when the selected evidence
supports those categories. Use `docs/templates/domain-context.md`, store the
result under `context/domains/<domain>/`, and keep exact workflow steps in the
linked workflow pages.

Generated domain documents record `status`, selected `sources`, source
revisions, `freshness`, `assumptions`, `unknowns`, `contradictions`, and an
`acceptance` state. They remain proposed until a human accepts them. A refresh
of accepted knowledge preserves accepted decisions and makes conflicting
evidence visible for review.

## Role Knowledge

Role Knowledge is a cross-domain view for one business or project role. Use
`docs/templates/role-context.md` and store the result at
`context/roles/<role>.md`. Its only domain relationship metadata is the simple
`domains` list; the body links to canonical domains and explains the role's
perspective, outcomes, surfaces, handoffs, and limitations without duplicating
domain facts.

Business/project roles are not agent execution roles. Coordinator, worker, and
verifier instructions remain under `agents/`, and role knowledge must not be
used as a technical access-control policy.

## Generation and selection

`cc-gather-context` handles request-scoped creation and refresh. It identifies
the requested domain or role, states the exact raw sources and repository
documents it will read, reads only those selections, and records provenance in
the result and `context/sources.yaml`. It updates the canonical existing
document instead of creating duplicates. A role request loads the role page
and only its relevant domain references; an unrelated session does not scan
the complete context tree or source inbox.

Raw sources remain in the passive `sources/` inbox. A source-based request
names and reads only the selected files, states why they will be read, and
records provenance in `context/sources.yaml` or the product artifact. Normal
session entry does not ingest the source inbox.
