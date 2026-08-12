# Work breakdown

| Work ID | Title | Parent | Depends on | Area | External reference |
| --- | --- | --- | --- | --- | --- |
| PKNOW-001 | Define Product Knowledge contracts, layout, templates, and validation | — | — | canonical-contracts | — |
| PKNOW-010 | Create the optional minimal Product Knowledge baseline during initialization | — | PKNOW-001 | initialization | — |
| PKNOW-020 | Discover source-backed knowledge candidates, contradictions, and gaps | — | PKNOW-001 | context-discovery | — |
| PKNOW-030 | Connect Product Knowledge references and impact declarations to plans | — | PKNOW-001, PKNOW-020 | planning | — |
| PKNOW-040 | Resolve bounded Product Knowledge packages for planless and planned tasks | — | PKNOW-020, PKNOW-030 | task-execution | — |
| PKNOW-050 | Carry Product Knowledge impact through worker, verifier, repair, and closeout evidence | — | PKNOW-040 | delivery-evidence | — |
| PKNOW-060 | Synchronize effective Product Knowledge through reviewed wrapper changes | — | PKNOW-020, PKNOW-050 | context-sync | — |
| PKNOW-070 | Generate revision-stamped role and domain onboarding packs | — | PKNOW-001, PKNOW-010 | onboarding | — |
| PKNOW-080 | Update host adapters, bundled commands, templates, and human guides | — | PKNOW-010, PKNOW-030, PKNOW-040, PKNOW-050, PKNOW-060, PKNOW-070 | host-adapters | — |
| PKNOW-090 | Prove bounded context, compatibility, lifecycle safety, and host parity end to end | — | PKNOW-080 | verification | — |

Live task status does not belong in this plan. Add confirmed external references only after an explicit publication action.
