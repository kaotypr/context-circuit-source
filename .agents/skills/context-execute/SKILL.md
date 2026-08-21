---
name: context-execute
description: Execute an explicitly requested approved plan or connected plan stack through isolated writer and verifier roles.
---

Use only after an explicit named run request. Preflight status, dependencies,
archive eligibility, dirty base, wrapper compatibility, lease ownership, child
primitives, and worktree safety. Acquire the atomic plan lease, create an
exclusive worktree, write a bounded delegation packet and receipt, direct one
writer, then create a separate read-only verifier. Sequential tasks share the
writer/worktree; connected plans use frozen graph/progress records.

Do not auto-approve, auto-finish, merge, push, publish, deploy, or discard.
Runtime shapes and ownership checks belong to the wrapper schemas and engine.
