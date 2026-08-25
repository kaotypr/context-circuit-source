# Template binding, versioning, and publication

This module defines how the source repository publishes the assembled product
into the product template repository
`git@github.com:kaotypr/context-circuit-template.git`, how that repository
carries its own independent version, and how a manually- or file-triggered
GitHub Actions workflow keeps its Git history and releases clean.

It builds on [09-source-and-template.md](./09-source-and-template.md), which
owns the source / template / workspace identity boundary and the assembly
boundary. This module does not restate that boundary; it defines the
*publication* mechanics on top of it. The release boundary owner remains
`wrapper/manifest.yaml` (`release_boundary`); the staged inventory owner remains
`scripts/release-manifest.txt`; the assembler is `scripts/release-artifact.sh`.
This module adds no competing policy owner.

## 1. Two independent version lines

Context Circuit tracks two version identities that must not be conflated.

~~~text
source / runtime line       e.g. 0.5.0        wrapper/manifest.yaml: runtime_version
  → the product implementation and runtime contract line in the source repo
  → advances with source development, independent of any publication

published template release   0.0.1-alpha.1 →   wrapper/manifest.yaml: template_version
  → the version of the distributable product-template repository itself
  → advances only when a source change reaches the published template
  → authoritative form is the template repository's git tag v<version>
~~~

The published template release version is the subject of this document. It is
**not** the source version and does not follow it. A source revision of the
`0.5.0` runtime line may first appear in published template release
`0.0.1-alpha.1`; a later source revision may only ever appear as
`0.0.3-alpha.4`. The two lines advance on their own cadence.

`wrapper/manifest.yaml` `template_version` is **reset to `0.0.1-alpha.1`** and
from now on means the published template release version. Its previous `0.5.0`
value carried no downstream consumer (nothing was published yet), so the reset
has no migration cost. The value is **stamped at publish time** (Section 6),
not hand-edited per release; the git tag is authoritative and the manifest field
is the in-tree mirror an installed workspace reads to self-report its version.
There is no separate `VERSION` file.

## 2. The two repositories

~~~text
context-circuit-source (this checkout)
  → owns implementation, contracts, seed, assembler, release ledger, and workflow
  → decides when the template is published

context-circuit-template (kaotypr/context-circuit-template)
  → the published product repository
  → default branch: main
  → carries its own SemVer line, its own history, its own tags and releases
  → authoritative source of truth for the current published version (latest tag)
~~~

Publication is one-directional: a workflow running in the source repository
assembles an artifact and writes it into the template repository. The template
repository is never edited by hand and never receives feature commits of its
own; every commit on `main` is a published release (Section 7).

## 3. Independent semantic versioning

The published template release version uses Semantic Versioning 2.0.0. It starts
at `0.0.1-alpha.1`.

Because the line is pre-`1.0.0` and in an `alpha` pre-release:

- **During the alpha line**, each publication that changes the template
  increments the pre-release counter: `0.0.1-alpha.1` → `0.0.1-alpha.2` → …
  This is the default bump and needs no core-triple decision.
- **Graduating the core triple** (off `0.0.1`, or off `alpha`) is a deliberate
  maintainer decision recorded in the release request: e.g.
  `0.0.1-alpha.N` → `0.0.1-beta.1` → `0.0.1` → `0.1.0`.
- **Pre-`1.0.0` core changes** follow the common convention: a
  backward-incompatible change bumps the *minor* (`0.y.0`), a compatible change
  or fix bumps the *patch* (`0.y.z`). Strict major/minor/patch meaning applies
  once the line reaches `1.0.0`.

The bump *level* is a human decision recorded in the release request; the bump
*trigger* is mechanical (Section 4).

### 3.1 Core-change mapping (guidance)

| What changed in the assembled artifact | Post-`1.0.0` | Pre-`1.0.0` |
| --- | --- | --- |
| Breaking contract/schema/invariant, removed skill or role, changed authority | major | minor |
| New capability, new skill, additive schema field, new guide | minor | patch |
| Fix, doc wording, non-behavioral cleanup | patch | patch |
| No effect on the assembled artifact | no release | no release |

The mapping is advisory; the maintainer records the chosen level and reason in
the release request, and the version in the request filename is final.

## 4. What "affects the template" means

A source change requires a template version bump **iff it changes the content of
the assembled artifact.** This is decided mechanically, not by inspecting which
source files a commit touched.

~~~text
affects-template(source) :=
  assemble(source), with the publish-stamped fields normalized,
  differs from the tree published at the template repo's latest tag
~~~

Procedure (the workflow's diff-gate, Section 6 step 5):

1. Assemble the artifact from the current source checkout
   (`scripts/build-dist.sh`). This already excludes source-only material per
   `scripts/release-manifest.txt` and `wrapper/manifest.yaml` `release_boundary`.
2. Normalize the publish-stamped `template_version` field on both the assembled
   tree and the latest published tree, so the comparison reflects real content
   changes and not the version stamp itself.
3. Compare the two trees. Identical → nothing to publish (fail loudly so a
   mistaken publish is caught). Different → the diff is the exact content of the
   new release and informs the bump level (Section 3.1).

**First release.** When the template repository has no tags yet there is no
prior tree to compare against; the diff-gate does not apply and the assembled
tree is published unconditionally as the initial release (Section 10). The gate
applies to every subsequent publication.

This makes the boundary self-checking: a change touching only maintainer files
(`sources/`, `test/`, `template-harness/`, `scripts/`, `plans/`, `.runtime/`,
`context/`) assembles to an identical tree and never bumps the template, while
any change to a shipped path (`wrapper/`, `.agents/skills/`, `agents/`, `docs/`,
the adapter root files, or the mutable seed under `template/`) surfaces as a
diff.

## 5. The release ledger (`release/`)

Publication is driven by a source-side ledger at the repository root. It is the
audit trail of every intended release and the input the workflow consumes.

~~~text
release/                              (source-only; never enters the artifact)
├── binding.yaml                      # static binding to the template repo
└── requests/
    ├── 0.0.1-alpha.1.md              # one file per version; filename IS the version
    ├── 0.0.1-alpha.2.md
    └── …
~~~

`release/binding.yaml`:

~~~yaml
destination_repo: kaotypr/context-circuit-template   # owner/repo slug
destination_ref:  main
product_line:     context-circuit
initial_version:  0.0.1-alpha.1
~~~

`release/requests/<version>.md` — one file per published version. The
**filename is the version** (`0.0.1-alpha.2.md` → `v0.0.1-alpha.2`); the body is
the release notes and records the chosen bump level and reason:

~~~markdown
---
level: patch            # major | minor | patch | prerelease (Section 3.1)
reason: <one line: why this bump>
---

<release notes; becomes the template CHANGELOG entry and the GitHub Release body>
~~~

Because the assembler stages an allowlist (`.agents agents docs wrapper` plus a
fixed set of copied adapter/seed files), `release/` is never picked up by
assembly; `scripts/release-manifest.txt` also lists it `exclude` as
belt-and-suspenders. The ledger stays source-only.

## 6. Publication via GitHub Actions

Publication runs as a workflow in `context-circuit-source`
(`.github/workflows/publish-template.yml`) that writes into
`context-circuit-template`. It is never performed by the runtime:
`wrapper/runtime/engine.sh` contains no push/publish behavior (INV-RUNTIME-01),
and the assembler stops at `remaining_gate: publication`.

### 6.1 Triggers

~~~yaml
on:
  workflow_dispatch:                  # manual: the primary, always-available path
    inputs:
      version:                        # optional; default = latest unpublished request
        required: false
  push:                               # convenience: publishes a newly-landed request
    branches: [ main ]                # source release branch (remote default)
    paths: [ 'release/requests/**' ]
~~~

Both triggers run the same job. `workflow_dispatch` publishes the input version,
or the single latest unpublished request when the input is omitted. A `push`
publishes the request file added or changed in that push. Either way the job
resolves **exactly one** target version; if resolution is ambiguous (several
unpublished requests, no input) it fails rather than guessing.

### 6.2 Job steps

~~~text
1. checkout source
2. resolve version  from the trigger (input, or the request filename)
3. verify           sh test/acceptance.sh
4. assemble         sh scripts/build-dist.sh
5. diff-gate        assembled tree vs. template repo's latest tag (Section 4)
                    → identical: fail ("nothing to publish for v<X>")
                    → no tags yet (first release): skip the gate, publish
6. guard            v<X> must not already be a tag in the template repo
7. clone template repo with the fine-grained PAT; if it is empty (no branch
                    yet), point the unborn HEAD at main so the first release
                    becomes main's first commit (release-only history)
8. full-tree replace the template working tree with the assembled tree
                    (copy with delete, so the tree equals the artifact exactly)
9. stamp            set template_version=<X> in wrapper/manifest.yaml;
                    prepend the request notes to the template CHANGELOG.md
10. commit          one commit on main:  chore(release): v<X>
                    author/committer = the maintainer identity (Section 8)
11. tag             annotated tag  v<X>  on that commit
12. push            main + tag to context-circuit-template
13. release         gh release create v<X> with the request notes as the body;
                    attach the artifact archive
~~~

## 7. One commit per version — template repo shape

The entire history of the template repository is a linear sequence of releases —
one commit, one tag, one GitHub Release, one version.

- **Full-tree replacement (step 8).** The published commit is the assembled
  artifact and nothing else. Files a previous release contained but the new
  artifact does not are deleted in the same commit, so the committed tree is
  byte-identical to `assemble(source)` with the version stamped. No incremental
  edit commits exist between releases.
- **One commit per bump (step 10).** No "prepare release" or "fix typo" commit.
  If a release is wrong, the fix is a **new** release (a new version), never an
  amendment or a force-rewrite of a published tag.
- **Release commit message.** Conventional Commits per INV-COMMIT-01, `chore`
  type, `release` scope: `chore(release): v<version>`. The body may summarize
  the change set; it never contains credentials, provider payloads, or host
  paths.
- **Tag immutability.** Published tags are never moved or deleted. `v0.0.1-alpha.1`
  always points at the commit that first published it.

## 8. Auth and identity

- **Cross-repo auth.** The workflow authenticates to
  `context-circuit-template` with a **fine-grained personal access token**,
  scoped to that single repository with Contents: read and write (covers push,
  tag, and `gh release create`), stored as a source-repo Actions secret (for
  example `TEMPLATE_REPO_TOKEN`). No SSH deploy key and no broad classic PAT.
- **Author identity.** The commit, tag, and release are authored and committed
  as the **human maintainer** (Adhitya Sanjaya,
  `adityasanjaya@sicepat.com`), configured via `git config user.name` /
  `user.email` in the workflow (kept in repo variables). The workflow **never**
  authors as `github-actions[bot]`.
- **No attribution.** Absolute, per `AGENTS.md`. No commit, tag, or GitHub
  Release text carries an AI/tool/bot attribution, co-author, or "generated by"
  trailer. The maintainer is the sole recorded author.
- **Secrets hygiene.** The token exists only in the Actions environment; it
  never enters a commit, tag, release note, changelog, or the assembled
  artifact (INV-SEC-01).

## 9. Records and source of truth

- **Authoritative current version:** the template repository's latest git tag
  `v<version>`.
- **In-tree mirror:** `wrapper/manifest.yaml` `template_version`, stamped at
  publish, equal to the tag.
- **Template repo notes/audit:** `CHANGELOG.md`, one entry per release, newest
  first, drawn from the release request.
- **Source-side ledger:** `release/requests/` — the intended-release history and
  the workflow input.
- **Traceability:** each release records the `source_revision` and
  `source_state` (`clean`/`dirty`) the assembler already reports, so any
  published template release traces to the exact source commit that produced it,
  and `assemble(source@revision)` can re-derive and byte-compare the tree.

## 10. Binding (first release)

The template repository is currently empty. The first workflow run performs the
binding — it is simply the first publication:

~~~text
Preconditions
  - context-circuit-template exists with default branch main and no tags
  - release/binding.yaml is present and correct
  - TEMPLATE_REPO_TOKEN and the maintainer identity variables are configured
  - release/requests/0.0.1-alpha.1.md exists

First run (workflow_dispatch or the push that adds the request)
  - resolves version 0.0.1-alpha.1
  - diff-gate is skipped: the template repo has no tags yet (Section 4)
  - assembles, stamps, and publishes it as the first commit on main:
      chore(release): v0.0.1-alpha.1  →  tag v0.0.1-alpha.1  →  GitHub Release
~~~

Binding is idempotent: a re-run for an already-published version stops at the
guard (step 6). After binding, all further publications follow Section 6 and
increment from `0.0.1-alpha.1` per Section 3.

## 11. Invariants and boundaries

- **Human gate.** Publication is a deliberate maintainer act — a manual dispatch
  or a committed request file. No `wrapper/runtime/` code pushes or publishes
  (INV-RUNTIME-01).
- **Commit convention.** Release commits and tags follow INV-COMMIT-01
  (`chore(release): v<version>`).
- **No attribution.** Absolute, per `AGENTS.md` — no tool/model/agent/bot credit
  anywhere in the record.
- **Assembly boundary.** Only template-owned files are published; source design,
  tests, harness, scripts, the `release/` ledger, plans, Product Knowledge,
  runtime evidence, credentials, and local bindings never enter the template
  repository (`wrapper/manifest.yaml` `release_boundary`,
  `scripts/release-manifest.txt`, `09-source-and-template.md` §6).
- **Clean history.** Exactly one commit per version bump; no non-release
  commits; no rewritten published tags.
- **Reproducibility.** `assemble(source@revision)` is deterministic, so a
  release can be re-derived from its recorded source revision and byte-compared.

## 12. Follow-up implementation items

Implied by this design; not performed by this document.

1. Reset `wrapper/manifest.yaml` `template_version` to `0.0.1-alpha.1` and
   document it as the published template release version (Section 1).
2. Add the root `release/` ledger: `binding.yaml` and
   `requests/0.0.1-alpha.1.md` (Section 5); add `exclude release` to
   `scripts/release-manifest.txt`.
3. Teach the diff-gate to normalize the stamped `template_version` before
   comparing (Section 4).
4. Have the assembler/workflow read destination and branch from
   `release/binding.yaml` instead of the hardcoded `DESTINATION_REPO` /
   `DESTINATION_REF` constants in `scripts/release-artifact.sh`; keep assembly
   source-versioned and apply the template version only at publish (stamp).
5. Add `.github/workflows/publish-template.yml` implementing Section 6, with the
   `TEMPLATE_REPO_TOKEN` secret and maintainer-identity variables (Section 8).
