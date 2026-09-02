# Context Circuit v1.0.0

Source design for the v1.0.0 template-harness host matrix. Read the preceding
v0.7.0 design for the direct-collaboration and execution-latency foundations;
this grouping defines the shared role-tiering fixture and concurrent
Codex/Claude Code/Cursor Agent scenario runner.

## Scopes

- [template-harness/](./template-harness/) — one host-grouped role-tiering
  fixture, host-specific live bindings, and a three-lane scenario matrix.

## Authority

This is passive source material. It adds no lifecycle gate, runtime policy, or
new invariant. The shipped role-tiering contract remains owned by
`docs/role-tiering.md`; host behavior remains bounded by `INV-HOST-01`,
`INV-RUNTIME-01`, and the verifier independence rules.
