# Getting started

Context Circuit stores wrapper context in `context/`, plans in `plans/`, and registered repositories in `workspace.yaml`.

## 0. Start with an Idea Brief when the workspace is empty

If there are no useful source documents, plans, or registered repositories yet, use `$cc-idea-brief`. Discuss the problem, users, desired outcome, scope, constraints, and unknowns, then create the human-reviewed `context/IDEA-BRIEF.md`.

After the human confirms the brief, import it as source-cited Product Knowledge:

```sh
node .agents/bin/cc.mjs import-product-knowledge --source context/IDEA-BRIEF.md
```

Configure a repository/domain and create a numbered plan only after the idea has enough shape.

## 1. Configure

Use `$cc-configure-workspace` or create a setup request and run:

```sh
node .agents/bin/cc.mjs configure-workspace --request setup.json
```

The setup request is ordinary JSON/YAML-shaped data. It names the workspace, repositories, default branches, and any explicitly authorized Git initialization. It must not contain credentials.

## 2. Import Product Knowledge

```sh
node .agents/bin/cc.mjs import-product-knowledge --source docs/product-requirements.md
node .agents/bin/cc.mjs refresh-product-knowledge
```

The source register is `context/sources.yaml`. A refresh prints a proposal; it never silently overwrites canonical pages.

## 3. Create and approve a plan

Prepare a request with a title, source, repository, and tasks, then run:

```sh
node .agents/bin/cc.mjs create-plan --input plan-request.json
node .agents/bin/cc.mjs validate-plan plans/my-repo-plans/0010-example
node .agents/bin/cc.mjs set-plan-state --plan plans/my-repo-plans/0010-example --status approved
```

The generated layout is:

```text
plans/my-repo-plans/0010-example/
  plan.yaml
  overview.md
  requirements.md
  acceptance-criteria.md
  solution.md
  delivery.md
  verification.md
  risks.md
  tasks/README.md
  tasks/MY-0001.md
```

## 4. Choose and run a plan

```sh
node .agents/bin/cc.mjs whats-next
node .agents/bin/cc.mjs publish-plan --plan plans/my-repo-plans/0010-example --references publication-urls.json
node .agents/bin/cc.mjs run-task --plan plans/my-repo-plans/0010-example
```

Publication is optional and must happen before execution if chosen. `run-task` creates or reuses the plan/domain worktree and gives the agent one prompt containing the whole plan and all unfinished tasks. The agent continues through the tasks in dependency order and returns one plan-level handoff.

## 5. Review and finish manually

After the agent finishes, optionally review the whole plan once:

```sh
node .agents/bin/cc.mjs review-plan --plan plans/my-repo-plans/0010-example --worktree .runtime/worktrees/my-repo/0010-example
```

Then explicitly mark tasks and the plan done when the human is satisfied. Archiving is a separate explicit action.
