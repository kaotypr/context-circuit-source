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
| Configure | `$cc-configure-workspace` | `/cc-configure-workspace` | `configure-workspace` |
| Initialize (compatibility) | `$cc-initialize-workspace` | `/cc-initialize-workspace` | `initialize-workspace` |
| Gather context (read-only) | `$cc-gather-context` | `/cc-gather-context` | host read-only tools |
| Import repository context | `$cc-import-context` | `/cc-import-context` | `import-context` |
| Plan | `$cc-create-plan` | `/cc-create-plan` | `create-plan` |
| Recommend work | `$cc-whats-next` | `/cc-whats-next` | `whats-next` |
| Execute approved plan | `$cc-execute-plan` | `/cc-execute-plan` | `execute-plan` |
| Run one manual task | `$cc-run-task` | `/cc-run-task` | `run-task` |
| Approve additional task files | human gate | human gate | `approve-scope-expansion` |
| Finish a run | `$cc-finish-work` | `/cc-finish-work` | `finish-work` |
| Curate context | `$cc-sync-context` | `/cc-sync-context` | `sync-context` |
| Onboarding pack | host skill | host skill | `onboarding-pack` |

`onboarding-pack --roles <role[,role...]>` generates a revision-stamped Product
Knowledge onboarding view. See [product-knowledge.md](product-knowledge.md) for the
full Product Knowledge lifecycle.

`import-context --request <request.json>` discovers bounded evidence from one
registered repository without modifying it and emits a source-cited contribution
snapshot. The `cc-import-context` skill reviews and curates that evidence before
using the existing `sync-context` and `prepare-context-review` commands. See
[context-sync.md](context-sync.md) for the complete handoff and safety boundaries.

Additional recorder and preparation subcommands are internal workflow steps.
Use them only when a skill or generated runtime instruction provides the exact
arguments. Run `validate --check-paths --check-documents` for the strongest
general workspace check.

The bundle does not call an AI model, store credentials, run a service, push,
merge, or deploy. Host adapters remain responsible for fresh agent sessions and
explicitly authorized external tools.
