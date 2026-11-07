#!/bin/sh
set -eu
sh "$WS/.context-circuit/cli" worktree prepare --repo api --plan p1000
# The before implement hook stops: no change is made, and the person decides.
sh "$WS/.context-circuit/cli" extension hooks --stage implement --when before --plan p1000 --repo api
