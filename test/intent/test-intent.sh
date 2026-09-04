#!/bin/sh
# Context Circuit v1.0 — intent lifecycle (Mechanism 1, INV-INTENT-01).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

# --- id allocation: distinct i-prefixed sequence, never reused ---
id1=$(sh "$ROOT/wrapper/runtime/engine.sh" intent-allocate-id "$ws" checkout-retries)
assert_eq "i001-checkout-retries" "$id1"

cc_fx_intent "$ws" "$id1" "Checkout retries" checkout-service "src/checkout test/checkout"
# a second allocation is the next number even before the first is approved
id2=$(sh "$ROOT/wrapper/runtime/engine.sh" intent-allocate-id "$ws" coupon-field)
assert_eq "i002-coupon-field" "$id2"

# --- validate: structure + fields ---
sh "$ROOT/wrapper/runtime/engine.sh" intent-validate "$ws/intent/$id1" >/dev/null

# --- old form is rejected and the exact three-digit ceiling is explicit ---
old_id=$(printf 'i%04d-old-form' 1)
expect_failure cc_intent_id_valid "$old_id"
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" intent-archive "$ws" "$old_id"
mkdir -p "$ws/intent/i999-ceiling"
overflow=$(sh "$ROOT/wrapper/runtime/engine.sh" intent-allocate-id "$ws" exhausted 2>&1 || true)
printf '%s\n' "$overflow" | grep -Fq 'INTENT_ID_EXHAUSTED' || fail 'intent overflow must fail clearly'
require_dir "$ws/intent/i999-ceiling"
rmdir "$ws/intent/i999-ceiling"
mkdir -p "$ws/intent/archive/i999-ceiling"
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" intent-allocate-id "$ws" archived-exhausted
require_dir "$ws/intent/archive/i999-ceiling"
rmdir "$ws/intent/archive/i999-ceiling"

# an intent with no acceptance criterion is refused (at least one outcome criterion).
# Built in a scratch dir so it does not consume an intent id number.
mkdir -p "$ws/scratch-intent"
cc_fx_intent "$ws" i002-nocrit "No criteria" api "src"
mv "$ws/intent/i002-nocrit" "$ws/scratch-intent/i002-nocrit"
awk '/^acceptance_criteria:/{print "acceptance_criteria: []"; skip=1; next} skip && /^  -/{next} skip && /^    /{next} skip && /^[A-Za-z]/{skip=0; print; next} {print}' \
	"$ws/scratch-intent/i002-nocrit/contract.yaml" >"$ws/scratch-intent/i002-nocrit/contract.yaml.new"
mv "$ws/scratch-intent/i002-nocrit/contract.yaml.new" "$ws/scratch-intent/i002-nocrit/contract.yaml"
expect_failure sh "$ROOT/wrapper/runtime/engine.sh" intent-validate "$ws/scratch-intent/i002-nocrit"

# an empty (or absent) scope is now VALID — scope is coarse and optional in v1.0;
# a lay human often draws no paths, and scope-safety is settled at delivery (Gate 2).
cc_fx_intent "$ws" i003-noscope "No scope" api "src"
awk '/^scope:/{print "scope:"; print "  repositories: []"; skip=1; next} skip && /^  repositories:/{next} skip && /^    /{next} skip && /^[A-Za-z]/{skip=0; print; next} {print}' \
	"$ws/intent/i003-noscope/contract.yaml" >"$ws/intent/i003-noscope/contract.yaml.new"
mv "$ws/intent/i003-noscope/contract.yaml.new" "$ws/intent/i003-noscope/contract.yaml"
sh "$ROOT/wrapper/runtime/engine.sh" intent-validate "$ws/intent/i003-noscope" >/dev/null

# --- approval is the single upstream gate; it freezes contract_digest ---
assert_eq "draft" "$(cc_scalar "$ws/intent/$id1/contract.yaml" status)"
before=$(cc_scalar "$ws/intent/$id1/contract.yaml" contract_digest)
assert_eq "" "$before"
# Gate 1 needs no pre-approval challenge record: the tracer reads the real code
# only AFTER approval, so approval depends on the contract alone (no adversary).
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
# and re-enters Gate 1: a re-approval on the changed criteria is allowed and re-freezes
# the digest, while an unchanged re-approval stays refused (above).
sed -i.bak 's/Checkout retries works./Checkout retries works differently./' "$ws/intent/$id1/contract.yaml"
changed=$(cc_intent_contract_digest "$ws/intent/$id1/contract.yaml")
test "$changed" != "$frozen" || fail "criteria change did not change the digest"
sh "$ROOT/wrapper/runtime/engine.sh" intent-approve "$ws" "$id1" >/dev/null
refrozen=$(cc_scalar "$ws/intent/$id1/contract.yaml" contract_digest)
assert_eq "$changed" "$refrozen"
rm -f "$ws/intent/$id1/contract.yaml.bak"

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
assert_eq "i004-another" "$id_next"

sh "$ROOT/wrapper/runtime/engine.sh" intent-restore "$ws" "$id1" >/dev/null
require_dir "$ws/intent/$id1"
contains "$ws/intent/INDEX.md" "| $id1 |"

# --- current-source completeness: historical migration evidence is the only old-form exception ---
stale=$(rg -n --hidden --pcre2 -g '!sources/**' -g '!.git/**' \
	-g '!intent/i003-intent-id-width/trace/context-circuit-source.yaml' \
	'i[0-9]{4}(?:-[a-z0-9]+(?:-[a-z0-9]+)*)?' "$ROOT" || true)
[ -z "$stale" ] || { printf '%s\n' "$stale" >&2; fail 'stale four-digit intent reference remains'; }
stale_derivation=$(rg -n --hidden -g '!sources/**' -g '!.git/**' \
	-g '!intent/i003-intent-id-width/trace/context-circuit-source.yaml' \
	'intent: i[0-9]{4}|intent/(?:archive/)?i[0-9]{4}|i\$(?:pid|cc_fxp_pid|cc_fxe_pid)' "$ROOT" || true)
[ -z "$stale_derivation" ] || { printf '%s\n' "$stale_derivation" >&2; fail 'stale dynamic intent derivation remains'; }

# Approved records retain their status and use the unchanged digest algorithm after renaming.
for dir in "$ROOT"/intent/i*-* "$ROOT"/intent/archive/i*-*; do
	[ -d "$dir" ] || continue
	id=$(basename -- "$dir")
	test "$(sed -n 's/^intent: //p' "$dir/contract.yaml")" = "$id" || fail "intent field does not match $id"
	if [ "$(cc_scalar "$dir/contract.yaml" status)" = approved ]; then
		assert_eq "$(cc_scalar "$dir/contract.yaml" contract_digest)" "$(cc_intent_contract_digest "$dir/contract.yaml")"
	fi
done

pass 'intent lifecycle'
