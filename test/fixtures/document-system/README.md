# Document-system foundation fixtures

These fixtures exercise the Plan 0016 foundation contract without installing a
repository parser or command runtime. The expected result categories use
`schemas/documents/validation-result-v1.yaml`.

The fixture set covers:

- valid and invalid base OKF concepts;
- the stricter Context Circuit profile;
- reserved `index.md` and `log.md`;
- legacy `kind`, type-only, matching `kind`/`type`, and conflicting values;
- `verified` mapping/list normalization;
- top-level and nested unknown extension preservation;
- initial front matter with a later Markdown horizontal rule;
- an unterminated front-matter extraction failure;
- YAML and artifact-schema failures;
- advisory-only broken-link and missing-index findings;
- non-OKF plan, runtime, instruction, and arbitrary-source boundaries.
- compact, complexity-triggered, and historical full plan bundles;
- invalid task delimiters, YAML, required fields, enum values, and nesting.

Expected outcomes are recorded in `expected-results.yaml`. They are evidence
fixtures, not a replacement for the approved host-provided deterministic YAML
and schema capability.
