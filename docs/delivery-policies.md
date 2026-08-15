# Delivery policies

Delivery is a post-verification choice. A configured policy helps the
coordinator prepare the next action; it never turns successful tests or a
verifier handoff into permission to merge, push, publish, or deploy.

## `team-review`

Use this policy when the result should remain reviewable on an implementation
branch. After implementation and independent verification, the coordinator
may prepare the commit and push path for the configured repository and branch.
The relevant authorization must be current and explicit before any commit or
push. The handoff names the branch, evidence, proposed review destination, and
the remaining human review gate. If authorization is missing, expired, or the
target changed, stop and use the manual fallback.

## `solo-local`

Use this policy when a verified result is intended for local integration. The
coordinator identifies the target from the registered repository's configured
default active branch (or a later branch explicitly confirmed by the user).
It preserves the plan worktree until the human merge gate is satisfied. It
must not silently merge, fast-forward, delete the worktree, or infer consent
from passing verification. A dirty base, branch mismatch, or changed risk
pauses delivery and requests focused confirmation.

## `manual`

Use this policy when the user wants to decide delivery later or when no safe
delivery capability is configured. The verified implementation remains in its
isolated plan worktree, with its evidence and handoff preserved. The
coordinator asks what to do only when a delivery decision is needed; it does
not guess between review, local integration, publication, or deployment.

## Common fallback and review rules

- Configuration is optional and is not requested by initialization.
- Missing, unavailable, denied, or stale configuration falls back to `manual`.
- A policy is reused only while repository, branch, risk, and authorization
  boundaries remain unchanged.
- Policy configuration is inspectable in `workspace.yaml`, but credentials and
  external activity records are never stored there.
- Human gates remain required for commit/push, merge, publication, deployment,
  completion, and any ambiguous ownership takeover.
