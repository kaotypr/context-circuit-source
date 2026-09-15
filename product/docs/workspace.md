# Workspace files

The executable carries the blank workspace as embedded assets. Initialization
creates the shared instruction and folders alongside these small records:

| File | Contents | Shared? |
| --- | --- | --- |
| .context-circuit/VERSION | Workspace template version | Yes |
| .context-circuit/CLI_VERSION | CLI version this workspace pins; installed side by side | Yes |
| .context-circuit/role-tiering.yaml | Per-host role model and effort preferences | Yes |
| .context-circuit/role-tiering.local.yaml | This machine's overrides of those preferences | No |
| workspace.yaml | Version, name, purpose, optional CLI release mirror, repository IDs with their URL and default branch, relationships | Yes |
| members.yaml | Member ID to display name and optional allocation band | Yes |
| .context-circuit/ids.yaml | Permanent intent and plan ID reservations | Yes |
| intent/iNNN-slug.md | Intent content, created_by, created_at, approved_at, approval note, linked plans | Yes |
| plans/pNNNN-slug.md | Plan, repositories, dependencies, created_by, created_at, completed_at, progress | Yes |
| context/ | Optional durable project notes and catalog | Yes |
| member.local.yaml | Active member ID on this machine | No |
| repositories.local.yaml | Repository ID to this machine's checkout path and base branch | No |
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
    url: https://git.example.com/acme/api.git
    default_branch: main
  web:
    url: https://git.example.com/acme/web.git
    default_branch: main
relationships:
  - from: web
    to: api
    description: Consumes the billing API
```

An organization that mirrors CLI releases into its own GitLab project records it
once as `cli_registry`, an https project URL. The value is shared, so everyone who
clones the workspace installs from that mirror without setting anything up on
their own machine; only the credential stays personal. Unset, the CLI installs
from the product's own releases.

Use a separate checkout for each execution environment (native Windows, WSL,
remote server, or container). Shared files synchronize through Git; local paths,
worktrees, installed dependencies, and uncommitted work do not migrate automatically.
A local binding holds this machine's checkout path, including `.` for the
workspace repository, and the branch that machine starts work from:

```yaml
bindings:
  api:
    path: ../api
    base_branch: release/24.4
```

The base branch is local because it differs per machine: one checkout follows a
release branch while another stays on the default. Worktrees branch from it, and
a binding that records none falls back to the repository's shared
`default_branch`. `repo base` changes this machine's base; `repo remote` changes
the shared URL and default branch.

After cloning a shared workspace onto another machine, use `member use` and
`repo connect` with the existing IDs. Do not initialize it again. `repo connect`
writes the shared record only when the ID is new, taking the URL from the
checkout's `origin` unless `--url` names one. A new repository can be registered
before its first commit; worktree preparation needs a commit. Setting a base
records the intended branch without creating or resetting it.

A record's instants — `created_at`, `approved_at`, `completed_at` — are canonical
ISO 8601 UTC timestamps, `2026-09-15T10:53:00Z`; every other date is an ISO 8601
calendar date, `YYYY-MM-DD`, in UTC. Both are written bare in YAML and unadorned
in prose, so an instant in frontmatter and the same instant heading its note read
identically.

YAML changes use goccy/go-yaml document edits to retain surrounding comments and
ordering. Presentation can normalize; byte-for-byte formatting preservation is
not promised. Unknown fields in structured control records and duplicate mapping
keys produce errors. User notes remain ordinary Markdown with no schema registry.

Atomic single-file replacements and a portable OS lock protect cooperating
commands in one directory. Multi-file operations may leave useful partial output
on interruption; IDs are reserved before file creation and are never rolled back.
Check the reported path and resume or repair the named files. Locks are released
when the process exits, including abnormal exit; the empty local lock file remains.

An optional per-member allocation band divides the numeric range so members can
allocate offline without colliding. A member holding band N takes intents from
`N*100` and plans from `N*1000`; unbanded members take the numbers no band
claims. Bands must be distinct, and a roster with two members on one band is
refused until it is resolved. Changing or clearing a band leaves existing records
and their reservations untouched.

```yaml
members:
  maya:
    name: Maya
    band: 1
  alex:
    name: Alex
    band: 2
```

Bands are the only offline collision prevention here, and they only cover members
who hold distinct ones against a current roster. Before allocating from separate
Git clones, synchronize workspace commits: a local file lock is not distributed
coordination, and a stale roster can hand two people the same band. Resolve
competing IDs and every reference before sharing either conflicting record;
preserve the reservation history.

An explicitly invoked `check` reports local bindings, broken record links,
duplicate IDs, dependency cycles, and stale worktree associations. It is diagnostic
information, not an execution gate or a substitute for code review.
