---
name: cc-workspace
description: Initialize a Context Circuit workspace, connect or clone its repositories, and manage members and their allocation bands.
---

# Set up a workspace and its repositories

Everything here happens on a request. None of it is inferred from a project that
merely looks uninitialized, and none of it runs a second time on a machine that
already has the shared files.

## Initialize

Initialize a workspace only when requested. Gather its name, purpose, and first
member's name; derive a readable member ID.

```sh
context-circuit-cli --workspace <root> init --name NAME --purpose TEXT \
  --member ID --member-name NAME
```

The executable creates embedded instructions, folders, and initial records, and
preserves existing files. A solo workspace still has a real member.

On another machine, select an existing member and connect its local checkouts;
do not reinitialize the shared workspace. Cloning a workspace onto a second
machine is member selection plus local bindings, never a second initialization:

```sh
context-circuit-cli --workspace <root> member use --id ID
```

## Connect repositories

Connect existing repositories, clone repositories, or initialize new ones when
requested. Resolve the correct destination and base branch from the request and
repository evidence; ask only when the choice is ambiguous.

```sh
context-circuit-cli --workspace <root> repo connect --id ID --path PATH --base BRANCH
context-circuit-cli --workspace <root> repo clone --id ID --path NEW_PATH --base BRANCH --url URL
context-circuit-cli --workspace <root> repo init --id ID --path NEW_PATH --base BRANCH
```

Logical IDs, repository URLs, and default branches are shared. Checkout paths and
base branches are this machine's, so worktrees branch from whatever this machine
bases work on and a second machine may base the same repository elsewhere. Record
meaningful relationships, and obey repository instructions within each checkout.

A repository can be registered before its first commit; worktree preparation
needs one. Setting a base records the intended branch without creating or
resetting it.

## Describe the workspace's own repository

The shared records travel through the Git repository carrying the workspace, so
a member reads plans from whatever commit that checkout sits on. Describe it
when the workspace is under version control, on the first machine and on every
machine that clones it:

```sh
context-circuit-cli --workspace <root> workspace connect --base BRANCH
context-circuit-cli --workspace <root> workspace base --branch BRANCH
context-circuit-cli --workspace <root> workspace remote --url URL --default-branch BRANCH
```

It splits like a repository: the URL and default branch are shared, while the
checkout root and this machine's base branch are local. `--path` defaults to the
workspace root and is needed only when the workspace sits inside a larger
repository. Once described, `status` reports the workspace's own branch and
uncommitted records.

This is not an entry in `repositories`, and `repo connect` refuses the
workspace's own checkout: a plan names repositories, a relationship joins two,
and a worktree is cut from one, none of which the workspace is. Members share
one branch of it; a machine that works from another records that in its local
binding alone.

## Members and allocation bands

A member may hold an allocation band: a numeric block that member allocates from
alone. Bands are what let two clones that cannot see each other allocate without
colliding, so give every member of a team workspace a distinct band before they
work apart. A solo workspace needs none, and an unbanded member allocates from
the numbers no band has claimed.

```sh
context-circuit-cli --workspace <root> member add --id ID --name NAME --band 100
context-circuit-cli --workspace <root> member band --id ID --band 200
```

The band decides which number comes next and never makes a reserved number
reusable. A band is not a namespace: the ID stays global.

Bands prevent collisions only between members who actually hold distinct ones.
Unbanded members working in separate clones, and any clone whose roster is
stale, can still allocate the same number. Synchronize the shared workspace
before allocating, and resolve conflicting allocations and references before
publishing their IDs. Do not claim distributed collision prevention beyond what
bands give.

`.context-circuit/docs/workspace.md` describes which files are shared, which
belong to one machine, and the guarantees editing them carries.
