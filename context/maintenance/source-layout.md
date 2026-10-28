# Source layout and what ships

This checkout is the maintainer source, not a generated workspace. Work happens
directly on the active branch; the product's own intent-and-plan flow governs
generated workspaces, not this one.

## Component ownership

| Component | Owns |
| --- | --- |
| `context-circuit-source@product/AGENTS.md.in` | The always-loaded gates and prohibitions, and the pointer to each stage's skill |
| `context-circuit-source@product/docs/` | Workspace files, agent-facing commands, subagents, working records, worktrees |
| `context-circuit-source@product/skills/` | One skill per lifecycle stage, plus CLI installation and subagent dispatch |
| `context-circuit-source@product/README.md`, `product/assets/readme/` | The published template repository's guide and the artwork it renders; restored at publication, shipped to no workspace |
| `context-circuit-source@product/workspace-README.md` | The front page initialization writes into a new workspace |
| `context-circuit-source@product/launcher/` | The workspace launcher, which runs exactly the CLI version the workspace pins |
| `context-circuit-source@product/host/` | Host settings the template ships into a new workspace |
| `context-circuit-source@internal/workspace/` | Bookkeeping, structured edits, Git and worktree operations |
| `context-circuit-source@internal/cli/` | The command surface and human or structured output |
| `context-circuit-source@agent-tests/` | The contributor tool that runs coding hosts against a fixture workspace and grades the result: a separate Go module, run by contributors outside CI, shipped nowhere |
| `context-circuit-source@internal/cow/` | Native filesystem cloning with independent-copy fallback |
| `context-circuit-source@internal/installer/` | Acceptance tests for the shipped install scripts; it carries no non-test source |
| `context-circuit-source@cmd/context-circuit/` | CLI entry point |
| `context-circuit-source@template/` | Blank workspace files embedded with the product instruction |
| `context-circuit-source@scripts/release-manifest.txt` | The exact source-to-output file mapping |
| `context-circuit-source@assets.go` | Embeds only product assets and materializes the manifest |
| `context-circuit-source@VERSION`, `CLI_VERSION` | The two independent version lines |
| `context-circuit-source@scripts/` | Build, checks, and explicitly invoked publication |
| `context-circuit-source@release/` | The release document (`release/README.md`), the page each CLI archive ships as its README (`release/cli/README.md`), the published destination identity, one release request per version, and the landing-page files and export rules that repository owns |
| `context-circuit-source@.changelog/` | One change file per user-visible change, waiting for the release that drafts it into a request |
| `context-circuit-source@LICENSE`, `product/LICENSE` | The two license lines: Apache-2.0 for this checkout and the CLI, 0BSD for everything a workspace receives and shown only at the template repository's root |

## Five classes of material

| Class | What it holds |
| --- | --- |
| Shipped instruction | The product's behavior: `context-circuit-source@product/` |
| Mutable seed | Files copied into a new workspace: `context-circuit-source@template/` |
| Destination-owned | The published repository's own landing page: `context-circuit-source@release/template-repo/` plus the guide and license restored from `product/`, marked `export-ignore` so the archive a project is created from carries none of it |
| Contributor tool | The agent tests, a separate module that runs coding hosts against a fixture workspace and is run by contributors outside CI: `context-circuit-source@agent-tests/` |
| Never shipped | This knowledge tree, the workspace's own records, change files, release requests, the Go implementation, scripts, and the maintainer skills under `.agents/skills/` (`cc-source-develop`, `audit-cc-skills`) |

Product history and maintainer data never reach a release asset. The
destination-owned class exists because a workspace root belongs to somebody
else's project; see [licensing](../product/licensing.md) for what that rules out
and why.

The manifest is the single place that decides which class a file is in, the
embedding reads only what it names, and the release check verifies the embedded
inventory against it. A new shipped file that is not in the manifest is simply
absent from the product, which is why the manifest and the file land together.
Artwork is shipped on the same terms as instruction — named once in the manifest
and once in the embed list — so removing an image means editing both, and an
image nothing references keeps shipping until someone does.

## One knowledge tree

`context/` is the product's one knowledge tree, and it stays public on purpose.
A generated workspace keeps its knowledge to itself; this tree is the deliberate
exception, because contributors, forks, and the maintainer all read it. It must
always match what the code does and how agents behave in the generated template,
so a pull request that changes either one updates this tree in the same pull
request.

## Working rules for this checkout

Preserve unrelated and dirty work. Historical knowledge, records, and publication material are source history
rather than a specification for the current line, and are retrieved only to
answer a specific question.

Never inspect credentials, private provider payloads, host-local configuration,
or another session's runtime state, and never store secrets.

Commit, push, merge, publication, deployment, and deletion of generated
workspace data require explicit authorization. A source-only commit can be
authorized by the request itself, using a typed, scoped, imperative subject with
the changed component as scope. No agent credit, attribution, co-author, or
generated-by text goes into commits, pull requests, reviews, or comments, and an
authored message is inspected for injected attribution before the work is
reported finished.

Refactoring may remove obsolete repository-owned product material within the
requested scope, as long as release assembly and continuous-integration callers
stay coherent. Validation uses fresh temporary directories and never overwrites
existing build output or a generated workspace.
