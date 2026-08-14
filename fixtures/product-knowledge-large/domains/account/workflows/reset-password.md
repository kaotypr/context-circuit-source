---
kind: workflow
title: Reset Password
owners:
  - Account team
sources:
  - Commerce platform PRD
review_date: 2026-02-10
implementation_ownership: account-service
known_gaps:
  - Edge cases for reset password are not fully confirmed.
---

# Reset Password

## Outcome

The business outcome of reset password in the account domain.

## Actors

Relevant roles and the account-service.

## Entry points

Triggered from the Account surfaces.

## Current flow

1. Start reset password.
2. Apply account rules.
3. Record the result.

## Variations

Behavior differs by role and account state.

## Business rules

Rules that govern reset password. Owned by this page.
