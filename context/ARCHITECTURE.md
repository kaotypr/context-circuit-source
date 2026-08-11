# Architecture

The wrapper owns workflow policy and cross-repository context. Registered
repositories own code and local conventions. Activity tools own live status,
approved plans own delivery intent, `.runtime/` holds private ignored evidence,
and append-only contributions hold durable outcomes.

Canonical contracts, agents, commands, and skills live under `.agents/`.
`.codex/` and `.claude/` are thin launch adapters. The distributable includes a
bundled Node 22 command at `.agents/bin/kao.mjs`; users install no npm packages.

Every run normalizes its source into a task brief with stable identity, scope,
acceptance evidence, test policy, repository ordering, authorization, and claim
limitations. Each repository receives its own branch, writable worktree, scoped
worker input, read-only verifier input, and result artifacts. A deterministic
recorder validates identity, Git history, scope, tests, and acceptance before
advancing per-repository and aggregate state.

For multi-repository work, the shared-contract owner runs first. Dependent
inputs remain locked until every dependency passes independent verification;
then both worker and verifier inputs receive verified contract state. Passing
work produces review handoffs only. Humans push, review, and merge normally.

Lifecycle actions are semantic and capability-based. The active host performs
authorized MCP, CLI, or manual actions and records only confirmed results with
stable idempotency keys. No credentials or provider configuration live here.
