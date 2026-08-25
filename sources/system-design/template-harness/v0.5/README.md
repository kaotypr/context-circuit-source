# Context Circuit v0.5 template test design

Status: maintainer design for a human-simulated product test harness
Related design: `sources/system-design/context-circuit/v0.5/`
Related runtime laboratory: `template-harness/` (deterministic engine-level)

This design defines how to test the **built** `context-circuit-template` as a
real project workspace, driven by natural human conversation rather than by
calling the runtime directly. Because the product is tri-host, the same test
runs against each host adapter — Codex CLI, Claude Code, and Cursor Agent — with
a host-neutral scenario library and a thin per-host binding (see
[host-matrix.md](./host-matrix.md)).

It introduces a maintainer-level **human-simulator** agent role that behaves like
an ordinary person who has never heard of Context Circuit, workspaces, plans,
worktrees, or the runtime. The human-simulator spawns the product **coordinator**
as its sub-agent inside a fresh, isolated workspace instantiated from the
assembled template, sends it plain prompts, reacts like a human, and — together
with a separate deterministic grader — judges whether the product experience is
correct.

## Why this is separate from the existing laboratory

`template-harness/` already assembles the template and exercises the
**runtime boundary** deterministically (it calls `cc_*` functions). It proves the
mechanism works. It does not prove that a real conversation with a person who
knows nothing about the internals produces the intended product experience.

This design adds that missing layer: a **semantic, conversational** test at the
human/product boundary. It answers "does the assembled product behave well when a
normal human talks to it?" — orientation, initialization, grounded planning,
plan review, the approval gate, refusal of unapproved execution, compound
approve-and-execute, human completion, archive/restore, and the delivery
boundary — all expressed only in natural language.

## Maintainer-only

The human-simulator role, the scenario prompt libraries, the grader, and every
generated workspace are **source-only maintainer material**. None of them ship in
`context-circuit-template`. The product still has exactly three roles:
coordinator, worker, verifier. The human-simulator is a test actor, not a product
role.

## Reading order

1. [roles/README.md](./roles/README.md) — the actors and their interaction
   contract, with a dedicated design for each:
   - [roles/human-simulator.md](./roles/human-simulator.md) — the maintainer
     agent that plays an ignorant human.
   - [roles/coordinator-under-test.md](./roles/coordinator-under-test.md) — the
     product coordinator as configured and constrained in the test context.
   - [roles/grader.md](./roles/grader.md) — the mechanical judge (state,
     access-discipline, transcript, and efficiency dimensions).
2. [scenario-library.md](./scenario-library.md) — the prompt libraries under
   `test/`, the case-file format, and the enumerated cases starting with a
   brand-new, no-repository, simple-idea project.
3. [harness-and-evaluation.md](./harness-and-evaluation.md) — how a run assembles
   and isolates a workspace, the turn protocol, grading, result artifacts, and
   boundaries.
4. [host-matrix.md](./host-matrix.md) — how the same harness runs the product
   under each host adapter (Codex CLI, Claude Code, Cursor Agent): the host axis,
   what stays host-neutral, the per-host bindings, and host-specific expectations.

## Authority

This is maintainer test design. It does not change the product contract in
`sources/system-design/context-circuit/v0.5/` or `wrapper/contracts/`. Where it observes
product behavior, the product design and `wrapper/contracts/invariants.yaml`
remain the authority for what "correct" means.
