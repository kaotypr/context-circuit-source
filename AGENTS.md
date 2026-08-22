# Context Circuit source-repository safety

This checkout is the Context Circuit product source, not an instantiated
workspace. The active source specification is
`sources/context-circuit-design/DESIGN-SPEC.md`; do not turn it into a
canonical workspace plan or restore the obsolete Context Circuit lifecycle,
routing model, or legacy skill behavior.

The shipped product contract is under `wrapper/`; the mutable blank seed is
under `template/`. Product behavior has one owner per rule: use the owner map
in `wrapper/contracts/invariants.yaml` and do not add parallel policy to a
skill or role file.

Repository safety still applies:

- preserve unrelated and dirty work;
- do not read or scan `sources/` unless a request names exact files;
- do not store credentials or provider payloads;
- do not commit, push, merge, publish, deploy, or delete user data implicitly;
- in this maintainer source checkout, an explicit user request may authorize a
  source-only commit; registered product repositories still require their
  separate delivery and publication gates;
- do not modify runtime state belonging to another session;
- use `wrapper/runtime/engine.sh` only as the host-neutral implementation
  library and keep the human interface conversational;
- run the semantic acceptance suite after meaningful phases.

Repository-local behavior belongs to product repositories registered in an
instantiated workspace. Maintainer files and release assembly are source-only.

Host adapters

Codex CLI, Claude Code, and Cursor Agent CLI use this shared instruction
surface. Host identity, observed version, capabilities, role, permission mode,
provider status, and offline fallback belong in the bounded `host_evidence`
shape owned by `wrapper/contracts/schemas/`; they never authorize a route or
gate. Native child features map to the existing coordinator, writer, and
independent verifier packets. If a required child is unavailable, preserve the
read-only `host-blocked` outcome and do not self-verify.
