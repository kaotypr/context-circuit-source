#!/bin/sh
# Resolve the single template version a publication run should target, and print
# it as "version=<v>" (for $GITHUB_OUTPUT) and "resolved: <v>" (for logs).
#
# Inputs (env):
#   INPUT_VERSION  explicit version from workflow_dispatch (optional)
#   EVENT_NAME     github.event_name (workflow_dispatch | push)
#   BEFORE, SHA    push range, to find the changed request file
# Arg:
#   $1  path to the checked-out context-circuit-template repository
#
# Resolution:
#   - explicit input wins;
#   - on push, the request file added/changed in this push;
#   - otherwise the single request whose version has no tag yet.
# Ambiguity (zero or many) is a hard failure, never a guess.
set -eu

fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }

template_dir=${1:-}
[ -n "$template_dir" ] || fail 'usage: resolve-release-version.sh <template-checkout-dir>'

source_root=$(git rev-parse --show-toplevel 2>/dev/null) || fail 'not a git source checkout'
requests_dir="$source_root/release/requests"

version_of() { f=${1##*/}; printf '%s' "${f%.md}"; }
tag_exists() { git -C "$template_dir" rev-parse -q --verify "refs/tags/v$1" >/dev/null 2>&1; }

resolved=''

if [ -n "${INPUT_VERSION:-}" ]; then
  resolved=$INPUT_VERSION
elif [ "${EVENT_NAME:-}" = push ]; then
  changed=$(git -C "$source_root" diff --name-only "${BEFORE:?}" "${SHA:?}" -- 'release/requests/*.md' \
            | sed -n 's|^release/requests/\(.*\)\.md$|\1|p' | sort -u)
  count=$(printf '%s\n' "$changed" | grep -c . || true)
  [ "$count" -eq 1 ] || fail "push changed $count request files; expected exactly one"
  resolved=$changed
else
  candidates=''
  for f in "$requests_dir"/*.md; do
    [ -f "$f" ] || continue
    v=$(version_of "$f")
    tag_exists "$v" || candidates="$candidates $v"
  done
  set -- $candidates
  [ "$#" -eq 1 ] || fail "expected exactly one unpublished request, found $#: $*"
  resolved=$1
fi

case "$resolved" in ''|v*|*[!A-Za-z0-9.+-]*|.*|-*) fail "invalid version: $resolved" ;; esac

[ -f "$requests_dir/$resolved.md" ] || fail "no release request for: $resolved"

printf 'resolved: %s\n' "$resolved"
if [ -n "${GITHUB_OUTPUT:-}" ]; then
  printf 'version=%s\n' "$resolved" >> "$GITHUB_OUTPUT"
fi
