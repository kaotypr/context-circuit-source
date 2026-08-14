# Architecture

Canonical behavior lives under `scripts/lib/` and is bundled into `.agents/bin/cc.mjs`. Plans are YAML plus Markdown under `plans/`. Product Knowledge is source-cited Markdown under `context/`, with lightweight provenance in `context/sources.yaml`.

Registered repositories own code and local conventions. `run-task` uses Git worktrees for isolation and writes a human-readable prompt under ignored `.runtime/`; it does not create runtime lifecycle manifests or result contracts. Publication providers are optional output adapters and do not own workflow state.
