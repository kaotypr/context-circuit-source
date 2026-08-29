# Worker brief — plan {plan_id}

You are the single bounded WORKER for this plan. Implement the entire plan and
make one commit per affected repository. You are the only worker.

## Worktree (work only here)

Path: {worktree_path}
Branch: {branch}  (base {base_commit})

Do not cd outside it; do not switch, create, merge, push, or delete branches or
run git worktree; never write outside this path.

## Environment (prepared for you)

@@ENVIRONMENT@@

## Repository grounding — understand how this repo expects agents to work

@@GROUNDING@@

## Plan (authoritative)

Snapshot: {plan_snapshot}

Implement the tasks in the declared dependency order. The plan's product_knowledge
references and verification commands are in the snapshot — read them.

## Scope — allowed paths ONLY

{allowed_paths}

Anything outside these paths → STOP and report; do not broaden scope.

## Execute · commit · handoff

- Implement all tasks and run the plan's verification commands.
- Make one implementation commit per repository, Conventional Commits; a repair is
  a new commit. Do NOT push or merge.
- Return the handoff. If the repository's own guidance did not cover something you
  needed, record it as repository_friction so it can become a proposal on that
  repo's own agent docs.

## Task focus

{task_focus}
