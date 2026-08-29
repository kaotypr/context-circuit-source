# design-layout-grouping — design

## Capability

Generalize the middle segment of the system-design source path from a strict
*version* to a *grouping* dimension, fix the agent's default for choosing it, and
settle version-grouping number depth — so design source is organized predictably
whether a team thinks in releases, phases, or quarters.

## Problem

The v0.6 layout convention names the path
`sources/system-design/<product>/<version>/<scope>/` and exemplifies `<version>`
as `v0.1`. Two gaps surfaced in use:

- **The middle level isn't always a version.** A team may organize design work
  by delivery phase, quarter, or milestone rather than by release number. The
  convention offered no room for that, so the grouping intent was implicit.
- **Version depth was inconsistent.** Existing folders are 2-number (`v0.5`,
  `v0.6`) while every shipped surface the design targets is 3-number semver
  (`template_version: 0.6.0`, dist `context-circuit-v0.6.0`). A patch-level
  design revision had no unambiguous home until `v0.6.1` was introduced ad hoc.

## Generalized path

```
sources/system-design/<product-or-project>/<grouping>/<scope>/
```

- **`<product-or-project>`** — the product or project as a whole. May collapse to
  a single segment in a single-project workspace. **Never one folder per git
  repository** (that fragments the cross-repo coherence the design exists to
  capture) — unchanged from the `cc-system-design` rule.
- **`<grouping>`** — the organizing dimension for a body of design work.
  **Default: a version.** Optionally a named grouping when a team does not
  organize design by release.
- **`<scope>`** — a concern, unchanged.

## Fixed decisions

1. **Version is the default grouping**, and version groupings use **3-number
   semver with a `v` prefix going forward** (`v0.6.1`, not `v0.6`). Existing
   `v0.5` and `v0.6` remain as legacy 2-number folders — not retroactively
   renamed (per the "going forward only" decision).
2. **Agent default derivation.** When the author does not specify a grouping, the
   agent derives it from the **version the design targets**: the project's
   **next** version for forthcoming design (the normal case), or the **current**
   version when documenting as-built state. The version token is read from the
   project's own source of truth (e.g. `workspace.yaml` `template_version`,
   release manifests), not invented.
3. **Named groupings are allowed** as an explicit author choice — a slug such as
   `phase-2`, `Q1`, or a milestone name. Within one product, pick one grouping
   dimension rather than mixing release and phase folders arbitrarily; a named
   grouping is a deliberate deviation from the version default, not an
   agent-picked one.
4. **`README.md` + `design.md` structure is unchanged** at every level; only the
   meaning and formatting of the middle segment are refined here.

## Example

```
sources/system-design/
  context-circuit/            # product
    v0.6.1/                   # grouping = version (default, 3-number going forward)
      design-layout-grouping/ # scope
    v0.5/  v0.6/              # legacy 2-number groupings, not renamed
  acme-platform/              # a different product/project
    phase-2/                  # grouping = named phase (explicit team choice)
      checkout-redesign/      # scope
```

## Relationship to the owner

This refines the layout convention owned by the `cc-system-design` skill
(v0.6 `system-design-authoring` scope). It adds no lifecycle, status, gate, or
runtime, and no new owner. It reaches the skill through the normal path: the
convention is proposed as a change to the skill's authoring guidance, a human
accepts it, and the skill text is updated then. Until then this scope is the
source of record for the refined rule.

## Constraints and edge cases

- **No retroactive rename.** Legacy `v0.5`/`v0.6` stay; the 3-number rule binds
  new groupings only. A tree may therefore hold both depths during transition —
  acceptable and expected.
- **Single-project workspace.** The `<product-or-project>` level may collapse, so
  the path can be `sources/system-design/<grouping>/<scope>/`.
- **Ambiguous "current or next."** If the targeted version is genuinely unclear,
  the agent surfaces the choice rather than guessing a number.
