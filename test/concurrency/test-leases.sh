#!/bin/sh
# Path leases (INV-CONCURRENCY-01): region overlap, serialization of unrelated
# plans, descendant exemption, release, and reentrancy. Deterministic — a live
# model cannot be made to contend on a path on command, so it is proven here.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# --- region overlap semantics (equal / ancestor / repo-wide / disjoint) ---
cc_region_overlap "src/a" "src/a"        || fail "equal regions must overlap"
cc_region_overlap "src" "src/a"          || fail "ancestor must overlap descendant"
cc_region_overlap "src/a" "src"          || fail "descendant must overlap ancestor"
cc_region_overlap "." "src/a"            || fail "repo-wide must overlap any"
cc_region_overlap "src/a" "."            || fail "any must overlap repo-wide"
expect_failure cc_region_overlap "src/app" "src/application"   # sibling prefixes disjoint
expect_failure cc_region_overlap "src/a" "src/b"               # siblings disjoint

# --- plans: P1 (root), P2 (unrelated), P3 (depends on P1) ---
cc_fx_plan_ex "$ws" 0001-alpha  "Alpha"  api src/a ""
cc_fx_plan_ex "$ws" 0002-beta   "Beta"   api src/a ""
cc_fx_plan_ex "$ws" 0003-child  "Child"  api src/a "0001-alpha"

# --- P1 holds src/a; an unrelated P2 is serialized on an overlapping region ---
cc_lease_acquire "$ws" api 0001-alpha "src/a" >/dev/null
expect_failure cc_lease_check "$ws" api 0002-beta "src/a"          # exact overlap
expect_failure cc_lease_check "$ws" api 0002-beta "src"           # ancestor overlap
expect_failure cc_lease_check "$ws" api 0002-beta "."             # repo-wide overlap
cc_lease_check "$ws" api 0002-beta "src/b" >/dev/null             # disjoint -> free

# --- descendant exemption: P3 depends on P1, so it builds on it, never blocked ---
cc_lease_check "$ws" api 0003-child "src/a" >/dev/null || fail "descendant must be exempt"

# --- reentrancy: the holder re-checks its own overlapping region freely ---
cc_lease_check "$ws" api 0001-alpha "src/a" >/dev/null || fail "own lease must be reentrant"

# --- another repository is independent ---
cc_lease_check "$ws" web 0002-beta "src/a" >/dev/null || fail "other repo must be free"

# --- release frees the region for the unrelated plan ---
cc_lease_release "$ws" api 0001-alpha >/dev/null
cc_lease_check "$ws" api 0002-beta "src/a" >/dev/null || fail "released region must be free"
# the record is preserved (released_at set), never deleted
require_file "$ws/.runtime/locks/paths/api/0001-alpha.yaml"
contains "$ws/.runtime/locks/paths/api/0001-alpha.yaml" "schema_version: 1"
contains "$ws/.runtime/locks/paths/api/0001-alpha.yaml" "released_at:"

# --- acquire refuses when a non-descendant holds an overlapping region ---
cc_lease_acquire "$ws" api 0002-beta "src/a" >/dev/null
expect_failure cc_lease_acquire "$ws" api 0001-alpha "src/a"      # 0001 not a descendant of 0002
# 0003 depends on 0001 (not 0002), so 0002's lease still blocks it
expect_failure cc_lease_acquire "$ws" api 0003-child "src/a"
cc_lease_release "$ws" api 0002-beta >/dev/null
cc_lease_acquire "$ws" api 0001-alpha "src/a" >/dev/null || fail "region should be free after release"

pass 'path leases'
