# Authorization boundaries

## Two decisions belong to a person

**Approving the intended outcome.** Before detailed code investigation the agent
writes the intent and presents it concretely. The rule that carries the weight:
neither a command, nor an editable approval note, nor another agent can supply
human consent. The approval operation records an actual decision and its date;
it does not constitute one. In the other direction, approval already supplied in
the conversation for that exact outcome is honored rather than re-requested —
asking twice for the same thing teaches people to approve reflexively.

Re-approval is required only when the outcome or its success criteria materially
change. A newly needed file or repository *within* the approved outcome is
explained and recorded, not re-gated, because recorded paths are descriptive
planning information rather than an enforcement contract.

**Authorizing an outward action.** Commit, push, pull-request creation, merge
into a base branch, deployment, external publication, and deletion of workspace
data each require explicit authorization, and authorization already given is
reused rather than re-requested.

## Nothing in between

There is no plan approval gate. A plan is *derived* from an approved outcome, so
approving it again asks a person to ratify a derivation they are not positioned
to check — which is how ratification becomes a rubber stamp. The diagnostic is
likewise not a gate; nothing waits on it.

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
