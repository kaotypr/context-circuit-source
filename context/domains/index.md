# Domain Knowledge

`context/domains/` is an explicitly declared OKF bundle. This lowercase
`index.md` is its reserved routing index; it does not duplicate domain bodies.

Each domain directory is the canonical home for knowledge about one bounded
project area. The domain page owns exact behavior, workflows, rules,
interfaces, data, constraints, edge cases, implementation references, and
verification when those facts are supported by evidence.

Generated domain pages use
[the domain template](../../docs/templates/domain-context.md) and normally
live at `context/domains/<domain>/README.md`. That `README.md` is a concept,
not a competing index. Detailed workflow concepts may live below that domain
at `workflows/<workflow>.md`.

## Selecting domain context

For a domain request, first read this index and then only the named domain page
and the workflow pages it links to. Do not recursively read every domain or
the raw source inbox. If the request spans domains, select the smallest set
that covers the stated outcome and record the selection in the plan or handoff.

Domain pages can be `status: draft`, `stable`, or `deprecated`. The separate
Context Circuit acceptance extension remains authoritative. A refresh of an
accepted page must preserve accepted decisions; conflicting evidence is
surfaced for review rather than silently replacing it.

## Domain entries

Add one link per generated domain page, with a short routing description. Keep
this index navigational; do not duplicate domain facts here.
