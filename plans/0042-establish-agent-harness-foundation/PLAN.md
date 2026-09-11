# Plan 0042 — Establish the standalone Agent Harness foundation

**Intent:** i029-establish-agent-harness-foundation  
**Repository:** agent-harness  
**Tier:** Standard  
**Status:** draft

## Objective

Build the greenfield repository into an independently packaged command-line
product. It resolves complete conditions from small versioned fixtures,
materializes one immutable Context Circuit artifact and deterministic Git
repositories into a unique disposable world, freezes exact inputs and baseline
observations, and safely retains, cleans, or reproduces that world. Host drivers
and grading remain outside this foundation.

## Grounding

The target is clean at `main` revision
`985acabffd575b58ed4e171f7d6bc7e99e039807`. Its only tracked file is
`README.md`, which already owns the product boundary: independent development
and consumption of Context Circuit only as an external immutable artifact.
There is no language, package, API, test, or release surface to migrate.

## Decisions

- **One plan:** contracts, resolution, construction, lifecycle, and the public
  end-to-end proof are one execution and verification seam.
- **Python 3.12 with canonical JSON:** the empty repository imposes no toolchain;
  standard-library filesystem, hashing, subprocess, and serialization keep the
  package independent, with Git as its only runtime executable.
- **Durable record plus disposable `world/`:** pass cleanup preserves manifest,
  observations, result, and reproduction inputs, while every questioned or kept
  outcome preserves the complete run.
- **Fresh reproduction:** recorded selectors and digests construct a new world;
  missing or changed artifact and host/tool facts are reported before building.

## Tasks

| id | title | depends on |
|---|---|---|
| AH-001 | Establish the package and declarative contracts | — |
| AH-002 | Build deterministic repositories and isolated worlds | AH-001 |
| AH-003 | Complete retention, cleanup, and reproduction | AH-002 |
| AH-004 | Prove the foundation end to end | AH-003 |

Canonical paths, changes, acceptance statements, and runnable checks are in
`plan.yaml`.

## Risks

- Cleanup owns recursive deletion: it must resolve a known run id and delete
  only that run's `world/`, never an arbitrary path.
- A fixture name or artifact path is not immutable; canonical bytes, revisions,
  and content digests must be frozen and the copied artifact reverified.
- Git determinism requires controlled authors, timestamps, branches, ordering,
  and file bytes, all asserted through Git itself.
- Host state can leak through diagnostics. All portable paths are run-relative,
  and sensitive fields and private payloads are rejected or scrubbed.
- Resolver, copy, Git, validation, and snapshot failures remain harness errors;
  this plan adds neither drivers nor synthetic agent results.

## Verification

Run `python3 -m unittest discover -s tests -v`,
`python3 -m compileall -q src tests`, public CLI help/version checks, and the
focused end-to-end story. The story builds multiple worlds from identical
inputs, proves isolation and deterministic history, applies both retention
branches, and reproduces into a fresh third world.
