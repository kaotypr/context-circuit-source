# Plan bundle prototype comparison

This is a bounded fixture measurement, not a claim that every historical
bundle has identical duplication. It compares the modeled full bundle with the
compact default while checking that rationale is not compressed away.

| Measure | Before: full companion prototype | After: compact prototype |
| --- | ---: | ---: |
| Bundle files | 9 | 3 |
| Canonical lifecycle copies outside `plan.yaml` | 2 | 0 |
| Canonical acceptance copies outside `plan.yaml` | 3 | 0 |
| Canonical verification copies outside `plan.yaml` | 3 | 0 |
| Rationale surfaces | 7 specialist narratives | 1 overview |

The before counts are an inventory of repeated authority surfaces in the
prototype, not a migration instruction. The after bundle retains human
rationale in `overview.md`, keeps bounded tasks, and adds specialist
companions when explicit complexity signals require them. The accompanying
`duplication-inventory.yaml` records the counts and safety checks.
