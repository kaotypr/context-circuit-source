# How agents call it, and how the prohibition goes away

## Mechanism

Today the adapter Runtime section is the owner of invoke-not-read: call
`sh .context-circuit/wrapper/runtime/engine.sh <action> <args>`, and do not
read that file. Skills reference that statement rather than restating a
second policy — and many still restate it anyway, including harness
forbidden paths.

After this change the owner says: invoke the compiled binary at the nested
runtime home, with the same `<action> <args>`, and consume printed results.
There is no "must not read" corollary in product instructions, because the
product has no readable implementation to forbid.

Skills, coordinator, role files, and tests retarget the invoke line. They do
not keep `sh … engine.sh`. They do not add a new prohibition aimed at
maintainer source; that source is not in the product tree.

Harness cases that forbid `engine.sh` as a stand-in for "do not load the
runtime as context" drop that path (it will not exist) rather than renaming
the forbid to a `.go` / `.rs` file that also must not be in the product.

## Interfaces

- Agent-facing contract: action names and printed results, unchanged.
- Invoke form: a binary under `.context-circuit/wrapper/runtime/` (exact
  filename is plan-level; it is not `engine.sh`).
- One owner: product adapter Runtime section. Skills name the invoke; they
  do not copy a read-ban.

## Edge cases

- **Maintainer checkout.** Instructions here may still say how to *build*
  and *change* the runtime. That is source-maintainer guidance, not a
  product-workspace "don't read the engine" rule.
- **Tests that sourced the bash library.** Semantic tests and harness
  drivers that expected to `. engine.sh` or call `sh engine.sh` must call
  the binary instead. They do not gain a back door into runtime source in
  the product artifact.
- **Host-blocked and permission prompts stay host-local.** A binary does not
  become a second router or a way around missing children.
