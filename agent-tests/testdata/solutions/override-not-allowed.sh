#!/bin/sh
set -eu
# Nothing is allowed, so the hook returns no override and the verifier is not
# dispatched; the CLI would refuse it.
cli="$WS/.context-circuit/cli"
sh "$cli" extension hooks --stage implement --when after --plan p1000 --repo api >/dev/null
if sh "$cli" --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --path "$WS/.worktrees/p1000/api" >/dev/null 2>&1; then
  echo "an unallowed override was dispatched" >&2
  exit 1
fi
