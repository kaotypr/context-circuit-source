#!/bin/sh
set -eu
# Two rounds, one fix between them, and the CLI refuses a fix after the last.
cli="$WS/.context-circuit/cli"
wt=$WS/.worktrees/p1000/api
sh "$cli" --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 1 --path "$wt" >/dev/null
sh "$cli" --json agent dispatch --host claude-code --role worker --plan p1000 --task 'Findings: no trailing count row.' --override repair-loop --round 1 --path "$wt" >/dev/null
sh "$cli" --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 2 --path "$wt" >/dev/null
if sh "$cli" --json agent dispatch --host claude-code --role worker --plan p1000 --task 'Findings: no trailing count row.' --override repair-loop --round 2 --path "$wt" >/dev/null 2>&1; then
  echo "a fix after the last round was accepted" >&2
  exit 1
fi
