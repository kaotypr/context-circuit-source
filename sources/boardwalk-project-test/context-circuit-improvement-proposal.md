# Context Circuit improvement proposal

This is a passive source artifact. It records a product-source proposal for
improving Context Circuit as a shipped wrapper. It is not accepted Product
Knowledge, not a canonical workspace plan, and not authorization to modify
`wrapper/`, `template/`, skills, or any instantiated workspace.

Drafted 2026-08-22 from named evidence only. Revised the same day after
reading the completed `repository-bootstrap` plan (CC-001). Do not scan the
rest of `sources/`.

## Provenance

| ID | Path | Why it was read |
| --- | --- | --- |
| BW-CONV | `sources/boardwalk-project-test/boardwalk-conversation-codex-context-circuit-v0.5.0.txt` | Live v0.5.0 wrapper session: naive human prompts, tool traces, lifecycle gates |
| BW-REVIEW | `sources/boardwalk-project-test/context-circuit-workspace-improvement-review.md` | Boardwalk workspace review of how the wrapper should improve |
| CC-001 | `plans/context-circuit-plans/repository-bootstrap/plan.yaml` plus `PLAN.md` and `tasks/001-contracts.md` … `005-verification.md` | Done maintainer plan, implemented after Boardwalk, for portable repo identity and explicit bootstrap |
| SRC-ENTRY | Product-source conversation, 2026-08-22 | Orientation and “create a plan” probe in this maintainer checkout, including isolated draft-plan agent `2171c3bf-a239-4577-9a01-5f59f01e83ca` |

Related but not copied into this artifact: experimental untracked draft
`plans/context-circuit-plans/local-binding-troubleshoot/` (CC-003) and approved
maintainer plan `multi-host-agent-support` (CC-002). CC-002 is host portability.
This proposal is action cost and remaining wrapper integrity. Do not fold it
into CC-001 or CC-002, and do not re-implement CC-001.

## Thesis

The lifecycle is the right product. The failure is how agents *find and
perform* each already-named action.

A human who does not know Context Circuit can still complete a useful loop:
initialize → draft a plan from a vague idea → approve → confirm → run. That
loop held in Boardwalk. What did not hold is bounded context, projection
integrity, atomic runtime, and verification that matches promised behavior.

`draft-plan` is one router action. Live agents treat it as a research project:
reread the router, glob the wrapper, and (in Boardwalk) read sibling
workspaces. The same pattern repeats on initialize, approve, and execute.

Do not add a second router, scheduler, or lifecycle. Make the existing
one-action path cheap, complete, and machine-enforced.

CC-001 is already done. Several Boardwalk repository failures are historical
against v0.5.0, not remaining product gaps. This proposal keeps only what
CC-001 did not own.

## Timeline

Boardwalk ran on the released v0.5.0 wrapper **before** CC-001 landed. The
conversation therefore shows an agent creating `repositories/boardwalk` as
the exclusive worktree during `Run approved plan`, with `git init` and a
local commit, while `workspace.yaml` started as `repositories: none`.

CC-001 (`status: done`) later shipped portable shared identity, host-local
bindings, an explicit clone/bootstrap gate, isolated runtime worktrees, and
dirty-source blocking. `INV-REPO-01` through `INV-REPO-07` and
`cc_prepare_bound_worktree` are now the contract.

Read Boardwalk as evidence of agent behavior and remaining non-repo gaps.
Do not treat “create a product git repo during execute” as still-open work.

## What to keep

- Review is read-only; approval changes status only; execution is a later
  named request; finish is a later status gate.
- Runtime evidence never marks a plan `done`.
- Writers stay path-bounded; verifiers stay independent and do not repair.
- Push, merge, publication, deployment, archive, takeover, and cleanup stay
  explicit human gates.
- No credentials or provider payloads in workspace files.
- `sources/` stays passive and request-scoped.
- CC-001 repository identity, local bindings, bootstrap confirmation,
  isolated worktrees, and credential-free Git. Do not weaken them.

## Already closed by CC-001

Do not propose these again:

| Boardwalk observation | Current contract |
| --- | --- |
| Writer created `repositories/boardwalk` as the execute worktree | `INV-REPO-06`: bound source stays untouched; execute prepares only `.runtime/worktrees/<repository-key>/<plan-id>/` |
| Execute did `git init` because the product path was not a repo | `cc_prepare_worktree` requires `$base/.git`; `cc_validate_repository_binding` returns `NOT_A_GIT_REPOSITORY` otherwise |
| Missing or invented checkout path, filesystem scan | `INV-REPO-03`: missing, ambiguous, or unsafe bindings fail without scanning or substitution (`BINDING_MISSING`, …) |
| Clone/create happened inside run-plan | `INV-REPO-04`: clone or destination creation needs a current `repository-bootstrap` confirmation |
| Machine paths in shared identity | `INV-REPO-01` / `INV-REPO-02`: `workspace.yaml` holds logical keys only; paths live in ignored `repositories.local.yaml` |
| Credentials in workspace files | `INV-REPO-05` |
| `repositories/` leaking into release | `INV-REPO-07` |

CC-001 also added a real `repository-bootstrap` context set. `plan-draft` and
`initialization` still have none. That contrast is the remaining packet gap.

## Remaining after CC-001

### Greenfield still has no first-class create-empty-repo gate

Boardwalk was a **new** local app with no remote. CC-001 bootstrap is
`git clone` of a selected remote to an absent destination. It does not
create an empty local repository.

After CC-001, `Run approved plan` must fail until a bound, clean Git
repository exists. That is correct. The missing human path is:

1. Register a logical repository in `workspace.yaml` (identity only).
2. Either bind an existing checkout, bootstrap a clone, **or** confirm
   creation of an empty local Git repo at an explicit destination.
3. Only then run the approved plan in `.runtime/worktrees/…`.

Without step 2 as its own gate, agents will fall back to the Boardwalk
habit: invent a path and `git init` during execute. The engine should
keep refusing that. The wrapper still needs a named empty-repo action so
a naive human is not stuck after `repositories: none`.

This is not a redo of CC-001. It is one additional gated effect, with the
same confirmation and credential rules.

### Wrapper workspace Git is still distinct from product Git

CC-001 versions and isolates the **product** repository. It does not
require the **instantiated workspace** (`workspace.yaml`, `context/`,
`plans/`, wrapper overlay) to be a Git repository.

Boardwalk mixed those two histories: the wrapper directory had no Git,
then execution created Git only inside the nested product path. CC-001
stops the nested-as-worktree part. It does not put accepted identity,
context, and plans on a wrapper history.

That remaining question is optional and smaller than BW-REVIEW stated:
should an instantiated workspace be versioned, or is filesystem evidence
enough? If versioned, it is a wrapper-bootstrap/workspace-git gate, not
product `git.init` during execute.

### Everything CC-001 never owned

Identity projections, omitted gate fields, overread, sibling-workspace
reads, handmade runtime records, verification below promised behavior,
ambiguous gate wording, duplicated child prompts, and long confirmation
phrases. Those remain.

## Evidence

### A. Instantiated wrapper (Boardwalk, pre-CC-001)

Human prompts were ordinary:

1. `initialize this workspace`
2. `Confirm initialize: name=simple-kanban-board, mode=solo, default_branch=main, repositories=none` (roles omitted)
3. A plain-language Boardwalk idea (local Trello, no login)
4. `Approve plan 0001-local-kanban-mvp`
5. `Confirm approval of 0001-local-kanban-mvp`
6. `Run approved plan 0001-local-kanban-mvp`

Observed goods: identity before other writes; plan before code; approval did
not start a writer; writer and verifier were separate children; plan status
stayed `approved` after implementation; finish/delivery remained gated.

Observed failures still live:

1. Identity projections drifted. `workspace.yaml` / `WORKSPACE.md` gained a
   repository; `PROJECT.md` still reported none; `INDEX.md` was left stale
   because it was outside the gate’s declared effects.
2. Omitted `roles` was silently accepted as `none`.
3. Initialize and draft-plan reread the same spine many times, used broad
   globs, and read sibling trees `claude-context-circuit-v0.5.0` and
   `tembiter-workspace` for plan examples. No scope-expansion gate.
4. Session, receipt, delegation, lease, and completion records were assembled
   by hand; a receipt used a placeholder digest, then was patched.
5. Verifier passed a production build. Acceptance promised browser
   drag-and-drop, edit, search, filter, and export. The finish card admitted
   UI clicks were not in the suite.
6. Approval said it would not change Git; later execution initialized a repo
   and committed. After CC-001 the *product* Git change is not an execute
   effect; the wording problem remains for whatever later gates own.
7. Writer prompt copied product requirements already owned by the plan and
   delegation.

Historical (closed by CC-001): wrapper had no product-repo binding; execute
created `repositories/boardwalk` as the worktree and ran `git init`.

### B. Product-source probe (this checkout, post-CC-001)

Request: create a plan to add one troubleshooting paragraph to
`docs/getting-started.md`.

Router decision was already one action: probe `plan-draft`, capability
`draft-plan`. The agent did not implement, approve, claim a lease, or spawn
children. It did spend ~80 tool steps and 70+ file reads before three writes.

`cc_route` already emits `context_set: plan-draft`.
`wrapper/contracts/context-sets.yaml` has `repository-bootstrap` but no
`plan-draft` packet. The agent filled the vacuum by hunting tests, fixtures,
and sibling maintainer plans — including CC-001 `plan.yaml` only as a v2
shape example.

`cc-next` already says to delegate routing to `wrapper/runtime/engine.sh`.
`cc-entry` still tells the agent to read adapters and reconstruct Stage A.

Byte-budget tests sum a pre-listed file set. They cannot see an agent opening
70 files to write 3, or leaving the workspace to read another project.

## Proposal

### P0 — Integrity the router does not yet enforce

1. **Greenfield empty-repo gate (new, beside CC-001).** After identity
   acceptance, a plan whose repository does not yet exist must not execute.
   Present a current card that creates an empty Git repository at an explicit
   destination and writes the host-local binding. Same shape as
   `repository-bootstrap`: logical key, destination, existing-path check,
   confirmation. No clone, no credentials, no execute. Reject `git.init`
   inside `execute-plan`.

2. **Identity summaries are projections.** `context/WORKSPACE.md`,
   `context/PROJECT.md`, and `context/INDEX.md` are derived from
   `workspace.yaml`. One engine operation updates all of them. Entry and
   preflight run `cc_validate_workspace_projections` and fail on disagreement.
   Registering a logical repository updates every identity projection in that
   same operation. This is the live form of Boardwalk’s stale `PROJECT.md`.

3. **Gates list every projection they will touch.** Omitted confirmation
   fields are either displayed defaults on the card or the confirmation is
   rejected as incomplete. No silent `roles: none` after the fact.

4. **Split immediate effects from later authorized effects.** Every gate
   shows two sections: *this confirmed action changes now* and *a later
   authorized execution may change*. Represent effects with canonical fields.
   After CC-001, `git.init` and `git.clone` are **not** execute-plan effects.
   Use `workspace.register_repository`, `repository-bootstrap`,
   `repository-create-empty`, `git.commit` (worktree only), `delivery.push`.
   Reject a delegation that introduces effects absent from approved intent.

5. **Optional: version the instantiated workspace.** Separate from product
   Git. If pursued, require wrapper-workspace Git or an explicit
   workspace-git bootstrap before treating `workspace.yaml` / `context/` /
   `plans/` as durable. Do not nest that history inside `repositories/`.
   This is lower priority than (1)–(4) because CC-001 already stopped
   execute from using an unversioned product path as the worktree.

### P1 — One action means one enforced packet

6. **Invoke the router; do not reread it.** Skills run:

   ```sh
   . wrapper/runtime/engine.sh && cc_route "$REQUEST" "$WORKSPACE_ROOT"
   ```

   Ordinary turns must not `Read` `engine.sh`, `routes.yaml`, or
   `invariants.yaml` unless the route is `blocked` or `recovery`. Cite
   invariant IDs; do not restate policy.

7. **Stage A returns an exact path allowlist and byte budget.** Loaders
   reject paths outside the selected set. Context-set overrun is a hard
   failure, not thoroughness.

8. **Fill missing context sets.** At least `plan-draft` and
   `initialization`, using `repository-bootstrap` as the packet pattern CC-001
   already shipped. A `plan-draft` packet should be small: identity, plan and
   task schemas, the two templates, the named target evidence, and one
   **shipped** example `plan.yaml`. Not `test/`, not sibling workspaces, not
   later gate docs, not rereading CC-001 for YAML shape.

9. **Ship canonical fixtures in the wrapper.** Plan, task, gate, and runtime
   examples live in the released artifact so agents never need
   `tembiter-workspace`, another Context Circuit checkout, or maintainer
   `plans/context-circuit-plans/` from a customer workspace. Cross-workspace
   reads require an explicit scope-expansion route with a human-visible reason.

10. **Generate runtime records in the engine.** Add operations such as
    `cc_create_session`, `cc_create_receipt`, `cc_create_delegation`,
    `cc_start_child`, `cc_record_completion`. Timestamps, versions, byte
    counts, and digests are engine-owned. Write to temp, validate, publish
    atomically. Reject placeholder revisions. Validate the ownership graph
    before child launch; launch requires a validation token. Reuse
    `cc_atomic_write` and binding-evidence helpers CC-001 already added;
    do not invent a parallel record layer.

11. **Delegation is the only child scope.** Launch prompts contain role,
    workspace root, delegation path, and handoff path. Scope, effects,
    acceptance, and stop conditions stay in the validated packet. The packet
    must be a subset of approved intent. Children stop if prompt text and
    delegation disagree. Coordinator writes `draft-plan` itself; no writer
    child for drafting.

12. **Execute preflight names the CC-001 blockers.** If binding is missing,
    the source is not Git, or the base is dirty, the next request is bind /
    bootstrap / create-empty / clean the source — not “continue execute.”
    The experimental CC-003 docs paragraph is not a substitute for this
    route.

### P2 — Evidence quality and human friction

13. **Evidence layers on acceptance criteria.** Classify each criterion as
    schema, store, API, process, browser, or human exploratory check. Reject
    verification mappings below the required layer. A production build cannot
    prove drag-and-drop. Missing host capability is an explicit waiver, not a
    pass. Process smoke (dev command + health) and browser journeys belong
    where the plan promised clicks.

14. **Digest tokens instead of long confirmation sentences.** Each card has
    a digest and a short `confirm <token>`. The token dies if target, state,
    or effects change. Unchanged Tier 0 is reused by receipt digest, not
    reloaded. Compact effect diffs; do not reprint the whole plan on every
    adjacent gate. Keep two-step confirmation for external and destructive
    work, including CC-001 bootstrap and the new empty-repo gate.

15. **Fail live overread.** Add transcript or fixture tests for: redundant
    Tier 0 reads on the same turn, broad wrapper/docs globs during
    `draft-plan` / `initialize`, and any read outside the workspace root.
    Static byte ledgers remain; they are not sufficient.

### Product-source fast path

This maintainer checkout is already identified. `WORKFLOW.md` already says
source work is not the product lifecycle. CC-001’s dirty-base rule already
special-cases product-source approval-status dirt (`INV-OWN-07`).

Keep the router. Add an identity branch:

- `product-source` + read-only → orientation packet
- `product-source` + create a plan → `plan-draft` packet
- `product-source` + run a named approved plan → full preflight, including
  CC-001 binding and worktree checks

Do not skip dirty-base, lease, or child checks on mutation. Do skip
recovery/lease/gate research for orientation and draft-plan.

## Suggested implementation order

1. Empty-repo gate beside CC-001; projection validation; split gate effects;
   execute preflight that names `BINDING_MISSING` / `NOT_A_GIT_REPOSITORY`
   instead of inventing a checkout.
2. `cc_route` as the only Stage A; allowlists; missing context sets modeled
   on `repository-bootstrap`; canonical fixtures; atomic runtime ops;
   delegation-only children.
3. Evidence-layer metadata; confirmation tokens; overread tests; thinned
   skills; product-source fast path.
4. Optional instantiated-workspace Git, only after (1)–(3).

That is one future maintainer plan, sourced from this file. It is not
CC-001, not CC-002, and not the experimental CC-003 docs draft.

## Success criteria

Context Circuit is improved when:

- execute cannot create, clone, or `git init` a product repository;
- a greenfield workspace with `repositories: none` has an explicit next
  gate (bind, bootstrap, or create-empty) rather than a silent execute
  fallback;
- all Tier 0 identity summaries agree with `workspace.yaml`;
- omitted gate fields cannot be decided after confirmation;
- one entry route loads one declared probe within budget;
- sibling workspace reads require explicit scope expansion;
- `draft-plan` can complete from its packet without reconstructing the
  router or reading CC-001 as an example;
- runtime records contain no placeholders and are atomic;
- no child starts before its delegation and ownership graph validate;
- browser-facing acceptance criteria have browser or explicit human evidence;
- gate text separates immediate effects from later authorized effects, and
  those effects match CC-001 ownership;
- the same lifecycle requires materially fewer file reads and prompt bytes;
- delivery, bootstrap, empty-repo, and destructive gates remain as strict
  as CC-001 bootstrap is today.

## Non-goals

- Re-implementing CC-001 (portable identity, local bindings, clone
  bootstrap, isolated worktrees, release exclusions, credential rules).
- Embedding or reimplementing a host CLI, model API, SDK, MCP server, or
  scheduler.
- Storing credentials, tokens, provider payloads, or transcripts in workspace
  files.
- Auto-installing, auto-authenticating, or auto-invoking a host provider.
- A second router, plan lifecycle, ownership authority, or completion path.
- Automatic cloning, fetching, pushing, merging, publication, or deployment.
- Implementing these changes inside a Boardwalk (or other sample) workspace.
- Turning this file into `context/` Product Knowledge or an approved plan
  without a later, explicit human request.
