# Kao Delivery Workspace agent instructions

While `PLAN.md` exists, read it completely as the development product authority.
Always read `WORKFLOW.md` and `workspace.yaml` before coordinating work. The final
0.1.0 completion change removes `PLAN.md` after its durable rules are embodied in
permanent documentation, contracts, and tests.
Canonical workflow behavior lives under `.agents/`; host adapters must remain thin.

- Product repositories own code and repository-local conventions.
- Plans and activity integrations are optional for an explicit task request.
- Never operate on a dirty base repository or discard unrecorded work.
- Use a separate branch and worktree for each repository worker.
- Workers may only modify their assigned worktree and must not mutate activity state.
- Verifiers are read-only and independent from workers.
- Do not merge, deploy, publish, or store credentials.
- Preserve `.runtime/` until a human invokes a safe closeout workflow.
