# Planning

Plans are numbered peer documents, not workflow contracts. A plan is located at:

```text
plans/<repository-key>-plans/<number>-<slug>/
```

`plan.yaml` contains the plan ID, number, title, track, status, source, repositories, plan dependencies, connections, and optional Product Knowledge references. The seven companion Markdown files explain the overview, requirements, acceptance criteria, solution, delivery, verification, and risks.

Task files use YAML frontmatter. Useful fields include `id`, `plan_id`, `title`, `status`, `description`, `repository`, `area`, `parent_task`, `dependencies`, `subtasks`, `connections`, `implementation_scope`, `test_scope`, `test_expectations`, `verification_commands`, `acceptance_criteria`, Product Knowledge references, and an optional external URL.

Fields may be omitted when they do not apply. Lightweight validation checks YAML parsing, core statuses, identity and reference resolution, registered repositories, local references, and dependency cycles. It does not demand placeholder evidence, hashes, digests, or empty ceremony.

New plans and tasks begin as `draft`. Approval and completion are explicit human actions. Editing approved content should prompt the human to consider reapproval, but Context Circuit does not manufacture digest or lifecycle evidence.
