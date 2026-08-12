# Publish approved plan tasks

Task publication is a separate explicit action. The host first searches the configured destination and writes a credential-free discovery JSON matching `.agents/contracts/plan-publication-discovery.schema.json`. Confirmed stable-work-ID mappings prevent duplicate creation.

```bash
node .agents/bin/cc.mjs prepare-plan-publication --plan <plan-id> --discovery <discovery.json>
```

Review the returned destination, hierarchy, dependency-first order, explicit
repository key, descriptive area, existing mappings, proposed creates, and
idempotency keys before authorizing writes. Preparation revalidates every
repository key against `workspace.yaml`; an unknown or removed registration
blocks publication before any external write. After each confirmed external
response, record it:

Publication record contract version 2 carries both `repository` and `area` for
each item so external task preparation preserves the same identity boundary.

```bash
node .agents/bin/cc.mjs record-plan-publication \
  --plan <plan-id> --work <work-id> --status created \
  --reference <external-ref> --evidence "Creation confirmed by the session tool."
```

Use `--status failed` with evidence when creation fails, then stop. The record remains `partial` when earlier tasks succeeded, preserves untouched proposed items, and supports a safe retry. A confirmed reference is written back to the plan work breakdown without revoking approval. Draft plans, provider mismatches, undeclared `create-tasks` capability, conflicting mappings, unknown work IDs, and credential-bearing evidence are rejected.
