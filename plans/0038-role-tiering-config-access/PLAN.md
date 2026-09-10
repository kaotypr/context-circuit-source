# Plan 0038 — Engine-verb access to role-tiering config, with a committed fallback

**Intent:** i028-role-tiering-config-access
**Repository:** context-circuit-source
**Tier:** Standard
**Status:** draft

## Objective

An agent stops trying to find the role-tiering config and asks for it instead.
One new `engine.sh` verb prints the effective config — the workspace-root
`role-tiering.local.yaml` when it exists, the new committed
`.context-circuit/role-tiering.fallback.yaml` when it does not — and says which
source it used. `.context-circuit/docs/role-tiering.md` becomes the one full
statement of that rule; ten host-facing files carry a short pointer to it and
nothing more. The verb prints and never selects: host-group matching and adapter
defaults stay a coordinator/host decision (INV-RUNTIME-01), and `(model, effort)`
remains bounded host evidence that authorizes nothing (INV-HOST-01).

## Grounding (HEAD e66d9da8)

| Surface | Current state |
|---------|---------------|
| `.context-circuit/wrapper/runtime/engine.sh` | 4179 lines, POSIX sh, one flat `cc_main` dispatch (4087-4173). `cc_member_identity_read` (line 850) is the precedent: reads a gitignored `*.local.yaml` from the workspace root, emits one field, fails with a reason code |
| `role-tiering.local.yaml` | Present at the workspace root; `hosts:` matrix for `codex`, `claude-code`, `cursor-agent` |
| `.gitignore` line 18 | Ignores it, and lines 15-17 still say coordinators MUST Read the path directly — non-goals forbid touching this file |
| `.context-circuit/docs/role-tiering.md` | 172 lines; already owns shape, adapter defaults, `escalate_on_repair`, verifier independence, host application. `invariants.yaml` line 557 already names it owner of `role_tiering` |
| Adapter surfaces | `.context-circuit/wrapper/adapters/{AGENTS,CLAUDE,CURSOR}.md` are the shipped sources (`release-artifact.sh` 40-42 copies them to the artifact root); the three root copies are maintainer-facing. Root `CURSOR.md` is byte-identical to its source; root `AGENTS.md` and `CLAUDE.md` are not |
| Role files | `coordinator.md` (203-213) and `worker.md` (45-55) name the file directly; `verifier.md` and `planner.md` mention only `(model, effort)` |
| Release boundary | `release-artifact.sh` line 37 tars all of `.context-circuit`, so the new file *does* reach the artifact — but `manifest.yaml` lists only `.context-circuit/{wrapper,agents,docs}/`, so the new path is classified by neither `shipped` nor `never_ship`, and an upgrade would never install or replace it |
| Acceptance | `test/acceptance.sh` fans out per area. `test/contracts/test-contracts.sh` 229-247 pins exact strings in `role-tiering.md` and the three adapter sources; `test/bands/test-member-roster.sh` 100-121 is the house pattern for exercising a new verb and its failure code |

`plan-allocate-id` fails `PLAN_PREFIX_COLLISION` on archived duplicate numeric
prefixes 0001-0016 (plan 0034 recorded the same blocker). This plan uses the next
in-band id **0038** for member `kao` (plan band 1-999, highest existing 0037).

## Decisions

**One plan, five tasks.** Runtime verb, fallback data, owner doc, pointer stubs,
boundary declarations, and acceptance proof share one execution and one
verification boundary. Splitting them would publish a doc naming a verb that does
not exist, or a shipped file no boundary classifies. No partition is independently
verifiable or carries its own failure surface, so no stack.

- **Verb:** `role-tiering-read ROOT`, matching `member-identity-read`. Stdout is
  `source: local|fallback`, `path: <resolved>`, a `---` line, then the chosen file
  verbatim. Failure code `ROLE_TIERING_MISSING`.
- **Fallback content:** the workspace-root `hosts:` matrix verbatim, header
  comment swapped for a committed-fallback note. It is the only content the human
  supplied — **confirm the exact bytes at plan review**, because this file ships to
  every template consumer.
- **Pointer lands twice:** in the three adapter sources (so the product ships it)
  and the three root maintainer copies (so this checkout's agents see it).
- **Boundary made explicit** in `manifest.yaml`, `scripts/release-manifest.txt`,
  and `test/release/test-release.sh`, rather than relying on the tar.
- **`cc-execute`, `cc-trace`, `cc-pair` keep their current wording.**
  AC-POINTERS-PRESENT enumerates its files exactly, the skills are not among them,
  and `test-contracts.sh` pins those strings. Recorded as a residual risk.

## Tasks

### RT-001 — Add the `role-tiering-read` verb and the committed fallback file

**Paths:** `.context-circuit/wrapper/runtime/engine.sh`,
`.context-circuit/role-tiering.fallback.yaml`

`cc_role_tiering_read ROOT` beside `cc_member_identity_read`: resolve with
`cc_root_abs`, prefer the local override, else the fallback, else
`cc_fail ROLE_TIERING_MISSING`. Emit `source:`/`path:` via `cc_emit`, then `---`,
then the file. Register in `cc_main`. Print-only — no host lookup, no defaults, no
writes.

**Acceptance:** RT-AC-ENGINE-VERB, RT-AC-FALLBACK-FILE, RT-AC-PRINT-ONLY

**Verification:** `RT-VT-001` local source and body · `RT-VT-002` fallback source
in a temp root · `RT-VT-003` `ROLE_TIERING_MISSING` when neither exists ·
`RT-VT-004` the fallback is not gitignored · `RT-VT-005` no model id leaks into
`engine.sh`

### RT-002 — Make `role-tiering.md` the sole owner of the access rule

**Paths:** `.context-circuit/docs/role-tiering.md` · **Depends on:** RT-001

Obtaining the config becomes one instruction: invoke the verb. State that an agent
must not `Read`, `find`, or `grep` the file, and why — on some hosts a search
wrapper never emits a gitignored path, so finding nothing is not evidence of
absence. Document the fallback and the reported source. Keep every pinned sentence
and all existing semantics.

**Acceptance:** RT-AC-OWNER-DOC, RT-AC-NO-SEMANTIC-DRIFT

**Verification:** `RT-VT-006` verb + fallback + never-read wording ·
`RT-VT-007` pinned strings intact, no "tracer" · `RT-VT-008` defaults, ladder, and
independence still taught

### RT-003 — Point the adapter surfaces and role files at the owner doc

**Paths:** the three `.context-circuit/wrapper/adapters/` files, the three root
copies, and `.context-circuit/agents/{coordinator,worker,verifier,planner}.md`
**Depends on:** RT-002

A short `Owner: …` pointer in each, in the shape
`.claude/rules/role-tiering-spawn.md` already uses. Replace only the self-read
instruction; leave every host's spawn-application guidance alone. No fallback path,
precedence, or search-invisibility reason appears outside the owner doc.

**Acceptance:** RT-AC-POINTERS-PRESENT, RT-AC-POINTERS-THIN, RT-AC-HOST-APPLY-INTACT

**Verification:** `RT-VT-009` all ten pointers present · `RT-VT-010` no pointer
names the fallback path · `RT-VT-011` root `CURSOR.md` still matches its source and
the `@AGENTS.md` imports hold · `RT-VT-012` `sh test/contracts/test-contracts.sh`

### RT-004 — Declare the fallback file in the release and upgrade boundary

**Paths:** `.context-circuit/wrapper/manifest.yaml`,
`scripts/release-manifest.txt`, `test/release/test-release.sh`
**Depends on:** RT-001

Add the path to `release_boundary.shipped` and `upgrade_boundary.wrapper_only`,
add a `required` manifest line, and assert it in the staged artifact. The
`never_ship` entry for `role-tiering.local.yaml` stays.

**Acceptance:** RT-AC-BOUNDARY-DECLARED, RT-AC-FALLBACK-FILE

**Verification:** `RT-VT-013` declarations present, never-ship intact ·
`RT-VT-014` `sh test/release/test-release.sh`

### RT-005 — Prove the verb, the owner doc, and the pointers in the acceptance suite

**Paths:** `test/runtime/`, `test/contracts/test-contracts.sh`,
`test/acceptance.sh` · **Depends on:** RT-001…RT-004

New `test/runtime/test-role-tiering-read.sh` covering local, fallback, and
missing; plus a print-only assertion that no model id or host selection lives in
`engine.sh`. Extend the contracts suite for the owner rule and the ten pointers,
and prove `.gitignore`/RTK config are untouched. Wire the new suite into
`test/acceptance.sh`.

**Acceptance:** RT-AC-PROVEN, RT-AC-NO-RTK-CHANGE

**Verification:** `RT-VT-015` new suite · `RT-VT-016` wired into acceptance ·
`RT-VT-017` no `.gitignore`/RTK diff · `RT-VT-018` `sh test/acceptance.sh`

## Risks

- **Pinned strings.** `test-contracts.sh` 229-247 hard-pins phrasing in
  `role-tiering.md` and the three adapter sources — including `"workspace root"`,
  `"Applying a configured tier on Cursor"`, `"Do not default to inherit"`, and
  `"A missing file in an isolated working copy is not an absent config"`. A clean
  rewrite of the access paragraph is the most likely way to break the suite.
  RT-VT-007 and RT-VT-012 catch it early.
- **`.gitignore` contradiction survives.** Lines 15-17 still instruct agents to
  Read the path directly, which the new rule supersedes. Non-goals forbid editing
  `.gitignore`, so the contradiction is deliberate and unresolved. A follow-up
  intent should reconcile it.
- **Three skills still teach the old access path.** `cc-execute`, `cc-trace`, and
  `cc-pair` say "read `role-tiering.local.yaml` from the workspace root". The
  owner doc's claim to sole ownership is therefore partial in practice; a
  verifier may read this as incomplete. Deliberate, out of the enumerated
  criteria, and flagged for a follow-up.
- **The fallback ships someone's model ids.** `claude-opus-5`, `gpt-5.6-luna`,
  `composer-2.5-fast` would become the product's committed default matrix. If that
  is not intended, the human should supply neutral content at plan review;
  `agent-harness/fixtures/role-tiering.yaml` is a ready alternative.
- **`rg` honors `.gitignore`.** Any done-check that greps for the local override
  will silently find nothing — the exact failure mode this intent exists to fix.
  Every verification command here targets tracked files or invokes the verb.
- **Boundary drift.** Adding the file without the `manifest.yaml` entries would
  leave a shipped file that no upgrade ever replaces. RT-004 exists for that and
  nothing else.

## Expected Product Knowledge impact

Review `context/domains/host-adapters/README.md` (how every adapter obtains the
per-role config) and
`context/domains/source-release-and-upgrade/README.md` (a new committed file inside
the product boundary) on completion.
