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
  for path in wrapper .agents/skills agents AGENTS.md WORKFLOW.md README.md; do printf '%s\n' "$path"; done
}
