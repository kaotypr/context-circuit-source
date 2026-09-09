# What ships, and what stays here

## Mechanism

Three identities stay distinct (source-release-and-upgrade): this maintainer
checkout, the distributable template, and an instantiated workspace.

After this change:

- **Maintainer checkout** holds `engine.sh` and the `shc` build that produces
  the binary. Changing the engine is still editing that script.
- **Assembled template / instantiated workspace** hold the compiled binary
  under `.context-circuit/wrapper/runtime/`. They do not hold `engine.sh`.

Release assembly already decides ships-vs-never-ships. `engine.sh` joins the
never-ships set (alongside maintainer plans, source-only tests, and
`sources/`). The binary joins the shipped product home.

An upgrade may replace the template-owned runtime. It must preserve
workspace-owned files. If the invoke path or the form of a record changes
meaning, the upgrade reports migration-needed and keeps the old state.

`shc` turns the script into C, encrypts the text, and compiles it. At run
time the binary decrypts and hands the script to a shell. The product still
needs a working shell on the host, the same as today.

## Interfaces

- Nested runtime home stays `.context-circuit/wrapper/runtime/`.
- Template seed `template/` receives the binary, not `engine.sh`.
- Agent-harness uses the same assembled template artifact it uses today, so
  cases see the binary the user would see.

## Edge cases

- **Source checkout vs product workspace.** Agents working *on* Context
  Circuit in this checkout still need `engine.sh` in order to change the
  engine. Opacity is a product-workspace property, not a claim that this
  repository has no runtime source.
- **A leftover `engine.sh` in the product is a bug.** A tiny shell stub that
  execs the binary is still a readable implementation and would keep the
  "don't read it" rule alive.
- **Platforms.** The product is POSIX today (macOS and Linux). `shc` produces
  a host-native binary, so release assembly must compile for those hosts.
  Windows is out of scope unless a later intent adds it.
- **Offline.** An instantiated workspace must be able to invoke the binary
  without a network fetch. The binary is present in the tree the template
  already ships, not downloaded at first use.
- **Not cryptographic.** `shc` is a wrapper, not a compiler in the C/Go
  sense. A determined reader can recover the script from some `shc` binaries.
  Success is: no `engine.sh` in the product for an agent to open, and no
  instruction tax to police reading it.
