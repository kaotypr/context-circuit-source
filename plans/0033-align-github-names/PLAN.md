# Plan 0033 — Align GitHub repository names with source and product identity

**Intent:** i024-align-github-names  
**Repository:** context-circuit-source  
**Tier:** Standard  
**Status:** done

## Objective

Swap GitHub slugs so this maintainer source lives at `kaotypr/context-circuit-source`,
the published template lives at `kaotypr/context-circuit`, and source-side
publication binding, publish Action, workspace canonical URL, and plans-github
issue records follow those names. The swap does not publish a new template version
or change the assembled product tree. GitLab mirror configuration is unchanged.

## Grounding (HEAD 3bc56d59)

Today `origin` is `git@github.com:kaotypr/context-circuit.git` (this source).
The template publishes to `kaotypr/context-circuit-template` via
`release/binding.yaml` and `.github/workflows/publish-template.yml`. Operational
GitHub slugs appear in:

- `workspace.yaml` — canonical URL still uses the short source name
- `release/binding.yaml` — `destination_repo: kaotypr/context-circuit-template`
- `.github/workflows/publish-template.yml` — checkout and `gh release create`
- `publication/plans-github/config.yaml` — `target_ref.repository`
- `publication/plans-github/published/*.yaml` — twenty issue URLs on the source repo

Scripts and tests refer to `context-circuit-template` only as the conceptual product
identity in comments; no hardcoded GitHub slugs there. `sources/` design history is
out of scope. `template/` and the assembled artifact are unchanged.

`plan-allocate-id` fails `PLAN_PREFIX_COLLISION` on archived duplicate numeric
prefixes; this plan uses the next in-band id **0033**.

## Open question disposition

**TEMPLATE_REPO_TOKEN (plan-resolution):** After the template repository is
renamed to `kaotypr/context-circuit`, verify the existing fine-grained PAT still
grants push and release access. If not, the operator re-grants or re-issues
`TEMPLATE_REPO_TOKEN` on the source repository before the next publish. This
rename does not run publish-template or cut a release.

## Decisions

**One plan, four tasks** — GitHub rename operator work, binding retarget,
publication record retarget, and static acceptance share one execution and
verification boundary.

- Rename source on GitHub first, then template, so the short name is free.
- Rewrite operational slugs only; keep conceptual `context-circuit-template`
  prose where it names the product identity, not a GitHub URL.
- Publication issue URLs explicitly target `kaotypr/context-circuit-source`.
- GitLab mirror vars/secrets and workflow steps stay unchanged.
- New `test/release/test-github-names.sh` for static slug checks; live GitHub
  state checked with `gh repo view` when credentials are available.

## Tasks

### GNA-001 — Rename GitHub repositories and refresh local remote

Using `gh`, rename `kaotypr/context-circuit` → `kaotypr/context-circuit-source`,
then `kaotypr/context-circuit-template` → `kaotypr/context-circuit`. Update
`origin` to the renamed source URL. Verify template PAT access; do not publish.

**Done when:** `gh repo view` confirms both new slugs; old template slug is gone;
`origin` points at `context-circuit-source`.

### GNA-002 — Retarget publication binding and publish workflow

Update `release/binding.yaml`, `.github/workflows/publish-template.yml`, and
`context/domains/source-release-and-upgrade/README.md` GitHub slug references.
GitLab mirror wiring unchanged.

**Done when:** binding and workflow target `kaotypr/context-circuit`; GitLab
mirror acceptance still passes.

### GNA-003 — Retarget workspace URL and plans-github publication records

Update `workspace.yaml` canonical URL, `publication/plans-github/config.yaml`,
and all twenty `published/*.yaml` issue URLs to `context-circuit-source`.

**Done when:** workspace and publication config name the source slug; no
published record still uses the pre-rename source issue URL pattern.

### GNA-004 — Static acceptance and no template release guard

Add `test/release/test-github-names.sh`, wire into `test/acceptance.sh`, confirm
`template_version` unchanged and no new `release/requests` entry.

**Done when:** `sh test/acceptance.sh` passes.

## Risks

- **Rename ordering:** Renaming the template before freeing the short name fails
  on GitHub; GNA-001 enforces source-first order.
- **Remote confusion:** Clones still on `kaotypr/context-circuit` must update
  `origin` after the source rename; this plan updates the maintainer checkout only.
- **PAT scope:** `TEMPLATE_REPO_TOKEN` may need re-grant after the template rename.
- **Accidental publish:** Any workflow run during the swap could target wrong slugs;
  do not dispatch publish-template until GNA-002 is merged.

## Verification summary

| Criterion | Primary check |
| --- | --- |
| AC-SOURCE-GITHUB | `gh repo view kaotypr/context-circuit-source` |
| AC-TEMPLATE-GITHUB | `gh repo view kaotypr/context-circuit` |
| AC-PUBLISH-TARGET | binding + workflow grep |
| AC-WORKSPACE-URL | workspace.yaml canonical_url |
| AC-SOURCE-LINKS | publication config + published URLs |
| AC-NO-TEMPLATE-RELEASE | manifest template_version + release/requests count + acceptance |
| AC-GITLAB-UNCHANGED | test-gitlab-mirror-workflow.sh |
