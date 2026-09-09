# What ships, and what stays here

## Mechanism

Three identities stay distinct (source-release-and-upgrade): this maintainer
checkout, the distributable template, and an instantiated workspace.

After this change:

- **Maintainer checkout** holds runtime source and the build that produces
  the binary. Changing the engine is source work here.
- **Assembled template / instantiated workspace** hold the compiled binary
  under `.context-circuit/wrapper/runtime/`. They do not hold a readable
  runtime implementation.

Release assembly already decides ships-vs-never-ships. Runtime source joins
the never-ships set (alongside maintainer plans, source-only tests, and
`sources/`). The binary joins the shipped product home.

An upgrade may replace the template-owned runtime. It must preserve
workspace-owned files. If the invoke path or the form of a record changes
meaning, the upgrade reports migration-needed and keeps the old state.

## Interfaces

- Nested runtime home stays `.context-circuit/wrapper/runtime/`.
- Template seed `template/` receives the binary, not runtime source.
- Agent-harness uses the same assembled template artifact it uses today, so
  cases see the binary the user would see.

## Edge cases

- **Source checkout vs product workspace.** Agents working *on* Context
  Circuit in this checkout still need source in order to change the engine.
  Opacity is a product-workspace property, not a claim that this repository
  has no runtime source.
- **A leftover `engine.sh` in the product is a bug.** A tiny shell stub that
  execs the binary is still a readable implementation and would keep the
  "don't read it" rule alive.
- **Platforms.** The product is POSIX today (macOS and Linux). The binary
  must run on those hosts. Windows is out of scope unless a later intent
  adds it.
- **Offline.** An instantiated workspace must be able to invoke the binary
  without a network fetch. The binary is present in the tree the template
  already ships, not downloaded at first use.
- **Not cryptographic.** A determined reader can still `strings` a binary.
  Success is: no useful source in the product, and no instruction tax to
  police reading it.
