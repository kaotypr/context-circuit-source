---
name: cc-source-develop
description: Orient, plan, implement, and verify direct maintainer changes in context-circuit-source. Applies to this source checkout, not generated Context Circuit workspaces.
---

# Develop Context Circuit source

Read root `AGENTS.md` and `WORKFLOW.md` for maintainer authority and safety.
Use `context/INDEX.md` to retrieve relevant accepted knowledge. Read `sources/`
only when the request names exact files.

Handle orientation and engineering planning conversationally. Implement by
reading and editing files directly on the current branch. Do not invoke the CC
intent, planning, pairing, execution, verification, or completion lifecycle;
do not create lifecycle branches, worktrees, leases, or runtime records.
No child agent is required for source maintenance.

Product implementation lives under `.context-circuit/`; `product/` packages
shipped skills and host integrations without activating them in this checkout.
Its nested paths become workspace-root paths during release assembly. Treat
product skills and role packets as implementation material, not instructions
for conducting source development. `template/` contains mutable blank seed data.

Before changing product behavior, identify the rule owner in
`.context-circuit/wrapper/contracts/invariants.yaml` and update its semantic
fixture. Run affected checks and `sh test/acceptance.sh` after meaningful phases;
report pre-existing failures separately. For packaging-only changes, assemble
artifacts in fresh temporary directories and compare their contents.

Preserve unrelated work. Root source skills and tools never ship. Commit,
push, publication, deployment, and destructive cleanup require explicit requests
under the maintainer safety instructions.
