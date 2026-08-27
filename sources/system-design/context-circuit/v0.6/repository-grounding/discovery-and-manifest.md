# Repository grounding — discovery and manifest

## When and where

Discovery runs during **worktree preparation**, right after the runtime creates
the worktree from **its execution base** (v0.5 §9). It scans **that worktree** —
the exact checkout the writer will use — so the result always matches what the
writer sees, with no cache and no staleness.

The execution base is not always the anchor tip: under the run-stack scope a
dependent plan's worktree is based on a **stacked predecessor branch or an
integration merge** (see [../run-stack/execution-bases.md](../run-stack/execution-bases.md)).
Discovery scanning the base — rather than the anchor — is what makes
**bootstrapping propagate through a stack**: if a scaffold plan authors the repo's
`AGENTS.md` and a dependent plan is stacked on it, the dependent's worktree
already contains that guidance, so its writer discovers and honors it.

It is a **deterministic runtime function**, `cc_discover_repo_grounding <worktree>`
(owned by `wrapper/runtime/engine.sh`). It contains no intelligence: a filesystem
scan plus a small frontmatter read. Per INV-RUNTIME-01 it emits **data**, not
prompt text.

## What it scans

A known, host-neutral set of agent-instruction locations:

- `AGENTS.md`, `CLAUDE.md` — the repo's agent contract (generic + Claude).
- `.cursor/rules/**` — Cursor rules.
- `.github/copilot-instructions.md` — Copilot instructions.
- `.agents/skills/*/SKILL.md` — skills; for each, read the `description:` from its
  frontmatter to build a menu.
- (Extensible: the set is a runtime constant, not per-repo configuration.)

## The manifest

Discovery emits a structured **grounding manifest** — the data the brief template
consumes and the execution evidence records:

```yaml
grounding_manifest:
  files:                       # agent guidance present in the worktree
    - AGENTS.md
    - CLAUDE.md
    - .cursor/rules/
  skills:                      # name + one-line description from each SKILL.md
    - name: gezit-ui
      description: DLS components via the private shadcn registry
    - name: gezit-illustration
      description: brand illustration package usage
  environment: ready           # ready | no-toolchain   (from worktree hardening)
```

- `files` lists only what exists (an empty list is valid — a doc-less repo).
- `skills` carries the **name + description** so the writer can judge relevance
  without opening each skill (see [writer-brief-template.md](./writer-brief-template.md)).
- `environment` is set by worktree hardening (see [worktree-hardening.md](./worktree-hardening.md)).

The manifest is recorded in the execution evidence (a snapshot of what was
discovered for that execution), so a run is reproducible and a later reviewer sees
exactly what the writer was pointed at.

## Two moments, one function

The same discovery function serves two moments, both reading the real checkout:

- **Plan creation** — the coordinator consults the bound repo's guidance so plans
  are written *with* the repo's conventions in view.
- **Worktree preparation** — the manifest fills the writer brief.

## Why deterministic discovery is the fix

The original failure was the coordinator sourcing repository facts from its own
reasoning and getting them wrong. A `find` over the worktree plus a frontmatter
read is ground truth: it cannot forget a skill that exists or invent one that does
not. The AI stops guessing; the shell reports.
