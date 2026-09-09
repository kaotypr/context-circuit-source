# Ship the runtime as a compiled binary

## Capability

Agents in a Context Circuit product workspace invoke a **compiled binary** for
every runtime action they already call today. The workspace they are working
in does not contain a readable runtime implementation, so there is nothing to
open instead of a skill or brief, and product instructions no longer repeat
"must not read `engine.sh`".

The runtime's job is unchanged: a small host-neutral deterministic library
(INV-RUNTIME-01) that writes atomic records (INV-RUNTIME-02). This intent
changes the **form** of that library in the product, not what it is allowed to
decide.

Runtime source remains in this maintainer checkout. Changing the engine is
still ordinary source work here.

## Problem

The product tells every role the same thing: invoke
`.context-circuit/wrapper/runtime/engine.sh`, and **do not read it**. That
prohibition is restated in the coordinator, in adapters, in skills, and in
harness forbidden lists. It exists because the engine is a readable bash
script sitting in the tree the agent can open.

A warning is not opacity. An agent that reads the script can treat it as a
substitute for a skill, a brief, or a gate. The repeated rule is load, drift,
and still fail-open if a role ignores it.

Compiling the current script in place, or shipping a binary **next to**
source, would not solve that. The readable implementation would still be in
the product, and the prohibition would just move.

## Principles

1. **Opacity is absence, not a warning.** The product an agent works in has
   no useful runtime source to read. See [shipping.md](shipping.md).
2. **Invoke, don't interpret.** Skills and briefs still name actions; the
   agent calls the binary and reads printed results. See [invoke.md](invoke.md).
3. **Same library, new form.** Model-blind, credential-free, no provider or
   routing inside the runtime. Gates and host-blocked stay.
4. **Source has a home.** People who change the engine read source in this
   checkout. Product workspaces do not.
5. **One owner per rule.** Drop the scattered "must not read" copies; do not
   retarget them at a new filename.

## Fixed decisions

- The thing agents invoke in a product workspace is a compiled binary under
  `.context-circuit/wrapper/runtime/`.
- The shipped template and an instantiated workspace contain that binary, not
  `engine.sh` and not compiled-language source beside it.
- Runtime source stays in the maintainer source checkout and is not a
  ships-with-template file.
- This is a rewrite into a compiled language, not a packer around the current
  bash script.
- CLI verbs (the actions skills already name) and printed results remain the
  agent-facing contract.
- INV-RUNTIME-01/02, Gate 1, Gate 2, planner, verifier floor, and
  host-blocked are unchanged.
- Assurance is **Critical**: the library that enforces every gate is rewritten
  and then shipped. A human may lower the tier.
- Language is **not** fixed here; it is the open question on the intent.
  Recommendation: Go.

## Shape of the whole

```mermaid
flowchart LR
  subgraph product ["Product workspace"]
    Skill["Skill / brief names an action"]
    Bin["Compiled runtime binary"]
    Skill --> Bin
  end
  subgraph source ["Maintainer checkout"]
    Src["Runtime source"]
    Build["Release assembly builds the binary"]
    Src --> Build
  end
  Build --> Bin
```

| Topic | What it settles | Detail |
| --- | --- | --- |
| What ships | Binary in the product; source only here; upgrade | [shipping.md](shipping.md) |
| How agents call it | Invoke path, drop the read prohibition | [invoke.md](invoke.md) |
