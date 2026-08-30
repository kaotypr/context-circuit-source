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
for concern in external_surface publication_config publication_intent publication_record; do
	contains "$inv" "$concern:"
done
# INV-EXTERNAL-02 clarified: a display-only preview read is not inbound flow, and a
# publication writes within its own publication/<name>/ folder (config, intent/, published/).
contains "$inv" "is not an inbound flow"
contains "$inv" "publication/<name>/ folder"
contains "$inv" "a read that ends at the screen"

# --- config schema: credential-free, export-only, manual-only ---
cfg="$W/contracts/schemas/publication-config.yaml"
require_file "$cfg"
contains "$cfg" "concern: publication_config"
contains "$cfg" "publication/<name>/config.yaml"
contains "$cfg" "values: [export]"
contains "$cfg" "values: [manual]"
contains "$cfg" "credential-free"
# instructions apply to every kind (language, tone, optional field enrichment)
contains "$cfg" "instructions:"
# preview mode: opt-in, display-only drift read (off by default)
contains "$cfg" "preview:"
contains "$cfg" "drift_read"
contains "$cfg" "display-only"
contains "$cfg" "intent/<plan-id>.yaml"

# --- intent schema: user-owned field layer, canonical minutes, no d/w ---
intent="$W/contracts/schemas/publication-intent.yaml"
require_file "$intent"
contains "$intent" "concern: publication_intent"
contains "$intent" "publication/<name>/intent/<plan-id>.yaml"
contains "$intent" "estimate_minutes"
contains "$intent" "then human-owned"
contains "$intent" "silently overwrites a human edit"
contains "$intent" "d and w are rejected"
contains "$intent" "INV-SEC-01"
# intent holds values, never identity or plan-owned facts
contains "$intent" "an external id or url"

# --- record schema (plan kind): under the publication, idempotent ---
rec="$W/contracts/schemas/publication-record.yaml"
require_file "$rec"
contains "$rec" "concern: publication_record"
contains "$rec" "publication/<name>/published/<plan-id>.yaml"
contains "$rec" "synced_digest"
contains "$rec" "INV-PLAN-01"
contains "$rec" "never under plans/"
contains "$rec" "not_an_execution_input"
contains "$rec" "any field on plan.yaml holding an external id"
# fields: snapshot — last-published values, additive/backward-compatible, drives field idempotency
contains "$rec" "Last-published snapshot"
contains "$rec" "estimate_minutes"
contains "$rec" "field-only"
contains "$rec" "backward-compatible"

# --- record schema (thread kind): discussion-safe, ids/timestamps only ---
trec="$W/contracts/schemas/publication-thread-record.yaml"
require_file "$trec"
contains "$trec" "concern: publication_thread_record"
contains "$trec" "publication/<name>/published/<plan-id>.yaml"
contains "$trec" "parent_ts"
contains "$trec" "discussion_safe"
contains "$trec" "not_an_execution_input"

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
# instructions honored (language, tone, field enrichment)
contains "$sk" "Authoring: instructions"
contains "$sk" "best-effort estimates"
# intent layer + estimate unit/format
contains "$sk" "Field intent and estimates"
contains "$sk" "publication/<name>/intent/<plan-id>.yaml"
contains "$sk" "estimate_minutes"
contains "$sk" "never silently overwrites a human edit"
contains "$sk" "input and display only"
contains "$sk" "are rejected"
# consult-before-publish preview mode: a mode of the skill, default local diff, opt-in drift read
contains "$sk" "Preview and consult"
contains "$sk" "dry-run"
contains "$sk" "consult-first"
contains "$sk" "intent vs last-published snapshot"
contains "$sk" "drift_read: true"
contains "$sk" "display only"
contains "$sk" "Publish on an explicit go"
# idempotency now covers fields, report notes field-only updates
contains "$sk" "field-only update"
contains "$sk" "nothing was pushed"
# the thread kind
contains "$sk" "Publish (the \`thread\` kind)"
contains "$sk" "[thread] [<plan-number>]"
contains "$sk" "One reply per question, fully described"
contains "$sk" "never delete a"

# --- manifest registers the schemas and the workspace-owned folder ---
man="$W/manifest.yaml"
contains "$man" "publication-config: [1]"
contains "$man" "publication-intent: [1]"
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
