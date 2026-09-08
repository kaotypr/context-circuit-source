# Context Circuit

A project setup you open in an AI coding agent and talk to in ordinary language.
It remembers what your project is and helps you make decisions based on that.
It turns a request into a grounded change across one or more Git repositories,
and keeps the two decisions that matter with you: what "correct" means, and when
to ship.

Works with **Claude Code**, **Codex CLI**, and **Cursor**.

## Start here

1. Copy or clone this repository into a new folder.
2. Open that folder in Claude Code, Codex, or Cursor.
3. Say:

> Initialize this for <project name and purpose>, and connect the <repo> repository.

That records the project and connects your code. It does not invent knowledge
or change any repository on its own.

Then either work live on a small change:

> Work with me on <small change> in <repository>.

Or describe a real change and approve it:

> I want to add <feature> — <what correct looks like>.
> Approve that.
> Build it, then open a pull request.

The longer walkthrough is [Getting started](.context-circuit/docs/getting-started.md).

## What you decide

You make two decisions. Everything else is mechanical.

- **What correct means.** You approve that from the plain ask. The agent then
  looks at the real code and writes the plan. When you ask it to build, it makes
  the change and independently checks it when the risk warrants it.
- **When to ship.** Opening a pull request, merging, or pushing is always a
  separate ask. Checking the work never ships it.

For a small change you want to judge as it happens, work together instead — no
plan and no independent check; you are watching it live. If it turns into real
work, it can be promoted in place.

## What you can read

- [Getting started](.context-circuit/docs/getting-started.md) — initialize, connect, change, ship
- `plans/` — readable plans for work in progress
- `intent/` — what you approved as correct

This folder is a blank, uninitialized workspace. It works offline and stores no
credentials. Branches and working copies are managed for you.
