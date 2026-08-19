# Complex fixture plan

This plan has multiple implementation surfaces and independent verification
paths. `plan.yaml` remains authoritative for lifecycle, scope, dependencies,
acceptance, and verification; this overview explains why specialist
companions are present.

## Companion rationale

- `requirements.md`: several stakeholder outcomes need a focused boundary
  record.
- `solution.md`: documentation and acceptance surfaces require a design
  trade-off.
- `risks.md`: a malformed task could block a readiness path.
- `delivery.md`: review and status-change remain separate human gates.
- `acceptance.md`: two canonical acceptance IDs need a readable review map.
- `verification.md`: shell syntax, acceptance, and independent review are
  separate checks.

See `plan.yaml#acceptance_criteria` and the task IDs for canonical references.
