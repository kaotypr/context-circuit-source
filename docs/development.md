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
at `.dist/context-circuit-0.2.0/`, and produces the authoritative release archive
at `.dist/context-circuit-0.2.0.tar.gz`.

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
