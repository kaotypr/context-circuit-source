# The executable and agent seam

The division that defines v2.

> Anything that must happen identically every time belongs to the executable.
> Anything requiring interpretation belongs to the agent, in the conversation,
> where a person can see it.

Neither side crosses. The executable does not interpret; the agent does not
hand-roll what the executable owns — it does not invent a number, improvise a
branch name where preparation has a convention, or decide work is complete
because it looks complete.

## What the executable owns

| Concern | Owner |
| --- | --- |
| Identity, reservations, record creation and notes | `context-circuit-source@internal/workspace/records.go` |
| Workspace and member records, local bindings | `context-circuit-source@internal/workspace/workspace.go` |
| Repository connect, clone, init, base, relate, fetch | `context-circuit-source@internal/workspace/repositories.go` |
| Git state inspection and the diagnostic | `context-circuit-source@internal/workspace/inspect.go` |
| Worktree preparation, listing, move, repair, removal | `context-circuit-source@internal/workspace/worktrees.go` |
| Ignored runtime reuse and filesystem cloning | `context-circuit-source@internal/workspace/reuse.go`, `internal/cow/` |
| Execution order derivation | `context-circuit-source@internal/workspace/order.go` |
| Role settings, native role files, dispatch specifications | `context-circuit-source@internal/workspace/agents.go` |
| Seed embedding and materialization | `context-circuit-source@assets.go` |

## What it refuses

The refusals are the product boundary, not gaps waiting to be filled. The
executable holds no model credentials and calls no model API. It launches
nothing — a dispatch specification returns a flag saying a launch is still
required, and is never evidence that an agent ran or finished. It runs no
application setup, executes no plan, and never commits, pushes, opens a pull
request, merges, deploys, or deletes a branch; it supplies the repository,
branch, and base facts that ordinary Git and provider tools act on. Its approval
and completion operations *record* a decision a person already made rather than
constituting one. Its diagnostic reports and exits, and nothing waits on it.

The shipped help text carries these refusals verbatim, so they are visible to
whoever is about to use the tool.

## What the agent owns

Everything interpretive, and the list is long precisely because the executable's
is short: understanding the request, retrieving the right knowledge, writing an
intent's goal and success criteria, deciding which repositories a change
touches, reading real code, writing plans, choosing whether isolation is
warranted, finishing the environment setup a reuse report left open, deciding
what to delegate, launching and integrating subagent results, running and
interpreting checks, judging which durable knowledge a completed plan changed,
and asking for authorization where a person decides.

It also owns the honesty obligations that cannot be compiled in: reporting what
was implemented, tested, and left uncertain; distinguishing a written role file
from a loaded one; and never calling its own inspection an independent review.

## Why a compiled executable

The predecessor was a large POSIX shell runtime. It forced Windows users through
WSL, edited structured records as text, had no practical test story, and grew as
one file with every capability — it was the ceiling on the product. Go supplies
one static binary per platform with no runtime to install, cross-compilation to
six targets from one machine, native access to the per-platform cloning calls
that environment reuse needs, and tests that establish file and Git behavior in
disposable fixtures.
