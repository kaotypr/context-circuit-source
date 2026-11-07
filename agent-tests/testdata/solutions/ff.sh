#!/bin/sh
set -eu
sh "$WS/.context-circuit/cli" --json repo sync
sh "$WS/.context-circuit/cli" --json workspace sync
