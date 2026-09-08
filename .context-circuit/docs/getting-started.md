# Getting started

Context Circuit is a universal project workspace. You talk to it in ordinary
language; it holds Product Knowledge and readable plans and executes changes safely
across one or more Git repositories.

There are **two decisions you make**, and only two: you approve an **intent** —
what "correct" means and what scope is in bounds (Gate 1) — and you authorize
**delivery** — the irreversible act of shipping (Gate 2). Everything between them
is mechanical and scaled to risk: a plan derives from the approved intent, an
independent check runs when the risk warrants it, you accept the result, and you
mark the plan done when you ask to.

## 1. Initialize

> Initialize this workspace for <project name and purpose>.

This records project identity in `workspace.yaml` and creates the minimal
structure. It does not invent Product Knowledge or create plans.

## 2. Connect repositories

> Connect the api repository at ../commerce-api on branch development.

Registration records portable logical identity in `workspace.yaml` and a
host-local binding (path + `base_branch`) in `repositories.local.yaml`. New
project repositories can be cloned or initialized under the git-ignored
`repositories/<id>/`. The workspace root itself may be bound as the reserved
`workspace` repository. `base_branch` — not `default_branch` — is the
execution base and default pull-request target.

Each machine also names **who you are** in the shared member roster once, in a
gitignored local identity file — the same host-local idea as repository
bindings. The committed roster maps members to number bands; pick an existing
member (or add the first one). You never type a block number.

## 3. Gather context (optional)

> Gather context about billing from sources/billing-requirements.md.

Context files are written in place and the retrieval catalog is kept consistent.
There is no separate accept step.

## 4. Describe and approve an intent (Gate 1)

> Add recurring billing.

The coordinator writes up what it understands you want — the goal, what's out of
scope, and the acceptance criteria — has those criteria challenged adversarially
first, and tells you the risk tier ("Standard risk, so it gets an independent
check"). You approve it in one conversational decision, or say what to change.
Approving the intent is the real decision; you are not asked to approve a plan
afterward — the plan derives from what you approved.

For a small change you want to judge live, you can work directly instead:

> Work with me on tightening the checkout screen in the web repository.

Direct collaboration is the **Explore tier**: one worker in a separate working copy
of one connected repository, no plan and no independent check; the result is
human-supervised, never verified. If it turns out to be real work, the coordinator
offers to **promote** it in place — attach an intent, raise the tier, and it gains
an independent check — without restarting.

## 5. It builds and checks itself

Once the intent is approved the coordinator builds it and, at Standard/Critical,
has it independently checked against your exact change. You are not asked to
approve each step. When it reports "built and independently checked," you try it and
accept it — "looks right, ship it."

If the work turns out to need scope you did not approve (a new repository or path),
the coordinator pauses and asks — it never widens what you approved on its own.

## 6. Deliver (Gate 2)

> Open a pull request for the billing change.

Delivery is the second and final decision — pull request, merge, push, deployment —
always separate and explicit. A pull request uses each change's branch as source
and the recorded `base_branch` as the default target. Asking for a pull request,
merge, or sync does not mark the plan done and does not start a knowledge
update. When you ask to mark the plan done, live context files are updated only
if that plan affected Product Knowledge. The next plan can start without waiting.
