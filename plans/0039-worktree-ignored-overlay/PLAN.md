# Plan 0039 — Overlay the checkout's gitignored paths into every new worktree

**Intent:** i027-runnable-worktrees
**Repository:** context-circuit-source
**Tier:** Standard
**Status:** draft

## Objective

`git worktree add` gives a tree tracked files and nothing else. Right after it,
the runtime copies every gitignored file and folder that exists on disk in the
bound checkout into the new tree — `.env` included — copy-on-write where the
filesystem supports it, a real copy otherwise, **never** a symlink back into the
checkout. Untracked files that are *not* gitignored stay out. This happens for
execution worktrees (`cc_worktree_prepare`, `cc_base_prepare`) and for Explore
worktrees (`cc_pair_begin`), before anyone is attached, and a failed overlay
blocks setup rather than handing over a half-ready tree.

## Grounding (HEAD 7d72f468)

| Surface | Current state |
|---------|---------------|
| `.context-circuit/wrapper/runtime/engine.sh` | 4206 lines, POSIX sh (dash-safe), one flat `cc_main` dispatch at 4106-4173; helpers `cc_emit`, `cc_fail`, `cc_atomic_write`, `cc_digest`, `cc_root_abs` |
| `cc_worktree_prepare` (1229) | Plain execution worktree. Early-returns `worktree_reused: true` at 1237-1243; two `worktree add` arms at 1246 / 1249 |
| `cc_base_prepare` (1571) | Stacked/integration base. Own stale-reuse branch at 1591-1607; three `worktree add` arms at 1618 / 1623 / 1628; runtime-authored integration merge at 1636 |
| `cc_pair_begin` (1305) | Explore worktree at 1334, then writes the pointer. **No** grounding, hardening, or environment report anywhere in the Explore path |
| `cc_execution_begin` (2501-2535) | Calls one of the first two, writes `repositories/<repo>.yaml`, then `cc_discover_repo_grounding`. Every setup failure sets `status: blocked` + a `blocked_reason` (2517, 2529, 2534) |
| Copying today | Nothing. `reflink`, `clonefile`, `ls-files --others`, `overlay`, `node_modules` all return zero hits in `engine.sh` |
| `repositories.local.yaml` | **`context-circuit-source` is bound at `path: .`** — the bound checkout *is* the workspace root |
| `.gitignore` (1-3, 6-35) | Ignores `.runtime/`, `dist/`, `/repositories/`, `/repositories.local.yaml`, `/role-tiering.local.yaml`, `/member.local.yaml`, and more |
| Cleanup | `cc_worktree_drop` (1464) uses `git worktree remove --force`, which discards ignored files; `cc_runtime_cleanup` (1495) drops dead Explore trees |
| Acceptance | `test/acceptance.sh` fans out per area; `test/lib/fixture.sh` builds workspaces and repos. No worktree suite exists yet |

### The hazard this plan must not walk into

Because the checkout is the workspace root, an unfiltered "copy every ignored
path" would copy **`.runtime/`** — the directory the new worktree lives inside —
into itself, along with `/repositories/` and the three host-local
`*.local.yaml` files the intent's non-goals already exclude. The overlay
therefore carries one fixed, top-level-only exclusion set. It is a workspace
self-protection constant, **not** a per-repository copy list, so the
no-CC-side-profile non-goal holds.

`plan-allocate-id` fails `PLAN_PREFIX_COLLISION` on archived duplicate prefixes
0001-0016 (plans 0034 and 0038 recorded the same blocker). Ids **0039, 0040,
0041** are the next consecutive in-band ids for member `kao` (band 1-999,
highest existing 0038), each checked free against `cc_plan_numbers`.

## Decisions

**Three plans, not one.** This one is filesystem overlay: offline, provable
without any package manager, failing on copy semantics and path exclusion.
[0040](../0040-worktree-toolchain-provision/PLAN.md) is dependency provisioning:
it may use the network, owns a runtime-constant toolchain table, and fails for
entirely different reasons.
[0041](../0041-honest-ready-and-proof/PLAN.md) is truth-telling — contracts, the
brief, Product Knowledge, end-to-end proof — and is only verifiable once both
behaviors exist. Distinct execution and verification boundaries, so they stack.

- **One function, three call sites.** `cc_overlay_ignored SRC DST`, exposed as
  the `worktree-overlay` verb, rather than three inlined copies.
- **Copy mode is probed, not assumed.** `cp -c` (APFS clonefile) →
  `cp --reflink=auto` (GNU) → plain recursive copy, emitted as `copy_mode` so a
  reviewer can see which one ran.
- **Git owns the listing.** `ls-files --others --ignored --exclude-standard
  --directory`, never a hand-maintained list.
- **Reuse paths are skipped.** Re-overlaying a live worktree would overwrite a
  worker's local state (INV-PRESERVE-01).
- **Explore fails all the way closed.** On overlay failure `cc_pair_begin`
  removes the worktree and branch *it just created* and writes no pointer. It
  removes nothing that existed beforehand.
- **No new persisted field.** `pairing-session.yaml` keeps its five fields
  (INV-PAIR-01); `copy_mode` is stdout evidence only.

## Tasks

| id | title | depends on |
|----|-------|------------|
| OV-001 | Copy-mode probe + `cc_overlay_ignored` + the `worktree-overlay` verb | — |
| OV-002 | Overlay freshly created execution worktrees; block on failure | OV-001 |
| OV-003 | Overlay the Explore worktree before the pointer exists | OV-001 |
| OV-004 | Prove it offline, including the workspace-root hazard | OV-002, OV-003 |

Full changes, acceptance, and runnable done-checks are in `plan.yaml`.

## Risks

- **Recursion into `.runtime/`.** The headline hazard above. Mitigated by the
  exclusion constant and proved directly by OV-004's workspace-root case.
- **Overlay cost.** A large ignored tree (`node_modules/`) copied without
  reflink support is slow. The probe prefers clone/reflink; the fallback is
  correct but slower, and `copy_mode` makes it visible instead of mysterious.
- **Secrets spread.** `.env` now exists in every worktree. The human accepted
  this at Gate 1; `cc_worktree_drop`'s `--force` removal already clears the
  copies, and OV-004 asserts `git status --porcelain` stays empty so nothing is
  ever staged.
- **`git worktree add` in a repo bound at the workspace root** places the tree
  under an ignored parent. Existing behavior, unchanged, but it is why the
  exclusion set exists.
- **Explore rollback deletes.** OV-003 is the only place this plan removes
  anything; it is scoped to the branch and worktree created microseconds
  earlier in the same call.

## Verification

`sh test/worktree/test-overlay.sh`, plus the existing
`test/execution/test-execution.sh`, `test/run-stack/test-run-stack.sh`,
`test/pairing/test-pairing.sh`, `test/security/test-boundaries.sh`, and finally
`sh test/acceptance.sh`.
