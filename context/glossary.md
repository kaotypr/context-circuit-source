# Glossary

Words this project uses with a specific meaning, and where each one lives in
code. Add a term the first time someone has to ask what it means. Keep rows
alphabetical, so two members adding terms from separate clones conflict on one
row rather than on one section.

Keep a meaning to a single sentence. A term needing more than that belongs in
its own note, linked from the catalog, with the row pointing there. Write code
anchors with the logical repository ID from `workspace.yaml`, never a local
checkout path, and name the identifier whenever the code calls the thing
something else. That mapping is the reason this file exists: a reader searching
the repositories for the business word finds nothing without it.

| Term | Meaning | Code |
| --- | --- | --- |
| Allocation band | A member's optional numeric block, so members allocate offline without colliding | `context-circuit-source@internal/workspace/records.go` `bandWidth`, `Store.allocate` |
| Catalog | The retrieval index a reader searches to choose a note, one unwrapped line per note | `context-circuit-source@internal/workspace/inspect.go` `catalogIssues`, `FindContext` |
| Check | The explicitly invoked diagnostic over records, bindings, and notes; never a gate | `context-circuit-source@internal/workspace/inspect.go` `Store.Check` |
| Coordinator | The main agent session that interprets, dispatches, integrates, and reports | `context-circuit-source@product/AGENTS.md.in` |
| Copy-on-write reuse | Cloning a checkout's ignored runtime entries into a new worktree by filesystem cloning | `context-circuit-source@internal/cow/`, `internal/workspace/reuse.go` |
| Dependency inputs | The tracked manifests and version files compared before dependencies are reused | `context-circuit-source@internal/workspace/reuse.go` `compareDependencyInputs` |
| Direct execution | Implementing a plan inside a bound checkout instead of a prepared worktree, only when a person asks for it | `context-circuit-source@product/AGENTS.md.in` |
| Dispatch specification | The resolved subagent invocation the CLI returns without launching anything | `context-circuit-source@internal/workspace/agents.go` `DispatchAgent` |
| Execution request | A person's separate decision to run presented plans, authorizing worktree preparation and implementation | `context-circuit-source@product/AGENTS.md.in` |
| Integration merge | The local merge assembling a dependent plan's base at a fan-in; implementation, not delivery | `context-circuit-source@internal/workspace/order.go` `OrderIntegration` |
| Intent | One approved statement of the outcome a change should produce, written before the code is read | `context-circuit-source@internal/workspace/records.go` `CreateRecord` |
| Ledger | The permanent record of every reserved intent and plan number | `context-circuit-source@internal/workspace/records.go` `Ledger` |
| Linear chain | An execution shape stacking every plan on the previous one, serial and merge-free | `context-circuit-source@internal/workspace/order.go` `Order.Chain` |
| Local binding | This machine's checkout path for a shared repository ID | `context-circuit-source@internal/workspace/workspace.go` `Bindings` |
| Logical repository ID | The shared name for a repository, independent of where anyone checked it out | `context-circuit-source@internal/workspace/workspace.go` `Config.Repositories` |
| Plan | How an approved intent maps to real code, plus what actually happened | `context-circuit-source@internal/workspace/records.go` `Record` |
| Reservation | A permanently allocated number that is never reused, including after archival | `context-circuit-source@internal/workspace/records.go` `Store.allocate` |
| Role tiering | The per-host model and effort preference for each subagent role, shared and optionally overridden per machine | `context-circuit-source@template/role-tiering.yaml`, `internal/workspace/agents.go` |
| Seed | The blank workspace embedded in the CLI for initialization and export | `context-circuit-source@assets.go`, `template/` |
| Start reference | Where a plan's branch begins in one repository, with any merges it needs first | `context-circuit-source@internal/workspace/order.go` `OrderStart` |
| Version store | The directory holding side-by-side CLI installs, one per version and platform | `context-circuit-source@product/skills/cc-cli/scripts/install.sh` |
| Wave | One layer of plans that may run concurrently, derived from dependencies and completion | `context-circuit-source@internal/workspace/order.go` `OrderWave` |

A word meaning different things in different repositories keeps one row per
repository, so the collision stays visible instead of resolving to whichever
row was written first.

| Term | Repository | Meaning |
| --- | --- | --- |

No term currently means different things in different repositories; this
workspace binds one.

Leave out general English, generic technology vocabulary, and the coordination
words this workspace product defines for itself.

_Reviewed: 2026-09-15._
