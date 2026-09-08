---
name: cc-workspace
description: Initialize or orient a Context Circuit workspace and register, connect, clone, or initialize project repositories.
---

Use for read-only orientation and for workspace/repository setup. This skill
never creates plans or executes work.

## Orientation (read-only)

Answer from `workspace.yaml` identity and the `context/INDEX.md` retrieval
catalog. Report the workspace name and purpose, registered repository keys and
whether their local bindings are available, important Product Knowledge areas,
the active plans in `plans/INDEX.md`, and current execution state only when
asked. Orient **member identity** the same way as repository bindings: report
whether `member.local.yaml` names a roster member, and if it is missing list
the members in committed `members.yaml`. State uncertainty instead of inventing
project facts. Do not read
`plans/archive/`. Recommending a next action is not authorization to take it.

## Initialize

When asked to initialize or set up a workspace, run the runtime
`workspace-init` to create the minimal deterministic structure, then record
project identity and register only the repositories the user names. Ask only for
facts that cannot be safely inferred. Do not invent Product Knowledge, create
plans, or clone repositories the user did not request. Do not create,
initialize, or register a repository the user has not explicitly named or
requested: when the code has no home yet, orient and ask whether to create a new
repository or connect an existing one rather than choosing a path and creating
it. Recording the project's identity from the stated goal is fine; materializing
a repository is a separate, explicitly requested action (doc 02 §2).

## Register / connect a repository

Registration has two parts: record the portable logical identity in
`workspace.yaml` (id, optional credential-free `canonical_url`, optional
`default_branch`), and add a host-local binding in `repositories.local.yaml`
with the local `path` and the user-selected `base_branch`. Connecting binds an
existing checkout; it does not clone or initialize. Validate the binding with the
runtime `repository-resolve`.

When you confirm a connection to the user, describe it in plain language by its
effect — "I've connected your <name> project; I'll work from the <branch> branch"
— and never name `workspace.yaml`/`repositories.local.yaml` or say "base
branch"/"binding".

The workspace root, when it is a Git repository, binds as the reserved id
`workspace` at path `.` with its own base branch. New project repositories
default to the git-ignored `repositories/<repository-id>/` path.

## Clone or initialize

Clone and `git init` are explicit, separately requested external Git actions.
For a clone, report source URL, destination (`repositories/<id>/` by default),
and branch before acting, then record `path` and `base_branch`. For an empty
repository, `git init` there, set the base branch, and create the initial
base commit required before execution can build a worktree. `default_branch`
may suggest the first branch but never becomes the base branch automatically.

Creating a remote, pushing, or publishing is a separate delivery action.

## Member identity (host-local)

Each machine names **one** committed roster member in gitignored
`member.local.yaml`, the same once-per-machine pattern as
`repositories.local.yaml`. Gitignore is not a read block — Read the file from
the workspace root.

When identity is missing, list the members already in `members.yaml` and write
`member.local.yaml` after a one-time human choice of an **existing** roster
member. If the roster is empty, bootstrap the first member using the default
band size (99 intents, 999 plans) with a name the human chooses; do not invent
overlapping bands. Never ask for numeric ranges or a block number in ordinary
chat — bands live in the roster, not in the conversation.

When a member's band is exhausted (`INTENT_ID_EXHAUSTED` / `PLAN_ID_EXHAUSTED`),
extend the committed roster with a new non-overlapping band for that member
(or add another member). Never wrap, recycle, or ask for a block number in
chat. Default new-member size remains 99 intents and 999 plans (INV-MEMBER-01).

Do not write `member.local.yaml` during `workspace-init`. Runtime verbs
`member-roster-validate`, `member-identity-read`, and `member-band-resolve`
fail closed (`MEMBER_IDENTITY_MISSING`) rather than inventing a band.
