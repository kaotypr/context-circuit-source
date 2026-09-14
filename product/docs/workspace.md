# Workspace files

The executable carries the blank workspace as embedded assets. Initialization
creates the shared instruction and folders alongside these small records:

| File | Contents | Shared? |
| --- | --- | --- |
| .context-circuit/VERSION | Workspace template version | Yes |
| .context-circuit/CLI_VERSION | Recommended independently released CLI version | Yes |
| .context-circuit/role-tiering.yaml | Per-host role model and effort preferences | Yes |
| workspace.yaml | Version, name, purpose, repository IDs, default base branches, relationships | Yes |
| members.yaml | Member ID to display name | Yes |
| .context-circuit/ids.yaml | Permanent intent and plan ID reservations | Yes |
| intent/iNNN-slug.md | Intent content, approval note, created_by, linked plans | Yes |
| plans/pNNNN-slug.md | Plan, repositories, dependencies, created_by, progress | Yes |
| context/ | Optional durable project notes and catalog | Yes |
| member.local.yaml | Active member ID on this machine | No |
| repositories.local.yaml | Repository ID to checkout path | No |
| .context-circuit/local/worktrees.yaml | Optional plan/worktree associations | No |
| .context-circuit/local/write.lock | OS-managed edit lock | No |
| repositories/, .worktrees/ | Local Git working copies | No |

Example shared repository definition:

```yaml
version: 2
name: Acme
purpose: Billing software
repositories:
  api:
    base_branch: main
  web:
    base_branch: main
relationships:
  - from: web
    to: api
    description: Consumes the billing API
```

Use a separate checkout for each execution environment (native Windows, WSL,
remote server, or container). Shared files synchronize through Git; local paths,
worktrees, installed dependencies, and uncommitted work do not migrate automatically.
Local bindings contain paths only, including `.` for the workspace repository.
After cloning a shared workspace onto another machine, use `member use` and
`repo connect` with the existing IDs. Do not initialize it again. A new repository
can be registered before its first commit; worktree preparation needs a commit.
Setting a default base records the intended branch without creating or resetting it.

YAML changes use goccy/go-yaml document edits to retain surrounding comments and
ordering. Presentation can normalize; byte-for-byte formatting preservation is
not promised. Unknown fields in structured control records and duplicate mapping
keys produce errors. User notes remain ordinary Markdown with no schema registry.

Atomic single-file replacements and a portable OS lock protect cooperating
commands in one directory. Multi-file operations may leave useful partial output
on interruption; IDs are reserved before file creation and are never rolled back.
Check the reported path and resume or repair the named files. Locks are released
when the process exits, including abnormal exit; the empty local lock file remains.

Before allocating from separate Git clones, synchronize workspace commits. A
local file lock is not distributed coordination. Resolve competing IDs and every
reference before sharing either conflicting record; preserve the reservation
history. Do not divide numeric ranges by member.

An explicitly invoked `check` reports local bindings, broken record links,
duplicate IDs, dependency cycles, and stale worktree associations. It is diagnostic
information, not an execution gate or a substitute for code review.
