# Scenario prompt libraries (maintainer-only)

These are the natural-language conversations the human-simulator plays against the
product coordinator. They are **source-only maintainer material**: the release
manifest excludes `test/`, so nothing here ships in `context-circuit-template`.

Design:
`sources/context-circuit-v0.5-template-test-design/scenario-library.md`.

## Layout

```
scenarios/
├── README.md
├── 01-new-project-simple-idea/
│   └── case.yaml
├── 02-connect-existing-repo/        (not yet scaffolded)
├── ...
└── 09-verifier-unavailable-host-blocked/   (not yet scaffolded)
```

Each case is a directory so it can carry optional fixtures (seed repositories,
seed source files) alongside `case.yaml`. Only case 01 is scaffolded so far;
cases 02–09 follow the same `case.yaml` shape.

## Case-file shape

A `case.yaml` has two top-level blocks handed to different actors so the persona
never sees the mechanism:

- `human:` — everything the human-simulator receives (persona, turns, reactions,
  visible expectations). The `.claude/agents/cc-human-simulator.md` sub-agent is
  given only this block.
- `grader:` — everything the deterministic grader (`../human/grade.sh`) receives:
  the AC/invariant mapping, state post-conditions, transcript checks, the
  per-action access policy, and the soft efficiency budgets. The human-simulator
  never reads this block.

See `01-new-project-simple-idea/case.yaml` for the fully annotated example and
`../human/README.md` for how a run is prepared, driven, and graded.

## Boundaries

The same host-neutral library runs against each product host (`codex`,
`claude-code`, `cursor-agent`); only the launch/trace binding is per-host. Cases
are host-neutral except where a case explicitly asserts host behavior (case 09).
