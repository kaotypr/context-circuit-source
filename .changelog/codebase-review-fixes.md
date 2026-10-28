---
product: cli
kind: fixed
authors: "@hallode"
pr: 169
---

Fixes from a review of the whole codebase. `check` and `record dependencies`
walk each plan's dependencies once, so plans that share ancestors no longer make
them exponentially slow. `agent dispatch --role worker` with neither `--plan`
nor `--task` is refused rather than launching a worker with nothing to
implement, and so is `--repo` without `--plan`. `init` stops before writing
anything when `README.md` is something it cannot read, such as a symlink or a
directory. `member language`, `member tone`, `member roles`, and `member band`
change only their own field, so the rest of the entry keeps its comments in
`members.yaml`; clearing a tone or the roles still rewrites the entry. `repo
remote`, `workspace remote`, and `knowledge remote` check the URL and the
default branch before writing either, so a refused change writes nothing.
`worktree prepare` skips and reports an environment entry that cannot be reused
by design, such as node_modules holding an `npm link`, instead of failing,
unless `--copy-mode required`; any other copy failure still stops it. `knowledge
sync` reports `unreachable`, not `diverged`, when origin no longer has the
default branch. `record create --kind intent` refuses `--intent`, `--repo`, and
`--depends-on` instead of dropping them.
