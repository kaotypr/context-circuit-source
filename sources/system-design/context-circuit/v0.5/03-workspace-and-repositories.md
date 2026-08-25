# Workspace and repositories

This document defines how a workspace represents one or more connected
repositories and how execution creates isolated branches and worktrees.

## 1. Workspace identity

workspace.yaml is the portable identity authority.

The example below assumes that the workspace root is itself a Git repository.
When a workspace is not versioned at its root, omit the reserved `workspace`
repository entry and its local binding; connected project repositories keep the
same identity and execution rules.

~~~yaml
schema_version: 1
workspace: acme-commerce
title: Acme Commerce platform
purpose: Customer commerce services and interfaces
repositories:
  - id: workspace
    purpose: Context Circuit workspace files, plans, and Product Knowledge
    default_branch: main # portable clone/setup hint; not the execution branch
  - id: api
    purpose: Public API and persistence
    canonical_url: https://github.com/acme/commerce-api.git
    default_branch: main # portable clone/setup hint; not the execution branch
  - id: web
    purpose: Customer web application
    canonical_url: https://github.com/acme/commerce-web.git
    default_branch: main # portable clone/setup hint; not the execution branch
  - id: contracts
    purpose: Shared API contracts
    canonical_url: https://github.com/acme/commerce-contracts.git
    default_branch: main # portable clone/setup hint; not the execution branch
~~~

Portable identity must not include local machine paths, access tokens, private
keys, provider payloads, or credentials.

## 2. Local bindings

repositories.local.yaml is host-local and ignored from portable artifacts. It
stores the user's actual local repository setup, including the branch they
have selected as the anchor for work.

~~~yaml
bindings:
  workspace:
    path: .
    anchor_branch: kao/development/v0.5
  api:
    path: ../commerce-api
    anchor_branch: development
  web:
    path: ../commerce-web
    anchor_branch: development
  contracts:
    path: ./repositories/contracts
    anchor_branch: development
~~~

Binding resolution is explicit. Reject:

- missing paths;
- non-Git paths;
- traversal outside permitted roots;
- unsafe symlinks;
- remote identity conflicts;
- duplicate or ambiguous bindings;
- dirty anchor checkouts when execution requires a clean anchor.

Never search the machine for a repository with a matching name.

The reserved logical repository ID `workspace` represents the workspace root
when that root is itself a Git repository. Its local binding points to `.` and
has its own user-selected `anchor_branch`, independent of the anchor branches
of project repositories. Workspace tasks may map changes to `workspace`, and
then execution creates a worktree for the workspace repository like any other
affected repository. The workspace repository is never placed under its own
`repositories/` directory.

## 3. Connected repository states

A repository binding can be:

- registered: logical identity exists but no local binding is available;
- bound: local path resolves and identity matches;
- unavailable: named path is missing or inaccessible;
- mismatched: Git identity conflicts with the workspace record;
- dirty: uncommitted changes exist in the configured anchor checkout;
- ready: bound, anchor branch configured and validated, clean, and suitable for
  worktree creation;
- active: assigned to an execution worktree.

The agent explains the state in project terms rather than exposing raw runtime
record names.

## 4. Registering a repository

Registration has two parts:

1. Record or update the portable logical identity.
2. Add a host-local binding to an existing checkout or an explicitly requested
   bootstrap destination.

Registration does not clone automatically. If the user asks to clone, report
source URL, destination, branch, and external Git effect before the separate
bootstrap action. `default_branch` may guide that clone or the first local
setup, but it does not become the execution base until the user sets or accepts
the local `anchor_branch`.

A plan may reference only registered logical repository IDs.

## 4.1 Clone or initialize a repository

Repository registration can bind an existing checkout, clone a remote
repository, or initialize a new local repository. Registration alone does not
perform any of these external Git actions.

The default destination for a new project-repository checkout is
`repositories/<repository-id>/` under the workspace root. The workspace
template creates or documents `repositories/` as a Git-ignored directory so
connected project source does not become part of the workspace repository.
Users may bind an existing checkout elsewhere when they explicitly provide its
path.

For a remote repository, an explicit clone setup:

1. resolves the logical repository ID and canonical URL;
2. proposes `repositories/<repository-id>/` as the destination;
3. clones the repository and checks out the requested branch;
4. records the local path and user-selected `anchor_branch` in
   `repositories.local.yaml`.

`default_branch` may suggest the first branch for this setup, but the user must
set or accept the actual `anchor_branch`. A clone does not automatically start
plan execution or create a worktree.

For a new local or empty repository, an explicit initialization setup:

1. creates `repositories/<repository-id>/`;
2. runs `git init` there, using the selected anchor branch as the initial branch
   when the Git host supports it;
3. records the new repository identity and local `anchor_branch`;
4. creates an initial anchor commit before the repository becomes executable.

The initial commit may be an explicitly requested empty bootstrap commit. A
repository with no commit cannot provide a worktree base and remains
registered but not ready for plan execution. Creating a remote, pushing the
initial commit, or publishing the repository is a separate delivery action.

## 5. Plan repository mapping

Every task names the repository or repositories it can change.

~~~yaml
tasks:
  - id: CONTRACT-001
    repositories: [contracts]
    paths: [openapi, generated]
    changes:
      - Add the subscription event schema.
  - id: API-001
    repositories: [api]
    paths: [src/subscriptions, test/subscriptions]
    depends_on: [CONTRACT-001]
  - id: WEB-001
    repositories: [web]
    paths: [app/subscriptions, test/subscriptions]
    depends_on: [CONTRACT-001, API-001]
~~~

A task that reads another repository but does not modify it declares a read
dependency. A task that may modify it lists it in repositories and paths.

The worker must not infer that all connected repositories are in scope.

## 6. Default and anchor branches

`default_branch` is an optional portable repository-identity hint. It describes
the branch a new team member may expect when cloning or setting up a connected
repository. It is not the user's current branch and is never, by itself, the
execution base.

`anchor_branch` is a required user-local setting for an executable repository
binding. It identifies the user's actual active branch in that local checkout.
Different users may set different anchor branches for the same logical
repository. A user who clones an active Cc workspace sets `anchor_branch` for
each local repository, usually using the current branch or the
`default_branch` as the starting suggestion. Typical values are a team branch
such as `development` or a personal development branch such as
`kao/development/v0.5`.

At execution start, the runtime validates the anchor branch, captures its tip
as the execution base commit, and records the repository ID, local binding,
`anchor_branch`, base commit, base cleanliness, and execution IDs. The runtime
creates each execution branch and worktree from that captured anchor tip; it
does not move, modify, or silently rebase the user's anchor branch.

If repository anchor branches differ in revision or name, those anchors are
recorded independently in the execution snapshot. A repair uses the existing
worktree and does not silently rebase onto a moving anchor branch.

## 7. Branch and worktree naming

The workspace runtime derives deterministic names from plan and repository IDs. Names must
be safe for Git, unique for simultaneous plans, stable across resume, independent
of machine paths, and easy to map back to plan and repository.

A conceptual layout is:

~~~text
.runtime/worktrees/<plan-id>/<repository-id>/
branch: cc/<plan-id>/<repository-id>
~~~

The physical location may differ when required by the host, but the logical
mapping remains in execution records.

## 8. Worktree lifecycle

At execution start:

1. Validate all repository bindings and each configured anchor branch.
2. Reject dirty anchor checkouts unless the plan explicitly authorizes a safe
   alternative.
3. Create one branch and worktree per affected repository.
4. Record each worktree root and initial revision.
5. Give the worker only the assigned worktree roots.

During repair, reuse the execution worktrees. Do not create an unrelated
worktree that makes the latest commit ambiguous.

After verification, preserve worktrees until a human or separate cleanup action
decides what to do.

## 8.1 Anchor branch as the pull-request target

The runtime records the selected `anchor_branch` and its captured base commit
in the execution record. The execution branch, such as
`cc/<plan-id>/<repository-id>`, is the source branch for later delivery. If the
human explicitly asks to open a pull request for the implemented plan, the
default target branch for each affected repository is that repository's
recorded `anchor_branch`.

The workspace does not open pull requests automatically after verification or
completion. A different target branch requires a separate explicit human
request; the agent must not infer one from `default_branch` or from a moving
remote branch. A multi-repository plan may produce one pull request per
affected repository.

Opening a pull request is a delivery action handled by the host or delivery
integration, not by the runtime engine. Before opening one, the action must
confirm that the execution branch and its commits are available to the
configured repository remote or provider, and that the requested target branch
exists. The action must not discover an unrelated remote, silently push a
branch, or silently substitute `default_branch`. If the source branch is not
published, the remote/provider is unavailable, or the recorded anchor branch
is missing or renamed, report the delivery as blocked and ask for an explicit
human decision. The captured anchor commit remains execution evidence; the
pull-request target is the recorded anchor branch name unless the human names
a different target.

## 9. Commit policy

The worker commits after implementation and before verifier start. The commit is
the revision the verifier must inspect; it is not a verification claim.

Each repair that changes a repository creates a new commit.

~~~yaml
commits:
  - attempt: 1
    repository: api
    revision: sha256:...
    git_commit: abc123
    kind: implementation
  - attempt: 2
    repository: api
    revision: sha256:...
    git_commit: def456
    kind: repair
~~~

The workspace runtime does not squash, amend, merge, or push as part of execution.

## 10. Multi-repository coordination

One worker may work across repositories in dependency order.

~~~yaml
repositories:
  api:
    worktree: .runtime/worktrees/0001-billing-v2/api
    branch: cc/0001-billing-v2/api
    anchor_branch: development
    base_commit: abc123
    allowed_paths: [src/billing, test/billing]
  web:
    worktree: .runtime/worktrees/0001-billing-v2/web
    branch: cc/0001-billing-v2/web
    anchor_branch: development
    base_commit: xyz789
    allowed_paths: [app/billing, test/billing]
~~~

The worker may run cross-repository tests when the plan declares them. It must
not make an undeclared change in a dependency repository just because a test
needs one.

If two independent plans affect the same repository, they use separate
worktrees. Integration or merge conflict is a delivery concern, not something
the execution worker silently resolves.

## 11. Repository evidence for context

Repository inspection can produce context proposals, but repository code is not
automatically copied into Product Knowledge. The agent extracts stable facts,
cites the repository revision and path, and proposes a focused context update.

Generated files, temporary output, credentials, and large code dumps do not
become Product Knowledge.

## 12. Repository safety failures

For a binding or worktree failure, report:

- logical repository ID;
- expected and observed identity;
- intended path or branch;
- whether product work has started;
- smallest human action needed.

Do not fall back to a similarly named repository or direct edits in the base
checkout.
