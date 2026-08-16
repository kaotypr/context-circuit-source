# Architecture

The target architecture is instruction- and filesystem-driven:

- AGENTS.md and WORKFLOW.md define normative agent behavior.
- Host integrations enter through the same workspace instructions and do not
  define a second command workflow.
- Product Knowledge is concise, source-cited Markdown under context/.
- Workspace identity is summarized in context/WORKSPACE.md. Project identity
  lives in context/PROJECT.md.
- Raw source material is a passive inbox under `sources/`; it is not part of
  ordinary session context.
- Idea Briefs and PRDs are product artifacts at a user/team-selected path
  under `sources/`, separate from accepted Product Knowledge. Do not prescribe
  subdirectories or filenames; record the exact path in provenance.
- Plans are human-reviewed YAML and Markdown under plans/.
- Runtime session records, leases, prompts, handoffs, and worktrees live under
  private .runtime/.
- Runtime record fields and lease ownership are defined in
  docs/runtime-contract.md.
- Registered repositories own code and repository-local conventions.
- Each writable plan execution has an exclusive Git worktree.
- `cc-run-plan` is the sole standard execution entry for approved plans. It
  directs child writers and verifiers in every workspace mode. The same
  writer-child and independent verifier-child topology applies regardless of
  workspace mode. `workspace.yaml` mode does not select an execution topology.
  Solo versus team remains identity; delivery policy IDs are
  `remote-review`, `local-target`, and `manual`.
- `cc-approve-plan`, `cc-finish-plan`, and `cc-cleanup-runtime` are the named
  plan-approval, status-change, and runtime-cleanup skills. They require
  current-session human confirmation and are discoverable, not mandatory
  ceremonies.

The workspace does not require a central database, activity provider, or
command runtime. Filesystem records are the inspectable runtime source of
truth.
