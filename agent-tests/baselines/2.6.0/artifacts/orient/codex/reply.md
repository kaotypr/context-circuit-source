## Workspace

This is the **Acme** workspace, for invoicing for small businesses. It works with one project repository, **`api`**, whose default and local base branch is `main`. The workspace also has its own repository for workspace records and coordination; its base branch is `main`.

## Running Context Circuit

From the workspace root, run commands through the launcher:

```sh
sh .context-circuit/cli status
```

The launcher uses the version pinned in `.context-circuit/CLI_VERSION` (`2.6.0`). For example, to search project knowledge:

```sh
sh .context-circuit/cli context find --query invoice
```

## When an invoice is voided or credited

The project knowledge says an invoice is **voided while it is still a draft**. Once it has been issued, it is **credited with a credit note**. This rule is documented in [the Invoices note]($PRIVATE/agent-test-3180962479/ws/context/domains/invoices.md).

I only read workspace files; I didn’t run any commands or change anything.