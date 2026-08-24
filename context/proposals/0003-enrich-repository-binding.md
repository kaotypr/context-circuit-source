---
id: 0003-enrich-repository-binding
target_context_unit: context/domains/repository-binding/README.md
operation: change
statement: >
  The repository-binding page is consistent but under-specified. Add the shipped
  facts the design chapter 03 enumerates and the wrapper implements: the
  connected-repository state set, the two-part registration model with pre-clone
  disclosure, git-init "no commit => not ready", the full workspace identity
  fields, omit-workspace-when-unversioned, the enumerated forbidden-content list,
  and the per-user anchor nuance.
evidence_refs:
  - .agents/skills/cc-workspace/SKILL.md
  - wrapper/contracts/schemas/workspace.yaml
  - wrapper/contracts/schemas/repositories-local.yaml
  - wrapper/runtime/engine.sh   # cc_repository_register, cc_repo_resolve, cc_repo_clean, cc_worktree_prepare
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Enrich repository-binding with shipped registration/binding detail

Add (all present in the shipped wrapper; design chapter 03 is the checklist):

- Connected-repository states: registered, bound, unavailable, mismatched,
  dirty, ready, active — explained in project terms, not raw record names.
- Registration is two parts: record portable logical identity + add a host-local
  binding. Registration does not clone; a clone/init reports source URL,
  destination, branch, and external Git effect before acting.
- `git init` sequence: create dir, init with anchor as initial branch, record
  identity/anchor, create an initial (optionally empty) anchor commit; a
  repository with no commit is "registered but not ready" (no worktree base).
- Workspace identity fields: `schema_version`, `workspace`, `title`, `purpose`;
  omit the reserved `workspace` entry when the root is unversioned.
- Portable identity forbidden content (enumerate): local machine paths, access
  tokens, private keys, provider payloads, credentials.
- Anchor nuance: `default_branch` is not the user's current branch; different
  users may set different anchors for the same logical repository.

Keep the worktree lifecycle cross-linked to `plan-execution`. `BINDING_MISSING`
is a shipped engine signal (keep it).

## Acceptance action

Fold these into Scope/Behavior/Data; ground each in the cited wrapper files.
