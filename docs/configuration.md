# Configuration

`workspace.yaml` records the wrapper name, solo/team mode, default branch,
registered repositories, and optional source-register location. An empty
`repositories: {}` registry is valid before initialization. Repository entries
have an exact key, path, mode, role, agent, and `default_branch`; initialization
recommends `development` only when it exists and lets the user choose another
branch.

Initialization records core identity only. Delivery behavior, merge policy,
publication, deployment, and external activity integrations are optional later
configuration and do not belong in the first-run questions.

Keep repository paths credential-free and specific. Use `ignored-clone` for a separate local clone or `submodule` for a tracked submodule. Context Circuit refuses dirty bases and never rewrites unrelated work.

The wrapper does not configure activity lifecycle actions, repair limits, external status synchronization, or workflow state machines. Providers are optional publication adapters only.
