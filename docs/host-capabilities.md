# Host capabilities

Codex CLI, Claude Code, and Cursor Agent CLI are host adapters around the same
filesystem contract. They discover the natural-language skills in
`.agents/skills/` and may provide child-session primitives or optional
integration adapters. They do not define a second route, lifecycle, ownership
authority, or authorization path.

All three hosts pass the final route `context_set` to the same bounded loader.
The loader is the filesystem enforcement boundary: it rejects undeclared
paths, unsafe links, stale receipts, and byte overruns. A host that cannot
intercept a direct read is not described as sandboxed; the live overread check
fails read-only or returns `host-blocked`.

The shipped discovery names are `cc-entry`, `cc-next`, `cc-plan`, `cc-execute`,
`cc-verify`, `cc-gates`, and `cc-upgrade`.
The human never needs to type these names.

| Host | Child primitive | Safe fallback |
| --- | --- | --- |
| Codex | native child agent or equivalent | filesystem-only workflow |
| Claude Code | subagent/Task equivalent | filesystem-only workflow |
| Cursor Agent | Task/subagent equivalent | filesystem-only workflow |

## Instruction and mode evidence

| Host | Instruction discovery | Bounded mode evidence |
| --- | --- | --- |
| Codex CLI | AGENTS.md, then cc-* skills | interactive entry, native child mapping, receipt-based resume |
| Claude Code | CLAUDE.md imports AGENTS.md | interactive or print entry, Task mapping, receipt-based resume |
| Cursor Agent CLI | AGENTS.md and CLAUDE.md | interactive, print, and resume; local permission settings stay subordinate |

Print mode and force-like host flags do not authorize execution, delivery, or
cleanup. A mutating probe needs a disposable fixture and the corresponding
human gate; the default host fixture is read-only.

## Provider-neutral host evidence

New session, writer, verifier, and handoff records may carry the same bounded
`host_evidence` shape. It records what the host exposed when the record was
created; it is not a permission grant and does not replace a human gate.

```yaml
host_evidence:
  host_id: codex | claude-code | cursor-agent
  observed_version: <version string> | unavailable
  instruction_surface: AGENTS.md | CLAUDE.md | shared-root
  session_role: root | writer | verifier
  root_capability: available | unavailable
  child_capability: available | unavailable
  verifier_capability: available | unavailable
  resume_capability: available | unavailable
  permission_mode: read-only | bounded-write | host-managed
  provider_status: enabled | disabled | denied | unavailable
  offline_fallback: filesystem-only | host-blocked
```

The shape intentionally has no credential, token, provider-payload, transcript,
or host-auth field. A provider marked `disabled`, `denied`, or `unavailable`
still leaves the filesystem workflow usable. A missing required child is
`host-blocked`; it never permits self-verification or a role downgrade.

The normalized host identifiers are `codex`, `claude-code`, and `cursor-agent`.
The root coordinator records `session_role: root`; bounded child packets record
`writer` or `verifier`. The root owns route, lease, worktree, cards, and
consolidation. The writer is the only implementation writer, and the verifier
is independent and read-only.

Host integration is a thin mapping around the engine-generated graph. The root
validates the graph and commit marker, then transports only
`cc_runtime_launch_projection`; writer and verifier roles consume the same
delegation, child-start, receipt, and handoff records. Host-local permission
flags never alter the packet or grant authorization.

If a host cannot create a required writer or verifier child, the route is
host-blocked and remains read-only. The agent does not skip independent
verification or silently self-verify.

If the provider is disabled, denied, or unavailable, retain the filesystem-only
workflow and its evidence checks. `host-blocked` is reserved for a missing
required child primitive; it is not a reason to downgrade the role or broaden
the launch projection.

## Optional question-prompt primitive

Cursor Agent `AskQuestion`, Claude Code `AskUserQuestion`, and Codex CLI
`request_user_input` (when listed) are optional native UI primitives some
hosts expose on the current session. `cc-plan` may use one, after the Review
Card in `docs/plan-review.md`, to present up to three focused human decisions.
The primitive is presentation only: it selects no route, satisfies no gate,
and is never a required child. A missing, denied, or failed prompt falls back
to the card text and is not `host-blocked`. No `host_evidence` field records
the prompt, the chosen option, or any transcript.

## Outcome labels

`available` describes an observed capability, not authorization. Use
`filesystem-only` for a provider that is disabled, denied, unavailable, or not
needed. Use `host-blocked` only when a required child primitive cannot be
created. Optional live smoke evidence is labeled `pass`, `unavailable`, or
`host-blocked`; offline tests never invoke a provider.
