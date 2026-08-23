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
    *'review plan'*|*'plan ready'*|*'plan'*'ready'*|*'approved-for-execution'*|*'read only the plan status'*|*'walk me through'*'plan'*) printf '%s\n' plan-review ;;
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
    plan-review)
      case "$cc_action_text" in
        *'this plan'*|*'plan status'*) printf '%s\n' clarify-target ;;
        *) printf '%s\n' review-plan ;;
      esac
      ;;
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

# The route contract owns action-to-packet selection.  A fixture or live
# workspace without that owner is a blocked boundary, never an invitation to
# maintain a second route table here.
cc_route_context_set() {
  cc_route_action=$1
  cc_route_root=${2:-.}
  cc_route_contract_root=$cc_route_root
  cc_route_contract="$cc_route_contract_root/wrapper/contracts/routes.yaml"
  test -f "$cc_route_contract" || { printf '%s\n' CONTEXT_ROUTE_OWNER_MISSING >&2; return 1; }
  cc_route_mapped=""
  cc_route_mapped=$(awk -v wanted="$cc_route_action" '
    /^context_set_map:/ { in_map=1; next }
    in_map && /^[^[:space:]]/ { exit }
    in_map && /^    [^:]+:/ {
      key=$1
      sub(":$", "", key)
      if (key == wanted) {
        value=$2
        print value
        exit
      }
    }
  ' "$cc_route_contract")
  test -n "$cc_route_mapped" || { printf 'CONTEXT_ROUTE_UNMAPPED: %s\n' "$cc_route_action" >&2; return 1; }
  cc_context_set_registered "$cc_route_contract_root" "$cc_route_mapped" || {
    printf 'CONTEXT_SET_UNREGISTERED: %s\n' "$cc_route_mapped" >&2
    return 1
  }
  printf '%s\n' "$cc_route_mapped"
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
  if cc_identity_projection_applies "$cc_workspace_root"; then
    if ! cc_validate_identity_projection "$cc_workspace_root" >/dev/null; then
      cc_stage_b=block-projection-mismatch
      cc_eligibility=blocked
      cc_gate=none
      cc_reason=projection-mismatch
      cc_authorization=absent
    fi
  fi
  cc_context_set=$(cc_route_context_set "$cc_stage_b" "$cc_workspace_root") || return 1
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
context_set: $cc_context_set
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
  cc_status_dest=$2
  cc_status_value=$3
  case "$cc_status_value" in
    draft|approved|ready|done) ;;
    *) return 1 ;;
  esac
  test -f "$cc_status_file" || return 1
  test -n "$cc_status_dest" || return 1
  cc_had_nl=$(tail -c 1 "$cc_status_file" | wc -l)
  awk -v status="$cc_status_value" '
    BEGIN { replaced=0 }
    {
      if (!replaced && $0 ~ /^status: /) { print "status: " status; replaced=1 }
      else print
    }
    END { exit replaced ? 0 : 1 }
  ' "$cc_status_file" > "$cc_status_dest" || { rm -f "$cc_status_dest"; return 1; }
  if test "$cc_had_nl" -eq 0; then
    cc_size=$(wc -c < "$cc_status_dest" | tr -d ' \t\n')
    if test "$cc_size" -gt 0; then
      dd if="$cc_status_dest" of="$cc_status_dest.raw" bs=1 count=$((cc_size - 1)) 2>/dev/null || {
        rm -f "$cc_status_dest" "$cc_status_dest.raw"
        return 1
      }
      mv "$cc_status_dest.raw" "$cc_status_dest"
    fi
  fi
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
  cc_temps=""
  cc_plan_tmp="$cc_plan_file.tmp.$$"
  if ! cc_replace_first_status "$cc_plan_file" "$cc_plan_tmp" "$cc_target"; then
    printf '%s\n' STATUS_LINE_MISSING >&2
    rm -f "$cc_plan_tmp"
    return 1
  fi
  cc_temps="$cc_plan_tmp"
  cc_task_list="$cc_plan_file.list.$$"
  : > "$cc_task_list"
  for cc_task in "$cc_task_dir"/*.md; do
    test -f "$cc_task" || continue
    cc_task_tmp="$cc_task.tmp.$$"
    if ! cc_replace_first_status "$cc_task" "$cc_task_tmp" "$cc_expected"; then
      printf '%s\n' STATUS_LINE_MISSING >&2
      for cc_tmp in $cc_temps; do
        rm -f "$cc_tmp" "$cc_tmp.raw"
      done
      rm -f "$cc_task_tmp" "$cc_task_list"
      return 1
    fi
    cc_temps="$cc_temps $cc_task_tmp"
    printf '%s\t%s\n' "$cc_task_tmp" "$cc_task" >> "$cc_task_list"
  done
  mv "$cc_plan_tmp" "$cc_plan_file"
  while IFS='	' read -r cc_task_tmp cc_task; do
    test -n "$cc_task_tmp" || continue
    mv "$cc_task_tmp" "$cc_task"
  done < "$cc_task_list"
  rm -f "$cc_task_list"
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

# Bounded context packet primitives.  The route contract selects the packet;
# this host-neutral library reads only the packet's declared paths, measures
# actual bytes, and records locators in a receipt.  It never broadens a packet
# to satisfy a missing file or a host request.
cc_digest_text() {
  cc_text=$1
  if command -v sha256sum >/dev/null 2>&1; then
    printf '%s' "$cc_text" | sha256sum | awk '{print "sha256:" $1}'
  else
    printf '%s' "$cc_text" | cksum | awk '{print "cksum:" $1}'
  fi
}

cc_context_set_contract() {
  cc_context_root=$1
  cc_context_file="$cc_context_root/wrapper/contracts/context-sets.yaml"
  test -f "$cc_context_file" || { printf '%s\n' CONTEXT_SET_CONTRACT_MISSING >&2; return 1; }
  test ! -L "$cc_context_file" || { printf '%s\n' CONTEXT_SET_CONTRACT_SYMLINK >&2; return 1; }
  printf '%s\n' "$cc_context_file"
}

cc_context_set_registered() {
  cc_context_root=$1
  cc_context_set=$2
  cc_safe_id "$cc_context_set" || return 1
  cc_context_file=$(cc_context_set_contract "$cc_context_root") || return 1
  awk -v wanted="$cc_context_set" '
    /^  - id:[[:space:]]*/ { if ($3 == wanted) count++ }
    END { exit count == 1 ? 0 : 1 }
  ' "$cc_context_file"
}

cc_context_set_budget() {
  cc_context_root=$1
  cc_context_set=$2
  cc_context_file=$(cc_context_set_contract "$cc_context_root") || return 1
  awk -v wanted="$cc_context_set" '
    /^  - id:[[:space:]]*/ { in_set=($3 == wanted); next }
    in_set && /^    budget_bytes:[[:space:]]*/ { print $2; exit }
  ' "$cc_context_file"
}

cc_context_set_paths() {
  cc_context_root=$1
  cc_context_set=$2
  cc_context_field=$3
  cc_context_file=$(cc_context_set_contract "$cc_context_root") || return 1
  case "$cc_context_field" in allowlist|conditional_allowlist) ;; *) return 1 ;; esac
  awk -v wanted="$cc_context_set" -v field="$cc_context_field" '
    /^  - id:[[:space:]]*/ { in_set=($3 == wanted); in_list=0; next }
    in_set && $0 ~ "^    " field ":[[:space:]]*$" { in_list=1; next }
    in_set && /^    [A-Za-z_][A-Za-z0-9_-]*:/ { in_list=0 }
    in_set && in_list && /^      -[[:space:]]*/ {
      sub("^      -[[:space:]]*", "")
      print
    }
  ' "$cc_context_file"
}

cc_context_path_pattern_matches() {
  cc_context_pattern=$1
  cc_context_path=$2
  case "$cc_context_pattern" in
    sources/\<selected-path\>)
      case "$cc_context_path" in sources/*) test "${cc_context_path#sources/}" != "" ;; *) return 1 ;; esac
      ;;
    context/domains/\<domain\>/README.md)
      case "$cc_context_path" in
        context/domains/*/README.md)
          cc_context_domain=${cc_context_path#context/domains/}
          cc_context_domain=${cc_context_domain%/README.md}
          cc_safe_id "$cc_context_domain"
          ;;
        *) return 1 ;;
      esac
      ;;
    context/roles/\<role\>.md)
      case "$cc_context_path" in
        context/roles/*.md)
          cc_context_role=${cc_context_path#context/roles/}
          cc_context_role=${cc_context_role%.md}
          cc_safe_id "$cc_context_role"
          ;;
        *) return 1 ;;
      esac
      ;;
    plans/\<repository\>-plans/\<plan\>/plan.yaml)
      case "$cc_context_path" in plans/*-plans/*/plan.yaml)
        cc_context_plan=${cc_context_path#plans/}; cc_context_plan=${cc_context_plan%/plan.yaml}
        cc_context_repo=${cc_context_plan%%/*}; cc_context_name=${cc_context_plan#*/}
        cc_context_repo=${cc_context_repo%-plans}
        cc_safe_id "$cc_context_repo" && cc_safe_id "$cc_context_name"
        ;; *) return 1 ;; esac
      ;;
    plans/\<repository\>-plans/\<plan\>/PLAN.md)
      case "$cc_context_path" in plans/*-plans/*/PLAN.md)
        cc_context_plan=${cc_context_path#plans/}; cc_context_plan=${cc_context_plan%/PLAN.md}
        cc_context_repo=${cc_context_plan%%/*}; cc_context_name=${cc_context_plan#*/}
        cc_context_repo=${cc_context_repo%-plans}
        cc_safe_id "$cc_context_repo" && cc_safe_id "$cc_context_name"
        ;; *) return 1 ;; esac
      ;;
    plans/\<repository\>-plans/\<plan\>/tasks/\<task\>.md)
      case "$cc_context_path" in plans/*-plans/*/tasks/*.md)
        cc_context_task=${cc_context_path#plans/}; cc_context_task=${cc_context_task%/tasks/*}
        cc_context_file=${cc_context_path##*/}; cc_context_file=${cc_context_file%.md}
        cc_context_repo=${cc_context_task%%/*}; cc_context_name=${cc_context_task#*/}
        cc_context_repo=${cc_context_repo%-plans}
        cc_safe_id "$cc_context_repo" && cc_safe_id "$cc_context_name" && cc_safe_id "$cc_context_file"
        ;; *) return 1 ;; esac
      ;;
    .runtime/plans/\<plan\>/lease.yaml|.runtime/plans/\<plan\>/completion.yaml)
      case "$cc_context_path" in .runtime/plans/*/lease.yaml|.runtime/plans/*/completion.yaml)
        cc_context_plan=${cc_context_path#*.runtime/plans/}
        cc_context_plan=${cc_context_plan%/lease.yaml}
        cc_context_plan=${cc_context_plan%/completion.yaml}
        cc_safe_id "$cc_context_plan" ;; *) return 1 ;; esac
      ;;
    .runtime/plans/\<plan\>/lease.lock/owner.yaml)
      case "$cc_context_path" in .runtime/plans/*/lease.lock/owner.yaml)
        cc_context_plan=${cc_context_path#*.runtime/plans/}; cc_context_plan=${cc_context_plan%/lease.lock/owner.yaml}
        cc_safe_id "$cc_context_plan" ;; *) return 1 ;; esac
      ;;
    .runtime/sessions/\<session-id\>/*)
      case "$cc_context_path" in .runtime/sessions/*/*)
        cc_context_session=${cc_context_path#*.runtime/sessions/}; cc_context_session=${cc_context_session%%/*}
        cc_safe_id "$cc_context_session" ;; *) return 1 ;; esac
      ;;
    *) test "$cc_context_pattern" = "$cc_context_path" ;;
  esac
}

cc_context_path_allowed() {
  cc_context_root=$1
  cc_context_set=$2
  cc_context_path=$3
  cc_safe_relative "$cc_context_path" || return 1
  cc_context_set_registered "$cc_context_root" "$cc_context_set" || return 1
  cc_context_set_paths "$cc_context_root" "$cc_context_set" allowlist | grep -Fx "$cc_context_path" >/dev/null 2>&1 && return 0
  while IFS= read -r cc_context_pattern; do
    test -n "$cc_context_pattern" || continue
    cc_context_path_pattern_matches "$cc_context_pattern" "$cc_context_path" && return 0
  done <<EOF
$(cc_context_set_paths "$cc_context_root" "$cc_context_set" conditional_allowlist)
EOF
  return 1
}

cc_context_path_validate() {
  cc_context_root=$1
  cc_context_set=$2
  cc_context_path=$3
  cc_context_path_allowed "$cc_context_root" "$cc_context_set" "$cc_context_path" || {
    printf 'CONTEXT_PATH_UNDECLARED: %s\n' "$cc_context_path" >&2
    return 1
  }
  case "$cc_context_path" in
    *credentials*|*.credentials|*.env|*secret*|*token*) printf 'CONTEXT_SENSITIVE_PATH\n' >&2; return 1 ;;
  esac
  cc_context_root=$(CDPATH= cd -- "$cc_context_root" 2>/dev/null && pwd -P) || return 1
  cc_context_target="$cc_context_root/$cc_context_path"
  test -f "$cc_context_target" || { printf 'CONTEXT_PATH_MISSING: %s\n' "$cc_context_path" >&2; return 1; }
  test ! -L "$cc_context_target" || { printf 'CONTEXT_UNSAFE_SYMLINK: %s\n' "$cc_context_path" >&2; return 1; }
  cc_context_real=$(realpath -e -- "$cc_context_target" 2>/dev/null) || return 1
  test "$cc_context_real" = "$cc_context_target" || { printf 'CONTEXT_UNSAFE_PATH: %s\n' "$cc_context_path" >&2; return 1; }
}

cc_context_packet_measure() {
  cc_context_root=$1
  cc_context_set=$2
  shift 2
  cc_context_set_registered "$cc_context_root" "$cc_context_set" || { printf 'CONTEXT_SET_UNREGISTERED: %s\n' "$cc_context_set" >&2; return 1; }
  cc_context_budget=$(cc_context_set_budget "$cc_context_root" "$cc_context_set")
  case "$cc_context_budget" in ''|*[!0-9]*) printf 'CONTEXT_BUDGET_MISSING: %s\n' "$cc_context_set" >&2; return 1 ;; esac
  cc_context_list=$(mktemp "${TMPDIR:-/tmp}/cc-context-paths.XXXXXX") || return 1
  cc_context_records=$(mktemp "${TMPDIR:-/tmp}/cc-context-records.XXXXXX") || { rm -f "$cc_context_list"; return 1; }
  if test "$#" -eq 0; then
    cc_context_set_paths "$cc_context_root" "$cc_context_set" allowlist > "$cc_context_list"
  else
    for cc_context_path in "$@"; do
      test -n "$cc_context_path" || { rm -f "$cc_context_list" "$cc_context_records"; printf 'CONTEXT_PATH_EMPTY\n' >&2; return 1; }
      if grep -Fx "$cc_context_path" "$cc_context_list" >/dev/null 2>&1; then
        rm -f "$cc_context_list" "$cc_context_records"
        printf 'CONTEXT_PATH_DUPLICATE: %s\n' "$cc_context_path" >&2
        return 1
      fi
      printf '%s\n' "$cc_context_path" >> "$cc_context_list"
    done
  fi
  cc_context_actual=0
  while IFS= read -r cc_context_path; do
    test -n "$cc_context_path" || continue
    cc_context_path_validate "$cc_context_root" "$cc_context_set" "$cc_context_path" || { rm -f "$cc_context_list" "$cc_context_records"; return 1; }
    cc_context_bytes=$(wc -c < "$cc_context_root/$cc_context_path" | tr -d ' \t\n')
    cc_context_revision=$(cc_digest "$cc_context_root/$cc_context_path")
    cc_context_actual=$((cc_context_actual + cc_context_bytes))
    printf '%s|%s|%s\n' "$cc_context_path" "$cc_context_revision" "$cc_context_bytes" >> "$cc_context_records"
  done < "$cc_context_list"
  if test "$cc_context_actual" -gt "$cc_context_budget"; then
    printf 'CONTEXT_BUDGET_EXCEEDED selected_set=%s actual_bytes=%s budget_bytes=%s newly_requested_evidence=%s\n' \
      "$cc_context_set" "$cc_context_actual" "$cc_context_budget" "${*:-none}" >&2
    rm -f "$cc_context_list" "$cc_context_records"
    return 1
  fi
  cat "$cc_context_records"
  rm -f "$cc_context_list" "$cc_context_records"
}

cc_load_context_packet() {
  cc_context_root=$1
  cc_context_set=$2
  cc_context_receipt=$3
  cc_context_route_digest=${4:-sha256:unavailable}
  cc_context_session=${5:-context-packet}
  shift 5
  cc_safe_id "$cc_context_session" || { printf 'INVALID_CONTEXT_SESSION\n' >&2; return 1; }
  test -n "$cc_context_receipt" || { printf 'CONTEXT_RECEIPT_REQUIRED\n' >&2; return 1; }
  test ! -L "$cc_context_receipt" || { printf 'CONTEXT_RECEIPT_SYMLINK\n' >&2; return 1; }
  cc_context_load_records=$(mktemp "${TMPDIR:-/tmp}/cc-context-load.XXXXXX") || return 1
  if ! cc_context_packet_measure "$cc_context_root" "$cc_context_set" "$@" > "$cc_context_load_records"; then
    rm -f "$cc_context_load_records"
    return 1
  fi
  cc_context_budget=$(cc_context_set_budget "$cc_context_root" "$cc_context_set")
  cc_context_actual=$(awk -F'|' '{sum += $3} END {print sum + 0}' "$cc_context_load_records")
  cc_context_material=$(awk -F'|' '{printf "%s|%s|%s\n", $1, $2, $3}' "$cc_context_load_records")
  cc_context_packet_digest=$(cc_digest_text "$cc_context_material")
  cc_context_parent=$(dirname -- "$cc_context_receipt")
  mkdir -p "$cc_context_parent" || { rm -f "$cc_context_load_records"; return 1; }
  cc_context_tmp="$cc_context_receipt.tmp.$$"
  {
    printf 'schema_version: 1\n'
    printf 'wrapper_version: %s\n' "${CC_WRAPPER_VERSION:-1.0.0}"
    printf 'session_id: %s\n' "$cc_context_session"
    printf 'route_decision_digest: %s\n' "$cc_context_route_digest"
    printf 'context_set: %s\n' "$cc_context_set"
    printf 'packet_id: %s\n' "$cc_context_set"
    printf 'declared_budget_bytes: %s\n' "$cc_context_budget"
    printf 'actual_bytes: %s\n' "$cc_context_actual"
    printf 'packet_digest: %s\n' "$cc_context_packet_digest"
    printf 'selected_paths:\n'
    while IFS='|' read -r cc_context_path cc_context_revision cc_context_bytes; do
      printf '  - %s\n' "$cc_context_path"
    done < "$cc_context_load_records"
    printf 'references:\n'
    while IFS='|' read -r cc_context_path cc_context_revision cc_context_bytes; do
      printf '  - path: %s\n' "$cc_context_path"
      printf '    revision: %s\n' "$cc_context_revision"
      printf '    bytes: %s\n' "$cc_context_bytes"
    done < "$cc_context_load_records"
    printf 'invariants: [INV-CTX-01, INV-CTX-02, INV-CTX-03, INV-CTX-04, INV-CTX-05]\n'
    printf 'created_at: %s\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || printf '%s' 1970-01-01T00:00:00Z)"
    } > "$cc_context_tmp" || { rm -f "$cc_context_load_records" "$cc_context_tmp"; return 1; }
  mv "$cc_context_tmp" "$cc_context_receipt" || { rm -f "$cc_context_load_records" "$cc_context_tmp"; return 1; }
  rm -f "$cc_context_load_records"
  printf '%s\n' "$cc_context_receipt"
}

cc_validate_context_receipt() {
  cc_context_root=$1
  cc_context_receipt=$2
  cc_context_expected_set=${3:-}
  cc_validate_receipt "$cc_context_receipt" || { printf 'INVALID_CONTEXT_RECEIPT\n' >&2; return 1; }
  for cc_context_field in packet_id declared_budget_bytes actual_bytes packet_digest; do
    grep -E "^$cc_context_field:" "$cc_context_receipt" >/dev/null 2>&1 || { printf 'CONTEXT_RECEIPT_FIELD_MISSING: %s\n' "$cc_context_field" >&2; return 1; }
  done
  cc_context_set=$(sed -n 's/^context_set: //p' "$cc_context_receipt" | head -n 1)
  test -z "$cc_context_expected_set" || test "$cc_context_expected_set" = "$cc_context_set" || { printf 'CONTEXT_SET_MISMATCH\n' >&2; return 1; }
  cc_context_packet_id=$(sed -n 's/^packet_id: //p' "$cc_context_receipt" | head -n 1)
  test "$cc_context_packet_id" = "$cc_context_set" || { printf 'CONTEXT_PACKET_ID_MISMATCH\n' >&2; return 1; }
  cc_context_expected_wrapper=${4:-${CC_WRAPPER_VERSION:-1.0.0}}
  cc_context_recorded_wrapper=$(sed -n 's/^wrapper_version: //p' "$cc_context_receipt" | head -n 1)
  test "$cc_context_recorded_wrapper" = "$cc_context_expected_wrapper" || { printf 'RELOAD_WRAPPER_CONTEXT\n' >&2; return 1; }
  cc_context_set_registered "$cc_context_root" "$cc_context_set" || { printf 'CONTEXT_SET_UNREGISTERED: %s\n' "$cc_context_set" >&2; return 1; }
  cc_context_budget=$(sed -n 's/^declared_budget_bytes: //p' "$cc_context_receipt" | head -n 1)
  cc_context_actual=$(sed -n 's/^actual_bytes: //p' "$cc_context_receipt" | head -n 1)
  case "$cc_context_budget:$cc_context_actual" in *[!0-9:]*|:) printf 'CONTEXT_RECEIPT_BYTES_INVALID\n' >&2; return 1 ;; esac
  test "$cc_context_actual" -le "$cc_context_budget" || { printf 'CONTEXT_BUDGET_EXCEEDED\n' >&2; return 1; }
  cc_context_records=$(mktemp "${TMPDIR:-/tmp}/cc-context-validate.XXXXXX") || return 1
  cc_context_sum=0
  while IFS='|' read -r cc_context_path cc_context_revision cc_context_bytes; do
    test -n "$cc_context_path" || continue
    cc_context_path_validate "$cc_context_root" "$cc_context_set" "$cc_context_path" || { rm -f "$cc_context_records"; return 1; }
    if grep -F "$cc_context_path|" "$cc_context_records" >/dev/null 2>&1; then
      rm -f "$cc_context_records"
      printf 'CONTEXT_RECEIPT_DUPLICATE: %s\n' "$cc_context_path" >&2
      return 1
    fi
    cc_context_current=$(cc_digest "$cc_context_root/$cc_context_path")
    test "$cc_context_current" = "$cc_context_revision" || { rm -f "$cc_context_records"; printf 'STALE_RECEIPT: %s\n' "$cc_context_path" >&2; return 1; }
    cc_context_now_bytes=$(wc -c < "$cc_context_root/$cc_context_path" | tr -d ' \t\n')
    test "$cc_context_now_bytes" = "$cc_context_bytes" || { rm -f "$cc_context_records"; printf 'STALE_RECEIPT: %s\n' "$cc_context_path" >&2; return 1; }
    cc_context_sum=$((cc_context_sum + cc_context_bytes))
    printf '%s|%s|%s\n' "$cc_context_path" "$cc_context_revision" "$cc_context_bytes" >> "$cc_context_records"
  done <<EOF
$(awk -F': ' '/^  - path: / { path=$2 } /^    revision: / { revision=$2 } /^    bytes: / { print path "|" revision "|" $2 }' "$cc_context_receipt")
EOF
  test "$cc_context_sum" -eq "$cc_context_actual" || { rm -f "$cc_context_records"; printf 'CONTEXT_RECEIPT_TOTAL_MISMATCH\n' >&2; return 1; }
  cc_context_material=$(cat "$cc_context_records")
  cc_context_expected_digest=$(cc_digest_text "$cc_context_material")
  cc_context_recorded_digest=$(sed -n 's/^packet_digest: //p' "$cc_context_receipt" | head -n 1)
  rm -f "$cc_context_records"
  test "$cc_context_expected_digest" = "$cc_context_recorded_digest" || { printf 'CONTEXT_RECEIPT_DIGEST_MISMATCH\n' >&2; return 1; }
  printf '%s\n' context-receipt-ok
}

cc_read_context_packet() {
  cc_context_root=$1
  cc_context_set=$2
  cc_context_path=$3
  cc_context_path_validate "$cc_context_root" "$cc_context_set" "$cc_context_path" || return 1
  cat "$cc_context_root/$cc_context_path"
}

# Descriptive aliases keep host adapters thin while the implementation remains
# in this host-neutral engine library.
cc_context_packet_load() { cc_load_context_packet "$@"; }
cc_measure_context_packet() { cc_context_packet_measure "$@"; }
cc_context_packet_read() { cc_read_context_packet "$@"; }

cc_validate_delegation() {
  cc_packet=$1
  for cc_field in schema_version wrapper_version session_id parent_session_id root_session_id role objective scope non_goals context_set context_receipt repository worktree permissions acceptance_criteria verification expected_evidence stop_conditions handoff_schema; do
    grep -E "^$cc_field:|^  $cc_field:" "$cc_packet" >/dev/null 2>&1 || return 1
  done
  if grep -F 'role: verifier' "$cc_packet" >/dev/null 2>&1; then
    grep -F 'write_worktree: false' "$cc_packet" >/dev/null 2>&1 || return 1
    grep -F 'write_plan: false' "$cc_packet" >/dev/null 2>&1 || return 1
    cc_validate_verifier_readonly "$cc_packet" >/dev/null || return 1
  fi
  cc_validate_delegation_evidence "$cc_packet" >/dev/null || return 1
}

# Evidence-layer comparison implements wrapper/contracts/schemas/plan.yaml.
# Skill and role files are not a second policy owner.

cc_evidence_layer_known() {
  case "$1" in
    schema|store|api|process|browser|human) return 0 ;;
    *) return 1 ;;
  esac
}

cc_runtime_yaml_field() {
  awk -F': ' -v field="$2" '
    {
      key=$1
      sub(/^[ \t]+/, "", key)
      if (key == field) { print $2; exit }
    }
  ' "$1"
}

cc_compare_evidence_layer() {
  cc_required=$1
  cc_observed=$2
  cc_claimed=${3:-}
  cc_source=${4:-independent}
  cc_capability=${5:-available}
  if test "$cc_claimed" = waived; then
    printf '%s\n' waived
    return 0
  fi
  if test "$cc_capability" = unavailable || test "$cc_claimed" = blocked; then
    printf '%s\n' blocked
    return 0
  fi
  if test -z "$cc_observed" || test "$cc_observed" = omitted || test "$cc_observed" = missing || test "$cc_observed" = null; then
    printf '%s\n' failed
    return 0
  fi
  if ! cc_evidence_layer_known "$cc_required" || ! cc_evidence_layer_known "$cc_observed"; then
    printf '%s\n' failed
    return 0
  fi
  if test "$cc_required" != "$cc_observed"; then
    printf '%s\n' failed
    return 0
  fi
  if test "$cc_source" = writer-claim; then
    printf '%s\n' failed
    return 0
  fi
  if test -n "$cc_claimed" && test "$cc_claimed" != passed; then
    printf '%s\n' failed
    return 0
  fi
  printf '%s\n' passed
}

cc_validate_evidence_mapping() {
  cc_mapping=$1
  test -f "$cc_mapping" || { printf '%s\n' EVIDENCE_LAYER_MISSING >&2; return 1; }
  cc_required=$(cc_runtime_yaml_field "$cc_mapping" required_layer)
  cc_observed=$(cc_runtime_yaml_field "$cc_mapping" observed_layer)
  cc_claimed=$(cc_runtime_yaml_field "$cc_mapping" outcome)
  cc_source=$(cc_runtime_yaml_field "$cc_mapping" evidence_source)
  cc_capability=$(cc_runtime_yaml_field "$cc_mapping" host_capability)
  test -n "$cc_source" || cc_source=independent
  test -n "$cc_capability" || cc_capability=available
  if test -z "$cc_required"; then
    printf '%s\n' EVIDENCE_LAYER_MISSING >&2
    return 1
  fi
  cc_computed=$(cc_compare_evidence_layer "$cc_required" "$cc_observed" "$cc_claimed" "$cc_source" "$cc_capability")
  if test "$cc_source" = writer-claim; then
    printf '%s\n' EVIDENCE_WRITER_CLAIM >&2
    return 1
  fi
  case "$cc_computed" in
    passed)
      printf '%s\n' EVIDENCE_LAYER_MATCH
      return 0
      ;;
    blocked)
      printf '%s\n' EVIDENCE_LAYER_BLOCKED >&2
      return 1
      ;;
    waived)
      printf '%s\n' EVIDENCE_LAYER_WAIVED >&2
      return 1
      ;;
    *)
      if test -z "$cc_observed" || test "$cc_observed" = omitted || test "$cc_observed" = missing || test "$cc_observed" = null; then
        printf '%s\n' EVIDENCE_LAYER_MISSING >&2
      else
        printf '%s\n' EVIDENCE_LAYER_MISMATCH >&2
      fi
      return 1
      ;;
  esac
}

cc_validate_delegation_evidence() {
  cc_packet=$1
  grep -E '^evidence_layers:' "$cc_packet" >/dev/null 2>&1 || return 0
  cc_required=$(cc_runtime_yaml_field "$cc_packet" required_layer)
  cc_observed=$(cc_runtime_yaml_field "$cc_packet" observed_layer)
  test -n "$cc_required" || test -n "$cc_observed" || return 0
  if test -n "$cc_required" && test -n "$cc_observed" && test "$cc_required" != "$cc_observed"; then
    printf '%s\n' EVIDENCE_LAYER_MISMATCH >&2
    return 1
  fi
  cc_claimed=$(cc_runtime_yaml_field "$cc_packet" outcome)
  if test "$cc_claimed" = passed || test -z "$cc_claimed"; then
    cc_validate_evidence_mapping "$cc_packet" || return 1
  else
    cc_source=$(cc_runtime_yaml_field "$cc_packet" evidence_source)
    cc_capability=$(cc_runtime_yaml_field "$cc_packet" host_capability)
    test -n "$cc_source" || cc_source=independent
    test -n "$cc_capability" || cc_capability=available
    cc_computed=$(cc_compare_evidence_layer "$cc_required" "$cc_observed" "$cc_claimed" "$cc_source" "$cc_capability")
    test "$cc_computed" = passed && { printf '%s\n' EVIDENCE_LAYER_MISMATCH >&2; return 1; }
  fi
  return 0
}

cc_validate_handoff_evidence() {
  cc_handoff=$1
  test -f "$cc_handoff" || return 1
  grep -E '^evidence_layers:|^  observed_layer:|^  required_layer:' "$cc_handoff" >/dev/null 2>&1 || return 0
  cc_validate_delegation_evidence "$cc_handoff"
}

cc_validate_completion_evidence() {
  cc_completion=$1
  test -f "$cc_completion" || return 1
  cc_status=$(cc_runtime_yaml_field "$cc_completion" status)
  if grep -E '^evidence_layers:|^  required_layer:|^  observed_layer:' "$cc_completion" >/dev/null 2>&1; then
    cc_required=$(cc_runtime_yaml_field "$cc_completion" required_layer)
    cc_observed=$(cc_runtime_yaml_field "$cc_completion" observed_layer)
    if test -n "$cc_required" && test -n "$cc_observed" && test "$cc_required" != "$cc_observed"; then
      printf '%s\n' EVIDENCE_LAYER_MISMATCH >&2
      return 1
    fi
    cc_claimed=$(cc_runtime_yaml_field "$cc_completion" outcome)
    if test "$cc_claimed" = passed || test -z "$cc_claimed"; then
      cc_validate_evidence_mapping "$cc_completion" || return 1
    fi
  fi
  case "$cc_status" in
    ready-for-human-status-change|completed)
      if grep -E '^evidence_layers:|^  required_layer:|^  observed_layer:' "$cc_completion" >/dev/null 2>&1; then
        cc_validate_evidence_mapping "$cc_completion" || {
          printf '%s\n' EVIDENCE_COMPLETION_NOT_READY >&2
          return 1
        }
      else
        printf '%s\n' EVIDENCE_LAYER_MISSING >&2
        return 1
      fi
      ;;
  esac
  return 0
}

cc_validate_verifier_readonly() {
  cc_packet=$1
  grep -F 'role: verifier' "$cc_packet" >/dev/null 2>&1 || return 0
  grep -F 'write_worktree: false' "$cc_packet" >/dev/null 2>&1 || { printf '%s\n' EVIDENCE_VERIFIER_WRITE >&2; return 1; }
  grep -F 'write_plan: false' "$cc_packet" >/dev/null 2>&1 || { printf '%s\n' EVIDENCE_VERIFIER_WRITE >&2; return 1; }
  grep -F 'write_activity: false' "$cc_packet" >/dev/null 2>&1 || { printf '%s\n' EVIDENCE_VERIFIER_WRITE >&2; return 1; }
  if grep -E '^repair: true' "$cc_packet" >/dev/null 2>&1; then
    printf '%s\n' EVIDENCE_VERIFIER_REPAIR >&2
    return 1
  fi
  printf '%s\n' EVIDENCE_VERIFIER_READONLY
}

cc_runtime_emit_evidence_layers() {
  cc_request=$1
  cc_required_layer=$(cc_runtime_request_value "$cc_request" required_layer)
  cc_produced_layer=$(cc_runtime_request_value "$cc_request" produced_layer)
  cc_observed_layer=$(cc_runtime_request_value "$cc_request" observed_layer)
  cc_evidence_outcome=$(cc_runtime_request_value "$cc_request" outcome)
  cc_evidence_ref=$(cc_runtime_request_value "$cc_request" evidence_ref)
  cc_evidence_source=$(cc_runtime_request_value "$cc_request" evidence_source)
  cc_host_capability=$(cc_runtime_request_value "$cc_request" host_capability)
  if test -z "$cc_required_layer" && test -z "$cc_produced_layer" && test -z "$cc_observed_layer"; then
    return 0
  fi
  printf 'evidence_layers:\n'
  if test -n "$cc_required_layer"; then printf '  required_layer: %s\n' "$cc_required_layer"; fi
  if test -n "$cc_produced_layer"; then printf '  produced_layer: %s\n' "$cc_produced_layer"; fi
  if test -n "$cc_observed_layer"; then printf '  observed_layer: %s\n' "$cc_observed_layer"; fi
  if test -n "$cc_evidence_outcome"; then printf '  outcome: %s\n' "$cc_evidence_outcome"; fi
  if test -n "$cc_evidence_ref"; then printf '  evidence_ref: %s\n' "$cc_evidence_ref"; fi
  if test -n "$cc_evidence_source"; then printf '  evidence_source: %s\n' "$cc_evidence_source"; fi
  if test -n "$cc_host_capability"; then printf '  host_capability: %s\n' "$cc_host_capability"; fi
  return 0
}

# Runtime record constructors and graph publication.  These functions accept
# only bounded caller intent and validated evidence; they do not invoke a host,
# provider, scheduler, or gate.  A transaction without commit.marker is never
# authoritative, even if every staged file happens to be readable.

cc_runtime_now() {
  date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || return 1
}

cc_runtime_request_value() {
  cc_request_file=$1
  cc_request_key=$2
  sed -n "s/^${cc_request_key}: //p" "$cc_request_file" | head -n 1
}

cc_runtime_require_request() {
  cc_request_value=$(cc_runtime_request_value "$1" "$2")
  test -n "$cc_request_value" && test "$cc_request_value" != null || {
    printf 'RUNTIME_INPUT_MISSING: %s\n' "$2" >&2
    return 1
  }
  printf '%s\n' "$cc_request_value"
}

cc_runtime_host_evidence() {
  cc_host_request=$1
  for cc_host_forbidden in credentials tokens provider_payload provider-payloads transcript transcripts auth_state auth-state; do
    if grep -E "^${cc_host_forbidden}:" "$cc_host_request" >/dev/null 2>&1; then
      printf 'HOST_EVIDENCE_FORBIDDEN: %s\n' "$cc_host_forbidden" >&2
      return 1
    fi
  done
  cc_host_id=$(cc_runtime_require_request "$cc_host_request" host_id) || return 1
  cc_host_version=$(cc_runtime_require_request "$cc_host_request" observed_version) || return 1
  cc_host_instruction=$(cc_runtime_require_request "$cc_host_request" instruction_surface) || return 1
  cc_host_role=$(cc_runtime_require_request "$cc_host_request" session_role) || return 1
  cc_host_root=$(cc_runtime_require_request "$cc_host_request" root_capability) || return 1
  cc_host_child=$(cc_runtime_require_request "$cc_host_request" child_capability) || return 1
  cc_host_verifier=$(cc_runtime_require_request "$cc_host_request" verifier_capability) || return 1
  cc_host_resume=$(cc_runtime_require_request "$cc_host_request" resume_capability) || return 1
  cc_host_permission=$(cc_runtime_require_request "$cc_host_request" permission_mode) || return 1
  cc_host_provider=$(cc_runtime_require_request "$cc_host_request" provider_status) || return 1
  cc_host_fallback=$(cc_runtime_require_request "$cc_host_request" offline_fallback) || return 1
  cc_host_smoke=$(cc_runtime_request_value "$cc_host_request" live_smoke_status)
  test -n "$cc_host_smoke" || cc_host_smoke=unavailable
  case "$cc_host_id" in codex|claude-code|cursor-agent) ;; *) printf '%s\n' HOST_ID_INVALID >&2; return 1 ;; esac
  case "$cc_host_instruction" in AGENTS.md|CLAUDE.md|shared-root) ;; *) printf '%s\n' HOST_INSTRUCTION_INVALID >&2; return 1 ;; esac
  case "$cc_host_role" in root|writer|verifier) ;; *) printf '%s\n' HOST_ROLE_INVALID >&2; return 1 ;; esac
  for cc_host_capability in "$cc_host_root" "$cc_host_child" "$cc_host_verifier" "$cc_host_resume"; do
    case "$cc_host_capability" in
      available|unavailable) ;;
      *) printf '%s\n' HOST_CAPABILITY_INVALID >&2; return 1 ;;
    esac
  done
  case "$cc_host_permission" in read-only|bounded-write|host-managed) ;; *) printf '%s\n' HOST_PERMISSION_INVALID >&2; return 1 ;; esac
  case "$cc_host_provider" in enabled|disabled|denied|unavailable) ;; *) printf '%s\n' HOST_PROVIDER_INVALID >&2; return 1 ;; esac
  case "$cc_host_fallback" in filesystem-only|host-blocked) ;; *) printf '%s\n' HOST_FALLBACK_INVALID >&2; return 1 ;; esac
  case "$cc_host_smoke" in pass|unavailable|host-blocked) ;; *) printf '%s\n' HOST_SMOKE_INVALID >&2; return 1 ;; esac
  cat <<EOF
host_evidence:
  host_id: $cc_host_id
  observed_version: $cc_host_version
  instruction_surface: $cc_host_instruction
  session_role: $cc_host_role
  root_capability: $cc_host_root
  child_capability: $cc_host_child
  verifier_capability: $cc_host_verifier
  resume_capability: $cc_host_resume
  permission_mode: $cc_host_permission
  provider_status: $cc_host_provider
  offline_fallback: $cc_host_fallback
  live_smoke_status: $cc_host_smoke
EOF
}

cc_runtime_record_strip_metadata() {
  awk '!/^(record_revision|record_bytes|record_digest|transaction_id|transaction_state|commit_marker|created_at|updated_at):[[:space:]]*/' "$1"
}

cc_runtime_record_digest_body() {
  cc_runtime_record_strip_metadata "$1" | awk '!/^ownership_graph_digest:[[:space:]]*/'
}

cc_runtime_record_seal() {
  cc_record_file=$1
  cc_record_transaction=$2
  cc_record_state=${3:-staged}
  cc_record_marker=${4:-null}
  test -f "$cc_record_file" || return 1
  case "$cc_record_state" in staged|valid|committed|interrupted|legacy-readable) ;; *) return 1 ;; esac
  cc_record_body="$cc_record_file.body.$$"
  cc_record_tmp="$cc_record_file.tmp.$$"
  cc_runtime_record_strip_metadata "$cc_record_file" > "$cc_record_body" || { rm -f "$cc_record_body"; return 1; }
  cc_record_digest_body="$cc_record_file.digest-body.$$"
  cc_runtime_record_digest_body "$cc_record_file" > "$cc_record_digest_body" || { rm -f "$cc_record_body" "$cc_record_digest_body"; return 1; }
  cc_record_revision=$(cc_digest "$cc_record_digest_body") || { rm -f "$cc_record_body" "$cc_record_digest_body"; return 1; }
  cc_record_bytes=$(wc -c < "$cc_record_digest_body" | tr -d ' \t\n')
  cc_record_now=$(cc_runtime_now) || { rm -f "$cc_record_body"; return 1; }
  {
    cat "$cc_record_body"
    printf 'created_at: %s\n' "$cc_record_now"
    printf 'updated_at: %s\n' "$cc_record_now"
    printf 'record_revision: %s\n' "$cc_record_revision"
    printf 'record_bytes: %s\n' "$cc_record_bytes"
    printf 'record_digest: %s\n' "$cc_record_revision"
    printf 'transaction_id: %s\n' "$cc_record_transaction"
    printf 'transaction_state: %s\n' "$cc_record_state"
    printf 'commit_marker: %s\n' "$cc_record_marker"
  } > "$cc_record_tmp" || { rm -f "$cc_record_body" "$cc_record_digest_body" "$cc_record_tmp"; return 1; }
  mv "$cc_record_tmp" "$cc_record_file" || { rm -f "$cc_record_body" "$cc_record_digest_body" "$cc_record_tmp"; return 1; }
  rm -f "$cc_record_body" "$cc_record_digest_body"
}

cc_construct_session() {
  cc_output=$1; cc_request=$2; cc_transaction=$3
  test -n "$cc_output" && test -f "$cc_request" || return 1
  cc_session=$(cc_runtime_require_request "$cc_request" session_id) || return 1
  cc_parent=$(cc_runtime_request_value "$cc_request" parent_session_id); test -n "$cc_parent" || cc_parent=null
  cc_root=$(cc_runtime_require_request "$cc_request" root_session_id) || return 1
  cc_kind=$(cc_runtime_require_request "$cc_request" kind) || return 1
  cc_role=$(cc_runtime_require_request "$cc_request" role) || return 1
  cc_objective=$(cc_runtime_require_request "$cc_request" objective) || return 1
  cc_scope=$(cc_runtime_require_request "$cc_request" scope) || return 1
  cc_worktree=$(cc_runtime_require_request "$cc_request" worktree) || return 1
  cc_safe_id "$cc_session" && cc_safe_id "$cc_root" || return 1
  case "$cc_kind" in root|child) ;; *) printf '%s\n' SESSION_KIND_INVALID >&2; return 1 ;; esac
  case "$cc_role" in root|writer|verifier) ;; *) printf '%s\n' SESSION_ROLE_INVALID >&2; return 1 ;; esac
  test "$cc_kind" = root && test "$cc_parent" = null || test "$cc_kind" = child && test "$cc_parent" != null || { printf '%s\n' SESSION_ANCESTRY_INVALID >&2; return 1; }
  cc_host_block=$(cc_runtime_host_evidence "$cc_request") || return 1
  cc_parent_dir=$(dirname "$cc_output"); mkdir -p "$cc_parent_dir" || return 1
  {
    printf 'schema_version: 1\nwrapper_version: %s\n' "${CC_WRAPPER_VERSION:-1.0.0}"
    printf 'session_id: %s\nparent_session_id: %s\nroot_session_id: %s\n' "$cc_session" "$cc_parent" "$cc_root"
    printf 'kind: %s\nrole: %s\nobjective: %s\nscope: %s\n' "$cc_kind" "$cc_role" "$cc_objective" "$cc_scope"
    printf 'permissions:\n  write_worktree: %s\n  write_plan: false\n  write_activity: false\n' "$(test "$cc_role" = writer && printf true || printf false)"
    printf 'worktree: %s\nstatus: created\n' "$cc_worktree"
    printf '%s\n' "$cc_host_block"
  } > "$cc_output" || return 1
  cc_runtime_record_seal "$cc_output" "$cc_transaction" staged
}

cc_construct_context_receipt() {
  cc_output=$1; cc_request=$2; cc_transaction=$3
  cc_session=$(cc_runtime_require_request "$cc_request" session_id) || return 1
  cc_route_digest=$(cc_runtime_require_request "$cc_request" route_decision_digest) || return 1
  cc_context_set=$(cc_runtime_require_request "$cc_request" context_set) || return 1
  cc_source=$(cc_runtime_request_value "$cc_request" receipt_source)
  if test -n "$cc_source"; then
    test -f "$cc_source" && test ! -L "$cc_source" || { printf '%s\n' RECEIPT_SOURCE_INVALID >&2; return 1; }
    cc_runtime_record_strip_metadata "$cc_source" > "$cc_output" || return 1
  else
    cc_parent_dir=$(dirname "$cc_output"); mkdir -p "$cc_parent_dir" || return 1
    {
      printf 'schema_version: 1\nwrapper_version: %s\n' "${CC_WRAPPER_VERSION:-1.0.0}"
      printf 'session_id: %s\nroute_decision_digest: %s\ncontext_set: %s\n' "$cc_session" "$cc_route_digest" "$cc_context_set"
      printf 'packet_id: %s\ndeclared_budget_bytes: 0\nactual_bytes: 0\npacket_digest: sha256:empty\nreferences: []\ninvariants: [INV-CTX-01, INV-CTX-02, INV-CTX-03, INV-CTX-04, INV-CTX-05]\n' "$cc_context_set"
    } > "$cc_output" || return 1
  fi
  cc_runtime_record_seal "$cc_output" "$cc_transaction" staged
}

cc_construct_delegation() {
  cc_output=$1; cc_request=$2; cc_transaction=$3
  cc_session_id=$(cc_runtime_require_request "$cc_request" session_id) || return 1
  cc_parent_session_id=$(cc_runtime_require_request "$cc_request" parent_session_id) || return 1
  cc_root_session_id=$(cc_runtime_require_request "$cc_request" root_session_id) || return 1
  cc_role=$(cc_runtime_require_request "$cc_request" role) || return 1
  cc_objective=$(cc_runtime_require_request "$cc_request" objective) || return 1
  cc_scope=$(cc_runtime_require_request "$cc_request" scope) || return 1
  cc_non_goals=$(cc_runtime_require_request "$cc_request" non_goals) || return 1
  cc_context_set=$(cc_runtime_require_request "$cc_request" context_set) || return 1
  cc_repository=$(cc_runtime_require_request "$cc_request" repository) || return 1
  cc_worktree=$(cc_runtime_require_request "$cc_request" worktree) || return 1
  cc_acceptance_criteria=$(cc_runtime_require_request "$cc_request" acceptance_criteria) || return 1
  cc_verification=$(cc_runtime_require_request "$cc_request" verification) || return 1
  cc_expected_evidence=$(cc_runtime_require_request "$cc_request" expected_evidence) || return 1
  cc_stop_conditions=$(cc_runtime_require_request "$cc_request" stop_conditions) || return 1
  cc_parent_dir=$(dirname "$cc_output"); mkdir -p "$cc_parent_dir" || return 1
  cc_host_block=$(cc_runtime_host_evidence "$cc_request") || return 1
  cc_permissions_writer=false; test "$cc_role" = writer && cc_permissions_writer=true
  {
    printf 'schema_version: 1\nwrapper_version: %s\n' "${CC_WRAPPER_VERSION:-1.0.0}"
    printf 'delegation_id: %s-delegation\n' "$cc_transaction"
    printf 'session_id: %s\nparent_session_id: %s\nroot_session_id: %s\nrole: %s\n' "$cc_session_id" "$cc_parent_session_id" "$cc_root_session_id" "$cc_role"
    printf 'objective: %s\nscope: %s\nnon_goals: %s\ncontext_set: %s\ncontext_receipt: context-receipt.yaml\n' "$cc_objective" "$cc_scope" "$cc_non_goals" "$cc_context_set"
    printf 'repository: %s\nworktree: %s\npermissions:\n  write_worktree: %s\n  write_plan: false\n  write_activity: false\n' "$cc_repository" "$cc_worktree" "$cc_permissions_writer"
    printf 'acceptance_criteria: %s\nverification: %s\nexpected_evidence: %s\nstop_conditions: %s\nhandoff_schema: wrapper/contracts/schemas/handoff.yaml\n' "$cc_acceptance_criteria" "$cc_verification" "$cc_expected_evidence" "$cc_stop_conditions"
    cc_runtime_emit_evidence_layers "$cc_request"
    printf '%s\n' "$cc_host_block"
  } > "$cc_output" || return 1
  cc_runtime_record_seal "$cc_output" "$cc_transaction" staged
}

cc_construct_handoff() {
  cc_output=$1; cc_request=$2; cc_transaction=$3
  cc_session=$(cc_runtime_require_request "$cc_request" session_id) || return 1
  cc_parent_dir=$(dirname "$cc_output"); mkdir -p "$cc_parent_dir" || return 1
  cc_now=$(cc_runtime_now) || return 1
  {
    printf '%s\n' '---'
    printf 'schema_version: 1\nwrapper_version: %s\nhandoff_id: %s-handoff\nsession_id: %s\n' "${CC_WRAPPER_VERSION:-1.0.0}" "$cc_transaction" "$cc_session"
    cc_runtime_emit_evidence_layers "$cc_request"
    printf '%s\n' '---' '# Runtime handoff skeleton' '' '## Observed state' '' 'Generated by the host-neutral engine.' '' '## Selected route and reason' '' 'Bounded runtime graph construction.' '' '## Evidence read' '' 'Validated session, receipt, delegation, and ownership inputs.' '' '## Host evidence' '' 'Provider-neutral host evidence remains in the delegation record.' '' '## Action performed or proposed' '' 'No provider or host action is performed by the engine.' '' '## Files/state changed' '' 'Runtime graph transaction only.' '' '## Verification' '' 'Independent verification required. Writer-recorded mappings are claims.' '' '## Blockers and human decisions' '' 'None recorded by the constructor.' '' '## One next safe action' '' 'Validate the complete graph before publication.'
    printf '<!-- created_at: %s -->\n' "$cc_now"
  } > "$cc_output" || return 1
  cc_runtime_record_seal "$cc_output" "$cc_transaction" staged
}

cc_construct_completion() {
  cc_output=$1; cc_request=$2; cc_transaction=$3
  cc_plan=$(cc_runtime_require_request "$cc_request" plan) || return 1
  cc_parent_dir=$(dirname "$cc_output"); mkdir -p "$cc_parent_dir" || return 1
  {
    printf 'schema_version: 1\nwrapper_version: %s\ncompletion_id: %s-completion\nplan: %s\nstatus: not-ready\ntask_evidence: pending\nverifier: pending\nworktree: pending\ngit: pending\nhuman_gate: required\n' "${CC_WRAPPER_VERSION:-1.0.0}" "$cc_transaction" "$cc_plan"
    cc_runtime_emit_evidence_layers "$cc_request"
  } > "$cc_output" || return 1
  cc_runtime_record_seal "$cc_output" "$cc_transaction" staged
}

cc_construct_child_start() {
  cc_output=$1; cc_request=$2; cc_transaction=$3; cc_delegation=$4; cc_receipt=$5
  test -f "$cc_delegation" && test -f "$cc_receipt" || return 1
  cc_validate_delegation "$cc_delegation" || { printf '%s\n' INVALID_DELEGATION >&2; return 1; }
  cc_validate_receipt "$cc_receipt" || { printf '%s\n' INVALID_RECEIPT >&2; return 1; }
  cc_session=$(cc_runtime_require_request "$cc_request" session_id) || return 1
  cc_parent=$(cc_runtime_require_request "$cc_request" parent_session_id) || return 1
  cc_root=$(cc_runtime_require_request "$cc_request" root_session_id) || return 1
  cc_role=$(cc_runtime_require_request "$cc_request" role) || return 1
  cc_repository=$(cc_runtime_require_request "$cc_request" repository) || return 1
  cc_worktree=$(cc_runtime_require_request "$cc_request" worktree) || return 1
  cc_delegation_digest=$(cc_digest "$cc_delegation") || return 1
  cc_receipt_digest=$(cc_digest "$cc_receipt") || return 1
  cc_delegation_bytes=$(wc -c < "$cc_delegation" | tr -d ' \t\n')
  cc_receipt_bytes=$(wc -c < "$cc_receipt" | tr -d ' \t\n')
  cc_host_digest=$(cc_digest_text "$(grep -A12 '^host_evidence:' "$cc_delegation")")
  cc_parent_dir=$(dirname "$cc_output"); mkdir -p "$cc_parent_dir" || return 1
  {
    printf 'schema_version: 1\nwrapper_version: %s\nchild_start_id: %s\n' "${CC_WRAPPER_VERSION:-1.0.0}" "$cc_transaction"
    printf 'session_id: %s\nparent_session_id: %s\nroot_session_id: %s\nrole: %s\nrepository: %s\nworktree: %s\n' "$cc_session" "$cc_parent" "$cc_root" "$cc_role" "$cc_repository" "$cc_worktree"
    printf 'delegation_locator: delegation.yaml\ndelegation_revision: %s\ndelegation_bytes: %s\ndelegation_digest: %s\n' "$cc_delegation_digest" "$cc_delegation_bytes" "$cc_delegation_digest"
    printf 'context_receipt_locator: context-receipt.yaml\ncontext_receipt_revision: %s\ncontext_receipt_bytes: %s\ncontext_receipt_digest: %s\n' "$cc_receipt_digest" "$cc_receipt_bytes" "$cc_receipt_digest"
    printf 'host_evidence_digest: %s\nownership_graph_digest: pending\nlaunch_locator: launch.yaml\nhandoff_locator: handoff.md\n' "$cc_host_digest"
  } > "$cc_output" || return 1
  cc_runtime_record_seal "$cc_output" "$cc_transaction" staged
}

cc_runtime_transaction_begin() {
  cc_runtime_root=$1; cc_transaction=$2
  cc_safe_id "$cc_transaction" || return 1
  cc_stage="$cc_runtime_root/.transactions/$cc_transaction"
  test ! -e "$cc_stage" || { printf '%s\n' TRANSACTION_EXISTS >&2; return 1; }
  mkdir -p "$cc_stage/records" || return 1
  cc_atomic_write "$cc_stage/transaction.yaml" "schema_version: 1" "transaction_id: $cc_transaction" 'transaction_state: staged' 'commit_marker: null'
}

cc_runtime_graph_files() {
  printf '%s\n' session.yaml context-receipt.yaml delegation.yaml child-start.yaml handoff.md completion.yaml lease.yaml
}

cc_runtime_graph_record_path() {
  printf '%s/records/%s\n' "$1" "$2"
}

cc_runtime_record_set_field() {
  cc_set_file=$1
  cc_set_field=$2
  cc_set_value=$3
  cc_set_tmp="$cc_set_file.set.$$"
  awk -v field="$cc_set_field" -v value="$cc_set_value" '
    BEGIN { replaced=0 }
    $0 ~ "^" field ":" { print field ": " value; replaced=1; next }
    { print }
    END { exit replaced ? 0 : 1 }
  ' "$cc_set_file" > "$cc_set_tmp" || { rm -f "$cc_set_tmp"; return 1; }
  mv "$cc_set_tmp" "$cc_set_file"
}

cc_runtime_record_validate_metadata() {
  cc_metadata_file=$1
  cc_metadata_transaction=$2
  cc_metadata_state=$3
  cc_metadata_marker=${4:-}
  test -f "$cc_metadata_file" || return 1
  for cc_metadata_field in schema_version wrapper_version record_revision record_bytes record_digest transaction_id transaction_state commit_marker created_at updated_at; do
    grep -E "^${cc_metadata_field}:" "$cc_metadata_file" >/dev/null 2>&1 || { printf 'RECORD_METADATA_MISSING: %s:%s\n' "$cc_metadata_file" "$cc_metadata_field" >&2; return 1; }
  done
  cc_metadata_id_field=
  case "$(basename "$cc_metadata_file")" in
    session.yaml) cc_metadata_id_field=session_id ;;
    context-receipt.yaml) cc_metadata_id_field=session_id ;;
    delegation.yaml) cc_metadata_id_field=delegation_id ;;
    child-start.yaml) cc_metadata_id_field=child_start_id ;;
    handoff.md) cc_metadata_id_field=handoff_id ;;
    completion.yaml) cc_metadata_id_field=completion_id ;;
  esac
  if test -n "$cc_metadata_id_field"; then
    cc_metadata_id=$(sed -n "s/^${cc_metadata_id_field}: //p" "$cc_metadata_file" | head -n 1)
    test -n "$cc_metadata_id" && test "$cc_metadata_id" != null && test "$cc_metadata_id" != pending || { printf 'RECORD_ID_MISSING: %s\n' "$cc_metadata_file" >&2; return 1; }
  fi
  cc_metadata_bytes=$(sed -n 's/^record_bytes: //p' "$cc_metadata_file" | head -n 1)
  case "$cc_metadata_bytes" in ''|*[!0-9]*) printf 'RECORD_BYTES_INVALID: %s\n' "$cc_metadata_file" >&2; return 1 ;; esac
  cc_metadata_record_transaction=$(sed -n 's/^transaction_id: //p' "$cc_metadata_file" | head -n 1)
  test "$cc_metadata_record_transaction" = "$cc_metadata_transaction" || { printf 'RECORD_TRANSACTION_MISMATCH: %s\n' "$cc_metadata_file" >&2; return 1; }
  cc_metadata_record_state=$(sed -n 's/^transaction_state: //p' "$cc_metadata_file" | head -n 1)
  test "$cc_metadata_record_state" = "$cc_metadata_state" || { printf 'RECORD_STATE_MISMATCH: %s\n' "$cc_metadata_file" >&2; return 1; }
  if test -n "$cc_metadata_marker"; then
    cc_metadata_record_marker=$(sed -n 's/^commit_marker: //p' "$cc_metadata_file" | head -n 1)
    test "$cc_metadata_record_marker" = "$cc_metadata_marker" || { printf 'RECORD_MARKER_MISMATCH: %s\n' "$cc_metadata_file" >&2; return 1; }
  fi
  cc_metadata_digest_body="$cc_metadata_file.verify-body.$$"
  cc_runtime_record_digest_body "$cc_metadata_file" > "$cc_metadata_digest_body" || { rm -f "$cc_metadata_digest_body"; return 1; }
  cc_metadata_current_digest=$(cc_digest "$cc_metadata_digest_body") || { rm -f "$cc_metadata_digest_body"; return 1; }
  cc_metadata_current_bytes=$(wc -c < "$cc_metadata_digest_body" | tr -d ' \t\n')
  rm -f "$cc_metadata_digest_body"
  test "$cc_metadata_current_digest" = "$(sed -n 's/^record_revision: //p' "$cc_metadata_file" | head -n 1)" || { printf 'RECORD_REVISION_STALE: %s\n' "$cc_metadata_file" >&2; return 1; }
  test "$cc_metadata_current_digest" = "$(sed -n 's/^record_digest: //p' "$cc_metadata_file" | head -n 1)" || { printf 'RECORD_DIGEST_STALE: %s\n' "$cc_metadata_file" >&2; return 1; }
  test "$cc_metadata_current_bytes" = "$cc_metadata_bytes" || { printf 'RECORD_BYTES_STALE: %s\n' "$cc_metadata_file" >&2; return 1; }
  grep -Eq '^created_at: [0-9]{4}-[0-9]{2}-[0-9]{2}T[^ ]+Z$' "$cc_metadata_file" || return 1
  grep -Eq '^updated_at: [0-9]{4}-[0-9]{2}-[0-9]{2}T[^ ]+Z$' "$cc_metadata_file" || return 1
  if test "$(basename "$cc_metadata_file")" = child-start.yaml && test "$cc_metadata_state" = committed; then
    test "$(sed -n 's/^ownership_graph_digest: //p' "$cc_metadata_file" | head -n 1)" = "$cc_metadata_marker" || { printf '%s\n' CHILD_START_GRAPH_DIGEST_MISMATCH >&2; return 1; }
  fi
}

cc_runtime_graph_digest() {
  cc_graph_dir=$1
  cc_graph_material_file="$cc_graph_dir.graph-material.$$"
  : > "$cc_graph_material_file" || return 1
  for cc_graph_name in $(cc_runtime_graph_files); do
    cc_graph_file=$(cc_runtime_graph_record_path "$cc_graph_dir" "$cc_graph_name")
    test -f "$cc_graph_file" || { rm -f "$cc_graph_material_file"; return 1; }
    cc_graph_body="$cc_graph_file.graph-body.$$"
    cc_runtime_record_digest_body "$cc_graph_file" > "$cc_graph_body" || { rm -f "$cc_graph_material_file" "$cc_graph_body"; return 1; }
    cc_graph_current_digest=$(cc_digest "$cc_graph_body") || { rm -f "$cc_graph_material_file" "$cc_graph_body"; return 1; }
    cc_graph_current_bytes=$(wc -c < "$cc_graph_body" | tr -d ' \t\n')
    rm -f "$cc_graph_body"
    printf '%s|%s|%s|%s\n' "$cc_graph_name" "$cc_graph_current_digest" "$(sed -n 's/^record_digest: //p' "$cc_graph_file" | head -n 1)" "$cc_graph_current_bytes" >> "$cc_graph_material_file"
  done
  cc_digest "$cc_graph_material_file"
  cc_graph_result=$?
  rm -f "$cc_graph_material_file"
  return "$cc_graph_result"
}

cc_runtime_validate_host_evidence() {
  cc_host_file=$1
  for cc_host_field in host_id observed_version instruction_surface session_role root_capability child_capability verifier_capability resume_capability permission_mode provider_status offline_fallback; do
    grep -E "^  ${cc_host_field}:" "$cc_host_file" >/dev/null 2>&1 || { printf 'HOST_EVIDENCE_MISSING: %s\n' "$cc_host_field" >&2; return 1; }
  done
  if grep -E '^[[:space:]]*(credentials|tokens|provider_payload|provider-payloads|transcript|transcripts|auth_state|auth-state):' "$cc_host_file" >/dev/null 2>&1; then
    printf '%s\n' HOST_EVIDENCE_FORBIDDEN >&2
    return 1
  fi
  return 0
}

cc_validate_launch_projection() {
  cc_projection_file=$1; cc_delegation_file=$2
  test -f "$cc_projection_file" && test -f "$cc_delegation_file" || return 1
  cc_validate_delegation "$cc_delegation_file" || return 1
  cc_projection_keys=$(awk -F: '/^[a-z_]+:/ { print $1 }' "$cc_projection_file")
  for cc_key in $cc_projection_keys; do
    case "$cc_key" in role|assigned_root|delegation_locator|handoff_locator) ;; *) printf 'LAUNCH_PROJECTION_EXPANDED: %s\n' "$cc_key" >&2; return 1 ;; esac
  done
  test "$(sed -n 's/^role: //p' "$cc_projection_file" | head -n 1)" = "$(sed -n 's/^role: //p' "$cc_delegation_file" | head -n 1)" || { printf '%s\n' LAUNCH_ROLE_MISMATCH >&2; return 1; }
  test "$(sed -n 's/^assigned_root: //p' "$cc_projection_file" | head -n 1)" = "$(sed -n 's/^worktree: //p' "$cc_delegation_file" | head -n 1)" || { printf '%s\n' LAUNCH_ROOT_MISMATCH >&2; return 1; }
  test "$(sed -n 's/^delegation_locator: //p' "$cc_projection_file" | head -n 1)" = delegation.yaml || { printf '%s\n' LAUNCH_DELEGATION_LOCATOR_INVALID >&2; return 1; }
  test "$(sed -n 's/^handoff_locator: //p' "$cc_projection_file" | head -n 1)" = handoff.md || { printf '%s\n' LAUNCH_HANDOFF_LOCATOR_INVALID >&2; return 1; }
}

cc_runtime_launch_projection() {
  cc_child_start=$1; cc_delegation=$2; cc_output=$3
  test -f "$cc_child_start" && test -f "$cc_delegation" || return 1
  cc_validate_delegation "$cc_delegation" || return 1
  {
    printf 'role: %s\n' "$(sed -n 's/^role: //p' "$cc_delegation" | head -n 1)"
    printf 'assigned_root: %s\n' "$(sed -n 's/^worktree: //p' "$cc_delegation" | head -n 1)"
    printf 'delegation_locator: delegation.yaml\nhandoff_locator: handoff.md\n'
  } > "$cc_output" || return 1
  cc_validate_launch_projection "$cc_output" "$cc_delegation"
}

cc_validate_runtime_graph() {
  cc_graph_dir=$1
  cc_graph_expected_state=${2:-}
  cc_records="$cc_graph_dir/records"
  test -d "$cc_records" || { printf '%s\n' RUNTIME_GRAPH_MISSING >&2; return 1; }
  test -f "$cc_graph_dir/transaction.yaml" || { printf '%s\n' TRANSACTION_RECORD_MISSING >&2; return 1; }
  for cc_graph_name in session.yaml context-receipt.yaml delegation.yaml child-start.yaml handoff.md completion.yaml lease.yaml; do
    test -f "$cc_records/$cc_graph_name" || { printf 'RUNTIME_RECORD_MISSING: %s\n' "$cc_graph_name" >&2; return 1; }
  done
  cc_validate_delegation "$cc_records/delegation.yaml" || { printf '%s\n' DELEGATION_INVALID >&2; return 1; }
  cc_validate_receipt "$cc_records/context-receipt.yaml" || { printf '%s\n' RECEIPT_INVALID >&2; return 1; }
  cc_runtime_validate_host_evidence "$cc_records/delegation.yaml" || return 1
  cc_graph_transaction=$(sed -n 's/^transaction_id: //p' "$cc_graph_dir/transaction.yaml" | head -n 1)
  test "$cc_graph_transaction" = "$(basename "$cc_graph_dir")" || { printf '%s\n' TRANSACTION_ID_MISMATCH >&2; return 1; }
  cc_session=$(sed -n 's/^session_id: //p' "$cc_records/session.yaml" | head -n 1)
  cc_delegation_session=$(sed -n 's/^session_id: //p' "$cc_records/delegation.yaml" | head -n 1)
  cc_child_session=$(sed -n 's/^session_id: //p' "$cc_records/child-start.yaml" | head -n 1)
  test "$cc_session" = "$cc_delegation_session" && test "$cc_session" = "$cc_child_session" || { printf '%s\n' OWNERSHIP_SESSION_MISMATCH >&2; return 1; }
  cc_root=$(sed -n 's/^root_session_id: //p' "$cc_records/session.yaml" | head -n 1)
  cc_delegation_parent=$(sed -n 's/^parent_session_id: //p' "$cc_records/delegation.yaml" | head -n 1)
  cc_session_parent=$(sed -n 's/^parent_session_id: //p' "$cc_records/session.yaml" | head -n 1)
  test "$cc_root" = "$(sed -n 's/^root_session_id: //p' "$cc_records/delegation.yaml" | head -n 1)" && test "$cc_root" = "$(sed -n 's/^root_session_id: //p' "$cc_records/child-start.yaml" | head -n 1)" || { printf '%s\n' OWNERSHIP_ROOT_MISMATCH >&2; return 1; }
  test "$cc_session_parent" = "$cc_delegation_parent" && test "$cc_delegation_parent" = "$(sed -n 's/^parent_session_id: //p' "$cc_records/child-start.yaml" | head -n 1)" || { printf '%s\n' OWNERSHIP_PARENT_MISMATCH >&2; return 1; }
  cc_role=$(sed -n 's/^role: //p' "$cc_records/delegation.yaml" | head -n 1)
  test "$cc_role" = "$(sed -n 's/^role: //p' "$cc_records/session.yaml" | head -n 1)" && test "$cc_role" = "$(sed -n 's/^role: //p' "$cc_records/child-start.yaml" | head -n 1)" || { printf '%s\n' OWNERSHIP_ROLE_MISMATCH >&2; return 1; }
  cc_repository=$(sed -n 's/^repository: //p' "$cc_records/delegation.yaml" | head -n 1)
  test -n "$cc_repository" && test "$cc_repository" = "$(sed -n 's/^repository: //p' "$cc_records/child-start.yaml" | head -n 1)" || { printf '%s\n' OWNERSHIP_REPOSITORY_MISMATCH >&2; return 1; }
  cc_worktree=$(sed -n 's/^worktree: //p' "$cc_records/delegation.yaml" | head -n 1)
  cc_session_worktree=$(sed -n 's/^worktree: //p' "$cc_records/session.yaml" | head -n 1)
  test "$cc_worktree" = "$cc_session_worktree" && test "$cc_worktree" = "$(sed -n 's/^worktree: //p' "$cc_records/child-start.yaml" | head -n 1)" || { printf '%s\n' OWNERSHIP_WORKTREE_MISMATCH >&2; return 1; }
  cc_receipt_session=$(sed -n 's/^session_id: //p' "$cc_records/context-receipt.yaml" | head -n 1)
  cc_receipt_packet_id=$(sed -n 's/^packet_id: //p' "$cc_records/context-receipt.yaml" | head -n 1)
  cc_receipt_context_set=$(sed -n 's/^context_set: //p' "$cc_records/context-receipt.yaml" | head -n 1)
  cc_delegation_context_set=$(sed -n 's/^context_set: //p' "$cc_records/delegation.yaml" | head -n 1)
  test "$cc_receipt_session" = "$cc_session" && test "$cc_receipt_packet_id" = "$cc_receipt_context_set" && test "$cc_receipt_context_set" = "$cc_delegation_context_set" || { printf '%s\n' RECEIPT_OWNERSHIP_MISMATCH >&2; return 1; }
  cc_validate_child_permissions "$cc_records/delegation.yaml" "$cc_role" || return 1
  cc_validate_launch_projection "$cc_graph_dir/launch.yaml" "$cc_records/delegation.yaml" || return 1
  for cc_handoff_section in '## Observed state' '## Selected route and reason' '## Evidence read' '## Host evidence' '## Action performed or proposed' '## Files/state changed' '## Verification' '## Blockers and human decisions' '## One next safe action'; do
    grep -F "$cc_handoff_section" "$cc_records/handoff.md" >/dev/null 2>&1 || { printf 'HANDOFF_SECTION_MISSING: %s\n' "$cc_handoff_section" >&2; return 1; }
  done
  cc_validate_delegation_evidence "$cc_records/delegation.yaml" || return 1
  cc_validate_handoff_evidence "$cc_records/handoff.md" || return 1
  cc_validate_completion_evidence "$cc_records/completion.yaml" || return 1
  cc_lease_session=$(sed -n 's/^session_id: //p' "$cc_records/lease.yaml" | head -n 1)
  cc_lease_root=$(sed -n 's/^root_session_id: //p' "$cc_records/lease.yaml" | head -n 1)
  cc_lease_repository=$(sed -n 's/^repository: //p' "$cc_records/lease.yaml" | head -n 1)
  cc_lease_worktree=$(sed -n 's/^worktree: //p' "$cc_records/lease.yaml" | head -n 1)
  test "$cc_lease_session" = "$cc_session" && test "$cc_lease_root" = "$cc_root" && test "$cc_lease_repository" = "$cc_repository" && test "$cc_lease_worktree" = "$cc_worktree" || { printf '%s\n' LEASE_OWNER_MISMATCH >&2; return 1; }
  cc_state=$(sed -n 's/^transaction_state: //p' "$cc_graph_dir/transaction.yaml" | head -n 1)
  if test -n "$cc_graph_expected_state"; then
    test "$cc_state" = "$cc_graph_expected_state" || { printf '%s\n' TRANSACTION_STATE_INVALID >&2; return 1; }
  else
    test "$cc_state" = staged || test "$cc_state" = valid || test "$cc_state" = committed || { printf '%s\n' TRANSACTION_STATE_INVALID >&2; return 1; }
  fi
  printf '%s\n' runtime-graph-valid
}

cc_validate_child_permissions() {
  cc_permission_file=$1; cc_permission_role=$2
  test "$cc_permission_role" = writer || test "$cc_permission_role" = verifier || return 1
  if test "$cc_permission_role" = writer; then
    grep -F 'write_worktree: true' "$cc_permission_file" >/dev/null 2>&1 || return 1
  else
    grep -F 'write_worktree: false' "$cc_permission_file" >/dev/null 2>&1 || return 1
  fi
  grep -F 'write_plan: false' "$cc_permission_file" >/dev/null 2>&1 && grep -F 'write_activity: false' "$cc_permission_file" >/dev/null 2>&1
}

cc_runtime_graph_validate_and_mark() {
  cc_graph_dir=$1
  cc_validate_runtime_graph "$cc_graph_dir" >/dev/null || return 1
  cc_graph_material=$(cc_runtime_graph_files | while IFS= read -r cc_graph_name; do
    cc_graph_file=$(cc_runtime_graph_record_path "$cc_graph_dir" "$cc_graph_name")
    test -f "$cc_graph_file" || continue
    case "$cc_graph_name" in *.md) cc_runtime_record_seal "$cc_graph_file" "$(sed -n 's/^transaction_id: //p' "$cc_graph_file" | head -n 1)" valid >/dev/null || return 1 ;; *) cc_runtime_record_seal "$cc_graph_file" "$(sed -n 's/^transaction_id: //p' "$cc_graph_file" | head -n 1)" valid >/dev/null || return 1 ;; esac
    printf '%s\n' "$cc_graph_name"
  done)
  cc_atomic_write "$cc_graph_dir/transaction.yaml" 'schema_version: 1' "transaction_id: $(basename "$cc_graph_dir")" 'transaction_state: valid' 'commit_marker: null'
  printf '%s\n' runtime-graph-valid
}

cc_runtime_graph_commit() {
  cc_graph_dir=$1
  cc_validate_runtime_graph "$cc_graph_dir" valid >/dev/null || return 1
  cc_state=$(sed -n 's/^transaction_state: //p' "$cc_graph_dir/transaction.yaml" | head -n 1)
  test "$cc_state" = valid || { printf '%s\n' TRANSACTION_NOT_VALID >&2; return 1; }
  cc_graph_digest=$(cc_runtime_graph_digest "$cc_graph_dir") || return 1
  cc_runtime_record_set_field "$cc_graph_dir/records/child-start.yaml" ownership_graph_digest "$cc_graph_digest" || return 1
  cc_runtime_graph_files | while IFS= read -r cc_graph_name; do
    cc_graph_file=$(cc_runtime_graph_record_path "$cc_graph_dir" "$cc_graph_name")
    test -f "$cc_graph_file" || continue
    cc_runtime_record_seal "$cc_graph_file" "$(sed -n 's/^transaction_id: //p' "$cc_graph_file" | head -n 1)" committed "$cc_graph_digest" || exit 1
  done || return 1
  cc_atomic_write "$cc_graph_dir/transaction.yaml" 'schema_version: 1' "transaction_id: $(basename "$cc_graph_dir")" 'transaction_state: committed' "commit_marker: $cc_graph_digest" || return 1
  cc_atomic_write "$cc_graph_dir/commit.marker" 'schema_version: 1' "transaction_id: $(basename "$cc_graph_dir")" 'transaction_state: committed' "ownership_graph_digest: $cc_graph_digest" "committed_at: $(cc_runtime_now)" || return 1
  printf '%s\n' runtime-graph-committed
}

cc_runtime_graph_authoritative() {
  cc_graph_dir=$1
  test -f "$cc_graph_dir/commit.marker" && test -f "$cc_graph_dir/transaction.yaml" || return 1
  cc_transaction_id=$(basename "$cc_graph_dir")
  cc_transaction_record_id=$(sed -n 's/^transaction_id: //p' "$cc_graph_dir/transaction.yaml" | head -n 1)
  cc_marker_transaction_id=$(sed -n 's/^transaction_id: //p' "$cc_graph_dir/commit.marker" | head -n 1)
  cc_transaction_state=$(sed -n 's/^transaction_state: //p' "$cc_graph_dir/transaction.yaml" | head -n 1)
  cc_marker_state=$(sed -n 's/^transaction_state: //p' "$cc_graph_dir/commit.marker" | head -n 1)
  cc_transaction_marker=$(sed -n 's/^commit_marker: //p' "$cc_graph_dir/transaction.yaml" | head -n 1)
  cc_marker_digest=$(sed -n 's/^ownership_graph_digest: //p' "$cc_graph_dir/commit.marker" | head -n 1)
  test "$cc_transaction_record_id" = "$cc_transaction_id" && test "$cc_marker_transaction_id" = "$cc_transaction_id" || return 1
  test "$cc_transaction_state" = committed && test "$cc_marker_state" = committed || return 1
  test -n "$cc_marker_digest" && test "$cc_transaction_marker" = "$cc_marker_digest" || return 1
  cc_validate_runtime_graph "$cc_graph_dir" committed >/dev/null || return 1
  for cc_authoritative_name in $(cc_runtime_graph_files); do
    cc_authoritative_file=$(cc_runtime_graph_record_path "$cc_graph_dir" "$cc_authoritative_name")
    cc_runtime_record_validate_metadata "$cc_authoritative_file" "$cc_transaction_id" committed "$cc_marker_digest" || return 1
  done
  cc_current_graph_digest=$(cc_runtime_graph_digest "$cc_graph_dir") || return 1
  test "$cc_current_graph_digest" = "$cc_marker_digest" || { printf '%s\n' GRAPH_DIGEST_STALE >&2; return 1; }
  printf '%s\n' runtime-graph-authoritative
}

cc_construct_runtime_graph() {
  cc_runtime_root=$1; cc_request=$2
  test -d "$cc_runtime_root" && test -f "$cc_request" || return 1
  cc_transaction=$(cc_runtime_require_request "$cc_request" transaction_id) || return 1
  cc_runtime_transaction_begin "$cc_runtime_root" "$cc_transaction" || return 1
  cc_graph_dir="$cc_runtime_root/.transactions/$cc_transaction"
  cc_records="$cc_graph_dir/records"
  cc_construct_session "$cc_records/session.yaml" "$cc_request" "$cc_transaction" || return 1
  cc_construct_context_receipt "$cc_records/context-receipt.yaml" "$cc_request" "$cc_transaction" || return 1
  cc_construct_delegation "$cc_records/delegation.yaml" "$cc_request" "$cc_transaction" || return 1
  cc_construct_child_start "$cc_records/child-start.yaml" "$cc_request" "$cc_transaction" "$cc_records/delegation.yaml" "$cc_records/context-receipt.yaml" || return 1
  cc_construct_handoff "$cc_records/handoff.md" "$cc_request" "$cc_transaction" || return 1
  cc_construct_completion "$cc_records/completion.yaml" "$cc_request" "$cc_transaction" || return 1
  cc_lease_source=$(cc_runtime_require_request "$cc_request" lease_source) || return 1
  test -f "$cc_lease_source" && test ! -L "$cc_lease_source" || { printf '%s\n' LEASE_SOURCE_INVALID >&2; return 1; }
  cp "$cc_lease_source" "$cc_records/lease.yaml" || return 1
  cc_runtime_record_seal "$cc_records/lease.yaml" "$cc_transaction" staged || return 1
  cc_runtime_launch_projection "$cc_records/child-start.yaml" "$cc_records/delegation.yaml" "$cc_graph_dir/launch.yaml" 2>/dev/null || return 1
  cc_runtime_graph_validate_and_mark "$cc_graph_dir" >/dev/null || return 1
  cc_runtime_graph_commit "$cc_graph_dir" >/dev/null || return 1
  printf '%s\n' "$cc_graph_dir"
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
  printf 'Immediate effects: repository-bootstrap\n'
  printf 'Later authorized effects: execute-plan\n'
  printf 'Not execute-plan effects: git.init, git.clone\n'
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

# Identity-region projection, displayed defaults, and descriptive effect
# identifiers. These primitives enforce agreement with workspace.yaml; they
# never derive Product Knowledge or grant authorization.

cc_identity_region_start() {
  printf '%s' '<!-- context-circuit:identity-region:start -->'
}

cc_identity_region_end() {
  printf '%s' '<!-- context-circuit:identity-region:end -->'
}

cc_identity_region_relpaths() {
  printf '%s\n' context/WORKSPACE.md context/PROJECT.md context/INDEX.md
}

cc_workspace_section_field() {
  cc_section_file=$1
  cc_section_name=$2
  cc_section_field=$3
  test -f "$cc_section_file" || return 0
  awk -v section="$cc_section_name" -v field="$cc_section_field" '
    $0 ~ "^" section ":" { in_section=1; next }
    /^[^[:space:]#]/ { in_section=0 }
    in_section && $0 ~ "^  " field ":[[:space:]]*" {
      sub("^  " field ":[[:space:]]*", "")
      gsub("[[:space:]]+$", "")
      print
      exit
    }
  ' "$cc_section_file"
}

cc_workspace_repository_keys() {
  cc_keys_file=$1
  test -f "$cc_keys_file" || return 0
  awk '
    /^repositories:[[:space:]]*\{\}[[:space:]]*$/ { exit }
    /^repositories:[[:space:]]*$/ { in_repositories=1; next }
    /^[^[:space:]#]/ { in_repositories=0 }
    in_repositories && /^  [A-Za-z0-9._-]+:[[:space:]]*$/ {
      key=$1
      sub(":", "", key)
      print key
    }
  ' "$cc_keys_file"
}

cc_identity_norm() {
  case "$1" in
    ''|null|[]) printf '%s\n' none ;;
    *) printf '%s\n' "$1" ;;
  esac
}

cc_roles_canonical() {
  cc_roles_raw=$(cc_identity_norm "$1")
  case "$cc_roles_raw" in
    none) printf '%s\n' none; return 0 ;;
  esac
  printf '%s\n' "$cc_roles_raw" | awk '{
    gsub(",[[:space:]]+", ",", $0)
    gsub("[[:space:]]+,", ",", $0)
    gsub("^[[:space:]]+|[[:space:]]+$", "", $0)
    if ($0 == "" || $0 == "none") print "none"
    else print $0
  }'
}

cc_workspace_roles() {
  cc_roles_file=$1
  test -f "$cc_roles_file" || { printf '%s\n' ''; return 0; }
  awk '
    /^workspace:[[:space:]]*$/ { in_ws=1; next }
    /^[^[:space:]#]/ {
      if (in_roles) { print acc; printed=1; exit }
      in_ws=0
      in_roles=0
      next
    }
    in_ws && /^  roles:[[:space:]]*\[\][[:space:]]*$/ { print ""; exit }
    in_ws && /^  roles:[[:space:]]*$/ { in_roles=1; acc=""; next }
    in_ws && /^  roles:[[:space:]]*/ {
      sub("^  roles:[[:space:]]*", "")
      gsub("[[:space:]]+$", "")
      print
      exit
    }
    in_roles && /^    -[[:space:]]*/ {
      item=$0
      sub("^    -[[:space:]]*", "", item)
      gsub("[[:space:]]+$", "", item)
      if (acc != "") acc = acc "," item
      else acc = item
      next
    }
    in_roles && /^  [A-Za-z0-9._-]+:/ {
      print acc
      printed=1
      exit
    }
    END { if (in_roles && !printed) print acc }
  ' "$cc_roles_file"
}

cc_proposed_default() {
  case "$1" in
    mode) printf '%s\n' solo ;;
    repositories|roles) printf '%s\n' none ;;
    default_branch|default_branches) printf '%s\n' main ;;
    canonical_url) printf '%s\n' none ;;
    *) printf '%s\n' INCOMPLETE_FIELD >&2; return 1 ;;
  esac
}

cc_resolve_displayed_field() {
  cc_displayed_field=$1
  cc_displayed_value=$2
  if test -n "$cc_displayed_value"; then
    printf '%s\n' "$cc_displayed_value"
    return 0
  fi
  cc_proposed_default "$cc_displayed_field"
}

cc_identity_canonical_from_workspace() {
  cc_identity_file=$1
  cc_identity_name=$(cc_workspace_section_field "$cc_identity_file" workspace name)
  cc_identity_mode=$(cc_workspace_section_field "$cc_identity_file" workspace mode)
  cc_identity_branch=$(cc_workspace_section_field "$cc_identity_file" workspace default_branch)
  cc_identity_roles=$(cc_roles_canonical "$(cc_workspace_roles "$cc_identity_file")")
  cc_identity_kind=$(cc_workspace_section_field "$cc_identity_file" identity kind)
  cc_identity_status=$(cc_workspace_section_field "$cc_identity_file" identity status)
  test -n "$cc_identity_mode" || cc_identity_mode=solo
  test -n "$cc_identity_branch" || cc_identity_branch=main
  {
    printf 'default_branch: %s\n' "$cc_identity_branch"
    printf 'kind: %s\n' "$cc_identity_kind"
    printf 'mode: %s\n' "$cc_identity_mode"
    printf 'name: %s\n' "$cc_identity_name"
    printf 'roles: %s\n' "$cc_identity_roles"
    printf 'status: %s\n' "$cc_identity_status"
    cc_repo_keys=$(cc_workspace_repository_keys "$cc_identity_file")
    if test -z "$cc_repo_keys"; then
      printf 'repositories: none\n'
    else
      printf '%s\n' "$cc_repo_keys" | while IFS= read -r cc_repo_key; do
        test -n "$cc_repo_key" || continue
        cc_repo_url=$(cc_repository_field "$cc_identity_file" "$cc_repo_key" canonical_url)
        cc_repo_branch=$(cc_repository_field "$cc_identity_file" "$cc_repo_key" default_branch)
        test -n "$cc_repo_url" || cc_repo_url=none
        test -n "$cc_repo_branch" || cc_repo_branch=none
        printf 'repository:%s:canonical_url:%s\n' "$cc_repo_key" "$cc_repo_url"
        printf 'repository:%s:default_branch:%s\n' "$cc_repo_key" "$cc_repo_branch"
      done
    fi
  } | sort
}

cc_identity_canonical_from_region_file() {
  cc_region_file=$1
  cc_region_name=$(awk '/^name:[[:space:]]*/ { sub("^name:[[:space:]]*", ""); print; exit }' "$cc_region_file")
  cc_region_mode=$(awk '/^mode:[[:space:]]*/ { sub("^mode:[[:space:]]*", ""); print; exit }' "$cc_region_file")
  cc_region_kind=$(awk '/^kind:[[:space:]]*/ { sub("^kind:[[:space:]]*", ""); print; exit }' "$cc_region_file")
  cc_region_status=$(awk '/^status:[[:space:]]*/ { sub("^status:[[:space:]]*", ""); print; exit }' "$cc_region_file")
  cc_region_branch=$(awk '/^default_branch:[[:space:]]*/ { sub("^default_branch:[[:space:]]*", ""); print; exit }' "$cc_region_file")
  cc_region_roles=$(cc_roles_canonical "$(awk '/^roles:[[:space:]]*/ { sub("^roles:[[:space:]]*", ""); print; exit }' "$cc_region_file")")
  test -n "$cc_region_mode" || cc_region_mode=solo
  test -n "$cc_region_branch" || cc_region_branch=main
  {
    printf 'default_branch: %s\n' "$cc_region_branch"
    printf 'kind: %s\n' "$cc_region_kind"
    printf 'mode: %s\n' "$cc_region_mode"
    printf 'name: %s\n' "$cc_region_name"
    printf 'roles: %s\n' "$cc_region_roles"
    printf 'status: %s\n' "$cc_region_status"
    if grep -Eq '^repositories:[[:space:]]*(none|\{\})[[:space:]]*$' "$cc_region_file"; then
      printf 'repositories: none\n'
    elif grep -Eq '^repositories:[[:space:]]*$' "$cc_region_file"; then
      awk '
        /^repositories:[[:space:]]*$/ { in_repositories=1; next }
        /^[^[:space:]]/ { in_repositories=0 }
        in_repositories && /^  [A-Za-z0-9._-]+:[[:space:]]*$/ {
          key=$1
          sub(":", "", key)
          print key
        }
      ' "$cc_region_file" | while IFS= read -r cc_repo_key; do
        test -n "$cc_repo_key" || continue
        cc_repo_url=$(cc_repository_field "$cc_region_file" "$cc_repo_key" canonical_url)
        cc_repo_branch=$(cc_repository_field "$cc_region_file" "$cc_repo_key" default_branch)
        test -n "$cc_repo_url" || cc_repo_url=none
        test -n "$cc_repo_branch" || cc_repo_branch=none
        printf 'repository:%s:canonical_url:%s\n' "$cc_repo_key" "$cc_repo_url"
        printf 'repository:%s:default_branch:%s\n' "$cc_repo_key" "$cc_repo_branch"
      done
    else
      printf 'repositories: none\n'
    fi
  } | sort
}

cc_identity_region_extract() {
  cc_extract_file=$1
  test -f "$cc_extract_file" || return 1
  cc_start=$(cc_identity_region_start)
  cc_end=$(cc_identity_region_end)
  cc_starts=$(grep -Fxc "$cc_start" "$cc_extract_file" || true)
  cc_ends=$(grep -Fxc "$cc_end" "$cc_extract_file" || true)
  test "$cc_starts" = 1 && test "$cc_ends" = 1 || return 1
  awk -v start="$cc_start" -v end="$cc_end" '
    $0 == start { in_region=1; next }
    $0 == end { in_region=0; next }
    in_region { print }
  ' "$cc_extract_file"
}

cc_identity_region_body_from_file() {
  cc_body_file=$1
  cc_body_name=$(cc_workspace_section_field "$cc_body_file" workspace name)
  cc_body_mode=$(cc_workspace_section_field "$cc_body_file" workspace mode)
  cc_body_branch=$(cc_workspace_section_field "$cc_body_file" workspace default_branch)
  cc_body_roles=$(cc_roles_canonical "$(cc_workspace_roles "$cc_body_file")")
  cc_body_kind=$(cc_workspace_section_field "$cc_body_file" identity kind)
  cc_body_status=$(cc_workspace_section_field "$cc_body_file" identity status)
  test -n "$cc_body_mode" || cc_body_mode=solo
  test -n "$cc_body_branch" || cc_body_branch=main
  printf 'name: %s\n' "$cc_body_name"
  printf 'mode: %s\n' "$cc_body_mode"
  printf 'kind: %s\n' "$cc_body_kind"
  printf 'status: %s\n' "$cc_body_status"
  printf 'default_branch: %s\n' "$cc_body_branch"
  printf 'roles: %s\n' "$cc_body_roles"
  cc_body_keys=$(cc_workspace_repository_keys "$cc_body_file")
  if test -z "$cc_body_keys"; then
    printf 'repositories: none\n'
    return 0
  fi
  printf 'repositories:\n'
  printf '%s\n' "$cc_body_keys" | while IFS= read -r cc_body_key; do
    test -n "$cc_body_key" || continue
    cc_body_url=$(cc_repository_field "$cc_body_file" "$cc_body_key" canonical_url)
    cc_body_repo_branch=$(cc_repository_field "$cc_body_file" "$cc_body_key" default_branch)
    printf '  %s:\n' "$cc_body_key"
    test -n "$cc_body_url" && test "$cc_body_url" != none && printf '    canonical_url: %s\n' "$cc_body_url"
    test -n "$cc_body_repo_branch" && test "$cc_body_repo_branch" != none && printf '    default_branch: %s\n' "$cc_body_repo_branch"
  done
}

cc_identity_region_replace() {
  cc_replace_file=$1
  cc_replace_body=$2
  cc_replace_dest=$3
  cc_start=$(cc_identity_region_start)
  cc_end=$(cc_identity_region_end)
  awk -v start="$cc_start" -v end="$cc_end" -v body_file="$cc_replace_body" '
    BEGIN {
      while ((getline line < body_file) > 0) {
        body = body line "\n"
      }
      close(body_file)
    }
    $0 == start { print; printf "%s", body; skip=1; next }
    $0 == end { skip=0; print; next }
    skip { next }
    { print }
  ' "$cc_replace_file" > "$cc_replace_dest"
}

cc_identity_projection_applies() {
  cc_applies_root=$1
  test -f "$cc_applies_root/workspace.yaml" || return 1
  cc_applies_kind=$(cc_workspace_section_field "$cc_applies_root/workspace.yaml" identity kind)
  test "$cc_applies_kind" = instantiated-workspace
}

cc_validate_identity_projection() {
  cc_validate_root=$1
  if test -f "$cc_validate_root/.runtime/identity-publish.journal"; then
    printf '%s\n' projection-mismatch
    return 1
  fi
  test -f "$cc_validate_root/workspace.yaml" || { printf '%s\n' projection-mismatch; return 1; }
  cc_expected=$(cc_identity_canonical_from_workspace "$cc_validate_root/workspace.yaml")
  cc_identity_region_relpaths | while IFS= read -r cc_rel; do
    cc_summary="$cc_validate_root/$cc_rel"
    test -f "$cc_summary" || { printf '%s\n' projection-mismatch; exit 1; }
    cc_region_tmp="$cc_summary.region.$$"
    if ! cc_identity_region_extract "$cc_summary" > "$cc_region_tmp"; then
      rm -f "$cc_region_tmp"
      printf '%s\n' projection-mismatch
      exit 1
    fi
    if grep -E 'path:[[:space:]]*/|path:[[:space:]]*\.\./' "$cc_region_tmp" >/dev/null 2>&1; then
      rm -f "$cc_region_tmp"
      printf '%s\n' projection-mismatch
      exit 1
    fi
    cc_actual=$(cc_identity_canonical_from_region_file "$cc_region_tmp")
    rm -f "$cc_region_tmp"
    test "$cc_expected" = "$cc_actual" || { printf '%s\n' projection-mismatch; exit 1; }
  done || { printf '%s\n' projection-mismatch; return 1; }
  printf '%s\n' identity-projection-ok
}

cc_entry_preflight() {
  cc_preflight_root=$1
  if cc_identity_projection_applies "$cc_preflight_root"; then
    cc_validate_identity_projection "$cc_preflight_root" || return 1
  fi
  printf '%s\n' entry-preflight-ok
}

cc_write_preflight() {
  cc_preflight_root=$1
  if cc_identity_projection_applies "$cc_preflight_root"; then
    cc_validate_identity_projection "$cc_preflight_root" || return 1
  fi
  printf '%s\n' write-preflight-ok
}

cc_identity_publish_files() {
  cc_publish_root=$1
  cc_workspace_source=${2:-$cc_publish_root/workspace.yaml}
  cc_publish_workspace_dest=${3:-}
  cc_journal="$cc_publish_root/.runtime/identity-publish.journal"
  mkdir -p "$cc_publish_root/.runtime"
  cc_body="$cc_journal.body.$$"
  cc_identity_region_body_from_file "$cc_workspace_source" > "$cc_body" || { rm -f "$cc_body"; printf '%s\n' projection-mismatch >&2; return 1; }
  cc_fail_after=${CC_IDENTITY_PUBLISH_FAIL_AFTER:-0}
  cc_moved=0
  cc_list="$cc_journal.list.$$"
  : > "$cc_list"
  cc_ok=1
  cc_identity_region_relpaths | while IFS= read -r cc_rel; do
    printf '%s\n' "$cc_publish_root/$cc_rel"
  done > "$cc_list.files.$$"
  while IFS= read -r cc_summary; do
    test -f "$cc_summary" || { cc_ok=0; break; }
    cc_tmp="$cc_summary.tmp.$$"
    cc_bak="$cc_summary.bak.$$"
    if ! cc_identity_region_extract "$cc_summary" >/dev/null; then
      cc_ok=0
      break
    fi
    cp "$cc_summary" "$cc_bak"
    if ! cc_identity_region_replace "$cc_summary" "$cc_body" "$cc_tmp"; then
      cc_ok=0
      break
    fi
    printf '%s\t%s\t%s\n' "$cc_tmp" "$cc_summary" "$cc_bak" >> "$cc_list"
  done < "$cc_list.files.$$"
  rm -f "$cc_list.files.$$" "$cc_body"
  if test -n "$cc_publish_workspace_dest" && test "$cc_ok" -eq 1; then
    cc_ws_bak="$cc_publish_workspace_dest.bak.$$"
    cp "$cc_publish_workspace_dest" "$cc_ws_bak"
    printf '%s\t%s\t%s\n' "$cc_workspace_source" "$cc_publish_workspace_dest" "$cc_ws_bak" >> "$cc_list"
  fi
  if test "$cc_ok" -ne 1; then
    while IFS='	' read -r cc_tmp cc_summary cc_bak; do
      rm -f "$cc_tmp"
      test -n "$cc_bak" && test -f "$cc_bak" && mv "$cc_bak" "$cc_summary"
    done < "$cc_list"
    rm -f "$cc_list"
    printf '%s\n' projection-mismatch >&2
    return 1
  fi
  printf 'status: committing\n' > "$cc_journal"
  while IFS='	' read -r cc_tmp cc_summary cc_bak; do
    test -n "$cc_tmp" || continue
    mv "$cc_tmp" "$cc_summary" || cc_ok=0
    cc_moved=$((cc_moved + 1))
    if test "$cc_fail_after" -gt 0 && test "$cc_moved" -ge "$cc_fail_after"; then
      cc_ok=0
      break
    fi
  done < "$cc_list"
  if test "$cc_ok" -ne 1; then
    while IFS='	' read -r cc_tmp cc_summary cc_bak; do
      rm -f "$cc_tmp"
      test -n "$cc_bak" && test -f "$cc_bak" && mv "$cc_bak" "$cc_summary"
    done < "$cc_list"
    rm -f "$cc_list" "$cc_journal"
    printf '%s\n' IDENTITY_PUBLISH_INTERRUPTED >&2
    return 1
  fi
  while IFS='	' read -r cc_tmp cc_summary cc_bak; do
    rm -f "$cc_bak"
  done < "$cc_list"
  rm -f "$cc_list" "$cc_journal"
}

cc_publish_identity_regions() {
  cc_identity_publish_files "$1"
}

cc_workspace_write_repository() {
  cc_ws_src=$1
  cc_ws_dest=$2
  cc_ws_key=$3
  cc_ws_url=$4
  cc_ws_branch=$5
  awk -v key="$cc_ws_key" -v url="$cc_ws_url" -v branch="$cc_ws_branch" '
    function emit_key() {
      print "  " key ":"
      if (url != "" && url != "none") print "    canonical_url: " url
      if (branch != "" && branch != "none") print "    default_branch: " branch
    }
    /^repositories:[[:space:]]*\{\}[[:space:]]*$/ {
      print "repositories:"
      emit_key()
      added=1
      next
    }
    /^repositories:[[:space:]]*$/ { in_repositories=1; print; next }
    in_repositories && /^  [A-Za-z0-9._-]+:[[:space:]]*$/ {
      current=$1
      sub(":", "", current)
      if (current == key) { skipping=1; next }
      skipping=0
      if (!added) { emit_key(); added=1 }
      print
      next
    }
    in_repositories && skipping && /^    / { next }
    in_repositories && /^[^[:space:]#]/ {
      if (!added) { emit_key(); added=1 }
      in_repositories=0
      skipping=0
      print
      next
    }
    { if (!(in_repositories && skipping)) print }
    END { if (in_repositories && !added) emit_key() }
  ' "$cc_ws_src" > "$cc_ws_dest"
}

cc_workspace_apply_displayed_identity() {
  cc_ws_src=$1
  cc_ws_dest=$2
  cc_ws_mode=$3
  cc_ws_roles=$4
  cc_ws_branch=$5
  cc_ws_status=$6
  awk -v mode="$cc_ws_mode" -v roles="$cc_ws_roles" -v branch="$cc_ws_branch" -v status="$cc_ws_status" '
    function emit_roles() {
      if (roles == "none" || roles == "") {
        print "  roles: none"
        return
      }
      print "  roles:"
      n = split(roles, items, ",")
      for (i = 1; i <= n; i++) {
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", items[i])
        if (items[i] != "") print "    - " items[i]
      }
    }
    function flush_workspace() {
      if (in_workspace) {
        if (!saw_mode) print "  mode: " mode
        if (!saw_branch) print "  default_branch: " branch
        if (!saw_roles) emit_roles()
      }
      in_workspace=0
    }
    function flush_identity() {
      if (in_identity && !saw_status) print "  status: " status
      in_identity=0
    }
    /^workspace:/ { flush_identity(); in_workspace=1; in_identity=0; skip_role_items=0; print; next }
    /^identity:/ { flush_workspace(); in_identity=1; skip_role_items=0; print; next }
    /^[^[:space:]#]/ { flush_workspace(); flush_identity(); skip_role_items=0; print; next }
    in_workspace && skip_role_items && /^    -/ { next }
    in_workspace && skip_role_items { skip_role_items=0 }
    in_workspace && /^  mode:/ { print "  mode: " mode; saw_mode=1; next }
    in_workspace && /^  default_branch:/ { print "  default_branch: " branch; saw_branch=1; next }
    in_workspace && /^  roles:/ { emit_roles(); saw_roles=1; skip_role_items=1; next }
    in_identity && /^  status:/ { print "  status: " status; saw_status=1; next }
    { print }
    END { flush_workspace(); flush_identity() }
  ' "$cc_ws_src" > "$cc_ws_dest"
}

cc_identity_acceptance_card() {
  cc_card_root=$1
  cc_card_mode=$(cc_resolve_displayed_field mode "$(cc_workspace_section_field "$cc_card_root/workspace.yaml" workspace mode)")
  cc_card_repos=$(cc_workspace_repository_keys "$cc_card_root/workspace.yaml")
  test -n "$cc_card_repos" || cc_card_repos=$(cc_resolve_displayed_field repositories "")
  cc_card_roles=$(cc_resolve_displayed_field roles "$(cc_roles_canonical "$(cc_workspace_roles "$cc_card_root/workspace.yaml")")")
  cc_card_branch=$(cc_resolve_displayed_field default_branch "$(cc_workspace_section_field "$cc_card_root/workspace.yaml" workspace default_branch)")
  printf 'Action: identity-acceptance\n'
  printf 'Target: workspace identity\n'
  printf 'Observed state: current session; identity remains uninitialized until confirmation\n'
  printf 'Proposed defaults:\n'
  printf '  mode: %s\n' "$cc_card_mode"
  printf '  repositories: %s\n' "$cc_card_repos"
  printf '  roles: %s\n' "$cc_card_roles"
  printf '  default branches: %s\n' "$cc_card_branch"
  printf 'Immediate effects: workspace.accept_identity\n'
  printf 'Later authorized effects: none from this confirmation\n'
  printf 'Will change after confirmation: accepted identity fields and identity regions in WORKSPACE.md, PROJECT.md, and INDEX.md\n'
  printf 'Will not change: Product Knowledge outside the identity region, Git, leases, clone, create-empty, execution, or delivery\n'
  printf 'Risks/open decisions: omitted fields are recorded as the displayed defaults; fields with no default remain incomplete\n'
  printf 'Confirmation requested: Confirm identity acceptance for this workspace in the current session.\n'
}

cc_repository_registration_card() {
  cc_card_root=$1
  cc_card_key=$2
  cc_card_url=${3:-}
  cc_card_branch=${4:-}
  test -n "$cc_card_url" || cc_card_url=$(cc_resolve_displayed_field canonical_url "")
  test -n "$cc_card_branch" || cc_card_branch=$(cc_resolve_displayed_field default_branch "")
  printf 'Action: repository-registration\n'
  printf 'Target: %s\n' "${cc_card_key:-incomplete}"
  printf 'Observed state: shared identity card for a later clone or create-empty gate; no destination is created now\n'
  printf 'Proposed defaults:\n'
  printf '  logical key: %s\n' "${cc_card_key:-incomplete}"
  printf '  canonical URL: %s\n' "$cc_card_url"
  printf '  default branch: %s\n' "$cc_card_branch"
  printf 'Immediate effects: workspace.register_repository\n'
  printf 'Later authorized effects: repository-bootstrap; repository-create-empty (reserved, not activated); execute-plan\n'
  printf 'Will change after confirmation: shared logical key, optional URL, optional branch, and identity regions\n'
  printf 'Will not change: repositories.local.yaml, clone, create-empty, execution, delivery, or authored Product Knowledge outside the identity region\n'
  printf 'Risks/open decisions: this card does not authorize clone or create-empty; those remain later gates\n'
  printf 'Confirmation requested: Confirm repository registration for this logical key in the current session.\n'
}

cc_accept_identity() {
  cc_accept_root=$1
  cc_accept_confirmation=${2:-}
  cc_accept_mode=${3:-}
  cc_accept_repos=${4:-}
  cc_accept_roles=${5:-}
  cc_accept_branch=${6:-}
  cc_accept_invented=${7:-}
  test "$cc_accept_confirmation" = confirmed || { printf '%s\n' GATE_REQUIRED >&2; return 1; }
  test -z "$cc_accept_invented" || { printf '%s\n' FIELD_INVENTED >&2; return 1; }
  cc_write_preflight "$cc_accept_root" >/dev/null || return 1
  cc_accept_name=$(cc_workspace_section_field "$cc_accept_root/workspace.yaml" workspace name)
  test -n "$cc_accept_name" || { printf '%s\n' INCOMPLETE_FIELD >&2; return 1; }
  cc_accept_mode=$(cc_resolve_displayed_field mode "$cc_accept_mode") || return 1
  if test -z "$cc_accept_repos"; then
    cc_accept_repos=$(cc_resolve_displayed_field repositories "") || return 1
  fi
  cc_accept_roles=$(cc_roles_canonical "$(cc_resolve_displayed_field roles "$cc_accept_roles")") || return 1
  cc_accept_branch=$(cc_resolve_displayed_field default_branch "$cc_accept_branch") || return 1
  cc_ws_tmp="$cc_accept_root/workspace.yaml.tmp.$$"
  cc_workspace_apply_displayed_identity "$cc_accept_root/workspace.yaml" "$cc_ws_tmp" "$cc_accept_mode" "$cc_accept_roles" "$cc_accept_branch" accepted
  cc_identity_publish_files "$cc_accept_root" "$cc_ws_tmp" "$cc_accept_root/workspace.yaml" || { rm -f "$cc_ws_tmp"; return 1; }
  rm -f "$cc_ws_tmp"
}

cc_register_repository() {
  cc_reg_root=$1
  cc_reg_key=$2
  cc_reg_url=${3:-}
  cc_reg_branch=${4:-}
  cc_reg_confirmation=${5:-}
  cc_reg_invented=${6:-}
  test "$cc_reg_confirmation" = confirmed || { printf '%s\n' GATE_REQUIRED >&2; return 1; }
  test -z "$cc_reg_invented" || { printf '%s\n' FIELD_INVENTED >&2; return 1; }
  test -n "$cc_reg_key" || { printf '%s\n' INCOMPLETE_FIELD >&2; return 1; }
  cc_safe_id "$cc_reg_key" || { printf '%s\n' INVALID_REPOSITORY_KEY >&2; return 1; }
  case "$cc_reg_key" in */*|.*|/*) printf '%s\n' INVALID_REPOSITORY_KEY >&2; return 1 ;; esac
  test -z "$cc_reg_url" || cc_repository_canonical_safe "$cc_reg_url" || { printf '%s\n' UNSAFE_CANONICAL_URL >&2; return 1; }
  test -n "$cc_reg_branch" || cc_reg_branch=$(cc_resolve_displayed_field default_branch "") || return 1
  cc_safe_id "$cc_reg_branch" || { printf '%s\n' INVALID_DEFAULT_BRANCH >&2; return 1; }
  case "$cc_reg_url" in *'/home/'*|*'path:'*) printf '%s\n' MACHINE_PATH_FORBIDDEN >&2; return 1 ;; esac
  cc_write_preflight "$cc_reg_root" >/dev/null || return 1
  cc_ws_tmp="$cc_reg_root/workspace.yaml.tmp.$$"
  cc_workspace_write_repository "$cc_reg_root/workspace.yaml" "$cc_ws_tmp" "$cc_reg_key" "$cc_reg_url" "$cc_reg_branch"
  if grep -E 'path:[[:space:]]' "$cc_ws_tmp" >/dev/null 2>&1; then
    rm -f "$cc_ws_tmp"
    printf '%s\n' MACHINE_PATH_FORBIDDEN >&2
    return 1
  fi
  cc_identity_publish_files "$cc_reg_root" "$cc_ws_tmp" "$cc_reg_root/workspace.yaml" || { rm -f "$cc_ws_tmp"; return 1; }
  rm -f "$cc_ws_tmp"
}

cc_locked_effect_identifiers() {
  printf '%s\n' \
    workspace.accept_identity \
    workspace.register_repository \
    repository-bootstrap \
    repository-create-empty \
    git.commit \
    delivery.push
}

cc_validate_delegated_effects() {
  cc_approved_file=$1
  cc_delegated_file=$2
  test -f "$cc_approved_file" && test -f "$cc_delegated_file" || { printf '%s\n' UNDECLARED_EFFECT; return 1; }
  while IFS= read -r cc_effect; do
    test -n "$cc_effect" || continue
    case "$cc_effect" in
      git.init|git.clone|repository-create-empty)
        printf '%s\n' UNDECLARED_EFFECT
        return 1
        ;;
    esac
    grep -Fx "$cc_effect" "$cc_approved_file" >/dev/null 2>&1 || { printf '%s\n' UNDECLARED_EFFECT; return 1; }
    cc_locked_effect_identifiers | grep -Fx "$cc_effect" >/dev/null 2>&1 || { printf '%s\n' UNDECLARED_EFFECT; return 1; }
  done < "$cc_delegated_file"
  printf '%s\n' delegated-effects-ok
}

cc_insert_identity_region() {
  cc_insert_file=$1
  cc_insert_body=$2
  test -f "$cc_insert_file" || return 1
  if cc_identity_region_extract "$cc_insert_file" >/dev/null; then
    return 0
  fi
  cc_start=$(cc_identity_region_start)
  cc_end=$(cc_identity_region_end)
  cc_tmp="$cc_insert_file.tmp.$$"
  awk -v start="$cc_start" -v end="$cc_end" -v body_file="$cc_insert_body" '
    BEGIN {
      while ((getline line < body_file) > 0) {
        body = body line "\n"
      }
      close(body_file)
    }
    NR == 1 {
      print
      print ""
      print start
      printf "%s", body
      print end
      print ""
      inserted=1
      next
    }
    { print }
  ' "$cc_insert_file" > "$cc_tmp"
  mv "$cc_tmp" "$cc_insert_file"
}

cc_remove_identity_region() {
  cc_remove_file=$1
  test -f "$cc_remove_file" || return 1
  cc_start=$(cc_identity_region_start)
  cc_end=$(cc_identity_region_end)
  cc_tmp="$cc_remove_file.tmp.$$"
  awk -v start="$cc_start" -v end="$cc_end" '
    $0 == start { skip=1; next }
    $0 == end { skip=0; next }
    skip { next }
    { print }
  ' "$cc_remove_file" > "$cc_tmp"
  mv "$cc_tmp" "$cc_remove_file"
}

# Stable aliases used by host adapters and semantic fixtures.
cc_resolve_repository() { cc_resolve_repository_binding "$@"; }
cc_validate_binding() { cc_validate_repository_binding "$@"; }
cc_repository_bootstrap() { cc_bootstrap_repository "$@"; }
