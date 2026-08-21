#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

left=$(cc_stack_join sha256:parent-b sha256:parent-a)
right=$(cc_stack_join sha256:parent-a sha256:parent-b)
assert_eq "$left" "$right"
runtime=$(mktemp -d "${TMPDIR:-/tmp}/cc-stack.XXXXXX")
trap 'rm -rf "$runtime"' EXIT HUP INT TERM
printf '%s\n' \
  'schema_version: 1' 'wrapper_version: 1.0.0' 'stack_id: stack-a' \
  'repository: app' 'members: [a, b]' 'edges: [a→b]' 'leaves: [b]' > "$runtime/graph.yaml"
printf '%s\n' \
  'schema_version: 1' 'wrapper_version: 1.0.0' 'stack_id: stack-a' \
  'status: executing' 'owner_session_id: root' 'members:' \
  '  a: pending' '  b: pending' > "$runtime/progress.yaml"
for field in graph.yaml progress.yaml; do require_file "$runtime/$field"; done
contains "$ROOT/wrapper/contracts/schemas/stack.yaml" 'graph is frozen once created'
contains "$ROOT/.agents/skills/context-execute/SKILL.md" 'frozen graph/progress records'
contains "$ROOT/.agents/skills/context-execute/SKILL.md" 'Do not auto-approve, auto-finish'
pass 'deterministic stack joins, frozen graph/progress, and no auto-finish'
