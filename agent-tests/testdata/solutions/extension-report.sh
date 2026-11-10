#!/bin/sh
set -eu
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
sh "$WS/.context-circuit/cli" worktree prepare --repo api --plan p1000
wt=$WS/.worktrees/p1000/api
cp "$here/files/export.go" "$wt/internal/invoice/export.go"
cp "$here/files/export_test.go" "$wt/internal/invoice/export_test.go"
git -C "$wt" add -A
git -C "$wt" commit -qm 'feat(invoice): export invoices as CSV'
# The after implement hook fails and only reports; the stage stands.
sh "$WS/.context-circuit/cli" extension hooks --stage implement --when after --plan p1000 --repo api
