# Repository grounding — worked briefs and acceptance

All three briefs come from the **same template**; only the discovery manifest and
the environment status differ. The coordinator authors only the task-focus line.

## Rich repo — guidance and skills present

```md
## Environment (prepared for you)
Ready: dependencies provisioned, commit hooks handled. Do NOT install or modify
dependencies. Use the repository's own build/test/lint commands.

## Repository grounding — understand how this repo expects agents to work
Before you write any code, read and apply this repository's own agent guidance.
It is authoritative on HOW to write code here, within the scope this brief sets.
- AGENTS.md, CLAUDE.md — read in full.
- Skills — all of this repo's skills are listed below; read only the ones
  RELEVANT TO YOUR TASK, chosen by description:
    - gezit-ui — DLS components via the private shadcn registry
    - gezit-illustration — brand illustration package usage
- .cursor/rules/ — honor any rule matching the files you change.
If anything conflicts with this brief's scope or safety, STOP and report.
```

## Doc-less repo — code exists, no agent guidance

```md
## Environment (prepared for you)
Ready: dependencies provisioned, commit hooks handled. Do NOT install or modify
dependencies. Use the repository's own build/test/lint commands.

## Repository grounding — understand how this repo expects agents to work
No repository agent guidance was discovered (no AGENTS.md, CLAUDE.md, skills, or
rules). Ground your work in the plan and Product Knowledge below, and follow the
patterns already present in the code.
```

## Empty / greenfield repo — nothing yet

```md
## Environment (prepared for you)
No dependency toolchain detected (no lockfile/manifest). If this plan scaffolds
the toolchain, create it within the allowed paths; add no dependencies beyond
what the plan specifies.

## Repository grounding — understand how this repo expects agents to work
No repository agent guidance was discovered. Ground your work in the plan and
Product Knowledge (and, for a greenfield build, the accepted system design). A
scaffold plan may author this repo's AGENTS.md/skills as one of its tasks; later
writers will then discover and honor them.
```

## The feedback loop in action

A writer in a doc-less repo cannot tell how the project runs its tests. It
completes what it can and returns, in its handoff:

```yaml
repository_friction:
  - "No AGENTS.md; test command not discoverable — inferred `npm test` from
     package.json scripts. Recommend the repo document its agent build/test flow."
```

Reconciliation turns that into a **proposal to add an `AGENTS.md` to that repo**
(a normal plan). Once delivered, the next writer discovers it — the gap closes at
the right owner.

## Acceptance criteria

The repository-grounding scope is acceptable when all v0.5 acceptance criteria
still hold and:

1. Worktree preparation runs a deterministic discovery scan and produces a
   grounding manifest recorded in the execution evidence.
2. The manifest lists every present agent file and every skill with its
   description; a doc-less repo yields an empty manifest.
3. The writer brief always contains the repository-grounding section, filled from
   the manifest; a brief missing it is refused by preflight before delivery.
4. The section instructs the writer to read AGENTS.md/CLAUDE.md fully and to
   select and follow only task-relevant skills.
5. The coordinator authors only the task-focus line of the brief; the fact
   sections are template + runtime data, not free-composed.
6. Precedence holds: a repo instruction never overrides a CC safety or scope rule;
   the writer stops and reports on conflict.
7. A hardened worktree lets the writer run the repo's normal commands with no
   symlink, dep-verify, or hook workarounds.
8. An empty/greenfield repo produces an honest brief (no toolchain, no guidance)
   and does not block; a scaffold plan may author the repo's agent docs.
9. Writer friction is returned as a handoff note and reconciled into a proposal on
   the repo's own agent docs, never a CC-side profile.
10. No `plan.yaml` field and no per-repo captured profile are introduced.
11. Discovery, hardening, and the preflight are deterministic runtime functions
    with no model prompts in the runtime.
12. The semantic acceptance suite covers: discovery + manifest, the three brief
    variants, the preflight refusal, precedence conflict, and the friction→proposal loop.
