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
- `context-circuit-source@internal/cli/documentation_test.go` — the check that
  holds one to the other.
