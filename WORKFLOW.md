# Workflow

This template is in team mode. The human gates are plan approval when a plan is
used, explicit task selection, and repository merge. A direct request to run
known work satisfies task selection. Plans and activity integrations are
optional.

When a plan is used, `create-plan` writes only a numbered draft. Approval must be
explicit and recorded in the plan index. Material revisions revoke approval,
increment the plan version, and preserve existing work IDs. Publishing external
tasks and beginning implementation are separate actions.

`whats-next` is read-only. It ranks urgent work, actionable work in progress,
review or verification failures, ready source tasks, and then dependency-ready
approved-plan items. Unknown dependency status does not count as completion.
If nothing is executable, it recommends the smallest source-backed enabling
action. Recommendation never claims work; explicit `run-task` selection is a
separate human gate.

`run-task` prepares isolated work but does not launch a provider SDK, push,
open a pull request, merge, or deploy. A fresh worker implements in its assigned
worktree. A separate fresh verifier inspects the result without modifying it.
The host-neutral result recorder validates both roles and advances the runtime
manifest through `running`, `verifying`, and a verifier outcome. Runtime evidence
remains ignored and preserved until deliberate closeout.

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

Configured activity actions use semantic lifecycle events. Required actions
block their transition on failure; optional actions warn and degrade; manual
actions present their exact description for human completion. The active host
uses authorized MCP/CLI tools and records confirmed results with stable
idempotency keys. Starting actions finish before branch or worktree creation,
and completion or cancellation actions finish before closeout.
