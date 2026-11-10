#!/bin/sh
set -eu
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
# The agent executes the plan, and the after-implement hook returns the
# override. Two rounds, one fix between them, and the CLI refuses a fix after
# the last.
cli="$WS/.context-circuit/cli"
sh "$cli" worktree prepare --repo api --plan p1000 >/dev/null
wt=$WS/.worktrees/p1000/api
cp "$here/files/export.go" "$wt/internal/invoice/export.go"
cp "$here/files/export_test.go" "$wt/internal/invoice/export_test.go"
git -C "$wt" add -A
git -C "$wt" commit -qm 'feat(invoice): export invoices as CSV'
hooks=$(sh "$cli" --json extension hooks --stage implement --when after --plan p1000 --repo api)
case $hooks in
*'"repair-loop"'*) ;;
*) echo "after implement returned no repair-loop override: $hooks" >&2; exit 1 ;;
esac
sh "$cli" --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 1 --path "$wt" >/dev/null
sh "$cli" --json agent dispatch --host claude-code --role worker --plan p1000 --task 'Findings: no trailing count row.' --override repair-loop --round 1 --path "$wt" >/dev/null
sh "$cli" --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 2 --path "$wt" >/dev/null
if sh "$cli" --json agent dispatch --host claude-code --role worker --plan p1000 --task 'Findings: no trailing count row.' --override repair-loop --round 2 --path "$wt" >/dev/null 2>&1; then
  echo "a fix after the last round was accepted" >&2
  exit 1
fi
