#!/bin/sh
# Context Circuit v0.6 external-surface (publication) semantic suite.
# The load-bearing check is isolation: the core workflow surfaces reference
# nothing here, which is what makes a publication orthogonal (INV-EXTERNAL-01).
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

W="$ROOT/wrapper"

# --- invariants and owners present ---
inv="$W/contracts/invariants.yaml"
require_file "$inv"
contains "$inv" "INV-EXTERNAL-01"
contains "$inv" "INV-EXTERNAL-02"
contains "$inv" "INV-EXTERNAL-03"
contains "$inv" "orthogonal to the core"
contains "$inv" "export-only"
contains "$inv" "self-contained"
contains "$inv" "no access to the workspace"
for concern in external_surface publication_config publication_record; do
	contains "$inv" "$concern:"
done

# --- config schema: credential-free, export-only, manual-only ---
cfg="$W/contracts/schemas/publication-config.yaml"
require_file "$cfg"
contains "$cfg" "concern: publication_config"
contains "$cfg" "publication/<name>/config.yaml"
contains "$cfg" "values: [export]"
contains "$cfg" "values: [manual]"
contains "$cfg" "credential-free"

# --- record schema: under the publication, idempotent, not an execution input ---
rec="$W/contracts/schemas/publication-record.yaml"
require_file "$rec"
contains "$rec" "concern: publication_record"
contains "$rec" "publication/<name>/published/<plan-id>.yaml"
contains "$rec" "synced_digest"
contains "$rec" "INV-PLAN-01"
contains "$rec" "never under plans/"
contains "$rec" "not_an_execution_input"
contains "$rec" "any field on plan.yaml holding an external id"

# --- cc-publish skill: manual, export, one-way, self-contained, host/MCP ---
sk="$ROOT/.agents/skills/cc-publish/SKILL.md"
require_file "$sk"
contains "$sk" "orthogonal to the core workflow"
contains "$sk" "never through the runtime engine"
contains "$sk" "host-blocked"
contains "$sk" "never encode it as nesting"
contains "$sk" "one-way"
contains "$sk" "self-contained"
contains "$sk" "heard of Context Circuit"
contains "$sk" "Never put in any external field"
contains "$sk" "[<plan-number>] <plan title>"
contains "$sk" "id slug appears nowhere"
contains "$sk" "never write under"

# --- manifest registers the schemas and the workspace-owned folder ---
man="$W/manifest.yaml"
contains "$man" "publication-config: [1]"
contains "$man" "publication-record: [1]"
contains "$man" "- publication/"

# --- ISOLATION: the core workflow references nothing here (INV-EXTERNAL-01) ---
# Runtime is host-neutral and gains no publish action.
eng="$W/runtime/engine.sh"
not_contains "$eng" "cc_publish"
not_contains "$eng" "publication/"
not_contains "$eng" "publication-config"
# The delivery boundary and coordinator route nothing to a publication.
not_contains "$W/adapters/WORKFLOW.md" "cc-publish"
not_contains "$W/adapters/WORKFLOW.md" "publication/"
not_contains "$ROOT/agents/coordinator.md" "cc-publish"
not_contains "$ROOT/agents/coordinator.md" "publication/"
# No core-workflow skill triggers or depends on the surface.
for sk_core in cc-plan cc-execute cc-run-stack cc-verify cc-complete cc-deliver cc-archive; do
	f="$ROOT/.agents/skills/$sk_core/SKILL.md"
	not_contains "$f" "cc-publish"
	not_contains "$f" "publication/"
done

# --- Path A: git delivery no longer says "publish/publication" ---
not_contains "$W/adapters/WORKFLOW.md" "publish"
not_contains "$ROOT/.agents/skills/cc-deliver/SKILL.md" "publication"

pass 'external-surface'
