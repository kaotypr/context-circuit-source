# Workflow

This template is in team mode. The human gates are plan approval when a plan is
used, explicit task selection, and repository merge. A direct request to run
known work satisfies task selection. Plans and activity integrations are
optional.

`run-task` prepares isolated work but does not launch a provider SDK, push,
open a pull request, merge, or deploy. A fresh worker implements in its assigned
worktree. A separate fresh verifier inspects the result without modifying it.
The host-neutral result recorder validates both roles and advances the runtime
manifest through `running`, `verifying`, and a verifier outcome. Runtime evidence
remains ignored and preserved until deliberate closeout.

New task briefs distinguish implementation scope from test expectations. Required
test paths become worker edit scope and must appear in the recorded Git diff.
Without authorized test scope, the default policy is verifier-only and requires
independent acceptance evidence.

Ignored clones are recommended. Each ignored clone must have an exact
`.gitignore` entry so tracked submodules can coexist under `repositories/`.
