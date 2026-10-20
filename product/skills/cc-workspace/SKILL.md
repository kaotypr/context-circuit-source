---
name: cc-workspace
description: Initialize a Context Circuit workspace, connect or clone its repositories, and manage members and their allocation bands.
---

# Set up a workspace and its repositories

Everything here happens on a request. None of it is inferred from a project that
merely looks uninitialized, and none of it runs a second time on a machine that
already has the shared files.

## Initialize

Initialize a workspace only when requested. Gather its name and purpose, and the
first member, who is the person initializing it — a person working alone is a
team of one. Suggest their name and a readable member ID from
`git config user.name` and `git config user.email`, and use them only after the
person confirms or changes them. Their band is 1 unless they pick another. Ask
what roles they contribute in, such as maintainer or docs; the question may be
skipped, and a skipped one records no roles.

```sh
context-circuit-cli --workspace <root> init --name NAME --purpose TEXT \
  --member ID --member-name NAME
context-circuit-cli --workspace <root> init --name NAME --purpose TEXT \
  --member ID --member-name NAME --member-band N --member-role maintainer
```

The CLI creates embedded instructions, folders, and initial records, and
preserves existing files. It writes the member to the roster with their band and
selects them on this machine, so nothing is left to set up afterwards. A
teammate brought in later is added with a band of their own.

On another machine, select an existing member and connect its local checkouts;
do not reinitialize the shared workspace. Cloning a workspace onto a second
machine is member selection plus local bindings, never a second initialization.
The next section is the whole of it.

## Join a workspace someone else created

A cloned workspace carries the shared records and none of this machine's state:
no active member, no bindings, and no role definitions, because those are
gitignored. `check` on a fresh clone names what is missing; work through it.

A person this machine has not identified is asked who they are before anything
changes, which is the entry instruction's gate. Show the members `members.yaml`
lists and ask whether they are one of them or someone new. An existing member is
selected with `member use`. Someone new is added with a readable ID, their name,
and a band: offer the `next_band` that `status` reports and let them confirm it
or pick another nobody holds, and ask what roles they contribute in, which they
may skip. Then select them. Take their answer as given;
nothing verifies it. Once selected, return to the request they made.

Unless the request says otherwise, put checkouts under `repositories/<id>`,
which is gitignored and is where the workspace expects its working copies.

```sh
context-circuit-cli --workspace <root> member add --id ID --name NAME --band N
context-circuit-cli --workspace <root> member use --id ID
context-circuit-cli --workspace <root> workspace connect --base BRANCH
context-circuit-cli --workspace <root> repo clone --id ID --path repositories/ID --base main
context-circuit-cli --workspace <root> knowledge clone --id ID
context-circuit-cli --workspace <root> agent setup
```

Add the member only when the roster does not already carry them; someone may
have added them already. Every member holds a band, `member add` refuses to run
without one, and the CLI refuses a band another member holds — propose the
suggested one and confirm it rather than choosing silently.

Cloning a repository the workspace already describes takes no `--url`: the
shared record carries it, and nothing shared is rewritten. A repository the
record gives no URL for has to be obtained by hand and bound with `repo
connect`.

`members.yaml` is the only shared file joining changes. Committing and pushing
it is the user's decision, as every commit is. Stop when `check` is clean and
say what changed.

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

## Mount knowledge this workspace reads but does not own

An organization that keeps a shared knowledge center — one repository many
workspaces read — mounts it here instead of copying from it. Do this on request,
naming the repository:

```sh
context-circuit-cli --workspace <root> knowledge clone --id core-service-knowledge --url GIT_URL
context-circuit-cli --workspace <root> knowledge connect --id core-service-knowledge --path ../core-service-knowledge
```

`clone` puts the checkout under `repositories/<id>`, beside the working copies
and already gitignored; `connect` binds one the user already has. Either way the shared record is written once and
the path stays local to this machine, the same split every other checkout makes.
On a second machine, clone or connect an ID the workspace already describes and
pass no URL.

The ID shares one namespace with `repositories`, because a note's anchor names a
repository by ID and an ID meaning two things makes every anchor carrying it
ambiguous. Connecting an ID the other side already holds is refused, as is
mounting a checkout where work happens.

The `index` is the file mapping its contents. Connecting detects the usual
spellings; where the layout is unfamiliar, read the repository and name the file
with `knowledge remote --id ID --index PATH`. Without one, nothing here can
retrieve from it, and `check` says so.

Keep it current with `knowledge sync`, which fetches and fast-forwards only when
the checkout is clean and on the shared branch. A dirty or diverged checkout is
reported and left untouched: say what it holds and where the work should go,
never resolve it by resetting or discarding. It is read-only here — the push URL
is disabled — so an improvement is a merge request in that repository, from a
separate checkout of it. `.agents/skills/cc-knowledge/SKILL.md` carries what may
and may not be done with what it says.

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

Every member holds an allocation band: a numeric block that member allocates
from alone. Band N takes intents from `N*100` and plans from `N*1000`, so band 5
starts at `i500` and `p5000`. Band 0 is a band like any other, covering
`i001`-`i099` and `p0001`-`p0999`. No two members hold the same band, which is
what lets two clones that cannot see each other allocate without colliding.

```sh
context-circuit-cli --workspace <root> member add --id ID --name NAME --band 2
context-circuit-cli --workspace <root> member band --id ID --band 3
```

The CLI suggests a band and the person confirms it or picks another. Suggestions
start at 1, never offer band 0, and never offer a band whose range already holds
records, so a suggested band starts empty. A person may still pick any unused
band by hand, band 0 or one holding old records included: allocation skips every
number already taken. The band decides which number comes next and never makes a
reserved number reusable; only `record delete` releases one. A band is not a
namespace: the ID stays global.

A roster written by an earlier version may hold members with no band. `status`
then reports `band_required` and a `band_suggestions` entry for each of them,
and every change is refused until each has one. Offer each member their
suggested band, record what the person confirms with `member band`, and return
to their request.

Two clones may band the same members at the same time. When their choices
clash — `members.yaml` left with Git conflict markers, two members on one band,
or one member given two bands — `check` reports it as needing a person. Show the
person both versions and let them choose; never pick a winner for them.

## The language a member writes records in

A member may record the language their intents and plans are written in. It is
the language's name as a person says it, because a dispatch brief quotes it into
a sentence:

```sh
context-circuit-cli --workspace <root> member language --id ID --language 'Bahasa Indonesia'
```

Unset means English. Knowledge is English whatever this says — a note outlives
the member who wrote it — and an approval or completion note keeps the words the
person actually used. The brief a worker receives names this language and tells
it that what goes into the repository is English, because a worker works under a
repository's instructions and never reads these.

Beside it, a member may record the register that language is written in:

```sh
context-circuit-cli --workspace <root> member tone --id ID --tone 'semi-formal; keep technical terms in English'
context-circuit-cli --workspace <root> member tone --id ID --clear
```

One line of prose, quoted into the instruction that writes a record rather than
parsed. Naming a language settles which words a record uses and nothing about
how they are put together, and an agent that composes in English and translates
produces prose that is grammatical and formal in a way nobody chose. Record it
when a team knows the register it wants; unset leaves the general composition
guidance in `.agents/skills/cc-intent/SKILL.md` to decide. Ask the member for
the wording rather than inventing a register on their behalf.

## The roles a member contributes in

A member may record roles: free text saying what they contribute, such as
maintainer or docs, so people reading the roster can see who works on what.

```sh
context-circuit-cli --workspace <root> member add --id ID --name NAME --band N --role docs
context-circuit-cli --workspace <root> member roles --id ID --role maintainer --role docs
context-circuit-cli --workspace <root> member roles --id ID --clear
```

`member roles` replaces whatever roles were recorded. Record the roles the person
names, in their words; there is no list of allowed values, and none is
invented for them. Roles are informational and decide nothing: never treat one
as authority to approve, plan, execute, deliver, complete, or delete. They are
not the subagent `--role` that `agent configure` and `agent dispatch` take.

Roles are written only at workspace schema 3. While a workspace written by an
earlier version is still mid-setup, `member add --role` is refused: add the
member without roles, finish setup, and set them with `member roles`, by which
time the next change has moved the workspace to schema 3.

Bands prevent collisions only between members who actually hold distinct ones.
Two clones of the same member, and any clone whose roster is stale, can still
allocate the same number. Synchronize the shared workspace
before allocating, and resolve conflicting allocations and references before
publishing their IDs. Do not claim distributed collision prevention beyond what
bands give.

`.context-circuit/docs/workspace.md` describes which files are shared, which
belong to one machine, and the guarantees editing them carries.
