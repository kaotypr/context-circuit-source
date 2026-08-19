---
type: Reference
title: Nested unknown extension
description: Nested extension fields survive an OKF metadata round trip.
status: stable
sources:
  - id: fixture-source
    resource: fixture/source.md
    title: Nested source
    x-source-extension:
      preserve: source-value
generated:
  by: process:document-check
  at: 2026-08-18T10:00:00Z
  x-generated-extension:
    preserve: generated-value
verified:
  - by: human:kaotypr
    at: 2026-08-18T10:05:00Z
    x-verified-extension:
      preserve: verified-value
---

# Nested extension preservation

Unknown fields in known nested metadata mappings remain opaque and survive
rewrites.
