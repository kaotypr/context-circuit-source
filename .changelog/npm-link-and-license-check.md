---
product: template
kind: fixed
authors: "@hallode"
pr: 169
---

The worktree reference also names an `npm link` inside node_modules as an entry
reuse rejects, and the `cc-cli` install script now requires the CLI archive's
`LICENSE` to be a regular file, as it already did for the binary, `README.md`,
and `THIRD_PARTY_NOTICES.txt`.
