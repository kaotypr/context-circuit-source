# Plan 0032 — Push published template to GitLab from GitHub Action

**Intent:** i023-gitlab-template-push  
**Repository:** context-circuit-source  
**Tier:** Standard  
**Status:** done

## Objective

After a successful template publication to GitHub (main, version tag, and GitHub
Release unchanged), the same `publish-template` workflow also pushes that same
commit and annotated tag from the `.template-repo` checkout to
`gitlab.sicepat.tech`, then creates a GitLab Release with the same notes and
archive as GitHub. GitHub remains canonical; GitLab receives a mirror of the
tree already published, not a second assembly.

## Grounding (HEAD 4575c02)

Publication today is one job in `.github/workflows/publish-template.yml`:

1. Checkout source and `kaotypr/context-circuit-template` into `.template-repo`
2. Run acceptance, resolve version, assemble/commit/tag via `scripts/publish-template.sh`
3. Push `main` and `v*` tag to GitHub `origin`
4. Create GitHub Release with the dist archive

`scripts/publish-template.sh` stops at `remaining_gate`; push and Release are
workflow-only. `release/binding.yaml` names only the GitHub destination.
`test/release/test-publish.sh` covers the script hermetically; no workflow YAML
tests exist yet.

`plan-allocate-id` fails `PLAN_PREFIX_COLLISION` on archived duplicate numeric
prefixes; this plan uses the next in-band id **0032**.

## Open question disposition

**GitLab project path (plan-resolution):** Do not hardcode the path in any
workspace file. Operators configure `vars.GITLAB_TEMPLATE_PROJECT`
(`group/project`, no host prefix) and `secrets.GITLAB_TEMPLATE_TOKEN` in GitHub
Actions settings before the first publish after merge. Static acceptance proves
wiring; end-to-end GitLab push succeeds once the project exists and values are
set.

## Decisions

**One plan, three tasks** — workflow change, GitLab Release, and static proof
share one execution and verification boundary.

- Add a workflow step that pushes `.template-repo` `main` and the resolved `v*`
  tag to GitLab using a transient remote and the configured variable/secret.
- Do not extend `publish-template.sh` or `release/binding.yaml` with GitLab
  credentials or paths.
- GitLab step must not use `continue-on-error`.
- After the GitLab git push, create a GitLab Release (same notes and archive
  as GitHub) via the GitLab API.
- New `test/release/test-gitlab-mirror-workflow.sh` for structural checks; wire
  into `test/acceptance.sh`.

## Tasks

### GLP-001 — GitLab mirror push in workflow

Extend `.github/workflows/publish-template.yml` with a GitLab push step after the
GitHub push, using the same `.template-repo` checkout. Document the mirror in
`context/domains/source-release-and-upgrade/README.md` (names only, no literal
path).

**Done when:** `sh test/release/test-gitlab-mirror-workflow.sh` passes.

### GLP-002 — Static acceptance and credential hygiene

Add `test/release/test-gitlab-mirror-workflow.sh`, wire it into
`test/acceptance.sh`, and confirm `test/release/test-publish.sh` is unchanged.

**Done when:** `sh test/acceptance.sh` passes.

### GLP-003 — GitLab Release after the mirror push

After the GitLab git push, create a GitLab Release for `v$V` using the same
`release-notes.md` and dist archive as the GitHub Release. Document that the
project access token needs `write_repository` and `api`.

**Done when:** `sh test/release/test-gitlab-mirror-workflow.sh` passes.

## Risks

- GitLab token echoed in logs — mitigated by GitHub Actions secret masking and
  avoiding `set -x` with token-bearing URLs in the step script.
- Hardcoding a guessed project path — mitigated by variable-only wiring and grep
  checks in GLP-002.
- Accidentally reordering or weakening GitHub push/Release — mitigated by
  workflow structure tests asserting those steps remain.

## Verification

Independent verifier runs task verification commands and full semantic acceptance.
End-to-end GitLab push requires operator configuration of
`GITLAB_TEMPLATE_PROJECT` and `GITLAB_TEMPLATE_TOKEN` on the repository; that
configuration is outside this plan's source diff. No delivery, merge, or
publication in this plan.
