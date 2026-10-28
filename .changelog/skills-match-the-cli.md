---
product: template
kind: fixed
pr: 203
---

Several shipped skills now say what the CLI and installers actually do.
`cc-intent` and `cc-plan` say that `record approve` also returns
`commit_recommended`, to pass on with the person's authorization; `cc-complete`
reads the completed record from the `completed` key; `cc-stacked` records an
integration merge's resolution in the merge commit and the report to the person
rather than in a workspace record; and `cc-cli` explains a pin with a reserved
label such as `2.5.0-rc.1`, a prerelease of the stable line that is never
offered as an update, and that the installers accept only `2.*` versions without
a `v` prefix.
