# Product Knowledge

Product Knowledge is concise, source-cited Markdown under `context/`. A typical tree contains `PROJECT.md`, role pages under `roles/`, domain summaries under `domains/`, and workflow pages under `domains/<domain>/workflows/`. When a workspace starts with only an idea, `$cc-idea-brief` first creates the human-reviewed `context/IDEA-BRIEF.md`; that brief can then be imported as a source.

Source provenance lives in `context/sources.yaml`:

```yaml
sources:
  - id: storefront-prd
    kind: prd
    location: docs/storefront-prd.md
    revision: sha256:...
    product_knowledge:
      - context/PROJECT.md
```

Use `import-product-knowledge` for an explicit import. Use `refresh-product-knowledge` when a source may have changed. The refresh compares lightweight content revisions and proposes a source-cited Markdown diff. Humans accept, revise, or reject the proposal; refresh never happens silently.

Plans and tasks may reference Product Knowledge pages directly. The references provide context to the human and agent.
