---
name: cc-source-develop
description: Maintain Context Circuit source, Go runtime, and release assembly; not generated workspaces.
---

# Maintain Context Circuit source

Use root AGENTS.md and CONTRIBUTING.md. Edit directly on the active branch from the
maintainer request. Product entry instructions are packaging material.
Keep generated behavior in product/AGENTS.md.in and the release mapping exact.
Exercise YAML, records, Git, and packaging in fresh temporary directories with
go test ./... and scripts/check-release.sh. Do not run retired product lifecycles.
Report native-host behavior as measured on a named host, model, and date by the
agent tests in agent-tests/, and as not exercised when it was not.
