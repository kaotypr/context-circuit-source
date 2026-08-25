# Scenario prompt libraries (maintainer-only)

Natural-language conversations the human-simulator plays against the product
coordinator. Source-only: the release manifest excludes `test/`, so nothing here
ships in `context-circuit-template`.

Design: `sources/system-design/template-harness/v0.5/scenario-library.md`.

## Cases

```
scenarios/
├── 01-new-project-simple-idea/        conv         (claude -p driver)
├── 02-connect-existing-repo/          conv         (claude -p driver)
├── 03-plan-review-and-approve/        conv         (claude -p driver)
├── 04-refuse-unapproved-execution/    conv         (claude -p driver)
├── 05-approve-and-execute/            full-exec    (in-session /cc-test-case)
├── 06-repair-then-complete/           full-exec    (in-session /cc-test-case)
├── 07-archive-plan/                   conv         (claude -p driver)
├── 07b-restore-archived-plan/         conv         (claude -p driver)
├── 08-delivery-boundary/              full-exec    (claude -p driver; no nesting needed)
└── 09-verifier-unavailable-host-blocked/  full-exec + fault  (in-session, host-can't-spawn)
```

Case 07 (design "archive-and-restore") is split into `07-archive-plan` +
`07b-restore-archived-plan` because an archive→restore round-trip is undetectable
from final state (active == never-archived); split, each half's outcome is
deterministic.

## Driver choice
- **Conversation-only** cases run via `../human/run-scenario.sh --live` (the
  `claude -p` driver), which captures a file-access trace so dimension C is a hard
  gate.
- **Full-execution** cases that need nested worker/verifier sub-agents run via the
  `/cc-test-case` skill (in-session Task driver); C degrades to warning-only there.

## Case-file shape
`case.yaml` has a `human:` block (persona, turns, reactions, visible_expectations —
all the human-simulator sees) and a `grader:` block (AC/invariant mapping,
post_conditions, transcript_checks, access_policy, budgets — the human-simulator
never sees it). See any case and `../human/README.md`.
