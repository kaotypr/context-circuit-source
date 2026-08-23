#!/bin/sh

fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
pass() { printf 'PASS: %s\n' "$1"; }
require_file() { test -f "$1" || fail "missing file: $1"; }
require_dir() { test -d "$1" || fail "missing directory: $1"; }
contains() { grep -F -- "$2" "$1" >/dev/null 2>&1 || fail "expected '$2' in $1"; }
not_contains() { grep -F -- "$2" "$1" >/dev/null 2>&1 && fail "unexpected '$2' in $1" || :; }
assert_eq() { test "$1" = "$2" || fail "expected '$1' = '$2'"; }
expect_failure() { if "$@" >/dev/null 2>&1; then fail "expected failure: $*"; fi; }
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
. "$ROOT/wrapper/runtime/engine.sh"
