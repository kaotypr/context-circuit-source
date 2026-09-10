# Normal role envelopes

## Coordinator

The coordinator receives the human conversation and small project entry surface.
It reads the Product Knowledge index, only the accepted knowledge selected for
the request, and the single procedure needed for the routed action. It invokes
deterministic runtime interfaces without reading their implementation.

During intent creation it does not read target code. During context gathering it
reads only evidence explicitly bounded by the human. Ordinary requests do not
authorize archived plans or unnamed source material.

## Planner

The planner receives the approved intent and contract, optional intent detail,
relevant selected Product Knowledge, and one repository identity. It reads that
repository's agent guidance and the code, tests, configuration, and documentation
needed to ground the approved outcome.

It does not receive the full coordinator conversation, unrelated knowledge,
another repository, archived plans, or unrelated source evidence. Separate
repositories receive separate planner children and envelopes.

## Worker

The worker receives one immutable execution packet, only the Product Knowledge
named by the plan, and one assigned repository worktree. It reads the repository
guidance, planned implementation areas, necessary same-repository dependencies,
and tests needed for the work.

It does not receive unrelated plans, workspace history, coordinator conversation,
another repository, or verifier context. A necessary same-repository read beyond
a task anchor remains attributable and must satisfy the role's scope-recording
contract; it is not permission to inventory the repository indiscriminately.

## Verifier

The verifier receives the immutable plan snapshot, repository map, latest worker
commit, worker handoff as a claim, and named verification commands and evidence
identifiers. It has read-only repository access and reads the complete committed
diff plus dependencies needed to verify the declared outcomes.

It does not receive the worker's private session, coordinator conversation,
unrelated Product Knowledge, other plans, or writable product access. Reads
beyond original anchors must be justified by the committed diff, a recorded
scope expansion, or a verification dependency.

