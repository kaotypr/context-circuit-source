---
kind: domain
title: Checkout
owners:
  - Checkout team
sources:
  - Acme Storefront PRD v3
review_date: 2026-01-05
workflows:
  - workflows/place-order.md
known_gaps:
  - Split payments are not covered.
---

# Checkout

## Summary

Everything involved in turning a cart into a paid, confirmed order. Route here for
payment and order-confirmation behavior.

## Workflows

- [Place an order](workflows/place-order.md) — pay for a cart and confirm the order.
