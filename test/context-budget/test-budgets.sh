#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

sum_profile() {
  cc_budget=$1; cc_paths=$2; cc_sum=0
  old_ifs=$IFS; IFS=,
  for cc_path in $cc_paths; do
    require_file "$ROOT/$cc_path"
    cc_bytes=$(wc -c < "$ROOT/$cc_path")
    cc_sum=$((cc_sum + cc_bytes))
  done
  IFS=$old_ifs
  test "$cc_sum" -le "$cc_budget" || fail "budget exceeded: $cc_sum > $cc_budget"
  printf '%s\n' "$cc_sum"
}

while IFS='|' read -r profile budget paths; do
  case "$profile" in ''|'#'*) continue ;; esac
  measured=$(sum_profile "$budget" "$paths")
  printf 'BUDGET %s: %s/%s bytes\n' "$profile" "$measured" "$budget"
done < "$ROOT/test/context-budget/ledger.tsv"

profile_measured() {
  cc_name=$1
  cc_budget=$(awk -F'|' -v name="$cc_name" '$1 == name {print $2}' "$ROOT/test/context-budget/ledger.tsv")
  cc_paths=$(awk -F'|' -v name="$cc_name" '$1 == name {print $3}' "$ROOT/test/context-budget/ledger.tsv")
  sum_profile "$cc_budget" "$cc_paths"
}
root=$(profile_measured run-plan)
writer=$(profile_measured writer)
verifier=$(profile_measured verifier)
l1=$((root + writer + verifier))
test "$l1" -le 46080 || fail "L1 budget exceeded: $l1"
resume=$(profile_measured resume)
task_delta=$(wc -c < "$ROOT/docs/templates/task.md")
l2=$((l1 + (task_delta * 3) + resume))
test "$l2" -le 58368 || fail "L2 budget exceeded: $l2"
printf 'BUDGET l2-three-task-resume: %s/58368 bytes\n' "$l2"
baseline=154368
reduction=$(( (baseline - l1) * 100 / baseline ))
test "$reduction" -ge 70 || fail "static reduction below target: $reduction%"
host_overlay=$(wc -c < "$ROOT/wrapper/adapters/AGENTS.md")
host_overlay=$((host_overlay + $(wc -c < "$ROOT/wrapper/adapters/CLAUDE.md")))
host_overlay=$((host_overlay + $(wc -c < "$ROOT/wrapper/adapters/WORKFLOW.md")))
test "$host_overlay" -le 8192 || fail "host overlay exceeds tier-0 budget: $host_overlay"
test "$host_overlay" -le 12288 || fail "host overlay exceeds resume budget: $host_overlay"
host_fixtures=0
for host_fixture in "$ROOT"/test/hosts/fixtures/*.yaml; do
  host_fixtures=$((host_fixtures + $(wc -c < "$host_fixture")))
done
test "$host_fixtures" -le 8192 || fail "host fixture evidence exceeds tier-0 budget: $host_fixtures"
printf 'BUDGET host-overlay: %s/8192 bytes; fixtures: %s/8192 bytes\n' "$host_overlay" "$host_fixtures"
packet_bytes=$(cc_context_packet_measure "$ROOT" verifier \
  wrapper/contracts/schemas/delegation.yaml \
  wrapper/contracts/schemas/completion.yaml \
  agents/verifier.md | awk -F'|' '{sum += $3} END {print sum + 0}')
packet_budget=$(cc_context_set_budget "$ROOT" verifier)
test "$packet_bytes" -le "$packet_budget" || fail "measured verifier packet exceeded: $packet_bytes > $packet_budget"
printf 'PACKET verifier: %s/%s bytes\n' "$packet_bytes" "$packet_budget"
pass "context budgets and static reduction ${reduction}% (L1=$l1 bytes)"
