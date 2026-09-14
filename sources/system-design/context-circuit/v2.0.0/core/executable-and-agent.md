# The executable and the agent

The seam that defines v2. `design.md` states the principle (P2); this file
specifies where the boundary falls, why it falls there, and what the CLI surface
guarantees to each side.

## The rule

> Anything that must happen identically every time belongs to the executable.
> Anything requiring interpretation belongs to the agent, in the conversation,
> where a human can see it.

Neither side may cross. The executable does not interpret; the agent does not
hand-roll what the executable owns.

## Why a compiled executable

v1's `wrapper/runtime/engine.sh` was ~2,000 lines of shell. It worked, and the
read that produced the v1.0.0 study concluded its *mechanics* were sound. What
was not sound was the substrate:

| Shell runtime | Consequence |
| --- | --- |
| POSIX shell + `yq`-style external tooling | Windows required WSL |
| Text manipulation of YAML | Comment and ordering loss, fragile edits |
| No type system, no test framework worth the name | Behavior changes landed unverified |
| Every capability grew one file | The engine was the ceiling on the product |

Go was chosen for a specific combination: a single static binary per platform
with no runtime to install, a real standard library for filesystem and process
work, cross-compilation to six targets from one machine, native access to the
per-OS cloning syscalls CoW reuse needs (`clonefile`, `FICLONE`, ReFS block
cloning), and a test story that lets file and Git behavior be established in
disposable fixtures.

The user-visible payoff is stated in one line of the release notes: *Windows no
longer needs WSL.*

## What the executable owns

| Concern | Operations | Package |
| --- | --- | --- |
| Identity and IDs | Allocate, reserve permanently, never reuse | `internal/workspace/records.go` |
| Records | Create, link plan↔intent, append notes, stamp `completed` | `internal/workspace/records.go` |
| Workspace state | Read/edit `workspace.yaml`, `members.yaml`, local bindings | `internal/workspace/store.go`, `workspace.go` |
| Repositories | Connect, clone, init, set base, relate, fetch, inspect | `internal/workspace/repositories.go` |
| Git state | Branch, HEAD, dirtiness, worktree inventory | `internal/workspace/inspect.go` |
| Worktrees | Prepare, list, inspect, move, repair, remove | `internal/workspace/worktrees.go` |
| Environment reuse | Discover ignored entries, CoW clone, fall back, report | `internal/workspace/reuse.go`, `internal/cow/` |
| Ordering | Derive waves, chain, start refs, integration merges | `internal/workspace/order.go` |
| Role settings | Resolve tiering, write native files, build dispatch specs | `internal/workspace/agents.go` |
| Diagnostics | `check` — links, duplicates, cycles, knowledge boundary | `internal/workspace/inspect.go` |
| Seed | Embed and materialize the blank workspace | `assets.go` |

### What it refuses

The refusals are the product boundary, not gaps to be filled later:

- **No LLM credentials, no model API calls.** The executable is model-blind.
- **No launching.** `agent dispatch` resolves an invocation and returns it with
  `launch_required: true`. The CLI has not started an agent, and a dispatch
  specification is never evidence that one completed.
- **No application setup.** It does not run package managers, migrations,
  services, or arbitrary setup hooks.
- **No plan execution.** No command executes a plan or a task.
- **No delivery.** No command commits, pushes, opens a PR, merges, deploys, or
  deletes a branch. It supplies the repository, branch, and base facts that
  ordinary Git and provider tools then act on.
- **No authorization.** Approval and completion commands *record* a decision a
  human already made; they do not constitute one.
- **No gate.** `check` reports and exits. Nothing in the product waits on it.

The help text carries the refusals verbatim, so they are visible to whoever is
about to use the tool:

> Approval, completion, fetch, and removal are explicit operations. The caller
> must hold the user's authorization. No command executes plans, starts review,
> commits, pushes, merges, deploys, or automatically removes completed worktrees.

## What the agent owns

Everything interpretive, and it is a long list precisely because the executable's
list is short: understanding the request; retrieving the right knowledge; writing
the intent's goal, non-goals, constraints, and success criteria; deciding which
repositories a change touches; reading real code; writing plans; choosing whether
isolation is warranted; finishing environment setup the reuse report left open;
deciding what to delegate; launching, waiting, and integrating subagent results;
running and interpreting the repositories' checks; judging which durable
knowledge a completed plan actually changed; and asking for authorization where
the design says a human decides.

The agent also owns the honesty obligations that cannot be compiled in: reporting
what was implemented, tested, and left uncertain; distinguishing a written role
file from a loaded one; and refusing to call its own inspection an independent
review.

## The CLI as interface contract

The command surface is the seam made concrete. Three properties make it usable
by an agent rather than only by a person:

**Structured output on demand.** `--json` yields machine-readable values;
default output is readable YAML. Global flags precede the command. Errors go to
stderr with a nonzero exit status, never mixed into a structured result.

**Every verb is one bounded operation.** `worktree prepare` prepares; it does not
also install dependencies or start work. `record complete` appends a note, stamps
a date, and returns candidate catalog entries; it does not edit knowledge.
Composition is the agent's job, in the open.

**Refusals are explicit and typed.** An ambiguous binding, a cycle, an unknown
repository, a branch already holding work, a symlink inside `context/`, a reuse
path pointing at control files — each is an error with a reason, never a silent
best-effort substitution.

### Surface at a glance

```
init · status · check
member add|use|list
repo connect|clone|init|base|relate|fetch|inspect
record create|show|list|note|approve|complete|dependencies|order
context find
worktree prepare|list|inspect|move|repair|remove
template export
agent settings|configure|setup|dispatch
version
```

`product/docs/commands.md` is the shipped reference and
`internal/cli/documentation_test.go` holds it to the implemented surface.

## Data handling

Two rules the executable enforces so the agent cannot accidentally break them:

- **Environment contents are opaque.** `.env` files are cloned or copied without
  being read into any report. The reuse report names entries and counts files;
  it never prints content. Contents never enter a prompt, a log, or a shared
  record.
- **Control files are out of reach.** A reuse path containing `.git` or
  `.context-circuit` is rejected, as is any non-repository-relative path. The
  executable never copies credential stores or host configuration.

## File-level guarantees

- **Atomic single-file replacement** plus a portable OS lock protect cooperating
  commands in one directory. The lock is released on process exit including
  abnormal exit; the empty lock file remains.
- **YAML edits are document edits** through `goccy/go-yaml`, retaining
  surrounding comments and ordering. Presentation may normalize; byte-for-byte
  preservation is not promised.
- **Unknown fields and duplicate keys in structured control records are
  errors.** User prose stays ordinary Markdown with no schema registry.
- **Multi-file operations may leave partial output** on interruption, and IDs are
  reserved before file creation and never rolled back. The command reports the
  path so the work can be resumed or repaired — the alternative, silently
  reusing a reserved number, is worse.

## Testing the seam

Deterministic tests in `internal/` use temporary directories and disposable Git
repositories to establish file and Git behavior, the documented command surface,
order derivation, CoW copying, and the installer including its upgrade and
rollback paths. `scripts/check-release.sh` builds and exercises a native binary,
verifies the embedded seed inventory against the manifest, cross-compiles every
supported target, and tests the publication guards. CI runs the native suite on
Linux, macOS, and Windows, because cross-compilation alone is not proof of
behavior on another operating system.

What no test establishes is the other side of the seam. Whether a coding host
loads `AGENTS.md` and follows it is not verified by any of this, and the design
requires saying so rather than letting green CI imply it.
