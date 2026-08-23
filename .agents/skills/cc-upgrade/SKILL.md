---
name: cc-upgrade
description: Classify and safely migrate or roll back a wrapper in a live workspace.
---

Use when wrapper files change or a live workspace resumes on a new version.
Compare manifest and receipt versions, classify compatible/migration-needed/
blocked, preserve legacy-unknown records, and reload only changed evidence.
Require confirmation for breaking migrations. Replace wrapper-owned files only;
preserve workspace data, plans, runtime, dirty worktrees, and repositories.
