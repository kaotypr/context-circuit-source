# Initialization host proof results

The initialization workflow was exercised on 2026-08-11 in fresh, isolated
copies of the wrapper with the same clean ignored-clone fixture. The proof used
Codex CLI 0.147.0 and Claude Code 2.1.220 through their canonical thin adapters.

## Shared baseline

Each proof began from an uninitialized wrapper committed on `main`, an ignored
`repositories/frontend/` Git repository on `main`, no configured remote, and no
activity provider. The hosts were asked to initialize team mode without
committing, pushing, changing remotes, invoking activity tools, or editing
product code.

Both hosts:

- ran inspection-only initialization before writing;
- produced valid `workspace.yaml` and required workspace documents;
- reconciled the exact managed ignored-clone block;
- left the product repository clean on its original branch;
- passed `npm run validate -- --check-paths --check-documents` and
  `git diff --check`;
- reported `gitignore_changed: false` on the final check-only run; and
- left the wrapper at its baseline commit with only reviewable working-tree
  changes.

Codex changed six wrapper files: `.gitignore`, `workspace.yaml`, and four context
documents. Claude Code changed eight: those files plus `WORKFLOW.md` and the
frontend domain agent. These prose differences are acceptable host judgment;
the machine configuration, safety policy, repository mode, and deterministic
validation result were equivalent.

## Compatibility findings

The first Codex run exposed that the `tsx` command-line entry point creates a
temporary IPC socket that is rejected by its workspace sandbox. The shared npm
commands now run TypeScript with `node --import tsx`, which avoids that socket
without adding host-specific workflow logic. A regression test protects this
command surface.

Codex streamed a detailed trace and reread the product plan while executing.
Claude Code's noninteractive print mode buffered output for several minutes, so
filesystem changes were the useful progress surface until completion. Neither
behavior changed the resulting workflow contract.

Warnings about the fixture's missing remote and disabled activity provider were
expected and were preserved for human review. The automated suite covers solo
mode and submodule validation; manual two-host proof for those optional modes is
deferred until a later compatibility increment.
