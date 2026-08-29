# Context Circuit — runtime rearchitecture

A **named grouping**, not a version. It collects the design work around one
question that is too foundational for a single release delta: **what the
deterministic runtime *is* — its form, its invocation surface, how it is tested,
and how it is distributed** — currently `wrapper/runtime/engine.sh`, a POSIX sh
library that is both *sourced* (its `cc_*` functions called directly) and *invoked*
(a thin CLI).

This grouping exists because changing the runtime's substance touches the core of
the product and the model by which the whole product is verified, so it spans more
than one release's worth of change and is sequenced contract-first. Read the v0.6
design ([../v0.6/](../v0.6/)) and the v0.5 base ([../v0.5/](../v0.5/)) first; this
grouping changes only what its scopes name.

## Why a grouping and not a v0.7.0 scope

It **broadens** the question the v0.7.0 `runtime-opacity` scope opened.
runtime-opacity is scoped around a discipline — "the coordinator *must not read*
`engine.sh`" — which presumes the engine stays a readable shell file. The moment
the runtime's *form* is in question, opacity stops being only a rule to enforce and
may instead become a property of the artifact, and two things runtime-opacity took
as fixed reopen: what the runtime is made of, and how the suite tests it (today by
*sourcing* the library and *grepping* its source text). That is a larger,
core-touching decision than a v0.7.0 delta, so it lives in this named grouping
rather than as a v0.7.0 scope.

**This grouping does not retire, delete, or modify `runtime-opacity`.** That scope
stays in place and active under [../v0.7.0/](../v0.7.0/), governing the
readable-shell status quo; this grouping only reframes opacity as one facet of the
runtime's *form*. If a later form decision makes the engine opaque by construction,
the disposition of runtime-opacity is revisited then; until then the two coexist.

## Scopes

- [runtime-form/](./runtime-form/) — the substance and shape of the runtime:
  shell library vs compiled binary, the black-box action contract that lets the
  form change without rewriting every test, opacity as a property of the form, and
  the distribution that follows. Start at
  [runtime-form/design.md](./runtime-form/design.md); its testing concern is
  [runtime-form/test-strategy.md](./runtime-form/test-strategy.md).

## Layout convention

`sources/system-design/<product>/<grouping>/<scope>/`. Every folder has a
`README.md` landing/index; each scope's normative design is `design.md`, with
detail split into files or sub-folders as it grows. This grouping uses a named slug
(`runtime-rearchitecture`) rather than a semver folder because the change spans
more than one release.

## Authority

This grouping asserts **no new invariant and no new role or skill**. The runtime's
canonical rule owners stay under `wrapper/` (`wrapper/contracts/invariants.yaml`
for INV-RUNTIME-01 and the path-lease / base-selection / repository-grounding
owners; `wrapper/adapters/AGENTS.md` for the invoke-not-read prose). A change of
*form* changes the artifact those owners point at and the corollary wording around
reading it — never the ownership map or the meaning of any invariant.
