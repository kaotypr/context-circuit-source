# Integrate Context Circuit through host-native route folders

Plan ID: 0006-host-native-routes
Intent: i020-host-conventions-preserved
Status: draft

## Original request and coverage

- What you asked: Context Circuit should plug into each host's own conventions
  — Claude `.claude/`, Codex `.codex/` plus `AGENTS.md`, Cursor `.cursor/` —
  as thin routes into Context Circuit's single-owned files, with the same
  integration in maintainer source and the blank template.
- Covered below: wire native folders (HOST-001), ship them through release
  assembly and template seed (HOST-002), prove they stay routes and
  fail-closed still holds (HOST-003).
- Still allowed: root `AGENTS.md` / `CLAUDE.md` / `CURSOR.md` as instruction
  adapters; skills owned at `.agents/skills/cc-*`; source-only maintainer
  extras (for example `cc-human-simulator`).

## Objective and desired behavior

After this plan, each supported host discovers worker, verifier, and planner
roles and standing rules through that host's native project tree. Host files
tell the host where the owner is; they do not become a second policy copy.
A newly assembled workspace receives the same host-native routes as the
maintained product set.

## Constraints and non-goals

- Route, do not copy — `.context-circuit/agents` and invariants stay owners.
- Use only documented host surfaces; no `.codex/rules/`.
- Do not embed or replace a host CLI or SDK.
- Do not add a second router, lifecycle, or host-specific authorization.
- Do not expand to additional hosts in this plan.
- Preserve maintainer source-checkout distinctions for root `AGENTS.md` and
  source-only Claude test agents.
- Do not deliver, publish, merge, or commit as part of this plan.

## Product Knowledge grounding

- **host-adapters** — today describes root-only entry; must reflect committed
  native folders after completion.
- **source-release-and-upgrade** — assembly and template seed must include the
  host-native product set and exclude maintainer-only files.

Grounding summary: at HEAD `8d638de`, `.codex/` is missing, `.claude/` and
`.cursor/` are partial, product `CLAUDE.md` still calls `.claude/` unshipped,
and release assembly does not stage any host folder. Intent detail under
`intent/i020-host-conventions-preserved/detail/` confirms the product set.

## Repositories and source evidence

- **context-circuit-source** — owns host folders, adapters, assembly, template,
  and tests. Evidence: `.codex/` absent; `.claude/agents/` has only
  `cc-human-simulator`; `.cursor/agents/` absent; `scripts/release-artifact.sh`
  stages only `.agents` and `.context-circuit`; `template/` has no host trees;
  product adapter contradicts approved intent on shipping `.claude/`.

## Plan decomposition

One plan with three ordered tasks. Wiring, shipping, and proof share one
worker lifecycle and one `test/acceptance.sh` verification surface — a move
without assembly fails AC-HOST-05; assembly without stubs fails manifest
checks; proof without both cannot close the intent.

## Tasks

### HOST-001 — Wire Claude, Codex, and Cursor native folders as thin routes

Create the product host-native set in the source checkout:

| Host | Agents | Rules / instructions | Skills |
| --- | --- | --- | --- |
| Claude | `.claude/agents/{worker,verifier,planner}.md` | `.claude/rules/` spawn + commit stubs | symlinks under `.claude/skills/` → `cc-*` |
| Codex | `.codex/agents/*.toml` | pointers in root `AGENTS.md` where needed | already at `.agents/skills/` |
| Cursor | `.cursor/agents/{worker,verifier,planner}.md` | existing + commit `.cursor/rules/*.mdc` | already at `.agents/skills/` |

Update shipped product adapters to describe committed host routes and remove
"host-local unshipped" language. Keep maintainer root pointer behavior.

**Done when:** HOST-VT-01 through HOST-VT-03 pass.

### HOST-002 — Ship host-native routes in assembly and template seed

Extend `release-artifact.sh` and `release-manifest.txt` to ship the product
host set; filter out `cc-human-simulator` and `cc-test-case`. Mirror the
product set under `template/`. Document upgrade notes. Update release tests.

**Done when:** HOST-VT-04 through HOST-VT-06 pass.

### HOST-003 — Prove routes stay pointers and host-blocked still holds

Add `test/host-native/test-host-routes.sh`, wire it into acceptance, update
host-adapters Product Knowledge, and run the full semantic suite.

**Done when:** HOST-VT-07 through HOST-VT-09 pass.

## Risks

- **Adapter contradiction** — product `CLAUDE.md` currently forbids shipping
  `.claude/`; must be updated in HOST-001 before assembly changes land.
- **Source-only leakage** — `cc-human-simulator` must not appear in artifacts;
  assembly needs explicit filtering.
- **Symlink portability** — Claude skill links must survive template clone on
  supported platforms; fall back to thinnest owner-pointing stub if needed.
- **Codex trust** — project `.codex/` applies only in trusted projects; that
  is host behavior, not a Context Circuit lifecycle change.
- **Cursor compatibility loader** — Cursor reads `.claude/` and `.codex/`
  too; Cursor-owned copies under `.cursor/agents/` must win on name collision.

## Intent acceptance criteria mapping

| Criterion | Task | Verification |
| --- | --- | --- |
| AC-HOST-01 native surfaces | HOST-001 | HOST-VT-01 |
| AC-HOST-02 routes not copies | HOST-001, HOST-003 | HOST-VT-02, HOST-VT-07 |
| AC-HOST-03 child mapping | HOST-001 | HOST-VT-01, HOST-VT-02 |
| AC-HOST-04 host-blocked | HOST-003 | HOST-VT-08 |
| AC-HOST-05 source/template align | HOST-002 | HOST-VT-05, HOST-VT-06 |
| AC-HOST-06 bounded evidence | HOST-003 | HOST-VT-07 |

## Verification summary

- Per-task checks: HOST-VT-01 through HOST-VT-08
- Full suite: `sh test/acceptance.sh` (HOST-VT-09)
