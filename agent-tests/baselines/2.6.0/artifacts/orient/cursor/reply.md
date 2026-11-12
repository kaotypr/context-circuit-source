I'll look up the workspace identity, its repositories, how the Context Circuit command is run here, and the knowledge note on voiding versus crediting an invoice — read-only.This workspace is **Acme**, for invoicing for small businesses. Nothing was changed.

**Repositories.** The product code lives in one logical repository, `api`, on the default branch `main`. Its remote is `$PRIVATE/agent-test-4087327256/remotes/api.git`. There are no recorded relationships between repositories. This workspace itself is a Git repository on `main`, with remote `$PRIVATE/agent-test-4087327256/remotes/ws.git`.

**Running Context Circuit.** From the workspace root, run every command through this workspace’s launcher so you get the pinned CLI (`2.6.0`), not whatever happens to be on `PATH`:

```sh
sh .context-circuit/cli COMMAND …
```

From anywhere else, use `sh <workspace-root>/.context-circuit/cli COMMAND …`. Add `--json` when you want structured output, for example `sh .context-circuit/cli status --json`. `sh .context-circuit/cli help` lists the full surface.

**When an invoice is voided.** Project knowledge says an invoice is voided only while it is still a draft. Once it has been issued, it is credited with a credit note. Invoices export only as a plain-text summary. That rule is recorded in `context/domains/invoices.md` and applies to `api`, last reviewed 2026-01-15. The owning code is `api@internal/invoice/` (voiding, credit notes, and the plain-text summary).

The Context Circuit plugin here is 0.5.0 (0.2.2 is installed). You can update it with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.