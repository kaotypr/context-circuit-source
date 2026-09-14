# Context Circuit v2.0.0

Source design for Context Circuit v2.0.0 — the **rewrite** that replaces the v1
shell runtime with a separately released Go executable and a smaller instruction
surface. Unlike [../v1.0.0/](../v1.0.0/), which was an argument to test before it
was built, this grouping documents a design **as built**: it was written against
the 2.0.0-rc.1 candidate and describes decisions already expressed in
`product/`, `internal/`, and `template/`.

Read [../v1.0.0/](../v1.0.0/) first if you need to understand what v2 declined to
carry forward; `core/retired-machinery.md` is the accounting of exactly that.
This README is the version index; each scope owns its own design.

## Scopes

- [core/](./core/) — the whole v2 coordination model: the seam between the Go
  executable and the agent, records and IDs, the knowledge circuit, worktrees
  and copy-on-write reuse, stacked execution, roles and dispatch, authority and
  delivery, versioning and distribution, and what v1 machinery was retired.
  Start at [core/design.md](./core/design.md).

## Layout convention

`sources/system-design/<product>/<version>/<scope>/`. A **scope** is a bounded
area of the design, cut by concern and never by repository. The version folder
holds one scope folder per area it covers; this README indexes them and carries
the reading order.

Every folder has a **`README.md`** as its landing and index. Each scope's
normative design is **`design.md`**; the scope's `README.md` states purpose and
reading order. Detail splits into `<concern>.md` files as a concern outgrows a
section of `design.md`.

## Source boundary

These documents are maintainer design material. They have no status, no
authority, and change nothing on their own. They are not copied into an
instantiated workspace, are not Product Knowledge, and are not normal agent
context. Where a document here disagrees with `product/AGENTS.md.in`,
`product/docs/`, or the code, the owner wins and this material is stale.
