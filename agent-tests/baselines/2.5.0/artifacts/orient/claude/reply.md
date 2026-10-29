**What the workspace is for**

This is a Context Circuit workspace for **Acme**, whose stated purpose is "Invoicing for small businesses." It is a Git repository that holds the project's shared knowledge, intents, and plans, separate from the code. The code lives in the repositories it binds to. The current roster has one member, Maya (band 1), and this machine is set up as Maya.

**Repositories**

- **`api`**: the only logical repository listed in `workspace.yaml`, default branch `main`. This machine's checkout is at `repositories/api`, with base branch `main`.
- The workspace repository itself is also on `main`.
- There are no relationships recorded between repositories.

**Running the Context Circuit command**

Run every command through the workspace launcher, from the workspace root:

```sh
sh .context-circuit/cli status
sh .context-circuit/cli --json status      # structured output
sh .context-circuit/cli help               # full command surface
```

Don't use a `cli` found on your PATH for this workspace. The launcher runs exactly the version pinned in `.context-circuit/CLI_VERSION`, which is **2.5.0-rc.1**. **That version is not installed on this machine**, so the launcher will refuse and name the install step. Installing it is covered by the `cc-cli` skill (`.agents/skills/cc-cli/SKILL.md`). I didn't install anything, since you asked for no changes. Until the CLI is installed, I can't run the CLI commands here. Reading files still works.

**When an invoice is voided rather than credited**

The project's note is `context/domains/invoices.md`, owned by `api@internal/invoice/`. It says:

> An invoice is voided only while it is a draft; once issued it is credited with a credit note.

So the rule is:
- **Draft:** void.
- **Issued:** credit it with a credit note.

The note doesn't give the reasoning behind the split, and it doesn't cover edge cases such as partial credits. The catalog entry for it (`context/INDEX.md`) says it was last reviewed 2026-01-15. The code is the place to confirm the behavior, starting with `api`'s `internal/invoice/` package.

Two things in `context/INDEX.md` to be aware of: the lines near the top that mention `invoice-lifecycle.md` are template examples, and that note does not exist. The real invoice entry is the one near the bottom, which points to `domains/invoices.md`.