#!/bin/sh
set -eu
sh "$WS/.context-circuit/cli" context find --query "invoice numbers" >/dev/null
