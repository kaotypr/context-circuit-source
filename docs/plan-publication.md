# Publish a plan before execution

Publication is optional and human-controlled. After a plan is approved and before execution, the human may publish the plan and selected tasks through a provider adapter.

```sh
node .agents/bin/cc.mjs publish-plan \
  --plan plans/api-plans/0010-checkout \
  --provider github \
  --references publication-urls.json
```

The references file contains stable IDs and URLs for the plan or its tasks. Publication preserves plan/task IDs and may store the current external URLs in `plan.yaml` and task frontmatter. It does not change statuses, start execution, monitor external state, or synchronize completion back into Context Circuit.
