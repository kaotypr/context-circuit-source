# Create-plan host compatibility results

Date: 2026-08-11

The canonical `create-plan` workflow was exercised in isolated, ignored wrapper
copies using Codex CLI 0.147.0 and Claude Code 2.1.220. Both wrappers contained
the same retry-observability brief, registered frontend fixture, contracts,
scripts, and thin host adapter.

## Outcomes

Both hosts:

- Read their thin adapter and the canonical `.agents/skills/w-create-plan/SKILL.md`.
- Normalized the same source into a valid draft request.
- Created exactly eight plan files without touching product code.
- Produced a draft version 1 with `RETRY-001`, `RETRY-010`, and `RETRY-020`.
- Preserved the requested parent and dependency edges without a live-status
  column or fabricated external references.
- Passed plan, workspace-document, and diff validation.
- Recorded explicit approval as version 1 with matching approved and material
  digests, preserving a runtime evidence copy of the approved index.
- Added the same reduced-motion acceptance criterion as a material revision.
- Returned the plan to draft version 2, cleared approval metadata, refreshed the
  material digest, and preserved all work IDs.
- Refused a second creation attempt and left all eight plan-file hashes
  unchanged after the refusal.
- Avoided task publication, activity mutation, implementation, branches,
  worktrees, commits, and pushes.

Independent reconciliation of both final plan directories and captured approved
indexes returned zero errors. Git state in each wrapper contained only the eight
untracked files under `context/plans/retry-observability/`; both product fixture
repositories remained clean.

## Host observations

- Plan prose varied appropriately between hosts while the machine-readable
  lifecycle, filenames, work IDs, and relationships remained compatible.
- Codex streamed its actions and Claude Code buffered output until completion.
- Temporary normalized requests were stored outside durable plan content: Codex
  used a temporary path, while Claude used ignored `.runtime/` evidence.
- An ambiguously phrased proof approver produced different inferred identities:
  Codex recorded `kao`; Claude Code recorded the session email. Both values were
  contract-valid, but identity inference is not portable. The canonical skill
  now requires an exact human-provided `approved_by` value and forbids inferring
  it from host, Git, or email identity.

The host-specific adapters required no workflow logic changes.
