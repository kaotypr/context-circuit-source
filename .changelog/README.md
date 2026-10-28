# Change files

A pull request that changes something someone using Context Circuit would
notice adds one file here, `.changelog/<slug>.md`. The note is written with the
change and reviewed in the same diff. When a release is cut,
`scripts/release-notes.sh` drafts the release request from these files, and the
maintainer edits the draft before it is published.

```markdown
---
product: cli
kind: fixed
upgrade: Run `agent setup` once after upgrading.
authors: "@your-handle"
---

One paragraph on what changes for someone using Context Circuit.
```

| Field | Values |
| --- | --- |
| `product` | `cli`, `template`, or `both`. Required. |
| `kind` | `added`, `fixed`, `changed`, or `breaking`. Required. `changed` is anything that is neither a feature nor a fix, such as reworded guidance in a shipped skill. |
| `upgrade` | Optional, one line: what the owner of an existing workspace has to do. It is listed under Upgrade notes. |
| `authors` | Optional GitHub handles, starting with `@` and separated by spaces. Without it, the draft credits the account the pull request came from, which is `kaotypr` for any branch pushed to this repository. |

The body is the entry. Write it for a user, not a reviewer: what they will see
or have to do, not how the code changed. It is published, so it carries no
private names and no AI-assistance credit.

Name the file after the change, using only letters, digits, `.`, `_` and `-`,
for example `record-archive-flag.md`. Two pull requests that pick the same name
conflict when the second one merges, and that one renames its file.

A pull request that changes nothing users see adds no file, and writes `None`
under Release note in its description. The Change files check fails a pull
request that changes CLI or template files with neither. Go tests and testdata
count for neither product.

When a release is cut, the draft lists every file for that product under its
kind, links the pull request that brought it to the release branch, and credits
its authors. A release candidate leaves the files in place, so the stable
release's draft covers them again. A stable release removes the files it used;
a `both` file is rewritten to name the other product, so that product's next
stable release still lists it.
