#!/bin/sh
set -eu
. "$(git -C "$(dirname -- "$0")" rev-parse --show-toplevel)/test/lib/assert.sh"

prep=$(sh "$ROOT/agent-harness/human/run-scenario.sh" --host codex --no-grade 14-codex-direct-collaboration)
run=$(printf '%s\n' "$prep" | sed -n 's/^prepared run: //p')
[ -n "$run" ] || fail 'scenario 14 did not produce a run directory'
trap 'rm -rf "$run"' EXIT HUP INT TERM

ws="$run/workspace"
require_dir "$ws"
contains "$ws/role-tiering.local.yaml" '  codex:'
contains "$ws/role-tiering.local.yaml" '  claude-code:'
contains "$ws/role-tiering.local.yaml" '  cursor-agent:'
contains "$ws/role-tiering.local.yaml" 'model: gpt-5.6-luna'
contains "$ws/role-tiering.local.yaml" 'effort: medium'
assert_eq 9 "$(grep -c '^      effort: medium$' "$ws/role-tiering.local.yaml")"
not_contains "$ROOT/agent-harness/scenarios/14-codex-direct-collaboration/case.yaml" 'role_tiering:'

# Seed the exact ground-truth state a successful live coordinator + worker must
# leave. This exercises scenario preparation and every deterministic grade rule
# without invoking a model.
# shellcheck disable=SC1090
. "$ws/.context-circuit/wrapper/runtime/engine.sh"
cc_pair_begin "$ws" widgets scenario-worker >/dev/null
pointer="$ws/.runtime/pairing/scenario-worker/pointer.yaml"
wt=$(cc_scalar "$pointer" worktree)
mkdir -p "$wt/src"
printf 'paired by Codex\n' > "$wt/src/pair-marker.txt"

printf '%s\t%s\t%s\n' worker gpt-5.6-luna medium > "$run/role-evidence.tsv"
cat > "$run/transcript.txt" <<'EOF'
# transcript: 14-codex-direct-collaboration (host=codex)

human: Work with me directly on the widgets project and leave the change uncommitted.
coordinator: We made the requested marker change together and checked its exact content. It remains uncommitted and available to continue later; nothing was delivered.
EOF

sh "$ROOT/agent-harness/human/grade.sh" "$run" >/dev/null
pass 'Codex direct-collaboration harness scenario'
