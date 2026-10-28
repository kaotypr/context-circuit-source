---
product: template
kind: changed
pr: 173
---

The worktree reference now says what copy-on-write reuse does with an entry it
cannot reuse: under `--copy-mode auto` or `copy` it is skipped and reported,
while `required` fails on it and `off` reuses nothing, and a rerun skips it
again until its source changes, while any other failure, such as a full disk,
stops preparation and keeps what was finished. A Windows file smaller than one
cluster is copied under `auto` and fails preparation under `required`.
