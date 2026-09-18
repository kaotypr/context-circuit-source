<p align="center">
  <img src=".context-circuit/assets/readme/context-circuit-logo.png" alt="Context Circuit" width="480">
</p>

<p align="center">
  <strong>Give your coding agent project context and a grounded plan—while you stay in control across every repository in your project.</strong>
</p>

<p align="center">
  <a href="#quick-start">Quick start</a> ·
  <a href="#how-to-use">How to use</a> ·
  <a href="#stacked-plans-parallel-when-safe-ordered-when-necessary">Stacked plans</a> ·
  <a href="#what-context-circuit-adds">Capabilities</a> ·
  <a href="#documentation">Documentation</a>
</p>

Coding agents are good at changing code. The harder problem is carrying the
right project context across sessions, repositories, people, and machines.

Context Circuit gives an agent a shared place to learn the project, help prepare the
goal before coding, coordinate work across repositories, update product knowledge
when the code changes, and preserve what the team learned afterward. You stay in
control of the decisions that matter.

<p align="center">
  <img src=".context-circuit/assets/readme/workflow-overview.png" alt="How Context Circuit helps an agent finish safely" width="840">
</p>

## Why Context Circuit

A prompt can be very specific about what to change while leaving out why the
current behavior exists, which product rules it must preserve, or what else the
change will affect. Without that wider understanding, an agent can follow the
request faithfully and still produce a locally correct change that conflicts
with the product or breaks a business flow.

Context Circuit does not promise to prevent every mistake. It helps the agent
read the project's accepted product knowledge and repository relationships
before the goal is approved. The agent writes down the outcome it understood,
points out conflicts or missing decisions, and turns questions it should not
answer alone into explicit open questions for a person.

Without this shared project context, every new session starts by rediscovering
the same architecture and conventions. A change that crosses repositories
becomes a collection of disconnected tasks. Plans drift away from the outcome,
and useful decisions disappear into chat history.

Context Circuit keeps those pieces connected:

- **A grounded goal** — the requested change is read against what the project
  already knows, so assumptions, contradictions, and open questions are visible
  before implementation begins.
- **One outcome across repositories** — a solo developer or team can coordinate
  several repository-specific plans as one intended product change, with their
  dependencies and execution order kept explicit.
- **Living product knowledge** — completed changes are reconciled with the
  project's architecture, rules, decisions, conventions, and vocabulary, giving
  the next goal more accurate context.
- **Shared project context** — product knowledge, repository relationships,
  intens (goals), plans, dependencies, and current work remain available across sessions,
  people, and machines.
- **Human control** — approval, execution, delivery, cleanup, and other consequential
  actions remain explicit human choices; unresolved product decisions are not
  quietly made by the agent.

The workspace is ordinary Markdown and YAML committed with your project. The
companion CLI handles dependable bookkeeping and Git mechanics; the coding agent
still does the reasoning, implementation, and project-specific checks.

## Quick start

### 1. Create and open the workspace

Use [Tembiter](https://www.npmjs.com/package/tembiter) to create a clean project
from a tagged Context Circuit release:

```sh
npx tembiter init \
  --template https://github.com/kaotypr/context-circuit.git \
  --tag v2.0.0 \
  --target acme-workspace

cd acme-workspace
```
> Replace "acme-workspace" with your preferred workspace folder name.

Tembiter copies that release without carrying over the template repository's Git
history, initializes a new repository, and records which template version it came
from. It requires Node.js 20 or later.

Open the new directory in Codex, Claude Code, or Cursor. The workspace includes
instructions and skills for each host.

### 2. Ask the agent to install and initialize

```text
Install the Context Circuit CLI, then initialize this workspace as Acme Billing
for our billing platform. Add me as Maya.
```

The bundled installation skill selects the correct macOS, Linux, or Windows
package, verifies its checksum, and installs it without administrator access.

### 3. Connect the repositories

Tell the agent which repositories belong to the project, where each one can be
found, and how they relate. You do not need to decide where Context Circuit
stores a repository obtained from a Git URL:

Think of the **`default branch`** as the repository's home branch on its origin—
usually `main` or `master`. It is shared because it normally means the same thing
for everyone and is the usual starting point for new work.

The **`base branch`** is this machine's current working anchor. It may be the
default branch, or a branch belonging to a person or feature already in progress.
Context Circuit starts new worktrees and plan execution from that anchor. You can
change it when you move to another feature or work on several efforts in
parallel. Base branch settings stay local and Git-ignored, so changing
your anchor never changes a teammate's.

```text
Add these repositories to the project:

- The billing API is at https://github.com/acme/billing-api.git. Its default
  branch on the origin is main. For the recurring-billing feature I am working
  on, use feature/recurring-billing as this machine's base branch. Creates one
  from main if not exist
- The billing web app is already on this machine at
  /Users/maya/Work/acme/billing-web. Its origin also defaults to main, and this
  machine base branch should be feature/recurring-billing.

The web app consumes the API.
```

When you provide a Git URL, the agent obtains and organizes the working copy for
you. When you provide a full local path, it uses that existing checkout in
place. It can also initialize an empty repository when you are starting
something new.

### 4. Gather the project context

Before starting the first change, ask the agent to learn the durable knowledge
already present across every connected repository:

```text
Gather project context from all connected repositories and writes product knowledge.
```

This gives later goals a useful starting point instead of making each new agent
session rediscover the product from code alone. The knowledge stays readable and
versioned in the workspace, where it can be corrected and kept current as the
project changes.

### 5. Optional: tune the subagent roles

Context Circuit can delegate bounded work to four roles. The defaults inherit
whatever your coding host chooses, so configuration is optional. To make the
tradeoffs explicit, ask the agent to recommend supported settings before it
writes them:

```text
Set up shared subagent role tiering for Codex. Check which models and reasoning
effort levels this host actually supports, then recommend exact settings using
these priorities:

- Explorer: a fast, economical model with low effort for narrow codebase research.
- Planner: a strong reasoning model with high or extra-high effort because
  every downstream plan depends on its analysis.
- Worker: a capable, balanced model with medium or high effort for implementation
  and normal repository checks.
- Reviewer: a strong model with high or extra-high effort

Show me the recommendation before applying it. If I approve, write the shared
configuration and set up the native role definitions for this host.
```

Shared preferences live in `.context-circuit/role-tiering.yaml` and travel with
the workspace. A member whose host, account, or budget differs can override only
the affected roles on their machine:

```text
Keep the team's shared role tiering, but configure a local worker override for
this machine. Show me the supported choices first, write the approved override
locally, and refresh the native role definitions.
```

Local overrides live in the Git-ignored
`.context-circuit/role-tiering.local.yaml`; roles absent from that file continue
tracking the shared configuration. Replace `Codex` in the prompt with Claude Code
or Cursor when that is the host you use.

Applying the settings also creates or updates the coding host's native role files:

- Codex: `.codex/agents/cc-*.toml`
- Claude Code: `.claude/agents/cc-*.md`
- Cursor: `.cursor/agents/cc-*.md`

These files are ignored by default because they are generated from the role
tiering and may include model or effort choices specific to one person's host,
account, or budget. Ignoring them prevents one member's local setup from becoming
everyone's setup; another member generates the right files for their own host.

A solo developer who wants the generated definitions to travel with the
workspace can remove the matching `<host>/agents/cc-*` rule from `.gitignore`
and commit those files. They then become project-owned files: review changes to
them like any other configuration, and remember that a fresh clone may preserve
them as existing customizations rather than replacing them automatically.

### 6. Start working on your project

The workspace is ready. Describe the goal you want to achieve to your coding
agent in ordinary language; you do not need to translate it into Context Circuit
commands or split it into repository tasks yourself.

Continue to [How to use](#how-to-use) to see how the agent grounds that goal,
surfaces decisions it should not make alone, prepares dependent plans, and waits
for your authorization before implementation or delivery.

## How to use

Work from the Context Circuit workspace and talk to your coding agent in ordinary
language. You normally do not need to run Context Circuit commands yourself. The
agent uses the workspace's pinned CLI for records and Git mechanics, then explains
the decisions that remain yours.

### 1. Ground the request

Describe the outcome, not a list of internal Context Circuit steps:

```text
Add recurring billing to the API and web app. Monthly plans only; keep the
existing payment provider.
```

The agent retrieves the relevant product knowledge and repository relationships,
writes the goal, calls out assumptions or missing decisions, and stops. Answer its
numbered questions and correct anything it misunderstood:

```text
1. Existing customers stay on their current plan.
2. Failed renewals get a seven-day grace period.
The rest of the goal looks right. I approve it.
```

Approval lets the agent inspect the actual repositories and prepare grounded
plans. It does not start implementation.

### 2. Read the plan—or ask for the summary

The agent shows which repositories will change, which plan depends on another,
what checks it expects to run, and what remains uncertain. You can inspect the
plan files or stay in the conversation:

```text
Summarize the plans, their dependencies, and the main risks.
```

If the outcome has changed, say so before execution. The agent updates the goal
and returns to approval instead of quietly expanding the work.

### 3. Start the work explicitly

When the plans look ready, ask the agent to run them:

```text
Run the billing plans using the recommended execution order.
```

That request starts implementation. The agent prepares isolated working copies,
runs independent plans together where safe, waits for unfinished dependencies,
follows each repository's instructions, and runs the relevant tests, lint, and
builds. It reports actual results and labels anything it could not verify.

### 4. Continue or recover without losing work

You can leave and return to the workspace later:

```text
Show me the status of the billing plans and resume the unfinished work.
```

The agent checks the real branches, working copies, commits, and diffs before it
continues. Failed or partial work stays available; Context Circuit does not reset,
stash, or discard it to make the records look clean.

### 5. Request review, delivery, and completion separately

Implementation does not silently become publication. Ask for each consequential
next step when you want it:

```text
Review the completed changes independently.
Open pull requests for the API and web plans.
Mark the plans complete now that they have landed.
```

Review is optional and read-only. Delivery uses the real work branches and target
branches. Completion records the result and updates any product knowledge the
change made stale; it does not automatically delete branches or working copies.

For a small, already-specific edit, you can ask the agent to work directly in a
connected repository without creating a goal or plan:

```text
Make this small copy change directly in the web checkout. Do not create a plan.
```

The same rules for preserving existing work and asking before delivery still
apply.

<p align="center">
  <img src=".context-circuit/assets/readme/knowledge-circuit.png" alt="The project learns from every change" width="840">
</p>

For the behavior behind these conversations, including the division between the
coding agent and CLI, read [How Context Circuit works](.context-circuit/docs/how-it-works.md).

## Stacked plans: parallel when safe, ordered when necessary

Large changes rarely fit into one repository or one uninterrupted agent run.
Context Circuit turns one approved goal into several linked plans, makes their
dependencies explicit, and runs the AI agents responsible for those plans in the
right order.

This is more than spawning several agents at once. Before an agent starts,
Context Circuit gives it the complete plan it owns, the repositories it may
change, its dependency information, and the correct working copy and starting
point. A plan whose dependency is unfinished stays blocked instead of running
against the wrong code.

For example, a billing launch might become:

```text
Launch recurring billing
├── Add the billing API
│   ├── Build the web checkout
│   └── Add the admin controls
└── Update the notification worker
```

Context Circuit can run the independent API and notification work together,
then release the web and admin plans when the API they need is complete:

```text
Wave 1  Billing API ─────────────┐
        Notification worker     │ run together
                                │
Wave 2  Web checkout            ├ start after the API
        Admin controls ─────────┘
```

It supports two execution shapes:

- **Waves** overlap independent plans to shorten wall-clock time. When branches
  converge, Context Circuit identifies the integration work that must happen
  before a dependent plan starts.
- **A linear chain** stacks every plan on the previous one. It is slower but
  avoids integration merges, which is often the safer choice when several plans
  modify the same repository.

The recommended shape is derived from every plan's recorded dependencies and
repository overlap, with its cost explained before anything runs. You confirm
the shape once. After each wave, Context Circuit recalculates what is ready from
the work that actually completed, then gives the next agents the dependency and
starting-point information they need.

If one plan fails, completed plans are not unwound. Their branches and commits
remain available, the blocked plan is named, and everything waiting behind it
stays paused. That makes parallel agent work faster without making it disposable.

## Workspace repository structure

A Context Circuit workspace is itself a small Git repository. It holds the
shared understanding and coordination files for the project; the application
repositories it coordinates can live beside it, inside its ignored
`repositories/` directory, or elsewhere on the machine.

```text
context-circuit/
├── README.md                    Product guide
├── AGENTS.md                    Standing workflow and safety rules
├── CLAUDE.md                    Claude Code entry instruction
├── CURSOR.md                    Cursor entry instruction
├── workspace.yaml              Project identity and repository relationships
├── members.yaml                People who share this workspace
├── context/                    Product knowledge shared through Git
│   ├── INDEX.md                Searchable knowledge catalog
│   └── glossary.md             Shared project vocabulary
├── intent/                     Goals, boundaries, and open questions
├── plans/                      Repository-specific plans and dependencies
├── sources/                    Optional evidence supplied for a task
├── repositories/               Optional local application checkouts (ignored)
├── member.local.yaml           Active person on this machine (ignored)
├── repositories.local.yaml     Local paths and starting branches (ignored)
├── .context-circuit/
│   ├── docs/                   Detailed workspace and CLI documentation
│   ├── assets/readme/          README artwork
│   ├── role-tiering.yaml       Shared defaults for AI-agent roles
│   ├── role-tiering.local.yaml Optional machine overrides (ignored)
│   ├── local/                  Host-local settings and locks (ignored)
│   ├── VERSION                 Workspace-template version
│   └── CLI_VERSION             CLI version this workspace expects
├── .agents/skills/             Canonical Context Circuit agent skills
├── .claude/                    Claude Code skills and agent roles
├── .codex/                     Codex agent roles
└── .cursor/                    Cursor rules, skills, and agent roles
```

Files marked as ignored are recreated or connected on each machine. Everything
else can be reviewed and shared like normal project documentation.

## What Context Circuit adds

| Capability | What it means for you |
| --- | --- |
| Product knowledge | Architecture, conventions, decisions, and domain rules are retrieved when relevant instead of rediscovered every session. |
| Stacked multi-repository plans | Every AI agent receives its plan dependencies, waits for prerequisites, and runs in the correct parallel wave or linear order. |
| Isolated implementation | Each plan gets a separate working copy by default, protecting unrelated and unfinished work. |
| Human control points | The agent waits for approval of the goal and a later request to begin implementation. Delivery and cleanup remain separate choices. |
| Preserved failures | Partial work and failed checks are kept and reported instead of reset, stashed, or hidden. |
| Optional sub-agents | Exploration, implementation, and independent review can be delegated when useful; they are not mandatory ceremony. |
| Actionable diagnostics | Workspace checks explain both the problem and the command, edit, or human decision needed to resolve it. |
| Reviewable files | Shared state is Markdown and YAML that your team can diff, discuss, and version in Git. |

## What is shared

<p align="center">
  <img src=".context-circuit/assets/readme/shared-vs-local.png" alt="What is shared and what is set locally" width="840">
</p>

Project purpose, repository relationships, approved goals, plans, and product
knowledge travel through Git. Machine paths, active AI-agent roles,
temporary locks, dependencies, and unfinished edits stay on the machine where
the work is running.

That split lets a team share the truth about the project without pretending that
every developer, container, or remote agent has the same local setup.

## Safety model

Context Circuit is designed to preserve user work and make consequential actions
visible:

- No implementation begins before you explicitly request it.
- No force checkout, hard reset, automatic stash, or silent overwrite is used to
  clear a path.
- Failed and partial implementation remains available for inspection and resume.
- Pushes, pull requests, merges, deployment, publication, and cleanup require
  explicit authorization.
- Independent review happens only when requested and does not silently modify
  code.
- Checks report what actually ran; anything not exercised is reported as
  unverified.

## Documentation

- [How Context Circuit works](.context-circuit/docs/how-it-works.md) — agent
  behavior, human decisions, and the CLI mechanisms supporting each stage
- [Workspace files](.context-circuit/docs/workspace.md) — what the workspace
  stores and which information stays local
- [Commands](.context-circuit/docs/commands.md) — the complete CLI surface
- [Working records](.context-circuit/docs/working.md) — the detailed goal and
  plan file formats
- [Worktrees](.context-circuit/docs/worktrees.md) — isolated work, environment
  reuse, recovery, and cleanup
- [Sub-agents](.context-circuit/docs/agents.md) — available roles, configuration,
  and dispatch

`AGENTS.md` carries the standing safety rules your agent reads whenever it works
in this workspace. Skills under `.agents/skills/` provide the procedure for the
stage currently in progress.

## Requirements and limits

- Git is required. Users do not need Go, Python, or a package manager.
- CLI packages are available for macOS, Linux, and Windows on amd64 and arm64.
- Each operating system, container, or remote host installs the CLI and records
  its local repository paths independently.
- Context Circuit v2 is for fresh workspaces. It does not automatically migrate
  v1 workspace data.
