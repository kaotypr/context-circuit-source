#!/bin/sh
set -eu
sh "$WS/.context-circuit/cli" member add --id dana --name "Dana Reyes" --band 2 --role "backend developer"
sh "$WS/.context-circuit/cli" member use --id dana
sh "$WS/.context-circuit/cli" workspace connect --base main
sh "$WS/.context-circuit/cli" repo connect --id api --path "$WS/repositories/api" --base main
sh "$WS/.context-circuit/cli" agent setup
