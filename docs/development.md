# Develop Context Circuit

This guide is for maintainers of the Context Circuit source repository, not
users of a generated wrapper.

## Source checkout

The source repository contains TypeScript, npm metadata, fixtures, and tests used
to develop the deterministic command:

```bash
npm ci
npm run check
npm run build:template
```

The build writes `.agents/bin/cc.mjs`, creates an inspectable staging directory
at `.dist/context-circuit-0.2.1/`, and produces the authoritative release archive
at `.dist/context-circuit-0.2.1.tar.gz`.

## Source versus output

Maintainer-only inputs:

- `package.json`, `package-lock.json`, and `tsconfig.json`
- `node_modules/`
- `scripts/`
- `test/`
- `fixtures/`
- historical host-proof reports

The distributable contains the static wrapper, canonical contracts and skills,
thin host adapters, human documentation, and the bundled command. A wrapper user
does not run `npm install`.

Before a release, run the complete checks, validate every canonical skill,
inspect the generated file inventory and complete Git diff, and confirm the
archive excludes credentials, local metadata, maintainer inputs, and placeholder
repository registration. Publish the generated archive rather than packaging the
staging directory manually.

## Release synchronization

`.github/workflows/sync-context-circuit-release.yml` runs only for `v*` tags or
an explicit dispatch. It validates the source, builds and inspects the archive,
then synchronizes the version-matched tarball and npm metadata to
`kaotypr/context-circuit-release` after that repository's tests and package
inventory pass. A repeated identical release is a no-op; a different tarball for
an existing version is rejected.

A repository administrator must create the
`CONTEXT_CIRCUIT_RELEASE_TOKEN` Actions secret before activation. Scope the token
to pushing only `kaotypr/context-circuit-release`; the workflow does not use the
source repository's `GITHUB_TOKEN` for that cross-repository push and does not
persist the release token in either checkout.

The destination repository must already exist with a `main` branch, a
`context-circuit` package manifest and lockfile, and `npm test` plus
`npm run test:package` scripts. Until those prerequisites and the secret exist,
the source workflow is safe to merge but release synchronization will stop at
the destination checkout or validation step without publishing anything.
