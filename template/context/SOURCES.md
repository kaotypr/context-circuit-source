# Sources

`sources/` is a passive, user-controlled inbox for raw evidence, briefs, and
PRDs. It is not Product Knowledge and is never scanned wholesale to fill a
context gap. An agent reads a source only when the human names the exact file or
a plan explicitly includes it.

Record provenance for each accepted source in `context/sources.yaml`.

`sources/README.md` documents a suggested folder layout (`raw/`,
`system-design/<domain>/<version>/`, `research/`, `decisions/`, `data-samples/`,
`reports/`, `archive/`). It is an ergonomic convention only: the agent infers
no status or authority from folder or file names, and passivity above still
applies.
