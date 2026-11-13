#!/bin/sh
# The host starts in the plan's worktree, which setup prepared, so the work
# lands in the current directory.
set -eu
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
test "$(git rev-parse --abbrev-ref HEAD)" = cc/p1000/api
cp "$here/files/export.go" internal/invoice/export.go
cp "$here/files/export_test.go" internal/invoice/export_test.go
git add -A
git commit -qm 'feat(invoice): export invoices as CSV'
