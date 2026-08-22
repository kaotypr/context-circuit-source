# Maintainer plans

This source-only directory contains implementation plans for improving the
Context Circuit product itself. It is not the released workspace plan area.

Maintainer plans live under a repository-specific namespace:

```text
plans/context-circuit-plans/<plan-id>/
├── plan.yaml
├── PLAN.md
└── tasks/
```

The release assembler does not copy this directory. Released workspaces get
the blank user-facing plan seed from `template/plans/` at their root instead.

Repository identity is shared in `workspace.yaml`; machine-specific checkout
paths belong in the ignored root `repositories.local.yaml`. Neither local
bindings nor repository contents are copied into a release artifact.
Maintainer plans still follow the same plan, task, lifecycle, gate, ownership,
and verification contracts so the product can be tested on itself.
