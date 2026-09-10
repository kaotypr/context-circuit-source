# Plan 0040 — Make the worktree's toolchain match this commit

**Intent:** i027-runnable-worktrees
**Repository:** context-circuit-source
**Tier:** Standard
**Status:** draft
**Depends on:** [0039-worktree-ignored-overlay](../0039-worktree-ignored-overlay/PLAN.md)

## Objective

Once the overlay has cloned the checkout's ignored paths — install tree
included — the runtime decides whether that tree actually belongs to *this*
commit. Same lockfile content: keep it. Different: throw it away and run a
frozen install from the worktree's own lockfile. Detection stops being "a
lockfile exists"; every detected key gets a real row saying where its install
tree lives and how to install it frozen, or that it has no per-worktree tree at
all. `environment: ready` is emitted only after that succeeded, and a failure
blocks setup instead of attaching someone to a tree that cannot run.

## Grounding (HEAD 7d72f468)

| Surface | Current state |
|---------|---------------|
| `cc_harden_worktree` (1668-1684) | The entire current behavior: loop 15 filenames, emit `environment: ready` + `toolchain: <file>` the moment one exists, else `no-toolchain`. Its own comment (1666-1667) says "no network, no per-worktree install here (full provisioning is a later phase)" |
| Detection list (1672-1674) | `package-lock.json`, `pnpm-lock.yaml`, `yarn.lock`, `bun.lockb`, `bun.lock`, `go.mod`, `Cargo.lock`, `Gemfile.lock`, `requirements.txt`, `poetry.lock`, `composer.lock`, `pom.xml`, `build.gradle`, `build.gradle.kts` |
| `cc_discover_repo_grounding` (1706) | Takes hardening's `environment:` line, and falls back to `no-toolchain` when hardening fails — a failure would become a cheerful manifest |
| `cc_worker_brief_assemble` (1788-1789) | Prints "Ready: dependencies provisioned, commit hooks handled. Do NOT install or modify dependencies." — a claim nothing currently makes true |
| `grounding-manifest.yaml` | Already fixes `environment` to `[ready, no-toolchain]` and describes ready as "a toolchain was detected **and prepared**" |
| `cc_execution_begin` (2530-2534) | Already knows how to turn a grounding failure into `status: blocked` + `blocked_reason` |
| `cc_pair_begin` (1305) | Calls neither hardening nor grounding at all |
| `test/grounding/test-grounding.sh` | Asserts `environment: ready` from an empty `package-lock.json` placeholder, and calls `cc_harden_worktree` directly. Both assertions change meaning here |
| `engine.sh` | POSIX sh, dash-safe. No network call, no PATH probe, no per-ecosystem command exists anywhere today |

## Decisions

- **Match by lockfile content, not by sniffing the tree.** Compare `cc_digest`
  of the source checkout's on-disk lockfile with the worktree's, and write a
  one-line **stamp inside the install tree** recording the digest it was
  provisioned from — 0039's overlay then carries the stamp along. Stamp wins
  when present; the source-lockfile digest is the first-run fallback; unstamped
  beside a differing lockfile is always a mismatch. This is right for `cc-pair`
  too, where the checkout may be dirty: the on-disk lockfile is exactly what the
  on-disk install tree was built from.
- **The table is total.** All 15 keys get a row. Global-cache ecosystems
  (`go.mod`, `pom.xml`, `build.gradle`, `build.gradle.kts`) get an explicit
  *no per-worktree install tree — ready without installing* outcome instead of
  falling through. **The exact command strings are flagged for confirmation at
  plan review** — they ship to every template consumer, and several ecosystems
  have more than one defensible frozen form.
- **One indirection point.** `cc_toolchain_install` is the single place a table
  command runs, so the suite can observe whether it ran without executing a
  package manager. A testability seam in the runtime — it reads no file and
  takes no flag — not a configuration hook.
- **Missing tool fails closed.** If the table's command is not on PATH,
  provisioning fails. Degrading to `ready` would be precisely the workaround
  INV-GROUND-03 forbids.
- **`environment` keeps two values.** A third would invite callers to treat
  "prepared" and "detected" as two kinds of ready. The contract wants failure to
  *block*, not to be reported as a state.
- **Explore stays a pointer.** `cc_pair_begin` provisions inline and emits
  `environment` on stdout only — no manifest, no new pointer field
  (INV-PAIR-01).

## Tasks

| id | title | depends on |
|----|-------|------------|
| TC-001 | Toolchain table + the lockfile-match decision | — |
| TC-002 | `cc_provision_worktree`; report ready only after preparing | TC-001 |
| TC-003 | Block execution and Explore setup on provisioning failure | TC-002 |
| TC-004 | Prove keep / install / skip / fail-closed with no network | TC-003 |

Full changes, acceptance, and runnable done-checks are in `plan.yaml`.

## Risks

- **The command table is opinionated.** `yarn install --immutable` is Berry-only;
  `pip install -r requirements.txt` needs a virtualenv decision; Gradle and Maven
  have no honest frozen form. Mitigated by the explicit no-tree rows and by
  flagging the strings for review before execution.
- **Install cost and network.** A mismatching lockfile means a real install
  during setup — slower, and it can fail offline. That is the contract's
  intent (fail closed beats a tree that cannot run), but it makes setup
  fallible in a way it never was.
- **Reworking existing assertions.** `test/grounding/test-grounding.sh`
  currently *depends* on the dishonest meaning of ready. Those two assertions
  must be reworked, not deleted, or the change silently loses coverage.
- **Stamp trust.** A stamp is only as good as the tree beside it; a human who
  hand-edits `node_modules` after an install defeats it. Accepted: the same is
  true of every package manager's own integrity file.
- **Cascade with 0039.** If the overlay's exclusion set ever grows to cover an
  install tree, every worktree silently takes the install branch. TC-004's keep
  path is the assertion that catches that.

## Verification

`sh test/worktree/test-provision.sh`, the reworked
`sh test/grounding/test-grounding.sh`, plus
`test/execution/test-execution.sh`, `test/pairing/test-pairing.sh`, and finally
`sh test/acceptance.sh`.
