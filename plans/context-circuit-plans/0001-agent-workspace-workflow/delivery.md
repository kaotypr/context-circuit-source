# Delivery

- Review and land the normative workflow contract and context updates as the first implementation slice.
- Implement and document the runtime records, session lifecycle, leases, handoffs, and ownership invariants.
- Implement root entry, resume, delegation, worker, and verifier behavior as agent-facing instructions and filesystem conventions.
- Integrate plan-scoped worktree ownership and lease checks without removing or bypassing the existing command adapter.
- Add and run the complete acceptance suite, including adversarial contention and recovery cases, then obtain explicit human review of results.
- Only after the replacement passes, draft a separate cleanup plan for legacy Node/JavaScript removal; do not perform cleanup in this plan.
