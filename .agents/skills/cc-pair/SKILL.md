---
name: cc-pair
description: Work directly with the user in one connected repository through a coordinator-to-worker loop with live human supervision and no verifier — the Explore tier of the assurance ladder, promotable in place to a proper checked change.
---

## When to use

Direct collaboration is the **Explore tier** of the one assurance ladder
(INV-ASSURE-01), not a separate world: the bottom rung where a small, live,
reversible change is worked with the user directly, with no independent verifier.
It answers "do I even need a verifier for this?" — at Explore the answer is "no,"
structurally: coordinator (root) + **one worker child**, and no verifier and no
tracer (the human reads the code alongside the agent live). Use it whenever one repository is connected, by ordinary intent or
`/cc-pair`, and it may be offered as an optional next step after a plan or stack
execution. Never enter it automatically.

Before promotion it is outside the plan lifecycle (INV-PAIR-01): it neither invokes
nor is invoked by intent approval, execution, verification, completion, or delivery.
When the work turns out to be real, **promote it in place** (below) rather than
stopping and restarting as a plan — the cliff is a ramp.

## Start or resume

Resolve exactly one connected repository. A session cannot span repositories; if
the request is ambiguous among several connected repositories, ask one focused
question before creating state.

Before any write, confirm that the host can create a worker child. If it cannot,
report `host-blocked` in plain language and remain read-only. The coordinator must
never perform the worker's edits itself.

Ask the human for a short session name and use that name. Map it to a safe
lowercase kebab slug (`cc_safe_slug` rules: lowercase kebab, at most 80
characters). If the name is missing or unsafe, ask again — **never invent** the
folder name.

Choose a base commit:

- fresh work defaults to the connected repository's recorded base tip;
- when offered after a completed plan or stack execution, use that repository's
  completed execution tip;
- use another commit only when the user names it explicitly.

Invoke `pair-begin . <repo> <session> [base]`. It creates a fresh
`cc-pair/<session>` branch and isolated worktree at
`.runtime/explore/<session>/<repo>`, not under `.runtime/worktrees/cc-pair/`. It
never reuses the active checkout or edits a plan branch in place. To resume,
invoke `pair-inspect . <session>` and attach a worker to the reported worktree
only when `resumable: true`.

## Model & effort per role

Spawn the Explore worker at the concrete `(model, effort)` configured for the
`worker` role in the host-local role-tiering config, with adapter-shipped defaults
for an unset role (`docs/role-tiering.md` owns the shape, defaults, and
escalation ladder). This is a coordinator/host decision — the runtime is
model-blind (INV-RUNTIME-01) and `(model, effort)` authorizes nothing
(INV-HOST-01). It changes cost and speed, never meaning. Explore still has no
verifier; do not launch one because a yaml file exists.

- Read `role-tiering.local.yaml` from the workspace root when present — the same
  directory as `repositories.local.yaml`, never a repository working copy —
  and apply the current host's `worker` entry on the spawn. A missing file in
  an isolated working copy is not an absent config.
- Adapter defaults still apply when no local file, host group, or worker role is
  present.
- Applying the setting does not change Explore's meaning, drop or add a verifier,
  or alter independence or the failure limit.

## The interactive loop

Launch one worker using the existing `agents/worker.md` role at that configured
`(model, effort)` and give it the
reported repository, worktree, branch, base, and the user's current intent. Tell
the worker to:

- work only in that pairing worktree and never create, switch, merge, push, or
  remove branches or worktrees;
- discover, read, and honor the repository's own agent guidance before writing;
- make the concrete change delegated for this turn, run relevant checks, and
  report what changed and any unresolved concern;
- never call its work verified, complete a plan, or deliver;
- commit only when the user explicitly asks, using INV-COMMIT-01.

Keep the same worker attached when the host supports it; otherwise a replacement
worker may resume from the pointer and worktree. For each user reaction, interpret
the ordinary-language intent, delegate one concrete next change, and narrate the
effect. The user—not the coordinator or worker—judges whether the result is right.
Do not launch a verifier and do not create leases, execution records, evidence,
handoffs, failure counters, or plan state.

Relevant build or test checks are feedback for the live loop; they do not create
independent verification evidence and must never be reported as verified status.

## Convergence

The session converges only when the user explicitly says the result is finished.
Inspect first. If the worktree is dirty, do not commit automatically and do not
close: report the uncommitted work and ask whether the worker should commit it or
the session should remain open. Once clean, invoke `pair-close . <session>`.

Closing preserves the branch and worktree. Cleanup is a separate explicit human
action (`runtime-cleanup`); closing a clean session does not delete its worktree,
and a still-live session is not removed unless cleanup was explicitly asked. Describe the result as **human-supervised, not independently verified**.
Offer delivery only as a separate action through `cc-deliver`.

## Promote — the ramp into the trust system

When an Explore session turns out to be real work — the change matters, or it
touches logic beyond a quick tweak — offer to **promote it in place** rather than
restart. Promotion is the moment an intent and a plan first exist for the work and
it enters the trust pipeline (candidate → verifier → delivery → reconciliation):

1. **Attach an intent.** Run `cc-intent` to author the goal, criteria, scope, and
   tier for what the session is actually doing, and the human approves it (Gate 1).
2. **Raise the tier.** Standard/Critical, per the risk signals — this is the moment a
   tracer reads the code (`cc-trace`) and the independent verifier appears
   (INV-ASSURE-01). A risk surface refuses to stay Explore.
3. **Author a lightweight plan of record.** `cc-plan` creates a `plans/<id>/` whose
   tasks describe the change already made; its execution binds to the existing
   pairing-branch commits, which produces the candidate. Everything downstream then
   runs on the standard scaffolding, unchanged.

Do not fabricate criteria or scope silently — surface them and let the human
approve. Un-promoted Explore work never produces a plan file or a candidate; that
is what keeps the fast path fast.

## Runtime actions — invoke, never read the engine

Invoke from the workspace directory:

- `sh wrapper/runtime/engine.sh pair-begin . <repo> <session> [base]`
- `sh wrapper/runtime/engine.sh pair-inspect . <session>`
- `sh wrapper/runtime/engine.sh pair-close . <session>`
- `sh wrapper/runtime/engine.sh runtime-cleanup .` — explicit leftover cleanup;
  removes closed Explore worktrees under `.runtime/explore/` as well as idle
  plan-execution worktrees. Never run this from `pair-close`.

The runtime is an opaque deterministic library (`wrapper/adapters/AGENTS.md` →
Runtime). The skill owns the conversation; the runtime owns only Git isolation
and the light pointer.

## Report — plain language only

Narrate what changed and what the user still needs to decide. Do not expose the
pairing branch name, worktree path, base commit, pointer, or runtime commands
unless the user explicitly asks for diagnostics. Say "a separate working copy"
and "the changes we made together," following `docs/terminology.md`. Name changed
files by repository-relative path (for example, `src/widget.ts`), never with an
absolute clickable target that reveals the hidden working-copy path.

## Boundaries

Never create a plan, approve, independently verify, mark complete, rebase on
delivery drift, push, merge, open a pull request, or clean up as part of pairing.
If the worktree cannot be created, preserve any reported state and report
`host-blocked`; the coordinator remains read-only. Do not invent the Explore
folder name.
