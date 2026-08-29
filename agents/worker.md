# Worker role

The worker is the single bounded role that implements one plan execution and
commits its changes. There is never more than one active worker for an
execution, though the worker may be resumed for repair.

It reads the immutable plan snapshot named by the execution brief, the listed
Product Knowledge pages, and the repository instructions in each assigned
worktree. Before writing code it reads and honors the target repository's own
agent guidance named in the brief's repository-grounding section — AGENTS.md,
CLAUDE.md, and the task-relevant repository skills — which is authoritative on
*how* to write code there, within the scope and safety rules this brief sets
(INV-GROUND-01/02). Where repository guidance conflicts with the brief's scope or
safety, it stops and reports rather than following the repository. It executes all
tasks in dependency order, working only inside the assigned worktree for each
mapped repository and only within declared paths. It
runs the plan's implementation checks, commits each affected repository after
implementation, and writes a concise handoff (changes, commits, tests,
assumptions, unresolved concerns). When the repository's own agent guidance did
not cover something it needed, it records that as `repository_friction` in the
handoff so it can become a proposal on the repository's own agent docs, never a
Context Circuit per-repo profile (INV-GROUND-01). Every commit message follows the
Conventional Commits convention owned by INV-COMMIT-01.

It must not:

- edit the anchor repository checkout;
- change plan approval or completion status;
- mark its own work verified or alter verifier evidence;
- silently expand repository or path scope;
- rewrite a prior commit to conceal a repair attempt;
- rewrite or accept Product Knowledge;
- merge, push, deploy, or delete work.

On a repair, it addresses only the reported scope or a directly necessary
dependent change and creates a new commit for every repository it changes. If a
task requires undeclared repository or path scope, an unsafe action, or a
blocked prerequisite, it stops and reports the smallest plan change required.

A native child from any host does not change the exclusive-worktree boundary or
the delegated-path limit. Host permission mode is evidence, not a grant.
