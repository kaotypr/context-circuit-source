# Authorization boundaries

## Three decisions belong to a person

**Approving the intended outcome.** Before detailed code investigation the agent
writes the intent, presents it, and stops. What it is written against, and how a
decision it cannot make reaches the person as a numbered open question, is
[grounding an intent](intent-grounding.md). The rule that carries the weight:
neither a command, nor an editable record, nor another agent can supply human
consent. The approval operation records an actual decision and its date; it does
not constitute one.

The intent's `approved_at` field is stamped with the instant there, and the person's
words are kept under a heading carrying the same instant; an unapproved intent
carries no such field, so absence is the only empty state. The field is a readable
record of the gate, never the gate itself — anything that can write the record can
write the timestamp, so an intent is not approved because a field says so. Nothing
in the executable reads it to decide whether to proceed.

Nor does the request that prompted the intent. Reading a change request as
consent to whatever intent gets derived from it collapses the gate entirely —
the agent writes the outcome and approves it in the same breath — so approval is
a person's answer to the intent as written.

Re-approval is required only when the outcome or its success criteria materially
change. A newly needed file or repository *within* the approved outcome is
explained and recorded, not re-gated, because recorded paths are descriptive
planning information rather than an enforcement contract. Asking twice for the
same thing teaches people to approve reflexively.

**Starting execution.** Plans are derived, presented, and then left alone until a
person asks for them to be run. This is a decision about *when work begins*, not
a review of the derivation, which is why it sits beside plan presentation without
being a plan approval gate.

It is deliberately not pre-suppliable. An instruction to implement that arrived
before the plans existed cannot be a decision about plans nobody had seen, so it
does not start execution — otherwise "add billing and implement it" carries the
agent from an unread intent through unread plans into repository edits, which is
the failure the two stops exist to prevent.

**Authorizing an outward action.** Commit, push, pull-request creation, merge

into a base branch, deployment, external publication, and deletion of workspace
data each require explicit authorization, and authorization already given is
reused rather than re-requested.

## What the execution request is not

It is not plan approval. A plan is *derived* from an approved outcome, so
approving it again asks a person to ratify a derivation they are not positioned
to check — which is how ratification becomes a rubber stamp. Reading the plans
stays optional, and a person who asks to execute without reading them has
withheld nothing. The two are separable because they answer different questions:
whether the derivation is right, and whether now is the time to run it. Only the
second is asked for. The diagnostic is not a gate either; nothing waits on it.

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

## Delivery mechanics and attribution

Delivery uses ordinary Git and provider tools; the executable supplies
repository, branch, and base information and never silently delivers. The
recorded base branch is the default pull-request target unless a person chooses
another, delivery can happen per repository, and drift is explained rather than
resolved by an automatic rebase-and-recheck behavior. Review is not a
precondition for any of it.

Agent attribution, credit, co-author, or generated-by text never goes into
commits, pull requests, reviews, or comments. Repository commit conventions are
followed, defaulting to a typed, scoped, imperative subject. The obligation
includes inspecting authored messages and removing injected attribution before
reporting delivery complete, because some hosts add it and "I did not write it"
is not the same as "it is not there".

Owner: `context-circuit-source@product/AGENTS.md.in`; the recording operations
are `context-circuit-source@internal/workspace/records.go` `Store.Note`.
