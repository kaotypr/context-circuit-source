# Plan 0027 — Member roster and local identity

**Intent:** i021-id-number-blocks  
**Repository:** context-circuit-source  
**Tier:** Standard  
**Status:** done

## Objective

The shared workspace records who may use which number band in a committed roster,
and each machine names its member once in a gitignored local identity file — the
same host-local pattern as repository bindings — so later allocation can resolve
bands without anyone typing a block number in chat.

## Grounding (HEAD 7db954f)

Today `engine.sh` allocates from a single workspace-wide sequence. There is no
roster or `member.local.yaml`. Host-local files already follow
`repositories.local.yaml` and `role-tiering.local.yaml` at the workspace root with
matching gitignore and manifest entries.

Existing intents (through i021) and plans (through 0026) stay as-is; this plan only
adds roster infrastructure and identity resolution. Band-scoped allocation is plan
0028.

## Decisions

**Stack plan 1 of 3** — roster and identity must exist before allocation changes
and before contract docs describe band behavior.

- **members.yaml** — committed roster at workspace root.
- **member.local.yaml** — gitignored local identity at workspace root.
- **Bootstrap** — first member band i001–i099 / 0001–0999 covers existing ids.
- **Missing identity** — runtime fails closed; cc-workspace handles one-time setup.

## Tasks

### MR-001 — Contract schemas and bootstrap roster

Author `members.yaml` and `member.local.yaml` schemas; gitignore and manifest
entries; bootstrap roster for this workspace and template seed.

**Done when:** schemas validate; bootstrap roster exists; template ignores local identity.

### MR-002 — Runtime validate and band resolve

Add roster validate, identity read, and band-resolve helpers in `engine.sh` without
changing allocate-id yet.

**Done when:** overlap and missing-identity cases fail deterministically in tests.

### MR-003 — Workspace guidance and acceptance

Update cc-workspace, cc-intent, coordinator, and getting-started; wire
`test/bands/test-member-roster.sh`.

**Done when:** guidance matches one-time identity setup; acceptance suite green.

## Risks

- Bootstrap roster that does not cover existing numeric prefixes would block the
  maintainer workspace on first band allocation (0028).
- Gitignore omission would ship host-local identity in template releases.

## Acceptance criteria mapping

| Intent criterion | Task evidence |
| --- | --- |
| AC-BAND-03 | MR-003 guidance; MR-002 no numeric prompt on CLI |
| AC-BAND-04 | MR-002 missing identity error; MR-003 setup flow |
| AC-BAND-08 | MR-001 bootstrap preserves existing id prefixes |
