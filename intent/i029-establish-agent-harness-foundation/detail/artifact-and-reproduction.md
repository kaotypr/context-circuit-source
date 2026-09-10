# Artifact identity and reproduction

## Artifact selection

A run accepts exactly one immutable Context Circuit input: a released artifact,
an assembled artifact tied to a source revision, or an explicitly supplied local
development artifact. The manifest records its source, version when available,
revision, and content digest.

The materialized copy is the only Context Circuit product surface visible to the
test workspace. Changing the selected artifact produces a different run identity.

## Retention

- A passing run normally retains its manifest, normalized observations, result,
  and host-independent evidence.
- A failed, inconclusive, interrupted, or explicitly kept run retains the entire
  run root.
- Cleanup targets only the exact run it owns and never deletes another session's
  state.
- Credentials, host configuration, and private provider payloads are never
  retained with a run.

## Reproduction

Every retained result supplies a command or equivalent invocation that selects
the scenario revision, resolved fixture digests, Context Circuit artifact,
prompt, and supported deterministic seeds. Reproduction constructs a fresh test
world; it does not resume the original host conversation or mutate the preserved
evidence.

If an exact host version or artifact is unavailable, reproduction reports that
difference before launch rather than claiming an identical run.

