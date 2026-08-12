---
kind: workflow
title: Create Shipment
owners:
  - Shipping team
sources:
  - Commerce platform PRD
review_date: 2026-02-10
implementation_ownership: shipping-service
known_gaps:
  - Edge cases for create shipment are not fully confirmed.
---

# Create Shipment

## Outcome

The business outcome of create shipment in the shipping domain.

## Actors

Relevant roles and the shipping-service.

## Entry points

Triggered from the Shipping surfaces.

## Current flow

1. Start create shipment.
2. Apply shipping rules.
3. Record the result.

## Variations

Behavior differs by role and account state.

## Business rules

Rules that govern create shipment. Owned by this page.
