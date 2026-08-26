# Repository grounding — schema and contracts

Each change names its canonical owner; this scope specifies intent, the owner
files carry the rule. No rule is duplicated (v0.5 principle 5.7). **No `plan.yaml`
change and no per-repo profile** — grounding is discovered, not captured.

| Owner | Change |
| --- | --- |
| `wrapper/runtime/engine.sh` | add `cc_discover_repo_grounding <worktree>` (deterministic scan → manifest); worktree-hardening functions (toolchain detect, dependency provisioning, hook handling); a brief preflight check for the required grounding slot |
| `wrapper/contracts/schemas/grounding-manifest.yaml` **(new)** | the manifest shape: `files[]`, `skills[]{name, description}`, `environment` |
| `wrapper/contracts/schemas/execution.yaml` | record the grounding manifest snapshot used for the execution (evidence + reproducibility) |
| `wrapper/adapters/` (writer-brief template) **(new/updated)** | the writer-brief template with the **required** repository-grounding slot and directive, filled by deterministic slot substitution from the manifest + `plan.yaml` |
| `wrapper/contracts/invariants.yaml` | add INV-GROUND-01 (writer must read and honor discovered repo guidance; discovery is deterministic and live from the worktree), INV-GROUND-02 (precedence: CC safety/scope wins on what/where; repo guidance authoritative on how within scope), INV-GROUND-03 (the brief's fact sections are template + runtime data, the grounding slot is required, a preflight enforces it); extend the worker-behavior rules with the read-and-honor obligation |
| `wrapper/contracts/schemas/worker-handoff.yaml` | add a `repository_friction` field routing "the repo's guidance did not cover X" to a proposal on the repo's own agent docs |
| `agents/coordinator.md` | the coordinator **delivers** the brief and adds only the task focus; it does not author the fact sections |
| `wrapper/manifest.yaml` | no schema-version change owned by this scope; `runtime_version 0.6.0` is set by the run-stack scope in the same release |

## Invariants (intent)

- **INV-GROUND-01** — Before implementing, the writer reads and honors the target
  repository's own agent guidance, discovered live from the worktree by a
  deterministic runtime scan. The runtime emits data (a manifest), not prompt text.
- **INV-GROUND-02** — Precedence: the brief's scope and safety rules are
  authoritative on what and where; the repository's guidance is authoritative on
  how to write code within that scope; it never overrides a CC safety or scope rule.
- **INV-GROUND-03** — The brief's fact-bearing sections are a fixed template
  filled from the manifest and the plan. The grounding slot is required; a
  preflight refuses to deliver a brief that omits it.

## The feedback loop

When a writer hits repository friction its guidance did not cover, the handoff
carries a `repository_friction` note. Reconciliation routes it to a **proposal to
improve that repository's own agent docs** (a normal plan targeting its
`AGENTS.md` / skills) — never a CC-side profile. The failure becomes a durable
improvement to the repo, and the next writer reads it via discovery.

## What this scope does NOT change

- No `plan.yaml` field and no per-repo captured profile.
- No change to worker/verifier roles, the repair loop, the completion gate, or
  execution authority — this scope grounds the writer and prepares its worktree.
- No weakening of worktree isolation or the one-writer lock.
