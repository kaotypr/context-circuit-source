# Intention — i030-integrate-real-host-sessions

_Status: approved, look complete, feasible.

## Intention

Give Agent Harness **real Codex, Claude Code, and Cursor session drivers** so its
scenarios exercise genuine host behavior rather than replayed transcripts or
mocked agents.

Each driver follows one shared run lifecycle while preserving the host's native
session, child-agent, event, authentication, usage, and termination behavior.

## Expectations

- A scenario can launch a new real session on each supported host.
- Native child agents remain attached to their real root session and are visible
  as distinct roles.
- Host differences and unavailable measurements remain visible instead of being
  normalized into misleading results.
- Existing host authentication is used safely without copying credentials into
  test workspaces or evidence.
- The same scenario can run on one selected host or across the supported-host
  matrix without duplicating its definition.

## The plans

1. **Integrate real coding-agent hosts.**
   _After this:_ Agent Harness can probe, launch, observe, wait for, terminate,
   and collect evidence from genuine Codex, Claude Code, and Cursor sessions.

## How carefully this is checked

**`Standard`**

Explanations:
- **Explore:** you check it yourself as you work alongside the agent — no separate
  verification. Meant for trying things out.
- **Standard:** a second, independent agent verifies the finished work against your
  definition of done before it's called done.
- **Critical:** the strictest — independent verification plus a repair-and-recheck
  loop. For anything risky or impossible to undo (money, security, production, data
  you can't get back).

## Open questions

_At draft time: no known unresolved human decisions._
