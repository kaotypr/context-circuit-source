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
