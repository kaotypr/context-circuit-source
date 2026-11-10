#!/bin/sh
set -eu
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
# The agent executes the plan: it prepares the worktree, implements without
# the header row the plan requires, and the after-implement hook returns the
# allowed override. The coordinator runs the loop: round 1 finds no header
# row, a worker fixes it on the worktree branch, and round 2 is clean. Each
# dispatch is accepted only under the allowed override.
cli="$WS/.context-circuit/cli"
sh "$cli" worktree prepare --repo api --plan p1000 >/dev/null
wt=$WS/.worktrees/p1000/api
sh "$cli" extension hooks --stage implement --when before --plan p1000 --repo api >/dev/null
git -C "$wt" apply "$here/../../fixtures/patches/csv-export-no-header.patch"
git -C "$wt" add -A
git -C "$wt" commit -qm 'feat(invoice): export invoices as CSV'
hooks=$(sh "$cli" --json extension hooks --stage implement --when after --plan p1000 --repo api)
case $hooks in
*'"repair-loop"'*) ;;
*) echo "after implement returned no repair-loop override: $hooks" >&2; exit 1 ;;
esac
sh "$cli" --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 1 --path "$wt" >/dev/null
sh "$cli" --json agent dispatch --host claude-code --role worker --plan p1000 --task 'Findings: the CSV output has no header row.' --override repair-loop --round 1 --path "$wt" >/dev/null
git -C "$wt" apply "$here/files/csv-header.patch"
git -C "$wt" add -A
git -C "$wt" commit -qm 'fix(invoice): write the CSV header row'
sh "$cli" --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 2 --path "$wt" >/dev/null
