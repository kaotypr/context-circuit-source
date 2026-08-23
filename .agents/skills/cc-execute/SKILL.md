---
name: cc-execute
description: Execute an explicitly requested approved plan or connected plan stack through isolated writer and verifier roles.
---

Use only after an explicit named run request. Preflight status, dependencies,
archive eligibility, dirty base, wrapper compatibility, lease ownership, child
primitives, and worktree safety. Acquire the atomic plan lease, create an
exclusive worktree, write a bounded delegation packet and receipt, direct one
writer, then create a separate read-only verifier. Sequential tasks share the
writer/worktree; connected plans use frozen graph/progress records.

If the workspace identity is `product-source` and the only dirty changes are
the exact approved-plan projection (`plan.yaml: draft → approved` and task
frontmatter `draft → ready`), classify the state as
`MAINTAINER_APPROVAL_COMMIT_REQUIRED` and show a focused commit card. Do not
discard, revert, or auto-commit it. Any other dirty path remains
`DIRTY_BASE_BLOCKED`. After the maintainer explicitly commits those exact
approval changes, rerun the named execution request.

Do not auto-approve, auto-finish, merge, push, publish, deploy, or discard.
Runtime shapes and ownership checks belong to the wrapper schemas and engine.

When the host supplies a native child, map it to the existing bounded writer or
independent verifier delegation packet and include provider-neutral
`host_evidence`. Host permission flags are observations, not authorization. If
the verifier child is unavailable, stop with the canonical host-blocked outcome
instead of changing the role or verification path. Resume uses the same
receipt, wrapper, Git, lease, and worktree checks on every host.

The root runtime path consumes the engine graph through
`cc_construct_runtime_graph`, `cc_validate_runtime_graph`, and the commit
marker. Project only `cc_runtime_launch_projection` to the child; the writer
and verifier receive the same generated delegation, child-start, receipt, and
handoff records. Do not reconstruct a record, broaden the delegation from
launch text, or treat a host permission flag as authorization.
