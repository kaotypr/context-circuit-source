# Sources inbox

Place raw evidence here: authored briefs, PRDs, exported requirements, and other
human-selected source material. This directory is passive.

An agent reads a source only when you name the exact file or a plan explicitly
includes it. The workspace never scans all sources to fill a context gap.

## Suggested layout (recommendation, not a rule)

The structure below is an ergonomic convention for humans. It is not a contract:
the agent does not infer status, ordering, or authority from folder or file
names. A file under `archive/` is not auto-ignored, and a file under
`system-design/` is not auto-loaded — you still name the file or cite it in a
plan.

```
sources/
├── raw/                     # untouched originals — never edited
│   ├── 2026-08-20-prd-export.pdf
│   └── figma-mockups/
├── system-design/           # synthesized, reviewable pre-plan thinking
│   └── <domain>/<version>/  # e.g. billing/v0.1/, search/v0.2/
│       ├── 01-overview.md
│       ├── 02-data-model.md
│       └── 03-api-surface.md
├── research/                # competitor analysis, benchmarks, spike write-ups
├── decisions/               # short decision records: chose X over Y because…
├── data-samples/            # sample payloads, CSV rows, fixtures-as-evidence
├── reports/                 # generated analyses, audits, review outputs
└── archive/                 # superseded sources, kept for provenance
```

- `raw/` holds originals you never edit. `system-design/` holds design you
  refine freely *before* it earns a governed plan under `plans/`.
- The flow is left to right: `raw/` (immutable inbox) → `system-design/` (draft
  and review) → `plans/` (contract-governed, executable). Keeping design in
  `sources/` means you can iterate without touching plan state.
- `system-design/<domain>/<version>/` separates concerns by domain, then
  versions each domain independently (the harness and the source spec need not
  move in lockstep). Prefer domain names that echo the workspace vocabulary
  (for example the names under `context/domains/`), but nothing enforces this.

## Naming conventions (suggested)

- **Version folders** for evolving design: `v0.1/`, `v0.2/`.
- **Numeric prefixes** for reading order within a set: `01-overview.md`,
  `02-data-model.md`.
- **Date prefixes** for time-ordered inbox items: `2026-08-20-prd.pdf`.
- **`archive/`** mirrors `plans/archive/`: superseded but kept for provenance.
  When archiving, **preserve the original subpath** so
  provenance stays legible — for example a retired
  `system-design/billing/v0.1/` moves to
  `archive/system-design/billing/v0.1/`.
