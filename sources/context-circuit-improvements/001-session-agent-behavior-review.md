# 001 — Session agent-behavior review

This is a passive source artifact. It records what agents actually did in one
product-source Cursor session that walked plan `local-binding-troubleshoot`
(CC-003) through approve → execute → finish → delivery, then reviewed the
trace. It is not accepted Product Knowledge, not a canonical workspace plan,
and not authorization to modify `wrapper/`, `template/`, skills, or runtime.

Drafted 2026-08-22 from named session evidence only. Do not scan the rest of
`sources/`.

## Provenance

| ID | Path or locator | Why it was read |
| --- | --- | --- |
| CC-003 | `plans/context-circuit-plans/local-binding-troubleshoot/plan.yaml`, `PLAN.md`, `tasks/001-getting-started-troubleshooting.md` | The plan that was approved, executed, finished, and delivered |
| RT-ROOT | `.runtime/sessions/sess-local-binding-troubleshoot-20260822-01/` | Root coordinator session, receipt, and handoff |
| RT-WRITER | `.runtime/sessions/child-writer-local-binding-troubleshoot-20260822-01/` | Writer packet, session, and handoff |
| RT-VERIFIER | `.runtime/sessions/child-verifier-local-binding-troubleshoot-20260822-01/` | Verifier packet, session, and handoff |
| RT-LEASE | `.runtime/plans/local-binding-troubleshoot/lease.yaml` and `completion.yaml` | Lease left active; completion after finish |
| WT | `.runtime/worktrees/context-circuit-source/local-binding-troubleshoot` | Exclusive writer worktree on `codex/local-binding-troubleshoot` |
| PR-59 | https://github.com/kaotypr/context-circuit/pull/59 | Delivery of the docs change targeting `codex/development/v0.5` |
| HOST | Cursor Agent session, 2026-08-22 | Human turns, coordinator behavior, Task children `06f590e4` (writer) and `faeab665` (verifier) |
| ENGINE | `wrapper/runtime/engine.sh`, `wrapper/contracts/routes.yaml`, `wrapper/contracts/context-sets.yaml` | Probe, worktree, and context-set mismatches against live behavior |
| ROLES | `agents/coordinator.md`, `agents/writer.md`, `agents/verifier.md` | Intended role boundaries |

Related but not copied here: Boardwalk improvement sources under
`sources/boardwalk-project-test/`, and done maintainer plans
`repository-bootstrap` (CC-001) and `multi-host-agent-support` (CC-002). This
note is about live agent cost and host mapping after those plans, observed on
CC-003.

## Thesis

The three roles are the right product. Isolation of writer and verifier held
in this session only as **prompt discipline**, not as a host or engine
constraint. The router and worktree primitives did not match what the
coordinator actually had to do, so the root agent worked around them.

A one-task docs plan took nine human turns. Children ran only once, on the
second execute request. Do not add a second router, scheduler, or lifecycle.
Make the existing host mapping and primitives match the roles that already
exist.

## Intended model

One two-stage router. Root coordinates. Writer is a bounded child in an
exclusive worktree. Verifier is a separate read-only child. Host Task and
subagent features are only the transport. Human gates are exact and
session-bound. Host capability is evidence, never authorization.

## Turn-by-turn trace

| # | Request | Who spoke | Who acted | What actually happened |
| --- | --- | --- | --- | --- |
| 1 | Orient / later approve-execute-finish then review agents | Human | Root | Engine said `draft-plan` / `PLAN_REQUIRED`. Coordinator overrode to recommend-next. |
| 2 | `Approve plan local-binding-troubleshoot` | Human | Root | Card only. No mutation. |
| 3 | `Confirm approval of plan local-binding-troubleshoot` | Human | Root | `plan.yaml` `draft` → `approved`; LBT-001 `draft` → `ready`. |
| 4 | `Run approved plan local-binding-troubleshoot` | Human | Root | Blocked: `MAINTAINER_APPROVAL_COMMIT_REQUIRED`. No lease, no writer. |
| 5 | `Confirm commit of the approved plan state.` | Human | Root | Commit `8e4ca0b` on `codex/development/v0.5`. Execution still not started. |
| 6 | `Run approved plan local-binding-troubleshoot` | Human | Root → writer → verifier | Lease, named-branch worktree, Cursor Task children, independent verifier pass. |
| 7 | `Confirm finish plan local-binding-troubleshoot` | Human | Root | Status-only `approved` → `done`. Lease left `active`. |
| 8 | Commit workspace changes, push, create PR targeting `codex/development/v0.5` | Human | Root | PR #59. Finish status forked onto development and the PR branch. |
| 9 | Review this session and improve from the agent trace | Human | Root | This source note. No new plan drafted. |

Children exist only in turn 6.

### Role scorecard

**Root / coordinator.** Owned routing, cards, lease, packets, and delivery.
Did not write `docs/getting-started.md` and did not self-verify. Also authored
the writer’s paragraph into the Task prompt, and used
`git worktree add -b` because `cc_prepare_worktree` cannot name `codex/…`
branches.

**Writer child** (`06f590e4`). Wrote only the assigned worktree, committed
`31501e019ca6e9c57e80ade319165e4963982201`, and handed off. The Cursor Task
workspace root was still the base checkout. Isolation was a prompt, not a
host sandbox.

**Verifier child** (`faeab665`). Reran LBT-VT-01 through LBT-VT-03, audited
declared paths, mapped LBT-AC-01 through LBT-AC-03, and refused to repair.
Host type was still `generalPurpose` (can write). Read-only was a prompt, not
a permission mode.

## What already worked

- Gates stayed separate: approve, execute, finish, and delivery were not
  inferred. Vague “yes” was rejected.
- Status-only transitions used `cc_transition_plan_status`.
- Root did not replace the writer or verifier.
- Implementation stayed in
  `.runtime/worktrees/context-circuit-source/local-binding-troubleshoot`.
- The verifier treated the writer handoff as a claim, not as truth.

## Findings

### P1 — Child isolation is prompt-only

Cursor Task inherits the base workspace. The writer was told not to touch it.
The verifier was told not to write. If the host cannot bind the writer to the
worktree root and the verifier to read-only, the honest result is
`host-blocked` (`INV-HOST-02`), not a long prompt.

### P1 — Router probe is substring-fragile

The first human message was a meta intent (“we will later approve, execute,
then review”). `cc_probe` / `cc_route` selected `draft-plan` with
`PLAN_REQUIRED` because the text was not classified read-only. “Approve some
plan” would not have matched `approve plan`. Route on named plan IDs and
explicit lifecycle verbs. Meta future-tense lifecycle talk must stay
orientation (`ROUTE-PRECEDENCE-05` / `06`).

`context_set` names emitted by the engine (`approval`, `execution-preflight`,
`plan-draft`) are not ids in `wrapper/contracts/context-sets.yaml`.

### P1 — Worktree primitive does not match real execution

`cc_prepare_worktree` uses `git worktree add --detach`. `cc_safe_id` rejects
`/`, so `codex/local-binding-troubleshoot` cannot pass as a branch. This
session created a named branch from clean HEAD with
`git worktree add -b`, matching prior maintainer executions
(`codex/repository-bootstrap`, `codex/multi-host-agent-support`). The engine
should create that named branch from the clean HEAD the coordinator actually
used (`INV-REPO-06`).

### P2 — Product-source commit gate surprises at execute

`INV-OWN-07` is correct: the exact approval projection is not an execution
exemption. Presenting the commit card only when execute is requested cost a
full extra human round trip (turns 4–5). Finish then forked the same status
commit: `6239305` on local `codex/development/v0.5` and cherry-pick `87d030c`
on the PR branch. Present the maintainer commit card immediately after
approval, or keep the status projection on the worktree branch only.

### P2 — Coordinator wrote the writer’s solution

Root pasted the `PLAN.md` paragraph into the Task prompt. The packet should
name plan and task paths. The writer should read them. Root should not
pre-author the edit.

### P2 — Runtime records went stale after finish

After confirm-finish:

- Root handoff still says `awaiting-human-gate` and that the plan remains
  `approved`.
- Root `session.yaml` was later set to `completed`, but the handoff was not
  refreshed.
- Lease `status` is still `active`; `released_at` is still `null`.
- Human-facing `PLAN.md` still says `Status: draft`.

Finish should refresh the handoff, decide lease release, and project
human-facing plan status. Completion evidence must not itself mark the plan
done (`INV-LIFE-02`); the human gate already did that.

### P3 — Weak verification lock

The optional `contains` assertion locks `'create the ignored root'`, which
also matches the older sentence “Each host may create the ignored root file”.
The verifier recorded this limitation and still passed. Nobody entered the
repair loop. Lock a unique phrase from the new paragraph. Failed uniqueness
should be bounded repair, not a footnote.

### P3 — Card text and required reply diverge

`cc_confirmation_card` ends with “Confirm this named action in the current
session.” Routing fixtures require exact phrases such as
`Confirm approval of plan <id>` and `Confirm finish plan <id>`. One owner
should print the exact confirmation phrase on the card.

## Time and what agents read

Counts are from this Cursor transcript
`6b906ead-aaec-40c9-89cc-35da5e913762` and its two Task children. Wall clock
for the lifecycle walk was about 34 minutes (21:02 orient through 21:36
finish). Delivery and this review are extra. Almost all of that time is
**human round-trips and coordinator reconstruction**, not the docs edit.

### Tool counts

| Agent | Reads | Shell | Grep | Glob | Other | Needed for the actual work |
| --- | ---: | ---: | ---: | ---: | --- | --- |
| Root coordinator | 83 (67 unique) | 37 | 16 | 15 | 2 Task, 10 Write | A `cc_route`, status headers, one plan bundle, then generate packets |
| Writer `06f590e4` | 14 (13 unique) | 9 | 2 | 4 | tried `move_agent_to_root` | Packet + worktree `docs/getting-started.md` + `test/contracts/test-contracts.sh` |
| Verifier `faeab665` | 15 (11 unique) | 5 | 5 | 4 | — | Packet, plan acceptance, writer handoff, worktree diff, three commands |

Root re-read `.agents/skills/cc-gates/SKILL.md` five times and
`cc-entry/SKILL.md` four times. `local-binding-troubleshoot/plan.yaml` and the
task file were each read three times. `docs/gates.md` three times. Those
repeats are skill-entry ritual, not new evidence.

### Context budgets vs bytes actually loaded

| Set | Budget | What happened |
| --- | ---: | --- |
| tier-0 | 8,192 | `wrapper/runtime/engine.sh` alone is **33,706 bytes** (719 lines). Entry always `Read`s it. That is already 4× the whole tier-0 ceiling, before AGENTS, WORKFLOW, identity, or summaries. |
| orientation | 20,480 | Turn 1 also read three **unrelated** `plan.yaml` files (CC-001, CC-002, CC-004) plus `cc-next`. Orientation prohibits unrelated plans. |
| run-plan | 22,528 | Execute copied eight prior multi-host/bootstrap runtime files to learn YAML shape, plus all five runtime schemas, `docs/runtime-contract.md`, and `invariants.yaml`. |
| writer | 12,288 | Writer also read `cc-execute` (a coordinator skill), globbed `**/handoff.md`, grepped `## Observed state` across the repo, and opened `test/upgrades/legacy-live/handoff.md`. |
| verifier | 11,264 | Verifier re-read `cc-verify` and the plan/task, globbed handoff examples and `repositories.local.yaml`, and grepped other sessions for handoff headings. |

No agent reported `CONTEXT_BUDGET_EXCEEDED` (`INV-CTX-02`). The budgets exist
on paper; live agents do not measure them.

### Was each read needed?

**Needed.** Selected CC-003 `plan.yaml` / task headers; `cc_route` via sourced
engine (shell, not a full file read); Git dirty/lease/worktree preflight;
writer packet + the two delegated worktree files; verifier’s three canonical
commands and path audit.

**Not needed, but expensive.**

1. Loading all of `engine.sh` into the model on entry. Shell can
   `. wrapper/runtime/engine.sh && cc_route "…"`. The 33kB file is an
   implementation library, not a context set (`AGENTS.md` already says that).
2. Reading sibling plans in full to discover which one was `draft`. A status
   grep would have been enough.
3. Copying previous session YAML (`sess-multi-host-agent-support-*`,
   `child-writer-multi-host-*`, `child-verifier-multi-host-*`,
   `completion.yaml`) because the engine does not emit session, lease,
   receipt, or delegation records.
4. Writer searching the tree for a handoff template. The packet already names
   `handoff_schema`. A ten-line skeleton in the packet would remove those
   globs and the legacy-live handoff read.
5. Writer inspecting `git show 3ba60a1` (the plan-draft commit, not the docs).
   The parent prompt already said the docs change was still to do.
6. Verifier globbing `repositories.local.yaml` and grepping other sessions.
   Non-goals already forbade creating that file; other sessions are prohibited
   on the writer set and are not verifier evidence.
7. Re-reading skills on every human turn. `cc-gates` is 595 bytes and does not
   change between approve, commit, and finish.

**Needed for safety, but they dominate wall clock.** Nine human gates. The
blocked first execute (turns 4–5) added a full human wait for a two-line
commit. Sequential writer-then-verifier is required; the waste is what each
child reads *before* the first delegated edit or first canonical command.

### Where time can come down

Dominant order, largest first:

1. **Human round-trips** — nine turns for one docs task. Combine the
   product-source approval-commit card with approval. Keep execute → finish
   as two turns, not three.
2. **Coordinator reconstruction** — ~150 root tool calls to re-learn routing,
   schemas, and prior packet shape. Generate runtime records from the engine;
   keep `engine.sh` out of the model context; cache skills in the host
   session.
3. **Child template search** — writer 14 reads + 4 globs before a 10-line
   docs insert; verifier 15 reads + 4 globs before three greps. Put the
   handoff skeleton and exact commands in the packet. Prohibit repo-wide
   glob from child prompts.
4. **Unrelated-plan orientation** — do not open other `plan.yaml` bodies.
   Index status in one listing.

The writer’s actual mutation was two `StrReplace`s and one commit. The
verifier’s actual proof was three commands and `git diff HEAD~1`. Everything
else was finding out how to behave.

## Highest-leverage next plan

Start with action cost: keep `engine.sh` out of the model, generate packets,
and stop child repo-wide template search. Then host child binding, named-plan
routing, and the worktree primitive. Do not add more lifecycle ceremony.

Suggested implementation order if a later plan is requested:

1. Stop loading `wrapper/runtime/engine.sh` into the model. Source it in
   shell for `cc_route` / lease / worktree. Generate session, receipt, lease,
   and delegation YAML from the engine so root does not copy prior packets.
   Put a handoff skeleton in the child packet so children do not glob the
   repo. Measure context-set bytes and report `CONTEXT_BUDGET_EXCEEDED`.
2. Bind Cursor/Codex/Claude writer children to the assigned worktree, and
   verifier children to read-only. Missing bind → `host-blocked`.
3. Make Stage A accept named plan ids and refuse `PLAN_REQUIRED` on
   orientation / future-tense lifecycle language.
4. Make `cc_prepare_bound_worktree` create a named branch from clean HEAD,
   including `codex/<plan-id>` style names.
5. Offer the product-source approval-commit card at approval time; avoid
   forking finish status onto both development and the implementation branch.
6. Refresh handoff, lease, and `PLAN.md` status on finish; keep packets as
   the only child scope; lock unique verification strings.

## Session facts (for later plans)

- Plan: `local-binding-troubleshoot` / CC-003 / LBT-001
- Wrapper: `1.0.0`
- Host: `cursor-agent` (`observed_version: unavailable`)
- Root session: `sess-local-binding-troubleshoot-20260822-01`
- Worktree commit: `31501e019ca6e9c57e80ade319165e4963982201`
- Delivery: https://github.com/kaotypr/context-circuit/pull/59
- Local `codex/development/v0.5` was left ahead (approval + finish commits)
  and not pushed, so the PR target would not receive finish before docs.
