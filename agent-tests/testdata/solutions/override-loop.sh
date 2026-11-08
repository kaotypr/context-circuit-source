#!/bin/sh
set -eu
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
# The coordinator runs the loop: round 1 finds no header row, a worker fixes
# it on the worktree branch, and round 2 is clean. Each dispatch is accepted
# only under the allowed override.
cli="$WS/.context-circuit/cli"
wt=$WS/.worktrees/p1000/api
sh "$cli" extension hooks --stage implement --when after --plan p1000 --repo api >/dev/null
sh "$cli" --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 1 --path "$wt" >/dev/null
sh "$cli" --json agent dispatch --host claude-code --role worker --plan p1000 --task 'Findings: the CSV output has no header row.' --override repair-loop --round 1 --path "$wt" >/dev/null
git -C "$wt" apply "$here/files/csv-header.patch"
git -C "$wt" add -A
git -C "$wt" commit -qm 'fix(invoice): write the CSV header row'
sh "$cli" --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 2 --path "$wt" >/dev/null
