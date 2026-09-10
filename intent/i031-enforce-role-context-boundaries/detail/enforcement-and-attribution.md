# Context enforcement and attribution

## Context dimensions

The context ledger keeps four dimensions separate:

- **available:** the resources technically reachable by the role;
- **provided:** content inserted into its initial or follow-up packet;
- **active read:** content requested through a tool or command;
- **charged:** host-reported input tokens or equivalent usage.

Each event records the role, resource identity, context kind, measured size,
normal or expanded authority, and host evidence reference.

```yaml
role: planner
kind: active-read
resource: repository:reporting/src/export.ts
bytes: 4821
authority: normal-planner-envelope
evidence: tool-event-184
```

## Enforcement preference

The strongest driver gives a role a filesystem or connector view containing only
its declared envelope. Where a native host requires a broader shared workspace,
the driver uses structured host events and supported process-tree or
operating-system access tracing to attribute reads.

Strict grading requires either enforced unavailability of forbidden resources or
positive role-attributed evidence. A host with neither produces an inconclusive
boundary result.

## Host activity

Background indexing is recorded separately when the host exposes it. It counts
as agent context only when evidence shows that indexed content entered the role's
session. If the driver cannot distinguish background access from role access, it
does not guess.

The root and every native child have separate ledgers. Flattening all activity
into one session would prevent the scenario from proving role correctness and is
therefore insufficient for a strict pass.

