# Agent-facing commands

Run the CLI version this workspace pins in `.context-circuit/CLI_VERSION`. Versions
install side by side, so the shared command on PATH may select a different one; the
`cc-cli` skill resolves the pinned path. The CLI warns on stderr when it is run
against a workspace pinning another version.

Use `context-circuit-cli --workspace PATH --json COMMAND ...` for structured output.
The global flags `--workspace PATH` and `--json` may precede the command or follow
it, as in `context-circuit-cli status --workspace PATH --json`. Run
`context-circuit-cli help` for the full surface, and `help` with a command name for
that command's signatures and what it decides, refuses, or leaves to the caller. A
fully resolved command also takes `--help`, which prints that same help; a flag it
does not take is named, followed by the command's signature and the `help` topic
to read, with exit status 2. A group such as `member` needs its subcommand first,
and reports which option it could not take. Errors go to stderr with a nonzero
exit status.
Normal output is readable YAML; `--json` gives the agent structured values.
The examples assume the current directory is the initialized workspace.

```sh
context-circuit-cli status
context-circuit-cli member add --id alex --name Alex --band 2
context-circuit-cli member band --id alex --band 3
context-circuit-cli member language --id alex --language 'Bahasa Indonesia'
context-circuit-cli member tone --id alex --tone 'semi-formal; keep technical terms in English'
context-circuit-cli member roles --id alex --role maintainer --role docs
context-circuit-cli member use --id alex
context-circuit-cli member list
context-circuit-cli repo connect --id api --path ../api --base main
context-circuit-cli repo clone --id web --path repositories/web --base main --url GIT_URL
context-circuit-cli repo clone --id api --path repositories/api --base main
context-circuit-cli repo init --id docs --path repositories/docs --base main
context-circuit-cli repo base --id api --branch development
context-circuit-cli repo remote --id api --url GIT_URL --default-branch main
context-circuit-cli repo relate --from web --to api --description 'Consumes the API'
context-circuit-cli repo inspect --id api
context-circuit-cli repo fetch --id api --remote origin
context-circuit-cli workspace connect --base main
context-circuit-cli workspace base --branch development
context-circuit-cli workspace remote --url GIT_URL --default-branch main
context-circuit-cli knowledge clone --id core-service-knowledge --url GIT_URL
context-circuit-cli knowledge connect --id core-service-knowledge --path ../core-service-knowledge --index index.md
context-circuit-cli knowledge sync --id core-service-knowledge
context-circuit-cli knowledge remote --id core-service-knowledge --index docs/index.md
context-circuit-cli knowledge list
context-circuit-cli context find --query billing
context-circuit-cli context find --repo api
```

A knowledge repository is one this workspace reads and never owns. It records no
base branch, takes no worktree, and is never named by a plan, and it shares one
ID namespace with `repositories` so a note's anchor means one thing. `clone`
puts it under `repositories/<id>` unless `--path` says otherwise, and takes no
`--url` for an ID the workspace already describes. `sync` fetches, then
fast-forwards only when the checkout is clean and on the shared branch; a dirty
or diverged checkout is reported and left exactly as it is, and nothing is ever
merged, reset, or discarded. Obtaining one disables its push URL and every sync
reasserts that, so an edit made here fails rather than half-succeeding —
improvements go upstream from a separate checkout.

Records:

```sh
context-circuit-cli record create --kind intent --slug add-billing --title 'Add billing'
context-circuit-cli record approve --id i100 --text 'User approved the billing outcome'
context-circuit-cli record create --kind plan --slug billing-api --title 'Billing API' \
  --intent i100 --repo api
context-circuit-cli record create --kind plan --slug billing-web --title 'Billing UI' \
  --intent i100 --repo web --depends-on p1000
context-circuit-cli record create --kind plan --slug specified-change --title 'Specified change' \
  --repo api
context-circuit-cli record show --id p1000
context-circuit-cli record list
context-circuit-cli record dependencies --id p1001 --depends-on p1000
context-circuit-cli record complete --id p1000 --text 'User requested completion.'
context-circuit-cli record order --intent i100 --mode waves
context-circuit-cli record order --plan p1002 --mode waves
context-circuit-cli record delete --member alex
context-circuit-cli record delete --member alex --confirm
context-circuit-cli record delete --all
context-circuit-cli record delete --id p1002
context-circuit-cli record delete --id i101 --id p1003 --confirm
```

`context find` takes one of the two: a query searches this workspace's catalog
and every borrowed index, marking which side each match came from and naming any
index it could not read, and a
repository returns the entries claiming it together with `reconcile_required`,
which is what completion names for a plan and what a change made without one
still owes. Those entries are always this workspace's own: a borrowed entry
names an obligation nobody here can discharge, so it is never offered as
knowledge to reconcile. An unknown repository is refused rather than answered
with an empty list, because nothing to reconcile and nothing found read the
same.

Every member holds a band, and `member add` requires `--band`; without it, it
refuses and names the next free band. A member holding band N allocates intents
from `N*100` and plans from `N*1000`, so band 2 writes `i200` and `p2000`. Band 0
is a band like any other and starts at `i001` and `p0001`, never `i000` or
`p0000`. `member band --band 0` sets band 0; there is no clearing a band. Bands
must be distinct across the roster, band 0 included, and every command refuses
to run while two members share one. Allocation skips every number already
reserved in any member's ID file or held by a record file, inside any band, and
a full band is refused naming the member, the band, and its range. The new
number is reserved in the author's own `.context-circuit/ids/MEMBER.yaml` and in
no other file. Assigning or changing a band never renumbers an existing record or
releases its reservation.

`record delete` is the one command that releases reservations. `--member` takes
every intent and plan that member created, active and archived. `--all` takes
every record, including any an earlier version wrote without an author, and
empties every ID file while keeping the files; it is refused unless the roster
lists exactly one member. `--id` names the records instead, and is repeatable:
each intent or plan with that ID, active or archived, approved or completed,
whoever wrote it or with no author at all, and a reserved ID no record file holds,
shown as `reservation only`. It cannot be combined with `--member` or `--all`,
and an ID matching neither a record nor a reservation is refused. Without
`--confirm` it returns the records, each with its author, path, and `status`
(`draft` or `approved` for an intent, `open` or `completed` for a plan), and
`released_ids` with `confirmation_required`, and deletes nothing; with it, the
record files and plan folders are removed, plans before intents and each
`plan.md` after the rest of its folder, then each released ID is removed from
every ID file holding it, and `sync_required` says the deletion lives in this
checkout until it is committed and pulled everywhere, and that every other
machine should then run `check` and remove or prune what it names before
creating plans. It refuses, deleting nothing, while a
record outside the selection would be left naming a deleted one — a plan under a
deleted intent, a dependency, or an intent linking a deleted plan — or while a
worktree on this machine was prepared for a deleted plan. With `--id`, links
between the named records do not block, and a deleted plan is taken out of its
own intent's `plans:` list and reported under `unlinked` instead of refusing;
any other link refuses and the refusal names the `--id` values to add so the
records are deleted together. A refusal, with or without `--confirm`, still
returns the preview on the result stream, in YAML or `--json`: the records asked
for with their author, path, and status; `blocked_by`, each record outside the
selection that points at one of them, with its `id`, `created_by`, `path`,
`status`, and `link` (such as `belongs to intent i100`, `depends on p1000`, or
`links plan p1000`); `refused`, the reason; and, by `--id`, `add_ids`, the IDs to
add. The reason is also printed on stderr, nothing is deleted or released, and
the exit status is 1. A plan its intent lists that no record file holds
any more is unlinked the same way. In every mode, `warnings` also lists what each
bound checkout still holds for a deleted plan number — local and remote-tracking
`cc/pNNNN/*` branches, worktrees on such a branch, and a leftover
`.worktrees/pNNNN` directory — without refusing and without touching them; a
plan given that number again meets the old branch. The 2.3.0 CLI does not know
`--id`, rejects it, and deletes nothing. No ID is released until every file is
gone, so an interrupted deletion leaves every selected number still reserved:
a failed removal names the path and the call that finishes the job, `record
delete --id ID ... --confirm` with exactly the IDs still reserved, or `record
delete --all --confirm` after `--all`, and running the same `--id` or `--all`
deletion again finishes it too. A plan folder that has already lost its
`plan.md` is no longer a record and no re-run reaches it, so the error names it
for removal by hand, and `check` reports it as an empty plan folder. A failure while
rewriting the ID files, after every record file is gone, names the IDs still
reserved and the same finishing call, or says nothing is left when none is.

A workspace written by a CLI before 2.4.0 is at schema 2, with every reservation
in one shared `.context-circuit/ids.yaml`. The first changing command run once
setup is complete splits it into one file per member, and so does the setup step
that completes setup — `member band`, `member add --band`, or `member use` —
right after it succeeds, rather than leaving the split to the next command. A
split that setup step cannot make is a warning, not a failure of the step, and
the next changing command refuses with the reason. Each reservation goes to
the file of its record's `created_by`, and one with no author or no record file
goes to the member selected on this machine. The shared file is removed and
`workspace.yaml` moves to `version: 3`. A shared file that comes back later, from
a clone that had not pulled the split, is split again on the next changing
command, skipping numbers a member file already holds. The split is reported as
one `notice:` paragraph on the error stream, in `--json` mode too, so no
command's result changes shape. The notice names every shared file the command
changed, from what actually changed: `workspace.yaml`, the ID files written or
removed, and `members.yaml` when the command changed it, as `member band` does.
Records the command writes are in its own result. `status`, `check`, the other reading commands,
and `record delete` without `--confirm` never split. Nor do a few refusals
checked before the split: `member add` without `--band`, a malformed `member
roles` call, and a confirmed `record delete` its preview would refuse. A refusal
found while the command runs, such as an unknown member or record ID, comes
after the split, which then stands; it is the same split any later change makes. `init` writes the per-member
layout itself.

The split is the one time a CLI rewrites an existing workspace's files on its
own, and it is a cut-over: CLIs before 2.4.0 refuse a schema 3 workspace with
"requires a newer CLI", so every member needs 2.4.0 before they pull it. It
writes a file for every author, so only one person on a team upgrades first:
they make the first changing command, then commit and push the files the notice
names; committing is the person's decision. Everyone else pulls
before making any change. Teammates who upgrade at the same time get a Git
conflict in `.context-circuit/ids/*.yaml`, which is theirs to resolve by hand:
keep the union of both sides' entries in each member file, then run `check`,
which reports any number reserved twice. If a merge leaves Git conflict markers
in `ids.yaml`, the split refuses until a person resolves them, keeping every
entry. A member or record author whose ID is a Windows device name stops the
split, which names them; `status` and `check` report a roster member like that
beforehand, and the rename is made by hand.

Setup is complete when the roster lists at least one member, every member holds a
band, and this machine has selected one of them. Until then every changing
command is refused with a message naming the step, except the setup steps
themselves: `init`, `member add`, `member use`, and `member band`. `status`,
`check`, and the other read-only commands keep working. `status` reports
`identity_required` when the roster is empty or unreadable, or this machine has
selected nobody on it or someone it does not carry; `band_required` with
`band_suggestions`, one suggested band per member without one, when a member
holds no band; and `next_band`, the band to offer a member added next.
Suggestions start at 1, never offer band 0, and skip any band whose range already
holds a taken number. `check` reports a `members.yaml` holding Git conflict
markers, one it cannot parse — a member listed twice or given two `band:` keys —
and two members holding one band as needing a person.

Use the actual returned IDs, not the example numbers. `--repo` and `--depends-on`
are repeatable. Creation reserves IDs and links a new plan from its intent when
`--intent` is supplied. Otherwise it creates a standalone plan folder.
Approval and completion notes require the actual corresponding user decision.
Approval records the instant in the intent's frontmatter `approved_at` and appends
the person's words under a heading carrying the same instant; an unapproved
intent has no such field. It records a decision and never establishes one, and
commands do not independently authorize anything. Approval returns
`planning_recommended`, a recommendation to put to the person rather than a duty
to start planning: plans follow a planning request, which may come with the
approval. Completion returns `reconcile_required` whenever the plan's
repositories hold catalog entries. Each names what comes due after its gate,
because a gate reporting success has recorded a decision and not finished a
request. A CLI older than 2.4.0 returns `planning_required` from approval; that
CLI is behind the template, and `AGENTS.md` wins over its text. Dependency updates detect cycles; omit `--depends-on` to
clear a plan's dependencies. Completion also
records `completed_at` in the plan's frontmatter beside its written note. Every
record also carries the `created_at` instant of its creation. Instants are
canonical ISO 8601 UTC timestamps, `2026-09-15T10:53:00Z`; other dates are ISO
8601 `YYYY-MM-DD` in UTC.

After it succeeds, approval reads the workspace's Git status for the approved
intent's file, and completion for the plan's files and its intent's file,
wherever each lives, archives included; both also read every ID file, the
shared `.context-circuit/ids.yaml` and each per-member file. Whatever is not
committed is returned as `commit_recommended`: `files`, each a `path` relative
to the workspace root with `/` separators and a `change` of `new`, `modified`,
or `deleted`, and a `note`. A renamed record is listed as its new path `new` and
its old path `deleted`. Other uncommitted records are not listed. Approval and
completion always rewrite the record they touch, so in a Git workspace that
record is always listed. The key is omitted when nothing is listed. When Git is
not installed, or the workspace is not in a Git repository, the result carries
`commit_unchecked` with the reason instead. Neither ever blocks or undoes the
command, and the CLI never stages, commits, or pushes: committing the records
still needs the person's authorization. A workspace inside another repository
is reported relative to its own root, leaving out files outside it; one that
repository ignores lists nothing.

`record order` derives execution shape from recorded dependencies and completion.
Repeat `--plan ID` to select an exact run, or use `--intent ID` to select linked
plans; the two filters are exclusive. With neither, all active plans are selected.
The result names `plan_ids` and `selection`. Unfinished dependencies outside an
explicit selection are blocked, never added silently.
`--mode auto` reports a recommendation and follows it; `waves` and `linear`
select explicitly without changing the recommendation. Waves list each layer, its
concurrency, every plan's start reference per repository, any integration merge a
fan-in needs, and which plans in the layer share a repository. A linear chain
stacks each plan on the previous one and never needs a merge. Completed plans
drop out and their dependents become ready, so an unfinished plan holds
everything behind it; recompute after each wave. Start references name a
predecessor's branch — adjust them yourself once that work is merged and its
branch is gone. The report is advice: it runs, merges, and reserves nothing, and
`shares_repository_with` reflects repositories, not paths.

Worktrees:

```sh
context-circuit-cli worktree prepare --repo api --plan p0001 --copy-mode auto
context-circuit-cli worktree prepare --repo web --plan p0002 --start origin/main
context-circuit-cli worktree list --repo api
context-circuit-cli worktree inspect --repo api --path .worktrees/p0001/api
context-circuit-cli worktree move --repo api --path .worktrees/p0001/api --to ../api-billing
context-circuit-cli worktree repair --repo api --path ../api-billing
context-circuit-cli worktree remove --repo api --path ../api-billing
context-circuit-cli worktree prune --repo api
context-circuit-cli check
```

`worktree prune` drops the plan-to-worktree associations of one repository whose
path Git no longer lists, such as a worktree removed with `git worktree remove` or
by hand. It never removes a worktree, branch, or file. Its `dropped` list names
each association it dropped, with `repository`, `plan`, `path`, `branch`, and
`start_commit`. `worktree remove` also removes the folders under `.worktrees/`
that the removal left empty, up to and including `.worktrees/`, and lists them as
`removed_directories`; a non-empty folder, or one outside `.worktrees/`, is left
alone, and one that could not be removed is a warning rather than a failure.
`worktree prepare` refuses a plan branch already checked out
elsewhere, and says where and what is safe to do next.

A prepared worktree reports `start`, what it actually started from: `--start` as
given, the base branch or `origin/<base>` it fell back to, or the existing branch
`--reuse` checked out. `base_branch` stays this machine's base, the branch the
work is delivered to, which a worktree started elsewhere need not share. Neither
form of start fetches. A start derived from a base behind its fetched remote is
reported in `sync_required`, and so is an explicit `--start` naming a local branch
that is behind its remote-tracking branch or has diverged from it: the worktree
still starts from the local branch, and `sync_required` gives the counts, how to
start from the remote instead (`--start origin/<branch>`), and that `repo fetch`
updates the remote ref. A start ahead of its remote, a remote-tracking ref, or a
commit is not questioned.

`check` reports findings as an `issue` and the `resolve` that discharges it:

```yaml
issues:
  - issue: knowledge repository core-service-knowledge is not obtained on this machine
    resolve: "`knowledge clone --id core-service-knowledge`, which needs no --url because the workspace already describes it"
  - issue: "duplicate record ID: p0003"
    resolve: "needs a person: two clones allocated the same number; renumber one record
      and every reference to it before either is shared, and keep both reservations"
  - issue: "i050 is reserved in more than one ID file: .context-circuit/ids/alex.yaml, .context-circuit/ids/maya.yaml"
    resolve: "needs a person: two clones each moved it into their own member's file; keep it
      in the file of the member who created the record and remove the extra entry by hand"
ok: false
```

For ID files, `check` reports a file it cannot read, a malformed or misplaced
ID, a number reserved twice in one file or across member files, a record whose
ID no file reserves, and a file in `.context-circuit/ids/` not named after a
member ID, whose numbers still count as taken. Numbers outside a member's band,
and files of members no longer on the roster, are valid and not reported.

`check` also reads each catalog entry's `reviewed` date against the code its note
anchors to, and reports a note whose anchors have moved since — with how many
commits. A note whose code has not moved is not reported, however old its date
is. `knowledge_review_days` in `workspace.yaml` additionally reports a note that
anchors to no code and has gone that long unconfirmed; unset, those are silent.

A resolution names the command where one exists, the edit where an edit is the
whole of it, and opens with `needs a person` where the next step is somebody's
judgment. That last kind is the one worth reading closely: it means no command
discharges the finding, so none should be improvised in its place. `check` still
reports and exits, and is never a gate.

`--reuse` explicitly selects an existing branch/worktree, without resetting it.
Without it, `worktree prepare` refuses a worktree or branch left from an earlier
plan that held the same number, and the refusal says how to remove it; preparing
again for the same plan resumes its worktree. Preparation derives no order, so a plan that records a dependency in the same
repository is refused without `--start`, and the refusal names the predecessor
branch to pass — or says which predecessor has not been implemented and committed
yet. An explicit start is always honored, including one that is not the
predecessor. `--discard` on removal explicitly authorizes disposing of its files;
never add it merely to bypass a failure. No command auto-deletes branches, commits, delivers,
runs application setup, executes a plan, or starts independent review.

`template export --path NEW_DIRECTORY` writes the exact blank embedded seed for
inspection or packaging. `init` can initialize that blank seed or a fresh directory.
It cannot migrate an existing v1 or reinitialize an active v2 workspace.

Subagent setup and dispatch:

```sh
context-circuit-cli agent settings
context-circuit-cli agent configure --host codex --role worker --model MODEL_ID --effort high
context-circuit-cli agent setup
context-circuit-cli --json agent dispatch --host codex --role planner \
  --intent i001 --path /actual/repository
context-circuit-cli --json agent dispatch --host codex --role planner \
  --task 'Plan the specified outcome' --path /actual/repository
context-circuit-cli --json agent dispatch --host codex --role worker \
  --plan p0001 --repo api --path /actual/worktree
context-circuit-cli --json agent dispatch --host codex --role explorer \
  --path /actual/repository --task 'Which module owns invoice voiding?'
```

Settings default to host inheritance. `agent setup` applies settings to native
agent files for every host unless `--host` narrows it, because one workspace is
opened in several; the host may need to reload them. Initialization runs it, so a
new workspace needs it again only after `agent configure` or in a fresh clone,
where the ignored definitions never arrive. `agent dispatch` returns the
resolved invocation for the cc-dispatch skill to launch through the host's tools.
`--task` is the assignment for a role no record assigns: an explorer's question,
a reviewer's diff and criteria, or the one slice a worker owns when several share
a plan. A planner is dispatched against a whole intent and refuses it. A worker
owning its whole plan needs none, and a brief carries no task section when none
was given. For a reviewer, `--review-requested` represents a real user request.
No CLI command launches an LLM or treats a dispatch specification as evidence of
completion.

A specification also reports `definition_path` and `definition_installed` for the
host's native role file, and `setup_required` naming the exact `agent setup` run
when that file is absent. Role files are host-local and gitignored, so a workspace
that never ran setup on this host names an agent type the host cannot resolve. It
is a report, not a refusal: the prompt remains complete, so a live spawn tool can
still carry it when native roles are unavailable.

The returned `prompt` is the complete brief and is launched unmodified: it opens
with the working directory, then ownership, then the quoted record, then the task,
and ends with what the role must return. Quoting the record in full is what stops
an agent from searching a repository for a file it was never handed.

`--plan` adds the plan's file path, repositories, dependencies, and Markdown body to
the brief, and states that its dependencies' work is already in the branch's
ancestry and that a concurrent sibling plan is invisible. For a worker, `--repo`
selects the assigned repository's required supporting files along with shared
files; the brief lists workspace-relative paths under the absolute workspace root
and tells the worker to read them before editing. A multi-repository plan requires
`--repo`. `--intent` does the same
for an approved intent and is how a planner is dispatched; a planner refuses
`--plan`, because the plan shape is the answer it returns rather than a number it
is handed, and it refuses an intent that is not yet approved. `--shared` is for
several workers inside one worktree: it replaces sole ownership of the directory
with a duty to preserve the other workers' edits. Omit it when each worker has its
own worktree.

Worktree `--copy-path` adds a specific ignored runtime entry. `--copy-mode` accepts
`auto`, `required`, `copy`, or `off`. The reuse report distinguishes native CoW
from fallback copying and reports existing entries or dependency mismatches.
