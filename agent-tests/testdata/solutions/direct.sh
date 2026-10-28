#!/bin/sh
set -eu
readme=$WS/repositories/api/README.md
tmp=$TMPDIR/readme
awk '{gsub(/recieve/, "receive"); print}' "$readme" > "$tmp"
mv "$tmp" "$readme"
