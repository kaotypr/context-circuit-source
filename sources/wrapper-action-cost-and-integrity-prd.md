---
kind: prd
status: draft
title: Wrapper action cost and remaining integrity
---

# Wrapper action cost and remaining integrity

This is a draft PRD for reducing how expensive it is to find and perform
each already-named wrapper action, and for closing remaining integrity
gaps after CC-001. It is not a lifecycle redesign, not accepted Product
Knowledge, not a canonical workspace plan, and not authorization to
modify `wrapper/`, `template/`, skills, or any instantiated workspace.

Acceptance of this PRD does not approve an implementation plan, merge,
publication, or deployment. A later maintainer plan remains a separate
artifact and human approval gate.

## Problem and desired outcome

The Context Circuit lifecycle is the right product: identity first, plan
before code, approval changes status only, execution is a later named
request, and finish, delivery, archive, and cleanup stay explicit human
gates.

The failure is how agents find and perform each already-named action. A
human who does not know Context Circuit can still complete a useful loop
(initialize → draft a plan from a vague idea → approve → confirm → run).
That loop held in the Boardwalk v0.5.0 session. What did not hold is
bounded context, projection integrity, atomic runtime records, and
verification that matches promised behavior.

`draft-plan` is one router action. Live agents treat it as a research
project: they reread the router, glob the wrapper, and (in Boardwalk) read
sibling workspaces. The same pattern repeats on initialize, approve, and
execute. Context sets for `repository-bootstrap` already exist; `plan-draft`
and `initialization` still do not, so agents fill the vacuum by hunting
tests, fixtures, and sibling trees.

CC-001 (`repository-bootstrap`) already closed portable shared identity,
host-local bindings, clone/bootstrap confirmation, isolated runtime
worktrees, and credential-free Git (`INV-REPO-01` through `INV-REPO-07`).
Greenfield still has no first-class create-empty-repo gate beside that
clone path. After CC-001, `Run approved plan` correctly fails until a
bound, clean Git repository exists. Without a named empty-repo action,
agents fall back to inventing a path and running `git init` during
execute, and a naive human is stuck after `repositories: none`.

**Desired outcome:** keep the existing one-action lifecycle. Make that
path cheap, complete, and machine-enforced. Execute cannot create, clone,
or `git init` a product repository. A greenfield workspace has an explicit
next gate (bind, bootstrap, or create-empty). Identity summaries always
agree with `workspace.yaml`. One entry route loads one declared probe
within budget. Runtime records are engine-owned and atomic. Verification
cannot pass below the evidence layer the plan promised.

## Users and scenarios

### Users

| User | Need |
| --- | --- |
| Naive human in an instantiated wrapper | Ordinary prompts complete the lifecycle without knowing Context Circuit internals |
| Host agent (Codex, Claude Code, Cursor) | One router decision, one packet, no reconstruction of Stage A |
| Independent verifier child | Read-only evidence check; no repair |
| Maintainer in this product-source checkout | Orientation or draft-plan without recovery/lease/gate research; mutation still runs full preflight |

### Primary scenarios

1. **Greenfield initialize.** Human: `initialize this workspace`, then
   confirms identity. Omitted fields are shown as defaults or rejected as
   incomplete. All identity projections update together. `repositories:
   none` is accepted identity, not an execute-time Git invention.
2. **Draft a plan from a vague idea.** Human states an outcome. Coordinator
   invokes `cc_route`, loads the `plan-draft` packet only, and writes the
   plan bundle itself. No writer child. No sibling-workspace reads.
3. **Approve then run.** Approval changes status only. Run fails closed
   until the target repository is bound and clean. Missing binding, missing
   Git, or dirty source names the next gate (bind / bootstrap /
   create-empty / clean), not “continue execute.”
4. **Create-empty then execute.** After identity, a plan whose repository
   does not yet exist presents a current confirmation card. Confirmation
   creates an empty Git repository at an explicit destination and writes
   the host-local binding. Only then may execute prepare
   `.runtime/worktrees/<repository-key>/<plan-id>/`.
5. **Verify what was promised.** Browser-facing acceptance cannot pass on a
   production build alone. Missing host capability is an explicit waiver,
   not a pass.
6. **Product-source fast path.** In this maintainer checkout, read-only
   orientation uses the orientation packet; “create a plan” uses
   `plan-draft`; running a named approved plan uses full preflight,
   including CC-001 binding and worktree checks.

### Edge cases

- Omitted confirmation fields (Boardwalk: `roles` silently became `none`).
- Identity mutation that updates `workspace.yaml` / `WORKSPACE.md` but
  leaves `PROJECT.md` or `INDEX.md` stale.
- Cross-workspace reads (`tembiter-workspace`, another Context Circuit
  checkout, sibling `plans/context-circuit-plans/` from a customer
  workspace).
- Handmade session, receipt, delegation, lease, or completion records,
  including placeholder digests.
- Child launch whose prompt disagrees with the validated delegation.
- Product-source dirty-base special case (`INV-OWN-07`) must not skip
  lease or child checks on mutation.

### Success signals

A naive human completes initialize → draft → approve → (bind or
create-empty if needed) → run without the agent reconstructing the
router. The same lifecycle uses materially fewer file reads and prompt
bytes than the Boardwalk and product-source probes recorded.

## Requirements

Priority: **Must** ships before **Should**. **Could** may follow without
blocking Must. Optional items are not requirements until a later explicit
decision.

### Must — integrity the router does not yet enforce

1. **Greenfield empty-repo gate.** After identity acceptance, a plan whose
   repository does not yet exist must not execute. Present a current card
   that creates an empty Git repository at an explicit destination and
   writes the host-local binding. Same shape as `repository-bootstrap`:
   logical key, destination, existing-path check, confirmation. No clone,
   no credentials, no execute. Reject `git.init` inside `execute-plan`.
   Register a new named action (for example `repository-create-empty`)
   beside existing `repository-bootstrap`; do not overload clone
   bootstrap.

2. **Identity summaries are projections.** `context/WORKSPACE.md`,
   `context/PROJECT.md`, and `context/INDEX.md` are derived from
   `workspace.yaml`. One engine operation updates all of them. Entry and
   preflight run `cc_validate_workspace_projections` and fail on
   disagreement. Registering a logical repository updates every identity
   projection in that same operation.

3. **Gates list every projection they will touch.** Omitted confirmation
   fields are either displayed defaults on the card or the confirmation is
   rejected as incomplete. No silent field invention after confirmation
   (including `roles: none`).

4. **Split immediate effects from later authorized effects.** Every gate
   shows two sections: *this confirmed action changes now* and *a later
   authorized execution may change*. Represent effects with canonical
   fields. `git.init` and `git.clone` are not `execute-plan` effects. Use
   at least `workspace.register_repository`, `repository-bootstrap`,
   `repository-create-empty`, `git.commit` (worktree only), and
   `delivery.push`. Reject a delegation that introduces effects absent
   from approved intent.

5. **Execute cannot provision product Git.** `execute-plan` cannot create,
   clone, or `git init` a product repository. Isolated worktrees remain
   under `.runtime/worktrees/<repository-key>/<plan-id>/` per
   `INV-REPO-06`. Bound source stays untouched.

6. **Execute preflight names CC-001 blockers.** If binding is missing, the
   source is not Git, or the base is dirty, the next request is bind /
   bootstrap / create-empty / clean the source — not “continue execute.”
   Surface existing reason codes such as `BINDING_MISSING` and
   `NOT_A_GIT_REPOSITORY`. A docs paragraph is not a substitute for this
   route.

### Should — one action means one enforced packet

7. **Invoke the router; do not reread it.** Skills run
   `. wrapper/runtime/engine.sh && cc_route "$REQUEST" "$WORKSPACE_ROOT"`.
   Ordinary turns must not `Read` `engine.sh`, `routes.yaml`, or
   `invariants.yaml` unless the route is `blocked` or `recovery`. Cite
   invariant IDs; do not restate policy.

8. **Stage A returns an exact path allowlist and byte budget.** Loaders
   reject paths outside the selected set. Context-set overrun is a hard
   failure (`CONTEXT_BUDGET_EXCEEDED`), not thoroughness.

9. **Fill missing context sets.** Add at least `plan-draft` and
   `initialization` to `wrapper/contracts/context-sets.yaml`, modeled on
   the shipped `repository-bootstrap` packet. A `plan-draft` packet is
   small: identity, plan and task schemas, the two templates, the named
   target evidence, and one **shipped** example `plan.yaml`. It must not
   include `test/`, sibling workspaces, later gate docs, or maintainer
   `plans/context-circuit-plans/` as the YAML-shape example.

10. **Ship canonical fixtures in the wrapper.** Plan, task, gate, and
    runtime examples live in the released artifact so agents never need
    another Context Circuit checkout or a sibling customer workspace.
    Cross-workspace reads require an explicit scope-expansion route with a
    human-visible reason (`block-scope-expansion` unless expanded).

11. **Generate runtime records in the engine.** Add operations such as
    `cc_create_session`, `cc_create_receipt`, `cc_create_delegation`,
    `cc_start_child`, and `cc_record_completion`. Timestamps, versions,
    byte counts, and digests are engine-owned. Write to temp, validate,
    publish atomically. Reject placeholder revisions. Validate the
    ownership graph before child launch; launch requires a validation
    token. Reuse `cc_atomic_write` and binding-evidence helpers from
    CC-001; do not invent a parallel record layer.

12. **Delegation is the only child scope.** Launch prompts contain role,
    workspace root, delegation path, and handoff path. Scope, effects,
    acceptance, and stop conditions stay in the validated packet. The
    packet must be a subset of approved intent. Children stop if prompt
    text and delegation disagree. Coordinator writes `draft-plan` itself;
    no writer child for drafting.

13. **Product-source identity branch.** Keep the router. Add an identity
    branch for this maintainer checkout: `product-source` + read-only →
    orientation packet; `product-source` + create a plan → `plan-draft`
    packet; `product-source` + run a named approved plan → full preflight,
    including CC-001 binding and worktree checks. Do not skip dirty-base,
    lease, or child checks on mutation. Do skip recovery/lease/gate
    research for orientation and draft-plan.

### Could — evidence quality and human friction

14. **Evidence layers on acceptance criteria.** Classify each criterion as
    schema, store, API, process, browser, or human exploratory check.
    Reject verification mappings below the required layer. A production
    build cannot prove drag-and-drop. Missing host capability is an
    explicit waiver, not a pass. Process smoke (dev command + health) and
    browser journeys belong where the plan promised clicks.

15. **Digest tokens instead of long confirmation sentences.** Each card
    has a digest and a short `confirm <token>`. The token dies if target,
    state, or effects change. Unchanged Tier 0 is reused by receipt
    digest, not reloaded. Compact effect diffs; do not reprint the whole
    plan on every adjacent gate. Keep two-step confirmation for external
    and destructive work, including CC-001 bootstrap and the empty-repo
    gate.

16. **Fail live overread.** Add transcript or fixture tests for redundant
    Tier 0 reads on the same turn, broad wrapper/docs globs during
    `draft-plan` / `initialize`, and any read outside the workspace root.
    Static byte ledgers remain; they are not sufficient.

## Non-goals

- Re-implementing CC-001: portable identity, local bindings, clone
  bootstrap, isolated worktrees, release exclusions, credential rules
  (`INV-REPO-01` through `INV-REPO-07`).
- Embedding or reimplementing a host CLI, model API, SDK, MCP server, or
  scheduler.
- Storing credentials, tokens, provider payloads, or transcripts in
  workspace files.
- Auto-installing, auto-authenticating, or auto-invoking a host provider.
- A second router, plan lifecycle, ownership authority, or completion
  path.
- Automatic cloning, fetching, pushing, merging, publication, or
  deployment.
- Implementing these changes inside a Boardwalk (or other sample)
  workspace.
- Turning this file, or the source proposal, into `context/` Product
  Knowledge or an approved plan without a later, explicit human request.
- Folding this work into CC-001, CC-002 (`multi-host-agent-support`), or
  the experimental CC-003 docs draft
  (`plans/context-circuit-plans/local-binding-troubleshoot/`).
- Requiring the instantiated wrapper workspace itself to be a Git
  repository as part of this PRD (see open questions).

## Acceptance criteria

Wrapper action cost and remaining integrity is accepted for later planning when a
human can observe all of the following in the shipped wrapper (Must
first; Should and Could as scoped by the later plan):

1. Execute cannot create, clone, or `git init` a product repository.
2. A greenfield workspace with `repositories: none` has an explicit next
   gate (bind, bootstrap, or create-empty) rather than a silent execute
   fallback.
3. All Tier 0 identity summaries agree with `workspace.yaml`; disagreement
   fails preflight.
4. Omitted gate fields cannot be decided after confirmation.
5. One entry route loads one declared probe within the selected
   context-set budget.
6. Sibling workspace reads require explicit scope expansion.
7. `draft-plan` can complete from its packet without reconstructing the
   router or reading CC-001 as an example.
8. Runtime records contain no placeholders and are published atomically.
9. No child starts before its delegation and ownership graph validate.
10. Browser-facing acceptance criteria have browser or explicit human
    evidence, or an explicit waiver.
11. Gate text separates immediate effects from later authorized effects,
    and those effects match CC-001 ownership.
12. Delivery, bootstrap, empty-repo, and destructive gates remain as
    strict as CC-001 bootstrap is today.
13. The same lifecycle requires materially fewer file reads and prompt
    bytes than the Boardwalk v0.5.0 session and the 2026-08-22
    product-source `draft-plan` probe.

Acceptance defines requirements. It does not approve a plan or authorize
execution.

## Constraints and dependencies

- Shipped product contract lives under `wrapper/`; mutable blank seed
  under `template/`. Product behavior has one owner per rule via
  `wrapper/contracts/invariants.yaml`. Do not add parallel policy to a
  skill or role file.
- Keep: review is read-only; approval changes status only; execution is a
  later named request; finish is a later status gate. Runtime evidence
  never marks a plan `done`. Writers stay path-bounded; verifiers stay
  independent and do not repair. Push, merge, publication, deployment,
  archive, takeover, and cleanup stay explicit human gates.
- `sources/` stays passive and request-scoped.
- CC-001 contracts that this PRD must not weaken: `INV-REPO-01` through
  `INV-REPO-07`, `cc_prepare_bound_worktree`, dirty-source blocking, and
  credential-free Git.
- Existing context sets to extend, not replace: `tier-0`, `orientation`,
  `plan-review`, `run-plan`, `repository-bootstrap`, `writer`, `verifier`,
  `resume`. Overrun already reports `CONTEXT_BUDGET_EXCEEDED`.
- Existing routes already include `repository-bootstrap` and
  `present-repository-bootstrap-card`. Empty-repo is an additional named
  action and human gate, not a reuse of clone bootstrap.
- Host adapters (Codex CLI, Claude Code, Cursor Agent CLI) share one
  instruction surface. Host evidence never authorizes a route or gate. If
  a required child is unavailable, preserve read-only `host-blocked`;
  never self-verify.
- `INV-OWN-07` already special-cases product-source approval-status dirt.
  `WORKFLOW.md` already says source work is not the product lifecycle.
- Implementation belongs in a later maintainer plan in this source
  checkout, not in a Boardwalk sample workspace.

## Assumptions and open questions

### Assumptions

- The lifecycle in `wrapper/contracts/routes.yaml` remains the single
  router. This PRD improves action cost and remaining wrapper integrity
  after CC-001.
- Boardwalk conversation evidence is treated as agent-behavior evidence
  against released v0.5.0, not as remaining repository-identity work.
  “Create a product git repo during execute” is closed.
- Canonical effect field names in requirement 4 may be refined in the
  later plan as long as execute still cannot own `git.init` / `git.clone`
  and empty-repo remains a distinct confirmed action.
- “Materially fewer file reads” in acceptance criterion 13 will be given
  a numeric budget in the later plan (the product-source probe was ~80
  tool steps and 70+ file reads before three writes).

### Open questions

1. **Should an instantiated workspace be versioned?** CC-001 versions the
   product repository. It does not require the instantiated workspace
   (`workspace.yaml`, `context/`, `plans/`, wrapper overlay) to be a Git
   repository. If pursued, it is a wrapper-bootstrap / workspace-git gate,
   not product `git.init` during execute, and must not nest that history
   inside `repositories/`. This PRD treats it as out of scope until an
   explicit later decision. Lower priority than Must and Should.
2. **Exact action and context-set identifiers** for empty-repo
   (`repository-create-empty` vs another name) and for `plan-draft` vs
   `draft-plan` alignment with existing route ids.
3. **Numeric byte budgets** for new `plan-draft` and `initialization`
   packets (existing sets: tier-0 8192, orientation 20480, run-plan
   22528, repository-bootstrap 12288).
4. **Whether confirmation tokens (requirement 15)** replace or sit beside
   today’s two-step phrase cards for non-destructive gates.

## Product Knowledge references

None. The source proposal and this PRD are not accepted Product
Knowledge. Do not copy them into `context/` without a later, explicit
human request.

## Provenance

Chosen artifact path:
`sources/wrapper-action-cost-and-integrity-prd.md`.

| ID | Path | Why it was read |
| --- | --- | --- |
| PROP | `sources/boardwalk-project-test/context-circuit-improvement-proposal.md` | Named source. Thesis, remaining gaps after CC-001, P0–P2 proposals, success criteria, and non-goals. |
| PRD-DOC | `docs/prd.md` | PRD contract, artifact home under `sources/`, required sections. |
| PRD-TPL | `docs/templates/prd.md` | Frontmatter and section skeleton. |
| SRC-README | `sources/README.md` | Sources remain passive and request-scoped. |
| CS | `wrapper/contracts/context-sets.yaml` | Confirm shipped sets; `plan-draft` and `initialization` are absent; `repository-bootstrap` is the packet pattern. |
| RT | `wrapper/contracts/routes.yaml` | Confirm single router, existing `repository-bootstrap` action/gate, missing empty-repo action. |
| T0 | `wrapper/contracts/tier0.yaml` | Confirm probes include `initialization` as a catalog name without a matching context set. |
| INV | `wrapper/contracts/invariants.yaml` | Confirm `INV-REPO-01`–`INV-REPO-07` and `INV-OWN-07` ids for constraints. |

Upstream evidence cited by PROP and **not re-read** for this draft:
`sources/boardwalk-project-test/boardwalk-conversation-codex-context-circuit-v0.5.0.txt`
(BW-CONV),
`sources/boardwalk-project-test/context-circuit-workspace-improvement-review.md`
(BW-REVIEW), CC-001 plan bundle, and product-source conversation SRC-ENTRY.

Related and excluded: CC-002 `multi-host-agent-support` (host portability);
experimental untracked CC-003
`plans/context-circuit-plans/local-binding-troubleshoot/`.

Drafted 2026-08-22 from the named proposal and the contract files listed
above.
