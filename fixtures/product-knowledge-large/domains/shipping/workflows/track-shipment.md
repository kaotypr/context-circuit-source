---
kind: workflow
title: Track Shipment
owners:
  - Shipping team
sources:
  - Commerce platform PRD
review_date: 2026-02-10
implementation_ownership: shipping-service
known_gaps:
  - Edge cases for track shipment are not fully confirmed.
---

# Track Shipment

## Outcome

The business outcome of track shipment in the shipping domain.

## Actors

Relevant roles and the shipping-service.

## Entry points

Triggered from the Shipping surfaces.

## Current flow

1. Start track shipment.
2. Apply shipping rules.
3. Record the result.

## Variations

Behavior differs by role and account state.

## Business rules

Rules that govern track shipment. Owned by this page.
