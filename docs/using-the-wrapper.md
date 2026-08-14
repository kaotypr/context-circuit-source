# Using the workspace

Context Circuit is a workspace for AI agents, not a command console for users.

## Shared durable context

The workspace keeps Product Knowledge, decisions, source provenance, plans, and
task definitions in ordinary Markdown and YAML. These files are reviewable and
should contain only durable project truth or human-reviewed intent.

## Private runtime state

The .runtime directory keeps local execution state:

- root and child session records;
- parent-child delegation handoffs;
- plan leases and heartbeats;
- generated context and delegation packets;
- prompts;
- isolated repository worktrees.

Runtime state is resumable but is not authoritative over context, plans,
human decisions, or repository instructions. Preserve it until a human chooses
cleanup.

## Repository boundaries

Product repositories own code and repository-local conventions. A writing
session modifies only its assigned worktree. Different plans may use different
worktrees concurrently. Research and verification sessions are read-only unless
their delegation explicitly grants write access.

## Human control

Agents may inspect, draft, test, delegate, create isolated worktrees, and
implement approved scope. Humans control Product Knowledge acceptance, plan
approval, material scope changes, merge, publication, deployment, completion,
and ambiguous session takeover.

## Host adapters

Codex, Claude, and other host integrations should expose the same workspace
entry and delegation behavior. Host-specific commands are adapters and must not
become a second workflow or source of truth.
