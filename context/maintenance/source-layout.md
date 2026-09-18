# Source layout and what ships

This checkout is the maintainer source, not a generated workspace. Work happens
directly on the active branch; the product's own intent-and-plan flow governs
generated workspaces, not this one.

## Component ownership

| Component | Owns |
| --- | --- |
| `context-circuit-source@product/AGENTS.md.in` | The always-loaded gates and prohibitions, and the pointer to each stage's skill |
| `context-circuit-source@product/docs/` | Workspace files, agent-facing commands, subagents, working records, worktrees |
| `context-circuit-source@product/skills/` | One skill per lifecycle stage, plus executable installation and subagent dispatch |
| `context-circuit-source@product/assets/readme/` | The artwork the workspace README and the template repository listing use |
| `context-circuit-source@internal/workspace/` | Bookkeeping, structured edits, Git and worktree operations |
| `context-circuit-source@internal/cli/` | The command surface and human or structured output |
| `context-circuit-source@internal/cow/` | Native filesystem cloning with independent-copy fallback |
| `context-circuit-source@internal/installer/` | Acceptance tests for the shipped install scripts; it carries no non-test source |
| `context-circuit-source@cmd/context-circuit/` | Executable entry point |
| `context-circuit-source@template/` | Blank workspace files embedded with the product instruction |
| `context-circuit-source@scripts/release-manifest.txt` | The exact source-to-output file mapping |
| `context-circuit-source@assets.go` | Embeds only product assets and materializes the manifest |
| `context-circuit-source@VERSION`, `CLI_VERSION` | The two independent version lines |
| `context-circuit-source@scripts/` | Build, checks, and explicitly invoked publication |
| `context-circuit-source@release/` | The published destination identity, and one release request per version |

## Three classes of material

| Class | What it holds |
| --- | --- |
| Shipped instruction | The product's behavior: `context-circuit-source@product/` |
| Mutable seed | Files copied into a new workspace: `context-circuit-source@template/` |
| Never shipped | This knowledge tree, the workspace's own records, release requests, the Go implementation, scripts, and the maintainer design and evidence material |

Product history and maintainer data never reach a release asset.

The manifest is the single place that decides which class a file is in, the
embedding reads only what it names, and the release check verifies the embedded
inventory against it. A new shipped file that is not in the manifest is simply
absent from the product, which is why the manifest and the file land together.
Artwork is shipped on the same terms as instruction — named once in the manifest
and once in the embed list — so removing an image means editing both, and an
image nothing references keeps shipping until someone does.

## Working rules for this checkout

Preserve unrelated and dirty work. Maintainer design and evidence material is
passive: read only exact files a request names, and never scan it by default.
Historical knowledge, records, and publication material are source history
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
