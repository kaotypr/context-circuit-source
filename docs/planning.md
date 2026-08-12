# Numbered planning workflow

Plans are optional. Use a plan when a PRD or substantial idea benefits from a
reviewed delivery model; a direct task request may still use `run-task` without
one. Each plan lives at `context/plans/<plan-id>/`, and its `README.md` contains
machine-readable YAML frontmatter plus links to sparse, four-digit Markdown
documents.

## Create a draft

Invoke `$create-plan` in Codex or `/create-plan` in Claude Code with a PRD,
document, issue, pull request, or explicit idea. Both adapters delegate to
`.agents/skills/create-plan/SKILL.md`; planning rules do not live in host files.

The host gathers relevant context and normalizes it into a temporary JSON file
matching `.agents/contracts/plan-draft-request.schema.json`. A minimal example is:

```json
{
  "contract_version": 1,
  "plan_id": "billing-v2",
  "title": "Billing v2",
  "source": { "kind": "prd", "reference": "docs/billing-v2.md" },
  "work_prefix": "BILLING",
  "summary": "Make retry behavior explicit and visible.",
  "affected_repositories": ["frontend"],
  "assumptions": [],
  "open_questions": ["Which retry states are customer-visible?"],
  "requirements": ["Retry state has explicit acceptance criteria."],
  "solution": ["Define the contract before implementing user-interface state."],
  "delivery": ["Approve the contract, then implement repository work."],
  "verification": ["Verify every acceptance criterion independently."],
  "risks": ["Unresolved states could create inconsistent behavior."],
  "work_items": [
    { "key": "retry-contract", "title": "Establish retry contract", "area": "architecture", "repository": "frontend", "scope": ["src/retry-contract.ts"], "test_scope": [], "test_policy": "verifier-only", "verification_commands": ["npm test"], "acceptance_criteria": ["The retry contract is explicit."] },
    { "key": "retry-ui", "title": "Display retry state", "area": "frontend", "repository": "frontend", "scope": ["src/App.tsx"], "test_scope": ["src/App.test.tsx"], "test_policy": "required", "verification_commands": ["npm test"], "acceptance_criteria": ["The retry state is visible."], "parent": "retry-contract", "depends_on": ["retry-contract"] }
  ]
}
```

Create and validate it with:

```bash
node .agents/bin/cc.mjs create-plan --input .runtime/plan-drafts/billing-v2.json
node .agents/bin/cc.mjs validate-plan context/plans/billing-v2
node .agents/bin/cc.mjs validate --check-documents
git diff --check
```

Creation is exclusive: an existing plan directory is never overwritten. The
first item receives `<PREFIX>-001`; later items receive sparse `-010`, `-020`,
and subsequent IDs. Those IDs must be preserved through revisions.
Each item also receives a canonical execution contract in the numbered work
breakdown. Repository, scope, test expectation, verification commands, and
acceptance criteria are approved plan material and are the authoritative input
to plan-linked execution.

## Approval and revision

New plans always start with `status: draft`, `plan_version: 1`, and null approval
fields. Human approval must explicitly cover scope, solution, delivery order,
risks, and acceptance criteria. Only then run the explicit transition to record
the actual UTC approval time, approver, and matching material fingerprints
before validating again. The human must provide the exact `approved_by` value;
do not infer it from the host session, Git identity, or email:

```bash
node .agents/bin/cc.mjs set-plan-state --plan context/plans/billing-v2 --approve-by delivery-lead
```

The validator recomputes the material digest so an edited numbered document
cannot silently retain stale approval.

A material revision increments `plan_version`, records `revision_reason`, returns
the plan to draft, clears approval fields and `approved_digest`, and refreshes
`material_digest`. Spelling, formatting, and link-only repairs do not revoke
approval after the human consciously classifies them as non-material. After
changing an approved plan materially, run
`node .agents/bin/cc.mjs set-plan-state --plan context/plans/billing-v2
--material-revision "<reason>"` to revoke approval and increment its version.
For a non-material repair, replace the final option with
`--non-material-repair`. Task publication and execution are separate explicit
workflows.

## Manual two-host proof

For Codex, invoke `$create-plan` with the same source and inspect the generated
files and validation output. For Claude Code, invoke `/create-plan` with that
source in a clean wrapper copy. Confirm both hosts create the same frontmatter
contract, document names, table columns, and work-ID allocation; confirm a
second creation attempt refuses to overwrite the plan. Approval remains a human
gate in both hosts.
