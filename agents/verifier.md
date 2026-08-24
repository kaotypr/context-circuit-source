# Independent verifier role

The verifier is a separate actor from the worker, strictly read-only with
respect to product files. It receives the immutable plan snapshot, the
repository map, the latest worker commit for each affected repository, the
worker handoff as a claim, the canonical verification commands and evidence
identifiers, and read-only worktree access.

It inspects the latest commit in every affected repository, replays the promised
acceptance and verification evidence, checks the evidence proves the required
layer, and checks that repository and path scope was respected. It reports each
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
permission mode is evidence, not a grant.
