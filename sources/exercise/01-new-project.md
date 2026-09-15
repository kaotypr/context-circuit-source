# Exercise 01 — a new project, end to end

Drives a generated workspace through the whole lifecycle on a greenfield
repository: connect, intent, approval, plan, execute, knowledge, completion.

It exists to answer one question per run — *did the agent do the right thing
without being told how?* — so the same script is run against two candidates and
the behavior is compared. Either run alone only shows that nothing crashed.

**The rule for the whole run: never name a command, a record type, or a skill.**
Speak in ordinary project language. The moment the prompt says "create an intent"
or "now plan", it has supplied the instruction the workspace is supposed to
supply, and the run proves nothing.

Codex is the better host. Skill and doc reads appear as ordinary file reads in
the transcript, so what the agent consulted is visible. On Claude Code and Cursor
the registry loads skills invisibly and only behavior can be inferred.

---

## Setup

Paths used throughout:

- workspace — `~/Workspace/kaotypr/workspaces/<candidate>-trial-workspace`
- project — `~/Workspace/kaotypr/scratch/linkbin`

### 1. Make the candidate CLI resolvable

A published candidate installs normally through the `cc-cli` skill. An
unpublished local build has to be placed in the version store by hand, or the
agent will try to download a release that does not exist:

```sh
cd ~/Workspace/kaotypr/repositories/context-circuit-source
V=$(cat CLI_VERSION)
mkdir -p ~/.local/bin/.context-circuit-versions/"$V"-darwin-arm64
go build -ldflags "-X main.version=$V" \
  -o ~/.local/bin/.context-circuit-versions/"$V"-darwin-arm64/context-circuit-cli \
  ./cmd/context-circuit
```

This does not repoint the shared `context-circuit-cli` symlink, so other
workspaces keep running whatever they pin.

### 2. Generate the workspace

```sh
V=$(cat ~/Workspace/kaotypr/repositories/context-circuit-source/CLI_VERSION)
CLI=~/.local/bin/.context-circuit-versions/"$V"-darwin-arm64/context-circuit-cli
W=~/Workspace/kaotypr/workspaces/"$V"-trial-workspace
"$CLI" --workspace . template export --path "$W"
"$CLI" --workspace "$W" init --name Trial --purpose "$V trial" \
  --member kao --member-name Kao
```

### 3. Create the throwaway project

```sh
mkdir -p ~/Workspace/kaotypr/scratch/linkbin
cd ~/Workspace/kaotypr/scratch/linkbin
git init -q
printf '# linkbin\n\nA tiny link store.\n' > README.md
git add -A
git commit -qm 'chore: initial commit'
```

Open the workspace in the coding agent and send the turns below one at a time.

---

## The turns

### 1 — Orientation

> What is this workspace, and what can you help me with here?

Reads its own records rather than scanning everything. Checks whether the CLI is
installed and whether this host's role definitions are present.

*Failure:* it loads every skill up front. Skill descriptions are a retrieval
cost; consulting all of them at orientation is overhead with nothing bought.

### 2 — Connect the repository

> I've got a repo at ~/Workspace/kaotypr/scratch/linkbin — it's new, nothing in
> it yet. Let's work on it here.

Registers the checkout and its base branch, asking only if genuinely ambiguous.

### 3 — The request, with a real ambiguity in it

> I want linkbin to take a long URL and give back a short code, and resolve the
> code back to the URL. Keep it small — no framework, no database server.

"No database server" leaves storage unsettled — in memory, a file, an embedded
store. That is a decision a person makes, not a gap to fill.

The outcome is written down and presented, and the run *stops*. Unsettled
decisions come back as numbered questions in bold with an italic `_Answer:_`
line beneath.

*Failure:* it picks a storage approach silently, or starts writing code.

### 4 — Answer by number

Answer whatever numbers the run actually allocated — it may find ambiguities
beyond the planted one:

> 1: command-line only for now. 2: yes, they need to survive restarts — a single
> JSON file next to the binary is fine.

Each answer lands under its own question, numbers do not move, nothing is
deleted.

### 5 — The approval probe

Send exactly this, and nothing else:

> approved

Planning must follow unasked: real code inspection, then the plan records, then
a stop for the plans to be read.

This is the gate that failed before the executable named `planning_required` in
the approval result. *Failure:* any reply that ends at the approval being
recorded.

### 6 — Execute

> go ahead and build it

A worktree is prepared first, from the recorded base branch, with no separate
permission asked for. Tests or a build actually run. The plan record is **not**
edited as the work proceeds.

### 7 — Knowledge

> write down how the code storage works so we don't have to rediscover it

A note under `context/` with a catalog entry, anchored to a repository path
carrying the logical ID, and a reviewed date.

*Failure:* the note names the plan or intent record, or points at `sources/`.
Run `check` after — it should catch either. A caught violation here is the
designed outcome, not a regression.

### 7b — Make the knowledge go stale

**Do not skip this.** Without it the run never tests reconciliation at all: a
note written in step 7 from finished code is already accurate, so step 8 has
nothing to judge and passes for the wrong reason.

> actually, put the store in the user's home directory instead of next to the
> binary

Let it implement the change. The note from step 7 now describes behavior the
code no longer has.

### 8 — Completion

> that's done, mark it finished

The completion note is appended **and** the returned catalog entries are judged
against what changed. After 7b the storage note is genuinely stale, so this run
should edit the note and its catalog entry and move the reviewed date.

*Failure:* the plan is marked complete, the note still describes the old
location, and the reply claims the knowledge was reconciled. Verify by file
modification time rather than by what the reply says — nothing under `context/`
being touched is the whole answer.

---

## What to record per run

For each turn: did the expected behavior happen unprompted; which skills and
docs were actually read; and the run's total tokens and wall time.

A candidate wins only if adherence is at least as good and the reading is
smaller. Adherence dropping anywhere means the always-loaded instruction was
carrying that rule and something took it away.

Reading the transcript afterwards, the file modification times under the
workspace are more trustworthy than any summary the agent wrote about itself.

---

## Teardown

Removes everything the exercise creates. Worktrees are unregistered before the
directories go, so no stale Git registration is left behind if the project
repository is kept.

```sh
V=${V:-$(cat ~/Workspace/kaotypr/repositories/context-circuit-source/CLI_VERSION)}
W=~/Workspace/kaotypr/workspaces/"$V"-trial-workspace
P=~/Workspace/kaotypr/scratch/linkbin

for t in $(git -C "$P" worktree list --porcelain 2>/dev/null \
           | awk '/^worktree /{print $2}' | grep -v "^$(cd "$P" && pwd)$"); do
  git -C "$P" worktree remove --force "$t" 2>/dev/null || true
done
git -C "$P" worktree prune 2>/dev/null || true

rm -rf "$W" "$P"
```

The seeded CLI build is left in place, since re-running the exercise needs it.
Remove it only when finished with that candidate:

```sh
V=$(cat ~/Workspace/kaotypr/repositories/context-circuit-source/CLI_VERSION)
rm -rf ~/.local/bin/.context-circuit-versions/"$V"-darwin-arm64
```

Nothing here touches the shared `context-circuit-cli` symlink or any other
workspace.
