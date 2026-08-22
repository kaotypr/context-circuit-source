# Wrapper migrations

The manifest version is compared with the live session and context receipt.
`compatible` updates are wrapper-only; `migration-needed` changes require a
current human migration confirmation; `blocked` state remains read-only.
Records without a version are `legacy-unknown` and stay readable.

The v1-to-v2 plan migration in `v1-to-v2.yaml` is additive. It preserves plan
intent, status, context, runtime, dirty work, repository identity, local
bindings, and registered repositories. Rollback restores only wrapper-owned
files, including the shipped root CLAUDE.md adapter, and never overwrites
`repositories.local.yaml` or repository contents. Optional host evidence
fields are additive; unknown or contradictory host state remains read-only.
