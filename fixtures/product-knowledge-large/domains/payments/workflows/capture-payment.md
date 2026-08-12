---
kind: workflow
title: Capture Payment
owners:
  - Payments team
sources:
  - Commerce platform PRD
review_date: 2026-02-10
implementation_ownership: payments-service
known_gaps:
  - Edge cases for capture payment are not fully confirmed.
---

# Capture Payment

## Outcome

The business outcome of capture payment in the payments domain.

## Actors

Relevant roles and the payments-service.

## Entry points

Triggered from the Payments surfaces.

## Current flow

1. Start capture payment.
2. Apply payments rules.
3. Record the result.

## Variations

Behavior differs by role and account state.

## Business rules

Rules that govern capture payment. Owned by this page.
