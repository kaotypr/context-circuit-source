#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

contains "$ROOT/docs/host-capabilities.md" 'host_id: codex | claude-code | cursor-agent'

fixture_root="$ROOT/test/hosts/fixtures"
for host in codex claude-code cursor-agent; do
  fixture="$fixture_root/$host.yaml"
  require_file "$fixture"
  contains "$fixture" "host_id: $host"
  contains "$fixture" 'host_evidence:'
  contains "$fixture" 'session_role: root'
  contains "$fixture" 'session_role: writer'
  contains "$fixture" 'session_role: verifier'
  contains "$fixture" 'permission_mode: bounded-write'
  contains "$fixture" 'permission_mode: read-only'
  contains "$fixture" 'worktree: exclusive'
  contains "$fixture" 'worktree: independent'
  contains "$fixture" 'write_worktree: true'
  contains "$fixture" 'write_worktree: false'
  contains "$fixture" 'checks: [receipt, wrapper, git, ownership]'
  contains "$fixture" 'offline_fallback: filesystem-only'
  contains "$fixture" 'outcome: host-blocked'
  contains "$fixture" 'live_smoke_status: unavailable'
done

for path in \
  AGENTS.md CLAUDE.md wrapper/adapters/AGENTS.md wrapper/adapters/CLAUDE.md \
  wrapper/adapters/WORKFLOW.md .agents/skills/cc-entry/SKILL.md \
  .agents/skills/cc-execute/SKILL.md .agents/skills/cc-verify/SKILL.md \
  agents/coordinator.md agents/writer.md agents/verifier.md; do
  require_file "$ROOT/$path"
done

for path in CLAUDE.md wrapper/adapters/AGENTS.md wrapper/adapters/CLAUDE.md; do
  contains "$ROOT/$path" 'AGENTS.md'
done
contains "$ROOT/AGENTS.md" 'Host adapters'
contains "$ROOT/wrapper/adapters/AGENTS.md" 'wrapper/contracts/invariants.yaml'
contains "$ROOT/wrapper/adapters/AGENTS.md" 'shared packet'
contains "$ROOT/wrapper/adapters/CLAUDE.md" 'shared bounded'
contains "$ROOT/wrapper/adapters/WORKFLOW.md" 'host-blocked'
contains "$ROOT/wrapper/adapters/WORKFLOW.md" 'shared packet loader'
contains "$ROOT/wrapper/adapters/WORKFLOW.md" 'cc_runtime_launch_projection'
contains "$ROOT/wrapper/adapters/WORKFLOW.md" 'cc_runtime_graph_authoritative'
contains "$ROOT/.agents/skills/cc-entry/SKILL.md" 'block-missing-child-primitive'
contains "$ROOT/.agents/skills/cc-execute/SKILL.md" 'host_evidence'
contains "$ROOT/.agents/skills/cc-execute/SKILL.md" 'cc_construct_runtime_graph'
contains "$ROOT/.agents/skills/cc-entry/SKILL.md" 'cc_validate_runtime_graph'
contains "$ROOT/.agents/skills/cc-verify/SKILL.md" 'independent'
contains "$ROOT/.agents/skills/cc-verify/SKILL.md" 'cc_validate_runtime_graph'
contains "$ROOT/agents/writer.md" 'exclusive worktree'
contains "$ROOT/agents/verifier.md" 'write_worktree: false'
contains "$ROOT/agents/coordinator.md" 'cc_runtime_graph_authoritative'
contains "$ROOT/docs/runtime-contract.md" 'commit.marker'
contains "$ROOT/docs/host-capabilities.md" 'cc_runtime_launch_projection'

for request in \
  'Start or resume work in this workspace with Codex CLI.' \
  'Start or resume work in this workspace with Claude Code.' \
  'Start or resume work in this workspace with Cursor Agent CLI.'; do
  route=$(cc_route "$request")
  printf '%s\n' "$route" | grep -F 'capability: recommend-next' >/dev/null ||
    fail "host entry changed canonical route: $request"
done
offline=$(cc_route 'The Claude Code provider is unavailable; continue offline.')
printf '%s\n' "$offline" | grep -F 'capability: offline-fallback' >/dev/null ||
  fail 'provider-unavailable fallback changed canonical route'
blocked=$(cc_route 'The Cursor Agent host cannot create a verifier child.')
printf '%s\n' "$blocked" | grep -F 'capability: block-missing-child-primitive' >/dev/null ||
  fail 'missing-child fallback changed canonical route'

if [ -d "$ROOT/test/hosts/fixtures" ]; then
  for fixture in "$ROOT"/test/hosts/fixtures/*.yaml; do
    [ -f "$fixture" ] || continue
    contains "$fixture" 'host_evidence:'
    not_contains "$fixture" 'password:'
    not_contains "$fixture" 'api_key:'
    not_contains "$fixture" 'provider_payload:'
    not_contains "$fixture" 'transcript:'
  done
fi

# Cursor host-local permission files are optional and are not part of the
# release surface. Existing maintainer-local files remain outside this task.
not_contains "$ROOT/scripts/release-manifest.txt" '.cursorrules'
not_contains "$ROOT/scripts/release-manifest.txt" '.cursor/rules'
git -C "$ROOT" diff --name-only --diff-filter=A -- .cursorrules .cursor/rules |
  grep . >/dev/null 2>&1 && fail 'task added a Cursor policy file' || :

# cc-plan remains the sole discovery adapter for both drafting and reviewing a
# plan; there is no separate cc-review-plan skill directory.
contains "$ROOT/.agents/skills/cc-plan/SKILL.md" 'review'
require_file "$ROOT/.agents/skills/cc-plan/SKILL.md"
test ! -d "$ROOT/.agents/skills/cc-review-plan" || fail 'a separate cc-review-plan skill directory exists'

# Optional native question prompts are host UI only; they never become a
# second router or a second gate.
contains "$ROOT/docs/host-capabilities.md" 'AskUserQuestion'
contains "$ROOT/docs/host-capabilities.md" 'selects no route, satisfies no gate'
contains "$ROOT/docs/plan-review.md" 'required writer or verifier child'
not_contains "$ROOT/wrapper/contracts/schemas/session.yaml" 'question_prompt_capability'

review_named=$(cc_route 'Review plan checkout-validation.')
printf '%s\n' "$review_named" | grep -F 'capability: review-plan' >/dev/null ||
  fail 'named-plan review route regressed'
review_unnamed=$(cc_route 'Read only the plan status and recommend the next action.')
printf '%s\n' "$review_unnamed" | grep -F 'capability: clarify-target' >/dev/null ||
  fail 'unnamed review request stopped clarifying the target'

# This is deliberately a label-only boundary. Offline CI never invokes a
# provider, even when a caller asks for the optional live smoke label.
for host in codex claude-code cursor-agent; do
  if [ "${CC_LIVE_HOST_SMOKE:-0}" = 1 ]; then
    printf 'LIVE %s: host-blocked (optional probe not invoked by offline suite)\n' "$host"
  else
    printf 'LIVE %s: unavailable (optional probe not invoked)\n' "$host"
  fi
done

contains "$ROOT/docs/host-capabilities.md" 'wrapper/contracts/schemas/plan.yaml'
contains "$ROOT/docs/host-capabilities.md" 'human-visible `waived` limitation'
contains "$ROOT/docs/host-capabilities.md" 'Live smoke remains optional and label-only.'
limit_root="$ROOT/test/hosts/fixtures/limitations"
for host in codex claude-code cursor-agent; do
  for case_name in unavailable host-blocked waived; do
    fixture="$limit_root/$host-$case_name.yaml"
    require_file "$fixture"
    contains "$fixture" "host_id: $host"
    contains "$fixture" 'converted_to_pass: false'
    contains "$fixture" 'live_check: none'
    not_contains "$fixture" 'password:'
    not_contains "$fixture" 'api_key:'
    not_contains "$fixture" 'provider_payload:'
    not_contains "$fixture" 'transcript:'
    result=$(cc_validate_evidence_mapping "$fixture" 2>&1 || true)
    printf '%s\n' "$result" | grep -F EVIDENCE_LAYER_MATCH >/dev/null && fail "$host $case_name became a false pass"
    case "$case_name" in
      waived)
        printf '%s\n' "$result" | grep -F EVIDENCE_LAYER_WAIVED >/dev/null || fail "$host waiver was not reported as waived"
        ;;
      *)
        printf '%s\n' "$result" | grep -F EVIDENCE_LAYER_BLOCKED >/dev/null || fail "$host $case_name was not blocked"
        ;;
    esac
  done
done
pass 'Codex, Claude Code, and Cursor adapters share route, role, resume, and offline boundaries'
