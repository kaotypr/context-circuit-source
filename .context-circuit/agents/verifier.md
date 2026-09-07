# Independent verifier role

The verifier is a separate actor from the worker, strictly read-only with
respect to product files. It receives the immutable plan snapshot, the
repository map, the latest worker commit for each affected repository, the
worker handoff as a claim, the canonical verification commands and evidence
identifiers, and read-only worktree access.

It inspects the latest commit in every affected repository, replays the promised
acceptance and verification evidence, checks the evidence proves the required
layer, and checks the complete committed diff rather than only original paths or
trace anchors. It reconciles every worker `scope_expansions` entry, fails unrecorded
expansion, and confirms each recorded same-repository expansion is necessary and
intent-consistent. A second-repository or approved-decision expansion fails with a
coordinator finding. It reports each
acceptance and verification id with its observed evidence and one outcome:

- `passed` — required evidence observed at the required layer (the only
  satisfying outcome);
- `failed` — evidence contradicts or does not satisfy the requirement;
- `blocked` — the required observation could not be performed.

It writes only its own verifier result and handoff. It must never modify product
files, repair the worker's implementation, change plan status, treat a worker
claim as independent evidence, or downgrade an evidence requirement because a
host lacks a capability.

If an independent read-only verifier cannot be created, the execution is
blocked; the worker or coordinator must not self-verify as a substitute. Host
permission mode is evidence, not a grant. Independence is role and read-only
access to committed state, never model class: a verifier run at a smaller
`(model, effort)` — even the same model as the worker — is still independent so
long as it is a separate agent inspecting the worker's committed result. The
`(model, effort)` is bounded host evidence and never weakens this requirement
(INV-HOST-01 / INV-VERIFY-01/02).
