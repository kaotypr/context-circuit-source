# Host capabilities

Codex, Claude Code, and Cursor Agent are host adapters around the same
Context Circuit contract. They expose the same capability names and preserve
the same intent, Product Knowledge, plan, runtime ownership, verification, and
human-gate semantics.

Hosts discover the same shared skill catalog: `cc-session-entry`,
`cc-initialize-workspace`, `cc-idea-brief`, `cc-create-prd`,
`cc-gather-context`, `cc-create-plan`, `cc-review-plan`, `cc-run-plan`,
`cc-run-stack`, `cc-whats-next`, `cc-configure-workspace`, `cc-approve-plan`,
`cc-finish-plan`, and `cc-cleanup-runtime`. The table below is a
safe-fallback example for optional configuration, not a per-host command
matrix. Hosts may improve discovery or credential storage. They must not
define a second workflow.

| Host | Capability discovery | Safe fallback |
| --- | --- | --- |
| Codex | `cc-configure-workspace` through the agent skill catalog | Core conversational workflow |
| Claude Code | `cc-configure-workspace` through the agent skill catalog | Core conversational workflow |
| Cursor Agent | `cc-configure-workspace` through the agent skill catalog | Core conversational workflow |

Host mappings may improve discovery or expose native secure credential storage,
but they must not create a second user-facing command workflow. A host that
cannot provide a requested integration reports `unavailable` and continues
with the filesystem-only path.

Host support never changes:

- the portable `workspace.yaml` configuration shape;
- the exclusive plan lease and worktree rules;
- independent verification;
- human approval for plan, delivery, publication, deployment, completion, or
  ambiguous ownership; or
- the prohibition on credentials and external activity records in workspace
  files.

## Child-session primitives

The filesystem packet is the coordination record. The host supplies the child
execution context. Name primitives behaviorally so a later host rename does
not require a new command layer.

| Host | Child-session primitive |
| --- | --- |
| Cursor Agent | Task / subagent tool |
| Claude Code | subagent / Task tool |
| Codex | native child-agent or equivalent |

Cursor's Task/subagent tool is a valid child-session primitive. Do not treat
a Cursor session as having no child primitive.

If the host cannot spawn a child, report the missing host primitive to the
human and ask how to proceed. A missing primitive is not a reason to skip
children.

## Deterministic YAML and schema validation

Context Circuit has an approved portable host capability for deterministic
YAML parsing and artifact-schema validation. The capability is host-provided;
the repository remains a pure filesystem product and does not add a required
Node, Python, Ruby, package-manager, or other command runtime.

The capability contract is:

| Field | Requirement |
| --- | --- |
| Capability name | `deterministic-yaml-schema-validation` |
| Provider | The active agent host or its native structured-data facility |
| Input | A Markdown path, the extracted initial front-matter block, the selected versioned schema, and whether the artifact is in an explicitly declared OKF bundle |
| YAML behavior | Parse the initial block as one deterministic YAML document; return structured data or a stable parse failure; do not approximate YAML with text matching |
| Schema behavior | Check required fields, types, enums, nested structures, `type`/legacy `kind` compatibility, and unknown-field preservation |
| Output | The independent `front_matter_extraction`, `yaml_parse`, `base_okf`, `context_circuit_profile`, `artifact_schema`, and `advisories` categories from `schemas/documents/validation-result-v1.yaml` |
| Ready rule | Output is ready only when every applicable hard category passes; advisories never change readiness |
| Runtime requirement | `repository_runtime_required: false`; `package_manager_required: false`; `command_runtime_required: false` |

Generation applies the capability before presenting or finalizing any
front-matter-bearing Markdown:

1. require `---` on line 1;
2. close at the first later standalone `---`;
3. parse only that initial block, leaving later Markdown horizontal rules in the
   body;
4. validate the selected artifact schema and, only where explicitly declared,
   base OKF and the Context Circuit profile;
5. repair hard failures and rerun the complete sequence.

If a host cannot provide this approved capability, it reports the capability
as unavailable and leaves the artifact not-ready. It must not silently fall
back to `grep`, a shell-only YAML approximation, or an undeclared repository
runtime. The filesystem model and all human gates remain unchanged.

The capability decision is recorded in the workspace runtime decision
`validation-capability-decision.yaml`; this document records only the portable
contract and its boundary.
