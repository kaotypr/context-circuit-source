# Proposal 0017 — add domain: repository-grounding

- id: 0017-domain-repository-grounding
- target_context_unit: context/domains/repository-grounding/README.md
- operation: add
- confidence: high
- status: review-needed
- related_plan: —
- affected_commits: [e3e95cd]
- evidence_refs:
  - sources/system-design/context-circuit/v0.6/repository-grounding/
  - wrapper/contracts/invariants.yaml (INV-GROUND-01, INV-GROUND-02, INV-GROUND-03)
  - wrapper/contracts/schemas/grounding-manifest.yaml
  - wrapper/contracts/schemas/execution.yaml (grounding manifest recorded per execution)
  - wrapper/contracts/schemas/worker-handoff.yaml (repository_friction)
  - wrapper/adapters/writer-brief.md (shipped brief template, promoted to workspace root)
  - wrapper/runtime/engine.sh (cc_discover_repo_grounding, cc_harden_worktree, cc_grounding_directive, cc_writer_brief_assemble, cc_brief_preflight)
  - agents/writer.md, agents/coordinator.md
  - test/grounding/test-grounding.sh
  - template-harness/scenarios/11-repo-grounding/case.yaml

## Statement

v0.6 ships **repository grounding**: before implementing, the writer reads and
honors the *target repository's own* agent guidance (AGENTS.md, CLAUDE.md,
`.cursor/rules`, `.github/copilot-instructions.md`, `.agents/skills/*/SKILL.md`),
discovered live from the execution worktree by a deterministic runtime scan and
delivered through a generated brief — never hand-authored into each writer prompt.
This is a bounded concern around the writer's grounding and its prepared worktree,
warranting its own retrieval unit.

## Proposed unit (summary to author on acceptance)

- **Discovery is deterministic; it emits data, not prompt text (INV-GROUND-01,
  INV-RUNTIME-01).** `cc_discover_repo_grounding` scans the worktree and records a
  grounding manifest (`files`, `skills` with descriptions, `environment`) as
  execution evidence under `.runtime/executions/<plan>/<exec>/grounding/<repo>.yaml`.
  Scanning the *execution base* means bootstrapping propagates through a stack (a
  scaffold plan's AGENTS.md is present in a stacked dependent's worktree).
- **Worktree hardening.** `cc_harden_worktree` detects the toolchain and reports
  `environment: ready | no-toolchain`, so CC-induced execution-environment
  workarounds are eliminated rather than documented (full dependency provisioning
  is a later phase).
- **The writer brief is delivered, not authored (INV-GROUND-03).** A fixed shipped
  template (`writer-brief.md`, promoted to the workspace root) is filled by
  deterministic slot substitution (`cc_writer_brief_assemble`) from the manifest
  and the plan; the coordinator adds only a one-line task focus; a preflight
  (`cc_brief_preflight`) refuses a brief missing the required grounding section.
- **Precedence (INV-GROUND-02).** The brief's scope and safety rules win on
  *what/where*; the repository's guidance is authoritative on *how* within that
  scope; on conflict the writer stops and reports — repository guidance never
  overrides a CC safety or scope rule.
- **Reference, not capture.** No per-repo profile and no `plan.yaml` field. Writer
  friction its guidance did not cover returns as `repository_friction` in the
  handoff and is reconciled into a proposal on the repository's *own* agent docs.

## Verification

`sh test/acceptance.sh` (repository-grounding suite); a section in
`template-harness/test-template-runtime.sh` proves the shipped template carries it;
live case `11-repo-grounding` (grade.sh PASS incl. `file_grounded` +
`grounding_manifest_recorded`, and human-simulator pass).
