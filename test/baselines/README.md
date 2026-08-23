# FINAL-001 baseline

This directory is the frozen pre-redesign evidence for Context Circuit 1.0.
`profile-baseline.yaml` records the byte values supplied by the design and the
target ceilings. `assertion-ledger.yaml` maps the protections in the former
monolithic acceptance harness to the semantic owner and test that replace it.
The ledger is evidence, not a second source of normative behavior.

The legacy acceptance run on 2026-08-21 passed before normative redesign edits:

```text
sh test/acceptance.sh  → PASS (all legacy scenarios)
```

Later acceptance is `sh test/acceptance.sh`; it runs the semantic suites. The
old harness is not a product implementation mechanism.
