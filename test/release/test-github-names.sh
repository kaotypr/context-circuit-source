#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

binding="$ROOT/release/binding.yaml"
workspace="$ROOT/workspace.yaml"
workflow="$ROOT/.github/workflows/publish-template.yml"
pubcfg="$ROOT/publication/plans-github/config.yaml"

require_file "$binding"
require_file "$workspace"
require_file "$workflow"
require_file "$pubcfg"

# Operational GitHub slugs after the source/template name swap.
contains "$binding" 'destination_repo: kaotypr/context-circuit'
not_contains "$binding" 'kaotypr/context-circuit-template'
not_contains "$binding" 'context-circuit-template'

contains "$workspace" 'canonical_url: https://github.com/kaotypr/context-circuit-source.git'

contains "$workflow" 'repository: kaotypr/context-circuit'
contains "$workflow" '--repo kaotypr/context-circuit'
not_contains "$workflow" 'kaotypr/context-circuit-template'
not_contains "$workflow" 'context-circuit-template'

contains "$pubcfg" 'repository: kaotypr/context-circuit-source'

# Live GitHub names when credentials are available. Redirects from a retired
# slug are acceptable; the resolved nameWithOwner must be the new destination.
if command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then
	src=$(gh repo view kaotypr/context-circuit-source --json nameWithOwner -q .nameWithOwner)
	assert_eq "$src" 'kaotypr/context-circuit-source'
	tpl=$(gh repo view kaotypr/context-circuit --json nameWithOwner -q .nameWithOwner)
	assert_eq "$tpl" 'kaotypr/context-circuit'
	old=$(gh repo view kaotypr/context-circuit-template --json nameWithOwner -q .nameWithOwner 2>/dev/null || true)
	if [ -n "$old" ]; then
		test "$old" != 'kaotypr/context-circuit-template' \
			|| fail 'retired slug kaotypr/context-circuit-template still names a distinct repository'
	fi
fi

pass 'github operational repository names'
