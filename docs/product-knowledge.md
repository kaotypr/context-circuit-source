# Product Knowledge

Product Knowledge is concise, source-cited project context under context/. It
describes durable intent, users, domains, workflows, architecture, conventions,
decisions, and known gaps.

A typical tree contains:

- PROJECT.md;
- ARCHITECTURE.md;
- CONVENTIONS.md;
- DECISIONS.md;
- roles/;
- domains/<domain>/;
- domains/<domain>/workflows/;
- SOURCES.md and sources.yaml.

Product Knowledge is not an instruction layer and does not override AGENTS.md,
WORKFLOW.md, repository-local instructions, approved plans, or human decisions.

When an agent receives a PRD or other source, it should:

1. Read and cite the source.
2. Separate requirements from current implementation.
3. Draft concise, useful context.
4. Present changes for human confirmation.
5. Record provenance and source revision.
6. Refresh only through a human-reviewed proposal when the source changes.

Product Knowledge should be scoped to the relevant domain. Agents should load
the smallest useful set rather than copying the entire project into every
session packet.
