# dist-build-version

Stop `scripts/build-dist.sh` from stamping a stale, hardcoded version onto the
local dist artifact; derive its default version from the single source of truth.

- [design.md](./design.md) — the normative design: why the dev dist build ships
  the wrong version today, the fixed decision to derive the default from
  `wrapper/manifest.yaml`, and the exact surface that changes.

Read [design.md](./design.md) first; it is complete on its own.
