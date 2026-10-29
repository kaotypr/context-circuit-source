This workspace is for **Acme**, a project for invoicing small businesses.

It works with one product repository, **`api`**, whose default branch is `main`. The workspace itself is also stored in a repository, on `main`.

Run Context Circuit commands from the workspace root through its launcher:

```sh
sh .context-circuit/cli COMMAND …
```

For example, `sh .context-circuit/cli status` shows workspace status. The launcher uses the version pinned for this workspace.

The project knowledge says an invoice is **voided only while it’s a draft**. Once issued, it’s **credited with a credit note**.