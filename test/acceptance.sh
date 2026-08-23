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
printf 'PASS: Context Circuit semantic acceptance\n'
