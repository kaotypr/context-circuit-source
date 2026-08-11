# Kao Delivery Workspace

Kao Delivery Workspace is a static, inspectable wrapper for coordinating
AI-assisted delivery across isolated Git worktrees. This repository currently
contains the version 0.1.0 Phase 0 foundation: validated contracts, planless task
preparation, a React fixture, and thin Codex and Claude Code adapters.

## Quick start

```bash
npm install
npm run check
```

Create the ignored fixture repository, then prepare a planless run:

```bash
npm run fixture:create
npm run run-task -- \
  --request "Add a reset button to the counter" \
  --repository frontend \
  --acceptance "A Reset button sets the displayed count to zero" \
  --scope "src/App.tsx" \
  --test-scope "src/App.test.tsx"
```

The command prints paths to a task brief, runtime manifest, worker input,
verifier input, and isolated worktree under `.runtime/`. It never invokes an
agent or changes an activity tool. See [the manual host proof](docs/phase-0-proof.md)
for Codex and Claude Code steps and [the recorded compatibility results](docs/phase-0-host-results.md)
for the completed first proof.

## Contract validation

`npm run validate` validates `workspace.yaml`. Validate another supported
document by passing its schema explicitly:

```bash
npm run validate -- --schema task-brief .runtime/tasks/<run-id>.json
```

Schemas live in `.agents/contracts/` and are the machine-readable contract authority.
The root TypeScript package, tests, and fixture source are development-only
maintainer tooling. They are excluded from the final distributable template.

## Runtime progression

Record lifecycle stages around fresh worker and verifier sessions:

```bash
npm run record-result -- --run-id <run-id> --repository frontend --stage worker-started
npm run record-result -- --run-id <run-id> --repository frontend --stage worker-result
npm run record-result -- --run-id <run-id> --repository frontend --stage verifier-result
```

The recorder is idempotent. It validates result contracts, task/run/repository
identity, allowed scope, branch and commit state, acceptance coverage, and clean
worktrees before atomically advancing the manifest.

## Test expectations

`--scope` declares implementation paths. `--test-scope` declares tests the worker
must add or update and defaults the test policy to `required`. Other explicit
policies are `existing-coverage`, `verifier-only`, and `not-required`; use
`--test-rationale` when the reason is task-specific. Omitting test scope defaults
to verifier-only evidence rather than silently implying that tests are covered.

## Development authority

`PLAN.md` is the product authority while version 0.1.0 is being built. It is a
development artifact, not part of the finished reusable template. Remove it in
the final completion change after its decisions and operating rules are embodied
in the permanent documentation, contracts, and tests.
