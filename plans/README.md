# Plans roadmap

Numbered peer plans live under `plans/<repository-key>-plans/`. Each plan is a human-reviewed set of YAML metadata, Markdown explanation, and task files.

Active plans use `<number>-<slug>`. Manual archival moves the exact directory to `archives/plans/<repository-key>-plans/<number>-<slug>/`; archived plans are ignored by `whats-next`.

Agents draft plans as Markdown and YAML, present them for explicit human
approval, and execute approved work through session records, leases, and
worktrees. There is no command-specific plan lifecycle.
