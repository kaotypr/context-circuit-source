#!/bin/sh
# Host-neutral runtime primitives. This file is sourced by adapters and tests;
# it is not a user-facing command or scheduler.

cc_lower() {
  printf '%s' "$1" | tr '[:upper:]' '[:lower:]'
}

cc_safe_id() {
  case "$1" in
    ''|*[!a-zA-Z0-9._-]*) return 1 ;;
    [0-9a-zA-Z]*) [ "${#1}" -le 96 ] ;;
    *) return 1 ;;
  esac
}

cc_safe_relative() {
  case "$1" in
    ''|/*|../*|*/../*|*/..|.|./*|*' '*) return 1 ;;
    *) return 0 ;;
  esac
}

cc_probe() {
  cc_probe_text=$(cc_lower "$1")
  case "$cc_probe_text" in
    *'what is this workspace'*|*'what is this project'*|*'do not change anything'*|*'explain the current state'*) printf '%s\n' orientation ;;
    *'bootstrap'*|*'clone repository'*|*'clone '* ) printf '%s\n' repository-bootstrap ;;
    *'initialize'*|*'set up this workspace'*) printf '%s\n' initialization ;;
    *'idea'*|*'uncertain idea'*) printf '%s\n' idea-brief ;;
    *'draft a plan'*|*'create a plan'*|*'plan from this'*) printf '%s\n' plan-draft ;;
    *'prd'*|*'requirements from'*|*'requirements'*) printf '%s\n' prd ;;
    *'refresh'*'context'*|*'domain context'*|*'role context'*) printf '%s\n' context-refresh ;;
    *'sources/'*|*'selected files'*|*'selected evidence'*) printf '%s\n' selected-source ;;
    *'contradiction'*|*'accepted context'*) printf '%s\n' context-review ;;
    *'review plan'*|*'plan ready'*|*'plan'*'ready'*|*'approved-for-execution'*|*'read only the plan status'*) printf '%s\n' plan-review ;;
    *'commit approved plan'*|*'commit the approved plan'*|*'commit plan approval'*|*'commit'*'approved plan'*) printf '%s\n' maintainer-commit ;;
    *'approve plan'*|*'approve this plan'*|*'confirm approval'*) printf '%s\n' approval ;;
    *'base checkout'*'dirty'*|*'dirty base'*) printf '%s\n' ownership ;;
    *'verification passed'*|*'verification'*'what happens next'*) printf '%s\n' completion ;;
    *'run draft plan'*|*'run approved plan'*|*'execute '*|*'execute'*'plan'*|*'run the plan'*) printf '%s\n' execution-preflight ;;
    *'connected approved plans'*|*'plan stack'*|*'connected plan stack'*) printf '%s\n' stack-preflight ;;
    *'dependency'*'archived'*|*'unresolved dependency'*) printf '%s\n' dependency ;;
    *'start or resume'*) printf '%s\n' session-preflight ;;
    *'resume'*|*'interrupted work'*|*'continue my'*'work'*) printf '%s\n' resume-preflight ;;
    *'take over'*|*'stale lease'*|*'takeover'*) printf '%s\n' recovery ;;
    *'live plan'*'another session'*|*'foreign owner'*) printf '%s\n' ownership ;;
    *'base checkout'*'dirty'*|*'dirty base'*) printf '%s\n' ownership ;;
    *'child packet'*|*'delegation'*|*'verifier child'*|*'host cannot create'*) printf '%s\n' child-entry ;;
    *'unapproved path'*|*'scope expansion'*) printf '%s\n' scope-recovery ;;
    *'verify'*|*'verification'*|*'repair the failed'*) printf '%s\n' verification ;;
    *'finish plan'*|*'confirm finish'*|*'verification passed'**) printf '%s\n' completion ;;
    *'archive plan'*|*'confirm archive'*|*'restore archived'*|*'confirm restore'*) printf '%s\n' archive ;;
    *'configure'*'delivery'*|*'configuration'*) printf '%s\n' configuration ;;
    *'remote review'*|*'commit and push'*|*'delivery'*) printf '%s\n' delivery ;;
    *'provider'*'unavailable'*|*'offline'*) printf '%s\n' integration ;;
    *'handoff'*|*'blockers'*) printf '%s\n' handoff ;;
    *'clean runtime'*|*'cleanup'*|*'runtime cleanup'*) printf '%s\n' cleanup ;;
    *'upgrade'*'wrapper'*|*'wrapper migration'*) printf '%s\n' compatibility ;;
    *'eligible next'*|*'what can i'*'next'*) printf '%s\n' session-preflight ;;
    *) printf '%s\n' orientation ;;
  esac
}

cc_action() {
  cc_action_text=$(cc_lower "$1")
  cc_action_probe=$2
  case "$cc_action_probe" in
    orientation) printf '%s\n' orient ;;
    initialization) printf '%s\n' initialize ;;
    repository-bootstrap)
      case "$cc_action_text" in
        *'confirm'*|*'confirmed'*) printf '%s\n' repository-bootstrap ;;
        *) printf '%s\n' present-repository-bootstrap-card ;;
      esac
      ;;
    idea-brief) printf '%s\n' draft-idea-brief ;;
    prd) printf '%s\n' draft-prd ;;
    selected-source) printf '%s\n' inspect-source ;;
    context-refresh) printf '%s\n' propose-context ;;
    context-review) printf '%s\n' review-context ;;
    plan-draft) printf '%s\n' draft-plan ;;
    plan-review) printf '%s\n' review-plan ;;
    approval)
      case "$cc_action_text" in *'approve this plan'*) printf '%s\n' clarify-target ;; *'confirm approval'*) printf '%s\n' approve-plan ;; *'yes'*|*'confirmed'*) printf '%s\n' approve-plan ;; *) printf '%s\n' present-approval-card ;; esac
      ;;
    maintainer-commit)
      case "$cc_action_text" in *'confirm'*|*'confirmed'*) printf '%s\n' commit-approved-plan ;; *) printf '%s\n' present-commit-approved-plan-card ;; esac
      ;;
    execution-preflight)
      case "$cc_action_text" in *'draft plan'*) printf '%s\n' block-unapproved ;; *) printf '%s\n' execute-plan ;; esac
      ;;
    stack-preflight) printf '%s\n' execute-stack ;;
    resume-preflight) printf '%s\n' resume-session ;;
    recovery)
      case "$cc_action_text" in *'confirm'*|*'takeover'*'confirmed'*) printf '%s\n' takeover-lease ;; *) printf '%s\n' present-takeover-card ;; esac
      ;;
    ownership)
      case "$cc_action_text" in *'base checkout'*|*'dirty base'*) printf '%s\n' block-dirty-base ;; *) printf '%s\n' block-foreign-owner ;; esac
      ;;
    child-entry)
      case "$cc_action_text" in *'host cannot'*|*'cannot create'* ) printf '%s\n' block-missing-child-primitive ;; *) printf '%s\n' block-incomplete-delegation ;; esac
      ;;
    scope-recovery) printf '%s\n' block-scope-expansion ;;
    verification)
      case "$cc_action_text" in *'repair'*) printf '%s\n' repair-plan ;; *) printf '%s\n' verify-plan ;; esac
      ;;
    completion)
      case "$cc_action_text" in *'confirm finish'*) printf '%s\n' finish-plan ;; *'verification passed'*) printf '%s\n' recommend-finish ;; *) printf '%s\n' present-finish-card ;; esac
      ;;
    archive)
      case "$cc_action_text" in *'confirm restore'*|*'restore'*'confirm'*) printf '%s\n' restore-plan ;; *'restore archived'*) printf '%s\n' present-restore-card ;; *'confirm archive'*) printf '%s\n' archive-plan ;; *) printf '%s\n' present-archive-card ;; esac
      ;;
    configuration)
      case "$cc_action_text" in *'confirm'*) printf '%s\n' configure-delivery ;; *) printf '%s\n' present-configuration-card ;; esac
      ;;
    delivery) printf '%s\n' present-delivery-card ;;
    integration) printf '%s\n' offline-fallback ;;
    handoff) printf '%s\n' show-handoff ;;
    cleanup)
      case "$cc_action_text" in *'confirm'*) printf '%s\n' cleanup-runtime ;; *) printf '%s\n' inspect-cleanup ;; esac
      ;;
    compatibility)
      case "$cc_action_text" in *'confirm'*|*'required migration'*) printf '%s\n' migrate-wrapper ;; *) printf '%s\n' classify-upgrade ;; esac
      ;;
    session-preflight) printf '%s\n' recommend-next ;;
    dependency) printf '%s\n' block-unresolved-dependency ;;
    *) printf '%s\n' orient ;;
  esac
}

cc_workspace_uninitialized() {
  cc_workspace_root=$1
  cc_workspace_file="$cc_workspace_root/workspace.yaml"
  test -f "$cc_workspace_file" || return 0
  cc_workspace_kind=$(awk '
    /^identity:/ { in_identity=1; next }
    /^[^[:space:]]/ { in_identity=0 }
    in_identity && /^  kind: / { print $2; exit }
  ' "$cc_workspace_file")
  case "$cc_workspace_kind" in
    product-source) return 1 ;;
    instantiated-workspace)
      cc_workspace_status=$(awk '
        /^identity:/ { in_identity=1; next }
        /^[^[:space:]]/ { in_identity=0 }
        in_identity && /^  status: / { print $2; exit }
      ' "$cc_workspace_file")
      test "$cc_workspace_status" = accepted || return 0
      return 1
      ;;
    *) return 0 ;;
  esac
}

cc_request_read_only() {
  cc_read_only_text=$(cc_lower "$1")
  case "$cc_read_only_text" in
    *'what is this workspace'*|*'what is this project'*|*'explain'*|*'show'*|*'review'*|*'inspect'*|*'read-only'*|*'do not change'*|*'current state'*|*'status'*) return 0 ;;
    *) return 1 ;;
  esac
}

cc_route() {
  cc_request=$1
  cc_workspace_root=${2:-.}
  cc_stage_a=$(cc_probe "$cc_request")
  cc_stage_b=$(cc_action "$cc_request" "$cc_stage_a")
  cc_initialization_guard=0
  cc_plan_required=0
  if cc_workspace_uninitialized "$cc_workspace_root" && ! cc_request_read_only "$cc_request"; then
    cc_stage_a=initialization
    cc_stage_b=initialize
    cc_initialization_guard=1
  elif test "$cc_stage_a" = orientation && test "$cc_stage_b" = orient && ! cc_request_read_only "$cc_request"; then
    cc_stage_a=plan-draft
    cc_stage_b=draft-plan
    cc_plan_required=1
  fi
  if test "$cc_initialization_guard" -eq 1; then
    cc_eligibility=eligible-with-gate
    cc_gate=identity-acceptance
    cc_reason=IDENTITY_ACCEPTANCE_REQUIRED
    cc_authorization=confirmed-gate-required
  elif test "$cc_plan_required" -eq 1; then
    cc_eligibility=ready
    cc_gate=none
    cc_reason=PLAN_REQUIRED
    cc_authorization=explicitly-requested
  else
    case "$cc_stage_b" in
      block-*) cc_eligibility=blocked; cc_gate=none; cc_reason=SAFETY_BLOCK ;;
      present-repository-bootstrap-card) cc_eligibility=eligible-with-gate; cc_gate=repository-bootstrap; cc_reason=HUMAN_CONFIRMATION_REQUIRED ;;
      present-*) cc_eligibility=eligible-with-gate; cc_gate=$cc_stage_b; cc_reason=HUMAN_CONFIRMATION_REQUIRED ;;
      approve-plan|commit-approved-plan|finish-plan|archive-plan|restore-plan|configure-delivery|cleanup-runtime|takeover-lease|migrate-wrapper|repository-bootstrap)
        cc_eligibility=ready; cc_gate=none; cc_reason=CONFIRMED_HUMAN_GATE ;;
      *) cc_eligibility=ready; cc_gate=none; cc_reason=ROUTE_SELECTED ;;
    esac
    case "$cc_stage_b" in
      orient|review-*|show-handoff|recommend-*|classify-upgrade|offline-fallback) cc_authorization=read-only ;;
      present-*) cc_authorization=confirmed-gate-required ;;
      block-*) cc_authorization=absent ;;
      approve-plan|commit-approved-plan|finish-plan|archive-plan|restore-plan|configure-delivery|cleanup-runtime|takeover-lease|migrate-wrapper|repository-bootstrap) cc_authorization=confirmed-gate ;;
      *) cc_authorization=explicitly-requested ;;
    esac
  fi
  cat <<EOF
intent: $cc_stage_b
session_kind: root
phase: $cc_stage_a
probe: $cc_stage_a
eligibility: $cc_eligibility
capability: $cc_stage_b
authorization: $cc_authorization
reason_codes:
  - $cc_reason
context_set: $cc_stage_a
human_gate: $cc_gate
EOF
}

cc_atomic_write() {
  cc_target=$1
  shift
  cc_parent=$(dirname "$cc_target")
  mkdir -p "$cc_parent" || return 1
  cc_tmp="$cc_target.tmp.$$"
  (printf '%s\n' "$@") > "$cc_tmp" || return 1
  mv "$cc_tmp" "$cc_target"
}

cc_digest() {
  if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | awk '{print "sha256:" $1}'; else cksum "$1" | awk '{print "cksum:" $1}'; fi
}

cc_acquire_lease() {
  cc_runtime=$1; cc_plan=$2; cc_session=$3; cc_root=$4; cc_repo=$5; cc_worktree=$6; cc_branch=$7
  cc_safe_id "$cc_plan" && cc_safe_id "$cc_session" && cc_safe_id "$cc_root" || return 1
  cc_plan_dir="$cc_runtime/plans/$cc_plan"
  cc_lock="$cc_plan_dir/lease.lock"
  mkdir -p "$cc_plan_dir"
  mkdir "$cc_lock" 2>/dev/null || { printf '%s\n' 'LEASE_CONFLICT' >&2; return 1; }
  cc_atomic_write "$cc_lock/owner.yaml" \
    'schema_version: 1' "wrapper_version: ${CC_WRAPPER_VERSION:-1.0.0}" "plan: $cc_plan" \
    "session_id: $cc_session" "root_session_id: $cc_root" "repository: $cc_repo" \
    "worktree: $cc_worktree" "branch: $cc_branch" || return 1
  cc_atomic_write "$cc_plan_dir/lease.yaml" \
    'schema_version: 1' "wrapper_version: ${CC_WRAPPER_VERSION:-1.0.0}" "plan: $cc_plan" \
    "session_id: $cc_session" "root_session_id: $cc_root" "repository: $cc_repo" \
    "worktree: $cc_worktree" "branch: $cc_branch" 'status: active' \
    'acquired_at: 2026-08-21T00:00:00Z' 'heartbeat_at: 2026-08-21T00:00:00Z' \
    'released_at: null' 'stale_after_seconds: 1800'
}

cc_prepare_worktree() {
  cc_base=$1; cc_target=$2; cc_branch=$3
  test -d "$cc_base/.git" || return 1
  git -C "$cc_base" diff --quiet --ignore-submodules -- . || { printf '%s\n' DIRTY_BASE_BLOCKED >&2; return 1; }
  git -C "$cc_base" diff --cached --quiet -- . || { printf '%s\n' DIRTY_BASE_BLOCKED >&2; return 1; }
  test -z "$(git -C "$cc_base" status --porcelain --untracked-files=all)" || { printf '%s\n' DIRTY_BASE_BLOCKED >&2; return 1; }
  test ! -e "$cc_target" || { printf '%s\n' WORKTREE_ALREADY_EXISTS >&2; return 1; }
  mkdir -p "$(dirname "$cc_target")"
  git -C "$cc_base" worktree add --detach "$cc_target" "$cc_branch" >/dev/null
}

cc_is_product_source() {
  cc_workspace_root=$1
  test -f "$cc_workspace_root/workspace.yaml" || return 1
  awk '
    /^identity:/ { in_identity=1; next }
    /^[^[:space:]]/ { in_identity=0 }
    in_identity && /^  kind: product-source$/ { found=1 }
    END { exit found ? 0 : 1 }
  ' "$cc_workspace_root/workspace.yaml"
}

cc_maintainer_approval_commit_required() {
  cc_base=$1
  cc_plan_file=$2
  cc_task_dir=$3
  cc_base=$(CDPATH= cd -- "$cc_base" 2>/dev/null && pwd) || return 1
  case "$cc_plan_file" in /*) ;; *) cc_plan_file="$cc_base/$cc_plan_file" ;; esac
  case "$cc_task_dir" in /*) ;; *) cc_task_dir="$cc_base/$cc_task_dir" ;; esac
  cc_is_product_source "$cc_base" || return 1
  test -f "$cc_plan_file" && test -d "$cc_task_dir" || return 1
  test ! -L "$cc_plan_file" && test ! -L "$cc_task_dir" || return 1
  cc_plan_rel=${cc_plan_file#"$cc_base"/}
  cc_task_rel=${cc_task_dir#"$cc_base"/}
  cc_safe_relative "$cc_plan_rel" && cc_safe_relative "$cc_task_rel" || return 1
  case "$cc_plan_rel:$cc_task_rel" in plans/*:plans/*) ;; *) return 1 ;; esac
  git -C "$cc_base" diff --cached --quiet -- . || return 1
  test -z "$(git -C "$cc_base" ls-files --others --exclude-standard)" || return 1

  cc_changed=$(git -C "$cc_base" diff --name-only -- .)
  cc_task_count=0
  for cc_task_file in "$cc_task_dir"/*.md; do
    test -f "$cc_task_file" || continue
    test ! -L "$cc_task_file" || return 1
    cc_task_count=$((cc_task_count + 1))
  done
  test "$cc_task_count" -gt 0 || return 1
  cc_changed_count=$(printf '%s\n' "$cc_changed" | awk 'NF { count++ } END { print count + 0 }')
  test "$cc_changed_count" -eq $((cc_task_count + 1)) || return 1

  cc_plan_before=$(git -C "$cc_base" show "HEAD:$cc_plan_rel" 2>/dev/null | sed -n 's/^status: //p' | head -n 1) || return 1
  cc_plan_after=$(sed -n 's/^status: //p' "$cc_plan_file" | head -n 1)
  test "$cc_plan_before" = draft && test "$cc_plan_after" = approved || return 1
  printf '%s\n' "$cc_changed" | grep -Fx "$cc_plan_rel" >/dev/null 2>&1 || return 1
  cc_plan_diff=$(git -C "$cc_base" diff --unified=0 -- "$cc_plan_rel" | grep -E '^[+-][^-+]' || :)
  test "$(printf '%s\n' "$cc_plan_diff" | awk 'NF { count++ } END { print count + 0 }')" -eq 2 || return 1
  printf '%s\n' "$cc_plan_diff" | grep -Fx -- '-status: draft' >/dev/null 2>&1 || return 1
  printf '%s\n' "$cc_plan_diff" | grep -Fx -- '+status: approved' >/dev/null 2>&1 || return 1

  for cc_task_file in "$cc_task_dir"/*.md; do
    test -f "$cc_task_file" || continue
    cc_task_rel_file=${cc_task_file#"$cc_base"/}
    cc_task_before=$(git -C "$cc_base" show "HEAD:$cc_task_rel_file" 2>/dev/null | sed -n 's/^status: //p' | head -n 1) || return 1
    cc_task_after=$(sed -n 's/^status: //p' "$cc_task_file" | head -n 1)
    test "$cc_task_before" = draft && test "$cc_task_after" = ready || return 1
    printf '%s\n' "$cc_changed" | grep -Fx "$cc_task_rel_file" >/dev/null 2>&1 || return 1
    cc_task_diff=$(git -C "$cc_base" diff --unified=0 -- "$cc_task_rel_file" | grep -E '^[+-][^-+]' || :)
    test "$(printf '%s\n' "$cc_task_diff" | awk 'NF { count++ } END { print count + 0 }')" -eq 2 || return 1
    printf '%s\n' "$cc_task_diff" | grep -Fx -- '-status: draft' >/dev/null 2>&1 || return 1
    printf '%s\n' "$cc_task_diff" | grep -Fx -- '+status: ready' >/dev/null 2>&1 || return 1
  done
  printf '%s\n' MAINTAINER_APPROVAL_COMMIT_REQUIRED
}

cc_release_lease() {
  cc_runtime=$1; cc_plan=$2; cc_session=$3
  cc_lease="$cc_runtime/plans/$cc_plan/lease.yaml"
  cc_owner="$cc_runtime/plans/$cc_plan/lease.lock/owner.yaml"
  test -f "$cc_lease" && test -f "$cc_owner" || return 1
  grep -F "session_id: $cc_session" "$cc_owner" >/dev/null 2>&1 || return 1
  cc_atomic_write "$cc_lease" \
    'schema_version: 1' "wrapper_version: ${CC_WRAPPER_VERSION:-1.0.0}" "plan: $cc_plan" \
    "session_id: $cc_session" 'status: released' 'released_at: 2026-08-21T00:00:00Z'
  rmdir "$cc_runtime/plans/$cc_plan/lease.lock" 2>/dev/null || return 1
}

cc_projection() {
  case "$1" in draft) printf '%s\n' draft ;; approved) printf '%s\n' ready ;; done) printf '%s\n' done ;; *) return 1 ;; esac
}

cc_replace_first_status() {
  cc_status_file=$1
  cc_status_value=$2
  case "$cc_status_value" in
    draft|approved|ready|done) ;;
    *) return 1 ;;
  esac
  test -f "$cc_status_file" || return 1
  cc_status_tmp="$cc_status_file.tmp.$$"
  awk -v status="$cc_status_value" '
    BEGIN { replaced=0 }
    {
      if (!replaced && $0 ~ /^status: /) {
        print "status: " status
        replaced=1
      } else {
        print
      }
    }
  ' "$cc_status_file" > "$cc_status_tmp" || { rm -f "$cc_status_tmp"; return 1; }
  if test -s "$cc_status_file"; then
    cc_orig_end=$(tail -c 1 "$cc_status_file" | od -An -tx1 | tr -d ' \n')
    if test "$cc_orig_end" != 0a; then
      cc_size=$(wc -c < "$cc_status_tmp")
      cc_size=$(printf '%s' "$cc_size" | tr -d ' \t\n')
      if test "$cc_size" -gt 0; then
        dd if="$cc_status_tmp" of="$cc_status_tmp.raw" bs=1 count=$((cc_size - 1)) 2>/dev/null || {
          rm -f "$cc_status_tmp" "$cc_status_tmp.raw"
          return 1
        }
        mv "$cc_status_tmp.raw" "$cc_status_tmp"
      fi
    fi
  fi
  mv "$cc_status_tmp" "$cc_status_file"
}

cc_transition_plan_status() {
  cc_plan_file=$1; cc_task_dir=$2; cc_target=$3; cc_confirmation=${4:-}
  test "$cc_confirmation" = confirmed || { printf '%s\n' GATE_REQUIRED >&2; return 1; }
  cc_current=$(sed -n 's/^status: //p' "$cc_plan_file" | head -n 1)
  case "$cc_current:$cc_target" in
    draft:approved|approved:done) ;;
    *) printf '%s\n' INVALID_LIFECYCLE_TRANSITION >&2; return 1 ;;
  esac
  cc_expected=$(cc_projection "$cc_target") || return 1
  cc_replace_first_status "$cc_plan_file" "$cc_target" || return 1
  for cc_task in "$cc_task_dir"/*.md; do
    test -f "$cc_task" || continue
    cc_replace_first_status "$cc_task" "$cc_expected" || return 1
  done
}

cc_archive_event() {
  cc_archive_file=$1; cc_plan=$2; cc_event=$3; cc_actor=$4; cc_reason=$5
  case "$cc_event" in archived|restored) ;; *) return 1 ;; esac
  if test -f "$cc_archive_file"; then
    grep -F "plan: $cc_plan" "$cc_archive_file" >/dev/null 2>&1 || return 1
  else
    cc_atomic_write "$cc_archive_file" 'schema_version: 1' "plan: $cc_plan" 'events:'
  fi
  cc_tmp="$cc_archive_file.tmp.$$"
  cat "$cc_archive_file" > "$cc_tmp"
  cat >> "$cc_tmp" <<EOF
  - event: $cc_event
    plan: $cc_plan
    actor: $cc_actor
    timestamp: 2026-08-21T00:00:00Z
    reason: $cc_reason
    observed_status: preserved
EOF
  mv "$cc_tmp" "$cc_archive_file"
}

cc_validate_receipt() {
  cc_receipt=$1
  for cc_field in schema_version wrapper_version session_id route_decision_digest context_set references invariants created_at; do
    grep -E "^$cc_field:" "$cc_receipt" >/dev/null 2>&1 || return 1
  done
  grep -F 'references:' "$cc_receipt" >/dev/null 2>&1 || return 1
}

cc_receipt_delta() {
  cc_receipt=$1
  cc_expected_wrapper=${2:-${CC_WRAPPER_VERSION:-1.0.0}}
  grep -F "wrapper_version: $cc_expected_wrapper" "$cc_receipt" >/dev/null 2>&1 || { printf '%s\n' RELOAD_WRAPPER_CONTEXT; return 0; }
  awk '/^  - path: / { path=$3 } /^    revision: / { print path "|" $2 }' "$cc_receipt" | while IFS='|' read -r cc_path cc_revision; do
    test -f "$cc_path" || { printf 'RELOAD: %s\n' "$cc_path"; continue; }
    cc_current_revision=$(cc_digest "$cc_path")
    if test "$cc_current_revision" = "$cc_revision"; then printf 'UNCHANGED: %s\n' "$cc_path"; else printf 'RELOAD: %s\n' "$cc_path"; fi
  done
}

cc_validate_delegation() {
  cc_packet=$1
  for cc_field in schema_version wrapper_version session_id parent_session_id root_session_id role objective scope non_goals context_set context_receipt repository worktree permissions acceptance_criteria verification expected_evidence stop_conditions handoff_schema; do
    grep -E "^$cc_field:|^  $cc_field:" "$cc_packet" >/dev/null 2>&1 || return 1
  done
  if grep -F 'role: verifier' "$cc_packet" >/dev/null 2>&1; then
    grep -F 'write_worktree: false' "$cc_packet" >/dev/null 2>&1 || return 1
    grep -F 'write_plan: false' "$cc_packet" >/dev/null 2>&1 || return 1
  fi
}

cc_scope_allows() {
  case "$1/" in "$2/"*) return 0 ;; *) return 1 ;; esac
}

cc_stack_join() {
  cc_stack_values=$(printf '%s\n' "$@" | sort | tr '\n' ' ')
  if command -v sha256sum >/dev/null 2>&1; then printf '%s' "$cc_stack_values" | sha256sum | awk '{print "sha256:" $1}'; else printf '%s' "$cc_stack_values" | cksum | awk '{print "cksum:" $1}'; fi
}

cc_upgrade_classify() {
  case "$1" in
    legacy-unknown|'') printf '%s\n' legacy-unknown ;;
    1.0.*) printf '%s\n' compatible ;;
    0.*) printf '%s\n' migration-needed ;;
    *) printf '%s\n' blocked ;;
  esac
}

cc_confirmation_card() {
  cat <<EOF
Action: $1
Target: $2
Observed state: $3
Will change: $4
Will not change: $5
Risks/open decisions: $6
Confirmation requested: Confirm this named action in the current session.
EOF
}

cc_approval_card() {
  cc_card_plan=$1
  cat <<EOF
Action: present-approval-card
Target: $cc_card_plan
Observed state: current session; nothing has changed yet; plan remains draft and included tasks remain draft
Will change after confirmation: plan.yaml status draft→approved; included task projections draft→ready
Will not change: Git, leases, worktrees, execution, delivery, publication, or runtime state
Risks/open decisions: confirmation does not commit Git and does not start Run approved plan
Confirmation requested: Confirm approval of plan $cc_card_plan.
EOF
}

cc_commit_approved_plan_card() {
  cc_card_plan=$1
  cat <<EOF
Action: commit-approved-plan
Target: $cc_card_plan
Observed state: exact approval projection is the only dirty source change
Will change: commit the listed plan.yaml and task frontmatter status changes
Will not change: implementation files, registered repositories, runtime state,
  external remotes, or delivery state
Risks/open decisions: this is a maintainer-source commit; inspect the file list
Confirmation requested: Confirm commit of the approved plan state.
EOF
}

cc_confirm_approval() {
  cc_workspace_root=$1
  cc_plan_file=$2
  cc_task_dir=$3
  cc_plan_id=$4
  cc_confirmation=${5:-}
  cc_transition_plan_status "$cc_plan_file" "$cc_task_dir" approved "$cc_confirmation" || return 1
  printf '%s\n' 'Approval is complete. Execution has not started.'
  if cc_is_product_source "$cc_workspace_root"; then
    cc_follow_on=$(cc_maintainer_approval_commit_required "$cc_workspace_root" "$cc_plan_file" "$cc_task_dir") || cc_follow_on=DIRTY_BASE_BLOCKED
    if test "$cc_follow_on" = MAINTAINER_APPROVAL_COMMIT_REQUIRED; then
      cc_commit_approved_plan_card "$cc_plan_id"
      printf '%s\n' "$cc_follow_on"
    else
      printf '%s\n' DIRTY_BASE_BLOCKED
    fi
  else
    printf 'Next action: Run approved plan %s\n' "$cc_plan_id"
  fi
}

cc_cleanup_disposition() {
  if [ "${1:-clean}" = dirty ] || [ "${2:-clean}" = unpushed ]; then
    printf '%s\n' CLEANUP_REQUIRES_DISCARD_CONFIRMATION
  else
    printf '%s\n' CLEANUP_REQUIRES_HUMAN_CONFIRMATION
  fi
}

cc_select_sources() {
  for cc_source in "$@"; do
    case "$cc_source" in
      sources/*) cc_safe_relative "$cc_source" || return 1; printf '%s\n' "$cc_source" ;;
      *) printf '%s\n' SOURCE_SELECTION_REQUIRED >&2; return 1 ;;
    esac
  done
}

cc_archive_latest() {
  cc_archive_file=$1
  test -f "$cc_archive_file" || { printf '%s\n' active; return 0; }
  awk '/^  - event: / { latest=$3 } END { if (latest != "") print latest; else print invalid }' "$cc_archive_file"
}

cc_archive_eligible() {
  cc_archive_file=$1
  case "$(cc_archive_latest "$cc_archive_file")" in archived) printf '%s\n' blocked-archived ;; restored|active) printf '%s\n' eligible ;; *) printf '%s\n' blocked-invalid ;; esac
}

cc_takeover_lease() {
  cc_runtime=$1; cc_plan=$2; cc_old=$3; cc_new=$4; cc_reason=$5; cc_confirmation=${6:-}
  test "$cc_confirmation" = confirmed || { printf '%s\n' TAKEOVER_CONFIRMATION_REQUIRED >&2; return 1; }
  cc_owner="$cc_runtime/plans/$cc_plan/lease.lock/owner.yaml"
  test -f "$cc_owner" || return 1
  grep -F "session_id: $cc_old" "$cc_owner" >/dev/null 2>&1 || return 1
  cc_atomic_write "$cc_runtime/plans/$cc_plan/takeover.yaml" \
    'schema_version: 1' "plan: $cc_plan" "previous_owner: $cc_old" \
    "new_owner: $cc_new" "reason: $cc_reason" 'confirmed: true' \
    'timestamp: 2026-08-21T00:00:00Z'
}

cc_cleanup_inspect() {
  cc_runtime=$1
  test -d "$cc_runtime/worktrees" || { printf '%s\n' NO_RUNTIME_WORKTREES; return 0; }
  cc_found=0
  for cc_marker in "$cc_runtime"/worktrees/*/*/.dirty "$cc_runtime"/worktrees/*/*/.unpushed; do
    test -e "$cc_marker" || continue
    cc_found=1
    printf 'DIRTY_OR_UNPUSHED: %s\n' "$cc_marker"
  done
  if test "$cc_found" -eq 1; then printf '%s\n' CLEANUP_REQUIRES_DISCARD_CONFIRMATION; else printf '%s\n' CLEANUP_REQUIRES_HUMAN_CONFIRMATION; fi
}

cc_provider_fallback() {
  case "$1" in
    unavailable) printf '%s\n' offline-fallback ;;
    denied) printf '%s\n' provider-denied-manual-fallback ;;
    disabled) printf '%s\n' provider-disabled-manual-fallback ;;
    *) printf '%s\n' provider-state-not-canonical ;;
  esac
}

# Repository bootstrap primitives. These functions deliberately use a small
# YAML subset matching workspace.yaml and repositories.local.yaml; hosts own
# richer parsing while this layer enforces the safety boundary.

cc_repository_field() {
  cc_repository_file=$1
  cc_repository_key=$2
  cc_repository_field_name=$3
  test -f "$cc_repository_file" || return 1
  cc_safe_id "$cc_repository_key" || return 1
  awk -v key="$cc_repository_key" -v field="$cc_repository_field_name" '
    /^repositories:[[:space:]]*$/ { in_repositories=1; next }
    /^[^[:space:]]/ { in_repositories=0; in_key=0 }
    in_repositories && $0 ~ "^  " key ":[[:space:]]*$" { in_key=1; next }
    in_repositories && in_key && $0 ~ /^  [A-Za-z0-9._-]+:[[:space:]]*$/ { in_key=0 }
    in_key && $0 ~ "^    " field ":[[:space:]]*" {
      sub("^    " field ":[[:space:]]*", "")
      print
      exit
    }
  ' "$cc_repository_file"
}

cc_repository_remote_safe() {
  case "$1" in
    ''|-[!-]*|*' '*|*'	'*|*password*|*PASSWORD*|*token*|*TOKEN*|*secret*|*SECRET*|*api_key*|*API_KEY*) return 1 ;;
    *://*'@'*) return 1 ;;
  esac
  return 0
}

cc_repository_canonical_safe() {
  test -n "$1" || return 0
  cc_repository_remote_safe "$1" || return 1
  case "$1" in
    https://*|ssh://*|git@*|file://*) return 0 ;;
    *) return 1 ;;
  esac
}

cc_repository_path_safe() {
  cc_repository_path=$1
  test -d "$cc_repository_path" || return 1
  case "$cc_repository_path" in *' '*|*'	'*) return 1 ;; esac
  cc_repository_real=$(realpath -e -- "$cc_repository_path" 2>/dev/null) || return 1
  test "$cc_repository_real" = "$cc_repository_path"
}

cc_resolve_repository_binding() {
  cc_workspace_root=$1
  cc_repository_key=$2
  cc_binding_file=${3:-$cc_workspace_root/repositories.local.yaml}
  cc_safe_id "$cc_repository_key" || { printf '%s\n' INVALID_REPOSITORY_KEY >&2; return 1; }
  cc_workspace_root=$(CDPATH= cd -- "$cc_workspace_root" 2>/dev/null && pwd -P) || { printf '%s\n' WORKSPACE_NOT_FOUND >&2; return 1; }
  test -f "$cc_binding_file" || { printf '%s\n' BINDING_MISSING >&2; return 1; }
  test ! -L "$cc_binding_file" || { printf '%s\n' UNSAFE_BINDING_SYMLINK >&2; return 1; }
  cc_repository_binding=$(cc_repository_field "$cc_binding_file" "$cc_repository_key" path)
  test -n "$cc_repository_binding" || { printf '%s\n' BINDING_PATH_MISSING >&2; return 1; }
  case "$cc_repository_binding" in
    /*) cc_repository_candidate=$cc_repository_binding ;;
    *) cc_safe_relative "$cc_repository_binding" || { printf '%s\n' UNSAFE_BINDING_PATH >&2; return 1; }
       cc_repository_candidate="$cc_workspace_root/$cc_repository_binding" ;;
  esac
  cc_repository_candidate=${cc_repository_candidate%/}
  cc_repository_path_safe "$cc_repository_candidate" || { printf '%s\n' UNSAFE_OR_MISSING_REPOSITORY_PATH >&2; return 1; }
  printf '%s\n' "$cc_repository_candidate"
}

cc_repository_identity_token() {
  cc_repository_identity=$1
  case "$cc_repository_identity" in
    *://*) cc_repository_identity=${cc_repository_identity#*://} ;;
    *@*:*/*) cc_repository_identity=${cc_repository_identity#*@}; cc_repository_identity=$(printf '%s' "$cc_repository_identity" | sed 's/:/\//') ;;
  esac
  cc_repository_identity=${cc_repository_identity%/}
  cc_repository_identity=${cc_repository_identity%.git}
  cc_lower "$cc_repository_identity"
}

cc_repository_identity_matches() {
  test -n "$1" || return 0
  test -n "$2" || return 1
  test "$(cc_repository_identity_token "$1")" = "$(cc_repository_identity_token "$2")"
}

cc_repository_is_clean() {
  cc_repository_clean_root=$1
  test -z "$(git -C "$cc_repository_clean_root" status --porcelain --untracked-files=all 2>/dev/null)"
}

cc_validate_repository_binding() {
  cc_workspace_root=$1
  cc_repository_key=$2
  cc_binding_file=${3:-$cc_workspace_root/repositories.local.yaml}
  cc_repository_path=$(cc_resolve_repository_binding "$cc_workspace_root" "$cc_repository_key" "$cc_binding_file") || return 1
  cc_shared_file="$cc_workspace_root/workspace.yaml"
  test -f "$cc_shared_file" || { printf '%s\n' SHARED_WORKSPACE_MISSING >&2; return 1; }
  cc_canonical_url=$(cc_repository_field "$cc_shared_file" "$cc_repository_key" canonical_url)
  cc_repository_canonical_safe "$cc_canonical_url" || { printf '%s\n' UNSAFE_CANONICAL_URL >&2; return 1; }
  cc_selected_remote=$(cc_repository_field "$cc_binding_file" "$cc_repository_key" remote)
  test -z "$cc_selected_remote" || cc_repository_remote_safe "$cc_selected_remote" || { printf '%s\n' UNSAFE_REPOSITORY_REMOTE >&2; return 1; }
  cc_origin=$(git -C "$cc_repository_path" remote get-url origin 2>/dev/null || true)
  if test -n "$cc_canonical_url"; then
    test -n "$cc_origin" || { printf '%s\n' REPOSITORY_REMOTE_MISSING >&2; return 1; }
    cc_repository_identity_matches "$cc_canonical_url" "$cc_origin" || { printf '%s\n' REPOSITORY_IDENTITY_MISMATCH >&2; return 1; }
  fi
  if test -n "$cc_selected_remote"; then
    test -n "$cc_origin" || { printf '%s\n' REPOSITORY_REMOTE_MISSING >&2; return 1; }
    cc_repository_identity_matches "$cc_selected_remote" "$cc_origin" || { printf '%s\n' REPOSITORY_REMOTE_MISMATCH >&2; return 1; }
  fi
  git -C "$cc_repository_path" rev-parse --show-toplevel >/dev/null 2>&1 || { printf '%s\n' NOT_A_GIT_REPOSITORY >&2; return 1; }
  cc_repository_is_clean "$cc_repository_path" || { printf '%s\n' DIRTY_SOURCE_BLOCKED >&2; return 1; }
  printf '%s\n' "$cc_repository_path"
}

cc_record_repository_binding_evidence() {
  cc_evidence_file=$1
  cc_repository_key=$2
  cc_repository_path=$3
  cc_repository_status=$4
  cc_repository_reason=${5:-none}
  cc_atomic_write "$cc_evidence_file" \
    'schema_version: 1' "repository: $cc_repository_key" \
    "path: $cc_repository_path" "status: $cc_repository_status" \
    "reason: $cc_repository_reason" \
    'credentials: never-recorded' 'source_preserved: true'
}

cc_prepare_bound_worktree() {
  cc_workspace_root=$1
  cc_repository_key=$2
  cc_plan_id=$3
  cc_branch=${4:-}
  cc_safe_id "$cc_repository_key" && cc_safe_id "$cc_plan_id" || { printf '%s\n' INVALID_WORKTREE_ID >&2; return 1; }
  cc_repository_path=$(cc_validate_repository_binding "$cc_workspace_root" "$cc_repository_key") || return 1
  test -n "$cc_branch" || cc_branch=$(cc_repository_field "$cc_workspace_root/workspace.yaml" "$cc_repository_key" default_branch)
  test -n "$cc_branch" || cc_branch=main
  cc_safe_id "$cc_branch" || { printf '%s\n' INVALID_DEFAULT_BRANCH >&2; return 1; }
  cc_worktree_target="$cc_workspace_root/.runtime/worktrees/$cc_repository_key/$cc_plan_id"
  cc_prepare_worktree "$cc_repository_path" "$cc_worktree_target" "$cc_branch" || return 1
  printf '%s\n' "$cc_worktree_target"
}

cc_repository_bootstrap_card() {
  cc_workspace_root=$1
  cc_repository_key=$2
  cc_canonical_url=${3:-$(cc_repository_field "$cc_workspace_root/workspace.yaml" "$cc_repository_key" canonical_url)}
  cc_selected_remote=$4
  cc_branch=${5:-$(cc_repository_field "$cc_workspace_root/workspace.yaml" "$cc_repository_key" default_branch)}
  cc_destination=$6
  cc_existing_check=${7:-not-checked}
  printf 'Action: repository-bootstrap\n'
  printf 'Repository: %s\n' "$cc_repository_key"
  printf 'Canonical URL: %s\n' "${cc_canonical_url:-none}"
  printf 'Selected remote: %s\n' "$cc_selected_remote"
  printf 'Branch: %s\n' "${cc_branch:-main}"
  printf 'Destination: %s\n' "$cc_destination"
  printf 'Existing-path check: %s\n' "$cc_existing_check"
  printf 'Will change: create the exact destination and clone the selected remote after confirmation.\n'
  printf 'Will not change: existing paths, bound repositories, credentials, or delivery state.\n'
  printf 'Risks/open decisions: host Git credentials and provider availability remain external.\n'
  printf 'Confirmation requested: Confirm repository-bootstrap for this exact target in the current session.\n'
}

cc_bootstrap_repository() {
  cc_workspace_root=$1
  cc_repository_key=$2
  cc_selected_remote=$3
  cc_branch=${4:-}
  cc_destination=$5
  cc_confirmation=${6:-}
  cc_safe_id "$cc_repository_key" || { printf '%s\n' INVALID_REPOSITORY_KEY >&2; return 1; }
  test -n "$cc_selected_remote" || { printf '%s\n' REMOTE_REQUIRED >&2; return 1; }
  cc_repository_remote_safe "$cc_selected_remote" || { printf '%s\n' UNSAFE_REPOSITORY_REMOTE >&2; return 1; }
  test -n "$cc_branch" || cc_branch=$(cc_repository_field "$cc_workspace_root/workspace.yaml" "$cc_repository_key" default_branch)
  test -n "$cc_branch" || cc_branch=main
  cc_safe_id "$cc_branch" || { printf '%s\n' INVALID_DEFAULT_BRANCH >&2; return 1; }
  case "$cc_destination" in
    /*) cc_bootstrap_target=$cc_destination ;;
    *) cc_safe_relative "$cc_destination" || { printf '%s\n' UNSAFE_BOOTSTRAP_DESTINATION >&2; return 1; }
       cc_bootstrap_target="$cc_workspace_root/$cc_destination" ;;
  esac
  cc_bootstrap_target=${cc_bootstrap_target%/}
  cc_existing_check=absent
  test ! -e "$cc_bootstrap_target" && test ! -L "$cc_bootstrap_target" || cc_existing_check=exists
  cc_repository_bootstrap_card "$cc_workspace_root" "$cc_repository_key" "$(cc_repository_field "$cc_workspace_root/workspace.yaml" "$cc_repository_key" canonical_url)" "$cc_selected_remote" "$cc_branch" "$cc_bootstrap_target" "$cc_existing_check"
  test "$cc_confirmation" = confirmed || { printf '%s\n' BOOTSTRAP_CONFIRMATION_REQUIRED >&2; return 1; }
  test "$cc_existing_check" = absent || { printf '%s\n' BOOTSTRAP_TARGET_EXISTS >&2; return 1; }
  if test "${CC_OFFLINE:-0}" = 1; then
    cc_record_repository_binding_evidence "$cc_workspace_root/.runtime/bootstrap/$cc_repository_key.yaml" "$cc_repository_key" "$cc_bootstrap_target" offline OFFLINE_PROVIDER
    printf '%s\n' offline-fallback
    return 1
  fi
  cc_bootstrap_parent=$(dirname -- "$cc_bootstrap_target")
  if test -e "$cc_bootstrap_parent"; then
    cc_repository_path_safe "$cc_bootstrap_parent" || { printf '%s\n' UNSAFE_BOOTSTRAP_PARENT >&2; return 1; }
  else
    mkdir -p -- "$cc_bootstrap_parent" || return 1
  fi
  if GIT_TERMINAL_PROMPT=0 git clone --branch "$cc_branch" -- "$cc_selected_remote" "$cc_bootstrap_target"; then
    cc_record_repository_binding_evidence "$cc_workspace_root/.runtime/bootstrap/$cc_repository_key.yaml" "$cc_repository_key" "$cc_bootstrap_target" cloned clone-succeeded
    printf '%s\n' "$cc_bootstrap_target"
  else
    cc_record_repository_binding_evidence "$cc_workspace_root/.runtime/bootstrap/$cc_repository_key.yaml" "$cc_repository_key" "$cc_bootstrap_target" failed clone-failed
    printf '%s\n' BOOTSTRAP_FAILED >&2
    return 1
  fi
}

# Stable aliases used by host adapters and semantic fixtures.
cc_resolve_repository() { cc_resolve_repository_binding "$@"; }
cc_validate_binding() { cc_validate_repository_binding "$@"; }
cc_repository_bootstrap() { cc_bootstrap_repository "$@"; }
