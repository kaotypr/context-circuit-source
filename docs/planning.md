# Planning model

The authored planning surface is the root `plans/` roadmap. Repository
collections are named exactly `<repository-key>-plans`; `__BAU__` is a separate
track inside a collection. Numbered plans are peer review units with a README,
fixed unnumbered documents, and detailed task files under `tasks/`. Archived
plans move explicitly under `archived/plans/` and are never selected for new
execution.

Plan status is durable lifecycle metadata. Task attempt state belongs in the
runtime manifest. Human approval and archive remain explicit gates; execution,
review, merge confirmation, and closeout update status only after evidence is
validated.
