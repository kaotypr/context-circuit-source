# Workflow

Workspace mode, human gates, and wrapper change policy are configured in
`workspace.yaml`. A direct request to run known work satisfies task selection.
Plans and activity integrations are optional.

When a plan is used, `create-plan` writes only a numbered draft. Approval must be
explicit and recorded in the plan index. Material revisions revoke approval,
increment the plan version, and preserve existing work IDs. Publishing external
tasks and beginning implementation are separate actions.

`whats-next` is read-only. Its work-state projection combines configured
activity facts with validated plan-linked runtime manifests, closeout records,
and tracked completed-work contributions. Durable outcomes take precedence over
closeout, runtime, and activity evidence, but disagreement produces an explicit
reconciliation action instead of silently choosing a state. It ranks urgent
work, actionable work in progress or closeout, review or verification failures,
ready source tasks, and then dependency-ready approved-plan items. Passed work
is review-ready, closing work is closeout-ready, completed and cancelled work is
excluded, and unknown dependency status does not count as completion. If
nothing is executable, it recommends the smallest source-backed enabling
action. Recommendation never claims work; explicit `run-task` selection is a
separate human gate.

`run-task` prepares isolated work but does not launch a provider SDK, push,
open a pull request, merge, or deploy. A fresh worker implements in its assigned
worktree. A separate fresh verifier inspects the result without modifying it.
The host-neutral result recorder validates both roles and advances the runtime
manifest through `running`, `verifying`, and a verifier outcome. Runtime evidence
remains ignored and preserved until deliberate closeout.

Passing work advances through local review, optional explicitly authorized
publication, human merge, and verified merge confirmation. Missing `origin`
still permits local review. Review handoffs use exact recorded commits and
shell-safe commands. Only a human runs the merge; `confirm-merge` is read-only
with respect to Git and requires the exact reviewed head and reported merge to
be reachable from the configured default target before closeout is ready.

After human-confirmed merge or deliberate abandonment, `finish-work` first writes
an append-only contribution and recoverable closeout record. Worktree cleanup is
a separate explicit rerun after that contribution is tracked, committed, and
clean. Cleanup never deletes branches or runtime evidence and refuses dirty or
unpushed work.

New task briefs distinguish implementation scope from test expectations. Required
test paths become worker edit scope and must appear in the recorded Git diff.
Without authorized test scope, the default policy is verifier-only and requires
independent acceptance evidence.

Ignored clones are recommended. Each ignored clone must have an exact
`.gitignore` entry so tracked submodules can coexist under `repositories/`.

Product Knowledge is an optional, incremental business-context view under
`context/` (product map, roles, domains, and workflow pages). It is validated when
present and absent-safe when not. Plans declare Product Knowledge impact, tasks
carry a bounded immutable context package to worker and verifier, evidence reports
observed impact, and canonical current-behavior pages change only after a
dedicated confirming role declares the behavior effective. See
`docs/product-knowledge.md`.

Configured activity actions use semantic lifecycle events. Required actions
block their transition on failure; optional actions warn and degrade; manual
actions present their exact description for human completion. The active host
uses authorized MCP/CLI tools and records confirmed results with stable
idempotency keys. Starting actions finish before branch or worktree creation,
and completion or cancellation actions finish before closeout.
