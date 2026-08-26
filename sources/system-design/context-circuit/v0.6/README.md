# Context Circuit v0.6

Source design for Context Circuit v0.6 — a **delta on v0.5**. Read the v0.5 design
first ([../v0.5/](../v0.5/)); v0.6 changes only what it names. This README is the
version index; each scope owns its own design.

## Scopes

- [run-stack/](./run-stack/) — execute a set of approved plans in one run (a *plan
  stack*): inter-plan dependencies, path leases, execution bases, the scheduling
  loop, failure containment, and the delivery drift guard. Start at
  [run-stack/design.md](./run-stack/design.md).
- [system-design-stage/](./system-design-stage/) — make a **System Design** a
  first-class, optional lifecycle stage between Product Knowledge and Plans: design
  the change, get it accepted, then let plans slice the accepted design. Start at
  [system-design-stage/design.md](./system-design-stage/design.md).
- [repository-grounding/](./repository-grounding/) — ground the **writer** in the
  target repository's own agent guidance (discovered live) and hand it a prepared
  worktree, delivered through a generated brief instead of hand-authored prose.
  Start at [repository-grounding/design.md](./repository-grounding/design.md).

These three scopes are independent capabilities released together in v0.6. They
share a coordinated contract bump (plan schema `[1, 2]`, `runtime_version 0.6.0`).

## Layout convention

`sources/system-design/<product>/<version>/<scope>/`. Every folder has a
`README.md` landing/index; each scope's normative design is `design.md`, with
detail split into files or sub-folders as it grows. See the v0.5 version README's
"Layout convention" for the full rule.

## Authority

v0.6 adds no owner and duplicates no rule. It specifies intent and points to the
canonical owners under `wrapper/`. Contract deltas are summarized in each scope's
`schema-and-*` file.
