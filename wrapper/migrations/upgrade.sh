#!/bin/sh
# Sourced migration helpers. They classify state and describe a wrapper-only
# change; callers retain the human gate and preserve workspace data.

cc_migration_classify() {
  case "${1:-legacy-unknown}" in
    legacy-unknown|'') printf '%s\n' legacy-unknown ;;
    1.0.*) printf '%s\n' compatible ;;
    0.*) printf '%s\n' migration-needed ;;
    *) printf '%s\n' blocked ;;
  esac
}

cc_migration_preserves() {
  for path in workspace.yaml context sources plans .runtime repositories.local.yaml repositories; do printf '%s\n' "$path"; done
}

cc_rollback_scope() {
  for path in wrapper .agents/skills agents AGENTS.md CLAUDE.md WORKFLOW.md README.md; do printf '%s\n' "$path"; done
}

cc_migration_packet_fixtures() {
  for path in \
    wrapper/contracts/schemas/context-receipt.yaml \
    wrapper/contracts/schemas/delegation.yaml \
    wrapper/contracts/schemas/child-start.yaml \
    wrapper/contracts/schemas/plan.yaml \
    wrapper/contracts/schemas/task.yaml \
    docs/gates.md docs/templates/plan.md docs/templates/plan.yaml \
    docs/templates/task.md docs/templates/prd.md; do
    printf '%s\n' "$path"
  done
}

cc_migration_receipt_status() {
  cc_migration_receipt=${1:-}
  test -f "$cc_migration_receipt" || { printf '%s\n' legacy-unknown; return 0; }
  cc_migration_receipt_version=$(sed -n 's/^wrapper_version: //p' "$cc_migration_receipt" | head -n 1)
  cc_migration_classify "${cc_migration_receipt_version:-legacy-unknown}"
}

cc_migrate_identity_regions() {
  cc_migrate_root=$1
  test -f "$cc_migrate_root/workspace.yaml" || return 1
  cc_migrate_status=$(cc_workspace_section_field "$cc_migrate_root/workspace.yaml" identity status)
  cc_body="$cc_migrate_root/.identity-region.body.$$"
  cc_identity_region_body_from_file "$cc_migrate_root/workspace.yaml" > "$cc_body" || { rm -f "$cc_body"; return 1; }
  for cc_summary in WORKSPACE.md PROJECT.md INDEX.md; do
    test -f "$cc_migrate_root/context/$cc_summary" || continue
    cc_insert_identity_region "$cc_migrate_root/context/$cc_summary" "$cc_body"
  done
  rm -f "$cc_body"
  cc_migrate_status_after=$(cc_workspace_section_field "$cc_migrate_root/workspace.yaml" identity status)
  test "$cc_migrate_status" = "$cc_migrate_status_after" || return 1
}

cc_rollback_identity_regions() {
  cc_rollback_root=$1
  test -f "$cc_rollback_root/workspace.yaml" || return 1
  cc_rollback_status=$(cc_workspace_section_field "$cc_rollback_root/workspace.yaml" identity status)
  for cc_summary in WORKSPACE.md PROJECT.md INDEX.md; do
    test -f "$cc_rollback_root/context/$cc_summary" || continue
    cc_remove_identity_region "$cc_rollback_root/context/$cc_summary"
  done
  cc_rollback_status_after=$(cc_workspace_section_field "$cc_rollback_root/workspace.yaml" identity status)
  test "$cc_rollback_status" = "$cc_rollback_status_after" || return 1
}

cc_migration_yaml_field() {
  awk -F': ' -v field="$2" '
    {
      key=$1
      sub(/^[ \t]+/, "", key)
      if (key == field) { print $2; exit }
    }
  ' "$1"
}

cc_migration_evidence_fields() {
  printf '%s\n' required_layer produced_layer observed_layer evidence_ref
}

cc_migration_evidence_rules() {
  printf '%s\n' 'completed-historical: preserve-without-rewrite'
  printf '%s\n' 'unfinished: explicit-mapping-required'
  printf '%s\n' 'invent-evidence: forbidden'
}

cc_legacy_evidence_classify() {
  cc_legacy_file=$1
  test -f "$cc_legacy_file" || { printf '%s\n' mapping-required; return 0; }
  cc_legacy_status=$(cc_migration_yaml_field "$cc_legacy_file" plan_status)
  test -n "$cc_legacy_status" || cc_legacy_status=$(cc_migration_yaml_field "$cc_legacy_file" status)
  cc_legacy_required=$(cc_migration_yaml_field "$cc_legacy_file" required_layer)
  case "$cc_legacy_status" in
    done|completed)
      printf '%s\n' legacy-completed-readable
      return 0
      ;;
  esac
  if test -z "$cc_legacy_required" || test "$cc_legacy_required" = omitted; then
    printf '%s\n' mapping-required
    return 0
  fi
  printf '%s\n' mapping-present
}

cc_migrate_legacy_evidence() {
  cc_legacy_src=$1
  cc_legacy_dest=$2
  test -f "$cc_legacy_src" || return 1
  cp "$cc_legacy_src" "$cc_legacy_dest"
}

cc_legacy_evidence_invented() {
  cc_legacy_before=$1
  cc_legacy_after=$2
  test -f "$cc_legacy_before" && test -f "$cc_legacy_after" || { printf '%s\n' EVIDENCE_INVENTED; return 0; }
  for cc_legacy_field in required_layer produced_layer observed_layer evidence_ref; do
    cc_legacy_before_value=$(cc_migration_yaml_field "$cc_legacy_before" "$cc_legacy_field")
    cc_legacy_after_value=$(cc_migration_yaml_field "$cc_legacy_after" "$cc_legacy_field")
    if test -z "$cc_legacy_before_value" || test "$cc_legacy_before_value" = omitted; then
      if test -n "$cc_legacy_after_value" && test "$cc_legacy_after_value" != omitted; then
        printf '%s\n' EVIDENCE_INVENTED
        return 0
      fi
    fi
  done
  printf '%s\n' EVIDENCE_NOT_INVENTED
}

cc_require_explicit_evidence_mapping() {
  cc_legacy_class=$(cc_legacy_evidence_classify "$1")
  case "$cc_legacy_class" in
    mapping-required)
      printf '%s\n' EVIDENCE_MAPPING_REQUIRED
      return 1
      ;;
    legacy-completed-readable)
      printf '%s\n' LEGACY_COMPLETED_READABLE
      return 0
      ;;
    *)
      printf '%s\n' EVIDENCE_MAPPING_PRESENT
      return 0
      ;;
  esac
}
