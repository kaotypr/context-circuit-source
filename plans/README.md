# Plans roadmap

Numbered peer plans live under `plans/<repository-key>-plans/`. Each plan is a human-reviewed set of YAML metadata, Markdown explanation, and task files.

Active plans use `<number>-<slug>`. Manual archival moves the exact directory to `archives/plans/<repository-key>-plans/<number>-<slug>/`; archived plans are ignored by `whats-next`.

Create a draft with `cc create-plan`, approve it explicitly with `cc set-plan-state`, and execute one selected task with `cc run-task`.
