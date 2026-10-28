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
| `pr` | Optional pull request numbers, separated by spaces, such as `pr: 177 180`. The entry links these instead of the pull request that added the file. Use it for a change several pull requests made, an entry written after its change merged, or a cherry-pick, which names the original fix. |

The body is the entry. Write it for a user, not a reviewer: what they will see
or have to do, not how the code changed. It is published, so it carries no
private names and no AI-assistance credit.

Name the file after the change, using only letters, digits, `.`, `_` and `-`,
for example `record-archive-flag.md`. Two pull requests that pick the same name
conflict when the second one merges, and that one renames its file.

## Describe the change since the last release

A change file says what will be different for someone using Context Circuit
between the last release and the next one. It is not a record of how the change
was built. When a pull request reworks, extends or undoes something an earlier
change file describes, and that file has not shipped yet, it edits or deletes
that file in the same pull request rather than adding another one:

- a feature reworked before its release is described once, as it finally works,
  and names every pull request that made it in `pr`, and everyone who wrote
  them in `authors`, because `pr` alone credits no one new;
- a change added and then reverted before its release appears nowhere;
- a fix for something that never shipped is not a fix, and needs no file.

So `.changelog/` always reads as the next release's notes would, and the draft
does not have to be untangled when the release is cut.

## When a pull request needs none

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
