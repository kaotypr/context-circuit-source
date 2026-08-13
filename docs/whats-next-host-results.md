# `whats-next` host compatibility results

Proof completed on 2026-08-11 with Codex CLI 0.147.0 and Claude Code 2.1.220.
The proof used a disposable wrapper copy committed at a clean baseline. It
contained one approved, dependency-free frontend plan item (`PROOF-001`) and one
explicit urgent fake-activity candidate. The real workspace, activity systems,
and product repository history were not mutated.

## Results

| Scenario | Codex | Claude Code | Shared deterministic result |
| --- | --- | --- | --- |
| Approved plans only | `$w-whats-next` passed | `/w-whats-next` passed | `plan:host-proof:PROOF-001` selected; no alternatives |
| Explicit `activity.json` fixture | `$w-whats-next` passed | `/w-whats-next` passed | `activity:urgent-proof` selected; approved plan item returned as the sole alternative |

Both hosts reported the same readiness facts: ready lifecycle state, applicable
plan approval, completed or absent dependencies, sufficient scope and
acceptance criteria, and available frontend repository access. Both used the
shared TypeScript selector and returned `no_state_changed: true`. Before and
after checks showed a clean wrapper branch, the same single wrapper worktree,
no `.runtime/` directory, and no task claim, branch, commit, push, pull request,
or external write.

## Host observations

- Codex resolved the project skill through `.codex/skills/w-whats-next/SKILL.md`
  and delegated to the canonical `.agents` skill. It ran in a read-only sandbox.
  Its first discovery command used raw `sed` before reading the included RTK
  instruction; subsequent commands used RTK. This ordering difference made no
  mutation and does not change the shared recommendation contract.
- Claude Code resolved `.claude/commands/w-whats-next.md` and delegated to the same
  canonical skill. In non-interactive `dontAsk` mode it refused the deterministic
  Bash command, correctly stopping instead of reproducing ranking in the host.
  The disposable proof was rerun with permission bypass and write/edit tools
  denied. Normal interactive use can approve the read-only `npm run whats-next`
  command instead; the adapter requires no bypass behavior.
- Host narration and source-link formatting differ, but candidate IDs, ordering,
  evidence, blockers, risks, alternatives, and mutation guarantees remained
  host-neutral.

## Forward-test finding

Both hosts noticed that the generated plan body said “This plan is a draft”
after its frontmatter had been approved. The selector correctly treated the
validated metadata and matching digests as authoritative, but the prose was
misleading. The plan renderer now uses timeless approval-gate wording and states
that frontmatter status is authoritative; an automated regression test covers
the approved rendering.
