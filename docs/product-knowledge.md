# Product Knowledge

Product Knowledge is concise, source-cited project context under `context/`.
It describes durable intent, users, domains, workflows, architecture,
conventions, decisions, and known gaps.

A typical tree contains:

- `WORKSPACE.md` — this workspace: identity after initialization and the
  registered-repository map as workspace structure, not the product being
  built;
- `PROJECT.md` — the project being built: purpose, boundaries, and what
  registered repositories contain, or a record that the project is not yet
  defined;
- `ARCHITECTURE.md`;
- `CONVENTIONS.md`;
- `DECISIONS.md`;
- `roles/<role>.md` for business or project perspectives;
- `domains/<domain>/README.md` for canonical bounded-area knowledge;
- `domains/<domain>/workflows/` for exact workflow behavior;
- `SOURCES.md` and `sources.yaml`.

`WORKSPACE.md` and `PROJECT.md` are separate identities. Workspace
behavior and the repository map belong in `WORKSPACE.md`. Whole-project
purpose belongs in `PROJECT.md`. There is no every-session product-direction
page.

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

## Document-system contract

The document-system foundation uses OKF v0.2 as an interoperability envelope
only for explicitly declared knowledge bundles. Product, Workspace, Domain,
Role, Architecture, Convention, Decision, Idea Brief, PRD, and similar
knowledge concepts may use the Context Circuit OKF profile. The profile is
stricter than base OKF and is reported separately.

The following remain outside OKF and retain their existing authority:

- `plan.yaml` and task contracts;
- `workspace.yaml` and provenance registries;
- `.runtime/` session, delegation, lease, prompt, evidence, handoff,
  completion, and stack records;
- `AGENTS.md`, `WORKFLOW.md`, skills, and other normative instructions;
- arbitrary raw files in the passive, request-scoped `sources/` inbox.

New knowledge writers emit `type`. Compatibility readers may normalize a
legacy `kind` only when `type` is absent; conflicting values fail the hard
check. Readers normalize a bare `verified` mapping to a one-item list and
preserve unknown extension fields. `verified` remains advisory and does not
satisfy a human acceptance, plan approval, status-change, publication, merge,
deployment, or cleanup gate.

Markdown front matter is validated during generation: the first standalone
`---` after line 1 closes the initial YAML block, later horizontal rules remain
body content, and the host-provided deterministic YAML/schema capability runs
before output is ready. No repository command runtime is added or silently
required. See `docs/document-system.md`, `docs/okf-profile.md`, and
`docs/host-capabilities.md`.

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
