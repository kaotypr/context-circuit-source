---
name: cc-configure-workspace
description: Configure an optional delivery policy, host capability, or external activity adapter after workspace identity is established.
---

# Configure workspace behavior

Use this skill when the user explicitly asks to configure delivery behavior,
an optional integration, or a host-specific convenience. It is a
capability-led conversational workflow, not a command console and not part of
`cc-initialize-workspace`.

## Read and discover

Read `AGENTS.md`, `WORKFLOW.md`, `workspace.yaml`, `context/INDEX.md`,
`docs/configuration.md`, and only the relevant policy, integration, or host
document. Inspect the registered repository and current branch only when the
requested capability needs a concrete target. Do not scan `sources/` or ask
unrelated optional questions.

Recognize requests for:

- a delivery policy: `remote-review`, `local-target`, or `manual`. When
  reading an existing file, treat `team-review` as `remote-review` and
  `solo-local` as `local-target`;
- an external activity record integration; or
- a host capability mapping for Codex, Claude Code, or Cursor Agent.

If the request is ambiguous, ask one focused question about which capability
the user wants. Do not configure a capability merely because it is available.

## Ask the minimum relevant questions

For delivery, determine only the selected policy, affected repository or
workspace scope, target branch when it is not already the registered default,
and the current authorization state needed by that policy. Ask with path
language: prepare a reviewable remote path toward the target branch
(`remote-review`), prepare local integration into the target branch
(`local-target`), or decide the delivery path later (`manual`). Do not
present workspace `mode` values as delivery policy IDs. Confirming this
configuration does not grant commit, push, PR, or merge.

For an integration, determine only the provider/capability, data it may read,
data it may write, whether the user wants it enabled, and whether the
requested external mutation is authorized. Explain that credentials stay in
host-managed secure storage and are never entered into workspace files.

For a host mapping, determine only the host and requested convenience. Keep
the portable capability name, ownership boundary, verification behavior, and
human gates unchanged.

## Explain before applying

Before changing `workspace.yaml`, show a concise proposal containing:

1. capability and scope;
2. behavior and affected repositories/branches;
3. required human authorization and what remains gated; and
4. safe fallback if the host, provider, network, or authorization is absent.

Wait for explicit user confirmation before writing the optional
`configuration` block. A configured policy is durable and reused on later
runs, but a changed repository, target branch, risk, or authorization state
requires focused reconfirmation.

## Durable record and safety

Write only portable intent under `workspace.yaml`:

- delivery policy, scope, target branches, authorization state, and `manual`
  fallback;
- integration enablement, provider-neutral read/write boundary,
  authorization state, and `core-filesystem` fallback; and
- host capability name and whether native support or the core fallback is
  used.

Never write credentials, tokens, secrets, provider payloads, external activity
records, or transient runtime ownership into the configuration. Do not change
plan/task status, leases, sessions, worktrees, or verification evidence.

## Delivery behavior

Apply delivery only after implementation and independent verification:

- `remote-review` prepares a reviewable remote path toward the target branch
  only after the relevant authorization and leaves a review handoff;
- `local-target` identifies the configured target branch and pauses at the
  human merge gate; and
- `manual` leaves the verified implementation in the plan worktree and asks
  the user's delivery choice only when needed.

No configuration grants automatic merge, push, publication, deployment, or
completion authority. If the configuration is absent, denied, stale, or the
capability is unavailable, use `manual` and continue the core filesystem
workflow.

## Integration behavior

External activity adapters are opt-in and disabled by default. They may
publish only the explicitly authorized plan summary, task progress, or handoff
metadata described in the configuration. They must report `disabled`,
`denied`, or `unavailable` without blocking the core plan, runtime, or
verification evidence. Host-neutral fallback remains available offline.

## Output

Report the requested capability, selected scope, proposal and confirmation,
durable configuration path, authorization state, fallback, and any remaining
human gate. If no confirmation was given, leave configuration unchanged and
report the next focused decision.
