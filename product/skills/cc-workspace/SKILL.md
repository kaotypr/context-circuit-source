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
repository evidence; ask only when the choice is ambiguous. The root workspace
may itself be connected as `.`.

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
