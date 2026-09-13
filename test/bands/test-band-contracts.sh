#!/bin/sh
# Product surfaces and criteria-map coverage for band allocation (AC-BAND-01..09).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

W="$ROOT/.context-circuit/wrapper"
case_name=all
if [ "${1:-}" = --case ]; then case_name=${2:-}; fi

product_files="
$ROOT/product/.agents/skills/cc-intent/SKILL.md
$ROOT/product/.agents/skills/cc-trace/SKILL.md
$ROOT/product/.agents/skills/cc-plan/SKILL.md
$ROOT/product/.agents/skills/cc-workspace/SKILL.md
$ROOT/.context-circuit/agents/coordinator.md
$ROOT/.context-circuit/agents/planner.md
$ROOT/.context-circuit/docs/getting-started.md
$ROOT/.context-circuit/docs/terminology.md
$ROOT/.context-circuit/docs/planning.md
$ROOT/.context-circuit/wrapper/contracts/invariants.yaml
$ROOT/.context-circuit/wrapper/contracts/schemas/intent-contract.yaml
$ROOT/.context-circuit/wrapper/contracts/schemas/plan.yaml
$ROOT/context/domains/plan-review/README.md
$ROOT/context/domains/repository-binding/README.md
$ROOT/sources/system-design/context-circuit/v1.0.0/core/engine-and-seam.md
"

case_docs_skills() {
	for f in $product_files; do
		require_file "$f"
		not_contains "$f" "until band-scoped allocation"
		not_contains "$f" "workspace-wide sequence"
	done

	# AC-BAND-03: ordinary authoring never asks for a block number
	contains "$ROOT/product/.agents/skills/cc-intent/SKILL.md" "member-band-resolve"
	contains "$ROOT/product/.agents/skills/cc-intent/SKILL.md" "block number"
	contains "$ROOT/product/.agents/skills/cc-intent/SKILL.md" "INV-MEMBER-01"
	contains "$ROOT/product/.agents/skills/cc-intent/SKILL.md" "current member's intent band"
	contains "$ROOT/product/.agents/skills/cc-trace/SKILL.md" "block number"
	contains "$ROOT/product/.agents/skills/cc-trace/SKILL.md" "member's band"
	contains "$ROOT/product/.agents/skills/cc-plan/SKILL.md" "current member's plan band"
	contains "$ROOT/product/.agents/skills/cc-plan/SKILL.md" "block number"
	contains "$ROOT/product/.agents/skills/cc-workspace/SKILL.md" "members.yaml"
	contains "$ROOT/product/.agents/skills/cc-workspace/SKILL.md" "member.local.yaml"
	contains "$ROOT/product/.agents/skills/cc-workspace/SKILL.md" "block number"
	contains "$ROOT/product/.agents/skills/cc-workspace/SKILL.md" "exhausted"
	contains "$ROOT/.context-circuit/agents/coordinator.md" "member identity"
	contains "$ROOT/.context-circuit/agents/coordinator.md" "block number"
	contains "$ROOT/.context-circuit/agents/planner.md" "member-band-resolve"
	contains "$ROOT/.context-circuit/agents/planner.md" "block number"
	contains "$ROOT/.context-circuit/docs/getting-started.md" "block number"
	contains "$ROOT/.context-circuit/docs/terminology.md" "members.yaml"
	contains "$ROOT/.context-circuit/docs/terminology.md" "member.local.yaml"
	contains "$ROOT/.context-circuit/docs/terminology.md" "Number band"
	contains "$ROOT/sources/system-design/context-circuit/v1.0.0/core/engine-and-seam.md" "member-band-resolve"
	contains "$ROOT/sources/system-design/context-circuit/v1.0.0/core/engine-and-seam.md" "current member's intent band"

	# AC-BAND-02 / AC-BAND-07: short ids, branch form, no member token
	contains "$ROOT/context/domains/plan-review/README.md" "current member's plan band"
	contains "$ROOT/context/domains/plan-review/README.md" "cc/<plan-id>/<repository-id>"
	contains "$ROOT/context/domains/repository-binding/README.md" "members.yaml"
	contains "$ROOT/context/domains/repository-binding/README.md" "member.local.yaml"
	contains "$W/contracts/invariants.yaml" "INV-MEMBER-01"
	contains "$W/contracts/invariants.yaml" "cc/<plan-id>/<repository-id>"
	contains "$W/contracts/schemas/execution.yaml" "cc/<plan-id>/<repository-id>"
	contains "$W/contracts/schemas/intent-contract.yaml" "no member, machine, or workspace token"
	contains "$W/contracts/schemas/plan.yaml" "no member, machine, or workspace token"
	not_contains "$W/contracts/schemas/execution.yaml" "cc/<member"
	not_contains "$W/contracts/invariants.yaml" "cc/<member"

	# template seed is example-only; local identity never ships
	contains "$ROOT/template/members.yaml" "Example-only"
	contains "$ROOT/template/members.yaml" "example:"
	contains "$ROOT/scripts/release-manifest.txt" "required members.yaml"
	contains "$ROOT/scripts/release-manifest.txt" "exclude member.local.yaml"
}

case_criteria_map() {
	map="$ROOT/test/acceptance/criteria-map.yaml"
	require_file "$map"
	n=1
	while [ "$n" -le 9 ]; do
		id=$(printf 'AC-BAND-%02d' "$n")
		grep -q "id: $id" "$map" || fail "criteria map missing $id"
		suite=$(awk -v id="$id" '
			$0 ~ "id: " id {found=1; next}
			found && $0 ~ /^    suite:/ { print $2; exit }
		' "$map")
		[ -n "$suite" ] || fail "criteria map $id has no suite"
		case "$suite" in
			test/bands/test-band-allocation.sh|test/bands/test-member-roster.sh|test/bands/test-band-contracts.sh) ;;
			*) fail "criteria map $id suite is not a band suite: $suite" ;;
		esac
		[ -f "$ROOT/$suite" ] || fail "criteria map $id references missing suite $suite"
		n=$((n + 1))
	done
	contains "$ROOT/test/acceptance.sh" "test/bands/test-band-contracts.sh"
	contains "$ROOT/test/acceptance.sh" "test/bands/test-member-roster.sh"
	contains "$ROOT/test/acceptance.sh" "test/bands/test-band-allocation.sh"
	contains "$ROOT/test/acceptance.sh" "AC-BAND-"
}

run_case() {
	case "$1" in
		docs-skills) case_docs_skills ;;
		criteria-map) case_criteria_map ;;
		*) fail "unknown case $1" ;;
	esac
}

if [ "$case_name" = all ]; then
	run_case docs-skills
	run_case criteria-map
else
	run_case "$case_name"
fi

pass 'band contracts'
