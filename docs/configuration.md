# Configuration

`workspace.yaml` records the wrapper name, mode, default branch, registered repositories, and optional source-register location. Repository entries have an exact key, path, mode, role, agent, and default branch.

Keep repository paths credential-free and specific. Use `ignored-clone` for a separate local clone or `submodule` for a tracked submodule. Context Circuit refuses dirty bases and never rewrites unrelated work.

The wrapper does not configure activity lifecycle actions, repair limits, external status synchronization, or workflow state machines. Providers are optional publication adapters only.
