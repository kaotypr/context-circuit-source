# Command surface and output contract

The seam made concrete. An agent learns this surface from the shipped
instruction and docs, so a command that drifts from its documentation is a
command the agent will get wrong.

## Shape

```
init · status · check
member add|band|use|list
repo connect|clone|init|base|remote|relate|fetch|inspect
workspace connect|base|remote
record create|show|list|approve|complete|archive|dependencies|order|delete
context find
extension add|update|remove|enable|disable|sync|list|hooks|allow
worktree prepare|list|inspect|move|repair|remove
template export
agent settings|configure|setup|dispatch
version
```

Global flags — a workspace root and a structured-output switch — go before the
command or after it, because a caller adding one to a line already typed puts it
at the end; a command defining a flag of the same name would keep its own
meaning, and none does. Every command accepts its own help, and a flag it does
not take is answered with that command's signature and help topic rather than
the bare flag list a parser prints.

## Properties that make it usable by an agent

**Structured output on demand.** The structured switch yields machine-readable
values; the default is readable YAML derived from the same public names. Errors
go to the error stream with a nonzero exit status and are never mixed into a
structured result.

**Every verb is one bounded operation.** Preparation prepares and does not also
install dependencies or start work. Completion appends a note, stamps a date,
and returns candidate catalog entries, but edits no knowledge. Composition is
the agent's job, in the open.

**Refusals are explicit and typed.** An ambiguous binding, a dependency cycle,
an unknown repository, a branch already holding work, a symlink inside the
knowledge tree, a reuse path pointing at control files — each is an error naming
its reason, never a silent best-effort substitution. A refusal whose answer is a
list the person must read, such as a deletion blocked by linked records, still
returns that list as its result beside the reason, and still exits nonzero.

**The caller's own work comes first.** Orientation opens with the active
member's open work and the count of their finished records not archived, and
the record listing defaults to that member's records, with another member's or
everyone's on request. A gate that leaves files uncommitted, archiving
included, names them as a commit recommendation, and completion carries an
offer to archive; both are said to the person, never acted on unasked.

**Read-only commands take no lock.** Inspection, listing, settings, dispatch
specification, export, and the diagnostic run without serializing the workspace;
everything that writes runs under the workspace lock.

Four commands fire an extension stage themselves — approve, prepare, complete,
and archive. Each runs its `before` point after its setup gate and precheck and
before taking the write lock, and its `after` point once it has succeeded and
released the lock. A `before` point that stops leaves the command unrun: it
exits 1, having written nothing, and still prints the block.

The command adds the `extensions` list of hook blocks to its result. The list
is absent where the workspace lists no extensions or nothing applied, so a
result keeps the shape it always had. `worktree prepare` returns its worktree embedded beside that
list for the same reason. `check` names a remote extension this machine has
not fetched under `unchecked`, apart from `issues`, without changing the exit
code.

## Documentation parity

The shipped reference cannot drift from the implemented surface. A test extracts
every complete invocation from the reference's fenced examples and checks it
against the real command table; a usage template teaching shape rather than a
runnable line is recognized and skipped.

This is why a new command lands with its documentation in the same change: the
test fails otherwise, which is the intended coupling.

Owner:

- `context-circuit-source@internal/cli/cli.go` — the command table, its flags,
  and the output contract.
- `context-circuit-source@product/docs/commands.md` — the reference a workspace
  is shipped.
- `context-circuit-source@internal/cli/extensions.go` — the stage hooks a
  command fires and `extension hooks`.
- `context-circuit-source@internal/cli/documentation_test.go` — the check that
  holds one to the other.
