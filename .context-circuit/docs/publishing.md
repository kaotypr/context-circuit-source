# Publishing

Publishing reflects your Context Circuit data outward to a system your team already
uses — a task tracker (ClickUp, Jira, GitHub, Notion) or a chat space (Slack). It is
opt-in, manual, and separate from the core workflow: publishing never approves,
executes, verifies, or delivers anything, and it never changes a plan. You ask for
it in plain language and the `cc-publish` skill does the work.

In this product "publish" means sending data to an external system. Landing code in
git is "deliver" / "open a pull request" (see `cc-deliver`) — a different thing.

## Before you start

Connect the provider once, at the host: the ClickUp / Slack / … integration or
credentials. Configuration files never hold tokens, so the credential stays with the
host. If the provider is not reachable, publishing stops and tells you so.

## Setting up a publication

There is nothing to edit by hand. The first time you ask to publish, if no
publication is set up yet, `cc-publish` creates one and asks only for what it needs
(which list or channel, and any authoring instructions). After that it is reused.

Examples:

- "Publish plan 0100 to ClickUp." — sets up publishing to a ClickUp list, then
  publishes.
- "Set up publishing our plans to the Engineering list in ClickUp." — sets it up
  without publishing yet.
- "Open a Slack discussion for plan 0100's open questions in #eng." — sets up a chat
  discussion, then posts it.

You can have several publications at once (for example one to ClickUp and one to
GitHub), each named for what it publishes and where.

## Publishing and re-publishing

- "Publish plan 0100 to ClickUp." / "Push 0100 to our tracker."
- "Re-publish 0100." — updates the same items in place; it never creates duplicates.
- "Open the discussion thread for 0100's open questions."

The external copy is one-way: it reflects your plan, and edits people make in the
tracker or thread never flow back into your workspace.

## What gets published

- **A plan** becomes a work item with one child item per task, dependencies as
  "blocked by" links, and each task's acceptance criteria as a checklist.
- **Open questions** become a chat discussion — one message that opens the thread and
  one reply per question, for the team to answer in place.

Everything reads as ordinary project work: no file names, internal ids, or tooling
details appear in what you publish. Each item carries your plan's id (for example
`[0100]`) so you can tell what came from where.

## Instructions (optional)

You can tell `cc-publish` how to write, in one line:

- Language — "write in Bahasa Indonesia."
- Tone and terms — "everyday conversational tone; keep technical terms in English."
- Extra fields — "estimate a time estimate per task and a target date range," which
  it fills on the tracker as best-effort estimates.

## Where it is kept

Each publication lives under `publication/<name>/` — its `config.yaml` and a record
of what was published. You can open these, but you never have to; the conversation
maintains them, and nothing about publishing is written under `plans/`.
