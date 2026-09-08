#!/bin/sh

fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
pass() { printf 'PASS: %s\n' "$1"; }
require_file() { test -f "$1" || fail "missing file: $1"; }
require_dir() { test -d "$1" || fail "missing directory: $1"; }
contains() { grep -F -- "$2" "$1" >/dev/null 2>&1 || fail "expected '$2' in $1"; }
not_contains() { grep -F -- "$2" "$1" >/dev/null 2>&1 && fail "unexpected '$2' in $1" || :; }
assert_eq() { test "$1" = "$2" || fail "expected '$1' = '$2'"; }
expect_failure() { if "$@" >/dev/null 2>&1; then fail "expected failure: $*"; fi; }
# Repo root, resolved depth-independently so suites may live at any nesting
# (e.g. test/<suite>/ or the top-level agent-harness/). Prefer git toplevel
# when it actually contains the engine; otherwise walk up from this test file.
# A nested checkout (publish-template's .template-repo) or a leaked GIT_DIR
# can make `git rev-parse --show-toplevel` name the wrong tree.
_cc_engine_rel=".context-circuit/wrapper/runtime/engine.sh"
_cc_git_root=$(git -C "$(dirname -- "$0")" rev-parse --show-toplevel 2>/dev/null || true)
if [ -n "$_cc_git_root" ] && [ -f "$_cc_git_root/$_cc_engine_rel" ]; then
	ROOT=$_cc_git_root
else
	_cc_d=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
	ROOT=
	while [ -n "$_cc_d" ] && [ "$_cc_d" != / ]; do
		if [ -f "$_cc_d/$_cc_engine_rel" ]; then
			ROOT=$_cc_d
			break
		fi
		_cc_d=$(dirname "$_cc_d")
	done
fi
[ -n "$ROOT" ] && [ -f "$ROOT/$_cc_engine_rel" ] \
	|| fail "cannot locate $_cc_engine_rel from $0"
unset _cc_engine_rel _cc_git_root _cc_d
. "$ROOT/.context-circuit/wrapper/runtime/engine.sh"
