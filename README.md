# Kao Delivery Workspace

Kao Delivery Workspace 0.1.0 is a static, inspectable wrapper for coordinating
AI-assisted delivery across isolated Git worktrees. It includes initialization,
optional numbered planning and activity hooks, read-only task recommendation,
single- and contract-first multi-repository execution, independent verification,
bounded repair, review handoff, human-invoked closeout, context curation, and
thin Codex and Claude Code adapters.

## Quick start

Initialize the wrapper with `$initialize-workspace` in Codex or
`/initialize-workspace` in Claude Code. After repository registration, validate
it with the bundled dependency-free command:

```bash
node .agents/bin/kao.mjs validate --check-paths --check-documents
node .agents/bin/kao.mjs run-task \
  --request "Add a reset button to the counter" \
  --repository frontend \
  --acceptance "A Reset button sets the displayed count to zero" \
  --scope "src/App.tsx" \
  --test-scope "src/App.test.tsx"
```

The command prints paths to a task brief, runtime manifest, worker input,
verifier input, and isolated worktree under `.runtime/`. It never invokes an
agent or changes an activity tool. See [run-task](docs/run-task.md) for the
complete Codex and Claude Code workflow.

## Initialize a wrapper

Invoke `$initialize-workspace` in Codex or `/initialize-workspace` in Claude Code
after creating a wrapper and placing its product repositories at their intended
paths. The canonical skill inspects existing material, resolves only material
unknowns, updates reviewable wrapper files, and uses deterministic validation for
repository modes and exact ignore entries. See [initialization](docs/initialization.md).

## Contract validation

`node .agents/bin/kao.mjs validate` validates `workspace.yaml`. Validate another supported
document by passing its schema explicitly:

```bash
node .agents/bin/kao.mjs validate --schema task-brief .runtime/tasks/<run-id>.json
```

Schemas live in `.agents/contracts/` and are the machine-readable contract authority.
The root TypeScript package, tests, and fixture source are development-only
maintainer tooling. They are excluded from the final distributable template.

## Runtime progression

Record lifecycle stages around fresh worker and verifier sessions:

```bash
node .agents/bin/kao.mjs record-result --run-id <run-id> --repository frontend --stage worker-started
node .agents/bin/kao.mjs record-result --run-id <run-id> --repository frontend --stage worker-result
node .agents/bin/kao.mjs record-result --run-id <run-id> --repository frontend --stage verifier-result
```

The recorder is idempotent. It validates result contracts, task/run/repository
identity, allowed scope, branch and commit state, acceptance coverage, and clean
worktrees before atomically advancing the manifest.

For cross-repository work, pass a validated `run-task-request` JSON document.
Each repository receives a separate worktree and scoped inputs. Dependency
workers remain locked until the shared-contract repository passes independent
verification; see [run-task](docs/run-task.md) and the
[example request](docs/examples/contract-first-run.json).

When an activity provider is configured, semantic actions are prepared and
recorded without a provider SDK. Required, optional, and manual policies control
whether execution continues, and starting actions complete before worktree
creation. See [activity lifecycle hooks](docs/activity-lifecycle.md).

Failing verification can enter a bounded fresh-repair loop, while passing work
can produce a validated draft-PR handoff without external mutation. See
[verification, repair, and review preparation](docs/review-lifecycle.md).

After human-confirmed merge or deliberate abandonment, use `$finish-work` in
Codex or `/finish-work` in Claude Code. Closeout first records an append-only
contribution; an explicit second step removes only a clean, safely preserved
worktree. See [finish-work](docs/finish-work.md).

## Optional planning

Invoke `$create-plan` in Codex or `/create-plan` in Claude Code to turn a PRD or
substantial idea into a numbered draft under `context/plans/`. Plans use
machine-readable approval metadata and stable work IDs, but remain optional for
direct task execution. See [the numbered planning workflow](docs/planning.md).

Approved plan work can be published only after an explicit request and read-only duplicate discovery. See [duplicate-safe plan publication](docs/plan-publication.md).

Completed-work learnings can be classified and curated through an isolated
wrapper worktree with `$sync-context` or `/sync-context`. Team mode prepares a
wrapper pull-request handoff; repository-local details and future tasks are
routed without being silently written as canonical facts. See
[context synchronization](docs/context-sync.md).

## Recommend the next action

Invoke `$whats-next` in Codex or `/whats-next` in Claude Code, or run
`node .agents/bin/kao.mjs whats-next`. The selector reads approved plans and optional test-only activity
fixtures, applies transparent readiness and ranking rules, and returns one
recommendation without changing state. See [the read-only recommendation
workflow](docs/whats-next.md).

## Test expectations

`--scope` declares implementation paths. `--test-scope` declares tests the worker
must add or update and defaults the test policy to `required`. Other explicit
policies are `existing-coverage`, `verifier-only`, and `not-required`; use
`--test-rationale` when the reason is task-specific. Omitting test scope defaults
to verifier-only evidence rather than silently implying that tests are covered.

Root `package.json`, `tsconfig.json`, `node_modules/`, `fixtures/`, `scripts/`,
and `test/` are maintainer inputs and are excluded from the distributable
template. General users run the checked-in `.agents/bin/kao.mjs` bundle and do
not install wrapper dependencies.
