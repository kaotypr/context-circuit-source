# Use a Context Circuit wrapper

This guide is for a human starting a new wrapper from the downloadable archive.
The archive is deliberately neutral: it is not yet a Git repository, contains no
registered product repositories, and makes no assumption about technology or
repository names.

## Requirements

- Git
- Node.js 22 or newer
- Codex or Claude Code

Extract the release into the directory that should own cross-repository context,
then open that directory in the host. Do not manually initialize Git or create a
placeholder repository first.

## Initialize interactively

In Codex, invoke `$initialize-workspace`. In Claude Code, invoke
`/initialize-workspace`. Describe whether the project is new or existing and how
each repository should be sourced.

For a greenfield application:

```text
Create a team workspace named commerce-platform.
Create a new repository named app at repositories/app on main.
Use ignored-repository mode and no activity provider.
```

For existing repositories:

```text
Create a team workspace named commerce-platform.
Clone web from git@github.com:example/web.git into repositories/web.
Register the existing local repository at repositories/api as api.
Both use main. Keep activity integration disabled.
```

The agent will ask only for missing material choices, then present the exact
bootstrap action plan. A new wrapper and each new product repository require
explicit initial-commit authorization. Clone URLs must not contain credentials.
The [greenfield bootstrap example](examples/bootstrap-new.json) shows the
machine-readable request produced after that interview; humans normally do not
write it by hand.

After approval, the initializer can:

- configure `workspace.yaml`, workflow policy, context, and domain agents;
- initialize the wrapper Git repository;
- create a new repository with an empty base commit;
- clone an existing remote;
- register an existing local Git root without changing it;
- add a tracked submodule;
- add only exact ignored-repository paths to `.gitignore`;
- make the configured wrapper state its first commit; and
- validate contracts, documents, branches, paths, and repository cleanliness.

It never replaces an existing path, embeds credentials, pushes, opens a pull
request, mutates an activity system, merges, or deploys without separate human
authorization.

## Verify the result

The initializer runs these checks, which can also be repeated manually:

```bash
node .agents/bin/cc.mjs initialize-workspace --check-only
node .agents/bin/cc.mjs validate --check-paths --check-documents
git status --short
```

Missing repository remotes are warnings. Local execution can proceed, but draft
pull-request publication remains unavailable until the relevant remote exists.

## Start work

Use `$run-task` for a concrete request with known scope and acceptance criteria.
Use `$create-plan` first for broad greenfield work or decisions that need human
review. See [run-task](run-task.md), [planning](planning.md), and
[what's next](whats-next.md).

Runtime evidence is stored under ignored `.runtime/`. Preserve it until a human
invokes `$finish-work` after merge or deliberate abandonment.
