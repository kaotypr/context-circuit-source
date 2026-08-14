---
name: cc-import-product-knowledge
description: Explicitly import source-cited Product Knowledge and record lightweight source provenance.
---

Read the selected Idea Brief, PRD, document, repository evidence, or other source. Keep Product Knowledge concise and source-cited. Use `node .agents/bin/cc.mjs import-product-knowledge --source <workspace-relative-source>` or a JSON request. The command records `context/sources.yaml` and creates only missing Product Knowledge pages; existing pages receive a proposal instead of being silently overwritten.
