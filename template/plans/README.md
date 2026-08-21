# Human plans

Plans are reviewed intent, not execution sessions. A plan bundle lives at
`plans/<repository>-plans/<plan>/`:

```text
plan.yaml       # canonical machine fields and lifecycle status
PLAN.md         # sole human review entry point
tasks/*.md      # bounded task details and references
archive.yaml    # optional append-only eligibility history
```

`plan.yaml` owns IDs, status (`draft`, `approved`, `done`), repository, source
and Product Knowledge references, bounded paths, dependencies, tasks,
acceptance criteria, verification IDs/commands, and human gates. `PLAN.md`
explains those fields but does not become a second authority. Task status is a
projection: draft → draft, approved → ready, done → done.

Human actions remain separate:

1. Review is read-only and produces a Review Card.
2. Approve changes status only after a current explicit approval card.
3. Run requires a separate explicit named request.
4. Verify produces runtime evidence, not completion status.
5. Finish changes approved → done only after a current status-change card.
6. Delivery, archive, publication, deployment, merge, takeover, and cleanup
   each have their own confirmation.

Approval never starts execution.

Schema v1 plans remain readable. v2 adds fields and references without silently
rewriting old intent or status. Runtime reviews, leases, handoffs, evidence,
and completion records belong under `.runtime/`, not inside the plan bundle.
