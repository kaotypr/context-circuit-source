# Requirements

## Requirements and acceptance criteria

- Define a compact Product Knowledge layout consisting of PROJECT.md, GLOSSARY.md, context/roles/README.md, one Markdown file per business role, domain README files, and domain workflow pages; validation must accept this layout without requiring optional role or navigation expansions.
- Keep PROJECT.md short enough to route readers to relevant roles and domains; a validator or fixture must demonstrate that an agent can select a domain without loading the complete knowledge tree.
- Require every role file to describe the role definition, primary outcomes, relevant domains, product surfaces, end-to-end role story, related workflows, role-specific behavior, limitations, owners, sources, review date, and known gaps.
- Require every workflow page to describe its outcome, actors, entry points, current flow, important role or condition variations, business rules, implementation ownership, owners, sources, review date, and known gaps.
- Prevent duplicated authority: role pages own the cross-domain user perspective, domain summaries own business-area routing, and workflow pages own exact current behavior and rules.
- Extend initialization so a workspace may establish a minimal reviewed baseline containing product purpose, major user roles, major domains, critical journeys or workflows, authoritative sources, and explicit unknowns without attempting a complete application inventory.
- Extend context gathering to resolve only the relevant product map, role, domain, workflow, and repository context for a request; missing or contradictory knowledge must produce source-backed runtime candidates or gaps and must not mutate canonical files.
- Produce a compact immutable task-context package that records the selected context paths and Git revision or digest so workers and verifiers can use the same business baseline without receiving the complete Product Knowledge tree.
- Extend plans to reference relevant role, domain, and workflow paths, declare product impact as none, documentation-correction, implementation-only, behavior-change, new-workflow, or retired-workflow, and summarize the proposed behavior change.
- Extend planless task normalization to perform the same Product Knowledge impact assessment when the request can change user-visible or business behavior.
- Require worker and verifier evidence to report whether Product Knowledge impact is absent, matches the declared change, is unexpectedly broader, or reveals incorrect or contradictory current context.
- Synchronize canonical Product Knowledge only through a source-backed, human-reviewed wrapper change after behavior is confirmed effective; merged but disabled or unreleased behavior must remain proposed or recorded in contribution evidence.
- Generate role- or domain-oriented onboarding reading packs from canonical Product Knowledge without making the generated view an independent source of truth.
- Preserve existing context files, planning, planless execution, host neutrality, ignored runtime evidence, and human gates; adoption of Product Knowledge must be incremental and backward-compatible.
- Provide deterministic validation and automated tests for required sections, valid relative references, missing sources or owners, context-package scope, plan impact declarations, contradiction preservation, and effective-state synchronization.
