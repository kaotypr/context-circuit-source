# v1.0.0 scenario prompt library (maintainer-only)

Natural-language conversations the human-simulator plays against the product
coordinator. Source-only: the release manifest excludes `test/`, so nothing here
ships in `context-circuit-template`.

The cases below exercise the shipped v1.0.0 wrapper contract. They are not a
second lifecycle or a source-design archive.

## Cases

```
scenarios/
├── 01-new-project-simple-idea/        conv         (claude -p driver)
├── 02-connect-existing-repo/          conv         (claude -p driver)
├── 03-plan-review-and-approve/        conv         (draft intent → Gate 1 → derived plan)
├── 04-refuse-unapproved-execution/    conv         (draft intent, fail closed)
├── 05-approve-and-execute/            full-exec    (Standard worker + verifier)
├── 06-repair-then-complete/           full-exec    (Critical explicit completion)
├── 07-archive-plan/                   conv         (claude -p driver)
├── 07b-restore-archived-plan/         conv         (claude -p driver)
├── 08-delivery-boundary/              full-exec    (claude -p driver; no nesting needed)
├── 09-verifier-unavailable-host-blocked/  full-exec + fault  (in-session, host-can't-spawn)
├── 10-run-approved-stack/              full-exec    (single-repo run-stack)
├── 11-repo-grounding/                  full-exec    (repository guidance)
├── 12-run-multi-repo-stack/            full-exec    (multi-repo run-stack)
├── 13-execution-latency/               full-exec    (host evidence, hidden from user)
├── 14-codex-direct-collaboration/      full-exec    (Explore/direct collaboration)
├── 15-intent-adversary-revise-and-approve/ conv     (feasibility question from the real code)
├── 16-scope-expansion-regate/          conv         (out-of-scope reach → feasibility question)
├── 17-stale-candidate-requires-reverification/ full-exec (candidate staleness)
├── 18-standard-delivery-inferred-completion/ full-exec (Standard inference)
├── 19-knowledge-debt-blocks-next-plan/ full-exec    (closed knowledge loop)
├── 20-tier-escalation-fails-upward/    conv         (assurance floor)
└── 21-change-set-one-verification/     full-exec    (one integrated candidate)
```

The v1 flow represented by the library is:

```text
ordinary request
  -> intent proposal
  -> explicit Gate 1 approval
  -> post-approval tracing + feasibility check over the real code
  -> automatic schema-3 plan derivation
  -> one worker + independent verifier (Standard/Critical)
  -> candidate-bound human acceptance
  -> explicit Gate 2 delivery
  -> Standard inferred completion / Critical explicit completion
```

Cases 01 and 02 cover safe orientation and repository binding. Cases 03–04 and
15 cover the real intent gate, including refusal and a post-approval feasibility
question resolved before the plan is written. Case 16 covers an out-of-scope reach
surfaced as a feasibility question rather than an automated re-gate. Cases 05–06 and 17
cover Standard/Critical assurance, repair, explicit completion, and stale
candidates. Cases 07 and 07b cover status-preserving organization. Cases 08,
18, and 21 cover separate delivery, inferred Standard completion, and one
change-set candidate. Case 09 covers host-blocked verification. Cases 10–13
cover dependency ordering, repository grounding, cross-repository stacking, and
bounded host evidence. Case 19 covers reconciliation debt, and case 20 covers
the fail-upward assurance floor. Case 14 covers the Explore tier: one
human-supervised worker, no plan, no verifier, and no completion record.

The deterministic runtime suites remain the authoritative coverage for the exact
candidate digests, completion predicates, reconciliation-debt mechanics,
change-set merge records, and tier classifier. The conversational cases now also
exercise how those contracts are explained and invoked at the product boundary;
they do not duplicate every engine-level mutation. Legacy anchor_branch
readability/migration remains intentionally outside this scenario library.

Case 10 exercises the single-repository run-stack: ten intent-authorized, inter-dependent plans of a
small task-tracker CLI, built in one request. The runtime runs them in dependency
order, stacks each dependent on its predecessor, and builds runtime-authored
integration bases for the fan-ins (0005, 0009, 0010); the grader proves each plan
was independently verified and that the fan-in bases were built on their
predecessors' verified commits (based_on / built_on), with nothing marked done.

Case 11 exercises repository grounding: the connected repo ships its own
`AGENTS.md` stating a convention the plan never mentions (every file starts with
`// @grounded`). The grader proves the grounding manifest was recorded and that the
committed file carries that header — the observable proof that the worker
discovered, read, and honored the repository's own agent guidance.

Case 12 exercises run-stack across TWO repositories: seven intent-authorized plans
spanning a backend and a web app. The grader proves cross-repo concurrency, that a
cross-repo dependency is an ordering gate with no git base (`based_on_absent`), and
that same-repo dependencies still stack/integrate (`based_on` + `built_on`), with
all seven verified and nothing marked done.

Case 07 (design "archive-and-restore") is split into `07-archive-plan` +
`07b-restore-archived-plan` because an archive→restore round-trip is undetectable
from final state (active == never-archived); split, each half's outcome is
deterministic.

## Role tiering and driver choice

Role tiering is defined once in `../fixtures/role-tiering.yaml`, grouped under
`hosts.codex`, `hosts.claude-code`, and `hosts.cursor-agent`. The runner copies
the complete fixture into each generated workspace as the ignored
`role-tiering.local.yaml`; cases do not carry model or effort settings.

- **Conversation-only** cases run via `../human/run-scenario.sh --live` (the
  `claude -p` driver), which captures a file-access trace so dimension C is a hard
  gate.
- **Full-execution** cases that need nested worker/verifier sub-agents run via the
  `/cc-test-case` skill (in-session Task driver); C degrades to warning-only there.
- Built-in live bindings cover Claude Code, Codex, and Cursor Agent. The complete
  matrix is launched with `../human/run-matrix.sh`; it runs one isolated case per
  host lane at a time.

## Case-file shape
`case.yaml` has a `human:` block (persona, turns, reactions, visible_expectations —
all the human-simulator sees) and a `grader:` block (AC/invariant mapping,
post_conditions, transcript_checks, access_policy, budgets — the human-simulator
never sees it). See any case and `../human/README.md`.
