---
id: 0011-domain-source-release-and-upgrade
target_context_unit: context/domains/source-release-and-upgrade/README.md
operation: add
statement: >
  Add a domain page "source-release-and-upgrade" covering the source ->
  distributable template -> instantiated workspace boundaries: the three named
  identities, template assembly (what ships vs never ships), the mandatory
  template .gitignore, upgrade preservation with migration-needed reporting, the
  release-as-gate rule, and the template-harness isolation contract. No current
  domain owns this; design chapter 09 is the checklist and the wrapper implements
  it via the manifest and scripts.
evidence_refs:
  - wrapper/manifest.yaml
  - scripts/release-manifest.txt
  - scripts/release-artifact.sh
  - wrapper/migrations/README.md
  - WORKFLOW.md
  - .gitignore
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
related_plan:
confidence: high
status: review-needed
---

# Propose domain: source-release-and-upgrade

## Why a new domain

The source repository's release and upgrade behavior — assembling the
distributable product, preserving workspace data on upgrade — is invariant- and
manifest-backed but has no domain page. It is distinct from the runtime
`delivery` domain (which is a user's PR/merge/push for their repos).

## Proposed scope (shipped facts; design ch09 is the checklist)

- Three identities: `context-circuit-source` (this maintainer checkout),
  `context-circuit-template` (the distributable universal project workspace), and
  the instantiated user workspace.
- Assembly ships template-owned files + blank `template/` seed and NEVER ships:
  source design, maintainer plans/logs, source tests/evidence, source `.runtime/`,
  `template-harness/`, credentials, local bindings, connected repositories,
  project Product Knowledge, customer plans, archived-plan contents
  (`wrapper/manifest.yaml`, `scripts/release-manifest.txt`).
- Mandatory template `.gitignore` excludes at least `/repositories/`,
  `/repositories.local.yaml`, `/.runtime/`.
- Upgrade replaces template-owned files/runtime/contracts/role guidance but
  preserves workspace-owned files (identity, accepted Product Knowledge,
  sources/provenance, plans + status, local bindings, connected repos, runtime
  evidence, active worktrees); it reports migration-needed and preserves old
  state when a template change changes the meaning of a plan/context/runtime
  record (`wrapper/migrations/README.md`).
- Release-as-gate: a source change is not a product-template change until release
  assembly includes it and the template acceptance checks pass.
- template-harness isolation: it assembles/selects the same
  `context-circuit-template` artifact that would ship and must not import the
  source repo's Product Knowledge, plans, `.runtime/`, or implementation state.

## Acceptance action

Create `context/domains/source-release-and-upgrade/README.md` from the domain
template, add its index row to `context/domains/README.md`, and note it in the
maintainer role's related domains.
