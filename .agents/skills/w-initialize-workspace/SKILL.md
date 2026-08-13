---
name: w-initialize-workspace
description: Compatibility alias for configuring a Context Circuit wrapper. Detects fresh versus existing state and routes to the canonical $w-configure-workspace workflow.
---

# Initialize workspace compatibility alias

Read and follow `../w-configure-workspace/SKILL.md` in full. `$w-initialize-workspace` remains accepted for existing callers, but it is not a separate configuration workflow.

- A fresh or unborn wrapper routes to configuration and its explicit internal bootstrap phase.
- A wrapper with a commit routes to reviewable reconfiguration or inspection.
- State detection is reported clearly before mutation.
- Legacy `initialize-workspace --bootstrap <request>` remains accepted and delegates to the same configuration engine.

Do not treat ordinary configuration as an upgrade. Do not recreate the prior initializer rules here.
