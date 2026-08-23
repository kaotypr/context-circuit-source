#!/bin/sh
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)

run_suite() {
  suite=$1
  printf 'SUITE %s\n' "$suite"
  sh "$ROOT/$suite"
}

run_suite test/contracts/test-contracts.sh
run_suite test/routing/test-router.sh
run_suite test/hosts/test-host-adapters.sh
run_suite test/context-budget/test-budgets.sh
run_suite test/lifecycle/test-lifecycle.sh
run_suite test/approval/test-approval-flow.sh
run_suite test/ownership/test-ownership.sh
run_suite test/runtime/test-runtime.sh
run_suite test/gates/test-gates.sh
run_suite test/recovery/test-recovery.sh
run_suite test/stacks/test-stacks.sh
run_suite test/upgrades/test-upgrades.sh
run_suite test/security/test-boundaries.sh
run_suite test/behavior-matrix/test-matrix.sh
run_suite test/release/test-release.sh
test -f "$ROOT/wrapper/contracts/schemas/workspace.yaml"
grep -F 'missing_or_disagreeing_result: projection-mismatch' "$ROOT/wrapper/contracts/schemas/workspace.yaml" >/dev/null
grep -F 'authorization: never' "$ROOT/wrapper/contracts/routes.yaml" >/dev/null
grep -F 'wrapper/contracts/schemas/child-start.yaml' "$ROOT/scripts/release-manifest.txt" >/dev/null
mapping="$ROOT/test/acceptance/evidence-layer-mapping.yaml"
plan="$ROOT/plans/context-circuit-plans/verification-evidence-layers/plan.yaml"
test -f "$mapping" || { printf 'FAIL: missing evidence-layer acceptance mapping\n' >&2; exit 1; }
grep -F 'implicit_external_checks: false' "$mapping" >/dev/null || { printf 'FAIL: acceptance mapping allows implicit external checks\n' >&2; exit 1; }
grep -F 'credential_free: true' "$mapping" >/dev/null || { printf 'FAIL: acceptance mapping is not credential-free\n' >&2; exit 1; }
for vt in VEL-VT-01 VEL-VT-02 VEL-VT-03 VEL-VT-04 VEL-VT-05 VEL-VT-06 VEL-VT-07; do
  grep -F "id: $vt" "$mapping" >/dev/null || { printf 'FAIL: mapping missing %s\n' "$vt" >&2; exit 1; }
  grep -F "id: $vt" "$plan" >/dev/null || { printf 'FAIL: plan missing %s\n' "$vt" >&2; exit 1; }
done
grep -F 'sh test/contracts/test-contracts.sh' "$mapping" >/dev/null
grep -F 'sh test/runtime/test-runtime.sh' "$mapping" >/dev/null
grep -F 'sh test/hosts/test-host-adapters.sh' "$mapping" >/dev/null
grep -F 'sh test/security/test-boundaries.sh' "$mapping" >/dev/null
grep -F 'sh test/upgrades/test-upgrades.sh' "$mapping" >/dev/null
grep -F 'sh test/release/test-release.sh' "$mapping" >/dev/null
grep -F 'sh test/acceptance.sh' "$mapping" >/dev/null
printf 'PASS: Context Circuit semantic acceptance\n'
