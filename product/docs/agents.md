# Roles and model/effort settings

The workspace agent is the coordinator. It chooses bounded tasks, prepares working
copies, dispatches subagents through the coding host, waits for results, and
integrates their work. The CLI resolves settings and writes native role definitions;
it has no LLM credentials and does not invoke model APIs.

| Role | Responsibility | Access |
| --- | --- | --- |
| explorer | Answer a specific codebase question with evidence | Read-only |
| planner | Investigate an approved intent and return a grounded plan | Read-only; coordinator writes the plan |
| worker | Implement assigned work and run normal checks | Repository edits within the assigned task |
| reviewer | Independently examine a diff on user request | Read-only; no fixes or automatic follow-up |

Role tiering is a concrete `(model, effort)` pair per role and host in
`.context-circuit/role-tiering.yaml`. All pairs default to `inherit`: the host
chooses its supported defaults. No hardcoded model catalog or cost ladder ships.
These settings describe agent capability, not plan ownership or risk gates.

```sh
context-circuit-cli agent settings
context-circuit-cli agent configure --host codex --role worker --model MODEL_ID --effort high
context-circuit-cli agent setup --host codex
```

Configure the actual host (`codex`, `claude-code`, or `cursor`) and verify that its
account supports the requested pair. Changes are shared workspace preferences.
`agent setup` materializes four native files under that host's agents directory;
these files and their update inventory are local and ignored. A customized native
file is preserved and reported instead of overwritten. Restart/reload the coding
host when it does not discover newly written definitions in the current session.

A worker's brief depends on how it is parallelized. One worker per plan owns its
whole worktree and cannot collide with a sibling plan, whose work arrives later
through an integration merge; several workers inside one worktree share files and
must preserve each other's edits. `--shared` selects the second. `--plan` adds the
plan record, its dependencies, and the fact that their work is already in the
branch's ancestry.

The `cc-dispatch` skill resolves each task, invokes the host's actual agent tool,
waits, and integrates the result. `agent dispatch` alone produces a specification;
`launch_required: true` means it has not launched anything. Available host APIs may
also accept the resolved settings directly without native-file loading.

Codex native TOML uses `model`, `model_reasoning_effort`, and `sandbox_mode` for
read-only roles. Claude Code uses `model`, `effort`, and a read-only tool allowlist;
its reviewer receives the diff from the coordinator. Cursor uses `model` with an
optional `[effort=...]` parameter and `readonly`. For Cursor an explicit effort
requires an explicit model. Inherit omits effort overrides. Actual settings remain
subject to host availability and policy; report substitutions or unsupported pins.

See the maintained host documentation for
[Codex](https://learn.chatgpt.com/docs/agent-configuration/subagents),
[Claude Code](https://code.claude.com/docs/en/sub-agents), and
[Cursor](https://cursor.com/docs/subagents). These mappings were checked against
their documentation; generating a definition does not prove a host loaded it.

Review remains manual. A fresh reviewer gets the full requested diff, exact
base/head references, success criteria, and surrounding code. It never modifies
code or launches repairs. The coordinator asks the user how to handle findings.
