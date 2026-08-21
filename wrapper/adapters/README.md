# Context Circuit

Start or resume work in this workspace in natural language. The agent reads a
small entry spine, selects the smallest safe evidence packet, and tells you
what it can do next.

The human surfaces are:

- `README.md` for the quick start;
- `plans/<repository>-plans/<plan>/PLAN.md` for the complete human plan;
- the conversational “what’s next” card;
- the latest surfaced runtime handoff.

Approval never starts execution. Say `Approve plan <name>` and later
`Run approved plan <name>`. Finishing, delivery, publication, deployment,
archive, takeover, and runtime cleanup are separate confirmations.

The workspace is filesystem-first and works offline. Never put credentials in
workspace files. Inspect `wrapper/manifest.yaml` for version and budgets.
