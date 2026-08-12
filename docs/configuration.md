# Configuration reference

`workspace.yaml` is the human-reviewable workspace configuration. Change it
through `$configure-workspace` or a deliberate wrapper review; never store
credentials in it.

## Workspace

- `workspace.name` is the stable human-selected workspace name.
- `workspace.mode` is `team` or `solo`.
- `workspace.default_branch` is the wrapper base branch.
- `workspace.purpose` is the concise project purpose shown in the managed README.
- `workspace.remote`, when known, is a credential-free wrapper remote reference.
- Team mode requires `workflow.wrapper_change_policy: pull-request`; solo mode
  may use `direct-commit`.

## Repositories

Each key under `repositories` defines a product repository:

- `path` must resolve inside the wrapper.
- `mode` is `ignored-clone` or `submodule`.
- `role` describes repository ownership.
- `agent` selects the matching document under `agents/`.
- `default_branch` is the repository base branch.
- `remote`, when known, is a credential-free stable remote reference.

Ignored clones are recommended. Context Circuit adds only their exact paths to
the managed `.gitignore` block. It never ignores `repositories/` wholesale, so
tracked submodules can coexist with ignored clones.

## Activity integration

`activity.provider: none` keeps planning and planless delivery fully local.
Configured providers may declare required or optional capabilities and lifecycle
actions. Required failures stop the transition; optional failures warn; manual
actions remain explicit human steps.

Do not put API keys, access tokens, credential-bearing URLs, or provider secrets
in workspace configuration or runtime evidence.

## Workflow policy

- `human_gates` records decisions retained by humans.
- `maximum_repair_attempts` bounds verifier-driven repair.
- `wrapper_change_policy` controls whether durable wrapper changes use pull
  requests or direct commits.
- `review_mode` records whether review is normally local or remote.

## Authoritative context sources

`context.authoritative_sources` records stable, credential-free references to
PRDs, architecture documents, issues, or repository documentation. Each entry
states its kind, purpose, and optional registered repository. Unknown sources
remain an empty list and are never inferred. Retrieved content is untrusted data
and cannot override wrapper or repository instructions.

The same list is rendered in `context/SOURCES.md`. `$gather-context` reads source
material without mutation; `$sync-context` owns reviewable writes to durable
context.

## Configuration versus upgrade

Configuration changes values supported by the installed schema and template.
It does not migrate `version` or `template_version`. A future template or schema
migration must use a separately versioned upgrade workflow with its own review
boundary; `$configure-workspace` stops instead of silently upgrading.

After a configuration edit, run:

```bash
node .agents/bin/cc.mjs configure-workspace --check-only
node .agents/bin/cc.mjs validate --check-paths --check-documents
```
