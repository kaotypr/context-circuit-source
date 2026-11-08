#!/bin/sh
set -eu
# The coordinator obtains the role's brief; the scripted reply stands in for
# what the role returned.
sh "$WS/.context-circuit/cli" --json agent dispatch --host claude-code --role x-parity-parity-verifier --plan p1000 --path "$WS/.worktrees/p1000/api" >/dev/null
