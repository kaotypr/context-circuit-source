# Work breakdown

| Work ID | Title | Parent | Depends on | Repository | Area | External reference |
| --- | --- | --- | --- | --- | --- | --- |
| PKNOW-001 | Define Product Knowledge contracts, layout, templates, and validation | — | — | context-circuit | canonical-contracts | — |
| PKNOW-010 | Create the optional minimal Product Knowledge baseline during initialization | — | PKNOW-001 | context-circuit | initialization | — |
| PKNOW-020 | Discover source-backed knowledge candidates, contradictions, and gaps | — | PKNOW-001 | context-circuit | context-discovery | — |
| PKNOW-030 | Connect Product Knowledge references and impact declarations to plans | — | PKNOW-001, PKNOW-020 | context-circuit | planning | — |
| PKNOW-040 | Resolve bounded Product Knowledge packages for planless and planned tasks | — | PKNOW-020, PKNOW-030 | context-circuit | task-execution | — |
| PKNOW-050 | Carry Product Knowledge impact through worker, verifier, repair, and closeout evidence | — | PKNOW-040 | context-circuit | delivery-evidence | — |
| PKNOW-060 | Synchronize effective Product Knowledge through reviewed wrapper changes | — | PKNOW-020, PKNOW-050 | context-circuit | context-sync | — |
| PKNOW-070 | Generate revision-stamped role and domain onboarding packs | — | PKNOW-001, PKNOW-010 | context-circuit | onboarding | — |
| PKNOW-080 | Update host adapters, bundled commands, templates, and human guides | — | PKNOW-010, PKNOW-030, PKNOW-040, PKNOW-050, PKNOW-060, PKNOW-070 | context-circuit | host-adapters | — |
| PKNOW-090 | Prove bounded context, compatibility, lifecycle safety, and host parity end to end | — | PKNOW-080 | context-circuit | verification | — |

## Execution contracts

```json
{
  "contract_version": 2,
  "items": [
    {
      "work_id": "PKNOW-001",
      "repository": "context-circuit",
      "scope": [
        ".agents/contracts/product-knowledge-project.schema.json",
        ".agents/contracts/product-knowledge-role.schema.json",
        ".agents/contracts/product-knowledge-workflow.schema.json",
        ".agents/contracts/product-knowledge-domain.schema.json",
        ".agents/templates/product-knowledge/",
        "scripts/lib/product-knowledge.ts",
        "scripts/lib/validation.ts",
        "scripts/lib/types.ts",
        "fixtures/product-knowledge-small/",
        "fixtures/product-knowledge-large/",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/product-knowledge.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "The Product Knowledge layout (PROJECT.md, GLOSSARY.md, context/roles/README.md, one role file per business role, domain README files, and domain workflow pages) validates without requiring optional role or navigation expansions.",
        "Role page validation requires role definition, primary outcomes, relevant domains, product surfaces, end-to-end role story, related workflows, role-specific behavior, limitations, owners, sources, review date, and known gaps.",
        "Workflow page validation requires outcome, actors, entry points, current flow, role or condition variations, business rules, implementation ownership, owners, sources, review date, and known gaps.",
        "Validation rejects broken relative references and missing sources or owners, and enforces that role pages own the cross-domain perspective, domain summaries own business-area routing, and workflow pages own exact current behavior and rules.",
        "A small fixture and a large multi-domain fixture both validate, and an agent can select a domain from PROJECT.md without loading the complete knowledge tree."
      ]
    },
    {
      "work_id": "PKNOW-010",
      "repository": "context-circuit",
      "scope": [
        "scripts/lib/initialize-workspace.ts",
        "scripts/lib/product-knowledge.ts",
        ".agents/skills/initialize-workspace/SKILL.md",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/initialize-workspace.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "Initialization can create a minimal reviewed Product Knowledge baseline containing product purpose, major user roles, major domains, critical journeys or workflows, authoritative sources, and explicit unknowns without attempting a complete application inventory.",
        "Explicit unknowns are preserved rather than invented, and the generated baseline validates against the PKNOW-001 contracts.",
        "An existing initialized wrapper with no Product Knowledge remains valid; adoption is incremental and backward-compatible."
      ]
    },
    {
      "work_id": "PKNOW-020",
      "repository": "context-circuit",
      "scope": [
        ".agents/contracts/product-knowledge-candidate.schema.json",
        "scripts/lib/product-knowledge-discovery.ts",
        "scripts/lib/product-knowledge.ts",
        ".agents/skills/gather-context/SKILL.md",
        "scripts/lib/validation.ts",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/product-knowledge-discovery.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "Gather-context resolves only the relevant product map, role, domain, workflow, and repository context for a request rather than the whole tree.",
        "Missing or contradictory knowledge produces source-backed runtime candidates, gaps, or contradictions and never mutates canonical files.",
        "Discovery stores only credential-free candidates, gaps, contradictions, and context packages under ignored runtime state."
      ]
    },
    {
      "work_id": "PKNOW-030",
      "repository": "context-circuit",
      "scope": [
        ".agents/contracts/product-knowledge-impact.schema.json",
        ".agents/contracts/plan-draft-request.schema.json",
        ".agents/contracts/plan-index.schema.json",
        "scripts/lib/plans.ts",
        "scripts/lib/product-knowledge.ts",
        "scripts/lib/types.ts",
        ".agents/skills/create-plan/SKILL.md",
        "docs/planning.md",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/plans.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "Plans can reference relevant role, domain, and workflow paths and declare product impact as none, documentation-correction, implementation-only, behavior-change, new-workflow, or retired-workflow with a concise proposed behavior-change summary.",
        "Planning leaves canonical workflow and role pages unchanged; the proposed delta is recorded only in plan material.",
        "The Product Knowledge plan fields are additive with safe defaults so existing plans and the product-knowledge plan itself remain valid."
      ]
    },
    {
      "work_id": "PKNOW-040",
      "repository": "context-circuit",
      "scope": [
        ".agents/contracts/task-context-package.schema.json",
        ".agents/contracts/task-brief.schema.json",
        ".agents/contracts/run-task-request.schema.json",
        "scripts/lib/run-task.ts",
        "scripts/lib/product-knowledge.ts",
        "scripts/lib/validation.ts",
        "scripts/lib/types.ts",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/run-task.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "A large synthetic application resolves a compact, immutable task-context package containing only the selected product map, affected role, domain, workflow, and repository instructions.",
        "The package records the selected context paths and a Git revision or digest so workers and verifiers share the same business baseline without receiving the complete tree.",
        "Planless task normalization performs the same Product Knowledge impact assessment when a request can change user-visible or business behavior."
      ]
    },
    {
      "work_id": "PKNOW-050",
      "repository": "context-circuit",
      "scope": [
        ".agents/contracts/worker-result.schema.json",
        ".agents/contracts/verifier-result.schema.json",
        ".agents/contracts/closeout-record.schema.json",
        "scripts/lib/record-result.ts",
        "scripts/lib/finish-work.ts",
        "scripts/lib/run-task.ts",
        "scripts/lib/types.ts",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/record-result.test.ts", "test/finish-work.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "Worker and verifier evidence report whether Product Knowledge impact is absent, matches the declared change, is unexpectedly broader, or reveals incorrect or contradictory current context.",
        "Unexpected worker or verifier Product Knowledge impact blocks silent synchronization and is retained for human review.",
        "Closeout contributions record the Product Knowledge impact of the completed work."
      ]
    },
    {
      "work_id": "PKNOW-060",
      "repository": "context-circuit",
      "scope": [
        ".agents/contracts/context-sync-request.schema.json",
        ".agents/contracts/context-sync-record.schema.json",
        ".agents/contracts/workspace.schema.json",
        "scripts/lib/context-sync.ts",
        "scripts/lib/product-knowledge.ts",
        ".agents/skills/sync-context/SKILL.md",
        "scripts/lib/validation.ts",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/context-sync.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "Canonical Product Knowledge synchronizes only through a source-backed, human-reviewed wrapper change after a dedicated workspace-configured confirming role (distinct from the per-page owner and the merge-gate human) declares the behavior effective.",
        "Merged but disabled or unreleased behavior remains proposed or recorded in contribution evidence and is not written into current-behavior pages.",
        "Synchronization updates only the referenced current-behavior pages and preserves recorded contradictions for human review."
      ]
    },
    {
      "work_id": "PKNOW-070",
      "repository": "context-circuit",
      "scope": [
        ".agents/contracts/onboarding-pack.schema.json",
        "scripts/lib/product-knowledge-onboarding.ts",
        "scripts/onboarding-pack.ts",
        "scripts/cc.ts",
        "scripts/lib/validation.ts",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/product-knowledge-onboarding.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "Onboarding packs select the product map, one or more role files, relevant domain summaries and workflows, architecture context, and known gaps.",
        "Every pack is reproducible from a single Git revision, lists its included source paths, and is clearly marked as a generated view that is not an independent source of truth."
      ]
    },
    {
      "work_id": "PKNOW-080",
      "repository": "context-circuit",
      "scope": [
        ".codex/skills/",
        ".claude/commands/",
        "scripts/lib/validation.ts",
        "scripts/build-template.ts",
        "docs/product-knowledge.md",
        "docs/using-the-wrapper.md",
        "docs/command-reference.md",
        "README.md",
        "WORKFLOW.md",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/packaging.test.ts", "test/release-artifact.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "Every advertised Product Knowledge skill, command, and contract exists in source and in the built neutral template, and Codex and Claude adapters remain thin delegates to canonical .agents behavior.",
        "The deterministic bundle regenerates and the release archive is version-matched with a matching bundle digest.",
        "Human guides document Product Knowledge usage and migration for existing wrappers without mandating a complete inventory."
      ]
    },
    {
      "work_id": "PKNOW-090",
      "repository": "context-circuit",
      "scope": [
        "docs/product-knowledge-delivery.md",
        ".agents/bin/cc.mjs"
      ],
      "test_scope": ["test/product-knowledge-end-to-end.test.ts"],
      "test_policy": "required",
      "verification_commands": ["npm run typecheck", "npm test", "npm run validate -- --check-paths --check-documents", "npm run build", "git diff --check"],
      "acceptance_criteria": [
        "A large synthetic application resolves a bounded task context, and role stories coexist with workflow-owned rules under link validation.",
        "Initialization baseline, existing-wrapper compatibility, read-only discovery, plan and planless impact classification, blocked silent synchronization, effective-only canonical updates, and reproducible onboarding packs are all proven end to end.",
        "Codex and Claude host parity is verified and every prior safety property is preserved."
      ]
    }
  ]
}
```

Live task status does not belong in this plan. Add confirmed external references only after an explicit publication action.
