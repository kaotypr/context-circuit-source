# Publish approved plan tasks

Task publication is a separate explicit action. The host first searches the configured destination and writes a credential-free discovery JSON matching `.agents/contracts/plan-publication-discovery.schema.json`. Confirmed stable-work-ID mappings prevent duplicate creation.

```bash
node .agents/bin/cc.mjs prepare-plan-publication --plan <plan-id> --discovery <discovery.json>
```

Review the returned destination, hierarchy, dependency-first order, existing mappings, proposed creates, and idempotency keys before authorizing writes. After each confirmed external response, record it:

```bash
node .agents/bin/cc.mjs record-plan-publication \
  --plan <plan-id> --work <work-id> --status created \
  --reference <external-ref> --evidence "Creation confirmed by the session tool."
```

Use `--status failed` with evidence when creation fails, then stop. The record remains `partial` when earlier tasks succeeded, preserves untouched proposed items, and supports a safe retry. A confirmed reference is written back to the plan work breakdown without revoking approval. Draft plans, provider mismatches, undeclared `create-tasks` capability, conflicting mappings, unknown work IDs, and credential-bearing evidence are rejected.
