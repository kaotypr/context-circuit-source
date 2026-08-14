# Legacy cleanup gate

The Node/JavaScript command layer is transitional compatibility material. This
implementation plan must preserve it. Removing it requires a separate,
human-reviewed cleanup plan after the replacement workflow passes.

## Required evidence before cleanup is proposed

All of the following must be true:

1. The filesystem acceptance suite passes for root entry, child delegation,
   concurrent plans, lease contention, exclusive worktrees, interruption and
   recovery, scope safety, contradiction handling, verification failure, and
   completion gating.
2. The existing typecheck, test suite, workspace validation, plan validation,
   and template build pass with the legacy layer still present.
3. An independent read-only verifier confirms that the acceptance evidence
   inspects runtime records and ownership invariants directly rather than
   merely checking legacy command output.
4. The root session records the evidence, remaining risks, active leases,
   unresolved worktrees, and any host-adapter limitations.
5. A human explicitly accepts the replacement result and requests a separate
   cleanup plan.

## Cleanup plan boundary

The cleanup plan must decide, explicitly and separately:

- whether to remove the user-facing legacy adapters first or all Node/JavaScript
  implementation files;
- which compatibility tests are retained as filesystem acceptance tests;
- how a rollback or temporary compatibility path works;
- what shipped template and host adapters must change;
- how removal is verified without deleting runtime evidence.

No file deletion, branch deletion, merge, publication, deployment, or status
change is authorized by this gate. It is evidence for a future human request,
not a cleanup action.
