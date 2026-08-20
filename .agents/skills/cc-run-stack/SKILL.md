---
name: cc-run-stack
description: Execute a connected set of already-approved plans in one root session through a frozen runtime graph, parent-worktree bases, in-run joins, and progress resume.
---

# Run an approved plan stack

Use this skill when the user asks to execute a connected set of already-approved
plans, when session-entry or `cc-whats-next` identifies approved unimplemented
connected plans, or when `.runtime/stacks/<stack-id>/` can be resumed.

## Route reads

Use the `execution` manifest for a new member and the `resume` manifest when
continuing a stack in
`docs/agent-workspace-workflow.md#route-read-manifests`. This skill retains the
stack-specific guard that `graph.yaml` and `progress.yaml` are runtime
authority for the run, while canonical plan status, leases, verifier isolation,
and human gates remain governed by the shared owners.

Invoking the skill starts or resumes execution. There is no stack-approval
gate. Refuse draft member plans, an invalid DAG, or a live foreign stack or
plan lease.

`cc-run-stack` is not a one-plan replacement for `cc-run-plan`. A one-plan
request still enters `cc-run-plan`. Do not treat a stack run as one
`cc-run-plan`. This skill is not a scheduler, daemon, or queue.

## Preflight

Read the named member `plan.yaml` files, declared prose dependencies, Product
Knowledge, repository-local instructions, active sessions, leases, worktrees,
and any existing `.runtime/stacks/<stack-id>/graph.yaml` and `progress.yaml`.
For member resume, prefer `handoff.yaml` and use historical Markdown-only
`handoff.md` only when no structured handoff exists.
Confirm:

- every member plan status is exactly `approved`;
- the interpreted DAG is acyclic and IDs resolve;
- no live foreign stack lease or plan lease blocks the run;
- the requested work remains within each member's approved scope.

Each member is read through its canonical `plan.yaml` first. A member may use
the compact `plan.yaml`/`overview.md`/`tasks/` bundle or a historical bundle
with specialist companions; companion presence never changes the frozen
dependency graph or lifecycle authority. New task Markdown must have passed
the approved host-provided deterministic front-matter and task-schema checks
before readiness. Validation is a document-readiness check only and cannot
approve, execute, verify, or complete a member plan.

A draft member, unknown plan, cycle, or ownership conflict is blocked. Do not
create `plans/<repository-key>-stacks/`, a durable `stack.yaml`, or a hidden
replacement plan. Do not change canonical status to make execution possible.
Runtime must not override `plan.yaml`.

## Claim, freeze, and loop

On start, create `.runtime/stacks/<stack-id>/` with no global current-stack
pointer. Freeze `graph.yaml` once: members, directed edges, and leaves. Claim
the stack lease under `lease.lock/owner.yaml` and `lease.yaml`. Update
`progress.yaml` as members move: `pending`, `waiting-parents`, `ready`,
`running`, `joining`, `implemented`, `failed`, or `blocked`, plus `frozen_sha`,
`worktree`, `base`, and `join_parents`.

Follow `cc-run-plan` internals per ready member: lease, exclusive worktree,
writer child, independent verifier child. Sequential tasks inside one member
still share one writer child. Independent ready members get separate children
and worktrees. Change the worktree base, continue across the DAG, freeze
`graph.yaml`, and update `progress.yaml`.

```text
Root session
  -> interpret approved-plan dependencies into graph.yaml
  -> claim stack lease
  -> compute ready frontier from graph.yaml + progress.yaml
  -> for each ready member:
       claim plan lease
       create exclusive worktree from parent frozen SHA or default branch
       writer child, then independent verifier child
       repair within scope until verifier passes or hard-stop
       require a local commit, freeze SHA, mark implemented in progress.yaml
  -> recompute frontier (including joins)
  -> stop when every member is implemented
  -> hand leaf worktrees; do not run cc-finish-plan
```

Spawn each child through the host child-session primitive. Cursor's
Task/subagent tool is a valid primitive. If the host cannot spawn a child,
report the missing host primitive to the human and ask how to proceed. Do not
quietly skip children.

This skill has a no-finish-plan stop: it must not run `cc-finish-plan` and
must not write `plan.yaml` status `done`. Do not push, merge, publish, or deploy.

## Parent-worktree bases and in-run joins

Each member keeps `.runtime/worktrees/<repository-key>/<plan-id>/`. Use
parent-worktree bases and in-run joins:

1. No parent: `git worktree add` from the repository default or active branch
   (same as `cc-run-plan`).
2. One parent: new branch at the parent frozen SHA.
3. Several parents: sort parent IDs, add from the first frozen SHA, merge the
   remaining SHAs in that order, then the writer runs. Routine joins do not
   pause for a human and do not wait for `default_branch`.

The new worktree is a Git checkout of that commit-ish, not a copy of the parent
directory. Writable worktrees stay exclusive.

Require a local commit on the exclusive member branch before freeze. That
commit is stack-internal. It is not merge to the default branch and not a
delivery authorization.

## Implemented-versus-done

Canonical plan status stays `draft`, `approved`, and `done`. Record
implemented-versus-done:

- Implemented is runtime: writer finished, independent verifier passed,
  worktree HEAD committed and clean,
  `.runtime/plans/<plan-id>/completion.yaml` is
  `ready-for-human-status-change`, and `plan.yaml` remains `approved`.
- Freeze the parent SHA on `progress.yaml`.
- Dependents wait on implemented parents, not `done`.
- Verifier findings are repaired in-session without human interaction, within
  that member's approved scope, until the verifier passes or a hard stop.

Do not rebase dependents if a human later wants changes; that later work is a
new plan.

Motivating ready frontier:

```text
0001 implemented → 0002, 0005, 0006
0002 implemented → 0003
0003 implemented → 0004
0004 implemented → 0007, 0008
0007 implemented → 0009
0005 and 0006 implemented → 0010 (join)
```

Leaves 0008, 0009, and 0010 are merge units.

Hard-stop and hand off when: a change is out of scope; a join cannot be merged
without leaving approved scope; a host child-session primitive is missing;
ownership or lease conflicts; or verifier repair would require a material
scope change.

## Resume

The same session resumes from `progress.yaml`: skip `implemented` members,
continue the frontier, reuse existing exclusive worktrees. Do not rebuild
`graph.yaml`.

A new session reads `graph.yaml` and `progress.yaml` and must not rebuild a
different tree or recreate implemented worktrees. A live stack lease blocks
silent takeover. If takeover is authorized, record `replaced_session_id` on
`progress.yaml` and become owner.

Progress resume uses those runtime files. After `cc-cleanup-runtime` deletes
`.runtime/`, there is no stack to resume; a new run writes a new graph.

## Output

Report the stack, frozen graph, progress cursor, delegated members, worktree
bases and joins, files changed, tests and verification, decisions, assumptions,
blockers, and next safe action. When every member is implemented, hand the
human the leaf worktrees as merge units and name per-plan `cc-finish-plan`.
This skill must not mark plans done.
