# Authorization boundaries

What a person decides, and what one decision covers afterwards.

## Decisions belong to a person

A fourth is theirs too, below: whether any of this applies to the change at
hand.

**Approving the intended outcome.** Before detailed code investigation the agent
writes the intent, presents it, and stops. What it is written against, and how a
decision it cannot make reaches the person as a numbered open question, is
[grounding an intent](intent-grounding.md). The rule that carries the weight:
neither a command, nor an editable record, nor another agent can supply human
consent. The approval operation records an actual decision and its date; it does
not constitute one.

The intent's `approved_at` field is stamped with the instant there, and the
person's words are kept under a heading carrying the same instant; an unapproved
intent carries no such field, so absence is the only empty state. The field is a
readable record of the gate, never the gate itself — anything that can write the
record can write the timestamp, so an intent is not approved because a field
says so. Nothing in the CLI reads it to decide whether to proceed.

Nor does the request that prompted the intent. Reading a change request as
consent to whatever intent gets derived from it collapses the gate entirely —
the agent writes the outcome and approves it in the same breath — so approval is
a person's answer to the intent as written.

Re-approval is required only when the outcome or its success criteria materially
change. A newly needed file or repository *within* the approved outcome is
explained and recorded, not re-gated, because recorded paths are descriptive
planning information rather than an enforcement contract. Asking twice for the
same thing teaches people to approve reflexively.

An explicit request for a detailed plan of an already specified outcome may
create a standalone plan without an intent. It authorizes investigation and
drafting, not implementation. The agent presents its complete entry and
supporting files before execution is requested.

**Starting execution.** Plans are derived, presented, and then left alone until a
person asks for them to be run. This is a decision about *when work begins*, not
a review of the derivation, which is why it sits beside plan presentation without
being a plan approval gate.

It is deliberately not pre-suppliable. An instruction to implement that arrived
before the plans existed cannot be a decision about plans nobody had seen, so it
does not start execution — otherwise "add billing and implement it" carries the
agent from an unread intent through unread plans into repository edits, which is
the failure the two stops exist to prevent.

**Authorizing an outward action.** Push, pull-request creation, merge into a base
branch, deployment, external publication, and deletion of workspace data each
require explicit authorization, and authorization already given is reused rather
than re-requested.

Committing is split by where it lands. A commit on a prepared worktree's own
branch is implementation: it leaves the machine no more than an edit does, the
branch belongs to the workspace, and three things downstream read commits and
nothing else — a dependent plan's starting point, an integration merge, and a
resume that inspects real diffs rather than trusting a record. So implementation
ends with its work committed, and nothing is left uncommitted when it is
reported. A commit in a bound checkout is a change to the person's own working
copy and stays authorized-only, the same carve-out direct execution already has.

The rule exists because the alternative was observed: a plan's work sat
uncommitted, its dependent plan prepared a worktree from the base branch that
therefore held none of it, and the agent copied files between worktrees to
stand in for ancestry it had no other way to get.

## A fourth decision: which path to take at all

A person may skip the intent and the plans entirely and ask for the change
directly. That choice is theirs alone — the agent may offer the lane where a
request is already its own specification, and never takes it unasked. What the
bypass does and does not relax is [changing something
directly](direct-changes.md).

It relaxes nothing here. The request stands in for the approval only because
nothing was derived to approve; every outward action, and a commit in the bound
checkout, still needs the same word from the same person.

## What the execution request is not

It is not a machine-recorded plan approval. The person receives the complete
plan and may revise it; a separate request after that presentation starts work.
There is no approval field or command on a plan. The diagnostic is not a gate.

## What the execution request covers

Preparing a worktree for each repository a plan names, and implementing there.
Preparation is a step of execution rather than a strategy chosen per task, so it
is authorized by the same request and never gated separately. Asking to work
directly in a bound checkout is the stated exception, and the only one.

## Scoped run authorization

Confirming a stacked run authorizes, **for that run**, worktree preparation,
implementation commits on the run's own branches, and local integration merges
that assemble a dependent plan's base. It does not authorize push, pull-request
creation, merge into a base branch, deployment, or deletion.

The integration merge is the case worth naming: it is a merge, and merges are
usually delivery, but this one assembles a starting point inside the workspace's
own branches and never touches a base branch. It is classified as implementation
in advance, so nobody has to reason their way to either answer mid-run. That is
the general pattern — an authorization has a scope, the scope is stated when it
is granted, and within it nothing is re-asked.

## Four separate acts

| Act | Trigger | Effect | Does not |
| --- | --- | --- | --- |
| Delivery | Explicit authorization | Push, pull request, merge | Mark work complete |
| Completion | Explicit request | Completion note, date, knowledge reconciliation | Remove a worktree or branch |
| Worktree removal | Explicit request | Remove one worktree, preserving files | Delete the branch |
| Branch deletion | Ordinary Git, separately | Delete the branch | — |

Collapsing any two of these produced recurring surprise — a merged pull request
that silently marked work done, or a completion that removed a worktree still in
use — so they stay apart and say so in each direction. Delivery alone does not
mark work complete; an explicit completion request appends the note *and*
reconciles the durable knowledge that work changed, and the request is finished
only when both are done.

## A gate says what it opened

Approval and completion both hand work back rather than finish it, and both have
been read as finished. An agent that recorded an approval stopped there and
created no plan; one that recorded a completion wrote that no durable knowledge
had changed, then reconciled four notes an hour later. In each case the rule was
in the always-loaded entry instruction and the command had answered with a bare
success, so success is what the agent acted on.

So each gate names its own consequence in its own output: `planning_required` on
approval, and `reconcile_required` on completion whenever the plan's repositories
hold catalog entries, the way a dispatch names `setup_required`. A field whose
name is the obligation arrives at the moment that obligation comes due, which a
paragraph read once at session start does not.

This is the same reason [workspace files and safe
editing](../architecture/workspace-files.md) reports a resolved base branch
instead of leaving it in a file: what the CLI knows and never states is
reconstructed by the agent from whatever is nearest, and a gate that states only
`ok` leaves the agent to infer that nothing follows.

## Delivery mechanics and attribution

Delivery uses ordinary Git and provider tools; the CLI supplies repository,
branch, and base information and never silently delivers. The branch the work
started from — this machine's recorded base, or the repository's shared default
branch when the binding records none — is the default pull-request target unless
a person chooses another, delivery can happen per repository, and drift is
explained rather than resolved by an automatic rebase-and-recheck behavior.
Review is not a precondition for any of it.

Agent attribution, credit, co-author, or generated-by text never goes into
commits, pull requests, reviews, or comments. Repository commit conventions are
followed, defaulting to a typed, scoped, imperative subject. The obligation
includes inspecting authored messages and removing injected attribution before
reporting delivery complete, because some hosts add it and "I did not write it"
is not the same as "it is not there".

Owner:

- `context-circuit-source@product/AGENTS.md.in` — which acts need a person and
  what one authorization covers.
- `context-circuit-source@internal/workspace/records.go` `Store.Note` — the
  operations that record approval and completion.
- `context-circuit-source@internal/cli/cli.go` — the obligation each gate's
  result names, beside the record it just wrote.
