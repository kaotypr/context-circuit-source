# Context Circuit

Start or resume work in this workspace in natural language. The agent reads a
small entry spine, selects the smallest safe evidence packet, and tells you
what it can do next.

The human surfaces are:

- `README.md` for the quick start;
- `plans/<repository>-plans/<plan>/PLAN.md` for the complete human plan;
- the conversational “what’s next” card;
- the latest surfaced runtime handoff.

Approval never starts execution and never commits Git. Say `Approve plan
<name>` to see a session-bound card; nothing has changed yet. The exact
confirmation is `Confirm approval of plan <name>`. On an instantiated or
wrapped workspace, later say `Run approved plan <name>`. On product-source,
if confirmation leaves only the approval status projection dirty, the existing
maintainer commit card is presented next; it still requires
`Confirm commit of the approved plan state.` Finishing, delivery, publication,
deployment, archive, takeover, and runtime cleanup are separate confirmations.

The workspace is filesystem-first and works offline. Never put credentials in
workspace files. Inspect `wrapper/manifest.yaml` for version and budgets.
