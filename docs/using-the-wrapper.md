# Using the wrapper

Context Circuit coordinates delivery; product repositories continue to own
their code, build systems, and repository-specific instructions. Work begins
only after a human selects a concrete request or approves a plan.

## Choose the next action

Use `$cc-create-plan` for broad or ordered work. After a numbered plan is
explicitly approved, use `$cc-execute-plan` as the core implementation workflow.
Use `$cc-run-task` only for a rare, manually selected single-task escape hatch.

`$cc-whats-next` can recommend one source-backed action from approved plans and
configured read-only work sources. It does not claim or start work.

## Run scoped work

`$cc-execute-plan` checks that every affected base repository is clean, derives
the approved task graph, and
and prepares one branch and isolated Git worktree per repository. A fresh worker
receives only its authorized scope and acceptance criteria. A separate verifier
reviews the committed result without modifying it.

For work spanning repositories, `execute-plan` orders dependencies and keeps a
dependent worker locked until its declared prerequisites pass independently.
`run-task` never advances plan-level lifecycle, review, or closeout.

Failed verification may produce a bounded repair attempt. The repair remains in
the same isolated worktree and cannot expand the authorized scope silently.

## Review and merge

A passing run prepares exact local review commands even when no `origin` exists.
Remote publication remains optional and explicitly authorized; the workflow
records only confirmed publication results. It does not push, open a pull
request, merge, or deploy automatically. Review the branch, verification
results, and runtime handoff, then use the normal repository review process.
Humans retain the merge decision and report its full commit through the emitted
`confirm-merge` command. Closeout becomes ready only after the exact reviewed
head is verified on the configured default target.

## Finish safely

After confirmed merge or deliberate abandonment, invoke `$cc-finish-work`. It
records an append-only contribution and closeout evidence before any optional
worktree cleanup. Dirty, unpushed, or otherwise unrecorded work is preserved.

Runtime evidence lives under ignored `.runtime/`. Do not delete it manually.
Successful cleanup removes only the registered clean worktree and retains its
branch and all runtime evidence. Cleanup blockers are returned as an ordered
checklist ending with the exact safe rerun. Runtime pruning is a separate future
operation, not part of normal closeout.

## Keep durable context useful

Use `$cc-gather-context` to resolve and cite the minimum authoritative source
material without changing repositories, external systems, or durable context.
Use `$cc-sync-context` to turn completed-work contributions into a focused,
reviewable update to `context/`. Product facts belong in `PROJECT.md`, structural
facts in `ARCHITECTURE.md`, established practices in `CONVENTIONS.md`, and durable
choices with rationale in `DECISIONS.md`.

Live task status belongs in the configured activity system, not canonical
context. Repository-specific details belong in the corresponding product
repository.

## Product Knowledge

Optionally, keep a lean business view — product map, roles, domains, and workflow
pages under `context/` — that gives humans and agents task-relevant business
context without loading the whole application. It is additive and incremental: a
wrapper without it stays valid, and plans, tasks, evidence, and onboarding packs
carry Product Knowledge impact when it exists. Canonical current-behavior pages
change only after a dedicated confirming role declares the behavior effective. See
[product-knowledge.md](product-knowledge.md).

## Safety boundaries

Context Circuit does not discard dirty work, embed credentials, run a background
service, merge, or deploy. Creating remotes, pushing branches, publishing tasks,
opening pull requests, and cleanup require the relevant human authorization.
