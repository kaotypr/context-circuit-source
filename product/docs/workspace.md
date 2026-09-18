# Workspace files

The executable carries the blank workspace as embedded assets. Initialization
creates the shared instruction and folders alongside these small records:

| File | Contents | Shared? |
| --- | --- | --- |
| .context-circuit/VERSION | Workspace template version | Yes |
| .context-circuit/CLI_VERSION | CLI version this workspace pins; installed side by side | Yes |
| .context-circuit/role-tiering.yaml | Per-host role model and effort preferences | Yes |
| .context-circuit/role-tiering.local.yaml | This machine's overrides of those preferences | No |
| workspace.yaml | Version, name, purpose, optional CLI release mirror, the workspace's own repository, repository IDs with their URL and default branch, relationships | Yes |
| members.yaml | Member ID to display name and optional allocation band | Yes |
| .context-circuit/ids.yaml | Permanent intent and plan ID reservations | Yes |
| intent/iNNN-slug.md | Intent content, created_by, created_at, approved_at, approval note, linked plans | Yes |
| plans/pNNNN-slug.md | Plan, repositories, dependencies, created_by, created_at, completed_at, progress | Yes |
| context/ | Optional durable project notes and catalog | Yes |
| member.local.yaml | Active member ID on this machine | No |
| repositories.local.yaml | This machine's workspace checkout, and repository ID to its checkout path and base branch | No |
| .context-circuit/local/worktrees.yaml | Optional plan/worktree associations | No |
| .context-circuit/local/write.lock | OS-managed edit lock | No |
| repositories/, .worktrees/ | Local Git working copies | No |

Example shared repository definition:

```yaml
version: 2
name: Acme
purpose: Billing software
workspace_repository:
  url: git@github.com:acme/acme-workspace.git
  default_branch: main
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
A local binding holds this machine's checkout path and the branch that machine
starts work from:

```yaml
workspace:
  path: .
  base_branch: main
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

The two are easy to confuse, and confusing them sends a pull request to the wrong
branch, so `status`, `repo inspect`, `worktree prepare`, and `worktree inspect`
all report `base_branch` beside the branch currently checked out. Read the base
from one of those rather than from `workspace.yaml`, whose `default_branch` is
the repository's default and not what this machine delivers to.

## The workspace's own repository

The shared records travel through the Git repository carrying the workspace
itself, so a member reads plans from whatever commit that checkout sits on.
`workspace connect` describes it: where it lives and which branch it defaults to
go into `workspace_repository`, and this machine's checkout root and base branch
into the `workspace` binding. `workspace base` changes this machine's base and
`workspace remote` changes the shared record, the same split `repo base` and
`repo remote` make. Once described, `status` reports the workspace's branch,
base, and uncommitted records beside every repository's, and `check` stops asking
for it.

`--path` defaults to the workspace root and needs naming only when the workspace
sits inside a larger repository, where it names the containing root. The URL is
read from that checkout's `origin` unless `--url` gives one.

It is deliberately not an entry in `repositories`. Every consumer of that map
reads an entry as somewhere work happens — a plan names repositories, a
relationship joins two, a worktree is cut from one — and a worktree of the
workspace would duplicate the records the CLI is reading. `repo connect` refuses
the workspace's own checkout and names `workspace connect` instead.

No command commits, pushes, or merges the workspace. Members share one branch of
it; a member who keeps the workspace on another branch records that in the local
binding, where it describes that machine and nobody else.

## Joining a workspace

A clone carries the shared records and none of this machine's state. `check`
names what is missing: an active member, the workspace binding, and a binding
per repository. Add the member to the roster if nobody has, select them, run
`workspace connect`, obtain each repository, and run `agent setup` to write the
role definitions, which are gitignored and so never arrive with the clone.

Working copies belong under `repositories/<id>`, which the shipped `.gitignore`
excludes. `repo clone` takes no `--url` for an ID the workspace already
describes — the shared record carries it, and nothing shared is rewritten, which
is what obtaining a described repository should mean. `repo init` is refused for
such an ID, because an ID already naming a repository is not one to start empty.

After cloning a shared workspace onto another machine, use `member use`,
`workspace connect`, and `repo clone` or `repo connect` with the existing IDs. Do not initialize it again. `repo connect`
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
