---
product: cli
kind: added
pr: 172
---

`scripts/release-version.sh <cli|template>` prints the version a product builds
and publishes as, and every release script now asks it. A public checkout prints
`CLI_VERSION` or `VERSION` unchanged; a downstream builder sets a label in
`CC_VERSION_SUFFIX` and a counter in `release/variant.env` and gets
`<base>-<label>.<counter>`.
