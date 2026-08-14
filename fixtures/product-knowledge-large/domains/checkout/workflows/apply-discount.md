---
kind: workflow
title: Apply Discount
owners:
  - Checkout team
sources:
  - Commerce platform PRD
review_date: 2026-02-10
implementation_ownership: checkout-service
known_gaps:
  - Edge cases for apply discount are not fully confirmed.
---

# Apply Discount

## Outcome

The business outcome of apply discount in the checkout domain.

## Actors

Relevant roles and the checkout-service.

## Entry points

Triggered from the Checkout surfaces.

## Current flow

1. Start apply discount.
2. Apply checkout rules.
3. Record the result.

## Variations

Behavior differs by role and account state.

## Business rules

Rules that govern apply discount. Owned by this page.
