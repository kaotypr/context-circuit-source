# Command reference

`.agents/bin/cc.mjs` is the dependency-free deterministic engine shared by
Codex and Claude Code. Skills decide what should happen; the command validates
and records how it happens.

It provides workspace validation and bootstrap, task normalization, Git
worktree preparation, result recording, bounded repair, review handoff,
planning, lifecycle records, verified human merge confirmation, closeout, recommendations, publication records,
and context synchronization.

Run it from the wrapper root:

```bash
node .agents/bin/cc.mjs validate
node .agents/bin/cc.mjs configure-workspace --check-only
node .agents/bin/cc.mjs whats-next
```

Humans normally invoke host skills rather than assembling low-level command
arguments. The principal skill-to-command mappings are:

| Human action | Codex | Claude Code | Deterministic command |
| --- | --- | --- | --- |
| Configure | `$configure-workspace` | `/configure-workspace` | `configure-workspace` |
| Initialize (compatibility) | `$initialize-workspace` | `/initialize-workspace` | `initialize-workspace` |
| Gather context (read-only) | `$gather-context` | `/gather-context` | host read-only tools |
| Plan | `$create-plan` | `/create-plan` | `create-plan` |
| Recommend work | `$whats-next` | `/whats-next` | `whats-next` |
| Run scoped work | `$run-task` | `/run-task` | `run-task` |
| Finish a run | `$finish-work` | `/finish-work` | `finish-work` |
| Curate context | `$sync-context` | `/sync-context` | `sync-context` |

Additional recorder and preparation subcommands are internal workflow steps.
Use them only when a skill or generated runtime instruction provides the exact
arguments. Run `validate --check-paths --check-documents` for the strongest
general workspace check.

The bundle does not call an AI model, store credentials, run a service, push,
merge, or deploy. Host adapters remain responsible for fresh agent sessions and
explicitly authorized external tools.
