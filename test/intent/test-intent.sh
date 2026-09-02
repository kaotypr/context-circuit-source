#!/bin/sh
# Context Circuit v1.0 — intent lifecycle (Mechanism 1, INV-INTENT-01).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# --- id allocation: distinct i-prefixed sequence, never reused ---
id1=$(sh "$ROOT/wrapper/runtime/engine.sh" intent-allocate-id "$ws" checkout-retries)
assert_eq "i0001-checkout-retries" "$id1"

cc_fx_intent "$ws" "$id1" "Checkout retries" checkout-service "src/checkout test/checkout"
# a second allocation is the next number even before the first is approved
id2=$(sh "$ROOT/wrapper/runtime/engine.sh" intent-allocate-id "$ws" coupon-field)
assert_eq "i0002-coupon-field" "$id2"

# --- validate: structure + fields ---
sh "$ROOT/wrapper/runtime/engine.sh" intent-validate "$ws/intent/$id1" >/dev/null

# a criterion with an invalid method is refused (executable or explicitly manual only)
cp "$ws/intent/$id1/contract.yaml" "$ws/intent/$id1/contract.yaml.bak"
sed 's/method: test/method: vibes/' "$ws/intent/$id1/contract.yaml.bak" >"$ws/intent/$id1/contract.yaml"
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" intent-validate "$ws/intent/$id1"
mv "$ws/intent/$id1/contract.yaml.bak" "$ws/intent/$id1/contract.yaml"

# an empty scope envelope is refused
cc_fx_intent "$ws" i0003-noscope "No scope" api "src"
awk '/^scope:/{print "scope:"; print "  repositories: []"; skip=1; next} skip && /^  repositories:/{next} skip && /^    /{next} skip && /^[A-Za-z]/{skip=0; print; next} {print}' \
	"$ws/intent/i0003-noscope/contract.yaml" >"$ws/intent/i0003-noscope/contract.yaml.new"
mv "$ws/intent/i0003-noscope/contract.yaml.new" "$ws/intent/i0003-noscope/contract.yaml"
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" intent-validate "$ws/intent/i0003-noscope"

# --- approval is the single upstream gate; it freezes contract_digest ---
assert_eq "draft" "$(cc_scalar "$ws/intent/$id1/contract.yaml" status)"
before=$(cc_scalar "$ws/intent/$id1/contract.yaml" contract_digest)
assert_eq "" "$before"
# Gate 1 cannot be reached without an independent, digest-bound spec-adversary
# result; an absent result is not silently treated as a pass.
mv "$ws/intent/$id1/adversary.md" "$ws/intent/$id1/adversary.md.bak"
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" intent-approve "$ws" "$id1"
mv "$ws/intent/$id1/adversary.md.bak" "$ws/intent/$id1/adversary.md"
sh "$ROOT/wrapper/runtime/engine.sh" intent-approve "$ws" "$id1" >/dev/null
assert_eq "approved" "$(cc_scalar "$ws/intent/$id1/contract.yaml" status)"
frozen=$(cc_scalar "$ws/intent/$id1/contract.yaml" contract_digest)
test -n "$frozen" || fail "approval did not freeze a contract_digest"
case "$frozen" in sha256:*|cksum:*) : ;; *) fail "contract_digest not a digest: $frozen" ;; esac

# the frozen digest is stable and equals a recompute over criteria-bearing content
recompute=$(cc_intent_contract_digest "$ws/intent/$id1/contract.yaml")
assert_eq "$frozen" "$recompute"

# re-approving an already-approved intent is refused (no draft<-approved)
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" intent-approve "$ws" "$id1"

# a criteria change after approval breaks the frozen digest (basis of INV-CANDIDATE-01)
cp "$ws/intent/$id1/contract.yaml" "$ws/intent/$id1/contract.yaml.bak"
sed 's/method: test/method: manual/' "$ws/intent/$id1/contract.yaml.bak" >"$ws/intent/$id1/contract.yaml"
changed=$(cc_intent_contract_digest "$ws/intent/$id1/contract.yaml")
test "$changed" != "$frozen" || fail "criteria change did not change the digest"
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" intent-approve "$ws" "$id1"
mv "$ws/intent/$id1/contract.yaml.bak" "$ws/intent/$id1/contract.yaml"

# the index carries the approved row
contains "$ws/intent/INDEX.md" "| $id1 |"
contains "$ws/intent/INDEX.md" "approved"

# --- archive is a status-blind move; restore returns it; id is never reused ---
sh "$ROOT/wrapper/runtime/engine.sh" intent-archive "$ws" "$id1" >/dev/null
require_dir "$ws/intent/archive/$id1"
test ! -d "$ws/intent/$id1" || fail "archived intent still in active area"
not_contains "$ws/intent/INDEX.md" "| $id1 |"
# next id still advances past the archived one
id_next=$(sh "$ROOT/wrapper/runtime/engine.sh" intent-allocate-id "$ws" another)
assert_eq "i0004-another" "$id_next"

sh "$ROOT/wrapper/runtime/engine.sh" intent-restore "$ws" "$id1" >/dev/null
require_dir "$ws/intent/$id1"
contains "$ws/intent/INDEX.md" "| $id1 |"

pass 'intent lifecycle'
