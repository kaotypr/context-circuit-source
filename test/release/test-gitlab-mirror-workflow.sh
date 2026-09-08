#!/bin/sh
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"

wf="$ROOT/.github/workflows/publish-template.yml"
require_file "$wf"

# GitHub push of main and the version tag remains.
contains "$wf" 'git -C .template-repo push origin main'
contains "$wf" 'git -C .template-repo push origin "v$V"'

# GitHub Release remains.
contains "$wf" 'gh release create'

# GitLab mirror uses the same .template-repo checkout and host.
contains "$wf" 'gitlab.sicepat.tech'
contains "$wf" 'git -C .template-repo remote add gitlab'
contains "$wf" 'git -C .template-repo push gitlab main'
contains "$wf" 'git -C .template-repo push gitlab "v$V"'

# Operator wiring: GitHub Actions variable and secret names only.
contains "$wf" 'vars.GITLAB_TEMPLATE_PROJECT'
contains "$wf" 'secrets.GITLAB_TEMPLATE_TOKEN'
contains "$wf" 'oauth2'

# Same tree: assembly runs once; GitLab does not re-invoke it.
count=$(grep -cF 'scripts/publish-template.sh' "$wf")
assert_eq "$count" 1

# GitLab Release uses the same host and token, after the git push.
contains "$wf" 'api/v4/projects'
contains "$wf" '/releases'
contains "$wf" 'packages/generic/context-circuit'
contains "$wf" 'Create GitLab Release'

# GitHub origin push precedes the GitLab mirror.
gh_push=$(grep -nF 'git -C .template-repo push origin main' "$wf" | head -n1 | cut -d: -f1)
gl_push=$(grep -nF 'git -C .template-repo push gitlab main' "$wf" | head -n1 | cut -d: -f1)
test "$gh_push" -lt "$gl_push" || fail "GitLab push must follow GitHub origin push"

# GitLab Release follows the GitLab git push.
gl_rel=$(grep -nF 'Create GitLab Release' "$wf" | head -n1 | cut -d: -f1)
test "$gl_push" -lt "$gl_rel" || fail "GitLab Release must follow GitLab git push"

# GitHub Release remains in the job (order vs GitLab is not load-bearing).
grep -qF 'gh release create' "$wf" || fail 'missing gh release create'

# Every GitLab-host step must fail the job on error and must not trace token URLs.
gitlab_steps=$(awk '
	/^      - name:/ {
		if (found) print buf
		buf = $0
		found = 0
		next
	}
	{
		if (buf != "") buf = buf "\n" $0
	}
	/gitlab\.sicepat\.tech/ { found = 1 }
	END { if (found) print buf }
' "$wf")
test -n "$gitlab_steps" || fail 'could not extract GitLab workflow steps'

printf '%s\n' "$gitlab_steps" | grep -F 'continue-on-error' >/dev/null \
	&& fail 'GitLab steps must not set continue-on-error' || :

printf '%s\n' "$gitlab_steps" | grep -F 'set -x' >/dev/null \
	&& fail 'GitLab steps must not set -x (would print a token-bearing URL)' || :

# GitLab git remote interpolates the project variable; API uses /api/v4/.
# Fail a hardcoded group/project (host/group/repo), not the API prefix.
printf '%s\n' "$gitlab_steps" | grep -E 'gitlab\.sicepat\.tech/[[:alnum:]_.-]+/[[:alnum:]_.-]+' | grep -v '/api/' >/dev/null \
	&& fail 'GitLab URL must not hardcode a group/project path' || :

# Token in the remote URL is the env expansion, not a stored secret.
printf '%s\n' "$gitlab_steps" | grep -E 'oauth2:[^$]' >/dev/null \
	&& fail 'oauth2 username must use ${GITLAB_TEMPLATE_TOKEN}, not a literal' || :

# No repo-stored GitLab token patterns in the workflow or domain page.
doc="$ROOT/context/domains/source-release-and-upgrade/README.md"
require_file "$doc"
contains "$doc" 'GITLAB_TEMPLATE_PROJECT'
contains "$doc" 'GITLAB_TEMPLATE_TOKEN'
contains "$doc" 'gitlab.sicepat.tech'
contains "$doc" 'GitLab Release'
contains "$doc" 'write_repository'
contains "$doc" '`api`'
not_contains "$wf" 'glpat-'
not_contains "$doc" 'glpat-'
not_contains "$wf" 'private_token='
not_contains "$doc" 'private_token='

# Domain page documents names only — no gitlab.sicepat.tech/<literal-path>.
grep -E 'gitlab\.sicepat\.tech/[[:alnum:]_.-]' "$doc" >/dev/null \
	&& fail 'domain page must not hardcode a GitLab group/project path' || :

pass 'gitlab mirror workflow wiring'
