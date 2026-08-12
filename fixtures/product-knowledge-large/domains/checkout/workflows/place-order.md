---
kind: workflow
title: Place Order
owners:
  - Checkout team
sources:
  - Commerce platform PRD
review_date: 2026-02-10
implementation_ownership: checkout-service
known_gaps:
  - Edge cases for place order are not fully confirmed.
---

# Place Order

## Outcome

The business outcome of place order in the checkout domain.

## Actors

Relevant roles and the checkout-service.

## Entry points

Triggered from the Checkout surfaces.

## Current flow

1. Start place order.
2. Apply checkout rules.
3. Record the result.

## Variations

Behavior differs by role and account state.

## Business rules

Rules that govern place order. Owned by this page.
