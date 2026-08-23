# Cc template-runtime laboratory

This source-only test area checks how the assembled `context-circuit-template`
works as a universal project workspace. It is deliberately separate from the
source repository's own ignored `.runtime/`, which stores execution state for
plans about the Cc source project.

Each scenario should exercise the built template through realistic project
prompts and small fixture repositories. A scenario may define:

- the project prompt or conversation;
- one or more fixture repositories and their starting revisions;
- starting Product Knowledge and source evidence;
- expected plan details and repository mappings;
- approval, execution, verification, repair, and failure-limit outcomes;
- expected human-facing results and preserved runtime evidence.

The harness must initialize an isolated temporary workspace from the assembled
template. It must not import the source repository's Product Knowledge, plans,
`.runtime/`, credentials, or implementation state. Generated workspaces,
worktrees, commits, and runtime records are disposable test artifacts.

The minimum scenario set covers context gathering, detailed plan creation,
`review plan`, approval, multi-repository execution, independent verification,
worker repair, the three-failure limit, explicit human completion, archive or
restore without plan-status validation, optional workspace-root repository
binding, clone and `git init` setup under the ignored `repositories/`
directory, personal or team anchor branches, and the separate pull-request
action using execution branches as sources and anchor branches as default
targets.

This directory belongs to `context-circuit-source` and is excluded from
`context-circuit-template` assembly.
