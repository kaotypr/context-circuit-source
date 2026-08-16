# Delivery policies

Delivery is a post-verification choice. A configured policy helps the
coordinator prepare the next action; it never turns successful tests or a
verifier handoff into permission to merge, push, publish, or deploy.

The live policy IDs are `remote-review`, `local-target`, and `manual`. They
name the path to prepare, not workspace identity and not an authorized git
action. When reading `workspace.yaml`, `team-review` remains a readable alias
for `remote-review` and `solo-local` remains a readable alias for
`local-target`. An old ID is not missing configuration and does not fall back
to `manual`. New writes use the new IDs.

## `remote-review`

Use this policy to prepare a reviewable remote path toward the configured
target branch. After implementation and independent verification, the
coordinator may prepare that remote path for the configured repository and
branch. The relevant authorization must be current and explicit before any
commit or push. The handoff names the branch, evidence, proposed review
destination, and the remaining human review gate. If authorization is
missing, expired, or the target changed, stop and use the manual fallback.

## `local-target`

Use this policy to prepare local integration into the configured target
branch. The coordinator identifies the target from the registered
repository's configured default active branch (or a later branch explicitly
confirmed by the user). It preserves the plan worktree until the human merge gate
is satisfied. It must not silently merge, fast-forward, delete the
worktree, or infer consent from passing verification. A dirty base, branch
mismatch, or changed risk pauses delivery and requests focused confirmation.

## `manual`

Use this policy to decide the delivery path later, or when no safe delivery
capability is configured. The verified implementation remains in its
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

## Publication

Publication is an optional human-authorized action after a plan is approved and
before external execution or collaboration. The human explicitly authorizes it.
The agent must not publish, merge, deploy, create external issues, or
synchronize external status implicitly.

Publication does not approve a plan, change plan or task status, start a
session, monitor external status, or synchronize completion back into the
workspace. It stores no implicit external status. Publication is outside the
filesystem execution protocol. If a host or provider integration is later
used, it remains an explicitly authorized adapter and must not redefine
session, plan, or completion state.
