I'll look up the workspace identity, its repositories, how the Context Circuit command is run here, and the knowledge note on voiding an invoice versus crediting it — read-only.This workspace is **Acme**, for invoicing for small businesses. Nothing was changed.

**Repositories.** The product repository is `api` (default branch `main`). This workspace has its own Git repository as well, also on `main`. There are no recorded relationships between repositories, and no borrowed knowledge centers.

**Running Context Circuit.** From this workspace root, run commands through the pinned launcher, not a `context-circuit` binary on your PATH:

```sh
sh .context-circuit/cli COMMAND …
```

From anywhere else, use `sh <workspace-root>/.context-circuit/cli COMMAND …`. Add `--json` when you want structured output, for example `sh .context-circuit/cli status --json`. `sh .context-circuit/cli help` lists the commands.

**Void versus credit.** An invoice is voided only while it is a draft. Once it has been issued, it is credited with a credit note. Invoices export only as a plain-text summary. That rule is recorded in `context/domains/invoices.md` and applies to `api` (`api@internal/invoice/`).

The context-circuit plugin here is 0.2.4 (installed as 0.2.2). To update: `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.