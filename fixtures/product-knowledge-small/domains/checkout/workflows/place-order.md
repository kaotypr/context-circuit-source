---
kind: workflow
title: Place an order
owners:
  - Checkout team
sources:
  - Acme Storefront PRD v3
review_date: 2026-01-05
implementation_ownership: storefront-api (orders service)
known_gaps:
  - Retry behavior after a declined card is not confirmed.
---

# Place an order

## Outcome

A shopper's cart becomes a paid, confirmed order.

## Actors

Shopper; orders service; payment provider.

## Entry points

The shopper selects "Place order" on the checkout screen.

## Current flow

1. Validate the cart is non-empty and prices are current.
2. Authorize payment with the payment provider.
3. Create the order and mark it paid.
4. Send an order confirmation.

## Variations

Registered shoppers reuse a saved address; guests enter one inline.

## Business rules

An order is only confirmed after payment authorization succeeds. Totals include
tax computed from the shipping address.
