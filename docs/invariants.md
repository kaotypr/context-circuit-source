# Invariant ownership

`wrapper/contracts/invariants.yaml` is the only normative invariant catalog.
This page explains how to use it and intentionally does not restate the rules.

When prose or a skill needs a safety rule, it names the invariant ID and links
to the owner. A semantic test must name the same ID. If two files appear to
define different values for one rule, the owner wins and the duplicate is
removed or reduced to a consequence.

The authority order is host/system safety, the shipped safety adapter,
workflow/router contract, repository-local instructions, human-approved plans
and decisions, accepted Product Knowledge, selected sources/repository
evidence, runtime evidence, then assumptions. Runtime cannot override a higher
authority.
