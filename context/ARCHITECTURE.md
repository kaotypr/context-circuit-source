# Architecture

The target architecture is instruction- and filesystem-driven:

- AGENTS.md and WORKFLOW.md define normative agent behavior.
- Host integrations enter through the same workspace instructions and do not
  define a second command workflow.
- Product Knowledge is concise, source-cited Markdown under context/.
- Raw source material is a passive inbox under `sources/`; it is not part of
  ordinary session context.
- Idea Briefs and PRDs are product artifacts under `contributions/`, separate
  from accepted Product Knowledge.
- Plans are human-reviewed YAML and Markdown under plans/.
- Runtime session records, leases, prompts, handoffs, and worktrees live under
  private .runtime/.
- Runtime record fields and lease ownership are defined in
  docs/runtime-contract.md.
- Registered repositories own code and repository-local conventions.
- Each writable plan execution has an exclusive Git worktree.
- `cc-run-plan` is the sole standard execution entry for approved plans. It
  directs a writer child and an independent verifier child in every workspace
  mode. `workspace.yaml` mode does not select an execution topology. Solo
  versus team remains identity; `solo-local` and `team-review` remain
  delivery policies.
- `cc-approve-plan`, `cc-finish-plan`, and `cc-cleanup-runtime` are the named
  plan-approval, status-change, and runtime-cleanup skills. They require
  current-session human confirmation and are discoverable, not mandatory
  ceremonies.

The workspace does not require a central database, activity provider, or
command runtime. Filesystem records are the inspectable runtime source of
truth.
