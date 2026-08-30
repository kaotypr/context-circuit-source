#!/bin/sh
# Context Circuit semantic acceptance suite.
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)

run_suite() {
	printf '\n--- %s ---\n' "$1"
	sh "$ROOT/$1" || { printf 'FAIL: suite %s\n' "$1" >&2; exit 1; }
}

run_suite test/contracts/test-contracts.sh
run_suite test/runtime/test-runtime.sh
run_suite test/repositories/test-repositories.sh
run_suite test/plans/test-plans.sh
run_suite test/execution/test-execution.sh
run_suite test/concurrency/test-leases.sh
run_suite test/run-stack/test-run-stack.sh
run_suite test/latency/test-latency.sh
run_suite test/grounding/test-grounding.sh
run_suite test/pairing/test-pairing.sh
run_suite test/external-surface/test-external-surface.sh
run_suite test/completion/test-completion.sh
run_suite test/archive/test-archive.sh
run_suite test/delivery/test-delivery.sh
run_suite test/security/test-boundaries.sh
run_suite test/release/test-release.sh
run_suite test/release/test-publish.sh
run_suite template-harness/test-template-runtime.sh

# --- semantic criteria mapping cross-check ---
printf '\n--- acceptance criteria mapping ---\n'
map="$ROOT/test/acceptance/criteria-map.yaml"
[ -f "$map" ] || { printf 'FAIL: missing criteria map\n' >&2; exit 1; }
grep -q '^credential_free: true' "$map" || { printf 'FAIL: criteria map not credential-free\n' >&2; exit 1; }
grep -q '^implicit_external_checks: false' "$map" || { printf 'FAIL: criteria map declares implicit external checks\n' >&2; exit 1; }

n=1
while [ "$n" -le 29 ]; do
	id=$(printf 'AC-%02d' "$n")
	grep -q "id: $id" "$map" || { printf 'FAIL: criteria map missing %s\n' "$id" >&2; exit 1; }
	n=$((n + 1))
done

# every referenced suite exists
grep '^    suite:' "$map" | sed 's/^    suite:[[:space:]]*//' | sort -u | while IFS= read -r s; do
	[ -f "$ROOT/$s" ] || { printf 'FAIL: criteria map references missing suite %s\n' "$s" >&2; exit 1; }
done

printf '\nPASS: Context Circuit semantic acceptance\n'
