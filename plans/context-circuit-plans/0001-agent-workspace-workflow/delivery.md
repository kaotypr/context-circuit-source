# Delivery

- Review and land the normative workflow contract and context updates as the first implementation slice.
- Implement and document the runtime records, session lifecycle, leases, handoffs, and ownership invariants.
- Implement root entry, resume, delegation, worker, and verifier behavior as agent-facing instructions and filesystem conventions.
- Integrate plan-scoped worktree ownership and lease checks as filesystem conventions owned by the session coordinator and workers.
- Remove the old Node/JavaScript scripts, generated bundle, command-specific host adapters, and release workflow from the repository.
- Add and run a command-independent acceptance suite, including adversarial contention and recovery cases, then obtain explicit human review of results.
