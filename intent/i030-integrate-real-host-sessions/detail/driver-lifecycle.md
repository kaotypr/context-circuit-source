# Driver lifecycle and host differences

## Shared lifecycle

Every driver implements six operations with the same meanings:

1. **Probe:** identify the host, version, authentication readiness, supported
   session controls, child-agent capability, event visibility, model controls,
   and usage reporting.
2. **Launch:** create one new root session in the prepared project and submit the
   unchanged human prompt.
3. **Observe:** stream human-visible output and structured native events where
   available.
4. **Wait:** follow the session until its native terminal state or the declared
   harness limit.
5. **Terminate:** stop only the root session and descendants owned by this run.
6. **Collect:** persist the native transcript, event records, child topology,
   usage, and termination evidence the host exposes.

## Capability negotiation

Probe results are run evidence, not permission. A scenario declares the
capabilities needed for its assertions. If the host cannot expose one, the
affected assertion is inconclusive unless an independent observation proves the
same fact.

Drivers translate native concepts into a small normalized vocabulary but retain
the raw host name, version, identifiers, terminal reason, and capability map.
They do not disguise top-level sessions as children or infer a role from a model
name.

## Matrix behavior

A matrix expands one resolved scenario into independent runs, each with its own
workspace and native session. Results are reported individually before any
cross-host summary. One host's success cannot replace another host's missing or
failed evidence.

