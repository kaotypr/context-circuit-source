Nothing was changed. I only read files.

**What this workspace is for**

It's the Context Circuit workspace for Acme, described as "Invoicing for small businesses." It's a separate Git repository that holds the project's shared context, goals (`intent/`), plans (`plans/`), member roster, and repository bindings. The code itself lives in other repositories.

**Which repositories it covers**

- `api`, default branch `main`. This is the only code repository listed in `workspace.yaml`. Its remote is a local path, so the URL won't work on your machine.
- The workspace's own repository, also on `main`.
- No relationships are recorded between repositories.

The roster in `members.yaml` has one member, Maya, band 1. Your own identity is in `member.local.yaml`, which I didn't open.

**How to run the Context Circuit command here**

Run commands through the workspace's launcher from the workspace root, not a `cli` found on your PATH:

```sh
sh .context-circuit/cli status
sh .context-circuit/cli help
sh .context-circuit/cli --json status
```

On Windows PowerShell or cmd, use `powershell -NoProfile -ExecutionPolicy Bypass -File .context-circuit/cli.ps1 COMMAND`. The launcher runs the version pinned in `.context-circuit/CLI_VERSION`, which is **2.6.0**. If that version isn't installed, the launcher refuses and names the install step. I haven't checked whether 2.6.0 is installed on this machine. Running `status` would tell you.

**When an invoice is voided rather than credited**

The project's knowledge note `context/domains/invoices.md` says:

- An invoice is **voided only while it is a draft**.
- Once it is **issued**, it is **credited with a credit note** instead.

The note's anchor is `api@internal/invoice/`, which covers voiding, credit notes, and the plain-text summary. I didn't open that code, so this answer rests on the note alone.

The note is short. It doesn't say why the split exists, what happens to an issued invoice that was never sent, or whether anything else can void an invoice. If you need those rules, check the `api` code or ask the team. The catalog entry in `context/INDEX.md` also says the note was reviewed on 2026-01-15, but the note itself carries no review date.