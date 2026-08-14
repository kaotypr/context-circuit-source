---
name: cc-refresh-product-knowledge
description: Detect changed Product Knowledge sources and propose a human-reviewed Markdown refresh.
---

Run `node .agents/bin/cc.mjs refresh-product-knowledge`. It compares current source content with the lightweight revision in `context/sources.yaml`, reports affected pages, and proposes a source-cited diff. It never refreshes canonical pages silently.
