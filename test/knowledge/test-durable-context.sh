#!/bin/sh
# Durable-only Product Knowledge (INV-KNOWLEDGE-03): live context files never
# name a particular plan, intent file, or sources file, and never name a
# sources/ path.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

inv="$ROOT/.context-circuit/wrapper/contracts/invariants.yaml"
docs="$ROOT/.context-circuit/docs/product-knowledge.md"
contains "$inv" "INV-KNOWLEDGE-03"
contains "$docs" "INV-KNOWLEDGE-03"
contains "$docs" "durable-only"

test ! -f "$ROOT/context/sources.yaml" || fail "context/sources.yaml must be absent"
test ! -f "$ROOT/template/context/sources.yaml" || fail "template/context/sources.yaml must be absent"

scan_tree() {
	tree=$1
	label=$2
	if grep -REn -- 'sources/' "$tree" >/dev/null 2>&1; then
		grep -REn -- 'sources/' "$tree" >&2 || :
		fail "forbidden sources/ path in $label"
	fi
	if grep -REn -- 'sources\.yaml' "$tree" >/dev/null 2>&1; then
		grep -REn -- 'sources\.yaml' "$tree" >&2 || :
		fail "forbidden sources.yaml in $label"
	fi
	if grep -REn -- 'intent/i[0-9]{3}-' "$tree" >/dev/null 2>&1; then
		grep -REn -- 'intent/i[0-9]{3}-' "$tree" >&2 || :
		fail "forbidden particular intent file in $label"
	fi
	if grep -REn -- 'plans/[0-9]{4}-' "$tree" >/dev/null 2>&1; then
		grep -REn -- 'plans/[0-9]{4}-' "$tree" >&2 || :
		fail "forbidden particular plan id in $label"
	fi
	if grep -REn -- 'plans/context-circuit' "$tree" >/dev/null 2>&1; then
		grep -REn -- 'plans/context-circuit' "$tree" >&2 || :
		fail "forbidden particular plan path in $label"
	fi
}

scan_tree "$ROOT/context" "context/"
scan_tree "$ROOT/template/context" "template/context/"

contains "$ROOT/context/domains/intent/README.md" "Gate 1"
contains "$ROOT/context/domains/intent/README.md" "intent/<id>"
contains "$ROOT/context/SOURCES.md" "passive source"

pass 'durable-only Product Knowledge (INV-KNOWLEDGE-03)'
