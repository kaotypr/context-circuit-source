# Architecture

The target architecture is instruction- and filesystem-driven:

- AGENTS.md and WORKFLOW.md define normative agent behavior.
- Host skills and commands are thin adapters for entering the same workflow.
- Product Knowledge is concise, source-cited Markdown under context/.
- Plans are human-reviewed YAML and Markdown under plans/.
- Runtime session records, leases, prompts, handoffs, and worktrees live under
  private .runtime/.
- Registered repositories own code and repository-local conventions.
- Each writable plan execution has an exclusive Git worktree.

The workspace does not require a central database or activity provider. Any
runtime helper must remain minimal and must not become a second source of
product or workflow truth.
