---
kind: workflow
title: Issue Refund
owners:
  - Payments team
sources:
  - Commerce platform PRD
review_date: 2026-02-10
implementation_ownership: payments-service
known_gaps:
  - Edge cases for issue refund are not fully confirmed.
---

# Issue Refund

## Outcome

The business outcome of issue refund in the payments domain.

## Actors

Relevant roles and the payments-service.

## Entry points

Triggered from the Payments surfaces.

## Current flow

1. Start issue refund.
2. Apply payments rules.
3. Record the result.

## Variations

Behavior differs by role and account state.

## Business rules

Rules that govern issue refund. Owned by this page.
