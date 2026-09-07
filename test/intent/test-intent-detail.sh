#!/bin/sh
# Optional intent detail: dual-root authoring, recommend/request, no second gate,
# short five-section surface, planning consumes detail, no v1.1.0 spec tree.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

sd="$ROOT/.agents/skills/cc-system-design/SKILL.md"
ci="$ROOT/.agents/skills/cc-intent/SKILL.md"
cp="$ROOT/.agents/skills/cc-plan/SKILL.md"
tr="$ROOT/.agents/skills/cc-trace/SKILL.md"
intent_dom="$ROOT/context/domains/intent/README.md"
sda="$ROOT/context/domains/system-design-authoring/README.md"
tracing="$ROOT/context/domains/tracing/README.md"
auth="$ROOT/context/domains/plan-authorization/README.md"
coord="$ROOT/.context-circuit/agents/coordinator.md"
planning="$ROOT/.context-circuit/docs/planning.md"
tmpl="$ROOT/.context-circuit/docs/templates/intent.md"

# AC-BOTH-ROOTS: same rubric, both homes
contains "$sd" "sources/system-design/"
contains "$sd" "intent/<id>/detail/"
contains "$sd" "three-tier"
contains "$sd" "by concern"
contains "$sda" "intent/<id>/detail/"
contains "$sda" "sources/system-design/"

# AC-OPTIONAL-RECOMMENDED
contains "$ci" "recommend"
contains "$ci" "request"
contains "$ci" "several concerns"
contains "$ci" "does not block"
contains "$ci" "never required"
contains "$intent_dom" "Skipping it does not block"
contains "$coord" "Skipping it does not block approval"

# AC-NO-SECOND-GATE
contains "$ci" "no second approval of the detail"
contains "$ci" "no separate approval"
contains "$cp" "not a second approval"
contains "$intent_dom" "no separate approval"
contains "$auth" "no approval status of its own"
contains "$planning" "no second approval gate"

# AC-SKILL-NOT-SOURCE-DESIGN
test ! -e "$ROOT/sources/system-design/context-circuit/v1.1.0" || fail "v1.1.0 spec tree must be absent"
test ! -e "$ROOT/.agents/skills/cc-intent-detail" || fail "cc-intent-detail must not exist"
contains "$ci" "intent/<id>/detail/"
not_contains "$sd" "v1.1.0/intent-detail"
not_contains "$ci" "v1.1.0/intent-detail"
not_contains "$cp" "v1.1.0/intent-detail"
not_contains "$sda" "v1.1.0/intent-detail"
not_contains "$intent_dom" "v1.1.0/intent-detail"

# AC-PLANNING-USES-DETAIL
contains "$cp" "confirmed shape"
contains "$cp" "without replacing"
contains "$planning" "confirmed shape"
contains "$tracing" "not a tracing output"
contains "$tr" "confirmed topic shape"

# AC-SHORT-SURFACE
contains "$tmpl" "exactly the five"
contains "$ci" "five human-facing sections"
contains "$intent_dom" "exactly five sections"
not_contains "$tmpl" "## Detail"
not_contains "$tmpl" "sixth section"
not_contains "$ROOT/.context-circuit/docs/templates/intent.example.md" "## Detail"

# Extra detail files do not affect digest or block a short intent from skipping them
ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM
iid=i001-short-change
cc_fx_intent "$ws" "$iid" "Short change" api "src/short"
sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-validate "$ws/intent/$iid" >/dev/null
sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-approve "$ws" "$iid" >/dev/null
frozen=$(cc_scalar "$ws/intent/$iid/contract.yaml" contract_digest)
mkdir -p "$ws/intent/$iid/detail"
printf '# overview\n' >"$ws/intent/$iid/detail/README.md"
assert_eq "$frozen" "$(cc_intent_contract_digest "$ws/intent/$iid/contract.yaml")"
# skipping detail still leaves an approved short intent
iid2=i002-skip-detail
cc_fx_intent "$ws" "$iid2" "Skip detail" api "src/skip"
sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-approve "$ws" "$iid2" >/dev/null
assert_eq "approved" "$(cc_scalar "$ws/intent/$iid2/contract.yaml" status)"
test ! -d "$ws/intent/$iid2/detail" || fail "small intent must be allowed to skip detail/"

pass 'intent detail'
